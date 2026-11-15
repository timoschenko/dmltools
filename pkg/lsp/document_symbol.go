// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lsp

import (
	"context"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"go.lsp.dev/protocol"
)

// ComputeDocumentSymbols computes hierarchical LSP document symbols for the
// file identified by uriStr.
func ComputeDocumentSymbols(
	ctx context.Context,
	fsCache *FsCache,
	uriStr string,
) (protocol.DocumentSymbolSlice, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	symbols, err := fsCache.getOrComputeDocumentSymbols(uriStr)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if symbols == nil {
		return protocol.DocumentSymbolSlice{}, nil
	}
	return protocol.DocumentSymbolSlice(symbols), nil
}

func normalizeRange(r protocol.Range) protocol.Range {
	if positionLess(r.End, r.Start) {
		r.End = r.Start
	}
	return r
}

func ensureRangeContains(
	outer protocol.Range,
	inner protocol.Range,
) protocol.Range {
	if positionLess(inner.Start, outer.Start) {
		outer.Start = inner.Start
	}
	if positionLess(outer.End, inner.End) {
		outer.End = inner.End
	}
	return outer
}

func sanitizeDocumentSymbol(sym *protocol.DocumentSymbol) {
	if sym == nil {
		return
	}
	sym.Range = normalizeRange(sym.Range)
	sym.SelectionRange = normalizeRange(sym.SelectionRange)

	if sym.SelectionRange.Start == sym.SelectionRange.End &&
		sym.Range.Start != sym.Range.End {
		sym.SelectionRange = sym.Range
	} else {
		sym.Range = ensureRangeContains(sym.Range, sym.SelectionRange)
	}

	for i := range sym.Children {
		sanitizeDocumentSymbol(&sym.Children[i])
		sym.Range = ensureRangeContains(sym.Range, sym.Children[i].Range)
	}
}

// ExtractDocumentSymbols walks a parsed DML AST and extracts the document
// symbol hierarchy.
func ExtractDocumentSymbols(tree antlr.ParseTree) []protocol.DocumentSymbol {
	if tree == nil {
		return nil
	}
	if isDML12Tree(tree) {
		return []protocol.DocumentSymbol{}
	}

	dmlCtx, ok := tree.(*dmlex.DmlContext)
	if !ok || dmlCtx.DmlCode() == nil {
		return []protocol.DocumentSymbol{}
	}

	codeCtx := dmlCtx.DmlCode().(*dmlex.DmlCodeContext)
	rootSymbols := make([]protocol.DocumentSymbol, 0)

	if codeCtx.DeviceTag() != nil {
		if dt, ok := codeCtx.DeviceTag().(*dmlex.DeviceTagContext); ok {
			name, nameRange := extractObjident(dt.Objident())
			if name != "" && name != "_" {
				sym := protocol.DocumentSymbol{
					Name:           name,
					Detail:         new("device"),
					Kind:           protocol.SymbolKindClass,
					Range:          ruleRange(dt),
					SelectionRange: nameRange,
				}
				sanitizeDocumentSymbol(&sym)
				rootSymbols = append(rootSymbols, sym)
			}
		}
	}

	if codeCtx.DeviceStatements() != nil {
		if ds, ok := codeCtx.DeviceStatements().(*dmlex.DeviceStatementsContext); ok {
			rootSymbols = append(
				rootSymbols,
				extractDeviceStatements(ds.AllDeviceStatement())...)
		}
	}

	return rootSymbols
}

func extractDeviceStatements(
	stmts []dmlex.IDeviceStatementContext,
) []protocol.DocumentSymbol {
	symbols := make([]protocol.DocumentSymbol, 0, len(stmts))
	for _, stmt := range stmts {
		symbols = append(symbols, extractDeviceStatement(stmt)...)
	}
	return symbols
}

