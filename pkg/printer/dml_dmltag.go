// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer

import (
	"fmt"
	"github.com/timoschenko/dmltools/pkg/dmlex"

	"github.com/antlr4-go/antlr/v4"
)

// g4-tag: dml
func (p *Printer) printDmlContext(node dmlex.IDmlContext) error {
	// format: DmlTag [SepMandatoryWrap] DmlCode [SepMandatoryWrap] EOF [SepMandatoryWrap]

	tree := node.GetChildren()
	hadStatement := false
	for _, child := range tree {
		var startIdx int
		var handler func() error
		newStatement := false
		pushGap := true
		switch child := child.(type) {
		case dmlex.IDmlTagContext:
			newStatement = true
			startIdx = child.GetStart().GetTokenIndex()
			handler = func() error { return p.printDmlTag(child) }
		case dmlex.IDmlCodeContext:
			// Don't push gap here — printCodeContext will handle
			// the gap between dmlTag and its first child, ensuring
			// the empty lines come before any comments.
			pushGap = false
			handler = func() error { return p.printCodeContext(child) }
		case antlr.TerminalNode:
			// Expect here only EOF
			startIdx = child.GetSymbol().GetTokenIndex()
			handler = func() error { return p.printEof(child) }
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}

		// Insert empty lines BEFORE the gap comments and the statement.
		if newStatement && hadStatement {
			if err := p.writeNewLines(
				CFG_NEW_LINES_ON_TOP_LEVEL + 1,
			); err != nil {
				return err
			}
		}

		// Push everything between the last processed token and this rule's start.
		if pushGap {
			if err := p.pushTokenRange(TokenIndex(startIdx - 1)); err != nil {
				return err
			}
		}

		if err := handler(); err != nil {
			return err
		}

		if newStatement {
			hadStatement = true
			// Ensure the statement ends with exactly one newline.
			if err := p.writeNewLines(1); err != nil {
				return err
			}
		}
	}
	return nil
}

// g4-tag: dmlCode
func (p *Printer) printCodeContext(node dmlex.IDmlCodeContext) error {
	// format: ((DeviceTag | BitorderTag | ProvisionalTag) [SepMandatoryWrap])* Statements [SepMandatoryWrap]

	tree := node.GetChildren()
	hadStatement := p.LayoutEngine.GetLastChar_Workaround() != 0
	var prevStopToken antlr.Token
	if p.lastProcessedTokenIndex > 0 {
		prevStopToken = p.extractToken(p.lastProcessedTokenIndex - 1)
	}

	for _, child := range tree {
		var startIdx int
		var handler func() error
		var stopNode antlr.ParserRuleContext
		newStatement := false
		pushGap := true
		switch child := child.(type) {
		case dmlex.IDeviceTagContext:
			newStatement = true
			startIdx = child.GetStart().GetTokenIndex()
			stopNode = child
			handler = func() error { return p.printDeviceTag(child) }
		case dmlex.IBitorderTagContext:
			newStatement = true
			startIdx = child.GetStart().GetTokenIndex()
			stopNode = child
			handler = func() error { return p.printBitorderTag(child) }
		case dmlex.IProvisionalTagContext:
			newStatement = true
			startIdx = child.GetStart().GetTokenIndex()
			stopNode = child
			handler = func() error { return p.printProvisionalTag(child) }
		case dmlex.IDeviceStatementsContext:
			pushGap = false
			startIdx = child.GetStart().GetTokenIndex()
			handler = func() error { return p.printDeviceStatements(child) }
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}

		// Insert empty lines BEFORE the gap comments and the statement.
		if newStatement && hadStatement {
			var nextStartTok antlr.Token
			if p.lastProcessedTokenIndex < TokenIndex(startIdx) {
				nextStartTok = p.extractToken(p.lastProcessedTokenIndex)
			} else {
				nextStartTok = child.(antlr.ParserRuleContext).GetStart()
			}
			blankLines := countInputBlankLines(prevStopToken, nextStartTok)
			emptyLines := min(blankLines, 2)
			if err := p.writeNewLines(emptyLines + 1); err != nil {
				return err
			}
		}

		// Push everything between the last processed token and this rule's start.
		if pushGap {
			if err := p.pushTokenRange(TokenIndex(startIdx - 1)); err != nil {
				return err
			}
		}

		if err := handler(); err != nil {
			return err
		}

		if newStatement {
			hadStatement = true
			if stopNode != nil {
				prevStopIdx := p.gatherTokenIndexFromRightiestToken(
					stopNode.GetStop(),
				)
				prevStopToken = p.extractToken(prevStopIdx)
			}
			// Ensure the statement ends with exactly one newline.
			if err := p.writeNewLines(1); err != nil {
				return err
			}
		}
	}
	return nil
}

// g4-tag: <EOF>
func (p *Printer) printEof(node antlr.TerminalNode) error {
	assertLazy(func() bool {
		text := node.GetText()
		return text == "<EOF>"
	}, "unknown terminal node")

	tokenIndexBegin := p.lastProcessedTokenIndex
	tokenIndexEnd := p.gatherTokenIndexFromRightiestTerminalNode(node)

	if tokenIndexBegin >= tokenIndexEnd {
		// EOF is the last token, so we don't need to print it.
		return nil
	}

	if err := p.pushTokenRange(tokenIndexEnd); err != nil {
		return err
	}
	return p.Flush()
}

