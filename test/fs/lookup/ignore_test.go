// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test_lookup

import (
	"errors"
	iofs "io/fs"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/timoschenko/dmltools/pkg/fs/lookup"
)

func TestIgnoreMatcher_PatternsAndCaching(t *testing.T) {
	t0 := time.Unix(1700000000, 0)

	t.Run("glob patterns escapes comments and CRLF", func(t *testing.T) {
		m := buildMockRepo(t0)
		rules := strings.Join([]string{
			"# comment",
			"   ",
			`\!bang.dml`,
			`\#hash.dml`,
			"trail.dml   ",
			`esc\ .dml`,
			"file[12].dml",
			"only[!1].dml",
			"[]a].dml",
			"unclosed[.dml",
			"char?.dml",
			"**/deep.dml",
			"a/**/b.dml",
			"a**b.dml",
			"box/**",
			"!box/keep.dml",
			"ignored_dir/",
			"!ignored_dir/cannot_rescue.dml",
			"rescued_dir/",
			"!rescued_dir/",
			"rescued_dir/drop.dml",
			"!rescued_dir/keep.dml",
		}, "\r\n") + "\r\n"

		m.setFileWithContent("/repo/.gitignore", []byte(rules), 1, 1000, t0)
		matcher := lookup.NewIgnoreMatcher(lookup.WithFileSystem(m))

		rs, err := matcher.LoadFile(lookup.DiscoveredFile{
			Path: "/repo/.gitignore",
			Dir:  "/repo",
			Name: lookup.FileGitIgnore,
		})
		if err != nil {
			t.Fatalf("LoadFile failed: %v", err)
		}
		stack := []lookup.IgnoreRuleSet{rs}

		shouldIgnore := []string{
			"/repo/!bang.dml",
			"/repo/#hash.dml",
			"/repo/trail.dml",
			"/repo/esc .dml",
			"/repo/file1.dml",
			"/repo/file2.dml",
			"/repo/only2.dml",
			"/repo/].dml",
			"/repo/a.dml",
			"/repo/sub/].dml",
			"/repo/unclosed[.dml",
			"/repo/charx.dml",
			"/repo/deep.dml",
			"/repo/sub/deep.dml",
			"/repo/a/b.dml",
			"/repo/a/x/y/b.dml",
			"/repo/axxxyb.dml",
			"/repo/box/drop.dml",
			"/repo/box/sub/nested.dml",
			"/repo/ignored_dir/file.dml",
			"/repo/ignored_dir/sub/nested.dml",
			"/repo/ignored_dir/cannot_rescue.dml",
			"/repo/rescued_dir/drop.dml",
		}
		for _, p := range shouldIgnore {
			if got := matcher.Match(
				stack,
				p,
				false,
			); got != lookup.FileGitIgnore {
				t.Errorf(
					"Match(%q) = %q, want %q",
					p,
					got,
					lookup.FileGitIgnore,
				)
			}
		}

		shouldNotIgnore := []string{
			"/repo/# comment",
			"/repo/   ",
			"/repo/trail.dml   ",
			"/repo/file3.dml",
			"/repo/only1.dml",
			"/repo/]a.dml",
			"/repo/only/.dml",
			"/repo/charab.dml",
			"/repo/sub/a/x/b.dml",
			"/repo/a/sub/b.dml.extra",
			"/repo/ax/yb.dml",
			"/repo/box",
			"/repo/box/keep.dml",
			"/repo/rescued_dir/keep.dml",
		}
		for _, p := range shouldNotIgnore {
			if got := matcher.Match(stack, p, false); got != "" {
				t.Errorf("Match(%q) = %q, want empty", p, got)
			}
		}
	})

	t.Run("cache hit and invalidation on file mutation", func(t *testing.T) {
		m := buildMockRepo(t0)
		m.setFileWithContent(
			"/repo/.gitignore",
			[]byte("skip.dml\n"),
			1,
			1000,
			t0,
		)
		m.setFileWithContent("/repo/skip.dml", []byte("x"), 1, 1000, t0)
		m.setFileWithContent("/repo/keep.dml", []byte("x"), 1, 1000, t0)

		finder, err := lookup.NewFinder(
			lookup.WithFileSystem(m),
			lookup.WithInclude(`\.dml$`),
			lookup.WithIgnoreFiles(lookup.FileGitIgnore, lookup.FileIgnore),
		)
		if err != nil {
			t.Fatalf("NewFinder failed: %v", err)
		}

		files, err := finder.Collect([]string{"/repo"})
		if err != nil {
			t.Fatalf("Collect 1 failed: %v", err)
		}
		if len(files) != 1 || files[0] != filepath.Clean("/repo/keep.dml") {
			t.Fatalf("Collect 1 got %v, want [/repo/keep.dml]", files)
		}
		if calls := m.totalReadFileCalls(); calls != 1 {
			t.Fatalf("expected 1 ReadFile call on cold cache, got %d", calls)
		}

		// Second Collect without FS changes should reuse cached IgnoreRuleSet.
		files2, err := finder.Collect([]string{"/repo"})
		if err != nil {
			t.Fatalf("Collect 2 failed: %v", err)
		}
		if !slices.Equal(files2, files) {
			t.Fatalf("Collect 2 got %v, want %v", files2, files)
		}
		if calls := m.totalReadFileCalls(); calls != 1 {
			t.Fatalf("expected cached ReadFile count=1, got %d", calls)
		}

		// Modify .gitignore mtime/size -> invalidates cache and re-parses.
		t1 := t0.Add(time.Second)
		m.setFileWithContent(
			"/repo/.gitignore",
			[]byte("keep.dml\n"),
			1,
			1000,
			t1,
		)
		files3, err := finder.Collect([]string{"/repo"})
		if err != nil {
			t.Fatalf("Collect 3 failed: %v", err)
		}
		if len(files3) != 1 || files3[0] != filepath.Clean("/repo/skip.dml") {
			t.Fatalf("Collect 3 got %v, want [/repo/skip.dml]", files3)
		}
		if calls := m.totalReadFileCalls(); calls != 2 {
			t.Fatalf("expected ReadFile count=2 after mutation, got %d", calls)
		}
	})

	t.Run("strict security and relative source in Collect", func(t *testing.T) {
		m := buildMockRepo(t0)
		m.setFileWithContent(
			"/repo/sub/.gitignore",
			[]byte("skip.dml\n"),
			1,
			1000,
			t0,
		)
		m.setFileWithContent(
			"/repo/sub/.ignore",
			[]byte("!skip.dml\n"),
			1,
			1000,
			t0,
		)
		m.setFileWithContent("/repo/sub/skip.dml", []byte("x"), 1, 1000, t0)
		m.setFileWithContent(
			"/repo/sub/deep/main.dml",
			[]byte("x"),
			1,
			1000,
			t0,
		)
		m.cwd = "/repo"

		matcher := lookup.NewIgnoreMatcher(lookup.WithFileSystem(m))
		finder, err := lookup.NewFinder(
			lookup.WithFileSystem(m),
			lookup.WithInclude(`\.dml$`),
			lookup.WithIgnoreMatcher(matcher),
			lookup.WithIgnoreFiles(lookup.FileGitIgnore, lookup.FileIgnore),
		)
		if err != nil {
			t.Fatalf("NewFinder failed: %v", err)
		}

		files, err := finder.Collect([]string{"."})
		if err != nil {
			t.Fatalf("Collect('.') failed: %v", err)
		}
		want := []string{
			filepath.Join("sub", "deep", "main.dml"),
			filepath.Join("sub", "skip.dml"),
		}
		slices.Sort(want)
		if !slices.Equal(files, want) {
			t.Fatalf("Collect('.') got %v, want %v", files, want)
		}

		// FIFO as .gitignore in strict mode must surface *SecurityError during Collect.
		m.setFile(
			"/repo/sub/deep/.gitignore",
			iofs.ModeNamedPipe|0644,
			0,
			1,
			1000,
			t0,
		)
		strictFinder, err := lookup.NewFinder(
			lookup.WithFileSystem(m),
			lookup.WithInclude(`\.dml$`),
			lookup.WithStrictSecurity(true),
			lookup.WithIgnoreFiles(lookup.FileGitIgnore),
		)
		if err != nil {
			t.Fatalf("NewFinder strict failed: %v", err)
		}
		_, err = strictFinder.Collect([]string{"/repo"})
		var secErr *lookup.SecurityError
		if !errors.As(err, &secErr) || secErr.Reason != lookup.SkipNotRegular {
			t.Fatalf("expected SkipNotRegular SecurityError, got %v", err)
		}
	})

	t.Run("collect starting inside ignored directory", func(t *testing.T) {
		m := buildMockRepo(t0)
		m.setDir("/repo/ignored_dir", 1, 1000, t0)
		m.setFileWithContent(
			"/repo/.gitignore",
			[]byte("ignored_dir/\n"),
			1,
			1000,
			t0,
		)
		m.setFileWithContent(
			"/repo/ignored_dir/main.dml",
			[]byte("x"),
			1,
			1000,
			t0,
		)

		type ignoreEvent struct {
			path   string
			isDir  bool
			filter string
		}
		var events []ignoreEvent

		finder, err := lookup.NewFinder(
			lookup.WithFileSystem(m),
			lookup.WithInclude(`\.dml$`),
			lookup.WithIgnoreFiles(lookup.FileGitIgnore, lookup.FileIgnore),
			lookup.WithOnIgnore(
				func(path string, isDir bool, filterName string) {
					events = append(
						events,
						ignoreEvent{path, isDir, filterName},
					)
				},
			),
		)
		if err != nil {
			t.Fatalf("NewFinder failed: %v", err)
		}

		files, err := finder.Collect([]string{"/repo/ignored_dir"})
		if err != nil {
			t.Fatalf("Collect failed: %v", err)
		}
		if len(files) != 0 {
			t.Fatalf("expected no files from ignored_dir, got %v", files)
		}
		wantEvent := ignoreEvent{
			path:   filepath.Clean("/repo/ignored_dir"),
			isDir:  true,
			filter: lookup.FileGitIgnore,
		}
		if !slices.Contains(events, wantEvent) {
			t.Fatalf("expected ignore event %+v, got %+v", wantEvent, events)
		}
	})
}

