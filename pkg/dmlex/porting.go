// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package dmlex

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"unicode"

	"github.com/antlr4-go/antlr/v4"
)

// NOTE: make sure to have `go install golang.org/x/tools/cmd/stringer@latest`

//go:generate stringer -type=TokenType
type TokenType int

const (
	TokenUnknown          TokenType = iota
	TokenEOF              TokenType = antlr.TokenEOF
	TokenT__0             TokenType = dmlLexerT__0
	TokenAfter            TokenType = dmlLexerAFTER
	TokenAs               TokenType = dmlLexerAS
	TokenAssert           TokenType = dmlLexerASSERT
	TokenAsync            TokenType = dmlLexerASYNC
	TokenAttribute        TokenType = dmlLexerATTRIBUTE
	TokenAuto             TokenType = dmlLexerAUTO
	TokenAwait            TokenType = dmlLexerAWAIT
	TokenBank             TokenType = dmlLexerBANK
	TokenBitfields        TokenType = dmlLexerBITFIELDS
	TokenBitorder         TokenType = dmlLexerBITORDER
	TokenBreak            TokenType = dmlLexerBREAK
	TokenCall             TokenType = dmlLexerCALL
	TokenCase             TokenType = dmlLexerCASE
	TokenCast             TokenType = dmlLexerCAST
	TokenCatch            TokenType = dmlLexerCATCH
	TokenChar             TokenType = dmlLexerCHAR
	TokenClass            TokenType = dmlLexerCLASS
	TokenConnect          TokenType = dmlLexerCONNECT
	TokenConst            TokenType = dmlLexerCONST
	TokenConstant         TokenType = dmlLexerCONSTANT
	TokenContinue         TokenType = dmlLexerCONTINUE
	TokenCritical         TokenType = dmlLexerCRITICAL
	TokenData             TokenType = dmlLexerDATA
	TokenDefault          TokenType = dmlLexerDEFAULT
	TokenDefaultNew       TokenType = dmlLexerDEFAULT_NEW
	TokenDefined          TokenType = dmlLexerDEFINED
	TokenDelete           TokenType = dmlLexerDELETE
	TokenDevice           TokenType = dmlLexerDEVICE
	TokenDml              TokenType = dmlLexerDML
	TokenDo               TokenType = dmlLexerDO
	TokenDouble           TokenType = dmlLexerDOUBLE
	TokenEach             TokenType = dmlLexerEACH
	TokenElse             TokenType = dmlLexerELSE
	TokenEndian           TokenType = dmlLexerENDIAN
	TokenEnum             TokenType = dmlLexerENUM
	TokenError            TokenType = dmlLexerERROR
	TokenEvent            TokenType = dmlLexerEVENT
	TokenExport           TokenType = dmlLexerEXPORT
	TokenExtern           TokenType = dmlLexerEXTERN
	TokenField            TokenType = dmlLexerFIELD
	TokenFloat            TokenType = dmlLexerFLOAT
	TokenFooter           TokenType = dmlLexerFOOTER
	TokenFor              TokenType = dmlLexerFOR
	TokenForeach          TokenType = dmlLexerFOREACH
	TokenGoto             TokenType = dmlLexerGOTO
	TokenGroup            TokenType = dmlLexerGROUP
	TokenHashElse         TokenType = dmlLexerHASH_ELSE
	TokenHashForeach      TokenType = dmlLexerHASH_FOREACH
	TokenHashIf           TokenType = dmlLexerHASH_IF
	TokenHashSelect       TokenType = dmlLexerHASH_SELECT
	TokenHeader           TokenType = dmlLexerHEADER
	TokenHeader2          TokenType = dmlLexerHEADER2
	TokenHook             TokenType = dmlLexerHOOK
	TokenIf               TokenType = dmlLexerIF
	TokenImplement        TokenType = dmlLexerIMPLEMENT
	TokenImport           TokenType = dmlLexerIMPORT
	TokenIn               TokenType = dmlLexerIN
	TokenIndependent      TokenType = dmlLexerINDEPENDENT
	TokenInfo             TokenType = dmlLexerINFO
	TokenInline           TokenType = dmlLexerINLINE
	TokenInt              TokenType = dmlLexerINT
	TokenInterface        TokenType = dmlLexerINTERFACE
	TokenIs               TokenType = dmlLexerIS
	TokenLayout           TokenType = dmlLexerLAYOUT
	TokenLocal            TokenType = dmlLexerLOCAL
	TokenLog              TokenType = dmlLexerLOG
	TokenLoggroup         TokenType = dmlLexerLOGGROUP
	TokenLong             TokenType = dmlLexerLONG
	TokenMemoized         TokenType = dmlLexerMEMOIZED
	TokenMethod           TokenType = dmlLexerMETHOD
	TokenNamespace        TokenType = dmlLexerNAMESPACE
	TokenNew              TokenType = dmlLexerNEW
	TokenNothrow          TokenType = dmlLexerNOTHROW
	TokenParam            TokenType = dmlLexerPARAM
	TokenPort             TokenType = dmlLexerPORT
	TokenPrivate          TokenType = dmlLexerPRIVATE
	TokenProtected        TokenType = dmlLexerPROTECTED
	TokenProvisional      TokenType = dmlLexerPROVISIONAL
	TokenPublic           TokenType = dmlLexerPUBLIC
	TokenRegister         TokenType = dmlLexerREGISTER
	TokenRestrict         TokenType = dmlLexerRESTRICT
	TokenReturn           TokenType = dmlLexerRETURN
	TokenSaved            TokenType = dmlLexerSAVED
	TokenSelect           TokenType = dmlLexerSELECT
	TokenSequence         TokenType = dmlLexerSEQUENCE
	TokenSession          TokenType = dmlLexerSESSION
	TokenShared           TokenType = dmlLexerSHARED
	TokenShort            TokenType = dmlLexerSHORT
	TokenSigned           TokenType = dmlLexerSIGNED
	TokenSize             TokenType = dmlLexerSIZE
	TokenSizeof           TokenType = dmlLexerSIZEOF
	TokenSizeoftype       TokenType = dmlLexerSIZEOFTYPE
	TokenSpecViol         TokenType = dmlLexerSPEC_VIOL
	TokenStartup          TokenType = dmlLexerSTARTUP
	TokenStatic           TokenType = dmlLexerSTATIC
	TokenStringify        TokenType = dmlLexerSTRINGIFY
	TokenStruct           TokenType = dmlLexerSTRUCT
	TokenSubdevice        TokenType = dmlLexerSUBDEVICE
	TokenSwitch           TokenType = dmlLexerSWITCH
	TokenTemplate         TokenType = dmlLexerTEMPLATE
	TokenThen             TokenType = dmlLexerTHEN
	TokenThis             TokenType = dmlLexerTHIS
	TokenThrow            TokenType = dmlLexerTHROW
	TokenThrows           TokenType = dmlLexerTHROWS
	TokenTry              TokenType = dmlLexerTRY
	TokenTypedef          TokenType = dmlLexerTYPEDEF
	TokenTypeof           TokenType = dmlLexerTYPEOF
	TokenUndefined        TokenType = dmlLexerUNDEFINED
	TokenUnimpl           TokenType = dmlLexerUNIMPL
	TokenUnion            TokenType = dmlLexerUNION
	TokenUnsigned         TokenType = dmlLexerUNSIGNED
	TokenUsing            TokenType = dmlLexerUSING
	TokenVect             TokenType = dmlLexerVECT
	TokenVirtual          TokenType = dmlLexerVIRTUAL
	TokenVoid             TokenType = dmlLexerVOID
	TokenVolatile         TokenType = dmlLexerVOLATILE
	TokenWarning          TokenType = dmlLexerWARNING
	TokenWhere            TokenType = dmlLexerWHERE
	TokenWhile            TokenType = dmlLexerWHILE
	TokenWith             TokenType = dmlLexerWITH
	TokenIntegerLiteral   TokenType = dmlLexerINTEGER_LITERAL
	TokenHexLiteral       TokenType = dmlLexerHEX_LITERAL
	TokenBinaryLiteral    TokenType = dmlLexerBINARY_LITERAL
	TokenFloatLiteral     TokenType = dmlLexerFLOAT_LITERAL
	TokenCharLiteral      TokenType = dmlLexerCHAR_LITERAL
	TokenStringLiteral    TokenType = dmlLexerSTRING_LITERAL
	TokenCCodeBlock       TokenType = dmlLexerC_CODE_BLOCK
	TokenIdentifier       TokenType = dmlLexerIDENTIFIER
	TokenLparen           TokenType = dmlLexerLPAREN
	TokenRparen           TokenType = dmlLexerRPAREN
	TokenLbrace           TokenType = dmlLexerLBRACE
	TokenLbraceNew        TokenType = dmlLexerLBRACE_NEW
	TokenRbrace           TokenType = dmlLexerRBRACE
	TokenLbrack           TokenType = dmlLexerLBRACK
	TokenRbrack           TokenType = dmlLexerRBRACK
	TokenSemi             TokenType = dmlLexerSEMI
	TokenComma            TokenType = dmlLexerCOMMA
	TokenDot              TokenType = dmlLexerDOT
	TokenEllipsis         TokenType = dmlLexerELLIPSIS
	TokenAssign           TokenType = dmlLexerASSIGN
	TokenNewAssign        TokenType = dmlLexerNEW_ASSIGN
	TokenGt               TokenType = dmlLexerGT
	TokenLt               TokenType = dmlLexerLT
	TokenBang             TokenType = dmlLexerBANG
	TokenTilde            TokenType = dmlLexerTILDE
	TokenQuestion         TokenType = dmlLexerQUESTION
	TokenColon            TokenType = dmlLexerCOLON
	TokenEqual            TokenType = dmlLexerEQUAL
	TokenLe               TokenType = dmlLexerLE
	TokenGe               TokenType = dmlLexerGE
	TokenNotequal         TokenType = dmlLexerNOTEQUAL
	TokenAnd              TokenType = dmlLexerAND
	TokenOr               TokenType = dmlLexerOR
	TokenInc              TokenType = dmlLexerINC
	TokenDec              TokenType = dmlLexerDEC
	TokenAdd              TokenType = dmlLexerADD
	TokenSub              TokenType = dmlLexerSUB
	TokenMul              TokenType = dmlLexerMUL
	TokenDiv              TokenType = dmlLexerDIV
	TokenBitand           TokenType = dmlLexerBITAND
	TokenBitor            TokenType = dmlLexerBITOR
	TokenCaret            TokenType = dmlLexerCARET
	TokenMod              TokenType = dmlLexerMOD
	TokenAddAssign        TokenType = dmlLexerADD_ASSIGN
	TokenSubAssign        TokenType = dmlLexerSUB_ASSIGN
	TokenMulAssign        TokenType = dmlLexerMUL_ASSIGN
	TokenDivAssign        TokenType = dmlLexerDIV_ASSIGN
	TokenModAssign        TokenType = dmlLexerMOD_ASSIGN
	TokenAndAssign        TokenType = dmlLexerAND_ASSIGN
	TokenOrAssign         TokenType = dmlLexerOR_ASSIGN
	TokenXorAssign        TokenType = dmlLexerXOR_ASSIGN
	TokenLshiftAssign     TokenType = dmlLexerLSHIFT_ASSIGN
	TokenRshiftAssign     TokenType = dmlLexerRSHIFT_ASSIGN
	TokenLshift           TokenType = dmlLexerLSHIFT
	TokenRshift           TokenType = dmlLexerRSHIFT
	TokenArrow            TokenType = dmlLexerARROW
	TokenAt               TokenType = dmlLexerAT
	TokenAltQuestion      TokenType = dmlLexerALT_QUESTION
	TokenAltColon         TokenType = dmlLexerALT_COLON
	TokenWs               TokenType = dmlLexerWS
	TokenComment          TokenType = dmlLexerCOMMENT
	TokenCommentBlockEol  TokenType = dmlLexerCOMMENT_BLOCK_EOL
	TokenCommentBlock     TokenType = dmlLexerCOMMENT_BLOCK
	TokenCommentInlineEol TokenType = dmlLexerCOMMENT_INLINE_EOL
	TokenCommentInline    TokenType = dmlLexerCOMMENT_INLINE
	TokenLastTokenType
)

