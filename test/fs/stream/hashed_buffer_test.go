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

	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/fs"
	"github.com/timoschenko/dmltools/pkg/fs/stream"
)

func TestCalculateStreamBufferSize(t *testing.T) {
	tests := []struct {
		name     string
		size     int
		expected int
	}{
		{
			name:     "size 0 returns 4KB",
			size:     0,
			expected: 4096,
		},
		{
			name:     "negative size returns 4KB",
			size:     -10,
			expected: 4096,
		},
		{
			name:     "small size 16 aligns to 4KB",
			size:     16,
			expected: 4096, // 16 + 1 = 17 -> aligned to 4096
		},
		{
			name:     "boundary size 3724 aligns to 4KB",
			size:     3724, // 3724 + 372 = 4096 -> aligned to 4096
			expected: 4096,
		},
		{
			name:     "boundary size 3725 aligns to 8KB",
			size:     3725, // 3725 + 372 = 4097 -> aligned to 8192
			expected: 8192,
		},
		{
			name:     "size 10000 aligns to 12KB",
			size:     10000, // 10000 + 1000 = 11000 -> 12288 (3*4096)
			expected: 12288,
		},
		{
			name:     "size exactly at 64KB grow limit (655360)",
			size:     655360, // grow = 65536, growed = 720896 = 176 * 4096
			expected: 720896,
		},
		{
			name:     "large size 1MB caps grow at 64KB",
			size:     1000000, // 10% is 100000 > 65536, growed = 1065536 -> aligned to 1069056
			expected: 1069056,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := fs.CalculateStreamBufferSize(tc.size)
			if got != tc.expected {
				t.Errorf(
					"calculateStreamBufferSize(%d) = %d, want %d",
					tc.size,
					got,
					tc.expected,
				)
			}
		})
	}
}

type opaqueReader struct {
	r io.Reader
}

func (o *opaqueReader) Read(p []byte) (int, error) {
	return o.r.Read(p)
}

func TestSizeFromReader(t *testing.T) {
	t.Run("strings.Reader", func(t *testing.T) {
		r := strings.NewReader("hello")
		size, ok := fs.SizeFromReader(r)
		if !ok || size != 5 {
			t.Errorf(
				"fs.SizeFromReader(strings.Reader) = (%d, %v), want (5, true)",
				size,
				ok,
			)
		}
	})

	t.Run("bytes.Buffer", func(t *testing.T) {
		buf := bytes.NewBufferString("foobar")
		size, ok := fs.SizeFromReader(buf)
		if !ok || size != 6 {
			t.Errorf(
				"fs.SizeFromReader(bytes.Buffer) = (%d, %v), want (6, true)",
				size,
				ok,
			)
		}
	})

	t.Run("bytes.Reader", func(t *testing.T) {
		r := bytes.NewReader([]byte("test byte reader"))
		size, ok := fs.SizeFromReader(r)
		if !ok || size != 16 {
			t.Errorf(
				"fs.SizeFromReader(bytes.Reader) = (%d, %v), want (16, true)",
				size,
				ok,
			)
		}
	})

	t.Run("os.File", func(t *testing.T) {
		tmpFile, err := os.CreateTemp(t.TempDir(), "stat-test-*.dml")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer func() {
			_ = tmpFile.Close()
		}()
		content := "dml 1.4;\ndevice dev;"
		if _, err := tmpFile.WriteString(content); err != nil {
			t.Fatalf("failed to write to temp file: %v", err)
		}

		size, ok := fs.SizeFromReader(tmpFile)
		if !ok || size != len(content) {
			t.Errorf(
				"fs.SizeFromReader(os.File) = (%d, %v), want (%d, true)",
				size,
				ok,
				len(content),
			)
		}
	})

	t.Run("directory non-regular file", func(t *testing.T) {
		dir, err := os.Open(t.TempDir())
		if err != nil {
			t.Fatalf("failed to open temp dir: %v", err)
		}
		defer func() {
			_ = dir.Close()
		}()
		size, ok := fs.SizeFromReader(dir)
		if ok || size != 0 {
			t.Errorf(
				"fs.SizeFromReader(dir) = (%d, %v), want (0, false)",
				size,
				ok,
			)
		}
	})

	t.Run("pipe non-regular file", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("failed to create pipe: %v", err)
		}
		defer func() {
			_ = r.Close()
			_ = w.Close()
		}()
		size, ok := fs.SizeFromReader(r)
		if ok || size != 0 {
			t.Errorf(
				"fs.SizeFromReader(pipe) = (%d, %v), want (0, false)",
				size,
				ok,
			)
		}
	})

	t.Run("opaque reader without sizing interface", func(t *testing.T) {
		r := &opaqueReader{r: strings.NewReader("hidden")}
		size, ok := fs.SizeFromReader(r)
		if ok {
			t.Errorf(
				"fs.SizeFromReader(opaqueReader) = (%d, %v), want (0, false)",
				size,
				ok,
			)
		}
	})

	t.Run("nil reader", func(t *testing.T) {
		size, ok := fs.SizeFromReader(nil)
		if ok || size != 0 {
			t.Errorf(
				"fs.SizeFromReader(nil) = (%d, %v), want (0, false)",
				size,
				ok,
			)
		}
	})
}

