// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

// Kind represents the classification code of a DML tag symbol.
type Kind byte

const (
	KindDevice    Kind = 'd'
	KindTemplate  Kind = 't'
	KindMethod    Kind = 'm'
	KindParam     Kind = 'p'
	KindBank      Kind = 'b'
	KindRegister  Kind = 'r'
	KindField     Kind = 'f'
	KindGroup     Kind = 'g'
	KindPort      Kind = 'P'
	KindAttribute Kind = 'a'
	KindConnect   Kind = 'c'
	KindEvent     Kind = 'e'
	KindImplement Kind = 'i'
	KindInterface Kind = 'I'
	KindSubdevice Kind = 'S'
	KindVariable  Kind = 'v'
	KindConstant  Kind = 'C'
	KindTypedef   Kind = 'T'
	KindStruct    Kind = 's'
	KindMember    Kind = 'M'
	KindHook      Kind = 'h'
	KindLogGroup  Kind = 'l'
	KindExtern    Kind = 'x'
)

// Letter returns the single-byte ASCII character representing this Kind.
func (k Kind) Letter() byte {
	return byte(k)
}

// Name returns the human-readable scope and kind name (e.g. "register", "bank").
func (k Kind) Name() string {
	switch k {
	case KindDevice:
		return "device"
	case KindTemplate:
		return "template"
	case KindMethod:
		return "method"
	case KindParam:
		return "param"
	case KindBank:
		return "bank"
	case KindRegister:
		return "register"
	case KindField:
		return "field"
	case KindGroup:
		return "group"
	case KindPort:
		return "port"
	case KindAttribute:
		return "attribute"
	case KindConnect:
		return "connect"
	case KindEvent:
		return "event"
	case KindImplement:
		return "implement"
	case KindInterface:
		return "interface"
	case KindSubdevice:
		return "subdevice"
	case KindVariable:
		return "variable"
	case KindConstant:
		return "constant"
	case KindTypedef:
		return "typedef"
	case KindStruct:
		return "struct"
	case KindMember:
		return "member"
	case KindHook:
		return "hook"
	case KindLogGroup:
		return "loggroup"
	case KindExtern:
		return "extern"
	default:
		return "unknown"
	}
}

// ExCmdMode specifies how the tag address field is formatted.
type ExCmdMode int

const (
	// ExCmdPattern formats the address as a search regex pattern /^...$/.
	ExCmdPattern ExCmdMode = iota
	// ExCmdNumber formats the address as a 1-based line number.
	ExCmdNumber
	// ExCmdCombine formats the address as line_number;/^...$/.
	ExCmdCombine
)

// SortMode specifies how tag entries are sorted in the generated tag file.
type SortMode int

const (
	// SortYes sorts tags in standard ASCII order.
	SortYes SortMode = iota
	// SortNo keeps tags in natural extraction order.
	SortNo
	// SortFoldCase sorts tags case-insensitively, tie-breaking case-sensitively.
	SortFoldCase
)

// Tag represents a declared DML symbol with location and scope metadata.
type Tag struct {
	Name       string
	File       string
	Line       int
	SourceLine string
	Kind       Kind
	ScopeKind  Kind
	ScopePath  string
	Signature  string
	Inherits   []string
}

// EscapePattern escapes backslashes and forward slashes in a source line
// for use inside an Ex search pattern (/^...$/).
func EscapePattern(s string) string {
	s = strings.TrimRight(s, "\r\n")
	var b strings.Builder
	b.Grow(len(s) + 4)
	for i := range len(s) {
		ch := s[i]
		if ch == '\\' || ch == '/' {
			b.WriteByte('\\')
		}
		b.WriteByte(ch)
	}
	return b.String()
}

// FormatExCmd formats the address portion of a tag line according to mode.
func FormatExCmd(line int, sourceLine string, mode ExCmdMode) string {
	escaped := EscapePattern(sourceLine)
	switch mode {
	case ExCmdNumber:
		if line <= 0 {
			line = 1
		}
		return strconv.Itoa(line)
	case ExCmdCombine:
		if line <= 0 {
			line = 1
		}
		return strconv.Itoa(line) + ";/^" + escaped + "$/"
	default:
		return "/^" + escaped + "$/"
	}
}

