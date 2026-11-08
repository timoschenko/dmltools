// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test_stream

import (
	"hash/crc64"
	"hash/fnv"
	"testing"

	"github.com/timoschenko/dmltools/pkg/fs/stream"
)

func TestFileHasher(t *testing.T) {
	content := "dml 1.4;\ndevice test_device;\n"
	hasher := stream.NewFileHasher()

	n, err := hasher.Write([]byte(content[:10]))
	if err != nil || n != 10 {
		t.Fatalf("Write returned n=%d, err=%v", n, err)
	}
	n, err = hasher.Write([]byte(content[10:]))
	if err != nil || n != len(content)-10 {
		t.Fatalf("Write returned n=%d, err=%v", n, err)
	}

	hash := hasher.Sum()

	// Verify size matches total byte count
	if hash.Size != uint64(len(content)) {
		t.Errorf("hash.Size = %d, want %d", hash.Size, len(content))
	}

	// Verify FNV-1a component
	fnvHasher := fnv.New64a()
	_, _ = fnvHasher.Write([]byte(content))
	if hash.FNV64 != fnvHasher.Sum64() {
		t.Errorf("hash.FNV64 = %x, want %x", hash.FNV64, fnvHasher.Sum64())
	}

	// Verify CRC64 component
	crcHasher := crc64.New(crc64.MakeTable(crc64.ISO))
	_, _ = crcHasher.Write([]byte(content))
	if hash.CRC64 != crcHasher.Sum64() {
		t.Errorf("hash.CRC64 = %x, want %x", hash.CRC64, crcHasher.Sum64())
	}

	// Verify Reset restores zero state
	hasher.Reset()
	zeroHash := hasher.Sum()
	emptyFNV := fnv.New64a().Sum64()
	if zeroHash.Size != 0 || zeroHash.CRC64 != 0 || zeroHash.FNV64 != emptyFNV {
		t.Errorf("Reset did not restore initial state: got %+v", zeroHash)
	}
}

func TestFileHash_Equality(t *testing.T) {
	h1 := stream.Hash{FNV64: 100, Size: 200, CRC64: 300}
	h2 := stream.Hash{FNV64: 100, Size: 200, CRC64: 300}
	if h1 != h2 {
		t.Errorf("identical hashes should be equal")
	}

	// Different FNV
	hDiffFNV := stream.Hash{FNV64: 101, Size: 200, CRC64: 300}
	if h1 == hDiffFNV {
		t.Errorf("hashes with different FNV should not be equal")
	}

	// Different Size
	hDiffSize := stream.Hash{FNV64: 100, Size: 201, CRC64: 300}
	if h1 == hDiffSize {
		t.Errorf("hashes with different Size should not be equal")
	}

	// Different CRC64
	hDiffCRC := stream.Hash{FNV64: 100, Size: 200, CRC64: 301}
	if h1 == hDiffCRC {
		t.Errorf("hashes with different CRC64 should not be equal")
	}
}
