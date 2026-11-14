// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lsp

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/fs/stream"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

// SymbolKind represents the classification of a DML symbol.
type SymbolKind int

const (
	SymbolKindUnknown SymbolKind = iota
	SymbolKindTemplate
	SymbolKindMethod
	SymbolKindParam
	SymbolKindConstant
	SymbolKindTypedef
	SymbolKindObject
	SymbolKindVariable
	SymbolKindSession
	SymbolKindSaved
	SymbolKindHook
	SymbolKindExtern
	SymbolKindLogGroup
	SymbolKindParameter
	SymbolKindStruct
)

// Symbol represents a declared DML symbol with location and hierarchy.
type Symbol struct {
	Name         string
	Kind         SymbolKind
	URI          uri.URI
	NameRange    protocol.Range
	FullRange    protocol.Range
	Parent       *Symbol
	Children     map[string]*Symbol
	Scope        *Scope
	Templates    []string
	IsConst      bool
	IsStatic     bool
	IsAbstract   bool
	IsFunc       bool
	IsDefinition bool
}

// ScopeKind defines whether a scope is global, object, method, or block.
type ScopeKind int

const (
	ScopeGlobal ScopeKind = iota
	ScopeObject
	ScopeMethod
	ScopeBlock
)

// Scope represents a lexical or hierarchical scope in DML source.
type Scope struct {
	Kind     ScopeKind
	Owner    *Symbol
	Parent   *Scope
	Symbols  map[string]*Symbol
	Children []*Scope
	Range    protocol.Range
}

// ImportInfo records an import statement in a DML document.
type ImportInfo struct {
	Path  string
	Range protocol.Range
}

// FileSymbols aggregates symbols, scopes, and imports for a single file.
type FileSymbols struct {
	URI             uri.URI
	GlobalScope     *Scope
	Imports         []ImportInfo
	DeviceSymbol    *Symbol
	DeviceTemplates []string
}

func positionLessOrEqual(a, b protocol.Position) bool {
	if a.Line < b.Line {
		return true
	}
	if a.Line > b.Line {
		return false
	}
	return a.Character <= b.Character
}

func rangeContains(r protocol.Range, pos protocol.Position) bool {
	return positionLessOrEqual(r.Start, pos) &&
		positionLessOrEqual(pos, r.End)
}

// EnclosingScope locates the innermost scope containing pos.
func (s *Scope) EnclosingScope(pos protocol.Position) *Scope {
	if !rangeContains(s.Range, pos) {
		return nil
	}
	for _, child := range s.Children {
		if enc := child.EnclosingScope(pos); enc != nil {
			return enc
		}
	}
	return s
}

func tokenUTF16Col(tok antlr.Token) uint32 {
	if tok == nil {
		return 0
	}
	if hb, ok := tok.GetInputStream().(*stream.HashedBuffer); ok && hb != nil {
		data := hb.GetRawString()
		start := tok.GetStart()
		if start >= 0 && start <= len(data) {
			lineStart := start
			for lineStart > 0 && data[lineStart-1] != '\n' {
				lineStart--
			}
			return utf16Len(data[lineStart:start])
		}
	}
	if tok.GetColumn() >= 0 {
		return uint32(tok.GetColumn()) //nolint:gosec
	}
	return 0
}

func tokenRange(tok antlr.Token) protocol.Range {
	if tok == nil {
		return protocol.Range{}
	}
	line := uint32(0)
	if tok.GetLine() > 0 {
		line = uint32(tok.GetLine() - 1) //nolint:gosec
	}
	col := tokenUTF16Col(tok)
	length := utf16Len(tok.GetText())
	return protocol.Range{
		Start: protocol.Position{Line: line, Character: col},
		End:   protocol.Position{Line: line, Character: col + length},
	}
}

func ruleRange(ctx antlr.ParserRuleContext) protocol.Range {
	if ctx == nil || ctx.GetStart() == nil {
		return protocol.Range{}
	}
	start := tokenRange(ctx.GetStart()).Start
	end := start
	if stop := ctx.GetStop(); stop != nil {
		end = tokenRange(stop).End
	}
	return protocol.Range{Start: start, End: end}
}

func extractObjident(ctx dmlex.IObjidentContext) (string, protocol.Range) {
	if ctx == nil {
		return "", protocol.Range{}
	}
	return dmlex.IdentName(ctx), ruleRange(ctx)
}

