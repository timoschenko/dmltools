// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLI_VersionAndHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"--version"}, nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("--version exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "dmltag version") {
		t.Errorf(
			"--version output missing version string:\n%s",
			stdout.String(),
		)
	}
	if !strings.Contains(stdout.String(), "GPL-2.0-only") {
		t.Errorf("--version output missing license info:\n%s", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = run([]string{"--help"}, nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("--help exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "Usage: dmltag") {
		t.Errorf("--help output missing Usage string:\n%s", stdout.String())
	}
}

func TestCLI_Stdout(t *testing.T) {
	dir := t.TempDir()
	dmlPath := filepath.Join(dir, "device.dml")
	dmlContent := `dml 1.4;
device my_dev;
bank regs {
    register r0;
}
`
	if err := os.WriteFile(dmlPath, []byte(dmlContent), 0o600); err != nil {
		t.Fatalf("failed to write dml file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"-f", "-", dmlPath}, nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf(
			"run exit code = %d, want 0 (stderr: %s)",
			code,
			stderr.String(),
		)
	}

	out := stdout.String()
	if !strings.Contains(out, "!_TAG_FILE_FORMAT\t2") {
		t.Errorf("stdout missing pseudo-tag header:\n%s", out)
	}
	if !strings.Contains(out, "my_dev\t") {
		t.Errorf("stdout missing my_dev tag:\n%s", out)
	}
	if !strings.Contains(out, "regs\t") {
		t.Errorf("stdout missing regs tag:\n%s", out)
	}
	if !strings.Contains(out, "r0\t") {
		t.Errorf("stdout missing r0 tag:\n%s", out)
	}
}

func TestCLI_FileOutputAndAppend(t *testing.T) {
	dir := t.TempDir()
	file1 := filepath.Join(dir, "f1.dml")
	file2 := filepath.Join(dir, "f2.dml")
	tagsPath := filepath.Join(dir, "tags")

	if err := os.WriteFile(
		file1,
		[]byte("dml 1.4;\ndevice dev1;\n"),
		0o600,
	); err != nil {
		t.Fatalf("failed to write f1: %v", err)
	}
	if err := os.WriteFile(
		file2,
		[]byte("dml 1.4;\ndevice dev2;\n"),
		0o600,
	); err != nil {
		t.Fatalf("failed to write f2: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"-f", tagsPath, file1}, nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf(
			"run exit code = %d, want 0 (stderr: %s)",
			code,
			stderr.String(),
		)
	}

	data, err := os.ReadFile(tagsPath) //nolint:gosec // test file
	if err != nil {
		t.Fatalf("failed to read generated tags file: %v", err)
	}
	if !strings.Contains(string(data), "dev1\t") {
		t.Errorf("tags file missing dev1:\n%s", string(data))
	}
	if strings.Contains(string(data), "dev2\t") {
		t.Errorf("tags file should not contain dev2 yet")
	}

	// Append file2
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"-a", "-f", tagsPath, file2}, nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf(
			"append exit code = %d, want 0 (stderr: %s)",
			code,
			stderr.String(),
		)
	}

	data2, err := os.ReadFile(tagsPath) //nolint:gosec // test file
	if err != nil {
		t.Fatalf("failed to read appended tags file: %v", err)
	}
	if !strings.Contains(string(data2), "dev1\t") {
		t.Errorf("appended tags file missing dev1:\n%s", string(data2))
	}
	if !strings.Contains(string(data2), "dev2\t") {
		t.Errorf("appended tags file missing dev2:\n%s", string(data2))
	}
}

