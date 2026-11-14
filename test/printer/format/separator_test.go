// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package format_test

import (
	"encoding/json"
	"testing"

	"github.com/timoschenko/dmltools/pkg/printer/format"
)

func TestSeparatorString(t *testing.T) {
	tests := []struct {
		s    format.Separator
		want string
	}{
		{format.SepSpace, "SepSpace"},
		{format.SepNone, "SepNone"},
		{format.SepSpaceMaybeWrap, "SepSpaceMaybeWrap"},
		{format.SepMandatoryWrap, "SepMandatoryWrap"},
		{format.Separator(10), "Separator(10)"},
	}
	for _, tt := range tests {
		if got := tt.s.String(); got != tt.want {
			t.Errorf("Separator.String() = %v, want %v", got, tt.want)
		}
	}
}

func TestSeparator_JSON(t *testing.T) {
	t.Run("MarshalAll", func(t *testing.T) {
		expected := map[format.Separator]string{
			format.SepSpace:            `"space"`,
			format.SepNone:             `"none"`,
			format.SepSpaceMaybeWrap:   `"wrap"`,
			format.SepNoSpaceMaybeWrap: `"wrap-no-space"`,
			format.SepMandatoryWrap:    `"wrap-mandatory"`,
			format.SepWrapButNoSpace:   `"wrap-mandatory-no-space"`,
		}
		for s, want := range expected {
			data, err := json.Marshal(s)
			if err != nil {
				t.Fatalf("Marshal(%v) error: %v", s, err)
			}
			if string(data) != want {
				t.Errorf("Marshal(%v) = %s, want %s", s, data, want)
			}
		}
	})

	t.Run("UnmarshalStringsAndAliases", func(t *testing.T) {
		cases := []struct {
			input string
			want  format.Separator
		}{
			{`"space"`, format.SepSpace},
			{`"SepSpace"`, format.SepSpace},
			{`"none"`, format.SepNone},
			{`"SepNone"`, format.SepNone},
			{`"wrap"`, format.SepSpaceMaybeWrap},
			{`"optional-wrap"`, format.SepSpaceMaybeWrap},
			{`"SepSpaceMaybeWrap"`, format.SepSpaceMaybeWrap},
			{`"wrap-no-space"`, format.SepNoSpaceMaybeWrap},
			{`"optional-wrap-no-space"`, format.SepNoSpaceMaybeWrap},
			{`"SepNoSpaceMaybeWrap"`, format.SepNoSpaceMaybeWrap},
			{`"wrap-mandatory"`, format.SepMandatoryWrap},
			{`"mandatory-wrap"`, format.SepMandatoryWrap},
			{`"SepMandatoryWrap"`, format.SepMandatoryWrap},
			{`"wrap-mandatory-no-space"`, format.SepWrapButNoSpace},
			{`"mandatory-wrap-no-space"`, format.SepWrapButNoSpace},
			{`"SepWrapButNoSpace"`, format.SepWrapButNoSpace},
		}
		for _, tc := range cases {
			var s format.Separator
			if err := json.Unmarshal([]byte(tc.input), &s); err != nil {
				t.Fatalf("Unmarshal(%s) error: %v", tc.input, err)
			}
			if s != tc.want {
				t.Errorf("Unmarshal(%s) = %v, want %v", tc.input, s, tc.want)
			}
		}
	})

	t.Run("UnmarshalNumbers", func(t *testing.T) {
		for i := range 6 {
			b, _ := json.Marshal(i)
			var s format.Separator
			if err := json.Unmarshal(b, &s); err != nil {
				t.Fatalf("Unmarshal(%d) error: %v", i, err)
			}
			if int(s) != i {
				t.Errorf("Unmarshal(%d) = %v, want %d", i, s, i)
			}
		}
	})

	t.Run("UnmarshalInvalid", func(t *testing.T) {
		invalid := []string{`"invalid-sep"`, `6`, `-1`, `true`, `null`}
		for _, in := range invalid {
			var s format.Separator
			if err := json.Unmarshal([]byte(in), &s); err == nil {
				t.Errorf("Unmarshal(%s) expected error, got %v", in, s)
			}
		}
	})
}
