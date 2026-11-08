// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package fs

import (
	"bytes"
	"io"
	"math"
	"os"
)

const (
	growPercent     = 10
	maxGrowBytes    = 64 * 1024 // 64 KB
	alignBlockBytes = 4 * 1024  // 4 KB
)

// CalculateStreamBufferSize calculates the buffer allocation size by applying:
// 1) 10% growth placeholder
// 2) capping the 10% growth to 64KB
// 3) aligning up to a 4KB boundary using padding
func CalculateStreamBufferSize(size int) int {
	if size <= 0 {
		return alignBlockBytes
	}

	grow := min(size/growPercent, maxGrowBytes)
	if size > math.MaxInt-maxGrowBytes-alignBlockBytes {
		return math.MaxInt
	}
	allocSize := size + grow
	if rem := allocSize % alignBlockBytes; rem != 0 {
		allocSize += alignBlockBytes - rem
	}

	return allocSize
}

// SizeFromReader attempts to retrieve the stream or file size by casting r to known sizing interfaces.
func SizeFromReader(r io.Reader) (int, bool) {
	type fileStater interface {
		Stat() (os.FileInfo, error)
	}
	type size64er interface {
		Size() int64
	}
	type sizer interface {
		Size() int
	}

	switch v := r.(type) {
	case fileStater:
		if fi, err := v.Stat(); err == nil && fi.Mode().IsRegular() {
			if fi.Size() >= 0 && fi.Size() <= int64(math.MaxInt) {
				return int(fi.Size()), true
			}
		}
	case size64er:
		if s := v.Size(); s >= 0 && s <= int64(math.MaxInt) {
			return int(s), true
		}
	case sizer:
		if s := v.Size(); s >= 0 {
			return s, true
		}
	case *bytes.Buffer:
		if s := v.Len(); s >= 0 {
			return s, true
		}
	}
	return 0, false
}
