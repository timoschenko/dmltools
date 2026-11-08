// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer

import (
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"strings"
)

// g4-tag: cCodeBlock
func (p *Printer) printCCodeBlock(node dmlex.ICCodeBlockContext) error {
	headerToken := node.HEADER()
	if headerToken == nil {
		headerToken = node.HEADER2()
	}
	if headerToken == nil {
		headerToken = node.FOOTER()
	}

	headerIdx := p.gatherTokenIndexFromRightiestTerminalNode(headerToken)

	// push up to header token
	if err := p.pushTokenRange(headerIdx); err != nil {
		return err
	}

	ccbNode := node.C_CODE_BLOCK()
	ccbIdx := p.gatherTokenIndexFromRightiestTerminalNode(ccbNode)

	// push any tokens (like spaces/comments) between header and C_CODE_BLOCK
	if err := p.pushTokenRange(ccbIdx - 1); err != nil {
		return err
	}

	token := p.extractToken(ccbIdx)
	text := token.GetText()

	if len(text) >= 4 && strings.HasPrefix(text, "%{") &&
		strings.HasSuffix(text, "%}") {
		inner := text[2 : len(text)-2]
		lines := strings.Split(inner, "\n")
		for i, line := range lines {
			lines[i] = strings.TrimRight(line, " \t\r")
		}
		processedInner := strings.Join(lines, "\n")
		if !strings.HasPrefix(processedInner, "\n") {
			processedInner = "\n" + processedInner
		}
		if !strings.HasSuffix(processedInner, "\n") {
			processedInner = processedInner + "\n"
		}
		text = "%{" + processedInner + "%}"
	}

	synthToken := &trimmedToken{Token: token, text: text}
	if err := p.processToken(synthToken, p.context); err != nil {
		return err
	}
	p.lastProcessedTokenIndex = ccbIdx + 1

	return p.Flush()
}