var (
	tokenTypeLookupOnce sync.Once
	tokenTypeLookup     map[string]TokenType
)

func camelToKebab(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 && s[i-1] != '-' && s[i-1] != '_' &&
				!unicode.IsUpper(rune(s[i-1])) {
				result.WriteRune('-')
			} else if i > 0 && unicode.IsUpper(rune(s[i-1])) && i+1 < len(s) &&
				unicode.IsLower(rune(s[i+1])) {
				result.WriteRune('-')
			}
			result.WriteRune(unicode.ToLower(r))
		} else if r == '_' {
			result.WriteRune('-')
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func initTokenTypeLookup() {
	tokenTypeLookup = make(map[string]TokenType, 2000)
	for tt := TokenEOF; tt <= TokenLastTokenType; tt++ {
		name := tt.String()
		if strings.HasPrefix(name, "TokenType(") {
			continue
		}

		addVariants := func(s string) {
			if s == "" {
				return
			}
			tokenTypeLookup[s] = tt
			tokenTypeLookup[strings.ToLower(s)] = tt
			tokenTypeLookup[strings.ToUpper(s)] = tt
			kebab := camelToKebab(s)
			tokenTypeLookup[kebab] = tt
			tokenTypeLookup[strings.ToLower(kebab)] = tt
			snake := strings.ReplaceAll(kebab, "-", "_")
			tokenTypeLookup[snake] = tt
			tokenTypeLookup[strings.ToUpper(snake)] = tt
		}

		addVariants(name)

		// Without "Token" prefix
		stripped := strings.TrimPrefix(name, "Token")
		if stripped != name {
			addVariants(stripped)
		}
	}
}

// ParseTokenType parses a string into a TokenType, supporting exact names ("TokenIdentifier"),
// case-insensitive variants ("identifier", "IDENTIFIER"), kebab-case ("token-identifier", "comment-inline"),
// snake_case ("COMMENT_INLINE"), and names without the "Token" prefix ("eof", "semi", "assign").
func ParseTokenType(s string) (TokenType, error) {
	tokenTypeLookupOnce.Do(initTokenTypeLookup)
	norm := strings.TrimSpace(s)
	if tt, ok := tokenTypeLookup[norm]; ok {
		return tt, nil
	}
	normLower := strings.ToLower(norm)
	if tt, ok := tokenTypeLookup[normLower]; ok {
		return tt, nil
	}
	kebab := camelToKebab(norm)
	if tt, ok := tokenTypeLookup[kebab]; ok {
		return tt, nil
	}
	normNoPrefix := strings.TrimPrefix(
		strings.TrimPrefix(kebab, "token-"),
		"token",
	)
	if tt, ok := tokenTypeLookup[normNoPrefix]; ok {
		return tt, nil
	}
	return TokenUnknown, fmt.Errorf("unknown token type: %q", s)
}

func (t TokenType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t TokenType) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *TokenType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return fmt.Errorf("cannot unmarshal null into TokenType")
	}
	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		tt := TokenType(num)
		if tt < TokenEOF || tt > TokenLastTokenType {
			return fmt.Errorf("invalid token type number: %d", num)
		}
		*t = tt
		return nil
	}
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	parsed, err := ParseTokenType(str)
	if err != nil {
		return err
	}
	*t = parsed
	return nil
}

