// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lookup

import (
	iofs "io/fs"
	"math"
	"path/filepath"
	"strings"
)

func (c *securityCache) ResolveBoundary(
	startPath string,
) (BoundaryInfo, error) {
	startDir, _, _, startMeta, err := c.ResolveStart(startPath)
	if err != nil {
		return BoundaryInfo{}, err
	}
	_, boundary, _, err := c.ResolveDirChain(startDir, startMeta)
	if err != nil {
		return BoundaryInfo{}, err
	}
	return boundary, nil
}

func (c *securityCache) ResolveStart(
	startPath string,
) (string, string, bool, FileMeta, error) {
	if startPath == "" {
		startPath = "."
	}
	absPath, err := c.fsys.Abs(startPath)
	if err != nil {
		return "", "", false, FileMeta{}, err
	}
	absPath = filepath.Clean(absPath)

	meta, err := c.fsys.Stat(absPath)
	if err == nil {
		if meta.Mode.IsDir() {
			return absPath, absPath, false, meta, nil
		}
		startDir := filepath.Dir(absPath)
		dirMeta, dirErr := c.fsys.Stat(startDir)
		if dirErr != nil {
			return "", "", false, FileMeta{}, dirErr
		}
		return startDir, absPath, true, dirMeta, nil
	}

	lmeta, lerr := c.fsys.Lstat(absPath)
	if lerr != nil {
		return "", "", false, FileMeta{}, err
	}
	if lmeta.Mode.IsDir() {
		return absPath, absPath, false, lmeta, nil
	}
	startDir := filepath.Dir(absPath)
	dirMeta, dirErr := c.fsys.Stat(startDir)
	if dirErr != nil {
		return "", "", false, FileMeta{}, dirErr
	}
	return startDir, absPath, true, dirMeta, nil
}

func (c *securityCache) ResolveDirChain(
	startDir string,
	startMeta FileMeta,
) ([]string, BoundaryInfo, FileMeta, error) {
	dirs := make([]string, 0, 8)
	currentDir := startDir
	depth := 0
	euid := c.fsys.Geteuid()

	curMeta, isVCS, vcsPath, err := c.getOrRefreshDir(currentDir)
	if err != nil {
		return nil, BoundaryInfo{}, FileMeta{}, err
	}
	startMeta = curMeta

	for {
		dirs = append(dirs, currentDir)

		if c.isCeiling(currentDir) {
			return dirs, BoundaryInfo{
				RootDir: currentDir,
				Reason:  BoundaryCeiling,
			}, startMeta, nil
		}
		if c.stopAtVCSRoot && isVCS {
			return dirs, BoundaryInfo{
				RootDir: currentDir,
				Reason:  BoundaryVCSRoot,
				VCSPath: vcsPath,
			}, startMeta, nil
		}
		if depth >= c.maxDepth {
			return dirs, BoundaryInfo{
				RootDir: currentDir,
				Reason:  BoundaryMaxDepth,
			}, startMeta, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return dirs, BoundaryInfo{
				RootDir: currentDir,
				Reason:  BoundaryFilesystemRoot,
			}, startMeta, nil
		}

		parentMeta, parentIsVCS, parentVCSPath, parentErr := c.getOrRefreshDir(
			parentDir,
		)
		if parentErr != nil {
			return dirs, BoundaryInfo{
				RootDir: currentDir,
				Reason:  BoundaryOwnership,
			}, startMeta, nil
		}

		if c.stopAtMount && startMeta.HasDevice && parentMeta.HasDevice &&
			parentMeta.DeviceID != startMeta.DeviceID {
			return dirs, BoundaryInfo{
				RootDir: currentDir,
				Reason:  BoundaryMountPoint,
			}, startMeta, nil
		}

		if c.verifyOwnership && startMeta.HasUID && parentMeta.HasUID {
			if !isTrustedOwner(parentMeta.UID, startMeta.UID, euid) {
				return dirs, BoundaryInfo{
					RootDir: currentDir,
					Reason:  BoundaryOwnership,
				}, startMeta, nil
			}
		}

		currentDir = parentDir
		isVCS = parentIsVCS
		vcsPath = parentVCSPath
		depth++
	}
}

func isTrustedOwner(parentUID, startUID uint32, euid int) bool {
	if parentUID == startUID || parentUID == 0 {
		return true
	}
	if euid >= 0 && int64(euid) <= math.MaxUint32 {
		// #nosec G115 -- guarded against negative and > MaxUint32
		return parentUID == uint32(euid)
	}
	return false
}

func (c *securityCache) isCeiling(dir string) bool {
	if len(c.ceilingDirs) == 0 {
		return false
	}
	if _, ok := c.ceilingDirs[dir]; ok {
		return true
	}
	if evalDir, err := c.fsys.EvalSymlinks(dir); err == nil {
		_, ok := c.ceilingDirs[filepath.Clean(evalDir)]
		return ok
	}
	return false
}

