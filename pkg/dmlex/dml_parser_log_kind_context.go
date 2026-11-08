// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ILogKindContext is an interface to support dynamic dispatch.
type ILogKindContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INFO() antlr.TerminalNode
	ERROR() antlr.TerminalNode
	SPEC_VIOL() antlr.TerminalNode
	WARNING() antlr.TerminalNode
	UNIMPL() antlr.TerminalNode
	CRITICAL() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	Utf8_sconst() IUtf8_sconstContext

	// IsLogKindContext differentiates from other interfaces.
	IsLogKindContext()
}

type LogKindContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*LogKindContext) IsLogKindContext() {}

func (s *LogKindContext) GetParser() antlr.Parser { return s.parser }

func (s *LogKindContext) INFO() antlr.TerminalNode {
	return s.GetToken(dmlParserINFO, 0)
}

func (s *LogKindContext) ERROR() antlr.TerminalNode {
	return s.GetToken(dmlParserERROR, 0)
}

func (s *LogKindContext) SPEC_VIOL() antlr.TerminalNode {
	return s.GetToken(dmlParserSPEC_VIOL, 0)
}

func (s *LogKindContext) WARNING() antlr.TerminalNode {
	return s.GetToken(dmlParserWARNING, 0)
}

func (s *LogKindContext) UNIMPL() antlr.TerminalNode {
	return s.GetToken(dmlParserUNIMPL, 0)
}

func (s *LogKindContext) CRITICAL() antlr.TerminalNode {
	return s.GetToken(dmlParserCRITICAL, 0)
}

func (s *LogKindContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(dmlParserIDENTIFIER, 0)
}

func (s *LogKindContext) Utf8_sconst() IUtf8_sconstContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUtf8_sconstContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUtf8_sconstContext)
}

func (s *LogKindContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogKindContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LogKindContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterLogKind(s)
	}
}

func (s *LogKindContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitLogKind(s)
	}
}

func NewEmptyLogKindContext() *LogKindContext {
	var p = new(LogKindContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_logKind
	return p
}

func InitEmptyLogKindContext(p *LogKindContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_logKind
}

func NewLogKindContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LogKindContext {
	var p = new(LogKindContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_logKind

	return p
}
