// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/timoschenko/dmltools/pkg/dmlex"
)

// ExtractOptions configures AST symbol extraction.
type ExtractOptions struct {
	// EmitQualified emits additional qualified tag names (<scope>.<name>)
	// alongside unqualified tags.
	EmitQualified bool
}

type scopeEntry struct {
	kind            Kind
	name            string
	path            string
	tagIndices      []int
	inEachTemplates []string
}

type extractor struct {
	filePath      string
	sourceLines   []string
	opts          ExtractOptions
	scopes        []scopeEntry
	tags          []Tag
	deviceTagSeen bool
}

func (e *extractor) currentScope() *scopeEntry {
	if len(e.scopes) == 0 {
		return nil
	}
	return &e.scopes[len(e.scopes)-1]
}

func (e *extractor) currentScopeKind() Kind {
	if s := e.currentScope(); s != nil {
		return s.kind
	}
	return 0
}

func (e *extractor) currentScopePath() string {
	if s := e.currentScope(); s != nil {
		return s.path
	}
	return ""
}

func (e *extractor) currentInEachTemplates() []string {
	if s := e.currentScope(); s != nil {
		return s.inEachTemplates
	}
	return nil
}

func (e *extractor) addCurrentScopeInherits(templates []string) {
	if len(templates) == 0 {
		return
	}
	for i := len(e.scopes) - 1; i >= 0; i-- {
		s := &e.scopes[i]
		if len(s.tagIndices) > 0 {
			for _, idx := range s.tagIndices {
				if idx >= 0 && idx < len(e.tags) {
					e.tags[idx].Inherits = append(
						e.tags[idx].Inherits,
						templates...,
					)
				}
			}
			return
		}
	}
}

func (e *extractor) pushScope(
	kind Kind,
	name string,
	inEachTemplates []string,
	tagIndices []int,
) {
	parentPath := e.currentScopePath()
	fullPath := name
	if parentPath != "" && name != "" {
		fullPath = parentPath + "." + name
	} else if parentPath != "" && name == "" {
		fullPath = parentPath
	}

	mergedInEach := inEachTemplates
	if len(inEachTemplates) == 0 && len(e.scopes) > 0 {
		// Inherit in-each templates from parent scope if any
		mergedInEach = e.currentScope().inEachTemplates
	}

	e.scopes = append(e.scopes, scopeEntry{
		kind:            kind,
		name:            name,
		path:            fullPath,
		tagIndices:      tagIndices,
		inEachTemplates: mergedInEach,
	})
}

func (e *extractor) popScope() {
	if len(e.scopes) > 0 {
		e.scopes = e.scopes[:len(e.scopes)-1]
	}
}

func (e *extractor) getSourceLine(line int) string {
	if line >= 1 && line <= len(e.sourceLines) {
		return e.sourceLines[line-1]
	}
	return ""
}

func (e *extractor) emitTag(
	name string,
	tok antlr.Token,
	kind Kind,
	inherits []string,
	signature string,
) []int {
	if name == "" || name == "_" || tok == nil {
		return nil
	}

	line := tok.GetLine()
	sourceLine := e.getSourceLine(line)

	scopeKind := e.currentScopeKind()
	scopePath := e.currentScopePath()

	allInherits := inherits
	if len(allInherits) == 0 {
		allInherits = e.currentInEachTemplates()
	}

	idx1 := len(e.tags)
	e.tags = append(e.tags, Tag{
		Name:       name,
		File:       e.filePath,
		Line:       line,
		SourceLine: sourceLine,
		Kind:       kind,
		ScopeKind:  scopeKind,
		ScopePath:  scopePath,
		Signature:  signature,
		Inherits:   allInherits,
	})

	indices := []int{idx1}
	if e.opts.EmitQualified && scopePath != "" {
		idx2 := len(e.tags)
		e.tags = append(e.tags, Tag{
			Name:       scopePath + "." + name,
			File:       e.filePath,
			Line:       line,
			SourceLine: sourceLine,
			Kind:       kind,
			ScopeKind:  scopeKind,
			ScopePath:  scopePath,
			Signature:  signature,
			Inherits:   allInherits,
		})
		indices = append(indices, idx2)
	}
	return indices
}

