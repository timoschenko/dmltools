// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/fs"
	"github.com/timoschenko/dmltools/pkg/fs/stream"
	"github.com/timoschenko/dmltools/pkg/printer"

	"github.com/antlr4-go/antlr/v4"
)

const copyTextIntoTokenOnParse = true

func init() {
	// NOTE: Enable "writing" all token content inside AST during parsing.
	//       It is expected that ANTLR runtime will call HashedBuffer.GetText()
	//       and save a pointer to string from the buffer instead of copying it.
	antlr.CommonTokenFactoryDEFAULT = antlr.NewCommonTokenFactory(
		copyTextIntoTokenOnParse,
	)
}

// InputJob represents a file read into memory and ready for parsing and formatting.
type InputJob struct {
	Path      string
	OrigMode  os.FileMode
	Size      int
	Buffer    *stream.HashedBuffer
	InputHash stream.Hash
}

// OutputJob represents the formatting result of an InputJob.
type OutputJob struct {
	Path          string
	OrigMode      os.FileMode
	InputHash     stream.Hash
	OutputHash    stream.Hash
	OutputBuffer  []byte
	SyntaxErrors  bool
	Err           error
	WouldReformat bool
	Ignored       bool
	IgnoreReason  string
}

// parseDmlStream parses input into a DmlFile and AST tree using an ErrorListener configured for path.
func parseDmlStream(
	input antlr.CharStream,
	path string,
) (*dmlex.File, antlr.ParseTree, *dmlex.ErrorListener) {
	// TODO: AI slop
	errorListener := dmlex.NewErrorListener(path, 0)
	file := dmlex.ParseStream(input, errorListener)
	tree := file.Parser.Dml()
	return file, tree, errorListener
}

// parseInput parses the job buffer, records token stats, and checks for unsupported versions or syntax errors.
// If the job cannot or should not be formatted, an early OutputJob is returned.
func parseInput(
	job InputJob,
	ui *UI,
	stats StatsCollector,
) (*dmlex.File, antlr.ParseTree, *OutputJob) {
	// TODO: AI slop
	startTime := time.Now()
	file, tree, errorCounter := parseDmlStream(job.Buffer, job.Path)
	stats.RecordTokenize(time.Since(startTime))

	if cts, ok := file.Parser.GetTokenStream().(*antlr.CommonTokenStream); ok {
		stats.RecordTokensAndLines(cts.GetAllTokens())
	}
	stats.SnapshotMemory("after_parse")

	if errorCounter.HasError(dmlex.EUNIMPL) {
		ver := errorCounter.UnimplementedVersion()
		if ver == "" {
			ver = "1.2"
		}
		return nil, nil, &OutputJob{
			Path:         job.Path,
			Ignored:      true,
			IgnoreReason: ver,
		}
	}

	if len(errorCounter.Errors) > 0 || file.Parser.HasError() {
		for _, e := range errorCounter.Errors {
			ui.SyntaxError(job.Path, e.Line, e.Column, e.Msg)
		}
		return nil, nil, &OutputJob{
			Path:         job.Path,
			SyntaxErrors: true,
		}
	}

	return file, tree, nil
}

// formatAST formats an AST tree to buf and calculates its stream hash.
// It catches any printer panics and formats them using panicPrefix.
func formatAST(
	buf *bytes.Buffer,
	tree antlr.ParseTree,
	file *dmlex.File,
	lineCap int,
	panicPrefix string,
) (stream.Hash, error) {
	// TODO: AI slop
	hasher := stream.NewFileHasher()
	writer := bufio.NewWriter(io.MultiWriter(buf, hasher))

	var printErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				printErr = fmt.Errorf("%s: %v", panicPrefix, r)
			}
		}()
		printErr = printer.PrinterWriteTo(tree, file, writer, lineCap)
	}()

	if printErr == nil {
		printErr = writer.Flush()
	}
	if printErr != nil {
		return stream.Hash{}, printErr
	}

	return hasher.Sum(), nil
}

// validateFormatting performs a safety check by re-parsing and re-formatting the output
// to ensure the formatting introduced no syntax errors and is idempotent.
func validateFormatting(
	path string,
	formatted *bytes.Buffer,
	expectedHash stream.Hash,
	lineCap int,
) error {
	// TODO: AI slop
	valInputStream := stream.NewHashedBufferFromString(path, formatted.String())
	valFile, valTree, valErrorCounter := parseDmlStream(valInputStream, path)

	if len(valErrorCounter.Errors) > 0 || valFile.Parser.HasError() {
		return fmt.Errorf(
			"validating file %s: syntax error during safety check",
			path,
		)
	}

	var valBuf bytes.Buffer
	valBuf.Grow(formatted.Len())
	valOutputHash, valPrintErr := formatAST(
		&valBuf,
		valTree,
		valFile,
		lineCap,
		"printer panic during safety check",
	)
	if valPrintErr != nil {
		return fmt.Errorf(
			"validating file %s: %w",
			path,
			valPrintErr,
		)
	}

	if valOutputHash != expectedHash {
		return fmt.Errorf(
			"validating file %s: additional formatting detected (formatting is not idempotent)",
			path,
		)
	}

	return nil
}

