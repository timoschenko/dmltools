// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test_dmlex

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/timoschenko/dmltools/pkg/dmlex"
)

func TestTokenType_JSON_Marshal(t *testing.T) {
	cases := []struct {
		token dmlex.TokenType
		want  string
	}{
		{dmlex.TokenEOF, `"TokenEOF"`},
		{dmlex.TokenUnknown, `"TokenUnknown"`},
		{dmlex.TokenIdentifier, `"TokenIdentifier"`},
		{dmlex.TokenSemi, `"TokenSemi"`},
		{dmlex.TokenCommentInline, `"TokenCommentInline"`},
		{dmlex.TokenEndian, `"TokenEndian"`},
	}

	for _, tc := range cases {
		data, err := json.Marshal(tc.token)
		if err != nil {
			t.Fatalf("Marshal(%v) error: %v", tc.token, err)
		}
		if string(data) != tc.want {
			t.Errorf("Marshal(%v) = %s, want %s", tc.token, data, tc.want)
		}
	}
}

func TestTokenType_JSON_UnmarshalStringsAndAliases(t *testing.T) {
	cases := []struct {
		input string
		want  dmlex.TokenType
	}{
		{`"TokenIdentifier"`, dmlex.TokenIdentifier},
		{`"Identifier"`, dmlex.TokenIdentifier},
		{`"identifier"`, dmlex.TokenIdentifier},
		{`"IDENTIFIER"`, dmlex.TokenIdentifier},
		{`"token-identifier"`, dmlex.TokenIdentifier},
		{`"token_identifier"`, dmlex.TokenIdentifier},
		{`"TokenEOF"`, dmlex.TokenEOF},
		{`"EOF"`, dmlex.TokenEOF},
		{`"eof"`, dmlex.TokenEOF},
		{`"TokenUnknown"`, dmlex.TokenUnknown},
		{`"unknown"`, dmlex.TokenUnknown},
		{`"TokenEndian"`, dmlex.TokenEndian},
		{`"endian"`, dmlex.TokenEndian},
		{`"ENDIAN"`, dmlex.TokenEndian},
		{`"TokenCommentInline"`, dmlex.TokenCommentInline},
		{`"comment-inline"`, dmlex.TokenCommentInline},
		{`"TokenSemi"`, dmlex.TokenSemi},
		{`"semi"`, dmlex.TokenSemi},
		{`"TokenAssign"`, dmlex.TokenAssign},
		{`"assign"`, dmlex.TokenAssign},
	}

	for _, tc := range cases {
		var got dmlex.TokenType
		if err := json.Unmarshal([]byte(tc.input), &got); err != nil {
			t.Fatalf("Unmarshal(%s) error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("Unmarshal(%s) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestTokenType_JSON_UnmarshalNumbers(t *testing.T) {
	for i := int(dmlex.TokenEOF); i <= int(dmlex.TokenLastTokenType); i++ {
		input := []byte(strconv.Itoa(i))
		var got dmlex.TokenType
		if err := json.Unmarshal(input, &got); err != nil {
			t.Fatalf("Unmarshal(%s) error: %v", input, err)
		}
		if int(got) != i {
			t.Errorf("Unmarshal(%s) = %v, want %d", input, got, i)
		}
	}
}

func TestTokenType_JSON_RoundTrip(t *testing.T) {
	for tt := dmlex.TokenEOF; tt <= dmlex.TokenLastTokenType; tt++ {
		data, err := json.Marshal(tt)
		if err != nil {
			t.Fatalf("Marshal(%v) error: %v", tt, err)
		}

		var decoded dmlex.TokenType
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Unmarshal(%s) error: %v", data, err)
		}

		if decoded != tt {
			t.Errorf(
				"Roundtrip mismatch: got %v (%d), want %v (%d)",
				decoded,
				decoded,
				tt,
				tt,
			)
		}
	}
}

func TestTokenType_JSON_Invalid(t *testing.T) {
	invalid := []string{
		`"NonExistentTokenType"`,
		`"random_string"`,
		`2000`,
		`-2`,
		`null`,
		`true`,
		`[]`,
	}

	for _, in := range invalid {
		var got dmlex.TokenType
		if err := json.Unmarshal([]byte(in), &got); err == nil {
			t.Errorf("Unmarshal(%s) expected error, got %v", in, got)
		}
	}
}
