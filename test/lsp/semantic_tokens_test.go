// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lsp_test

import (
	"context"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
	"unicode/utf16"

	"github.com/timoschenko/dmltools/pkg/lsp"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

type decodedSemanticToken struct {
	line      uint32
	char      uint32
	length    uint32
	tokenType string
	modifiers []string
	text      string
}

func utf16Slice(s string, startCol, length uint32) string {
	units := utf16.Encode([]rune(s))
	start := int(startCol)
	end := int(startCol + length)
	if start < 0 || end > len(units) || start > end {
		return ""
	}
	return string(utf16.Decode(units[start:end]))
}

func decodeSemanticTokens(
	t *testing.T,
	data []uint32,
	content string,
) []decodedSemanticToken {
	t.Helper()
	if len(data)%5 != 0 {
		t.Fatalf("semantic token data len = %d, want multiple of 5", len(data))
	}

	legend := lsp.SemanticTokensLegend()
	lines := strings.Split(content, "\n")
	out := make([]decodedSemanticToken, 0, len(data)/5)

	var line uint32
	var char uint32
	for i := 0; i < len(data); i += 5 {
		deltaLine := data[i]
		deltaChar := data[i+1]
		length := data[i+2]
		typeIdx := int(data[i+3])
		modMask := data[i+4]

		if deltaLine > 0 {
			line += deltaLine
			char = deltaChar
		} else {
			char += deltaChar
		}

		if typeIdx < 0 || typeIdx >= len(legend.TokenTypes) {
			t.Fatalf("token type index %d out of legend range", typeIdx)
		}

		var mods []string
		for bit, name := range legend.TokenModifiers {
			if modMask&(1<<bit) != 0 {
				mods = append(mods, name)
			}
		}

		var text string
		if int(line) < len(lines) {
			text = utf16Slice(lines[line], char, length)
		}

		out = append(out, decodedSemanticToken{
			line:      line,
			char:      char,
			length:    length,
			tokenType: legend.TokenTypes[typeIdx],
			modifiers: mods,
			text:      text,
		})
	}
	return out
}

func findToken(
	tokens []decodedSemanticToken,
	line uint32,
	text string,
) (decodedSemanticToken, bool) {
	for _, tok := range tokens {
		if tok.line == line && tok.text == text {
			return tok, true
		}
	}
	return decodedSemanticToken{}, false
}

func lineOfText(t *testing.T, content, needle string) uint32 {
	t.Helper()
	for i, line := range strings.Split(content, "\n") {
		if strings.Contains(line, needle) {
			return uint32(i) //nolint:gosec // i is non-negative
		}
	}
	t.Fatalf("line containing %q not found", needle)
	return 0
}

func assertToken(
	t *testing.T,
	tokens []decodedSemanticToken,
	line uint32,
	text string,
	wantType string,
	wantMods ...string,
) {
	t.Helper()
	tok, ok := findToken(tokens, line, text)
	if !ok {
		t.Fatalf(
			"token %q on line %d not found (tokens on line: %v)",
			text,
			line,
			tokensOnLine(tokens, line),
		)
	}
	if tok.tokenType != wantType {
		t.Errorf(
			"line %d %q: tokenType = %q, want %q",
			line,
			text,
			tok.tokenType,
			wantType,
		)
	}
	slices.Sort(wantMods)
	gotMods := slices.Clone(tok.modifiers)
	slices.Sort(gotMods)
	if !slices.Equal(gotMods, wantMods) {
		t.Errorf(
			"line %d %q: modifiers = %v, want %v",
			line,
			text,
			gotMods,
			wantMods,
		)
	}
}

func tokensOnLine(
	tokens []decodedSemanticToken,
	line uint32,
) []decodedSemanticToken {
	var out []decodedSemanticToken
	for _, tok := range tokens {
		if tok.line == line {
			out = append(out, tok)
		}
	}
	return out
}

func initSyncHarness(t *testing.T, rootDir string) *syncHarness {
	t.Helper()
	h := setupSyncHarness(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rootURI := uri.File(rootDir)
	_, err := h.client.Initialize(ctx, &protocol.InitializeParams{
		RootURI: &rootURI,
	})
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	if err := h.client.Initialized(
		ctx,
		&protocol.InitializedParams{},
	); err != nil {
		t.Fatalf("Initialized failed: %v", err)
	}
	return h
}

func openAndWait(
	t *testing.T,
	ctx context.Context,
	h *syncHarness,
	docURI uri.URI,
	content string,
) {
	t.Helper()
	if err := h.client.DidOpen(ctx, &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        docURI,
			LanguageID: "dml",
			Version:    1,
			Text:       content,
		},
	}); err != nil {
		t.Fatalf("DidOpen failed: %v", err)
	}
	waitFor(t, 2*time.Second, func() bool {
		data, err := h.server.FsCache().ReadFile(string(docURI))
		return err == nil && string(data) == content
	})
}

