// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer

import (
	"fmt"
	"io"

	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer/format"
	"github.com/timoschenko/dmltools/pkg/printer/layout"

	"github.com/antlr4-go/antlr/v4"
)

func PrinterWriteTo(
	tree any,
	file *dmlex.File,
	wr io.Writer,
	lineCap int,
) error {
	return PrinterWriteToDebug(tree, file, wr, lineCap, nil)
}

func PrinterWriteToDebug(
	tree any,
	file *dmlex.File,
	wr io.Writer,
	lineCap int,
	tracer layout.TokenTracer,
) error {
	if file.Parser.HasError() {
		return fmt.Errorf("file has syntax errors")
	}
	p := NewPrinterDebug(lineCap, file, wr, tracer)
	if err := p.Print(tree); err != nil {
		return err
	}
	return p.Flush()
}

// Printer is responsible for traversing the ANTLR AST and mapping grammatical
// constructs to formatting rules. It acts as a Context manager, creating and
// pushing PrintContexts that contain FormatFlag overrides for specific TokenTypes.
//
// Architectural Constraints:
//   - Passes tokens from the AST to the LayoutEngine via processToken or pushTokenRange.
//   - MUST NOT perform manual string processing, length calculations, or manual line wrapping.
//   - MUST NOT manually write spacing or newlines for line wrapping within a statement.
//   - MUST NOT inject synthetic tokens or print new tokens out of nowhere until the
//     Synthetic Token Injection feature is implemented.
//
// Exception for Statement Boundaries:
//   - It IS allowed to insert empty lines between statements using writeNewLines(),
//     because the LayoutEngine lacks top-level statement context. This is the ONLY
//     feature where the Printer is allowed to insert new lines manually.
type Printer struct {
	file     *dmlex.File
	wr       io.Writer
	cts      *antlr.CommonTokenStream
	settings Settings

	lastProcessedTokenIndex TokenIndex

	// TODO: The LayoutEngine is currently exposed here for testing purposes.
	LayoutEngine *layout.Engine
	context      *PrintContext
	source       *dmlex.TokenSource
}

func NewPrinter(maxLineCap int, file *dmlex.File, wr io.Writer) *Printer {
	return NewPrinterDebug(maxLineCap, file, wr, nil)
}

func NewPrinterDebug(
	maxLineCap int,
	file *dmlex.File,
	wr io.Writer,
	tracer layout.TokenTracer,
) *Printer {
	ts := file.Parser.GetTokenStream()
	cts, ok := ts.(*antlr.CommonTokenStream)
	if !ok {
		panic("unreachable: " + fmt.Sprintf("%T", ts))
	}

	settings := DefaultSettings(maxLineCap)
	le := layout.NewEngine(settings, wr, nil)
	le.SetTracer(tracer)

	ctx := NewPrintContext(settings.TokenFormatDefaults(), nil)
	p := &Printer{
		file:                    file,
		wr:                      wr,
		cts:                     cts,
		settings:                settings,
		lastProcessedTokenIndex: 0,
		LayoutEngine:            le,
		context:                 ctx,
		source:                  dmlex.NewTokenSouce(cts.Get(0)),
	}
	return p
}

func (p *Printer) Flush() error {
	if p.LayoutEngine.HasActiveGroups() {
		return nil
	}
	if err := p.LayoutEngine.Flush(); err != nil {
		return err
	}
	return nil
}

func (p *Printer) Print(tree any) error {
	switch v := tree.(type) {
	case dmlex.IDmlTagContext:
		return p.printDmlTag(v)
	case dmlex.IDeviceTagContext:
		return p.printDeviceTag(v)
	case dmlex.IBitorderTagContext:
		return p.printBitorderTag(v)
	case dmlex.IProvisionalTagContext:
		return p.printProvisionalTag(v)
	case dmlex.IDmlContext:
		return p.printDmlContext(v)
	case dmlex.IDmlCodeContext:
		return p.printCodeContext(v)
	case dmlex.IDeviceStatementsContext:
		return p.printDeviceStatements(v)
	case dmlex.IDeviceStatementContext:
		return p.printDeviceStatement(v)
	default:
		panic("Printer.Print is not implemented")
	}
}

func (p *Printer) newContext(
	overrides format.TokenMap,
) {
	p.context = NewPrintContext(
		overrides, p.context,
	)
}

func (p *Printer) freeContext() {
	prev := p.context.Parent()
	if prev == nil {
		panic("Context stack is empty")
	}
	p.context = prev
}