func extractDeviceStatement(
	stmt dmlex.IDeviceStatementContext,
) []protocol.DocumentSymbol {
	if stmt == nil {
		return nil
	}
	dsCtx, ok := stmt.(*dmlex.DeviceStatementContext)
	if !ok {
		return nil
	}

	switch {
	case dsCtx.TemplateStatement() != nil:
		if ts, ok := dsCtx.TemplateStatement().(*dmlex.TemplateStatementContext); ok {
			return extractTemplateStatement(ts)
		}
	case dsCtx.Object() != nil:
		return extractObject(dsCtx.Object())
	case dsCtx.Param() != nil:
		if p, ok := dsCtx.Param().(*dmlex.ParamContext); ok {
			return extractParam(p)
		}
	case dsCtx.MethodUnified() != nil:
		return extractMethodUnified(dsCtx.MethodUnified())
	case dsCtx.ConstantStatement() != nil:
		if cs, ok := dsCtx.ConstantStatement().(*dmlex.ConstantStatementContext); ok {
			return extractConstantStatement(cs)
		}
	case dsCtx.LogGroupStatement() != nil:
		if lgs, ok := dsCtx.LogGroupStatement().(*dmlex.LogGroupStatementContext); ok {
			return extractLogGroupStatement(lgs)
		}
	case dsCtx.TypedefStatement() != nil:
		if tds, ok := dsCtx.TypedefStatement().(*dmlex.TypedefStatementContext); ok {
			return extractTypedefStatement(tds)
		}
	case dsCtx.ExternStatement() != nil:
		if es, ok := dsCtx.ExternStatement().(*dmlex.ExternStatementContext); ok {
			return extractExternStatement(es)
		}
	case dsCtx.IneachStatement() != nil:
		if ie, ok := dsCtx.IneachStatement().(*dmlex.IneachStatementContext); ok {
			return extractIneachStatement(ie)
		}
	case dsCtx.ToplevelIf() != nil:
		if tif, ok := dsCtx.ToplevelIf().(*dmlex.ToplevelIfContext); ok {
			return extractToplevelIf(tif)
		}
	}
	return nil
}

func extractTemplateStatement(
	ctx *dmlex.TemplateStatementContext,
) []protocol.DocumentSymbol {
	name, nameRange := extractObjident(ctx.Objident())
	if name == "" || name == "_" {
		return nil
	}

	var children []protocol.DocumentSymbol
	for _, body := range ctx.AllTemplateBody() {
		tb, ok := body.(*dmlex.TemplateBodyContext)
		if !ok {
			continue
		}
		if tb.ObjectStatement() != nil {
			children = append(
				children,
				extractObjectStatement(tb.ObjectStatement())...)
		} else if tb.Hook_decl() != nil {
			if hd, ok := tb.Hook_decl().(*dmlex.Hook_declContext); ok {
				children = append(children, extractHookDecl(hd)...)
			}
		}
	}

	sym := protocol.DocumentSymbol{
		Name:           name,
		Detail:         new("template"),
		Kind:           protocol.SymbolKindInterface,
		Range:          ruleRange(ctx),
		SelectionRange: nameRange,
		Children:       children,
	}
	sanitizeDocumentSymbol(&sym)
	return []protocol.DocumentSymbol{sym}
}

func extractIneachStatement(
	ctx *dmlex.IneachStatementContext,
) []protocol.DocumentSymbol {
	name, nameRange := formatIneachName(ctx)
	children := extractObjectStatements(ctx.AllObjectStatement())
	sym := protocol.DocumentSymbol{
		Name:           name,
		Detail:         new("in each"),
		Kind:           protocol.SymbolKindNamespace,
		Range:          ruleRange(ctx),
		SelectionRange: nameRange,
		Children:       children,
	}
	sanitizeDocumentSymbol(&sym)
	return []protocol.DocumentSymbol{sym}
}

func formatIneachName(
	ctx *dmlex.IneachStatementContext,
) (string, protocol.Range) {
	if ctx.IstemplateList() == nil {
		return "in each", ruleRange(ctx)
	}
	nameRange := ruleRange(ctx.IstemplateList())
	il, ok := ctx.IstemplateList().(*dmlex.IstemplateListContext)
	if !ok {
		return "in each " + ctx.IstemplateList().GetText(), nameRange
	}
	var names []string
	for _, oid := range il.AllObjident() {
		name, _ := extractObjident(oid)
		if name != "" {
			names = append(names, name)
		}
	}
	if il.LPAREN() != nil || len(names) > 1 {
		return "in each (" + strings.Join(names, ", ") + ")", nameRange
	}
	if len(names) == 1 {
		return "in each " + names[0], nameRange
	}
	return "in each " + ctx.IstemplateList().GetText(), nameRange
}

