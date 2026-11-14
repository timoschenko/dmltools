// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

//go:build windows

package lookup

import (
	"hash/fnv"
	"os"
	"path/filepath"
	"strings"
)

func fillSysMeta(path string, _ os.FileInfo, meta *FileMeta) {
	vol := filepath.VolumeName(filepath.Clean(path))
	if vol == "" {
		return
	}
	h := fnv.New64a()
	_, _ = h.Write([]byte(strings.ToUpper(vol)))
	meta.DeviceID = h.Sum64()
	meta.HasDevice = true
}

func osGeteuid() int {
	return -1
}
