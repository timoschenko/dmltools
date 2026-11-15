// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lsp_test

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/timoschenko/dmltools/pkg/fs/stream"
	"github.com/timoschenko/dmltools/pkg/lsp"
	"go.lsp.dev/protocol"
)

func TestFsCache_OpenGetClose(t *testing.T) {
	cache := lsp.NewFsCache()
	testPath := "/tmp/test_device.dml"

	if cache.HasOverlay(testPath) {
		t.Fatal("expected no overlay initially")
	}

	initialText := "dml 1.4;\ndevice dev;\n"
	if err := cache.Open(testPath, 1, initialText); err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	if !cache.HasOverlay(testPath) {
		t.Fatal("expected overlay to be present")
	}

	entry, ok := cache.Get(testPath)
	if !ok || entry.Content != initialText || entry.Version != 1 {
		t.Fatalf("unexpected entry: %+v", entry)
	}
	if entry.Buffer == nil || entry.Buffer.GetRawString() != initialText {
		t.Fatalf("expected HashedBuffer to match initial content")
	}

	data, err := cache.ReadFile(testPath)
	if err != nil || string(data) != initialText {
		t.Fatalf("unexpected ReadFile data: %s, err: %v", string(data), err)
	}

	if err := cache.Close(testPath); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	if cache.HasOverlay(testPath) {
		t.Fatal("expected overlay to be removed after Close")
	}
}

func TestFsCache_ReadFileAndLookupFallback(t *testing.T) {
	cache := lsp.NewFsCache()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "disk_file.dml")
	diskContent := "dml 1.4;\n// on disk\n"
	if err := os.WriteFile(filePath, []byte(diskContent), 0o600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// 1. Read and Stat without overlay -> disk fallback
	data, err := cache.ReadFile(filePath)
	if err != nil || string(data) != diskContent {
		t.Fatalf(
			"expected disk content %q, got %q, err: %v",
			diskContent,
			string(data),
			err,
		)
	}

	meta, err := cache.LookupCache().Stat(filePath)
	if err != nil || !meta.Mode.IsRegular() ||
		meta.Size != int64(len(diskContent)) {
		t.Fatalf(
			"unexpected disk LookupCache.Stat meta: %+v, err: %v",
			meta,
			err,
		)
	}

	// 2. Open overlay with modified content
	memoryContent := "dml 1.4;\n// in memory buffer with extra text\n"
	if err := cache.Open(filePath, 2, memoryContent); err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	data, err = cache.ReadFile(filePath)
	if err != nil || string(data) != memoryContent {
		t.Fatalf(
			"expected memory content %q, got %q, err: %v",
			memoryContent,
			string(data),
			err,
		)
	}

	meta, err = cache.LookupCache().Stat(filePath)
	if err != nil || !meta.Mode.IsRegular() ||
		meta.Size != int64(len(memoryContent)) {
		t.Fatalf(
			"expected overlay size %d in LookupCache.Stat, got %+v, err: %v",
			len(memoryContent),
			meta,
			err,
		)
	}

	// 3. Close overlay -> falls back to disk content again
	if err := cache.Close(filePath); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	data, err = cache.ReadFile(filePath)
	if err != nil || string(data) != diskContent {
		t.Fatalf(
			"expected disk content %q after close, got %q, err: %v",
			diskContent,
			string(data),
			err,
		)
	}

	meta, err = cache.LookupCache().Stat(filePath)
	if err != nil || meta.Size != int64(len(diskContent)) {
		t.Fatalf("expected disk size after close, got %+v, err: %v", meta, err)
	}
}

func TestFsCache_ApplyEdit_Incremental(t *testing.T) {
	cache := lsp.NewFsCache()
	filePath := "/tmp/incremental.dml"

	initial := "dml 1.4;\ndevice dev;\n"
	if err := cache.Open(filePath, 1, initial); err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	rng := protocol.Range{
		Start: protocol.Position{Line: 1, Character: 7},
		End:   protocol.Position{Line: 1, Character: 10},
	}

	if err := cache.ApplyPartial(filePath, 2, rng, "my_device"); err != nil {
		t.Fatalf("ApplyPartial failed: %v", err)
	}

	expected := "dml 1.4;\ndevice my_device;\n"
	entry, ok := cache.Get(filePath)
	if !ok || entry.Content != expected {
		t.Fatalf("expected content %q, got %q", expected, entry.Content)
	}
	if entry.Version != 2 {
		t.Fatalf("expected version 2, got %d", entry.Version)
	}
}

