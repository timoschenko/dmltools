// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package format

import (
	"encoding/json"
	"fmt"

	"github.com/timoschenko/dmltools/pkg/dmlex"
)

// TokenMap provides format flag lookups for token types.
type TokenMap interface {
	GetFormatFlag(tt dmlex.TokenType) FormatFlag
}

// TokenMapImpl is a concrete map of token types to format flags.
type TokenMapImpl map[dmlex.TokenType]FormatFlag

// GetFormatFlag returns the FormatFlag for a given TokenType.
func (m TokenMapImpl) GetFormatFlag(tt dmlex.TokenType) FormatFlag {
	return m[tt]
}

// MarshalJSON serializes TokenMapImpl into a JSON object with token names as keys.
func (m TokenMapImpl) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("{}"), nil
	}
	raw := make(map[string]FormatFlag, len(m))
	for k, v := range m {
		raw[k.String()] = v
	}
	return json.Marshal(raw)
}

// UnmarshalJSON deserializes a JSON object into a TokenMapImpl.
func (m *TokenMapImpl) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*m = make(TokenMapImpl)
		return nil
	}
	var raw map[string]FormatFlag
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	res := make(TokenMapImpl, len(raw))
	for k, v := range raw {
		tt, err := dmlex.ParseTokenType(k)
		if err != nil {
			var num int
			if _, errInt := fmt.Sscanf(k, "%d", &num); errInt == nil {
				candidate := dmlex.TokenType(num)
				if candidate >= dmlex.TokenEOF &&
					candidate <= dmlex.TokenLastTokenType {
					res[candidate] = v
					continue
				}
			}
			return fmt.Errorf(
				"invalid token type key in TokenMap: %w",
				err,
			)
		}
		res[tt] = v
	}
	*m = res
	return nil
}

// NewTokenMap creates a new TokenMap from a map of token types to format flags.
func NewTokenMap(m map[dmlex.TokenType]FormatFlag) TokenMap {
	if m == nil {
		return TokenMapImpl{}
	}
	return TokenMapImpl(m)
}

// UnmarshalTokenMap unmarshals JSON data into a TokenMap interface.
func UnmarshalTokenMap(data []byte) (TokenMap, error) {
	var m TokenMapImpl
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// EmptyTokenMap is a TokenMap with no overrides.
var EmptyTokenMap TokenMap = TokenMapImpl{}
