// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lsp

import (
	"context"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"go.lsp.dev/protocol"

	"github.com/timoschenko/dmltools/pkg/dmlex"
)

const (
	semTokComment    uint32 = 0
	semTokKeyword    uint32 = 1
	semTokModifier   uint32 = 2
	semTokString     uint32 = 3
	semTokNumber     uint32 = 4
	semTokOperator   uint32 = 5
	semTokInterface  uint32 = 6
	semTokClass      uint32 = 7
	semTokProperty   uint32 = 8
	semTokMethod     uint32 = 9
	semTokFunction   uint32 = 10
	semTokParameter  uint32 = 11
	semTokVariable   uint32 = 12
	semTokEnumMember uint32 = 13
	semTokEvent      uint32 = 14
	semTokType       uint32 = 15
	semTokStruct     uint32 = 16
	semTokLabel      uint32 = 17
)

const (
	semModDeclaration    uint32 = 1 << 0
	semModDefinition     uint32 = 1 << 1
	semModReadonly       uint32 = 1 << 2
	semModStatic         uint32 = 1 << 3
	semModAbstract       uint32 = 1 << 4
	semModDefaultLibrary uint32 = 1 << 5
)

type semanticSpan struct {
	line      uint32
	startChar uint32
	length    uint32
	tokenType uint32
	modifiers uint32
}

type tokenClass struct {
	tokenType uint32
	modifiers uint32
	active    bool
}

// SemanticTokensLegend returns the semantic token types and modifiers legend
// advertised during LSP initialization.
func SemanticTokensLegend() protocol.SemanticTokensLegend {
	return protocol.SemanticTokensLegend{
		TokenTypes: []string{
			"comment",
			"keyword",
			"modifier",
			"string",
			"number",
			"operator",
			"interface",
			"class",
			"property",
			"method",
			"function",
			"parameter",
			"variable",
			"enumMember",
			"event",
			"type",
			"struct",
			"label",
		},
		TokenModifiers: []string{
			"declaration",
			"definition",
			"readonly",
			"static",
			"abstract",
			"defaultLibrary",
		},
	}
}

// ComputeSemanticTokens computes encoded LSP semantic tokens for the full
// document identified by uriStr.
func ComputeSemanticTokens(
	ctx context.Context,
	fsCache *FsCache,
	uriStr string,
) (*protocol.SemanticTokens, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	spans, err := fsCache.getOrComputeSemanticSpans(uriStr)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return &protocol.SemanticTokens{
		Data: encodeSemanticSpans(spans),
	}, nil
}

// ComputeSemanticTokensRange computes encoded LSP semantic tokens overlapping
// rng for the document identified by uriStr.
func ComputeSemanticTokensRange(
	ctx context.Context,
	fsCache *FsCache,
	uriStr string,
	rng protocol.Range,
) (*protocol.SemanticTokens, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	spans, err := fsCache.getOrComputeSemanticSpans(uriStr)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	filtered := filterSpansInRange(spans, rng)
	return &protocol.SemanticTokens{
		Data: encodeSemanticSpans(filtered),
	}, nil
}

func buildSemanticSpans(
	tokens *antlr.CommonTokenStream,
	tree antlr.ParseTree,
	syms *FileSymbols,
) []semanticSpan {
	if tokens == nil {
		return nil
	}
	tokens.Fill()
	allTokens := tokens.GetAllTokens()
	if len(allTokens) == 0 {
		return nil
	}

	classes := make([]tokenClass, len(allTokens))
	slotMap := make(map[int]int, len(allTokens))
	for i, tok := range allTokens {
		if tok == nil {
			continue
		}
		tokIdx := tok.GetTokenIndex()
		if tokIdx >= 0 && tokIdx != i {
			slotMap[tokIdx] = i
		}
		classes[i] = classifyLexerToken(tok)
	}

	if tree != nil && !isDML12Tree(tree) {
		collector := &semanticTokenCollector{
			allTokens: allTokens,
			classes:   classes,
			slotMap:   slotMap,
			syms:      syms,
		}
		antlr.ParseTreeWalkerDefault.Walk(collector, tree)
	}

	spans := make([]semanticSpan, 0, len(allTokens))
	for i, tok := range allTokens {
		cls := classes[i]
		if !cls.active || tok == nil {
			continue
		}
		spans = appendTokenSpans(spans, tok, cls)
	}
	return spans
}

func isDML12Tree(tree antlr.ParseTree) bool {
	dmlCtx, ok := tree.(*dmlex.DmlContext)
	if !ok || dmlCtx.DmlTag() == nil {
		return false
	}
	fl := dmlCtx.DmlTag().FLOAT_LITERAL()
	return fl != nil && strings.TrimSpace(fl.GetText()) == "1.2"
}