func TestCLI_RecurseAndListFile(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "sub")
	if err := os.MkdirAll(subDir, 0o750); err != nil {
		t.Fatalf("failed to create subDir: %v", err)
	}

	f1 := filepath.Join(dir, "f1.dml")
	f2 := filepath.Join(subDir, "f2.dml")
	listFile := filepath.Join(dir, "dml_list.txt")

	if err := os.WriteFile(
		f1,
		[]byte("dml 1.4;\ndevice dev1;\n"),
		0o600,
	); err != nil {
		t.Fatalf("failed to write f1: %v", err)
	}
	if err := os.WriteFile(
		f2,
		[]byte("dml 1.4;\ndevice dev2;\n"),
		0o600,
	); err != nil {
		t.Fatalf("failed to write f2: %v", err)
	}
	if err := os.WriteFile(
		listFile,
		[]byte(f1+"\n"+f2+"\n"),
		0o600,
	); err != nil {
		t.Fatalf("failed to write listFile: %v", err)
	}

	// Test -R
	var stdout, stderr bytes.Buffer
	code := run([]string{"-f", "-", "-R", dir}, nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf(
			"-R exit code = %d, want 0 (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	outR := stdout.String()
	if !strings.Contains(outR, "dev1\t") || !strings.Contains(outR, "dev2\t") {
		t.Errorf("-R failed to find both devices:\n%s", outR)
	}

	// Test -L
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"-f", "-", "-L", listFile}, nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf(
			"-L exit code = %d, want 0 (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	outL := stdout.String()
	if !strings.Contains(outL, "dev1\t") || !strings.Contains(outL, "dev2\t") {
		t.Errorf("-L failed to process both files:\n%s", outL)
	}
}

func TestCLI_AddressModes(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.dml")
	if err := os.WriteFile(
		f,
		[]byte("dml 1.4;\ndevice my_dev;\n"),
		0o600,
	); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Line number mode (-n)
	var stdout, stderr bytes.Buffer
	code := run([]string{"-f", "-", "-n", f}, nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("-n exit code = %d, want 0", code)
	}
	if !strings.Contains(
		stdout.String(),
		"my_dev\t"+filepath.ToSlash(f)+"\t2;\"\td",
	) {
		t.Errorf("-n unexpected tag line:\n%s", stdout.String())
	}

	// Combine mode (--excmd=combine)
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"-f", "-", "--excmd=combine", f}, nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("--excmd=combine exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "2;/^device my_dev;$/;\"\td") {
		t.Errorf("--excmd=combine unexpected tag line:\n%s", stdout.String())
	}
}

