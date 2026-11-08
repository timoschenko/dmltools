// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func createTestTree(t *testing.T, baseDir string, files []string) {
	t.Helper()
	for _, f := range files {
		fullPath := filepath.Join(baseDir, filepath.FromSlash(f))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0750); err != nil {
			t.Fatalf("failed to create dir for %s: %v", f, err)
		}
		if err := os.WriteFile(fullPath, []byte("content"), 0600); err != nil {
			t.Fatalf("failed to write file %s: %v", f, err)
		}
	}
}

func TestCollector_Defaults(t *testing.T) {
	tmpDir := t.TempDir()
	createTestTree(t, tmpDir, []string{
		"root.dml",
		"notes.txt",
		".git/hook.dml",
		"venv/lib.dml",
		"linux64/bin.dml",
		"sub/nested.dml",
		"sub/ignored.go",
	})

	cfg := &CLIConfig{
		Include: DefaultIncludePattern,
		Exclude: DefaultExcludePattern,
	}
	collector, err := NewFileCollector(cfg)
	if err != nil {
		t.Fatalf("unexpected error creating collector: %v", err)
	}

	files, err := collector.Collect([]string{tmpDir})
	if err != nil {
		t.Fatalf("unexpected collect error: %v", err)
	}

	expected := []string{
		filepath.Join(tmpDir, "root.dml"),
		filepath.Join(tmpDir, "sub", "nested.dml"),
	}
	slices.Sort(expected)

	if !slices.Equal(files, expected) {
		t.Fatalf("expected files %v, got %v", expected, files)
	}
}

