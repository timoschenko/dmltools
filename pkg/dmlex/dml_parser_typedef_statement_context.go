// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ITypedefStatementContext is an interface to support dynamic dispatch.
type ITypedefStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TYPEDEF() antlr.TerminalNode
	Cdecl() ICdeclContext
	SEMI() antlr.TerminalNode

	// IsTypedefStatementContext differentiates from other interfaces.
	IsTypedefStatementContext()
}

type TypedefStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*TypedefStatementContext) IsTypedefStatementContext() {}

func (s *TypedefStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *TypedefStatementContext) TYPEDEF() antlr.TerminalNode {
	return s.GetToken(dmlParserTYPEDEF, 0)
}

func (s *TypedefStatementContext) Cdecl() ICdeclContext {
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

func (s *TypedefStatementContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *TypedefStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypedefStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypedefStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterTypedefStatement(s)
	}
}

func (s *TypedefStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitTypedefStatement(s)
	}
}

func NewEmptyTypedefStatementContext() *TypedefStatementContext {
	var p = new(TypedefStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_typedefStatement
	return p
}

func InitEmptyTypedefStatementContext(p *TypedefStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_typedefStatement
}

func NewTypedefStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypedefStatementContext {
	var p = new(TypedefStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_typedefStatement

	return p
}
