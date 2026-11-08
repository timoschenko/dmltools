// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IObjectDescSpecContext is an interface to support dynamic dispatch.
type IObjectDescSpecContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Composed_string_literal() IComposed_string_literalContext

	// IsObjectDescSpecContext differentiates from other interfaces.
	IsObjectDescSpecContext()
}

type ObjectDescSpecContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ObjectDescSpecContext) IsObjectDescSpecContext() {}

func (s *ObjectDescSpecContext) GetParser() antlr.Parser { return s.parser }

func (s *ObjectDescSpecContext) Composed_string_literal() IComposed_string_literalContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IComposed_string_literalContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IComposed_string_literalContext)
}

func (s *ObjectDescSpecContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ObjectDescSpecContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ObjectDescSpecContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterObjectDescSpec(s)
	}
}

func (s *ObjectDescSpecContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitObjectDescSpec(s)
	}
}

func NewEmptyObjectDescSpecContext() *ObjectDescSpecContext {
	var p = new(ObjectDescSpecContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_objectDescSpec
	return p
}

func InitEmptyObjectDescSpecContext(p *ObjectDescSpecContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_objectDescSpec
}

func NewObjectDescSpecContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ObjectDescSpecContext {
	var p = new(ObjectDescSpecContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_objectDescSpec

	return p
}
