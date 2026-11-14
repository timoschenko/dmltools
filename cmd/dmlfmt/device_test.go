// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

func writeTestDmlFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func readTestDmlFile(t *testing.T, path string) string {
	t.Helper()
	//nolint:gosec // test helper reads controlled temp path
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func runDmlfmtTest(
	t *testing.T,
	args ...string,
) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := runWithIO(args, strings.NewReader(""), &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func TestParseDmlcFlags(t *testing.T) {
	dirs, errMsg := parseDmlcFlags([]string{
		"-I", "inc1",
		"-Iinc2",
		"-I=inc3",
		"--include-dir=inc4",
		"--include-dir", "inc5",
		"-D", "FOO=1",
		"-DBAR",
		"-P", "header.h",
		"--dep", "out.d",
		"--dep-target", "out.o",
		"--warn", "WUNUSED",
		"--nowarn", "WEXPERIMENTAL",
		"--simics-api", "6",
		"--max-errors", "10",
		"--no-compat", "dml12_int",
		"--breaking-change", "bc",
		"--split-c-file", "2",
		"--mindeps",
		"-g",
		"--werror",
	})
	if errMsg != "" {
		t.Fatalf("unexpected error: %s", errMsg)
	}
	want := []string{"inc1", "inc2", "inc3", "inc4", "inc5"}
	if !slices.Equal(dirs, want) {
		t.Fatalf("got dirs %v, want %v", dirs, want)
	}

	errCases := [][]string{
		{"-I"},
		{"-I="},
		{"--include-dir"},
		{"--include-dir="},
		{"-D"},
		{"--simics-api"},
		{"bare.dml"},
		{"-"},
		{"-I", "inc", "extra.dml"},
	}
	for _, tc := range errCases {
		if _, msg := parseDmlcFlags(tc); msg == "" {
			t.Errorf("parseDmlcFlags(%v) expected error, got empty", tc)
		}
	}
}

func TestDeviceCLIValidation(t *testing.T) {
	dir := t.TempDir()
	devFile := filepath.Join(dir, "dev.dml")
	writeTestDmlFile(t, devFile, "dml 1.4;\ndevice dev;\n")

	cases := []struct {
		name     string
		args     []string
		wantCode int
		wantErr  string
	}{
		{
			name:     "device_with_src",
			args:     []string{"--device", devFile, devFile},
			wantCode: ExitCodeInternalError,
			wantErr:  "'SRC' and 'device'",
		},
		{
			name:     "device_with_stdin_dash",
			args:     []string{"--device", devFile, "-"},
			wantCode: ExitCodeInternalError,
			wantErr:  "'SRC' and 'device'",
		},
		{
			name:     "device_with_code",
			args:     []string{"--device", devFile, "-c", "dml 1.4;"},
			wantCode: ExitCodeInternalError,
			wantErr:  "'device' and 'code'",
		},
		{
			name:     "device_empty",
			args:     []string{"--device", ""},
			wantCode: ExitCodeInternalError,
			wantErr:  "'--device' requires a DML file path",
		},
		{
			name:     "device_missing_file",
			args:     []string{"--device", filepath.Join(dir, "nope.dml")},
			wantCode: ExitCodeInternalError,
			wantErr:  "nope.dml",
		},
		{
			name:     "device_is_dir",
			args:     []string{"--device", dir},
			wantCode: ExitCodeInternalError,
			wantErr:  "is a directory",
		},
		{
			name:     "follow_all_import_without_device",
			args:     []string{"--follow-all-import", devFile},
			wantCode: ExitCodeInternalError,
			wantErr:  "'--follow-all-import' requires '--device'",
		},
		{
			name:     "positional_after_double_dash",
			args:     []string{"--device", devFile, "--", "extra.dml"},
			wantCode: ExitCodeInternalError,
			wantErr:  "unexpected argument 'extra.dml'",
		},
		{
			name:     "dmlc_flag_without_double_dash",
			args:     []string{"--device", devFile, "-I", dir},
			wantCode: ExitCodeInternalError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			code, _, stderr := runDmlfmtTest(t, tc.args...)
			if code != tc.wantCode {
				t.Fatalf(
					"got exit code %d, want %d; stderr=%q",
					code,
					tc.wantCode,
					stderr,
				)
			}
			if tc.wantErr != "" && !strings.Contains(stderr, tc.wantErr) {
				t.Fatalf(
					"stderr %q does not contain %q",
					stderr,
					tc.wantErr,
				)
			}
		})
	}
}

func TestDeviceMode_FollowImportsAndSearchOrder(t *testing.T) {
	ws := t.TempDir()
	if err := os.MkdirAll(filepath.Join(ws, ".git"), 0o750); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	inc1 := filepath.Join(ws, "inc1")
	inc2 := filepath.Join(ws, "inc2")
	devDir := filepath.Join(ws, "devices", "mydev")

	unfmt := func(body string) string {
		return "dml 1.4;\nconstant   x   =   1;\n" + body
	}

	devPath := filepath.Join(devDir, "dev.dml")
	writeTestDmlFile(
		t,
		devPath,
		"dml 1.4;\n"+
			"device mydev;\n"+
			"constant   x   =   1;\n"+
			"import \"from_inc.dml\";\n"+
			"import \"prefer_ver.dml\";\n"+
			"import \"sub/child.dml\";\n",
	)

	// Implicit dml-builtins.dml inside inc1/1.4
	builtinsPath := filepath.Join(inc1, "1.4", "dml-builtins.dml")
	writeTestDmlFile(t, builtinsPath, unfmt(""))

	// inc1 vs inc2 precedence
	inc1File := filepath.Join(inc1, "from_inc.dml")
	inc2Shadow := filepath.Join(inc2, "from_inc.dml")
	writeTestDmlFile(t, inc1File, unfmt(""))
	writeTestDmlFile(t, inc2Shadow, unfmt(""))

	// 1.4/ subdir vs direct subdir inside inc1
	verPreferred := filepath.Join(inc1, "1.4", "prefer_ver.dml")
	directShadow := filepath.Join(inc1, "prefer_ver.dml")
	writeTestDmlFile(t, verPreferred, unfmt(""))
	writeTestDmlFile(t, directShadow, unfmt(""))

	// ./ and ../ relative import checks from sub/child.dml
	childPath := filepath.Join(devDir, "sub", "child.dml")
	writeTestDmlFile(t, childPath, unfmt(
		"import \"./sibling.dml\";\n"+
			"import \"../parent_rel.dml\";\n"+
			"import \"../fallback_inc.dml\";\n",
	))

	siblingPath := filepath.Join(devDir, "sub", "sibling.dml")
	writeTestDmlFile(t, siblingPath, unfmt(""))
	// Should NOT be resolved via ./ from inc1
	incSiblingShadow := filepath.Join(inc1, "sibling.dml")
	writeTestDmlFile(t, incSiblingShadow, unfmt(""))

	parentRelPath := filepath.Join(devDir, "parent_rel.dml")
	writeTestDmlFile(t, parentRelPath, unfmt(""))

	// ../fallback_inc.dml doesn't exist in devDir, falls back to inc1/1.4
	fallbackIncPath := filepath.Join(inc1, "1.4", "../fallback_inc.dml")
	// filepath.Join cleans .., so place at inc1/fallback_inc.dml and also
	// place one at inc1/other_fallback.dml referenced via fallback.dml.
	_ = fallbackIncPath
	fallbackInInc1 := filepath.Join(inc1, "fallback_inc.dml")
	writeTestDmlFile(t, fallbackInInc1, unfmt(""))

	// Wait: `../fallback_inc.dml` against `inc1/1.4/` cleans to
	// `inc1/fallback_inc.dml` cleanly! Let's also verify `check` mode first.
	checkCode, checkOut, checkErr := runDmlfmtTest(
		t,
		"--check",
		"--device", devPath,
		"--", "-I", inc1, "-I"+inc2, "-D", "SIMICS=6",
	)
	if checkCode != ExitCodeWouldReformat {
		t.Fatalf(
			"--check got %d, want %d; out=%q err=%q",
			checkCode,
			ExitCodeWouldReformat,
			checkOut,
			checkErr,
		)
	}

	code, stdout, stderr := runDmlfmtTest(
		t,
		"--device", devPath,
		"--", "-I", inc1, "-I="+inc2,
	)
	if code != ExitCodeOK {
		t.Fatalf("got code %d; stdout=%q stderr=%q", code, stdout, stderr)
	}

	shouldBeFormatted := []string{
		devPath,
		builtinsPath,
		inc1File,
		verPreferred,
		childPath,
		siblingPath,
		parentRelPath,
		fallbackInInc1,
	}
	for _, p := range shouldBeFormatted {
		got := readTestDmlFile(t, p)
		if strings.Contains(got, "   ") {
			t.Errorf("expected %s to be formatted, got:\n%s", p, got)
		}
	}

	shouldStayUnformatted := []string{
		inc2Shadow,
		directShadow,
		incSiblingShadow,
	}
	for _, p := range shouldStayUnformatted {
		got := readTestDmlFile(t, p)
		if !strings.Contains(got, "constant   x   =   1;") {
			t.Errorf("expected %s NOT to be formatted, got:\n%s", p, got)
		}
	}
}

func TestDeviceMode_DotSlashStrictRelative(t *testing.T) {
	ws := t.TempDir()
	if err := os.MkdirAll(filepath.Join(ws, ".git"), 0o750); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	inc := filepath.Join(ws, "inc")
	devPath := filepath.Join(ws, "dev.dml")
	writeTestDmlFile(
		t,
		devPath,
		"dml 1.4;\ndevice dev;\nimport \"./in_inc_only.dml\";\n",
	)
	writeTestDmlFile(
		t,
		filepath.Join(inc, "in_inc_only.dml"),
		"dml 1.4;\nconstant   x   =   1;\n",
	)

	// Without --follow-all-import: ./ must NOT match file in -I inc.
	code, _, _ := runDmlfmtTest(
		t,
		"--follow-all-import",
		"--device", devPath,
		"--", "-I", inc,
	)
	if code != ExitCodeInternalError {
		t.Fatalf("expected EIMPORT for ./ matching only -I dir, got %d", code)
	}
}

func TestDeviceMode_TopLevelIfElse(t *testing.T) {
	ws := t.TempDir()
	if err := os.MkdirAll(filepath.Join(ws, ".git"), 0o750); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	devPath := filepath.Join(ws, "dev.dml")
	writeTestDmlFile(
		t,
		devPath,
		"dml 1.4;\n"+
			"device dev;\n"+
			"#if (true) {\n"+
			"    import \"if_branch.dml\";\n"+
			"} #else #if (false) {\n"+
			"    import \"elif_branch.dml\";\n"+
			"} #else {\n"+
			"    import \"else_branch.dml\";\n"+
			"}\n",
	)

	branches := []string{
		filepath.Join(ws, "if_branch.dml"),
		filepath.Join(ws, "elif_branch.dml"),
		filepath.Join(ws, "else_branch.dml"),
	}
	for _, p := range branches {
		writeTestDmlFile(t, p, "dml 1.4;\nconstant   b   =   2;\n")
	}

	code, stdout, stderr := runDmlfmtTest(t, "--device", devPath)
	if code != ExitCodeOK {
		t.Fatalf("got code %d; stdout=%q stderr=%q", code, stdout, stderr)
	}

	for _, p := range branches {
		if got := readTestDmlFile(t, p); strings.Contains(got, "   ") {
			t.Errorf("branch file %s was not formatted:\n%s", p, got)
		}
	}
}

func TestDeviceMode_SemanticAndImportErrors(t *testing.T) {
	t.Run("EDEVICE_missing_device_tag", func(t *testing.T) {
		ws := t.TempDir()
		devPath := filepath.Join(ws, "dev.dml")
		orig := "dml 1.4;\nconstant   x   =   1;\n"
		writeTestDmlFile(t, devPath, orig)

		code, _, stderr := runDmlfmtTest(t, "--device", devPath)
		if code != ExitCodeInternalError {
			t.Fatalf("got %d, want %d; stderr=%q", code, ExitCodeInternalError,
				stderr)
		}
		if !strings.Contains(stderr, "EDEVICE") {
			t.Fatalf("expected EDEVICE in stderr, got %q", stderr)
		}
		if got := readTestDmlFile(t, devPath); got != orig {
			t.Fatalf("root file with EDEVICE should not be modified")
		}
	})

	t.Run("EDEVIMP_imported_device_file", func(t *testing.T) {
		ws := t.TempDir()
		devPath := filepath.Join(ws, "dev.dml")
		impPath := filepath.Join(ws, "imp.dml")
		impOrig := "dml 1.4;\ndevice   subdev;\n"
		writeTestDmlFile(
			t,
			devPath,
			"dml 1.4;\ndevice dev;\nimport \"imp.dml\";\n",
		)
		writeTestDmlFile(t, impPath, impOrig)

		code, _, stderr := runDmlfmtTest(t, "--device", devPath)
		if code != ExitCodeInternalError {
			t.Fatalf("got %d, want %d; stderr=%q", code, ExitCodeInternalError,
				stderr)
		}
		if !strings.Contains(stderr, "EDEVIMP") {
			t.Fatalf("expected EDEVIMP in stderr, got %q", stderr)
		}
		if got := readTestDmlFile(t, impPath); got != impOrig {
			t.Fatalf("imported file with EDEVIMP should not be modified")
		}
	})

	t.Run("unresolved_import_default_vs_follow_all", func(t *testing.T) {
		ws := t.TempDir()
		devPath := filepath.Join(ws, "dev.dml")
		presentPath := filepath.Join(ws, "present.dml")
		writeTestDmlFile(
			t,
			devPath,
			"dml 1.4;\ndevice dev;\n"+
				"import \"missing.dml\";\n"+
				"import \"present.dml\";\n",
		)
		writeTestDmlFile(t, presentPath, "dml 1.4;\nconstant   p   =   1;\n")

		// Default mode skips missing import and formats present.dml.
		code, stdout, stderr := runDmlfmtTest(t, "-v", "--device", devPath)
		if code != ExitCodeOK {
			t.Fatalf("default mode got %d; out=%q err=%q", code, stdout, stderr)
		}
		if !strings.Contains(
			stdout,
			"imported file not found, skipping: missing.dml",
		) {
			t.Fatalf("expected verbose skip message, got:\n%s", stdout)
		}
		if got := readTestDmlFile(
			t,
			presentPath,
		); strings.Contains(
			got,
			"   ",
		) {
			t.Fatalf("present.dml should be formatted")
		}

		// With --follow-all-import: reports EIMPORT and exits 123.
		code, _, stderr = runDmlfmtTest(
			t,
			"--follow-all-import",
			"--device", devPath,
		)
		if code != ExitCodeInternalError {
			t.Fatalf("--follow-all-import got %d, want %d; stderr=%q",
				code, ExitCodeInternalError, stderr)
		}
		if !strings.Contains(stderr, "EIMPORT") ||
			!strings.Contains(
				stderr,
				"cannot find file to import: missing.dml",
			) {
			t.Fatalf("expected EIMPORT error in stderr, got %q", stderr)
		}
	})
}

func TestDeviceMode_WorkspaceBoundaryAndInclude(t *testing.T) {
	base := t.TempDir()
	ws := filepath.Join(base, "repo")
	ext := filepath.Join(base, "external")
	if err := os.MkdirAll(filepath.Join(ws, ".git"), 0o750); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	devPath := filepath.Join(ws, "dev.dml")
	internalPath := filepath.Join(ws, "internal.dml")
	extAPath := filepath.Join(ext, "ext_a.dml")
	extBPath := filepath.Join(ext, "ext_b.dml")
	extBouncePath := filepath.Join(ws, "bounced.dml")

	unfmt := "dml 1.4;\nconstant   x   =   1;\n"
	writeTestDmlFile(
		t,
		devPath,
		"dml 1.4;\n"+
			"device dev;\n"+
			"import \"internal.dml\";\n"+
			"import \"ext_a.dml\";\n",
	)
	writeTestDmlFile(t, internalPath, unfmt)
	writeTestDmlFile(
		t,
		extAPath,
		"dml 1.4;\nconstant   a   =   1;\n"+
			"import \"ext_b.dml\";\n"+
			"import \"bounced.dml\";\n",
	)
	writeTestDmlFile(t, extBPath, unfmt)
	writeTestDmlFile(t, extBouncePath, unfmt)

	// 1. Default without --include: external files and bounced transitive
	// imports are skipped immediately.
	code, stdout, stderr := runDmlfmtTest(
		t,
		"-v",
		"--device", devPath,
		"--", "-I", ext, "-I", ws,
	)
	if code != ExitCodeOK {
		t.Fatalf("got %d; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if strings.Contains(readTestDmlFile(t, internalPath), "   ") {
		t.Errorf("internal.dml inside workspace should be formatted")
	}
	for _, p := range []string{extAPath, extBPath, extBouncePath} {
		if !strings.Contains(readTestDmlFile(t, p), "   ") {
			t.Errorf("%s should NOT be formatted without --include", p)
		}
	}

	// 2. With --include matching external files (using non-matching pattern for
	// internal.dml to verify workspace files still format even with --include).
	writeTestDmlFile(t, internalPath, unfmt)
	extPattern := regexp.QuoteMeta("ext_")
	code, stdout, stderr = runDmlfmtTest(
		t,
		"--device", devPath,
		"--include", extPattern,
		"--", "-I", ext, "-I", ws,
	)
	if code != ExitCodeOK {
		t.Fatalf("with --include got %d; out=%q err=%q", code, stdout, stderr)
	}
	for _, p := range []string{
		internalPath,
		extAPath,
		extBPath,
		extBouncePath,
	} {
		if strings.Contains(readTestDmlFile(t, p), "   ") {
			t.Errorf("%s should be formatted with --include %q", p, extPattern)
		}
	}
}

func TestDeviceMode_FiltersAndCycles(t *testing.T) {
	ws := t.TempDir()
	if err := os.MkdirAll(filepath.Join(ws, ".git"), 0o750); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}
	writeTestDmlFile(t, filepath.Join(ws, ".gitignore"), "gen/\nignored.dml\n")

	devPath := filepath.Join(ws, "dev.dml")
	aPath := filepath.Join(ws, "a.dml")
	bPath := filepath.Join(ws, "b.dml")
	genPath := filepath.Join(ws, "gen", "g.dml")
	ignoredPath := filepath.Join(ws, "ignored.dml")

	writeTestDmlFile(
		t,
		devPath,
		"dml 1.4;\ndevice dev;\n"+
			"import \"a.dml\";\n"+
			"import \"gen/g.dml\";\n"+
			"import \"ignored.dml\";\n",
	)
	writeTestDmlFile(
		t,
		aPath,
		"dml 1.4;\nconstant   a   =   1;\nimport \"b.dml\";\n",
	)
	writeTestDmlFile(
		t,
		bPath,
		"dml 1.4;\nconstant   b   =   2;\nimport \"a.dml\";\n",
	)
	writeTestDmlFile(t, genPath, "dml 1.4;\nconstant   g   =   3;\n")
	writeTestDmlFile(t, ignoredPath, "dml 1.4;\nconstant   i   =   4;\n")

	code, stdout, stderr := runDmlfmtTest(
		t,
		"--follow-all-import",
		"--device", devPath,
	)
	if code != ExitCodeOK {
		t.Fatalf("got %d; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if strings.Contains(readTestDmlFile(t, aPath), "   ") ||
		strings.Contains(readTestDmlFile(t, bPath), "   ") {
		t.Errorf("cyclic files a.dml and b.dml should be formatted")
	}
	if !strings.Contains(readTestDmlFile(t, genPath), "   ") {
		t.Errorf("gitignored dir file %s should NOT be formatted", genPath)
	}
	if !strings.Contains(readTestDmlFile(t, ignoredPath), "   ") {
		t.Errorf("gitignored file %s should NOT be formatted", ignoredPath)
	}

	// --force-exclude on root device file skips everything cleanly.
	code, stdout, stderr = runDmlfmtTest(
		t,
		"--force-exclude", `dev\.dml$`,
		"--device", devPath,
	)
	if code != ExitCodeOK {
		t.Fatalf("force-exclude got %d; out=%q err=%q", code, stdout, stderr)
	}
	if !strings.Contains(stdout, "No DML files are present") {
		t.Fatalf("expected no-files message on stdout, got %q", stdout)
	}

	// --exclude skips imported files matching pattern without disabling
	// explicit non-excluded files.
	excludedImp := filepath.Join(ws, "excluded_imp.dml")
	writeTestDmlFile(
		t,
		devPath,
		"dml 1.4;\ndevice dev;\nimport \"excluded_imp.dml\";\n",
	)
	writeTestDmlFile(t, excludedImp, "dml 1.4;\nconstant   e   =   1;\n")
	code, stdout, stderr = runDmlfmtTest(
		t,
		"--exclude", `excluded_imp\.dml$`,
		"--device", devPath,
	)
	if code != ExitCodeOK {
		t.Fatalf("--exclude got %d; out=%q err=%q", code, stdout, stderr)
	}
	if !strings.Contains(readTestDmlFile(t, excludedImp), "   ") {
		t.Errorf("excluded_imp.dml should NOT be formatted under --exclude")
	}
}