// processFormattingJob performs ANTLR parsing and printer layout formatting on an InputJob.
func processFormattingJob(
	job InputJob,
	lineCap int,
	safe bool,
	ui *UI,
	stats StatsCollector,
) OutputJob {
	// TODO: AI slop
	stats.StartFile(job.Path)
	defer func() {
		stats.SnapshotMemory("end")
		stats.EndFile()
	}()
	stats.SnapshotMemory("start")

	file, tree, earlyJob := parseInput(job, ui, stats)
	if earlyJob != nil {
		return *earlyJob
	}

	if lineCap <= 0 {
		lineCap = DefaultMaxLineCap
	}

	formatingStartTime := time.Now()
	var outBuf bytes.Buffer
	outBuf.Grow(job.Size)

	outputHash, printErr := formatAST(
		&outBuf,
		tree,
		file,
		lineCap,
		"printer panic",
	)
	stats.RecordPrinting(time.Since(formatingStartTime))
	stats.SnapshotMemory("after_format")

	if printErr != nil {
		return OutputJob{
			Path: job.Path,
			Err: fmt.Errorf(
				"formatting file %s: %w",
				job.Path,
				printErr,
			),
		}
	}

	wouldReformat := outputHash != job.InputHash
	if !wouldReformat {
		outBuf.Reset()
	} else if safe {
		validationStartTime := time.Now()
		defer func() {
			stats.RecordValidate(time.Since(validationStartTime))
		}()

		if err := validateFormatting(
			job.Path,
			&outBuf,
			outputHash,
			lineCap,
		); err != nil {
			return OutputJob{
				Path: job.Path,
				Err:  err,
			}
		}
	}

	return OutputJob{
		Path:          job.Path,
		OrigMode:      job.OrigMode,
		InputHash:     job.InputHash,
		OutputHash:    outputHash,
		OutputBuffer:  outBuf.Bytes(),
		WouldReformat: wouldReformat,
	}
}

// readInputJob opens a file, reads its contents into a HashingStream, and closes
// the file handle before returning.
func readInputJob(path string, stats StatsCollector) (InputJob, error) {
	openTime := time.Now()

	//nolint:gosec // path is cleaned during the walk&collect
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return InputJob{}, fmt.Errorf("opening file %s: %w", path, err)
	}
	defer func() {
		_ = f.Close()
	}()

	info, err := f.Stat()
	if err != nil {
		return InputJob{}, fmt.Errorf(
			"reading stat for file %s: %w",
			path,
			err,
		)
	}
	if info.Size() >= math.MaxInt {
		return InputJob{}, fmt.Errorf(
			"file %s exceeds maximum supported size (%d bytes)",
			path,
			info.Size(),
		)
	}

	size := int(info.Size())
	buf, err := stream.NewHashedBufferFromStream(path, f, size)
	if err != nil {
		return InputJob{}, fmt.Errorf("reading file %s: %w", path, err)
	}

	stats.RecordReadIo(time.Since(openTime))
	stats.RecordReadBytes(int64(size))

	return InputJob{
		Path:      path,
		OrigMode:  info.Mode().Perm(),
		Size:      size,
		Buffer:    buf,
		InputHash: buf.Hash(),
	}, nil
}

