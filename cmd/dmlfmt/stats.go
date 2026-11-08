// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"

	"github.com/antlr4-go/antlr/v4"
)

// StatsCollector collects execution, token, line, and memory metrics for dmlfmt.
type StatsCollector interface {
	RecordWalkIo(spend time.Duration)
	RecordReadIo(spend time.Duration)
	RecordWriteIo(spend time.Duration)
	RecordTokenize(spend time.Duration)
	RecordPrinting(spend time.Duration)
	RecordValidate(spend time.Duration)
	StartFile(filename string)
	RecordReadBytes(bytesRead int64)
	RecordTokensAndLines(tokens []antlr.Token)
	SnapshotMemory(stage string)
	EndFile()
	PrintSummary(w io.Writer)
}

type noopStatsCollector struct{}

func (n *noopStatsCollector) RecordWalkIo(time.Duration)         {}
func (n *noopStatsCollector) RecordReadIo(time.Duration)         {}
func (n *noopStatsCollector) RecordWriteIo(time.Duration)        {}
func (n *noopStatsCollector) RecordTokenize(time.Duration)       {}
func (n *noopStatsCollector) RecordPrinting(time.Duration)       {}
func (n *noopStatsCollector) RecordValidate(time.Duration)       {}
func (n *noopStatsCollector) StartFile(string)                   {}
func (n *noopStatsCollector) RecordReadBytes(int64)              {}
func (n *noopStatsCollector) RecordTokensAndLines([]antlr.Token) {}
func (n *noopStatsCollector) SnapshotMemory(string)              {}
func (n *noopStatsCollector) EndFile()                           {}
func (n *noopStatsCollector) PrintSummary(io.Writer)             {}

type nerdStatsCollector struct {
	mu sync.Mutex

	startTime time.Time
	startMem  runtime.MemStats

	spend struct {
		walk time.Duration
		rd   time.Duration
		wr   time.Duration

		tokenize time.Duration
		printing time.Duration
		validate time.Duration
	}
	peakAlloc   uint64
	totalBytes  int64
	totalTokens int64
	totalLines  int64
}

func newNerdStatsCollector() *nerdStatsCollector {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return &nerdStatsCollector{
		startTime: time.Now(),
		startMem:  m,
		peakAlloc: m.Alloc,
	}
}

func (c *nerdStatsCollector) RecordWalkIo(spend time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.spend.walk += spend
}

func (c *nerdStatsCollector) RecordReadIo(spend time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.spend.rd += spend
}

func (c *nerdStatsCollector) RecordWriteIo(spend time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.spend.wr += spend
}

func (c *nerdStatsCollector) RecordTokenize(spend time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.spend.tokenize += spend
}

func (c *nerdStatsCollector) RecordPrinting(spend time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.spend.printing += spend
}

func (c *nerdStatsCollector) RecordValidate(spend time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.spend.validate += spend
}

func (c *nerdStatsCollector) StartFile(filename string) {
	// Ready for per-file tracking if needed in future multi-worker modes.
}

func (c *nerdStatsCollector) RecordReadBytes(bytesRead int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.totalBytes += bytesRead
}