func classifyLexerToken(tok antlr.Token) tokenClass {
	tokType := dmlex.GetTokenType(tok)
	if tokType <= 0 {
		return tokenClass{}
	}

	switch tokType {
	case dmlex.TokenWs,
		dmlex.TokenCCodeBlock,
		dmlex.TokenT__0,
		dmlex.TokenLparen,
		dmlex.TokenRparen,
		dmlex.TokenLbrace,
		dmlex.TokenLbraceNew,
		dmlex.TokenRbrace,
		dmlex.TokenLbrack,
		dmlex.TokenRbrack,
		dmlex.TokenSemi,
		dmlex.TokenComma,
		dmlex.TokenDot:
		return tokenClass{}

	case dmlex.TokenComment,
		dmlex.TokenCommentBlock,
		dmlex.TokenCommentBlockEol,
		dmlex.TokenCommentInline,
		dmlex.TokenCommentInlineEol:
		return tokenClass{tokenType: semTokComment, active: true}

	case dmlex.TokenStringLiteral, dmlex.TokenCharLiteral:
		return tokenClass{tokenType: semTokString, active: true}

	case dmlex.TokenIntegerLiteral,
		dmlex.TokenHexLiteral,
		dmlex.TokenBinaryLiteral,
		dmlex.TokenFloatLiteral:
		return tokenClass{tokenType: semTokNumber, active: true}

	case dmlex.TokenThis:
		return tokenClass{
			tokenType: semTokVariable,
			modifiers: semModDefaultLibrary,
			active:    true,
		}

	case dmlex.TokenInline,
		dmlex.TokenShared,
		dmlex.TokenIndependent,
		dmlex.TokenStartup,
		dmlex.TokenMemoized,
		dmlex.TokenThrows,
		dmlex.TokenNothrow,
		dmlex.TokenConst,
		dmlex.TokenStatic,
		dmlex.TokenVect,
		dmlex.TokenSigned,
		dmlex.TokenUnsigned,
		dmlex.TokenVolatile,
		dmlex.TokenRestrict,
		dmlex.TokenPublic,
		dmlex.TokenPrivate,
		dmlex.TokenProtected,
		dmlex.TokenVirtual,
		dmlex.TokenAsync,
		dmlex.TokenAwait:
		return tokenClass{tokenType: semTokModifier, active: true}

	case dmlex.TokenInt,
		dmlex.TokenChar,
		dmlex.TokenFloat,
		dmlex.TokenDouble,
		dmlex.TokenShort,
		dmlex.TokenLong,
		dmlex.TokenVoid:
		return tokenClass{
			tokenType: semTokType,
			modifiers: semModDefaultLibrary,
			active:    true,
		}

	case dmlex.TokenIdentifier:
		txt := tok.GetText()
		switch {
		case txt == "true" || txt == "false" || txt == "NULL":
			return tokenClass{tokenType: semTokKeyword, active: true}
		case txt == "parent" || txt == "dev":
			return tokenClass{
				tokenType: semTokVariable,
				modifiers: semModDefaultLibrary,
				active:    true,
			}
		case isBuiltinTypeName(txt):
			return tokenClass{
				tokenType: semTokType,
				modifiers: semModDefaultLibrary,
				active:    true,
			}
		default:
			return tokenClass{}
		}
	}

	if tokType == dmlex.TokenWith || tokType.IsKeyword() {
		return tokenClass{tokenType: semTokKeyword, active: true}
	}

	if isOperatorToken(tokType) {
		return tokenClass{tokenType: semTokOperator, active: true}
	}

	return tokenClass{}
}

func isOperatorToken(t dmlex.TokenType) bool {
	switch t {
	case dmlex.TokenEllipsis,
		dmlex.TokenAssign,
		dmlex.TokenNewAssign,
		dmlex.TokenGt,
		dmlex.TokenLt,
		dmlex.TokenBang,
		dmlex.TokenTilde,
		dmlex.TokenQuestion,
		dmlex.TokenColon,
		dmlex.TokenEqual,
		dmlex.TokenLe,
		dmlex.TokenGe,
		dmlex.TokenNotequal,
		dmlex.TokenAnd,
		dmlex.TokenOr,
		dmlex.TokenInc,
		dmlex.TokenDec,
		dmlex.TokenAdd,
		dmlex.TokenSub,
		dmlex.TokenMul,
		dmlex.TokenDiv,
		dmlex.TokenBitand,
		dmlex.TokenBitor,
		dmlex.TokenCaret,
		dmlex.TokenMod,
		dmlex.TokenAddAssign,
		dmlex.TokenSubAssign,
		dmlex.TokenMulAssign,
		dmlex.TokenDivAssign,
		dmlex.TokenModAssign,
		dmlex.TokenAndAssign,
		dmlex.TokenOrAssign,
		dmlex.TokenXorAssign,
		dmlex.TokenLshiftAssign,
		dmlex.TokenRshiftAssign,
		dmlex.TokenLshift,
		dmlex.TokenRshift,
		dmlex.TokenArrow,
		dmlex.TokenAt,
		dmlex.TokenAltQuestion,
		dmlex.TokenAltColon:
		return true
	default:
		return false
	}
}

func isBuiltinTypeName(name string) bool {
	switch name {
	case "bool", "size_t", "ssize_t", "uintptr_t", "intptr_t":
		return true
	case "uint8_t", "uint16_t", "uint32_t", "uint64_t",
		"int8_t", "int16_t", "int32_t", "int64_t":
		return true
	}
	return isBitWidthIntType(name)
}

func isBitWidthIntType(name string) bool {
	var rest string
	switch {
	case strings.HasPrefix(name, "uint"):
		rest = name[4:]
	case strings.HasPrefix(name, "int"):
		rest = name[3:]
	default:
		return false
	}
	if len(rest) == 0 || rest[0] == '0' {
		return false
	}
	bits, err := strconv.Atoi(rest)
	return err == nil && bits >= 1 && bits <= 64
}