func extractObject(obj dmlex.IObjectContext) []protocol.DocumentSymbol {
	if obj == nil {
		return nil
	}
	objCtx, ok := obj.(*dmlex.ObjectContext)
	if !ok {
		return nil
	}

	switch {
	case objCtx.Object_decl() != nil:
		if od, ok := objCtx.Object_decl().(*dmlex.Object_declContext); ok {
			return extractObjectDecl(od)
		}
	case objCtx.SessionOrSavedStatement() != nil:
		if sss, ok := objCtx.SessionOrSavedStatement().(*dmlex.SessionOrSavedStatementContext); ok {
			return extractSessionOrSavedStatement(sss)
		}
	case objCtx.Hook_decl() != nil:
		if hd, ok := objCtx.Hook_decl().(*dmlex.Hook_declContext); ok {
			return extractHookDecl(hd)
		}
	}
	return nil
}

func extractObjectDecl(
	ctx *dmlex.Object_declContext,
) []protocol.DocumentSymbol {
	name, nameRange := extractObjident(ctx.Objident())
	if name == "" || name == "_" {
		return nil
	}

	detail := "object"
	switch {
	case ctx.REGISTER() != nil:
		detail = "register"
	case ctx.FIELD() != nil:
		detail = "field"
	case ctx.Object_type() != nil:
		detail = strings.ToLower(ctx.Object_type().GetText())
	}

	var children []protocol.DocumentSymbol
	if ctx.ObjectBody() != nil {
		if ob, ok := ctx.ObjectBody().(*dmlex.ObjectBodyContext); ok {
			children = extractObjectStatements(ob.AllObjectStatement())
		}
	}

	sym := protocol.DocumentSymbol{
		Name:           name,
		Detail:         new(detail),
		Kind:           protocol.SymbolKindClass,
		Range:          ruleRange(ctx),
		SelectionRange: nameRange,
		Children:       children,
	}
	sanitizeDocumentSymbol(&sym)
	return []protocol.DocumentSymbol{sym}
}

func extractObjectStatements(
	stmts []dmlex.IObjectStatementContext,
) []protocol.DocumentSymbol {
	symbols := make([]protocol.DocumentSymbol, 0, len(stmts))
	for _, stmt := range stmts {
		symbols = append(symbols, extractObjectStatement(stmt)...)
	}
	return symbols
}

func extractObjectStatement(
	stmt dmlex.IObjectStatementContext,
) []protocol.DocumentSymbol {
	if stmt == nil {
		return nil
	}
	osCtx, ok := stmt.(*dmlex.ObjectStatementContext)
	if !ok {
		return nil
	}

	switch {
	case osCtx.Object() != nil:
		return extractObject(osCtx.Object())
	case osCtx.Param() != nil:
		if p, ok := osCtx.Param().(*dmlex.ParamContext); ok {
			return extractParam(p)
		}
	case osCtx.MethodUnified() != nil:
		return extractMethodUnified(osCtx.MethodUnified())
	case osCtx.IneachStatement() != nil:
		if ie, ok := osCtx.IneachStatement().(*dmlex.IneachStatementContext); ok {
			return extractIneachStatement(ie)
		}
	case osCtx.ObjectIf() != nil:
		if oif, ok := osCtx.ObjectIf().(*dmlex.ObjectIfContext); ok {
			return extractObjectIf(oif)
		}
	}
	return nil
}

func extractParam(ctx *dmlex.ParamContext) []protocol.DocumentSymbol {
	name, nameRange := extractObjident(ctx.Objident())
	if name == "" || name == "_" {
		return nil
	}

	var children []protocol.DocumentSymbol
	if ctx.Ctypedecl() != nil {
		if ct, ok := ctx.Ctypedecl().(*dmlex.CtypedeclContext); ok &&
			ct.Basetype() != nil {
			children = extractCompositeMembers(ct.Basetype())
		}
	}

	sym := protocol.DocumentSymbol{
		Name:           name,
		Detail:         new("param"),
		Kind:           protocol.SymbolKindProperty,
		Range:          ruleRange(ctx),
		SelectionRange: nameRange,
		Children:       children,
	}
	sanitizeDocumentSymbol(&sym)
	return []protocol.DocumentSymbol{sym}
}

