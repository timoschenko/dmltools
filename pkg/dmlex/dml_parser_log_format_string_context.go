// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ILogFormatStringContext is an interface to support dynamic dispatch.
type ILogFormatStringContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Bracketed_string_literal() IBracketed_string_literalContext

	// IsLogFormatStringContext differentiates from other interfaces.
	IsLogFormatStringContext()
}

type LogFormatStringContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*LogFormatStringContext) IsLogFormatStringContext() {}

func (s *LogFormatStringContext) GetParser() antlr.Parser { return s.parser }

func (s *LogFormatStringContext) Bracketed_string_literal() IBracketed_string_literalContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBracketed_string_literalContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBracketed_string_literalContext)
}

func (s *LogFormatStringContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogFormatStringContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LogFormatStringContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterLogFormatString(s)
	}
}

func (s *LogFormatStringContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitLogFormatString(s)
	}
}

func NewEmptyLogFormatStringContext() *LogFormatStringContext {
	var p = new(LogFormatStringContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_logFormatString
	return p
}

func InitEmptyLogFormatStringContext(p *LogFormatStringContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_logFormatString
}

func NewLogFormatStringContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LogFormatStringContext {
	var p = new(LogFormatStringContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_logFormatString

	return p
}
