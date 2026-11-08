// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"sync"
	"testing"
)

func makeReport(check bool, ref, unch, ign, fail int) *Report {
	r := NewReport(check)
	if ref > 0 {
		r.AddReformatted(ref)
	}
	if unch > 0 {
		r.AddUnchanged(unch)
	}
	if ign > 0 {
		r.AddIgnoredDml12(ign)
	}
	if fail > 0 {
		r.AddFailed(fail)
	}
	return r
}

func TestReport_GettersAndIncrements(t *testing.T) {
	r := NewReport(false)
	if r.Reformatted() != 0 || r.Unchanged() != 0 || r.IgnoredDml12() != 0 ||
		r.Failed() != 0 {
		t.Fatalf("expected all initial counts to be 0")
	}

	r.IncrementReformatted()
	r.IncrementUnchanged()
	r.IncrementIgnoredDml12()
	r.IncrementFailed()

	if r.Reformatted() != 1 {
		t.Errorf("expected Reformatted() == 1, got %d", r.Reformatted())
	}
	if r.Unchanged() != 1 {
		t.Errorf("expected Unchanged() == 1, got %d", r.Unchanged())
	}
	if r.IgnoredDml12() != 1 {
		t.Errorf("expected IgnoredDml12() == 1, got %d", r.IgnoredDml12())
	}
	if r.Failed() != 1 {
		t.Errorf("expected Failed() == 1, got %d", r.Failed())
	}

	r.AddReformatted(5)
	r.AddUnchanged(10)
	r.AddIgnoredDml12(2)
	r.AddFailed(3)

	if r.Reformatted() != 6 {
		t.Errorf("expected Reformatted() == 6, got %d", r.Reformatted())
	}
	if r.Unchanged() != 11 {
		t.Errorf("expected Unchanged() == 11, got %d", r.Unchanged())
	}
	if r.IgnoredDml12() != 3 {
		t.Errorf("expected IgnoredDml12() == 3, got %d", r.IgnoredDml12())
	}
	if r.Failed() != 4 {
		t.Errorf("expected Failed() == 4, got %d", r.Failed())
	}
}

func TestReport_ConcurrentIncrements(t *testing.T) {
	r := NewReport(false)
	var wg sync.WaitGroup

	numGoroutines := 50
	incrementsPerGoroutine := 100

	for range numGoroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range incrementsPerGoroutine {
				r.IncrementReformatted()
				r.IncrementUnchanged()
				r.IncrementIgnoredDml12()
				r.IncrementFailed()

				// Concurrent reads during mutations
				_ = r.Reformatted()
				_ = r.Unchanged()
				_ = r.IgnoredDml12()
				_ = r.Failed()
				_ = r.IsEmpty()
				_ = r.Render(false)
			}
		}()
	}

	wg.Wait()

	expectedTotal := numGoroutines * incrementsPerGoroutine
	if r.Reformatted() != expectedTotal {
		t.Errorf(
			"expected Reformatted() == %d, got %d",
			expectedTotal,
			r.Reformatted(),
		)
	}
	if r.Unchanged() != expectedTotal {
		t.Errorf(
			"expected Unchanged() == %d, got %d",
			expectedTotal,
			r.Unchanged(),
		)
	}
	if r.IgnoredDml12() != expectedTotal {
		t.Errorf(
			"expected IgnoredDml12() == %d, got %d",
			expectedTotal,
			r.IgnoredDml12(),
		)
	}
	if r.Failed() != expectedTotal {
		t.Errorf("expected Failed() == %d, got %d", expectedTotal, r.Failed())
	}
}

func TestReport_IsEmpty(t *testing.T) {
	r := NewReport(false)
	if !r.IsEmpty() {
		t.Errorf("expected empty report to return IsEmpty() == true")
	}

	r.IncrementReformatted()
	if r.IsEmpty() {
		t.Errorf(
			"expected report with Reformatted > 0 to return IsEmpty() == false",
		)
	}
}