func appendTokenSpans(
	spans []semanticSpan,
	tok antlr.Token,
	cls tokenClass,
) []semanticSpan {
	if tok.GetLine() <= 0 || tok.GetStop() < tok.GetStart() {
		return spans
	}
	text := tok.GetText()
	if text == "" {
		return spans
	}

	baseLine := uint32(tok.GetLine() - 1) //nolint:gosec // Line >= 1 checked
	firstCol := tokenUTF16Col(tok)
	isComment := cls.tokenType == semTokComment

	if !strings.Contains(text, "\n") {
		seg := strings.TrimRight(text, "\r")
		if isComment {
			seg = strings.TrimRight(seg, " \t")
		}
		if segLen := utf16Len(seg); segLen > 0 {
			spans = append(spans, semanticSpan{
				line:      baseLine,
				startChar: firstCol,
				length:    segLen,
				tokenType: cls.tokenType,
				modifiers: cls.modifiers,
			})
		}
		return spans
	}

	lines := strings.Split(text, "\n")
	for i, rawLine := range lines {
		seg := strings.TrimRight(rawLine, "\r")
		if isComment {
			seg = strings.TrimRight(seg, " \t")
		}
		segLen := utf16Len(seg)
		if segLen == 0 {
			continue
		}
		col := uint32(0)
		if i == 0 {
			col = firstCol
		}
		spans = append(spans, semanticSpan{
			line:      baseLine + uint32(i), //nolint:gosec // i >= 0
			startChar: col,
			length:    segLen,
			tokenType: cls.tokenType,
			modifiers: cls.modifiers,
		})
	}
	return spans
}

func filterSpansInRange(
	spans []semanticSpan,
	rng protocol.Range,
) []semanticSpan {
	var out []semanticSpan
	for _, sp := range spans {
		spStart := protocol.Position{Line: sp.line, Character: sp.startChar}
		spEnd := protocol.Position{
			Line:      sp.line,
			Character: sp.startChar + sp.length,
		}
		if positionLess(rng.Start, spEnd) && positionLess(spStart, rng.End) {
			out = append(out, sp)
		}
	}
	return out
}

func positionLess(a, b protocol.Position) bool {
	if a.Line != b.Line {
		return a.Line < b.Line
	}
	return a.Character < b.Character
}

func encodeSemanticSpans(spans []semanticSpan) []uint32 {
	data := make([]uint32, 0, len(spans)*5)
	var prevLine uint32
	var prevChar uint32
	hasPrev := false

	for _, sp := range spans {
		if sp.length == 0 {
			continue
		}
		if hasPrev {
			if sp.line < prevLine {
				continue
			}
			if sp.line == prevLine && sp.startChar < prevChar {
				continue
			}
		}
		deltaLine := sp.line - prevLine
		deltaStart := sp.startChar
		if deltaLine == 0 && hasPrev {
			deltaStart = sp.startChar - prevChar
		}
		data = append(
			data,
			deltaLine,
			deltaStart,
			sp.length,
			sp.tokenType,
			sp.modifiers,
		)
		prevLine = sp.line
		prevChar = sp.startChar
		hasPrev = true
	}
	return data
}

type semanticTokenCollector struct {
	*dmlex.BasedmlListener
	allTokens []antlr.Token
	classes   []tokenClass
	slotMap   map[int]int
	syms      *FileSymbols
}

func (c *semanticTokenCollector) slotForToken(tok antlr.Token) (int, bool) {
	if tok == nil {
		return 0, false
	}
	tokIdx := tok.GetTokenIndex()
	if tokIdx < 0 {
		return 0, false
	}
	if slot, ok := c.slotMap[tokIdx]; ok {
		if c.allTokens[slot] == tok {
			return slot, true
		}
		return 0, false
	}
	if tokIdx < len(c.allTokens) && c.allTokens[tokIdx] == tok {
		return tokIdx, true
	}
	return 0, false
}

func (c *semanticTokenCollector) setTokenClass(
	tok antlr.Token,
	tokenType uint32,
	modifiers uint32,
) {
	slot, ok := c.slotForToken(tok)
	if !ok {
		return
	}
	c.classes[slot] = tokenClass{
		tokenType: tokenType,
		modifiers: modifiers,
		active:    true,
	}
}

func (c *semanticTokenCollector) clearToken(tok antlr.Token) {
	slot, ok := c.slotForToken(tok)
	if !ok {
		return
	}
	c.classes[slot] = tokenClass{}
}

func (c *semanticTokenCollector) isActiveToken(tok antlr.Token) bool {
	slot, ok := c.slotForToken(tok)
	return ok && c.classes[slot].active
}

func identOrObjidentToken(node antlr.Tree) antlr.Token {
	return findIdentTokenInTree(node)
}

func tokenStartPos(tok antlr.Token) protocol.Position {
	line := max(tok.GetLine()-1, 0)
	return protocol.Position{
		Line:      uint32(line), //nolint:gosec // line >= 0 checked
		Character: tokenUTF16Col(tok),
	}
}

func (c *semanticTokenCollector) EnterIdent(ctx *dmlex.IdentContext) {
	tok := ctx.GetStart()
	if tok == nil {
		return
	}
	tokType := dmlex.GetTokenType(tok)
	if tokType != dmlex.TokenIdentifier {
		c.clearToken(tok)
		return
	}
	if isBuiltinTypeName(tok.GetText()) {
		c.clearToken(tok)
	}
}

func (c *semanticTokenCollector) EnterObjident(ctx *dmlex.ObjidentContext) {
	if ctx.REGISTER() != nil {
		c.clearToken(ctx.REGISTER().GetSymbol())
	}
}

func (c *semanticTokenCollector) ExitDeviceTag(ctx *dmlex.DeviceTagContext) {
	if tok := identOrObjidentToken(ctx.Objident()); tok != nil {
		c.setTokenClass(tok, semTokClass, semModDefinition)
	}
}

