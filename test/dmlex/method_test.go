// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test_dmlex

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	lex "github.com/timoschenko/dmltools/pkg/dmlex"

	"github.com/antlr4-go/antlr/v4"
)

const CODE_ALL_METHODS_METHODS_COUNT = 32
const CODE_ALL_LINENO = 16
const CODE_ALL_METHODS = `
//     method commented_m     ()                                                                                                      { }
/*     method commentedi_m    ()                                                                                                      { } */
       method m               ()                                                                                                      { }
       method m_1a            (int a)                                                                                                 { }
       method m_2a            (int a, int b)                                                                                          { }
       method m_3a            (int a, int b, int b)                                                                                   { }
       method m_1r            ()                                             -> (int)                                                 { }
       method m_2r            ()                                             -> (int, int)                                            { }
       method m_3r            ()                                             -> (int, int, int)                                       { }
       method m_1a_1r         (int a)                                        -> (int)                                                 { }
       method m_1a_2r         (int a)                                        -> (int, int)                                            { }
       method m_1a_3r         (int a)                                        -> (int, int, int)                                       { }
       method m_3ac           (const foo_t a1, foo_t *a2, foo_t ** const a3)                                                          { }
       method m_3rc           ()                                             -> (const foo_t, foo_t *, foo_t ** const)                { }
       method m_t             ()                                                                                       throws         { }
       method m_1r_t          ()                                             -> (int)                                  throws         { }
       method m_1a_2r_t       (int a)                                        -> (int, int)                             throws         { }
       method m_1r_d          ()                                             -> (int)                                         default { }
       method m_1a_2r_d       (int a)                                        -> (int, int)                                    default { }
       method m_1r_d_t        ()                                             -> (int)                                  throws default { }
       method m_1a_2r_d_t     (int a)                                        -> (int, int)                             throws default { }
inline method inline_m        ()                                                                                                      { }
inline method inline_m_1ia    (inline ai)                                                                                             { }
inline method inline_m_2ia    (inline ai, inline bi)                                                                                  { }
inline method inline_m_2a_2ia (int a, inline ai, int b, inline bi)                                                                    { }
       method BAD_m_1a_1ia    (int a, inline ai)                                                                                      { }

       independent startup memoized  method m_memoized()                                                                              { }
       independent startup           method m_startup()                                                                               { }
       independent                   method m_independent()                                                                           { }

shared independent startup memoized  method m_shared_memoized()                                                                       { }
shared independent startup           method m_shared_startup()                                                                        { }
shared independent                   method m_shared_independent()                                                                    { }

template tmpl {
	method t_m() { }

	shared method t_shared();
}
`

func TestSimpleMethod(t *testing.T) {
	t.Parallel()

	const EXPECTED_METHOD_NAME = "fancyname"
	_, filename, linenoOffset, _ := runtime.Caller(0)
	const CODE = `dml 1.4;

	// A simple method
	method fancyname() {
	}`

	tree := ParseDmlText(t, CODE, filename, linenoOffset+1)
	v := NewMethodCollector(t)
	antlr.ParseTreeWalkerDefault.Walk(v, tree)
	methods_count := len(v.Methods)
	if methods_count != 1 {
		if methods_count == 0 {
			fmt.Println("No methods were found.")
		} else {
			for name := range v.Methods {
				t.Fatalf("Too many methods: '%s'", name)
			}
		}
	} else {
		for name := range v.Methods {
			if name != EXPECTED_METHOD_NAME {
				t.Fatalf(
					"Method '%s' is unexpected, need '%s'",
					name,
					EXPECTED_METHOD_NAME,
				)
			}
		}
	}
}

func TestAllMethodsSimple(t *testing.T) {
	t.Parallel()

	_, filename, _, _ := runtime.Caller(0)
	tree := ParseDmlText(t, CODE_ALL_METHODS, filename, CODE_ALL_LINENO)
	v := NewMethodCollector(t)
	antlr.ParseTreeWalkerDefault.Walk(v, tree)
	methods_count := len(v.Methods)
	if methods_count != CODE_ALL_METHODS_METHODS_COUNT {
		if methods_count == 0 {
			fmt.Println("No methods were found.")
		} else {
			t.Fatalf(
				"Too many methods %d, but expected %d",
				methods_count,
				CODE_ALL_METHODS_METHODS_COUNT,
			)
			for name := range v.Methods {
				t.Fatalf("Too many methods: '%s'", name)
			}
		}
	}
}

func TestAllMethodsCheckThrows(t *testing.T) {
	t.Parallel()

	_, filename, _, _ := runtime.Caller(0)
	tree := ParseDmlText(t, CODE_ALL_METHODS, filename, CODE_ALL_LINENO)
	v := NewMethodCollector(t)
	antlr.ParseTreeWalkerDefault.Walk(v, tree)

	const EXPECTED = 5
	found := 0
	for _, m := range v.Methods {
		if lex.MethodThrows(m) {
			found++
		}
	}
	if found != EXPECTED {
		t.Errorf("Found %d 'throws' tags but expected %d", found, EXPECTED)
	}
}

func TestAllMethodsCheckDefault(t *testing.T) {
	t.Parallel()

	_, filename, _, _ := runtime.Caller(0)
	tree := ParseDmlText(t, CODE_ALL_METHODS, filename, CODE_ALL_LINENO)
	v := NewMethodCollector(t)
	antlr.ParseTreeWalkerDefault.Walk(v, tree)

	const EXPECTED = 4
	found := 0
	for _, m := range v.Methods {
		if lex.MethodDefault(m) {
			found++
		}
	}
	if found != EXPECTED {
		t.Errorf("Found %d 'throws' tags but expected %d", found, EXPECTED)
	}
}

