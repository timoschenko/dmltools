// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lookup

import (
	"path/filepath"
	"testing"
	"time"
)

func TestDirStore_BasicsAndStats(t *testing.T) {
	store := newDirStore(100 * time.Millisecond)

	initialStats := store.stats()
	if initialStats != (CacheStats{}) {
		t.Fatalf("expected zero stats, got %+v", initialStats)
	}

	t0 := time.Unix(1000, 0)
	meta := FileMeta{Name: "dir", Mode: 0755}
	store.saveNewDir("/repo", meta, true, "/repo/.git", nil, t0, false)

	stats := store.stats()
	if stats.DirMisses != 1 {
		t.Errorf("expected 1 DirMiss, got %d", stats.DirMisses)
	}

	snap, ttlHit, ok := store.lookupDir("/repo", t0.Add(50*time.Millisecond))
	if !ok || !ttlHit {
		t.Fatalf("expected TTL hit, got ok=%v, ttlHit=%v", ok, ttlHit)
	}
	if snap.meta.Name != "dir" || !snap.isVCSRoot ||
		snap.vcsPath != "/repo/.git" {
		t.Errorf("unexpected snapshot: %+v", snap)
	}

	stats = store.stats()
	if stats.DirHits != 1 {
		t.Errorf("expected 1 DirHit, got %d", stats.DirHits)
	}

	// After TTL expires:
	_, ttlHit, ok = store.lookupDir("/repo", t0.Add(200*time.Millisecond))
	if !ok || ttlHit {
		t.Errorf("expected expired TTL hit, got ok=%v, ttlHit=%v", ok, ttlHit)
	}

	store.touchDir("/repo", snap.entry, t0.Add(250*time.Millisecond))
	stats = store.stats()
	if stats.DirHits != 2 {
		t.Errorf("expected 2 DirHits after touchDir, got %d", stats.DirHits)
	}

	store.updateDirVCS(
		"/repo",
		snap.entry,
		false,
		"",
		nil,
		t0.Add(260*time.Millisecond),
	)
	stats = store.stats()
	if stats.DirHits != 3 {
		t.Errorf("expected 3 DirHits after updateDirVCS, got %d", stats.DirHits)
	}

	store.reset()
	stats = store.stats()
	if stats.Invalidations != 1 {
		t.Errorf(
			"expected 1 invalidation after reset, got %d",
			stats.Invalidations,
		)
	}

	_, _, ok = store.lookupDir("/repo", t0.Add(300*time.Millisecond))
	if ok {
		t.Error("expected dir to be evicted after reset")
	}
}

func TestDirStore_InvalidateTargets(t *testing.T) {
	store := newDirStore(0)
	t0 := time.Unix(1000, 0)
	m := FileMeta{Mode: 0755}

	store.saveNewDir(filepath.Clean("/repo"), m, false, "", nil, t0, false)
	store.saveNewDir(filepath.Clean("/repo/sub"), m, false, "", nil, t0, false)
	store.saveNewDir(
		filepath.Clean("/repo/sub/deep"),
		m,
		false,
		"",
		nil,
		t0,
		false,
	)
	store.saveNewDir(filepath.Clean("/other"), m, false, "", nil, t0, false)

	// Invalidate /repo/sub/deep -> affects /repo/sub/deep and its parent /repo/sub
	store.invalidateTargets([]string{filepath.Clean("/repo/sub/deep")})

	_, _, okDeep := store.lookupDir(filepath.Clean("/repo/sub/deep"), t0)
	_, _, okSub := store.lookupDir(filepath.Clean("/repo/sub"), t0)
	_, _, okRepo := store.lookupDir(filepath.Clean("/repo"), t0)
	_, _, okOther := store.lookupDir(filepath.Clean("/other"), t0)

	if okDeep || okSub {
		t.Errorf(
			"expected deep and sub to be evicted: okDeep=%v, okSub=%v",
			okDeep,
			okSub,
		)
	}
	if !okRepo || !okOther {
		t.Errorf(
			"expected repo and other to remain: okRepo=%v, okOther=%v",
			okRepo,
			okOther,
		)
	}
}

func TestDirStore_FileCacheAndCounters(t *testing.T) {
	store := newDirStore(100 * time.Millisecond)
	t0 := time.Unix(1000, 0)

	store.saveNewDir("/repo", FileMeta{Mode: 0755}, false, "", nil, t0, false)
	dirEntry, _, ok := store.lookupFile("/repo", "file.txt")
	if ok || dirEntry == nil {
		t.Fatalf(
			"expected no file, but dirEntry present: ok=%v, dirEntry=%v",
			ok,
			dirEntry,
		)
	}

	fe := fileCacheEntry{
		validatedAt: t0,
		exists:      true,
		discovered:  DiscoveredFile{Path: "/repo/file.txt", Name: "file.txt"},
	}
	store.storeFile("/repo", dirEntry, "file.txt", fe)

	_, gotFE, ok := store.lookupFile("/repo", "file.txt")
	if !ok || gotFE.discovered.Path != "/repo/file.txt" {
		t.Fatalf("expected file to be found, got ok=%v, fe=%+v", ok, gotFE)
	}

	store.recordFileHit()
	store.recordFileMiss(true)
	stats := store.stats()
	if stats.FileHits != 1 || stats.FileMisses != 1 ||
		stats.Invalidations != 1 {
		t.Errorf("unexpected file stats: %+v", stats)
	}
}
