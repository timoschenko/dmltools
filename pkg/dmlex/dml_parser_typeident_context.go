// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ITypeidentContext is an interface to support dynamic dispatch.
type ITypeidentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CHAR() antlr.TerminalNode
	DOUBLE() antlr.TerminalNode
	FLOAT() antlr.TerminalNode
	INT() antlr.TerminalNode
	LONG() antlr.TerminalNode
	SHORT() antlr.TerminalNode
	SIGNED() antlr.TerminalNode
	UNSIGNED() antlr.TerminalNode
	VOID() antlr.TerminalNode
	REGISTER() antlr.TerminalNode
	Ident() IIdentContext

	// IsTypeidentContext differentiates from other interfaces.
	IsTypeidentContext()
}

type TypeidentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*TypeidentContext) IsTypeidentContext() {}

func (s *TypeidentContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeidentContext) CHAR() antlr.TerminalNode {
	return s.GetToken(dmlParserCHAR, 0)
}

func (s *TypeidentContext) DOUBLE() antlr.TerminalNode {
	return s.GetToken(dmlParserDOUBLE, 0)
}

func (s *TypeidentContext) FLOAT() antlr.TerminalNode {
	return s.GetToken(dmlParserFLOAT, 0)
}

func (s *TypeidentContext) INT() antlr.TerminalNode {
	return s.GetToken(dmlParserINT, 0)
}

func (s *TypeidentContext) LONG() antlr.TerminalNode {
	return s.GetToken(dmlParserLONG, 0)
}

func (s *TypeidentContext) SHORT() antlr.TerminalNode {
	return s.GetToken(dmlParserSHORT, 0)
}

func (s *TypeidentContext) SIGNED() antlr.TerminalNode {
	return s.GetToken(dmlParserSIGNED, 0)
}

func (s *TypeidentContext) UNSIGNED() antlr.TerminalNode {
	return s.GetToken(dmlParserUNSIGNED, 0)
}

func (s *TypeidentContext) VOID() antlr.TerminalNode {
	return s.GetToken(dmlParserVOID, 0)
}

func (s *TypeidentContext) REGISTER() antlr.TerminalNode {
	return s.GetToken(dmlParserREGISTER, 0)
}

func (s *TypeidentContext) Ident() IIdentContext {
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

func (s *TypeidentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeidentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeidentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterTypeident(s)
	}
}

func (s *TypeidentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitTypeident(s)
	}
}

func NewEmptyTypeidentContext() *TypeidentContext {
	var p = new(TypeidentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_typeident
	return p
}

func InitEmptyTypeidentContext(p *TypeidentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_typeident
}

func NewTypeidentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeidentContext {
	var p = new(TypeidentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_typeident

	return p
}