func (c *securityCache) getOrRefreshDir(
	dir string,
) (FileMeta, bool, string, error) {
	now := c.clock.Now()

	snap, ttlHit, ok := c.store.lookupDir(dir, now)
	if ok && ttlHit {
		return snap.meta, snap.isVCSRoot, snap.vcsPath, nil
	}

	curMeta, err := c.fsys.Stat(dir)
	if err != nil {
		if ok {
			c.store.evictDir(dir, snap.entry)
		}
		return FileMeta{}, false, "", err
	}

	if ok && snap.meta.sameFingerprint(curMeta) {
		if snap.vcsProbed && c.markersUnchanged(dir, snap.markers) {
			c.store.touchDir(dir, snap.entry, now)
			return snap.meta, snap.isVCSRoot, snap.vcsPath, nil
		}
		if !snap.vcsProbed {
			isVCS, vcsPath, markers := c.probeVCSMarkers(dir)
			c.store.updateDirVCS(dir, snap.entry, isVCS, vcsPath, markers, now)
			return snap.meta, isVCS, vcsPath, nil
		}
	}

	isVCS, vcsPath, markers := c.probeVCSMarkers(dir)
	c.store.saveNewDir(dir, curMeta, isVCS, vcsPath, markers, now, ok)
	return curMeta, isVCS, vcsPath, nil
}

func (c *securityCache) probeVCSMarkers(
	dir string,
) (bool, string, []vcsMarkerState) {
	if !c.stopAtVCSRoot || len(c.vcsMarkers) == 0 {
		return false, "", nil
	}
	states := make([]vcsMarkerState, 0, len(c.vcsMarkers))
	isVCS := false
	vcsPath := ""

	for _, marker := range c.vcsMarkers {
		p := filepath.Join(dir, marker)
		lmeta, err := c.fsys.Lstat(p)
		if err != nil {
			states = append(states, vcsMarkerState{
				marker: marker,
				exists: false,
			})
			continue
		}
		states = append(states, vcsMarkerState{
			marker: marker,
			exists: true,
			meta:   lmeta,
		})
		if isVCS {
			continue
		}
		if lmeta.Mode.IsDir() || lmeta.Mode.IsRegular() {
			isVCS = true
			vcsPath = p
			continue
		}
		if lmeta.Mode&iofs.ModeSymlink != 0 {
			if smeta, serr := c.fsys.Stat(p); serr == nil &&
				(smeta.Mode.IsDir() || smeta.Mode.IsRegular()) {
				isVCS = true
				vcsPath = p
			}
		}
	}

	return isVCS, vcsPath, states
}

func (c *securityCache) markersUnchanged(
	dir string,
	states []vcsMarkerState,
) bool {
	if !c.stopAtVCSRoot {
		return len(states) == 0
	}
	if len(states) != len(c.vcsMarkers) {
		return false
	}
	for _, st := range states {
		p := filepath.Join(dir, st.marker)
		lmeta, err := c.fsys.Lstat(p)
		if err != nil {
			if st.exists {
				return false
			}
			continue
		}
		if !st.exists || !st.meta.sameFingerprint(lmeta) {
			return false
		}
	}
	return true
}

func (c *securityCache) EvalBoundaryRoot(boundaryRoot string) string {
	if evalRoot, err := c.fsys.EvalSymlinks(boundaryRoot); err == nil {
		return filepath.Clean(evalRoot)
	}
	return filepath.Clean(boundaryRoot)
}

func (c *securityCache) absClean(path string) string {
	if abs, err := c.fsys.Abs(path); err == nil {
		return filepath.Clean(abs)
	}
	return filepath.Clean(path)
}

func (c *securityCache) IsWithinBoundary(
	boundaryRoot string,
	targetPath string,
) bool {
	if boundaryRoot == "" || targetPath == "" {
		return false
	}
	absBoundary := c.absClean(boundaryRoot)
	evalBoundary := c.EvalBoundaryRoot(absBoundary)
	absTarget := c.absClean(targetPath)

	if evalTarget, err := c.fsys.EvalSymlinks(absTarget); err == nil {
		cleanEval := filepath.Clean(evalTarget)
		if isWithinBoundary(evalBoundary, cleanEval) {
			return true
		}
		if !c.allowSymlinkEsc {
			return false
		}
	} else if evalDir, dirErr := c.fsys.EvalSymlinks(
		filepath.Dir(absTarget),
	); dirErr == nil {
		cand := filepath.Join(
			filepath.Clean(evalDir),
			filepath.Base(absTarget),
		)
		if isWithinBoundary(evalBoundary, cand) {
			return true
		}
		if !c.allowSymlinkEsc {
			return false
		}
	}

	return isWithinBoundary(absBoundary, absTarget) ||
		isWithinBoundary(evalBoundary, absTarget)
}

func isWithinBoundary(boundaryRoot, target string) bool {
	cleanRoot := filepath.Clean(boundaryRoot)
	cleanTarget := filepath.Clean(target)
	rel, err := filepath.Rel(cleanRoot, cleanTarget)
	if err != nil {
		return false
	}
	if filepath.IsAbs(rel) || rel == ".." {
		return false
	}
	return !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