func (c *nerdStatsCollector) RecordTokensAndLines(tokens []antlr.Token) {
	if len(tokens) == 0 {
		return
	}

	tokenCount := len(tokens)
	totalLines := 0
	if tokenCount > 0 {
		latestToken := tokens[len(tokens)-1]
		totalLines = latestToken.GetLine()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.totalTokens += int64(tokenCount)
	c.totalLines += int64(totalLines)
}

func (c *nerdStatsCollector) SnapshotMemory(stage string) {
	// TODO: AI slop
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	c.mu.Lock()
	defer c.mu.Unlock()
	if m.Alloc > c.peakAlloc {
		c.peakAlloc = m.Alloc
	}
}

func (c *nerdStatsCollector) EndFile() {
	// Ready for per-file finalization.
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	val := float64(b) / float64(div)
	units := []string{"KiB", "MiB", "GiB", "TiB"}
	return fmt.Sprintf("%.2f %s", val, units[exp])
}

func formatReadBytes(b int64) string {
	if b < 1024 {
		return fmt.Sprintf("%d bytes", b)
	}
	return fmt.Sprintf("%d bytes (%s)", b, formatBytes(uint64(b)))
}

func (c *nerdStatsCollector) PrintSummary(w io.Writer) {
	// TODO: AI slop
	c.mu.Lock()
	defer c.mu.Unlock()

	totalDuration := time.Since(c.startTime)
	totalMs := float64(totalDuration.Nanoseconds()) / 1e6

	var finalMem runtime.MemStats
	runtime.ReadMemStats(&finalMem)
	if finalMem.Alloc > c.peakAlloc {
		c.peakAlloc = finalMem.Alloc
	}

	var linesPerMs, tokensPerMs float64
	if totalMs > 0 {
		linesPerMs = float64(c.totalLines) / totalMs
		tokensPerMs = float64(c.totalTokens) / totalMs
	}

	totalIo := c.spend.walk + c.spend.rd + c.spend.wr
	totalIoSeconds := totalIo.Seconds()

	var walkPct, rdPct, wrPct float64
	if totalIo > 0 {
		walkPct = float64(c.spend.walk) / float64(totalIo) * 100
		rdPct = float64(c.spend.rd) / float64(totalIo) * 100
		wrPct = float64(c.spend.wr) / float64(totalIo) * 100
	}

	totalFormat := c.spend.tokenize + c.spend.printing + c.spend.validate
	totalFormatSeconds := totalFormat.Seconds()

	var tokenizePct, printingPct, validatePct float64
	if totalFormat > 0 {
		tokenizePct = float64(c.spend.tokenize) / float64(totalFormat) * 100
		printingPct = float64(c.spend.printing) / float64(totalFormat) * 100
		validatePct = float64(c.spend.validate) / float64(totalFormat) * 100
	}

	allocsDelta := finalMem.Mallocs - c.startMem.Mallocs
	freesDelta := finalMem.Frees - c.startMem.Frees
	gcDelta := finalMem.NumGC - c.startMem.NumGC

	_, _ = fmt.Fprintln(w, "--- Nerd Stats ---")
	_, _ = fmt.Fprintf(
		w,
		"Total read bytes:        %s\n",
		formatReadBytes(c.totalBytes),
	)
	_, _ = fmt.Fprintf(w, "Total processed tokens:  %d\n", c.totalTokens)
	_, _ = fmt.Fprintf(w, "Total processed lines:   %d\n", c.totalLines)
	_, _ = fmt.Fprintf(w, "Total time spent:        %.2f ms\n", totalMs)
	_, _ = fmt.Fprintf(
		w,
		"Total IO time spent:     %.2f s (walk %.1f%%, read %.1f%%, write %.1f%%)\n",
		totalIoSeconds,
		walkPct,
		rdPct,
		wrPct,
	)
	_, _ = fmt.Fprintf(
		w,
		"Total format time spent: %.2f s "+
			"(tokenize %.1f%%, printing %.1f%%, validate %.1f%%)\n",
		totalFormatSeconds,
		tokenizePct,
		printingPct,
		validatePct,
	)
	_, _ = fmt.Fprintf(
		w,
		"Throughput:              %.2f lines/ms, %.2f tokens/ms\n",
		linesPerMs,
		tokensPerMs,
	)
	_, _ = fmt.Fprintf(
		w,
		"GC stats:                %d allocs, %d frees, %d GC cycles\n",
		allocsDelta,
		freesDelta,
		gcDelta,
	)
	_, _ = fmt.Fprintf(
		w,
		"Memory consumption:      final %s, peak %s\n",
		formatBytes(finalMem.Alloc),
		formatBytes(c.peakAlloc),
	)
	_, _ = fmt.Fprintln(w, "------------------")
}
