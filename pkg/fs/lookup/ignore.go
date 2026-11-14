// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lookup

import (
	"bytes"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
)

// IgnoreRuleSet holds the compiled ignore patterns parsed from a single ignore
// file (such as .gitignore or .ignore).
type IgnoreRuleSet struct {
	BaseDir    string
	SourceName string
	rules      []ignoreRule
}

// IgnoreMatcher parses, caches, and evaluates hierarchical ignore rule files
// (such as .gitignore and .ignore).
type IgnoreMatcher interface {
	LoadFile(df DiscoveredFile) (IgnoreRuleSet, error)
	Match(stack []IgnoreRuleSet, fullPath string, isDir bool) string
}

type ignoreRule struct {
	regex   *regexp.Regexp
	negate  bool
	dirOnly bool
}

type cachedRuleSet struct {
	meta         FileMeta
	resolvedPath string
	baseDir      string
	sourceName   string
	ruleSet      IgnoreRuleSet
}

type gitignoreMatcher struct {
	fsys  FileSystem
	mu    sync.RWMutex
	cache map[string]cachedRuleSet
}

// NewIgnoreMatcher creates a standalone [IgnoreMatcher] with stat-validated
// rule caching.
func NewIgnoreMatcher(opts ...Option) IgnoreMatcher {
	cfg := defaultConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return newIgnoreMatcherFromConfig(&cfg)
}

func newIgnoreMatcherFromConfig(cfg *config) *gitignoreMatcher {
	return &gitignoreMatcher{
		fsys:  cfg.fsys,
		cache: make(map[string]cachedRuleSet),
	}
}

func (m *gitignoreMatcher) LoadFile(
	df DiscoveredFile,
) (IgnoreRuleSet, error) {
	readPath := df.Path
	if readPath == "" {
		readPath = df.ResolvedPath
	}
	cacheKey := filepath.Clean(readPath)

	sourceName := df.Name
	if sourceName == "" {
		sourceName = filepath.Base(readPath)
	}
	baseDir := df.Dir
	if baseDir == "" {
		baseDir = filepath.Dir(readPath)
	}
	baseDir = filepath.Clean(baseDir)

	meta, err := m.resolveFileMeta(readPath, df)
	if err != nil {
		return IgnoreRuleSet{}, err
	}

	m.mu.RLock()
	entry, ok := m.cache[cacheKey]
	m.mu.RUnlock()
	if ok && entry.resolvedPath == df.ResolvedPath &&
		entry.baseDir == baseDir &&
		entry.sourceName == sourceName &&
		entry.meta.sameFingerprint(meta) {
		return entry.ruleSet, nil
	}

	data, err := m.fsys.ReadFile(readPath)
	if err != nil {
		return IgnoreRuleSet{}, err
	}

	rs := IgnoreRuleSet{
		BaseDir:    baseDir,
		SourceName: sourceName,
		rules:      parseIgnoreRules(data),
	}

	m.mu.Lock()
	m.cache[cacheKey] = cachedRuleSet{
		meta:         meta,
		resolvedPath: df.ResolvedPath,
		baseDir:      baseDir,
		sourceName:   sourceName,
		ruleSet:      rs,
	}
	m.mu.Unlock()

	return rs, nil
}

func (m *gitignoreMatcher) resolveFileMeta(
	readPath string,
	df DiscoveredFile,
) (FileMeta, error) {
	if df.Meta.Mode != 0 {
		return df.Meta, nil
	}
	if df.Mode != 0 || !df.ModTime.IsZero() || df.Size != 0 {
		return FileMeta{
			Name:    df.Name,
			Size:    df.Size,
			Mode:    df.Mode,
			ModTime: df.ModTime,
		}, nil
	}
	return m.fsys.Stat(readPath)
}

func (m *gitignoreMatcher) Match(
	stack []IgnoreRuleSet,
	fullPath string,
	isDir bool,
) string {
	if len(stack) == 0 {
		return ""
	}

	cleanFull := filepath.Clean(fullPath)
	for _, ancDir := range ancestorDirsUnderStack(stack, cleanFull) {
		if src := matchSinglePath(stack, ancDir, true); src != "" {
			return src
		}
	}

	return matchSinglePath(stack, cleanFull, isDir)
}

func ancestorDirsUnderStack(
	stack []IgnoreRuleSet,
	cleanFull string,
) []string {
	var ancestors []string
	for dir := filepath.Dir(cleanFull); dir != "" &&
		dir != "." &&
		dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
		underAny := false
		for _, rs := range stack {
			if len(rs.rules) == 0 {
				continue
			}
			if _, ok := relUnderBaseDir(rs.BaseDir, dir); ok {
				underAny = true
				break
			}
		}
		if !underAny {
			break
		}
		ancestors = append(ancestors, dir)
	}
	slices.Reverse(ancestors)
	return ancestors
}

func matchSinglePath(
	stack []IgnoreRuleSet,
	cleanFull string,
	isDir bool,
) string {
	ignored := false
	matchedSource := ""

	var lastBase string
	var lastRel string
	var lastValid bool

	for _, rs := range stack {
		if len(rs.rules) == 0 {
			continue
		}
		var slashRel string
		var valid bool
		if rs.BaseDir == lastBase {
			slashRel = lastRel
			valid = lastValid
		} else {
			slashRel, valid = relUnderBaseDir(rs.BaseDir, cleanFull)
			lastBase = rs.BaseDir
			lastRel = slashRel
			lastValid = valid
		}
		if !valid {
			continue
		}

		for _, r := range rs.rules {
			if r.dirOnly && !isDir {
				continue
			}
			if r.regex.MatchString(slashRel) {
				ignored = !r.negate
				if ignored {
					matchedSource = rs.SourceName
				} else {
					matchedSource = ""
				}
			}
		}
	}

	return matchedSource
}