func (c *semanticTokenCollector) ExitTemplateStatement(
	ctx *dmlex.TemplateStatementContext,
) {
	if tok := identOrObjidentToken(ctx.Objident()); tok != nil {
		c.setTokenClass(tok, semTokInterface, semModDefinition)
	}
}

func (c *semanticTokenCollector) ExitObject_decl(
	ctx *dmlex.Object_declContext,
) {
	if tok := identOrObjidentToken(ctx.Objident()); tok != nil {
		c.setTokenClass(tok, semTokClass, semModDefinition)
	}
}

func (c *semanticTokenCollector) ExitParam(ctx *dmlex.ParamContext) {
	tok := identOrObjidentToken(ctx.Objident())
	if tok == nil {
		return
	}
	mods := semModDeclaration
	if ctx.ASSIGN() != nil || ctx.NEW_ASSIGN() != nil ||
		ctx.DEFAULT() != nil || ctx.DEFAULT_NEW() != nil {
		mods = semModDefinition
	}
	c.setTokenClass(tok, semTokProperty, mods)
}

func (c *semanticTokenCollector) ExitMethodUnifiedDeclaration(
	ctx *dmlex.MethodUnifiedDeclarationContext,
) {
	if tok := identOrObjidentToken(ctx.Objident()); tok != nil {
		mods := semModDeclaration
		if ctx.MethodQualifiers() != nil &&
			ctx.MethodQualifiers().SHARED() != nil {
			mods |= semModAbstract
		}
		c.setTokenClass(tok, semTokMethod, mods)
	}
	c.classifyMethodOutParams(dmlex.MethodReturns(ctx))
}

func (c *semanticTokenCollector) ExitMethodUnifiedDefinition(
	ctx *dmlex.MethodUnifiedDefinitionContext,
) {
	if tok := identOrObjidentToken(ctx.Objident()); tok != nil {
		c.setTokenClass(tok, semTokMethod, semModDefinition)
	}
	c.classifyMethodOutParams(dmlex.MethodReturns(ctx))
}

func (c *semanticTokenCollector) classifyMethodOutParams(
	outParams []dmlex.ICdeclContext,
) {
	for _, outParam := range outParams {
		if outParam == nil {
			continue
		}
		tok := findCdeclIdentToken(outParam)
		if tok == nil || tok.GetText() == "_" {
			continue
		}
		mods := semModDeclaration
		if isConstCdecl(outParam) {
			mods |= semModReadonly
		}
		c.setTokenClass(tok, semTokParameter, mods)
	}
}

func (c *semanticTokenCollector) ExitCdeclOrInline(
	ctx *dmlex.CdeclOrInlineContext,
) {
	var tok antlr.Token
	mods := semModDeclaration
	if ctx.Ident() != nil {
		tok = findIdentTokenInTree(ctx.Ident())
	} else if ctx.Cdecl() != nil {
		cd := ctx.Cdecl()
		if isConstCdecl(cd) {
			mods |= semModReadonly
		}
		tok = findCdeclIdentToken(cd)
	}
	if tok == nil || tok.GetText() == "_" {
		return
	}
	c.setTokenClass(tok, semTokParameter, mods)
}

func (c *semanticTokenCollector) ExitCdeclFuncArgs(
	ctx *dmlex.CdeclFuncArgsContext,
) {
	v := ctx.Cdecl_list_variadic()
	if v == nil || v.Cdecl_list() == nil {
		return
	}
	for _, cd := range v.Cdecl_list().AllCdecl() {
		if cd == nil {
			continue
		}
		tok := findCdeclIdentToken(cd)
		if tok == nil || tok.GetText() == "_" {
			continue
		}
		mods := semModDeclaration
		if isConstCdecl(cd) {
			mods |= semModReadonly
		}
		c.setTokenClass(tok, semTokParameter, mods)
	}
}

func (c *semanticTokenCollector) classifyCvarDecl(
	cvar dmlex.ICvarDeclContext,
	isStatic bool,
	isDef bool,
) {
	if cvar == nil {
		return
	}
	classifyOne := func(cd dmlex.ICdeclContext) {
		if cd == nil {
			return
		}
		tok := findCdeclIdentToken(cd)
		if tok == nil || tok.GetText() == "_" {
			return
		}
		mods := semModDeclaration
		if isDef {
			mods = semModDefinition
		}
		if isStatic {
			mods |= semModStatic
		}
		isFunc := isDirectFuncCdecl(cd.Cdecl_qualifier_ident_array_func())
		if isConstCdecl(cd) && !isFunc {
			mods |= semModReadonly
		}
		tokType := semTokVariable
		if isFunc {
			tokType = semTokFunction
		}
		c.setTokenClass(tok, tokType, mods)
	}

	classifyOne(cvar.Cdecl())
	if cvar.Cdecl_list() != nil {
		for _, cd := range cvar.Cdecl_list().AllCdecl() {
			classifyOne(cd)
		}
	}
}

func (c *semanticTokenCollector) ExitVariableDecl(
	ctx *dmlex.VariableDeclContext,
) {
	isStatic := false
	if lk := ctx.Local_decl_kind(); lk != nil {
		if lk.SESSION() != nil || lk.SAVED() != nil || lk.STATIC() != nil {
			isStatic = true
		}
	}
	isDef := ctx.CvarDecl() != nil && ctx.CvarDecl().ASSIGN() != nil
	c.classifyCvarDecl(ctx.CvarDecl(), isStatic, isDef)
}

