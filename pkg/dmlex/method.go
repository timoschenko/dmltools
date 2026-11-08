// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package dmlex

import "fmt"

const InlineType = "<inline>"

func MethodName(m IMethodUnifiedContext) string {
	switch ctx := m.(type) {
	case *MethodUnifiedDeclarationContext:
		if ctx.Objident() != nil {
			return ctx.Objident().GetText()
		}
	case *MethodUnifiedDefinitionContext:
		if ctx.Objident() != nil {
			return ctx.Objident().GetText()
		}
	}
	panic("unreachable: MethodName " + fmt.Sprintf("%T", m))
}

func getMethodQualifiers(m IMethodUnifiedContext) IMethodQualifiersContext {
	switch ctx := m.(type) {
	case *MethodUnifiedDeclarationContext:
		return ctx.MethodQualifiers()
	case *MethodUnifiedDefinitionContext:
		return ctx.MethodQualifiers()
	default:
		return nil
	}
}

func MethodIsInline(m IMethodUnifiedContext) bool {
	qualifiers := getMethodQualifiers(m)
	if qualifiers == nil {
		return false
	}

	return qualifiers.INLINE() != nil
}

func MethodIsShared(m IMethodUnifiedContext) bool {
	qualifiers := getMethodQualifiers(m)
	if qualifiers == nil {
		return false
	}

	return qualifiers.SHARED() != nil
}

func MethodIsIndependent(m IMethodUnifiedContext) bool {
	qualifiers := getMethodQualifiers(m)
	if qualifiers == nil {
		return false
	}

	return qualifiers.INDEPENDENT() != nil
}

func MethodIsStartup(m IMethodUnifiedContext) bool {
	qualifiers := getMethodQualifiers(m)
	if qualifiers == nil {
		return false
	}

	return qualifiers.STARTUP() != nil
}

func MethodIsMemoized(m IMethodUnifiedContext) bool {
	qualifiers := getMethodQualifiers(m)
	if qualifiers == nil {
		return false
	}

	return qualifiers.MEMOIZED() != nil
}

func MethodThrows(m IMethodUnifiedContext) bool {
	switch ctx := m.(type) {
	case *MethodUnifiedDeclarationContext:
		return ctx.THROWS() != nil
	case *MethodUnifiedDefinitionContext:
		return ctx.THROWS() != nil
	default:
		return false
	}
}

func MethodDefault(m IMethodUnifiedContext) bool {
	switch ctx := m.(type) {
	case *MethodUnifiedDefinitionContext:
		return ctx.DEFAULT() != nil
	default:
		return false
	}
}

func MethodBody(m IMethodUnifiedContext) []IStatementContext {
	switch ctx := m.(type) {
	case *MethodUnifiedDefinitionContext:
		return ctx.AllStatement()
	default:
		return nil
	}
}

func MethodReturns(m IMethodUnifiedContext) []ICdeclContext {
	var outputs IMethodOutParamsContext
	switch ctx := m.(type) {
	case *MethodUnifiedDeclarationContext:
		outputs = ctx.MethodOutParams()
	case *MethodUnifiedDefinitionContext:
		outputs = ctx.MethodOutParams()
	}
	if outputs == nil {
		return nil
	}
	cdecl_list := outputs.Cdecl_list()
	if cdecl_list == nil {
		return nil
	}
	return cdecl_list.AllCdecl()
}

// TODO: Move to a dedicated file
func IdentName(node IObjidentContext) string {
	switch upcast := node.(type) {
	case IObjidentContext:
	case IIdentContext:
	case ITypeidentContext:
		_ = upcast /* nothing */
	default:
		panic("IdentName does not support selected node type")
	}
	return node.GetText()
}
