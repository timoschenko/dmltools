// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"runtime/pprof"
	"slices"
	"strings"
	"sync"

	"github.com/antlr4-go/antlr/v4"
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/fs"
	"github.com/timoschenko/dmltools/pkg/fs/lookup"
	"github.com/timoschenko/dmltools/pkg/fs/stream"
)

var (
	Version   = "dev"
	GitCommit = "none"
	BuildDate = "unknown"
)

func init() {
	antlr.CommonTokenFactoryDEFAULT = antlr.NewCommonTokenFactory(true)
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
			"dmltag version %s (commit: %s, built: %s)\n",
			ver,
			commit,
			date,
		)
	} else {
		_, _ = fmt.Fprintf(w, "dmltag version %s\n\n", ver)
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
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintln(w, previewWarning)
}

func readListFile(path string, stdin io.Reader) ([]string, error) {
	var r io.Reader
	if path == "-" {
		if stdin == nil {
			stdin = os.Stdin
		}
		r = stdin
	} else {
		f, err := os.Open(filepath.Clean(path))
		if err != nil {
			return nil, fmt.Errorf("opening list file %s: %w", path, err)
		}
		defer func() { _ = f.Close() }()
		r = f
	}

	var files []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			files = append(files, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading list file %s: %w", path, err)
	}
	return files, nil
}

func collectSources(
	cfg *CLIConfig,
	finder lookup.Finder,
	stdin io.Reader,
	stderr io.Writer,
) ([]string, error) {
	if cfg.DeviceSet {
		return collectDeviceFiles(cfg, finder, stderr)
	}

	var candidates []string

	if cfg.ListFile != "" {
		listCandidates, err := readListFile(cfg.ListFile, stdin)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, listCandidates...)
	}

	candidates = append(candidates, cfg.Sources...)

	if len(candidates) == 0 && cfg.Recurse {
		candidates = []string{"."}
	}

	if !cfg.Recurse {
		var nonDirs []string
		for _, cand := range candidates {
			fi, err := os.Stat(cand)
			if err == nil && fi.IsDir() {
				if !cfg.Quiet {
					_, _ = fmt.Fprintf(
						stderr,
						"dmltag: warning: %q is a directory and -R was not specified\n",
						cand,
					)
				}
				continue
			}
			nonDirs = append(nonDirs, cand)
		}
		candidates = nonDirs
	}

	if len(candidates) == 0 {
		return nil, nil
	}

	return finder.Collect(candidates)
}

type fileResult struct {
	path   string
	tags   []Tag
	hasErr bool
}