func TestSemanticTokensCapabilities(t *testing.T) {
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

	prov := res.Capabilities.SemanticTokensProvider
	opts, ok := prov.(*protocol.SemanticTokensOptions)
	if !ok || opts == nil {
		t.Fatalf(
			"SemanticTokensProvider = %T, want *SemanticTokensOptions",
			res.Capabilities.SemanticTokensProvider,
		)
	}

	wantLegend := lsp.SemanticTokensLegend()
	if !slices.Equal(opts.Legend.TokenTypes, wantLegend.TokenTypes) {
		t.Errorf("Legend.TokenTypes = %v, want %v",
			opts.Legend.TokenTypes, wantLegend.TokenTypes)
	}
	if !slices.Equal(opts.Legend.TokenModifiers, wantLegend.TokenModifiers) {
		t.Errorf("Legend.TokenModifiers = %v, want %v",
			opts.Legend.TokenModifiers, wantLegend.TokenModifiers)
	}
	if opts.Full != protocol.Boolean(true) {
		t.Errorf("opts.Full = %v, want true", opts.Full)
	}
	if opts.Range != protocol.Boolean(true) {
		t.Errorf("opts.Range = %v, want true", opts.Range)
	}
}

func TestSemanticTokensFullComprehensive(t *testing.T) {
	dir := t.TempDir()
	h := initSyncHarness(t, dir)
	docURI := uri.File(filepath.Join(dir, "device.dml"))

	content := strings.Join([]string{
		"dml 1.4;",
		"provisional explicit_param_decls;",
		"device my_dev;",
		"bitorder be;",
		"import \"external_missing.dml\";",
		"// line comment",
		"/* multi line",
		"",
		" * block comment */",
		"constant MAX_SZ = 0x10;",
		"loggroup lg_io;",
		"typedef uint32 my_u32_t;",
		"typedef struct my_s { uint32 f; } my_s_t;",
		"typedef layout \"big-endian\" { uint32 lv; } my_lay_t;",
		"typedef bitfields 32 { uint16 lo @ [15:0]; } my_bf_t;",
		"extern void ext_fn(const int arg);",
		"extern const int ext_val;",
		"header %{ int raw_c_ident = 1; %}",
		"template abs_tpl {",
		"    shared method abs_m(int p) -> (int ret);",
		"}",
		"template base_tpl is abs_tpl {",
		"    param p_def = 42;",
		"    param p_decl : int;",
		"    param size = 4;",
		"}",
		"bank regs[i < 4] is base_tpl {",
		"    session int ses_x = 1;",
		"    saved int sav_y;",
		"    hook(int) ev_hook;",
		"    register r size 4 {",
		"        param size = 8;",
		"    }",
		"    inline method run(const int inp) -> (int outp) throws {",
		"        local const int c = MAX_SZ;",
		"        local int uninit;",
		"        static int st = 0b101;",
		"        local unsigned bare_u = 1;",
		"        local signed int si = 2;",
		"        local struct my_s ms = { .f = 1 };",
		"        local my_s_t ms2;",
		"        local my_u32_t u = inp + outp + bare_u + si;",
		"        local bool flag = true;",
		"        local ext_unknown_t ext_var;",
		"        uninit += c + st + ses_x + sav_y + p_def + i;",
		"        #if (flag) {",
		"            log info, 1, lg_io: \"msg %d\", regs[0].r.size;",
		"        } #else {",
		"            log critical, 1, other_lg: \"err\";",
		"        }",
		"        for (local int k = 0; k < 2; k++) {",
		"            uninit += k;",
		"        }",
		"        foreach obj in (each abs_tpl in (dev.regs)) {",
		"            obj.abs_m(1);",
		"        }",
		"        local sequence(abs_tpl) seq;",
		"        after 0.1 s: ext_fn(1);",
		"        after ev_hook -> (mx, my): cb_fn(mx + my);",
		"        unknown_fn(this.p_def, parent.name, ms.f);",
		"        ms.unknown_obj.do_thing();",
		"        outp = ms.missing_prop;",
		"        if (unknown_ident == 0) goto cleanup;",
		"        return;",
		"    }",
		"}",
	}, "\n") + "\n"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	openAndWait(t, ctx, h, docURI, content)

	res, err := h.client.SemanticTokensFull(ctx, &protocol.SemanticTokensParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: docURI},
	})
	if err != nil {
		t.Fatalf("SemanticTokensFull failed: %v", err)
	}

	toks := decodeSemanticTokens(t, res.Data, content)

	// Header, device, bitorder, provisional
	assertToken(t, toks, lineOfText(t, content, "dml 1.4"), "dml", "keyword")
	assertToken(t, toks, lineOfText(t, content, "dml 1.4"), "1.4", "number")
	assertToken(t, toks, lineOfText(t, content, "device my_dev"),
		"my_dev", "class", "definition")
	assertToken(t, toks, lineOfText(t, content, "bitorder be"), "be", "keyword")
	assertToken(t, toks, lineOfText(t, content, "provisional"),
		"explicit_param_decls", "property")

	// Comments (including multiline block comment split and empty line skipped)
	assertToken(t, toks, lineOfText(t, content, "// line comment"),
		"// line comment", "comment")
	assertToken(t, toks, lineOfText(t, content, "/* multi line"),
		"/* multi line", "comment")
	assertToken(t, toks, lineOfText(t, content, "* block comment"),
		" * block comment */", "comment")

	// Constants, loggroup, typedefs, structs, layouts, bitfields, externs
	assertToken(t, toks, lineOfText(t, content, "constant MAX_SZ"),
		"MAX_SZ", "enumMember", "definition", "readonly")
	assertToken(t, toks, lineOfText(t, content, "constant MAX_SZ"),
		"0x10", "number")
	assertToken(t, toks, lineOfText(t, content, "loggroup lg_io"),
		"lg_io", "enumMember", "definition")
	assertToken(t, toks, lineOfText(t, content, "typedef uint32"),
		"uint32", "type", "defaultLibrary")
	assertToken(t, toks, lineOfText(t, content, "typedef uint32"),
		"my_u32_t", "type", "definition")

	sLine := lineOfText(t, content, "typedef struct my_s")
	assertToken(t, toks, sLine, "my_s", "struct", "definition")
	assertToken(t, toks, sLine, "f", "property", "declaration")
	assertToken(t, toks, sLine, "my_s_t", "struct", "definition")

	lLine := lineOfText(t, content, "typedef layout")
	assertToken(t, toks, lLine, "lv", "property", "declaration")
	assertToken(t, toks, lLine, "my_lay_t", "struct", "definition")

	bLine := lineOfText(t, content, "typedef bitfields")
	assertToken(t, toks, bLine, "uint16", "type", "defaultLibrary")
	assertToken(t, toks, bLine, "lo", "property", "declaration")
	assertToken(t, toks, bLine, "my_bf_t", "struct", "definition")

	assertToken(t, toks, lineOfText(t, content, "extern void ext_fn"),
		"ext_fn", "function", "declaration")
	assertToken(t, toks, lineOfText(t, content, "extern void ext_fn"),
		"arg", "parameter", "declaration", "readonly")
	assertToken(t, toks, lineOfText(t, content, "extern const int ext_val"),
		"ext_val", "variable", "declaration", "readonly")

	// C_CODE_BLOCK untokenized
	cLine := lineOfText(t, content, "raw_c_ident")
	assertToken(t, toks, cLine, "header", "keyword")
	if _, found := findToken(toks, cLine, "raw_c_ident"); found {
		t.Errorf("C block content raw_c_ident should remain untokenized")
	}

	// Templates, abstract shared method, params, objects, hooks, session/saved
	assertToken(t, toks, lineOfText(t, content, "template abs_tpl"),
		"abs_tpl", "interface", "definition")
	assertToken(t, toks, lineOfText(t, content, "shared method abs_m"),
		"shared", "modifier")
	assertToken(t, toks, lineOfText(t, content, "shared method abs_m"),
		"abs_m", "method", "declaration", "abstract")
	assertToken(t, toks, lineOfText(t, content, "shared method abs_m"),
		"p", "parameter", "declaration")
	assertToken(t, toks, lineOfText(t, content, "shared method abs_m"),
		"ret", "parameter", "declaration")

	assertToken(t, toks, lineOfText(t, content, "template base_tpl"),
		"abs_tpl", "interface")
	assertToken(t, toks, lineOfText(t, content, "param p_def"),
		"p_def", "property", "definition")
	assertToken(t, toks, lineOfText(t, content, "param p_decl"),
		"p_decl", "property", "declaration")
	assertToken(t, toks, lineOfText(t, content, "param size = 4"),
		"size", "property", "definition")

	bankLine := lineOfText(t, content, "bank regs[i < 4]")
	assertToken(t, toks, bankLine, "regs", "class", "definition")
	assertToken(t, toks, bankLine, "i", "variable", "definition")
	assertToken(t, toks, bankLine, "base_tpl", "interface")

	assertToken(t, toks, lineOfText(t, content, "session int ses_x"),
		"ses_x", "variable", "definition", "static")
	assertToken(t, toks, lineOfText(t, content, "saved int sav_y"),
		"sav_y", "variable", "declaration", "static")
	assertToken(t, toks, lineOfText(t, content, "hook(int) ev_hook"),
		"ev_hook", "event", "declaration")

	// Keyword used as identifier vs keyword
	regLine := lineOfText(t, content, "register r size 4")
	assertToken(t, toks, regLine, "r", "class", "definition")
	assertToken(t, toks, regLine, "size", "keyword")
	assertToken(t, toks, lineOfText(t, content, "param size = 8"),
		"size", "property", "definition")

	// Method definition, params, locals, types, modifiers, refs
	methLine := lineOfText(t, content, "inline method run")
	assertToken(t, toks, methLine, "inline", "modifier")
	assertToken(t, toks, methLine, "run", "method", "definition")
	assertToken(
		t, toks, methLine,
		"inp", "parameter", "declaration", "readonly",
	)
	assertToken(t, toks, methLine, "outp", "parameter", "declaration")
	assertToken(t, toks, methLine, "throws", "modifier")

	assertToken(t, toks, lineOfText(t, content, "local const int c ="),
		"c", "variable", "definition", "readonly")
	assertToken(t, toks, lineOfText(t, content, "local const int c ="),
		"MAX_SZ", "enumMember", "readonly")
	assertToken(t, toks, lineOfText(t, content, "local int uninit"),
		"uninit", "variable", "declaration")
	assertToken(t, toks, lineOfText(t, content, "static int st ="),
		"st", "variable", "definition", "static")
	assertToken(t, toks, lineOfText(t, content, "unsigned bare_u"),
		"unsigned", "type", "defaultLibrary")
	assertToken(t, toks, lineOfText(t, content, "signed int si"),
		"signed", "modifier")

	msLine := lineOfText(t, content, "local struct my_s ms")
	assertToken(t, toks, msLine, "my_s", "struct")
	assertToken(t, toks, msLine, "f", "property")
	assertToken(t, toks, lineOfText(t, content, "local my_s_t ms2"),
		"my_s_t", "struct")

	uLine := lineOfText(t, content, "local my_u32_t u =")
	assertToken(t, toks, uLine, "my_u32_t", "type")
	assertToken(t, toks, uLine, "inp", "parameter", "readonly")
	assertToken(t, toks, uLine, "outp", "parameter")

	assertToken(t, toks, lineOfText(t, content, "local bool flag"),
		"bool", "type", "defaultLibrary")
	assertToken(t, toks, lineOfText(t, content, "local bool flag"),
		"true", "keyword")
	assertToken(t, toks, lineOfText(t, content, "ext_unknown_t"),
		"ext_unknown_t", "type")

	refLine := lineOfText(t, content, "uninit += c + st")
	assertToken(t, toks, refLine, "uninit", "variable")
	assertToken(t, toks, refLine, "c", "variable", "readonly")
	assertToken(t, toks, refLine, "st", "variable", "static")
	assertToken(t, toks, refLine, "ses_x", "variable", "static")
	assertToken(t, toks, refLine, "sav_y", "variable", "static")
	assertToken(t, toks, refLine, "p_def", "property")
	assertToken(t, toks, refLine, "i", "variable")

	ifLine := lineOfText(t, content, "#if (flag)")
	assertToken(t, toks, ifLine, "#if", "keyword")
	assertToken(t, toks, ifLine, "flag", "variable")

	logLine := lineOfText(t, content, "lg_io:")
	assertToken(t, toks, logLine, "lg_io", "enumMember")
	assertToken(t, toks, logLine, "regs", "class")
	assertToken(t, toks, logLine, "r", "class")
	assertToken(t, toks, logLine, "size", "property")

	elseLogLine := lineOfText(t, content, "other_lg:")
	assertToken(t, toks, elseLogLine, "other_lg", "enumMember")

	forLine := lineOfText(t, content, "for (local int k")
	assertToken(t, toks, forLine, "k", "variable", "definition")
	assertToken(t, toks, lineOfText(t, content, "uninit += k"),
		"k", "variable")

	feLine := lineOfText(t, content, "foreach obj")
	assertToken(t, toks, feLine, "obj", "variable", "definition")
	assertToken(t, toks, feLine, "abs_tpl", "interface")
	assertToken(t, toks, feLine, "dev", "variable", "defaultLibrary")
	assertToken(t, toks, feLine, "regs", "class")

	feBodyLine := lineOfText(t, content, "obj.abs_m(1)")
	assertToken(t, toks, feBodyLine, "obj", "variable")
	assertToken(t, toks, feBodyLine, "abs_m", "method")

	seqLine := lineOfText(t, content, "sequence(abs_tpl)")
	assertToken(t, toks, seqLine, "abs_tpl", "interface")

	afterTimeLine := lineOfText(t, content, "after 0.1 s:")
	assertToken(t, toks, afterTimeLine, "s", "keyword")
	assertToken(t, toks, afterTimeLine, "ext_fn", "function")

	afterMsgLine := lineOfText(t, content, "after ev_hook")
	assertToken(t, toks, afterMsgLine, "ev_hook", "event")
	assertToken(t, toks, afterMsgLine, "mx", "parameter", "declaration")
	assertToken(t, toks, afterMsgLine, "my", "parameter", "declaration")
	assertToken(t, toks, afterMsgLine, "cb_fn", "function")

	// Calls, unresolved members/idents, this/parent, goto label
	callLine := lineOfText(t, content, "unknown_fn(this.p_def")
	assertToken(t, toks, callLine, "unknown_fn", "function")
	assertToken(t, toks, callLine, "this", "variable", "defaultLibrary")
	assertToken(t, toks, callLine, "p_def", "property")
	assertToken(t, toks, callLine, "parent", "variable", "defaultLibrary")

	mCallLine := lineOfText(t, content, "ms.unknown_obj.do_thing()")
	assertToken(t, toks, mCallLine, "unknown_obj", "property")
	assertToken(t, toks, mCallLine, "do_thing", "method")

	assertToken(t, toks, lineOfText(t, content, "ms.missing_prop"),
		"missing_prop", "property")

	gotoLine := lineOfText(t, content, "goto cleanup")
	if _, found := findToken(toks, gotoLine, "unknown_ident"); found {
		t.Errorf(
			"unresolved bare identifier unknown_ident should not be tokenized",
		)
	}
	assertToken(t, toks, gotoLine, "cleanup", "label")
}

