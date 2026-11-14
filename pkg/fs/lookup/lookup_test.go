// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lookup

import (
	"errors"
	iofs "io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

type mockClock struct {
	mu  sync.Mutex
	now time.Time
}

func newMockClock(t time.Time) *mockClock {
	return &mockClock{now: t}
}

func (m *mockClock) Now() time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.now
}

func (m *mockClock) Advance(d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.now = m.now.Add(d)
}

type mockNode struct {
	meta         FileMeta
	symlinkTo    string
	content      []byte
	statErr      error
	lstatErr     error
	readDirErr   error
	readFileErr  error
	evalOverride string
}

type mockDirEntry struct {
	name string
	mode iofs.FileMode
	meta FileMeta
}

func (e mockDirEntry) Name() string        { return e.name }
func (e mockDirEntry) IsDir() bool         { return e.mode.IsDir() }
func (e mockDirEntry) Type() iofs.FileMode { return e.mode.Type() }

func (e mockDirEntry) Info() (iofs.FileInfo, error) { return mockFileInfo{meta: e.meta}, nil }

type mockFileInfo struct {
	meta FileMeta
}

func (i mockFileInfo) Name() string        { return i.meta.Name }
func (i mockFileInfo) Size() int64         { return i.meta.Size }
func (i mockFileInfo) Mode() iofs.FileMode { return i.meta.Mode }
func (i mockFileInfo) ModTime() time.Time  { return i.meta.ModTime }
func (i mockFileInfo) IsDir() bool         { return i.meta.Mode.IsDir() }
func (i mockFileInfo) Sys() any            { return nil }

type mockFS struct {
	mu            sync.Mutex
	cwd           string
	euid          int
	nodes         map[string]*mockNode
	statCalls     map[string]int
	readFileCalls map[string]int
}

func newMockFS(cwd string, euid int) *mockFS {
	return &mockFS{
		cwd:           filepath.Clean(cwd),
		euid:          euid,
		nodes:         make(map[string]*mockNode),
		statCalls:     make(map[string]int),
		readFileCalls: make(map[string]int),
	}
}

func (m *mockFS) setDir(
	path string,
	dev uint64,
	uid uint32,
	mtime time.Time,
) {
	m.mu.Lock()
	defer m.mu.Unlock()
	clean := filepath.Clean(path)
	m.nodes[clean] = &mockNode{
		meta: FileMeta{
			Name:      filepath.Base(clean),
			Mode:      iofs.ModeDir | 0755,
			ModTime:   mtime,
			DeviceID:  dev,
			Inode:     uint64(len(m.nodes) + 1),
			UID:       uid,
			HasDevice: true,
			HasInode:  true,
			HasUID:    true,
		},
	}
}

func (m *mockFS) setFile(
	path string,
	mode iofs.FileMode,
	size int64,
	dev uint64,
	uid uint32,
	mtime time.Time,
) {
	m.mu.Lock()
	defer m.mu.Unlock()
	clean := filepath.Clean(path)
	m.nodes[clean] = &mockNode{
		meta: FileMeta{
			Name:      filepath.Base(clean),
			Mode:      mode,
			Size:      size,
			ModTime:   mtime,
			DeviceID:  dev,
			Inode:     uint64(len(m.nodes) + 1),
			UID:       uid,
			HasDevice: true,
			HasInode:  true,
			HasUID:    true,
		},
	}
}

func (m *mockFS) setFileWithContent(
	path string,
	content []byte,
	dev uint64,
	uid uint32,
	mtime time.Time,
) {
	m.mu.Lock()
	defer m.mu.Unlock()
	clean := filepath.Clean(path)
	m.nodes[clean] = &mockNode{
		meta: FileMeta{
			Name:      filepath.Base(clean),
			Mode:      0644,
			Size:      int64(len(content)),
			ModTime:   mtime,
			DeviceID:  dev,
			Inode:     uint64(len(m.nodes) + 1),
			UID:       uid,
			HasDevice: true,
			HasInode:  true,
			HasUID:    true,
		},
		content: slices.Clone(content),
	}
}

