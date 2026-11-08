// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"errors"
	"fmt"
	"strings"
)

// Sentinel errors representing mask categories, inspected via errors.Is.
var (
	ErrMaskIncorrectFormatting = errors.New("incorrect formatting")
	ErrMaskSyntaxError         = errors.New("syntax error")
	ErrMaskSkip                = errors.New("skip test")
)

// Masked represents the base interface for a masked test case.
// It implements the standard library error interface and exposes the issue code.
type Masked interface {
	error
	Issue() int
}

// MaskedFormatter is implemented by masks that expect a specific buggy formatted output.
type MaskedFormatter interface {
	Masked
	Expected() string
	WithExpected(expect string) MaskedFormatter
}

// MaskedSkip is implemented by masks that skip test execution upfront.
type MaskedSkip interface {
	Masked
	Reason() string
}

// MaskedSyntax is implemented by masks representing known parser syntax errors.
type MaskedSyntax interface {
	Masked
	SyntaxMsg() string
}

// MaskIncorrectFormatting captures an expected buggy formatting result associated with an issue.
type MaskIncorrectFormatting struct {
	issue  int
	expect string
}

func (m MaskIncorrectFormatting) Error() string {
	if m.issue == 0 {
		return fmt.Sprintf("[Uncategorized] incorrect formatting: %q", m.expect)
	}
	return fmt.Sprintf(
		"[Issue #%d] incorrect formatting: %q",
		m.issue,
		m.expect,
	)
}

func (m MaskIncorrectFormatting) Issue() int {
	return m.issue
}

func (m MaskIncorrectFormatting) Expected() string {
	return m.expect
}

func (m MaskIncorrectFormatting) WithExpected(expect string) MaskedFormatter {
	return MaskIncorrectFormatting{
		issue:  m.issue,
		expect: expect,
	}
}

func (m MaskIncorrectFormatting) Unwrap() error {
	return ErrMaskIncorrectFormatting
}

// MaskSkip captures an upfront skip reason associated with an issue.
type MaskSkip struct {
	issue  int
	reason string
}

func (m MaskSkip) Error() string {
	if m.issue == 0 {
		return fmt.Sprintf("[Uncategorized] skip test: %s", m.reason)
	}
	return fmt.Sprintf("[Issue #%d] skip test: %s", m.issue, m.reason)
}

func (m MaskSkip) Issue() int {
	return m.issue
}

func (m MaskSkip) Reason() string {
	return m.reason
}

func (m MaskSkip) Unwrap() error {
	return ErrMaskSkip
}

// MaskSyntaxError captures a known syntax/parsing failure associated with an issue.
type MaskSyntaxError struct {
	issue int
	msg   string
}

func (m MaskSyntaxError) Error() string {
	if m.issue <= 0 {
		return fmt.Sprintf("[Uncategorized] syntax error: %s", m.msg)
	}
	return fmt.Sprintf("[Issue #%d] syntax error: %s", m.issue, m.msg)
}

func (m MaskSyntaxError) Issue() int {
	return m.issue
}

func (m MaskSyntaxError) SyntaxMsg() string {
	return m.msg
}

func (m MaskSyntaxError) Unwrap() error {
	return ErrMaskSyntaxError
}

// Generic constructors

func NewMaskIncorrectFormatting(
	issue int,
	expect string,
) MaskIncorrectFormatting {
	return MaskIncorrectFormatting{
		issue:  issue,
		expect: expect,
	}
}

func NewMaskSkip(issue int, reason string) MaskSkip {
	return MaskSkip{
		issue:  issue,
		reason: reason,
	}
}

func NewMaskSyntaxError(issue int, msg string) MaskSyntaxError {
	return MaskSyntaxError{
		issue: issue,
		msg:   msg,
	}
}

// Helpers and tools

// normalizeMask adapts nil, string, or Masked into a canonical Masked representation.
// Legacy non-empty strings are converted to MaskUncategorizedFormatting.
func normalizeMask(raw any) (Masked, error) {
	if raw == nil {
		return nil, nil
	}

	switch v := raw.(type) {
	case string:
		if v == "" {
			return nil, nil
		}
		return MaskUncategorizedFormatting(
			strings.ReplaceAll(v, "`", "\""),
		), nil
	case Masked:
		if fmtMask, ok := v.(MaskedFormatter); ok {
			return fmtMask.WithExpected(
				strings.ReplaceAll(fmtMask.Expected(), "`", "\""),
			), nil
		}
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported mask type: %T", raw)
	}
}

// MaskTransform applies a string transformation function fn to a masked value.
// - If m is nil, it returns nil.
// - If m is a string, it returns fn(m).
// - If m is a MaskedFormatter, it returns WithExpected(fn(formatter.Expected())).
// - If m is any other Masked (e.g. MaskSkip, MaskSyntaxError), it returns m unchanged.
func MaskTransform(m Masked, fn func(string) string) Masked {
	if m == nil {
		return nil
	}
	switch v := m.(type) {
	case MaskedFormatter:
		return v.WithExpected(fn(v.Expected()))
	default:
		return m
	}
}

