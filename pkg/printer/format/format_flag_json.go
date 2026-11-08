// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package format

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SideFlagJSON represents the grouped boundary properties (before or after) of a token.
type SideFlagJSON struct {
	TokenPriority Priority  `json:"token_priority"`
	Separator     Separator `json:"separator"`
	BreakPriority Priority  `json:"break_priority"`
}

// UnmarshalJSON unmarshals a SideFlagJSON, supporting aliases for keys.
func (s *SideFlagJSON) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	for k, v := range raw {
		norm := strings.ToLower(strings.ReplaceAll(k, "_", "-"))
		switch norm {
		case "token-priority", "token-prio", "tokenpriority", "tokenprio":
			if err := json.Unmarshal(v, &s.TokenPriority); err != nil {
				return fmt.Errorf("invalid token_priority: %w", err)
			}
		case "separator", "sep":
			if err := json.Unmarshal(v, &s.Separator); err != nil {
				return fmt.Errorf("invalid separator: %w", err)
			}
		case "break-priority", "break-prio", "breakpriority", "breakprio",
			"sep-priority", "sep-prio", "seppriority", "sepprio":
			if err := json.Unmarshal(v, &s.BreakPriority); err != nil {
				return fmt.Errorf("invalid break_priority: %w", err)
			}
		default:
			return fmt.Errorf("unknown side flag field: %q", k)
		}
	}
	return nil
}

type formatFlagJSON struct {
	Before    SideFlagJSON `json:"before"`
	After     SideFlagJSON `json:"after"`
	Indent    string       `json:"indent"`
	ForceWrap bool         `json:"force_wrap"`
}

// IndentString returns the string representation of the indent state ("none", "increase", "decrease", "both").
func (f FormatFlag) IndentString() string {
	inc := f.HasIndentIncrease()
	dec := f.HasIndentDecrease()
	if inc && dec {
		return "both"
	}
	if inc {
		return "increase"
	}
	if dec {
		return "decrease"
	}
	return "none"
}

// MarshalJSON serializes FormatFlag into verbose grouped JSON.
func (f FormatFlag) MarshalJSON() ([]byte, error) {
	raw := formatFlagJSON{
		Before: SideFlagJSON{
			TokenPriority: f.GetTokenBeforePriority(),
			Separator:     f.GetSeparatorBefore(),
			BreakPriority: f.GetBreakBeforePriority(),
		},
		After: SideFlagJSON{
			TokenPriority: f.GetTokenAfterPriority(),
			Separator:     f.GetSeparatorAfter(),
			BreakPriority: f.GetBreakAfterPriority(),
		},
		Indent:    f.IndentString(),
		ForceWrap: f.IsForceWrap(),
	}
	return json.Marshal(raw)
}

// UnmarshalJSON deserializes verbose grouped JSON into a FormatFlag bitfield.
func (f *FormatFlag) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var before SideFlagJSON
	var after SideFlagJSON
	var indentStr string
	var forceWrap bool

	for k, v := range raw {
		norm := strings.ToLower(strings.ReplaceAll(k, "_", "-"))
		switch norm {
		case "before":
			if err := json.Unmarshal(v, &before); err != nil {
				return fmt.Errorf("invalid 'before' field: %w", err)
			}
		case "after":
			if err := json.Unmarshal(v, &after); err != nil {
				return fmt.Errorf("invalid 'after' field: %w", err)
			}
		case "indent", "ident":
			if err := json.Unmarshal(v, &indentStr); err != nil {
				return fmt.Errorf("invalid 'indent' field: %w", err)
			}
		case "force-wrap", "forcewrap":
			if err := json.Unmarshal(v, &forceWrap); err != nil {
				return fmt.Errorf("invalid 'force_wrap' field: %w", err)
			}
		default:
			return fmt.Errorf("unknown FormatFlag field: %q", k)
		}
	}

	res := FormatFlag(0).
		WithSeparatorBeforePriority(before.TokenPriority).
		WithSeparatorBefore(before.Separator).
		WithBreakBeforePriority(before.BreakPriority).
		WithSeparatorAfterPriority(after.TokenPriority).
		WithSeparatorAfter(after.Separator).
		WithBreakAfterPriority(after.BreakPriority)

	normIndent := strings.ToLower(strings.TrimSpace(indentStr))
	switch normIndent {
	case "", "none":
		// no indent bits
	case "increase", "inc":
		res = res.WithIndentIncrease()
	case "decrease", "dec":
		res = res.WithIndentDecrease()
	case "both":
		res = res.WithIndentIncrease().WithIndentDecrease()
	default:
		return fmt.Errorf("unknown indent value: %q", indentStr)
	}

	if forceWrap {
		res = res.WithForceWrap()
	}

	*f = res
	return nil
}
