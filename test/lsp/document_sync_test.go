// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lsp_test

import (
	"context"
	"io"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/timoschenko/dmltools/pkg/lsp"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

type syncHarness struct {
	server *lsp.Server
	client protocol.Server
	done   chan error
}

func setupSyncHarness(t *testing.T) *syncHarness {
	t.Helper()

	clientPipe, serverPipe := net.Pipe()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := lsp.NewServer(lsp.WithLogger(logger))

	serverStream := jsonrpc2.NewStream(serverPipe)
	clientStream := jsonrpc2.NewStream(clientPipe)

	ctx := context.Background()
	done := make(chan error, 1)

	go func() {
		done <- server.Serve(ctx, serverStream)
	}()

	_, _, clientDispatcher := protocol.NewClient(
		ctx,
		protocol.UnimplementedClient{},
		clientStream,
	)

	return &syncHarness{
		server: server,
		client: clientDispatcher,
		done:   done,
	}
}

func waitFor(t *testing.T, timeout time.Duration, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("timed out waiting for condition")
}

func TestDocumentSync_OpenEditClose(t *testing.T) {
	h := setupSyncHarness(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testURI := uri.File("/workspace/test.dml")

	initParams := &protocol.InitializeParams{
		ClientInfo: protocol.ClientInfo{
			Name:    "test-client",
			Version: protocol.NewOptional("1.0.0"),
		},
		WorkspaceFolders: protocol.NewNullable([]protocol.WorkspaceFolder{
			{URI: uri.File("/workspace"), Name: "workspace"},
		}),
	}

	res, err := h.client.Initialize(ctx, initParams)
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	if res.Capabilities.TextDocumentSync == nil {
		t.Fatal("expected TextDocumentSync capability to be advertised")
	}

	if err := h.client.Initialized(
		ctx,
		&protocol.InitializedParams{},
	); err != nil {
		t.Fatalf("Initialized failed: %v", err)
	}

	// 1. Send didOpen
	initialContent := "dml 1.4;\ndevice dev;\n"
	openParams := &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        testURI,
			LanguageID: "dml",
			Version:    1,
			Text:       initialContent,
		},
	}

	if err := h.client.DidOpen(ctx, openParams); err != nil {
		t.Fatalf("DidOpen failed: %v", err)
	}

	waitFor(t, 2*time.Second, func() bool {
		data, err := h.server.FsCache().ReadFile(string(testURI))
		return err == nil && string(data) == initialContent
	})

	// 2. Send didChange (incremental: replace "dev" with "my_device")
	changeParams := &protocol.DidChangeTextDocumentParams{
		TextDocument: protocol.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: protocol.TextDocumentIdentifier{
				URI: testURI,
			},
			Version: 2,
		},
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			&protocol.TextDocumentContentChangePartial{
				Range: protocol.Range{
					Start: protocol.Position{Line: 1, Character: 7},
					End:   protocol.Position{Line: 1, Character: 10},
				},
				Text: "my_device",
			},
		},
	}

	if err := h.client.DidChange(ctx, changeParams); err != nil {
		t.Fatalf("DidChange failed: %v", err)
	}

	expectedUpdated := "dml 1.4;\ndevice my_device;\n"
	waitFor(t, 2*time.Second, func() bool {
		data, err := h.server.FsCache().ReadFile(string(testURI))
		return err == nil && string(data) == expectedUpdated
	})

	// 3. Send didSave (with Text) and verify version 2 is preserved
	savedContent := "dml 1.4;\ndevice saved_device;\n"
	saveParams := &protocol.DidSaveTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: testURI,
		},
		Text: &savedContent,
	}
	if err := h.client.DidSave(ctx, saveParams); err != nil {
		t.Fatalf("DidSave failed: %v", err)
	}

	waitFor(t, 2*time.Second, func() bool {
		entry, ok := h.server.FsCache().Get(string(testURI))
		return ok && entry.Content == savedContent
	})

	if entry, ok := h.server.FsCache().Get(string(testURI)); !ok ||
		entry.Version != 2 {
		t.Fatalf("expected version 2 preserved after DidSave, got %+v", entry)
	}

	// 4. Send didClose
	closeParams := &protocol.DidCloseTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: testURI,
		},
	}
	if err := h.client.DidClose(ctx, closeParams); err != nil {
		t.Fatalf("DidClose failed: %v", err)
	}

	waitFor(t, 2*time.Second, func() bool {
		return !h.server.FsCache().HasOverlay(string(testURI))
	})

	_ = h.client.Shutdown(ctx)
	_ = h.client.Exit(ctx)
	_ = <-h.done
}

func TestDocumentSync_WholeDocumentReplacement(t *testing.T) {
	h := setupSyncHarness(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testURI := uri.File("/workspace/whole.dml")

	initParams := &protocol.InitializeParams{
		WorkspaceFolders: protocol.NewNullable([]protocol.WorkspaceFolder{
			{URI: uri.File("/workspace"), Name: "workspace"},
		}),
	}

	if _, err := h.client.Initialize(ctx, initParams); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	if err := h.client.DidOpen(ctx, &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:     testURI,
			Version: 1,
			Text:    "original text\n",
		},
	}); err != nil {
		t.Fatalf("DidOpen failed: %v", err)
	}

	// Send whole document change
	if err := h.client.DidChange(ctx, &protocol.DidChangeTextDocumentParams{
		TextDocument: protocol.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: protocol.TextDocumentIdentifier{
				URI: testURI,
			},
			Version: 2,
		},
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			&protocol.TextDocumentContentChangeWholeDocument{
				Text: "replacement text\n",
			},
		},
	}); err != nil {
		t.Fatalf("DidChange failed: %v", err)
	}

	waitFor(t, 2*time.Second, func() bool {
		data, err := h.server.FsCache().ReadFile(string(testURI))
		return err == nil && string(data) == "replacement text\n"
	})

	_ = h.client.Shutdown(ctx)
	_ = h.client.Exit(ctx)
	_ = <-h.done
}

func TestDocumentSync_DiskFallback(t *testing.T) {
	h := setupSyncHarness(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "disk.dml")
	diskText := "dml 1.4;\n// written to disk\n"
	if err := os.WriteFile(filePath, []byte(diskText), 0o600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	initParams := &protocol.InitializeParams{
		WorkspaceFolders: protocol.NewNullable([]protocol.WorkspaceFolder{
			{URI: uri.File(tmpDir), Name: "tmp"},
		}),
	}

	if _, err := h.client.Initialize(ctx, initParams); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Reading an unopened file through the server's FsCache must fall back to disk
	data, err := h.server.FsCache().ReadFile(filePath)
	if err != nil || string(data) != diskText {
		t.Fatalf(
			"expected disk text %q, got %q, err: %v",
			diskText,
			string(data),
			err,
		)
	}

	_ = h.client.Shutdown(ctx)
	_ = h.client.Exit(ctx)
	_ = <-h.done
}
