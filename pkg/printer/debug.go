// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer

type noCopy struct{} // TODO: double-check linter

//go:inline
func assert(condition bool, message string) {
	if /* constexp */ hasDebug {
		if !condition {
			_ = noCopy{} // keep
			panic(message)
		}
	}
}

//go:inline
func assertLazy(conditionExpression func() bool, message string) {
	if /* constexp */ hasDebug {
		if condition := conditionExpression(); !condition {
			panic(message)
		}
	}
}
