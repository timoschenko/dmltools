// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test_lookup

import (
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/timoschenko/dmltools/pkg/fs/lookup"
)

func TestSOLIDCompositionAndInjection(t *testing.T) {
	now := time.Unix(1700000000, 0)
	m := buildMockRepo(now)
	m.setFile("/repo/.editorconfig", 0644, 20, 1, 1000, now)
	m.setFile("/repo/sub/main.dml", 0644, 30, 1, 1000, now)

	cache := lookup.NewCache(lookup.WithFileSystem(m))
	filter, err := lookup.NewFilter(
		lookup.WithFileSystem(m),
		lookup.WithInclude(`\.dml$`),
		lookup.WithExclude(`/\.git/`),
	)
	if err != nil {
		t.Fatalf("NewFilter failed: %v", err)
	}
	if !filter.HasInclude() {
		t.Fatal("expected filter.HasInclude() to be true")
	}

	finder1, err := lookup.NewFinder(
		lookup.WithCache(cache),
		lookup.WithFilter(filter),
	)
	if err != nil {
		t.Fatalf("NewFinder 1 failed: %v", err)
	}
	finder2, err := lookup.NewFinder(
		lookup.WithCache(cache),
		lookup.WithFilter(filter),
	)
	if err != nil {
		t.Fatalf("NewFinder 2 failed: %v", err)
	}

	if finder1.Cache() != cache || finder1.Filter() != filter {
		t.Fatal("expected finder1 to expose injected Cache and Filter")
	}

	collected, err := finder1.Collect([]string{"/repo"})
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}
	if len(collected) != 1 ||
		collected[0] != filepath.Clean("/repo/sub/main.dml") {
		t.Fatalf("unexpected collected files: %v", collected)
	}

	// finder2 shares the same Cache instance as finder1.
	chain, err := finder2.FindChain(
		"/repo/sub/main.dml",
		lookup.FileEditorConfig,
	)
	if err != nil {
		t.Fatalf("FindChain failed: %v", err)
	}
	if len(chain) != 1 ||
		chain[0].Path != filepath.Clean("/repo/.editorconfig") {
		t.Fatalf("unexpected chain: %+v", chain)
	}
	if cache.Stats().DirHits == 0 {
		t.Errorf("expected shared cache hits between finder1 and finder2")
	}
}

func TestFindChainAndFindNearest_OrderingAndNames(t *testing.T) {
	now := time.Unix(1700000000, 0)
	m := buildMockRepo(now)
	m.setFile("/outside/.editorconfig", 0644, 10, 1, 1000, now)
	m.setFile("/repo/.editorconfig", 0644, 20, 1, 1000, now)
	m.setFile("/repo/.gitignore", 0644, 15, 1, 1000, now)
	m.setFile("/repo/sub/.editorconfig", 0644, 25, 1, 1000, now)
	m.setFile("/repo/sub/deep/main.dml", 0644, 40, 1, 1000, now)

	finder, err := lookup.NewFinder(lookup.WithFileSystem(m))
	if err != nil {
		t.Fatalf("NewFinder failed: %v", err)
	}

	chain, err := finder.FindChain(
		"/repo/sub/deep/main.dml",
		lookup.FileEditorConfig,
		lookup.FileGitIgnore,
	)
	if err != nil {
		t.Fatalf("FindChain failed: %v", err)
	}
	wantPaths := []string{
		filepath.Clean("/repo/.editorconfig"),
		filepath.Clean("/repo/.gitignore"),
		filepath.Clean("/repo/sub/.editorconfig"),
	}
	if len(chain) != len(wantPaths) {
		t.Fatalf(
			"got %d files (%v), want %d",
			len(chain),
			chain,
			len(wantPaths),
		)
	}
	for i, want := range wantPaths {
		if chain[i].Path != want {
			t.Errorf("chain[%d].Path = %q, want %q", i, chain[i].Path, want)
		}
	}

	nearest, found, err := finder.FindNearest(
		"/repo/sub/deep",
		lookup.FileEditorConfig,
		lookup.FileGitIgnore,
	)
	if err != nil {
		t.Fatalf("FindNearest failed: %v", err)
	}
	if !found {
		t.Fatal("expected FindNearest to find a file")
	}
	if nearest.Path != filepath.Clean("/repo/sub/.editorconfig") {
		t.Errorf(
			"nearest.Path = %q, want /repo/sub/.editorconfig",
			nearest.Path,
		)
	}
}

