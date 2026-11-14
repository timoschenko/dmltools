// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package dmlex

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

// ExprKind clearly distinguishes the high-level kind of an expression.
type ExprKind int

const (
	ExprKindLiteral ExprKind = iota
	ExprKindBinaryOp
	ExprKindUnaryOp
	ExprKindTernary
	ExprKindFunctionCall
	ExprKindArrayAccess
	ExprKindBitSlice
	ExprKindMemberAccess
	ExprKindCast
	ExprKindSizeof
	ExprKindDefined
	ExprKindStringify
	ExprKindNew
	ExprKindEach
	ExprKindSizeoftype
	ExprKindArrayLiteral
	ExprKindParen
)

// ISimplifiedExpr is the clean interface all wrapped expressions will implement
//
// Usage:
//
//	rawExpr := ctx.Expression() // some IExpressionContext from ANTLR
//	simplified := dmlex.WrapExpression(rawExpr)
//
//	switch e := simplified.(type) {
//	case *dmlex.BinaryExpr:
//	    fmt.Printf("Left: %v, Op: %s, Right: %v", e.Left, e.Op, e.Right)
//	case *dmlex.LiteralExpr:
//	    fmt.Printf("Literal: %s", e.Value)
//	}
type ISimplifiedExpr interface {
	Kind() ExprKind
	RawNode() IExpressionContext // Keep a reference to the raw ANTLR node
}

type basenodeExpr struct {
	node IExpressionContext
}

func (b *basenodeExpr) RawNode() IExpressionContext { return b.node }

// LiteralExpr represents atomic literal values like identifiers, strings, or numbers.
//
// Examples:
//   - "hello"
//   - 123
//   - myVariable
//   - true
//   - undefined
type LiteralExpr struct {
	basenodeExpr
	Value antlr.Token
}

func (l *LiteralExpr) Kind() ExprKind { return ExprKindLiteral }

// BinaryExpr represents any operation involving two expressions and an operator.
//
// Examples:
//   - a + b
//   - a && b
//   - a == b
//   - a << 2
type BinaryExpr struct {
	basenodeExpr
	Left  ISimplifiedExpr
	Right ISimplifiedExpr
	Op    antlr.Token
}

func (b *BinaryExpr) Kind() ExprKind { return ExprKindBinaryOp }

// UnaryExpr represents a unary operation applied to a single expression.
//
// Examples:
//   - !a
//   - ~b
//   - &c
//   - *d
//   - a++ (where IsPost=true)
//   - ++a (where IsPost=false)
type UnaryExpr struct {
	basenodeExpr
	Expr   ISimplifiedExpr
	Op     antlr.Token
	IsPost bool
}

func (u *UnaryExpr) Kind() ExprKind { return ExprKindUnaryOp }

// TernaryExpr represents an inline if-else conditional expression.
//
// Examples:
//   - a ? b : c
//   - a #? b #: c (where IsConst=true)
type TernaryExpr struct {
	basenodeExpr
	Cond      ISimplifiedExpr
	TrueExpr  ISimplifiedExpr
	FalseExpr ISimplifiedExpr
	IsConst   bool
}

func (t *TernaryExpr) Kind() ExprKind { return ExprKindTernary }

// FunctionCallExpr represents a method or function invocation with arguments.
//
// Examples:
//   - myFunc()
//   - myFunc(1, 2, a + b)
type FunctionCallExpr struct {
	basenodeExpr
	Expr ISimplifiedExpr
	Args []ISimplifiedExpr
}

func (f *FunctionCallExpr) Kind() ExprKind { return ExprKindFunctionCall }

// ArrayAccessExpr represents a standard array indexing operation.
//
// Examples:
//   - myArray[5]
//   - myArray[index + 1]
type ArrayAccessExpr struct {
	basenodeExpr
	Expr  ISimplifiedExpr
	Index ISimplifiedExpr
}

func (a *ArrayAccessExpr) Kind() ExprKind { return ExprKindArrayAccess }

// BitSliceExpr represents DML bit slicing and endian-flagged array access operations.
//
// Examples:
//   - reg[15:8]        (Left=15, Right=8, Endian=nil)
//   - reg[31:0, be]    (Left=31, Right=0, Endian=be)
//   - myArray[5, le]   (Left=5, Right=nil, Endian=le)
type BitSliceExpr struct {
	basenodeExpr
	Expr   ISimplifiedExpr
	Left   ISimplifiedExpr
	Right  ISimplifiedExpr // nil when no colon part (endian-only form)
	Endian antlr.Token     // nil when no endian flag
}

