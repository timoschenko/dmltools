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
	"sync"
	"testing"
	"time"

	"github.com/timoschenko/dmltools/pkg/lsp"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

type diagCollector struct {
	protocol.UnimplementedClient
	mu           sync.Mutex
	diagnostics  map[string][]protocol.Diagnostic
	versions     map[string]int32
	publishCount map[string]int
}

func newDiagCollector() *diagCollector {
	return &diagCollector{
		diagnostics:  make(map[string][]protocol.Diagnostic),
		versions:     make(map[string]int32),
		publishCount: make(map[string]int),
	}
}

func (c *diagCollector) PublishDiagnostics(
	ctx context.Context,
	params *protocol.PublishDiagnosticsParams,
) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := string(params.URI)
	c.diagnostics[key] = params.Diagnostics
	c.publishCount[key]++
	if v, ok := params.Version.Get(); ok {
		c.versions[key] = v
	} else {
		delete(c.versions, key)
	}
	return nil
}

func (c *diagCollector) getDiagnostics(
	uriStr string,
) ([]protocol.Diagnostic, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	d, ok := c.diagnostics[uriStr]
	return d, ok
}

func (c *diagCollector) getVersion(uriStr string) (int32, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.versions[uriStr]
	return v, ok
}

func (c *diagCollector) getPublishCount(uriStr string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.publishCount[uriStr]
}

type diagHarness struct {
	server    *lsp.Server
	client    protocol.Server
	collector *diagCollector
	done      chan error
}

func setupDiagHarness(t *testing.T, debounce time.Duration) *diagHarness {
	t.Helper()

	clientPipe, serverPipe := net.Pipe()

	logOut := io.Discard
	if testing.Verbose() {
		logOut = os.Stderr
	}
	logger := slog.New(
		slog.NewTextHandler(
			logOut,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		),
	)
	server := lsp.NewServer(
		lsp.WithLogger(logger),
		lsp.WithDiagnosticsDebounce(debounce),
	)

	serverStream := jsonrpc2.NewStream(serverPipe)
	clientStream := jsonrpc2.NewStream(clientPipe)

	ctx := context.Background()
	done := make(chan error, 1)

	go func() {
		done <- server.Serve(ctx, serverStream)
	}()

	collector := newDiagCollector()
	_, _, clientDispatcher := protocol.NewClient(
		ctx,
		collector,
		clientStream,
	)

	return &diagHarness{
		server:    server,
		client:    clientDispatcher,
		collector: collector,
		done:      done,
	}
}

func TestDiagnostics_PublishOnError(t *testing.T) {
	h := setupDiagHarness(t, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testURI := uri.File("/workspace/error.dml")

	initParams := &protocol.InitializeParams{
		WorkspaceFolders: protocol.NewNullable([]protocol.WorkspaceFolder{
			{URI: uri.File("/workspace"), Name: "workspace"},
		}),
	}

	if _, err := h.client.Initialize(ctx, initParams); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Open file with syntax error (missing identifier after device)
	brokenContent := "dml 1.4;\ndevice;\n"
	openParams := &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        testURI,
			LanguageID: "dml",
			Version:    1,
			Text:       brokenContent,
		},
	}

	if err := h.client.DidOpen(ctx, openParams); err != nil {
		t.Fatalf("DidOpen failed: %v", err)
	}

	// Await diagnostics arriving at client
	waitFor(t, 2*time.Second, func() bool {
		diags, ok := h.collector.getDiagnostics(string(testURI))
		return ok && len(diags) > 0
	})

	diags, _ := h.collector.getDiagnostics(string(testURI))
	if len(diags) == 0 {
		t.Fatal("expected at least one diagnostic")
	}

	if diags[0].Severity != protocol.DiagnosticSeverityError {
		t.Errorf("expected DiagnosticSeverityError, got %v", diags[0].Severity)
	}

	_ = h.client.Shutdown(ctx)
	_ = h.client.Exit(ctx)
	_ = <-h.done
}

