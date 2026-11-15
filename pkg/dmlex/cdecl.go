// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package dmlex

import "strings"

const TypeInline = "<inline>"

// Extracts base type
// Examples:
//
//	int a           -> int
//	const int a     -> int
//	int *a          -> int
//	int ** const a  -> int
func CdeclBasetypeName(t ICdeclContext) string {
	if t == nil || t.Basetype() == nil {
		return ""
	}
	basetype := t.Basetype()
	if su := basetype.Signed_or_unsigned_typeident(); su != nil {
		sign := ""
		if su.UNSIGNED() != nil {
			sign = "unsigned"
		} else if su.SIGNED() != nil {
			sign = "signed"
		}
		if su.Typeident() != nil {
			typeText := su.Typeident().GetText()
			if sign != "" {
				return sign + " " + typeText
			}
			return typeText
		}
		return sign
	}
	return basetype.GetText()
}

// Extracts full type
// Examples:
//
//	int a           -> int
//	const int a     -> const int
//	int *a          -> int*
//	int ** const a  -> int**const
func CdeclTypeName(c ICdeclContext) string {
	if c == nil || c.Basetype() == nil {
		return ""
	}
	var typeName strings.Builder
	if typeConst := c.CONST(); typeConst != nil {
		typeName.WriteString("const ")
	}
	typeName.WriteString(CdeclBasetypeName(c))
	cdeclQIAF := c.Cdecl_qualifier_ident_array_func()
	if cdeclQIAF != nil {
		for _, tq := range cdeclQIAF.AllTypeQualifier() {
			if found := tq.CONST(); found != nil {
				typeName.WriteString("const ")
			} else if found := tq.MUL(); found != nil {
				typeName.WriteString("*")
			} else if found := tq.VECT(); found != nil {
				typeName.WriteString("vect ")
			} else {
				panic("AllTypeQualifier")
			}
		}
	}

	return strings.TrimRight(typeName.String(), " ")
}

// Extracts full type for inline method
// Examples:
//
//	int a           -> int
//	const int a     -> const int
//	int *a          -> int*
//	int ** const a  -> int**const
//	inline a        -> <inline>
func CdeclTypeNameInline(c ICdeclOrInlineContext) string {
	if c == nil {
		return ""
	}
	if c.INLINE() != nil {
		return TypeInline
	}
	cdecl := c.Cdecl()
	return CdeclTypeName(cdecl)
}
