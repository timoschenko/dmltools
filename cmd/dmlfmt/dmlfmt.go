// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"slices"

	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/fs/lookup"
	"github.com/timoschenko/dmltools/pkg/fs/stream"
	"github.com/timoschenko/dmltools/pkg/printer"

	"github.com/antlr4-go/antlr/v4"
)

func PrintSyntaxError(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func run(args []string, stdout, stderr io.Writer) int {
	return runWithIO(args, os.Stdin, stdout, stderr)
}

// TODO: AI Slop
func runWithIO(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	cfg, err := parseArgs(args, stderr)
	if err != nil {
		return ExitCodeInternalError
	}

	isStdin := len(cfg.Sources) == 1 && cfg.Sources[0] == "-"
	stdoutWriter := stdout
	if cfg.Quiet && !isStdin && !cfg.CodeSet && !cfg.ShowHelp &&
		!cfg.ShowVersion {
		stdoutWriter = io.Discard
	}

	ui := NewUI(stdoutWriter, stderr, cfg.ColorMode)

	if cfg.ShowHelp {
		ui.PrintUsage()
		return ExitCodeOK
	}

	if cfg.ShowVersion {
		ui.PrintVersion()
		return ExitCodeOK
	}

	var stats StatsCollector = &noopStatsCollector{}
	if cfg.NerdStats {
		stats = newNerdStatsCollector()
		defer func() {
			stats.PrintSummary(ui.Err())
		}()
	}

	// Reject any unimplemented flags with an explicit error and exit code 123
	for _, flagName := range cfg.ExplicitFlags {
		if !isFlagImplemented(flagName) {
			ui.PrintError("Error: flag -%s is not implemented yet.", flagName)
			return ExitCodeInternalError
		}
	}

	if slices.Contains(cfg.ExplicitFlags, "cpuprofile") &&
		cfg.CPUProfile == "" {
		ui.PrintError("Error: flag -cpuprofile requires a non-empty file path.")
		return ExitCodeInternalError
	}
	if slices.Contains(cfg.ExplicitFlags, "memprofile") &&
		cfg.MemProfile == "" {
		ui.PrintError("Error: flag -memprofile requires a non-empty file path.")
		return ExitCodeInternalError
	}

	if cfg.WorkersErr != "" {
		ui.PrintError("Error: %s", cfg.WorkersErr)
		return ExitCodeInternalError
	}

	if cfg.CodeSet && len(cfg.Sources) > 0 {
		ui.PrintError("'SRC' and 'code' cannot be passed simultaneously.")
		ui.PrintErrorUsage()
		return ExitCodeInternalError
	}

	if !cfg.CodeSet && len(cfg.Sources) == 0 {
		ui.PrintNoSourcesError()
		return ExitCodeInternalError
	}

	if cfg.CPUProfile != "" {
		f, err := os.Create(filepath.Clean(cfg.CPUProfile))
		if err != nil {
			ui.PrintError("Error creating CPU profile: %v", err)
			return ExitCodeInternalError
		}
		defer func() {
			_ = f.Close()
		}()
		if err := pprof.StartCPUProfile(f); err != nil {
			ui.PrintError("Error starting CPU profile: %v", err)
			return ExitCodeInternalError
		}
		defer pprof.StopCPUProfile()
	}

	if cfg.MemProfile != "" {
		f, err := os.Create(filepath.Clean(cfg.MemProfile))
		if err != nil {
			ui.PrintError("Error creating memory profile: %v", err)
			return ExitCodeInternalError
		}
		defer func() {
			runtime.GC()
			if err := pprof.WriteHeapProfile(f); err != nil {
				ui.PrintError("Error writing memory profile: %v", err)
			}
			_ = f.Close()
		}()
	}

	if cfg.CodeSet {
		return formatCode(
			cfg.StdinFilename,
			cfg.Code,
			cfg.LineLength,
			cfg.Check,
			cfg.Quiet,
			ui,
			stats,
		)
	}

	finder, err := newFinderFromConfig(cfg)
	if err != nil {
		ui.PrintError("Error: %v", err)
		return ExitCodeInternalError
	}

	if len(cfg.Sources) == 1 && cfg.Sources[0] == "-" {
		if cfg.StdinFilename != "" &&
			finder.Filter().IsForceExcluded(cfg.StdinFilename) {
			if !cfg.Check {
				if stdin == nil {
					stdin = os.Stdin
				}
				_, _ = io.Copy(ui.Out(), stdin)
			}
			return ExitCodeOK
		}
		return formatStdin(
			cfg.StdinFilename,
			cfg.LineLength,
			cfg.Check,
			cfg.Quiet,
			stdin,
			ui,
			stats,
		)
	}

	if slices.Contains(cfg.Sources, "-") {
		ui.PrintError("Error: Cannot mix stdin ('-') with other file paths.")
		return ExitCodeInternalError
	}

	return runPipeline(cfg, finder, ui, stats)
}

func newFinderFromConfig(cfg *CLIConfig) (lookup.Finder, error) {
	opts := []lookup.Option{
		lookup.WithInclude(cfg.Include),
		lookup.WithExclude(cfg.Exclude),
		lookup.WithExtendExclude(cfg.ExtendExclude),
		lookup.WithForceExclude(cfg.ForceExclude),
	}
	if !slices.Contains(cfg.ExplicitFlags, "exclude") &&
		cfg.Exclude == DefaultExcludePattern {
		opts = append(
			opts,
			lookup.WithIgnoreFiles(lookup.FileGitIgnore, lookup.FileIgnore),
		)
	}
	return lookup.NewFinder(opts...)
}

func formatCode(
	stdinFilename string,
	code string,
	lineCap int,
	check bool,
	quiet bool,
	ui *UI,
	stats StatsCollector,
) int {
	// TODO: AI slop
	displayName := "<string>"
	if stdinFilename != "" {
		displayName = stdinFilename
	}
	return formatBuffer(
		[]byte(code),
		displayName,
		lineCap,
		check,
		quiet,
		ui,
		stats,
		true,
	)
}

func formatStdin(
	stdinFilename string,
	lineCap int,
	check bool,
	quiet bool,
	stdin io.Reader,
	ui *UI,
	stats StatsCollector,
) int {
	// TODO: AI slop
	if stdin == nil {
		stdin = os.Stdin
	}
	buffer, err := io.ReadAll(stdin)
	if err != nil {
		ui.PrintError("Error reading stdin: %v", err)
		return ExitCodeInternalError
	}
	displayName := "-"
	if stdinFilename != "" {
		displayName = stdinFilename
	}
	return formatBuffer(
		buffer,
		displayName,
		lineCap,
		check,
		quiet,
		ui,
		stats,
		false,
	)
}

func formatBuffer(
	buffer []byte,
	displayName string,
	lineCap int,
	check bool,
	quiet bool,
	ui *UI,
	stats StatsCollector,
	isCode bool,
) int {
	// TODO: AI slop
	stats.StartFile(displayName)
	defer func() {
		stats.SnapshotMemory("end")
		stats.EndFile()
	}()
	stats.SnapshotMemory("start")

	stats.RecordReadBytes(int64(len(buffer)))
	stats.SnapshotMemory("after_read")

	if len(buffer) == 0 {
		return ExitCodeOK
	}

	buf := stream.NewHashedBufferFromString(displayName, string(buffer))
	inputHash := buf.Hash()

	errorCounter := ui.NewSyntaxErrorListener(displayName)

	file := dmlex.ParseStream(buf, errorCounter)
	tree := file.Parser.Dml()

	if cts, ok := file.Parser.GetTokenStream().(*antlr.CommonTokenStream); ok {
		stats.RecordTokensAndLines(cts.GetAllTokens())
	}
	stats.SnapshotMemory("after_parse")

	if len(errorCounter.Errors) > 0 || file.Parser.HasError() {
		if !check {
			_, _ = ui.Write(buffer)
		}
		return ExitCodeInternalError
	}

	if lineCap <= 0 {
		lineCap = DefaultMaxLineCap
	}

	var outBuf bytes.Buffer
	outHasher := stream.NewFileHasher()
	writer := io.MultiWriter(&outBuf, outHasher)

	var printErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				printErr = fmt.Errorf("printer panic: %v", r)
			}
		}()
		printErr = printer.PrinterWriteTo(tree, file, writer, lineCap)
	}()
	stats.SnapshotMemory("after_format")

	if printErr != nil {
		ui.PrintError("Error formatting %s: %v", displayName, printErr)
		if !check {
			_, _ = ui.Write(buffer)
		}
		return ExitCodeInternalError
	}

	outputHash := outHasher.Sum()
	if check {
		if outputHash != inputHash {
			if !isCode && !quiet {
				ui.WouldReformat(displayName)
			}
			return ExitCodeWouldReformat
		}
		return ExitCodeOK
	}

	_, _ = ui.Write(outBuf.Bytes())

	if !isCode && !quiet {
		if outputHash != inputHash {
			ui.ReformattedStdin(displayName)
		}
	}

	return ExitCodeOK
}