// g4-tag: dmlTag
func (p *Printer) printDmlTag(node dmlex.IDmlTagContext) error {
	// format: DML [SepSpace] FLOAT_LITERAL [SepNone] SEMI [SepMandatoryWrap]
	return p.printSimpleRuleContextEndsSemi(node)
}

// g4-tag: provisionalTag
func (p *Printer) printProvisionalTag(node dmlex.IProvisionalTagContext) error {
	// format: PROVISIONAL [SepSpaceMaybeWrap] (LPAREN [SepSpaceMaybeWrap])     \
	//         Ident [SepNone] (COMMA [SepSpaceMaybeWrap] Ident [SepNone])+   \
	//         ([SepSpaceMaybeWrap] RPAREN) [SepNone] SEMI [SepMandatoryWrap]
	identList := node.IdentList()
	if identList == nil {
		return nil
	}
	idents := identList.AllIdent()
	if len(idents) == 0 {
		return nil
	}

	var softWrapHandle func() = nil
	softWrap := func() {
		if softWrapHandle != nil {
			softWrapHandle()
			softWrapHandle = nil
		}
	}
	var innerSoftWrapHandle func() = nil
	innerSoftWrap := func() {
		if innerSoftWrapHandle != nil {
			innerSoftWrapHandle()
			innerSoftWrapHandle = nil
		}
	}
	defer func() {
		innerSoftWrap()
		softWrap()
	}()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			killWrap, softIdent := false, false
			switch dmlex.GetTokenType(child.GetSymbol()) {
			case dmlex.TokenProvisional:
				softIdent = true
			case dmlex.TokenSemi:
				// nothing

			case dmlex.TokenLparen:
				assert(softWrapHandle == nil, "softWrap must be nil")
				softWrapHandle = p.LayoutEngine.BeginWrapGroup(
					"provisionaLRParens",
				)
			case dmlex.TokenRparen:
				killWrap = true
			default:
				panic("unreachable: " + fmt.Sprintf("%T", child))
			}
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			if killWrap && softWrapHandle != nil {
				softWrap()
			} else if softIdent {
				defer p.LayoutEngine.BeginSoftIndent()()
			}
		case dmlex.IIdentListContext:
			assert(innerSoftWrapHandle == nil, "innerSoftWrap must be nil")
			innerSoftWrapHandle = p.LayoutEngine.BeginWrapGroup(
				"provisionaInnerArgs",
			)
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			innerSoftWrap()
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return p.Flush()
}

// g4-tag: deviceTag
func (p *Printer) printDeviceTag(node dmlex.IDeviceTagContext) error {
	// format: DEVICE [SepSpace] Objident [SepNone] SEMI [SepMandatoryWrap]
	return p.printSimpleRuleContextEndsSemi(node)
}

// g4-tag: bitorderTag
func (p *Printer) printBitorderTag(node dmlex.IBitorderTagContext) error {
	// format: BITORDER [SepSpace] Bitorder_endian [SepNone] SEMI [SepMandatoryWrap]
	return p.printSimpleRuleContextEndsSemi(node)
}

// g4-tag: deviceStatements
func (p *Printer) printDeviceStatements(
	node dmlex.IDeviceStatementsContext,
) error {
	tree := node.AllDeviceStatement()
	if len(tree) == 0 {
		return nil
	}

	hadStatement := p.LayoutEngine.GetLastChar_Workaround() != 0
	prevKind := StmtKindHeader
	var prevStopToken antlr.Token
	if p.lastProcessedTokenIndex > 0 {
		prevStopToken = p.extractToken(p.lastProcessedTokenIndex - 1)
	}

	for _, child := range tree {
		currKind := categorizeDeviceStatement(child)
		startIdx := child.GetStart().GetTokenIndex()

		var nextStartTok antlr.Token
		if p.lastProcessedTokenIndex < TokenIndex(startIdx) {
			nextStartTok = p.extractToken(p.lastProcessedTokenIndex)
		} else {
			nextStartTok = child.GetStart()
		}
		blankLines := countInputBlankLines(prevStopToken, nextStartTok)

		if hadStatement {
			emptyLines := computeTopLevelSpacing(
				prevKind,
				currKind,
				blankLines,
			)
			if err := p.writeNewLines(emptyLines + 1); err != nil {
				return err
			}
		}

		if err := p.pushTokenRange(TokenIndex(startIdx - 1)); err != nil {
			return err
		}

		if err := p.printDeviceStatement(child); err != nil {
			return err
		}

		hadStatement = true
		prevKind = currKind
		prevStopIdx := p.gatherTokenIndexFromRightiestToken(child.GetStop())
		prevStopToken = p.extractToken(prevStopIdx)

		if err := p.writeNewLines(1); err != nil {
			return err
		}
	}

	return nil
}
