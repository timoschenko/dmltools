// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test_lookup

import (
	"errors"
	iofs "io/fs"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/timoschenko/dmltools/pkg/fs/lookup"
)

func TestVCSBoundaries(t *testing.T) {
	now := time.Unix(1700000000, 0)

	t.Run("git worktree regular file", func(t *testing.T) {
		m := newMockFS("/outer/worktree/sub", 1000)
		m.setDir("/", 1, 0, now)
		m.setDir("/outer", 1, 1000, now)
		m.setDir("/outer/.git", 1, 1000, now)
		m.setFile("/outer/.editorconfig", 0644, 10, 1, 1000, now)
		m.setDir("/outer/worktree", 1, 1000, now)
		m.setFile("/outer/worktree/.git", 0644, 32, 1, 1000, now)
		m.setFile("/outer/worktree/.editorconfig", 0644, 12, 1, 1000, now)
		m.setDir("/outer/worktree/sub", 1, 1000, now)

		finder, err := lookup.NewFinder(lookup.WithFileSystem(m))
		if err != nil {
			t.Fatalf("NewFinder failed: %v", err)
		}
		b, err := finder.Cache().ResolveBoundary("/outer/worktree/sub")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if b.Reason != lookup.BoundaryVCSRoot ||
			b.RootDir != filepath.Clean("/outer/worktree") {
			t.Fatalf("unexpected boundary: %+v", b)
		}

		chain, err := finder.FindChain(
			"/outer/worktree/sub",
			lookup.FileEditorConfig,
		)
		if err != nil {
			t.Fatalf("FindChain failed: %v", err)
		}
		if len(chain) != 1 ||
			chain[0].Path != filepath.Clean("/outer/worktree/.editorconfig") {
			t.Fatalf("unexpected chain: %+v", chain)
		}
	})

	t.Run("hg directory and custom markers", func(t *testing.T) {
		m := newMockFS("/proj/src", 1000)
		m.setDir("/", 1, 0, now)
		m.setDir("/proj", 1, 1000, now)
		m.setDir("/proj/.hg", 1, 1000, now)
		m.setDir("/proj/src", 1, 1000, now)
		m.setDir("/proj/src/.jj", 1, 1000, now)

		c1 := lookup.NewCache(lookup.WithFileSystem(m))
		b, err := c1.ResolveBoundary("/proj/src")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if b.Reason != lookup.BoundaryVCSRoot ||
			b.RootDir != filepath.Clean("/proj") {
			t.Fatalf("expected /proj vcs boundary, got %+v", b)
		}

		c2 := lookup.NewCache(
			lookup.WithFileSystem(m),
			lookup.WithVCSMarkers(".jj"),
		)
		cb, err := c2.ResolveBoundary("/proj/src")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if cb.Reason != lookup.BoundaryVCSRoot ||
			cb.RootDir != filepath.Clean("/proj/src") {
			t.Fatalf("expected /proj/src .jj boundary, got %+v", cb)
		}
	})

	t.Run("disable VCS root stop", func(t *testing.T) {
		m := buildMockRepo(now)
		c := lookup.NewCache(
			lookup.WithFileSystem(m),
			lookup.WithStopAtVCSRoot(false),
		)
		b, err := c.ResolveBoundary("/repo/sub")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if b.Reason != lookup.BoundaryFilesystemRoot {
			t.Fatalf("expected BoundaryFilesystemRoot, got %+v", b)
		}
	})
}

