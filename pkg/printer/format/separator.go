// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package format

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Separator represents mutually exclusive whitespace and line break states between tokens.
type Separator uint8

const (
	// SepSpace indicates a single space between adjacent tokens (0b000).
	SepSpace Separator = iota
	// SepNone indicates no whitespace between tokens (0b001, e.g. clinging tokens like parens or dots).
	SepNone
	// SepSpaceMaybeWrap emits a space on a single line, or wraps with an indent if linecap is exceeded (0b010).
	SepSpaceMaybeWrap
	// SepNoSpaceMaybeWrap emits no space on a single line, or wraps if linecap is exceeded (0b011).
	SepNoSpaceMaybeWrap
	// SepMandatoryWrap indicates a forced line break between tokens (0b100).
	SepMandatoryWrap
	// SepWrapButNoSpace indicates a forced line break without leading or trailing space (0b101).
	SepWrapButNoSpace
)

func (s Separator) String() string {
	switch s {
	case SepSpace:
		return "SepSpace"
	case SepNone:
		return "SepNone"
	case SepSpaceMaybeWrap:
		return "SepSpaceMaybeWrap"
	case SepNoSpaceMaybeWrap:
		return "SepNoSpaceMaybeWrap"
	case SepMandatoryWrap:
		return "SepMandatoryWrap"
	case SepWrapButNoSpace:
		return "SepWrapButNoSpace"
	default:
		return fmt.Sprintf("Separator(%d)", s)
	}
}

// JSONString returns the lowercase kebab-case JSON string representation of the Separator.
func (s Separator) JSONString() string {
	switch s {
	case SepSpace:
		return "space"
	case SepNone:
		return "none"
	case SepSpaceMaybeWrap:
		return "wrap"
	case SepNoSpaceMaybeWrap:
		return "wrap-no-space"
	case SepMandatoryWrap:
		return "wrap-mandatory"
	case SepWrapButNoSpace:
		return "wrap-mandatory-no-space"
	default:
		return fmt.Sprintf("separator(%d)", s)
	}
}

// ParseSeparator parses a string into a Separator, supporting standard names, aliases, and Go identifier forms.
func ParseSeparator(s string) (Separator, error) {
	norm := strings.ToLower(strings.TrimSpace(s))
	norm = strings.ReplaceAll(norm, "_", "-")
	switch norm {
	case "space", "sepspace":
		return SepSpace, nil
	case "none", "sepnone":
		return SepNone, nil
	case "wrap",
		"optional-wrap",
		"sepoptionalwrap",
		"sepspacemaybewrap",
		"space-maybe-wrap",
		"spacemaybewrap":
		return SepSpaceMaybeWrap, nil
	case "wrap-no-space",
		"wrapnospace",
		"optional-wrap-no-space",
		"sepoptionalwrapnospace",
		"sepnospacemaybewrap",
		"no-space-maybe-wrap",
		"nospacemaybewrap":
		return SepNoSpaceMaybeWrap, nil
	case "wrap-mandatory",
		"wrapmandatory",
		"mandatory-wrap",
		"sepmandatorywrap":
		return SepMandatoryWrap, nil
	case "wrap-mandatory-no-space",
		"wrapmandatorynospace",
		"mandatory-wrap-no-space",
		"sepmandatorywrapnospace",
		"sepwrapbutnospace",
		"wrapbutnospace",
		"wrap-but-no-space":
		return SepWrapButNoSpace, nil
	default:
		return SepSpace, fmt.Errorf("unknown separator: %q", s)
	}
}

func (s Separator) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.JSONString())
}

func (s *Separator) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return fmt.Errorf("cannot unmarshal null into Separator")
	}
	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		if num < int(SepSpace) || num > int(SepWrapButNoSpace) {
			return fmt.Errorf("invalid separator number: %d", num)
		}
		*s = Separator(num)
		return nil
	}
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	parsed, err := ParseSeparator(str)
	if err != nil {
		return err
	}
	*s = parsed
	return nil
}