// MaskReplaceAll replaces all occurrences of old with new in a string or MaskedFormatter mask.
func MaskReplaceAll(m Masked, old, new string) Masked {
	return MaskTransform(m, func(s string) string {
		return strings.ReplaceAll(s, old, new)
	})
}

// Convenience & issue-specific helpers

// MaskUncategorizedFormatting marks untracked or uncategorized formatting discrepancies
// where the printer produces buggy output different from expected.
func MaskUncategorizedFormatting(expect string) MaskIncorrectFormatting {
	return NewMaskIncorrectFormatting(0, expect)
}

// MaskUnimplemented marks unimplemented syntax nodes where the printer produces
// buggy or unformatted output different from expected.
func MaskUnimplemented(expect ...string) MaskIncorrectFormatting {
	exp := ""
	if len(expect) > 0 {
		exp = expect[0]
	}
	return NewMaskIncorrectFormatting(0, exp)
}

// MaskIncorrectParenWhitespace marks formatting issues involving incorrect whitespace
// around parentheses.
func MaskIncorrectParenWhitespace(expect string) MaskIncorrectFormatting {
	return NewMaskIncorrectFormatting(-1, expect)
}

// MaskInputNewlinePreservation marks Issue #14: User-authored line breaks inside compile-time
// arrays [...], template is (...), or in each (...) statements are collapsed into a single line
// when fitting within line cap instead of preserving multiline layout.
//
// Issue #14: https://git.home.timoschenko.dev/stimoshc/dmltools/issues/14
func MaskInputNewlinePreservation(expect string) MaskIncorrectFormatting {
	return NewMaskIncorrectFormatting(14, expect)
}

// MaskArgumentOverWrappingWithLineComments marks Issue #16: In argument lists (func(...),
// cast(...)), initializers ({...}), or arrays ([...]), placing a line comment on the trailing
// element unexpectedly forces all preceding arguments onto separate lines.
//
// Issue #16: https://git.home.timoschenko.dev/stimoshc/dmltools/issues/16
func MaskArgumentOverWrappingWithLineComments(
	expect string,
) MaskIncorrectFormatting {
	return NewMaskIncorrectFormatting(16, expect)
}

// MaskSingleElementTupleParenthesesRemoval marks Issue #17: Redundant parentheses in
// single-variable tuple declarations (e.g. local (int x);) are retained instead of unwrapping
// to local int x;.
//
// Issue #17: https://git.home.timoschenko.dev/stimoshc/dmltools/issues/17
func MaskSingleElementTupleParenthesesRemoval(
	expect string,
) MaskIncorrectFormatting {
	return NewMaskIncorrectFormatting(17, expect)
}

// MaskTask006 marks Issue #18: Test failure due to compiler workaround / synthetic token injection.
//
// Issue #18: https://git.home.timoschenko.dev/stimoshc/dmltools/issues/18
func MaskTask006(expect string) MaskIncorrectFormatting {
	return NewMaskIncorrectFormatting(18, expect)
}

// MaskTask006Skip marks Issue #18 test cases that must be skipped upfront due to AST visitor panics.
//
// Issue #18: https://git.home.timoschenko.dev/stimoshc/dmltools/issues/18
func MaskTask006Skip(reason ...string) MaskSkip {
	r := "Inject synthetic tokens"
	if len(reason) > 0 {
		r = reason[0]
	}
	return NewMaskSkip(18, r)
}

// MaskSyntheticParens marks Issue #18: Missing synthetic parentheses wrapping multiline
// expressions (e.g. member access, each ... in, or binary expressions).
//
// Issue #18: https://git.home.timoschenko.dev/stimoshc/dmltools/issues/18
func MaskSyntheticParens(expect string) MaskIncorrectFormatting {
	return NewMaskIncorrectFormatting(18, expect)
}

// MaskSyntheticProvisionalParens marks Issue #18: Missing synthetic parentheses in multiline
// provisional declarations or failure to strip redundant parentheses.
//
// Issue #18: https://git.home.timoschenko.dev/stimoshc/dmltools/issues/18
func MaskSyntheticProvisionalParens(expect string) MaskIncorrectFormatting {
	return NewMaskIncorrectFormatting(18, expect)
}

// MaskSyntheticSizeofParens marks Issue #18: Missing synthetic parentheses around
// sizeof/defined/sizeoftype operands.
//
// Issue #18: https://git.home.timoschenko.dev/stimoshc/dmltools/issues/18
func MaskSyntheticSizeofParens(expect string) MaskIncorrectFormatting {
	return NewMaskIncorrectFormatting(18, expect)
}

// MaskSyntheticTrailingComma marks Issue #18: Trailing comma not injected in multiline arguments
// or not stripped in single-line initializers.
//
// Issue #18: https://git.home.timoschenko.dev/stimoshc/dmltools/issues/18
func MaskSyntheticTrailingComma(expect string) MaskIncorrectFormatting {
	return NewMaskIncorrectFormatting(18, expect)
}

