// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer

import (
	dmldefaults "github.com/timoschenko/dmltools/pkg/printer/dml_defaults"
	"github.com/timoschenko/dmltools/pkg/printer/format"
)

// Deprecated: use Settings.TabWidth instead
const TabSize = 4

const (
	CFG_NEW_LINES_DEFAULT      = 1
	CFG_NEW_LINES_ON_TOP_LEVEL = 2
)

// Settings encapsulates the configuration parameters that influence
// the formatting behavior of the printer.
type Settings interface {

	// Maximum recommended line length.
	// The printer will attempt to wrap lines that exceed this length.
	LineCap() int

	// Number of spaces that a tab character should be expanded to.
	// This is used for calculating the visual width of tokens containing tabs.
	TabWidth() int

	TokenFormatDefaults() format.TokenMap
	ExpressionDefaults() format.TokenMap
	CdeclDefaults() format.TokenMap
}

// Reflects the `Settings` inteface
type settings struct {
	lineCap  int
	tabWidth int

	tokenFormatDefaults format.TokenMap
	expressionDefaults  format.TokenMap
	cdeclDefaults       format.TokenMap
}

func (s *settings) LineCap() int {
	return s.lineCap
}

func (s *settings) TabWidth() int {
	return s.tabWidth
}

func (s *settings) TokenFormatDefaults() format.TokenMap {
	return s.tokenFormatDefaults
}

func (s *settings) ExpressionDefaults() format.TokenMap {
	return s.expressionDefaults
}

func (s *settings) CdeclDefaults() format.TokenMap {
	return s.cdeclDefaults
}

func DefaultSettings(lineCap int) Settings {
	return &settings{
		lineCap:             lineCap,
		tabWidth:            TabSize, // Default tab width
		tokenFormatDefaults: dmldefaults.Tokens(),
		expressionDefaults:  dmldefaults.Expression(),
		cdeclDefaults:       dmldefaults.Cdecl(),
	}
}

func NewSettingsDebug(
	lineCap int,
	tokenFormatDefaults, expressionDefaults, cdeclDefaults format.TokenMap,
) Settings {
	return &settings{
		lineCap:             lineCap,
		tabWidth:            TabSize, // Default tab width
		tokenFormatDefaults: tokenFormatDefaults,
		expressionDefaults:  expressionDefaults,
		cdeclDefaults:       cdeclDefaults,
	}
}