func (c *semanticTokenCollector) ExitSessionOrSavedDecl(
	ctx *dmlex.SessionOrSavedDeclContext,
) {
	isDef := ctx.CvarDecl() != nil && ctx.CvarDecl().ASSIGN() != nil
	c.classifyCvarDecl(ctx.CvarDecl(), true, isDef)
}

func (c *semanticTokenCollector) ExitConstantStatement(
	ctx *dmlex.ConstantStatementContext,
) {
	if tok := findIdentTokenInTree(ctx.Ident()); tok != nil {
		c.setTokenClass(
			tok,
			semTokEnumMember,
			semModDefinition|semModReadonly,
		)
	}
}

func (c *semanticTokenCollector) ExitLogGroupStatement(
	ctx *dmlex.LogGroupStatementContext,
) {
	if tok := findIdentTokenInTree(ctx.Ident()); tok != nil {
		c.setTokenClass(tok, semTokEnumMember, semModDefinition)
	}
}

func (c *semanticTokenCollector) ExitTypedefStatement(
	ctx *dmlex.TypedefStatementContext,
) {
	if ctx.Cdecl() == nil {
		return
	}
	tok := findCdeclIdentToken(ctx.Cdecl())
	if tok == nil || tok.GetText() == "_" {
		return
	}
	tokType := semTokType
	if isStructLikeBasetype(ctx.Cdecl().Basetype()) {
		tokType = semTokStruct
	}
	c.setTokenClass(tok, tokType, semModDefinition)
}

func (c *semanticTokenCollector) ExitExternStatement(
	ctx *dmlex.ExternStatementContext,
) {
	if ctx.Cdecl() == nil {
		return
	}
	tok := findCdeclIdentToken(ctx.Cdecl())
	if tok == nil || tok.GetText() == "_" {
		return
	}
	if ctx.TYPEDEF() != nil {
		tokType := semTokType
		if isStructLikeBasetype(ctx.Cdecl().Basetype()) {
			tokType = semTokStruct
		}
		c.setTokenClass(tok, tokType, semModDeclaration)
		return
	}
	isFunc := isDirectFuncCdecl(ctx.Cdecl().Cdecl_qualifier_ident_array_func())
	tokType := semTokVariable
	if isFunc {
		tokType = semTokFunction
	}
	mods := semModDeclaration
	if isConstCdecl(ctx.Cdecl()) && !isFunc {
		mods |= semModReadonly
	}
	c.setTokenClass(tok, tokType, mods)
}

func (c *semanticTokenCollector) ExitHook_decl(ctx *dmlex.Hook_declContext) {
	if tok := findIdentTokenInTree(ctx.Ident()); tok != nil {
		c.setTokenClass(tok, semTokEvent, semModDeclaration)
	}
	if ctx.Cdecl_list() != nil {
		for _, cd := range ctx.Cdecl_list().AllCdecl() {
			if cd == nil {
				continue
			}
			tok := findCdeclIdentToken(cd)
			if tok == nil || tok.GetText() == "_" {
				continue
			}
			mods := semModDeclaration
			if isConstCdecl(cd) {
				mods |= semModReadonly
			}
			c.setTokenClass(tok, semTokParameter, mods)
		}
	}
}

func (c *semanticTokenCollector) ExitStruct(ctx *dmlex.StructContext) {
	if ctx.Ident() == nil {
		return
	}
	if tok := findIdentTokenInTree(ctx.Ident()); tok != nil {
		var mods uint32
		if ctx.LBRACE() != nil {
			mods = semModDefinition
		}
		c.setTokenClass(tok, semTokStruct, mods)
	}
}

func (c *semanticTokenCollector) ExitUnion(ctx *dmlex.UnionContext) {
	if ctx.Ident() == nil {
		return
	}
	if tok := findIdentTokenInTree(ctx.Ident()); tok != nil {
		var mods uint32
		if ctx.LBRACE() != nil {
			mods = semModDefinition
		}
		c.setTokenClass(tok, semTokStruct, mods)
	}
}

func (c *semanticTokenCollector) classifyMemberFieldDecl(
	cd dmlex.ICdeclContext,
) {
	if cd == nil {
		return
	}
	tok := findCdeclIdentToken(cd)
	if tok == nil || tok.GetText() == "_" {
		return
	}
	mods := semModDeclaration
	if isConstCdecl(cd) {
		mods |= semModReadonly
	}
	c.setTokenClass(tok, semTokProperty, mods)
}

func (c *semanticTokenCollector) ExitStructDecls(
	ctx *dmlex.StructDeclsContext,
) {
	c.classifyMemberFieldDecl(ctx.Cdecl())
}

func (c *semanticTokenCollector) ExitLayoutDecls(
	ctx *dmlex.LayoutDeclsContext,
) {
	c.classifyMemberFieldDecl(ctx.Cdecl())
}

func (c *semanticTokenCollector) ExitBitfields_decls(
	ctx *dmlex.Bitfields_declsContext,
) {
	c.classifyMemberFieldDecl(ctx.Cdecl())
}

func (c *semanticTokenCollector) ExitArraydef(ctx *dmlex.ArraydefContext) {
	if tok := findIdentTokenInTree(ctx.Ident()); tok != nil &&
		tok.GetText() != "_" {
		c.setTokenClass(tok, semTokVariable, semModDefinition)
	}
}

func (c *semanticTokenCollector) ExitForeachStatement(
	ctx *dmlex.ForeachStatementContext,
) {
	if tok := findIdentTokenInTree(ctx.Ident()); tok != nil &&
		tok.GetText() != "_" {
		c.setTokenClass(tok, semTokVariable, semModDefinition)
	}
}