func TestDiagnostics_ClearOnFix(t *testing.T) {
	h := setupDiagHarness(t, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testURI := uri.File("/workspace/fix.dml")

	initParams := &protocol.InitializeParams{
		WorkspaceFolders: protocol.NewNullable([]protocol.WorkspaceFolder{
			{URI: uri.File("/workspace"), Name: "workspace"},
		}),
	}

	if _, err := h.client.Initialize(ctx, initParams); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// 1. Open with error
	if err := h.client.DidOpen(ctx, &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:     testURI,
			Version: 1,
			Text:    "dml 1.4;\ndevice;\n",
		},
	}); err != nil {
		t.Fatalf("DidOpen failed: %v", err)
	}

	waitFor(t, 2*time.Second, func() bool {
		diags, ok := h.collector.getDiagnostics(string(testURI))
		return ok && len(diags) > 0
	})

	// 2. Fix the error via whole document replacement
	if err := h.client.DidChange(ctx, &protocol.DidChangeTextDocumentParams{
		TextDocument: protocol.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: protocol.TextDocumentIdentifier{
				URI: testURI,
			},
			Version: 2,
		},
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			&protocol.TextDocumentContentChangeWholeDocument{
				Text: "dml 1.4;\ndevice dev;\n",
			},
		},
	}); err != nil {
		t.Fatalf("DidChange failed: %v", err)
	}

	// Await diagnostics being cleared
	waitFor(t, 2*time.Second, func() bool {
		diags, ok := h.collector.getDiagnostics(string(testURI))
		return ok && len(diags) == 0
	})

	// Verify AST is now cached in FsCache
	tree, ok := h.server.FsCache().GetAST(string(testURI))
	if !ok || tree == nil {
		t.Fatal("expected valid AST to be cached in FsCache after fix")
	}

	_ = h.client.Shutdown(ctx)
	_ = h.client.Exit(ctx)
	_ = <-h.done
}

func TestDiagnostics_ClearOnClose(t *testing.T) {
	h := setupDiagHarness(t, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testURI := uri.File("/workspace/close.dml")

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
			Text:    "dml 1.4;\ndevice;\n",
		},
	}); err != nil {
		t.Fatalf("DidOpen failed: %v", err)
	}

	waitFor(t, 2*time.Second, func() bool {
		diags, ok := h.collector.getDiagnostics(string(testURI))
		return ok && len(diags) > 0
	})

	// Close the file
	if err := h.client.DidClose(ctx, &protocol.DidCloseTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: testURI,
		},
	}); err != nil {
		t.Fatalf("DidClose failed: %v", err)
	}

	// Await empty diagnostics sent to clear editor markers
	waitFor(t, 2*time.Second, func() bool {
		diags, ok := h.collector.getDiagnostics(string(testURI))
		return ok && len(diags) == 0
	})

	_ = h.client.Shutdown(ctx)
	_ = h.client.Exit(ctx)
	_ = <-h.done
}

func TestDiagnostics_Debounce(t *testing.T) {
	h := setupDiagHarness(t, 50*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testURI := uri.File("/workspace/debounce.dml")

	initParams := &protocol.InitializeParams{
		WorkspaceFolders: protocol.NewNullable([]protocol.WorkspaceFolder{
			{URI: uri.File("/workspace"), Name: "workspace"},
		}),
	}

	if _, err := h.client.Initialize(ctx, initParams); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// 1. Open with valid code
	if err := h.client.DidOpen(ctx, &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:     testURI,
			Version: 1,
			Text:    "dml 1.4;\ndevice dev;\n",
		},
	}); err != nil {
		t.Fatalf("DidOpen failed: %v", err)
	}

	waitFor(t, 2*time.Second, func() bool {
		diags, ok := h.collector.getDiagnostics(string(testURI))
		return ok && len(diags) == 0
	})

	// 2. Send two rapid changes
	_ = h.client.DidChange(ctx, &protocol.DidChangeTextDocumentParams{
		TextDocument: protocol.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: protocol.TextDocumentIdentifier{
				URI: testURI,
			},
			Version: 2,
		},
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			&protocol.TextDocumentContentChangeWholeDocument{
				Text: "dml 1.4;\ndevice;\n",
			},
		},
	})

	_ = h.client.DidChange(ctx, &protocol.DidChangeTextDocumentParams{
		TextDocument: protocol.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: protocol.TextDocumentIdentifier{
				URI: testURI,
			},
			Version: 3,
		},
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			&protocol.TextDocumentContentChangeWholeDocument{
				Text: "dml 1.4;\ndevice final_dev;\n",
			},
		},
	})

	// Debounce timer fires after 50ms, final state is valid DML
	waitFor(t, 2*time.Second, func() bool {
		diags, ok := h.collector.getDiagnostics(string(testURI))
		return ok && len(diags) == 0
	})

	_ = h.client.Shutdown(ctx)
	_ = h.client.Exit(ctx)
	_ = <-h.done
}

