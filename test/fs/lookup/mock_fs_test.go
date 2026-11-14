// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test_lookup

import (
	"errors"
	iofs "io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/timoschenko/dmltools/pkg/fs/lookup"
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
	meta         lookup.FileMeta
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
	meta lookup.FileMeta
}

func (e mockDirEntry) Name() string        { return e.name }
func (e mockDirEntry) IsDir() bool         { return e.mode.IsDir() }
func (e mockDirEntry) Type() iofs.FileMode { return e.mode.Type() }

func (e mockDirEntry) Info() (iofs.FileInfo, error) { return mockFileInfo{meta: e.meta}, nil }

type mockFileInfo struct {
	meta lookup.FileMeta
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
		meta: lookup.FileMeta{
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
		meta: lookup.FileMeta{
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
		meta: lookup.FileMeta{
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
		meta: lookup.FileMeta{
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

func (m *mockFS) Stat(path string) (lookup.FileMeta, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	clean := m.absLocked(path)
	m.statCalls[clean]++
	for range 8 {
		n, ok := m.nodes[clean]
		if !ok {
			return lookup.FileMeta{}, os.ErrNotExist
		}
		if n.statErr != nil {
			return lookup.FileMeta{}, n.statErr
		}
		if n.symlinkTo == "" {
			return n.meta, nil
		}
		clean = n.symlinkTo
	}
	return lookup.FileMeta{}, errors.New("too many levels of symbolic links")
}

func (m *mockFS) Lstat(path string) (lookup.FileMeta, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	clean := m.absLocked(path)
	m.statCalls[clean]++
	n, ok := m.nodes[clean]
	if !ok {
		return lookup.FileMeta{}, os.ErrNotExist
	}
	if n.lstatErr != nil {
		return lookup.FileMeta{}, n.lstatErr
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
