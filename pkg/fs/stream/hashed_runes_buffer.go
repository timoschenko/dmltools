// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package stream

import (
	"bufio"
	"github.com/timoschenko/dmltools/pkg/fs"
	"io"
	"math"
	"os"

	"github.com/antlr4-go/antlr/v4"
)

// HashedRunesBuffer is an antlr.CharStream implementation that stores input data as []rune.
//
// Deprecated: Use HashedBuffer instead. HashedRunesBuffer eagerly decodes the
// entire stream into []rune upon initialization, consuming ~4x memory compared to
// raw UTF-8 file sizes.
type HashedRunesBuffer struct {
	baseStreamHash
	antlrStreamWrapper
	data []rune
}

// NewHashedRunesBufferFromStream constructs a HashedRunesBuffer from an io.Reader.
//
// Deprecated: Use NewHashedBufferFromStream or NewHashedBufferFromString instead.
func NewHashedRunesBufferFromStream(
	filename string, r io.Reader, size int,
) (*HashedRunesBuffer, error) {
	if size <= 0 {
		if s, ok := fs.SizeFromReader(r); ok {
			size = s
		} else if filename != "" {
			if fi, err := os.Stat(
				filename,
			); err == nil &&
				fi.Mode().IsRegular() {
				if fi.Size() >= 0 && fi.Size() <= int64(math.MaxInt) {
					size = int(fi.Size())
				}
			}
		}
	}

	allocSize := fs.CalculateStreamBufferSize(size)

	inHasher := NewFileHasher()
	teeStream := io.TeeReader(r, inHasher)
	reader := bufio.NewReaderSize(teeStream, allocSize)

	hfs := HashedRunesBuffer{
		baseStreamHash: baseStreamHash{
			hash: Hash{},
		},
		antlrStreamWrapper: antlrStreamWrapper{
			filename: filename,
			index:    0,
			size:     0,
		},
		data: make([]rune, 0, allocSize),
	}

	// Pre-build the buffer and read runes efficiently
	//
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break // Normal end of stream
			}
			return nil, err
		}
		hfs.data = append(hfs.data, r)
	}

	hfs.size = len(hfs.data) // Size in runes
	hfs.hash = inHasher.Sum()

	return &hfs, nil
}

// Consume moves the input pointer to the next character in the input stream
func (hrb *HashedRunesBuffer) Consume() {
	if hrb.index >= hrb.size {
		// assert is.LA(1) == TokenEOF
		panic("unreachable: cannot consume EOF")
	}
	hrb.index++
}

// LA returns the character at the given offset from the start of the input stream
func (hrb *HashedRunesBuffer) LA(offset int) int {
	if offset == 0 {
		return 0 // nil
	}
	if offset < 0 {
		offset++ // e.g., translate LA(-1) to use offset=0
	}
	pos := hrb.index + offset - 1

	if pos < 0 || pos >= hrb.size { // invalid
		return antlr.TokenEOF
	}

	return int(hrb.data[pos])
}

func (hrb *HashedRunesBuffer) GetText(start int, stop int) string {
	if stop < start || start < 0 {
		return ""
	}
	if stop >= hrb.size {
		stop = hrb.size - 1
	}
	if start >= hrb.size {
		return ""
	}

	return string(hrb.data[start : stop+1])
}

// Cap returns the allocated capacity of the internal rune buffer.
func (hrb *HashedRunesBuffer) Cap() int {
	return cap(hrb.data)
}

// LT returns the character at the given offset from the start of the input stream.
func (hrb *HashedRunesBuffer) LT(offset int) int {
	return hrb.LA(offset)
}

// GetTextFromTokens returns the text between the start and stop tokens.
func (hrb *HashedRunesBuffer) GetTextFromTokens(
	start, stop antlr.Token,
) string {
	if start == nil || stop == nil {
		return ""
	}

	return hrb.GetText(start.GetStart(), stop.GetStop())
}

// GetTextFromInterval returns the text from the interval [i.Start, i.Stop].
func (hrb *HashedRunesBuffer) GetTextFromInterval(i antlr.Interval) string {
	return hrb.GetText(i.Start, i.Stop)
}