func TestFinderCollect(t *testing.T) {
	now := time.Unix(1700000000, 0)
	m := newMockFS("/ws", 1000)
	m.setDir("/", 1, 0, now)
	m.setDir("/ws", 1, 1000, now)
	m.setDir("/ws/.git", 1, 1000, now)
	m.setFile("/ws/.git/hook.dml", 0644, 5, 1, 1000, now)
	m.setDir("/ws/venv", 1, 1000, now)
	m.setFile("/ws/venv/lib.dml", 0644, 5, 1, 1000, now)
	m.setDir("/ws/sub", 1, 1000, now)
	m.setFile("/ws/sub/a.dml", 0644, 5, 1, 1000, now)
	m.setFile("/ws/sub/b.txt", 0644, 5, 1, 1000, now)
	m.setFile("/ws/sub/skip_gen.dml", 0644, 5, 1, 1000, now)
	m.setFile("/ws/root.dml", 0644, 5, 1, 1000, now)
	m.setSymlink("/ws/dirlink", "/ws/sub", 1, 1000, now)
	m.setSymlink("/ws/broken.dml", "/ws/missing.dml", 1, 1000, now)

	type ignoredEvent struct {
		path   string
		isDir  bool
		filter string
	}
	var ignored []ignoredEvent

	finder, err := lookup.NewFinder(
		lookup.WithFileSystem(m),
		lookup.WithInclude(`\.dml$`),
		lookup.WithExclude(`/(\.git|venv)/`),
		lookup.WithForceExclude(`skip_gen`),
	)
	if err != nil {
		t.Fatalf("NewFinder failed: %v", err)
	}
	finder.SetOnIgnore(func(path string, isDir bool, filterName string) {
		ignored = append(ignored, ignoredEvent{path, isDir, filterName})
	})

	files, err := finder.Collect([]string{
		"/ws/root.dml",
		"/ws",
		"/ws/venv/lib.dml",
		"/ws/sub/skip_gen.dml",
	})
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	// /ws/venv/lib.dml passed explicitly bypasses --exclude, while skip_gen.dml
	// is skipped by --force-exclude. Directory symlink and broken symlink are skipped.
	want := []string{
		filepath.Clean("/ws/root.dml"),
		filepath.Clean("/ws/sub/a.dml"),
		filepath.Clean("/ws/venv/lib.dml"),
	}
	if !slices.Equal(files, want) {
		t.Fatalf("got %v, want %v", files, want)
	}

	if !slices.Contains(ignored, ignoredEvent{
		path:   filepath.Clean("/ws/sub/b.txt"),
		isDir:  false,
		filter: lookup.FilterInclude,
	}) {
		t.Errorf(
			"expected b.txt to trigger %s ignore event, got %+v",
			lookup.FilterInclude,
			ignored,
		)
	}
	if !slices.Contains(ignored, ignoredEvent{
		path:   filepath.Clean("/ws/venv"),
		isDir:  true,
		filter: lookup.FilterExclude,
	}) {
		t.Errorf(
			"expected venv dir to trigger %s ignore event, got %+v",
			lookup.FilterExclude,
			ignored,
		)
	}

	if _, err := finder.Collect([]string{"/ws/nonexistent.dml"}); err == nil {
		t.Fatal("expected error for non-existent source")
	}

	m.mu.Lock()
	subGitProbes := m.statCalls[filepath.Clean("/ws/sub/.git")]
	subHgProbes := m.statCalls[filepath.Clean("/ws/sub/.hg")]
	m.mu.Unlock()
	if subGitProbes != 0 || subHgProbes != 0 {
		t.Errorf(
			"expected no VCS marker probes on subdirectory during Collect, got .git=%d .hg=%d",
			subGitProbes,
			subHgProbes,
		)
	}
}
