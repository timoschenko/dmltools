// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IIdentContext is an interface to support dynamic dispatch.
type IIdentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ATTRIBUTE() antlr.TerminalNode
	ASYNC() antlr.TerminalNode
	AUTO() antlr.TerminalNode
	AWAIT() antlr.TerminalNode
	BANK() antlr.TerminalNode
	BITORDER() antlr.TerminalNode
	CALL() antlr.TerminalNode
	CLASS() antlr.TerminalNode
	CONNECT() antlr.TerminalNode
	CONSTANT() antlr.TerminalNode
	CRITICAL() antlr.TerminalNode
	DATA() antlr.TerminalNode
	DEVICE() antlr.TerminalNode
	ENDIAN() antlr.TerminalNode
	ENUM() antlr.TerminalNode
	EVENT() antlr.TerminalNode
	FIELD() antlr.TerminalNode
	FOOTER() antlr.TerminalNode
	GROUP() antlr.TerminalNode
	HEADER() antlr.TerminalNode
	HEADER2() antlr.TerminalNode
	IMPLEMENT() antlr.TerminalNode
	IMPORT() antlr.TerminalNode
	INDEPENDENT() antlr.TerminalNode
	INFO() antlr.TerminalNode
	INTERFACE() antlr.TerminalNode
	LOGGROUP() antlr.TerminalNode
	MEMOIZED() antlr.TerminalNode
	METHOD() antlr.TerminalNode
	NAMESPACE() antlr.TerminalNode
	NOTHROW() antlr.TerminalNode
	PARAM() antlr.TerminalNode
	PORT() antlr.TerminalNode
	PRIVATE() antlr.TerminalNode
	PROTECTED() antlr.TerminalNode
	PROVISIONAL() antlr.TerminalNode
	PUBLIC() antlr.TerminalNode
	RESTRICT() antlr.TerminalNode
	SAVED() antlr.TerminalNode
	SELECT() antlr.TerminalNode
	SIZE() antlr.TerminalNode
	SPEC_VIOL() antlr.TerminalNode
	STARTUP() antlr.TerminalNode
	STATIC() antlr.TerminalNode
	SUBDEVICE() antlr.TerminalNode
	THEN() antlr.TerminalNode
	THROWS() antlr.TerminalNode
	UNIMPL() antlr.TerminalNode
	UNION() antlr.TerminalNode
	USING() antlr.TerminalNode
	VIRTUAL() antlr.TerminalNode
	VOLATILE() antlr.TerminalNode
	WITH() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode

	// IsIdentContext differentiates from other interfaces.
	IsIdentContext()
}

type IdentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*IdentContext) IsIdentContext() {}

func (s *IdentContext) GetParser() antlr.Parser { return s.parser }

func (s *IdentContext) ATTRIBUTE() antlr.TerminalNode {
	return s.GetToken(dmlParserATTRIBUTE, 0)
}

func (s *IdentContext) ASYNC() antlr.TerminalNode {
	return s.GetToken(dmlParserASYNC, 0)
}

func (s *IdentContext) AUTO() antlr.TerminalNode {
	return s.GetToken(dmlParserAUTO, 0)
}

func (s *IdentContext) AWAIT() antlr.TerminalNode {
	return s.GetToken(dmlParserAWAIT, 0)
}

func (s *IdentContext) BANK() antlr.TerminalNode {
	return s.GetToken(dmlParserBANK, 0)
}

func (s *IdentContext) BITORDER() antlr.TerminalNode {
	return s.GetToken(dmlParserBITORDER, 0)
}

func (s *IdentContext) CALL() antlr.TerminalNode {
	return s.GetToken(dmlParserCALL, 0)
}

func (s *IdentContext) CLASS() antlr.TerminalNode {
	return s.GetToken(dmlParserCLASS, 0)
}

func (s *IdentContext) CONNECT() antlr.TerminalNode {
	return s.GetToken(dmlParserCONNECT, 0)
}

func (s *IdentContext) CONSTANT() antlr.TerminalNode {
	return s.GetToken(dmlParserCONSTANT, 0)
}

func (s *IdentContext) CRITICAL() antlr.TerminalNode {
	return s.GetToken(dmlParserCRITICAL, 0)
}

func (s *IdentContext) DATA() antlr.TerminalNode {
	return s.GetToken(dmlParserDATA, 0)
}

func (s *IdentContext) DEVICE() antlr.TerminalNode {
	return s.GetToken(dmlParserDEVICE, 0)
}

func (s *IdentContext) ENDIAN() antlr.TerminalNode {
	return s.GetToken(dmlParserENDIAN, 0)
}

func (s *IdentContext) ENUM() antlr.TerminalNode {
	return s.GetToken(dmlParserENUM, 0)
}

func (s *IdentContext) EVENT() antlr.TerminalNode {
	return s.GetToken(dmlParserEVENT, 0)
}

func (s *IdentContext) FIELD() antlr.TerminalNode {
	return s.GetToken(dmlParserFIELD, 0)
}

func (s *IdentContext) FOOTER() antlr.TerminalNode {
	return s.GetToken(dmlParserFOOTER, 0)
}