func extractMethodUnified(
	m dmlex.IMethodUnifiedContext,
) []protocol.DocumentSymbol {
	if m == nil {
		return nil
	}
	name := ""
	var nameRange protocol.Range
	switch ctx := m.(type) {
	case *dmlex.MethodUnifiedDeclarationContext:
		if ctx.Objident() != nil {
			name, nameRange = extractObjident(ctx.Objident())
		}
	case *dmlex.MethodUnifiedDefinitionContext:
		if ctx.Objident() != nil {
			name, nameRange = extractObjident(ctx.Objident())
		}
	}
	if name == "" || name == "_" {
		return nil
	}

	sym := protocol.DocumentSymbol{
		Name:           name,
		Detail:         new("method"),
		Kind:           protocol.SymbolKindMethod,
		Range:          ruleRange(m),
		SelectionRange: nameRange,
	}
	sanitizeDocumentSymbol(&sym)
	return []protocol.DocumentSymbol{sym}
}

func extractConstantStatement(
	ctx *dmlex.ConstantStatementContext,
) []protocol.DocumentSymbol {
	name, nameRange := extractIdent(ctx.Ident())
	if name == "" || name == "_" {
		return nil
	}

	sym := protocol.DocumentSymbol{
		Name:           name,
		Detail:         new("constant"),
		Kind:           protocol.SymbolKindConstant,
		Range:          ruleRange(ctx),
		SelectionRange: nameRange,
	}
	sanitizeDocumentSymbol(&sym)
	return []protocol.DocumentSymbol{sym}
}

func extractLogGroupStatement(
	ctx *dmlex.LogGroupStatementContext,
) []protocol.DocumentSymbol {
	name, nameRange := extractIdent(ctx.Ident())
	if name == "" || name == "_" {
		return nil
	}

	sym := protocol.DocumentSymbol{
		Name:           name,
		Detail:         new("loggroup"),
		Kind:           protocol.SymbolKindEnumMember,
		Range:          ruleRange(ctx),
		SelectionRange: nameRange,
	}
	sanitizeDocumentSymbol(&sym)
	return []protocol.DocumentSymbol{sym}
}

func extractHookDecl(ctx *dmlex.Hook_declContext) []protocol.DocumentSymbol {
	name, nameRange := extractIdent(ctx.Ident())
	if name == "" || name == "_" {
		return nil
	}

	sym := protocol.DocumentSymbol{
		Name:           name,
		Detail:         new("hook"),
		Kind:           protocol.SymbolKindEvent,
		Range:          ruleRange(ctx),
		SelectionRange: nameRange,
	}
	sanitizeDocumentSymbol(&sym)
	return []protocol.DocumentSymbol{sym}
}

func extractSessionOrSavedStatement(
	ctx *dmlex.SessionOrSavedStatementContext,
) []protocol.DocumentSymbol {
	if ctx.SessionOrSavedDecl() == nil {
		return nil
	}
	decl, ok := ctx.SessionOrSavedDecl().(*dmlex.SessionOrSavedDeclContext)
	if !ok || decl.CvarDecl() == nil {
		return nil
	}

	detail := "session"
	if decl.Global_decl_kind() != nil &&
		decl.Global_decl_kind().SAVED() != nil {
		detail = "saved"
	}

	return extractCvarDecl(decl.CvarDecl(), detail, ruleRange(ctx))
}