func findIdentToken(t antlr.Tree) antlr.Token {
	if t == nil {
		return nil
	}
	switch node := t.(type) {
	case dmlex.ICdecl_qualifier_ident_array_funcContext:
		return findIdentToken(node.Cdecl_ident_array_func())
	case *dmlex.CdeclIdentContext:
		if node.Ident() != nil {
			return node.Ident().GetStart()
		}
		return nil
	case *dmlex.CdeclArrayContext:
		return findIdentToken(node.Cdecl_ident_array_func())
	case *dmlex.CdeclFuncArgsContext:
		return findIdentToken(node.Cdecl_ident_array_func())
	case *dmlex.CdeclFuncIdentContext:
		return findIdentToken(node.Cdecl_qualifier_ident_array_func())
	case *dmlex.CdeclNoneContext:
		return nil
	case dmlex.IIdentContext:
		return node.GetStart()
	case dmlex.IObjidentContext:
		return node.GetStart()
	}
	childCount := t.GetChildCount()
	for i := range childCount {
		if tok := findIdentToken(t.GetChild(i)); tok != nil {
			return tok
		}
	}
	return nil
}

func findCdeclIdent(cd dmlex.ICdeclContext) (string, antlr.Token) {
	if cd == nil {
		return "", nil
	}
	if tok := findIdentToken(
		cd.Cdecl_qualifier_ident_array_func(),
	); tok != nil {
		return tok.GetText(), tok
	}
	b := cd.Basetype()
	if b == nil {
		return "", nil
	}
	su := b.Signed_or_unsigned_typeident()
	if su == nil || (su.SIGNED() == nil && su.UNSIGNED() == nil) {
		return "", nil
	}
	if su.Typeident() != nil && su.Typeident().Ident() != nil {
		tok := findIdentToken(su.Typeident().Ident())
		if tok != nil {
			return tok.GetText(), tok
		}
	}
	return "", nil
}

func extractIstemplateList(il dmlex.IIstemplateListContext) []string {
	if il == nil {
		return nil
	}
	var res []string
	for _, oid := range il.AllObjident() {
		name := dmlex.IdentName(oid)
		if name != "" && name != "_" {
			res = append(res, name)
		}
	}
	return res
}

func extractTemplateExpr(expr dmlex.IIstemplateExprContext) []string {
	if expr == nil {
		return nil
	}
	return extractIstemplateList(expr.IstemplateList())
}

func formatMethodSignature(m dmlex.IMethodUnifiedContext) string {
	if m == nil {
		return "()"
	}

	var params []string
	var paramCtx dmlex.IMethodParamsMaybeUntypedContext

	switch ctx := m.(type) {
	case *dmlex.MethodUnifiedDeclarationContext:
		paramCtx = ctx.MethodParamsMaybeUntyped()
	case *dmlex.MethodUnifiedDefinitionContext:
		paramCtx = ctx.MethodParamsMaybeUntyped()
	}

	if paramCtx != nil {
		for _, p := range paramCtx.AllCdeclOrInline() {
			if p == nil {
				continue
			}
			if p.Ident() != nil {
				params = append(params, p.Ident().GetText())
				continue
			}
			if p.Cdecl() != nil {
				cd := p.Cdecl()
				name, _ := findCdeclIdent(cd)
				typeStr := dmlex.CdeclTypeName(cd)
				if name != "" && name != "_" {
					params = append(params, typeStr+" "+name)
				} else if typeStr != "" {
					params = append(params, typeStr)
				} else {
					params = append(params, cd.GetText())
				}
			}
		}
	}

	paramStr := "(" + strings.Join(params, ", ") + ")"

	returns := dmlex.MethodReturns(m)
	if len(returns) == 0 {
		return paramStr
	}

	var rets []string
	for _, r := range returns {
		if r == nil {
			continue
		}
		name, _ := findCdeclIdent(r)
		typeStr := dmlex.CdeclTypeName(r)
		if name != "" && name != "_" {
			rets = append(rets, typeStr+" "+name)
		} else if typeStr != "" {
			rets = append(rets, typeStr)
		} else {
			rets = append(rets, r.GetText())
		}
	}

	if len(rets) == 0 {
		return paramStr
	}
	return paramStr + " -> (" + strings.Join(rets, ", ") + ")"
}

func (e *extractor) processDeviceStatements(
	stmts dmlex.IDeviceStatementsContext,
) {
	if stmts == nil {
		return
	}
	for _, stmt := range stmts.AllDeviceStatement() {
		e.processDeviceStatement(stmt)
	}
}

