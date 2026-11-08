// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package dmlex

import (
	"errors"

	"github.com/antlr4-go/antlr/v4"
)

type SyntaxErrorInfo struct {
	Line   int
	Column int
	Msg    string
	Symbol string
	Length int
	Err    error
}

type ErrorListener struct {
	*antlr.DefaultErrorListener
	Errors       []SyntaxErrorInfo
	Filename     string
	LinenoOffset int
	Handler      func(e SyntaxErrorInfo)
}

func NewErrorListener(filename string, linenoOffset int) *ErrorListener {
	return &ErrorListener{
		DefaultErrorListener: antlr.NewDefaultErrorListener(),
		Errors:               make([]SyntaxErrorInfo, 0),
		Filename:             filename,
		LinenoOffset:         linenoOffset,
	}
}

func (l *ErrorListener) SyntaxError(
	recognizer antlr.Recognizer,
	offendingSymbol any,
	line, column int,
	msg string,
	e antlr.RecognitionException,
) {
	symText := ""
	length := 1
	if tok, ok := offendingSymbol.(antlr.Token); ok && tok != nil {
		symText = tok.GetText()
		if symText != "<EOF>" && len(symText) > 0 {
			length = len(symText)
		}
	}
	errInfo := SyntaxErrorInfo{
		Line:   line,
		Column: column,
		Msg:    msg,
		Symbol: symText,
		Length: length,
	}
	if err, ok := e.(error); ok {
		errInfo.Err = err
	}
	l.Errors = append(l.Errors, errInfo)
	if l.Handler != nil {
		l.Handler(errInfo)
	}
}

func (l *ErrorListener) HasError(target error) bool {
	for _, e := range l.Errors {
		if errors.Is(e.Err, target) {
			return true
		}
	}
	return false
}
