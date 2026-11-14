// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

// Package lookup provides multipurpose file discovery, hierarchical
// configuration lookup, recursive source collection, path filtering, and
// stat-validated security caching.
//
// # SOLID Architecture
//
// Responsibilities are decomposed into focused interfaces:
//
//   - [Cache]: Manages stat-fingerprint-validated filesystem metadata, directory
//     boundary resolution (VCS roots like .git/.hg, mount-point boundaries,
//     directory UID ownership checks, ceiling directories, and max depth), and
//     candidate file security validation (regular-file checks, size limits,
//     and symlink containment).
//   - [Filter]: Evaluates path inclusion and exclusion rules (--include,
//     --exclude, --extend-exclude, --force-exclude) normalized via
//     [fs.NormalizePath].
//   - [IgnoreMatcher]: Parses, caches, and evaluates hierarchical ignore rule
//     files (.gitignore, .ignore) with standard glob, directory, and negation
//     semantics.
//   - [Finder]: Composes [Cache], [Filter], and [IgnoreMatcher] to support both
//     upward configuration discovery ([Finder.FindChain], [Finder.FindNearest])
//     and downward recursive source collection ([Finder.Collect]).
//
// # Usage Example
//
//	finder, err := lookup.NewFinder(
//		lookup.WithInclude(`\.dml$`),
//		lookup.WithExclude(fs.DefaultExcludePattern),
//	)
//	if err != nil {
//		log.Fatalf("init lookup finder: %v", err)
//	}
//
//	sources, err := finder.Collect([]string{"."})
//	if err != nil {
//		log.Fatalf("collect sources: %v", err)
//	}
//	configs, err := finder.FindChain(
//		"path/to/device.dml",
//		lookup.FileEditorConfig,
//		lookup.FileGitIgnore,
//	)
package lookup