func extractIdent(ctx dmlex.IIdentContext) (string, protocol.Range) {
	if ctx == nil {
		return "", protocol.Range{}
	}
	return ctx.GetText(), ruleRange(ctx)
}

func findIdentTokenInTree(t antlr.Tree) antlr.Token {
	if t == nil {
		return nil
	}
	switch node := t.(type) {
	case dmlex.ICdecl_qualifier_ident_array_funcContext:
		return findIdentTokenInTree(node.Cdecl_ident_array_func())
	case *dmlex.CdeclIdentContext:
		if node.Ident() != nil {
			return node.Ident().GetStart()
		}
		return nil
	case *dmlex.CdeclArrayContext:
		return findIdentTokenInTree(node.Cdecl_ident_array_func())
	case *dmlex.CdeclFuncArgsContext:
		return findIdentTokenInTree(node.Cdecl_ident_array_func())
	case *dmlex.CdeclFuncIdentContext:
		return findIdentTokenInTree(node.Cdecl_qualifier_ident_array_func())
	case *dmlex.CdeclNoneContext:
		return nil
	case dmlex.IIdentContext:
		return node.GetStart()
	}
	childCount := t.GetChildCount()
	for i := range childCount {
		if tok := findIdentTokenInTree(t.GetChild(i)); tok != nil {
			return tok
		}
	}
	return nil
}

func isConstCdecl(c dmlex.ICdeclContext) bool {
	if c == nil {
		return false
	}
	seenMul := false
	constAfterLastMul := false
	if q := c.Cdecl_qualifier_ident_array_func(); q != nil {
		for _, tq := range q.AllTypeQualifier() {
			if tq.MUL() != nil {
				seenMul = true
				constAfterLastMul = false
			} else if tq.CONST() != nil {
				constAfterLastMul = true
			}
		}
	}
	if !seenMul {
		return c.CONST() != nil || constAfterLastMul
	}
	return constAfterLastMul
}

func isDirectFuncCdecl(
	q dmlex.ICdecl_qualifier_ident_array_funcContext,
) bool {
	if q == nil {
		return false
	}
	fnCtx, ok := q.Cdecl_ident_array_func().(*dmlex.CdeclFuncArgsContext)
	if !ok {
		return false
	}
	inner := fnCtx.Cdecl_ident_array_func()
	if _, isDirect := inner.(*dmlex.CdeclIdentContext); isDirect {
		return true
	}
	if _, isNone := inner.(*dmlex.CdeclNoneContext); isNone {
		return true
	}
	return false
}

func isStructLikeBasetype(b dmlex.IBasetypeContext) bool {
	if b == nil {
		return false
	}
	return b.Struct_() != nil || b.Union() != nil ||
		b.Layout() != nil || b.Bitfields() != nil
}

func findCdeclIdentToken(cd dmlex.ICdeclContext) antlr.Token {
	if cd == nil {
		return nil
	}
	if tok := findIdentTokenInTree(
		cd.Cdecl_qualifier_ident_array_func(),
	); tok != nil {
		return tok
	}
	b := cd.Basetype()
	if b == nil {
		return nil
	}
	su := b.Signed_or_unsigned_typeident()
	if su == nil || (su.SIGNED() == nil && su.UNSIGNED() == nil) {
		return nil
	}
	if su.Typeident() != nil && su.Typeident().Ident() != nil {
		return findIdentTokenInTree(su.Typeident().Ident())
	}
	return nil
}

func findCdeclIdent(cd dmlex.ICdeclContext) (string, protocol.Range) {
	tok := findCdeclIdentToken(cd)
	if tok == nil {
		return "", protocol.Range{}
	}
	return tok.GetText(), tokenRange(tok)
}

type symbolCollector struct {
	dmlex.BasedmlListener

	fileSymbols  *FileSymbols
	currentScope *Scope
	currentObj   *Symbol
}

func (c *symbolCollector) parentObj() *Symbol {
	if c.currentObj != nil {
		return c.currentObj
	}
	return c.fileSymbols.DeviceSymbol
}