func (t *TokenType) UnmarshalText(text []byte) error {
	s := string(text)
	parsed, err := ParseTokenType(s)
	if err != nil {
		var num int
		if _, errInt := fmt.Sscanf(s, "%d", &num); errInt == nil {
			tt := TokenType(num)
			if tt >= TokenEOF && tt <= TokenLastTokenType {
				*t = tt
				return nil
			}
		}
		return err
	}
	*t = parsed
	return nil
}

func GetTokenType(token antlr.Token) TokenType {
	return TokenType(token.GetTokenType())
}

func (t TokenType) IsOperand() bool {
	switch t {
	case TokenIdentifier,
		TokenIntegerLiteral,
		TokenHexLiteral,
		TokenBinaryLiteral,
		TokenFloatLiteral,
		TokenCharLiteral,
		TokenStringLiteral,
		TokenRparen,
		TokenRbrack:
		return true
	default:
		return false
	}
}

// IsBlockComment returns true for any block comment variant.
func (t TokenType) IsBlockComment() bool {
	switch t {
	case TokenCommentInline,
		TokenCommentInlineEol,
		TokenCommentBlock,
		TokenCommentBlockEol:
		return true
	default:
		return false
	}
}

func (t TokenType) IsMultiLineBlockComment() bool {
	switch t {
	case TokenCommentBlock, TokenCommentBlockEol:
		return true
	default:
		return false
	}
}

func (t TokenType) IsComment() bool {
	return t == TokenComment || t.IsBlockComment()
}

func (t TokenType) IsKeyword() bool {
	return t >= TokenAfter && t <= TokenWhile
}

func (t TokenType) Int() int {
	return int(t)
}

// ISimpleRuleContextEndsSemi is an interface to support dynamic dispatch.
type ISimpleRuleContextEndsSemi interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	SEMI() antlr.TerminalNode
}