func (m *mockFS) setSymlink(
	path string,
	target string,
	dev uint64,
	uid uint32,
	mtime time.Time,
) {
	m.mu.Lock()
	defer m.mu.Unlock()
	clean := filepath.Clean(path)
	m.nodes[clean] = &mockNode{
		meta: FileMeta{
			Name:      filepath.Base(clean),
			Mode:      iofs.ModeSymlink | 0777,
			ModTime:   mtime,
			DeviceID:  dev,
			Inode:     uint64(len(m.nodes) + 1),
			UID:       uid,
			HasDevice: true,
			HasInode:  true,
			HasUID:    true,
		},
		symlinkTo: filepath.Clean(target),
	}
}

func (m *mockFS) deleteNode(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.nodes, filepath.Clean(path))
}

func (m *mockFS) setStatErr(path string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	clean := filepath.Clean(path)
	if n, ok := m.nodes[clean]; ok {
		n.statErr = err
		n.lstatErr = err
	} else {
		m.nodes[clean] = &mockNode{
			statErr:  err,
			lstatErr: err,
		}
	}
}

func (m *mockFS) totalStatCalls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	total := 0
	for _, c := range m.statCalls {
		total += c
	}
	return total
}

func (m *mockFS) resetStatCalls() {
	m.mu.Lock()
	defer m.mu.Unlock()
	clear(m.statCalls)
}

func (m *mockFS) totalReadFileCalls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	total := 0
	for _, c := range m.readFileCalls {
		total += c
	}
	return total
}

func (m *mockFS) absLocked(path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Clean(filepath.Join(m.cwd, path))
}

func (m *mockFS) Stat(path string) (FileMeta, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	clean := m.absLocked(path)
	m.statCalls[clean]++
	for range 8 {
		n, ok := m.nodes[clean]
		if !ok {
			return FileMeta{}, os.ErrNotExist
		}
		if n.statErr != nil {
			return FileMeta{}, n.statErr
		}
		if n.symlinkTo == "" {
			return n.meta, nil
		}
		clean = n.symlinkTo
	}
	return FileMeta{}, errors.New("too many levels of symbolic links")
}

func (m *mockFS) Lstat(path string) (FileMeta, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	clean := m.absLocked(path)
	m.statCalls[clean]++
	n, ok := m.nodes[clean]
	if !ok {
		return FileMeta{}, os.ErrNotExist
	}
	if n.lstatErr != nil {
		return FileMeta{}, n.lstatErr
	}
	return n.meta, nil
}

func (m *mockFS) ReadDir(path string) ([]iofs.DirEntry, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	clean := m.absLocked(path)
	n, ok := m.nodes[clean]
	if !ok {
		return nil, os.ErrNotExist
	}
	if n.readDirErr != nil {
		return nil, n.readDirErr
	}
	prefix := clean + string(filepath.Separator)
	var entries []iofs.DirEntry
	for p, node := range m.nodes {
		if !strings.HasPrefix(p, prefix) {
			continue
		}
		rest := strings.TrimPrefix(p, prefix)
		if rest == "" || strings.Contains(rest, string(filepath.Separator)) {
			continue
		}
		entries = append(entries, mockDirEntry{
			name: rest,
			mode: node.meta.Mode,
			meta: node.meta,
		})
	}
	return entries, nil
}

func (m *mockFS) ReadFile(path string) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	clean := m.absLocked(path)
	m.readFileCalls[clean]++
	for range 8 {
		n, ok := m.nodes[clean]
		if !ok {
			return nil, os.ErrNotExist
		}
		if n.readFileErr != nil {
			return nil, n.readFileErr
		}
		if n.symlinkTo == "" {
			if n.meta.Mode.IsDir() {
				return nil, errors.New("is a directory")
			}
			return slices.Clone(n.content), nil
		}
		clean = n.symlinkTo
	}
	return nil, errors.New("too many levels of symbolic links")
}

func (m *mockFS) EvalSymlinks(path string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	clean := m.absLocked(path)
	for range 8 {
		n, ok := m.nodes[clean]
		if !ok {
			return "", os.ErrNotExist
		}
		if n.evalOverride != "" {
			return n.evalOverride, nil
		}
		if n.symlinkTo == "" {
			return clean, nil
		}
		clean = n.symlinkTo
	}
	return "", errors.New("too many levels of symbolic links")
}

func (m *mockFS) Abs(path string) (string, error) {
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	return filepath.Clean(filepath.Join(m.cwd, path)), nil
}

func (m *mockFS) Geteuid() int {
	return m.euid
}

