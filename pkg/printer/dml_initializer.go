// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer

import (
	"fmt"
	"github.com/timoschenko/dmltools/pkg/dmlex"

	"github.com/antlr4-go/antlr/v4"
)

// g4-tag: singleInitializerArray
func (p *Printer) printSingleInitializerArray(
	node dmlex.ISingleInitializerArrayContext,
) error {
	singleInitializerElementGroup := p.LayoutEngine.BeginWrapGroup(
		"singleInitializerArray",
	)
	defer singleInitializerElementGroup()

	for _, childRaw := range node.GetChildren() {
		switch child := childRaw.(type) {
		case dmlex.ISingleInitializerContext:
			if err := p.printSingleInitializer(child); err != nil {
				return nil
			}
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenComma:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: designatedStructInitializer
func (p *Printer) printDesignatedStructInitializer(
	node dmlex.IDesignatedStructInitializerContext,
) error {
	for _, childRaw := range node.GetChildren() {
		switch child := childRaw.(type) {
		case dmlex.IIdentContext:
			if err := p.printIdent(child); err != nil {
				return err
			}
		case dmlex.ISingleInitializerContext:
			if err := p.printSingleInitializer(child); err != nil {
				return err
			}
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenDot, dmlex.TokenAssign:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: designated_struct_initializer_list
func (p *Printer) printDesignatedStructInitializerList(
	node dmlex.IDesignated_struct_initializer_listContext,
) error {
	p.newContext(p.settings.ExpressionDefaults())
	defer p.freeContext()

	singleInitializerElementGroup := p.LayoutEngine.BeginWrapGroup(
		"designated_struct_initializer_list",
	)
	defer singleInitializerElementGroup()

	for _, childRaw := range node.GetChildren() {
		switch child := childRaw.(type) {
		case dmlex.IDesignatedStructInitializerContext:
			if err := p.printDesignatedStructInitializer(child); err != nil {
				return nil
			}
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenComma, dmlex.TokenEllipsis:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: cdeclInitializer
func (p *Printer) printCdeclInitializer(
	node dmlex.ICdeclInitializerContext,
) error {
	// TODO: add assert about context is expressions default
	p.newContext(p.settings.ExpressionDefaults())
	defer p.freeContext()

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
			case dmlex.TokenLbrace:
				assert(softWrapHandle == nil, "softWrap must be nil")
				softWrapHandle = p.LayoutEngine.BeginWrapGroup(
					"cdeclInitializer",
				)
			case dmlex.TokenRbrace:
				killWrap = true
			default:
				panic("unreachable: " + fmt.Sprintf("%T", child))
			}
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			if killWrap && softWrapHandle != nil {
				softWrap()
			}
		case dmlex.ISingleInitializerArrayContext:
			if err := p.printSingleInitializerArray(child); err != nil {
				return err
			}
		case dmlex.IDesignated_struct_initializer_listContext:
			if err := p.printDesignatedStructInitializerList(
				child,
			); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: singleInitializer
//
//go:inline
func (p *Printer) printSingleInitializer(
	node dmlex.ISingleInitializerContext,
) error {
	// Can be only one child item at this level (excl Comments)
	if node.GetChildCount() != 1 {
		panic("unreachable: " + fmt.Sprintf("%T", node))
	}

	childNode := node.GetChild(0)
	switch child := childNode.(type) {
	case dmlex.IExpressionContext:
		return p.printExpression(child)
	case dmlex.ICdeclInitializerContext:
		return p.printCdeclInitializer(child)
	default:
		panic("unreachable: " + fmt.Sprintf("%T", child))
	}

}

// g4-tag: argumentsList
func (p *Printer) printArgumentsList(node dmlex.IArgumentsListContext) error {
	p.newContext(p.settings.ExpressionDefaults())
	defer p.freeContext()

	endGroup := p.LayoutEngine.BeginWrapGroup("ArgumentsList")
	defer endGroup()

	for _, childRaw := range node.GetChildren() {
		switch child := childRaw.(type) {
		case dmlex.ISingleInitializerArrayContext:
			if err := p.printSingleInitializerArray(child); err != nil {
				return err
			}
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
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: initializer
//
//go:inline
func (p *Printer) printInitializer(
	node dmlex.IInitializerContext,
) error {
	// Can be only one child item at this level (excl Comments)
	if node.GetChildCount() != 1 {
		panic("unreachable: " + fmt.Sprintf("%T", node))
	}

	childNode := node.GetChild(0)
	switch child := childNode.(type) {
	case dmlex.ISingleInitializerContext:
		return p.printSingleInitializer(child)
	case dmlex.IArgumentsListContext:
		return p.printArgumentsList(child)
	default:
		panic("unreachable: " + fmt.Sprintf("%T", child))
	}
}
