// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test_stream

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/antlr4-go/antlr/v4"

	"github.com/timoschenko/dmltools/pkg/fs/stream"
)

func TestNewHashingStream(t *testing.T) {
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
			wantSize: 19,
			wantText: "log \"Hello, 世界! 🚀\";",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.input)
			sbuf, err := stream.NewHashedRunesBufferFromStream(
				tc.filename,
				r,
				tc.bufSize,
			)
			if err != nil {
				t.Fatalf("unexpected error creating HashingFileStream: %v", err)
			}

			if sbuf.GetSourceName() != tc.filename {
				t.Errorf(
					"GetSourceName() = %q, want %q",
					sbuf.GetSourceName(),
					tc.filename,
				)
			}
			if sbuf.Size() != tc.wantSize {
				t.Errorf("Size() = %d, want %d", sbuf.Size(), tc.wantSize)
			}

			hasher := stream.NewFileHasher()
			_, _ = hasher.Write([]byte(tc.input))
			expectedHash := hasher.Sum()
			if sbuf.Hash() != expectedHash {
				t.Errorf("Hash() = %+v, want %+v", sbuf.Hash(), expectedHash)
			}
			if sbuf.Hash().Size != uint64(len(tc.input)) {
				t.Errorf(
					"Hash().Size = %d, want %d",
					sbuf.Hash().Size,
					len(tc.input),
				)
			}

			if tc.wantSize > 0 {
				got := sbuf.GetText(0, tc.wantSize-1)
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

func TestHashingFileStream_Seek(t *testing.T) {
	content := "abcdef"
	tests := []struct {
		name      string
		initial   int
		seekTo    int
		wantIndex int
	}{
		{
			name:      "seek negative from start bounds to 0",
			initial:   0,
			seekTo:    -5,
			wantIndex: 0,
		},
		{
			name:      "seek negative from middle bounds to 0",
			initial:   3,
			seekTo:    -1,
			wantIndex: 0,
		},
		{
			name:      "seek to zero",
			initial:   2,
			seekTo:    0,
			wantIndex: 0,
		},
		{
			name:      "seek forward within bounds",
			initial:   1,
			seekTo:    4,
			wantIndex: 4,
		},
		{
			name:      "seek backward within bounds",
			initial:   5,
			seekTo:    2,
			wantIndex: 2,
		},
		{
			name:      "seek to exact size (EOF position)",
			initial:   0,
			seekTo:    6,
			wantIndex: 6,
		},
		{
			name:      "seek beyond size bounds to size",
			initial:   2,
			seekTo:    100,
			wantIndex: 6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stream, err := stream.NewHashedRunesBufferFromStream(
				"test.dml",
				strings.NewReader(content),
				16,
			)
			if err != nil {
				t.Fatalf("failed to create stream: %v", err)
			}

			stream.Seek(tc.initial)
			if stream.Index() != tc.initial {
				t.Fatalf(
					"initial Index() = %d, want %d",
					stream.Index(),
					tc.initial,
				)
			}

			stream.Seek(tc.seekTo)
			if stream.Index() != tc.wantIndex {
				t.Errorf(
					"Seek(%d): Index() = %d, want %d",
					tc.seekTo,
					stream.Index(),
					tc.wantIndex,
				)
			}
		})
	}
}