func TestSemanticTokensRange(t *testing.T) {
	dir := t.TempDir()
	h := initSyncHarness(t, dir)
	docURI := uri.File(filepath.Join(dir, "range.dml"))

	content := strings.Join([]string{
		"dml 1.4;",
		"template t1 {",
		"    param x = 1;",
		"}",
		"template t2 {",
		"    param y = 2;",
		"}",
	}, "\n") + "\n"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	openAndWait(t, ctx, h, docURI, content)

	res, err := h.client.SemanticTokensRange(
		ctx,
		&protocol.SemanticTokensRangeParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: docURI},
			Range: protocol.Range{
				Start: protocol.Position{Line: 2, Character: 0},
				End:   protocol.Position{Line: 5, Character: 0},
			},
		},
	)
	if err != nil {
		t.Fatalf("SemanticTokensRange failed: %v", err)
	}

	toks := decodeSemanticTokens(t, res.Data, content)
	for _, tok := range toks {
		if tok.line < 2 || tok.line >= 5 {
			t.Errorf("unexpected token outside range: %+v", tok)
		}
	}
	assertToken(t, toks, 2, "x", "property", "definition")
	assertToken(t, toks, 4, "t2", "interface", "definition")

	// Empty range must return empty Data slice without error
	emptyRes, err := h.client.SemanticTokensRange(
		ctx,
		&protocol.SemanticTokensRangeParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: docURI},
			Range: protocol.Range{
				Start: protocol.Position{Line: 2, Character: 4},
				End:   protocol.Position{Line: 2, Character: 4},
			},
		},
	)
	if err != nil {
		t.Fatalf("SemanticTokensRange empty failed: %v", err)
	}
	if len(emptyRes.Data) != 0 {
		t.Errorf("empty range Data len = %d, want 0", len(emptyRes.Data))
	}
}