func TestHashedBuffer_LexerParity(t *testing.T) {
	dmlSample := `
// This is a test file with Unicode comments: 🚀 and 世界
dml 1.4;

device test_dev {
	param desc = "Multi-byte string: Привет мир, 日本語, 🌍!";
	register r0 size 4 @ 0x1000 {
		method read() -> (uint32) {
			log "Reading register r0";
			return 0xCAFE;
		}
	}
}
`
	// 1. Lex with HashedRunesBuffer
	runesStream, err := stream.NewHashedRunesBufferFromStream(
		"sample.dml",
		strings.NewReader(dmlSample),
		len(dmlSample),
	)
	if err != nil {
		t.Fatalf("failed to create HashedRunesBuffer: %v", err)
	}
	runesLexer := dmlex.NewdmlLexer(runesStream)
	runesTokens := runesLexer.GetAllTokens()

	// 2. Lex with HashedBuffer
	byteStream, err := stream.NewHashedBufferFromStream(
		"sample.dml",
		strings.NewReader(dmlSample),
		len(dmlSample),
	)
	if err != nil {
		t.Fatalf("failed to create HashedBuffer: %v", err)
	}
	byteLexer := dmlex.NewdmlLexer(byteStream)
	byteTokens := byteLexer.GetAllTokens()

	// 3. Verify hashes match
	if runesStream.Hash() != byteStream.Hash() {
		t.Errorf(
			"Stream hashes do not match: runes=%+v, bytes=%+v",
			runesStream.Hash(),
			byteStream.Hash(),
		)
	}

	// 4. Verify token counts match
	if len(runesTokens) != len(byteTokens) {
		t.Fatalf(
			"Token counts differ: runes got %d tokens, bytes got %d tokens",
			len(runesTokens),
			len(byteTokens),
		)
	}

	// 5. Verify each token has identical type, text, line, and column
	for i := range runesTokens {
		rt := runesTokens[i]
		bt := byteTokens[i]

		if rt.GetTokenType() != bt.GetTokenType() {
			t.Errorf(
				"Token [%d] Type mismatch: runes=%d (%s), bytes=%d (%s)",
				i,
				rt.GetTokenType(),
				dmlex.TokenType(rt.GetTokenType()).String(),
				bt.GetTokenType(),
				dmlex.TokenType(bt.GetTokenType()).String(),
			)
		}

		// HashedBuffer preserves raw source bytes, whereas HashedRunesBuffer converts
		// invalid UTF-8 bytes to RuneError (U+FFFD) at construction time.
		// When decoded to runes, both streams represent the exact same rune sequence.
		if rt.GetText() != bt.GetText() {
			decodedRunesText := string([]rune(bt.GetText()))
			if decodedRunesText != rt.GetText() {
				t.Errorf(
					"Token [%d] Text mismatch: runes=%q, bytes raw=%q, bytes decoded=%q",
					i,
					rt.GetText(),
					bt.GetText(),
					decodedRunesText,
				)
			}
		}

		if rt.GetLine() != bt.GetLine() {
			t.Errorf(
				"Token [%d] Line mismatch: runes=%d, bytes=%d",
				i,
				rt.GetLine(),
				bt.GetLine(),
			)
		}

		if rt.GetColumn() != bt.GetColumn() {
			t.Errorf(
				"Token [%d] Column mismatch: runes=%d, bytes=%d",
				i,
				rt.GetColumn(),
				bt.GetColumn(),
			)
		}
	}
}