func TestHashingFileStream_GetText(t *testing.T) {
	content := "hello 世界"
	// Runes:
	// 0: 'h'
	// 1: 'e'
	// 2: 'l'
	// 3: 'l'
	// 4: 'o'
	// 5: ' '
	// 6: '世'
	// 7: '界'
	// total size: 8 runes

	tests := []struct {
		name     string
		input    string
		start    int
		stop     int
		expected string
	}{
		{
			name:     "valid full range",
			input:    content,
			start:    0,
			stop:     7,
			expected: "hello 世界",
		},
		{
			name:     "valid subslice",
			input:    content,
			start:    1,
			stop:     4,
			expected: "ello",
		},
		{
			name:     "single rune",
			input:    content,
			start:    6,
			stop:     6,
			expected: "世",
		},
		{
			name:     "inverted interval (stop < start)",
			input:    content,
			start:    5,
			stop:     3,
			expected: "",
		},
		{
			name:     "antlr epsilon/missing interval (start=0, stop=-1)",
			input:    content,
			start:    0,
			stop:     -1,
			expected: "",
		},
		{
			name:     "both start and stop negative",
			input:    content,
			start:    -5,
			stop:     -2,
			expected: "",
		},
		{
			name:     "negative start with valid stop",
			input:    content,
			start:    -1,
			stop:     4,
			expected: "",
		},
		{
			name:     "start at size",
			input:    content,
			start:    8,
			stop:     10,
			expected: "",
		},
		{
			name:     "start beyond size",
			input:    content,
			start:    10,
			stop:     12,
			expected: "",
		},
		{
			name:     "stop exceeds size (clamps to end)",
			input:    content,
			start:    6,
			stop:     20,
			expected: "世界",
		},
		{
			name:     "empty stream: any interval returns empty",
			input:    "",
			start:    0,
			stop:     0,
			expected: "",
		},
		{
			name:     "empty stream: negative bounds",
			input:    "",
			start:    -1,
			stop:     0,
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sbuf, err := stream.NewHashedRunesBufferFromStream(
				"test.dml",
				strings.NewReader(tc.input),
				16,
			)
			if err != nil {
				t.Fatalf("failed to create stream: %v", err)
			}

			got := sbuf.GetText(tc.start, tc.stop)
			if got != tc.expected {
				t.Errorf(
					"GetText(%d, %d) = %q, want %q",
					tc.start,
					tc.stop,
					got,
					tc.expected,
				)
			}
		})
	}
}

func TestHashingFileStream_GetTextFromInterval(t *testing.T) {
	sbuf, err := stream.NewHashedRunesBufferFromStream(
		"test.dml",
		strings.NewReader("dml 1.4;"),
		16,
	)
	if err != nil {
		t.Fatalf("failed to create stream: %v", err)
	}

	tests := []struct {
		name     string
		interval antlr.Interval
		expected string
	}{
		{
			name:     "normal interval",
			interval: antlr.NewInterval(0, 2),
			expected: "dml",
		},
		{
			name:     "inverted interval",
			interval: antlr.NewInterval(3, 1),
			expected: "",
		},
		{
			name:     "negative interval",
			interval: antlr.NewInterval(-2, -1),
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := sbuf.GetTextFromInterval(tc.interval)
			if got != tc.expected {
				t.Errorf(
					"GetTextFromInterval(%v) = %q, want %q",
					tc.interval,
					got,
					tc.expected,
				)
			}
		})
	}
}

func TestHashingFileStream_GetTextFromTokens(t *testing.T) {
	sbuf, err := stream.NewHashedRunesBufferFromStream(
		"test.dml",
		strings.NewReader("dml 1.4;"),
		16,
	)
	if err != nil {
		t.Fatalf("failed to create stream: %v", err)
	}

	t.Run("nil tokens return empty string", func(t *testing.T) {
		tok := &antlr.CommonToken{}
		tok.SetTokenIndex(0)
		if got := sbuf.GetTextFromTokens(nil, tok); got != "" {
			t.Errorf("GetTextFromTokens(nil, tok) = %q, want empty", got)
		}
		if got := sbuf.GetTextFromTokens(tok, nil); got != "" {
			t.Errorf("GetTextFromTokens(tok, nil) = %q, want empty", got)
		}
		if got := sbuf.GetTextFromTokens(nil, nil); got != "" {
			t.Errorf("GetTextFromTokens(nil, nil) = %q, want empty", got)
		}
	})
}