func TestSemanticTokensUTF16SurrogatePairs(t *testing.T) {
	dir := t.TempDir()
	h := initSyncHarness(t, dir)
	docURI := uri.File(filepath.Join(dir, "utf16.dml"))

	content := "dml 1.4;\nmethod foo() {\n    /* 😀 */ local uint32 val = 7;\n}\n"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	openAndWait(t, ctx, h, docURI, content)

	res, err := h.client.SemanticTokensFull(
		ctx,
		&protocol.SemanticTokensParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: docURI},
		},
	)
	if err != nil {
		t.Fatalf("SemanticTokensFull failed: %v", err)
	}

	toks := decodeSemanticTokens(t, res.Data, content)
	// "    /* 😀 */ local uint32 val = 7;"
	// 😀 takes 2 UTF-16 code units, so "local" starts at col 13.
	assertToken(t, toks, 2, "/* 😀 */", "comment")
	localTok, ok := findToken(toks, 2, "local")
	if !ok {
		t.Fatalf("token 'local' not found on line 2")
	}
	if localTok.char != 13 {
		t.Errorf("UTF-16 col for 'local' = %d, want 13", localTok.char)
	}
	assertToken(t, toks, 2, "uint32", "type", "defaultLibrary")
	assertToken(t, toks, 2, "val", "variable", "definition")
}