func TestHashedBuffer_MalformedRunes_LexerParity(t *testing.T) {
	dmlSample := "// Comment with truncated UTF-8 prefix: \xE4\xB8\n" +
		"dml 1.4;\n\n" +
		"/* Block comment with illegal bytes: \xFF \xFE and stray \x80 */\n" +
		"device test_device {\n" +
		"\tparam desc = \"String with malformed byte \xFF and overlong \xC0\x80 inside\";\n" +
		"\tparam desc2 = \"Another string \xED\xBF\xBF end\";\n" +
		"\tregister r0 size 4 @ 0x1000 {\n" +
		"\t\tmethod read() -> (uint32) {\n" +
		"\t\t\t// In-method comment \x80\xBF\x90\n" +
		"\t\t\tlog \"Log with \xF4\x90\x80\x80\";\n" +
		"\t\t\treturn 0xCAFE;\n" +
		"\t\t}\n" +
		"\t}\n" +
		"}\n"

	// 1. Lex with HashedRunesBuffer
	runesStream, err := stream.NewHashedRunesBufferFromStream(
		"malformed_sample.dml",
		strings.NewReader(dmlSample),
		len(dmlSample),
	)
	if err != nil {
		t.Fatalf("failed to create HashedRunesBuffer: %v", err)
	}
	runesLexer := dmlex.NewdmlLexer(runesStream)
	runesTokens := runesLexer.GetAllTokens()

	// 2. Lex with HashedBuffer
	byteStream, err := stream.NewHashedBufferFromStream(
		"malformed_sample.dml",
		strings.NewReader(dmlSample),
		len(dmlSample),
	)
	if err != nil {
		t.Fatalf("failed to create HashedBuffer: %v", err)
	}
	byteLexer := dmlex.NewdmlLexer(byteStream)
	byteTokens := byteLexer.GetAllTokens()

	// 3. Verify hashes match
	if runesStream.Hash() != byteStream.Hash() {
		t.Errorf(
			"Stream hashes do not match: runes=%+v, bytes=%+v",
			runesStream.Hash(),
			byteStream.Hash(),
		)
	}

	// 4. Verify token counts match
	if len(runesTokens) != len(byteTokens) {
		t.Fatalf(
			"Token counts differ: runes got %d tokens, bytes got %d tokens",
			len(runesTokens),
			len(byteTokens),
		)
	}

	// 5. Verify each token has identical type, text, line, and column
	for i := range runesTokens {
		rt := runesTokens[i]
		bt := byteTokens[i]

		if rt.GetTokenType() != bt.GetTokenType() {
			t.Errorf(
				"Token [%d] Type mismatch: runes=%d (%s), bytes=%d (%s)",
				i,
				rt.GetTokenType(),
				dmlex.TokenType(rt.GetTokenType()).String(),
				bt.GetTokenType(),
				dmlex.TokenType(bt.GetTokenType()).String(),
			)
		}

		// HashedBuffer preserves raw source bytes, whereas HashedRunesBuffer converts
		// invalid UTF-8 bytes to RuneError (U+FFFD) at construction time.
		// When decoded to runes, both streams represent the exact same rune sequence.
		if rt.GetText() != bt.GetText() {
			decodedRunesText := string([]rune(bt.GetText()))
			if decodedRunesText != rt.GetText() {
				t.Errorf(
					"Token [%d] Text mismatch: runes=%q, bytes raw=%q, bytes decoded=%q",
					i,
					rt.GetText(),
					bt.GetText(),
					decodedRunesText,
				)
			}
		}

		if rt.GetLine() != bt.GetLine() {
			t.Errorf(
				"Token [%d] Line mismatch: runes=%d, bytes=%d",
				i,
				rt.GetLine(),
				bt.GetLine(),
			)
		}

		if rt.GetColumn() != bt.GetColumn() {
			t.Errorf(
				"Token [%d] Column mismatch: runes=%d, bytes=%d",
				i,
				rt.GetColumn(),
				bt.GetColumn(),
			)
		}
	}
}