func (c *symbolCollector) attachToObjOrDevice(name string, sym *Symbol) {
	if c.currentObj != nil {
		if c.currentObj.Children != nil {
			c.currentObj.Children[name] = sym
		}
		return
	}
	if c.fileSymbols.DeviceSymbol != nil {
		c.fileSymbols.DeviceSymbol.Children[name] = sym
	}
}

func (c *symbolCollector) popObj() {
	if c.currentObj == nil {
		return
	}
	if c.currentObj.Parent == c.fileSymbols.DeviceSymbol {
		c.currentObj = nil
		return
	}
	c.currentObj = c.currentObj.Parent
}

func (c *symbolCollector) EnterDml(ctx *dmlex.DmlContext) {
	c.fileSymbols.GlobalScope.Range = ruleRange(ctx)
	c.currentScope = c.fileSymbols.GlobalScope
}

func (c *symbolCollector) EnterDeviceTag(ctx *dmlex.DeviceTagContext) {
	name, r := extractObjident(ctx.Objident())
	if name == "" {
		return
	}
	sym := &Symbol{
		Name:         name,
		Kind:         SymbolKindObject,
		URI:          c.fileSymbols.URI,
		NameRange:    r,
		FullRange:    ruleRange(ctx),
		Children:     make(map[string]*Symbol),
		Templates:    c.fileSymbols.DeviceTemplates,
		IsDefinition: true,
	}
	c.fileSymbols.DeviceSymbol = sym
	c.currentScope.Symbols[name] = sym
}

func (c *symbolCollector) EnterImportStatement(
	ctx *dmlex.ImportStatementContext,
) {
	lit := ctx.STRING_LITERAL()
	if lit == nil {
		return
	}
	raw := lit.GetText()
	clean := dmlex.UnquoteDmlString(raw)
	c.fileSymbols.Imports = append(c.fileSymbols.Imports, ImportInfo{
		Path:  clean,
		Range: ruleRange(ctx),
	})
}

func (c *symbolCollector) EnterTemplateStatement(
	ctx *dmlex.TemplateStatementContext,
) {
	name, r := extractObjident(ctx.Objident())
	sym := &Symbol{
		Name:         name,
		Kind:         SymbolKindTemplate,
		URI:          c.fileSymbols.URI,
		NameRange:    r,
		FullRange:    ruleRange(ctx),
		Parent:       c.currentObj,
		Children:     make(map[string]*Symbol),
		IsDefinition: true,
	}
	if name != "" {
		c.currentScope.Symbols[name] = sym
		if c.currentObj != nil {
			c.currentObj.Children[name] = sym
		}
	}

	scope := &Scope{
		Kind:    ScopeObject,
		Owner:   sym,
		Parent:  c.currentScope,
		Symbols: make(map[string]*Symbol),
		Range:   ruleRange(ctx),
	}
	c.currentScope.Children = append(c.currentScope.Children, scope)
	sym.Scope = scope

	c.currentScope = scope
	c.currentObj = sym
}

func (c *symbolCollector) ExitTemplateStatement(
	_ *dmlex.TemplateStatementContext,
) {
	if c.currentScope != nil && c.currentScope.Parent != nil {
		c.currentScope = c.currentScope.Parent
	}
	c.popObj()
}

func (c *symbolCollector) EnterIneachStatement(
	ctx *dmlex.IneachStatementContext,
) {
	sym := &Symbol{
		Kind:      SymbolKindObject,
		URI:       c.fileSymbols.URI,
		FullRange: ruleRange(ctx),
		Parent:    c.parentObj(),
		Children:  make(map[string]*Symbol),
	}
	scope := &Scope{
		Kind:    ScopeObject,
		Owner:   sym,
		Parent:  c.currentScope,
		Symbols: make(map[string]*Symbol),
		Range:   ruleRange(ctx),
	}
	c.currentScope.Children = append(c.currentScope.Children, scope)
	sym.Scope = scope

	c.currentScope = scope
	c.currentObj = sym
}

func (c *symbolCollector) ExitIneachStatement(
	_ *dmlex.IneachStatementContext,
) {
	if c.currentScope != nil && c.currentScope.Parent != nil {
		c.currentScope = c.currentScope.Parent
	}
	c.popObj()
}

