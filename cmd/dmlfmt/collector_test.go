// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/timoschenko/dmltools/pkg/fs/lookup"
)

func newTestFinder(cfg *CLIConfig) (lookup.Finder, error) {
	return newFinderFromConfig(cfg)
}

func createTestTree(t *testing.T, baseDir string, files []string) {
	t.Helper()
	withContent := make(map[string]string, len(files))
	for _, f := range files {
		withContent[f] = "content"
	}
	createTestTreeWithContent(t, baseDir, withContent)
}

func createTestTreeWithContent(
	t *testing.T,
	baseDir string,
	files map[string]string,
) {
	t.Helper()
	for f, content := range files {
		if strings.HasSuffix(f, "/") {
			dirPath := filepath.Join(baseDir, filepath.FromSlash(f))
			if err := os.MkdirAll(dirPath, 0750); err != nil {
				t.Fatalf("failed to create dir %s: %v", f, err)
			}
			continue
		}
		fullPath := filepath.Join(baseDir, filepath.FromSlash(f))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0750); err != nil {
			t.Fatalf("failed to create dir for %s: %v", f, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0600); err != nil {
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
	finder, err := newTestFinder(cfg)
	if err != nil {
		t.Fatalf("unexpected error creating finder: %v", err)
	}

	files, err := finder.Collect([]string{tmpDir})
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
	finder, err := newTestFinder(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, err := finder.Collect([]string{tmpDir})
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
	finderEmpty, err := newTestFinder(cfgEmpty)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	filesEmpty, err := finderEmpty.Collect([]string{tmpDir})
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
	finder, err := newTestFinder(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, err := finder.Collect([]string{tmpDir})
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
	finder, err := newTestFinder(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, err := finder.Collect([]string{tmpDir})
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
	finder, err := newTestFinder(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test 1: Recursive directory search
	files, err := finder.Collect([]string{tmpDir})
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
	explicitFiles, err := finder.Collect([]string{keepFile, skipFile})
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
	finder, err := newTestFinder(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	target := filepath.Join(tmpDir, "venv", "test.txt")
	files, err := finder.Collect([]string{target})
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
	finder, err := newTestFinder(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fileA := filepath.Join(tmpDir, "a.dml")
	fileB := filepath.Join(tmpDir, "b.dml")

	// Pass redundant directory and explicit files out of order
	files, err := finder.Collect([]string{fileB, tmpDir, fileA})
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
			_, err := newTestFinder(tc.cfg)
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
	finder, err := newTestFinder(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	missing := filepath.Join(t.TempDir(), "does_not_exist.dml")
	_, err = finder.Collect([]string{missing})
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
	finder, err := newTestFinder(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !finder.Filter().IsForceExcluded("path/to/generated/foo.dml") {
		t.Errorf("expected path to be force-excluded")
	}
	if finder.Filter().IsForceExcluded("path/to/normal/foo.dml") {
		t.Errorf("expected normal path not to be force-excluded")
	}

	// When force-exclude is not set
	nilFinder, _ := newTestFinder(&CLIConfig{})
	if nilFinder.Filter().IsForceExcluded("any/file.dml") {
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
	finderInc, err := newTestFinder(cfgInclude)
	if err != nil {
		t.Fatalf("unexpected error creating finder: %v", err)
	}

	filesInc, err := finderInc.Collect([]string{tmpDir})
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
	finderExc, err := newTestFinder(cfgExclude)
	if err != nil {
		t.Fatalf("unexpected error creating finder: %v", err)
	}

	filesExc, err := finderExc.Collect([]string{tmpDir})
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

	t.Chdir(tmpDir)

	cfg := &CLIConfig{
		Include: DefaultIncludePattern,
		Exclude: DefaultExcludePattern,
	}
	finder, err := newTestFinder(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, err := finder.Collect([]string{"."})
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

func TestCollector_GitIgnore_Basic(t *testing.T) {
	tmpDir := t.TempDir()
	createTestTreeWithContent(t, tmpDir, map[string]string{
		".git/": "",
		".gitignore": strings.Join([]string{
			"# comment should be ignored",
			"exact_skip.dml",
			"*_gen.dml",
			"**/*.bak.dml",
			"ignored_dir/",
			"/root_only.dml",
		}, "\n") + "\n",
		"keep.dml":                   "content",
		"exact_skip.dml":             "content",
		"dev_gen.dml":                "content",
		"top.bak.dml":                "content",
		"root_only.dml":              "content",
		"sub/keep_sub.dml":           "content",
		"sub/nested_gen.dml":         "content",
		"sub/deep/inner.bak.dml":     "content",
		"sub/root_only.dml":          "content",
		"ignored_dir/skipped.dml":    "content",
		"sub/ignored_dir/nested.dml": "content",
	})

	cfg := &CLIConfig{
		Include: DefaultIncludePattern,
		Exclude: DefaultExcludePattern,
	}
	finder, err := newTestFinder(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, err := finder.Collect([]string{tmpDir})
	if err != nil {
		t.Fatalf("unexpected collect error: %v", err)
	}

	expected := []string{
		filepath.Join(tmpDir, "keep.dml"),
		filepath.Join(tmpDir, "sub", "keep_sub.dml"),
		filepath.Join(tmpDir, "sub", "root_only.dml"),
	}
	slices.Sort(expected)
	if !slices.Equal(files, expected) {
		t.Fatalf("expected %v, got %v", expected, files)
	}
}

func TestCollector_IgnoreFile(t *testing.T) {
	tmpDir := t.TempDir()
	createTestTreeWithContent(t, tmpDir, map[string]string{
		".git/": "",
		".gitignore": strings.Join([]string{
			"*.gen.dml",
			"override_dir/",
		}, "\n") + "\n",
		".ignore": strings.Join([]string{
			"ignore_only.dml",
			"temp_ignore_dir/",
			"!keep.gen.dml",
			"!override_dir/",
		}, "\n") + "\n",
		"normal.dml":               "content",
		"ignore_only.dml":          "content",
		"skip.gen.dml":             "content",
		"keep.gen.dml":             "content",
		"temp_ignore_dir/a.dml":    "content",
		"override_dir/rescued.dml": "content",
	})

	cfg := &CLIConfig{
		Include: DefaultIncludePattern,
		Exclude: DefaultExcludePattern,
	}
	finder, err := newTestFinder(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, err := finder.Collect([]string{tmpDir})
	if err != nil {
		t.Fatalf("unexpected collect error: %v", err)
	}

	expected := []string{
		filepath.Join(tmpDir, "keep.gen.dml"),
		filepath.Join(tmpDir, "normal.dml"),
		filepath.Join(tmpDir, "override_dir", "rescued.dml"),
	}
	slices.Sort(expected)
	if !slices.Equal(files, expected) {
		t.Fatalf("expected %v, got %v", expected, files)
	}
}

func TestCollector_NestedGitIgnore_AndNegation(t *testing.T) {
	tmpDir := t.TempDir()
	createTestTreeWithContent(t, tmpDir, map[string]string{
		".git/":      "",
		".gitignore": "parent_skip_*.dml\n",
		"sub/.gitignore": strings.Join([]string{
			"/local_only.dml",
			"!parent_skip_kept.dml",
		}, "\n") + "\n",
		"sub/.ignore": strings.Join([]string{
			"sub_ignored.dml",
			"nested_glob/**",
			"!nested_glob/allow.dml",
		}, "\n") + "\n",
		"local_only.dml":            "content",
		"parent_skip_root.dml":      "content",
		"sub/ok.dml":                "content",
		"sub/local_only.dml":        "content",
		"sub/parent_skip_drop.dml":  "content",
		"sub/parent_skip_kept.dml":  "content",
		"sub/sub_ignored.dml":       "content",
		"sub/deep/local_only.dml":   "content",
		"sub/nested_glob/skip.dml":  "content",
		"sub/nested_glob/allow.dml": "content",
		"sub/nested_glob/dir/x.dml": "content",
	})

	cfg := &CLIConfig{
		Include: DefaultIncludePattern,
		Exclude: DefaultExcludePattern,
	}
	finder, err := newTestFinder(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, err := finder.Collect([]string{tmpDir})
	if err != nil {
		t.Fatalf("unexpected collect error: %v", err)
	}

	expected := []string{
		filepath.Join(tmpDir, "local_only.dml"),
		filepath.Join(tmpDir, "sub", "deep", "local_only.dml"),
		filepath.Join(tmpDir, "sub", "nested_glob", "allow.dml"),
		filepath.Join(tmpDir, "sub", "ok.dml"),
		filepath.Join(tmpDir, "sub", "parent_skip_kept.dml"),
	}
	slices.Sort(expected)
	if !slices.Equal(files, expected) {
		t.Fatalf("expected %v, got %v", expected, files)
	}
}

func TestCollector_GitIgnore_OverriddenByCustomExclude(t *testing.T) {
	tmpDir := t.TempDir()
	createTestTreeWithContent(t, tmpDir, map[string]string{
		".git/":           "",
		".gitignore":      "git_skip.dml\n",
		".ignore":         "ignore_skip.dml\n",
		"keep.dml":        "content",
		"git_skip.dml":    "content",
		"ignore_skip.dml": "content",
		"custom/skip.dml": "content",
		"venv/lib.dml":    "content",
	})

	t.Run("custom exclude value disables ignore files", func(t *testing.T) {
		cfg := &CLIConfig{
			Include: DefaultIncludePattern,
			Exclude: `/(\.git|custom)/`,
		}
		finder, err := newTestFinder(cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		files, err := finder.Collect([]string{tmpDir})
		if err != nil {
			t.Fatalf("unexpected collect error: %v", err)
		}

		expected := []string{
			filepath.Join(tmpDir, "git_skip.dml"),
			filepath.Join(tmpDir, "ignore_skip.dml"),
			filepath.Join(tmpDir, "keep.dml"),
			filepath.Join(tmpDir, "venv", "lib.dml"),
		}
		slices.Sort(expected)
		if !slices.Equal(files, expected) {
			t.Fatalf("expected %v, got %v", expected, files)
		}
	})

	t.Run(
		"explicit default exclude flag disables ignore files",
		func(t *testing.T) {
			cfg := &CLIConfig{
				Include:       DefaultIncludePattern,
				Exclude:       DefaultExcludePattern,
				ExplicitFlags: []string{"exclude"},
			}
			finder, err := newTestFinder(cfg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			files, err := finder.Collect([]string{tmpDir})
			if err != nil {
				t.Fatalf("unexpected collect error: %v", err)
			}

			expected := []string{
				filepath.Join(tmpDir, "custom", "skip.dml"),
				filepath.Join(tmpDir, "git_skip.dml"),
				filepath.Join(tmpDir, "ignore_skip.dml"),
				filepath.Join(tmpDir, "keep.dml"),
			}
			slices.Sort(expected)
			if !slices.Equal(files, expected) {
				t.Fatalf("expected %v, got %v", expected, files)
			}
		},
	)

	t.Run("empty exclude disables ignore files", func(t *testing.T) {
		cfg := &CLIConfig{
			Include: DefaultIncludePattern,
			Exclude: "",
		}
		finder, err := newTestFinder(cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		files, err := finder.Collect([]string{tmpDir})
		if err != nil {
			t.Fatalf("unexpected collect error: %v", err)
		}

		expected := []string{
			filepath.Join(tmpDir, "custom", "skip.dml"),
			filepath.Join(tmpDir, "git_skip.dml"),
			filepath.Join(tmpDir, "ignore_skip.dml"),
			filepath.Join(tmpDir, "keep.dml"),
			filepath.Join(tmpDir, "venv", "lib.dml"),
		}
		slices.Sort(expected)
		if !slices.Equal(files, expected) {
			t.Fatalf("expected %v, got %v", expected, files)
		}
	})
}

func TestCollector_GitIgnore_CombinedWithExtendAndForceExclude(t *testing.T) {
	tmpDir := t.TempDir()
	createTestTreeWithContent(t, tmpDir, map[string]string{
		".git/":              "",
		".gitignore":         "git_skip.dml\n",
		".ignore":            "ignore_skip.dml\n",
		"keep.dml":           "content",
		"git_skip.dml":       "content",
		"ignore_skip.dml":    "content",
		"forced_git.dml":     "content",
		"extra_dir/skip.dml": "content",
		"venv/lib.dml":       "content",
	})

	cfg := &CLIConfig{
		Include:       DefaultIncludePattern,
		Exclude:       DefaultExcludePattern,
		ExtendExclude: `/extra_dir/`,
		ForceExclude:  `forced_git\.dml$`,
	}
	finder, err := newTestFinder(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Recursive collection honors default excludes, --extend-exclude, --force-exclude,
	// .gitignore, and .ignore.
	files, err := finder.Collect([]string{tmpDir})
	if err != nil {
		t.Fatalf("unexpected collect error: %v", err)
	}
	expected := []string{
		filepath.Join(tmpDir, "keep.dml"),
	}
	if !slices.Equal(files, expected) {
		t.Fatalf("recursive collect expected %v, got %v", expected, files)
	}

	// Explicit file arguments bypass .gitignore/.ignore unless matched by --force-exclude.
	gitSkipFile := filepath.Join(tmpDir, "git_skip.dml")
	ignoreSkipFile := filepath.Join(tmpDir, "ignore_skip.dml")
	forcedFile := filepath.Join(tmpDir, "forced_git.dml")
	explicitFiles, err := finder.Collect([]string{
		gitSkipFile,
		ignoreSkipFile,
		forcedFile,
	})
	if err != nil {
		t.Fatalf("explicit collect error: %v", err)
	}
	expectedExplicit := []string{gitSkipFile, ignoreSkipFile}
	slices.Sort(expectedExplicit)
	if !slices.Equal(explicitFiles, expectedExplicit) {
		t.Fatalf(
			"explicit collect expected %v, got %v",
			expectedExplicit,
			explicitFiles,
		)
	}
}

func TestCollector_GitIgnore_VCSBoundary(t *testing.T) {
	tmpDir := t.TempDir()
	outerDir := filepath.Join(tmpDir, "outer")
	repoDir := filepath.Join(outerDir, "repo")

	createTestTreeWithContent(t, outerDir, map[string]string{
		".gitignore":              "outer_ignored.dml\n",
		".ignore":                 "outer_dir/\n",
		"repo/.git/":              "",
		"repo/.gitignore":         "repo_ignored.dml\n",
		"repo/main.dml":           "content",
		"repo/outer_ignored.dml":  "content",
		"repo/outer_dir/keep.dml": "content",
		"repo/repo_ignored.dml":   "content",
	})

	cfg := &CLIConfig{
		Include: DefaultIncludePattern,
		Exclude: DefaultExcludePattern,
	}
	finder, err := newTestFinder(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, err := finder.Collect([]string{repoDir})
	if err != nil {
		t.Fatalf("unexpected collect error: %v", err)
	}

	expected := []string{
		filepath.Join(repoDir, "main.dml"),
		filepath.Join(repoDir, "outer_dir", "keep.dml"),
		filepath.Join(repoDir, "outer_ignored.dml"),
	}
	slices.Sort(expected)
	if !slices.Equal(files, expected) {
		t.Fatalf("expected %v, got %v", expected, files)
	}
}

func TestCollector_GitIgnore_OnIgnoreCallback(t *testing.T) {
	tmpDir := t.TempDir()
	createTestTreeWithContent(t, tmpDir, map[string]string{
		".git/":               "",
		".gitignore":          "git_skip.dml\ngit_dir/\n",
		".ignore":             "dotignore_skip.dml\ndotignore_dir/\n",
		"keep.dml":            "content",
		"git_skip.dml":        "content",
		"git_dir/a.dml":       "content",
		"dotignore_skip.dml":  "content",
		"dotignore_dir/b.dml": "content",
	})

	cfg := &CLIConfig{
		Include: DefaultIncludePattern,
		Exclude: DefaultExcludePattern,
	}
	finder, err := newTestFinder(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	type ignoreEvent struct {
		path   string
		isDir  bool
		filter string
	}
	var events []ignoreEvent
	finder.SetOnIgnore(func(path string, isDir bool, filterName string) {
		events = append(events, ignoreEvent{
			path:   path,
			isDir:  isDir,
			filter: filterName,
		})
	})

	files, err := finder.Collect([]string{tmpDir})
	if err != nil {
		t.Fatalf("unexpected collect error: %v", err)
	}
	expectedFiles := []string{filepath.Join(tmpDir, "keep.dml")}
	if !slices.Equal(files, expectedFiles) {
		t.Fatalf("expected %v, got %v", expectedFiles, files)
	}

	wantEvents := []ignoreEvent{
		{
			path:   filepath.Join(tmpDir, "git_skip.dml"),
			isDir:  false,
			filter: lookup.FileGitIgnore,
		},
		{
			path:   filepath.Join(tmpDir, "git_dir"),
			isDir:  true,
			filter: lookup.FileGitIgnore,
		},
		{
			path:   filepath.Join(tmpDir, "dotignore_skip.dml"),
			isDir:  false,
			filter: lookup.FileIgnore,
		},
		{
			path:   filepath.Join(tmpDir, "dotignore_dir"),
			isDir:  true,
			filter: lookup.FileIgnore,
		},
	}
	for _, want := range wantEvents {
		if !slices.Contains(events, want) {
			t.Errorf("missing ignore event %+v in %+v", want, events)
		}
	}
}

func TestCollector_MaxFileSize_50MiB(t *testing.T) {
	tmpDir := t.TempDir()

	gitIgnorePath := filepath.Join(tmpDir, ".gitignore")
	f, err := os.Create(filepath.Clean(gitIgnorePath))
	if err != nil {
		t.Fatalf("failed to create .gitignore: %v", err)
	}
	if _, err := f.WriteString("skip.dml\n"); err != nil {
		_ = f.Close()
		t.Fatalf("failed to write rule to .gitignore: %v", err)
	}
	if err := f.Truncate(2 << 20); err != nil {
		_ = f.Close()
		t.Fatalf("failed to truncate .gitignore to 2 MiB: %v", err)
	}
	_ = f.Close()

	createTestTree(t, tmpDir, []string{
		"keep.dml",
		"skip.dml",
	})

	cfg := &CLIConfig{
		Include: DefaultIncludePattern,
		Exclude: DefaultExcludePattern,
	}
	finder, err := newTestFinder(cfg)
	if err != nil {
		t.Fatalf("failed to create test finder: %v", err)
	}

	files, err := finder.Collect([]string{tmpDir})
	if err != nil {
		t.Fatalf("unexpected collect error: %v", err)
	}

	expectedFiles := []string{filepath.Join(tmpDir, "keep.dml")}
	if !slices.Equal(files, expectedFiles) {
		t.Fatalf("expected %v, got %v", expectedFiles, files)
	}
}