func (b *BitSliceExpr) Kind() ExprKind { return ExprKindBitSlice }

// MemberAccessExpr represents dot or arrow notation to access struct or object members.
// Supports chained access (e.g., a.b->c.d) where Members and IsArrow are parallel slices.
//
// Examples:
//   - obj.field (Members=[field], IsArrow=[false])
//   - ptr->field (Members=[field], IsArrow=[true])
//   - a.b->c.d (Members=[b, c, d], IsArrow=[false, true, false])
type MemberAccessExpr struct {
	basenodeExpr
	Expr    ISimplifiedExpr
	Members []antlr.Token
	IsArrow []bool
}

func (m *MemberAccessExpr) Kind() ExprKind { return ExprKindMemberAccess }

// CastExpr represents a DML explicit type cast.
//
// Examples:
//   - cast(a, int)
//   - cast(ptr, const foo_t*)
type CastExpr struct {
	basenodeExpr
	Expr ISimplifiedExpr
	Type ICtypedeclContext
}

func (c *CastExpr) Kind() ExprKind { return ExprKindCast }

// SizeofExpr represents the sizeof operator.
//
// Examples:
//   - sizeof a
//   - sizeof(a)
type SizeofExpr struct {
	basenodeExpr
	Expr ISimplifiedExpr
}

func (s *SizeofExpr) Kind() ExprKind { return ExprKindSizeof }

// DefinedExpr represents the defined operator.
//
// Examples:
//   - defined a
type DefinedExpr struct {
	basenodeExpr
	Expr ISimplifiedExpr
}

func (d *DefinedExpr) Kind() ExprKind { return ExprKindDefined }

// StringifyExpr represents the stringify operator in DML.
//
// Examples:
//   - stringify(a)
type StringifyExpr struct {
	basenodeExpr
	Expr ISimplifiedExpr
}

func (s *StringifyExpr) Kind() ExprKind { return ExprKindStringify }

// NewExpr represents the new operator to instantiate types or arrays.
//
// Examples:
//   - new int
//   - new int[10] (where Size != nil)
type NewExpr struct {
	basenodeExpr
	Type ICtypedeclContext
	Size ISimplifiedExpr
}

func (n *NewExpr) Kind() ExprKind { return ExprKindNew }

// EachExpr represents the each-in generator in DML.
//
// Examples:
//   - each i in (myList)
type EachExpr struct {
	basenodeExpr
	Ident antlr.Token
	List  ISimplifiedExpr
}

func (e *EachExpr) Kind() ExprKind { return ExprKindEach }

// SizeoftypeExpr represents the sizeoftype operator.
//
// Examples:
//   - sizeoftype(int)
//   - sizeoftype foo_t
type SizeoftypeExpr struct {
	basenodeExpr
	Type ITypeopargContext
}

func (s *SizeoftypeExpr) Kind() ExprKind { return ExprKindSizeoftype }

// ParenExpr represents an expression wrapped in parentheses.
//
// Examples:
//   - (a + b)
type ParenExpr struct {
	basenodeExpr
	Expr ISimplifiedExpr
}

func (p *ParenExpr) Kind() ExprKind { return ExprKindParen }

// ArrayLiteralExpr represents an inline array definition.
//
// Examples:
//   - [1, 2, 3]
//   - [a, a + b]
type ArrayLiteralExpr struct {
	basenodeExpr
	Exprs []ISimplifiedExpr
}

func (a *ArrayLiteralExpr) Kind() ExprKind { return ExprKindArrayLiteral }

