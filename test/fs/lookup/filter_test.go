// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test_lookup

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/timoschenko/dmltools/pkg/fs/lookup"
)

func TestRegexFiltersAndIncludeDiscovery(t *testing.T) {
	now := time.Unix(1700000000, 0)

	t.Run("invalid regex errors", func(t *testing.T) {
		bad := []struct {
			opt  lookup.Option
			flag string
		}{
			{lookup.WithInclude("["), lookup.FilterInclude},
			{lookup.WithExclude("["), lookup.FilterExclude},
			{lookup.WithExtendExclude("["), lookup.FilterExtendExclude},
			{lookup.WithForceExclude("["), lookup.FilterForceExclude},
		}
		for _, tc := range bad {
			_, err := lookup.NewFilter(tc.opt)
			if err == nil || !strings.Contains(err.Error(), tc.flag) ||
				!strings.Contains(err.Error(), "invalid regular expression") {
				t.Errorf("unexpected error for %s: %v", tc.flag, err)
			}
		}
	})

	t.Run("filter matching methods", func(t *testing.T) {
		m := buildMockRepo(now)
		flt, err := lookup.NewFilter(
			lookup.WithFileSystem(m),
			lookup.WithInclude(`\.dml$`),
			lookup.WithExclude(`/venv/`),
			lookup.WithExtendExclude(`/custom/`),
			lookup.WithForceExclude(`generated/.*\.dml$`),
		)
		if err != nil {
			t.Fatalf("NewFilter failed: %v", err)
		}

		if !flt.IsForceExcluded("src/generated/foo.dml") {
			t.Error("expected generated/foo.dml to be force-excluded")
		}
		if flt.IsForceExcluded("src/normal/foo.dml") {
			t.Error("expected normal/foo.dml not to be force-excluded")
		}

		if got := flt.MatchExclude(
			"a/generated/x.dml",
			"a/generated/x.dml",
			false,
		); got != lookup.FilterForceExclude {
			t.Errorf("got %q, want %q", got, lookup.FilterForceExclude)
		}
		if got := flt.MatchExclude(
			"venv/x.dml",
			"/repo/venv/x.dml",
			false,
		); got != lookup.FilterExclude {
			t.Errorf("got %q, want %q", got, lookup.FilterExclude)
		}
		if got := flt.MatchExclude(
			"custom/x.dml",
			"/repo/custom/x.dml",
			false,
		); got != lookup.FilterExtendExclude {
			t.Errorf("got %q, want %q", got, lookup.FilterExtendExclude)
		}

		if !flt.IsIncluded("sub/a.dml", "/repo/sub/a.dml") {
			t.Error("expected .dml file to be included")
		}
		if flt.IsIncluded("sub/a.txt", "/repo/sub/a.txt") {
			t.Error("expected .txt file not to be included")
		}
	})

	t.Run("filters in FindChain and Include fallback", func(t *testing.T) {
		m := buildMockRepo(now)
		m.setDir("/repo/venv", 1, 1000, now)
		m.setFile("/repo/venv/.editorconfig", 0644, 10, 1, 1000, now)
		m.setFile("/repo/.editorconfig", 0644, 10, 1, 1000, now)
		m.setFile("/repo/sub/a.cfg", 0644, 10, 1, 1000, now)
		m.setFile("/repo/sub/b.cfg", 0644, 10, 1, 1000, now)
		m.setFile("/repo/sub/skip.cfg", 0644, 10, 1, 1000, now)
		m.setFile("/repo/sub/other.txt", 0644, 10, 1, 1000, now)

		var skipped []string
		finder, err := lookup.NewFinder(
			lookup.WithFileSystem(m),
			lookup.WithInclude(`\.cfg$`),
			lookup.WithExclude(`/venv/`),
			lookup.WithExtendExclude(`skip\.cfg$`),
			lookup.WithStrictSecurity(true),
			lookup.WithOnSkip(
				func(path string, reason lookup.SkipReason, _ error) {
					if reason == lookup.SkipExcluded {
						skipped = append(skipped, path)
					}
				},
			),
		)
		if err != nil {
			t.Fatalf("NewFinder failed: %v", err)
		}

		chain, err := finder.FindChain(
			"/repo/venv",
			lookup.FileEditorConfig,
		)
		if err != nil {
			t.Fatalf("FindChain failed: %v", err)
		}
		if len(chain) != 1 ||
			chain[0].Path != filepath.Clean("/repo/.editorconfig") {
			t.Fatalf("expected only /repo/.editorconfig, got %+v", chain)
		}
		if !slices.Contains(skipped, filepath.Clean("/repo/venv")) {
			t.Errorf("expected /repo/venv to be reported as SkipExcluded")
		}

		incFiles, err := finder.FindChain("/repo/sub/deep")
		if err != nil {
			t.Fatalf("FindChain with Include failed: %v", err)
		}
		if len(incFiles) != 2 ||
			incFiles[0].Name != "a.cfg" ||
			incFiles[1].Name != "b.cfg" {
			t.Fatalf("unexpected Include discovery results: %+v", incFiles)
		}

		m.setFile("/repo/venv/test.dml", 0644, 5, 1, 1000, now)
		res, found, err := finder.FindNearest(
			"/repo/venv/test.dml",
			lookup.FileEditorConfig,
		)
		if err != nil || found || res.Path != "" {
			t.Fatalf(
				"expected excluded start file to return empty, got (%+v, %v, %v)",
				res,
				found,
				err,
			)
		}
	})
}
