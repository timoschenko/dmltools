// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IExternStatementContext is an interface to support dynamic dispatch.
type IExternStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EXTERN() antlr.TerminalNode
	Cdecl() ICdeclContext
	SEMI() antlr.TerminalNode
	TYPEDEF() antlr.TerminalNode

	// IsExternStatementContext differentiates from other interfaces.
	IsExternStatementContext()
}

type ExternStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ExternStatementContext) IsExternStatementContext() {}

func (s *ExternStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *ExternStatementContext) EXTERN() antlr.TerminalNode {
	return s.GetToken(dmlParserEXTERN, 0)
}

func (s *ExternStatementContext) Cdecl() ICdeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICdeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICdeclContext)
}

func (s *ExternStatementContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *ExternStatementContext) TYPEDEF() antlr.TerminalNode {
	return s.GetToken(dmlParserTYPEDEF, 0)
}

func (s *ExternStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExternStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExternStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterExternStatement(s)
	}
}

func (s *ExternStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitExternStatement(s)
	}
}

func NewEmptyExternStatementContext() *ExternStatementContext {
	var p = new(ExternStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_externStatement
	return p
}

func InitEmptyExternStatementContext(p *ExternStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_externStatement
}

func NewExternStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternStatementContext {
	var p = new(ExternStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_externStatement

	return p
}
