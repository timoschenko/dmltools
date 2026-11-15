// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lsp_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

func TestDocumentSymbol_Capabilities(t *testing.T) {
	h := setupSyncHarness(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rootURI := uri.File(t.TempDir())
	res, err := h.client.Initialize(ctx, &protocol.InitializeParams{
		RootURI: &rootURI,
	})
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	if res.Capabilities.DocumentSymbolProvider != protocol.Boolean(true) {
		t.Fatalf(
			"DocumentSymbolProvider = %v, want true",
			res.Capabilities.DocumentSymbolProvider,
		)
	}
}

func assertRangeEncloses(
	t *testing.T,
	outer, inner protocol.Range,
	name string,
) {
	t.Helper()
	if inner.Start.Line < outer.Start.Line ||
		(inner.Start.Line == outer.Start.Line && inner.Start.Character < outer.Start.Character) {
		t.Errorf(
			"%s: inner start %+v is before outer start %+v",
			name,
			inner.Start,
			outer.Start,
		)
	}
	if inner.End.Line > outer.End.Line ||
		(inner.End.Line == outer.End.Line && inner.End.Character > outer.End.Character) {
		t.Errorf(
			"%s: inner end %+v is after outer end %+v",
			name,
			inner.End,
			outer.End,
		)
	}
}

func verifySymbolHierarchy(t *testing.T, syms []protocol.DocumentSymbol) {
	t.Helper()
	for _, sym := range syms {
		if sym.Name == "" {
			t.Errorf("found symbol with empty Name: %+v", sym)
		}
		assertRangeEncloses(t, sym.Range, sym.SelectionRange, sym.Name)
		for _, child := range sym.Children {
			assertRangeEncloses(
				t,
				sym.Range,
				child.Range,
				sym.Name+" -> "+child.Name,
			)
		}
		if len(sym.Children) > 0 {
			verifySymbolHierarchy(t, sym.Children)
		}
	}
}

func findDocSymbol(
	syms []protocol.DocumentSymbol,
	name string,
) (protocol.DocumentSymbol, bool) {
	for _, s := range syms {
		if s.Name == name {
			return s, true
		}
	}
	return protocol.DocumentSymbol{}, false
}

func TestDocumentSymbol_Comprehensive(t *testing.T) {
	dir := t.TempDir()
	h := initSyncHarness(t, dir)
	docURI := uri.File(filepath.Join(dir, "symbols.dml"))

	content := strings.Join([]string{
		"dml 1.4;",
		"device test_device;",
		"template base_tpl {",
		"    param p_base = 1;",
		"    method m_base();",
		"    in each sub_tpl {",
		"        param p_sub = 2;",
		"    }",
		"}",
		"in each (tpl1, tpl2) {",
		"    param ineach_p = 10;",
		"}",
		"bank regs {",
		"    register ctrl size 4 @ 0x0 {",
		"        field en @ [0];",
		"        method read() -> (uint32) {",
		"            local int ignored_local = 5;",
		"            return 1;",
		"        }",
		"    }",
		"}",
		"group grp {}",
		"port p {}",
		"connect conn {}",
		"implement impl {}",
		"interface iface {}",
		"attribute attr {}",
		"event ev {}",
		"subdevice subdev {}",
		"method top_method(int arg) -> (int res) {",
		"    local int ignored_val = 10;",
		"    return arg;",
		"}",
		"param top_param = 100;",
		"session int ses_single = 1;",
		"session (int ses_a, int ses_b) = (2, 3);",
		"saved int sav_val;",
		"hook(int) on_event;",
		"constant MY_CONST = 42;",
		"loggroup lg_debug;",
		"typedef uint32 custom_u32;",
		"typedef struct point_s { int x; int y; } point_t;",
		"typedef layout \"big-endian\" { bitfields 32 { uint16 lo @ [15:0]; } reg; } my_layout_t;",
		"extern void ext_func(int a);",
		"extern int ext_var;",
		"extern typedef struct ext_s { int ext_f; } ext_s_t;",
		"#if (1) {",
		"    param if_p = 1;",
		"} #else {",
		"    param else_p = 2;",
		"}",
	}, "\n") + "\n"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	openAndWait(t, ctx, h, docURI, content)

	res, err := h.client.DocumentSymbol(ctx, &protocol.DocumentSymbolParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: docURI},
	})
	if err != nil {
		t.Fatalf("DocumentSymbol request failed: %v", err)
	}

	symSlice, ok := res.(protocol.DocumentSymbolSlice)
	if !ok {
		t.Fatalf("expected protocol.DocumentSymbolSlice, got %T", res)
	}

	verifySymbolHierarchy(t, symSlice)

	// 1. device
	dev, ok := findDocSymbol(symSlice, "test_device")
	if !ok {
		t.Fatalf("test_device not found")
	}
	if dev.Kind != protocol.SymbolKindClass || dev.Detail == nil ||
		*dev.Detail != "device" {
		t.Errorf("test_device: Kind=%v Detail=%v", dev.Kind, dev.Detail)
	}

	// 2. template base_tpl
	tpl, ok := findDocSymbol(symSlice, "base_tpl")
	if !ok {
		t.Fatalf("base_tpl not found")
	}
	if tpl.Kind != protocol.SymbolKindInterface || tpl.Detail == nil ||
		*tpl.Detail != "template" {
		t.Errorf("base_tpl: Kind=%v Detail=%v", tpl.Kind, tpl.Detail)
	}
	if len(tpl.Children) != 3 {
		t.Fatalf("base_tpl children count = %d, want 3", len(tpl.Children))
	}
	pBase, ok := findDocSymbol(tpl.Children, "p_base")
	if !ok || pBase.Kind != protocol.SymbolKindProperty ||
		pBase.Detail == nil ||
		*pBase.Detail != "param" {
		t.Errorf("p_base: %+v", pBase)
	}
	mBase, ok := findDocSymbol(tpl.Children, "m_base")
	if !ok || mBase.Kind != protocol.SymbolKindMethod || mBase.Detail == nil ||
		*mBase.Detail != "method" {
		t.Errorf("m_base: %+v", mBase)
	}
	ieSub, ok := findDocSymbol(tpl.Children, "in each sub_tpl")
	if !ok || ieSub.Kind != protocol.SymbolKindNamespace ||
		ieSub.Detail == nil ||
		*ieSub.Detail != "in each" {
		t.Errorf("in each sub_tpl: %+v", ieSub)
	}
	if len(ieSub.Children) != 1 || ieSub.Children[0].Name != "p_sub" {
		t.Errorf("in each sub_tpl children = %+v", ieSub.Children)
	}

	// 3. Top-level in each
	ieTop, ok := findDocSymbol(symSlice, "in each (tpl1, tpl2)")
	if !ok {
		t.Fatalf("in each (tpl1, tpl2) not found")
	}
	if ieTop.Kind != protocol.SymbolKindNamespace || len(ieTop.Children) != 1 ||
		ieTop.Children[0].Name != "ineach_p" {
		t.Errorf("in each top: %+v", ieTop)
	}

	// 4. bank regs -> register ctrl -> field en, method read
	regs, ok := findDocSymbol(symSlice, "regs")
	if !ok {
		t.Fatalf("regs not found")
	}
	if regs.Kind != protocol.SymbolKindClass || regs.Detail == nil ||
		*regs.Detail != "bank" {
		t.Errorf("regs: Kind=%v Detail=%v", regs.Kind, regs.Detail)
	}
	ctrl, ok := findDocSymbol(regs.Children, "ctrl")
	if !ok {
		t.Fatalf("ctrl not found in regs")
	}
	if ctrl.Kind != protocol.SymbolKindClass || ctrl.Detail == nil ||
		*ctrl.Detail != "register" {
		t.Errorf("ctrl: Kind=%v Detail=%v", ctrl.Kind, ctrl.Detail)
	}
	en, ok := findDocSymbol(ctrl.Children, "en")
	if !ok || en.Kind != protocol.SymbolKindClass || en.Detail == nil ||
		*en.Detail != "field" {
		t.Errorf("en: %+v", en)
	}
	readMeth, ok := findDocSymbol(ctrl.Children, "read")
	if !ok || readMeth.Kind != protocol.SymbolKindMethod {
		t.Errorf("read: %+v", readMeth)
	}
	// Method is a leaf symbol: locals inside method are omitted
	if len(readMeth.Children) != 0 {
		t.Errorf("read children count = %d, want 0", len(readMeth.Children))
	}

	// 5. Various objects
	for name, expectedDetail := range map[string]string{
		"grp":    "group",
		"p":      "port",
		"conn":   "connect",
		"impl":   "implement",
		"iface":  "interface",
		"attr":   "attribute",
		"ev":     "event",
		"subdev": "subdevice",
	} {
		obj, ok := findDocSymbol(symSlice, name)
		if !ok {
			t.Errorf("object %s not found", name)
			continue
		}
		if obj.Kind != protocol.SymbolKindClass || obj.Detail == nil ||
			*obj.Detail != expectedDetail {
			t.Errorf("object %s: Kind=%v Detail=%v", name, obj.Kind, obj.Detail)
		}
	}

	// 6. Top method
	topMeth, ok := findDocSymbol(symSlice, "top_method")
	if !ok || topMeth.Kind != protocol.SymbolKindMethod ||
		len(topMeth.Children) != 0 {
		t.Errorf("top_method: %+v", topMeth)
	}

	// 7. Session single and multi-variable
	sesSingle, ok := findDocSymbol(symSlice, "ses_single")
	if !ok || sesSingle.Kind != protocol.SymbolKindVariable ||
		sesSingle.Detail == nil ||
		*sesSingle.Detail != "session" {
		t.Errorf("ses_single: %+v", sesSingle)
	}
	sesA, ok := findDocSymbol(symSlice, "ses_a")
	if !ok || sesA.Kind != protocol.SymbolKindVariable {
		t.Errorf("ses_a: %+v", sesA)
	}
	sesB, ok := findDocSymbol(symSlice, "ses_b")
	if !ok || sesB.Kind != protocol.SymbolKindVariable {
		t.Errorf("ses_b: %+v", sesB)
	}

	// 8. Saved, Hook, Constant, Loggroup
	sav, ok := findDocSymbol(symSlice, "sav_val")
	if !ok || sav.Kind != protocol.SymbolKindVariable || sav.Detail == nil ||
		*sav.Detail != "saved" {
		t.Errorf("sav_val: %+v", sav)
	}
	hook, ok := findDocSymbol(symSlice, "on_event")
	if !ok || hook.Kind != protocol.SymbolKindEvent || hook.Detail == nil ||
		*hook.Detail != "hook" {
		t.Errorf("on_event: %+v", hook)
	}
	cst, ok := findDocSymbol(symSlice, "MY_CONST")
	if !ok || cst.Kind != protocol.SymbolKindConstant || cst.Detail == nil ||
		*cst.Detail != "constant" {
		t.Errorf("MY_CONST: %+v", cst)
	}
	lg, ok := findDocSymbol(symSlice, "lg_debug")
	if !ok || lg.Kind != protocol.SymbolKindEnumMember || lg.Detail == nil ||
		*lg.Detail != "loggroup" {
		t.Errorf("lg_debug: %+v", lg)
	}

	// 9. Typedef scalar vs struct vs layout
	scalarType, ok := findDocSymbol(symSlice, "custom_u32")
	if !ok || scalarType.Kind != protocol.SymbolKindTypeParameter ||
		len(scalarType.Children) != 0 {
		t.Errorf("custom_u32: %+v", scalarType)
	}

	pointType, ok := findDocSymbol(symSlice, "point_t")
	if !ok || pointType.Kind != protocol.SymbolKindStruct ||
		len(pointType.Children) != 2 {
		t.Fatalf("point_t: %+v", pointType)
	}
	fieldX, ok := findDocSymbol(pointType.Children, "x")
	if !ok || fieldX.Kind != protocol.SymbolKindField || fieldX.Detail == nil ||
		*fieldX.Detail != "int" {
		t.Errorf("point_t.x: %+v", fieldX)
	}
	fieldY, ok := findDocSymbol(pointType.Children, "y")
	if !ok || fieldY.Kind != protocol.SymbolKindField || fieldY.Detail == nil ||
		*fieldY.Detail != "int" {
		t.Errorf("point_t.y: %+v", fieldY)
	}

	layoutType, ok := findDocSymbol(symSlice, "my_layout_t")
	if !ok || layoutType.Kind != protocol.SymbolKindStruct ||
		len(layoutType.Children) != 1 {
		t.Fatalf("my_layout_t: %+v", layoutType)
	}
	regField, ok := findDocSymbol(layoutType.Children, "reg")
	if !ok || regField.Kind != protocol.SymbolKindStruct ||
		len(regField.Children) != 1 {
		t.Fatalf("reg: %+v", regField)
	}
	loField, ok := findDocSymbol(regField.Children, "lo")
	if !ok || loField.Kind != protocol.SymbolKindField ||
		loField.Detail == nil ||
		*loField.Detail != "uint16" {
		t.Errorf("lo: %+v", loField)
	}

	// 10. Externs
	extFunc, ok := findDocSymbol(symSlice, "ext_func")
	if !ok || extFunc.Kind != protocol.SymbolKindFunction ||
		extFunc.Detail == nil ||
		*extFunc.Detail != "extern" {
		t.Errorf("ext_func: %+v", extFunc)
	}
	extVar, ok := findDocSymbol(symSlice, "ext_var")
	if !ok || extVar.Kind != protocol.SymbolKindVariable ||
		extVar.Detail == nil ||
		*extVar.Detail != "extern" {
		t.Errorf("ext_var: %+v", extVar)
	}
	extStruct, ok := findDocSymbol(symSlice, "ext_s_t")
	if !ok || extStruct.Kind != protocol.SymbolKindStruct ||
		extStruct.Detail == nil ||
		*extStruct.Detail != "extern typedef" {
		t.Errorf("ext_s_t: %+v", extStruct)
	}
	if len(extStruct.Children) != 1 || extStruct.Children[0].Name != "ext_f" {
		t.Errorf("ext_s_t children: %+v", extStruct.Children)
	}

	// 11. #if / #else flattened into top-level parent
	ifP, ok := findDocSymbol(symSlice, "if_p")
	if !ok || ifP.Kind != protocol.SymbolKindProperty {
		t.Errorf("if_p: %+v", ifP)
	}
	elseP, ok := findDocSymbol(symSlice, "else_p")
	if !ok || elseP.Kind != protocol.SymbolKindProperty {
		t.Errorf("else_p: %+v", elseP)
	}
}