func TestSemanticTokensIncrementalEditInvalidation(t *testing.T) {
	dir := t.TempDir()
	h := initSyncHarness(t, dir)
	docURI := uri.File(filepath.Join(dir, "edit.dml"))

	v1 := "dml 1.4;\ntemplate alpha {}\n"
	v2 := "dml 1.4;\nbank beta {}\n"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	openAndWait(t, ctx, h, docURI, v1)

	res1, err := h.client.SemanticTokensFull(
		ctx,
		&protocol.SemanticTokensParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: docURI},
		},
	)
	if err != nil {
		t.Fatalf("SemanticTokensFull v1 failed: %v", err)
	}
	assertToken(t, decodeSemanticTokens(t, res1.Data, v1),
		1, "alpha", "interface", "definition")

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

	res2, err := h.client.SemanticTokensFull(
		ctx,
		&protocol.SemanticTokensParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: docURI},
		},
	)
	if err != nil {
		t.Fatalf("SemanticTokensFull v2 failed: %v", err)
	}
	assertToken(t, decodeSemanticTokens(t, res2.Data, v2),
		1, "beta", "class", "definition")
}

func TestSemanticTokensSyntaxErrorsAndDML12(t *testing.T) {
	dir := t.TempDir()
	h := initSyncHarness(t, dir)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Active syntax error mid-edit
	brokenURI := uri.File(filepath.Join(dir, "broken.dml"))
	broken := "dml 1.4;\nmethod broken( {\n    local uint32 x = 5;\n    // trailing\n"
	openAndWait(t, ctx, h, brokenURI, broken)
	brokenRes, err := h.client.SemanticTokensFull(
		ctx,
		&protocol.SemanticTokensParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: brokenURI},
		},
	)
	if err != nil {
		t.Fatalf("SemanticTokensFull broken failed: %v", err)
	}
	brokenToks := decodeSemanticTokens(t, brokenRes.Data, broken)
	assertToken(t, brokenToks, 1, "method", "keyword")
	assertToken(t, brokenToks, 2, "uint32", "type", "defaultLibrary")
	assertToken(t, brokenToks, 2, "5", "number")
	assertToken(t, brokenToks, 3, "// trailing", "comment")

	// DML 1.2 falls back cleanly to lexer baseline
	dml12URI := uri.File(filepath.Join(dir, "legacy.dml"))
	dml12 := strings.Join([]string{
		"dml 1.2;",
		"// legacy file",
		"method legacy() {",
		"    static uint32 cnt = 10;",
		"    log \"info\": \"ok\";",
		"}",
	}, "\n") + "\n"
	openAndWait(t, ctx, h, dml12URI, dml12)
	dml12Res, err := h.client.SemanticTokensFull(
		ctx,
		&protocol.SemanticTokensParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: dml12URI},
		},
	)
	if err != nil {
		t.Fatalf("SemanticTokensFull dml12 failed: %v", err)
	}
	dml12Toks := decodeSemanticTokens(t, dml12Res.Data, dml12)
	assertToken(t, dml12Toks, 0, "dml", "keyword")
	assertToken(t, dml12Toks, 0, "1.2", "number")
	assertToken(t, dml12Toks, 1, "// legacy file", "comment")
	assertToken(t, dml12Toks, 3, "static", "modifier")
	assertToken(t, dml12Toks, 3, "uint32", "type", "defaultLibrary")
	assertToken(t, dml12Toks, 4, "\"ok\"", "string")
}

