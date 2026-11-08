// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package layout

import (
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"strings"
	"unicode/utf8"

	"github.com/antlr4-go/antlr/v4"
)

// WillTokensMalform returns true if printing prevToken and currToken adjacent
// without any intervening whitespace would cause the ANTLR lexer to produce
// an unintended compound token, comment start, or otherwise malform the token stream.
func WillTokensMalform(prevToken, currToken antlr.Token) bool {
	if prevToken == nil || currToken == nil {
		return false
	}

	prevType := dmlex.GetTokenType(prevToken)
	currType := dmlex.GetTokenType(currToken)
	if prevType == dmlex.TokenUnknown || currType == dmlex.TokenUnknown {
		return false
	}

	prevText := prevToken.GetText()
	currText := currToken.GetText()
	if len(prevText) == 0 || len(currText) == 0 {
		return false
	}

	rPrev, _ := utf8.DecodeLastRuneInString(prevText)
	rCurr, _ := utf8.DecodeRuneInString(currText)

	// 1. Identifiers / Keywords / Numbers continuation
	if isIdentOrDigit(rPrev) && isIdentOrDigit(rCurr) {
		return true
	}

	// 2. Division & Slash operators ('/') - exclude comment endings ('*/')
	if !prevType.IsComment() && rPrev == '/' {
		if rCurr == '*' || rCurr == '/' || rCurr == '=' {
			return true
		}
	}

	// 3. Multiplication & Pointer asterisk ('*')
	if !prevType.IsComment() && rPrev == '*' {
		if (!currType.IsComment() && rCurr == '/') || rCurr == '=' {
			return true
		}
	}

	// 4. Plus ('+')
	if rPrev == '+' {
		if rCurr == '+' || rCurr == '=' {
			return true
		}
	}

	// 5. Minus ('-')
	if rPrev == '-' {
		if rCurr == '-' || rCurr == '=' || rCurr == '>' {
			return true
		}
	}

	// 6. Colon (':')
	if rPrev == ':' {
		if rCurr == '{' || rCurr == '=' ||
			strings.HasPrefix(currText, "default") {
			return true
		}
	}

	// 7. Modulo ('%')
	if rPrev == '%' {
		if rCurr == '{' || rCurr == '=' {
			return true
		}
	}

	// 8. Bitwise AND ('&')
	if rPrev == '&' {
		if rCurr == '&' || rCurr == '=' {
			return true
		}
	}

	// 9. Bitwise OR ('|')
	if rPrev == '|' {
		if rCurr == '|' || rCurr == '=' {
			return true
		}
	}

	// 10. Bitwise XOR ('^')
	if rPrev == '^' && rCurr == '=' {
		return true
	}

	// 11. Logical NOT / BANG ('!')
	if rPrev == '!' && rCurr == '=' {
		return true
	}

	// 12. Equal ('=')
	if rPrev == '=' && rCurr == '=' {
		return true
	}

	// 13. Less-Than ('<')
	if rPrev == '<' {
		if rCurr == '<' || rCurr == '=' {
			return true
		}
	}

	// 14. Greater-Than ('>')
	if rPrev == '>' {
		if rCurr == '>' || rCurr == '=' {
			return true
		}
	}

	// 15. Dot ('.') and digits (Float literals & ellipsis)
	if (rPrev == '.' && isDigit(rCurr)) ||
		(isDigit(rPrev) && rCurr == '.') ||
		(rPrev == '.' && rCurr == '.') {
		return true
	}

	// 16. Hash ('#')
	if rPrev == '#' && (rCurr == '?' || rCurr == ':') {
		return true
	}

	return false
}

func isIdentOrDigit(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '_'
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
