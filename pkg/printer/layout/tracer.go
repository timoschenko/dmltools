// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package layout

import (
	"bytes"
	"encoding/json"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/timoschenko/dmltools/pkg/printer/format"
)

// TokenInfo captures the full state of a token and its formatting flag.
type TokenInfo struct {
	Text       string            `json:"text,omitempty"`
	TokenType  int               `json:"tokenType,omitempty"`
	Line       int               `json:"line,omitempty"`
	Column     int               `json:"column,omitempty"`
	TokenIndex int               `json:"tokenIndex,omitempty"`
	FormatFlag format.FormatFlag `json:"formatFlag"`
	GroupPath  string            `json:"groupPath,omitempty"`
}

// TraceRecord contains information about the formatting decision made between two tokens.
type TraceRecord struct {
	Left          TokenInfo        `json:"left"`
	Right         TokenInfo        `json:"right"`
	ResolvedSep   format.Separator `json:"resolvedSep"`
	BreakPriority format.Priority  `json:"breakPriority"`
	Winner        format.Winner    `json:"winner,omitempty"` // e.g., 1 (left), 2 (right), 3 (equal)
}

// TokenTracer is an optional interface to trace token layout resolutions.
type TokenTracer interface {
	Record(record *TraceRecord)
}

// TraceEntry wraps a TraceRecord with system information (timestamp, PID, GID)
// formatted for Chrome's trace event format.
type TraceEntry struct {
	Name      string       `json:"name"`
	Category  string       `json:"cat"`
	Phase     string       `json:"ph"`
	Scope     string       `json:"s,omitempty"`
	Timestamp int64        `json:"ts"` // Microseconds
	PID       int          `json:"pid"`
	TID       uint64       `json:"tid"`
	Args      *TraceRecord `json:"args"`
}

// JSONTokenTracer is a basic tracer that collects trace records into memory
// and can serialize them to JSON.
type JSONTokenTracer struct {
	mu        sync.Mutex
	Entries   []TraceEntry
	startTime time.Time
	pid       int
}

func NewJSONTokenTracer() *JSONTokenTracer {
	return &JSONTokenTracer{
		Entries:   make([]TraceEntry, 0),
		startTime: time.Now(),
		pid:       os.Getpid(),
	}
}

func (t *JSONTokenTracer) Record(record *TraceRecord) {
	if record == nil {
		panic("TraceRecord cannot be nil")
	}
	ts := time.Since(t.startTime).Microseconds() // make sure to capture it asap
	tid := getGID()
	entry := TraceEntry{
		Name:      "TokenFormatDecision",
		Category:  "formatter",
		Phase:     "I",
		Scope:     "t",
		Timestamp: ts,
		PID:       t.pid,
		TID:       tid,
		Args:      record,
	}
	t.mu.Lock()
	t.Entries = append(t.Entries, entry)
	t.mu.Unlock()
}

func (t *JSONTokenTracer) ToJSON() ([]byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return json.MarshalIndent(t.Entries, "", "  ")
}

// getGID parses the current goroutine ID from runtime.Stack.
func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