// runPipeline executes file discovery, single-worker IO reading/writing,
// and parallel processing (-W workers).
func runPipeline(
	cfg *CLIConfig,
	collector *FileCollector,
	ui *UI,
	stats StatsCollector,
) int {
	// TODO: AI Slop

	if cfg.Verbose {
		collector.SetOnIgnore(func(path string, isDir bool, filterName string) {
			ui.IgnoredFilter(path, isDir, filterName)
		})
	}

	startTime := time.Now()
	paths, err := collector.Collect(cfg.Sources)
	stats.RecordWalkIo(time.Since(startTime))
	if err != nil {
		ui.PrintError("Error: %v", err)
		return ExitCodeInternalError
	}
	if len(paths) == 0 {
		ui.PrintNoFiles()
		return ExitCodeOK
	}

	report := NewReport(cfg.Check)

	workers := cfg.Workers
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	processChan := make(chan InputJob, workers*2)
	writeChan := make(chan OutputJob, workers*2)

	// Processing Pool (CPU bound, workers concurrency)
	var processWg sync.WaitGroup
	for i := 0; i < max(workers, 1); i++ {
		processWg.Add(1)
		go func() {
			defer processWg.Done()
			for job := range processChan {
				writeChan <- processFormattingJob(
					job,
					cfg.LineLength,
					cfg.Safe,
					ui,
					stats,
				)
			}
		}()
	}

	go func() {
		processWg.Wait()
		close(writeChan)
	}()

	// Single IO Worker: interleaves file reading and atomic file writing
	hasErrors := false
	wouldReformat := false

	var (
		pathIdx       int
		nextJob       InputJob
		hasPending    bool
		processClosed bool
	)
	writeOut := writeChan

	for writeOut != nil {
		if !hasPending && pathIdx < len(paths) {
			path := paths[pathIdx]
			pathIdx++
			if cfg.Verbose {
				ui.Processing(path)
			}
			job, readErr := readInputJob(path, stats)
			if readErr != nil {
				ui.PrintError("Error: %v", readErr)
				hasErrors = true
				report.IncrementFailed()
				continue
			}
			nextJob = job
			hasPending = true
		}

		if !hasPending && pathIdx == len(paths) && !processClosed {
			close(processChan)
			processClosed = true
		}

		var processIn chan<- InputJob
		if hasPending {
			processIn = processChan
		}

		select {
		case job, ok := <-writeOut:
			if !ok {
				writeOut = nil
				continue
			}
			if job.SyntaxErrors {
				hasErrors = true
				report.IncrementFailed()
				continue
			}
			if job.Err != nil {
				ui.PrintError("Error: %v", job.Err)
				hasErrors = true
				report.IncrementFailed()
				continue
			}
			if job.Ignored {
				ui.Ignored(job.Path, job.IgnoreReason)
				if job.IgnoreReason == "1.2" || job.IgnoreReason == "" {
					report.IncrementIgnoredDml12()
				}
				continue
			}
			if !job.WouldReformat {
				if cfg.Verbose {
					ui.Unchanged(job.Path)
				}
				report.IncrementUnchanged()
				continue
			}
			if cfg.Check {
				ui.WouldReformat(job.Path)
				wouldReformat = true
				report.IncrementReformatted()
				continue
			}

			if err := flushReformattedFile(ui, &job, stats); err != nil {
				ui.PrintError("Error: %v", err)
				hasErrors = true
				report.IncrementFailed()
				continue
			}
			report.IncrementReformatted()
		case processIn <- nextJob:
			hasPending = false
			nextJob = InputJob{}
		}
	}

	ui.PrintReport(report)

	if hasErrors {
		return ExitCodeInternalError
	}
	if wouldReformat {
		return ExitCodeWouldReformat
	}
	return ExitCodeOK
}

func flushReformattedFile(
	ui *UI,
	job *OutputJob,
	stats StatsCollector,
) error {
	// TODO: AI Slop
	startTime := time.Now()

	if job == nil {
		panic("unreachable: flushReformattedFile job == nil")
	}

	dir := filepath.Dir(job.Path)
	tempFile, err := os.CreateTemp(dir, ".dmlfmt-*.tmp")
	if err != nil {
		return fmt.Errorf(
			"cannot create temp file in %s: %w",
			dir,
			err,
		)
	}
	closed := false
	replaced := false
	tempPath := tempFile.Name()
	defer func() {
		if !closed {
			if err := tempFile.Close(); err != nil {
				ui.PrintError("Error closing temp file: %v",
					err,
				)
			}
		}
		if !replaced {
			if err := os.Remove(tempPath); err != nil && !os.IsNotExist(err) {
				ui.PrintError("Error removing temp file: %v",
					err,
				)
			}
		}
	}()

	if _, err := tempFile.Write(job.OutputBuffer); err != nil {
		return fmt.Errorf(
			"cannot write to temp file %s: %w",
			tempPath,
			err,
		)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf(
			"cannot close temp file %s: %w",
			tempPath,
			err,
		)
	}
	closed = true

	if err := fs.ReplaceFileAtomically(
		job.Path, tempPath, job.OrigMode,
	); err != nil {
		return fmt.Errorf(
			"cannot replace file %s: %w",
			job.Path,
			err,
		)
	}
	replaced = true
	stats.RecordWriteIo(time.Since(startTime))

	ui.Reformatted(job.Path)

	return nil
}
