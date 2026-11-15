// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/fs/stream"
)

func parseDML(code string) (*dmlex.File, []string) {
	lines := strings.Split(code, "\n")
	buf := stream.NewHashedBufferFromString("test.dml", code)
	file := dmlex.ParseStream(buf)
	return file, lines
}

func TestKindsAndNames(t *testing.T) {
	cases := []struct {
		kind   Kind
		letter byte
		name   string
	}{
		{KindDevice, 'd', "device"},
		{KindTemplate, 't', "template"},
		{KindMethod, 'm', "method"},
		{KindParam, 'p', "param"},
		{KindBank, 'b', "bank"},
		{KindRegister, 'r', "register"},
		{KindField, 'f', "field"},
		{KindGroup, 'g', "group"},
		{KindPort, 'P', "port"},
		{KindAttribute, 'a', "attribute"},
		{KindConnect, 'c', "connect"},
		{KindEvent, 'e', "event"},
		{KindImplement, 'i', "implement"},
		{KindInterface, 'I', "interface"},
		{KindSubdevice, 'S', "subdevice"},
		{KindVariable, 'v', "variable"},
		{KindConstant, 'C', "constant"},
		{KindTypedef, 'T', "typedef"},
		{KindStruct, 's', "struct"},
		{KindMember, 'M', "member"},
		{KindHook, 'h', "hook"},
		{KindLogGroup, 'l', "loggroup"},
		{KindExtern, 'x', "extern"},
	}

	for _, tc := range cases {
		if tc.kind.Letter() != tc.letter {
			t.Errorf(
				"Kind %v: Letter() = %c, want %c",
				tc.kind,
				tc.kind.Letter(),
				tc.letter,
			)
		}
		if tc.kind.Name() != tc.name {
			t.Errorf(
				"Kind %v: Name() = %q, want %q",
				tc.kind,
				tc.kind.Name(),
				tc.name,
			)
		}
	}
}

func TestEscapePattern(t *testing.T) {
	input := `    method foo(int a / b) { return a \ b; }` + "\r\n"
	expected := `    method foo(int a \/ b) { return a \\ b; }`
	got := EscapePattern(input)
	if got != expected {
		t.Errorf("EscapePattern() = %q, want %q", got, expected)
	}
}

func TestFormatTag(t *testing.T) {
	tag := Tag{
		Name:       "my_reg",
		File:       "dev.dml",
		Line:       42,
		SourceLine: "    register my_reg @ 0x10 is (read_only);",
		Kind:       KindRegister,
		ScopeKind:  KindBank,
		ScopePath:  "regs",
		Signature:  "",
		Inherits:   []string{"read_only"},
	}

	linePattern := FormatTag(tag, ExCmdPattern)
	expectedPattern := "my_reg\tdev.dml\t/^    register my_reg @ 0x10 is (read_only);$/;\"\tr\tline:42\tbank:regs\tinherits:read_only"
	if linePattern != expectedPattern {
		t.Errorf(
			"FormatTag(ExCmdPattern) = %q, want %q",
			linePattern,
			expectedPattern,
		)
	}

	lineNumber := FormatTag(tag, ExCmdNumber)
	expectedNumber := "my_reg\tdev.dml\t42;\"\tr\tline:42\tbank:regs\tinherits:read_only"
	if lineNumber != expectedNumber {
		t.Errorf(
			"FormatTag(ExCmdNumber) = %q, want %q",
			lineNumber,
			expectedNumber,
		)
	}

	lineCombine := FormatTag(tag, ExCmdCombine)
	expectedCombine := "my_reg\tdev.dml\t42;/^    register my_reg @ 0x10 is (read_only);$/;\"\tr\tline:42\tbank:regs\tinherits:read_only"
	if lineCombine != expectedCombine {
		t.Errorf(
			"FormatTag(ExCmdCombine) = %q, want %q",
			lineCombine,
			expectedCombine,
		)
	}
}

func TestSortTagsAndLines(t *testing.T) {
	lines := []string{
		"zebra\tfile.dml\t10;\"\tm",
		"Apple\tfile.dml\t1;\"\tr",
		"apple\tfile.dml\t5;\"\tf",
		"Banana\tfile.dml\t20;\"\tp",
	}

	// ASCII sort (SortYes)
	SortLines(lines, SortYes)
	expectedYes := []string{
		"Apple\tfile.dml\t1;\"\tr",
		"Banana\tfile.dml\t20;\"\tp",
		"apple\tfile.dml\t5;\"\tf",
		"zebra\tfile.dml\t10;\"\tm",
	}
	for i, l := range lines {
		if l != expectedYes[i] {
			t.Errorf("SortYes[%d] = %q, want %q", i, l, expectedYes[i])
		}
	}

	// Foldcase sort (SortFoldCase)
	SortLines(lines, SortFoldCase)
	expectedFold := []string{
		"Apple\tfile.dml\t1;\"\tr",
		"apple\tfile.dml\t5;\"\tf",
		"Banana\tfile.dml\t20;\"\tp",
		"zebra\tfile.dml\t10;\"\tm",
	}
	for i, l := range lines {
		if l != expectedFold[i] {
			t.Errorf("SortFoldCase[%d] = %q, want %q", i, l, expectedFold[i])
		}
	}
}