func (c *symbolCollector) EnterObject_decl(
	ctx *dmlex.Object_declContext,
) {
	name, r := extractObjident(ctx.Objident())
	sym := &Symbol{
		Name:         name,
		Kind:         SymbolKindObject,
		URI:          c.fileSymbols.URI,
		NameRange:    r,
		FullRange:    ruleRange(ctx),
		Parent:       c.parentObj(),
		Children:     make(map[string]*Symbol),
		IsDefinition: true,
	}
	if name != "" {
		c.currentScope.Symbols[name] = sym
		c.attachToObjOrDevice(name, sym)
	}

	scope := &Scope{
		Kind:    ScopeObject,
		Owner:   sym,
		Parent:  c.currentScope,
		Symbols: make(map[string]*Symbol),
		Range:   ruleRange(ctx),
	}
	c.currentScope.Children = append(c.currentScope.Children, scope)
	sym.Scope = scope

	c.currentScope = scope
	c.currentObj = sym
}

func (c *symbolCollector) ExitObject_decl(
	_ *dmlex.Object_declContext,
) {
	if c.currentScope != nil && c.currentScope.Parent != nil {
		c.currentScope = c.currentScope.Parent
	}
	c.popObj()
}

func (c *symbolCollector) EnterIstemplateList(
	ctx *dmlex.IstemplateListContext,
) {
	for _, oid := range ctx.AllObjident() {
		name := dmlex.IdentName(oid)
		if name == "" {
			continue
		}
		if c.currentObj != nil {
			c.currentObj.Templates = append(c.currentObj.Templates, name)
		} else {
			c.fileSymbols.DeviceTemplates = append(
				c.fileSymbols.DeviceTemplates,
				name,
			)
			if c.fileSymbols.DeviceSymbol != nil {
				c.fileSymbols.DeviceSymbol.Templates = c.fileSymbols.DeviceTemplates
			}
		}
	}
}

func (c *symbolCollector) EnterMethodUnifiedDeclaration(
	ctx *dmlex.MethodUnifiedDeclarationContext,
) {
	if ctx.Objident() == nil {
		return
	}
	name := dmlex.MethodName(ctx)
	if name == "" {
		return
	}
	isAbstract := ctx.MethodQualifiers() != nil &&
		ctx.MethodQualifiers().SHARED() != nil
	r := ruleRange(ctx.Objident())
	sym := &Symbol{
		Name:       name,
		Kind:       SymbolKindMethod,
		URI:        c.fileSymbols.URI,
		NameRange:  r,
		FullRange:  ruleRange(ctx),
		Parent:     c.parentObj(),
		Children:   make(map[string]*Symbol),
		IsAbstract: isAbstract,
	}
	c.currentScope.Symbols[name] = sym
	c.attachToObjOrDevice(name, sym)
}

func (c *symbolCollector) EnterMethodUnifiedDefinition(
	ctx *dmlex.MethodUnifiedDefinitionContext,
) {
	name := ""
	var r protocol.Range
	if ctx.Objident() != nil {
		name = dmlex.MethodName(ctx)
		r = ruleRange(ctx.Objident())
	}
	sym := &Symbol{
		Name:         name,
		Kind:         SymbolKindMethod,
		URI:          c.fileSymbols.URI,
		NameRange:    r,
		FullRange:    ruleRange(ctx),
		Parent:       c.parentObj(),
		Children:     make(map[string]*Symbol),
		IsDefinition: true,
	}
	if name != "" {
		c.currentScope.Symbols[name] = sym
		c.attachToObjOrDevice(name, sym)
	}

	scope := &Scope{
		Kind:    ScopeMethod,
		Owner:   sym,
		Parent:  c.currentScope,
		Symbols: make(map[string]*Symbol),
		Range:   ruleRange(ctx),
	}
	for _, outParam := range dmlex.MethodReturns(ctx) {
		if outParam == nil {
			continue
		}
		outName, outRange := findCdeclIdent(outParam)
		if outName != "" && outName != "_" {
			scope.Symbols[outName] = &Symbol{
				Name:      outName,
				Kind:      SymbolKindParameter,
				URI:       c.fileSymbols.URI,
				NameRange: outRange,
				FullRange: ruleRange(outParam),
				Parent:    sym,
				IsConst:   isConstCdecl(outParam),
			}
		}
	}

	c.currentScope.Children = append(c.currentScope.Children, scope)
	sym.Scope = scope

	c.currentScope = scope
}