func TestSemanticTokensEndianFlags(t *testing.T) {
	dir := t.TempDir()
	h := initSyncHarness(t, dir)
	docURI := uri.File(filepath.Join(dir, "endian.dml"))

	content := strings.Join([]string{
		"dml 1.4;",
		"method test_endian(uint32 v) {",
		"    local uint16 b = v[15:0, be];",
		"    local uint16 l = v[15:0, le];",
		"}",
	}, "\n") + "\n"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	openAndWait(t, ctx, h, docURI, content)

	res, err := h.client.SemanticTokensFull(ctx, &protocol.SemanticTokensParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: docURI},
	})
	if err != nil {
		t.Fatalf("SemanticTokensFull failed: %v", err)
	}

	toks := decodeSemanticTokens(t, res.Data, content)
	assertToken(t, toks, lineOfText(t, content, "v[15:0, be]"), "be", "keyword")
	assertToken(t, toks, lineOfText(t, content, "v[15:0, le]"), "le", "keyword")
}

func TestSemanticTokensHookParams(t *testing.T) {
	dir := t.TempDir()
	h := initSyncHarness(t, dir)
	docURI := uri.File(filepath.Join(dir, "hook.dml"))

	content := strings.Join([]string{
		"dml 1.4;",
		"hook(int a, const bool b) my_hook;",
	}, "\n") + "\n"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	openAndWait(t, ctx, h, docURI, content)

	res, err := h.client.SemanticTokensFull(ctx, &protocol.SemanticTokensParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: docURI},
	})
	if err != nil {
		t.Fatalf("SemanticTokensFull failed: %v", err)
	}

	toks := decodeSemanticTokens(t, res.Data, content)
	line := lineOfText(t, content, "my_hook")
	assertToken(t, toks, line, "my_hook", "event", "declaration")
	assertToken(t, toks, line, "a", "parameter", "declaration")
	assertToken(t, toks, line, "b", "parameter", "declaration", "readonly")
}

