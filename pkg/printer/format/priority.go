// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package format

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Priority represents a 4-bit priority level (0-15) governing both line break
// candidate selection and token pair conflict resolution.
//
// Scale reference:
//   - Priority 14-15 (Critical/Max): Structural delimiters and mandatory boundaries.
//   - Priority 10-12 (High/Highest): Comma separators, outer list elements, statement split points.
//   - Priority 8-9 (AboveNormal/MediumHigh): Major binary operators (assignment, logical &&, ||).
//   - Priority 6-7 (Normal/Default): Comparison and relational operators (==, !=, <, >).
//   - Priority 4-5 (BelowNormal/MediumLow): Additive arithmetic operators (+, -).
//   - Priority 2-3 (Low/VeryLow): Multiplicative operators (*, /, %).
//   - Priority 1 (Lowest): Clinging operators and member access (., ->).
type Priority uint8

const (
	PriorityBelowLowest  Priority = iota // 0
	PriorityLowest                       // 1
	PriorityVeryLow                      // 2
	PriorityLow                          // 3
	PriorityMediumLow                    // 4
	PriorityBelowNormal                  // 5
	PriorityNormal                       // 6
	PriorityDefault                      // 7
	PriorityAboveNormal                  // 8
	PriorityMediumHigh                   // 9
	PriorityHigh                         // 10
	PriorityVeryHigh                     // 11
	PriorityHighest                      // 12
	PriorityAboveHighest                 // 13
	PriorityCritical                     // 14
	PriorityMax                          // 15
)

var priorityJSONNames = [...]string{
	"below-lowest", "lowest", "very-low", "low",
	"medium-low", "below-normal", "normal", "default",
	"above-normal", "medium-high", "high", "very-high",
	"highest", "above-highest", "critical", "max",
}

// String makes the priority pretty-print nicely.
func (p Priority) String() string {
	names := [...]string{
		"BelowLowest", "Lowest", "VeryLow", "Low",
		"MediumLow", "BelowNormal", "Normal", "Default",
		"AboveNormal", "MediumHigh", "High", "VeryHigh",
		"Highest", "AboveHighest", "Critical", "Max",
	}

	if p > PriorityMax {
		return fmt.Sprintf("Priority(%d)", p)
	}
	return names[p]
}

// JSONString returns the lowercase kebab-case JSON string representation of the Priority.
func (p Priority) JSONString() string {
	if p > PriorityMax {
		return fmt.Sprintf("priority(%d)", p)
	}
	return priorityJSONNames[p]
}

// ParsePriority parses a string into a Priority, supporting case-insensitive names, kebab-case, snake_case, and Go identifier forms.
func ParsePriority(s string) (Priority, error) {
	norm := strings.ToLower(strings.TrimSpace(s))
	norm = strings.TrimPrefix(norm, "priority")
	norm = strings.ReplaceAll(norm, "_", "-")
	switch norm {
	case "below-lowest", "belowlowest":
		return PriorityBelowLowest, nil
	case "lowest":
		return PriorityLowest, nil
	case "very-low", "verylow":
		return PriorityVeryLow, nil
	case "low":
		return PriorityLow, nil
	case "medium-low", "mediumlow":
		return PriorityMediumLow, nil
	case "below-normal", "belownormal":
		return PriorityBelowNormal, nil
	case "normal":
		return PriorityNormal, nil
	case "default":
		return PriorityDefault, nil
	case "above-normal", "abovenormal":
		return PriorityAboveNormal, nil
	case "medium-high", "mediumhigh":
		return PriorityMediumHigh, nil
	case "high":
		return PriorityHigh, nil
	case "very-high", "veryhigh":
		return PriorityVeryHigh, nil
	case "highest":
		return PriorityHighest, nil
	case "above-highest", "abovehighest":
		return PriorityAboveHighest, nil
	case "critical":
		return PriorityCritical, nil
	case "max":
		return PriorityMax, nil
	default:
		return PriorityDefault, fmt.Errorf("unknown priority: %q", s)
	}
}

func (p Priority) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.JSONString())
}

func (p *Priority) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return fmt.Errorf("cannot unmarshal null into Priority")
	}
	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		if num < int(PriorityBelowLowest) || num > int(PriorityMax) {
			return fmt.Errorf("invalid priority number: %d", num)
		}
		*p = Priority(num)
		return nil
	}
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	parsed, err := ParsePriority(str)
	if err != nil {
		return err
	}
	*p = parsed
	return nil
}