func (e *extractor) processDeviceStatement(stmt dmlex.IDeviceStatementContext) {
	if stmt == nil {
		return
	}
	dsCtx, ok := stmt.(*dmlex.DeviceStatementContext)
	if !ok {
		return
	}

	switch {
	case dsCtx.TemplateStatement() != nil:
		if ts, ok := dsCtx.TemplateStatement().(*dmlex.TemplateStatementContext); ok {
			e.processTemplateStatement(ts)
		}
	case dsCtx.Object() != nil:
		e.processObject(dsCtx.Object())
	case dsCtx.Param() != nil:
		if p, ok := dsCtx.Param().(*dmlex.ParamContext); ok {
			e.processParam(p)
		}
	case dsCtx.MethodUnified() != nil:
		e.processMethodUnified(dsCtx.MethodUnified())
	case dsCtx.ConstantStatement() != nil:
		if cs, ok := dsCtx.ConstantStatement().(*dmlex.ConstantStatementContext); ok {
			e.processConstantStatement(cs)
		}
	case dsCtx.LogGroupStatement() != nil:
		if lgs, ok := dsCtx.LogGroupStatement().(*dmlex.LogGroupStatementContext); ok {
			e.processLogGroupStatement(lgs)
		}
	case dsCtx.TypedefStatement() != nil:
		if tds, ok := dsCtx.TypedefStatement().(*dmlex.TypedefStatementContext); ok {
			e.processTypedefStatement(tds)
		}
	case dsCtx.ExternStatement() != nil:
		if es, ok := dsCtx.ExternStatement().(*dmlex.ExternStatementContext); ok {
			e.processExternStatement(es)
		}
	case dsCtx.IneachStatement() != nil:
		if ie, ok := dsCtx.IneachStatement().(*dmlex.IneachStatementContext); ok {
			e.processIneachStatement(ie)
		}
	case dsCtx.ToplevelIf() != nil:
		if tif, ok := dsCtx.ToplevelIf().(*dmlex.ToplevelIfContext); ok {
			e.processToplevelIf(tif)
		}
	case dsCtx.IstemplateStatement() != nil:
		if is, ok := dsCtx.IstemplateStatement().(*dmlex.IstemplateStatementContext); ok {
			templates := extractTemplateExpr(is.IstemplateExpr())
			e.addCurrentScopeInherits(templates)
		}
	}
}

func (e *extractor) processTemplateStatement(
	ctx *dmlex.TemplateStatementContext,
) {
	if ctx == nil || ctx.Objident() == nil {
		return
	}
	name := dmlex.IdentName(ctx.Objident())
	tok := ctx.Objident().GetStart()
	templates := extractTemplateExpr(ctx.IstemplateExpr())

	indices := e.emitTag(name, tok, KindTemplate, templates, "")

	e.pushScope(KindTemplate, name, nil, indices)
	defer e.popScope()

	for _, body := range ctx.AllTemplateBody() {
		tb, ok := body.(*dmlex.TemplateBodyContext)
		if !ok {
			continue
		}
		if tb.ObjectStatement() != nil {
			e.processObjectStatement(tb.ObjectStatement())
		} else if tb.Hook_decl() != nil {
			if hd, ok := tb.Hook_decl().(*dmlex.Hook_declContext); ok {
				e.processHookDecl(hd)
			}
		}
	}
}

func (e *extractor) processObject(obj dmlex.IObjectContext) {
	if obj == nil {
		return
	}
	objCtx, ok := obj.(*dmlex.ObjectContext)
	if !ok {
		return
	}

	switch {
	case objCtx.Object_decl() != nil:
		if od, ok := objCtx.Object_decl().(*dmlex.Object_declContext); ok {
			e.processObjectDecl(od)
		}
	case objCtx.SessionOrSavedStatement() != nil:
		if sss, ok := objCtx.SessionOrSavedStatement().(*dmlex.SessionOrSavedStatementContext); ok {
			e.processSessionOrSavedStatement(sss)
		}
	case objCtx.Hook_decl() != nil:
		if hd, ok := objCtx.Hook_decl().(*dmlex.Hook_declContext); ok {
			e.processHookDecl(hd)
		}
	}
}