func parseAndExtractTags(
	path string,
	cfg *CLIConfig,
	stderr io.Writer,
) fileResult {
	data, err := os.ReadFile(
		filepath.Clean(path),
	) //nolint:gosec // user-specified path
	if err != nil {
		if !cfg.Quiet {
			_, _ = fmt.Fprintf(stderr, "dmltag: reading %s: %v\n", path, err)
		}
		return fileResult{path: path, hasErr: true}
	}

	buf := stream.NewHashedBufferFromString(path, string(data))
	errorCounter := dmlex.NewErrorListener(path, 0)
	file := dmlex.ParseStream(buf, errorCounter)
	tree := file.Parser.Dml()

	hasErr := false

	if unimpl := errorCounter.UnimplementedVersion(); unimpl != "" {
		hasErr = true
		if !cfg.Quiet {
			_, _ = fmt.Fprintf(
				stderr,
				"dmltag: warning: skipping unsupported DML %s file: %s\n",
				unimpl,
				path,
			)
		}
		return fileResult{path: path, hasErr: hasErr}
	} else if len(errorCounter.Errors) > 0 || file.Parser.HasError() {
		hasErr = true
		if !cfg.Quiet && cfg.Verbose {
			_, _ = fmt.Fprintf(
				stderr,
				"dmltag: warning: %s has %d syntax error(s)\n",
				path,
				len(errorCounter.Errors),
			)
		}
	}

	sourceLines := strings.Split(buf.GetRawString(), "\n")

	displayPath := path
	if cfg.TagRelative && cfg.OutputFile != "-" {
		targetDir := filepath.Dir(cfg.OutputFile)
		if rel, err := filepath.Rel(targetDir, path); err == nil {
			displayPath = rel
		}
	}
	displayPath = filepath.ToSlash(displayPath)

	tags := ExtractFileTags(
		tree,
		displayPath,
		sourceLines,
		ExtractOptions{
			EmitQualified: cfg.Qualified,
		},
	)

	return fileResult{
		path:   path,
		tags:   tags,
		hasErr: hasErr,
	}
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	cfg, err := parseArgs(args, stderr)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "dmltag: %v\n", err)
		return 2
	}

	if cfg.ShowHelp {
		printUsage(stdout)
		return 0
	}

	if cfg.ShowVersion {
		printVersion(stdout)
		return 0
	}

	if !cfg.Quiet {
		_, _ = fmt.Fprintln(stderr, previewWarning)
	}

	if cfg.WorkersErr != "" {
		_, _ = fmt.Fprintf(stderr, "dmltag: %s\n", cfg.WorkersErr)
		return 2
	}
	if cfg.DmlcFlagsErr != "" {
		_, _ = fmt.Fprintf(stderr, "dmltag: %s\n", cfg.DmlcFlagsErr)
		return 2
	}
	if cfg.DeviceSet && len(cfg.Sources) > 0 {
		_, _ = fmt.Fprintln(
			stderr,
			"dmltag: sources and --device cannot be passed simultaneously",
		)
		return 2
	}
	if cfg.FollowAllImport && !cfg.DeviceSet {
		_, _ = fmt.Fprintln(
			stderr,
			"dmltag: --follow-all-import requires --device",
		)
		return 2
	}

	if cfg.CPUProfile != "" {
		f, err := os.Create(filepath.Clean(cfg.CPUProfile))
		if err != nil {
			_, _ = fmt.Fprintf(
				stderr,
				"dmltag: creating CPU profile: %v\n",
				err,
			)
			return 1
		}
		defer func() { _ = f.Close() }()
		if err := pprof.StartCPUProfile(f); err != nil {
			_, _ = fmt.Fprintf(
				stderr,
				"dmltag: starting CPU profile: %v\n",
				err,
			)
			return 1
		}
		defer pprof.StopCPUProfile()
	}

	if cfg.MemProfile != "" {
		f, err := os.Create(filepath.Clean(cfg.MemProfile))
		if err != nil {
			_, _ = fmt.Fprintf(
				stderr,
				"dmltag: creating memory profile: %v\n",
				err,
			)
			return 1
		}
		defer func() {
			if err := pprof.WriteHeapProfile(f); err != nil {
				_, _ = fmt.Fprintf(
					stderr,
					"dmltag: writing memory profile: %v\n",
					err,
				)
			}
			_ = f.Close()
		}()
	}

	finderOpts := []lookup.Option{
		lookup.WithInclude(cfg.Include),
		lookup.WithExclude(cfg.Exclude),
		lookup.WithExtendExclude(cfg.ExtendExclude),
		lookup.WithForceExclude(cfg.ForceExclude),
	}
	if !slices.Contains(cfg.ExplicitFlags, "exclude") &&
		cfg.Exclude == DefaultExcludePattern {
		finderOpts = append(
			finderOpts,
			lookup.WithIgnoreFiles(lookup.FileGitIgnore, lookup.FileIgnore),
		)
	}

	finder, err := lookup.NewFinder(finderOpts...)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "dmltag: %v\n", err)
		return 1
	}

	sourceFiles, err := collectSources(cfg, finder, stdin, stderr)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "dmltag: %v\n", err)
		return 1
	}

	if len(sourceFiles) == 0 {
		if !cfg.Quiet {
			_, _ = fmt.Fprintln(stderr, "dmltag: no DML files found to process")
		}
		return 0
	}

	type fileTask struct {
		index int
		path  string
	}
	type indexedResult struct {
		index int
		res   fileResult
	}

	numWorkers := max(1, cfg.Workers)
	tasks := make(chan fileTask, len(sourceFiles))
	for i, sf := range sourceFiles {
		tasks <- fileTask{index: i, path: sf}
	}
	close(tasks)

	results := make(chan indexedResult, len(sourceFiles))
	var wg sync.WaitGroup

	for range numWorkers {
		wg.Add(1)
		go func() { //nolint:modernize // compatibility
			defer wg.Done()
			for task := range tasks {
				results <- indexedResult{
					index: task.index,
					res:   parseAndExtractTags(task.path, cfg, stderr),
				}
			}
		}()
	}

	wg.Wait()
	close(results)

	fileResults := make([]fileResult, len(sourceFiles))
	for r := range results {
		fileResults[r.index] = r.res
	}

	var allTags []Tag
	hadStrictError := false

	for _, res := range fileResults {
		if res.hasErr && cfg.Strict {
			hadStrictError = true
		}
		allTags = append(allTags, res.tags...)
	}

	if hadStrictError {
		return 1
	}

	var existing []string
	if cfg.Append && cfg.OutputFile != "-" {
		if existingFile, err := os.Open(
			filepath.Clean(cfg.OutputFile),
		); err == nil {
			existing, _ = ParseExistingTags(existingFile)
			_ = existingFile.Close()
		}
	}

	writeOpts := WriteOptions{
		Version:  Version,
		Sort:     cfg.Sort,
		ExCmd:    cfg.ExCmd,
		Append:   cfg.Append,
		Existing: existing,
	}

	if cfg.OutputFile == "-" {
		if err := WriteTags(stdout, allTags, writeOpts); err != nil {
			_, _ = fmt.Fprintf(stderr, "dmltag: writing tags: %v\n", err)
			return 1
		}
		return 0
	}

	outPath := filepath.Clean(cfg.OutputFile)
	outDir := filepath.Dir(outPath)
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		_, _ = fmt.Fprintf(
			stderr,
			"dmltag: creating output directory %s: %v\n",
			outDir,
			err,
		)
		return 1
	}

	tmpFile, err := os.CreateTemp(outDir, ".tags-*.tmp")
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "dmltag: creating temporary file: %v\n", err)
		return 1
	}
	tmpName := tmpFile.Name()

	if err := WriteTags(tmpFile, allTags, writeOpts); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpName)
		_, _ = fmt.Fprintf(stderr, "dmltag: writing tags: %v\n", err)
		return 1
	}
	_ = tmpFile.Close()

	if err := fs.ReplaceFileAtomically(outPath, tmpName, 0o644); err != nil {
		_ = os.Remove(tmpName)
		_, _ = fmt.Fprintf(
			stderr,
			"dmltag: saving tag file %s: %v\n",
			outPath,
			err,
		)
		return 1
	}

	return 0
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