func TestAllMethodsCheckShared(t *testing.T) {
	t.Parallel()

	_, filename, _, _ := runtime.Caller(0)
	tree := ParseDmlText(t, CODE_ALL_METHODS, filename, CODE_ALL_LINENO)
	v := NewMethodCollector(t)
	antlr.ParseTreeWalkerDefault.Walk(v, tree)

	const EXPECTED = 3 + 1
	found := 0
	for _, m := range v.Methods {
		if lex.MethodIsShared(m) {
			found++
		}
	}
	if found != EXPECTED {
		t.Errorf("Found %d 'throws' tags but expected %d", found, EXPECTED)
	}
}

func TestAllMethodsCheckIndependent(t *testing.T) {
	t.Parallel()

	_, filename, _, _ := runtime.Caller(0)
	tree := ParseDmlText(t, CODE_ALL_METHODS, filename, CODE_ALL_LINENO)
	v := NewMethodCollector(t)
	antlr.ParseTreeWalkerDefault.Walk(v, tree)

	const EXPECTED = 3 + 3
	found := 0
	for _, m := range v.Methods {
		if lex.MethodIsIndependent(m) {
			found++
		}
	}
	if found != EXPECTED {
		t.Errorf("Found %d 'throws' tags but expected %d", found, EXPECTED)
	}
}

func TestAllMethodsCheckStartup(t *testing.T) {
	t.Parallel()

	_, filename, _, _ := runtime.Caller(0)
	tree := ParseDmlText(t, CODE_ALL_METHODS, filename, CODE_ALL_LINENO)
	v := NewMethodCollector(t)
	antlr.ParseTreeWalkerDefault.Walk(v, tree)

	const EXPECTED = 2 + 2
	found := 0
	for _, m := range v.Methods {
		if lex.MethodIsStartup(m) {
			found++
		}
	}
	if found != EXPECTED {
		t.Errorf("Found %d 'throws' tags but expected %d", found, EXPECTED)
	}
}

func TestAllMethodsCheckMemoized(t *testing.T) {
	t.Parallel()

	_, filename, _, _ := runtime.Caller(0)
	tree := ParseDmlText(t, CODE_ALL_METHODS, filename, CODE_ALL_LINENO)
	v := NewMethodCollector(t)
	antlr.ParseTreeWalkerDefault.Walk(v, tree)

	const EXPECTED = 2
	found := 0
	for _, m := range v.Methods {
		if lex.MethodIsMemoized(m) {
			found++
		}
	}
	if found != EXPECTED {
		t.Errorf("Found %d 'throws' tags but expected %d", found, EXPECTED)
	}
}

func buildLookupMap(list []string) map[string]struct{} {
	lookup := make(map[string]struct{}, len(list))
	for _, v := range list {
		lookup[v] = struct{}{}
	}
	return lookup
}

func TestAllMethodsRetuns(t *testing.T) {
	t.Parallel()

	EXPECTED_OUTPUTS := map[string][]string{
		"m_1r":        {"int"},
		"m_2r":        {"int", "int"},
		"m_3r":        {"int", "int", "int"},
		"m_1a_1r":     {"int"},
		"m_1a_2r":     {"int", "int"},
		"m_1a_3r":     {"int", "int", "int"},
		"m_3rc":       {"const foo_t", "foo_t*", "foo_t**const"},
		"m_1r_t":      {"int"},
		"m_1a_2r_t":   {"int", "int"},
		"m_1r_d":      {"int"},
		"m_1a_2r_d":   {"int", "int"},
		"m_1r_d_t":    {"int"},
		"m_1a_2r_d_t": {"int", "int"},
	}

	_, filename, _, _ := runtime.Caller(0)
	tree := ParseDmlText(t, CODE_ALL_METHODS, filename, CODE_ALL_LINENO)
	v := NewMethodCollector(t)
	antlr.ParseTreeWalkerDefault.Walk(v, tree)

	for methodName, m := range v.Methods {
		expect_rets := 0
		if strings.Contains(methodName, "_1r") {
			expect_rets = 1
		} else if strings.Contains(methodName, "_2r") {
			expect_rets = 2
		} else if strings.Contains(methodName, "_3r") {
			expect_rets = 3
		}
		outputs := lex.MethodReturns(m)
		if outputs == nil {
			if expect_rets != 0 {
				t.Errorf(
					"Method '%s' has zero returns but expects %d",
					methodName,
					expect_rets,
				)
			}
			continue
		}
		if len(outputs) != expect_rets {
			t.Errorf(
				"Method '%s' has %d returns but expects %d",
				methodName,
				len(outputs),
				expect_rets,
			)
		}
		expectedOutputs, ok := EXPECTED_OUTPUTS[methodName]
		if !ok {
			t.Errorf("Method '%s' has not in EXPECTED_OUTPUTS", methodName)
			continue
		}
		lookupMap := buildLookupMap(expectedOutputs)
		for i, output := range outputs {
			name := lex.CdeclTypeName(output)
			if _, ok := lookupMap[name]; !ok {
				t.Errorf(
					"Method '%s' return '%s' but expected '%s'",
					methodName,
					name,
					expectedOutputs[i],
				)
			}
		}
	}
}