func (c *symbolCollector) ExitMethodUnifiedDefinition(
	_ *dmlex.MethodUnifiedDefinitionContext,
) {
	if c.currentScope != nil && c.currentScope.Parent != nil {
		c.currentScope = c.currentScope.Parent
	}
}

func (c *symbolCollector) EnterCdeclOrInline(
	ctx *dmlex.CdeclOrInlineContext,
) {
	if c.currentScope == nil || c.currentScope.Kind != ScopeMethod {
		return
	}
	name := ""
	var r protocol.Range
	var isConst bool
	if ctx.Ident() != nil {
		name = ctx.Ident().GetText()
		r = ruleRange(ctx.Ident())
	} else if ctx.Cdecl() != nil {
		cdecl := ctx.Cdecl()
		isConst = isConstCdecl(cdecl)
		name, r = findCdeclIdent(cdecl)
	}
	if name != "" && name != "_" {
		paramSym := &Symbol{
			Name:      name,
			Kind:      SymbolKindParameter,
			URI:       c.fileSymbols.URI,
			NameRange: r,
			FullRange: ruleRange(ctx),
			Parent:    c.currentScope.Owner,
			IsConst:   isConst,
		}
		c.currentScope.Symbols[name] = paramSym
	}
}

func (c *symbolCollector) extractCvarSymbols(
	cvar dmlex.ICvarDeclContext,
	kind SymbolKind,
	fullRange protocol.Range,
	parent *Symbol,
	attachToObj bool,
	isStatic bool,
	isDefinition bool,
) {
	if c.currentScope == nil || cvar == nil {
		return
	}
	registerCdecl := func(cd dmlex.ICdeclContext) {
		if cd == nil {
			return
		}
		name, r := findCdeclIdent(cd)
		if name == "" || name == "_" {
			return
		}
		q := cd.Cdecl_qualifier_ident_array_func()
		sym := &Symbol{
			Name:         name,
			Kind:         kind,
			URI:          c.fileSymbols.URI,
			NameRange:    r,
			FullRange:    fullRange,
			Parent:       parent,
			IsConst:      isConstCdecl(cd),
			IsStatic:     isStatic,
			IsFunc:       isDirectFuncCdecl(q),
			IsDefinition: isDefinition,
		}
		c.currentScope.Symbols[name] = sym
		if attachToObj {
			c.attachToObjOrDevice(name, sym)
		}
	}

	registerCdecl(cvar.Cdecl())
	if cvar.Cdecl_list() != nil {
		for _, cd := range cvar.Cdecl_list().AllCdecl() {
			registerCdecl(cd)
		}
	}
}

func (c *symbolCollector) EnterVariableDecl(
	ctx *dmlex.VariableDeclContext,
) {
	if c.currentScope == nil || ctx.CvarDecl() == nil {
		return
	}
	kind := SymbolKindVariable
	isStatic := false
	if lk := ctx.Local_decl_kind(); lk != nil {
		switch {
		case lk.SESSION() != nil:
			kind = SymbolKindSession
			isStatic = true
		case lk.SAVED() != nil:
			kind = SymbolKindSaved
			isStatic = true
		case lk.STATIC() != nil:
			isStatic = true
		}
	}
	isDef := ctx.CvarDecl().ASSIGN() != nil
	c.extractCvarSymbols(
		ctx.CvarDecl(),
		kind,
		ruleRange(ctx),
		c.currentScope.Owner,
		false,
		isStatic,
		isDef,
	)
}

func (c *symbolCollector) EnterSessionOrSavedDecl(
	ctx *dmlex.SessionOrSavedDeclContext,
) {
	kind := SymbolKindSession
	if gk := ctx.Global_decl_kind(); gk != nil && gk.SAVED() != nil {
		kind = SymbolKindSaved
	}
	isDef := ctx.CvarDecl() != nil && ctx.CvarDecl().ASSIGN() != nil
	c.extractCvarSymbols(
		ctx.CvarDecl(),
		kind,
		ruleRange(ctx),
		c.parentObj(),
		true,
		true,
		isDef,
	)
}