func extractCvarDecl(
	cvar dmlex.ICvarDeclContext,
	detail string,
	stmtRange protocol.Range,
) []protocol.DocumentSymbol {
	if cvar == nil {
		return nil
	}

	var cdecls []dmlex.ICdeclContext
	if cvar.Cdecl() != nil {
		cdecls = append(cdecls, cvar.Cdecl())
	}
	if cvar.Cdecl_list() != nil {
		cdecls = append(cdecls, cvar.Cdecl_list().AllCdecl()...)
	}

	var symbols []protocol.DocumentSymbol
	for _, cd := range cdecls {
		if cd == nil {
			continue
		}
		name, nameRange := findCdeclIdent(cd)
		if name == "" || name == "_" {
			continue
		}

		kind := protocol.SymbolKindVariable
		var children []protocol.DocumentSymbol
		if isDirectFuncCdecl(cd.Cdecl_qualifier_ident_array_func()) {
			kind = protocol.SymbolKindFunction
		} else if isStructLikeBasetype(cd.Basetype()) {
			kind = protocol.SymbolKindStruct
			children = extractCompositeMembers(cd.Basetype())
		}

		sym := protocol.DocumentSymbol{
			Name:           name,
			Detail:         new(detail),
			Kind:           kind,
			Range:          stmtRange,
			SelectionRange: nameRange,
			Children:       children,
		}
		sanitizeDocumentSymbol(&sym)
		symbols = append(symbols, sym)
	}
	return symbols
}

func extractTypedefStatement(
	ctx *dmlex.TypedefStatementContext,
) []protocol.DocumentSymbol {
	if ctx.Cdecl() == nil {
		return nil
	}
	cd := ctx.Cdecl()
	name, nameRange := findCdeclIdent(cd)
	if name == "" || name == "_" {
		return nil
	}

	kind := protocol.SymbolKindTypeParameter
	var children []protocol.DocumentSymbol
	if isStructLikeBasetype(cd.Basetype()) {
		kind = protocol.SymbolKindStruct
		children = extractCompositeMembers(cd.Basetype())
	}

	sym := protocol.DocumentSymbol{
		Name:           name,
		Detail:         new("typedef"),
		Kind:           kind,
		Range:          ruleRange(ctx),
		SelectionRange: nameRange,
		Children:       children,
	}
	sanitizeDocumentSymbol(&sym)
	return []protocol.DocumentSymbol{sym}
}

func extractExternStatement(
	ctx *dmlex.ExternStatementContext,
) []protocol.DocumentSymbol {
	if ctx.Cdecl() == nil {
		return nil
	}
	cd := ctx.Cdecl()
	name, nameRange := findCdeclIdent(cd)
	if name == "" || name == "_" {
		return nil
	}

	isTypedef := ctx.TYPEDEF() != nil
	detail := "extern"
	if isTypedef {
		detail = "extern typedef"
	}

	kind := protocol.SymbolKindVariable
	var children []protocol.DocumentSymbol
	if isTypedef {
		kind = protocol.SymbolKindTypeParameter
		if isStructLikeBasetype(cd.Basetype()) {
			kind = protocol.SymbolKindStruct
			children = extractCompositeMembers(cd.Basetype())
		}
	} else if isDirectFuncCdecl(cd.Cdecl_qualifier_ident_array_func()) {
		kind = protocol.SymbolKindFunction
	} else if isStructLikeBasetype(cd.Basetype()) {
		kind = protocol.SymbolKindStruct
		children = extractCompositeMembers(cd.Basetype())
	}

	sym := protocol.DocumentSymbol{
		Name:           name,
		Detail:         new(detail),
		Kind:           kind,
		Range:          ruleRange(ctx),
		SelectionRange: nameRange,
		Children:       children,
	}
	sanitizeDocumentSymbol(&sym)
	return []protocol.DocumentSymbol{sym}
}

func extractToplevelIf(ctx *dmlex.ToplevelIfContext) []protocol.DocumentSymbol {
	var symbols []protocol.DocumentSymbol
	if ctx.ToplevelIfStatement() != nil {
		symbols = append(
			symbols,
			extractToplevelIfStatement(ctx.ToplevelIfStatement())...)
	}
	if ctx.ToplevelElseStatement() != nil {
		symbols = append(
			symbols,
			extractToplevelElseStatement(ctx.ToplevelElseStatement())...)
	}
	return symbols
}

func extractToplevelIfStatement(
	stmt dmlex.IToplevelIfStatementContext,
) []protocol.DocumentSymbol {
	if stmt == nil {
		return nil
	}
	tis, ok := stmt.(*dmlex.ToplevelIfStatementContext)
	if !ok || tis.DeviceStatements() == nil {
		return nil
	}
	ds, ok := tis.DeviceStatements().(*dmlex.DeviceStatementsContext)
	if !ok {
		return nil
	}
	return extractDeviceStatements(ds.AllDeviceStatement())
}