func buildMockRepo(baseTime time.Time) *mockFS {
	m := newMockFS("/repo/sub/deep", 1000)
	m.setDir("/", 1, 0, baseTime)
	m.setDir("/outside", 1, 1000, baseTime)
	m.setDir("/repo", 1, 1000, baseTime)
	m.setDir("/repo/.git", 1, 1000, baseTime)
	m.setDir("/repo/sub", 1, 1000, baseTime)
	m.setDir("/repo/sub/deep", 1, 1000, baseTime)
	return m
}

func TestSOLIDCompositionAndInjection(t *testing.T) {
	now := time.Unix(1700000000, 0)
	m := buildMockRepo(now)
	m.setFile("/repo/.editorconfig", 0644, 20, 1, 1000, now)
	m.setFile("/repo/sub/main.dml", 0644, 30, 1, 1000, now)

	cache := NewCache(WithFileSystem(m))
	filter, err := NewFilter(
		WithFileSystem(m),
		WithInclude(`\.dml$`),
		WithExclude(`/\.git/`),
	)
	if err != nil {
		t.Fatalf("NewFilter failed: %v", err)
	}
	if !filter.HasInclude() {
		t.Fatal("expected filter.HasInclude() to be true")
	}

	finder1, err := NewFinder(WithCache(cache), WithFilter(filter))
	if err != nil {
		t.Fatalf("NewFinder 1 failed: %v", err)
	}
	finder2, err := NewFinder(WithCache(cache), WithFilter(filter))
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
	chain, err := finder2.FindChain("/repo/sub/main.dml", FileEditorConfig)
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

	finder, err := NewFinder(WithFileSystem(m))
	if err != nil {
		t.Fatalf("NewFinder failed: %v", err)
	}

	chain, err := finder.FindChain(
		"/repo/sub/deep/main.dml",
		FileEditorConfig,
		FileGitIgnore,
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
		FileEditorConfig,
		FileGitIgnore,
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

	finder, err := NewFinder(
		WithFileSystem(m),
		WithInclude(`\.dml$`),
		WithExclude(`/(\.git|venv)/`),
		WithForceExclude(`skip_gen`),
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
		filter: FilterInclude,
	}) {
		t.Errorf(
			"expected b.txt to trigger %s ignore event, got %+v",
			FilterInclude,
			ignored,
		)
	}
	if !slices.Contains(ignored, ignoredEvent{
		path:   filepath.Clean("/ws/venv"),
		isDir:  true,
		filter: FilterExclude,
	}) {
		t.Errorf(
			"expected venv dir to trigger %s ignore event, got %+v",
			FilterExclude,
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

		finder, err := NewFinder(WithFileSystem(m))
		if err != nil {
			t.Fatalf("NewFinder failed: %v", err)
		}
		b, err := finder.Cache().ResolveBoundary("/outer/worktree/sub")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if b.Reason != BoundaryVCSRoot ||
			b.RootDir != filepath.Clean("/outer/worktree") {
			t.Fatalf("unexpected boundary: %+v", b)
		}

		chain, err := finder.FindChain("/outer/worktree/sub", FileEditorConfig)
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

		c1 := NewCache(WithFileSystem(m))
		b, err := c1.ResolveBoundary("/proj/src")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if b.Reason != BoundaryVCSRoot || b.RootDir != filepath.Clean("/proj") {
			t.Fatalf("expected /proj vcs boundary, got %+v", b)
		}

		c2 := NewCache(WithFileSystem(m), WithVCSMarkers(".jj"))
		cb, err := c2.ResolveBoundary("/proj/src")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if cb.Reason != BoundaryVCSRoot ||
			cb.RootDir != filepath.Clean("/proj/src") {
			t.Fatalf("expected /proj/src .jj boundary, got %+v", cb)
		}
	})

	t.Run("disable VCS root stop", func(t *testing.T) {
		m := buildMockRepo(now)
		c := NewCache(WithFileSystem(m), WithStopAtVCSRoot(false))
		b, err := c.ResolveBoundary("/repo/sub")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if b.Reason != BoundaryFilesystemRoot {
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

		c := NewCache(WithFileSystem(m))
		b, err := c.ResolveBoundary("/mnt/usb/proj/sub")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if b.Reason != BoundaryMountPoint ||
			b.RootDir != filepath.Clean("/mnt/usb") {
			t.Fatalf("expected BoundaryMountPoint at /mnt/usb, got %+v", b)
		}

		noMount := NewCache(WithFileSystem(m), WithStopAtMountPoint(false))
		b2, err := noMount.ResolveBoundary("/mnt/usb/proj/sub")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if b2.Reason != BoundaryFilesystemRoot {
			t.Fatalf("expected BoundaryFilesystemRoot, got %+v", b2)
		}
	})

	t.Run("ownership boundary", func(t *testing.T) {
		m := newMockFS("/shared/otheruser/myproj", 1000)
		m.setDir("/", 1, 0, now)
		m.setDir("/shared", 1, 0, now)
		m.setDir("/shared/otheruser", 1, 2000, now)
		m.setDir("/shared/otheruser/myproj", 1, 1000, now)

		c := NewCache(WithFileSystem(m))
		b, err := c.ResolveBoundary("/shared/otheruser/myproj")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if b.Reason != BoundaryOwnership ||
			b.RootDir != filepath.Clean("/shared/otheruser/myproj") {
			t.Fatalf("expected BoundaryOwnership, got %+v", b)
		}

		noOwn := NewCache(WithFileSystem(m), WithVerifyOwnership(false))
		b2, err := noOwn.ResolveBoundary("/shared/otheruser/myproj")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if b2.Reason != BoundaryFilesystemRoot {
			t.Fatalf("expected BoundaryFilesystemRoot, got %+v", b2)
		}
	})

	t.Run("ceiling dir and max depth", func(t *testing.T) {
		m := buildMockRepo(now)
		ceilCache := NewCache(WithFileSystem(m), WithCeilingDirs("/repo/sub"))
		cb, err := ceilCache.ResolveBoundary("/repo/sub/deep")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if cb.Reason != BoundaryCeiling ||
			cb.RootDir != filepath.Clean("/repo/sub") {
			t.Fatalf("expected BoundaryCeiling at /repo/sub, got %+v", cb)
		}

		depthCache := NewCache(WithFileSystem(m), WithMaxDepth(0))
		db, err := depthCache.ResolveBoundary("/repo/sub/deep")
		if err != nil {
			t.Fatalf("ResolveBoundary failed: %v", err)
		}
		if db.Reason != BoundaryMaxDepth ||
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
		opts       []Option
		wantReason SkipReason
	}{
		{
			name: "directory instead of regular file",
			setup: func(m *mockFS) string {
				m.setDir("/repo/.editorconfig", 1, 1000, now)
				return "/repo/.editorconfig"
			},
			wantReason: SkipNotRegular,
		},
		{
			name: "fifo named pipe",
			setup: func(m *mockFS) string {
				p := "/repo/.editorconfig"
				m.setFile(p, iofs.ModeNamedPipe|0644, 0, 1, 1000, now)
				return p
			},
			wantReason: SkipNotRegular,
		},
		{
			name: "unix socket",
			setup: func(m *mockFS) string {
				p := "/repo/.editorconfig"
				m.setFile(p, iofs.ModeSocket|0644, 0, 1, 1000, now)
				return p
			},
			wantReason: SkipNotRegular,
		},
		{
			name: "device node",
			setup: func(m *mockFS) string {
				p := "/repo/.editorconfig"
				m.setFile(p, iofs.ModeDevice|0644, 0, 1, 1000, now)
				return p
			},
			wantReason: SkipNotRegular,
		},
		{
			name: "oversized file",
			setup: func(m *mockFS) string {
				p := "/repo/.editorconfig"
				m.setFile(p, 0644, 2048, 1, 1000, now)
				return p
			},
			opts:       []Option{WithMaxFileSize(1024)},
			wantReason: SkipFileTooLarge,
		},
		{
			name: "symlinks disabled",
			setup: func(m *mockFS) string {
				m.setFile("/repo/shared.cfg", 0644, 16, 1, 1000, now)
				p := "/repo/.editorconfig"
				m.setSymlink(p, "/repo/shared.cfg", 1, 1000, now)
				return p
			},
			opts:       []Option{WithAllowSymlinks(false)},
			wantReason: SkipSymlinkDisabled,
		},
		{
			name: "broken symlink",
			setup: func(m *mockFS) string {
				p := "/repo/.editorconfig"
				m.setSymlink(p, "/repo/missing.cfg", 1, 1000, now)
				return p
			},
			wantReason: SkipBrokenSymlink,
		},
		{
			name: "symlink escaping boundary",
			setup: func(m *mockFS) string {
				m.setFile("/outside/secret.cfg", 0644, 16, 1, 1000, now)
				p := "/repo/.editorconfig"
				m.setSymlink(p, "/outside/secret.cfg", 1, 1000, now)
				return p
			},
			wantReason: SkipSymlinkEscape,
		},
		{
			name: "symlink crossing mount device",
			setup: func(m *mockFS) string {
				m.setFile("/repo/other.cfg", 0644, 16, 99, 1000, now)
				p := "/repo/.editorconfig"
				m.setSymlink(p, "/repo/other.cfg", 1, 1000, now)
				return p
			},
			wantReason: SkipSymlinkCrossDevice,
		},
		{
			name: "permission denied",
			setup: func(m *mockFS) string {
				p := "/repo/.editorconfig"
				m.setStatErr(p, iofs.ErrPermission)
				return p
			},
			wantReason: SkipPermissionDenied,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := buildMockRepo(now)
			target := tc.setup(m)

			var skips []SkipReason
			opts := append([]Option{}, tc.opts...)
			opts = append(
				opts,
				WithFileSystem(m),
				WithOnSkip(func(p string, r SkipReason, _ error) {
					if p == filepath.Clean(target) {
						skips = append(skips, r)
					}
				}),
			)

			finder, err := NewFinder(opts...)
			if err != nil {
				t.Fatalf("NewFinder failed: %v", err)
			}
			got, err := finder.FindChain("/repo/sub", FileEditorConfig)
			if err != nil {
				t.Fatalf("non-strict mode returned error: %v", err)
			}
			if len(got) != 0 {
				t.Fatalf("expected 0 files, got %+v", got)
			}
			if len(skips) != 1 || skips[0] != tc.wantReason {
				t.Fatalf("skips = %v, want [%v]", skips, tc.wantReason)
			}

			strictOpts := append([]Option{}, tc.opts...)
			strictOpts = append(
				strictOpts,
				WithFileSystem(m),
				WithStrictSecurity(true),
			)
			strictFinder, err := NewFinder(strictOpts...)
			if err != nil {
				t.Fatalf("NewFinder strict failed: %v", err)
			}
			_, err = strictFinder.FindChain("/repo/sub", FileEditorConfig)
			var secErr *SecurityError
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

	finder, err := NewFinder(WithFileSystem(m))
	if err != nil {
		t.Fatalf("NewFinder failed: %v", err)
	}
	files, err := finder.FindChain("/repo/sub", FileEditorConfig)
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

	allowEsc, err := NewFinder(
		WithFileSystem(m),
		WithAllowSymlinkEscape(true),
	)
	if err != nil {
		t.Fatalf("NewFinder failed: %v", err)
	}
	escFiles, err := allowEsc.FindChain("/repo/sub", FileEditorConfig)
	if err != nil {
		t.Fatalf("FindChain failed: %v", err)
	}
	if len(escFiles) != 2 {
		t.Fatalf("expected 2 files with escape allowed, got %d", len(escFiles))
	}
}

func TestRegexFiltersAndIncludeDiscovery(t *testing.T) {
	now := time.Unix(1700000000, 0)

	t.Run("invalid regex errors", func(t *testing.T) {
		bad := []struct {
			opt  Option
			flag string
		}{
			{WithInclude("["), FilterInclude},
			{WithExclude("["), FilterExclude},
			{WithExtendExclude("["), FilterExtendExclude},
			{WithForceExclude("["), FilterForceExclude},
		}
		for _, tc := range bad {
			_, err := NewFilter(tc.opt)
			if err == nil || !strings.Contains(err.Error(), tc.flag) ||
				!strings.Contains(err.Error(), "invalid regular expression") {
				t.Errorf("unexpected error for %s: %v", tc.flag, err)
			}
		}
	})

	t.Run("filter matching methods", func(t *testing.T) {
		m := buildMockRepo(now)
		flt, err := NewFilter(
			WithFileSystem(m),
			WithInclude(`\.dml$`),
			WithExclude(`/venv/`),
			WithExtendExclude(`/custom/`),
			WithForceExclude(`generated/.*\.dml$`),
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
		); got != FilterForceExclude {
			t.Errorf("got %q, want %q", got, FilterForceExclude)
		}
		if got := flt.MatchExclude(
			"venv/x.dml",
			"/repo/venv/x.dml",
			false,
		); got != FilterExclude {
			t.Errorf("got %q, want %q", got, FilterExclude)
		}
		if got := flt.MatchExclude(
			"custom/x.dml",
			"/repo/custom/x.dml",
			false,
		); got != FilterExtendExclude {
			t.Errorf("got %q, want %q", got, FilterExtendExclude)
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
		finder, err := NewFinder(
			WithFileSystem(m),
			WithInclude(`\.cfg$`),
			WithExclude(`/venv/`),
			WithExtendExclude(`skip\.cfg$`),
			WithStrictSecurity(true),
			WithOnSkip(func(path string, reason SkipReason, _ error) {
				if reason == SkipExcluded {
					skipped = append(skipped, path)
				}
			}),
		)
		if err != nil {
			t.Fatalf("NewFinder failed: %v", err)
		}

		chain, err := finder.FindChain("/repo/venv", FileEditorConfig)
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
			FileEditorConfig,
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

func TestCachingMutationsTTLAndConcurrency(t *testing.T) {
	t0 := time.Unix(1700000000, 0)
	clk := newMockClock(t0)

	t.Run("real-time stat validation detects FS mutations", func(t *testing.T) {
		m := buildMockRepo(t0)
		finder, err := NewFinder(WithFileSystem(m), WithClock(clk))
		if err != nil {
			t.Fatalf("NewFinder failed: %v", err)
		}

		files, err := finder.FindChain("/repo/sub", FileEditorConfig)
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

		_, err = finder.FindChain("/repo/sub", FileEditorConfig)
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
		files, err = finder.FindChain("/repo/sub", FileEditorConfig)
		if err != nil || len(files) != 1 || files[0].Size != 12 {
			t.Fatalf(
				"expected newly created file size 12, got %+v, %v",
				files,
				err,
			)
		}

		m.setFile("/repo/.editorconfig", 0644, 99, 1, 1000, t1.Add(time.Second))
		files, err = finder.FindChain("/repo/sub", FileEditorConfig)
		if err != nil || len(files) != 1 || files[0].Size != 99 {
			t.Fatalf("expected updated file size 99, got %+v, %v", files, err)
		}

		m.setFile("/repo/target.cfg", 0644, 7, 1, 1000, t1)
		m.setSymlink("/repo/.editorconfig", "/repo/target.cfg", 1, 1000, t1)
		files, err = finder.FindChain("/repo/sub", FileEditorConfig)
		if err != nil || len(files) != 1 || !files[0].IsSymlink ||
			files[0].Size != 7 {
			t.Fatalf("expected symlink size 7, got %+v, %v", files, err)
		}
		m.setFile("/repo/target.cfg", 0644, 33, 1, 1000, t1.Add(2*time.Second))
		files, err = finder.FindChain("/repo/sub", FileEditorConfig)
		if err != nil || len(files) != 1 || files[0].Size != 33 {
			t.Fatalf(
				"expected updated symlink target size 33, got %+v, %v",
				files,
				err,
			)
		}

		m.deleteNode("/repo/.editorconfig")
		files, err = finder.FindChain("/repo/sub", FileEditorConfig)
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
			finder, err := NewFinder(
				WithFileSystem(m),
				WithClock(clk),
				WithCacheTTL(5*time.Second),
			)
			if err != nil {
				t.Fatalf("NewFinder failed: %v", err)
			}

			if _, err := finder.FindChain(
				"/repo/sub",
				FileEditorConfig,
			); err != nil {
				t.Fatalf("FindChain failed: %v", err)
			}
			m.resetStatCalls()

			clk.Advance(2 * time.Second)
			if _, err := finder.FindChain(
				"/repo/sub",
				FileEditorConfig,
			); err != nil {
				t.Fatalf("FindChain failed: %v", err)
			}
			if calls := m.totalStatCalls(); calls != 1 {
				t.Errorf("expected 1 start stat call within TTL, got %d", calls)
			}

			finder.Cache().Invalidate("/repo/.editorconfig")
			if finder.Cache().Stats().Invalidations == 0 {
				t.Error("expected Invalidations > 0 after Invalidate")
			}
			if _, err := finder.FindChain(
				"/repo/sub",
				FileEditorConfig,
			); err != nil {
				t.Fatalf("FindChain failed: %v", err)
			}

			clk.Advance(10 * time.Second)
			finder.Cache().Reset()
			if _, err := finder.FindChain(
				"/repo/sub",
				FileEditorConfig,
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

		finder, err := NewFinder(
			WithFileSystem(m),
			WithClock(clk),
			WithInclude(`\.dml$`),
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
						FileEditorConfig,
						FileGitIgnore,
					)
					_, _, _ = finder.FindNearest(
						"/repo/sub/deep",
						FileEditorConfig,
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
	reasons := []BoundaryReason{
		BoundaryNone,
		BoundaryFilesystemRoot,
		BoundaryVCSRoot,
		BoundaryMountPoint,
		BoundaryOwnership,
		BoundaryCeiling,
		BoundaryMaxDepth,
		BoundaryReason(99),
	}
	for _, r := range reasons {
		if r.String() == "" {
			t.Errorf("empty String for BoundaryReason %d", r)
		}
	}

	skips := []SkipReason{
		SkipNotRegular,
		SkipFileTooLarge,
		SkipSymlinkDisabled,
		SkipBrokenSymlink,
		SkipSymlinkEscape,
		SkipSymlinkCrossDevice,
		SkipPermissionDenied,
		SkipExcluded,
		SkipReason(99),
	}
	for _, s := range skips {
		if s.String() == "" {
			t.Errorf("empty String for SkipReason %d", s)
		}
	}

	inner := errors.New("boom")
	se := &SecurityError{Path: "/a", Reason: SkipNotRegular, Err: inner}
	if !errors.Is(se, inner) {
		t.Error("SecurityError.Unwrap did not match inner error")
	}
	if !strings.Contains(se.Error(), "boom") {
		t.Errorf("unexpected Error(): %q", se.Error())
	}
	seNoErr := &SecurityError{Path: "/a", Reason: SkipNotRegular}
	if seNoErr.Error() == "" {
		t.Error("expected non-empty Error() without inner err")
	}
}

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
		matcher := NewIgnoreMatcher(WithFileSystem(m))

		rs, err := matcher.LoadFile(DiscoveredFile{
			Path: "/repo/.gitignore",
			Dir:  "/repo",
			Name: FileGitIgnore,
		})
		if err != nil {
			t.Fatalf("LoadFile failed: %v", err)
		}
		stack := []IgnoreRuleSet{rs}

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
			if got := matcher.Match(stack, p, false); got != FileGitIgnore {
				t.Errorf("Match(%q) = %q, want %q", p, got, FileGitIgnore)
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

		finder, err := NewFinder(
			WithFileSystem(m),
			WithInclude(`\.dml$`),
			WithIgnoreFiles(FileGitIgnore, FileIgnore),
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

		matcher := NewIgnoreMatcher(WithFileSystem(m))
		finder, err := NewFinder(
			WithFileSystem(m),
			WithInclude(`\.dml$`),
			WithIgnoreMatcher(matcher),
			WithIgnoreFiles(FileGitIgnore, FileIgnore),
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
		strictFinder, err := NewFinder(
			WithFileSystem(m),
			WithInclude(`\.dml$`),
			WithStrictSecurity(true),
			WithIgnoreFiles(FileGitIgnore),
		)
		if err != nil {
			t.Fatalf("NewFinder strict failed: %v", err)
		}
		_, err = strictFinder.Collect([]string{"/repo"})
		var secErr *SecurityError
		if !errors.As(err, &secErr) || secErr.Reason != SkipNotRegular {
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

		finder, err := NewFinder(
			WithFileSystem(m),
			WithInclude(`\.dml$`),
			WithIgnoreFiles(FileGitIgnore, FileIgnore),
			WithOnIgnore(func(path string, isDir bool, filterName string) {
				events = append(events, ignoreEvent{path, isDir, filterName})
			}),
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
			filter: FileGitIgnore,
		}
		if !slices.Contains(events, wantEvent) {
			t.Fatalf("expected ignore event %+v, got %+v", wantEvent, events)
		}
	})
}