func (c *symbolCollector) EnterHook_decl(
	ctx *dmlex.Hook_declContext,
) {
	name, r := extractIdent(ctx.Ident())
	if name == "" {
		return
	}
	hookSym := &Symbol{
		Name:      name,
		Kind:      SymbolKindHook,
		URI:       c.fileSymbols.URI,
		NameRange: r,
		FullRange: ruleRange(ctx),
		Parent:    c.parentObj(),
	}
	c.currentScope.Symbols[name] = hookSym
	c.attachToObjOrDevice(name, hookSym)
}

func (c *symbolCollector) EnterParam(ctx *dmlex.ParamContext) {
	name, r := extractObjident(ctx.Objident())
	if name == "" {
		return
	}
	isDef := ctx.ASSIGN() != nil || ctx.NEW_ASSIGN() != nil ||
		ctx.DEFAULT() != nil || ctx.DEFAULT_NEW() != nil
	paramSym := &Symbol{
		Name:         name,
		Kind:         SymbolKindParam,
		URI:          c.fileSymbols.URI,
		NameRange:    r,
		FullRange:    ruleRange(ctx),
		Parent:       c.parentObj(),
		IsDefinition: isDef,
	}
	c.currentScope.Symbols[name] = paramSym
	c.attachToObjOrDevice(name, paramSym)
}

func (c *symbolCollector) EnterConstantStatement(
	ctx *dmlex.ConstantStatementContext,
) {
	name, r := extractIdent(ctx.Ident())
	if name == "" {
		return
	}
	constSym := &Symbol{
		Name:         name,
		Kind:         SymbolKindConstant,
		URI:          c.fileSymbols.URI,
		NameRange:    r,
		FullRange:    ruleRange(ctx),
		Parent:       c.currentObj,
		IsConst:      true,
		IsDefinition: true,
	}
	c.currentScope.Symbols[name] = constSym
	if c.currentObj != nil {
		c.currentObj.Children[name] = constSym
	}
}

func (c *symbolCollector) EnterTypedefStatement(
	ctx *dmlex.TypedefStatementContext,
) {
	if ctx.Cdecl() == nil {
		return
	}
	name, r := findCdeclIdent(ctx.Cdecl())
	if name == "" {
		return
	}
	kind := SymbolKindTypedef
	if isStructLikeBasetype(ctx.Cdecl().Basetype()) {
		kind = SymbolKindStruct
	}
	q := ctx.Cdecl().Cdecl_qualifier_ident_array_func()
	typeSym := &Symbol{
		Name:         name,
		Kind:         kind,
		URI:          c.fileSymbols.URI,
		NameRange:    r,
		FullRange:    ruleRange(ctx),
		IsConst:      isConstCdecl(ctx.Cdecl()),
		IsFunc:       isDirectFuncCdecl(q),
		IsDefinition: true,
	}
	c.currentScope.Symbols[name] = typeSym
}

func (c *symbolCollector) EnterExternStatement(
	ctx *dmlex.ExternStatementContext,
) {
	if ctx.Cdecl() == nil {
		return
	}
	name, r := findCdeclIdent(ctx.Cdecl())
	if name == "" {
		return
	}
	kind := SymbolKindExtern
	if ctx.TYPEDEF() != nil {
		kind = SymbolKindTypedef
		if isStructLikeBasetype(ctx.Cdecl().Basetype()) {
			kind = SymbolKindStruct
		}
	}
	q := ctx.Cdecl().Cdecl_qualifier_ident_array_func()
	extSym := &Symbol{
		Name:      name,
		Kind:      kind,
		URI:       c.fileSymbols.URI,
		NameRange: r,
		FullRange: ruleRange(ctx),
		IsConst:   isConstCdecl(ctx.Cdecl()),
		IsFunc:    isDirectFuncCdecl(q),
	}
	c.currentScope.Symbols[name] = extSym
}

func (c *symbolCollector) EnterLogGroupStatement(
	ctx *dmlex.LogGroupStatementContext,
) {
	name, r := extractIdent(ctx.Ident())
	if name == "" {
		return
	}
	logSym := &Symbol{
		Name:         name,
		Kind:         SymbolKindLogGroup,
		URI:          c.fileSymbols.URI,
		NameRange:    r,
		FullRange:    ruleRange(ctx),
		IsDefinition: true,
	}
	c.currentScope.Symbols[name] = logSym
}

