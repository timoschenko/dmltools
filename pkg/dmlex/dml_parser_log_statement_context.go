// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ILogStatementContext is an interface to support dynamic dispatch.
type ILogStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LOG() antlr.TerminalNode
	LogKind() ILogKindContext
	COLON() antlr.TerminalNode
	LogFormatString() ILogFormatStringContext
	SEMI() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode
	LogLevel() ILogLevelContext
	LogGroup() ILogGroupContext
	Expression_list_ntc_nonempty() IExpression_list_ntc_nonemptyContext

	// IsLogStatementContext differentiates from other interfaces.
	IsLogStatementContext()
}

type LogStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*LogStatementContext) IsLogStatementContext() {}

func (s *LogStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *LogStatementContext) LOG() antlr.TerminalNode {
	return s.GetToken(dmlParserLOG, 0)
}

func (s *LogStatementContext) LogKind() ILogKindContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILogKindContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILogKindContext)
}

func (s *LogStatementContext) COLON() antlr.TerminalNode {
	return s.GetToken(dmlParserCOLON, 0)
}

func (s *LogStatementContext) LogFormatString() ILogFormatStringContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILogFormatStringContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILogFormatStringContext)
}

func (s *LogStatementContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *LogStatementContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(dmlParserCOMMA)
}

func (s *LogStatementContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserCOMMA, i)
}

func (s *LogStatementContext) LogLevel() ILogLevelContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILogLevelContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILogLevelContext)
}

func (s *LogStatementContext) LogGroup() ILogGroupContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILogGroupContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILogGroupContext)
}

func (s *LogStatementContext) Expression_list_ntc_nonempty() IExpression_list_ntc_nonemptyContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpression_list_ntc_nonemptyContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpression_list_ntc_nonemptyContext)
}

func (s *LogStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LogStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterLogStatement(s)
	}
}

func (s *LogStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitLogStatement(s)
	}
}

func NewEmptyLogStatementContext() *LogStatementContext {
	var p = new(LogStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_logStatement
	return p
}

func InitEmptyLogStatementContext(p *LogStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_logStatement
}

func NewLogStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LogStatementContext {
	var p = new(LogStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_logStatement

	return p
}
