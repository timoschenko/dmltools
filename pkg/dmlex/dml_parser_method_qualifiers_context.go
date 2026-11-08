// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IMethodQualifiersContext is an interface to support dynamic dispatch.
type IMethodQualifiersContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INDEPENDENT() antlr.TerminalNode
	SHARED() antlr.TerminalNode
	STARTUP() antlr.TerminalNode
	MEMOIZED() antlr.TerminalNode
	INLINE() antlr.TerminalNode

	// IsMethodQualifiersContext differentiates from other interfaces.
	IsMethodQualifiersContext()
}

type MethodQualifiersContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*MethodQualifiersContext) IsMethodQualifiersContext() {}

func (s *MethodQualifiersContext) GetParser() antlr.Parser { return s.parser }

func (s *MethodQualifiersContext) INDEPENDENT() antlr.TerminalNode {
	return s.GetToken(dmlParserINDEPENDENT, 0)
}

func (s *MethodQualifiersContext) SHARED() antlr.TerminalNode {
	return s.GetToken(dmlParserSHARED, 0)
}

func (s *MethodQualifiersContext) STARTUP() antlr.TerminalNode {
	return s.GetToken(dmlParserSTARTUP, 0)
}

func (s *MethodQualifiersContext) MEMOIZED() antlr.TerminalNode {
	return s.GetToken(dmlParserMEMOIZED, 0)
}

func (s *MethodQualifiersContext) INLINE() antlr.TerminalNode {
	return s.GetToken(dmlParserINLINE, 0)
}

func (s *MethodQualifiersContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MethodQualifiersContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MethodQualifiersContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterMethodQualifiers(s)
	}
}

func (s *MethodQualifiersContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitMethodQualifiers(s)
	}
}

func NewEmptyMethodQualifiersContext() *MethodQualifiersContext {
	var p = new(MethodQualifiersContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_methodQualifiers
	return p
}

func InitEmptyMethodQualifiersContext(p *MethodQualifiersContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_methodQualifiers
}

func NewMethodQualifiersContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MethodQualifiersContext {
	var p = new(MethodQualifiersContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_methodQualifiers

	return p
}
