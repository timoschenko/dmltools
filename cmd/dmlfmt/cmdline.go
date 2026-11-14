// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/timoschenko/dmltools/pkg/fs"
)

const DefaultMaxLineCap = 99

var (
	Version   = "dev"
	GitCommit = "none"
	BuildDate = "unknown"
)

const (
	DefaultExcludePattern = fs.DefaultExcludePattern
	DefaultIncludePattern = `\.dml$`

	ExitCodeOK            = 0
	ExitCodeWouldReformat = 1
	ExitCodeInternalError = 123

	// MaxFileSize defines the maximum permitted file size in bytes for disk files (50 MiB).
	MaxFileSize int64 = 50 << 20
)

// ColorMode specifies the coloring behavior for terminal outputs.
type ColorMode int

const (
	ColorAuto ColorMode = iota
	ColorAlways
	ColorNever
)

// CLIConfig represents the baseline command-line options for dmlfmt,
// adapted from black's CLI for Device Modeling Language (DML).
type CLIConfig struct {
	Code                string
	LineLength          int
	TargetVersion       string
	SkipSourceFirstLine bool
	Preview             bool
	Unstable            bool
	UnstableFeature     string
	Check               bool
	Diff                bool
	Color               bool
	NoColor             bool
	ColorMode           ColorMode
	LineRanges          string
	Fast                bool
	Safe                bool
	Exclude             string
	ExtendExclude       string
	ForceExclude        string
	StdinFilename       string
	Include             string
	Workers             int
	WorkersErr          string
	Quiet               bool
	Verbose             bool
	ShowVersion         bool
	NoCache             bool
	NerdStats           bool
	ShowHelp            bool

	CPUProfile string
	MemProfile string

	// Positional source arguments (files, directories, or "-" for stdin)
	Sources []string

	// Flags explicitly provided by the user on the command line
	ExplicitFlags []string

	// CodeSet indicates whether -c or --code was explicitly provided
	CodeSet bool
}

var implementedFlags = map[string]bool{
	"W":              true,
	"c":              true,
	"check":          true,
	"code":           true,
	"color":          true,
	"exclude":        true,
	"extend-exclude": true,
	"fast":           true,
	"force-exclude":  true,
	"h":              true,
	"help":           true,
	"include":        true,
	"l":              true,
	"line-length":    true,
	"nerd-stats":     true,
	"no-color":       true,
	"q":              true,
	"quiet":          true,
	"safe":           true,
	"stdin-filename": true,
	"v":              true,
	"verbose":        true,
	"version":        true,
	"workers":        true,
}

func isFlagImplemented(flagName string) bool {
	if isDebugEnabled {
		switch flagName {
		case "cpuprofile", "memprofile":
			return true
		}
	}
	return implementedFlags[flagName]
}

// -t, --target-version [1.2|1.4]  DML versions that should be supported by
//                                 dmlfmt's output. By default, dmlfmt will
//                                 infer target versions from the file header.
// -x, --skip-source-first-line    Skip the first line of the source code.
// --preview                       Enable potentially disruptive style changes
//                                 that may be added to dmlfmt's main
//                                 functionality in the next major release.
// --unstable                      Enable potentially disruptive style changes
//                                 that have known bugs or are not currently
//                                 expected to make it into the stable style
//                                 dmlfmt's next major release. Implies
//                                 --preview.
// --enable-unstable-feature TEXT  Enable specific features included in the
//                                 `+"`--unstable`"+` style. Requires `+"`--preview`"+`. No
//                                 compatibility guarantees are provided on the
//                                 behavior or existence of any unstable
//                                 features.
// -v, --verbose                   Emit messages about files that were not
//                                 changed or were ignored due to exclusion
//                                 patterns. If dmlfmt is using a configuration
//                                 file, a message detailing which one it is
//                                 using will be emitted.
// --no-cache                      Skip reading and writing the cache, forcing
//                                 dmlfmt to reformat all included files.
// --diff                          Don't write the files back, just output a
//                                 diff to indicate what changes dmlfmt would've
//                                 made. They are printed to stdout so
//                                 capturing them is simple.
//  Only
//                                 applies when --diff is given.
// --line-ranges START-END         When specified, dmlfmt will try its best to
//                                 only format these lines. This option can be
//                                 specified multiple times, and a union of the
//                                 lines will be formatted. Each range must be
//                                 specified as two integers connected by a
//                                 `+"`-`"+`: `+"`<START>-<END>`"+`. The `+"`<START>`"+` and
//                                 `+"`<END>`"+` integer indices are 1-based and
//                                 inclusive on both ends.

