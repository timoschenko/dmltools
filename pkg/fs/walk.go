// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package fs

import (
	"path/filepath"
	"strings"
)

// DefaultExcludePattern matches common directories to exclude during recursive searches.
const DefaultExcludePattern = `/(\.git|\.hg|linux64|win64|venv|\.venv)/`

func isWindowsVolume(p string) bool {
	if len(p) < 2 {
		return false
	}
	isLetter := (p[0] >= 'a' && p[0] <= 'z') || (p[0] >= 'A' && p[0] <= 'Z')
	return isLetter && p[1] == ':' && (len(p) == 2 || p[2] == '/')
}

// NormalizePath converts a path into a standardized forward-slash
// representation without injecting leading slashes onto relative paths or
// Windows volume identifiers.
// Directory paths have a trailing slash (e.g. "dir/", "/dir/", "C:/dir/").
// File paths do not have a trailing slash
// (e.g. "dir/file.dml", "/dir/file.dml", "C:/dir/file.dml").
func NormalizePath(path string, isDir bool) string {
	if path == "" {
		return ""
	}
	isUNC := strings.HasPrefix(path, "//") || strings.HasPrefix(path, `\\`)
	clean := filepath.Clean(path)
	slash := filepath.ToSlash(clean)
	slash = strings.ReplaceAll(slash, "\\", "/")
	if isUNC && !strings.HasPrefix(slash, "//") {
		slash = "/" + slash
	}
	if slash == "." {
		return ""
	}
	if len(slash) >= 3 && slash[0] == '/' && isWindowsVolume(slash[1:]) {
		slash = slash[1:]
	}
	if isDir && !strings.HasSuffix(slash, "/") {
		slash += "/"
	}
	return slash
}
