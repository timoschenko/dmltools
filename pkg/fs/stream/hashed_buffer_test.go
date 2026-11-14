// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package stream

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/antlr4-go/antlr/v4"
)

func TestNewHashedBufferFromStream(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		input    string
		bufSize  int
		wantSize int
		wantText string
	}{
		{
			name:     "empty input",
			filename: "empty.dml",
			input:    "",
			bufSize:  16,
			wantSize: 0,
			wantText: "",
		},
		{
			name:     "ascii text",
			filename: "simple.dml",
			input:    "dml 1.4;\ndevice dev;",
			bufSize:  64,
			wantSize: 20,
			wantText: "dml 1.4;\ndevice dev;",
		},
		{
			name:     "unicode multi-byte runes",
			filename: "unicode.dml",
			input:    "log \"Hello, 世界! 🚀\";",
			bufSize:  8,
			wantSize: len("log \"Hello, 世界! 🚀\";"),
			wantText: "log \"Hello, 世界! 🚀\";",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.input)
			buf, err := NewHashedBufferFromStream(
				tc.filename,
				r,
				tc.bufSize,
			)
			if err != nil {
				t.Fatalf("unexpected error creating HashedBuffer: %v", err)
			}

			if buf.GetSourceName() != tc.filename {
				t.Errorf(
					"GetSourceName() = %q, want %q",
					buf.GetSourceName(),
					tc.filename,
				)
			}
			if buf.Size() != tc.wantSize {
				t.Errorf("Size() = %d, want %d", buf.Size(), tc.wantSize)
			}

			hasher := NewFileHasher()
			_, _ = hasher.Write([]byte(tc.input))
			expectedHash := hasher.Sum()
			if buf.Hash() != expectedHash {
				t.Errorf("Hash() = %+v, want %+v", buf.Hash(), expectedHash)
			}
			if buf.Hash().Size != uint64(len(tc.input)) {
				t.Errorf(
					"Hash().Size = %d, want %d",
					buf.Hash().Size,
					len(tc.input),
				)
			}

			if tc.wantSize > 0 {
				got := buf.GetText(0, tc.wantSize-1)
				if got != tc.wantText {
					t.Errorf(
						"GetText(0, %d) = %q, want %q",
						tc.wantSize-1,
						got,
						tc.wantText,
					)
				}
			}
		})
	}
}

func TestHashedBuffer_NavigationAndUTF8(t *testing.T) {
	// 'A' -> 1 byte (index 0)
	// '世' -> 3 bytes (index 1..3)
	// '🚀' -> 4 bytes (index 4..7)
	// '!' -> 1 byte (index 8)
	input := "A世🚀!"
	buf, err := NewHashedBufferFromStream(
		"utf8.dml",
		strings.NewReader(input),
		0,
	)
	if err != nil {
		t.Fatalf("failed to create HashedBuffer: %v", err)
	}

	expectedRunes := []struct {
		runeVal   rune
		byteIndex int
		width     int
	}{
		{'A', 0, 1},
		{'世', 1, 3},
		{'🚀', 4, 4},
		{'!', 8, 1},
	}

	for i, exp := range expectedRunes {
		if buf.Index() != exp.byteIndex {
			t.Fatalf(
				"[%d] Index() = %d, want %d",
				i,
				buf.Index(),
				exp.byteIndex,
			)
		}
		gotRune := buf.LA(1)
		if gotRune != int(exp.runeVal) {
			t.Fatalf("[%d] LA(1) = %U, want %U", i, gotRune, exp.runeVal)
		}
		if buf.LT(1) != gotRune {
			t.Fatalf("[%d] LT(1) = %U, want %U", i, buf.LT(1), exp.runeVal)
		}
		buf.Consume()
	}

	if buf.Index() != len(input) {
		t.Errorf(
			"Index() after consumes = %d, want %d",
			buf.Index(),
			len(input),
		)
	}
	if buf.LA(1) != antlr.TokenEOF {
		t.Errorf("LA(1) at EOF = %d, want %d", buf.LA(1), antlr.TokenEOF)
	}

	// Consuming past EOF should panic
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("expected panic on Consume() at EOF")
		}
	}()
	buf.Consume()
}