func TestMountOwnershipCeilingAndDepthBoundaries(t *testing.T) {
	now := time.Unix(1700000000, 0)

	t.Run("mount point boundary", func(t *testing.T) {
		m := newMockFS("/mnt/usb/proj/sub", 1000)
		m.setDir("/", 1, 1000, now)
		m.setDir("/mnt", 1, 1000, now)
		m.setDir("/mnt/usb", 42, 1000, now)
		m.setDir("/mnt/usb/proj", 42, 1000, now)
		m.setDir("/mnt/usb/proj/sub", 42, 1000, now)

		c := lookup.NewCache(lookup.WithFileSystem(m))
		b, err := c.ResolveBoundary("/mnt/usb/proj/sub")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if b.Reason != lookup.BoundaryMountPoint ||
			b.RootDir != filepath.Clean("/mnt/usb") {
			t.Fatalf("expected BoundaryMountPoint at /mnt/usb, got %+v", b)
		}

		noMount := lookup.NewCache(
			lookup.WithFileSystem(m),
			lookup.WithStopAtMountPoint(false),
		)
		b2, err := noMount.ResolveBoundary("/mnt/usb/proj/sub")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if b2.Reason != lookup.BoundaryFilesystemRoot {
			t.Fatalf("expected BoundaryFilesystemRoot, got %+v", b2)
		}
	})

	t.Run("ownership boundary", func(t *testing.T) {
		m := newMockFS("/shared/otheruser/myproj", 1000)
		m.setDir("/", 1, 0, now)
		m.setDir("/shared", 1, 0, now)
		m.setDir("/shared/otheruser", 1, 2000, now)
		m.setDir("/shared/otheruser/myproj", 1, 1000, now)

		c := lookup.NewCache(lookup.WithFileSystem(m))
		b, err := c.ResolveBoundary("/shared/otheruser/myproj")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if b.Reason != lookup.BoundaryOwnership ||
			b.RootDir != filepath.Clean("/shared/otheruser/myproj") {
			t.Fatalf("expected BoundaryOwnership, got %+v", b)
		}

		noOwn := lookup.NewCache(
			lookup.WithFileSystem(m),
			lookup.WithVerifyOwnership(false),
		)
		b2, err := noOwn.ResolveBoundary("/shared/otheruser/myproj")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if b2.Reason != lookup.BoundaryFilesystemRoot {
			t.Fatalf("expected BoundaryFilesystemRoot, got %+v", b2)
		}
	})

	t.Run("ceiling dir and max depth", func(t *testing.T) {
		m := buildMockRepo(now)
		ceilCache := lookup.NewCache(
			lookup.WithFileSystem(m),
			lookup.WithCeilingDirs("/repo/sub"),
		)
		cb, err := ceilCache.ResolveBoundary("/repo/sub/deep")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if cb.Reason != lookup.BoundaryCeiling ||
			cb.RootDir != filepath.Clean("/repo/sub") {
			t.Fatalf("expected BoundaryCeiling at /repo/sub, got %+v", cb)
		}

		depthCache := lookup.NewCache(
			lookup.WithFileSystem(m),
			lookup.WithMaxDepth(0),
		)
		db, err := depthCache.ResolveBoundary("/repo/sub/deep")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if db.Reason != lookup.BoundaryMaxDepth ||
			db.RootDir != filepath.Clean("/repo/sub/deep") {
			t.Fatalf("expected BoundaryMaxDepth at deep, got %+v", db)
		}
	})
}

