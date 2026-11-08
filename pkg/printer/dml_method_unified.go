// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer

import (
	"fmt"
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer/format"

	"github.com/antlr4-go/antlr/v4"
)

// g4-tag: methodUnified
func (p *Printer) printMethodUnified(
	node dmlex.IMethodUnifiedContext,
) error {
	p.newContext(p.settings.CdeclDefaults())
	defer p.freeContext()

	switch realNode := node.(type) {
	case *dmlex.MethodUnifiedDeclarationContext:
		return p.printMethodUnifiedDeclaration(realNode)
	case *dmlex.MethodUnifiedDefinitionContext:
		return p.printMethodUnifiedDefinition(realNode)
	default:
		panic("unreachable: " + fmt.Sprintf("%T", realNode))
	}

}

func (p *Printer) methodArgsFormatMap() format.TokenMapImpl {
	m := make(format.TokenMapImpl)
	m[dmlex.TokenLparen] = p.context.GetFormatFlag(dmlex.TokenLparen).
		WithBreakAfterPriority(format.PriorityVeryHigh)
	m[dmlex.TokenRparen] = p.context.GetFormatFlag(dmlex.TokenRparen).
		WithBreakBeforePriority(format.PriorityVeryHigh)
	m[dmlex.TokenComma] = p.context.GetFormatFlag(dmlex.TokenComma).
		WithBreakAfterPriority(format.PriorityAboveNormal)
	return m
}

// g4-tag: methodQualifiers
func (p *Printer) printMethodQualifiers(
	node dmlex.IMethodQualifiersContext,
) error {
	endQualifiers := p.LayoutEngine.BeginWrapGroup(
		"methodQualifiers",
	)
	defer endQualifiers()

	tokenCount := node.GetChildCount()
	if tokenCount <= 0 {
		return nil
	}
	lastToken := node.GetChild(tokenCount - 1)
	if lastToken == nil {
		return nil
	}
	stopTok := castStopToken(lastToken)
	if stopTok == nil {
		return nil
	}
	tokenType := dmlex.GetTokenType(stopTok)
	// NOTE: Make sure that only last token has optional wrap
	//       to avoid wrapping on all of them
	// NOTE: Can be done via break priority?
	customMap := make(format.TokenMapImpl)
	customMap[tokenType] = p.context.GetFormatFlag(tokenType).
		WithSeparatorAfter(format.SepSpaceMaybeWrap).
		WithBreakAfterPriority(format.PriorityLow).
		WithSeparatorPriority(format.PriorityAboveNormal)

	p.newContext(customMap)
	defer p.freeContext()

	return p.pushTokenRangeFromRightiest(node)
}

