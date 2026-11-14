// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/fs/lookup"
)

// ANSI color escape codes
const (
	ansiReset      = "\033[0m"
	ansiBold       = "\033[1m"
	ansiRed        = "\033[31m"
	ansiBoldRed    = "\033[1;31m"
	ansiGreen      = "\033[32m"
	ansiBoldGreen  = "\033[1;32m"
	ansiYellow     = "\033[33m"
	ansiBoldYellow = "\033[1;33m"
	ansiBlue       = "\033[34m"
	ansiBoldBlue   = "\033[1;34m"
	ansiDarkGray   = "\033[90m"
)

// isTerminal checks if w is an open character device (terminal).
func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// isColorSupported checks whether an output stream supports ANSI colors.
func isColorSupported(w io.Writer) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	return isTerminal(w)
}

type ColoredWriter interface {
	io.Writer

	IsColorEnabled() bool
	SetColorMode(mode ColorMode)
	Unwrap() io.Writer
}

type coloredWriterWrapper struct {
	io.Writer
	hasColor bool
}

func wrapWriter(wr io.Writer) ColoredWriter {
	return &coloredWriterWrapper{
		Writer:   wr,
		hasColor: false,
	}
}

func (cw *coloredWriterWrapper) Unwrap() io.Writer {
	return cw.Writer
}

func (cw *coloredWriterWrapper) IsColorEnabled() bool {
	return cw.hasColor
}

func (cw *coloredWriterWrapper) SetColorMode(mode ColorMode) {
	if mode == ColorAuto {
		cw.hasColor = isColorSupported(cw.Writer)
	} else {
		cw.hasColor = mode == ColorAlways
	}
}

// UI manages user-facing terminal interactions, stream routing, and color styling.
type UI struct {
	mu  sync.Mutex
	out ColoredWriter
	err ColoredWriter
}

// NewUI creates a UI configured with target streams and color mode.
func NewUI(out, err io.Writer, mode ColorMode) *UI {
	if out == nil {
		out = os.Stdout
	}
	if err == nil {
		err = os.Stderr
	}
	u := &UI{
		out: wrapWriter(out),
		err: wrapWriter(err),
	}
	u.SetColorMode(mode)
	return u
}

func (u *UI) Out() ColoredWriter { return u.out }
func (u *UI) Err() ColoredWriter { return u.err }

func (u *UI) SetColorMode(mode ColorMode) {
	u.out.SetColorMode(mode)
	u.err.SetColorMode(mode)
}

func (u *UI) SetColor(color bool) {
	if color {
		u.SetColorMode(ColorAlways)
	} else {
		u.SetColorMode(ColorNever)
	}
}

func (u *UI) Write(p []byte) (int, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.out.Write(p)
}

// Print formats and writes a message to stdout.
func (u *UI) Print(format string, a ...any) {
	msg := fmt.Appendf([]byte{}, format, a...)
	_, _ = u.Write(msg)
}