func TestSecurityGuardsAndStrictMode(t *testing.T) {
	now := time.Unix(1700000000, 0)

	cases := []struct {
		name       string
		setup      func(m *mockFS) string
		opts       []lookup.Option
		wantReason lookup.SkipReason
	}{
		{
			name: "directory instead of regular file",
			setup: func(m *mockFS) string {
				m.setDir("/repo/.editorconfig", 1, 1000, now)
				return "/repo/.editorconfig"
			},
			wantReason: lookup.SkipNotRegular,
		},
		{
			name: "fifo named pipe",
			setup: func(m *mockFS) string {
				p := "/repo/.editorconfig"
				m.setFile(p, iofs.ModeNamedPipe|0644, 0, 1, 1000, now)
				return p
			},
			wantReason: lookup.SkipNotRegular,
		},
		{
			name: "unix socket",
			setup: func(m *mockFS) string {
				p := "/repo/.editorconfig"
				m.setFile(p, iofs.ModeSocket|0644, 0, 1, 1000, now)
				return p
			},
			wantReason: lookup.SkipNotRegular,
		},
		{
			name: "device node",
			setup: func(m *mockFS) string {
				p := "/repo/.editorconfig"
				m.setFile(p, iofs.ModeDevice|0644, 0, 1, 1000, now)
				return p
			},
			wantReason: lookup.SkipNotRegular,
		},
		{
			name: "oversized file",
			setup: func(m *mockFS) string {
				p := "/repo/.editorconfig"
				m.setFile(p, 0644, 2048, 1, 1000, now)
				return p
			},
			opts:       []lookup.Option{lookup.WithMaxFileSize(1024)},
			wantReason: lookup.SkipFileTooLarge,
		},
		{
			name: "symlinks disabled",
			setup: func(m *mockFS) string {
				m.setFile("/repo/shared.cfg", 0644, 16, 1, 1000, now)
				p := "/repo/.editorconfig"
				m.setSymlink(p, "/repo/shared.cfg", 1, 1000, now)
				return p
			},
			opts:       []lookup.Option{lookup.WithAllowSymlinks(false)},
			wantReason: lookup.SkipSymlinkDisabled,
		},
		{
			name: "broken symlink",
			setup: func(m *mockFS) string {
				p := "/repo/.editorconfig"
				m.setSymlink(p, "/repo/missing.cfg", 1, 1000, now)
				return p
			},
			wantReason: lookup.SkipBrokenSymlink,
		},
		{
			name: "symlink escaping boundary",
			setup: func(m *mockFS) string {
				m.setFile("/outside/secret.cfg", 0644, 16, 1, 1000, now)
				p := "/repo/.editorconfig"
				m.setSymlink(p, "/outside/secret.cfg", 1, 1000, now)
				return p
			},
			wantReason: lookup.SkipSymlinkEscape,
		},
		{
			name: "symlink crossing mount device",
			setup: func(m *mockFS) string {
				m.setFile("/repo/other.cfg", 0644, 16, 99, 1000, now)
				p := "/repo/.editorconfig"
				m.setSymlink(p, "/repo/other.cfg", 1, 1000, now)
				return p
			},
			wantReason: lookup.SkipSymlinkCrossDevice,
		},
		{
			name: "permission denied",
			setup: func(m *mockFS) string {
				p := "/repo/.editorconfig"
				m.setStatErr(p, iofs.ErrPermission)
				return p
			},
			wantReason: lookup.SkipPermissionDenied,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := buildMockRepo(now)
			target := tc.setup(m)

			var skips []lookup.SkipReason
			opts := append([]lookup.Option{}, tc.opts...)
			opts = append(
				opts,
				lookup.WithFileSystem(m),
				lookup.WithOnSkip(func(p string, r lookup.SkipReason, _ error) {
					if p == filepath.Clean(target) {
						skips = append(skips, r)
					}
				}),
			)

			finder, err := lookup.NewFinder(opts...)
			if err != nil {
				t.Fatalf("NewFinder failed: %v", err)
			}
			got, err := finder.FindChain("/repo/sub", lookup.FileEditorConfig)
			if err != nil {
				t.Fatalf("non-strict mode returned error: %v", err)
			}
			if len(got) != 0 {
				t.Fatalf("expected 0 files, got %+v", got)
			}
			if len(skips) != 1 || skips[0] != tc.wantReason {
				t.Fatalf("skips = %v, want [%v]", skips, tc.wantReason)
			}

			strictOpts := append([]lookup.Option{}, tc.opts...)
			strictOpts = append(
				strictOpts,
				lookup.WithFileSystem(m),
				lookup.WithStrictSecurity(true),
			)
			strictFinder, err := lookup.NewFinder(strictOpts...)
			if err != nil {
				t.Fatalf("NewFinder strict failed: %v", err)
			}
			_, err = strictFinder.FindChain(
				"/repo/sub",
				lookup.FileEditorConfig,
			)
			var secErr *lookup.SecurityError
			if !errors.As(err, &secErr) {
				t.Fatalf("expected *SecurityError, got %v", err)
			}
			if secErr.Reason != tc.wantReason {
				t.Errorf(
					"secErr.Reason = %v, want %v",
					secErr.Reason,
					tc.wantReason,
				)
			}
			if secErr.Error() == "" {
				t.Error("expected non-empty Error() string")
			}
		})
	}
}

