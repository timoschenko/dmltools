// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ICdeclOrInlineContext is an interface to support dynamic dispatch.
type ICdeclOrInlineContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Cdecl() ICdeclContext
	INLINE() antlr.TerminalNode
	Ident() IIdentContext

	// IsCdeclOrInlineContext differentiates from other interfaces.
	IsCdeclOrInlineContext()
}

type CdeclOrInlineContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*CdeclOrInlineContext) IsCdeclOrInlineContext() {}

func (s *CdeclOrInlineContext) GetParser() antlr.Parser { return s.parser }

func (s *CdeclOrInlineContext) Cdecl() ICdeclContext {
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

func (s *CdeclOrInlineContext) INLINE() antlr.TerminalNode {
	return s.GetToken(dmlParserINLINE, 0)
}

func (s *CdeclOrInlineContext) Ident() IIdentContext {
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

func (s *CdeclOrInlineContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CdeclOrInlineContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CdeclOrInlineContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterCdeclOrInline(s)
	}
}

func (s *CdeclOrInlineContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitCdeclOrInline(s)
	}
}

func NewEmptyCdeclOrInlineContext() *CdeclOrInlineContext {
	var p = new(CdeclOrInlineContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_cdeclOrInline
	return p
}

func InitEmptyCdeclOrInlineContext(p *CdeclOrInlineContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_cdeclOrInline
}

func NewCdeclOrInlineContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CdeclOrInlineContext {
	var p = new(CdeclOrInlineContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_cdeclOrInline

	return p
}