func (c *symbolCollector) EnterStruct(ctx *dmlex.StructContext) {
	if c.currentScope == nil || ctx.Ident() == nil || ctx.LBRACE() == nil {
		return
	}
	name, r := extractIdent(ctx.Ident())
	if name == "" {
		return
	}
	c.currentScope.Symbols[name] = &Symbol{
		Name:         name,
		Kind:         SymbolKindStruct,
		URI:          c.fileSymbols.URI,
		NameRange:    r,
		FullRange:    ruleRange(ctx),
		IsDefinition: true,
	}
}

func (c *symbolCollector) EnterUnion(ctx *dmlex.UnionContext) {
	if c.currentScope == nil || ctx.Ident() == nil || ctx.LBRACE() == nil {
		return
	}
	name, r := extractIdent(ctx.Ident())
	if name == "" {
		return
	}
	c.currentScope.Symbols[name] = &Symbol{
		Name:         name,
		Kind:         SymbolKindStruct,
		URI:          c.fileSymbols.URI,
		NameRange:    r,
		FullRange:    ruleRange(ctx),
		IsDefinition: true,
	}
}

func (c *symbolCollector) EnterArraydef(ctx *dmlex.ArraydefContext) {
	if c.currentScope == nil || ctx.Ident() == nil {
		return
	}
	name, r := extractIdent(ctx.Ident())
	if name == "" || name == "_" {
		return
	}
	c.currentScope.Symbols[name] = &Symbol{
		Name:         name,
		Kind:         SymbolKindVariable,
		URI:          c.fileSymbols.URI,
		NameRange:    r,
		FullRange:    ruleRange(ctx.Ident()),
		Parent:       c.parentObj(),
		IsDefinition: true,
	}
}

func (c *symbolCollector) EnterForLoopStatement(
	ctx *dmlex.ForLoopStatementContext,
) {
	if c.currentScope == nil {
		return
	}
	scope := &Scope{
		Kind:    ScopeBlock,
		Range:   ruleRange(ctx),
		Parent:  c.currentScope,
		Owner:   c.currentScope.Owner,
		Symbols: make(map[string]*Symbol),
	}
	c.currentScope.Children = append(c.currentScope.Children, scope)
	c.currentScope = scope
}

func (c *symbolCollector) ExitForLoopStatement(
	_ *dmlex.ForLoopStatementContext,
) {
	if c.currentScope != nil && c.currentScope.Parent != nil {
		c.currentScope = c.currentScope.Parent
	}
}

func (c *symbolCollector) EnterForeachStatement(
	ctx *dmlex.ForeachStatementContext,
) {
	if c.currentScope == nil {
		return
	}
	scopeRange := ruleRange(ctx)
	if ctx.LoopStatement() != nil {
		scopeRange = ruleRange(ctx.LoopStatement())
	}
	scope := &Scope{
		Kind:    ScopeBlock,
		Range:   scopeRange,
		Parent:  c.currentScope,
		Owner:   c.currentScope.Owner,
		Symbols: make(map[string]*Symbol),
	}
	name, r := extractIdent(ctx.Ident())
	if name != "" && name != "_" {
		scope.Symbols[name] = &Symbol{
			Name:         name,
			Kind:         SymbolKindVariable,
			URI:          c.fileSymbols.URI,
			NameRange:    r,
			FullRange:    ruleRange(ctx),
			Parent:       c.currentScope.Owner,
			IsDefinition: true,
		}
	}
	c.currentScope.Children = append(c.currentScope.Children, scope)
	c.currentScope = scope
}

func (c *symbolCollector) ExitForeachStatement(
	_ *dmlex.ForeachStatementContext,
) {
	if c.currentScope != nil && c.currentScope.Parent != nil {
		c.currentScope = c.currentScope.Parent
	}
}

