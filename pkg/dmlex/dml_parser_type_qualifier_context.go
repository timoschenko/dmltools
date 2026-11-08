// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ITypeQualifierContext is an interface to support dynamic dispatch.
type ITypeQualifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CONST() antlr.TerminalNode
	MUL() antlr.TerminalNode
	VECT() antlr.TerminalNode

	// IsTypeQualifierContext differentiates from other interfaces.
	IsTypeQualifierContext()
}

type TypeQualifierContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*TypeQualifierContext) IsTypeQualifierContext() {}

func (s *TypeQualifierContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeQualifierContext) CONST() antlr.TerminalNode {
	return s.GetToken(dmlParserCONST, 0)
}

func (s *TypeQualifierContext) MUL() antlr.TerminalNode {
	return s.GetToken(dmlParserMUL, 0)
}

func (s *TypeQualifierContext) VECT() antlr.TerminalNode {
	return s.GetToken(dmlParserVECT, 0)
}

func (s *TypeQualifierContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeQualifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeQualifierContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterTypeQualifier(s)
	}
}

func (s *TypeQualifierContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitTypeQualifier(s)
	}
}

func NewEmptyTypeQualifierContext() *TypeQualifierContext {
	var p = new(TypeQualifierContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_typeQualifier
	return p
}

func InitEmptyTypeQualifierContext(p *TypeQualifierContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_typeQualifier
}

func NewTypeQualifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeQualifierContext {
	var p = new(TypeQualifierContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_typeQualifier

	return p
}
