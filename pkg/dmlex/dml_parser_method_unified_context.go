// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IMethodUnifiedContext is an interface to support dynamic dispatch.
type IMethodUnifiedContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsMethodUnifiedContext differentiates from other interfaces.
	IsMethodUnifiedContext()
}

type MethodUnifiedContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*MethodUnifiedContext) IsMethodUnifiedContext() {}

func (s *MethodUnifiedContext) GetParser() antlr.Parser { return s.parser }

func (s *MethodUnifiedContext) CopyAll(ctx *MethodUnifiedContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *MethodUnifiedContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MethodUnifiedContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func NewEmptyMethodUnifiedContext() *MethodUnifiedContext {
	var p = new(MethodUnifiedContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_methodUnified
	return p
}

func InitEmptyMethodUnifiedContext(p *MethodUnifiedContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_methodUnified
}

func NewMethodUnifiedContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MethodUnifiedContext {
	var p = new(MethodUnifiedContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_methodUnified

	return p
}
