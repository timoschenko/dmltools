// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

//go:build unix

package lookup

import (
	"os"
	"syscall"
)

func fillSysMeta(_ string, info os.FileInfo, meta *FileMeta) {
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok || st == nil {
		return
	}
	// #nosec G115 -- st.Dev on Darwin is int32; preserve raw bit pattern as ID
	//nolint:unconvert // st.Dev is int32 on Darwin and uint64 on Linux
	meta.DeviceID = uint64(st.Dev)
	meta.Inode = st.Ino
	meta.UID = st.Uid
	meta.HasDevice = true
	meta.HasInode = true
	meta.HasUID = true
}

func osGeteuid() int {
	return os.Geteuid()
}
