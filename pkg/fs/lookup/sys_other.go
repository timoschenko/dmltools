// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

//go:build !unix && !windows

package lookup

import "os"

func fillSysMeta(_ string, _ os.FileInfo, _ *FileMeta) {}

func osGeteuid() int {
	return -1
}