// PrintError formats and writes an error message to stderr with a trailing newline.
func (u *UI) PrintError(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	// TODO: fmt.Appendf + hasSuffix for bytes
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	if u.err.IsColorEnabled() {
		msg = u.colorizeError(msg)
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	_, _ = io.WriteString(u.err.Unwrap(), msg)
}

func (u *UI) colorizeError(msg string) string {
	hasNewline := strings.HasSuffix(msg, "\n")
	content := strings.TrimSuffix(msg, "\n")

	var colored string
	if rest, ok := strings.CutPrefix(content, "Error: "); ok {
		colored = ansiBoldRed + "Error:" + ansiReset + " " + rest
	} else if rest, ok := strings.CutPrefix(content, "Error "); ok {
		colored = ansiBoldRed + "Error" + ansiReset + " " + rest
	} else {
		colored = ansiBoldRed + content + ansiReset
	}

	if hasNewline {
		colored += "\n"
	}
	return colored
}

// --- Formatting Status Messages ---

// Reformatted writes a file reformat success message to stdout.
func (u *UI) Reformatted(filename string) {
	wr := u.out

	var msg string
	if wr.IsColorEnabled() {
		msg = fmt.Sprintf(
			"%s %s\n",
			ansiBoldGreen+"reformatted"+ansiReset,
			ansiBold+filename+ansiReset,
		)
	} else {
		msg = fmt.Sprintf("reformatted %s\n", filename)
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	_, _ = io.WriteString(wr.Unwrap(), msg)
}

// ReformattedStdin writes a reformat status message to stderr (to avoid corrupting stdout).
func (u *UI) ReformattedStdin(displayName string) {
	wr := u.err

	var msg string
	if wr.IsColorEnabled() {
		msg = fmt.Sprintf(
			"%s %s\n",
			ansiBoldGreen+"reformatted"+ansiReset,
			ansiBold+displayName+ansiReset,
		)
	} else {
		msg = fmt.Sprintf("reformatted %s\n", displayName)
	}

	u.mu.Lock()
	defer u.mu.Unlock()
	_, _ = io.WriteString(wr.Unwrap(), msg)
}

// WouldReformat writes a check-mode notice to stdout.
func (u *UI) WouldReformat(filename string) {
	wr := u.out

	var msg string
	if wr.IsColorEnabled() {
		msg = fmt.Sprintf(
			"%s %s\n",
			ansiBoldBlue+"would reformat"+ansiReset,
			ansiBold+filename+ansiReset,
		)
	} else {
		msg = fmt.Sprintf("would reformat %s\n", filename)
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	_, _ = io.WriteString(wr.Unwrap(), msg)
}

// Ignored writes an ignored-file notice to stdout.
func (u *UI) Ignored(filename, version string) {
	wr := u.out
	var msg string
	if version != "" {
		if wr.IsColorEnabled() {
			msg = fmt.Sprintf(
				"%s %s (dml %s)\n",
				ansiBoldYellow+"ignored"+ansiReset,
				ansiBold+filename+ansiReset,
				version,
			)
		} else {
			msg = fmt.Sprintf("ignored %s (dml %s)\n", filename, version)
		}
	} else {
		if wr.IsColorEnabled() {
			msg = fmt.Sprintf(
				"%s %s\n",
				ansiBoldYellow+"ignored"+ansiReset,
				ansiBold+filename+ansiReset,
			)
		} else {
			msg = fmt.Sprintf("ignored %s\n", filename)
		}
	}

	u.mu.Lock()
	defer u.mu.Unlock()
	_, _ = io.WriteString(wr.Unwrap(), msg)
}

// --- Verbose Diagnostics ---

// Verbose writes a dimmed diagnostic message to stdout when color is enabled.
func (u *UI) Verbose(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	hasNewline := strings.HasSuffix(msg, "\n")
	msg = strings.TrimSuffix(msg, "\n")
	wr := u.out
	if wr.IsColorEnabled() {
		msg = ansiDarkGray + msg + ansiReset
	}
	if hasNewline {
		msg += "\n"
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	_, _ = io.WriteString(wr.Unwrap(), msg)
}

// Processing writes a file start processing notice to stdout.
func (u *UI) Processing(path string) {
	u.Verbose("Processing %s...\n", path)
}

// IgnoredFilter writes a filtered-out file or directory notice to stdout.
func (u *UI) IgnoredFilter(path string, isDir bool, filterName string) {
	reason := "regular expression"
	if filterName == lookup.FileGitIgnore || filterName == lookup.FileIgnore {
		reason = "file content"
	}
	if isDir {
		u.Verbose(
			"%s directory ignored: matches the %s %s\n",
			path,
			filterName,
			reason,
		)
	} else {
		u.Verbose(
			"%s ignored: matches the %s %s\n",
			path,
			filterName,
			reason,
		)
	}
}

// Unchanged writes an unchanged-file notice to stdout.
func (u *UI) Unchanged(path string) {
	u.Verbose("%s already well formatted, good job.\n", path)
}

// --- Syntax & Diagnostics ---

// SyntaxError outputs a parser syntax error to stderr.
func (u *UI) SyntaxError(filename string, line, column int, msg string) {
	wr := u.err

	var outStr string
	if wr.IsColorEnabled() {
		outStr = fmt.Sprintf(
			"%s:%d column %d: %s - %s\n",
			filename,
			line,
			column,
			ansiBoldRed+"Syntax error"+ansiReset,
			msg,
		)
	} else {
		outStr = fmt.Sprintf(
			"%s:%d column %d: Syntax error - %s\n",
			filename,
			line,
			column,
			msg,
		)
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	_, _ = io.WriteString(wr.Unwrap(), outStr)
}

// NewSyntaxErrorListener creates an antlr error listener that routes syntax errors to the UI.
func (u *UI) NewSyntaxErrorListener(filename string) *dmlex.ErrorListener {
	el := dmlex.NewErrorListener(filename, 0)
	el.Handler = func(e dmlex.SyntaxErrorInfo) {
		u.SyntaxError(filename, e.Line, e.Column, e.Msg)
	}
	return el
}

// --- Help & Version Info ---

// PrintUsage outputs CLI usage to stdout.
func (u *UI) PrintUsage() {
	printUsage(u.out)
}

// PrintErrorUsage outputs CLI usage to stderr.
func (u *UI) PrintErrorUsage() {
	printUsage(u.err)
}

// PrintVersion outputs CLI version information to stdout.
func (u *UI) PrintVersion() {
	printVersion(u.out)
}

// PrintNoFiles writes the notice when no DML files were found to process.
func (u *UI) PrintNoFiles() {
	u.mu.Lock()
	defer u.mu.Unlock()
	_, _ = io.WriteString(
		u.out.Unwrap(),
		"No DML files are present to be formatted. Nothing to do.\n",
	)
}

// PrintReport renders and writes the summary line to stdout if there are processed files.
func (u *UI) PrintReport(r *Report) {
	if r == nil || r.IsEmpty() {
		return
	}
	msg := r.Render(u.out.IsColorEnabled()) + "\n"
	u.mu.Lock()
	defer u.mu.Unlock()
	_, _ = io.WriteString(u.out.Unwrap(), msg)
}

// PrintNoSourcesError writes the usage error when neither SRC nor code is provided to stderr.
func (u *UI) PrintNoSourcesError() {
	u.mu.Lock()
	defer u.mu.Unlock()
	_, _ = io.WriteString(
		u.err.Unwrap(),
		"Usage: dmlfmt [OPTIONS] SRC ...\n\nOne of 'SRC' or 'code' is required.\n",
	)
}