func (e *extractor) processObjectDecl(ctx *dmlex.Object_declContext) {
	if ctx == nil || ctx.Objident() == nil {
		return
	}
	name := dmlex.IdentName(ctx.Objident())
	tok := ctx.Objident().GetStart()

	kind := KindBank
	switch {
	case ctx.REGISTER() != nil:
		kind = KindRegister
	case ctx.FIELD() != nil:
		kind = KindField
	case ctx.Object_type() != nil:
		ot := ctx.Object_type()
		switch {
		case ot.BANK() != nil:
			kind = KindBank
		case ot.GROUP() != nil:
			kind = KindGroup
		case ot.PORT() != nil:
			kind = KindPort
		case ot.ATTRIBUTE() != nil:
			kind = KindAttribute
		case ot.CONNECT() != nil:
			kind = KindConnect
		case ot.EVENT() != nil:
			kind = KindEvent
		case ot.IMPLEMENT() != nil:
			kind = KindImplement
		case ot.INTERFACE() != nil:
			kind = KindInterface
		case ot.SUBDEVICE() != nil:
			kind = KindSubdevice
		}
	}

	var templates []string
	var bodyStmts []dmlex.IObjectStatementContext
	if ctx.ObjectBody() != nil {
		if ob, ok := ctx.ObjectBody().(*dmlex.ObjectBodyContext); ok {
			templates = extractTemplateExpr(ob.IstemplateExpr())
			bodyStmts = ob.AllObjectStatement()
		}
	}

	indices := e.emitTag(name, tok, kind, templates, "")

	e.pushScope(kind, name, nil, indices)
	defer e.popScope()

	for _, stmt := range bodyStmts {
		e.processObjectStatement(stmt)
	}
}

func (e *extractor) processObjectStatement(stmt dmlex.IObjectStatementContext) {
	if stmt == nil {
		return
	}
	osCtx, ok := stmt.(*dmlex.ObjectStatementContext)
	if !ok {
		return
	}

	switch {
	case osCtx.Object() != nil:
		e.processObject(osCtx.Object())
	case osCtx.Param() != nil:
		if p, ok := osCtx.Param().(*dmlex.ParamContext); ok {
			e.processParam(p)
		}
	case osCtx.MethodUnified() != nil:
		e.processMethodUnified(osCtx.MethodUnified())
	case osCtx.IneachStatement() != nil:
		if ie, ok := osCtx.IneachStatement().(*dmlex.IneachStatementContext); ok {
			e.processIneachStatement(ie)
		}
	case osCtx.ObjectIf() != nil:
		if oif, ok := osCtx.ObjectIf().(*dmlex.ObjectIfContext); ok {
			e.processObjectIf(oif)
		}
	case osCtx.IstemplateStatement() != nil:
		if is, ok := osCtx.IstemplateStatement().(*dmlex.IstemplateStatementContext); ok {
			templates := extractTemplateExpr(is.IstemplateExpr())
			e.addCurrentScopeInherits(templates)
		}
	}
}

func (e *extractor) processParam(ctx *dmlex.ParamContext) {
	if ctx == nil || ctx.Objident() == nil {
		return
	}
	name := dmlex.IdentName(ctx.Objident())
	tok := ctx.Objident().GetStart()
	e.emitTag(name, tok, KindParam, nil, "")

	if ctx.Ctypedecl() != nil {
		if ct, ok := ctx.Ctypedecl().(*dmlex.CtypedeclContext); ok &&
			ct.Basetype() != nil {
			e.pushScope(KindParam, name, nil, nil)
			e.processCompositeMembers(ct.Basetype())
			e.popScope()
		}
	}
}

func (e *extractor) processMethodUnified(m dmlex.IMethodUnifiedContext) {
	if m == nil {
		return
	}
	var tok antlr.Token
	name := ""
	switch ctx := m.(type) {
	case *dmlex.MethodUnifiedDeclarationContext:
		if ctx.Objident() != nil {
			name = dmlex.IdentName(ctx.Objident())
			tok = ctx.Objident().GetStart()
		}
	case *dmlex.MethodUnifiedDefinitionContext:
		if ctx.Objident() != nil {
			name = dmlex.IdentName(ctx.Objident())
			tok = ctx.Objident().GetStart()
		}
	}
	if name == "" || tok == nil {
		return
	}

	sig := formatMethodSignature(m)
	e.emitTag(name, tok, KindMethod, nil, sig)
}

func (e *extractor) processConstantStatement(
	ctx *dmlex.ConstantStatementContext,
) {
	if ctx == nil || ctx.Ident() == nil {
		return
	}
	name := ctx.Ident().GetText()
	tok := ctx.Ident().GetStart()
	e.emitTag(name, tok, KindConstant, nil, "")
}

func (e *extractor) processLogGroupStatement(
	ctx *dmlex.LogGroupStatementContext,
) {
	if ctx == nil || ctx.Ident() == nil {
		return
	}
	name := ctx.Ident().GetText()
	tok := ctx.Ident().GetStart()
	e.emitTag(name, tok, KindLogGroup, nil, "")
}

