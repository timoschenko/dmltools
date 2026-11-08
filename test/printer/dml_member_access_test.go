// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestPrintExpressionMemberAccess(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	testCases := []printerTestCase{
		{
			pc:       func(f func()) { f() },
			name:     "01-dot-basic",
			code:     "constant T =   a  .  b;\n",
			expected: "constant T = a.b;\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "02-arrow-basic",
			code:     "constant T = a   ->   b   ;\n",
			expected: "constant T = a->b;\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "03-dot-chained",
			code:     "constant T = a   .   b   .   c;\n",
			expected: "constant T = a.b.c;\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "04-arrow-chained",
			code:     "constant T = a->b->c;\n",
			expected: "constant T = a->b->c;\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "05-mixed-chained",
			code:     "constant T = a.  b   ->c.d;\n",
			expected: "constant T = a.b->c.d;\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "06-messy-whitespace",
			code:     "constant T = a  .  b  ->  c  ;\n",
			expected: "constant T = a.b->c;\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "07-complex-left",
			code:     "constant T = (a+ b  ) .  member;\n",
			expected: "constant T = (a + b).member;\n",
		},
		// {
		// 	pc:       func(f func()) { f() },
		// 	name:     "08-func-call",
		// 	code:     "constant T = func  (  ) . member;\n",
		// 	expected: "constant T = func().member;\n",
		// },
		{
			pc:       func(f func()) { f() },
			name:     "09-array-access",
			code:     "constant T = array[  0  ].member;\n",
			expected: "constant T = array[0].member;\n",
		},
		// {
		// 	pc:       func(f func()) { f() },
		// 	name:     "10-multiline-chain",
		// 	code:     "constant T = object.very_long_member_name_that_should_stay_on_same_line_if_it_fits;\n",
		// 	expected: "constant T = (\n    object\n    .very_long_member_name_that_should_stay_on_same_line_if_it_fits\n);\n",
		// },
		{
			pc:       func(f func()) { f() },
			name:     "11-member-register",
			code:     "constant T =   this  .  object .  register;\n",
			expected: "constant T = this.object.register;\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "12-member-template",
			code:     "constant T =   this  .  template .  register;\n",
			expected: "constant T = this.template.register;\n",
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintExpressionMemberAccess_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	testCases := []printerTestCase{
		{
			pc:       func(f func()) { f() },
			name:     "01-inline-comment",
			code:     "constant T = a   /* comment */  .  b   .  c;\n",
			expected: "constant T = a /* comment */ .b.c;\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "02-inline-comment-after-dot",
			code:     "constant T = a  .  /* comment */   b.  c;\n",
			expected: "constant T = a. /* comment */ b.c;\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "03-inline-comment-both",
			code:     "constant T = a    /* c1 */   ./* c2 */   b  .c;\n",
			expected: "constant T = a /* c1 */ . /* c2 */ b.c;\n",
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintExpressionMemberAccess_Newlines(t *testing.T) {
	t.Parallel()

	const lineCap = 30

	testCases := []printerTestCase{
		{
			pc:       func(f func()) { f() },
			name:     "01-wrap-long-dot",
			code:     "constant T = long_object_name.very_long_member_name;\n",
			expected: "constant T = (\n    long_object_name\n    .very_long_member_name\n);\n",
			masked: MaskSyntheticParens(
				"constant T = long_object_name\n    .very_long_member_name;\n",
			),
		},
		{
			pc:       func(f func()) { f() },
			name:     "02-wrap-long-arrow",
			code:     "constant T = long_object_name->very_long_member_name;\n",
			expected: "constant T = (\n    long_object_name\n    ->very_long_member_name\n);\n",
			masked: MaskSyntheticParens(
				"constant T = long_object_name\n    ->very_long_member_name;\n",
			),
		},
		{
			pc:       func(f func()) { f() },
			name:     "03-wrap-chained-arrow",
			code:     "constant T = a.b.c.d.e.f->g.h->i.j.k.l->n;\n",
			expected: "constant T = (\n    a.b.c.d.e.f->g.h->i.j.k.l\n    ->n\n);\n",
			masked: MaskSyntheticParens(
				"constant T = a.b.c.d.e.f->g.h\n    ->i.j.k.l->n;\n",
			),
		},
		{
			pc:       func(f func()) { f() },
			name:     "04-wrap-chained-dot",
			code:     "constant T = a.b.c.d.e.f->g.h->i.j.k.l.n;\n",
			expected: "constant T = (\n    a.b.c.d.e.f->g.h->i.j.k.l\n    .n\n);\n",
			masked: MaskSyntheticParens(
				"constant T = a.b.c.d.e.f->g.h\n    ->i.j.k.l.n;\n",
			),
		},
		{
			pc:       func(f func()) { f() },
			name:     "04-wrap-chained-dot-w-array",
			code:     "constant T = a.b.c.d.e.f->g.h.i[1].j.k;\n",
			expected: "constant T = (\n    a.b.c.d.e.f->g.h.i[1].j.k\n);\n",
			masked: MaskSyntheticParens(
				"constant T = a.b.c.d.e.f->g.h\n    .i[1].j.k;\n",
			),
		},
		{
			pc:       func(f func()) { f() },
			name:     "05-wrap-chained-dot-w-array",
			code:     "constant T = a.b.c.d.e.f->g.h.i[669].j.k.l;\n",
			expected: "constant T = (\n    a.b.c.d.e.f->g.h.i[669].j\n    .k.l\n);\n",
			masked: MaskSyntheticParens(
				"constant T = a.b.c.d.e.f->g.h\n    .i[669].j.k.l;\n",
			),
		},
		{
			pc:       func(f func()) { f() },
			name:     "06-wrap-chained-dot-inside-array",
			code:     "constant T = werylongname[a.b.c.d.e.f.g.h].i.k;\n",
			expected: "constant T = werylongname[\n    a.b.c.d.e.f.g.h\n].i.k;\n",
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