func TestWriteTagsAndExisting(t *testing.T) {
	existingContent := "!_TAG_FILE_FORMAT\t2\t/comment/\n" +
		"!_TAG_FILE_SORTED\t1\t/comment/\n" +
		"old_tag\told.dml\t1;\"\tr\n"

	existing, err := ParseExistingTags(
		strings.NewReader(existingContent),
	)
	if err != nil {
		t.Fatalf("ParseExistingTags error: %v", err)
	}
	if len(existing) != 1 || existing[0] != "old_tag\told.dml\t1;\"\tr" {
		t.Fatalf("ParseExistingTags parsed: %#v", existing)
	}

	newTags := []Tag{
		{
			Name:       "new_tag",
			File:       "new.dml",
			Line:       10,
			SourceLine: "method new_tag() {}",
			Kind:       KindMethod,
		},
	}

	var buf bytes.Buffer
	err = WriteTags(&buf, newTags, WriteOptions{
		Version:  "1.0.0",
		Author:   "Author",
		Sort:     SortYes,
		ExCmd:    ExCmdNumber,
		Append:   true,
		Existing: existing,
	})
	if err != nil {
		t.Fatalf("WriteTags error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "!_TAG_FILE_FORMAT\t2") {
		t.Errorf("Expected pseudo header in output:\n%s", out)
	}
	if !strings.Contains(out, "new_tag\tnew.dml\t10;\"\tm\tline:10") {
		t.Errorf("Expected new tag in output:\n%s", out)
	}
	if !strings.Contains(out, "old_tag\told.dml\t1;\"\tr") {
		t.Errorf("Expected existing tag in output:\n%s", out)
	}
}

func TestExtractFileTags_Comprehensive(t *testing.T) {
	code := `dml 1.4;
device test_dev;

is (base_device);

constant MY_CONST = 100;
loggroup LG_ERR;
hook() my_hook;

typedef uint32 reg_t;

typedef struct {
    int id;
    bool active;
} info_s;

template my_tpl {
    is (read_only);
    param tpl_param = 1;
    method tpl_func(int x) -> (bool ok) {}
}

bank regs {
    register r0 @ 0x0 {
        field f0 @ [0];
        method read() -> (uint64 val) {
            return 0;
        }
    }
}

port my_port;
attribute my_attr;
connect my_conn;
event my_ev;
group my_grp;
implement my_impl;
interface my_iface;
subdevice my_subdev;

session int sess_var;
saved uint32 saved_var;

extern int ext_fn();
extern typedef uint32 ext_t;
extern typedef struct { int a; } ext_s;

in each (my_tpl) {
    param ineach_p = 2;
    method ineach_m() {}
}

#if (MY_CONST == 100) {
param if_param = 10;
} #else {
param else_param = 20;
}
`

	file, lines := parseDML(code)
	tree := file.Parser.Dml()

	tags := ExtractFileTags(
		tree,
		"test.dml",
		lines,
		ExtractOptions{
			EmitQualified: false,
		},
	)

	tagMap := make(map[string]Tag)
	for _, tag := range tags {
		tagMap[tag.Name] = tag
	}

	expected := []struct {
		name      string
		kind      Kind
		scopeKind Kind
		scopePath string
		hasSig    bool
		inherits  []string
	}{
		{
			name:      "test_dev",
			kind:      KindDevice,
			scopePath: "",
			inherits:  []string{"base_device"},
		},
		{
			name:      "MY_CONST",
			kind:      KindConstant,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
		{
			name:      "LG_ERR",
			kind:      KindLogGroup,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
		{
			name:      "my_hook",
			kind:      KindHook,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
		{
			name:      "reg_t",
			kind:      KindTypedef,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
		{
			name:      "info_s",
			kind:      KindStruct,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
		{
			name:      "id",
			kind:      KindMember,
			scopeKind: KindStruct,
			scopePath: "test_dev.info_s",
		},
		{
			name:      "active",
			kind:      KindMember,
			scopeKind: KindStruct,
			scopePath: "test_dev.info_s",
		},
		{
			name:      "my_tpl",
			kind:      KindTemplate,
			scopeKind: KindDevice,
			scopePath: "test_dev",
			inherits:  []string{"read_only"},
		},
		{
			name:      "tpl_param",
			kind:      KindParam,
			scopeKind: KindTemplate,
			scopePath: "test_dev.my_tpl",
		},
		{
			name:      "tpl_func",
			kind:      KindMethod,
			scopeKind: KindTemplate,
			scopePath: "test_dev.my_tpl",
			hasSig:    true,
		},
		{
			name:      "regs",
			kind:      KindBank,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
		{
			name:      "r0",
			kind:      KindRegister,
			scopeKind: KindBank,
			scopePath: "test_dev.regs",
		},
		{
			name:      "f0",
			kind:      KindField,
			scopeKind: KindRegister,
			scopePath: "test_dev.regs.r0",
		},
		{
			name:      "read",
			kind:      KindMethod,
			scopeKind: KindRegister,
			scopePath: "test_dev.regs.r0",
			hasSig:    true,
		},
		{
			name:      "my_port",
			kind:      KindPort,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
		{
			name:      "my_attr",
			kind:      KindAttribute,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
		{
			name:      "my_conn",
			kind:      KindConnect,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
		{
			name:      "my_ev",
			kind:      KindEvent,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
		{
			name:      "my_grp",
			kind:      KindGroup,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
		{
			name:      "my_impl",
			kind:      KindImplement,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
		{
			name:      "my_iface",
			kind:      KindInterface,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
		{
			name:      "my_subdev",
			kind:      KindSubdevice,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
		{
			name:      "sess_var",
			kind:      KindVariable,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
		{
			name:      "saved_var",
			kind:      KindVariable,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
		{
			name:      "ext_fn",
			kind:      KindExtern,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
		{
			name:      "ext_t",
			kind:      KindTypedef,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
		{
			name:      "ext_s",
			kind:      KindStruct,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
		{
			name:      "ineach_p",
			kind:      KindParam,
			scopeKind: KindDevice,
			scopePath: "test_dev",
			inherits:  []string{"my_tpl"},
		},
		{
			name:      "ineach_m",
			kind:      KindMethod,
			scopeKind: KindDevice,
			scopePath: "test_dev",
			inherits:  []string{"my_tpl"},
		},
		{
			name:      "if_param",
			kind:      KindParam,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
		{
			name:      "else_param",
			kind:      KindParam,
			scopeKind: KindDevice,
			scopePath: "test_dev",
		},
	}

	for _, exp := range expected {
		tag, ok := tagMap[exp.name]
		if !ok {
			t.Errorf("Missing expected tag %q", exp.name)
			continue
		}
		if tag.Kind != exp.kind {
			t.Errorf(
				"Tag %q: kind = %c (%s), want %c (%s)",
				exp.name,
				tag.Kind,
				tag.Kind.Name(),
				exp.kind,
				exp.kind.Name(),
			)
		}
		if tag.ScopePath != exp.scopePath {
			t.Errorf(
				"Tag %q: scopePath = %q, want %q",
				exp.name,
				tag.ScopePath,
				exp.scopePath,
			)
		}
		if exp.scopeKind != 0 && tag.ScopeKind != exp.scopeKind {
			t.Errorf(
				"Tag %q: scopeKind = %c (%s), want %c (%s)",
				exp.name,
				tag.ScopeKind,
				tag.ScopeKind.Name(),
				exp.scopeKind,
				exp.scopeKind.Name(),
			)
		}
		if exp.hasSig && tag.Signature == "" {
			t.Errorf("Tag %q: expected non-empty signature", exp.name)
		}
		if len(tag.Inherits) != len(exp.inherits) {
			t.Errorf(
				"Tag %q: inherits = %v, want %v",
				exp.name,
				tag.Inherits,
				exp.inherits,
			)
		} else {
			for i, inh := range exp.inherits {
				if tag.Inherits[i] != inh {
					t.Errorf(
						"Tag %q: inherits[%d] = %q, want %q",
						exp.name,
						i,
						tag.Inherits[i],
						inh,
					)
				}
			}
		}
	}
}

func TestExtractFileTags_Qualified(t *testing.T) {
	code := `dml 1.4;
device dev;
bank regs {
    register r0 {
        field f0;
    }
}
`
	file, lines := parseDML(code)
	tree := file.Parser.Dml()

	tags := ExtractFileTags(
		tree,
		"test.dml",
		lines,
		ExtractOptions{
			EmitQualified: true,
		},
	)

	names := make(map[string]bool)
	for _, tag := range tags {
		names[tag.Name] = true
	}

	expectedQualified := []string{
		"dev",
		"regs",
		"dev.regs",
		"r0",
		"dev.regs.r0",
		"f0",
		"dev.regs.r0.f0",
	}

	for _, name := range expectedQualified {
		if !names[name] {
			t.Errorf(
				"Expected tag name %q was not emitted in qualified mode",
				name,
			)
		}
	}
}
