// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IMethodParamsMaybeUntypedContext is an interface to support dynamic dispatch.
type IMethodParamsMaybeUntypedContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllCdeclOrInline() []ICdeclOrInlineContext
	CdeclOrInline(i int) ICdeclOrInlineContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsMethodParamsMaybeUntypedContext differentiates from other interfaces.
	IsMethodParamsMaybeUntypedContext()
}

type MethodParamsMaybeUntypedContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*MethodParamsMaybeUntypedContext) IsMethodParamsMaybeUntypedContext() {}

func (s *MethodParamsMaybeUntypedContext) GetParser() antlr.Parser { return s.parser }

func (s *MethodParamsMaybeUntypedContext) AllCdeclOrInline() []ICdeclOrInlineContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICdeclOrInlineContext); ok {
			len++
		}
	}

	tst := make([]ICdeclOrInlineContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICdeclOrInlineContext); ok {
			tst[i] = t.(ICdeclOrInlineContext)
			i++
		}
	}

	return tst
}

func (s *MethodParamsMaybeUntypedContext) CdeclOrInline(i int) ICdeclOrInlineContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICdeclOrInlineContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICdeclOrInlineContext)
}

func (s *MethodParamsMaybeUntypedContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(dmlParserCOMMA)
}

func (s *MethodParamsMaybeUntypedContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserCOMMA, i)
}

func (s *MethodParamsMaybeUntypedContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MethodParamsMaybeUntypedContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MethodParamsMaybeUntypedContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterMethodParamsMaybeUntyped(s)
	}
}

func (s *MethodParamsMaybeUntypedContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitMethodParamsMaybeUntyped(s)
	}
}

func NewEmptyMethodParamsMaybeUntypedContext() *MethodParamsMaybeUntypedContext {
	var p = new(MethodParamsMaybeUntypedContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_methodParamsMaybeUntyped
	return p
}

func InitEmptyMethodParamsMaybeUntypedContext(p *MethodParamsMaybeUntypedContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_methodParamsMaybeUntyped
}

func NewMethodParamsMaybeUntypedContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MethodParamsMaybeUntypedContext {
	var p = new(MethodParamsMaybeUntypedContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_methodParamsMaybeUntyped

	return p
}