func TestHashedBuffer_LA_LookaheadAndLookback(t *testing.T) {
	t.Run("ASCII lookahead and lookback", func(t *testing.T) {
		buf, err := NewHashedBufferFromStream(
			"test.dml",
			strings.NewReader("hello"),
			0,
		)
		if err != nil {
			t.Fatalf("failed to create HashedBuffer: %v", err)
		}

		if got := buf.LA(0); got != 0 {
			t.Errorf("LA(0) = %d, want 0", got)
		}

		// Initial position: index 0
		if got := buf.LA(1); got != int('h') {
			t.Errorf("LA(1) = %v, want 'h'", got)
		}
		if got := buf.LA(2); got != int('e') {
			t.Errorf("LA(2) = %v, want 'e'", got)
		}
		if got := buf.LA(3); got != int('l') {
			t.Errorf("LA(3) = %v, want 'l'", got)
		}
		if got := buf.LA(4); got != int('l') {
			t.Errorf("LA(4) = %v, want 'l'", got)
		}
		if got := buf.LA(5); got != int('o') {
			t.Errorf("LA(5) = %v, want 'o'", got)
		}
		if got := buf.LA(6); got != antlr.TokenEOF {
			t.Errorf("LA(6) = %d, want TokenEOF", got)
		}
		if got := buf.LA(100); got != antlr.TokenEOF {
			t.Errorf("LA(100) = %d, want TokenEOF", got)
		}

		// Lookback at start of stream must return TokenEOF
		if got := buf.LA(-1); got != antlr.TokenEOF {
			t.Errorf("LA(-1) at start = %d, want TokenEOF", got)
		}

		// Consume 'h'
		buf.Consume()
		if got := buf.LA(-1); got != int('h') {
			t.Errorf("LA(-1) after consume = %v, want 'h'", got)
		}
		if got := buf.LA(1); got != int('e') {
			t.Errorf("LA(1) after consume = %v, want 'e'", got)
		}
		if got := buf.LA(2); got != int('l') {
			t.Errorf("LA(2) after consume = %v, want 'l'", got)
		}

		// Deep negative lookback (LA(-2)) should panic
		func() {
			defer func() {
				r := recover()
				if r == nil {
					t.Errorf("expected panic on LA(-2)")
				}
			}()
			_ = buf.LA(-2)
		}()
	})

	t.Run("Multi-byte Unicode lookahead and lookback", func(t *testing.T) {
		buf, err := NewHashedBufferFromStream(
			"unicode.dml",
			strings.NewReader("A世🚀!"),
			0,
		)
		if err != nil {
			t.Fatalf("failed to create HashedBuffer: %v", err)
		}

		// Lookahead at start
		if got := buf.LA(1); got != int('A') {
			t.Errorf("LA(1) = %U, want 'A'", got)
		}
		if got := buf.LA(2); got != int('世') {
			t.Errorf("LA(2) = %U, want '世'", got)
		}
		if got := buf.LA(3); got != int('🚀') {
			t.Errorf("LA(3) = %U, want '🚀'", got)
		}
		if got := buf.LA(4); got != int('!') {
			t.Errorf("LA(4) = %U, want '!''", got)
		}
		if got := buf.LA(5); got != antlr.TokenEOF {
			t.Errorf("LA(5) = %d, want TokenEOF", got)
		}
		if got := buf.LA(-1); got != antlr.TokenEOF {
			t.Errorf("LA(-1) at start = %d, want TokenEOF", got)
		}

		// Consume 'A'
		buf.Consume()
		if got := buf.LA(-1); got != int('A') {
			t.Errorf("LA(-1) after 'A' = %U, want 'A'", got)
		}
		if got := buf.LA(1); got != int('世') {
			t.Errorf("LA(1) after 'A' = %U, want '世'", got)
		}
		if got := buf.LA(2); got != int('🚀') {
			t.Errorf("LA(2) after 'A' = %U, want '🚀'", got)
		}

		// Consume '世'
		buf.Consume()
		if got := buf.LA(-1); got != int('世') {
			t.Errorf("LA(-1) after '世' = %U, want '世'", got)
		}
		if got := buf.LA(1); got != int('🚀') {
			t.Errorf("LA(1) after '世' = %U, want '🚀'", got)
		}

		// Consume '🚀'
		buf.Consume()
		if got := buf.LA(-1); got != int('🚀') {
			t.Errorf("LA(-1) after '🚀' = %U, want '🚀'", got)
		}
		if got := buf.LA(1); got != int('!') {
			t.Errorf("LA(1) after '🚀' = %U, want '!'", got)
		}

		// Consume '!'
		buf.Consume()
		if got := buf.LA(-1); got != int('!') {
			t.Errorf("LA(-1) after '!' = %U, want '!'", got)
		}
		if got := buf.LA(1); got != antlr.TokenEOF {
			t.Errorf("LA(1) after '!' = %d, want TokenEOF", got)
		}
	})
}

