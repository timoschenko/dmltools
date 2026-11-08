// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

// Package fs provides file system utilities, safe atomic file writing,
// path normalization, and buffer sizing heuristics for DML tooling.
//
// # Key Capabilities
//
//   - Atomic File Replacement: [ReplaceFileAtomically] provides crash-safe,
//     cross-platform file replacement using temporary staging and backup
//     renaming to prevent partial writes and corruption.
//   - Path Normalization: [NormalizePath] produces uniform, forward-slash
//     file and directory paths across operating systems.
//   - Buffer Sizing: [CalculateStreamBufferSize] determines optimal I/O
//     buffer sizes based on file size thresholds.
//
// # Atomic Replacement Guarantees
//
// When formatting or rewriting source files in place, modifying a file directly
// risks data loss if the process is terminated mid-write. [ReplaceFileAtomically]
// writes output to a temporary staging file first, creates a temporary backup
// of the destination if it exists, renames the staged file over the target,
// and removes the backup upon success. If an error occurs, the original file
// is restored from backup.
//
// # Usage Example
//
//	target := "device.dml"
//	staged := "device.dml.tmp"
//
//	// Write new content to staged file...
//	if err := fs.ReplaceFileAtomically(target, staged, 0644); err != nil {
//	    log.Fatalf("failed to replace file: %v", err)
//	}
package fs
