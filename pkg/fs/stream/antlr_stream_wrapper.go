// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package stream

type antlrStreamWrapper struct {
	filename string
	index    int
	size     int
}

// GetSourceName returns the file name or source identifier.
func (as *antlrStreamWrapper) GetSourceName() string {
	return as.filename
}

// Mark does nothing here as we have the entire buffer in memory.
func (as *antlrStreamWrapper) Mark() int {
	return -1
}

// Release does nothing here as we have the entire buffer in memory.
func (as *antlrStreamWrapper) Release(_ int) {
}

// Index returns the current byte offset into the input stream.
func (as *antlrStreamWrapper) Index() int {
	return as.index
}

// Size returns the total number of bytes in the input stream.
func (as *antlrStreamWrapper) Size() int {
	return as.size
}

// Seek moves the input point to the provided byte offset.
func (as *antlrStreamWrapper) Seek(index int) {
	as.index = max(0, min(index, as.size))
}
