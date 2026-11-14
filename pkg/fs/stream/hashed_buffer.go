// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package stream

import (
	"fmt"
	"github.com/timoschenko/dmltools/pkg/fs"
	"io"
	"math"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/antlr4-go/antlr/v4"
)

// HashedBuffer is an antlr.CharStream implementation that stores input data
// as an immutable UTF-8 string and decodes runes lazily at runtime.
//
// Internal indexing (Index, Size, Seek) operates in raw byte offsets rather than
// rune counts, eliminating the need to expand multi-byte UTF-8 data into a []rune
// slice upon initialization.
type HashedBuffer struct {
	baseStreamHash
	antlrStreamWrapper
	data     string
	capacity int
}

// NewHashedBufferFromStream constructs a HashedBuffer from an io.Reader.
// It pre-allocates buffer space using size hints, calculates the composite hash
// while streaming, and converts to an immutable string with zero-copy.
func NewHashedBufferFromStream(
	filename string, r io.Reader, size int,
) (*HashedBuffer, error) {
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

	var sb strings.Builder
	sb.Grow(allocSize)

	if _, err := io.Copy(&sb, teeStream); err != nil {
		return nil, err
	}

	data := sb.String()

	hb := HashedBuffer{
		hash:     inHasher.Sum(),
		filename: filename,
		index:    0,
		size:     len(data),
		data:     data,
		capacity: sb.Cap(),
	}

	return &hb, nil
}

// NewHashedBufferFromString constructs a HashedBuffer directly from an existing string.
// It assigns data with zero memory copying and hashes the string content in-place.
func NewHashedBufferFromString(filename string, data string) *HashedBuffer {
	inHasher := NewFileHasher()
	if n, err := inHasher.WriteString(data); err != nil || n != len(data) {
		panic(fmt.Sprintf("failed to hash string buffer %q: %v", filename, err))
	}

	return &HashedBuffer{
		hash:     inHasher.Sum(),
		filename: filename,
		index:    0,
		size:     len(data),
		data:     data,
		capacity: len(data),
	}
}

// Consume moves the input pointer to the next UTF-8 character in the input stream.
func (hb *HashedBuffer) Consume() {
	if hb.index >= hb.size {
		// assert is.LA(1) == TokenEOF
		panic("unreachable: cannot consume EOF")
	}

	if hb.data[hb.index] < utf8.RuneSelf {
		hb.index++
		return
	}

	_, width := utf8.DecodeRuneInString(hb.data[hb.index:])
	hb.index += width
}

func (hb *HashedBuffer) previousRune() int {
	if hb.index <= 0 {
		return antlr.TokenEOF
	}
	// Fast path: previous byte is ASCII
	b := hb.data[hb.index-1]
	if b < utf8.RuneSelf {
		return int(b)
	}

	// Previous byte is multi-byte or continuation byte
	start := hb.index - 1
	for start > 0 && start > hb.index-4 && (hb.data[start]&0xC0) == 0x80 {
		start--
	}
	r, width := utf8.DecodeRuneInString(hb.data[start:hb.index])
	if start+width == hb.index {
		return int(r)
	}
	// Fallback for single malformed byte
	return int(utf8.RuneError)
}

// LA returns the character at the given offset from the current stream position.
//
// Special offsets per antlr.IntStream specification:
// - LA(0) returns 0 (undefined/nil).
// - LA(1) returns the next character to be consumed (O(1) fast path).
// - LA(k) for k > 1 returns the character k-1 steps ahead, or TokenEOF.
// - LA(-1) returns the character at index - 1 (the most recently consumed character), or TokenEOF.
// - Offsets < -1 panic as unsupported.
func (hb *HashedBuffer) LA(offset int) int {
	if offset == 0 {
		return 0 // nil
	}

	// Fast path: LA(1) is the vast majority of all lexer calls
	if offset == 1 {
		if hb.index >= hb.size {
			return antlr.TokenEOF
		}
		b := hb.data[hb.index]
		if b < utf8.RuneSelf {
			return int(b)
		}
		r, _ := utf8.DecodeRuneInString(hb.data[hb.index:])
		return int(r)
	}

	// Lookback: LA(-1)
	if offset == -1 {
		return hb.previousRune()
	}

	// Forward lookahead: LA(k) for k > 1
	if offset > 1 {
		idx := hb.index
		for k := 1; k < offset; k++ {
			if idx >= hb.size {
				return antlr.TokenEOF
			}
			if hb.data[idx] < utf8.RuneSelf {
				idx++
			} else {
				_, width := utf8.DecodeRuneInString(hb.data[idx:])
				idx += width
			}
		}
		if idx >= hb.size {
			return antlr.TokenEOF
		}
		b := hb.data[idx]
		if b < utf8.RuneSelf {
			return int(b)
		}
		r, _ := utf8.DecodeRuneInString(hb.data[idx:])
		return int(r)
	}

	panic(fmt.Sprintf(
		"HashedBuffer: unsupported negative LA offset %d"+
			" (only LA(-1) lookback is supported)",
		offset,
	))
}

// GetText returns the slice of string between start and stop byte offsets inclusive.
func (hb *HashedBuffer) GetText(start int, stop int) string {
	if stop < start || start < 0 || start >= hb.size {
		return ""
	}
	if stop >= hb.size {
		stop = hb.size - 1
	}

	return hb.data[start : stop+1]
}

// Cap returns the allocated capacity of the internal buffer.
func (hb *HashedBuffer) Cap() int {
	return hb.capacity
}

// LT returns the character at the given offset from the start of the input stream.
func (hb *HashedBuffer) LT(offset int) int {
	return hb.LA(offset)
}

// GetTextFromTokens returns the text between the start and stop tokens.
func (hb *HashedBuffer) GetTextFromTokens(
	start, stop antlr.Token,
) string {
	if start == nil || stop == nil {
		return ""
	}

	return hb.GetText(start.GetStart(), stop.GetStop())
}

// GetTextFromInterval returns the text from the interval [i.Start, i.Stop].
func (hb *HashedBuffer) GetTextFromInterval(i antlr.Interval) string {
	return hb.GetText(i.Start, i.Stop)
}

// GetRawString returns the underlying string content of the buffer.
func (hb *HashedBuffer) GetRawString() string {
	return hb.data
}
