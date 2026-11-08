// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"bytes"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sampleSource = `// Package sample demonstrates baker parsing.
package sample

import (
	"fmt"
	"io"
)

// Greeter interface handles greeting.
type Greeter interface {
	Greet(name string) string
}

// IUser is the user interface.
type IUser interface {
	Hello() string
}

// User struct represents a user.
type User struct {
	ID   int
	Name string
}

// NewUser creates a new User.
func NewUser(id int, name string) *User {
	return &User{ID: id, Name: name}
}

// InitUser initializes a user.
func InitUser(u *User) {
	u.ID = 1
}

// Hello returns a greeting from User.
func (u *User) Hello() string {
	return fmt.Sprintf("Hello %s", u.Name)
}

// IDValue returns the user ID.
func (u User) IDValue() int {
	return u.ID
}

// StandaloneFunc is an orphan package-level function.
func StandaloneFunc() {
	_ = io.EOF
}
`

func TestParseFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "sample.go")
	if err := os.WriteFile(filePath, []byte(sampleSource), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	output, err := ParseFile(filePath)
	if err != nil {
		t.Fatalf("ParseFile returned error: %v", err)
	}

	if len(output.Entries) != 1 {
		t.Fatalf("expected 1 file entry, got %d", len(output.Entries))
	}

	entry := output.Entries[0]
	expectedFileName := "sample_user.go"
	if entry.FileName != expectedFileName {
		t.Errorf(
			"expected file name '%s', got '%s'",
			expectedFileName,
			entry.FileName,
		)
	}

	// Check structures in sample_user.go
	if len(entry.Group.Structures) != 1 {
		t.Fatalf(
			"expected 1 structure in %s, got %d",
			expectedFileName,
			len(entry.Group.Structures),
		)
	}
	st := entry.Group.Structures[0]
	if st.Name != "User" {
		t.Errorf("expected struct name 'User', got '%s'", st.Name)
	}
	if len(st.Methods) != 2 {
		t.Fatalf("expected 2 methods on User, got %d", len(st.Methods))
	}

	// Check interfaces in sample_user.go
	if len(entry.Group.Interfaces) != 1 {
		t.Fatalf(
			"expected 1 interface in %s, got %d",
			expectedFileName,
			len(entry.Group.Interfaces),
		)
	}
	if entry.Group.Interfaces[0].Name != "IUser" {
		t.Errorf(
			"expected interface 'IUser', got '%s'",
			entry.Group.Interfaces[0].Name,
		)
	}

	// Check functions in sample_user.go
	if len(entry.Group.Functions) != 2 {
		t.Fatalf(
			"expected 2 functions in %s, got %d",
			expectedFileName,
			len(entry.Group.Functions),
		)
	}

	// Check <uncategorized>
	if output.Uncategorized == nil {
		t.Fatal("expected non-nil Uncategorized")
	}
	if len(output.Uncategorized.Interfaces) != 1 {
		t.Fatalf(
			"expected 1 interface in <uncategorized>, got %d",
			len(output.Uncategorized.Interfaces),
		)
	}
	if output.Uncategorized.Interfaces[0].Name != "Greeter" {
		t.Errorf(
			"expected interface 'Greeter' in <uncategorized>, got '%s'",
			output.Uncategorized.Interfaces[0].Name,
		)
	}
	if len(output.Uncategorized.Functions) != 1 {
		t.Fatalf(
			"expected 1 function in <uncategorized>, got %d",
			len(output.Uncategorized.Functions),
		)
	}
	if output.Uncategorized.Functions[0].Name != "StandaloneFunc" {
		t.Errorf(
			"expected function 'StandaloneFunc' in <uncategorized>, got '%s'",
			output.Uncategorized.Functions[0].Name,
		)
	}

	// Leftovers check
	if !output.HasLeftovers() {
		t.Error(
			"expected HasLeftovers to be true due to orphan functions/interfaces",
		)
	}

	// Check JSON serialization and ordering of <uncategorized>
	data, err := EncodeJSON(output, false)
	if err != nil {
		t.Fatalf("EncodeJSON failed: %v", err)
	}

	jsonStr := string(data)
	userIdx := strings.Index(jsonStr, `"sample_user.go"`)
	uncatIdx := strings.Index(jsonStr, `"<uncategorized>"`)

	if userIdx == -1 {
		t.Error("missing 'sample_user.go' in JSON output")
	}
	if uncatIdx == -1 {
		t.Error("missing '<uncategorized>' in JSON output")
	}
	if uncatIdx <= userIdx {
		t.Errorf(
			"<uncategorized> (idx %d) must appear after sample_user.go (idx %d)",
			uncatIdx,
			userIdx,
		)
	}
}

func TestGenerateFiles(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "sample.go")
	if err := os.WriteFile(filePath, []byte(sampleSource), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	output, err := ParseFile(filePath)
	if err != nil {
		t.Fatalf("ParseFile returned error: %v", err)
	}

	outDir := filepath.Join(tempDir, "out")
	var stdout bytes.Buffer
	if err := GenerateFiles(output, outDir, &stdout); err != nil {
		t.Fatalf("GenerateFiles returned error: %v", err)
	}

	// Verify filename was printed to stdout
	printed := strings.TrimSpace(stdout.String())
	if printed != "sample_user.go" {
		t.Errorf("expected stdout to have 'sample_user.go', got '%s'", printed)
	}

	// Verify generated file exists and is valid Go code
	targetFile := filepath.Join(outDir, "sample_user.go")
	content, err := os.ReadFile(filepath.Clean(targetFile))
	if err != nil {
		t.Fatalf("failed to read generated file %s: %v", targetFile, err)
	}

	fset := token.NewFileSet()
	fileNode, err := parser.ParseFile(
		fset,
		targetFile,
		content,
		parser.ParseComments,
	)
	if err != nil {
		t.Fatalf(
			"generated file %s failed to parse: %v\nContent:\n%s",
			targetFile,
			err,
			string(content),
		)
	}

	if fileNode.Name.Name != "sample" {
		t.Errorf("expected package name 'sample', got '%s'", fileNode.Name.Name)
	}

	// Check that 'fmt' is imported (since Hello uses fmt.Sprintf), but 'io' is NOT imported
	hasFmt := false
	hasIO := false
	for _, imp := range fileNode.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		if path == "fmt" {
			hasFmt = true
		}
		if path == "io" {
			hasIO = true
		}
	}

	if !hasFmt {
		t.Error("expected generated file to import 'fmt'")
	}
	if hasIO {
		t.Error(
			"expected generated file NOT to import 'io' because it is only used in orphan StandaloneFunc",
		)
	}
}

const sampleWithConstsVars = `package sample