func TestFsCache_ApplyEdit_MultiLine(t *testing.T) {
	cache := lsp.NewFsCache()
	filePath := "/tmp/multiline.dml"

	initial := "line 0\nline 1\nline 2\nline 3\n"
	if err := cache.Open(filePath, 1, initial); err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	rng := protocol.Range{
		Start: protocol.Position{Line: 1, Character: 5},
		End:   protocol.Position{Line: 2, Character: 6},
	}

	if err := cache.ApplyPartial(filePath, 2, rng, "X"); err != nil {
		t.Fatalf("ApplyPartial failed: %v", err)
	}

	expected := "line 0\nline X\nline 3\n"
	entry, ok := cache.Get(filePath)
	if !ok || entry.Content != expected {
		t.Fatalf("expected content %q, got %q", expected, entry.Content)
	}
}

func TestFsCache_ApplyEdit_UTF16(t *testing.T) {
	cache := lsp.NewFsCache()
	filePath := "/tmp/utf16.dml"

	// Emoji 🙂 (U+1F642) is 4 bytes in UTF-8, but 2 code units in UTF-16.
	initial := "hello 🙂 world\n"
	if err := cache.Open(filePath, 1, initial); err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	rng := protocol.Range{
		Start: protocol.Position{Line: 0, Character: 6},
		End:   protocol.Position{Line: 0, Character: 8},
	}

	if err := cache.ApplyPartial(filePath, 2, rng, "🌍"); err != nil {
		t.Fatalf("ApplyPartial failed: %v", err)
	}

	expected := "hello 🌍 world\n"
	entry, ok := cache.Get(filePath)
	if !ok || entry.Content != expected {
		t.Fatalf("expected content %q, got %q", expected, entry.Content)
	}

	// Verify UTF-16 LineAndByteCol translates UTF-16 col 9 ("hello 🌍 ") to
	// byte offset 11 ('w' in "world").
	line, byteCol, ok := entry.LineAndByteCol(
		protocol.Position{Line: 0, Character: 9},
	)
	if !ok || line != "hello 🌍 world" || line[byteCol] != 'w' {
		t.Fatalf(
			"unexpected LineAndByteCol: line=%q byteCol=%d ok=%v",
			line,
			byteCol,
			ok,
		)
	}
}

func TestFsCache_ApplyChange_WholeDocument(t *testing.T) {
	cache := lsp.NewFsCache()
	filePath := "/tmp/whole_doc.dml"

	if err := cache.Open(filePath, 1, "old text"); err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	change := &protocol.TextDocumentContentChangeWholeDocument{
		Text: "entirely new text",
	}

	if err := cache.ApplyChange(filePath, 2, change); err != nil {
		t.Fatalf("ApplyChange failed: %v", err)
	}

	entry, ok := cache.Get(filePath)
	if !ok || entry.Content != "entirely new text" || entry.Version != 2 {
		t.Fatalf("unexpected entry: %+v", entry)
	}
}

