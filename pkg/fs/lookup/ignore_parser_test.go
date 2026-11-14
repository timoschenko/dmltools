// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lookup

import (
	"slices"
	"testing"
)

func TestParseIgnoreLine(t *testing.T) {
	tests := []struct {
		name       string
		raw        string
		wantOK     bool
		wantNeg    bool
		wantDir    bool
		wantAnchor bool
	}{
		{"empty", "", false, false, false, false},
		{"comment", "# just a comment", false, false, false, false},
		{"spaces only", "      \r", false, false, false, false},
		{"simple file", "foo.txt", true, false, false, false},
		{"negated file", "!foo.txt", true, true, false, false},
		{"dir only", "build/", true, false, true, false},
		{"negated dir only", "!build/", true, true, true, false},
		{"anchored path", "pkg/foo.txt", true, false, false, false},
		{"leading slash anchored", "/root.txt", true, false, false, false},
		{"escaped space", `hello\ world.txt   `, true, false, false, false},
		{"trailing slash only", "/", false, false, false, false},
		{"negation only", "!", false, false, false, false},
		{"escaped hash", `\#hash.dml`, true, false, false, false},
		{"escaped exclamation", `\!bang.dml`, true, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule, ok := parseIgnoreLine(tt.raw)
			if ok != tt.wantOK {
				t.Fatalf(
					"parseIgnoreLine(%q) ok = %v, want %v",
					tt.raw,
					ok,
					tt.wantOK,
				)
			}
			if !ok {
				return
			}
			if rule.negate != tt.wantNeg {
				t.Errorf("rule.negate = %v, want %v", rule.negate, tt.wantNeg)
			}
			if rule.dirOnly != tt.wantDir {
				t.Errorf("rule.dirOnly = %v, want %v", rule.dirOnly, tt.wantDir)
			}
			if rule.regex == nil {
				t.Error("rule.regex should not be nil")
			}
		})
	}
}

func TestParseIgnoreRules(t *testing.T) {
	data := []byte("# comment\nfoo.dml\n\n!bar.dml\nbuild/\n")
	rules := parseIgnoreRules(data)
	if len(rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(rules))
	}
	if rules[0].negate || rules[0].dirOnly {
		t.Errorf("rule[0] unexpected flags: %+v", rules[0])
	}
	if !rules[1].negate || rules[1].dirOnly {
		t.Errorf("rule[1] unexpected flags: %+v", rules[1])
	}
	if rules[2].negate || !rules[2].dirOnly {
		t.Errorf("rule[2] unexpected flags: %+v", rules[2])
	}
}

func TestTrimUnescapedTrailingSpaces(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"abc   ", "abc"},
		{`abc\ `, `abc\ `},
		{`abc\\ `, `abc\\`},
		{`abc\\\ `, `abc\\\ `},
		{"abc", "abc"},
		{"   ", ""},
	}
	for _, tt := range tests {
		got := trimUnescapedTrailingSpaces(tt.input)
		if got != tt.want {
			t.Errorf(
				"trimUnescapedTrailingSpaces(%q) = %q, want %q",
				tt.input,
				got,
				tt.want,
			)
		}
	}
}

func TestCompileGitignoreGlob_Patterns(t *testing.T) {
	tests := []struct {
		pattern  string
		anchored bool
		match    []string
		noMatch  []string
	}{
		{
			pattern:  "**",
			anchored: false,
			match:    []string{"a", "a/b", "a/b/c"},
			noMatch:  nil,
		},
		{
			pattern:  "*.txt",
			anchored: false,
			match:    []string{"file.txt", "a/b/c.txt"},
			noMatch:  []string{"file.txt2", "file.log"},
		},
		{
			pattern:  "foo/*.txt",
			anchored: true,
			match:    []string{"foo/a.txt"},
			noMatch:  []string{"bar/foo/a.txt", "foo/sub/a.txt"},
		},
		{
			pattern:  "**/build",
			anchored: true,
			match:    []string{"build", "a/build", "a/b/build"},
			noMatch:  []string{"builder", "a/build/nested"},
		},
		{
			pattern:  "a/**/b",
			anchored: true,
			match:    []string{"a/b", "a/x/b", "a/x/y/b"},
			noMatch:  []string{"a/b/c", "x/a/b"},
		},
		{
			pattern:  "file[12].dml",
			anchored: false,
			match:    []string{"file1.dml", "file2.dml", "sub/file1.dml"},
			noMatch:  []string{"file3.dml", "file.dml"},
		},
		{
			pattern:  "only[!1].dml",
			anchored: false,
			match:    []string{"only2.dml", "only3.dml"},
			noMatch:  []string{"only1.dml"},
		},
		{
			pattern:  "char?.dml",
			anchored: false,
			match:    []string{"chara.dml", "char1.dml"},
			noMatch:  []string{"charab.dml", "char/.dml"},
		},
		{
			pattern:  `esc\ .dml`,
			anchored: false,
			match:    []string{"esc .dml"},
			noMatch:  []string{"esc.dml", "esc  .dml"},
		},
		{
			pattern:  `box/**`,
			anchored: true,
			match:    []string{"box/a", "box/a/b"},
			noMatch:  []string{"box", "boxer"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			re, err := compileGitignoreGlob(tt.pattern, tt.anchored)
			if err != nil {
				t.Fatalf("compileGitignoreGlob(%q) failed: %v", tt.pattern, err)
			}
			for _, m := range tt.match {
				if !re.MatchString(m) {
					t.Errorf("pattern %q should match %q", tt.pattern, m)
				}
			}
			for _, nm := range tt.noMatch {
				if re.MatchString(nm) {
					t.Errorf("pattern %q should NOT match %q", tt.pattern, nm)
				}
			}
		})
	}
}

func TestNormalizeIgnoreFileNames(t *testing.T) {
	input := []string{".gitignore", "", ".ignore", ".gitignore", ".custom"}
	got := normalizeIgnoreFileNames(input)
	want := []string{".gitignore", ".ignore", ".custom"}
	if !slices.Equal(got, want) {
		t.Errorf("normalizeIgnoreFileNames got %v, want %v", got, want)
	}

	if normalizeIgnoreFileNames(nil) != nil {
		t.Error("expected nil for empty slice")
	}
}