func TestSemanticTokensBareSignedUnsignedFunc(t *testing.T) {
	dir := t.TempDir()
	h := initSyncHarness(t, dir)
	docURI := uri.File(filepath.Join(dir, "funcs.dml"))

	content := strings.Join([]string{
		"dml 1.4;",
		"extern unsigned my_u_fn(int x);",
		"extern signed my_s_fn(int y);",
	}, "\n") + "\n"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	openAndWait(t, ctx, h, docURI, content)

	res, err := h.client.SemanticTokensFull(ctx, &protocol.SemanticTokensParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: docURI},
	})
	if err != nil {
		t.Fatalf("SemanticTokensFull failed: %v", err)
	}

	toks := decodeSemanticTokens(t, res.Data, content)
	uLine := lineOfText(t, content, "my_u_fn")
	assertToken(t, toks, uLine, "unsigned", "type", "defaultLibrary")
	assertToken(t, toks, uLine, "my_u_fn", "function", "declaration")
	assertToken(t, toks, uLine, "x", "parameter", "declaration")

	sLine := lineOfText(t, content, "my_s_fn")
	assertToken(t, toks, sLine, "signed", "type", "defaultLibrary")
	assertToken(t, toks, sLine, "my_s_fn", "function", "declaration")
	assertToken(t, toks, sLine, "y", "parameter", "declaration")
}