func printUsage(w io.Writer) {
	_, _ = fmt.Fprintf(w, `Usage: dmlfmt [OPTIONS] [SRC]...

  The uncompromising DML (Device Modeling Language) code formatter.

Options:
  -c, --code TEXT                 Format the code passed in as a string.
  -l, --line-length INTEGER       How many characters per line to allow.
                                  [default: %d]
  --check                         Don't write the files back, just return the
                                  status. Return code 0 means nothing would
                                  change. Return code 1 means some files would
                                  be reformatted. Return code 123 means there
                                  was an internal error.
  --color / --no-color            Show (or do not show) colored diff.
  --fast / --safe                 By default, dmlfmt performs formatting
                                  with a safety check that re-formats changes
                                  to verify parser and formatting stability
                                  (--safe). The --fast flag skips this safety
                                  check. [default: --safe]
  --exclude TEXT                  A regular expression that matches files and
                                  directories that should be excluded on
                                  recursive searches. An empty value means no
                                  paths are excluded. Use forward slashes for
                                  directories on all platforms (Windows, too).
                                  By default, dmlfmt also ignores all paths
                                  listed in .gitignore. Changing this value
                                  will override all default exclusions.
                                  [default: %s]
  --extend-exclude TEXT           Like --exclude, but adds additional files
                                  and directories on top of the default values
                                  instead of overriding them.
  --force-exclude TEXT            Like --exclude, but files and directories
                                  matching this regex will be excluded even
                                  when they are passed explicitly as
                                  arguments. This is useful when invoking
                                  dmlfmt programmatically on changed files,
                                  such as in a pre-commit hook or editor
                                  plugin.
  --stdin-filename TEXT           The name of the file when passing it through
                                  stdin. Useful to make sure dmlfmt will
                                  respect the --force-exclude option on some
                                  editors that rely on using stdin.
  --include TEXT                  A regular expression that matches files and
                                  directories that should be included on
                                  recursive searches. An empty value means all
                                  files are included regardless of the name.
                                  Use forward slashes for directories on all
                                  platforms (Windows, too). Overrides all
                                  exclusions, including from .gitignore and
                                  command line options. [default: %s]
  -W, --workers INTEGER           When dmlfmt formats multiple files, it may
                                  use a worker pool to speed up formatting.
                                  This option controls the number of parallel
                                  workers. This can also be specified via the
                                  DMLFMT_NUM_WORKERS environment variable.
                                  Defaults to the number of CPUs in the
                                  system. [x>=1]
  -q, --quiet                     Stop emitting all non-critical output. Error
                                  messages will still be emitted (which can be
                                  silenced by 2>/dev/null).
  -v, --verbose                   Emit messages about files that were not
                                  changed or were ignored due to exclusion
                                  patterns.
  --version                       Show the version and exit.
  --nerd-stats                    Record performance and memory metrics and
                                  print at exit.
  -h, --help                      Show this message and exit.
`, DefaultMaxLineCap, DefaultExcludePattern, DefaultIncludePattern)
	if isDebugEnabled {
		_, _ = fmt.Fprint(w, `
Debug Options:
  --cpuprofile FILE               Write CPU profile to FILE.
  --memprofile FILE               Write memory profile to FILE.
`)
	}
}

