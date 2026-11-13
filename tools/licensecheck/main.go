// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/licensecheck"
)

var defaultSkipDirs = map[string]bool{
	".git":        true,
	"venv":        true,
	"build":       true,
	".docs-cache": true,
	".vscode":     true,
}

func printf(w io.Writer, format string, a ...any) {
	_, _ = fmt.Fprintf(w, format, a...)
}

// NewScannerWithLicense creates a licensecheck.Scanner configured with
// built-in licenses and an explicit pattern for the expected license header.
func NewScannerWithLicense(
	expectedLicense string,
) (*licensecheck.Scanner, error) {
	licenses := licensecheck.BuiltinLicenses()
	customLicenses := []licensecheck.License{
		{
			ID:   expectedLicense,
			Type: licensecheck.ShareProgram,
			LRE: fmt.Sprintf(
				"SPDX-License-Identifier: %s",
				expectedLicense,
			),
		},
	}
	if expectedLicense == "GPL-2.0-only" {
		customLicenses = append(customLicenses, licensecheck.License{
			ID:   "GPL-2.0-or-later",
			Type: licensecheck.ShareProgram,
			LRE:  "SPDX-License-Identifier: GPL-2.0-or-later",
		})
	}
	customLicenses = append(customLicenses, licenses...)

	return licensecheck.NewScanner(customLicenses)
}

func isLicenseFile(relPath string) bool {
	base := filepath.Base(relPath)
	return base == "LICENSE" || base == "LICENSE.md" ||
		base == "LICENSE.txt" || base == "COPYING"
}

// IsGeneratedFile determines whether a file is generated and should be skipped.
func IsGeneratedFile(relPath string, data []byte) bool {
	base := filepath.Base(relPath)
	if strings.HasPrefix(base, "dml_parser_") ||
		base == "dml_lexer.go" ||
		base == "dml_base_listener.go" ||
		base == "dml_listener.go" ||
		strings.HasSuffix(base, "_string.go") ||
		base == "wasm_exec.js" {
		return true
	}

	limit := min(len(data), 1024)
	header := data[:limit]
	return bytes.Contains(header, []byte("DO NOT EDIT")) ||
		bytes.Contains(header, []byte("Code generated"))
}

// ExtractHeader extracts the leading comment/header block of a source file,
// avoiding false positive license matches inside function bodies or string literals.
func ExtractHeader(data []byte) []byte {
	if bytes.HasPrefix(data, []byte("package ")) {
		return nil
	}
	if before, _, found := bytes.Cut(data, []byte("\npackage ")); found {
		return before
	}
	if bytes.HasPrefix(data, []byte("grammar ")) {
		return nil
	}
	if before, _, found := bytes.Cut(data, []byte("\ngrammar ")); found {
		return before
	}
	limit := min(len(data), 4096)
	return data[:limit]
}

// CheckData scans source file bytes and reports whether it matches the expected license.
func CheckData(
	scanner *licensecheck.Scanner,
	data []byte,
	expectedLicense string,
) (bool, []string) {
	header := ExtractHeader(data)
	cov := scanner.Scan(header)
	matchedIDs := make([]string, 0, len(cov.Match))
	foundExpected := false
	hasDisallowed := false

	for _, m := range cov.Match {
		matchedIDs = append(matchedIDs, m.ID)
		if m.ID == expectedLicense {
			foundExpected = true
		} else if expectedLicense == "GPL-2.0-only" && m.ID == "GPL-2.0-or-later" {
			hasDisallowed = true
		}
	}

	return foundExpected && !hasDisallowed, matchedIDs
}

// CheckLicenseFileData validates that the root LICENSE file matches expected license
// with at least 90% coverage as required by pkg.go.dev.
func CheckLicenseFileData(
	scanner *licensecheck.Scanner,
	data []byte,
	expectedLicense string,
) (bool, float64, []string) {
	cov := scanner.Scan(data)
	matchedIDs := make([]string, 0, len(cov.Match))
	foundMatch := false

	// For GPL-2.0-only, google/licensecheck classifies the canonical GPLv2 license file text as "GPL-2.0".
	// "GPL-2.0-or-later" is strictly disallowed in this repository.
	expectedTarget := expectedLicense
	if expectedLicense == "GPL-2.0-only" {
		expectedTarget = "GPL-2.0"
	}

	for _, m := range cov.Match {
		matchedIDs = append(matchedIDs, m.ID)
		if m.ID == expectedTarget {
			foundMatch = true
		}
		if expectedLicense == "GPL-2.0-only" && m.ID == "GPL-2.0-or-later" {
			return false, cov.Percent, matchedIDs
		}
	}

	return foundMatch && cov.Percent >= 90.0, cov.Percent, matchedIDs
}

