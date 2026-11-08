// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IDmlCodeContext is an interface to support dynamic dispatch.
type IDmlCodeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DeviceStatements() IDeviceStatementsContext
	ProvisionalTag() IProvisionalTagContext
	DeviceTag() IDeviceTagContext
	BitorderTag() IBitorderTagContext

	// IsDmlCodeContext differentiates from other interfaces.
	IsDmlCodeContext()
}

type DmlCodeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*DmlCodeContext) IsDmlCodeContext() {}

func (s *DmlCodeContext) GetParser() antlr.Parser { return s.parser }

func (s *DmlCodeContext) DeviceStatements() IDeviceStatementsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDeviceStatementsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDeviceStatementsContext)
}

func (s *DmlCodeContext) ProvisionalTag() IProvisionalTagContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IProvisionalTagContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IProvisionalTagContext)
}

func (s *DmlCodeContext) DeviceTag() IDeviceTagContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDeviceTagContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDeviceTagContext)
}

func (s *DmlCodeContext) BitorderTag() IBitorderTagContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBitorderTagContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBitorderTagContext)
}

func (s *DmlCodeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DmlCodeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DmlCodeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterDmlCode(s)
	}
}

func (s *DmlCodeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitDmlCode(s)
	}
}

func NewEmptyDmlCodeContext() *DmlCodeContext {
	var p = new(DmlCodeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_dmlCode
	return p
}

func InitEmptyDmlCodeContext(p *DmlCodeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_dmlCode
}

func NewDmlCodeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DmlCodeContext {
	var p = new(DmlCodeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_dmlCode

	return p
}
