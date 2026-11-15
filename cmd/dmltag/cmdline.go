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

const (
	previewWarning = "WARNING: dmltag is a preview tool; " +
		"features and behavior may change without notice."

	DefaultExcludePattern = fs.DefaultExcludePattern
	DefaultIncludePattern = `\.dml$`
)

// CLIConfig holds parsed command-line configuration for dmltag.
type CLIConfig struct {
	OutputFile      string
	Append          bool
	ExCmd           ExCmdMode
	Sort            SortMode
	TagRelative     bool
	Qualified       bool
	Recurse         bool
	ListFile        string
	Device          string
	DeviceSet       bool
	FollowAllImport bool
	Exclude         string
	ExtendExclude   string
	ForceExclude    string
	Include         string
	IncludeSet      bool
	Workers         int
	WorkersErr      string
	Strict          bool
	Quiet           bool
	Verbose         bool
	ShowVersion     bool
	ShowHelp        bool
	CPUProfile      string
	MemProfile      string
	Sources         []string
	DmlcFlags       []string
	DmlcIncludeDirs []string
	DmlcFlagsErr    string
	ExplicitFlags   []string
}

func printUsage(out io.Writer) {
	_, _ = fmt.Fprintf(out, `Usage: dmltag [OPTIONS] [FILES/DIRS...]

  Generate ctags-compatible tag files for Device Modeling Language (DML).

Options:
  -f, -o FILE                     Write tags to FILE (use "-" for stdout) [default: tags].
  -a, --append                    Append tag entries to existing tag file.
  -n                              Format tag address using line numbers (--excmd=number).
  -N                              Format tag address using search patterns (--excmd=pattern).
  --excmd [number|pattern|combine]
                                  Address format for tag entries [default: pattern].
  --sort [yes|no|foldcase]        Sort mode for tag entries [default: yes].
  -u                              Equivalent to --sort=no.
  --tag-relative [yes|no]         File paths in tags file are relative to tag file directory.
  --qualified, --extras=+q        Emit extra qualified tag entries (<scope>.<name>).
  -R, --recurse                   Recursively scan directories for DML files.
  -L FILE                         Read source file names from FILE (use "-" for stdin).
  --device FILE                   Tag device DML file and imported files (pass dmlc flags after '--').
  --follow-all-import             Fail if an imported file cannot be found (requires --device).
  --exclude TEXT                  Regex matching files and directories to exclude [default: %s].
  --extend-exclude TEXT           Additional regex matching files and directories to exclude.
  --force-exclude TEXT            Regex matching files and directories to forcibly exclude.
  --include TEXT                  Regex matching files and directories to include [default: %s].
  -W, --workers INTEGER           Number of parallel parser workers [default: %d].
  --strict                        Exit with non-zero status on syntax errors or DML 1.2 files.
  -q, --quiet                     Suppress preview warning and non-critical messages.
  -v, --verbose                   Emit verbose progress messages.
  -V, --version                   Show version information and exit.
  -h, --help                      Show this help message and exit.
`, DefaultExcludePattern, DefaultIncludePattern, min(2, runtime.NumCPU()))

	if isDebugEnabled {
		_, _ = fmt.Fprint(out, `
Debug Options:
  --cpuprofile FILE               Write CPU profile to FILE.
  --memprofile FILE               Write memory profile to FILE.
`)
	}
	_, _ = fmt.Fprintf(out, "\n%s\n", previewWarning)
}

func parseExtrasFlag(current bool, s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return current
	}

	isAdd := false
	isRemove := false
	if strings.HasPrefix(s, "+") {
		isAdd = true
		s = s[1:]
	} else if strings.HasPrefix(s, "-") {
		isRemove = true
		s = s[1:]
	}

	hasQ := false
	if strings.Contains(s, "{qualified}") {
		hasQ = true
	} else {
		for part := range strings.SplitSeq(s, ",") {
			part = strings.TrimSpace(part)
			if part == "qualified" || strings.ContainsRune(part, 'q') {
				hasQ = true
				break
			}
		}
	}

	if !hasQ {
		if !isAdd && !isRemove {
			return false
		}
		return current
	}

	if isRemove {
		return false
	}
	return true
}

func takesArg(arg string) bool {
	switch arg {
	case "-f", "-o", "-L", "--device",
		"--exclude", "--extend-exclude", "--force-exclude", "--include",
		"-W", "--workers", "--excmd", "--sort", "--tag-relative",
		"--extras", "--extra", "--cpuprofile", "--memprofile":
		return true
	default:
		return false
	}
}