func extractToplevelElseStatement(
	stmt dmlex.IToplevelElseStatementContext,
) []protocol.DocumentSymbol {
	if stmt == nil {
		return nil
	}
	tes, ok := stmt.(*dmlex.ToplevelElseStatementContext)
	if !ok {
		return nil
	}
	if tes.DeviceStatements() != nil {
		ds, ok := tes.DeviceStatements().(*dmlex.DeviceStatementsContext)
		if ok {
			return extractDeviceStatements(ds.AllDeviceStatement())
		}
	} else if tes.ToplevelIf() != nil {
		tif, ok := tes.ToplevelIf().(*dmlex.ToplevelIfContext)
		if ok {
			return extractToplevelIf(tif)
		}
	}
	return nil
}

func extractObjectIf(ctx *dmlex.ObjectIfContext) []protocol.DocumentSymbol {
	var symbols []protocol.DocumentSymbol
	if ctx.ObjectIfStatement() != nil {
		symbols = append(
			symbols,
			extractObjectIfStatement(ctx.ObjectIfStatement())...)
	}
	if ctx.ObjectElseStatement() != nil {
		symbols = append(
			symbols,
			extractObjectElseStatement(ctx.ObjectElseStatement())...)
	}
	return symbols
}

func extractObjectIfStatement(
	stmt dmlex.IObjectIfStatementContext,
) []protocol.DocumentSymbol {
	if stmt == nil {
		return nil
	}
	ois, ok := stmt.(*dmlex.ObjectIfStatementContext)
	if !ok {
		return nil
	}
	return extractObjectStatements(ois.AllObjectStatement())
}

func extractObjectElseStatement(
	stmt dmlex.IObjectElseStatementContext,
) []protocol.DocumentSymbol {
	if stmt == nil {
		return nil
	}
	oes, ok := stmt.(*dmlex.ObjectElseStatementContext)
	if !ok {
		return nil
	}
	return extractObjectStatements(oes.AllObjectStatement())
}

func extractCompositeMembers(
	base dmlex.IBasetypeContext,
) []protocol.DocumentSymbol {
	if base == nil {
		return nil
	}

	switch {
	case base.Struct_() != nil:
		st, ok := base.Struct_().(*dmlex.StructContext)
		if ok {
			return extractStructDecls(st.AllStructDecls())
		}
	case base.Union() != nil:
		un, ok := base.Union().(*dmlex.UnionContext)
		if ok {
			return extractStructDecls(un.AllStructDecls())
		}
	case base.Layout() != nil:
		lay, ok := base.Layout().(*dmlex.LayoutContext)
		if ok {
			return extractLayoutDecls(lay.AllLayoutDecls())
		}
	case base.Bitfields() != nil:
		bf, ok := base.Bitfields().(*dmlex.BitfieldsContext)
		if ok {
			return extractBitfieldsDecls(bf.AllBitfields_decls())
		}
	}
	return nil
}

func extractStructDecls(
	decls []dmlex.IStructDeclsContext,
) []protocol.DocumentSymbol {
	var symbols []protocol.DocumentSymbol
	for _, d := range decls {
		sd, ok := d.(*dmlex.StructDeclsContext)
		if !ok || sd.Cdecl() == nil {
			continue
		}
		symbols = append(
			symbols,
			extractMemberField(sd.Cdecl(), ruleRange(sd))...)
	}
	return symbols
}

func extractLayoutDecls(
	decls []dmlex.ILayoutDeclsContext,
) []protocol.DocumentSymbol {
	var symbols []protocol.DocumentSymbol
	for _, d := range decls {
		ld, ok := d.(*dmlex.LayoutDeclsContext)
		if !ok || ld.Cdecl() == nil {
			continue
		}
		symbols = append(
			symbols,
			extractMemberField(ld.Cdecl(), ruleRange(ld))...)
	}
	return symbols
}

