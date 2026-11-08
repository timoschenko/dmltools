// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IObject_typeContext is an interface to support dynamic dispatch.
type IObject_typeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ATTRIBUTE() antlr.TerminalNode
	BANK() antlr.TerminalNode
	CONNECT() antlr.TerminalNode
	EVENT() antlr.TerminalNode
	GROUP() antlr.TerminalNode
	IMPLEMENT() antlr.TerminalNode
	INTERFACE() antlr.TerminalNode
	PORT() antlr.TerminalNode
	SUBDEVICE() antlr.TerminalNode

	// IsObject_typeContext differentiates from other interfaces.
	IsObject_typeContext()
}

type Object_typeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Object_typeContext) IsObject_typeContext() {}

func (s *Object_typeContext) GetParser() antlr.Parser { return s.parser }

func (s *Object_typeContext) ATTRIBUTE() antlr.TerminalNode {
	return s.GetToken(dmlParserATTRIBUTE, 0)
}

func (s *Object_typeContext) BANK() antlr.TerminalNode {
	return s.GetToken(dmlParserBANK, 0)
}

func (s *Object_typeContext) CONNECT() antlr.TerminalNode {
	return s.GetToken(dmlParserCONNECT, 0)
}

func (s *Object_typeContext) EVENT() antlr.TerminalNode {
	return s.GetToken(dmlParserEVENT, 0)
}

func (s *Object_typeContext) GROUP() antlr.TerminalNode {
	return s.GetToken(dmlParserGROUP, 0)
}

func (s *Object_typeContext) IMPLEMENT() antlr.TerminalNode {
	return s.GetToken(dmlParserIMPLEMENT, 0)
}

func (s *Object_typeContext) INTERFACE() antlr.TerminalNode {
	return s.GetToken(dmlParserINTERFACE, 0)
}

func (s *Object_typeContext) PORT() antlr.TerminalNode {
	return s.GetToken(dmlParserPORT, 0)
}

func (s *Object_typeContext) SUBDEVICE() antlr.TerminalNode {
	return s.GetToken(dmlParserSUBDEVICE, 0)
}

func (s *Object_typeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Object_typeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Object_typeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterObject_type(s)
	}
}

func (s *Object_typeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitObject_type(s)
	}
}

func NewEmptyObject_typeContext() *Object_typeContext {
	var p = new(Object_typeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_object_type
	return p
}

func InitEmptyObject_typeContext(p *Object_typeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_object_type
}

func NewObject_typeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Object_typeContext {
	var p = new(Object_typeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_object_type

	return p
}