func TestDocumentSymbol_IncrementalEdit(t *testing.T) {
	dir := t.TempDir()
	h := initSyncHarness(t, dir)
	docURI := uri.File(filepath.Join(dir, "edit.dml"))

	v1 := "dml 1.4;\ntemplate alpha {}\n"
	v2 := "dml 1.4;\nbank beta {}\n"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	openAndWait(t, ctx, h, docURI, v1)

	res1, err := h.client.DocumentSymbol(ctx, &protocol.DocumentSymbolParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: docURI},
	})
	if err != nil {
		t.Fatalf("DocumentSymbol v1 failed: %v", err)
	}
	syms1, ok := res1.(protocol.DocumentSymbolSlice)
	if !ok {
		t.Fatalf("expected DocumentSymbolSlice, got %T", res1)
	}
	if len(syms1) != 1 || syms1[0].Name != "alpha" {
		t.Fatalf("v1 symbols = %+v, want [alpha]", syms1)
	}

	if err := h.client.DidChange(
		ctx,
		&protocol.DidChangeTextDocumentParams{
			TextDocument: protocol.VersionedTextDocumentIdentifier{
				TextDocumentIdentifier: protocol.TextDocumentIdentifier{
					URI: docURI,
				},
				Version: 2,
			},
			ContentChanges: []protocol.TextDocumentContentChangeEvent{
				&protocol.TextDocumentContentChangeWholeDocument{Text: v2},
			},
		},
	); err != nil {
		t.Fatalf("DidChange failed: %v", err)
	}

	waitFor(t, 2*time.Second, func() bool {
		data, readErr := h.server.FsCache().ReadFile(string(docURI))
		return readErr == nil && string(data) == v2
	})

	res2, err := h.client.DocumentSymbol(ctx, &protocol.DocumentSymbolParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: docURI},
	})
	if err != nil {
		t.Fatalf("DocumentSymbol v2 failed: %v", err)
	}
	syms2, ok := res2.(protocol.DocumentSymbolSlice)
	if !ok {
		t.Fatalf("expected DocumentSymbolSlice, got %T", res2)
	}
	if len(syms2) != 1 || syms2[0].Name != "beta" {
		t.Fatalf("v2 symbols = %+v, want [beta]", syms2)
	}
}