func TestSemanticTokensDevInFileWithoutDevice(t *testing.T) {
	dir := t.TempDir()
	h := initSyncHarness(t, dir)
	docURI := uri.File(filepath.Join(dir, "no_dev.dml"))

	content := strings.Join([]string{
		"dml 1.4;",
		"template tpl {",
		"    param p_val = 10;",
		"}",
		"is tpl;",
		"bank regs {",
		"    register r size 4;",
		"}",
		"method test_dev() {",
		"    local int a = p_val;",
		"    local int b = dev.regs.r.size;",
		"}",
	}, "\n") + "\n"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	openAndWait(t, ctx, h, docURI, content)

	res, err := h.client.SemanticTokensFull(ctx, &protocol.SemanticTokensParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: docURI},
	})
	if err != nil {
		t.Fatalf("SemanticTokensFull failed: %v", err)
	}

	toks := decodeSemanticTokens(t, res.Data, content)
	aLine := lineOfText(t, content, "local int a = p_val")
	assertToken(t, toks, aLine, "p_val", "property")

	bLine := lineOfText(t, content, "dev.regs.r.size")
	assertToken(t, toks, bLine, "dev", "variable", "defaultLibrary")
	assertToken(t, toks, bLine, "regs", "class")
	assertToken(t, toks, bLine, "r", "class")
}

func TestSemanticTokensConcurrentAccess(t *testing.T) {
	dir := t.TempDir()
	h := initSyncHarness(t, dir)
	docURI := uri.File(filepath.Join(dir, "concurrent.dml"))

	content := strings.Join([]string{
		"dml 1.4;",
		"bank b {",
		"    register r size 4;",
		"}",
	}, "\n") + "\n"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	openAndWait(t, ctx, h, docURI, content)

	const goroutines = 8
	errCh := make(chan error, goroutines*2)

	for range goroutines {
		go func() {
			_, err := h.client.SemanticTokensFull(
				ctx,
				&protocol.SemanticTokensParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: docURI,
					},
				},
			)
			errCh <- err
		}()
		go func() {
			_, err := h.client.SemanticTokensRange(
				ctx,
				&protocol.SemanticTokensRangeParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: docURI,
					},
					Range: protocol.Range{
						Start: protocol.Position{Line: 1, Character: 0},
						End:   protocol.Position{Line: 3, Character: 0},
					},
				},
			)
			errCh <- err
		}()
	}

	for range goroutines * 2 {
		if err := <-errCh; err != nil {
			t.Errorf("concurrent request failed: %v", err)
		}
	}
}
