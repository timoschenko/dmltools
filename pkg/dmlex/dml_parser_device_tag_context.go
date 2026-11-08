// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IDeviceTagContext is an interface to support dynamic dispatch.
type IDeviceTagContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DEVICE() antlr.TerminalNode
	Objident() IObjidentContext
	SEMI() antlr.TerminalNode

	// IsDeviceTagContext differentiates from other interfaces.
	IsDeviceTagContext()
}

type DeviceTagContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*DeviceTagContext) IsDeviceTagContext() {}

func (s *DeviceTagContext) GetParser() antlr.Parser { return s.parser }

func (s *DeviceTagContext) DEVICE() antlr.TerminalNode {
	return s.GetToken(dmlParserDEVICE, 0)
}

func (s *DeviceTagContext) Objident() IObjidentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjidentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjidentContext)
}

func (s *DeviceTagContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *DeviceTagContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DeviceTagContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DeviceTagContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterDeviceTag(s)
	}
}

func (s *DeviceTagContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitDeviceTag(s)
	}
}

func NewEmptyDeviceTagContext() *DeviceTagContext {
	var p = new(DeviceTagContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_deviceTag
	return p
}

func InitEmptyDeviceTagContext(p *DeviceTagContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_deviceTag
}

func NewDeviceTagContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeviceTagContext {
	var p = new(DeviceTagContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_deviceTag

	return p
}
