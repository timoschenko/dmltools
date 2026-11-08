// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ILogGroupStatementContext is an interface to support dynamic dispatch.
type ILogGroupStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LOGGROUP() antlr.TerminalNode
	Ident() IIdentContext
	SEMI() antlr.TerminalNode

	// IsLogGroupStatementContext differentiates from other interfaces.
	IsLogGroupStatementContext()
}

type LogGroupStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*LogGroupStatementContext) IsLogGroupStatementContext() {}

func (s *LogGroupStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *LogGroupStatementContext) LOGGROUP() antlr.TerminalNode {
	return s.GetToken(dmlParserLOGGROUP, 0)
}

func (s *LogGroupStatementContext) Ident() IIdentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *LogGroupStatementContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *LogGroupStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogGroupStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LogGroupStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterLogGroupStatement(s)
	}
}

func (s *LogGroupStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitLogGroupStatement(s)
	}
}

func NewEmptyLogGroupStatementContext() *LogGroupStatementContext {
	var p = new(LogGroupStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_logGroupStatement
	return p
}

func InitEmptyLogGroupStatementContext(p *LogGroupStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_logGroupStatement
}

func NewLogGroupStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LogGroupStatementContext {
	var p = new(LogGroupStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_logGroupStatement

	return p
}
