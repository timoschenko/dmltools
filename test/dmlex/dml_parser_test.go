// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test_dmlex

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/timoschenko/dmltools/pkg/dmlex"

	"github.com/antlr4-go/antlr/v4"
)

func TestParseDmlFullSyntaxExample(t *testing.T) {
	t.Parallel()

	filename := "doc/example-all-syntax-1.4.dml"
	code, err := os.ReadFile(
		"../../grammar/dml/examples/example-all-syntax-1.4.dml",
	)
	if err != nil {
		t.Fatalf("cannot read file - %v", err)
	}

	ParseDmlText(t, string(code), filename, 0)
}

func TestParseDmlVersion_1_4(t *testing.T) {
	t.Parallel()

	code := "dml 1.4;\ndevice dev;\n"
	errorListener := dmlex.NewErrorListener("test.dml", 0)
	input := antlr.NewInputStream(code)
	file := dmlex.ParseStream(input, errorListener)
	_ = file.Parser.Dml()

	if len(errorListener.Errors) > 0 {
		t.Fatalf("unexpected errors for DML 1.4: %v", errorListener.Errors)
	}
}

func TestParseDmlVersion_1_2_ThrowsEUNIMPL(t *testing.T) {
	t.Parallel()

	// Contains DML 1.2 syntax that would fail in DML 1.4 parser if not aborted immediately.
	code := "dml 1.2;\nparameter foo = 1;\ncall $bar();\n"
	errorListener := dmlex.NewErrorListener("test.dml", 0)
	input := antlr.NewInputStream(code)
	file := dmlex.ParseStream(input, errorListener)
	_ = file.Parser.Dml()

	if len(errorListener.Errors) != 1 {
		t.Fatalf(
			"expected exactly 1 error, got %d: %v",
			len(errorListener.Errors),
			errorListener.Errors,
		)
	}

	errInfo := errorListener.Errors[0]
	if !errors.Is(errInfo.Err, dmlex.EUNIMPL) {
		t.Errorf("expected error to match EUNIMPL, got: %v", errInfo.Err)
	}
	expectedMsg := "Error Unimplemented: dml 1.2 syntax is not supported yet"
	if errInfo.Msg != expectedMsg {
		t.Errorf("expected msg %q, got %q", expectedMsg, errInfo.Msg)
	}
	if errInfo.Symbol != "1.2" {
		t.Errorf("expected symbol '1.2', got %q", errInfo.Symbol)
	}
	if !errorListener.HasError(dmlex.EUNIMPL) {
		t.Errorf("expected HasError(EUNIMPL) to be true")
	}
}

func TestParseDmlVersion_OtherUnsupported_ThrowsEUNIMPL(t *testing.T) {
	t.Parallel()

	for _, ver := range []string{"1.0", "1.3", "2.0"} {
		t.Run("version_"+ver, func(t *testing.T) {
			code := "dml " + ver + ";\ndevice dev;\n"
			errorListener := dmlex.NewErrorListener("test.dml", 0)
			input := antlr.NewInputStream(code)
			file := dmlex.ParseStream(input, errorListener)
			_ = file.Parser.Dml()

			if len(errorListener.Errors) != 1 {
				t.Fatalf(
					"expected 1 error for version %s, got %d",
					ver,
					len(errorListener.Errors),
				)
			}
			errInfo := errorListener.Errors[0]
			if !errors.Is(errInfo.Err, dmlex.EUNIMPL) {
				t.Errorf(
					"expected EUNIMPL for version %s, got %v",
					ver,
					errInfo.Err,
				)
			}
			if !strings.Contains(
				errInfo.Msg,
				"dml "+ver+" syntax is not supported yet",
			) {
				t.Errorf(
					"expected error message to mention version %s, got %q",
					ver,
					errInfo.Msg,
				)
			}
		})
	}
}

func TestParseDmlVersion_Omitted(t *testing.T) {
	t.Parallel()

	code := "device dev;\n"
	errorListener := dmlex.NewErrorListener("test.dml", 0)
	input := antlr.NewInputStream(code)
	file := dmlex.ParseStream(input, errorListener)
	_ = file.Parser.Dml()

	if len(errorListener.Errors) > 0 {
		t.Fatalf(
			"unexpected error for omitted dmlTag: %v",
			errorListener.Errors,
		)
	}
}

