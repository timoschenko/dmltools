// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lookup

import (
	"bytes"
	"regexp"
	"slices"
	"strings"
)

type ignoreRule struct {
	regex   *regexp.Regexp
	negate  bool
	dirOnly bool
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
