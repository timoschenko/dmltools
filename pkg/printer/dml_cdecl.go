// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer

import (
	"fmt"
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer/format"

	"github.com/antlr4-go/antlr/v4"
)

// g4-tag: typeof
func (p *Printer) printTypeof(
	node dmlex.ITypeofContext,
) error {
	// TODO: ensure the CdeclDefaults
	// TODO: Review AI output
	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				return dmlex.GetTokenType(
					child.GetSymbol(),
				) == dmlex.TokenTypeof
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IExpressionContext:
			if err := p.printExpression(child); err != nil {
				return err
			}
		default:
			panic("unreachable: printTypeof - " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: structDecls
func (p *Printer) printStructDecls(
	node dmlex.IStructDeclsContext,
) error {
	overrides := map[dmlex.TokenType]format.FormatFlag{
		dmlex.TokenColon: p.context.GetFormatFlag(dmlex.TokenColon).
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorPriority(format.PriorityHigh),
	}
	p.newContext(format.NewTokenMap(overrides))
	defer p.freeContext()

	// TODO: ensure the CdeclDefaults
	// TODO: Review AI output
	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenColon, dmlex.TokenSemi:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.ICdeclContext:
			if err := p.printCdecl(child); err != nil {
				return err
			}
		case dmlex.IExpressionContext:
			if err := p.printExpression(child); err != nil {
				return err
			}
		default:
			panic("unreachable: printStructDecls -" + fmt.Sprintf("%T", child))
		}
	}

	return p.Flush()
}