func TestHashedBuffer_Seek(t *testing.T) {
	content := "hello world"
	buf, err := NewHashedBufferFromStream(
		"test.dml",
		strings.NewReader(content),
		0,
	)
	if err != nil {
		t.Fatalf("failed to create buffer: %v", err)
	}

	buf.Seek(5)
	if buf.Index() != 5 {
		t.Errorf("Seek(5) = %d, want 5", buf.Index())
	}

	buf.Seek(-10)
	if buf.Index() != 0 {
		t.Errorf("Seek(-10) = %d, want 0", buf.Index())
	}

	buf.Seek(100)
	if buf.Index() != len(content) {
		t.Errorf("Seek(100) = %d, want %d", buf.Index(), len(content))
	}
}

func TestHashedBuffer_GetText(t *testing.T) {
	content := "hello 世界"
	// "hello " = 6 bytes (0..5)
	// "世" = 3 bytes (6..8)
	// "界" = 3 bytes (9..11)
	// total = 12 bytes

	buf, err := NewHashedBufferFromStream(
		"test.dml",
		strings.NewReader(content),
		0,
	)
	if err != nil {
		t.Fatalf("failed to create buffer: %v", err)
	}

	tests := []struct {
		name     string
		start    int
		stop     int
		expected string
	}{
		{"full range", 0, 11, "hello 世界"},
		{"subslice ascii", 1, 4, "ello"},
		{"subslice multibyte first", 6, 8, "世"},
		{"subslice multibyte second", 9, 11, "界"},
		{"inverted stop < start", 5, 3, ""},
		{"negative start", -1, 5, ""},
		{"start beyond size", 12, 15, ""},
		{"stop beyond size clamps", 6, 50, "世界"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := buf.GetText(tc.start, tc.stop)
			if got != tc.expected {
				t.Errorf(
					"GetText(%d, %d) = %q, want %q",
					tc.start,
					tc.stop,
					got,
					tc.expected,
				)
			}
			intervalGot := buf.GetTextFromInterval(
				antlr.NewInterval(tc.start, tc.stop),
			)
			if intervalGot != tc.expected {
				t.Errorf(
					"GetTextFromInterval(%d, %d) = %q, want %q",
					tc.start,
					tc.stop,
					intervalGot,
					tc.expected,
				)
			}
		})
	}

	// Nil tokens should return empty string
	if got := buf.GetTextFromTokens(nil, nil); got != "" {
		t.Errorf("GetTextFromTokens(nil, nil) = %q, want empty", got)
	}
}

func TestHashedBuffer_MarkReleaseCap(t *testing.T) {
	buf, err := NewHashedBufferFromStream(
		"test.dml",
		strings.NewReader("sample"),
		0,
	)
	if err != nil {
		t.Fatalf("failed to create buffer: %v", err)
	}

	if marker := buf.Mark(); marker != -1 {
		t.Errorf("Mark() = %d, want -1", marker)
	}
	buf.Release(0) // should not panic

	if buf.Cap() < len("sample") {
		t.Errorf("Cap() = %d, want >= %d", buf.Cap(), len("sample"))
	}
}