// Run executes the license check CLI logic with provided arguments and I/O streams.
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("licensecheck", flag.ContinueOnError)
	flags.SetOutput(stderr)

	rootDir := flags.String(
		"root",
		".",
		"Root directory of the repository to scan",
	)
	expectedLicense := flags.String(
		"license",
		"GPL-2.0-only",
		"Expected SPDX license identifier",
	)
	verbose := flags.Bool(
		"verbose",
		false,
		"Verbose output listing scanned and skipped files",
	)

	if err := flags.Parse(args); err != nil {
		return 2
	}

	scanner, err := NewScannerWithLicense(*expectedLicense)
	if err != nil {
		printf(stderr, "error: failed to create license scanner: %v\n", err)
		return 1
	}

	positional := flags.Args()
	var scannedFiles []string

	if len(positional) > 0 {
		scannedFiles = positional
	} else {
		err := filepath.WalkDir(
			*rootDir,
			func(path string, d fs.DirEntry, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}
				if d.IsDir() {
					if defaultSkipDirs[d.Name()] {
						return filepath.SkipDir
					}
					return nil
				}

				relPath, relErr := filepath.Rel(*rootDir, path)
				if relErr != nil {
					relPath = path
				}

				isGo := strings.HasSuffix(relPath, ".go")
				isG4 := relPath == filepath.Join(
					"grammar",
					"dml",
					"dml.g4",
				)
				isLic := isLicenseFile(relPath)

				if isGo || isG4 || isLic {
					scannedFiles = append(scannedFiles, path)
				}
				return nil
			},
		)
		if err != nil {
			printf(stderr, "error: walking root %s: %v\n", *rootDir, err)
			return 1
		}
	}

	var passedCount, failedCount, skippedCount int
	var failedFiles []string

	for _, filePath := range scannedFiles {
		var cleanPath string
		if filepath.IsAbs(filePath) {
			cleanPath = filepath.Clean(filePath)
		} else {
			candidate := filepath.Join(*rootDir, filePath)
			if _, err := os.Stat(candidate); err == nil {
				cleanPath = filepath.Clean(candidate)
			} else {
				cleanPath = filepath.Clean(filePath)
			}
		}

		// #nosec G304 -- CLI utility scanning repository files
		data, err := os.ReadFile(cleanPath)
		if err != nil {
			printf(
				stderr,
				"FAIL: cannot read file %s: %v\n",
				cleanPath,
				err,
			)
			failedCount++
			failedFiles = append(failedFiles, cleanPath)
			continue
		}

		relPath, relErr := filepath.Rel(*rootDir, cleanPath)
		if relErr != nil {
			relPath = cleanPath
		}

		if isLicenseFile(relPath) {
			ok, percent, matches := CheckLicenseFileData(
				scanner,
				data,
				*expectedLicense,
			)
			if ok {
				passedCount++
				if *verbose {
					printf(
						stdout,
						"PASS: %s (license file, coverage: %.1f%%, matches: %s)\n",
						relPath,
						percent,
						strings.Join(matches, ", "),
					)
				}
			} else {
				failedCount++
				failedFiles = append(failedFiles, relPath)
				printf(
					stderr,
					"FAIL: %s: license file check failed (coverage: %.1f%%, matches: %s, expected %s with >=90%% coverage)\n",
					relPath,
					percent,
					strings.Join(matches, ", "),
					*expectedLicense,
				)
			}
			continue
		}

		if IsGeneratedFile(relPath, data) {
			skippedCount++
			if *verbose {
				printf(
					stdout,
					"SKIP: %s (generated)\n",
					relPath,
				)
			}
			continue
		}

		ok, matches := CheckData(scanner, data, *expectedLicense)
		if ok {
			passedCount++
			if *verbose {
				printf(stdout, "PASS: %s\n", relPath)
			}
		} else {
			failedCount++
			failedFiles = append(failedFiles, relPath)
			if len(matches) > 0 {
				printf(
					stderr,
					"FAIL: %s: license mismatch (expected %s, found: %s)\n",
					relPath,
					*expectedLicense,
					strings.Join(matches, ", "),
				)
			} else {
				printf(
					stderr,
					"FAIL: %s: missing or unrecognized license header (expected %s)\n",
					relPath,
					*expectedLicense,
				)
			}
		}
	}

	if len(failedFiles) > 0 {
		printf(stderr, "\nFailed files (%d):\n", len(failedFiles))
		for _, f := range failedFiles {
			printf(stderr, "  - %s\n", f)
		}
	}

	printf(
		stdout,
		"License check summary: %d files checked, %d passed, %d failed, %d skipped generated\n",
		passedCount+failedCount,
		passedCount,
		failedCount,
		skippedCount,
	)

	if failedCount > 0 {
		return 1
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