func TestParsePlaceholderNames_Issue76(t *testing.T) {
	t.Parallel()

	code := "dml 1.4;\ndevice dev;\ntypedef layout \"little-endian\" {\n" +
		"    uint8 _;\n" +
		"    uint8 data;\n" +
		"    uint16 __;\n" +
		"    uint32 ___;\n" +
		"} foo_t;\n" +
		"method test() {\n" +
		"    local int __ = 123_456;\n" +
		"    local int _1 = 0x_eace_f9b6;\n" +
		"}\n"
	errorListener := dmlex.NewErrorListener("test.dml", 0)
	input := antlr.NewInputStream(code)
	file := dmlex.ParseStream(input, errorListener)
	_ = file.Parser.Dml()

	if len(errorListener.Errors) > 0 {
		t.Fatalf(
			"unexpected syntax errors for placeholder/underscore identifiers: %v",
			errorListener.Errors,
		)
	}
}

func TestParseUnion_Issue77(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		code string
	}{
		{
			name: "01-extern-typedef-union-issue-77",
			code: "dml 1.4;\ndevice dev;\n" +
				"extern typedef union {\n" +
				"    struct {\n" +
				"        uint32 a   :8;\n" +
				"        uint32 b   :8;\n" +
				"        uint32 ca  :4;\n" +
				"        uint32 cb  :4;\n" +
				"        uint32 d   :8;\n" +
				"    } field;\n" +
				"    uint32 value;\n" +
				"} foo_t;\n",
		},
		{
			name: "02-empty-union",
			code: "dml 1.4;\ndevice dev;\n" +
				"typedef union {} my_union_t;\n" +
				"extern typedef union {} extern_u_t;\n",
		},
		{
			name: "03-tagged-union-def-and-ref",
			code: "dml 1.4;\ndevice dev;\n" +
				"typedef union my_tag {\n" +
				"    int a;\n" +
				"    float b;\n" +
				"} my_union_t;\n" +
				"extern union my_tag extern_var;\n",
		},
		{
			name: "04-struct-bitfields",
			code: "dml 1.4;\ndevice dev;\n" +
				"typedef struct {\n" +
				"    uint32 a : 8;\n" +
				"    uint32 : 8;\n" +
				"    uint32 : 0;\n" +
				"    uint32 b : 16;\n" +
				"    uint32 c : (4 + 4);\n" +
				"} bitfield_struct_t;\n",
		},
		{
			name: "05-anonymous-nested-struct-and-union",
			code: "dml 1.4;\ndevice dev;\n" +
				"typedef union {\n" +
				"    struct {\n" +
				"        uint32 a : 8;\n" +
				"        uint32 b : 8;\n" +
				"    };\n" +
				"    uint32 value;\n" +
				"} anon_struct_in_union_t;\n" +
				"typedef struct {\n" +
				"    int id;\n" +
				"    union {\n" +
				"        int i;\n" +
				"        float f;\n" +
				"    };\n" +
				"} anon_union_in_struct_t;\n",
		},
		{
			name: "06-union-pointers-and-arrays",
			code: "dml 1.4;\ndevice dev;\n" +
				"typedef union {\n" +
				"    int *ptr;\n" +
				"    int arr[4];\n" +
				"    void (*fn)(int);\n" +
				"} complex_u_t;\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			errorListener := dmlex.NewErrorListener("test.dml", 0)
			input := antlr.NewInputStream(tc.code)
			file := dmlex.ParseStream(input, errorListener)
			_ = file.Parser.Dml()

			if len(errorListener.Errors) > 0 {
				t.Fatalf(
					"unexpected syntax errors (%d): %v",
					len(errorListener.Errors),
					errorListener.Errors,
				)
			}

			if strings.Contains(tc.code, "union") {
				tokens := file.Parser.GetTokenStream()
				hasUnion := false
				for i := range tokens.Size() {
					tok := tokens.Get(i)
					if dmlex.GetTokenType(tok) == dmlex.TokenUnion {
						hasUnion = true
						break
					}
				}
				if !hasUnion {
					t.Fatalf(
						"expected TokenUnion in token stream for %s",
						tc.name,
					)
				}
			}
		})
	}
}