func (s *IdentContext) GROUP() antlr.TerminalNode {
	return s.GetToken(dmlParserGROUP, 0)
}

func (s *IdentContext) HEADER() antlr.TerminalNode {
	return s.GetToken(dmlParserHEADER, 0)
}

func (s *IdentContext) HEADER2() antlr.TerminalNode {
	return s.GetToken(dmlParserHEADER2, 0)
}

func (s *IdentContext) IMPLEMENT() antlr.TerminalNode {
	return s.GetToken(dmlParserIMPLEMENT, 0)
}

func (s *IdentContext) IMPORT() antlr.TerminalNode {
	return s.GetToken(dmlParserIMPORT, 0)
}

func (s *IdentContext) INDEPENDENT() antlr.TerminalNode {
	return s.GetToken(dmlParserINDEPENDENT, 0)
}

func (s *IdentContext) INFO() antlr.TerminalNode {
	return s.GetToken(dmlParserINFO, 0)
}

func (s *IdentContext) INTERFACE() antlr.TerminalNode {
	return s.GetToken(dmlParserINTERFACE, 0)
}

func (s *IdentContext) LOGGROUP() antlr.TerminalNode {
	return s.GetToken(dmlParserLOGGROUP, 0)
}

func (s *IdentContext) MEMOIZED() antlr.TerminalNode {
	return s.GetToken(dmlParserMEMOIZED, 0)
}

func (s *IdentContext) METHOD() antlr.TerminalNode {
	return s.GetToken(dmlParserMETHOD, 0)
}

func (s *IdentContext) NAMESPACE() antlr.TerminalNode {
	return s.GetToken(dmlParserNAMESPACE, 0)
}

func (s *IdentContext) NOTHROW() antlr.TerminalNode {
	return s.GetToken(dmlParserNOTHROW, 0)
}

func (s *IdentContext) PARAM() antlr.TerminalNode {
	return s.GetToken(dmlParserPARAM, 0)
}

func (s *IdentContext) PORT() antlr.TerminalNode {
	return s.GetToken(dmlParserPORT, 0)
}

func (s *IdentContext) PRIVATE() antlr.TerminalNode {
	return s.GetToken(dmlParserPRIVATE, 0)
}

func (s *IdentContext) PROTECTED() antlr.TerminalNode {
	return s.GetToken(dmlParserPROTECTED, 0)
}

func (s *IdentContext) PROVISIONAL() antlr.TerminalNode {
	return s.GetToken(dmlParserPROVISIONAL, 0)
}

func (s *IdentContext) PUBLIC() antlr.TerminalNode {
	return s.GetToken(dmlParserPUBLIC, 0)
}

func (s *IdentContext) RESTRICT() antlr.TerminalNode {
	return s.GetToken(dmlParserRESTRICT, 0)
}

func (s *IdentContext) SAVED() antlr.TerminalNode {
	return s.GetToken(dmlParserSAVED, 0)
}

func (s *IdentContext) SELECT() antlr.TerminalNode {
	return s.GetToken(dmlParserSELECT, 0)
}

func (s *IdentContext) SIZE() antlr.TerminalNode {
	return s.GetToken(dmlParserSIZE, 0)
}

func (s *IdentContext) SPEC_VIOL() antlr.TerminalNode {
	return s.GetToken(dmlParserSPEC_VIOL, 0)
}

func (s *IdentContext) STARTUP() antlr.TerminalNode {
	return s.GetToken(dmlParserSTARTUP, 0)
}

func (s *IdentContext) STATIC() antlr.TerminalNode {
	return s.GetToken(dmlParserSTATIC, 0)
}

func (s *IdentContext) SUBDEVICE() antlr.TerminalNode {
	return s.GetToken(dmlParserSUBDEVICE, 0)
}

func (s *IdentContext) THEN() antlr.TerminalNode {
	return s.GetToken(dmlParserTHEN, 0)
}

func (s *IdentContext) THROWS() antlr.TerminalNode {
	return s.GetToken(dmlParserTHROWS, 0)
}

func (s *IdentContext) UNIMPL() antlr.TerminalNode {
	return s.GetToken(dmlParserUNIMPL, 0)
}

func (s *IdentContext) UNION() antlr.TerminalNode {
	return s.GetToken(dmlParserUNION, 0)
}

func (s *IdentContext) USING() antlr.TerminalNode {
	return s.GetToken(dmlParserUSING, 0)
}

func (s *IdentContext) VIRTUAL() antlr.TerminalNode {
	return s.GetToken(dmlParserVIRTUAL, 0)
}

func (s *IdentContext) VOLATILE() antlr.TerminalNode {
	return s.GetToken(dmlParserVOLATILE, 0)
}

func (s *IdentContext) WITH() antlr.TerminalNode {
	return s.GetToken(dmlParserWITH, 0)
}

func (s *IdentContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(dmlParserIDENTIFIER, 0)
}

func (s *IdentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IdentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterIdent(s)
	}
}

func (s *IdentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitIdent(s)
	}
}

func NewEmptyIdentContext() *IdentContext {
	var p = new(IdentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_ident
	return p
}

func InitEmptyIdentContext(p *IdentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_ident
}

func NewIdentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentContext {
	var p = new(IdentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_ident

	return p
}