func TestFsCache_GetOrParse_CachingAndInvalidation(t *testing.T) {
	cache := lsp.NewFsCache()
	tmpDir := t.TempDir()
	diskPath := filepath.Join(tmpDir, "disk_cached.dml")

	if err := os.WriteFile(
		diskPath,
		[]byte("dml 1.4;\nconstant C1 = 1;\n"),
		0o600,
	); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	buf1, tree1, syms1, err := cache.GetOrParse(diskPath)
	if err != nil || buf1 == nil || tree1 == nil || syms1 == nil {
		t.Fatalf("GetOrParse disk failed: err=%v", err)
	}
	if _, ok := syms1.GlobalScope.Symbols["C1"]; !ok {
		t.Fatalf("expected C1 in disk symbols")
	}

	// Second call with identical hash returns same cached AST/Symbols pointers.
	buf2, tree2, syms2, err := cache.GetOrParse(diskPath)
	if err != nil {
		t.Fatalf("second GetOrParse failed: %v", err)
	}
	if tree2 != tree1 || syms2 != syms1 || buf2.Hash() != buf1.Hash() {
		t.Fatalf("expected GetOrParse to return cached AST and Symbols")
	}

	// Modifying disk file invalidates by hash.
	if err := os.WriteFile(
		diskPath,
		[]byte("dml 1.4;\nconstant C2 = 2;\n"),
		0o600,
	); err != nil {
		t.Fatalf("WriteFile update failed: %v", err)
	}
	_, tree3, syms3, err := cache.GetOrParse(diskPath)
	if err != nil {
		t.Fatalf("third GetOrParse failed: %v", err)
	}
	if tree3 == tree1 || syms3 == syms1 {
		t.Fatalf("expected updated disk file to re-parse")
	}
	if _, ok := syms3.GlobalScope.Symbols["C2"]; !ok {
		t.Fatalf("expected C2 in updated disk symbols")
	}

	// Opening overlay caches overlay AST/Symbols and invalidates on ApplyFull.
	if err := cache.Open(
		diskPath,
		1,
		"dml 1.4;\nconstant C3 = 3;\n",
	); err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	_, tree4, syms4, err := cache.GetOrParse(diskPath)
	if err != nil {
		t.Fatalf("overlay GetOrParse failed: %v", err)
	}
	if _, ok := syms4.GlobalScope.Symbols["C3"]; !ok {
		t.Fatalf("expected C3 in overlay symbols")
	}
	_, tree5, syms5, err := cache.GetOrParse(diskPath)
	if err != nil || tree5 != tree4 || syms5 != syms4 {
		t.Fatalf("expected overlay cache hit")
	}

	if err := cache.ApplyFull(
		diskPath,
		2,
		"dml 1.4;\nconstant C4 = 4;\n",
	); err != nil {
		t.Fatalf("ApplyFull failed: %v", err)
	}
	_, tree6, syms6, err := cache.GetOrParse(diskPath)
	if err != nil || tree6 == tree4 || syms6 == syms4 {
		t.Fatalf("expected ApplyFull to invalidate cached AST and Symbols")
	}
	if _, ok := syms6.GlobalScope.Symbols["C4"]; !ok {
		t.Fatalf("expected C4 after ApplyFull")
	}
}

func TestFsCache_GetOrParse_SyntaxErrorNoStderr(t *testing.T) {
	origStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed: %v", err)
	}
	os.Stderr = w
	defer func() {
		os.Stderr = origStderr
		_ = r.Close()
	}()

	cache := lsp.NewFsCache()
	tmpDir := t.TempDir()

	broken := "dml 1.4;\ndevice;\n"
	diskPath := filepath.Join(tmpDir, "broken_disk.dml")
	if err := os.WriteFile(diskPath, []byte(broken), 0o600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	if _, _, _, err := cache.GetOrParse(diskPath); err != nil {
		t.Fatalf("disk GetOrParse failed: %v", err)
	}

	overlayPath := filepath.Join(tmpDir, "broken_overlay.dml")
	if err := cache.Open(overlayPath, 1, broken); err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	if _, _, _, err := cache.GetOrParse(overlayPath); err != nil {
		t.Fatalf("overlay GetOrParse failed: %v", err)
	}

	_ = w.Close()
	captured, readErr := io.ReadAll(r)
	if readErr != nil {
		t.Fatalf("ReadAll failed: %v", readErr)
	}
	if len(captured) != 0 {
		t.Fatalf("expected no stderr output, got: %q", string(captured))
	}
}

