// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test_fs

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/timoschenko/dmltools/pkg/fs"
	"github.com/timoschenko/dmltools/pkg/fs/lookup"
)

func writeTestFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestLookup_RealRepoHierarchyCollectAndMutations(t *testing.T) {
	tmpDir := t.TempDir()

	outerRepo := filepath.Join(tmpDir, "outer")
	innerRepo := filepath.Join(outerRepo, "inner")
	subDir := filepath.Join(innerRepo, "pkg", "dev")
	targetFile := filepath.Join(subDir, "device.dml")

	if err := os.MkdirAll(filepath.Join(outerRepo, ".git"), 0750); err != nil {
		t.Fatalf("mkdir outer .git: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(innerRepo, ".git"), 0750); err != nil {
		t.Fatalf("mkdir inner .git: %v", err)
	}
	writeTestFile(t, filepath.Join(outerRepo, ".editorconfig"), "root = true\n")
	writeTestFile(t, filepath.Join(innerRepo, ".editorconfig"), "indent = 4\n")
	writeTestFile(t, filepath.Join(innerRepo, ".gitignore"), "*.bak.dml\n")
	writeTestFile(
		t,
		filepath.Join(innerRepo, "pkg", ".ignore"),
		"generated/\n",
	)
	writeTestFile(t, targetFile, "dml 1.4;\ndevice dev;\n")
	writeTestFile(t, filepath.Join(innerRepo, "skip.bak.dml"), "dml 1.4;\n")
	writeTestFile(
		t,
		filepath.Join(innerRepo, "pkg", "generated", "gen.dml"),
		"dml 1.4;\n",
	)
	writeTestFile(t, filepath.Join(innerRepo, "venv", "skip.dml"), "dml 1.4;\n")

	finder, err := lookup.NewFinder(
		lookup.WithInclude(`\.dml$`),
		lookup.WithExclude(fs.DefaultExcludePattern),
		lookup.WithIgnoreFiles(lookup.FileGitIgnore, lookup.FileIgnore),
	)
	if err != nil {
		t.Fatalf("lookup.NewFinder: %v", err)
	}

	collected, err := finder.Collect([]string{innerRepo})
	if err != nil {
		t.Fatalf("Collect: %v", err)
	}
	if !slices.Equal(collected, []string{targetFile}) {
		t.Fatalf("Collect got %v, want [%s]", collected, targetFile)
	}

	boundary, err := finder.Cache().ResolveBoundary(targetFile)
	if err != nil {
		t.Fatalf("ResolveBoundary: %v", err)
	}
	if boundary.Reason != lookup.BoundaryVCSRoot ||
		boundary.RootDir != innerRepo {
		t.Fatalf("expected innerRepo VCS boundary, got %+v", boundary)
	}

	chain, err := finder.FindChain(
		targetFile,
		lookup.FileEditorConfig,
		lookup.FileGitIgnore,
		lookup.FileIgnore,
	)
	if err != nil {
		t.Fatalf("FindChain: %v", err)
	}

	wantPaths := []string{
		filepath.Join(innerRepo, lookup.FileEditorConfig),
		filepath.Join(innerRepo, lookup.FileGitIgnore),
		filepath.Join(innerRepo, "pkg", lookup.FileIgnore),
	}
	gotPaths := make([]string, len(chain))
	for i, f := range chain {
		gotPaths[i] = f.Path
	}
	if !slices.Equal(gotPaths, wantPaths) {
		t.Fatalf("got paths %v, want %v", gotPaths, wantPaths)
	}

	// Mutate local filesystem: add a deeper .editorconfig, edit existing, delete .ignore.
	deepEditorCfg := filepath.Join(subDir, lookup.FileEditorConfig)
	writeTestFile(t, deepEditorCfg, "indent = 2\n")
	writeTestFile(
		t,
		filepath.Join(innerRepo, lookup.FileEditorConfig),
		"indent = 8\n# updated\n",
	)
	if err := os.Remove(
		filepath.Join(innerRepo, "pkg", lookup.FileIgnore),
	); err != nil {
		t.Fatalf("remove .ignore: %v", err)
	}

	nearest, found, err := finder.FindNearest(
		targetFile,
		lookup.FileEditorConfig,
	)
	if err != nil {
		t.Fatalf("FindNearest after mutation: %v", err)
	}
	if !found || nearest.Path != deepEditorCfg {
		t.Fatalf(
			"expected nearest %s, got (%+v, %v)",
			deepEditorCfg,
			nearest,
			found,
		)
	}

	chainAfter, err := finder.FindChain(
		targetFile,
		lookup.FileEditorConfig,
		lookup.FileIgnore,
	)
	if err != nil {
		t.Fatalf("FindChain after mutation: %v", err)
	}
	if len(chainAfter) != 2 ||
		chainAfter[0].Size != int64(len("indent = 8\n# updated\n")) ||
		chainAfter[1].Path != deepEditorCfg {
		t.Fatalf("unexpected chain after mutation: %+v", chainAfter)
	}
}

func TestLookup_RealSymlinksAndRelativePaths(t *testing.T) {
	tmpDir := t.TempDir()
	repo := filepath.Join(tmpDir, "repo")
	sub := filepath.Join(repo, "sub")
	outsideFile := filepath.Join(tmpDir, "outside.cfg")
	inRepoTarget := filepath.Join(repo, "shared.cfg")

	if err := os.MkdirAll(filepath.Join(repo, ".git"), 0750); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}
	if err := os.MkdirAll(sub, 0750); err != nil {
		t.Fatalf("mkdir sub: %v", err)
	}
	writeTestFile(t, outsideFile, "secret=1\n")
	writeTestFile(t, inRepoTarget, "safe=1\n")

	if err := os.Symlink(
		outsideFile,
		filepath.Join(repo, lookup.FileEditorConfig),
	); err != nil {
		t.Fatalf("symlink outside: %v", err)
	}
	if err := os.Symlink(
		inRepoTarget,
		filepath.Join(sub, lookup.FileEditorConfig),
	); err != nil {
		t.Fatalf("symlink inside: %v", err)
	}
	if err := os.Symlink(
		filepath.Join(repo, "missing.cfg"),
		filepath.Join(sub, lookup.FileIgnore),
	); err != nil {
		t.Fatalf("symlink broken: %v", err)
	}

	var skips []lookup.SkipReason
	finder, err := lookup.NewFinder(
		lookup.WithOnSkip(func(_ string, r lookup.SkipReason, _ error) {
			skips = append(skips, r)
		}),
	)
	if err != nil {
		t.Fatalf("lookup.NewFinder: %v", err)
	}

	t.Chdir(sub)

	files, err := finder.FindChain(
		".",
		lookup.FileEditorConfig,
		lookup.FileIgnore,
	)
	if err != nil {
		t.Fatalf("FindChain: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 safe file, got %+v", files)
	}
	evalTarget, err := filepath.EvalSymlinks(inRepoTarget)
	if err != nil {
		t.Fatalf("EvalSymlinks: %v", err)
	}
	if !files[0].IsSymlink || files[0].ResolvedPath != evalTarget {
		t.Errorf(
			"unexpected symlink result: %+v (want %s)",
			files[0],
			evalTarget,
		)
	}
	if !slices.Contains(skips, lookup.SkipSymlinkEscape) ||
		!slices.Contains(skips, lookup.SkipBrokenSymlink) {
		t.Errorf("expected escape and broken symlink skips, got %v", skips)
	}

	strictFinder, err := lookup.NewFinder(lookup.WithStrictSecurity(true))
	if err != nil {
		t.Fatalf("NewFinder strict: %v", err)
	}
	_, err = strictFinder.FindChain(".", lookup.FileEditorConfig)
	var secErr *lookup.SecurityError
	if !errors.As(err, &secErr) ||
		secErr.Reason != lookup.SkipSymlinkEscape {
		t.Fatalf("expected SkipSymlinkEscape SecurityError, got %v", err)
	}
}