func (e *extractor) processHookDecl(ctx *dmlex.Hook_declContext) {
	if ctx == nil || ctx.Ident() == nil {
		return
	}
	name := ctx.Ident().GetText()
	tok := ctx.Ident().GetStart()
	e.emitTag(name, tok, KindHook, nil, "")
}

func (e *extractor) processSessionOrSavedStatement(
	ctx *dmlex.SessionOrSavedStatementContext,
) {
	if ctx == nil || ctx.SessionOrSavedDecl() == nil {
		return
	}
	decl, ok := ctx.SessionOrSavedDecl().(*dmlex.SessionOrSavedDeclContext)
	if !ok || decl.CvarDecl() == nil {
		return
	}
	e.processCvarDecl(decl.CvarDecl(), KindVariable)
}

func (e *extractor) processCvarDecl(cvar dmlex.ICvarDeclContext, kind Kind) {
	if cvar == nil {
		return
	}
	var cdecls []dmlex.ICdeclContext
	if cvar.Cdecl() != nil {
		cdecls = append(cdecls, cvar.Cdecl())
	}
	if cvar.Cdecl_list() != nil {
		cdecls = append(cdecls, cvar.Cdecl_list().AllCdecl()...)
	}

	for _, cd := range cdecls {
		if cd == nil {
			continue
		}
		name, tok := findCdeclIdent(cd)
		if name == "" || tok == nil {
			continue
		}
		e.emitTag(name, tok, kind, nil, "")
		if isStructLikeBasetype(cd.Basetype()) {
			e.pushScope(kind, name, nil, nil)
			e.processCompositeMembers(cd.Basetype())
			e.popScope()
		}
	}
}

func (e *extractor) processTypedefStatement(
	ctx *dmlex.TypedefStatementContext,
) {
	if ctx == nil || ctx.Cdecl() == nil {
		return
	}
	name, tok := findCdeclIdent(ctx.Cdecl())
	if name == "" || tok == nil {
		return
	}
	kind := KindTypedef
	if isStructLikeBasetype(ctx.Cdecl().Basetype()) {
		kind = KindStruct
	}
	e.emitTag(name, tok, kind, nil, "")

	if isStructLikeBasetype(ctx.Cdecl().Basetype()) {
		e.pushScope(kind, name, nil, nil)
		e.processCompositeMembers(ctx.Cdecl().Basetype())
		e.popScope()
	}
}

func (e *extractor) processExternStatement(ctx *dmlex.ExternStatementContext) {
	if ctx == nil || ctx.Cdecl() == nil {
		return
	}
	name, tok := findCdeclIdent(ctx.Cdecl())
	if name == "" || tok == nil {
		return
	}
	kind := KindExtern
	if ctx.TYPEDEF() != nil {
		kind = KindTypedef
		if isStructLikeBasetype(ctx.Cdecl().Basetype()) {
			kind = KindStruct
		}
	}
	e.emitTag(name, tok, kind, nil, "")

	if isStructLikeBasetype(ctx.Cdecl().Basetype()) {
		e.pushScope(kind, name, nil, nil)
		e.processCompositeMembers(ctx.Cdecl().Basetype())
		e.popScope()
	}
}

func (e *extractor) processIneachStatement(ctx *dmlex.IneachStatementContext) {
	if ctx == nil {
		return
	}
	templates := extractIstemplateList(ctx.IstemplateList())

	// Per design decision: attribute symbols inside 'in each' to their enclosing named scope
	// and record the 'in each' target templates in Inherits.
	e.pushScope(e.currentScopeKind(), "", templates, nil)
	defer e.popScope()

	for _, stmt := range ctx.AllObjectStatement() {
		e.processObjectStatement(stmt)
	}
}

func (e *extractor) processToplevelIf(ctx *dmlex.ToplevelIfContext) {
	if ctx == nil {
		return
	}
	if ctx.ToplevelIfStatement() != nil {
		if tis, ok := ctx.ToplevelIfStatement().(*dmlex.ToplevelIfStatementContext); ok {
			e.processDeviceStatements(tis.DeviceStatements())
		}
	}
	if ctx.ToplevelElseStatement() != nil {
		if tes, ok := ctx.ToplevelElseStatement().(*dmlex.ToplevelElseStatementContext); ok {
			if tes.DeviceStatements() != nil {
				e.processDeviceStatements(tes.DeviceStatements())
			} else if tes.ToplevelIf() != nil {
				if tif, ok := tes.ToplevelIf().(*dmlex.ToplevelIfContext); ok {
					e.processToplevelIf(tif)
				}
			}
		}
	}
}