// FormatTag serializes a single Tag into the extended ctags line format:
//
//	<name>\t<file>\t<excmd>;"\t<kind>[\tline:<n>][\t<scopeKind>:<scopePath>][\tinherits:<t1,t2>][\tsignature:<sig>]
func FormatTag(t Tag, mode ExCmdMode) string {
	var b strings.Builder
	b.WriteString(t.Name)
	b.WriteByte('\t')
	b.WriteString(t.File)
	b.WriteByte('\t')
	b.WriteString(FormatExCmd(t.Line, t.SourceLine, mode))
	b.WriteString(";\"\t")
	b.WriteByte(t.Kind.Letter())

	if t.Line > 0 {
		b.WriteString("\tline:")
		b.WriteString(strconv.Itoa(t.Line))
	}

	if t.ScopeKind != 0 && t.ScopePath != "" {
		b.WriteByte('\t')
		b.WriteString(t.ScopeKind.Name())
		b.WriteByte(':')
		b.WriteString(t.ScopePath)
	}

	if len(t.Inherits) > 0 {
		b.WriteString("\tinherits:")
		b.WriteString(strings.Join(t.Inherits, ","))
	}

	if t.Signature != "" {
		b.WriteString("\tsignature:")
		b.WriteString(t.Signature)
	}

	return b.String()
}

// PseudoHeader returns the standard ctags pseudo-tag header lines.
func PseudoHeader(version string, sort SortMode, author string) []string {
	if author == "" {
		author = "Sergei Timoshchenko"
	}
	sortVal := "1"
	sortComment := "0=unsorted, 1=sorted, 2=foldcase"
	switch sort {
	case SortNo:
		sortVal = "0"
	case SortFoldCase:
		sortVal = "2"
	}

	if version == "" {
		version = "dev"
	}

	return []string{
		"!_TAG_FILE_FORMAT\t2\t/extended format; --format=1 will not append ;\" to lines/",
		fmt.Sprintf("!_TAG_FILE_SORTED\t%s\t/%s/", sortVal, sortComment),
		fmt.Sprintf("!_TAG_PROGRAM_AUTHOR\t%s\t//", author),
		"!_TAG_PROGRAM_NAME\tdmltag\t//",
		fmt.Sprintf("!_TAG_PROGRAM_VERSION\t%s\t//", version),
	}
}

// WriteOptions configures tag emission and file formatting.
type WriteOptions struct {
	Version  string
	Author   string
	Sort     SortMode
	ExCmd    ExCmdMode
	Append   bool
	Existing []string
	NoHeader bool
}

// ParseExistingTags reads an existing tags file and extracts non-header lines.
func ParseExistingTags(r io.Reader) ([]string, error) {
	if r == nil {
		return nil, nil
	}
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(line, "!_TAG_") {
			continue
		}
		lines = append(lines, line)
	}
	return lines, scanner.Err()
}

// tagLineKey extracts the sort key (column 1: tag name) from a formatted tag line.
func tagLineKey(line string) string {
	before, _, found := strings.Cut(line, "\t")
	if found {
		return before
	}
	return line
}

// SortLines sorts raw formatted tag lines according to SortMode.
func SortLines(lines []string, mode SortMode) {
	if mode == SortNo {
		return
	}
	if mode == SortFoldCase {
		slices.SortFunc(lines, func(a, b string) int {
			ka := tagLineKey(a)
			kb := tagLineKey(b)
			kla := strings.ToLower(ka)
			klb := strings.ToLower(kb)
			if c := strings.Compare(kla, klb); c != 0 {
				return c
			}
			if c := strings.Compare(ka, kb); c != 0 {
				return c
			}
			return strings.Compare(a, b)
		})
		return
	}
	slices.SortFunc(lines, func(a, b string) int {
		ka := tagLineKey(a)
		kb := tagLineKey(b)
		if c := strings.Compare(ka, kb); c != 0 {
			return c
		}
		return strings.Compare(a, b)
	})
}

// WriteTags writes headers and formatted tag lines to w.
func WriteTags(w io.Writer, tags []Tag, opts WriteOptions) error {
	var lines []string
	if opts.Append && len(opts.Existing) > 0 {
		lines = append(lines, opts.Existing...)
	}

	for _, t := range tags {
		lines = append(lines, FormatTag(t, opts.ExCmd))
	}

	SortLines(lines, opts.Sort)

	writer := bufio.NewWriter(w)

	if !opts.NoHeader {
		for _, hdr := range PseudoHeader(opts.Version, opts.Sort, opts.Author) {
			if _, err := writer.WriteString(hdr + "\n"); err != nil {
				return err
			}
		}
	}

	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return writer.Flush()
}