func TestHashedBuffer_MalformedUTF8(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantRunes []rune
	}{
		{
			name:      "truncated 2-byte sequence at EOF",
			input:     "\xC2",
			wantRunes: []rune{utf8.RuneError},
		},
		{
			name:      "truncated 3-byte sequence at EOF",
			input:     "\xE4\xB8",
			wantRunes: []rune{utf8.RuneError, utf8.RuneError},
		},
		{
			name:      "truncated 4-byte sequence at EOF",
			input:     "\xF0\x9F\x98",
			wantRunes: []rune{utf8.RuneError, utf8.RuneError, utf8.RuneError},
		},
		{
			name:      "truncated 2-byte sequence followed by ASCII",
			input:     "\xC2a",
			wantRunes: []rune{utf8.RuneError, 'a'},
		},
		{
			name:      "truncated 3-byte sequence followed by ASCII",
			input:     "\xE4\xB8a",
			wantRunes: []rune{utf8.RuneError, utf8.RuneError, 'a'},
		},
		{
			name:  "truncated 4-byte sequence followed by ASCII",
			input: "\xF0\x9F\x98a",
			wantRunes: []rune{
				utf8.RuneError,
				utf8.RuneError,
				utf8.RuneError,
				'a',
			},
		},
		{
			name:      "stray single continuation byte 0x80",
			input:     "\x80",
			wantRunes: []rune{utf8.RuneError},
		},
		{
			name:      "stray single continuation byte 0xBF",
			input:     "\xBF",
			wantRunes: []rune{utf8.RuneError},
		},
		{
			name:      "consecutive stray continuation bytes",
			input:     "\x80\xBF\x90",
			wantRunes: []rune{utf8.RuneError, utf8.RuneError, utf8.RuneError},
		},
		{
			name:      "stray continuation byte between ASCII",
			input:     "a\x80b",
			wantRunes: []rune{'a', utf8.RuneError, 'b'},
		},
		{
			name:      "illegal lead byte 0xFF",
			input:     "\xFF",
			wantRunes: []rune{utf8.RuneError},
		},
		{
			name:      "illegal lead byte 0xFE",
			input:     "\xFE",
			wantRunes: []rune{utf8.RuneError},
		},
		{
			name:      "illegal lead byte 0xF5 (out of Unicode range)",
			input:     "\xF5",
			wantRunes: []rune{utf8.RuneError},
		},
		{
			name:      "overlong 2-byte NUL (0xC0 0x80)",
			input:     "\xC0\x80",
			wantRunes: []rune{utf8.RuneError, utf8.RuneError},
		},
		{
			name:      "overlong 2-byte ASCII 0x7F (0xC1 0xBF)",
			input:     "\xC1\xBF",
			wantRunes: []rune{utf8.RuneError, utf8.RuneError},
		},
		{
			name:      "overlong 3-byte (0xE0 0x80 0xAF)",
			input:     "\xE0\x80\xAF",
			wantRunes: []rune{utf8.RuneError, utf8.RuneError, utf8.RuneError},
		},
		{
			name:  "overlong 4-byte (0xF0 0x80 0x80 0xAF)",
			input: "\xF0\x80\x80\xAF",
			wantRunes: []rune{
				utf8.RuneError,
				utf8.RuneError,
				utf8.RuneError,
				utf8.RuneError,
			},
		},
		{
			name:      "surrogate half lead (0xED 0xA0 0x80, U+D800)",
			input:     "\xED\xA0\x80",
			wantRunes: []rune{utf8.RuneError, utf8.RuneError, utf8.RuneError},
		},
		{
			name:      "surrogate half trail (0xED 0xBF 0xBF, U+DFFF)",
			input:     "\xED\xBF\xBF",
			wantRunes: []rune{utf8.RuneError, utf8.RuneError, utf8.RuneError},
		},
		{
			name:  "out of Unicode range codepoint (> U+10FFFF, 0xF4 0x90 0x80 0x80)",
			input: "\xF4\x90\x80\x80",
			wantRunes: []rune{
				utf8.RuneError,
				utf8.RuneError,
				utf8.RuneError,
				utf8.RuneError,
			},
		},
		{
			name:  "mixed valid and invalid bytes",
			input: "ok\xFF世界\x80!",
			wantRunes: []rune{
				'o',
				'k',
				utf8.RuneError,
				'世',
				'界',
				utf8.RuneError,
				'!',
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf, err := NewHashedBufferFromStream(
				"malformed.dml",
				strings.NewReader(tc.input),
				len(tc.input),
			)
			if err != nil {
				t.Fatalf("failed to create HashedBuffer: %v", err)
			}

			// 1. Verify byte-by-byte traversal via LA(1) and Consume()
			var gotRunes []rune
			for {
				la := buf.LA(1)
				if la == antlr.TokenEOF {
					break
				}
				//nolint:gosec // la is a valid decoded rune (or RuneError) within Unicode range
				gotRunes = append(gotRunes, rune(la))
				buf.Consume()
			}

			if len(gotRunes) != len(tc.wantRunes) {
				t.Fatalf(
					"rune count mismatch: got %d runes (%+q), want %d runes (%+q)",
					len(gotRunes),
					gotRunes,
					len(tc.wantRunes),
					tc.wantRunes,
				)
			}

			for i := range gotRunes {
				if gotRunes[i] != tc.wantRunes[i] {
					t.Errorf(
						"rune [%d] mismatch: got %U (%c), want %U (%c)",
						i,
						gotRunes[i],
						gotRunes[i],
						tc.wantRunes[i],
						tc.wantRunes[i],
					)
				}
			}

			// 2. Verify Index() reached end of buffer
			if buf.Index() != len(tc.input) {
				t.Errorf("Index() = %d, want %d", buf.Index(), len(tc.input))
			}

			// 3. Verify GetText() returns exact raw bytes
			if len(tc.input) > 0 {
				fullText := buf.GetText(0, len(tc.input)-1)
				if fullText != tc.input {
					t.Errorf("GetText() = %q, want %q", fullText, tc.input)
				}
			}

			// 4. Verify parity against HashedRunesBuffer on identical input
			runesBuf, err := NewHashedRunesBufferFromStream(
				"malformed.dml",
				strings.NewReader(tc.input),
				len(tc.input),
			)
			if err != nil {
				t.Fatalf("failed to create HashedRunesBuffer: %v", err)
			}

			var runesGot []rune
			for {
				la := runesBuf.LA(1)
				if la == antlr.TokenEOF {
					break
				}
				//nolint:gosec // la is a valid decoded rune (or RuneError) within Unicode range
				runesGot = append(runesGot, rune(la))
				runesBuf.Consume()
			}

			if len(gotRunes) != len(runesGot) {
				t.Errorf(
					"parity mismatch with HashedRunesBuffer: HashedBuffer got %d runes, HashedRunesBuffer got %d runes",
					len(gotRunes),
					len(runesGot),
				)
			} else {
				for i := range gotRunes {
					if gotRunes[i] != runesGot[i] {
						t.Errorf(
							"parity rune [%d] mismatch: HashedBuffer=%U, HashedRunesBuffer=%U",
							i,
							gotRunes[i],
							runesGot[i],
						)
					}
				}
			}
		})
	}
}