func TestReport_Render_Uncolored(t *testing.T) {
	tests := []struct {
		name     string
		report   *Report
		expected string
	}{
		{
			name:   "All four categories plural (user request example)",
			report: makeReport(false, 186, 27, 3, 1),
			expected: "186 files reformatted, 27 files left unchanged, " +
				"3 files ignored due to dml 1.2, 1 file failed to reformat.",
		},
		{
			name:   "All four categories singular",
			report: makeReport(false, 1, 1, 1, 1),
			expected: "1 file reformatted, 1 file left unchanged, " +
				"1 file ignored due to dml 1.2, 1 file failed to reformat.",
		},
		{
			name:   "Check mode with plural and singular",
			report: makeReport(true, 186, 27, 3, 2),
			expected: "186 files would be reformatted, 27 files would be left unchanged, " +
				"3 files ignored due to dml 1.2, 2 files would fail to reformat.",
		},
		{
			name:   "Check mode singular",
			report: makeReport(true, 1, 1, 1, 1),
			expected: "1 file would be reformatted, 1 file would be left unchanged, " +
				"1 file ignored due to dml 1.2, 1 file would fail to reformat.",
		},
		{
			name:     "Only reformatted (skip unchanged, ignored, failed)",
			report:   makeReport(false, 5, 0, 0, 0),
			expected: "5 files reformatted.",
		},
		{
			name:     "Only reformatted singular",
			report:   makeReport(false, 1, 0, 0, 0),
			expected: "1 file reformatted.",
		},
		{
			name:     "Only unchanged (skip reformatted, ignored, failed)",
			report:   makeReport(false, 0, 10, 0, 0),
			expected: "10 files left unchanged.",
		},
		{
			name:     "Only unchanged singular",
			report:   makeReport(false, 0, 1, 0, 0),
			expected: "1 file left unchanged.",
		},
		{
			name:     "Only ignored due to dml 1.2",
			report:   makeReport(false, 0, 0, 4, 0),
			expected: "4 files ignored due to dml 1.2.",
		},
		{
			name:     "Only ignored singular",
			report:   makeReport(false, 0, 0, 1, 0),
			expected: "1 file ignored due to dml 1.2.",
		},
		{
			name:     "Only failed",
			report:   makeReport(false, 0, 0, 0, 2),
			expected: "2 files failed to reformat.",
		},
		{
			name:     "Only failed singular",
			report:   makeReport(false, 0, 0, 0, 1),
			expected: "1 file failed to reformat.",
		},
		{
			name:     "Reformatted and failed (skip unchanged and ignored)",
			report:   makeReport(false, 3, 0, 0, 1),
			expected: "3 files reformatted, 1 file failed to reformat.",
		},
		{
			name:     "Unchanged and ignored (skip reformatted and failed)",
			report:   makeReport(false, 0, 15, 2, 0),
			expected: "15 files left unchanged, 2 files ignored due to dml 1.2.",
		},
		{
			name:     "Empty report",
			report:   NewReport(false),
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.report.Render(false)
			if actual != tc.expected {
				t.Errorf("\nexpected: %q\ngot:      %q", tc.expected, actual)
			}
		})
	}
}

func TestReport_Render_Colored(t *testing.T) {
	t.Run("Standard mode colored", func(t *testing.T) {
		r := makeReport(false, 186, 27, 3, 1)
		expected := "\033[1;34m186 files \033[0m\033[1mreformatted\033[0m, " +
			"\033[34m27 files \033[0mleft unchanged, " +
			"\033[33m3 files \033[0mignored due to dml 1.2, " +
			"\033[31m1 file failed to reformat\033[0m."
		actual := r.Render(true)
		if actual != expected {
			t.Errorf("\nexpected: %q\ngot:      %q", expected, actual)
		}
	})

	t.Run("Check mode colored singular", func(t *testing.T) {
		r := makeReport(true, 1, 1, 1, 1)
		expected := "\033[1;34m1 file \033[0m\033[1mwould be reformatted\033[0m, " +
			"\033[34m1 file \033[0mwould be left unchanged, " +
			"\033[33m1 file \033[0mignored due to dml 1.2, " +
			"\033[31m1 file would fail to reformat\033[0m."
		actual := r.Render(true)
		if actual != expected {
			t.Errorf("\nexpected: %q\ngot:      %q", expected, actual)
		}
	})
}
