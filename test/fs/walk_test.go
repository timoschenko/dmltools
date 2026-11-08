// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test_fs

import (
	"regexp"
	"testing"

	"github.com/timoschenko/dmltools/pkg/fs"
)

func TestNormalizePath_InternalAlias(t *testing.T) {
	tests := []struct {
		path  string
		isDir bool
		want  string
	}{
		{
			path:  `C:\path\file.dml`,
			isDir: false,
			want:  "C:/path/file.dml",
		},
		{
			path:  `C:\path\dir`,
			isDir: true,
			want:  "C:/path/dir/",
		},
		{
			path:  "foo/bar.dml",
			isDir: false,
			want:  "foo/bar.dml",
		},
	}

	for _, tt := range tests {
		got := fs.NormalizePath(tt.path, tt.isDir)
		if got != tt.want {
			t.Errorf(
				"normalizePath(%q, %v) = %q, want %q",
				tt.path,
				tt.isDir,
				got,
				tt.want,
			)
		}
	}
}

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		isDir bool
		want  string
	}{
		{
			name:  "windows backslash absolute file",
			path:  `C:\path\file.dml`,
			isDir: false,
			want:  "C:/path/file.dml",
		},
		{
			name:  "windows slash absolute file",
			path:  "C:/path/file.dml",
			isDir: false,
			want:  "C:/path/file.dml",
		},
		{
			name:  "windows drive with leading slash",
			path:  "/C:/path/file.dml",
			isDir: false,
			want:  "C:/path/file.dml",
		},
		{
			name:  "windows backslash directory",
			path:  `C:\path\dir`,
			isDir: true,
			want:  "C:/path/dir/",
		},
		{
			name:  "windows slash directory",
			path:  "C:/path/dir",
			isDir: true,
			want:  "C:/path/dir/",
		},
		{
			name:  "windows drive root backslash",
			path:  `C:\`,
			isDir: true,
			want:  "C:/",
		},
		{
			name:  "windows drive root slash",
			path:  "C:/",
			isDir: true,
			want:  "C:/",
		},
		{
			name:  "relative forward slash file",
			path:  "foo/bar.dml",
			isDir: false,
			want:  "foo/bar.dml",
		},
		{
			name:  "relative backslash file",
			path:  `foo\bar.dml`,
			isDir: false,
			want:  "foo/bar.dml",
		},
		{
			name:  "relative directory",
			path:  "src/dir",
			isDir: true,
			want:  "src/dir/",
		},
		{
			name:  "unix absolute directory",
			path:  "/abs/dir",
			isDir: true,
			want:  "/abs/dir/",
		},
		{
			name:  "unix absolute file",
			path:  "/abs/file.dml",
			isDir: false,
			want:  "/abs/file.dml",
		},
		{
			name:  "dot file",
			path:  ".",
			isDir: false,
			want:  "",
		},
		{
			name:  "dot dir",
			path:  ".",
			isDir: true,
			want:  "",
		},
		{
			name:  "empty file",
			path:  "",
			isDir: false,
			want:  "",
		},
		{
			name:  "empty dir",
			path:  "",
			isDir: true,
			want:  "",
		},
		{
			name:  "unix root",
			path:  "/",
			isDir: true,
			want:  "/",
		},
		{
			name:  "unc slash file",
			path:  "//server/share/file.dml",
			isDir: false,
			want:  "//server/share/file.dml",
		},
		{
			name:  "unc backslash file",
			path:  `\\server\share\file.dml`,
			isDir: false,
			want:  "//server/share/file.dml",
		},
		{
			name:  "unc backslash directory",
			path:  `\\server\share\dir`,
			isDir: true,
			want:  "//server/share/dir/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fs.NormalizePath(tt.path, tt.isDir)
			if got != tt.want {
				t.Errorf(
					"NormalizePath(%q, %v) = %q, want %q",
					tt.path,
					tt.isDir,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestDefaultExcludePattern(t *testing.T) {
	re := regexp.MustCompile(fs.DefaultExcludePattern)
	matches := []string{
		"/path/to/.git/config",
		"/project/.hg/store",
		"/build/linux64/bin",
		"/build/win64/bin",
		"/env/venv/lib",
		"/env/.venv/lib",
	}
	for _, m := range matches {
		if !re.MatchString(m) {
			t.Errorf("expected DefaultExcludePattern to match %q", m)
		}
	}
}
