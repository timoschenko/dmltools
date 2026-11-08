// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test_fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/timoschenko/dmltools/pkg/fs"
)

func TestReplaceFileAtomically_ExistingFile(t *testing.T) {
	dir := t.TempDir()
	targetPath := filepath.Join(dir, "target.dml")
	tempPath := filepath.Join(dir, "temp.dml.tmp")

	if err := os.WriteFile(
		targetPath,
		[]byte("original content"),
		0600,
	); err != nil {
		t.Fatalf("failed to write target file: %v", err)
	}
	if err := os.WriteFile(tempPath, []byte("new content"), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	err := fs.ReplaceFileAtomically(targetPath, tempPath, 0644)
	if err != nil {
		t.Fatalf("ReplaceFileAtomically failed: %v", err)
	}

	// Verify target content
	content, err := os.ReadFile(filepath.Clean(targetPath))
	if err != nil {
		t.Fatalf("failed to read replaced target file: %v", err)
	}
	if string(content) != "new content" {
		t.Errorf(
			"got target content %q, want %q",
			string(content),
			"new content",
		)
	}

	// Verify tempPath was moved
	if _, err := os.Stat(tempPath); !os.IsNotExist(err) {
		t.Errorf("tempPath %s still exists", tempPath)
	}

	// Verify no backup or leftover files in directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read dir: %v", err)
	}
	if len(entries) != 1 || entries[0].Name() != "target.dml" {
		t.Errorf("unexpected directory entries after replace: %v", entries)
	}
}

func TestReplaceFileAtomically_NonExistentTarget(t *testing.T) {
	dir := t.TempDir()
	targetPath := filepath.Join(dir, "new_file.dml")
	tempPath := filepath.Join(dir, "temp.dml.tmp")

	if err := os.WriteFile(
		tempPath,
		[]byte("initial content"),
		0600,
	); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	err := fs.ReplaceFileAtomically(targetPath, tempPath, 0644)
	if err != nil {
		t.Fatalf("ReplaceFileAtomically failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Clean(targetPath))
	if err != nil {
		t.Fatalf("failed to read target file: %v", err)
	}
	if string(content) != "initial content" {
		t.Errorf(
			"got target content %q, want %q",
			string(content),
			"initial content",
		)
	}

	if _, err := os.Stat(tempPath); !os.IsNotExist(err) {
		t.Errorf("tempPath %s still exists", tempPath)
	}
}

func TestReplaceFileAtomically_IdenticalPaths(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "same.dml")

	if err := os.WriteFile(path, []byte("same content"), 0600); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	err := fs.ReplaceFileAtomically(path, path, 0644)
	if err != nil {
		t.Fatalf("ReplaceFileAtomically failed on identical paths: %v", err)
	}

	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(content) != "same content" {
		t.Errorf("got %q, want %q", string(content), "same content")
	}
}

func TestReplaceFileAtomically_MissingTempWithMode(t *testing.T) {
	dir := t.TempDir()
	targetPath := filepath.Join(dir, "target.dml")
	tempPath := filepath.Join(dir, "does_not_exist.tmp")

	if err := os.WriteFile(
		targetPath,
		[]byte("original content"),
		0600,
	); err != nil {
		t.Fatalf("failed to write target file: %v", err)
	}

	err := fs.ReplaceFileAtomically(targetPath, tempPath, 0644)
	if err == nil {
		t.Fatal("expected error when tempPath does not exist, got nil")
	}

	// Verify target file is completely untouched
	content, err := os.ReadFile(filepath.Clean(targetPath))
	if err != nil {
		t.Fatalf("failed to read target file: %v", err)
	}
	if string(content) != "original content" {
		t.Errorf(
			"got target content %q, want %q",
			string(content),
			"original content",
		)
	}
}

func TestReplaceFileAtomically_RollbackOnFailure(t *testing.T) {
	dir := t.TempDir()
	targetPath := filepath.Join(dir, "target.dml")
	tempPath := filepath.Join(dir, "nonexistent.tmp")

	if err := os.WriteFile(
		targetPath,
		[]byte("rollback original"),
		0600,
	); err != nil {
		t.Fatalf("failed to write target file: %v", err)
	}

	// origMode = 0 skips chmod so that the rename error triggers the rollback path
	err := fs.ReplaceFileAtomically(targetPath, tempPath, 0)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Verify target file was restored by rollback
	content, err := os.ReadFile(filepath.Clean(targetPath))
	if err != nil {
		t.Fatalf("failed to read restored target file: %v", err)
	}
	if string(content) != "rollback original" {
		t.Errorf("got %q, want %q", string(content), "rollback original")
	}

	// Verify no stray backup files remain
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read dir: %v", err)
	}
	if len(entries) != 1 || entries[0].Name() != "target.dml" {
		t.Errorf("unexpected directory entries after rollback: %v", entries)
	}
}

func TestReplaceFileAtomically_PreservesPermissions(t *testing.T) {
	dir := t.TempDir()
	targetPath := filepath.Join(dir, "target.dml")
	tempPath := filepath.Join(dir, "temp.dml.tmp")

	if err := os.WriteFile(targetPath, []byte("old"), 0600); err != nil {
		t.Fatalf("failed to write target: %v", err)
	}
	//nolint:gosec // testing preservation of permission modes
	if err := os.Chmod(targetPath, 0755); err != nil {
		t.Fatalf("failed to chmod target: %v", err)
	}
	if err := os.WriteFile(tempPath, []byte("new"), 0600); err != nil {
		t.Fatalf("failed to write temp: %v", err)
	}

	if err := fs.ReplaceFileAtomically(targetPath, tempPath, 0755); err != nil {
		t.Fatalf("ReplaceFileAtomically failed: %v", err)
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		t.Fatalf("failed to stat target: %v", err)
	}
	// On Unix, check permission bits (0755)
	if info.Mode().Perm() != 0755 {
		t.Errorf("got mode %v, want 0755", info.Mode().Perm())
	}
}