func TestDiagnostics_DidSavePreservesVersion(t *testing.T) {
	h := setupDiagHarness(t, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testURI := uri.File("/workspace/save_ver.dml")

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
			Text:    "dml 1.4;\ndevice dev;\n",
		},
	}); err != nil {
		t.Fatalf("DidOpen failed: %v", err)
	}
	waitFor(t, 2*time.Second, func() bool {
		return h.collector.getPublishCount(string(testURI)) == 1
	})

	if err := h.client.DidChange(ctx, &protocol.DidChangeTextDocumentParams{
		TextDocument: protocol.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: protocol.TextDocumentIdentifier{
				URI: testURI,
			},
			Version: 2,
		},
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			&protocol.TextDocumentContentChangeWholeDocument{
				Text: "dml 1.4;\ndevice dev2;\n",
			},
		},
	}); err != nil {
		t.Fatalf("DidChange failed: %v", err)
	}
	waitFor(t, 2*time.Second, func() bool {
		return h.collector.getPublishCount(string(testURI)) == 2
	})

	savedText := "dml 1.4;\ndevice;\n"
	if err := h.client.DidSave(ctx, &protocol.DidSaveTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: testURI},
		Text:         &savedText,
	}); err != nil {
		t.Fatalf("DidSave failed: %v", err)
	}

	waitFor(t, 2*time.Second, func() bool {
		return h.collector.getPublishCount(string(testURI)) == 3
	})

	entry, ok := h.server.FsCache().Get(string(testURI))
	if !ok || entry.Version != 2 {
		t.Fatalf("expected FsCache version 2 after DidSave, got %+v", entry)
	}
	ver, hasVer := h.collector.getVersion(string(testURI))
	if !hasVer || ver != 2 {
		t.Fatalf(
			"expected published diagnostics version 2, got %d (present=%v)",
			ver,
			hasVer,
		)
	}

	_ = h.client.Shutdown(ctx)
	_ = h.client.Exit(ctx)
	_ = <-h.done
}

func TestDiagnostics_NoPublishAfterDidClose(t *testing.T) {
	h := setupDiagHarness(t, 30*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tmpDir := t.TempDir()
	diskPath := filepath.Join(tmpDir, "closed.dml")
	if err := os.WriteFile(
		diskPath,
		[]byte("dml 1.4;\ndevice;\n"),
		0o600,
	); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	testURI := uri.File(diskPath)

	initParams := &protocol.InitializeParams{
		WorkspaceFolders: protocol.NewNullable([]protocol.WorkspaceFolder{
			{URI: uri.File(tmpDir), Name: "tmp"},
		}),
	}
	if _, err := h.client.Initialize(ctx, initParams); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	if err := h.client.DidOpen(ctx, &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:     testURI,
			Version: 1,
			Text:    "dml 1.4;\ndevice dev;\n",
		},
	}); err != nil {
		t.Fatalf("DidOpen failed: %v", err)
	}
	waitFor(t, 2*time.Second, func() bool {
		return h.collector.getPublishCount(string(testURI)) == 1
	})

	if err := h.client.DidChange(ctx, &protocol.DidChangeTextDocumentParams{
		TextDocument: protocol.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: protocol.TextDocumentIdentifier{
				URI: testURI,
			},
			Version: 2,
		},
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			&protocol.TextDocumentContentChangeWholeDocument{
				Text: "dml 1.4;\ndevice;\n",
			},
		},
	}); err != nil {
		t.Fatalf("DidChange failed: %v", err)
	}

	if err := h.client.DidClose(ctx, &protocol.DidCloseTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: testURI},
	}); err != nil {
		t.Fatalf("DidClose failed: %v", err)
	}

	waitFor(t, 2*time.Second, func() bool {
		return h.collector.getPublishCount(string(testURI)) == 2
	})

	time.Sleep(80 * time.Millisecond)

	if count := h.collector.getPublishCount(string(testURI)); count != 2 {
		t.Fatalf(
			"expected no diagnostics published after DidClose, got %d publishes",
			count,
		)
	}
	if diags, _ := h.collector.getDiagnostics(
		string(testURI),
	); len(diags) != 0 {
		t.Fatalf("expected cleared diagnostics after DidClose, got %+v", diags)
	}

	_ = h.client.Shutdown(ctx)
	_ = h.client.Exit(ctx)
	_ = <-h.done
}

