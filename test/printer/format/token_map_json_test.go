// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package format_test

import (
	"encoding/json"
	"testing"

	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer/format"
)

func TestTokenMap_JSON_MarshalUnmarshal(t *testing.T) {
	t.Run("EmptyMap", func(t *testing.T) {
		m := format.NewTokenMap(nil)
		data, err := json.Marshal(m)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}
		if string(data) != "{}" {
			t.Errorf("Marshal(empty) = %s, want {}", string(data))
		}

		decoded, err := format.UnmarshalTokenMap(data)
		if err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}
		if flag := decoded.GetFormatFlag(
			dmlex.TokenSemi,
		); flag != format.FormatFlag(
			0,
		) {
			t.Errorf("GetFormatFlag(TokenSemi) = %v, want 0", flag)
		}
	})

	t.Run("RoundTrip", func(t *testing.T) {
		flagSemi := format.FormatFlag(0).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorAfter(format.SepMandatoryWrap).
			WithBreakAfterPriority(format.PriorityHigh)
		flagComma := format.FormatFlag(0).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorAfter(format.SepSpaceMaybeWrap).
			WithIndentIncrease()

		original := format.NewTokenMap(
			map[dmlex.TokenType]format.FormatFlag{
				dmlex.TokenSemi:  flagSemi,
				dmlex.TokenComma: flagComma,
			},
		)

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}

		decoded, err := format.UnmarshalTokenMap(data)
		if err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}

		if got := decoded.GetFormatFlag(dmlex.TokenSemi); got != flagSemi {
			t.Errorf("TokenSemi flag = %v, want %v", got, flagSemi)
		}
		if got := decoded.GetFormatFlag(dmlex.TokenComma); got != flagComma {
			t.Errorf("TokenComma flag = %v, want %v", got, flagComma)
		}
	})

	t.Run("UnmarshalAliases", func(t *testing.T) {
		rawJSON := `{
			"semi": {
				"before": {"token-prio": "normal", "sep": "none", "break-prio": "default"},
				"after": {"token-prio": "normal", "sep": "wrap-mandatory", "break-prio": "high"},
				"indent": "none",
				"force-wrap": false
			},
			"token-comma": {
				"before": {"token-prio": "normal", "sep": "none", "break-prio": "default"},
				"after": {"token-prio": "normal", "sep": "wrap", "break-prio": "high"},
				"indent": "increase",
				"force-wrap": false
			}
		}`

		decoded, err := format.UnmarshalTokenMap([]byte(rawJSON))
		if err != nil {
			t.Fatalf("UnmarshalTokenMap error: %v", err)
		}

		semiFlag := decoded.GetFormatFlag(dmlex.TokenSemi)
		if semiFlag.GetSeparatorAfter() != format.SepMandatoryWrap {
			t.Errorf(
				"TokenSemi SepAfter = %v, want %v",
				semiFlag.GetSeparatorAfter(),
				format.SepMandatoryWrap,
			)
		}

		commaFlag := decoded.GetFormatFlag(dmlex.TokenComma)
		if !commaFlag.HasIndentIncrease() {
			t.Errorf("TokenComma HasIndentIncrease = false, want true")
		}
		if commaFlag.GetSeparatorAfter() != format.SepSpaceMaybeWrap {
			t.Errorf(
				"TokenComma SepAfter = %v, want %v",
				commaFlag.GetSeparatorAfter(),
				format.SepSpaceMaybeWrap,
			)
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		invalidCases := []string{
			`{"UnknownTokenType123": {}}`,
			`{"semi": "not-an-object"}`,
			`not-json`,
		}

		for _, in := range invalidCases {
			if _, err := format.UnmarshalTokenMap(
				[]byte(in),
			); err == nil {
				t.Errorf(
					"UnmarshalTokenMap(%s) expected error, got nil",
					in,
				)
			}
		}
	})
}
