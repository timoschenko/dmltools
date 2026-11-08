// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// ReplaceFileAtomically replaces targetPath with tempPath atomically
// and safely across OS platforms.
// If origMode is non-zero, it applies origMode to tempPath before moving it.
// To avoid overwriting issues (e.g. on Windows where replacing
// an existing file may fail), it first renames targetPath to a temporary
// backup file in the same directory, moves tempPath to targetPath,
// and removes the backup file on success.
// If moving tempPath fails, targetPath is restored from the temporary backup.
func ReplaceFileAtomically(
	targetPath, tempPath string,
	origMode os.FileMode,
) error {
	cleanTarget := filepath.Clean(targetPath)
	cleanTemp := filepath.Clean(tempPath)
	if cleanTarget == cleanTemp {
		return nil
	}

	if origMode != 0 {
		if err := os.Chmod(tempPath, origMode); err != nil {
			return fmt.Errorf("setting permissions on %s: %w", tempPath, err)
		}
	}

	dir := filepath.Dir(targetPath)
	base := filepath.Base(targetPath)

	_, err := os.Lstat(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Rename(tempPath, targetPath); err != nil {
				return fmt.Errorf(
					"moving %s to %s: %w",
					tempPath,
					targetPath,
					err,
				)
			}
			return nil
		}
		return fmt.Errorf("checking target file %s: %w", targetPath, err)
	}

	backupFile, err := os.CreateTemp(dir, "."+base+".*.tmp")
	if err != nil {
		return fmt.Errorf(
			"creating temp backup file for %s: %w",
			targetPath,
			err,
		)
	}
	targetTemp := backupFile.Name()
	_ = backupFile.Close()
	if err := os.Remove(targetTemp); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("preparing temp backup slot %s: %w", targetTemp, err)
	}

	if err := os.Rename(targetPath, targetTemp); err != nil {
		return fmt.Errorf(
			"backing up original file %s to %s: %w",
			targetPath,
			targetTemp,
			err,
		)
	}

	if err := os.Rename(tempPath, targetPath); err != nil {
		if rollbackErr := os.Rename(
			targetTemp,
			targetPath,
		); rollbackErr != nil {
			return fmt.Errorf(
				"moving %s to %s: %w (rollback also failed: %v)",
				tempPath, targetPath, err, rollbackErr,
			)
		}
		return fmt.Errorf("moving %s to %s: %w", tempPath, targetPath, err)
	}

	if err := os.Remove(targetTemp); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing temporary backup %s: %w", targetTemp, err)
	}

	return nil
}
