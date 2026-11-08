// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ICatchStatementContext is an interface to support dynamic dispatch.
type ICatchStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Statement() IStatementContext

	// IsCatchStatementContext differentiates from other interfaces.
	IsCatchStatementContext()
}

type CatchStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*CatchStatementContext) IsCatchStatementContext() {}

func (s *CatchStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *CatchStatementContext) Statement() IStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStatementContext)
}

func (s *CatchStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CatchStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CatchStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterCatchStatement(s)
	}
}

func (s *CatchStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitCatchStatement(s)
	}
}

func NewEmptyCatchStatementContext() *CatchStatementContext {
	var p = new(CatchStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_catchStatement
	return p
}

func InitEmptyCatchStatementContext(p *CatchStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_catchStatement
}

func NewCatchStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CatchStatementContext {
	var p = new(CatchStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_catchStatement

	return p
}
