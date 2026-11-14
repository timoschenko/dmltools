// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package format_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/timoschenko/dmltools/pkg/printer/format"
)

func TestPriority_JSON(t *testing.T) {
	t.Run("MarshalAll", func(t *testing.T) {
		expected := map[format.Priority]string{
			format.PriorityBelowLowest:  `"below-lowest"`,
			format.PriorityLowest:       `"lowest"`,
			format.PriorityVeryLow:      `"very-low"`,
			format.PriorityLow:          `"low"`,
			format.PriorityMediumLow:    `"medium-low"`,
			format.PriorityBelowNormal:  `"below-normal"`,
			format.PriorityNormal:       `"normal"`,
			format.PriorityDefault:      `"default"`,
			format.PriorityAboveNormal:  `"above-normal"`,
			format.PriorityMediumHigh:   `"medium-high"`,
			format.PriorityHigh:         `"high"`,
			format.PriorityVeryHigh:     `"very-high"`,
			format.PriorityHighest:      `"highest"`,
			format.PriorityAboveHighest: `"above-highest"`,
			format.PriorityCritical:     `"critical"`,
			format.PriorityMax:          `"max"`,
		}
		for p, want := range expected {
			data, err := json.Marshal(p)
			if err != nil {
				t.Fatalf("Marshal(%v) error: %v", p, err)
			}
			if string(data) != want {
				t.Errorf("Marshal(%v) = %s, want %s", p, data, want)
			}
		}
	})

	t.Run("UnmarshalStringsAndAliases", func(t *testing.T) {
		cases := []struct {
			input string
			want  format.Priority
		}{
			{`"normal"`, format.PriorityNormal},
			{`"Normal"`, format.PriorityNormal},
			{`"PriorityNormal"`, format.PriorityNormal},
			{`"below-lowest"`, format.PriorityBelowLowest},
			{`"below_lowest"`, format.PriorityBelowLowest},
			{`"belowlowest"`, format.PriorityBelowLowest},
			{`"critical"`, format.PriorityCritical},
			{`"max"`, format.PriorityMax},
		}
		for _, tc := range cases {
			var p format.Priority
			if err := json.Unmarshal([]byte(tc.input), &p); err != nil {
				t.Fatalf("Unmarshal(%s) error: %v", tc.input, err)
			}
			if p != tc.want {
				t.Errorf("Unmarshal(%s) = %v, want %v", tc.input, p, tc.want)
			}
		}
	})

	t.Run("UnmarshalNumbers", func(t *testing.T) {
		for i := range 16 {
			b, _ := json.Marshal(i)
			input := []byte(strings.TrimSpace(string(b)))
			var p format.Priority
			if err := json.Unmarshal(input, &p); err != nil {
				t.Fatalf("Unmarshal(%s) error: %v", input, err)
			}
			if int(p) != i {
				t.Errorf("Unmarshal(%s) = %v, want %d", input, p, i)
			}
		}
	})

	t.Run("UnmarshalInvalid", func(t *testing.T) {
		invalid := []string{`"nonexistent"`, `16`, `-1`, `true`, `null`}
		for _, in := range invalid {
			var p format.Priority
			if err := json.Unmarshal([]byte(in), &p); err == nil {
				t.Errorf("Unmarshal(%s) expected error, got %v", in, p)
			}
		}
	})
}