func TestDocumentSymbol_SyntaxErrorsAndDML12(t *testing.T) {
	dir := t.TempDir()
	h := initSyncHarness(t, dir)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. DML 1.2 file returns empty slice
	dml12URI := uri.File(filepath.Join(dir, "legacy.dml"))
	dml12Content := "dml 1.2;\ndevice leg_dev;\nbank b {}\n"
	openAndWait(t, ctx, h, dml12URI, dml12Content)

	res12, err := h.client.DocumentSymbol(ctx, &protocol.DocumentSymbolParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: dml12URI},
	})
	if err != nil {
		t.Fatalf("DocumentSymbol dml12 failed: %v", err)
	}
	syms12, ok := res12.(protocol.DocumentSymbolSlice)
	if !ok {
		t.Fatalf("expected DocumentSymbolSlice, got %T", res12)
	}
	if len(syms12) != 0 {
		t.Errorf("DML 1.2 symbols len = %d, want 0", len(syms12))
	}

	// 2. Syntax error in buffer: returns valid partial symbols
	brokenURI := uri.File(filepath.Join(dir, "broken.dml"))
	brokenContent := "dml 1.4;\ndevice my_dev;\nbank good_b {}\nbank broken_b { register r\n"
	openAndWait(t, ctx, h, brokenURI, brokenContent)

	brokenRes, err := h.client.DocumentSymbol(
		ctx,
		&protocol.DocumentSymbolParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: brokenURI},
		},
	)
	if err != nil {
		t.Fatalf("DocumentSymbol broken failed: %v", err)
	}
	brokenSyms, ok := brokenRes.(protocol.DocumentSymbolSlice)
	if !ok {
		t.Fatalf("expected DocumentSymbolSlice, got %T", brokenRes)
	}
	if len(brokenSyms) == 0 {
		t.Fatal("expected partial symbols for broken buffer, got none")
	}
	verifySymbolHierarchy(t, brokenSyms)
}