func TestValidSymlinkAndEscapeOverride(t *testing.T) {
	now := time.Unix(1700000000, 0)
	m := buildMockRepo(now)
	m.setFile("/repo/shared/editorconfig", 0644, 42, 1, 1000, now)
	m.setDir("/repo/shared", 1, 1000, now)
	m.setSymlink(
		"/repo/sub/.editorconfig",
		"/repo/shared/editorconfig",
		1,
		1000,
		now,
	)
	m.setFile("/outside/global.cfg", 0644, 12, 1, 1000, now)
	m.setSymlink("/repo/.editorconfig", "/outside/global.cfg", 1, 1000, now)

	finder, err := lookup.NewFinder(lookup.WithFileSystem(m))
	if err != nil {
		t.Fatalf("NewFinder failed: %v", err)
	}
	files, err := finder.FindChain("/repo/sub", lookup.FileEditorConfig)
	if err != nil {
		t.Fatalf("FindChain failed: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 valid in-repo symlink, got %+v", files)
	}
	if !files[0].IsSymlink ||
		files[0].ResolvedPath != filepath.Clean("/repo/shared/editorconfig") {
		t.Errorf("unexpected symlink metadata: %+v", files[0])
	}

	allowEsc, err := lookup.NewFinder(
		lookup.WithFileSystem(m),
		lookup.WithAllowSymlinkEscape(true),
	)
	if err != nil {
		t.Fatalf("NewFinder failed: %v", err)
	}
	escFiles, err := allowEsc.FindChain("/repo/sub", lookup.FileEditorConfig)
	if err != nil {
		t.Fatalf("FindChain failed: %v", err)
	}
	if len(escFiles) != 2 {
		t.Fatalf("expected 2 files with escape allowed, got %d", len(escFiles))
	}
}

