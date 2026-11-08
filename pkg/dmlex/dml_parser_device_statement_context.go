// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IDeviceStatementContext is an interface to support dynamic dispatch.
type IDeviceStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TemplateStatement() ITemplateStatementContext
	CCodeBlock() ICCodeBlockContext
	LogGroupStatement() ILogGroupStatementContext
	ConstantStatement() IConstantStatementContext
	TypedefStatement() ITypedefStatementContext
	ExternStatement() IExternStatementContext
	ImportStatement() IImportStatementContext
	ExportStatement() IExportStatementContext
	Object() IObjectContext
	Param() IParamContext
	MethodUnified() IMethodUnifiedContext
	IstemplateStatement() IIstemplateStatementContext
	ToplevelIf() IToplevelIfContext
	ErrorStatement() IErrorStatementContext
	IneachStatement() IIneachStatementContext

	// IsDeviceStatementContext differentiates from other interfaces.
	IsDeviceStatementContext()
}

type DeviceStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*DeviceStatementContext) IsDeviceStatementContext() {}

func (s *DeviceStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *DeviceStatementContext) TemplateStatement() ITemplateStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITemplateStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITemplateStatementContext)
}

func (s *DeviceStatementContext) CCodeBlock() ICCodeBlockContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICCodeBlockContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICCodeBlockContext)
}

func (s *DeviceStatementContext) LogGroupStatement() ILogGroupStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILogGroupStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILogGroupStatementContext)
}

func (s *DeviceStatementContext) ConstantStatement() IConstantStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConstantStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConstantStatementContext)
}

func (s *DeviceStatementContext) TypedefStatement() ITypedefStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypedefStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypedefStatementContext)
}

func (s *DeviceStatementContext) ExternStatement() IExternStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExternStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExternStatementContext)
}

func (s *DeviceStatementContext) ImportStatement() IImportStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IImportStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IImportStatementContext)
}

func (s *DeviceStatementContext) ExportStatement() IExportStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExportStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExportStatementContext)
}

func (s *DeviceStatementContext) Object() IObjectContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectContext)
}

func (s *DeviceStatementContext) Param() IParamContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParamContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *DeviceStatementContext) MethodUnified() IMethodUnifiedContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMethodUnifiedContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMethodUnifiedContext)
}

func (s *DeviceStatementContext) IstemplateStatement() IIstemplateStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIstemplateStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIstemplateStatementContext)
}

func (s *DeviceStatementContext) ToplevelIf() IToplevelIfContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IToplevelIfContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IToplevelIfContext)
}

func (s *DeviceStatementContext) ErrorStatement() IErrorStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IErrorStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IErrorStatementContext)
}

func (s *DeviceStatementContext) IneachStatement() IIneachStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIneachStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIneachStatementContext)
}

func (s *DeviceStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DeviceStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DeviceStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterDeviceStatement(s)
	}
}

func (s *DeviceStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitDeviceStatement(s)
	}
}

func NewEmptyDeviceStatementContext() *DeviceStatementContext {
	var p = new(DeviceStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_deviceStatement
	return p
}

func InitEmptyDeviceStatementContext(p *DeviceStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_deviceStatement
}

func NewDeviceStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeviceStatementContext {
	var p = new(DeviceStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_deviceStatement

	return p
}