func printVersion(w io.Writer) {
	ver := Version
	commit := GitCommit
	date := BuildDate

	if (ver == "dev" || ver == "") && commit == "none" {
		if info, ok := debug.ReadBuildInfo(); ok {
			if info.Main.Version != "" && info.Main.Version != "(devel)" {
				ver = info.Main.Version
			}
			for _, setting := range info.Settings {
				switch setting.Key {
				case "vcs.revision":
					if len(setting.Value) > 7 {
						commit = setting.Value[:7]
					} else {
						commit = setting.Value
					}
				case "vcs.time":
					date = setting.Value
				}
			}
		}
	}

	if commit != "none" || date != "unknown" {
		_, _ = fmt.Fprintf(
			w,
			"dmlfmt version %s (commit: %s, built: %s)\n",
			ver,
			commit,
			date,
		)
	} else {
		_, _ = fmt.Fprintf(w, "dmlfmt version %s\n", ver)
	}
	_, _ = fmt.Fprintln(
		w,
		"License: GNU General Public License v2.0 only (GPL-2.0-only)",
	)
	_, _ = fmt.Fprintln(
		w,
		"This is free software: you are free to change and redistribute it.",
	)
	_, _ = fmt.Fprintln(
		w,
		"There is NO WARRANTY, to the extent permitted by law.",
	)
}

func main() {
	os.Exit(runWithIO(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