func TestCachingMutationsTTLAndConcurrency(t *testing.T) {
	t0 := time.Unix(1700000000, 0)
	clk := newMockClock(t0)

	t.Run("real-time stat validation detects FS mutations", func(t *testing.T) {
		m := buildMockRepo(t0)
		finder, err := lookup.NewFinder(
			lookup.WithFileSystem(m),
			lookup.WithClock(clk),
		)
		if err != nil {
			t.Fatalf("NewFinder failed: %v", err)
		}

		files, err := finder.FindChain("/repo/sub", lookup.FileEditorConfig)
		if err != nil || len(files) != 0 {
			t.Fatalf("expected 0 files, got %v, %v", files, err)
		}
		st1 := finder.Cache().Stats()
		if st1.DirMisses == 0 || st1.FileMisses == 0 || st1.DirHits != 0 {
			t.Fatalf(
				"expected initial misses and 0 dir hits, got %+v",
				st1,
			)
		}

		_, err = finder.FindChain("/repo/sub", lookup.FileEditorConfig)
		if err != nil {
			t.Fatalf("FindChain failed: %v", err)
		}
		st2 := finder.Cache().Stats()
		if st2.DirHits <= st1.DirHits || st2.FileHits <= st1.FileHits {
			t.Fatalf(
				"expected dir/file hits, got before=%+v after=%+v",
				st1,
				st2,
			)
		}

		t1 := t0.Add(time.Second)
		m.setFile("/repo/.editorconfig", 0644, 12, 1, 1000, t1)
		files, err = finder.FindChain("/repo/sub", lookup.FileEditorConfig)
		if err != nil || len(files) != 1 || files[0].Size != 12 {
			t.Fatalf(
				"expected newly created file size 12, got %+v, %v",
				files,
				err,
			)
		}

		m.setFile(
			"/repo/.editorconfig",
			0644,
			99,
			1,
			1000,
			t1.Add(time.Second),
		)
		files, err = finder.FindChain("/repo/sub", lookup.FileEditorConfig)
		if err != nil || len(files) != 1 || files[0].Size != 99 {
			t.Fatalf("expected updated file size 99, got %+v, %v", files, err)
		}

		m.setFile("/repo/target.cfg", 0644, 7, 1, 1000, t1)
		m.setSymlink("/repo/.editorconfig", "/repo/target.cfg", 1, 1000, t1)
		files, err = finder.FindChain("/repo/sub", lookup.FileEditorConfig)
		if err != nil || len(files) != 1 || !files[0].IsSymlink ||
			files[0].Size != 7 {
			t.Fatalf("expected symlink size 7, got %+v, %v", files, err)
		}
		m.setFile(
			"/repo/target.cfg",
			0644,
			33,
			1,
			1000,
			t1.Add(2*time.Second),
		)
		files, err = finder.FindChain("/repo/sub", lookup.FileEditorConfig)
		if err != nil || len(files) != 1 || files[0].Size != 33 {
			t.Fatalf(
				"expected updated symlink target size 33, got %+v, %v",
				files,
				err,
			)
		}

		m.deleteNode("/repo/.editorconfig")
		files, err = finder.FindChain("/repo/sub", lookup.FileEditorConfig)
		if err != nil || len(files) != 0 {
			t.Fatalf("expected deleted file gone, got %+v, %v", files, err)
		}

		m.setDir("/repo/sub/.git", 1, 1000, t1)
		b, err := finder.Cache().ResolveBoundary("/repo/sub")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if b.RootDir != filepath.Clean("/repo/sub") {
			t.Fatalf("expected new .git boundary at /repo/sub, got %+v", b)
		}
	})

	t.Run(
		"TTL cache avoids syscalls and explicit Invalidate/Reset",
		func(t *testing.T) {
			m := buildMockRepo(t0)
			m.setFile("/repo/.editorconfig", 0644, 10, 1, 1000, t0)
			finder, err := lookup.NewFinder(
				lookup.WithFileSystem(m),
				lookup.WithClock(clk),
				lookup.WithCacheTTL(5*time.Second),
			)
			if err != nil {
				t.Fatalf("NewFinder failed: %v", err)
			}

			if _, err := finder.FindChain(
				"/repo/sub",
				lookup.FileEditorConfig,
			); err != nil {
				t.Fatalf("FindChain failed: %v", err)
			}
			m.resetStatCalls()

			clk.Advance(2 * time.Second)
			if _, err := finder.FindChain(
				"/repo/sub",
				lookup.FileEditorConfig,
			); err != nil {
				t.Fatalf("FindChain failed: %v", err)
			}
			if calls := m.totalStatCalls(); calls != 1 {
				t.Errorf(
					"expected 1 start stat call within TTL, got %d",
					calls,
				)
			}

			finder.Cache().Invalidate("/repo/.editorconfig")
			if finder.Cache().Stats().Invalidations == 0 {
				t.Error("expected Invalidations > 0 after Invalidate")
			}
			if _, err := finder.FindChain(
				"/repo/sub",
				lookup.FileEditorConfig,
			); err != nil {
				t.Fatalf("FindChain failed: %v", err)
			}

			clk.Advance(10 * time.Second)
			finder.Cache().Reset()
			if _, err := finder.FindChain(
				"/repo/sub",
				lookup.FileEditorConfig,
			); err != nil {
				t.Fatalf("FindChain after Reset failed: %v", err)
			}
		},
	)

	t.Run("concurrent access", func(t *testing.T) {
		m := buildMockRepo(t0)
		m.setFile("/repo/.editorconfig", 0644, 10, 1, 1000, t0)
		m.setFile("/repo/sub/.gitignore", 0644, 15, 1, 1000, t0)
		m.setFile("/repo/sub/deep/a.dml", 0644, 10, 1, 1000, t0)

		finder, err := lookup.NewFinder(
			lookup.WithFileSystem(m),
			lookup.WithClock(clk),
			lookup.WithInclude(`\.dml$`),
		)
		if err != nil {
			t.Fatalf("NewFinder failed: %v", err)
		}

		var wg sync.WaitGroup
		for i := range 16 {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := range 25 {
					_, _ = finder.FindChain(
						"/repo/sub/deep",
						lookup.FileEditorConfig,
						lookup.FileGitIgnore,
					)
					_, _, _ = finder.FindNearest(
						"/repo/sub/deep",
						lookup.FileEditorConfig,
					)
					_, _ = finder.Collect([]string{"/repo"})
					_, _ = finder.Cache().ResolveBoundary("/repo/sub/deep")
					if (id+j)%7 == 0 {
						finder.Cache().Invalidate("/repo/sub")
					}
					if (id+j)%13 == 0 {
						finder.Cache().Reset()
					}
					_ = finder.Cache().Stats()
				}
			}(i)
		}
		wg.Wait()
	})
}