func TestCollector_CustomInclude(t *testing.T) {
	tmpDir := t.TempDir()
	createTestTree(t, tmpDir, []string{
		"file.dml",
		"file.dml14",
		"file.txt",
	})

	// Case 1: Custom regex
	cfg := &CLIConfig{
		Include: `\.dml14$`,
		Exclude: DefaultExcludePattern,
	}
	collector, err := NewFileCollector(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, err := collector.Collect([]string{tmpDir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{filepath.Join(tmpDir, "file.dml14")}
	if !slices.Equal(files, expected) {
		t.Fatalf("expected %v, got %v", expected, files)
	}

	// Case 2: Empty include includes all file names
	cfgEmpty := &CLIConfig{
		Include: "",
		Exclude: DefaultExcludePattern,
	}
	collectorEmpty, err := NewFileCollector(cfgEmpty)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	filesEmpty, err := collectorEmpty.Collect([]string{tmpDir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filesEmpty) != 3 {
		t.Fatalf(
			"expected 3 files with empty include, got %d: %v",
			len(filesEmpty),
			filesEmpty,
		)
	}
}

func TestCollector_CustomExclude(t *testing.T) {
	tmpDir := t.TempDir()
	createTestTree(t, tmpDir, []string{
		"src/a.dml",
		"skip/b.dml",
		"venv/c.dml",
	})

	// Custom exclude overrides default exclusions
	cfg := &CLIConfig{
		Include: DefaultIncludePattern,
		Exclude: `/skip/`,
	}
	collector, err := NewFileCollector(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, err := collector.Collect([]string{tmpDir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// venv/c.dml should now be included because default exclude was overridden
	expected := []string{
		filepath.Join(tmpDir, "src", "a.dml"),
		filepath.Join(tmpDir, "venv", "c.dml"),
	}
	slices.Sort(expected)
	if !slices.Equal(files, expected) {
		t.Fatalf("expected %v, got %v", expected, files)
	}
}

func TestCollector_ExtendExclude(t *testing.T) {
	tmpDir := t.TempDir()
	createTestTree(t, tmpDir, []string{
		"src/a.dml",
		"custom/b.dml",
		"venv/c.dml",
	})

	// ExtendExclude keeps default exclusions and adds /custom/
	cfg := &CLIConfig{
		Include:       DefaultIncludePattern,
		Exclude:       DefaultExcludePattern,
		ExtendExclude: `/custom/`,
	}
	collector, err := NewFileCollector(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, err := collector.Collect([]string{tmpDir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		filepath.Join(tmpDir, "src", "a.dml"),
	}
	if !slices.Equal(files, expected) {
		t.Fatalf("expected %v, got %v", expected, files)
	}
}

func TestCollector_ForceExclude(t *testing.T) {
	tmpDir := t.TempDir()
	createTestTree(t, tmpDir, []string{
		"keep.dml",
		"skip_gen.dml",
		"sub/ignored.dml",
	})

	cfg := &CLIConfig{
		Include:      DefaultIncludePattern,
		Exclude:      DefaultExcludePattern,
		ForceExclude: `(skip_gen|_ignored)`,
	}
	collector, err := NewFileCollector(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test 1: Recursive directory search
	files, err := collector.Collect([]string{tmpDir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{
		filepath.Join(tmpDir, "keep.dml"),
		filepath.Join(tmpDir, "sub", "ignored.dml"),
	}
	slices.Sort(expected)
	if !slices.Equal(files, expected) {
		t.Fatalf("expected %v, got %v", expected, files)
	}

	// Test 2: Explicitly passed files (force-exclude skips even explicit files)
	keepFile := filepath.Join(tmpDir, "keep.dml")
	skipFile := filepath.Join(tmpDir, "skip_gen.dml")
	explicitFiles, err := collector.Collect([]string{keepFile, skipFile})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedExplicit := []string{keepFile}
	if !slices.Equal(explicitFiles, expectedExplicit) {
		t.Fatalf("expected %v, got %v", expectedExplicit, explicitFiles)
	}
}

func TestCollector_ExplicitFilesBypassExcludeAndInclude(t *testing.T) {
	tmpDir := t.TempDir()
	createTestTree(t, tmpDir, []string{
		"venv/test.txt",
	})

	cfg := &CLIConfig{
		Include: DefaultIncludePattern,
		Exclude: DefaultExcludePattern,
	}
	collector, err := NewFileCollector(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	target := filepath.Join(tmpDir, "venv", "test.txt")
	files, err := collector.Collect([]string{target})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{target}
	if !slices.Equal(files, expected) {
		t.Fatalf(
			"expected explicit file %v to bypass recursive filters, got %v",
			expected,
			files,
		)
	}
}

func TestCollector_DeduplicationAndSorting(t *testing.T) {
	tmpDir := t.TempDir()
	createTestTree(t, tmpDir, []string{
		"b.dml",
		"a.dml",
	})

	cfg := &CLIConfig{
		Include: DefaultIncludePattern,
		Exclude: DefaultExcludePattern,
	}
	collector, err := NewFileCollector(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fileA := filepath.Join(tmpDir, "a.dml")
	fileB := filepath.Join(tmpDir, "b.dml")

	// Pass redundant directory and explicit files out of order
	files, err := collector.Collect([]string{fileB, tmpDir, fileA})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{fileA, fileB}
	if !slices.Equal(files, expected) {
		t.Fatalf("expected deduplicated and sorted %v, got %v", expected, files)
	}
}

func TestCollector_InvalidRegex(t *testing.T) {
	cases := []struct {
		name string
		cfg  *CLIConfig
	}{
		{"invalid include", &CLIConfig{Include: "["}},
		{"invalid exclude", &CLIConfig{Exclude: "["}},
		{"invalid extend-exclude", &CLIConfig{ExtendExclude: "["}},
		{"invalid force-exclude", &CLIConfig{ForceExclude: "["}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewFileCollector(tc.cfg)
			if err == nil {
				t.Fatalf("expected error for %s, got nil", tc.name)
			}
			if !strings.Contains(err.Error(), "invalid regular expression") {
				t.Errorf(
					"expected 'invalid regular expression' in error, got: %v",
					err,
				)
			}
		})
	}
}

func TestCollector_NonExistentPath(t *testing.T) {
	cfg := &CLIConfig{
		Include: DefaultIncludePattern,
		Exclude: DefaultExcludePattern,
	}
	collector, err := NewFileCollector(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	missing := filepath.Join(t.TempDir(), "does_not_exist.dml")
	_, err = collector.Collect([]string{missing})
	if err == nil {
		t.Fatalf("expected error for non-existent path, got nil")
	}
	if !strings.Contains(err.Error(), "reading file") {
		t.Errorf("expected 'reading file' in error, got: %v", err)
	}
}

func TestCollector_IsForceExcluded(t *testing.T) {
	cfg := &CLIConfig{
		ForceExclude: `generated/.*\.dml$`,
	}
	collector, err := NewFileCollector(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !collector.IsForceExcluded("path/to/generated/foo.dml") {
		t.Errorf("expected path to be force-excluded")
	}
	if collector.IsForceExcluded("path/to/normal/foo.dml") {
		t.Errorf("expected normal path not to be force-excluded")
	}

	// When force-exclude is not set
	nilCollector, _ := NewFileCollector(&CLIConfig{})
	if nilCollector.IsForceExcluded("any/file.dml") {
		t.Errorf("expected false when force-exclude is not set")
	}
}

func TestCollector_AnchoredIncludeAndExclude(t *testing.T) {
	tmpDir := t.TempDir()
	createTestTree(t, tmpDir, []string{
		"src/main.dml",
		"src/sub/helper.dml",
		"test/test.dml",
		"other/lib.dml",
	})

	// Case 1: Anchored include regex (^src/.*)
	cfgInclude := &CLIConfig{
		Include: `^src/.*`,
		Exclude: DefaultExcludePattern,
	}
	collectorInc, err := NewFileCollector(cfgInclude)
	if err != nil {
		t.Fatalf("unexpected error creating collector: %v", err)
	}

	filesInc, err := collectorInc.Collect([]string{tmpDir})
	if err != nil {
		t.Fatalf("unexpected collect error: %v", err)
	}
	expectedInc := []string{
		filepath.Join(tmpDir, "src", "main.dml"),
		filepath.Join(tmpDir, "src", "sub", "helper.dml"),
	}
	slices.Sort(expectedInc)
	if !slices.Equal(filesInc, expectedInc) {
		t.Fatalf("expected files %v, got %v", expectedInc, filesInc)
	}

	// Case 2: Anchored exclude regex (^test/)
	cfgExclude := &CLIConfig{
		Include: DefaultIncludePattern,
		Exclude: `^test/`,
	}
	collectorExc, err := NewFileCollector(cfgExclude)
	if err != nil {
		t.Fatalf("unexpected error creating collector: %v", err)
	}

	filesExc, err := collectorExc.Collect([]string{tmpDir})
	if err != nil {
		t.Fatalf("unexpected collect error: %v", err)
	}
	expectedExc := []string{
		filepath.Join(tmpDir, "other", "lib.dml"),
		filepath.Join(tmpDir, "src", "main.dml"),
		filepath.Join(tmpDir, "src", "sub", "helper.dml"),
	}
	slices.Sort(expectedExc)
	if !slices.Equal(filesExc, expectedExc) {
		t.Fatalf("expected files %v, got %v", expectedExc, filesExc)
	}
}

func TestCollector_RelativeSourceDefaultExclusions(t *testing.T) {
	tmpDir := t.TempDir()
	createTestTree(t, tmpDir, []string{
		"root.dml",
		"venv/ignored.dml",
		".git/hook.dml",
		"sub/nested.dml",
	})

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir to %s: %v", tmpDir, err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origWd)
	})

	cfg := &CLIConfig{
		Include: DefaultIncludePattern,
		Exclude: DefaultExcludePattern,
	}
	collector, err := NewFileCollector(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, err := collector.Collect([]string{"."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		"root.dml",
		filepath.Join("sub", "nested.dml"),
	}
	slices.Sort(expected)
	if !slices.Equal(files, expected) {
		t.Fatalf("expected %v, got %v", expected, files)
	}
}