func parseArgs(args []string, stderr io.Writer) (*CLIConfig, error) {
	cfg := &CLIConfig{}
	fs := flag.NewFlagSet("dmlfmt", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		printUsage(stderr)
	}

	// TODO: Go's standard library `flag` package does not natively support:
	// 1. GNU-style combined short flags (e.g., `-qv`).
	// 2. Flags interspersed after positional arguments (e.g., `dmlfmt file.dml --check`).
	// 3. Repeatable flags accumulating into a slice without custom `flag.Value` types
	//    (needed for `--line-ranges`, `--enable-unstable-feature`, and `-t/--target-version`).
	// 4. Mutually exclusive boolean toggle pairs (`--color / --no-color`, `--fast / --safe`)
	//    where the last flag specified on the command line wins.

	// Code input (-c, --code)
	fs.StringVar(
		&cfg.Code,
		"code",
		"",
		"Format the code passed in as a string.",
	)
	fs.StringVar(&cfg.Code, "c", "", "Format the code passed in as a string.")

	// Line length (-l, --line-length)
	fs.IntVar(
		&cfg.LineLength,
		"line-length",
		DefaultMaxLineCap,
		"How many characters per line to allow.",
	)
	fs.IntVar(
		&cfg.LineLength,
		"l",
		DefaultMaxLineCap,
		"How many characters per line to allow.",
	)

	// Target DML version (-t, --target-version)
	// TODO: Support multiple target version flags and enum validation [1.2|1.4].
	fs.StringVar(
		&cfg.TargetVersion,
		"target-version",
		"",
		"DML version to target (1.2 or 1.4).",
	)
	fs.StringVar(
		&cfg.TargetVersion,
		"t",
		"",
		"DML version to target (1.2 or 1.4).",
	)

	// Skip source first line (-x, --skip-source-first-line)
	fs.BoolVar(
		&cfg.SkipSourceFirstLine,
		"skip-source-first-line",
		false,
		"Skip the first line of the source code.",
	)
	fs.BoolVar(
		&cfg.SkipSourceFirstLine,
		"x",
		false,
		"Skip the first line of the source code.",
	)

	// Preview & unstable style flags
	fs.BoolVar(&cfg.Preview, "preview", false, "Enable preview style changes.")
	fs.BoolVar(
		&cfg.Unstable,
		"unstable",
		false,
		"Enable unstable style changes. Implies --preview.",
	)
	// TODO: Support specifying --enable-unstable-feature multiple times.
	fs.StringVar(
		&cfg.UnstableFeature,
		"enable-unstable-feature",
		"",
		"Enable specific unstable features.",
	)

	// Check & diff modes
	fs.BoolVar(
		&cfg.Check,
		"check",
		false,
		"Don't write files back, just return status.",
	)
	fs.BoolVar(
		&cfg.Diff,
		"diff",
		false,
		"Don't write files back, just output a diff.",
	)

	// Color toggles (--color / --no-color)
	// TODO: Standard library `flag` cannot track ordering between --color and --no-color.
	fs.BoolVar(&cfg.Color, "color", false, "Show colored diff.")
	fs.BoolVar(&cfg.NoColor, "no-color", false, "Do not show colored diff.")

	// Line ranges (--line-ranges)
	// TODO: Support specifying --line-ranges START-END multiple times.
	fs.StringVar(
		&cfg.LineRanges,
		"line-ranges",
		"",
		"Format only the specified line range.",
	)

	// AST safety checks (--fast / --safe)
	fs.BoolVar(&cfg.Fast, "fast", false, "Turn off AST safety check.")
	fs.BoolVar(&cfg.Safe, "safe", true, "Explicitly enable AST safety check.")

	// File inclusion & exclusion patterns
	fs.StringVar(
		&cfg.Exclude,
		"exclude",
		DefaultExcludePattern,
		"Regex matching files and directories to exclude.",
	)
	fs.StringVar(
		&cfg.ExtendExclude,
		"extend-exclude",
		"",
		"Additional regex matching files and directories to exclude.",
	)
	fs.StringVar(
		&cfg.ForceExclude,
		"force-exclude",
		"",
		"Regex matching files and directories to forcibly exclude.",
	)
	fs.StringVar(
		&cfg.StdinFilename,
		"stdin-filename",
		"",
		"The name of the file when passing it through stdin.",
	)
	fs.StringVar(
		&cfg.Include,
		"include",
		DefaultIncludePattern,
		"Regex matching files and directories to include.",
	)

	// Concurrency (-W, --workers)
	// NOTE: Limit default workers to 2 due to ANLTR guts is mostly single threaded
	defaultWorkers := min(2, runtime.NumCPU())
	fs.IntVar(
		&cfg.Workers,
		"workers",
		defaultWorkers,
		"Number of parallel workers.",
	)
	fs.IntVar(&cfg.Workers, "W", defaultWorkers, "Number of parallel workers.")

	// Verbosity (-q, --quiet, -v, --verbose)
	fs.BoolVar(
		&cfg.Quiet,
		"quiet",
		false,
		"Stop emitting all non-critical output.",
	)
	fs.BoolVar(&cfg.Quiet, "q", false, "Stop emitting all non-critical output.")
	fs.BoolVar(
		&cfg.Verbose,
		"verbose",
		false,
		"Emit verbose diagnostic messages.",
	)
	fs.BoolVar(&cfg.Verbose, "v", false, "Emit verbose diagnostic messages.")

	// Meta & configuration flags
	fs.BoolVar(&cfg.ShowVersion, "version", false, "Show the version and exit.")
	fs.BoolVar(
		&cfg.NoCache,
		"no-cache",
		false,
		"Skip reading and writing the cache.",
	)
	fs.BoolVar(
		&cfg.NerdStats,
		"nerd-stats",
		false,
		"Record performance and memory metrics and print at exit.",
	)
	fs.BoolVar(&cfg.ShowHelp, "help", false, "Show help message and exit.")
	fs.BoolVar(&cfg.ShowHelp, "h", false, "Show help message and exit.")

	// Profiling flags (unlocked in debug mode)
	fs.StringVar(
		&cfg.CPUProfile,
		"cpuprofile",
		"",
		"Write CPU profile to file.",
	)
	fs.StringVar(
		&cfg.MemProfile,
		"memprofile",
		"",
		"Write memory profile to file.",
	)

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	fs.Visit(func(f *flag.Flag) {
		cfg.ExplicitFlags = append(cfg.ExplicitFlags, f.Name)
		if f.Name == "c" || f.Name == "code" {
			cfg.CodeSet = true
		}
	})

	workersExplicit := slices.Contains(cfg.ExplicitFlags, "W") ||
		slices.Contains(cfg.ExplicitFlags, "workers")
	if workersExplicit {
		if cfg.Workers < 1 {
			cfg.WorkersErr = fmt.Sprintf(
				"Invalid value for -W / --workers: %d is not in the range x>=1.",
				cfg.Workers,
			)
		}
	} else if envVal, ok := os.LookupEnv("DMLFMT_NUM_WORKERS"); ok {
		envVal = strings.TrimSpace(envVal)
		if envVal != "" {
			w, err := strconv.Atoi(envVal)
			if err != nil {
				cfg.WorkersErr = fmt.Sprintf(
					"Invalid value for DMLFMT_NUM_WORKERS: %q is not a valid integer.",
					envVal,
				)
			} else if w < 1 {
				cfg.WorkersErr = fmt.Sprintf(
					"Invalid value for DMLFMT_NUM_WORKERS: %d is not in the range x>=1.",
					w,
				)
			} else {
				cfg.Workers = w
			}
		}
	}

	if cfg.Unstable {
		cfg.Preview = true
	}
	if cfg.Quiet {
		cfg.Verbose = false
	}
	// Fast / safe toggles: track ordering so the last specified flag wins.
	lastSafeIdx := -1
	lastFastIdx := -1
	for i, arg := range args {
		switch arg {
		case "--safe", "-safe":
			lastSafeIdx = i
		case "--fast", "-fast":
			lastFastIdx = i
		}
	}
	if lastSafeIdx != -1 || lastFastIdx != -1 {
		if lastSafeIdx > lastFastIdx {
			cfg.Safe = true
			cfg.Fast = false
		} else {
			cfg.Safe = false
			cfg.Fast = true
		}
	} else {
		cfg.Safe = true
		cfg.Fast = false
	}
	// Color toggles: track ordering so the last specified flag wins.
	lastColorIdx := -1
	lastNoColorIdx := -1
	for i, arg := range args {
		switch arg {
		case "--color", "-color":
			lastColorIdx = i
		case "--no-color", "-no-color":
			lastNoColorIdx = i
		}
	}
	if lastColorIdx != -1 || lastNoColorIdx != -1 {
		if lastColorIdx > lastNoColorIdx {
			cfg.Color = true
			cfg.NoColor = false
			cfg.ColorMode = ColorAlways
		} else {
			cfg.Color = false
			cfg.NoColor = true
			cfg.ColorMode = ColorNever
		}
	} else if cfg.Color {
		cfg.ColorMode = ColorAlways
	} else if cfg.NoColor {
		cfg.ColorMode = ColorNever
	} else {
		cfg.ColorMode = ColorAuto
	}

	cfg.Sources = append(cfg.Sources, fs.Args()...)

	return cfg, nil
}