func TestHashingFileStream_LA_LT_Consume(t *testing.T) {
	sbuf, err := stream.NewHashedRunesBufferFromStream(
		"test.dml",
		strings.NewReader("ab"),
		16,
	)
	if err != nil {
		t.Fatalf("failed to create stream: %v", err)
	}

	// Initial position: 0
	if sbuf.LA(0) != 0 {
		t.Errorf("LA(0) = %d, want 0", sbuf.LA(0))
	}
	if sbuf.LT(1) != int('a') {
		t.Errorf("LT(1) = %q, want 'a'", sbuf.LT(1))
	}
	if sbuf.LA(1) != int('a') {
		t.Errorf("LA(1) = %q, want 'a'", sbuf.LA(1))
	}
	if sbuf.LA(2) != int('b') {
		t.Errorf("LA(2) = %q, want 'b'", sbuf.LA(2))
	}
	if sbuf.LA(3) != antlr.TokenEOF {
		t.Errorf("LA(3) = %d, want TokenEOF (%d)", sbuf.LA(3), antlr.TokenEOF)
	}
	// Lookbehind before start
	if sbuf.LA(-1) != antlr.TokenEOF {
		t.Errorf("LA(-1) = %d, want TokenEOF", sbuf.LA(-1))
	}

	// Consume 'a'
	sbuf.Consume()
	if sbuf.Index() != 1 {
		t.Errorf("Index() after Consume = %d, want 1", sbuf.Index())
	}
	if sbuf.LA(1) != int('b') {
		t.Errorf("LA(1) = %q, want 'b'", sbuf.LA(1))
	}
	if sbuf.LA(-1) != int('a') {
		t.Errorf("LA(-1) = %q, want 'a'", sbuf.LA(-1))
	}

	// Consume 'b'
	sbuf.Consume()
	if sbuf.Index() != 2 {
		t.Errorf("Index() after Consume = %d, want 2", sbuf.Index())
	}
	if sbuf.LA(1) != antlr.TokenEOF {
		t.Errorf("LA(1) at EOF = %d, want TokenEOF", sbuf.LA(1))
	}

	// Consuming at EOF must panic
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("expected panic on Consume() at EOF, but did not panic")
		}
	}()
	sbuf.Consume()
}

func TestHashingFileStream_Mark_Release(t *testing.T) {
	sbuf, err := stream.NewHashedRunesBufferFromStream(
		"test.dml",
		strings.NewReader("test"),
		16,
	)
	if err != nil {
		t.Fatalf("failed to create stream: %v", err)
	}

	mark := sbuf.Mark()
	if mark != -1 {
		t.Errorf("Mark() = %d, want -1", mark)
	}
	// Release should be a no-op and not panic
	sbuf.Release(mark)
}

type opaqueTestReader struct {
	r io.Reader
}

func (o *opaqueTestReader) Read(p []byte) (int, error) {
	return o.r.Read(p)
}