func (c *symbolCollector) EnterSelectStatement(
	ctx *dmlex.SelectStatementContext,
) {
	if c.currentScope == nil {
		return
	}
	scopeRange := ruleRange(ctx)
	switch {
	case ctx.WHERE() != nil:
		scopeRange.Start = tokenRange(ctx.WHERE().GetSymbol()).Start
	case ctx.ConditionExpression() != nil:
		scopeRange.Start = ruleRange(ctx.ConditionExpression()).Start
	}
	if ctx.LoopStatement() != nil {
		scopeRange.End = ruleRange(ctx.LoopStatement()).End
	}
	scope := &Scope{
		Kind:    ScopeBlock,
		Range:   scopeRange,
		Parent:  c.currentScope,
		Owner:   c.currentScope.Owner,
		Symbols: make(map[string]*Symbol),
	}
	name, r := extractIdent(ctx.Ident())
	if name != "" && name != "_" {
		scope.Symbols[name] = &Symbol{
			Name:         name,
			Kind:         SymbolKindVariable,
			URI:          c.fileSymbols.URI,
			NameRange:    r,
			FullRange:    ruleRange(ctx),
			Parent:       c.currentScope.Owner,
			IsDefinition: true,
		}
	}
	c.currentScope.Children = append(c.currentScope.Children, scope)
	c.currentScope = scope
}

func (c *symbolCollector) ExitSelectStatement(
	_ *dmlex.SelectStatementContext,
) {
	if c.currentScope != nil && c.currentScope.Parent != nil {
		c.currentScope = c.currentScope.Parent
	}
}

func (c *symbolCollector) EnterAfterStatement(
	ctx *dmlex.AfterStatementContext,
) {
	if c.currentScope == nil {
		return
	}
	cond := ctx.AfterDelayCondition()
	if cond == nil || cond.ARROW() == nil {
		return
	}
	scope := &Scope{
		Kind:    ScopeBlock,
		Parent:  c.currentScope,
		Owner:   c.currentScope.Owner,
		Symbols: make(map[string]*Symbol),
		Range: protocol.Range{
			Start: tokenRange(cond.ARROW().GetSymbol()).Start,
			End:   ruleRange(ctx).End,
		},
	}
	registerParam := func(id dmlex.IIdentContext) {
		name, r := extractIdent(id)
		if name == "" || name == "_" {
			return
		}
		scope.Symbols[name] = &Symbol{
			Name:      name,
			Kind:      SymbolKindParameter,
			URI:       c.fileSymbols.URI,
			NameRange: r,
			FullRange: r,
			Parent:    c.currentScope.Owner,
		}
	}
	registerParam(cond.Ident())
	if cond.IdentList() != nil {
		for _, id := range cond.IdentList().AllIdent() {
			registerParam(id)
		}
	}
	c.currentScope.Children = append(c.currentScope.Children, scope)
	c.currentScope = scope
}

func (c *symbolCollector) ExitAfterStatement(
	ctx *dmlex.AfterStatementContext,
) {
	cond := ctx.AfterDelayCondition()
	if cond == nil || cond.ARROW() == nil {
		return
	}
	if c.currentScope != nil && c.currentScope.Parent != nil {
		c.currentScope = c.currentScope.Parent
	}
}

// ExtractFileSymbols builds the symbol hierarchy and scopes from an AST.
func ExtractFileSymbols(docURI uri.URI, tree antlr.ParseTree) *FileSymbols {
	fs := &FileSymbols{
		URI: docURI,
		GlobalScope: &Scope{
			Kind:    ScopeGlobal,
			Symbols: make(map[string]*Symbol),
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   protocol.Position{Line: 1 << 30, Character: 1 << 30},
			},
		},
	}
	if tree == nil {
		return fs
	}
	collector := &symbolCollector{
		fileSymbols:  fs,
		currentScope: fs.GlobalScope,
	}
	antlr.ParseTreeWalkerDefault.Walk(collector, tree)
	return fs
}

// DeviceOrGlobalSymbol returns DeviceSymbol if present, or creates a synthetic
// device symbol representing the global scope and device templates for files
// that omit a 'device <name>;' declaration.
func (fs *FileSymbols) DeviceOrGlobalSymbol() *Symbol {
	if fs == nil {
		return nil
	}
	if fs.DeviceSymbol != nil {
		return fs.DeviceSymbol
	}
	sym := &Symbol{
		Name:      "dev",
		Kind:      SymbolKindObject,
		URI:       fs.URI,
		Children:  make(map[string]*Symbol),
		Templates: fs.DeviceTemplates,
	}
	if fs.GlobalScope != nil {
		for name, s := range fs.GlobalScope.Symbols {
			if isDeviceMemberKind(s.Kind) {
				sym.Children[name] = s
			}
		}
	}
	return sym
}