func extractBitfieldsDecls(
	decls []dmlex.IBitfields_declsContext,
) []protocol.DocumentSymbol {
	var symbols []protocol.DocumentSymbol
	for _, d := range decls {
		bd, ok := d.(*dmlex.Bitfields_declsContext)
		if !ok || bd.Cdecl() == nil {
			continue
		}
		symbols = append(
			symbols,
			extractMemberField(bd.Cdecl(), ruleRange(bd))...)
	}
	return symbols
}

func extractArrayBrackets(node dmlex.ICdecl_ident_array_funcContext) string {
	if node == nil {
		return ""
	}
	if arr, ok := node.(*dmlex.CdeclArrayContext); ok {
		inner := extractArrayBrackets(arr.Cdecl_ident_array_func())
		exprText := ""
		if arr.Expression() != nil {
			exprText = arr.Expression().GetText()
		}
		return inner + "[" + exprText + "]"
	}
	return ""
}

func formatCdeclType(cd dmlex.ICdeclContext) string {
	if cd == nil || cd.Basetype() == nil {
		return ""
	}
	var b strings.Builder
	if cd.CONST() != nil {
		b.WriteString("const ")
	}

	base := cd.Basetype()
	if su := base.Signed_or_unsigned_typeident(); su != nil {
		sign := ""
		if su.UNSIGNED() != nil {
			sign = "unsigned"
		} else if su.SIGNED() != nil {
			sign = "signed"
		}

		identInType := false
		if (su.SIGNED() != nil || su.UNSIGNED() != nil) &&
			su.Typeident() != nil && su.Typeident().Ident() != nil &&
			findIdentTokenInTree(cd.Cdecl_qualifier_ident_array_func()) == nil {
			identInType = true
		}

		if identInType {
			b.WriteString(sign)
		} else {
			if sign != "" {
				b.WriteString(sign)
				if su.Typeident() != nil {
					b.WriteString(" ")
				}
			}
			if su.Typeident() != nil {
				b.WriteString(su.Typeident().GetText())
			}
		}
	} else {
		b.WriteString(base.GetText())
	}

	q := cd.Cdecl_qualifier_ident_array_func()
	if q != nil {
		for _, tq := range q.AllTypeQualifier() {
			if tq.CONST() != nil {
				b.WriteString("const")
			} else if tq.MUL() != nil {
				b.WriteString("*")
			} else if tq.VECT() != nil {
				b.WriteString("vect ")
			}
		}
		b.WriteString(extractArrayBrackets(q.Cdecl_ident_array_func()))
	}

	return strings.TrimSpace(b.String())
}

func formatFieldType(cd dmlex.ICdeclContext) string {
	if cd == nil || cd.Basetype() == nil {
		return ""
	}
	b := cd.Basetype()
	switch {
	case b.Struct_() != nil:
		st, ok := b.Struct_().(*dmlex.StructContext)
		if ok && st.Ident() != nil {
			return "struct " + st.Ident().GetText()
		}
		return "struct"
	case b.Union() != nil:
		un, ok := b.Union().(*dmlex.UnionContext)
		if ok && un.Ident() != nil {
			return "union " + un.Ident().GetText()
		}
		return "union"
	case b.Layout() != nil:
		return "layout"
	case b.Bitfields() != nil:
		return "bitfields"
	case b.Typeof() != nil:
		return "typeof"
	default:
		return formatCdeclType(cd)
	}
}

func extractMemberField(
	cd dmlex.ICdeclContext,
	fieldRange protocol.Range,
) []protocol.DocumentSymbol {
	name, nameRange := findCdeclIdent(cd)
	if name == "" || name == "_" {
		return nil
	}

	detail := formatFieldType(cd)
	kind := protocol.SymbolKindField
	var children []protocol.DocumentSymbol
	if isStructLikeBasetype(cd.Basetype()) {
		kind = protocol.SymbolKindStruct
		children = extractCompositeMembers(cd.Basetype())
	}

	sym := protocol.DocumentSymbol{
		Name:           name,
		Detail:         new(detail),
		Kind:           kind,
		Range:          fieldRange,
		SelectionRange: nameRange,
		Children:       children,
	}
	sanitizeDocumentSymbol(&sym)
	return []protocol.DocumentSymbol{sym}
}
