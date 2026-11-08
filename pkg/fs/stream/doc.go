// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

// Package stream provides streaming buffers, composite hashing, and ANTLR
// character stream implementations for cached file processing.
//
// # Overview
//
// When parsing and formatting large DML source trees, file caching is vital
// for avoiding redundant AST reconstruction. This package provides high-performance
// I/O wrappers that simultaneously stream source data into ANTLR-compatible
// character streams and compute composite content fingerprints.
//
// # Key Components
//
//   - [Hash] and [Hasher]: Computes a composite 192-bit checksum consisting of
//     a 64-bit FNV-1a hash, the total byte size, and an ISO 64-bit CRC64.
//     This composite structure provides strong collision resistance and rapid
//     change detection.
//   - [HashedBuffer]: An [antlr.CharStream] implementation that stores input data
//     as an immutable UTF-8 string, decoding runes lazily on demand. It avoids
//     eager rune slice allocation, reducing memory footprint by ~4x compared to
//     traditional rune-buffered streams.
//   - [NewHashedBufferFromStream] and [NewHashedBufferFromString]: Factory
//     functions for initializing hashed buffers from readers or strings.
//
// # Usage Example
//
//	r := strings.NewReader("dml 1.4;\ndevice test;\n")
//	buf, err := stream.NewHashedBufferFromStream("test.dml", r, r.Len())
//	if err != nil {
//	    log.Fatalf("failed to create hashed buffer: %v", err)
//	}
//
//	fmt.Println("Hash:", buf.Hash().String())
//	// buf satisfies antlr.CharStream and can be fed directly to dmlex.ParseStream
package stream
