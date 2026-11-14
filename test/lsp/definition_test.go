// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lsp_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"

	"github.com/timoschenko/dmltools/pkg/lsp"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

type defHarness struct {
	server *lsp.Server
	client protocol.Server
	done   chan error
}

func setupDefHarness(t *testing.T) *defHarness {
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

	return &defHarness{
		server: server,
		client: clientDispatcher,
		done:   done,
	}
}

func TestDefinition_ImportNavigation(t *testing.T) {
	h := setupDefHarness(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wsURI := uri.File("/workspace")
	mainURI := uri.File("/workspace/main.dml")
	subURI := uri.File("/workspace/sub.dml")

	initParams := &protocol.InitializeParams{
		ClientInfo: protocol.ClientInfo{
			Name:    "test-client",
			Version: protocol.NewOptional("1.0.0"),
		},
		WorkspaceFolders: protocol.NewNullable([]protocol.WorkspaceFolder{
			{URI: wsURI, Name: "workspace"},
		}),
	}

	res, err := h.client.Initialize(ctx, initParams)
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	if res.Capabilities.DefinitionProvider == nil {
		t.Fatal("expected DefinitionProvider capability to be advertised")
	}

	if err := h.client.Initialized(
		ctx,
		&protocol.InitializedParams{},
	); err != nil {
		t.Fatalf("Initialized notification failed: %v", err)
	}

	mainContent := `dml 1.4;
device main_dev;
import "sub.dml";
import "missing.dml";

bank b {
	is imported_template;
}
`
	subContent := `dml 1.4;
template imported_template {
	param p = 42;
}
`

	err = h.client.DidOpen(ctx, &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        mainURI,
			LanguageID: "dml",
			Version:    1,
			Text:       mainContent,
		},
	})
	if err != nil {
		t.Fatalf("DidOpen main failed: %v", err)
	}

	err = h.client.DidOpen(ctx, &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        subURI,
			LanguageID: "dml",
			Version:    1,
			Text:       subContent,
		},
	})
	if err != nil {
		t.Fatalf("DidOpen sub failed: %v", err)
	}

	waitFor(t, 2*time.Second, func() bool {
		return h.server.FsCache().HasOverlay(string(mainURI)) &&
			h.server.FsCache().HasOverlay(string(subURI))
	})

	// 1. Definition on import "sub.dml"; -> line 3 (0-based: 2)
	defResult, err := h.client.Definition(ctx, &protocol.DefinitionParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: mainURI},
		Position:     protocol.Position{Line: 2, Character: 8},
	})
	if err != nil {
		t.Fatalf("Definition request failed: %v", err)
	}
	loc, ok := defResult.(*protocol.Location)
	if !ok || loc == nil {
		t.Fatalf(
			"expected *protocol.Location, got %T (%+v)",
			defResult,
			defResult,
		)
	}
	if loc.URI != subURI {
		t.Fatalf("expected target URI %s, got %s", subURI, loc.URI)
	}
	if loc.Range.Start.Line != 0 || loc.Range.Start.Character != 0 {
		t.Fatalf("expected start at 0:0, got %+v", loc.Range.Start)
	}

	// 2. Definition on missing import "missing.dml"; -> line 4 (0-based: 3)
	defMissing, err := h.client.Definition(ctx, &protocol.DefinitionParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: mainURI},
		Position:     protocol.Position{Line: 3, Character: 8},
	})
	if err != nil {
		t.Fatalf("Definition request on missing import failed: %v", err)
	}
	if defMissing != nil {
		t.Fatalf("expected nil for missing import, got %+v", defMissing)
	}

	// 3. Cross-file template navigation: line 7 (0-based: 6)
	defCross, err := h.client.Definition(ctx, &protocol.DefinitionParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: mainURI},
		Position:     protocol.Position{Line: 6, Character: 8},
	})
	if err != nil {
		t.Fatalf("Definition cross-file request failed: %v", err)
	}
	locCross, ok := defCross.(*protocol.Location)
	if !ok || locCross == nil {
		t.Fatalf("expected cross-file Location, got %+v", defCross)
	}
	if locCross.URI != subURI {
		t.Fatalf("expected target URI %s, got %s", subURI, locCross.URI)
	}
	if locCross.Range.Start.Line != 1 {
		t.Fatalf(
			"expected line 1 in sub.dml, got %d",
			locCross.Range.Start.Line,
		)
	}

	if err := h.client.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}
	if err := h.client.Exit(ctx); err != nil {
		t.Fatalf("Exit failed: %v", err)
	}
	select {
	case err := <-h.done:
		if err != nil && !errors.Is(err, io.EOF) {
			t.Logf("server terminated: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("server did not terminate after exit")
	}
}

func TestDefinition_SymbolsAndHierarchy(t *testing.T) {
	h := setupDefHarness(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testURI := uri.File("/workspace/symbols.dml")
	content := `dml 1.4;
device sym_dev;

template base_templ {
	param p = 100;
}

bank regs {
	is base_templ;
	register r1 {
		field f1;
		method test() {
			this.parent.r2;
		}
	}
	register r2 {
		field f2;
	}
}

method top() {
	regs.r1.f1;
}
`

	initParams := &protocol.InitializeParams{
		ClientInfo: protocol.ClientInfo{Name: "test-client"},
	}
	if _, err := h.client.Initialize(ctx, initParams); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	if err := h.client.Initialized(
		ctx,
		&protocol.InitializedParams{},
	); err != nil {
		t.Fatalf("Initialized failed: %v", err)
	}

	err := h.client.DidOpen(ctx, &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        testURI,
			LanguageID: "dml",
			Version:    1,
			Text:       content,
		},
	})
	if err != nil {
		t.Fatalf("DidOpen failed: %v", err)
	}

	waitFor(t, 2*time.Second, func() bool {
		return h.server.FsCache().HasOverlay(string(testURI))
	})

	// 1. Template reference: line 9 (0-based: 8): `\tis base_templ;`
	defTempl, err := h.client.Definition(ctx, &protocol.DefinitionParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: testURI},
		Position:     protocol.Position{Line: 8, Character: 6},
	})
	if err != nil {
		t.Fatalf("Definition failed: %v", err)
	}
	locTempl, ok := defTempl.(*protocol.Location)
	if !ok || locTempl == nil {
		t.Fatalf("expected template Location, got %+v", defTempl)
	}
	// `template base_templ` is line 4 (0-based: 3)
	if locTempl.Range.Start.Line != 3 {
		t.Fatalf("expected line 3, got %d", locTempl.Range.Start.Line)
	}

	// 2. Multi-level member chain: line 22 (0-based: 21): `\tregs.r1.f1;`
	defChain, err := h.client.Definition(ctx, &protocol.DefinitionParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: testURI},
		Position:     protocol.Position{Line: 21, Character: 10},
	})
	if err != nil {
		t.Fatalf("Definition failed: %v", err)
	}
	locChain, ok := defChain.(*protocol.Location)
	if !ok || locChain == nil {
		t.Fatalf("expected chain Location, got %+v", defChain)
	}
	// `field f1;` is line 11 (0-based: 10)
	if locChain.Range.Start.Line != 10 {
		t.Fatalf("expected line 10 for f1, got %d", locChain.Range.Start.Line)
	}

	// 3. Keyword navigation: line 13 (0-based: 12): `\t\t\tthis.parent.r2;`
	defKeyword, err := h.client.Definition(ctx, &protocol.DefinitionParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: testURI},
		Position:     protocol.Position{Line: 12, Character: 16},
	})
	if err != nil {
		t.Fatalf("Definition failed: %v", err)
	}
	locKeyword, ok := defKeyword.(*protocol.Location)
	if !ok || locKeyword == nil {
		t.Fatalf("expected keyword Location, got %+v", defKeyword)
	}
	// `register r2` is line 16 (0-based: 15)
	if locKeyword.Range.Start.Line != 15 {
		t.Fatalf("expected line 15 for r2, got %d", locKeyword.Range.Start.Line)
	}

	if err := h.client.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}
	if err := h.client.Exit(ctx); err != nil {
		t.Fatalf("Exit failed: %v", err)
	}
	select {
	case err := <-h.done:
		if err != nil && !errors.Is(err, io.EOF) {
			t.Logf("server terminated: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("server did not terminate after exit")
	}
}

func TestResolveDefinition_Template(t *testing.T) {
	cache := lsp.NewFsCache()
	uri := "file:///test/main.dml"
	content := `dml 1.4;
device test_dev;

template my_template {
	param val = 1;
}

bank b {
	is my_template;
}
`
	if err := cache.Open(uri, 1, content); err != nil {
		t.Fatalf("failed to open file in cache: %v", err)
	}

	// Line 9 (0-based: 8): `\tis my_template;`
	// Character 4 is on 'm' of 'my_template'
	pos := protocol.Position{Line: 8, Character: 6}
	ctx := context.Background()

	loc, err := lsp.ResolveDefinition(ctx, cache, uri, "", pos)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc == nil {
		t.Fatal("expected location, got nil")
	}

	// Line 4 (0-based: 3) is `template my_template {`
	if loc.Range.Start.Line != 3 {
		t.Fatalf("expected definition on line 3, got %d", loc.Range.Start.Line)
	}
}

func TestResolveDefinition_MethodAndParam(t *testing.T) {
	cache := lsp.NewFsCache()
	uri := "file:///test/methods.dml"
	content := `dml 1.4;
device test_dev;

method compute(int count) -> (int) {
	local int step = 1;
	return count + step;
}

bank b {
	method run() {
		compute(5);
	}
}
`
	if err := cache.Open(uri, 1, content); err != nil {
		t.Fatalf("failed to open file in cache: %v", err)
	}
	ctx := context.Background()

	// 1. Call to compute: line 11 (0-based: 10): `\t\tcompute(5);`
	posMethod := protocol.Position{Line: 10, Character: 4}
	locMethod, err := lsp.ResolveDefinition(ctx, cache, uri, "", posMethod)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if locMethod == nil {
		t.Fatal("expected location for method, got nil")
	}
	if locMethod.Range.Start.Line != 3 {
		t.Fatalf(
			"expected method on line 3, got %d",
			locMethod.Range.Start.Line,
		)
	}

	// 2. Local var inside compute: line 6 (0-based: 5): `\treturn count + step;`
	// Col 16 is on 'step'
	posLocal := protocol.Position{Line: 5, Character: 16}
	locLocal, err := lsp.ResolveDefinition(ctx, cache, uri, "", posLocal)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if locLocal == nil {
		t.Fatal("expected location for local var step, got nil")
	}
	// `local int step = 1;` is on line 5 (0-based: 4)
	if locLocal.Range.Start.Line != 4 {
		t.Fatalf("expected step on line 4, got %d", locLocal.Range.Start.Line)
	}

	// 3. Method parameter inside compute: line 6 (0-based: 5):
	// `\treturn count + step;` - Col 9 is on 'count'
	posParam := protocol.Position{Line: 5, Character: 9}
	locParam, err := lsp.ResolveDefinition(ctx, cache, uri, "", posParam)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if locParam == nil {
		t.Fatal("expected location for param count, got nil")
	}
	if locParam.Range.Start.Line != 3 {
		t.Fatalf("expected count on line 3, got %d", locParam.Range.Start.Line)
	}
}

func TestResolveDefinition_ObjectHierarchyAndKeywords(t *testing.T) {
	cache := lsp.NewFsCache()
	uri := "file:///test/hierarchy.dml"
	content := `dml 1.4;
device test_dev;

bank regs {
	register r1 {
		field f1;
		method test_access() {
			this.parent.r2;
		}
	}
	register r2 {
		field f2;
	}
}

method top() {
	regs.r1.f1;
}
`
	if err := cache.Open(uri, 1, content); err != nil {
		t.Fatalf("failed to open file in cache: %v", err)
	}
	ctx := context.Background()

	// 1. Dotted chain from top(): line 17 (0-based: 16): `\tregs.r1.f1;`
	// Col 10 is on 'f1'
	posF1 := protocol.Position{Line: 16, Character: 10}
	locF1, err := lsp.ResolveDefinition(ctx, cache, uri, "", posF1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if locF1 == nil {
		t.Fatal("expected location for f1, got nil")
	}
	// `field f1;` is line 6 (0-based: 5)
	if locF1.Range.Start.Line != 5 {
		t.Fatalf("expected field f1 on line 5, got %d", locF1.Range.Start.Line)
	}

	// 2. Hierarchy keywords: line 8 (0-based: 7): `\t\t\tthis.parent.r2;`
	// Col 16 is on 'r2'
	posR2 := protocol.Position{Line: 7, Character: 16}
	locR2, err := lsp.ResolveDefinition(ctx, cache, uri, "", posR2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if locR2 == nil {
		t.Fatal("expected location for r2 via this.parent.r2, got nil")
	}
	// `register r2` is line 11 (0-based: 10)
	if locR2.Range.Start.Line != 10 {
		t.Fatalf("expected r2 on line 10, got %d", locR2.Range.Start.Line)
	}
}

func TestResolveDefinition_RelativeImport_FoundAndMissing(t *testing.T) {
	cache := lsp.NewFsCache()
	mainURI := "file:///workspace/main.dml"
	subURI := "file:///workspace/sub.dml"

	mainContent := `dml 1.4;
device dev;
import "sub.dml";
import "missing.dml";

bank b {
	is imported_template;
}
`
	subContent := `dml 1.4;
template imported_template {
	param p = 42;
}
`
	if err := cache.Open(mainURI, 1, mainContent); err != nil {
		t.Fatalf("failed to open main file: %v", err)
	}
	if err := cache.Open(subURI, 1, subContent); err != nil {
		t.Fatalf("failed to open sub file: %v", err)
	}
	ctx := context.Background()

	// 1. Cursor on `import "sub.dml";` (line 3, 0-based: 2)
	posImport := protocol.Position{Line: 2, Character: 8}
	locImport, err := lsp.ResolveDefinition(
		ctx,
		cache,
		mainURI,
		"/workspace",
		posImport,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if locImport == nil {
		t.Fatal("expected location for found import, got nil")
	}
	if string(locImport.URI) != subURI {
		t.Fatalf("expected URI %s, got %s", subURI, string(locImport.URI))
	}
	if locImport.Range.Start.Line != 0 || locImport.Range.Start.Character != 0 {
		t.Fatalf(
			"expected jump to line 0 col 0, got %+v",
			locImport.Range.Start,
		)
	}

	// 2. Cursor on missing import: line 4 (0-based: 3)
	posMissing := protocol.Position{Line: 3, Character: 8}
	locMissing, err := lsp.ResolveDefinition(
		ctx,
		cache,
		mainURI,
		"/workspace",
		posMissing,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if locMissing != nil {
		t.Fatalf("expected nil for missing import, got %+v", locMissing)
	}

	// 3. Cursor on cross-file symbol `imported_template`: line 7 (0-based: 6)
	posCrossSym := protocol.Position{Line: 6, Character: 8}
	locCrossSym, err := lsp.ResolveDefinition(
		ctx,
		cache,
		mainURI,
		"/workspace",
		posCrossSym,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if locCrossSym == nil {
		t.Fatal("expected location for cross-file symbol, got nil")
	}
	if string(locCrossSym.URI) != subURI {
		t.Fatalf("expected URI %s, got %s", subURI, string(locCrossSym.URI))
	}
	// `template imported_template` is on line 2 (0-based: 1)
	if locCrossSym.Range.Start.Line != 1 {
		t.Fatalf(
			"expected symbol on line 1 of sub.dml, got %d",
			locCrossSym.Range.Start.Line,
		)
	}
}

func TestResolveDefinition_TransitiveAndCyclicImports(t *testing.T) {
	cache := lsp.NewFsCache()
	mainURI := "file:///workspace/main.dml"
	midURI := "file:///workspace/mid.dml"
	leafURI := "file:///workspace/leaf.dml"

	mainContent := `dml 1.4;
device dev;
import "mid.dml";

method use_leaf() -> (int) {
	return LEAF_CONST;
}

method use_missing() -> (int) {
	return MISSING_SYM;
}
`
	midContent := `dml 1.4;
import "leaf.dml";
`
	leafContent := `dml 1.4;
import "main.dml";
constant LEAF_CONST = 99;
`

	if err := cache.Open(mainURI, 1, mainContent); err != nil {
		t.Fatalf("failed to open main file: %v", err)
	}
	if err := cache.Open(midURI, 1, midContent); err != nil {
		t.Fatalf("failed to open mid file: %v", err)
	}
	if err := cache.Open(leafURI, 1, leafContent); err != nil {
		t.Fatalf("failed to open leaf file: %v", err)
	}
	ctx := context.Background()

	// Line 6 (0-based: 5): `\treturn LEAF_CONST;`
	loc, err := lsp.ResolveDefinition(
		ctx,
		cache,
		mainURI,
		"/workspace",
		protocol.Position{Line: 5, Character: 10},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc == nil {
		t.Fatal("expected location for transitive import symbol, got nil")
	}
	if string(loc.URI) != leafURI {
		t.Fatalf("expected URI %s, got %s", leafURI, string(loc.URI))
	}
	if loc.Range.Start.Line != 2 {
		t.Fatalf("expected constant on line 2, got %d", loc.Range.Start.Line)
	}

	// Ensure cyclic imports terminate cleanly for missing symbols.
	missingLoc, err := lsp.ResolveDefinition(
		ctx,
		cache,
		mainURI,
		"/workspace",
		protocol.Position{Line: 9, Character: 10},
	)
	if err != nil {
		t.Fatalf("unexpected error for missing symbol: %v", err)
	}
	if missingLoc != nil {
		t.Fatalf("expected nil for missing symbol, got %+v", missingLoc)
	}
}

func TestResolveDefinition_VersionedDirectoryImport(t *testing.T) {
	cache := lsp.NewFsCache()
	mainURI := "file:///workspace/dev.dml"
	stdURI := "file:///workspace/1.4/std.dml"

	mainContent := `dml 1.4;
device dev;
import "std.dml";

bank b {
	is std_tmpl;
}
`
	stdContent := `dml 1.4;
template std_tmpl {
	param ok = 1;
}
`

	if err := cache.Open(mainURI, 1, mainContent); err != nil {
		t.Fatalf("failed to open main file: %v", err)
	}
	if err := cache.Open(stdURI, 1, stdContent); err != nil {
		t.Fatalf("failed to open std file: %v", err)
	}
	ctx := context.Background()

	// Cursor on `import "std.dml";` (line 3, 0-based: 2)
	locImp, err := lsp.ResolveDefinition(
		ctx,
		cache,
		mainURI,
		"/workspace",
		protocol.Position{Line: 2, Character: 8},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if locImp == nil || string(locImp.URI) != stdURI {
		t.Fatalf("expected import to resolve to %s, got %+v", stdURI, locImp)
	}

	// Cursor on `std_tmpl` (line 6, 0-based: 5)
	locTmpl, err := lsp.ResolveDefinition(
		ctx,
		cache,
		mainURI,
		"/workspace",
		protocol.Position{Line: 5, Character: 6},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if locTmpl == nil || string(locTmpl.URI) != stdURI {
		t.Fatalf("expected symbol in %s, got %+v", stdURI, locTmpl)
	}
	if locTmpl.Range.Start.Line != 1 {
		t.Fatalf(
			"expected template on line 1, got %d",
			locTmpl.Range.Start.Line,
		)
	}
}

func TestResolveDefinition_TemplateMemberInheritance(t *testing.T) {
	cache := lsp.NewFsCache()
	mainURI := "file:///workspace/main.dml"
	tmplURI := "file:///workspace/tmpl.dml"

	tmplContent := `dml 1.4;
template root_tmpl {
	param root_param = 1;
}

template base_tmpl is root_tmpl {
	param base_param = 2;
}
`
	mainContent := `dml 1.4;
device dev;
import "tmpl.dml";

bank b is base_tmpl {
	method m() {
		this.base_param;
		root_param;
	}
}

method top() {
	b.base_param;
}
`

	if err := cache.Open(tmplURI, 1, tmplContent); err != nil {
		t.Fatalf("failed to open tmpl file: %v", err)
	}
	if err := cache.Open(mainURI, 1, mainContent); err != nil {
		t.Fatalf("failed to open main file: %v", err)
	}
	ctx := context.Background()

	// 1. `this.base_param;` -> line 6, col 8
	locThis, err := lsp.ResolveDefinition(
		ctx,
		cache,
		mainURI,
		"/workspace",
		protocol.Position{Line: 6, Character: 8},
	)
	if err != nil || locThis == nil {
		t.Fatalf("expected this.base_param to resolve, err=%v", err)
	}
	if string(locThis.URI) != tmplURI || locThis.Range.Start.Line != 6 {
		t.Fatalf("expected %s:6, got %+v", tmplURI, locThis)
	}

	// 2. Unqualified `root_param;` from inherited parent template -> line 7
	locRoot, err := lsp.ResolveDefinition(
		ctx,
		cache,
		mainURI,
		"/workspace",
		protocol.Position{Line: 7, Character: 4},
	)
	if err != nil || locRoot == nil {
		t.Fatalf("expected root_param to resolve, err=%v", err)
	}
	if string(locRoot.URI) != tmplURI || locRoot.Range.Start.Line != 2 {
		t.Fatalf("expected %s:2, got %+v", tmplURI, locRoot)
	}

	// 3. `b.base_param;` -> line 12, col 4
	locDot, err := lsp.ResolveDefinition(
		ctx,
		cache,
		mainURI,
		"/workspace",
		protocol.Position{Line: 12, Character: 4},
	)
	if err != nil || locDot == nil {
		t.Fatalf("expected b.base_param to resolve, err=%v", err)
	}
	if string(locDot.URI) != tmplURI || locDot.Range.Start.Line != 6 {
		t.Fatalf("expected %s:6, got %+v", tmplURI, locDot)
	}
}

func TestResolveDefinition_DmlDeclarationsAndOutParam(t *testing.T) {
	cache := lsp.NewFsCache()
	uri := "file:///workspace/decls.dml"
	content := `dml 1.4;
device dev;

extern int ext_fn(int x);
loggroup my_log;
session int sess_var = 0;
saved int saved_var = 1;
hook(int) my_hook;

method compute(int in_v) -> (int out_v) {
	out_v = in_v + sess_var + saved_var;
	ext_fn(out_v);
	log info, 1, my_log: "hi";
	my_hook;
}
`
	if err := cache.Open(uri, 1, content); err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	ctx := context.Background()

	cases := []struct {
		name     string
		pos      protocol.Position
		wantLine uint32
	}{
		{"out_v", protocol.Position{Line: 10, Character: 2}, 9},
		{"sess_var", protocol.Position{Line: 10, Character: 18}, 5},
		{"saved_var", protocol.Position{Line: 10, Character: 29}, 6},
		{"ext_fn", protocol.Position{Line: 11, Character: 2}, 3},
		{"my_log", protocol.Position{Line: 12, Character: 15}, 4},
		{"my_hook", protocol.Position{Line: 13, Character: 3}, 7},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			loc, err := lsp.ResolveDefinition(ctx, cache, uri, "", tc.pos)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if loc == nil {
				t.Fatalf("expected definition for %s, got nil", tc.name)
			}
			if loc.Range.Start.Line != tc.wantLine {
				t.Fatalf(
					"%s: expected line %d, got %d",
					tc.name,
					tc.wantLine,
					loc.Range.Start.Line,
				)
			}
		})
	}
}

func TestResolveDefinition_Utf16EmojiPosition(t *testing.T) {
	cache := lsp.NewFsCache()
	uri := "file:///workspace/utf16.dml"
	content := "dml 1.4;\n" +
		"device dev;\n\n" +
		"bank b {\n" +
		"\t/* 🔥🔥 */ param emoji_param = 42;\n" +
		"}\n\n" +
		"method test() {\n" +
		"\t/* 🔥🔥 */ b.emoji_param;\n" +
		"}\n"

	if err := cache.Open(uri, 1, content); err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	// UTF-16 col 14 points at 'e' in 'emoji_param' after `\t/* 🔥🔥 */ b.`.
	loc, err := lsp.ResolveDefinition(
		context.Background(),
		cache,
		uri,
		"",
		protocol.Position{Line: 8, Character: 14},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc == nil {
		t.Fatal("expected definition after UTF-16 surrogate pair, got nil")
	}
	if loc.Range.Start.Line != 4 {
		t.Fatalf("expected param on line 4, got %d", loc.Range.Start.Line)
	}
	// `\t` (1) + `/* ` (3) + `🔥🔥` (4) + ` */ ` (4) + `param ` (6) = 18.
	// `emoji_param` length = 11 -> [18, 29).
	if loc.Range.Start.Character != 18 || loc.Range.End.Character != 29 {
		t.Fatalf(
			"expected UTF-16 NameRange [18, 29), got [%d, %d)",
			loc.Range.Start.Character,
			loc.Range.End.Character,
		)
	}
}

func TestResolveDefinition_TopLevelBankParentAndDevScope(t *testing.T) {
	cache := lsp.NewFsCache()
	uri := "file:///workspace/top_parent.dml"
	content := `dml 1.4;
device test_dev;

template my_tmpl {
	param tp = 1;
}

bank regs {
	register r1 {
		method nested_m() {
			this.parent.parent;
		}
	}
	method m() {
		this.parent;
		parent;
	}
}

method top() {
	dev.regs;
	dev.my_tmpl;
}
`
	if err := cache.Open(uri, 1, content); err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	ctx := context.Background()

	// 1. `this.parent.parent;` from register r1 -> device test_dev (line 1)
	locNested, err := lsp.ResolveDefinition(
		ctx,
		cache,
		uri,
		"",
		protocol.Position{Line: 10, Character: 15},
	)
	if err != nil || locNested == nil || locNested.Range.Start.Line != 1 {
		t.Fatalf(
			"expected this.parent.parent to resolve to line 1, got %+v (err=%v)",
			locNested,
			err,
		)
	}

	// 2. `this.parent;` from top-level bank regs -> device test_dev (line 1)
	locThisParent, err := lsp.ResolveDefinition(
		ctx,
		cache,
		uri,
		"",
		protocol.Position{Line: 14, Character: 7},
	)
	if err != nil || locThisParent == nil ||
		locThisParent.Range.Start.Line != 1 {
		t.Fatalf(
			"expected this.parent to resolve to line 1, got %+v (err=%v)",
			locThisParent,
			err,
		)
	}

	// 3. `parent;` from top-level bank regs -> device test_dev (line 1)
	locParent, err := lsp.ResolveDefinition(
		ctx,
		cache,
		uri,
		"",
		protocol.Position{Line: 15, Character: 2},
	)
	if err != nil || locParent == nil || locParent.Range.Start.Line != 1 {
		t.Fatalf(
			"expected parent to resolve to line 1, got %+v (err=%v)",
			locParent,
			err,
		)
	}

	// 4. `dev.regs;` -> bank regs (line 7)
	locDevRegs, err := lsp.ResolveDefinition(
		ctx,
		cache,
		uri,
		"",
		protocol.Position{Line: 20, Character: 5},
	)
	if err != nil || locDevRegs == nil || locDevRegs.Range.Start.Line != 7 {
		t.Fatalf(
			"expected dev.regs to resolve to line 7, got %+v (err=%v)",
			locDevRegs,
			err,
		)
	}

	// 5. `dev.my_tmpl;` must NOT leak top-level templates onto dev
	locDevTmpl, err := lsp.ResolveDefinition(
		ctx,
		cache,
		uri,
		"",
		protocol.Position{Line: 21, Character: 5},
	)
	if err != nil {
		t.Fatalf("unexpected error for dev.my_tmpl: %v", err)
	}
	if locDevTmpl != nil {
		t.Fatalf(
			"expected nil for dev.my_tmpl, got %+v",
			locDevTmpl,
		)
	}
}

func TestResolveDefinition_InEachIsolation(t *testing.T) {
	cache := lsp.NewFsCache()
	uri := "file:///workspace/ineach.dml"
	content := `dml 1.4;
device test_dev;

template inner_tmpl {
	param tmpl_p = 1;
}

in each bank {
	is inner_tmpl;
	param leaked_param = 2;
}

bank b {
	in each register {
		param bank_inner_p = 3;
	}
	method m() {
		bank_inner_p;
		this.bank_inner_p;
	}
}

method top() {
	tmpl_p;
	leaked_param;
	dev.leaked_param;
}
`
	if err := cache.Open(uri, 1, content); err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	ctx := context.Background()

	positions := []struct {
		name string
		pos  protocol.Position
	}{
		{"bank_inner_p", protocol.Position{Line: 17, Character: 2}},
		{"this.bank_inner_p", protocol.Position{Line: 18, Character: 8}},
		{"tmpl_p", protocol.Position{Line: 23, Character: 2}},
		{"leaked_param", protocol.Position{Line: 24, Character: 2}},
		{"dev.leaked_param", protocol.Position{Line: 25, Character: 6}},
	}

	for _, tc := range positions {
		loc, err := lsp.ResolveDefinition(ctx, cache, uri, "", tc.pos)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", tc.name, err)
		}
		if loc != nil {
			t.Fatalf(
				"%s: expected in each symbol not to leak, got %+v",
				tc.name,
				loc,
			)
		}
	}
}

func TestResolveDefinition_IsKeywordParenthesizedList(t *testing.T) {
	cache := lsp.NewFsCache()
	uri := "file:///workspace/is_paren.dml"
	content := `dml 1.4;
device test_dev;

template tmpl_a {
	param a = 1;
}

template tmpl_b {
	param b = 2;
}

bank regs {
	is (tmpl_a, tmpl_b);
	is ;
}
`
	if err := cache.Open(uri, 1, content); err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	ctx := context.Background()

	// Cursor on `is` in `\tis (tmpl_a, tmpl_b);` (line 12, col 1)
	loc, err := lsp.ResolveDefinition(
		ctx,
		cache,
		uri,
		"",
		protocol.Position{Line: 12, Character: 1},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc == nil || loc.Range.Start.Line != 3 {
		t.Fatalf("expected tmpl_a on line 3, got %+v", loc)
	}

	// Cursor on bare `is ;` without template name must return nil
	locEmpty, err := lsp.ResolveDefinition(
		ctx,
		cache,
		uri,
		"",
		protocol.Position{Line: 13, Character: 1},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if locEmpty != nil {
		t.Fatalf("expected nil for bare is, got %+v", locEmpty)
	}
}