func TestCLI_ResilienceVsStrict(t *testing.T) {
	dir := t.TempDir()
	dml12File := filepath.Join(dir, "legacy.dml")
	if err := os.WriteFile(
		dml12File,
		[]byte("dml 1.2;\ndevice legacy;\n"),
		0o600,
	); err != nil {
		t.Fatalf("failed to write dml12: %v", err)
	}

	// Normal resilient run should skip unsupported file and return 0
	var stdout, stderr bytes.Buffer
	code := run([]string{"-f", "-", dml12File}, nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf(
			"resilient exit code = %d, want 0 (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if !strings.Contains(stderr.String(), "skipping unsupported DML 1.2 file") {
		t.Errorf("expected warning on stderr, got:\n%s", stderr.String())
	}
	if strings.Contains(stdout.String(), "legacy") {
		t.Errorf(
			"expected DML 1.2 file tags to be skipped, got:\n%s",
			stdout.String(),
		)
	}

	// Strict run should return 1
	stdout.Reset()
	stderr.Reset()
	code = run(
		[]string{"--strict", "-f", "-", dml12File},
		nil,
		&stdout,
		&stderr,
	)
	if code != 1 {
		t.Fatalf("strict exit code = %d, want 1", code)
	}
}

func TestCLI_Quiet(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "dev.dml")
	if err := os.WriteFile(
		f,
		[]byte("dml 1.4;\ndevice dev;\n"),
		0o600,
	); err != nil {
		t.Fatalf("failed to write dev file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"-q", "-f", "-", f}, nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("quiet run exit code = %d, want 0", code)
	}
	if strings.Contains(stderr.String(), "WARNING") {
		t.Errorf(
			"previewWarning should be silenced in quiet mode, got:\n%s",
			stderr.String(),
		)
	}
}

func TestCLI_DeviceMode(t *testing.T) {
	dir := t.TempDir()
	rootDml := filepath.Join(dir, "root.dml")
	subDml := filepath.Join(dir, "sub.dml")

	if err := os.WriteFile(
		rootDml,
		[]byte("dml 1.4;\ndevice root_dev;\nimport \"sub.dml\";\n"),
		0o600,
	); err != nil {
		t.Fatalf("failed to write root.dml: %v", err)
	}
	if err := os.WriteFile(
		subDml,
		[]byte("dml 1.4;\ntemplate sub_tpl {\n    param p = 1;\n}\n"),
		0o600,
	); err != nil {
		t.Fatalf("failed to write sub.dml: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"-f", "-", "--device", rootDml}, nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf(
			"device mode exit code = %d, want 0 (stderr: %s)",
			code,
			stderr.String(),
		)
	}

	out := stdout.String()
	if !strings.Contains(out, "root_dev\t") {
		t.Errorf("device mode output missing root_dev:\n%s", out)
	}
	if !strings.Contains(out, "sub_tpl\t") {
		t.Errorf("device mode output missing sub_tpl:\n%s", out)
	}
	if !strings.Contains(out, "p\t") {
		t.Errorf("device mode output missing sub_tpl.p:\n%s", out)
	}

	// Test missing import with --follow-all-import
	badRoot := filepath.Join(dir, "bad_root.dml")
	if err := os.WriteFile(
		badRoot,
		[]byte("dml 1.4;\ndevice bad_dev;\nimport \"nonexistent.dml\";\n"),
		0o600,
	); err != nil {
		t.Fatalf("failed to write bad_root.dml: %v", err)
	}

	stdout.Reset()
	stderr.Reset()
	code = run(
		[]string{"-f", "-", "--device", badRoot, "--follow-all-import"},
		nil,
		&stdout,
		&stderr,
	)
	if code == 0 {
		t.Fatalf(
			"expected error on missing import with --follow-all-import, got exit code 0",
		)
	}
}

func TestCLI_ExtrasFlag(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "dev.dml")
	if err := os.WriteFile(
		f,
		[]byte("dml 1.4;\ndevice dev;\nbank regs { register r0; }\n"),
		0o600,
	); err != nil {
		t.Fatalf("failed to write dev file: %v", err)
	}

	cases := []struct {
		flags       []string
		wantExtra   bool
		description string
	}{
		{[]string{"--extras=+q"}, true, "--extras=+q"},
		{[]string{"--extras=q"}, true, "--extras=q"},
		{[]string{"--extra=+fq"}, true, "--extra=+fq"},
		{[]string{"--extra=q"}, true, "--extra=q"},
		{[]string{"--extras=-q"}, false, "--extras=-q"},
		{
			[]string{"--qualified", "--extras=-q"},
			false,
			"--qualified with --extras=-q",
		},
		{[]string{"--extras=+{qualified}"}, true, "--extras=+{qualified}"},
	}

	for _, tc := range cases {
		var stdout, stderr bytes.Buffer
		args := append([]string{"-f", "-"}, tc.flags...)
		args = append(args, f)
		code := run(args, nil, &stdout, &stderr)
		if code != 0 {
			t.Fatalf(
				"%s exit code = %d, want 0 (stderr: %s)",
				tc.description,
				code,
				stderr.String(),
			)
		}
		hasQualified := strings.Contains(stdout.String(), "dev.regs.r0\t")
		if hasQualified != tc.wantExtra {
			t.Errorf(
				"%s: hasQualified = %v, want %v\nOutput:\n%s",
				tc.description,
				hasQualified,
				tc.wantExtra,
				stdout.String(),
			)
		}
	}
}

func TestCLI_SortNo_Determinism(t *testing.T) {
	dir := t.TempDir()
	files := make([]string, 0, 5)
	for i := range 5 {
		p := filepath.Join(dir, fmt.Sprintf("dev%d.dml", i))
		content := fmt.Sprintf(
			"dml 1.4;\ndevice dev%d;\nbank b%d { register r%d; }\n",
			i,
			i,
			i,
		)
		if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
			t.Fatalf("writing file: %v", err)
		}
		files = append(files, p)
	}

	var firstOutput string
	for iter := range 10 {
		var stdout, stderr bytes.Buffer
		args := append([]string{"-u", "-W", "4", "-f", "-"}, files...)
		code := run(args, nil, &stdout, &stderr)
		if code != 0 {
			t.Fatalf(
				"iter %d exit code = %d (stderr: %s)",
				iter,
				code,
				stderr.String(),
			)
		}
		if iter == 0 {
			firstOutput = stdout.String()
		} else if stdout.String() != firstOutput {
			t.Fatalf(
				"iter %d produced different tag order under --sort=no:\nGot:\n%s\nWant:\n%s",
				iter,
				stdout.String(),
				firstOutput,
			)
		}
	}
}
