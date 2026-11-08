// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"fmt"
	"strings"
	"sync"
)

// Report aggregates the formatting metrics across a dmlfmt run.
type Report struct {
	mu           sync.RWMutex
	Check        bool
	reformatted  int
	unchanged    int
	ignoredDml12 int
	failed       int
}

// NewReport constructs a new Report instance.
func NewReport(check bool) *Report {
	return &Report{
		Check: check,
	}
}

// IncrementReformatted increments the reformatted files counter by 1.
func (r *Report) IncrementReformatted() {
	r.AddReformatted(1)
}

// AddReformatted increments the reformatted files counter by delta.
func (r *Report) AddReformatted(delta int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reformatted += delta
}

// Reformatted returns the number of reformatted files.
func (r *Report) Reformatted() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.reformatted
}

// IncrementUnchanged increments the unchanged files counter by 1.
func (r *Report) IncrementUnchanged() {
	r.AddUnchanged(1)
}

// AddUnchanged increments the unchanged files counter by delta.
func (r *Report) AddUnchanged(delta int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.unchanged += delta
}

// Unchanged returns the number of files left unchanged.
func (r *Report) Unchanged() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.unchanged
}

// IncrementIgnoredDml12 increments the DML 1.2 ignored files counter by 1.
func (r *Report) IncrementIgnoredDml12() {
	r.AddIgnoredDml12(1)
}

// AddIgnoredDml12 increments the DML 1.2 ignored files counter by delta.
func (r *Report) AddIgnoredDml12(delta int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ignoredDml12 += delta
}

// IgnoredDml12 returns the number of files ignored due to DML 1.2.
func (r *Report) IgnoredDml12() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.ignoredDml12
}

// IncrementFailed increments the failed files counter by 1.
func (r *Report) IncrementFailed() {
	r.AddFailed(1)
}

// AddFailed increments the failed files counter by delta.
func (r *Report) AddFailed(delta int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.failed += delta
}

// Failed returns the number of files that failed formatting.
func (r *Report) Failed() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.failed
}

// IsEmpty reports whether no files were tracked in any status category.
func (r *Report) IsEmpty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.reformatted == 0 && r.unchanged == 0 && r.ignoredDml12 == 0 &&
		r.failed == 0
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// Render formats the summary sentence, omitting zero-count categories and applying ANSI colors if enabled.
func (r *Report) Render(colorEnabled bool) string {
	// TODO: AI slop
	r.mu.RLock()
	reformatted := r.reformatted
	unchanged := r.unchanged
	ignoredDml12 := r.ignoredDml12
	failed := r.failed
	check := r.Check
	r.mu.RUnlock()

	var parts []string

	if reformatted > 0 {
		action := "reformatted"
		if check {
			action = "would be reformatted"
		}
		s := plural(reformatted)
		if colorEnabled {
			parts = append(
				parts,
				fmt.Sprintf(
					"%s%d file%s %s%s%s%s",
					ansiBoldBlue,
					reformatted,
					s,
					ansiReset,
					ansiBold,
					action,
					ansiReset,
				),
			)
		} else {
			parts = append(
				parts,
				fmt.Sprintf("%d file%s %s", reformatted, s, action),
			)
		}
	}

	if unchanged > 0 {
		action := "left unchanged"
		if check {
			action = "would be left unchanged"
		}
		s := plural(unchanged)
		if colorEnabled {
			parts = append(
				parts,
				fmt.Sprintf(
					"%s%d file%s %s%s",
					ansiBlue,
					unchanged,
					s,
					ansiReset,
					action,
				),
			)
		} else {
			parts = append(
				parts,
				fmt.Sprintf("%d file%s %s", unchanged, s, action),
			)
		}
	}

	if ignoredDml12 > 0 {
		action := "ignored due to dml 1.2"
		s := plural(ignoredDml12)
		if colorEnabled {
			parts = append(
				parts,
				fmt.Sprintf(
					"%s%d file%s %s%s",
					ansiYellow,
					ignoredDml12,
					s,
					ansiReset,
					action,
				),
			)
		} else {
			parts = append(
				parts,
				fmt.Sprintf("%d file%s %s", ignoredDml12, s, action),
			)
		}
	}

	if failed > 0 {
		action := "failed to reformat"
		if check {
			action = "would fail to reformat"
		}
		s := plural(failed)
		if colorEnabled {
			parts = append(
				parts,
				fmt.Sprintf(
					"%s%d file%s %s%s",
					ansiRed,
					failed,
					s,
					action,
					ansiReset,
				),
			)
		} else {
			parts = append(
				parts,
				fmt.Sprintf("%d file%s %s", failed, s, action),
			)
		}
	}

	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, ", ") + "."
}