func TestEnumsAndErrorFormatting(t *testing.T) {
	reasons := []lookup.BoundaryReason{
		lookup.BoundaryNone,
		lookup.BoundaryFilesystemRoot,
		lookup.BoundaryVCSRoot,
		lookup.BoundaryMountPoint,
		lookup.BoundaryOwnership,
		lookup.BoundaryCeiling,
		lookup.BoundaryMaxDepth,
		lookup.BoundaryReason(99),
	}
	for _, r := range reasons {
		if r.String() == "" {
			t.Errorf("empty String for BoundaryReason %d", r)
		}
	}

	skips := []lookup.SkipReason{
		lookup.SkipNotRegular,
		lookup.SkipFileTooLarge,
		lookup.SkipSymlinkDisabled,
		lookup.SkipBrokenSymlink,
		lookup.SkipSymlinkEscape,
		lookup.SkipSymlinkCrossDevice,
		lookup.SkipPermissionDenied,
		lookup.SkipExcluded,
		lookup.SkipReason(99),
	}
	for _, s := range skips {
		if s.String() == "" {
			t.Errorf("empty String for SkipReason %d", s)
		}
	}

	inner := errors.New("boom")
	se := &lookup.SecurityError{
		Path:   "/a",
		Reason: lookup.SkipNotRegular,
		Err:    inner,
	}
	if !errors.Is(se, inner) {
		t.Error("SecurityError.Unwrap did not match inner error")
	}
	if !strings.Contains(se.Error(), "boom") {
		t.Errorf("unexpected Error(): %q", se.Error())
	}
	seNoErr := &lookup.SecurityError{
		Path:   "/a",
		Reason: lookup.SkipNotRegular,
	}
	if seNoErr.Error() == "" {
		t.Error("expected non-empty Error() without inner err")
	}
}

func TestCacheIsWithinBoundary(t *testing.T) {
	t0 := time.Unix(1700000000, 0)
	m := buildMockRepo(t0)
	m.setFile("/repo/sub/deep/main.dml", 0644, 10, 1, 1000, t0)
	m.setFile("/outside/ext.dml", 0644, 10, 1, 1000, t0)
	m.setSymlink("/repo/sub/escape.dml", "/outside/ext.dml", 1, 1000, t0)
	m.setSymlink(
		"/repo/sub/internal.dml",
		"/repo/sub/deep/main.dml",
		1,
		1000,
		t0,
	)

	c := lookup.NewCache(lookup.WithFileSystem(m))
	if !c.IsWithinBoundary("/repo", "/repo/sub/deep/main.dml") {
		t.Error("expected main.dml to be within /repo boundary")
	}
	if !c.IsWithinBoundary("/repo", "main.dml") {
		t.Error(
			"expected relative main.dml (cwd=/repo/sub/deep) to be within /repo",
		)
	}
	if c.IsWithinBoundary("/repo", "/outside/ext.dml") {
		t.Error("expected /outside/ext.dml to be outside /repo boundary")
	}
	if c.IsWithinBoundary("/repo", "/repo/sub/escape.dml") {
		t.Error("expected escaping symlink to be outside /repo boundary")
	}
	if !c.IsWithinBoundary("/repo", "/repo/sub/internal.dml") {
		t.Error("expected internal symlink to be within /repo boundary")
	}
	if c.IsWithinBoundary("", "/repo/sub/deep/main.dml") ||
		c.IsWithinBoundary("/repo", "") {
		t.Error("expected empty boundary or target to return false")
	}

	cEsc := lookup.NewCache(
		lookup.WithFileSystem(m),
		lookup.WithAllowSymlinkEscape(true),
	)
	if !cEsc.IsWithinBoundary("/repo", "/repo/sub/escape.dml") {
		t.Error("expected escaping symlink allowed when allowSymlinkEsc=true")
	}
}
