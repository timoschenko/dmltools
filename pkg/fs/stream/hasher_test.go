// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package stream

import (
	"errors"
	"hash"
	"io"
	"testing"
)

type stubHash64 struct {
	hash.Hash64
	writeFunc func(p []byte) (int, error)
}

func (s *stubHash64) Write(p []byte) (int, error) {
	if s.writeFunc != nil {
		return s.writeFunc(p)
	}
	return s.Hash64.Write(p)
}

func TestFileHasher_WriteErrors(t *testing.T) {
	errCustom := errors.New("custom write error")

	t.Run("fnv error", func(t *testing.T) {
		hasher := NewFileHasher()
		hasher.fnv = &stubHash64{
			Hash64: hasher.fnv,
			writeFunc: func(p []byte) (int, error) {
				return 3, errCustom
			},
		}

		n, err := hasher.Write([]byte("hello world"))
		if !errors.Is(err, errCustom) {
			t.Fatalf("expected errCustom, got %v", err)
		}
		if n != 3 {
			t.Fatalf("expected n=3, got %d", n)
		}
		if hasher.Sum().Size != 0 {
			t.Fatalf("expected size to remain 0, got %d", hasher.Sum().Size)
		}
	})

	t.Run("fnv short write", func(t *testing.T) {
		hasher := NewFileHasher()
		hasher.fnv = &stubHash64{
			Hash64: hasher.fnv,
			writeFunc: func(p []byte) (int, error) {
				return len(p) - 1, nil
			},
		}

		n, err := hasher.Write([]byte("hello world"))
		if !errors.Is(err, io.ErrShortWrite) {
			t.Fatalf("expected io.ErrShortWrite, got %v", err)
		}
		if n != len("hello world")-1 {
			t.Fatalf("expected n=%d, got %d", len("hello world")-1, n)
		}
		if hasher.Sum().Size != 0 {
			t.Fatalf("expected size to remain 0, got %d", hasher.Sum().Size)
		}
	})

	t.Run("crc error", func(t *testing.T) {
		hasher := NewFileHasher()
		hasher.crc = &stubHash64{
			Hash64: hasher.crc,
			writeFunc: func(p []byte) (int, error) {
				return 4, errCustom
			},
		}

		n, err := hasher.Write([]byte("hello world"))
		if !errors.Is(err, errCustom) {
			t.Fatalf("expected errCustom, got %v", err)
		}
		if n != 4 {
			t.Fatalf("expected n=4, got %d", n)
		}
		if hasher.Sum().Size != 0 {
			t.Fatalf("expected size to remain 0, got %d", hasher.Sum().Size)
		}
	})

	t.Run("crc short write", func(t *testing.T) {
		hasher := NewFileHasher()
		hasher.crc = &stubHash64{
			Hash64: hasher.crc,
			writeFunc: func(p []byte) (int, error) {
				return len(p) - 2, nil
			},
		}

		n, err := hasher.Write([]byte("hello world"))
		if !errors.Is(err, io.ErrShortWrite) {
			t.Fatalf("expected io.ErrShortWrite, got %v", err)
		}
		if n != len("hello world")-2 {
			t.Fatalf("expected n=%d, got %d", len("hello world")-2, n)
		}
		if hasher.Sum().Size != 0 {
			t.Fatalf("expected size to remain 0, got %d", hasher.Sum().Size)
		}
	})
}

func TestFileHasher_WriteString(t *testing.T) {
	input := "hello world \u2764 \U0001F680"

	h1 := NewFileHasher()
	n1, err1 := h1.Write([]byte(input))
	if err1 != nil {
		t.Fatalf("Write failed: %v", err1)
	}
	if n1 != len(input) {
		t.Fatalf("Write returned %d, want %d", n1, len(input))
	}

	h2 := NewFileHasher()
	n2, err2 := h2.WriteString(input)
	if err2 != nil {
		t.Fatalf("WriteString failed: %v", err2)
	}
	if n2 != len(input) {
		t.Fatalf("WriteString returned %d, want %d", n2, len(input))
	}

	if h1.Sum() != h2.Sum() {
		t.Errorf(
			"Sum() mismatch: Write=%+v, WriteString=%+v",
			h1.Sum(),
			h2.Sum(),
		)
	}
}
