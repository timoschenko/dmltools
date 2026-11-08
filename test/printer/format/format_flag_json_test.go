// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package format_test

import (
	"encoding/json"
	"testing"

	"github.com/timoschenko/dmltools/pkg/printer/format"
)

func TestFormatFlag_JSON_RoundTrip(t *testing.T) {
	flags := []format.FormatFlag{
		format.FormatFlag(0),
		format.FormatFlag(0).
			WithSeparatorBeforePriority(format.PriorityNormal).
			WithSeparatorBefore(format.SepNone).
			WithBreakBeforePriority(format.PriorityHigh).
			WithSeparatorAfterPriority(format.PriorityDefault).
			WithSeparatorAfter(format.SepMandatoryWrap).
			WithBreakAfterPriority(format.PriorityMediumLow).
			WithIndentIncrease().
			WithForceWrap(),
		format.FormatFlag(0).
			WithSeparatorBeforePriority(format.PriorityCritical).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithBreakBeforePriority(format.PriorityAboveHighest).
			WithSeparatorAfterPriority(format.PriorityMax).
			WithSeparatorAfter(format.SepWrapButNoSpace).
			WithBreakAfterPriority(format.PriorityVeryHigh).
			WithIndentDecrease(),
		format.FormatFlag(0).
			WithIndentIncrease().
			WithIndentDecrease(),
	}

	for _, original := range flags {
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal(%v) error: %v", original, err)
		}

		var decoded format.FormatFlag
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Unmarshal(%s) error: %v", data, err)
		}

		if decoded != original {
			t.Errorf(
				"Roundtrip mismatch:\ngot:  %s (%032b)\nwant: %s (%032b)\njson: %s",
				decoded,
				decoded,
				original,
				original,
				string(data),
			)
		}
	}
}

func TestFormatFlag_JSON_AliasesAndStructure(t *testing.T) {
	rawJSON := `{
		"before": {
			"token-prio": "normal",
			"sep": "none",
			"sep-prio": "high"
		},
		"after": {
			"token_priority": "default",
			"separator": "wrap-mandatory",
			"break_priority": "medium-high"
		},
		"ident": "increase",
		"force-wrap": true
	}`

	var f format.FormatFlag
	if err := json.Unmarshal([]byte(rawJSON), &f); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if f.GetTokenBeforePriority() != format.PriorityNormal {
		t.Errorf(
			"GetTokenBeforePriority() = %v, want %v",
			f.GetTokenBeforePriority(),
			format.PriorityNormal,
		)
	}
	if f.GetSeparatorBefore() != format.SepNone {
		t.Errorf(
			"GetSeparatorBefore() = %v, want %v",
			f.GetSeparatorBefore(),
			format.SepNone,
		)
	}
	if f.GetBreakBeforePriority() != format.PriorityHigh {
		t.Errorf(
			"GetBreakBeforePriority() = %v, want %v",
			f.GetBreakBeforePriority(),
			format.PriorityHigh,
		)
	}
	if f.GetTokenAfterPriority() != format.PriorityDefault {
		t.Errorf(
			"GetTokenAfterPriority() = %v, want %v",
			f.GetTokenAfterPriority(),
			format.PriorityDefault,
		)
	}
	if f.GetSeparatorAfter() != format.SepMandatoryWrap {
		t.Errorf(
			"GetSeparatorAfter() = %v, want %v",
			f.GetSeparatorAfter(),
			format.SepMandatoryWrap,
		)
	}
	if f.GetBreakAfterPriority() != format.PriorityMediumHigh {
		t.Errorf(
			"GetBreakAfterPriority() = %v, want %v",
			f.GetBreakAfterPriority(),
			format.PriorityMediumHigh,
		)
	}
	if !f.HasIndentIncrease() {
		t.Errorf("HasIndentIncrease() = false, want true")
	}
	if f.HasIndentDecrease() {
		t.Errorf("HasIndentDecrease() = true, want false")
	}
	if !f.IsForceWrap() {
		t.Errorf("IsForceWrap() = false, want true")
	}
}

func TestFormatFlag_JSON_IndentValues(t *testing.T) {
	cases := []struct {
		indent   string
		wantInc  bool
		wantDec  bool
		wantDesc string
	}{
		{`"none"`, false, false, "none"},
		{`""`, false, false, "none"},
		{`"increase"`, true, false, "increase"},
		{`"inc"`, true, false, "increase"},
		{`"decrease"`, false, true, "decrease"},
		{`"dec"`, false, true, "decrease"},
		{`"both"`, true, true, "both"},
	}

	for _, tc := range cases {
		jsonStr := `{"indent": ` + tc.indent + `}`
		var f format.FormatFlag
		if err := json.Unmarshal([]byte(jsonStr), &f); err != nil {
			t.Fatalf("Unmarshal(%s) error: %v", jsonStr, err)
		}
		if f.HasIndentIncrease() != tc.wantInc {
			t.Errorf(
				"Indent %s: HasIndentIncrease() = %v, want %v",
				tc.indent,
				f.HasIndentIncrease(),
				tc.wantInc,
			)
		}
		if f.HasIndentDecrease() != tc.wantDec {
			t.Errorf(
				"Indent %s: HasIndentDecrease() = %v, want %v",
				tc.indent,
				f.HasIndentDecrease(),
				tc.wantDec,
			)
		}
		if f.IndentString() != tc.wantDesc {
			t.Errorf(
				"Indent %s: IndentString() = %s, want %s",
				tc.indent,
				f.IndentString(),
				tc.wantDesc,
			)
		}
	}

	// Test invalid indent
	var f format.FormatFlag
	if err := json.Unmarshal([]byte(`{"indent": "invalid"}`), &f); err == nil {
		t.Errorf("Expected error for invalid indent, got nil")
	}
}

func TestFormatFlag_JSON_InvalidFields(t *testing.T) {
	invalid := []string{
		`{"unknown_field": 123}`,
		`{"before": {"unknown_side": 123}}`,
		`{"before": "not-an-object"}`,
		`{"force_wrap": "not-a-bool"}`,
	}
	for _, in := range invalid {
		var f format.FormatFlag
		if err := json.Unmarshal([]byte(in), &f); err == nil {
			t.Errorf("Unmarshal(%s) expected error, got nil", in)
		}
	}
}
