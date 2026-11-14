// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package dmlex

import (
	"errors"
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

const SupportedDmlVersion = "1.4"

// ErrUnimplemented is the error for unimplemented DML language versions/features.
var ErrUnimplemented = errors.New("Error Unimplemented")

// EUNIMPL is an alias for ErrUnimplemented matching the DML error diagnostic code.
//
//nolint:stylecheck,revive
var EUNIMPL = ErrUnimplemented

// UnimplementedVersionError represents an unsupported DML version encountered in dmlTag.
type UnimplementedVersionError struct {
	*antlr.BaseRecognitionException
	Version string
	Token   antlr.Token
	Message string
}

func NewUnimplementedVersionError(
	recognizer antlr.Parser,
	offendingToken antlr.Token,
	version string,
) *UnimplementedVersionError {
	msg := fmt.Sprintf(
		"%s: dml %s syntax is not supported yet",
		EUNIMPL.Error(),
		version,
	)
	base := antlr.NewBaseRecognitionException(
		msg,
		recognizer,
		recognizer.GetInputStream(),
		recognizer.GetParserRuleContext(),
	)
	return &UnimplementedVersionError{
		BaseRecognitionException: base,
		Version:                  version,
		Token:                    offendingToken,
		Message:                  msg,
	}
}

func (e *UnimplementedVersionError) Error() string {
	return e.Message
}

func (e *UnimplementedVersionError) Unwrap() error {
	return EUNIMPL
}

func (e *UnimplementedVersionError) GetOffendingToken() antlr.Token {
	return e.Token
}

func (e *UnimplementedVersionError) GetMessage() string {
	return e.Message
}

func (p *dmlParser) checkDmlVersion(localctx IDmlTagContext) {
	if localctx == nil || p.HasError() {
		return
	}
	floatNode := localctx.FLOAT_LITERAL()
	if floatNode == nil {
		return
	}
	tok := floatNode.GetSymbol()
	if tok == nil {
		return
	}
	version := tok.GetText()
	if version == SupportedDmlVersion {
		return
	}

	err := NewUnimplementedVersionError(p, tok, version)
	localctx.SetException(err)
	if parent, ok := localctx.GetParent().(antlr.ParserRuleContext); ok &&
		parent != nil {
		parent.SetException(err)
	}
	p.NotifyErrorListeners(err.GetMessage(), tok, err)

	// Remove lexer error listeners so subsequent unsupported tokens (e.g. DML 1.2 '$')
	// do not trigger noisy lexer recognition errors after EUNIMPL.
	if recognizer, ok := p.GetTokenStream().GetTokenSource().(antlr.Recognizer); ok {
		recognizer.RemoveErrorListeners()
	}

	// Immediately advance the token stream to EOF so dmlCode is skipped cleanly.
	for p.GetTokenStream().LA(1) != antlr.TokenEOF {
		p.GetTokenStream().Consume()
	}
}

// UnimplementedVersion returns the unsupported version string if any error is EUNIMPL.
func (l *ErrorListener) UnimplementedVersion() string {
	for _, e := range l.Errors {
		if verErr, ok := errors.AsType[*UnimplementedVersionError](e.Err); ok {
			return verErr.Version
		}
	}
	return ""
}