func TestDiagnose_Valid(t *testing.T) {
	code := "dml 1.4;\ndevice dev;\n"
	diags, tree := lsp.Diagnose(code)

	if len(diags) != 0 {
		t.Fatalf(
			"expected 0 diagnostics for valid code, got %d: %+v",
			len(diags),
			diags,
		)
	}

	if tree == nil {
		t.Fatal("expected non-nil AST for valid code")
	}
}

func TestDiagnose_SyntaxError(t *testing.T) {
	// Missing semicolon on line 1 (0-indexed: line 1)
	code := "dml 1.4;\ndevice dev\n"
	diags, _ := lsp.Diagnose(code)

	if len(diags) == 0 {
		t.Fatal("expected diagnostics for missing semicolon, got 0")
	}

	found := false
	for _, d := range diags {
		if d.Severity == protocol.DiagnosticSeverityError {
			found = true
			if d.Range.Start.Line < 1 {
				t.Errorf(
					"expected error around line 1, got line %d",
					d.Range.Start.Line,
				)
			}
		}
	}
	if !found {
		t.Fatal("expected at least one error severity diagnostic")
	}
}

func TestDiagnose_InvalidToken(t *testing.T) {
	code := "dml 1.4;\n@@@;\n"
	diags, _ := lsp.Diagnose(code)

	if len(diags) == 0 {
		t.Fatal("expected diagnostics for invalid token '@', got 0")
	}
}

func TestDiagnose_UnimplementedVersionDml12(t *testing.T) {
	code := "dml 1.2;\ndevice dev;\n"
	diags, _ := lsp.Diagnose(code)

	if len(diags) != 1 {
		t.Fatalf(
			"expected 1 diagnostic for dml 1.2, got %d: %+v",
			len(diags),
			diags,
		)
	}

	if diags[0].Severity != protocol.DiagnosticSeverityWarning {
		t.Errorf(
			"expected DiagnosticSeverityWarning, got %v",
			diags[0].Severity,
		)
	}

	if diags[0].Code != protocol.String("EUNIMPL") {
		t.Errorf(
			"expected diagnostic code 'EUNIMPL', got %v",
			diags[0].Code,
		)
	}
}

func TestDiagnose_UTF16Columns(t *testing.T) {
	// Line 1: `/* 🔥🔥 */ "😀";`
	// `/* ` (3) + `🔥🔥` (4 UTF-16) + ` */ ` (4) = 11 UTF-16 code units.
	// `"😀"` = 4 UTF-16 code units (`"` + surrogate pair + `"`).
	code := "dml 1.4;\n/* 🔥🔥 */ \"😀\";\n"
	diags, _ := lsp.Diagnose(code)

	if len(diags) == 0 {
		t.Fatal("expected syntax error diagnostic, got 0")
	}

	d := diags[0]
	if d.Range.Start.Line != 1 {
		t.Fatalf("expected error on line 1, got %d", d.Range.Start.Line)
	}
	if d.Range.Start.Character != 11 || d.Range.End.Character != 15 {
		t.Fatalf(
			"expected UTF-16 range [11, 15), got [%d, %d)",
			d.Range.Start.Character,
			d.Range.End.Character,
		)
	}
}