func TestDocumentSymbol_UTF16SurrogatePairs(t *testing.T) {
	dir := t.TempDir()
	h := initSyncHarness(t, dir)
	docURI := uri.File(filepath.Join(dir, "emoji.dml"))

	content := strings.Join([]string{
		"dml 1.4;",
		"/* 🚀✨ */",
		"device emoji_dev;",
		"bank b {",
		"    /* 😀 */ param p = 1;",
		"}",
	}, "\n") + "\n"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	openAndWait(t, ctx, h, docURI, content)

	res, err := h.client.DocumentSymbol(ctx, &protocol.DocumentSymbolParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: docURI},
	})
	if err != nil {
		t.Fatalf("DocumentSymbol request failed: %v", err)
	}

	syms, ok := res.(protocol.DocumentSymbolSlice)
	if !ok {
		t.Fatalf("expected DocumentSymbolSlice, got %T", res)
	}
	verifySymbolHierarchy(t, syms)

	dev, ok := findDocSymbol(syms, "emoji_dev")
	if !ok {
		t.Fatalf("emoji_dev not found")
	}
	if dev.Range.Start.Line != 2 {
		t.Errorf("emoji_dev line = %d, want 2", dev.Range.Start.Line)
	}
}

func TestDocumentSymbol_MemberFieldTypes(t *testing.T) {
	dir := t.TempDir()
	h := initSyncHarness(t, dir)
	docURI := uri.File(filepath.Join(dir, "fields.dml"))

	content := strings.Join([]string{
		"dml 1.4;",
		"device dev;",
		"typedef struct fields_s {",
		"    int plain_int;",
		"    unsigned implicit_uint;",
		"    signed implicit_sint;",
		"    unsigned int full_uint;",
		"    signed char full_schar;",
		"    const int* ptr_int;",
		"    foo_t ** const ptr_ptr;",
		"    uint32 arr[16];",
		"    struct { int nested_x; } anon_st;",
		"    union u_tag { int u_f; } named_un;",
		"    layout \"big-endian\" { int l_f; } l_inst;",
		"    bitfields 32 { int bf_f @ [0]; } bf_inst;",
		"} fields_t;",
	}, "\n") + "\n"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	openAndWait(t, ctx, h, docURI, content)

	res, err := h.client.DocumentSymbol(ctx, &protocol.DocumentSymbolParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: docURI},
	})
	if err != nil {
		t.Fatalf("DocumentSymbol request failed: %v", err)
	}

	syms, ok := res.(protocol.DocumentSymbolSlice)
	if !ok {
		t.Fatalf("expected DocumentSymbolSlice, got %T", res)
	}
	verifySymbolHierarchy(t, syms)

	fieldsT, ok := findDocSymbol(syms, "fields_t")
	if !ok || fieldsT.Kind != protocol.SymbolKindStruct {
		t.Fatalf("fields_t not found or wrong kind: %+v", fieldsT)
	}

	expectedFields := map[string]struct {
		kind   protocol.SymbolKind
		detail string
	}{
		"plain_int":     {protocol.SymbolKindField, "int"},
		"implicit_uint": {protocol.SymbolKindField, "unsigned"},
		"implicit_sint": {protocol.SymbolKindField, "signed"},
		"full_uint":     {protocol.SymbolKindField, "unsigned int"},
		"full_schar":    {protocol.SymbolKindField, "signed char"},
		"ptr_int":       {protocol.SymbolKindField, "const int*"},
		"ptr_ptr":       {protocol.SymbolKindField, "foo_t**const"},
		"arr":           {protocol.SymbolKindField, "uint32[16]"},
		"anon_st":       {protocol.SymbolKindStruct, "struct"},
		"named_un":      {protocol.SymbolKindStruct, "union u_tag"},
		"l_inst":        {protocol.SymbolKindStruct, "layout"},
		"bf_inst":       {protocol.SymbolKindStruct, "bitfields"},
	}

	if len(fieldsT.Children) != len(expectedFields) {
		t.Fatalf(
			"fields_t child count = %d, want %d",
			len(fieldsT.Children),
			len(expectedFields),
		)
	}

	for name, expected := range expectedFields {
		f, found := findDocSymbol(fieldsT.Children, name)
		if !found {
			t.Errorf("field %s not found", name)
			continue
		}
		if f.Kind != expected.kind {
			t.Errorf(
				"field %s: kind = %v, want %v",
				name,
				f.Kind,
				expected.kind,
			)
		}
		if f.Detail == nil || *f.Detail != expected.detail {
			t.Errorf(
				"field %s: detail = %v, want %s",
				name,
				f.Detail,
				expected.detail,
			)
		}
	}
}