func parseArgs(args []string, stderr io.Writer) (*CLIConfig, error) {
	cfg := &CLIConfig{
		OutputFile: "tags",
		ExCmd:      ExCmdPattern,
		Sort:       SortYes,
		Exclude:    DefaultExcludePattern,
		Include:    DefaultIncludePattern,
		Workers:    min(2, runtime.NumCPU()),
	}

	dashIdx := -1
	for i, arg := range args {
		if arg == "--" {
			dashIdx = i
			break
		}
	}

	parseTarget := args
	var afterDash []string
	if dashIdx >= 0 {
		parseTarget = args[:dashIdx]
		afterDash = args[dashIdx+1:]
	}

	var (
		excmdStr     string
		sortStr      string
		tagRelStr    string
		extrasStr    string
		extraStr     string
		flagN        bool
		flagCapN     bool
		flagU        bool
		flagHelpLong bool
	)

	fs := flag.NewFlagSet("dmltag", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		printUsage(stderr)
	}

	// Output flags
	fs.StringVar(&cfg.OutputFile, "f", "tags", "Write tags to FILE")
	fs.StringVar(&cfg.OutputFile, "o", "tags", "Write tags to FILE")
	fs.BoolVar(&cfg.Append, "a", false, "Append tags to existing file")
	fs.BoolVar(&cfg.Append, "append", false, "Append tags to existing file")

	// Address and sort mode flags
	fs.BoolVar(&flagN, "n", false, "Line numbers address mode")
	fs.BoolVar(&flagCapN, "N", false, "Pattern address mode")
	fs.StringVar(
		&excmdStr,
		"excmd",
		"",
		"Address mode: number, pattern, or combine",
	)
	fs.StringVar(&sortStr, "sort", "", "Sort mode: yes, no, or foldcase")
	fs.BoolVar(&flagU, "u", false, "Equivalent to --sort=no")
	fs.StringVar(
		&tagRelStr,
		"tag-relative",
		"",
		"Relative file paths: yes or no",
	)

	// Scope and qualification flags
	fs.BoolVar(&cfg.Qualified, "qualified", false, "Emit qualified tag names")
	fs.StringVar(&extrasStr, "extras", "", "Extra tag features, e.g. +q")
	fs.StringVar(&extraStr, "extra", "", "Extra tag features, e.g. +q")

	// Traversal and input flags
	fs.BoolVar(&cfg.Recurse, "R", false, "Recurse into directories")
	fs.BoolVar(&cfg.Recurse, "recurse", false, "Recurse into directories")
	fs.StringVar(&cfg.ListFile, "L", "", "List of files to process")
	fs.StringVar(&cfg.Device, "device", "", "Device DML file")
	fs.BoolVar(
		&cfg.FollowAllImport,
		"follow-all-import",
		false,
		"Fail if import missing",
	)

	// Filter and concurrency flags
	fs.StringVar(
		&cfg.Exclude,
		"exclude",
		DefaultExcludePattern,
		"Exclude regex",
	)
	fs.StringVar(
		&cfg.ExtendExclude,
		"extend-exclude",
		"",
		"Extend exclude regex",
	)
	fs.StringVar(&cfg.ForceExclude, "force-exclude", "", "Force exclude regex")
	fs.StringVar(
		&cfg.Include,
		"include",
		DefaultIncludePattern,
		"Include regex",
	)
	fs.IntVar(&cfg.Workers, "W", cfg.Workers, "Parallel workers")
	fs.IntVar(&cfg.Workers, "workers", cfg.Workers, "Parallel workers")
	fs.BoolVar(&cfg.Strict, "strict", false, "Fail on syntax errors")

	// Verbosity and info flags
	fs.BoolVar(&cfg.Quiet, "q", false, "Quiet mode")
	fs.BoolVar(&cfg.Quiet, "quiet", false, "Quiet mode")
	fs.BoolVar(&cfg.Verbose, "v", false, "Verbose output")
	fs.BoolVar(&cfg.Verbose, "verbose", false, "Verbose output")
	fs.BoolVar(&cfg.ShowVersion, "V", false, "Show version")
	fs.BoolVar(&cfg.ShowVersion, "version", false, "Show version")
	fs.BoolVar(&cfg.ShowHelp, "h", false, "Show help")
	fs.BoolVar(&flagHelpLong, "help", false, "Show help")

	// Debug flags
	fs.StringVar(&cfg.CPUProfile, "cpuprofile", "", "CPU profile path")
	fs.StringVar(&cfg.MemProfile, "memprofile", "", "Memory profile path")

	// Rearrange arguments so flags after positional arguments can be parsed cleanly
	var flagsOnly []string
	var positional []string
	for i := 0; i < len(parseTarget); i++ {
		arg := parseTarget[i]
		if strings.HasPrefix(arg, "-") && arg != "-" {
			flagsOnly = append(flagsOnly, arg)
			if strings.Contains(arg, "=") {
				continue
			}
			if takesArg(arg) {
				if i+1 < len(parseTarget) {
					next := parseTarget[i+1]
					if !strings.HasPrefix(next, "-") || next == "-" {
						i++
						flagsOnly = append(flagsOnly, next)
					}
				}
			}
		} else {
			positional = append(positional, arg)
		}
	}

	if err := fs.Parse(flagsOnly); err != nil {
		return nil, err
	}

	if flagHelpLong {
		cfg.ShowHelp = true
	}

	fs.Visit(func(f *flag.Flag) {
		cfg.ExplicitFlags = append(cfg.ExplicitFlags, f.Name)
		switch f.Name {
		case "device":
			cfg.DeviceSet = true
		case "include":
			cfg.IncludeSet = true
		}
	})

	if flagN {
		cfg.ExCmd = ExCmdNumber
	}
	if flagCapN {
		cfg.ExCmd = ExCmdPattern
	}
	if excmdStr != "" {
		switch strings.ToLower(excmdStr) {
		case "number":
			cfg.ExCmd = ExCmdNumber
		case "pattern":
			cfg.ExCmd = ExCmdPattern
		case "combine":
			cfg.ExCmd = ExCmdCombine
		default:
			return nil, fmt.Errorf(
				"invalid value for --excmd: %q (expected number, pattern, or combine)",
				excmdStr,
			)
		}
	}

	if flagU {
		cfg.Sort = SortNo
	}
	if sortStr != "" {
		switch strings.ToLower(sortStr) {
		case "yes", "1", "true":
			cfg.Sort = SortYes
		case "no", "0", "false":
			cfg.Sort = SortNo
		case "foldcase", "2":
			cfg.Sort = SortFoldCase
		default:
			return nil, fmt.Errorf(
				"invalid value for --sort: %q (expected yes, no, or foldcase)",
				sortStr,
			)
		}
	}

	if tagRelStr != "" {
		switch strings.ToLower(tagRelStr) {
		case "yes", "always", "true", "1":
			cfg.TagRelative = true
		case "no", "never", "false", "0":
			cfg.TagRelative = false
		default:
			return nil, fmt.Errorf(
				"invalid value for --tag-relative: %q (expected yes or no)",
				tagRelStr,
			)
		}
	}

	if extrasStr != "" {
		cfg.Qualified = parseExtrasFlag(cfg.Qualified, extrasStr)
	}
	if extraStr != "" {
		cfg.Qualified = parseExtrasFlag(cfg.Qualified, extraStr)
	}

	workersExplicit := slices.Contains(cfg.ExplicitFlags, "W") ||
		slices.Contains(cfg.ExplicitFlags, "workers")
	if workersExplicit {
		if cfg.Workers < 1 {
			cfg.WorkersErr = fmt.Sprintf(
				"Invalid value for -W / --workers: %d is not in the range x>=1.",
				cfg.Workers,
			)
		}
	} else if envVal, ok := os.LookupEnv("DMLTAG_NUM_WORKERS"); ok {
		envVal = strings.TrimSpace(envVal)
		if envVal != "" {
			w, err := strconv.Atoi(envVal)
			if err != nil {
				cfg.WorkersErr = fmt.Sprintf(
					"Invalid value for DMLTAG_NUM_WORKERS: %q is not a valid integer.",
					envVal,
				)
			} else if w < 1 {
				cfg.WorkersErr = fmt.Sprintf(
					"Invalid value for DMLTAG_NUM_WORKERS: %d is not in the range x>=1.",
					w,
				)
			} else {
				cfg.Workers = w
			}
		}
	}

	cfg.Sources = append(cfg.Sources, positional...)

	if cfg.DeviceSet {
		if len(afterDash) > 0 {
			cfg.DmlcFlags = afterDash
			cfg.DmlcIncludeDirs, cfg.DmlcFlagsErr = parseDmlcFlags(afterDash)
		}
	} else if dashIdx >= 0 {
		cfg.Sources = append(cfg.Sources, afterDash...)
	}

	return cfg, nil
}