func (c *semanticTokenCollector) ExitSelectStatement(
	ctx *dmlex.SelectStatementContext,
) {
	if tok := findIdentTokenInTree(ctx.Ident()); tok != nil &&
		tok.GetText() != "_" {
		c.setTokenClass(tok, semTokVariable, semModDefinition)
	}
}

func (c *semanticTokenCollector) ExitAfterDelayCondition(
	ctx *dmlex.AfterDelayConditionContext,
) {
	if ctx.ARROW() != nil {
		classifyParam := func(id dmlex.IIdentContext) {
			tok := findIdentTokenInTree(id)
			if tok != nil && tok.GetText() != "_" {
				c.setTokenClass(tok, semTokParameter, semModDeclaration)
			}
		}
		classifyParam(ctx.Ident())
		if ctx.IdentList() != nil {
			for _, id := range ctx.IdentList().AllIdent() {
				classifyParam(id)
			}
		}
		return
	}
	if ctx.IDENTIFIER() != nil {
		c.setTokenClass(ctx.IDENTIFIER().GetSymbol(), semTokKeyword, 0)
	}
}

func (c *semanticTokenCollector) ExitIstemplateList(
	ctx *dmlex.IstemplateListContext,
) {
	for _, objID := range ctx.AllObjident() {
		if tok := identOrObjidentToken(objID); tok != nil {
			c.setTokenClass(tok, semTokInterface, 0)
		}
	}
}

func (c *semanticTokenCollector) ExitEachExpr(ctx *dmlex.EachExprContext) {
	if tok := identOrObjidentToken(ctx.Objident()); tok != nil {
		c.setTokenClass(tok, semTokInterface, 0)
	}
}

func (c *semanticTokenCollector) ExitBasetype(ctx *dmlex.BasetypeContext) {
	if ctx.HOOK() != nil && ctx.Cdecl_list() != nil {
		for _, cd := range ctx.Cdecl_list().AllCdecl() {
			if cd == nil {
				continue
			}
			tok := findCdeclIdentToken(cd)
			if tok == nil || tok.GetText() == "_" {
				continue
			}
			mods := semModDeclaration
			if isConstCdecl(cd) {
				mods |= semModReadonly
			}
			c.setTokenClass(tok, semTokParameter, mods)
		}
	}
	if ctx.SEQUENCE() != nil && ctx.Signed_or_unsigned_typeident() != nil {
		if ti := ctx.Signed_or_unsigned_typeident().Typeident(); ti != nil {
			if tok := findIdentTokenInTree(ti); tok != nil {
				c.setTokenClass(tok, semTokInterface, 0)
			}
		}
	}
}

func (c *semanticTokenCollector) ExitSigned_or_unsigned_typeident(
	ctx *dmlex.Signed_or_unsigned_typeidentContext,
) {
	if ctx.Typeident() != nil && ctx.Typeident().Ident() == nil {
		return
	}
	if ctx.SIGNED() != nil {
		c.setTokenClass(
			ctx.SIGNED().GetSymbol(),
			semTokType,
			semModDefaultLibrary,
		)
	} else if ctx.UNSIGNED() != nil {
		c.setTokenClass(
			ctx.UNSIGNED().GetSymbol(),
			semTokType,
			semModDefaultLibrary,
		)
	}
}

func (c *semanticTokenCollector) ExitTypeident(ctx *dmlex.TypeidentContext) {
	switch {
	case ctx.INT() != nil,
		ctx.CHAR() != nil,
		ctx.DOUBLE() != nil,
		ctx.FLOAT() != nil,
		ctx.LONG() != nil,
		ctx.SHORT() != nil,
		ctx.VOID() != nil,
		ctx.SIGNED() != nil,
		ctx.UNSIGNED() != nil:
		c.setTokenClass(ctx.GetStart(), semTokType, semModDefaultLibrary)
	case ctx.REGISTER() != nil:
		c.setTokenClass(ctx.REGISTER().GetSymbol(), semTokType, 0)
	case ctx.Ident() != nil:
		tok := findIdentTokenInTree(ctx.Ident())
		if tok == nil {
			return
		}
		name := tok.GetText()
		if isBuiltinTypeName(name) {
			c.setTokenClass(tok, semTokType, semModDefaultLibrary)
			return
		}
		sym := lookupLocalTypeSymbol(name, tokenStartPos(tok), c.syms)
		if sym != nil && sym.Kind == SymbolKindStruct {
			c.setTokenClass(tok, semTokStruct, 0)
			return
		}
		c.setTokenClass(tok, semTokType, 0)
	}
}

func (c *semanticTokenCollector) ExitGotoStatement(
	ctx *dmlex.GotoStatementContext,
) {
	if tok := findIdentTokenInTree(ctx.Ident()); tok != nil {
		c.setTokenClass(tok, semTokLabel, 0)
	}
}

func (c *semanticTokenCollector) ExitLogKind(ctx *dmlex.LogKindContext) {
	if ctx.IDENTIFIER() != nil {
		c.setTokenClass(ctx.IDENTIFIER().GetSymbol(), semTokKeyword, 0)
	}
}

func (c *semanticTokenCollector) ExitLogGroup(ctx *dmlex.LogGroupContext) {
	la, ok := ctx.Expression().(*dmlex.LiteralAtomContext)
	if !ok {
		return
	}
	ei := la.ExpressionIdent()
	if ei == nil || ei.ExpressionAccessIdent() == nil {
		return
	}
	ea := ei.ExpressionAccessIdent()
	if ea.Objident() == nil {
		return
	}
	tok := identOrObjidentToken(ea.Objident())
	if tok == nil || c.isActiveToken(tok) {
		return
	}
	c.setTokenClass(tok, semTokEnumMember, 0)
}