func TestFinderMatchIgnore(t *testing.T) {
	t0 := time.Unix(1700000000, 0)
	m := buildMockRepo(t0)
	m.setFileWithContent(
		"/repo/.gitignore",
		[]byte("*.gen.dml\nbuild/\n"),
		1,
		1000,
		t0,
	)
	m.setFileWithContent(
		"/repo/sub/.gitignore",
		[]byte("!keep.gen.dml\n"),
		1,
		1000,
		t0,
	)
	m.setDir("/repo/build", 1, 1000, t0)
	m.setFile("/repo/build/a.dml", 0644, 10, 1, 1000, t0)
	m.setFile("/repo/sub/skip.gen.dml", 0644, 10, 1, 1000, t0)
	m.setFile("/repo/sub/keep.gen.dml", 0644, 10, 1, 1000, t0)
	m.setFile("/repo/sub/deep/main.dml", 0644, 10, 1, 1000, t0)

	finder, err := lookup.NewFinder(
		lookup.WithFileSystem(m),
		lookup.WithIgnoreFiles(lookup.FileGitIgnore, lookup.FileIgnore),
	)
	if err != nil {
		t.Fatalf("NewFinder failed: %v", err)
	}

	cases := []struct {
		path string
		want string
	}{
		{"/repo/sub/deep/main.dml", ""},
		{"/repo/sub/keep.gen.dml", ""},
		{"/repo/sub/skip.gen.dml", lookup.FileGitIgnore},
		{"/repo/build/a.dml", lookup.FileGitIgnore},
	}
	for _, tc := range cases {
		got, matchErr := finder.MatchIgnore(tc.path)
		if matchErr != nil {
			t.Fatalf("MatchIgnore(%s) failed: %v", tc.path, matchErr)
		}
		if got != tc.want {
			t.Errorf("MatchIgnore(%s) = %q, want %q", tc.path, got, tc.want)
		}
	}
}
