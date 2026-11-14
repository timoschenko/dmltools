// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

//go:build unix

package test_fs

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/timoschenko/dmltools/pkg/fs/lookup"
)

func TestLookup_UnixRejectsFIFO(t *testing.T) {
	tmpDir := t.TempDir()
	repo := filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(filepath.Join(repo, ".git"), 0750); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	fifoPath := filepath.Join(repo, lookup.FileEditorConfig)
	if err := syscall.Mkfifo(fifoPath, 0600); err != nil {
		t.Fatalf("mkfifo: %v", err)
	}

	finder, err := lookup.NewFinder(lookup.WithStrictSecurity(true))
	if err != nil {
		t.Fatalf("lookup.NewFinder: %v", err)
	}

	_, err = finder.FindChain(repo, lookup.FileEditorConfig)
	var secErr *lookup.SecurityError
	if !errors.As(err, &secErr) ||
		secErr.Reason != lookup.SkipNotRegular {
		t.Fatalf("expected SkipNotRegular SecurityError for FIFO, got %v", err)
	}
}
