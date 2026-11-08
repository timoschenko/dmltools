// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IDesignated_struct_initializer_listContext is an interface to support dynamic dispatch.
type IDesignated_struct_initializer_listContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllDesignatedStructInitializer() []IDesignatedStructInitializerContext
	DesignatedStructInitializer(i int) IDesignatedStructInitializerContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode
	ELLIPSIS() antlr.TerminalNode

	// IsDesignated_struct_initializer_listContext differentiates from other interfaces.
	IsDesignated_struct_initializer_listContext()
}

type Designated_struct_initializer_listContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Designated_struct_initializer_listContext) IsDesignated_struct_initializer_listContext() {}

func (s *Designated_struct_initializer_listContext) GetParser() antlr.Parser { return s.parser }

func (s *Designated_struct_initializer_listContext) AllDesignatedStructInitializer() []IDesignatedStructInitializerContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IDesignatedStructInitializerContext); ok {
			len++
		}
	}

	tst := make([]IDesignatedStructInitializerContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IDesignatedStructInitializerContext); ok {
			tst[i] = t.(IDesignatedStructInitializerContext)
			i++
		}
	}

	return tst
}

func (s *Designated_struct_initializer_listContext) DesignatedStructInitializer(i int) IDesignatedStructInitializerContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDesignatedStructInitializerContext); ok {
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

	return t.(IDesignatedStructInitializerContext)
}

func (s *Designated_struct_initializer_listContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(dmlParserCOMMA)
}

func (s *Designated_struct_initializer_listContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserCOMMA, i)
}

func (s *Designated_struct_initializer_listContext) ELLIPSIS() antlr.TerminalNode {
	return s.GetToken(dmlParserELLIPSIS, 0)
}

func (s *Designated_struct_initializer_listContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Designated_struct_initializer_listContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Designated_struct_initializer_listContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterDesignated_struct_initializer_list(s)
	}
}

func (s *Designated_struct_initializer_listContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitDesignated_struct_initializer_list(s)
	}
}

func NewEmptyDesignated_struct_initializer_listContext() *Designated_struct_initializer_listContext {
	var p = new(Designated_struct_initializer_listContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_designated_struct_initializer_list
	return p
}

func InitEmptyDesignated_struct_initializer_listContext(p *Designated_struct_initializer_listContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_designated_struct_initializer_list
}

func NewDesignated_struct_initializer_listContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Designated_struct_initializer_listContext {
	var p = new(Designated_struct_initializer_listContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_designated_struct_initializer_list

	return p
}