func TestDocumentSymbol_EmptyAndDML12Caching(t *testing.T) {
	dir := t.TempDir()
	h := initSyncHarness(t, dir)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Symbol-less DML 1.4 file returns non-nil empty slice
	emptyURI := uri.File(filepath.Join(dir, "import_only.dml"))
	emptyContent := "dml 1.4;\nimport \"unused.dml\";\n"
	openAndWait(t, ctx, h, emptyURI, emptyContent)

	res1, err := h.client.DocumentSymbol(ctx, &protocol.DocumentSymbolParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: emptyURI},
	})
	if err != nil {
		t.Fatalf("DocumentSymbol empty request failed: %v", err)
	}
	syms1, ok := res1.(protocol.DocumentSymbolSlice)
	if !ok || len(syms1) != 0 {
		t.Fatalf(
			"expected empty DocumentSymbolSlice, got %T (%d symbols)",
			res1,
			len(syms1),
		)
	}

	// Repeat request should hit cache
	res2, err := h.client.DocumentSymbol(ctx, &protocol.DocumentSymbolParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: emptyURI},
	})
	if err != nil {
		t.Fatalf("second DocumentSymbol empty request failed: %v", err)
	}
	syms2, ok := res2.(protocol.DocumentSymbolSlice)
	if !ok || len(syms2) != 0 {
		t.Fatalf(
			"expected empty DocumentSymbolSlice, got %T (%d symbols)",
			res2,
			len(syms2),
		)
	}
}
