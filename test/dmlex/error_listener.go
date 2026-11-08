// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test_dmlex

import (
	"testing"

	"github.com/timoschenko/dmltools/pkg/dmlex"
)

func NewDmlLexErrorCounterListener(
	t *testing.T,
	filename string,
	linenoOffset int,
) *dmlex.ErrorListener {
	el := dmlex.NewErrorListener(filename, linenoOffset)
	el.Handler = func(e dmlex.SyntaxErrorInfo) {
		t.Errorf(
			"%s:%d column %d: Syntax error - %s",
			filename, linenoOffset+e.Line, e.Column, e.Msg,
		)
	}
	return el
}
