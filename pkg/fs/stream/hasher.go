// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package stream

import (
	"hash"
	"hash/crc64"
	"hash/fnv"
	"io"
)

var crc64Table = crc64.MakeTable(crc64.ISO)

// Hash represents a composite checksum consisting of
// a 64-bit FNV-1a hash, a 64-bit byte size, and a 64-bit CRC64 checksum.
type Hash struct {
	FNV64 uint64
	Size  uint64
	CRC64 uint64
}

// Hasher computes a composite StreamHash across written bytes.
type Hasher struct {
	fnv  hash.Hash64
	crc  hash.Hash64
	size uint64
}

// NewFileHasher creates a new StreamHasher using FNV-1a and CRC64 (ISO).
func NewFileHasher() *Hasher {
	return &Hasher{
		fnv: fnv.New64a(),
		crc: crc64.New(crc64Table),
	}
}

func (sh *Hasher) Write(p []byte) (n int, err error) {
	expected := len(p)
	n, err = sh.fnv.Write(p)
	if err != nil {
		return n, err
	}
	if n != expected {
		return n, io.ErrShortWrite
	}

	n, err = sh.crc.Write(p)
	if err != nil {
		return n, err
	}
	if n != expected {
		return n, io.ErrShortWrite
	}

	sh.size += uint64(expected)
	return expected, nil
}

// WriteString computes the hash of s without heap allocating a byte slice.
func (sh *Hasher) WriteString(s string) (n int, err error) {
	expected := len(s)
	n, err = io.WriteString(sh.fnv, s)
	if err != nil {
		return n, err
	}
	if n != expected {
		return n, io.ErrShortWrite
	}

	n, err = io.WriteString(sh.crc, s)
	if err != nil {
		return n, err
	}
	if n != expected {
		return n, io.ErrShortWrite
	}

	sh.size += uint64(expected)
	return expected, nil
}

// Sum returns the composite StreamHash.
func (sh *Hasher) Sum() Hash {
	return Hash{
		FNV64: sh.fnv.Sum64(),
		Size:  sh.size,
		CRC64: sh.crc.Sum64(),
	}
}

// Reset resets the StreamHasher to its initial state.
func (sh *Hasher) Reset() {
	sh.fnv.Reset()
	sh.crc.Reset()
	sh.size = 0
}

type baseStreamHash struct {
	hash Hash
}

func (sh *baseStreamHash) Hash() Hash {
	return sh.hash
}