// WrapExpression takes a raw ANTLR expression context and returns our simplified AST node.
func WrapExpression(ctx IExpressionContext) ISimplifiedExpr {
	if ctx == nil {
		return nil
	}

	switch c := ctx.(type) {
	case *LiteralAtomContext:
		return &LiteralExpr{basenodeExpr: basenodeExpr{c}, Value: c.GetStart()}
	case *AdditiveExprContext:
		return &BinaryExpr{
			basenodeExpr: basenodeExpr{c},
			Op:           c.GetOp(),
			Left:         WrapExpression(c.Expression(0)),
			Right:        WrapExpression(c.Expression(1)),
		}
	case *MultiplicativeExprContext:
		return &BinaryExpr{
			basenodeExpr: basenodeExpr{c},
			Op:           c.GetOp(),
			Left:         WrapExpression(c.Expression(0)),
			Right:        WrapExpression(c.Expression(1)),
		}
	case *ShiftExprContext:
		return &BinaryExpr{
			basenodeExpr: basenodeExpr{c},
			Op:           c.GetOp(),
			Left:         WrapExpression(c.Expression(0)),
			Right:        WrapExpression(c.Expression(1)),
		}
	case *RelationalExprContext:
		return &BinaryExpr{
			basenodeExpr: basenodeExpr{c},
			Op:           c.GetOp(),
			Left:         WrapExpression(c.Expression(0)),
			Right:        WrapExpression(c.Expression(1)),
		}
	case *EqualityExprContext:
		return &BinaryExpr{
			basenodeExpr: basenodeExpr{c},
			Op:           c.GetOp(),
			Left:         WrapExpression(c.Expression(0)),
			Right:        WrapExpression(c.Expression(1)),
		}
	case *BitwiseExprContext:
		return &BinaryExpr{
			basenodeExpr: basenodeExpr{c},
			Op:           c.GetOp(),
			Left:         WrapExpression(c.Expression(0)),
			Right:        WrapExpression(c.Expression(1)),
		}
	case *LogicalExprContext:
		return &BinaryExpr{
			basenodeExpr: basenodeExpr{c},
			Op:           c.GetOp(),
			Left:         WrapExpression(c.Expression(0)),
			Right:        WrapExpression(c.Expression(1)),
		}
	case *TernaryIfExprContext:
		return &TernaryExpr{
			basenodeExpr: basenodeExpr{c},
			Cond:         WrapExpression(c.Expression(0)),
			TrueExpr:     WrapExpression(c.Expression(1)),
			FalseExpr:    WrapExpression(c.Expression(2)),
			IsConst:      false,
		}
	case *TernaryConstIfExprContext:
		return &TernaryExpr{
			basenodeExpr: basenodeExpr{c},
			Cond:         WrapExpression(c.Expression(0)),
			TrueExpr:     WrapExpression(c.Expression(1)),
			FalseExpr:    WrapExpression(c.Expression(2)),
			IsConst:      true,
		}
	case *UnaryExprContext:
		return &UnaryExpr{
			basenodeExpr: basenodeExpr{c},
			Op:           c.GetOp(),
			Expr:         WrapExpression(c.Expression()),
			IsPost:       false,
		}
	case *NotExprContext:
		return &UnaryExpr{
			basenodeExpr: basenodeExpr{c},
			Op:           c.BANG().GetSymbol(),
			Expr:         WrapExpression(c.Expression()),
			IsPost:       false,
		}
	case *BitwiseNotExprContext:
		return &UnaryExpr{
			basenodeExpr: basenodeExpr{c},
			Op:           c.TILDE().GetSymbol(),
			Expr:         WrapExpression(c.Expression()),
			IsPost:       false,
		}
	case *AddressofExprContext:
		return &UnaryExpr{
			basenodeExpr: basenodeExpr{c},
			Op:           c.BITAND().GetSymbol(),
			Expr:         WrapExpression(c.Expression()),
			IsPost:       false,
		}
	case *DerefExprContext:
		return &UnaryExpr{
			basenodeExpr: basenodeExpr{c},
			Op:           c.MUL().GetSymbol(),
			Expr:         WrapExpression(c.Expression()),
			IsPost:       false,
		}
	case *PreIncrementExprContext:
		return &UnaryExpr{
			basenodeExpr: basenodeExpr{c},
			Op:           c.GetOp(),
			Expr:         WrapExpression(c.Expression()),
			IsPost:       false,
		}
	case *PostIncrementExprContext:
		return &UnaryExpr{
			basenodeExpr: basenodeExpr{c},
			Op:           c.GetOp(),
			Expr:         WrapExpression(c.Expression()),
			IsPost:       true,
		}
	case *CastExprContext:
		return &CastExpr{
			basenodeExpr: basenodeExpr{c},
			Expr:         WrapExpression(c.Expression()),
			Type:         c.Ctypedecl(),
		}
	case *SizeofExprContext:
		return &SizeofExpr{
			basenodeExpr: basenodeExpr{c},
			Expr:         WrapExpression(c.Expression()),
		}
	case *DefinedExprContext:
		return &DefinedExpr{
			basenodeExpr: basenodeExpr{c},
			Expr:         WrapExpression(c.Expression()),
		}
	case *StringifyExprContext:
		return &StringifyExpr{
			basenodeExpr: basenodeExpr{c},
			Expr:         WrapExpression(c.Expression()),
		}
	case *ParenExprContext:
		return &ParenExpr{
			basenodeExpr: basenodeExpr{c},
			Expr:         WrapExpression(c.Expression()),
		}
	case *MemberAccessExprContext:
		objidents := c.AllExpressionAccessIdent()
		members := make([]antlr.Token, len(objidents))
		isArrowSlice := make([]bool, len(objidents))
		for i, oid := range objidents {
			members[i] = oid.GetStart()
			// Walk children to find the operator (DOT or ARROW) immediately preceding this objident.
			// The children of the context are: [Expression] [op0] [Objident0] [op1] [Objident1] ...
			for j := range c.GetChildCount() {
				if c.GetChild(j) == oid {
					// The operator is the child immediately before this objident
					if j > 0 {
						if tn, ok := c.GetChild(j - 1).(antlr.TerminalNode); ok {
							isArrowSlice[i] = tn.GetText() == "->"
						}
					}
					break
				}
			}
		}
		return &MemberAccessExpr{
			basenodeExpr: basenodeExpr{c},
			Expr:         WrapExpression(c.Expression()),
			Members:      members,
			IsArrow:      isArrowSlice,
		}
	case *ArrayAccessExprContext:
		return &ArrayAccessExpr{
			basenodeExpr: basenodeExpr{c},
			Expr:         WrapExpression(c.Expression(0)),
			Index:        WrapExpression(c.Expression(1)),
		}
	case *BitSliceExprContext:
		bs := c.BitSlice()
		var endian antlr.Token
		if ef := bs.Endianflag(); ef != nil {
			endian = ef.Ident().GetStart()
		}
		var right ISimplifiedExpr
		if bs.Expression() != nil {
			right = WrapExpression(bs.Expression())
		}
		return &BitSliceExpr{
			basenodeExpr: basenodeExpr{c},
			Expr:         WrapExpression(c.Expression(0)),
			Left:         WrapExpression(c.Expression(1)),
			Right:        right,
			Endian:       endian,
		}
	case *FunctionCallExprContext:
		var args []ISimplifiedExpr
		if argsList := c.ArgumentsList(); argsList != nil {
			if initArray := argsList.SingleInitializerArray(); initArray != nil {
				for _, init := range initArray.AllSingleInitializer() {
					if expr := init.Expression(); expr != nil {
						args = append(args, WrapExpression(expr))
					} else {
						args = append(
							args,
							&LiteralExpr{
								basenodeExpr: basenodeExpr{nil},
								Value:        init.GetStart(),
							},
						)
					}
				}
			}
		}
		return &FunctionCallExpr{
			basenodeExpr: basenodeExpr{c},
			Expr:         WrapExpression(c.Expression()),
			Args:         args,
		}
	case *NewExprContext:
		var size ISimplifiedExpr
		if c.Expression() != nil {
			size = WrapExpression(c.Expression())
		}
		return &NewExpr{
			basenodeExpr: basenodeExpr{c},
			Type:         c.Ctypedecl(),
			Size:         size,
		}
	case *EachExprContext:
		return &EachExpr{
			basenodeExpr: basenodeExpr{c},
			Ident:        c.Objident().GetStart(),
			List:         WrapExpression(c.Expression()),
		}
	case *SizeoftypeExprContext:
		return &SizeoftypeExpr{
			basenodeExpr: basenodeExpr{c},
			Type:         c.Typeoparg(),
		}
	case *ArrayLiteralExprContext:
		var exprs []ISimplifiedExpr
		if list := c.Expression_list(); list != nil {
			if inner := list.Expression_list_ntc_nonempty(); inner != nil {
				for _, expr := range inner.AllExpression() {
					exprs = append(exprs, WrapExpression(expr))
				}
			}
		}
		return &ArrayLiteralExpr{basenodeExpr: basenodeExpr{c}, Exprs: exprs}
	default:
		panic(fmt.Sprintf("Unimplemented expression context wrapper: %T", c))
	}
}
