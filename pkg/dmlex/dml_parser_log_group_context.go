// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ILogGroupContext is an interface to support dynamic dispatch.
type ILogGroupContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext

	// IsLogGroupContext differentiates from other interfaces.
	IsLogGroupContext()
}

type LogGroupContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*LogGroupContext) IsLogGroupContext() {}

func (s *LogGroupContext) GetParser() antlr.Parser { return s.parser }

func (s *LogGroupContext) Expression() IExpressionContext {
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

func (s *LogGroupContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogGroupContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LogGroupContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterLogGroup(s)
	}
}

func (s *LogGroupContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitLogGroup(s)
	}
}

func NewEmptyLogGroupContext() *LogGroupContext {
	var p = new(LogGroupContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_logGroup
	return p
}

func InitEmptyLogGroupContext(p *LogGroupContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_logGroup
}

func NewLogGroupContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LogGroupContext {
	var p = new(LogGroupContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_logGroup

	return p
}