func (c *semanticTokenCollector) ExitDesignatedStructInitializer(
	ctx *dmlex.DesignatedStructInitializerContext,
) {
	if tok := findIdentTokenInTree(ctx.Ident()); tok != nil {
		c.setTokenClass(tok, semTokProperty, 0)
	}
}

func (c *semanticTokenCollector) ExitBitorderTag(
	ctx *dmlex.BitorderTagContext,
) {
	if ctx.ENDIAN() != nil {
		c.setTokenClass(ctx.ENDIAN().GetSymbol(), semTokKeyword, 0)
	}
}

func (c *semanticTokenCollector) ExitEndianflag(
	ctx *dmlex.EndianflagContext,
) {
	if ctx.Ident() != nil {
		if tok := findIdentTokenInTree(ctx.Ident()); tok != nil {
			c.setTokenClass(tok, semTokKeyword, 0)
		}
	}
}

func (c *semanticTokenCollector) ExitProvisionalTag(
	ctx *dmlex.ProvisionalTagContext,
) {
	if ctx.IdentList() == nil {
		return
	}
	for _, id := range ctx.IdentList().AllIdent() {
		if tok := findIdentTokenInTree(id); tok != nil {
			c.setTokenClass(tok, semTokProperty, 0)
		}
	}
}

func isCallTarget(ctx antlr.ParserRuleContext) bool {
	fc, ok := ctx.GetParent().(*dmlex.FunctionCallExprContext)
	return ok && fc.Expression() == ctx
}

func (c *semanticTokenCollector) ExitLiteralAtom(
	ctx *dmlex.LiteralAtomContext,
) {
	if _, ok := ctx.GetParent().(*dmlex.MemberAccessExprContext); ok {
		return
	}
	c.resolveRootAtom(ctx, isCallTarget(ctx))
}

func (c *semanticTokenCollector) ExitMemberAccessExpr(
	ctx *dmlex.MemberAccessExprContext,
) {
	curSym := c.resolveExprSymbol(ctx.Expression())

	isCall := isCallTarget(ctx)
	members := ctx.AllExpressionAccessIdent()
	lastIdx := len(members) - 1
	for i, ea := range members {
		if ea == nil {
			continue
		}
		if ea.THIS() != nil || ea.DEFAULT() != nil || ea.TEMPLATE() != nil {
			continue
		}
		if ea.Objident() == nil {
			curSym = nil
			continue
		}
		tok := identOrObjidentToken(ea.Objident())
		if tok == nil {
			curSym = nil
			continue
		}
		name := tok.GetText()
		var memSym *Symbol
		if curSym != nil {
			memSym = lookupLocalMember(curSym, name, c.syms)
		}
		if memSym != nil {
			tokType, mods := classifySymbolRef(memSym)
			c.setTokenClass(tok, tokType, mods)
			curSym = memSym
			continue
		}
		if i == lastIdx && isCall {
			c.setTokenClass(tok, semTokMethod, 0)
		} else {
			c.setTokenClass(tok, semTokProperty, 0)
		}
		curSym = nil
	}
}

func (c *semanticTokenCollector) resolveExprSymbol(
	expr dmlex.IExpressionContext,
) *Symbol {
	switch e := expr.(type) {
	case *dmlex.LiteralAtomContext:
		return c.resolveRootAtom(e, false)
	case *dmlex.ArrayAccessExprContext:
		return c.resolveExprSymbol(e.Expression(0))
	case *dmlex.MemberAccessExprContext:
		cur := c.resolveExprSymbol(e.Expression())
		for _, ea := range e.AllExpressionAccessIdent() {
			if cur == nil || ea == nil || ea.Objident() == nil {
				return nil
			}
			tok := identOrObjidentToken(ea.Objident())
			if tok == nil {
				return nil
			}
			cur = lookupLocalMember(cur, tok.GetText(), c.syms)
		}
		return cur
	default:
		return nil
	}
}

func (c *semanticTokenCollector) resolveRootAtom(
	ctx *dmlex.LiteralAtomContext,
	isCall bool,
) *Symbol {
	if ctx == nil || ctx.ExpressionIdent() == nil {
		return nil
	}
	ea := ctx.ExpressionIdent().ExpressionAccessIdent()
	if ea == nil {
		return nil
	}
	pos := tokenStartPos(ea.GetStart())
	var enclosingScope *Scope
	if c.syms != nil && c.syms.GlobalScope != nil {
		enclosingScope = c.syms.GlobalScope.EnclosingScope(pos)
	}

	if ea.THIS() != nil {
		return findEnclosingObject(enclosingScope, c.syms)
	}
	if ea.DEFAULT() != nil || ea.TEMPLATE() != nil {
		return nil
	}
	if ea.Objident() == nil {
		return nil
	}
	tok := identOrObjidentToken(ea.Objident())
	if tok == nil {
		return nil
	}
	name := tok.GetText()
	sym := lookupLocalSymbol(name, enclosingScope, c.syms)
	if sym != nil {
		tokType, mods := classifySymbolRef(sym)
		c.setTokenClass(tok, tokType, mods)
		return sym
	}
	if name == "dev" {
		c.setTokenClass(tok, semTokVariable, semModDefaultLibrary)
		if c.syms != nil {
			return c.syms.DeviceOrGlobalSymbol()
		}
		return nil
	}
	if name == "parent" {
		c.setTokenClass(tok, semTokVariable, semModDefaultLibrary)
		if enc := findEnclosingObject(enclosingScope, c.syms); enc != nil {
			return enc.Parent
		}
		return nil
	}
	if dmlex.GetTokenType(tok) == dmlex.TokenIdentifier &&
		(name == "true" || name == "false" || name == "NULL") {
		c.setTokenClass(tok, semTokKeyword, 0)
		return nil
	}
	if isCall {
		c.setTokenClass(tok, semTokFunction, 0)
	}
	return nil
}