func (e *extractor) processObjectIf(ctx *dmlex.ObjectIfContext) {
	if ctx == nil {
		return
	}
	if ctx.ObjectIfStatement() != nil {
		if ois, ok := ctx.ObjectIfStatement().(*dmlex.ObjectIfStatementContext); ok {
			for _, stmt := range ois.AllObjectStatement() {
				e.processObjectStatement(stmt)
			}
		}
	}
	if ctx.ObjectElseStatement() != nil {
		if oes, ok := ctx.ObjectElseStatement().(*dmlex.ObjectElseStatementContext); ok {
			for _, stmt := range oes.AllObjectStatement() {
				e.processObjectStatement(stmt)
			}
		}
	}
}

func isStructLikeBasetype(b dmlex.IBasetypeContext) bool {
	if b == nil {
		return false
	}
	return b.Struct_() != nil || b.Union() != nil ||
		b.Layout() != nil || b.Bitfields() != nil
}

func (e *extractor) processCompositeMembers(base dmlex.IBasetypeContext) {
	if base == nil {
		return
	}
	switch {
	case base.Struct_() != nil:
		if st, ok := base.Struct_().(*dmlex.StructContext); ok {
			for _, d := range st.AllStructDecls() {
				if sd, ok := d.(*dmlex.StructDeclsContext); ok &&
					sd.Cdecl() != nil {
					e.emitCompositeMember(sd.Cdecl())
				}
			}
		}
	case base.Union() != nil:
		if un, ok := base.Union().(*dmlex.UnionContext); ok {
			for _, d := range un.AllStructDecls() {
				if sd, ok := d.(*dmlex.StructDeclsContext); ok &&
					sd.Cdecl() != nil {
					e.emitCompositeMember(sd.Cdecl())
				}
			}
		}
	case base.Layout() != nil:
		if lay, ok := base.Layout().(*dmlex.LayoutContext); ok {
			for _, d := range lay.AllLayoutDecls() {
				if ld, ok := d.(*dmlex.LayoutDeclsContext); ok &&
					ld.Cdecl() != nil {
					e.emitCompositeMember(ld.Cdecl())
				}
			}
		}
	case base.Bitfields() != nil:
		if bf, ok := base.Bitfields().(*dmlex.BitfieldsContext); ok {
			for _, d := range bf.AllBitfields_decls() {
				if bd, ok := d.(*dmlex.Bitfields_declsContext); ok &&
					bd.Cdecl() != nil {
					e.emitCompositeMember(bd.Cdecl())
				}
			}
		}
	}
}

func (e *extractor) emitCompositeMember(cd dmlex.ICdeclContext) {
	if cd == nil {
		return
	}
	name, tok := findCdeclIdent(cd)
	if name == "" || tok == nil {
		return
	}
	e.emitTag(name, tok, KindMember, nil, "")
}

// ExtractFileTags parses an ANTLR DML parse tree and extracts all tags.
func ExtractFileTags(
	tree antlr.ParseTree,
	filePath string,
	sourceLines []string,
	opts ExtractOptions,
) []Tag {
	if tree == nil {
		return nil
	}

	dmlCtx, ok := tree.(*dmlex.DmlContext)
	if !ok || dmlCtx.DmlCode() == nil {
		return nil
	}

	codeCtx, ok := dmlCtx.DmlCode().(*dmlex.DmlCodeContext)
	if !ok || codeCtx == nil {
		return nil
	}

	e := &extractor{
		filePath:    filePath,
		sourceLines: sourceLines,
		opts:        opts,
	}

	if codeCtx.DeviceTag() != nil {
		if dt, ok := codeCtx.DeviceTag().(*dmlex.DeviceTagContext); ok &&
			dt.Objident() != nil {
			name := dmlex.IdentName(dt.Objident())
			tok := dt.Objident().GetStart()
			indices := e.emitTag(name, tok, KindDevice, nil, "")
			e.deviceTagSeen = true
			e.pushScope(KindDevice, name, nil, indices)
		}
	}

	if codeCtx.DeviceStatements() != nil {
		if ds, ok := codeCtx.DeviceStatements().(*dmlex.DeviceStatementsContext); ok {
			e.processDeviceStatements(ds)
		}
	}

	if e.deviceTagSeen {
		e.popScope()
	}

	return e.tags
}