func TestNewHashingStream_SizeZeroDiscovery(t *testing.T) {
	t.Run("strings.Reader with size 0 discovers size", func(t *testing.T) {
		content := "dml 1.4;\ndevice dev;"
		r := strings.NewReader(content)
		sbuf, err := stream.NewHashedRunesBufferFromStream("virtual.dml", r, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sbuf.Size() != len(content) {
			t.Errorf("stream.Size() = %d, want %d", sbuf.Size(), len(content))
		}
		if sbuf.Cap() < 4096 {
			t.Errorf("stream.Cap() = %d, want >= 4096", sbuf.Cap())
		}
	})

	t.Run(
		"bytes.Buffer with size 0 discovers size via Len",
		func(t *testing.T) {
			content := "device dev;"
			buf := bytes.NewBufferString(content)
			sbuf, err := stream.NewHashedRunesBufferFromStream(
				"virtual.dml",
				buf,
				0,
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if sbuf.Size() != len(content) {
				t.Errorf(
					"stream.Size() = %d, want %d",
					sbuf.Size(),
					len(content),
				)
			}
			if sbuf.Cap() < 4096 {
				t.Errorf("stream.Cap() = %d, want >= 4096", sbuf.Cap())
			}
		},
	)

	t.Run("os.File with size 0 discovers size via Stat", func(t *testing.T) {
		content := "log \"hello world\";"
		tmpFile, err := os.CreateTemp(t.TempDir(), "file-stat-*.dml")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer func() {
			_ = tmpFile.Close()
		}()
		if _, err := tmpFile.WriteString(content); err != nil {
			t.Fatalf("failed to write to temp file: %v", err)
		}
		if _, err := tmpFile.Seek(0, 0); err != nil {
			t.Fatalf("failed to seek: %v", err)
		}

		sbuf, err := stream.NewHashedRunesBufferFromStream(
			tmpFile.Name(),
			tmpFile,
			0,
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sbuf.Size() != len(content) {
			t.Errorf("stream.Size() = %d, want %d", sbuf.Size(), len(content))
		}
		if sbuf.Cap() < 4096 {
			t.Errorf("stream.Cap() = %d, want >= 4096", sbuf.Cap())
		}
	})

	t.Run(
		"opaque reader with existing file falls back to os.Stat",
		func(t *testing.T) {
			content := "import \"simics\";"
			tmpFile, err := os.CreateTemp(t.TempDir(), "stat-fallback-*.dml")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer func() {
				_ = tmpFile.Close()
			}()
			if _, err := tmpFile.WriteString(content); err != nil {
				t.Fatalf("failed to write to temp file: %v", err)
			}

			r := &opaqueTestReader{r: strings.NewReader(content)}
			sbuf, err := stream.NewHashedRunesBufferFromStream(
				tmpFile.Name(),
				r,
				0,
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if sbuf.Size() != len(content) {
				t.Errorf(
					"stream.Size() = %d, want %d",
					sbuf.Size(),
					len(content),
				)
			}
			if sbuf.Cap() < 4096 {
				t.Errorf("stream.Cap() = %d, want >= 4096", sbuf.Cap())
			}
		},
	)

	t.Run(
		"opaque reader with nonexistent file suppresses stat error",
		func(t *testing.T) {
			content := "parameter p = 1;"
			r := &opaqueTestReader{r: strings.NewReader(content)}
			sbuf, err := stream.NewHashedRunesBufferFromStream(
				"nonexistent_path/virtual.dml",
				r,
				0,
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if sbuf.Size() != len(content) {
				t.Errorf(
					"stream.Size() = %d, want %d",
					sbuf.Size(),
					len(content),
				)
			}
			if sbuf.Cap() < 4096 {
				t.Errorf("stream.Cap() = %d, want >= 4096", sbuf.Cap())
			}
		},
	)
}

func TestNewHashingStream_StreamGrowthRace(t *testing.T) {
	// Simulate race condition: caller passed small size (e.g. 10), but content is 100 runes
	content := strings.Repeat("0123456789", 10)
	r := strings.NewReader(content)
	sbuf, err := stream.NewHashedRunesBufferFromStream("race.dml", r, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sbuf.Size() != 100 {
		t.Errorf("stream.Size() = %d, want 100", sbuf.Size())
	}
	if got := sbuf.GetText(0, 99); got != content {
		t.Errorf("GetText() mismatch: got %q, want %q", got, content)
	}
	if sbuf.Cap() < 4096 {
		t.Errorf("stream.Cap() = %d, want >= 4096", sbuf.Cap())
	}
}

func TestNewHashingStream_BufferCap(t *testing.T) {
	// Size 10000: 10% is 1000 -> 11000 -> aligned to 12288 (3*4096)
	content := strings.Repeat("a", 100)
	sbuf, err := stream.NewHashedRunesBufferFromStream(
		"test.dml",
		strings.NewReader(content),
		10000,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sbuf.Cap() != 12288 {
		t.Errorf("stream.Cap() = %d, want 12288", sbuf.Cap())
	}
}