func TestFsCache_SetASTIfMatch(t *testing.T) {
	cache := lsp.NewFsCache()
	path := "/tmp/ast_match.dml"

	if err := cache.Open(path, 1, "dml 1.4;\ndevice a;\n"); err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	entry1, ok := cache.Get(path)
	if !ok {
		t.Fatal("expected entry")
	}
	oldHash := entry1.Buffer.Hash()

	_, tree1, _, err := cache.GetOrParse(path)
	if err != nil {
		t.Fatalf("GetOrParse failed: %v", err)
	}

	if err := cache.ApplyFull(path, 2, "dml 1.4;\ndevice b;\n"); err != nil {
		t.Fatalf("ApplyFull failed: %v", err)
	}

	if cache.SetASTIfMatch(path, oldHash, tree1) {
		t.Fatal("expected SetASTIfMatch to reject stale hash")
	}
	if _, hasAST := cache.GetAST(path); hasAST {
		t.Fatal("expected AST to remain unset after stale hash rejection")
	}

	entry2, _ := cache.Get(path)
	newHash := entry2.Buffer.Hash()
	if !cache.SetASTIfMatch(path, newHash, tree1) {
		t.Fatal("expected SetASTIfMatch to succeed with matching hash")
	}
	if !cache.HasOverlayWithHash(path, newHash) {
		t.Fatal("expected HasOverlayWithHash to return true")
	}
	var wrongHash stream.Hash
	if cache.HasOverlayWithHash(path, wrongHash) {
		t.Fatal("expected HasOverlayWithHash to return false for wrong hash")
	}
}

func TestFsCache_Concurrency(t *testing.T) {
	cache := lsp.NewFsCache()
	filePath := "/tmp/concurrent.dml"
	if err := cache.Open(filePath, 0, "initial text\n"); err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_, _ = cache.ReadFile(filePath)
		}()
		go func(v int32) {
			defer wg.Done()
			_ = cache.ApplyFull(filePath, v, "updated text\n")
		}(int32(i))
	}
	wg.Wait()
}

func TestFsCache_DocumentSymbolsDiskCachingAndInvalidation(t *testing.T) {
	cache := lsp.NewFsCache()
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "disk_doc_syms.dml")
	v1Content := "dml 1.4;\nbank b { register r; }\n"
	if err := os.WriteFile(filePath, []byte(v1Content), 0o600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Compute document symbols on disk file
	syms1, err := lsp.ComputeDocumentSymbols(ctx, cache, filePath)
	if err != nil {
		t.Fatalf("ComputeDocumentSymbols failed: %v", err)
	}
	if len(syms1) != 1 || syms1[0].Name != "b" {
		t.Fatalf("expected [b], got %+v", syms1)
	}

	// 2. Call GetOrParse on same disk file: should reuse cached AST
	// and preserve docSymbols
	buf, tree, fileSyms, err := cache.GetOrParse(filePath)
	if err != nil {
		t.Fatalf("GetOrParse failed: %v", err)
	}
	if buf == nil || tree == nil || fileSyms == nil {
		t.Fatalf("unexpected nil from GetOrParse")
	}

	// Second ComputeDocumentSymbols call should hit cache
	syms2, err := lsp.ComputeDocumentSymbols(ctx, cache, filePath)
	if err != nil {
		t.Fatalf("second ComputeDocumentSymbols failed: %v", err)
	}
	if len(syms2) != 1 || syms2[0].Name != "b" {
		t.Fatalf("expected [b] after GetOrParse, got %+v", syms2)
	}

	// 3. Modify disk file to new content
	v2Content := "dml 1.4;\nbank b2 { register r2; }\n"
	if err := os.WriteFile(filePath, []byte(v2Content), 0o600); err != nil {
		t.Fatalf("WriteFile v2 failed: %v", err)
	}

	// Compute semantic tokens on updated disk file (must not carry stale
	// docSymbols into new hash entry)
	tokens, err := lsp.ComputeSemanticTokens(ctx, cache, filePath)
	if err != nil || tokens == nil {
		t.Fatalf("ComputeSemanticTokens failed: %v", err)
	}

	// ComputeDocumentSymbols on updated disk file must reflect v2 content ("b2"),
	// not stale "b"
	syms3, err := lsp.ComputeDocumentSymbols(ctx, cache, filePath)
	if err != nil {
		t.Fatalf("ComputeDocumentSymbols v2 failed: %v", err)
	}
	if len(syms3) != 1 || syms3[0].Name != "b2" {
		t.Fatalf("expected [b2], got %+v (stale cache detected)", syms3)
	}
}