func relUnderBaseDir(baseDir, fullPath string) (string, bool) {
	rel, err := filepath.Rel(baseDir, fullPath)
	if err != nil {
		return "", false
	}
	if rel == "." || rel == ".." ||
		strings.HasPrefix(rel, ".."+string(filepath.Separator)) ||
		filepath.IsAbs(rel) {
		return "", false
	}
	return filepath.ToSlash(rel), true
}

func parseIgnoreRules(data []byte) []ignoreRule {
	var rules []ignoreRule
	for line := range bytes.SplitSeq(data, []byte("\n")) {
		rule, ok := parseIgnoreLine(string(line))
		if ok {
			rules = append(rules, rule)
		}
	}
	return rules
}

func parseIgnoreLine(raw string) (ignoreRule, bool) {
	line := strings.TrimSuffix(raw, "\r")
	if line == "" || line[0] == '#' {
		return ignoreRule{}, false
	}

	line = trimUnescapedTrailingSpaces(line)
	if line == "" {
		return ignoreRule{}, false
	}

	negate := false
	if line[0] == '!' {
		negate = true
		line = line[1:]
		if line == "" {
			return ignoreRule{}, false
		}
	}

	dirOnly := false
	if strings.HasSuffix(line, "/") {
		dirOnly = true
		line = strings.TrimRight(line, "/")
		if line == "" {
			return ignoreRule{}, false
		}
	}

	anchored := false
	if strings.Contains(line, "/") {
		anchored = true
		line = strings.TrimPrefix(line, "/")
		if line == "" {
			return ignoreRule{}, false
		}
	}

	re, err := compileGitignoreGlob(line, anchored)
	if err != nil {
		return ignoreRule{}, false
	}

	return ignoreRule{
		regex:   re,
		negate:  negate,
		dirOnly: dirOnly,
	}, true
}

func trimUnescapedTrailingSpaces(s string) string {
	end := len(s)
	for end > 0 && s[end-1] == ' ' {
		backslashes := 0
		for i := end - 2; i >= 0 && s[i] == '\\'; i-- {
			backslashes++
		}
		if backslashes%2 == 1 {
			break
		}
		end--
	}
	return s[:end]
}

func compileGitignoreGlob(
	pat string,
	anchored bool,
) (*regexp.Regexp, error) {
	var b strings.Builder
	if !anchored {
		b.WriteString("^(?:[^/]+/)*")
	} else {
		b.WriteByte('^')
	}

	if pat == "**" {
		b.WriteString(".*$")
		return regexp.Compile(b.String())
	}

	i := 0
	for i < len(pat) {
		switch {
		case pat[i] == '\\':
			if i+1 < len(pat) {
				b.WriteString(regexp.QuoteMeta(string(pat[i+1])))
				i += 2
			} else {
				b.WriteString(`\\`)
				i++
			}
		case i == 0 && strings.HasPrefix(pat[i:], "**/"):
			b.WriteString("(?:[^/]+/)*")
			i += 3
			for strings.HasPrefix(pat[i:], "**/") {
				i += 3
			}
		case strings.HasPrefix(pat[i:], "/**/"):
			b.WriteString("/(?:[^/]+/)*")
			i += 4
			for strings.HasPrefix(pat[i:], "**/") {
				i += 3
			}
		case pat[i:] == "/**":
			b.WriteString("/.*")
			i += 3
		case pat[i] == '*':
			for i+1 < len(pat) && pat[i+1] == '*' {
				i++
			}
			b.WriteString("[^/]*")
			i++
		case pat[i] == '?':
			b.WriteString("[^/]")
			i++
		case pat[i] == '[':
			cls, next, ok := translateCharClass(pat, i)
			if ok {
				b.WriteString(cls)
				i = next
			} else {
				b.WriteString(`\[`)
				i++
			}
		default:
			b.WriteString(regexp.QuoteMeta(string(pat[i])))
			i++
		}
	}

	b.WriteByte('$')
	return regexp.Compile(b.String())
}

func translateCharClass(pat string, start int) (string, int, bool) {
	i := start + 1
	if i >= len(pat) {
		return "", start, false
	}

	var body strings.Builder
	negate := false
	if pat[i] == '!' || pat[i] == '^' {
		negate = true
		i++
	}

	if i < len(pat) && pat[i] == ']' {
		body.WriteString(`\]`)
		i++
	}

	for i < len(pat) {
		ch := pat[i]
		if ch == '/' {
			return "", start, false
		}
		if ch == '\\' && i+1 < len(pat) {
			if pat[i+1] == '/' {
				return "", start, false
			}
			body.WriteByte('\\')
			body.WriteByte(pat[i+1])
			i += 2
			continue
		}
		if ch == ']' {
			if body.Len() == 0 {
				return "", start, false
			}
			if negate {
				return "[^/" + body.String() + "]", i + 1, true
			}
			return "[" + body.String() + "]", i + 1, true
		}
		if ch == '[' {
			body.WriteString(`\[`)
		} else {
			body.WriteByte(ch)
		}
		i++
	}

	return "", start, false
}

func normalizeIgnoreFileNames(names []string) []string {
	if len(names) == 0 {
		return nil
	}
	out := make([]string, 0, len(names))
	for _, n := range names {
		if n == "" || slices.Contains(out, n) {
			continue
		}
		out = append(out, n)
	}
	return out
}