// g4-tag: methodUnifiedDeclaration
func (p *Printer) printMethodUnifiedDeclaration(
	node *dmlex.MethodUnifiedDeclarationContext,
) error {
	var endMethodArgumentsHandle func() = nil
	endMethodArguments := func() {
		if endMethodArgumentsHandle != nil {
			endMethodArgumentsHandle()
			endMethodArgumentsHandle = nil
		}
	}
	var endSignatureHandle func() = nil
	endSignature := func() {
		if endSignatureHandle != nil {
			endSignatureHandle()
			endSignatureHandle = nil
		}
	}
	inArgsContext := false
	defer func() {
		endMethodArguments()
		endSignature()
		if inArgsContext {
			p.freeContext()
			inArgsContext = false
		}
	}()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			killArgsWrap := false
			killSignatureWrap := false
			switch dmlex.GetTokenType(child.GetSymbol()) {
			case dmlex.TokenMethod, dmlex.TokenThrows:
				// nothing
			case dmlex.TokenSemi:
				killSignatureWrap = true
			case dmlex.TokenLparen:
				assert(
					endMethodArgumentsHandle == nil,
					"methodArguments wrap group already active",
				)
				endSignatureHandle = p.LayoutEngine.BeginWrapGroup(
					"methodSignature",
				)
				endMethodArgumentsHandle = p.LayoutEngine.BeginWrapGroup(
					"methodArguments",
				)
				inArgsContext = true
				p.newContext(p.methodArgsFormatMap())
			case dmlex.TokenRparen:
				killArgsWrap = true
			default:
				panic("unreachable: " + fmt.Sprintf("%T", child))
			}
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			if killArgsWrap && endMethodArgumentsHandle != nil {
				endMethodArguments()
				if inArgsContext {
					p.freeContext()
					inArgsContext = false
				}
			}
			if killSignatureWrap && endSignatureHandle != nil {
				endSignature()
			}

		case dmlex.IObjidentContext:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IMethodQualifiersContext:
			if err := p.printMethodQualifiers(child); err != nil {
				return err
			}

		case dmlex.IMethodParamsMaybeUntypedContext:
			if err := p.printMethodParamsMaybeUntyped(child); err != nil {
				return err
			}

		case dmlex.IMethodOutParamsContext:
			if err := p.printMethodOutParams(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: methodUnifiedDefinition
func (p *Printer) printMethodUnifiedDefinition(
	node *dmlex.MethodUnifiedDefinitionContext,
) error {
	var endMethodArgumentsHandle func() = nil
	endMethodArguments := func() {
		if endMethodArgumentsHandle != nil {
			endMethodArgumentsHandle()
			endMethodArgumentsHandle = nil
		}
	}
	var endSignatureHandle func() = nil
	endSignature := func() {
		if endSignatureHandle != nil {
			endSignatureHandle()
			endSignatureHandle = nil
		}
	}
	inArgsContext := false
	defer func() {
		endMethodArguments()
		endSignature()
		if inArgsContext {
			p.freeContext()
			inArgsContext = false
		}
	}()

	hadStatement := false
	var prevStopToken antlr.Token
	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case dmlex.IStatementContext:
			startIdx := child.GetStart().GetTokenIndex()
			var nextStartTok antlr.Token
			if p.lastProcessedTokenIndex < TokenIndex(startIdx) {
				nextStartTok = p.extractToken(p.lastProcessedTokenIndex)
			} else {
				nextStartTok = child.GetStart()
			}
			if hadStatement {
				blankLines := countInputBlankLines(prevStopToken, nextStartTok)
				emptyLines := computeMethodBodySpacing(blankLines)
				if err := p.writeNewLines(emptyLines + 1); err != nil {
					return err
				}
			}
			if err := p.pushTokenRange(TokenIndex(startIdx - 1)); err != nil {
				return err
			}
			if err := p.printStatement(child); err != nil {
				return err
			}
			hadStatement = true
			prevStopIdx := p.gatherTokenIndexFromRightiestToken(child.GetStop())
			prevStopToken = p.extractToken(prevStopIdx)

		case antlr.TerminalNode:
			killArgsWrap := false
			killSignatureWrap := false
			switch dmlex.GetTokenType(child.GetSymbol()) {
			case dmlex.TokenMethod, dmlex.TokenThrows, dmlex.TokenDefault,
				dmlex.TokenDefaultNew,
				dmlex.TokenRbrace: // Rbrace ({) has TokenFormat settings to handle
				// nothing
			case dmlex.TokenLparen:
				assert(
					endMethodArgumentsHandle == nil,
					"methodArguments wrap group already active",
				)
				endSignatureHandle = p.LayoutEngine.BeginWrapGroup(
					"methodSignature",
				)
				endMethodArgumentsHandle = p.LayoutEngine.BeginWrapGroup(
					"methodArguments",
				)
				inArgsContext = true
				p.newContext(p.methodArgsFormatMap())
			case dmlex.TokenRparen:
				killArgsWrap = true
			case dmlex.TokenLbrace, dmlex.TokenLbraceNew:
				killSignatureWrap = true
			default:
				panic("unreachable: " + fmt.Sprintf("%T", child))
			}
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			if killArgsWrap && endMethodArgumentsHandle != nil {
				endMethodArguments()
				if inArgsContext {
					p.freeContext()
					inArgsContext = false
				}
			}
			if killSignatureWrap && endSignatureHandle != nil {
				endSignature()
			}

		case dmlex.IObjidentContext:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IMethodQualifiersContext:
			if err := p.printMethodQualifiers(child); err != nil {
				return err
			}

		case dmlex.IMethodParamsMaybeUntypedContext:
			if err := p.printMethodParamsMaybeUntyped(child); err != nil {
				return err
			}

		case dmlex.IMethodOutParamsContext:
			if err := p.printMethodOutParams(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: methodParamsMaybeUntyped
func (p *Printer) printMethodParamsMaybeUntyped(
	node dmlex.IMethodParamsMaybeUntypedContext,
) error {
	endMethodArguments := p.LayoutEngine.BeginWrapGroup(
		"methodParamsMaybeUntyped",
	)
	defer endMethodArguments()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			switch dmlex.GetTokenType(child.GetSymbol()) {
			case dmlex.TokenComma:
				// nothing
			default:
				panic("unreachable: " + fmt.Sprintf("%T", child))
			}
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.ICdeclOrInlineContext:
			if err := p.printCdeclOrInline(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: methodOutParams
func (p *Printer) printMethodOutParams(
	node dmlex.IMethodOutParamsContext,
) error {
	methodOutMap := make(format.TokenMapImpl)
	methodOutMap[dmlex.TokenLparen] = p.context.GetFormatFlag(dmlex.TokenLparen).
		WithBreakAfterPriority(format.PriorityHigh)
	methodOutMap[dmlex.TokenRparen] = p.context.GetFormatFlag(dmlex.TokenRparen).
		WithBreakBeforePriority(format.PriorityHigh)
	methodOutMap[dmlex.TokenComma] = p.context.GetFormatFlag(dmlex.TokenComma).
		WithBreakAfterPriority(format.PriorityNormal)

	p.newContext(methodOutMap)
	defer p.freeContext()

	var endOutParamsHandle func() = nil
	endOutParams := func() {
		if endOutParamsHandle != nil {
			endOutParamsHandle()
			endOutParamsHandle = nil
		}
	}
	defer endOutParams()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			killWrap := false
			switch dmlex.GetTokenType(child.GetSymbol()) {
			case dmlex.TokenArrow:
				// nothing
			case dmlex.TokenLparen:
				assert(
					endOutParamsHandle == nil,
					"methodOutParamsParen wrap group already active",
				)
				endOutParamsHandle = p.LayoutEngine.BeginWrapGroup(
					"methodOutParamsParen",
				)
			case dmlex.TokenRparen:
				killWrap = true
			default:
				panic("unreachable: " + fmt.Sprintf("%T", child))
			}
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			if killWrap && endOutParamsHandle != nil {
				endOutParams()
			}
		case dmlex.ICdecl_listContext:
			if err := p.printCdecl_list(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}
