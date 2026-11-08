// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ICCodeBlockContext is an interface to support dynamic dispatch.
type ICCodeBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	C_CODE_BLOCK() antlr.TerminalNode
	HEADER() antlr.TerminalNode
	HEADER2() antlr.TerminalNode
	FOOTER() antlr.TerminalNode

	// IsCCodeBlockContext differentiates from other interfaces.
	IsCCodeBlockContext()
}

type CCodeBlockContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*CCodeBlockContext) IsCCodeBlockContext() {}

func (s *CCodeBlockContext) GetParser() antlr.Parser { return s.parser }

func (s *CCodeBlockContext) C_CODE_BLOCK() antlr.TerminalNode {
	return s.GetToken(dmlParserC_CODE_BLOCK, 0)
}

func (s *CCodeBlockContext) HEADER() antlr.TerminalNode {
	return s.GetToken(dmlParserHEADER, 0)
}

func (s *CCodeBlockContext) HEADER2() antlr.TerminalNode {
	return s.GetToken(dmlParserHEADER2, 0)
}

func (s *CCodeBlockContext) FOOTER() antlr.TerminalNode {
	return s.GetToken(dmlParserFOOTER, 0)
}

func (s *CCodeBlockContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CCodeBlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CCodeBlockContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterCCodeBlock(s)
	}
}

func (s *CCodeBlockContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitCCodeBlock(s)
	}
}

func NewEmptyCCodeBlockContext() *CCodeBlockContext {
	var p = new(CCodeBlockContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_cCodeBlock
	return p
}

func InitEmptyCCodeBlockContext(p *CCodeBlockContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_cCodeBlock
}

func NewCCodeBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CCodeBlockContext {
	var p = new(CCodeBlockContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_cCodeBlock

	return p
}