func TestNewHashedBufferFromString(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		input    string
	}{
		{"empty string", "empty.dml", ""},
		{"simple ascii", "simple.dml", "dml 1.4;\ndevice dev;"},
		{"unicode multibyte", "unicode.dml", "log \"Hello, 世界! 🚀\";"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := NewHashedBufferFromString(tc.filename, tc.input)

			if buf.GetSourceName() != tc.filename {
				t.Errorf(
					"GetSourceName() = %q, want %q",
					buf.GetSourceName(),
					tc.filename,
				)
			}
			if buf.Size() != len(tc.input) {
				t.Errorf("Size() = %d, want %d", buf.Size(), len(tc.input))
			}
			if buf.GetRawString() != tc.input {
				t.Errorf("String() = %q, want %q", buf.GetRawString(), tc.input)
			}
			if buf.Cap() != len(tc.input) {
				t.Errorf("Cap() = %d, want %d", buf.Cap(), len(tc.input))
			}

			// Verify hash matches NewHashedBufferFromStream
			streamBuf, err := NewHashedBufferFromStream(
				tc.filename,
				strings.NewReader(tc.input),
				len(tc.input),
			)
			if err != nil {
				t.Fatalf("NewHashedBufferFromStream failed: %v", err)
			}

			if buf.Hash() != streamBuf.Hash() {
				t.Errorf("Hash() = %+v, want %+v", buf.Hash(), streamBuf.Hash())
			}
		})
	}
}