// MaskUnaryIncorrectFormatting captures expected formatting for unary operators across
// prefix and postfix operator variations.
type MaskUnaryIncorrectFormatting struct {
	issue        int
	preWrapped   string
	preUnwrapped string
	post         string
}

func (m MaskUnaryIncorrectFormatting) Error() string {
	return fmt.Sprintf("[Issue #%d] unary op incorrect formatting", m.issue)
}

func (m MaskUnaryIncorrectFormatting) Issue() int {
	return m.issue
}

func (m MaskUnaryIncorrectFormatting) Resolve(
	op string,
	isPost bool,
) MaskIncorrectFormatting {
	if isPost {
		return NewMaskIncorrectFormatting(
			m.issue,
			strings.ReplaceAll(m.post, "%p", op),
		)
	}
	switch op {
	case "-", "+", "&", "*":
		return NewMaskIncorrectFormatting(
			m.issue,
			strings.ReplaceAll(m.preWrapped, "%b", op),
		)
	default:
		return NewMaskIncorrectFormatting(
			m.issue,
			strings.ReplaceAll(m.preUnwrapped, "%b", op),
		)
	}
}

// MaskSyntheticParensUnary marks Issue #18 for unary operator expression tests where
// binary-capable prefix ops (-/+/&/*), pure prefix ops (++/--/~/!), and postfix ops (++/--)
// format with different wrapping behavior.
//
// Issue #18: https://git.home.timoschenko.dev/stimoshc/dmltools/issues/18
func MaskSyntheticParensUnary(
	preWrapped, preUnwrapped, post string,
) MaskUnaryIncorrectFormatting {
	return MaskUnaryIncorrectFormatting{
		issue:        18,
		preWrapped:   preWrapped,
		preUnwrapped: preUnwrapped,
		post:         post,
	}
}

// MaskMultilineBlockCommentIndentation marks Issue #21: Internal indentation and leading
// whitespace inside trailing multiline block comments (/* ... */) are preserved verbatim rather
// than normalized to the comment header column.
//
// Issue #21: https://git.home.timoschenko.dev/stimoshc/dmltools/issues/21
func MaskMultilineBlockCommentIndentation(
	expect string,
) MaskIncorrectFormatting {
	return NewMaskIncorrectFormatting(21, expect)
}

// MaskMissingWrappingForMultilineBlockCommentsInIs marks Issue #22: A multiline block comment
// inside a template header or is (...) clause fails to trigger multiline wrapping and indentation
// of header arguments.
//
// Issue #22: https://git.home.timoschenko.dev/stimoshc/dmltools/issues/22
func MaskMissingWrappingForMultilineBlockCommentsInIs(
	expect string,
) MaskIncorrectFormatting {
	return NewMaskIncorrectFormatting(22, expect)
}

// MaskBlockCommentWrapping marks Issue #23: Incorrect line wrapping around inline block
// comments (/* ... */), such as an unexpected newline after /*b*/ in unary expressions or a
// missing break after qualifier comment /*d*/ before 'method' when exceeding line cap.
//
// Issue #23: https://git.home.timoschenko.dev/stimoshc/dmltools/issues/23
func MaskBlockCommentWrapping(expect string) MaskIncorrectFormatting {
	return NewMaskIncorrectFormatting(23, expect)
}

// MaskTernaryOperatorWrappingAndIndentation marks Issue #26: In chained ternary expressions or
// after multiline parenthesized conditions, an unexpected newline is emitted before '?' / '#?'
// instead of keeping '(cond) ? <val>' on the same line.
//
// Issue #26: https://git.home.timoschenko.dev/stimoshc/dmltools/issues/26
func MaskTernaryOperatorWrappingAndIndentation(
	expect string,
) MaskIncorrectFormatting {
	return NewMaskIncorrectFormatting(26, expect)
}

// MaskReturnTypeOverWrapping marks Issue #27: A single return type '(bool)' that is already on
// its own line after a line comment gets unexpectedly expanded across multiple lines.
//
// Issue #27: https://git.home.timoschenko.dev/stimoshc/dmltools/issues/27
func MaskReturnTypeOverWrapping(expect string) MaskIncorrectFormatting {
	return NewMaskIncorrectFormatting(27, expect)
}

// MaskBinaryOperatorWrappingAfterMultilineParen marks Issue #43: Binary/relational operators are
// forced onto a new line after the closing ')' of a multiline parenthesized operand due to
// ancestor wrap-group trigger propagation.
//
// Issue #43: https://git.home.timoschenko.dev/stimoshc/dmltools/issues/43
func MaskBinaryOperatorWrappingAfterMultilineParen(
	expect string,
) MaskIncorrectFormatting {
	return NewMaskIncorrectFormatting(43, expect)
}

// MaskOpenDelimiterFollowedByLineComment marks formatting discrepancies where a line comment
// appearing on the next line after an opening delimiter '(', '{', or '[' improperly collapses
// onto the same line as the delimiter instead of preserving its line position.
func MaskOpenDelimiterFollowedByLineComment(
	expect string,
) MaskIncorrectFormatting {
	return NewMaskIncorrectFormatting(0, expect)
}