func TestHashedBuffer_LA_Parity(t *testing.T) {
	sample := "dml 1.4;\n// Unicode comment: 🚀 and 世界\nparam desc = \"Привет мир, 🌍!\";"
	runesBuf, err := stream.NewHashedRunesBufferFromStream(
		"sample.dml",
		strings.NewReader(sample),
		len(sample),
	)
	if err != nil {
		t.Fatalf("failed to create HashedRunesBuffer: %v", err)
	}
	byteBuf, err := stream.NewHashedBufferFromStream(
		"sample.dml",
		strings.NewReader(sample),
		len(sample),
	)
	if err != nil {
		t.Fatalf("failed to create HashedBuffer: %v", err)
	}

	offsets := []int{-1, 0, 1, 2, 3, 4, 5, 10}
	step := 0
	for {
		for _, off := range offsets {
			rLA := runesBuf.LA(off)
			bLA := byteBuf.LA(off)
			if rLA != bLA {
				t.Fatalf(
					"step %d, LA(%d) mismatch: runes=%d (%U), bytes=%d (%U)",
					step,
					off,
					rLA,
					rLA,
					bLA,
					bLA,
				)
			}
		}
		if byteBuf.LA(1) == antlr.TokenEOF {
			break
		}
		runesBuf.Consume()
		byteBuf.Consume()
		step++
	}
}

func makeBenchmarkPayload(kbSize int) string {
	base := `// Benchmark payload with Unicode comments: 🚀 and 世界
dml 1.4;
device test_device {
	param desc = "Hello 世界 🌍 Привет";
	register reg size 4 @ 0x1000 {
		method read() -> (uint32) {
			log "read called";
			return 42;
		}
	}
}
`
	repeatCount := max(1, (kbSize*1024)/len(base))
	var sb strings.Builder
	for range repeatCount {
		sb.WriteString(base)
	}
	return sb.String()
}

func BenchmarkBufferInitialization_HashedRunesBuffer_1MB(b *testing.B) {
	payload := makeBenchmarkPayload(1024)
	size := len(payload)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r := strings.NewReader(payload)
		buf, err := stream.NewHashedRunesBufferFromStream("bench.dml", r, size)
		if err != nil {
			b.Fatal(err)
		}
		_ = buf.Size()
	}
}

func BenchmarkBufferInitialization_HashedBuffer_1MB(b *testing.B) {
	payload := makeBenchmarkPayload(1024)
	size := len(payload)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r := strings.NewReader(payload)
		buf, err := stream.NewHashedBufferFromStream("bench.dml", r, size)
		if err != nil {
			b.Fatal(err)
		}
		_ = buf.Size()
	}
}

func BenchmarkLexer_HashedRunesBuffer(b *testing.B) {
	payload := makeBenchmarkPayload(100)
	size := len(payload)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r := strings.NewReader(payload)
		buf, err := stream.NewHashedRunesBufferFromStream("bench.dml", r, size)
		if err != nil {
			b.Fatal(err)
		}
		lexer := dmlex.NewdmlLexer(buf)
		for {
			tok := lexer.NextToken()
			if tok.GetTokenType() == antlr.TokenEOF {
				break
			}
		}
	}
}

func BenchmarkLexer_HashedBuffer(b *testing.B) {
	payload := makeBenchmarkPayload(100)
	size := len(payload)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r := strings.NewReader(payload)
		buf, err := stream.NewHashedBufferFromStream("bench.dml", r, size)
		if err != nil {
			b.Fatal(err)
		}
		lexer := dmlex.NewdmlLexer(buf)
		for {
			tok := lexer.NextToken()
			if tok.GetTokenType() == antlr.TokenEOF {
				break
			}
		}
	}
}
