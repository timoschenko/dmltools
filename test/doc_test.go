// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test

import (
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestPackageDocumentation ensures all Go packages in pkg/ have package-level
// documentation comments to support 'go doc'.
func TestPackageDocumentation(t *testing.T) {
	root := filepath.Join("..", "pkg")
	if _, err := os.Stat(root); err != nil {
		// When tests are run from repository root
		root = "pkg"
	}

	var pkgDirs []string
	err := filepath.Walk(
		root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				return nil
			}

			// Check if directory contains any .go files
			entries, err := os.ReadDir(path)
			if err != nil {
				return err
			}

			hasGoFiles := false
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") &&
					!strings.HasSuffix(entry.Name(), "_test.go") {
					hasGoFiles = true
					break
				}
			}

			if hasGoFiles {
				pkgDirs = append(pkgDirs, path)
			}
			return nil
		},
	)

	if err != nil {
		t.Fatalf("failed to walk pkg directory: %v", err)
	}

	if len(pkgDirs) == 0 {
		t.Fatalf("no Go packages found under %s", root)
	}

	for _, dir := range pkgDirs {
		relDir, _ := filepath.Rel(root, dir)
		t.Run(relDir, func(t *testing.T) {
			pkg, err := build.ImportDir(dir, 0)
			if err != nil {
				t.Fatalf("failed to inspect package at %s: %v", dir, err)
			}

			doc := strings.TrimSpace(pkg.Doc)
			if doc == "" {
				t.Errorf(
					"package %s at %s is missing package-level documentation (doc.go or package comment)",
					pkg.Name,
					dir,
				)
			}

			// Verify that the doc starts with "Package <name>"
			expectedPrefix := "Package " + pkg.Name
			if !strings.HasPrefix(doc, expectedPrefix) {
				t.Errorf(
					"package doc for %s should start with %q, but got: %q",
					pkg.Name,
					expectedPrefix,
					doc,
				)
			}
		})
	}
}

// TestGoDocCommand verifies that 'go doc' executes cleanly on all pkg/ packages.
func TestGoDocCommand(t *testing.T) {
	pkgs := []string{
		"./pkg/dmlex",
		"./pkg/fs",
		"./pkg/fs/lookup",
		"./pkg/fs/stream",
		"./pkg/printer",
		"./pkg/printer/dml_defaults",
		"./pkg/printer/format",
		"./pkg/printer/layout",
	}

	for _, pkgPath := range pkgs {
		t.Run(pkgPath, func(t *testing.T) {
			// #nosec G204 -- trusted test execution of go doc
			c := exec.Command("go", "doc", pkgPath)
			if _, err := os.Stat("pkg"); err == nil {
				c.Dir = "."
			} else {
				c.Dir = ".."
			}
			out, err := c.CombinedOutput()
			if err != nil {
				t.Fatalf(
					"go doc %s failed: %v\nOutput:\n%s",
					pkgPath,
					err,
					string(out),
				)
			}
			outStr := string(out)
			if !strings.Contains(outStr, "Package ") {
				t.Errorf(
					"go doc %s output missing package summary:\n%s",
					pkgPath,
					outStr,
				)
			}
		})
	}
}