import (
	"fmt"
	"io"
)

// Suppress unused imports
var _ = fmt.Printf
var _ = io.EOF

const (
	UserDefaultID = 42
	UserNamePrefix = "user_"
)

const UnrelatedConst = 999

var UserStaticConfig = "cfg"
var UnrelatedVar = "orphan"

type User struct {
	ID int
}

func UserInit() {
}
`

func TestParseConstantsVariablesAndInit(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "sample_vars.go")
	if err := os.WriteFile(
		filePath,
		[]byte(sampleWithConstsVars),
		0600,
	); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	output, err := ParseFile(filePath)
	if err != nil {
		t.Fatalf("ParseFile returned error: %v", err)
	}

	if len(output.Entries) != 1 {
		t.Fatalf("expected 1 file entry, got %d", len(output.Entries))
	}

	entry := output.Entries[0]
	if len(entry.Group.Constants) != 1 ||
		entry.Group.Constants[0].Name != "UserDefaultID" {
		t.Errorf(
			"expected 1 constant block on User (UserDefaultID), got %v",
			entry.Group.Constants,
		)
	}

	if len(entry.Group.Variables) != 1 ||
		entry.Group.Variables[0].Name != "UserStaticConfig" {
		t.Errorf(
			"expected 1 variable on User (UserStaticConfig), got %v",
			entry.Group.Variables,
		)
	}

	foundUserInit := false
	for _, fn := range entry.Group.Functions {
		if fn.Name == "UserInit" {
			foundUserInit = true
			break
		}
	}
	if !foundUserInit {
		t.Errorf(
			"expected UserInit function in entry group, got %v",
			entry.Group.Functions,
		)
	}

	// Verify uncategorized constants and variables
	if len(output.Uncategorized.Constants) != 1 ||
		output.Uncategorized.Constants[0].Name != "UnrelatedConst" {
		t.Errorf(
			"expected 1 uncategorized const, got %v",
			output.Uncategorized.Constants,
		)
	}
	if len(output.Uncategorized.Variables) != 1 ||
		output.Uncategorized.Variables[0].Name != "UnrelatedVar" {
		t.Errorf(
			"expected 1 uncategorized var, got %v",
			output.Uncategorized.Variables,
		)
	}

	// Verify file generation parses cleanly
	outDir := filepath.Join(tempDir, "out")
	var stdout bytes.Buffer
	if err := GenerateFiles(output, outDir, &stdout); err != nil {
		t.Fatalf("GenerateFiles failed: %v", err)
	}

	targetFile := filepath.Join(outDir, "sample_vars_user.go")
	content, err := os.ReadFile(filepath.Clean(targetFile))
	if err != nil {
		t.Fatalf("failed to read generated file: %v", err)
	}

	fset := token.NewFileSet()
	if _, err := parser.ParseFile(
		fset,
		targetFile,
		content,
		parser.ParseComments,
	); err != nil {
		t.Fatalf(
			"generated file failed to parse: %v\nContent:\n%s",
			err,
			string(content),
		)
	}
}