func classifySymbolRef(sym *Symbol) (uint32, uint32) {
	switch sym.Kind {
	case SymbolKindTemplate:
		return semTokInterface, 0
	case SymbolKindObject:
		return semTokClass, 0
	case SymbolKindParam:
		return semTokProperty, 0
	case SymbolKindMethod:
		return semTokMethod, 0
	case SymbolKindParameter:
		var mods uint32
		if sym.IsConst {
			mods |= semModReadonly
		}
		return semTokParameter, mods
	case SymbolKindVariable:
		tokType := semTokVariable
		if sym.IsFunc {
			tokType = semTokFunction
		}
		var mods uint32
		if sym.IsStatic {
			mods |= semModStatic
		}
		if sym.IsConst && !sym.IsFunc {
			mods |= semModReadonly
		}
		return tokType, mods
	case SymbolKindSession, SymbolKindSaved:
		tokType := semTokVariable
		if sym.IsFunc {
			tokType = semTokFunction
		}
		mods := semModStatic
		if sym.IsConst && !sym.IsFunc {
			mods |= semModReadonly
		}
		return tokType, mods
	case SymbolKindConstant:
		return semTokEnumMember, semModReadonly
	case SymbolKindLogGroup:
		return semTokEnumMember, 0
	case SymbolKindHook:
		return semTokEvent, 0
	case SymbolKindStruct:
		return semTokStruct, 0
	case SymbolKindTypedef:
		return semTokType, 0
	case SymbolKindExtern:
		tokType := semTokVariable
		if sym.IsFunc {
			tokType = semTokFunction
		}
		var mods uint32
		if sym.IsConst && !sym.IsFunc {
			mods |= semModReadonly
		}
		return tokType, mods
	default:
		return semTokVariable, 0
	}
}

func findEnclosingObject(scope *Scope, syms *FileSymbols) *Symbol {
	for s := scope; s != nil; s = s.Parent {
		if s.Kind == ScopeObject && s.Owner != nil {
			return s.Owner
		}
	}
	if syms != nil {
		return syms.DeviceOrGlobalSymbol()
	}
	return nil
}

func lookupLocalTypeSymbol(
	name string,
	pos protocol.Position,
	syms *FileSymbols,
) *Symbol {
	if syms == nil || syms.GlobalScope == nil || name == "" {
		return nil
	}
	for s := syms.GlobalScope.EnclosingScope(pos); s != nil; s = s.Parent {
		if sym, ok := s.Symbols[name]; ok {
			if sym.Kind == SymbolKindStruct || sym.Kind == SymbolKindTypedef {
				return sym
			}
		}
	}
	return nil
}

func lookupLocalSymbol(
	name string,
	scope *Scope,
	syms *FileSymbols,
) *Symbol {
	if syms == nil || name == "" {
		return nil
	}
	for s := scope; s != nil; s = s.Parent {
		if sym, ok := s.Symbols[name]; ok {
			return sym
		}
		if s.Kind == ScopeObject && s.Owner != nil {
			visited := make(map[string]bool)
			if sym := lookupLocalTemplates(
				s.Owner.Templates,
				name,
				syms,
				visited,
			); sym != nil {
				return sym
			}
		}
	}
	if syms.DeviceSymbol != nil {
		if sym, ok := syms.DeviceSymbol.Children[name]; ok {
			return sym
		}
		visited := make(map[string]bool)
		if sym := lookupLocalTemplates(
			syms.DeviceTemplates,
			name,
			syms,
			visited,
		); sym != nil {
			return sym
		}
	} else if len(syms.DeviceTemplates) > 0 {
		visited := make(map[string]bool)
		if sym := lookupLocalTemplates(
			syms.DeviceTemplates,
			name,
			syms,
			visited,
		); sym != nil {
			return sym
		}
	}
	return nil
}

func lookupLocalMember(
	sym *Symbol,
	member string,
	syms *FileSymbols,
) *Symbol {
	if sym == nil || member == "" {
		return nil
	}
	if child, ok := sym.Children[member]; ok {
		return child
	}
	if len(sym.Templates) > 0 && syms != nil {
		visited := make(map[string]bool)
		if child := lookupLocalTemplates(
			sym.Templates,
			member,
			syms,
			visited,
		); child != nil {
			return child
		}
	}
	return nil
}

func lookupLocalTemplates(
	templates []string,
	member string,
	syms *FileSymbols,
	visited map[string]bool,
) *Symbol {
	if syms == nil || syms.GlobalScope == nil {
		return nil
	}
	for _, tplName := range templates {
		if visited[tplName] {
			continue
		}
		visited[tplName] = true
		tplSym, ok := syms.GlobalScope.Symbols[tplName]
		if !ok || tplSym.Kind != SymbolKindTemplate {
			continue
		}
		if child, ok := tplSym.Children[member]; ok {
			return child
		}
		if child := lookupLocalTemplates(
			tplSym.Templates,
			member,
			syms,
			visited,
		); child != nil {
			return child
		}
	}
	return nil
}
