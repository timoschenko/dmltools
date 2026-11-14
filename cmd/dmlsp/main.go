// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/timoschenko/dmltools/pkg/lsp"

	"go.lsp.dev/jsonrpc2"
)

const previewWarning = "WARNING: dmlsp is a preview tool; " +
	"features and behavior may change without notice."

var (
	Version   = "dev"
	GitCommit = "none"
	BuildDate = "unknown"
	Debug     = "false"
)

type readWriteCloser struct {
	io.Reader
	io.Writer
}

func (rwc *readWriteCloser) Close() error {
	if c, ok := rwc.Reader.(io.Closer); ok && rwc.Reader != os.Stdin {
		_ = c.Close()
	}
	if c, ok := rwc.Writer.(io.Closer); ok &&
		rwc.Writer != os.Stdout &&
		rwc.Writer != os.Stderr {
		_ = c.Close()
	}
	return nil
}

func printUsage(out io.Writer, fs *flag.FlagSet) {
	_, _ = fmt.Fprintf(out, "Usage: dmlsp [options]\n\n")
	_, _ = fmt.Fprintf(
		out,
		"Language Server Protocol (LSP) server for Device Modeling Language (DML).\n",
	)
	_, _ = fmt.Fprintf(
		out,
		"Communicates with IDE/editor clients via standard I/O (stdin/stdout).\n\n",
	)
	_, _ = fmt.Fprintf(out, "%s\n\n", previewWarning)
	_, _ = fmt.Fprintf(out, "Options:\n")
	fs.SetOutput(out)
	fs.PrintDefaults()
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
			"dmlsp version %s (commit: %s, built: %s)\n",
			ver,
			commit,
			date,
		)
	} else {
		_, _ = fmt.Fprintf(w, "dmlsp version %s\n\n", ver)
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

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("dmlsp", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var logFile string
	var verbose bool
	var showVersion bool
	var showHelp bool

	fs.StringVar(
		&logFile,
		"log-file",
		"",
		"Path to optional log file (logs to stderr by default)",
	)
	fs.StringVar(&logFile, "logfile", "", "Alias for --log-file")
	fs.BoolVar(&verbose, "v", false, "Enable verbose (debug) logging")
	fs.BoolVar(&verbose, "verbose", false, "Alias for -v")
	fs.BoolVar(
		&showVersion,
		"version",
		false,
		"Print version information and exit",
	)
	fs.BoolVar(&showVersion, "V", false, "Alias for --version")
	fs.BoolVar(&showHelp, "help", false, "Print help usage and exit")
	fs.BoolVar(&showHelp, "h", false, "Alias for --help")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if showHelp {
		printUsage(stdout, fs)
		return 0
	}

	if showVersion {
		printVersion(stdout)
		return 0
	}

	logWriter := stderr
	if logFile != "" {
		cleanPath := filepath.Clean(logFile)
		f, err := os.OpenFile( //nolint:gosec // user-specified log destination via CLI flag
			cleanPath,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0o600,
		)
		if err != nil {
			_, _ = fmt.Fprintf(
				stderr,
				"dmlsp: failed to open log file %q: %v\n",
				logFile,
				err,
			)
			return 1
		}
		defer func() {
			_ = f.Close()
		}()
		logWriter = f
	}

	logLevel := slog.LevelInfo
	if verbose || Debug == "true" {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(logWriter, &slog.HandlerOptions{
		Level: logLevel,
	}))

	logger.Info("starting dmlsp server",
		slog.String("version", Version),
		slog.String("commit", GitCommit),
	)
	logger.Warn(previewWarning)

	rwc := &readWriteCloser{Reader: stdin, Writer: stdout}
	stream := jsonrpc2.NewStream(rwc)
	server := lsp.NewServer(
		lsp.WithLogger(logger),
		lsp.WithVersion(Version),
	)

	if err := server.Serve(context.Background(), stream); err != nil {
		logger.Error("server error encountered", slog.Any("error", err))
		return server.ExitCode()
	}

	return server.ExitCode()
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