// g4-tag: struct
func (p *Printer) printStruct(
	node dmlex.IStructContext,
) error {
	// TODO: ensure the CdeclDefaults
	// TODO: Review AI output

	hadMember := false
	var prevStopToken antlr.Token
	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenStruct, dmlex.TokenLbrace, dmlex.TokenRbrace:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IIdentContext:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IStructDeclsContext:
			startIdx := child.GetStart().GetTokenIndex()
			var nextStartTok antlr.Token
			if p.lastProcessedTokenIndex < TokenIndex(startIdx) {
				nextStartTok = p.extractToken(p.lastProcessedTokenIndex)
			} else {
				nextStartTok = child.GetStart()
			}
			if hadMember {
				blankLines := countInputBlankLines(prevStopToken, nextStartTok)
				emptyLines := computeMemberSpacing(blankLines)
				if err := p.writeNewLines(emptyLines + 1); err != nil {
					return err
				}
			}
			if err := p.pushTokenRange(TokenIndex(startIdx - 1)); err != nil {
				return err
			}
			if err := p.printStructDecls(child); err != nil {
				return err
			}
			hadMember = true
			prevStopIdx := p.gatherTokenIndexFromRightiestToken(child.GetStop())
			prevStopToken = p.extractToken(prevStopIdx)
		default:
			panic("unreachable: printStruct - " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: union
func (p *Printer) printUnion(
	node dmlex.IUnionContext,
) error {
	// TODO: ensure the CdeclDefaults
	// TODO: Review AI output

	hadMember := false
	var prevStopToken antlr.Token
	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenUnion, dmlex.TokenLbrace, dmlex.TokenRbrace:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IIdentContext:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IStructDeclsContext:
			startIdx := child.GetStart().GetTokenIndex()
			var nextStartTok antlr.Token
			if p.lastProcessedTokenIndex < TokenIndex(startIdx) {
				nextStartTok = p.extractToken(p.lastProcessedTokenIndex)
			} else {
				nextStartTok = child.GetStart()
			}
			if hadMember {
				blankLines := countInputBlankLines(prevStopToken, nextStartTok)
				emptyLines := computeMemberSpacing(blankLines)
				if err := p.writeNewLines(emptyLines + 1); err != nil {
					return err
				}
			}
			if err := p.pushTokenRange(TokenIndex(startIdx - 1)); err != nil {
				return err
			}
			if err := p.printStructDecls(child); err != nil {
				return err
			}
			hadMember = true
			prevStopIdx := p.gatherTokenIndexFromRightiestToken(child.GetStop())
			prevStopToken = p.extractToken(prevStopIdx)
		default:
			panic("unreachable: printUnion - " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: layoutDecls
func (p *Printer) printLayoutDecls(
	node dmlex.ILayoutDeclsContext,
) error {
	// TODO: ensure the CdeclDefaults
	// TODO: Review AI output
	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				return dmlex.GetTokenType(child.GetSymbol()) == dmlex.TokenSemi
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.ICdeclContext:
			if err := p.printCdecl(child); err != nil {
				return err
			}
		default:
			panic("unreachable: printLayoutDecls - " + fmt.Sprintf("%T", child))
		}
	}

	return p.Flush()
}

// g4-tag: layout
func (p *Printer) printLayout(
	node dmlex.ILayoutContext,
) error {
	// TODO: ensure the CdeclDefaults
	// TODO: Review AI output

	hadMember := false
	var prevStopToken antlr.Token
	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenLayout,
					dmlex.TokenStringLiteral,
					dmlex.TokenLbrace,
					dmlex.TokenRbrace:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.ILayoutDeclsContext:
			startIdx := child.GetStart().GetTokenIndex()
			var nextStartTok antlr.Token
			if p.lastProcessedTokenIndex < TokenIndex(startIdx) {
				nextStartTok = p.extractToken(p.lastProcessedTokenIndex)
			} else {
				nextStartTok = child.GetStart()
			}
			if hadMember {
				blankLines := countInputBlankLines(prevStopToken, nextStartTok)
				emptyLines := computeMemberSpacing(blankLines)
				if err := p.writeNewLines(emptyLines + 1); err != nil {
					return err
				}
			}
			if err := p.pushTokenRange(TokenIndex(startIdx - 1)); err != nil {
				return err
			}
			if err := p.printLayoutDecls(child); err != nil {
				return err
			}
			hadMember = true
			prevStopIdx := p.gatherTokenIndexFromRightiestToken(child.GetStop())
			prevStopToken = p.extractToken(prevStopIdx)
		default:
			panic("unreachable: printLayout - " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: bitfields_decls
func (p *Printer) printBitfieldsDecls(
	node dmlex.IBitfields_declsContext,
) error {
	// TODO: ensure the CdeclDefaults
	// TODO: Review AI output
	//
	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				return dmlex.GetTokenType(child.GetSymbol()) == dmlex.TokenSemi
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.ICdeclContext:
			if err := p.printCdecl(child); err != nil {
				return err
			}
		case dmlex.IBitrangeDeclContext:
			if err := p.printBitrangeDecl(child); err != nil {
				return err
			}
		default:
			panic("unreachable: BitfieldsDecls - " + fmt.Sprintf("%T", child))
		}
	}

	return p.Flush()
}

// g4-tag: bitfields
func (p *Printer) printBitfields(
	node dmlex.IBitfieldsContext,
) error {
	// TODO: ensure the CdeclDefaults
	// TODO: Review AI output

	hadMember := false
	var prevStopToken antlr.Token
	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenBitfields,
					dmlex.TokenIntegerLiteral,
					dmlex.TokenLbrace,
					dmlex.TokenRbrace:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IBitfields_declsContext:
			startIdx := child.GetStart().GetTokenIndex()
			var nextStartTok antlr.Token
			if p.lastProcessedTokenIndex < TokenIndex(startIdx) {
				nextStartTok = p.extractToken(p.lastProcessedTokenIndex)
			} else {
				nextStartTok = child.GetStart()
			}
			if hadMember {
				blankLines := countInputBlankLines(prevStopToken, nextStartTok)
				emptyLines := computeMemberSpacing(blankLines)
				if err := p.writeNewLines(emptyLines + 1); err != nil {
					return err
				}
			}
			if err := p.pushTokenRange(TokenIndex(startIdx - 1)); err != nil {
				return err
			}
			if err := p.printBitfieldsDecls(child); err != nil {
				return err
			}
			hadMember = true
			prevStopIdx := p.gatherTokenIndexFromRightiestToken(child.GetStop())
			prevStopToken = p.extractToken(prevStopIdx)
		default:
			panic("unreachable: printBrintfields - " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: struct
// g4-tag: union
// g4-tag: layout
// g4-tag: bitfields
func (p *Printer) printStructLayoutOrBitfield(
	node antlr.Tree,
) error {
	assertLazy(func() bool {
		switch node.(type) {
		case dmlex.IStructContext,
			dmlex.IUnionContext,
			dmlex.ILayoutContext,
			dmlex.IBitfieldsContext:
			return true
		}
		return false
	}, "type is not supported")
	// TODO: ensure the CdeclDefaults
	// TODO: Combine all 3 processing at one method

	switch n := node.(type) {
	case dmlex.IStructContext:
		return p.printStruct(n)
	case dmlex.IUnionContext:
		return p.printUnion(n)
	case dmlex.ILayoutContext:
		return p.printLayout(n)
	case dmlex.IBitfieldsContext:
		return p.printBitfields(n)
	default:
		panic("unreachable: structLayoutOrBitfield " + fmt.Sprintf("%T", n))
	}
}

// g4-tag: signed_or_unsigned_typeident
func (p *Printer) printSignedOrUnsignedTypeident(
	node antlr.Tree,
) error {
	assertLazy(func() bool {
		switch node.(type) {
		case dmlex.ISigned_or_unsigned_typeidentContext,
			dmlex.ITypeidentContext:
			return true
		}
		return false
	}, "type is not supported")

	// TODO: ensure the CdeclDefaults

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenUnsigned, dmlex.TokenSigned:
					return true
				case dmlex.TokenChar,
					dmlex.TokenDouble,
					dmlex.TokenFloat,
					dmlex.TokenInt:
					return true
				case dmlex.TokenLong,
					dmlex.TokenShort,
					dmlex.TokenVoid,
					dmlex.TokenRegister:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.ITypeidentContext:
			if err := p.printSignedOrUnsignedTypeident(child); err != nil {
				return err
			}
		case dmlex.IIdentContext:
			// NOTE: Simply print
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		default:
			panic("unreachable: signedOrUnsid " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: basetype
func (p *Printer) printBasetype(node dmlex.IBasetypeContext) error {
	// TODO: ensure the CdeclDefaults
	//
	var softWrapHandle func() = nil
	softWrap := func() {
		if softWrapHandle != nil {
			softWrapHandle()
			softWrapHandle = nil
		}
	}
	defer softWrap()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			killWrap := false
			switch dmlex.GetTokenType(child.GetSymbol()) {
			case dmlex.TokenSequence, dmlex.TokenHook:
				// nothing

			case dmlex.TokenLparen:
				assert(softWrapHandle == nil, "softWrap must be nil")
				softWrapHandle = p.LayoutEngine.BeginWrapGroup(
					"provisionaLRParens",
				)
			case dmlex.TokenRparen:
				killWrap = true
			default:
				panic("unreachable: pBasetype " + fmt.Sprintf("%T", child))
			}
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			if killWrap && softWrapHandle != nil {
				softWrap()
			}

		case dmlex.ISigned_or_unsigned_typeidentContext:
			if err := p.printSignedOrUnsignedTypeident(child); err != nil {
				return err
			}

		case dmlex.IStructContext,
			dmlex.IUnionContext,
			dmlex.ILayoutContext,
			dmlex.IBitfieldsContext:
			if err := p.printStructLayoutOrBitfield(child); err != nil {
				return err
			}
		case dmlex.ITypeofContext:
			if err := p.printTypeof(child); err != nil {
				return err
			}
		case dmlex.ICdecl_listContext:
			if err := p.printCdecl_list(child); err != nil {
				return err
			}
		default:
			panic("unreachable: pBasetype " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: typeQualifier
func (p *Printer) printTypeQualifier(
	node dmlex.ITypeQualifierContext,
) error {
	// TODO: ensure the CdeclDefaults
	// TODO: Review AI output

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenConst, dmlex.TokenMul, dmlex.TokenVect:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		default:
			panic("unreachable: pTypeQualifier " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: cdecl_list_variadic
func (p *Printer) printCdeclListVariadic(
	node dmlex.ICdecl_list_variadicContext,
) error {
	// TODO: ensure the CdeclDefaults
	// TODO: Review AI output

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenEllipsis, dmlex.TokenComma:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.ICdecl_listContext:
			if err := p.printCdecl_list(child); err != nil {
				return err
			}
		default:
			panic("unreachable: pCdeclListVaria " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

func (p *Printer) printCdeclArray(
	node *dmlex.CdeclArrayContext,
) error {
	// cdecl_ident_array_func LBRACK expression RBRACK # cdeclArray

	var softWrapHandle func() = nil
	softWrap := func() {
		if softWrapHandle != nil {
			softWrapHandle()
			softWrapHandle = nil
		}
	}
	defer softWrap()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			killWrap := false
			switch dmlex.GetTokenType(child.GetSymbol()) {
			case dmlex.TokenLbrack:
				assert(softWrapHandle == nil, "softWrap must be nil")
				softWrapHandle = p.LayoutEngine.BeginWrapGroup("cdeclArgs")
			case dmlex.TokenRbrack:
				killWrap = true
			default:
				panic("unreachable: pCdeclArray " + fmt.Sprintf("%T", child))
			}
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			if killWrap && softWrapHandle != nil {
				softWrap()
			}
		case dmlex.ICdecl_ident_array_funcContext:
			if err := p.printCdeclIdentArrayOrFunc(child); err != nil {
				return err
			}
		case dmlex.IExpressionContext:
			if err := p.printExpression(child); err != nil {
				return err
			}
		default:
			panic("unreachable: pCdeclArray " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

func (p *Printer) printCdeclFuncArgs(
	node *dmlex.CdeclFuncArgsContext,
) error {
	// FIXME unimplemented

	var softWrapHandle func() = nil
	softWrap := func() {
		if softWrapHandle != nil {
			softWrapHandle()
			softWrapHandle = nil
		}
	}
	defer softWrap()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			killWrap := false
			switch dmlex.GetTokenType(child.GetSymbol()) {
			case dmlex.TokenLbrack, dmlex.TokenRbrack, dmlex.TokenIdentifier:
				// nothing
			case dmlex.TokenLparen:
				assert(softWrapHandle == nil, "softWrap must be nil")
				softWrapHandle = p.LayoutEngine.BeginWrapGroup("cdeclArgs")
			case dmlex.TokenRparen:
				killWrap = true
			default:
				panic("unreachable: pCdeclFuncA " + fmt.Sprintf("%T", child))
			}
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			if killWrap && softWrapHandle != nil {
				softWrap()
			}
		case dmlex.IIdentContext:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.ICdecl_ident_array_funcContext:
			if err := p.printCdeclIdentArrayOrFunc(child); err != nil {
				return err
			}
		case dmlex.IExpressionContext:
			if err := p.printExpression(child); err != nil {
				return err
			}
		case dmlex.ICdecl_list_variadicContext:
			if err := p.printCdeclListVariadic(child); err != nil {
				return err
			}
		case dmlex.ICdecl_qualifier_ident_array_funcContext:
			if err := p.printCdeclQualifierIdentArrayOrFunc(child); err != nil {
				return err
			}
		default:
			panic("unreachable: pCdeclFuncArgs " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

func (p *Printer) printCdeclFuncIdent(
	node *dmlex.CdeclFuncIdentContext,
) error {
	// LPAREN cdecl_qualifier_ident_array_func RPAREN # cdeclFuncIdent
	lflag := p.context.GetFormatFlag(dmlex.TokenLparen).
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorAfter(format.SepNone).
		WithSeparatorPriority(format.PriorityHighest)
	rflag := p.context.GetFormatFlag(dmlex.TokenRparen).
		WithSeparatorBefore(format.SepNone).
		WithSeparatorAfter(format.SepNone).
		WithSeparatorPriority(format.PriorityHighest)
	overrides := map[dmlex.TokenType]format.FormatFlag{
		dmlex.TokenLparen: lflag,
		dmlex.TokenRparen: rflag,
	}
	p.newContext(format.NewTokenMap(overrides))
	defer p.freeContext()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenLparen, dmlex.TokenRparen:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.ICdecl_qualifier_ident_array_funcContext:
			if err := p.printCdeclQualifierIdentArrayOrFunc(child); err != nil {
				return err
			}
		default:
			panic("unreachable: pCdeclFuncIdent " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: cdecl_ident_array_func
func (p *Printer) printCdeclIdentArrayOrFunc(
	node dmlex.ICdecl_ident_array_funcContext,
) error {
	// TODO: ensure the CdeclDefaults

	switch child := node.(type) {
	case *dmlex.CdeclIdentContext:
		if err := p.pushTokenRangeFromRightiest(child); err != nil {
			return err
		}
	case *dmlex.CdeclNoneContext:
		// nothing
	case *dmlex.CdeclFuncIdentContext:
		if err := p.printCdeclFuncIdent(child); err != nil {
			return err
		}
	case *dmlex.CdeclFuncArgsContext:
		if err := p.printCdeclFuncArgs(child); err != nil {
			return err
		}
	case *dmlex.CdeclArrayContext:
		if err := p.printCdeclArray(child); err != nil {
			return err
		}
	default:
		panic("unreachable: pCdeclIdentArrayOrFun " + fmt.Sprintf("%T", child))
	}

	return nil
}

// g4-tag: cdecl_qualifier_ident_array_func
func (p *Printer) printCdeclQualifierIdentArrayOrFunc(
	node dmlex.ICdecl_qualifier_ident_array_funcContext,
) error {
	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.ITypeQualifierContext:
			if err := p.printTypeQualifier(child); err != nil {
				return err
			}
		case dmlex.ICdecl_ident_array_funcContext:
			if err := p.printCdeclIdentArrayOrFunc(child); err != nil {
				return err
			}
		default:
			panic("unreachable: pCdeclIdentArrayO " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: cdecl
func (p *Printer) printCdecl(node dmlex.ICdeclContext) error {
	p.newContext(p.settings.CdeclDefaults())
	defer p.freeContext()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

		case dmlex.IBasetypeContext:
			if err := p.printBasetype(child); err != nil {
				return err
			}

		case dmlex.ICdecl_qualifier_ident_array_funcContext:
			if err := p.printCdeclQualifierIdentArrayOrFunc(child); err != nil {
				return err
			}
		default:
			panic("unreachable: printCdecl " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: cdeclOrInline
func (p *Printer) printCdeclOrInline(node dmlex.ICdeclOrInlineContext) error {
	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case dmlex.ICdeclContext:
			if err := p.printCdecl(child); err != nil {
				return err
			}
		case antlr.TerminalNode, dmlex.IIdentContext:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		default:
			panic("unreachable: pCdeclOrInline " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: typedefStatement
func (p *Printer) printTypedefStatement(
	node dmlex.ITypedefStatementContext,
) error {
	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenTypedef, dmlex.TokenSemi:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.ICdeclContext:
			if err := p.printCdecl(child); err != nil {
				return err
			}
		default:
			panic("unreachable: pTypedefStatement " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: cdecl_list
func (p *Printer) printCdecl_list(node dmlex.ICdecl_listContext) error {
	defer p.LayoutEngine.BeginWrapGroup("cdecl_list")()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				return dmlex.GetTokenType(child.GetSymbol()) == dmlex.TokenComma
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.ICdeclContext:
			if err := p.printCdecl(child); err != nil {
				return err
			}
		default:
			panic("unreachable: pCdecl_list " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: ctypedecl_simple
func (p *Printer) printCtypedecl_simple(
	node dmlex.ICtypedecl_simpleContext,
) error {
	var softWrapHandle func() = nil
	softWrap := func() {
		if softWrapHandle != nil {
			softWrapHandle()
			softWrapHandle = nil
		}
	}
	defer softWrap()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			killWrap := false
			switch dmlex.GetTokenType(child.GetSymbol()) {
			case dmlex.TokenLbrack, dmlex.TokenRbrack:
				// nothing
			case dmlex.TokenLparen:
				assert(softWrapHandle == nil, "softWrap must be nil")
				softWrapHandle = p.LayoutEngine.BeginWrapGroup(
					"ctypedecl_simple",
				)
			case dmlex.TokenRparen:
				killWrap = true
			default:
				panic("unreachable: pCtypedecl_s " + fmt.Sprintf("%T", child))
			}
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			if killWrap && softWrapHandle != nil {
				softWrap()
			}
		case dmlex.ICtypedecl_ptrContext:
			if err := p.printCtype_ptr(child); err != nil {
				return err
			}
		default:
			panic("unreachable: pCtypedecl_s " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: ctypedecl_ptr
func (p *Printer) printCtype_ptr(node dmlex.ICtypedecl_ptrContext) error {
	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenMul, dmlex.TokenConst:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.ICtypedecl_simpleContext:
			if err := p.printCtypedecl_simple(child); err != nil {
				return err
			}
		default:
			panic("unreachable: pCdecl_ptr " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: ctypedecl
func (p *Printer) printCtypedecl(node dmlex.ICtypedeclContext) error {
	p.newContext(p.settings.CdeclDefaults())
	defer p.freeContext()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenConst:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IBasetypeContext:
			if err := p.printBasetype(child); err != nil {
				return err
			}
		case dmlex.ICtypedecl_ptrContext:
			if err := p.printCtype_ptr(child); err != nil {
				return err
			}
		default:
			panic("unreachable: pCtypedecl " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}
