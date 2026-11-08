// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ILogLevelContext is an interface to support dynamic dispatch.
type ILogLevelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext
	THEN() antlr.TerminalNode
	ThenLevel() IThenLevelContext

	// IsLogLevelContext differentiates from other interfaces.
	IsLogLevelContext()
}

type LogLevelContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*LogLevelContext) IsLogLevelContext() {}

func (s *LogLevelContext) GetParser() antlr.Parser { return s.parser }

func (s *LogLevelContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *LogLevelContext) THEN() antlr.TerminalNode {
	return s.GetToken(dmlParserTHEN, 0)
}

func (s *LogLevelContext) ThenLevel() IThenLevelContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IThenLevelContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IThenLevelContext)
}

func (s *LogLevelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogLevelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LogLevelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterLogLevel(s)
	}
}

func (s *LogLevelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitLogLevel(s)
	}
}

func NewEmptyLogLevelContext() *LogLevelContext {
	var p = new(LogLevelContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_logLevel
	return p
}

func InitEmptyLogLevelContext(p *LogLevelContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_logLevel
}

func NewLogLevelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LogLevelContext {
	var p = new(LogLevelContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_logLevel

	return p
}
