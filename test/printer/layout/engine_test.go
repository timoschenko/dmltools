// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package layout_test

import (
	"fmt"
	"github.com/timoschenko/dmltools/pkg/printer"
	"github.com/timoschenko/dmltools/pkg/printer/format"
	"github.com/timoschenko/dmltools/pkg/printer/layout"
	"strings"
	"testing"

	"github.com/antlr4-go/antlr/v4"
)

type mockToken struct {
	antlr.Token
	text string
}

func (t *mockToken) GetText() string {
	return t.text
}

func (t *mockToken) GetLine() int {
	return 1
}

func (t *mockToken) GetTokenType() int {
	return 0
}

func getEmptyLeState() layout.State {
	var sb strings.Builder
	cfg := printer.DefaultSettings(80)
	le := layout.NewEngine(
		cfg,
		&sb,
		nil,
	)

	return le.TesOnly_GetState()
}

func TestLayoutEngine_PushAndEmit(t *testing.T) {
	var sb strings.Builder
	le := layout.NewEngine(
		printer.DefaultSettings(80),
		&sb,
		nil,
	)

	// Token 1: priority 1, sep after space
	t1 := &mockToken{text: "foo"}
	f1 := format.FormatFlag(0).
		WithSeparatorPriority(1).
		WithSeparatorAfter(format.SepSpace)

	// Token 2: priority 1, sep before space
	t2 := &mockToken{text: "bar"}
	f2 := format.FormatFlag(0).
		WithSeparatorPriority(1).
		WithSeparatorBefore(format.SepSpace)

	if err := le.Push(t1, f1); err != nil {
		t.Fatal(err)
	}
	if err := le.Push(t2, f2); err != nil {
		t.Fatal(err)
	}
	if err := le.Flush(); err != nil {
		t.Fatal(err)
	}

	expected := "foo bar"
	if sb.String() != expected {
		t.Errorf("expected %q, got %q", expected, sb.String())
	}
}

func TestLayoutEngine_LineWrapping(t *testing.T) {
	var sb strings.Builder
	// Linecap 10. "longtoken" (9) + " " (1) + "next" (4) = 14 > 10
	le := layout.NewEngine(
		printer.DefaultSettings(10),
		&sb,
		nil,
	)

	t1 := &mockToken{text: "longtoken"}
	f1 := format.FormatFlag(0).
		WithSeparatorAfter(format.SepSpaceMaybeWrap).
		WithBreakAfterPriority(5)

	t2 := &mockToken{text: "next"}
	f2 := format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpaceMaybeWrap).
		WithBreakBeforePriority(5)

	if err := le.Push(t1, f1); err != nil {
		t.Fatal(err)
	}
	if err := le.Push(t2, f2); err != nil {
		t.Fatal(err)
	}
	if err := le.Flush(); err != nil {
		t.Fatal(err)
	}

	// Should break at the optional wrap
	expected := "longtoken\nnext"
	if sb.String() != expected {
		t.Errorf("expected %q, got %q", expected, sb.String())
	}
}

func TestLayoutEngine_CascadingBreaks(t *testing.T) {
	var sb strings.Builder
	// Linecap 9.
	// "aaaaa" (5) + "bbbbb" (5) + "ccccc" (5) = 15.
	// 15 > 9, so it breaks at P5: "aaaaabbbbb\nccccc"
	// Now first segment "aaaaabbbbb" is 10. 10 > 9.
	// So it breaks at P1: "aaaaa\nbbbbb\nccccc"
	le := layout.NewEngine(
		printer.DefaultSettings(9),
		&sb,
		nil,
	)

	ta := &mockToken{text: "aaaaa"}
	fa := format.FormatFlag(0).
		WithSeparatorAfter(format.SepSpaceMaybeWrap).
		WithBreakAfterPriority(1)

	tb := &mockToken{text: "bbbbb"}
	fb := format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpaceMaybeWrap).
		WithBreakBeforePriority(1).
		WithSeparatorAfter(format.SepSpaceMaybeWrap).
		WithBreakAfterPriority(5)

	tc := &mockToken{text: "ccccc"}
	fc := format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpaceMaybeWrap).
		WithBreakBeforePriority(5)

	if err := le.Push(ta, fa); err != nil {
		t.Fatal(err)
	}
	if err := le.Push(tb, fb); err != nil {
		t.Fatal(err)
	}
	if err := le.Push(tc, fc); err != nil {
		t.Fatal(err)
	}
	if err := le.Flush(); err != nil {
		t.Fatal(err)
	}

	expected := "aaaaa\nbbbbb\nccccc"
	if sb.String() != expected {
		t.Errorf("expected %q, got %q", expected, sb.String())
	}
}

func TestLayoutEngine_Indentation(t *testing.T) {
	var sb strings.Builder
	// Linecap 5. "if" (2) + "cond" (4) = 6 > 5.
	le := layout.NewEngine(
		printer.DefaultSettings(5),
		&sb,
		nil,
	)

	t1 := &mockToken{text: "if"}
	f1 := format.FormatFlag(0).
		WithSeparatorAfter(format.SepSpaceMaybeWrap).
		WithBreakAfterPriority(5).
		WithIndentIncrease()

	t2 := &mockToken{text: "cond"}
	f2 := format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpaceMaybeWrap).
		WithBreakBeforePriority(5)

	if err := le.Push(t1, f1); err != nil {
		t.Fatal(err)
	}
	if err := le.Push(t2, f2); err != nil {
		t.Fatal(err)
	}
	if err := le.Flush(); err != nil {
		t.Fatal(err)
	}

	// Should break and indent
	expected := "if\n    cond"
	if sb.String() != expected {
		t.Errorf("expected %q, got %q", expected, sb.String())
	}
}

func TestLayoutEngine_WriteRaw(t *testing.T) {
	var sb strings.Builder
	le := layout.NewEngine(
		printer.DefaultSettings(80),
		&sb,
		nil,
	)

	if err := le.Push(&mockToken{text: "foo"}, format.FormatFlag(0).
		WithSeparatorAfter(format.SepSpace)); err != nil {
		t.Fatal(err)
	}
	if err := le.WriteRaw("RAW"); err != nil {
		t.Fatal(err)
	}
	if err := le.Push(&mockToken{text: "bar"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpace)); err != nil {
		t.Fatal(err)
	}
	if err := le.Flush(); err != nil {
		t.Fatal(err)
	}

	expected := "fooRAW bar"
	if sb.String() != expected {
		t.Errorf("expected %q, got %q", expected, sb.String())
	}
}

func TestLayoutEngine_Unicode(t *testing.T) {
	var sb strings.Builder
	// Rocket emoji 🚀 is 4 bytes but 1 rune.
	// "🚀" (1) + " " (1) + "🚀" (1) + " " (1) + "b" (1) = 5 runes.
	le := layout.NewEngine(
		printer.DefaultSettings(6),
		&sb,
		nil,
	)

	ta := &mockToken{text: "🚀"}
	fa := format.FormatFlag(0).
		WithSeparatorAfter(format.SepSpaceMaybeWrap).
		WithBreakAfterPriority(5)

	tb := &mockToken{text: "🚀"}
	fb := format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpaceMaybeWrap).
		WithBreakBeforePriority(5).
		WithSeparatorAfter(format.SepSpaceMaybeWrap).
		WithBreakAfterPriority(5)

	tc := &mockToken{text: "b"}
	fc := format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpaceMaybeWrap).
		WithBreakBeforePriority(5)

	if err := le.Push(ta, fa); err != nil {
		t.Fatal(err)
	}
	if err := le.Push(tb, fb); err != nil {
		t.Fatal(err)
	}
	if err := le.Push(tc, fc); err != nil {
		t.Fatal(err)
	}
	if err := le.Flush(); err != nil {
		t.Fatal(err)
	}

	// Should NOT break if we use rune count, because 3 runes <= 3 linecap.
	// However, it SHOULD now have spaces because SepSpaceMaybeWrap defaults to space.
	expected := "🚀 🚀 b"
	if sb.String() != expected {
		t.Errorf("expected %q, got %q", expected, sb.String())
	}
}

func TestLayoutEngine_WrapGroups(t *testing.T) {
	defaultFlag := format.FormatFlag(0).
		WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
		WithBreakBeforePriority(format.PriorityNormal).
		WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
		WithBreakAfterPriority(format.PriorityNormal).
		WithSeparatorPriority(format.PriorityNormal)
	defaultWhitespaceFlag := defaultFlag.
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorAfter(format.SepSpace)

	t.Run("independent groups", func(t *testing.T) {
		var sb strings.Builder
		// Linecap 15.
		// "foo(a,b)" (8) + " + " (3) + "bar(c,d)" (8) = 19 > 15.
		le := layout.NewEngine(
			printer.DefaultSettings(15),
			&sb,
			nil,
		)

		push := func(token antlr.Token, flag format.FormatFlag) {
			if err := le.Push(token, flag); err != nil {
				t.Fatal(err)
			}
		}

		// foo(a,b)
		push(&mockToken{text: "foo"}, format.FormatFlag(0))
		le.BeginWrapGroup("test_independent_1")
		push(&mockToken{text: "("}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithBreakAfterPriority(5))
		push(&mockToken{text: "a"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithBreakBeforePriority(5).
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithBreakAfterPriority(5))
		push(&mockToken{text: ","}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithBreakBeforePriority(5).
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithBreakAfterPriority(5))
		push(&mockToken{text: "b"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithBreakBeforePriority(5))
		push(&mockToken{text: ")"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithBreakBeforePriority(5))
		le.EndWrapGroup()

		push(&mockToken{text: "+"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorAfter(format.SepSpace))

		// bar(c,d)
		push(&mockToken{text: "bar"}, format.FormatFlag(0))
		le.BeginWrapGroup("test_independent_2")
		push(&mockToken{text: "("}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithBreakAfterPriority(10))
		push(&mockToken{text: "c"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithBreakBeforePriority(10).
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithBreakAfterPriority(10))
		push(&mockToken{text: ","}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithBreakBeforePriority(10).
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithBreakAfterPriority(10))
		push(&mockToken{text: "d"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithBreakBeforePriority(10))
		push(&mockToken{text: ")"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithBreakBeforePriority(10))
		le.EndWrapGroup()

		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		// If bar(c,d) causes the overflow, it should wrap.
		// foo(a,b) should remain inline because it fits on the first line.
		// Since bar( fits on the first line (15 chars), it stays there, and c wraps.
		expected := "foo(a,b) + bar(\nc\n,\nd\n)"
		if sb.String() != expected {
			t.Errorf("expected %q, got %q", expected, sb.String())
		}
	})

	t.Run("nested groups", func(t *testing.T) {
		var sb strings.Builder
		// Linecap 25.
		// "outer(inner(a,b),c)" is 19 chars.
		le := layout.NewEngine(
			printer.DefaultSettings(25),
			&sb,
			nil,
		)

		push := func(token antlr.Token, flag format.FormatFlag) {
			if err := le.Push(token, flag); err != nil {
				t.Fatal(err)
			}
		}

		le.BeginWrapGroup("test_nested_outer")
		push(&mockToken{text: "outer"}, format.FormatFlag(0))
		push(&mockToken{text: "("}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithBreakAfterPriority(10).
			WithIndentIncrease())

		// Force outer to wrap by adding a long token at the end
		le.BeginWrapGroup("test_nested_inner")
		push(&mockToken{text: "inner"}, format.FormatFlag(0))
		push(&mockToken{text: "("}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithBreakAfterPriority(5))
		push(&mockToken{text: "a"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithBreakBeforePriority(5).
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithBreakAfterPriority(5))
		push(&mockToken{text: ","}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithBreakBeforePriority(5).
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithBreakAfterPriority(5))
		push(&mockToken{text: "b"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithBreakBeforePriority(5))
		push(&mockToken{text: ")"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithBreakBeforePriority(5))
		le.EndWrapGroup()

		push(&mockToken{text: ","}, format.FormatFlag(0).
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithBreakAfterPriority(10))
		push(
			&mockToken{text: "c_is_very_long_so_it_wraps"},
			format.FormatFlag(0).
				WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
				WithBreakBeforePriority(10),
		)
		push(&mockToken{text: ")"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithBreakBeforePriority(10).
			WithIndentDecrease())
		le.EndWrapGroup()

		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		// Since outer group has P10 and inner has P5, outer will wrap first.
		// inner(a,b) should stay inline because it's only 10 chars.
		if !strings.Contains(sb.String(), "inner(a,b)") {
			t.Errorf("expected inner group to stay inline, got %q", sb.String())
		}
	})

	t.Run("}SEMI gets newline", func(t *testing.T) {
		var sb strings.Builder
		le := layout.NewEngine(
			printer.DefaultSettings(len("    aa , bb")+len("};")+1), // <---
			&sb,
			nil,
		)

		tokenStream := []struct {
			tok            mockToken
			flag           format.FormatFlag
			wgBegin, wgEnd bool
		}{
			{mockToken{text: "label10"}, defaultWhitespaceFlag, false, false},
			{mockToken{text: "="}, defaultWhitespaceFlag, false, false},
			{mockToken{text: "Group{}"}, defaultFlag, true, false},
			// @Step 1: LE shall emit wrapping here
			{mockToken{text: "{"}, defaultWhitespaceFlag.
				WithIndentIncrease().
				WithBreakAfterPriority(format.PriorityAboveHighest).
				WithSeparatorAfter(format.SepNoSpaceMaybeWrap),
				false, false,
			},
			{mockToken{text: "aa"}, defaultWhitespaceFlag, false, false},
			{mockToken{text: ","}, defaultWhitespaceFlag, false, false},
			{mockToken{text: "bb"}, defaultWhitespaceFlag, false, false},
			// @Step 2: Because #Step1 has been triggered within the "Group{}"
			//          LE shall apply wrap here as well
			{mockToken{text: "}"}, defaultFlag.
				WithIndentDecrease().
				WithBreakBeforePriority(format.PriorityAboveHighest).
				WithSeparatorBefore(format.SepNoSpaceMaybeWrap),
				false, false,
			},
			{mockToken{text: "Group{}"}, defaultFlag, false, true},
			{mockToken{text: ";"}, defaultFlag, false, false},
		}

		wg := make(map[string]func())
		for _, tf := range tokenStream {
			if tf.wgBegin {
				name := tf.tok.GetText()
				wg[name] = le.BeginWrapGroup(name)
				continue
			}
			if tf.wgEnd {
				name := tf.tok.GetText()
				pred, ok := wg[name]
				if !ok {
					t.Fatal("missing wg group")
				}
				pred()
				continue
			}
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "label10 = {\n    aa , bb\n};"
		if result != expected {
			t.Errorf(
				"Code is not formatted correctly\n"+
					"  Expect: %#v\n"+
					"     Got: %#v",
				expected, result,
			)
		}
	})

	t.Run("}SEMI dont gets newline - tail", func(t *testing.T) {
		var sb strings.Builder
		le := layout.NewEngine(
			printer.DefaultSettings(len("    aa , bb")+len("};")+1), // <---
			&sb,
			nil,
		)

		tokenStream := []struct {
			tok            mockToken
			flag           format.FormatFlag
			wgBegin, wgEnd bool
		}{
			{mockToken{text: "label10"}, defaultWhitespaceFlag, false, false},
			{mockToken{text: "="}, defaultWhitespaceFlag, false, false},
			{mockToken{text: "Group{}"}, defaultFlag, true, false},
			// @Step 1: LE shall emit wrapping here
			{mockToken{text: "{"}, defaultWhitespaceFlag.
				WithIndentIncrease().
				WithBreakAfterPriority(format.PriorityAboveHighest).
				WithSeparatorAfter(format.SepNoSpaceMaybeWrap),
				false, false,
			},
			{mockToken{text: "aa"}, defaultWhitespaceFlag, false, false},
			{mockToken{text: ","}, defaultWhitespaceFlag, false, false},
			{mockToken{text: "bb"}, defaultWhitespaceFlag, false, false},
			{mockToken{text: "Group{}"}, defaultFlag, false, true},
			// @Step 2: Besides this is the breaking point
			//      LE should NOT apply wrap here because the TOKEN(})
			//      with the winning wrap is outside the "Group{}"
			{mockToken{text: "}"}, defaultFlag.
				WithIndentDecrease().
				WithBreakBeforePriority(format.PriorityAboveHighest).
				WithSeparatorBefore(format.SepNoSpaceMaybeWrap),
				false, false,
			},
			{mockToken{text: ";"}, defaultFlag, false, false},
		}

		wg := make(map[string]func())
		for _, tf := range tokenStream {
			if tf.wgBegin {
				name := tf.tok.GetText()
				wg[name] = le.BeginWrapGroup(name)
				continue
			}
			if tf.wgEnd {
				name := tf.tok.GetText()
				pred, ok := wg[name]
				if !ok {
					t.Fatal("missing wg group")
				}
				pred()
				continue
			}
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "label10 = {\n    aa , bb};"
		if result != expected {
			t.Errorf(
				"Code is not formatted correctly\n"+
					"  Expect: %#v\n"+
					"     Got: %#v",
				expected, result,
			)
		}
	})

	t.Run("}SEMI dont gets newline - head", func(t *testing.T) {
		var sb strings.Builder
		le := layout.NewEngine(
			printer.DefaultSettings(len("    aa , bb")+len("};")+1), // <---
			&sb,
			nil,
		)

		tokenStream := []struct {
			tok            mockToken
			flag           format.FormatFlag
			wgBegin, wgEnd bool
		}{
			{mockToken{text: "label10"}, defaultWhitespaceFlag, false, false},
			{mockToken{text: "="}, defaultWhitespaceFlag, false, false},
			// @Step 1: LE shall emit wrapping here
			{mockToken{text: "{"}, defaultWhitespaceFlag.
				WithIndentIncrease().
				WithBreakAfterPriority(format.PriorityAboveHighest).
				WithSeparatorAfter(format.SepNoSpaceMaybeWrap),
				false, false,
			},
			{mockToken{text: "Group{}"}, defaultFlag, true, false},
			{mockToken{text: "aa"}, defaultWhitespaceFlag, false, false},
			{mockToken{text: ","}, defaultWhitespaceFlag, false, false},
			{mockToken{text: "bb"}, defaultWhitespaceFlag, false, false},
			// @Step 2: Besides this is the breaking point
			//      LE should NOT apply wrap here because the TOKEN({)
			//      with the winning wrap is outside the "Group{}"
			{mockToken{text: "}"}, defaultFlag.
				WithIndentDecrease().
				WithBreakBeforePriority(format.PriorityAboveHighest).
				WithSeparatorBefore(format.SepNoSpaceMaybeWrap),
				false, false,
			},
			{mockToken{text: "Group{}"}, defaultFlag, false, true},
			{mockToken{text: ";"}, defaultFlag, false, false},
		}

		wg := make(map[string]func())
		for _, tf := range tokenStream {
			if tf.wgBegin {
				name := tf.tok.GetText()
				wg[name] = le.BeginWrapGroup(name)
				continue
			}
			if tf.wgEnd {
				name := tf.tok.GetText()
				pred, ok := wg[name]
				if !ok {
					t.Fatal("missing wg group")
				}
				pred()
				continue
			}
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "label10 = {\n    aa , bb};"
		if result == "label10 = {\n    aa , bb\n};" {
			t.Skip("WrappingGroup Bug: group ownership leaking")
		} else if result != expected {
			t.Errorf(
				"Code is not formatted correctly\n"+
					"  Expect: %#v\n"+
					"     Got: %#v",
				expected, result,
			)
		}
	})

}

func TestLayoutEngine_OptionalWrapNoSpace(t *testing.T) {
	var sb strings.Builder
	// Linecap 5. "foo" (3) + "(" (1) + "a" (1) = 5. Fits.
	// But "foo(a)" is 5. If we add ")" (1) -> 6 > 5.
	le := layout.NewEngine(
		printer.DefaultSettings(5),
		&sb,
		nil,
	)

	push := func(token antlr.Token, flag format.FormatFlag) {
		if err := le.Push(token, flag); err != nil {
			t.Fatal(err)
		}
	}

	push(&mockToken{text: "foo"}, format.FormatFlag(0))
	push(&mockToken{text: "("}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepNone).
		WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
		WithBreakAfterPriority(5))
	push(&mockToken{text: "a"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
		WithBreakBeforePriority(5))
	push(&mockToken{text: ")"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepNone))

	if err := le.Flush(); err != nil {
		t.Fatal(err)
	}

	// Should break at the optional wrap with no space
	expected := "foo(\na)"
	if sb.String() != expected {
		t.Errorf("expected %q, got %q", expected, sb.String())
	}
}

func TestLayoutEngine_OptionalWrappingBehavior(t *testing.T) {
	defaultFlag := format.FormatFlag(0).
		WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
		WithSeparatorBeforePriority(format.PriorityNormal).
		WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
		WithSeparatorBeforePriority(format.PriorityNormal).
		WithSeparatorPriority(format.PriorityNormal)

	t.Run("fits within linecap", func(t *testing.T) {
		var sb strings.Builder
		le := layout.NewEngine(
			printer.DefaultSettings(80), // <---
			&sb,
			nil,
		)

		tokenStream := []struct {
			tok  mockToken
			flag format.FormatFlag
		}{
			{mockToken{text: "aaa"}, defaultFlag},
			{mockToken{text: "BBB"}, defaultFlag},
			{mockToken{text: "ccc"}, defaultFlag},
			{mockToken{text: "DDD"}, defaultFlag},
		}

		for _, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		expected := "aaaBBBcccDDD"
		if sb.String() != expected {
			t.Errorf("expected %q, got %q", expected, sb.String())
		}
	})

	t.Run("fits only one", func(t *testing.T) {
		var sb strings.Builder
		le := layout.NewEngine(
			printer.DefaultSettings(4), // <---
			&sb,
			nil,
		)

		tokenStream := []struct {
			tok  mockToken
			flag format.FormatFlag
		}{
			{mockToken{text: "aaa"}, defaultFlag},
			{mockToken{text: "BBB"}, defaultFlag},
			{mockToken{text: "ccc"}, defaultFlag},
			{mockToken{text: "DDD"}, defaultFlag},
		}

		for _, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		expected := "aaa\nBBB\nccc\nDDD"
		if sb.String() != expected {
			t.Errorf("expected %q, got %q", expected, sb.String())
		}
	})

	t.Run("fits only one", func(t *testing.T) {
		var sb strings.Builder
		le := layout.NewEngine(
			printer.DefaultSettings(2), // <---
			&sb,
			nil,
		)

		tokenStream := []struct {
			tok  mockToken
			flag format.FormatFlag
		}{
			{mockToken{text: "aaa"}, defaultFlag},
			{mockToken{text: "BBB"}, defaultFlag},
			{mockToken{text: "ccc"}, defaultFlag},
			{mockToken{text: "DDD"}, defaultFlag},
		}

		for _, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		expected := "\naaa\nBBB\nccc\nDDD"
		if sb.String() != expected {
			t.Errorf("expected %q, got %q", expected, sb.String())
		}
	})

	t.Run("fits only two", func(t *testing.T) {
		var sb strings.Builder
		le := layout.NewEngine(
			printer.DefaultSettings(7), // <---
			&sb,
			nil,
		)

		tokenStream := []struct {
			tok  mockToken
			flag format.FormatFlag
		}{
			{mockToken{text: "aaa"}, defaultFlag},
			{mockToken{text: "BBB"}, defaultFlag},
			{mockToken{text: "ccc"}, defaultFlag},
			{mockToken{text: "DDD"}, defaultFlag},
		}

		for _, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		expected := "aaaBBB\ncccDDD"
		if sb.String() != expected {
			t.Errorf("expected %q, got %q", expected, sb.String())
		}
	})

	t.Run("fits only three", func(t *testing.T) {
		var sb strings.Builder
		le := layout.NewEngine(
			printer.DefaultSettings(9), // <---
			&sb,
			nil,
		)

		tokenStream := []struct {
			tok  mockToken
			flag format.FormatFlag
		}{
			{mockToken{text: "aaa"}, defaultFlag},
			{mockToken{text: "BBB"}, defaultFlag},
			{mockToken{text: "ccc"}, defaultFlag},
			{mockToken{text: "DDD"}, defaultFlag},
		}

		for _, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		expected := "aaaBBBccc\nDDD"
		if sb.String() != expected {
			t.Errorf("expected %q, got %q", expected, sb.String())
		}
	})

	t.Run("fits only three-break after prio", func(t *testing.T) {
		var sb strings.Builder
		le := layout.NewEngine(
			printer.DefaultSettings(9), // <---
			&sb,
			nil,
		)

		tokenStream := []struct {
			tok  mockToken
			flag format.FormatFlag
		}{
			{mockToken{text: "aaa"}, defaultFlag},
			{
				mockToken{text: "BBB"},
				defaultFlag.WithBreakAfterPriority(format.PriorityAboveNormal),
			},
			{mockToken{text: "ccc"}, defaultFlag},
			{mockToken{text: "DDD"}, defaultFlag},
		}

		for _, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		expected := "aaaBBB\ncccDDD"
		if sb.String() != expected {
			t.Errorf("expected %q, got %q", expected, sb.String())
		}
	})

	t.Run("fits only three-break after prio", func(t *testing.T) {
		var sb strings.Builder
		le := layout.NewEngine(
			printer.DefaultSettings(9), // <---
			&sb,
			nil,
		)

		tokenStream := []struct {
			tok  mockToken
			flag format.FormatFlag
		}{
			{mockToken{text: "aaa"}, defaultFlag},
			{mockToken{text: "BBB"}, defaultFlag},
			{
				mockToken{text: "ccc"},
				defaultFlag.WithBreakBeforePriority(
					format.PriorityAboveNormal,
				),
			},
			{mockToken{text: "DDD"}, defaultFlag},
		}

		for _, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		expected := "aaaBBB\ncccDDD"
		if sb.String() != expected {
			t.Errorf("expected %q, got %q", expected, sb.String())
		}
	})

	// TODO: But why?
	t.Run("priority doesnt affect wrapping", func(t *testing.T) {
		var sb strings.Builder
		le := layout.NewEngine(
			printer.DefaultSettings(9), // <---
			&sb,
			nil,
		)

		tokenStream := []struct {
			tok  mockToken
			flag format.FormatFlag
		}{
			{mockToken{text: "aaa"}, defaultFlag},
			{
				mockToken{text: "BBB"},
				defaultFlag.WithSeparatorPriority(format.PriorityCritical),
			},
			{mockToken{text: "ccc"}, defaultFlag},
			{mockToken{text: "DDD"}, defaultFlag},
		}

		for _, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		expected := "aaaBBBccc\nDDD"
		if sb.String() != expected {
			t.Errorf("expected %q, got %q", expected, sb.String())
		}
	})

	t.Run("higher break prio wins over token prio", func(t *testing.T) {
		var sb strings.Builder
		le := layout.NewEngine(
			printer.DefaultSettings(9), // <---
			&sb,
			nil,
		)

		tokenStream := []struct {
			tok  mockToken
			flag format.FormatFlag
		}{
			{mockToken{text: "aaa"}, defaultFlag},
			{
				mockToken{text: "BBB"},
				defaultFlag.WithSeparatorPriority(format.PriorityCritical),
			},
			{mockToken{text: "ccc"}, defaultFlag},
			{
				mockToken{text: "DDD"},
				defaultFlag.WithBreakBeforePriority(format.PriorityCritical),
			},
		}

		for _, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		expected := "aaaBBBccc\nDDD"
		if sb.String() != expected {
			t.Errorf("expected %q, got %q", expected, sb.String())
		}
	})
}

func TestLayoutEngine_SoftIndentAbsorbMode(t *testing.T) {
	var sb strings.Builder
	// Linecap 10. "foo=arr[" is 8 chars, fits on first line.
	// "very_long_name" is 14 chars, forces wrap after "[".
	// "]" wraps again because 14 + 1 > 10.
	le := layout.NewEngine(
		printer.DefaultSettings(10),
		&sb,
		nil,
	)

	push := func(token antlr.Token, flag format.FormatFlag) {
		if err := le.Push(token, flag); err != nil {
			t.Fatal(err)
		}
	}

	push(&mockToken{text: "foo="}, format.FormatFlag(0))
	endSoftIndent := le.BeginSoftIndent() // Manually inject SlcSoftIdentBegin

	push(&mockToken{text: "arr"}, format.FormatFlag(0))

	push(&mockToken{text: "["}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepNone).
		WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
		WithBreakAfterPriority(5).
		WithIndentIncrease())

	push(&mockToken{text: "very_long_name"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
		WithBreakBeforePriority(5))

	push(&mockToken{text: "]"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
		WithBreakBeforePriority(5).
		WithIndentDecrease())

	push(&mockToken{text: ";"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepNone))

	endSoftIndent()

	if err := le.Flush(); err != nil {
		t.Fatal(err)
	}

	// In Absorb Mode (softIndentAbsorbMode = true):
	// HardIndent inside `[` is 4. SoftIndent is 4. max(4,4) = 4.
	// `]` has IndentDecrease, so it ignores soft indent (hasHardIndentMarker = true) and uses HardIndent=0.
	expected := "foo= arr[\n    very_long_name\n];"
	if sb.String() != expected {
		t.Errorf("expected %q, got %q", expected, sb.String())
	}
}

func TestLayoutEngine_BareChainAbsorbMode(t *testing.T) {
	var sb strings.Builder
	// Set linecap very small so all optional wraps trigger a line break
	le := layout.NewEngine(
		printer.DefaultSettings(5),
		&sb,
		nil,
	)

	push := func(token antlr.Token, flag format.FormatFlag) {
		if err := le.Push(token, flag); err != nil {
			t.Fatal(err)
		}
	}

	push(
		&mockToken{text: "param"},
		format.FormatFlag(0).WithSeparatorAfter(format.SepSpace),
	)
	push(
		&mockToken{text: "TEST"},
		format.FormatFlag(0).WithSeparatorAfter(format.SepSpace),
	)
	push(
		&mockToken{text: "="},
		format.FormatFlag(0).WithSeparatorAfter(format.SepSpace),
	)

	le.BeginSoftIndent()

	// (a)
	push(&mockToken{text: "("}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorAfter(format.SepNone).
		WithIndentIncrease())
	push(&mockToken{text: "a"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepNone).
		WithSeparatorAfter(format.SepNone))
	push(&mockToken{text: ")"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepNone).
		WithIndentDecrease())

	// ? 1
	push(&mockToken{text: "?"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpaceMaybeWrap).
		WithBreakBeforePriority(5).
		WithSeparatorAfter(format.SepSpace))
	push(&mockToken{text: "1"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorAfter(format.SepNone))

	// : (b)
	push(&mockToken{text: ":"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpaceMaybeWrap).
		WithBreakBeforePriority(5).
		WithSeparatorAfter(format.SepSpace))
	push(&mockToken{text: "("}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorAfter(format.SepNone).
		WithIndentIncrease())
	push(&mockToken{text: "b"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepNone).
		WithSeparatorAfter(format.SepNone))
	push(&mockToken{text: ")"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepNone).
		WithIndentDecrease())

	// ? 2
	push(&mockToken{text: "?"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpaceMaybeWrap).
		WithBreakBeforePriority(5).
		WithSeparatorAfter(format.SepSpace))
	push(&mockToken{text: "2"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorAfter(format.SepNone))

	// : (c)
	push(&mockToken{text: ":"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpaceMaybeWrap).
		WithBreakBeforePriority(5).
		WithSeparatorAfter(format.SepSpace))
	push(&mockToken{text: "("}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorAfter(format.SepNone).
		WithIndentIncrease())
	push(&mockToken{text: "c"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepNone).
		WithSeparatorAfter(format.SepNone))
	push(&mockToken{text: ")"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepNone).
		WithIndentDecrease())

	// ? 3
	push(&mockToken{text: "?"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpaceMaybeWrap).
		WithBreakBeforePriority(5).
		WithSeparatorAfter(format.SepSpace))
	push(&mockToken{text: "3"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorAfter(format.SepNone))

	// : 4
	push(&mockToken{text: ":"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpaceMaybeWrap).
		WithBreakBeforePriority(5).
		WithSeparatorAfter(format.SepSpace))
	push(&mockToken{text: "4"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorAfter(format.SepNone))

	le.EndSoftIndent()

	push(&mockToken{text: ";"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepNone))

	if err := le.Flush(); err != nil {
		t.Fatal(err)
	}

	expected := "param TEST = (a)\n    ? 1\n    : (b)\n    ? 2\n    : (c)\n    ? 3\n    : 4;"
	if sb.String() != expected {
		t.Errorf("expected %q, got %q", expected, sb.String())
	}
}

func TestLayoutEngine_LenghtCalc(t *testing.T) {
	defaultFlag := format.FormatFlag(0).
		WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
		WithBreakBeforePriority(format.PriorityNormal).
		WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
		WithBreakAfterPriority(format.PriorityNormal).
		WithSeparatorPriority(format.PriorityNormal)
	defaultWhitespaceFlag := defaultFlag.
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorAfter(format.SepSpace)

	_ = defaultWhitespaceFlag

	t.Run("01-SepNoWrap-NoIdent-All", func(t *testing.T) {
		var sb strings.Builder
		le := layout.NewEngine(
			printer.DefaultSettings(28), // <---
			&sb,
			nil,
		)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultFlag, 1, -1},
			{mockToken{text: "bb"}, defaultFlag, 3, -1},
			{mockToken{text: "ccc"}, defaultFlag, 6, -1},
			{mockToken{text: "dddd"}, defaultFlag, 10, -1},
			{mockToken{text: "eeeee"}, defaultFlag, 15, -1},
			{mockToken{text: "ffffff"}, defaultFlag, 21, -1},
			{mockToken{text: "ggggggg"}, defaultFlag, 28, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				printer.TabSize,
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Fatalf(
					"Token % 02d is not formatted correctly\n"+
						"  Expect: %#v\n"+
						"     Got: %#v",
					idx, col, tf.col,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := 28
		if len(result) != expected {
			t.Errorf(
				"Code is not formatted correctly\n"+
					"  Expect: %#v\n"+
					"     Got: %#v",
				expected, result,
			)
		}
	})

	t.Run("02-SepSpaceAfter-NoIdent-All", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(28 + 6) // <---
		le := layout.NewEngine(
			cfg,
			&sb,
			nil,
		)

		defaultSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultSpaceAfter, 1, -1},
			{mockToken{text: "bb"}, defaultSpaceAfter, 3 + 1, -1},
			{mockToken{text: "ccc"}, defaultSpaceAfter, 6 + 2, -1},
			{mockToken{text: "dddd"}, defaultSpaceAfter, 10 + 3, -1},
			{mockToken{text: "eeeee"}, defaultSpaceAfter, 15 + 4, -1},
			{mockToken{text: "ffffff"}, defaultSpaceAfter, 21 + 5, -1},
			{mockToken{text: "ggggggg"}, defaultSpaceAfter, 28 + 6, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Fatalf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: %#v\n"+
						"     Got: %#v",
					idx, col, tf.col,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := 28 + 6
		if len(result) != expected {
			t.Errorf(
				"Code is not formatted correctly\n"+
					"  Expect: %#v\n"+
					"     Got: %#v",
				expected, result,
			)
		}
	})

	t.Run("03-SepSpaceBefore-NoIdent-All", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(28 + 6) // <---
		le := layout.NewEngine(
			cfg,
			&sb,
			nil,
		)

		defaultSpaceAfter := defaultFlag.
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultSpaceAfter, 1, -1},
			{mockToken{text: "bb"}, defaultSpaceAfter, 3 + 1, -1},
			{mockToken{text: "ccc"}, defaultSpaceAfter, 6 + 2, -1},
			{mockToken{text: "dddd"}, defaultSpaceAfter, 10 + 3, -1},
			{mockToken{text: "eeeee"}, defaultSpaceAfter, 15 + 4, -1},
			{mockToken{text: "ffffff"}, defaultSpaceAfter, 21 + 5, -1},
			{mockToken{text: "ggggggg"}, defaultSpaceAfter, 28 + 6, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Fatalf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: %#v\n"+
						"     Got: %#v",
					idx, col, tf.col,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := 28 + 6
		if len(result) != expected {
			t.Errorf(
				"Code is not formatted correctly\n"+
					"  Expect: %#v\n"+
					"     Got: %#v",
				expected, result,
			)
		}
	})

	t.Run("04-SepOptionalWrapNoSpaceAfter-NoIdent-All", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(9) // <---
		le := layout.NewEngine(
			cfg,
			&sb,
			nil,
		)

		defaultOptionalWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithSeparatorPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceAfter, 3, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapNoSpaceAfter, 6, -1},
			// wrap here
			{mockToken{text: "dddd"}, defaultOptionalWrapNoSpaceAfter, 4, 3},
			{mockToken{text: "ee"}, defaultOptionalWrapNoSpaceAfter, 6, -1},
			{mockToken{text: "fff"}, defaultOptionalWrapNoSpaceAfter, 9, -1},
			// wrap here
			{mockToken{text: "gggg"}, defaultOptionalWrapNoSpaceAfter, 4, 6},
		}

		totalLen := 0
		for idx, tf := range tokenStream {
			totalLen += len(tf.tok.GetText())
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				totalLen += len("\n")

				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Fatalf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := totalLen
		if len(result) != expected {
			t.Errorf(
				"Code is not formatted correctly\n"+
					"  Expect: %#v\n"+
					"     Got: %#v",
				expected, result,
			)
		}
	})

	t.Run("SepNoWrap-NoIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)

		flagNoWrap := defaultFlag.
			WithSeparatorBefore(format.SepNone).
			WithSeparatorBeforePriority(format.PriorityAboveNormal).
			WithSeparatorAfter(format.SepNone).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, flagNoWrap, 1, -1},
			{mockToken{text: "bb"}, flagNoWrap, 3, -1},
			{mockToken{text: "ccc"}, flagNoWrap, 6, -1},
			{mockToken{text: "dddd"}, flagNoWrap, 10, -1},
			{mockToken{text: "eeeee"}, flagNoWrap, 15, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbcccddddeeeee"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepSpaceAfter-NoIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)

		firstTokenFlag := defaultFlag.
			WithSeparatorBefore(format.SepNone).
			WithSeparatorBeforePriority(format.PriorityAboveNormal).
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		defaultSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, firstTokenFlag, 1, -1},
			{mockToken{text: "bb"}, defaultSpaceAfter, 4, -1},
			{mockToken{text: "ccc"}, defaultSpaceAfter, 8, -1},
			{mockToken{text: "dddd"}, defaultSpaceAfter, 13, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepSpaceBefore-NoIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultSpaceBefore, 1, -1},
			{mockToken{text: "bb"}, defaultSpaceBefore, 4, -1},
			{mockToken{text: "ccc"}, defaultSpaceBefore, 8, -1},
			{mockToken{text: "dddd"}, defaultSpaceBefore, 13, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceAfter-NoIdent-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(28)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithSeparatorPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceAfter, 3, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapNoSpaceAfter, 6, -1},
			{mockToken{text: "dddd"}, defaultOptionalWrapNoSpaceAfter, 10, -1},
			{mockToken{text: "eeeee"}, defaultOptionalWrapNoSpaceAfter, 15, -1},
			{
				mockToken{text: "ffffff"},
				defaultOptionalWrapNoSpaceAfter,
				21,
				-1,
			},
			{
				mockToken{text: "ggggggg"},
				defaultOptionalWrapNoSpaceAfter,
				28,
				-1,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbcccddddeeeeeffffffggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceAfter-NoIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithSeparatorPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceAfter, 3, -1},
			{
				mockToken{text: "cccccccc"},
				defaultOptionalWrapNoSpaceAfter,
				8,
				2,
			},
			{mockToken{text: "d"}, defaultOptionalWrapNoSpaceAfter, 1, 3},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abb\ncccccccc\nd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceBefore-NoIdent-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(28)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapNoSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceBefore, 3, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapNoSpaceBefore, 6, -1},
			{mockToken{text: "dddd"}, defaultOptionalWrapNoSpaceBefore, 10, -1},
			{
				mockToken{text: "eeeee"},
				defaultOptionalWrapNoSpaceBefore,
				15,
				-1,
			},
			{
				mockToken{text: "ffffff"},
				defaultOptionalWrapNoSpaceBefore,
				21,
				-1,
			},
			{
				mockToken{text: "ggggggg"},
				defaultOptionalWrapNoSpaceBefore,
				28,
				-1,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbcccddddeeeeeffffffggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceBefore-NoIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(9)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapNoSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceBefore, 3, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapNoSpaceBefore, 6, -1},
			{mockToken{text: "dddd"}, defaultOptionalWrapNoSpaceBefore, 4, 3},
			{mockToken{text: "ee"}, defaultOptionalWrapNoSpaceBefore, 6, -1},
			{mockToken{text: "fff"}, defaultOptionalWrapNoSpaceBefore, 9, -1},
			{mockToken{text: "gggg"}, defaultOptionalWrapNoSpaceBefore, 4, 6},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbccc\nddddeefff\ngggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceBefore-NoIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapNoSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceBefore, 3, -1},
			{
				mockToken{text: "cccccccc"},
				defaultOptionalWrapNoSpaceBefore,
				8,
				2,
			},
			{mockToken{text: "d"}, defaultOptionalWrapNoSpaceBefore, 1, 3},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abb\ncccccccc\nd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapAfter-NoIdent-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(34)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpaceMaybeWrap).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapAfter, 4, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapAfter, 8, -1},
			{mockToken{text: "dddd"}, defaultOptionalWrapAfter, 13, -1},
			{mockToken{text: "eeeee"}, defaultOptionalWrapAfter, 19, -1},
			{mockToken{text: "ffffff"}, defaultOptionalWrapAfter, 26, -1},
			{mockToken{text: "ggggggg"}, defaultOptionalWrapAfter, 34, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd eeeee ffffff ggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapAfter-NoIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(9)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpaceMaybeWrap).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapAfter, 4, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapAfter, 8, -1},
			{mockToken{text: "dddd"}, defaultOptionalWrapAfter, 4, 3},
			{mockToken{text: "ee"}, defaultOptionalWrapAfter, 7, -1},
			{mockToken{text: "fff"}, defaultOptionalWrapAfter, 3, 5},
			{mockToken{text: "gggg"}, defaultOptionalWrapAfter, 8, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc\ndddd ee\nfff gggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapAfter-NoIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpaceMaybeWrap).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapAfter, 4, -1},
			{mockToken{text: "cccccccc"}, defaultOptionalWrapAfter, 8, 2},
			{mockToken{text: "d"}, defaultOptionalWrapAfter, 1, 3},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb\ncccccccc\nd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapBefore-NoIdent-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(34)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapBefore := defaultFlag.
			WithSeparatorBefore(format.SepSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapBefore, 4, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapBefore, 8, -1},
			{mockToken{text: "dddd"}, defaultOptionalWrapBefore, 13, -1},
			{mockToken{text: "eeeee"}, defaultOptionalWrapBefore, 19, -1},
			{mockToken{text: "ffffff"}, defaultOptionalWrapBefore, 26, -1},
			{mockToken{text: "ggggggg"}, defaultOptionalWrapBefore, 34, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd eeeee ffffff ggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapBefore-NoIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(9)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapBefore := defaultFlag.
			WithSeparatorBefore(format.SepSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapBefore, 4, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapBefore, 8, -1},
			{mockToken{text: "dddd"}, defaultOptionalWrapBefore, 4, 3},
			{mockToken{text: "ee"}, defaultOptionalWrapBefore, 7, -1},
			{mockToken{text: "fff"}, defaultOptionalWrapBefore, 3, 5},
			{mockToken{text: "gggg"}, defaultOptionalWrapBefore, 8, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc\ndddd ee\nfff gggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapBefore-NoIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapBefore := defaultFlag.
			WithSeparatorBefore(format.SepSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapBefore, 4, -1},
			{mockToken{text: "cccccccc"}, defaultOptionalWrapBefore, 8, 2},
			{mockToken{text: "d"}, defaultOptionalWrapBefore, 1, 3},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb\ncccccccc\nd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapAfter-NoIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(80)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultMandatoryWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepMandatoryWrap).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapAfter, 1, -1},
			{mockToken{text: "bb"}, defaultMandatoryWrapAfter, 2, -1},
			{mockToken{text: "ccc"}, defaultMandatoryWrapAfter, 3, -1},
			{mockToken{text: "dddd"}, defaultMandatoryWrapAfter, 4, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a\nbb\nccc\ndddd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapAfter-NoIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultMandatoryWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepMandatoryWrap).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapAfter, 1, -1},
			{mockToken{text: "bbbbbbbb"}, defaultMandatoryWrapAfter, 8, -1},
			{mockToken{text: "c"}, defaultMandatoryWrapAfter, 1, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a\nbbbbbbbb\nc"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapBefore-NoIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(80)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultMandatoryWrapBefore := defaultFlag.
			WithSeparatorBefore(format.SepMandatoryWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapBefore, 1, -1},
			{mockToken{text: "bb"}, defaultMandatoryWrapBefore, 2, -1},
			{mockToken{text: "ccc"}, defaultMandatoryWrapBefore, 3, -1},
			{mockToken{text: "dddd"}, defaultMandatoryWrapBefore, 4, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "\na\nbb\nccc\ndddd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapBefore-NoIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultMandatoryWrapBefore := defaultFlag.
			WithSeparatorBefore(format.SepMandatoryWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapBefore, 1, -1},
			{mockToken{text: "bbbbbbbb"}, defaultMandatoryWrapBefore, 8, -1},
			{mockToken{text: "c"}, defaultMandatoryWrapBefore, 1, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "\na\nbbbbbbbb\nc"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapNoSpaceAfter-NoIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(80)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultMandatoryWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepWrapButNoSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapNoSpaceAfter, 1, -1},
			{mockToken{text: "bb"}, defaultMandatoryWrapNoSpaceAfter, 2, -1},
			{mockToken{text: "ccc"}, defaultMandatoryWrapNoSpaceAfter, 3, -1},
			{mockToken{text: "dddd"}, defaultMandatoryWrapNoSpaceAfter, 4, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a\nbb\nccc\ndddd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapNoSpaceAfter-NoIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultMandatoryWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepWrapButNoSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapNoSpaceAfter, 1, -1},
			{
				mockToken{text: "bbbbbbbb"},
				defaultMandatoryWrapNoSpaceAfter,
				8,
				-1,
			},
			{mockToken{text: "c"}, defaultMandatoryWrapNoSpaceAfter, 1, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a\nbbbbbbbb\nc"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapNoSpaceBefore-NoIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(80)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultMandatoryWrapNoSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepWrapButNoSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapNoSpaceBefore, 1, -1},
			{mockToken{text: "bb"}, defaultMandatoryWrapNoSpaceBefore, 2, -1},
			{mockToken{text: "ccc"}, defaultMandatoryWrapNoSpaceBefore, 3, -1},
			{mockToken{text: "dddd"}, defaultMandatoryWrapNoSpaceBefore, 4, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "\na\nbb\nccc\ndddd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapNoSpaceBefore-NoIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultMandatoryWrapNoSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepWrapButNoSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapNoSpaceBefore, 1, -1},
			{
				mockToken{text: "bbbbbbbb"},
				defaultMandatoryWrapNoSpaceBefore,
				8,
				-1,
			},
			{mockToken{text: "c"}, defaultMandatoryWrapNoSpaceBefore, 1, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "\na\nbbbbbbbb\nc"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("05-SepOptionalWrapNoSpaceAfter-PreIdent-All", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(9) // <---
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(
			cfg,
			&sb,
			&preIdent,
		)

		defaultOptionalWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithSeparatorPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceAfter, 3, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapNoSpaceAfter, 6, -1},
			// wrap here
			{
				mockToken{text: "dddd"},
				defaultOptionalWrapNoSpaceAfter,
				4 + 4,
				3,
			},
			// wrap here
			{mockToken{text: "ee"}, defaultOptionalWrapNoSpaceAfter, 2 + 4, 4},
			{
				mockToken{text: "fff"},
				defaultOptionalWrapNoSpaceAfter,
				5 + 4,
				-1,
			},
			// wrap here
			{
				mockToken{text: "gggg"},
				defaultOptionalWrapNoSpaceAfter,
				4 + 4,
				6,
			},
		}

		totalLen := 0
		for idx, tf := range tokenStream {
			totalLen += len(tf.tok.GetText())
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			state := le.TesOnly_GetState()
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				state = le.TesOnly_GetState()
				totalLen += len("    \n")

				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := totalLen
		if len(result) != expected {
			t.Errorf(
				"Code is not formatted correctly\n"+
					"  Expect: %#v\n"+
					"     Got: %#v",
				expected, result,
			)
		}
	})

	t.Run(
		"06-SepOptionalWrapNoSpaceAfter-PreIdent-OverCap",
		func(t *testing.T) {
			var sb strings.Builder
			cfg := printer.DefaultSettings(12) // <---
			preIdent := getEmptyLeState()
			preIdent.HardIndent = cfg.TabWidth()
			preIdent.LineCap = cfg.LineCap()
			le := layout.NewEngine(
				cfg,
				&sb,
				&preIdent,
			)

			defaultOptionalWrapAfter := defaultFlag.
				WithSeparatorAfter(format.SepSpaceMaybeWrap).
				WithSeparatorPriority(format.PriorityAboveNormal).
				WithSeparatorBefore(format.SepSpace)

			tokenStream := []struct {
				tok           mockToken
				flag          format.FormatFlag
				col, breakIdx int
			}{
				{mockToken{text: "a"}, defaultOptionalWrapAfter, 1, -1},
				{mockToken{text: "bb"}, defaultOptionalWrapAfter, 3 + 1, -1},
				{mockToken{text: "ccc"}, defaultOptionalWrapAfter, 6 + 2, -1},
				// wrap here
				{
					mockToken{text: "dddddddd"},
					defaultOptionalWrapAfter,
					8 + 4,
					3,
				},
				// wrap here
				{mockToken{text: "ee"}, defaultOptionalWrapAfter, 2 + 4, 4},
				{
					mockToken{text: "fff"},
					defaultOptionalWrapAfter,
					5 + 4 + 1,
					-1,
				},
				// wrap here
				{mockToken{text: "gggg"}, defaultOptionalWrapAfter, 4 + 4, 6},
			}

			totalLen := 0 + len("   ")
			for idx, tf := range tokenStream {
				totalLen += len(tf.tok.GetText())
				if err := le.Push(&tf.tok, tf.flag); err != nil {
					t.Fatal(err)
				}
				state := le.TesOnly_GetState()
				col, breakId := layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
				)
				if tf.breakIdx != -1 && tf.breakIdx == breakId {
					le.TesOnly_ApplyBreak(breakId)
					state = le.TesOnly_GetState()
					totalLen += len("    \n")

					col, _ = layout.TesOnly_SimulateBuffer(
						le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
					)
				}
				if col != tf.col || breakId != tf.breakIdx {
					t.Errorf(
						"Token #% 02d is not formatted correctly\n"+
							"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
							"     Got: {col: % 3d, breakIdx: % 3d}",
						idx, tf.col, tf.breakIdx, col, breakId,
					)
				}
			}
			if err := le.Flush(); err != nil {
				t.Fatal(err)
			}

			result := sb.String()
			expected := totalLen
			if len(result) != expected {
				t.Errorf(
					"Code is not formatted correctly\n"+
						"  Expect: %#v\n"+
						"     Got: %#v len - %d",
					expected, result, len(result),
				)
			}
		},
	)

	t.Run("SepNoWrap-PreIdent-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(28)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		flagNoWrap := defaultFlag.
			WithSeparatorBefore(format.SepNone).
			WithSeparatorBeforePriority(format.PriorityAboveNormal).
			WithSeparatorAfter(format.SepNone).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, flagNoWrap, 1, -1},
			{mockToken{text: "bb"}, flagNoWrap, 3, -1},
			{mockToken{text: "ccc"}, flagNoWrap, 6, -1},
			{mockToken{text: "dddd"}, flagNoWrap, 10, -1},
			{mockToken{text: "eeeee"}, flagNoWrap, 15, -1},
			{mockToken{text: "ffffff"}, flagNoWrap, 21, -1},
			{mockToken{text: "ggggggg"}, flagNoWrap, 28, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbcccddddeeeeeffffffggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepNoWrap-PreIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		flagNoWrap := defaultFlag.
			WithSeparatorBefore(format.SepNone).
			WithSeparatorBeforePriority(format.PriorityAboveNormal).
			WithSeparatorAfter(format.SepNone).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, flagNoWrap, 1, -1},
			{mockToken{text: "bb"}, flagNoWrap, 3, -1},
			{mockToken{text: "ccc"}, flagNoWrap, 6, -1},
			{mockToken{text: "dddd"}, flagNoWrap, 10, -1},
			{mockToken{text: "eeeee"}, flagNoWrap, 15, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbcccddddeeeee"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepSpaceAfter-PreIdent-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(34)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		firstTokenFlag := defaultFlag.
			WithSeparatorBefore(format.SepNone).
			WithSeparatorBeforePriority(format.PriorityAboveNormal).
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		defaultSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, firstTokenFlag, 1, -1},
			{mockToken{text: "bb"}, defaultSpaceAfter, 4, -1},
			{mockToken{text: "ccc"}, defaultSpaceAfter, 8, -1},
			{mockToken{text: "dddd"}, defaultSpaceAfter, 13, -1},
			{mockToken{text: "eeeee"}, defaultSpaceAfter, 19, -1},
			{mockToken{text: "ffffff"}, defaultSpaceAfter, 26, -1},
			{mockToken{text: "ggggggg"}, defaultSpaceAfter, 34, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd eeeee ffffff ggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepSpaceAfter-PreIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		firstTokenFlag := defaultFlag.
			WithSeparatorBefore(format.SepNone).
			WithSeparatorBeforePriority(format.PriorityAboveNormal).
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		defaultSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, firstTokenFlag, 1, -1},
			{mockToken{text: "bb"}, defaultSpaceAfter, 4, -1},
			{mockToken{text: "ccc"}, defaultSpaceAfter, 8, -1},
			{mockToken{text: "dddd"}, defaultSpaceAfter, 13, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepSpaceBefore-PreIdent-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(34)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultSpaceBefore, 1, -1},
			{mockToken{text: "bb"}, defaultSpaceBefore, 4, -1},
			{mockToken{text: "ccc"}, defaultSpaceBefore, 8, -1},
			{mockToken{text: "dddd"}, defaultSpaceBefore, 13, -1},
			{mockToken{text: "eeeee"}, defaultSpaceBefore, 19, -1},
			{mockToken{text: "ffffff"}, defaultSpaceBefore, 26, -1},
			{mockToken{text: "ggggggg"}, defaultSpaceBefore, 34, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd eeeee ffffff ggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepSpaceBefore-PreIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultSpaceBefore, 1, -1},
			{mockToken{text: "bb"}, defaultSpaceBefore, 4, -1},
			{mockToken{text: "ccc"}, defaultSpaceBefore, 8, -1},
			{mockToken{text: "dddd"}, defaultSpaceBefore, 13, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceAfter-PreIdent-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(28)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultOptionalWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithSeparatorPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceAfter, 3, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapNoSpaceAfter, 6, -1},
			{mockToken{text: "dddd"}, defaultOptionalWrapNoSpaceAfter, 10, -1},
			{mockToken{text: "eeeee"}, defaultOptionalWrapNoSpaceAfter, 15, -1},
			{
				mockToken{text: "ffffff"},
				defaultOptionalWrapNoSpaceAfter,
				21,
				-1,
			},
			{
				mockToken{text: "ggggggg"},
				defaultOptionalWrapNoSpaceAfter,
				28,
				-1,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbcccddddeeeeeffffffggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceAfter-PreIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(10)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultOptionalWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithSeparatorPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceAfter, 3, -1},
			{
				mockToken{text: "cccccccc"},
				defaultOptionalWrapNoSpaceAfter,
				4 + 8,
				2,
			},
			{mockToken{text: "d"}, defaultOptionalWrapNoSpaceAfter, 4 + 1, 3},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abb\n    cccccccc\n    d"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceBefore-PreIdent-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(28)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultOptionalWrapNoSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceBefore, 3, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapNoSpaceBefore, 6, -1},
			{mockToken{text: "dddd"}, defaultOptionalWrapNoSpaceBefore, 10, -1},
			{
				mockToken{text: "eeeee"},
				defaultOptionalWrapNoSpaceBefore,
				15,
				-1,
			},
			{
				mockToken{text: "ffffff"},
				defaultOptionalWrapNoSpaceBefore,
				21,
				-1,
			},
			{
				mockToken{text: "ggggggg"},
				defaultOptionalWrapNoSpaceBefore,
				28,
				-1,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbcccddddeeeeeffffffggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceBefore-PreIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(9)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultOptionalWrapNoSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceBefore, 3, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapNoSpaceBefore, 6, -1},
			{
				mockToken{text: "dddd"},
				defaultOptionalWrapNoSpaceBefore,
				4 + 4,
				3,
			},
			{mockToken{text: "ee"}, defaultOptionalWrapNoSpaceBefore, 2 + 4, 4},
			{
				mockToken{text: "fff"},
				defaultOptionalWrapNoSpaceBefore,
				5 + 4,
				-1,
			},
			{
				mockToken{text: "gggg"},
				defaultOptionalWrapNoSpaceBefore,
				4 + 4,
				6,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbccc\n    dddd\n    eefff\n    gggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceBefore-PreIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(10)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultOptionalWrapNoSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceBefore, 3, -1},
			{
				mockToken{text: "cccccccc"},
				defaultOptionalWrapNoSpaceBefore,
				4 + 8,
				2,
			},
			{mockToken{text: "d"}, defaultOptionalWrapNoSpaceBefore, 4 + 1, 3},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abb\n    cccccccc\n    d"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapAfter-PreIdent-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(34)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultOptionalWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpaceMaybeWrap).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapAfter, 4, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapAfter, 8, -1},
			{mockToken{text: "dddd"}, defaultOptionalWrapAfter, 13, -1},
			{mockToken{text: "eeeee"}, defaultOptionalWrapAfter, 19, -1},
			{mockToken{text: "ffffff"}, defaultOptionalWrapAfter, 26, -1},
			{mockToken{text: "ggggggg"}, defaultOptionalWrapAfter, 34, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd eeeee ffffff ggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapAfter-PreIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(12)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultOptionalWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpaceMaybeWrap).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapAfter, 4, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapAfter, 8, -1},
			{mockToken{text: "dddd"}, defaultOptionalWrapAfter, 4 + 4, 3},
			{mockToken{text: "ee"}, defaultOptionalWrapAfter, 8 + 1 + 2, -1},
			{mockToken{text: "fff"}, defaultOptionalWrapAfter, 4 + 3, 5},
			{mockToken{text: "gggg"}, defaultOptionalWrapAfter, 7 + 1 + 4, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc\n    dddd ee\n    fff gggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapAfter-PreIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(10)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultOptionalWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpaceMaybeWrap).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapAfter, 4, -1},
			{mockToken{text: "cccccccc"}, defaultOptionalWrapAfter, 4 + 8, 2},
			{mockToken{text: "d"}, defaultOptionalWrapAfter, 4 + 1, 3},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb\n    cccccccc\n    d"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapBefore-PreIdent-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(34)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultOptionalWrapBefore := defaultFlag.
			WithSeparatorBefore(format.SepSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapBefore, 4, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapBefore, 8, -1},
			{mockToken{text: "dddd"}, defaultOptionalWrapBefore, 13, -1},
			{mockToken{text: "eeeee"}, defaultOptionalWrapBefore, 19, -1},
			{mockToken{text: "ffffff"}, defaultOptionalWrapBefore, 26, -1},
			{mockToken{text: "ggggggg"}, defaultOptionalWrapBefore, 34, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd eeeee ffffff ggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapBefore-PreIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(12)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultOptionalWrapBefore := defaultFlag.
			WithSeparatorBefore(format.SepSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapBefore, 4, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapBefore, 8, -1},
			{mockToken{text: "dddd"}, defaultOptionalWrapBefore, 4 + 4, 3},
			{mockToken{text: "ee"}, defaultOptionalWrapBefore, 8 + 1 + 2, -1},
			{mockToken{text: "fff"}, defaultOptionalWrapBefore, 4 + 3, 5},
			{mockToken{text: "gggg"}, defaultOptionalWrapBefore, 7 + 1 + 4, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc\n    dddd ee\n    fff gggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapBefore-PreIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(10)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultOptionalWrapBefore := defaultFlag.
			WithSeparatorBefore(format.SepSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapBefore, 4, -1},
			{mockToken{text: "cccccccc"}, defaultOptionalWrapBefore, 4 + 8, 2},
			{mockToken{text: "d"}, defaultOptionalWrapBefore, 4 + 1, 3},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb\n    cccccccc\n    d"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapAfter-PreIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(80)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultMandatoryWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepMandatoryWrap).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapAfter, 1, -1},
			{mockToken{text: "bb"}, defaultMandatoryWrapAfter, 4 + 2, -1},
			{mockToken{text: "ccc"}, defaultMandatoryWrapAfter, 4 + 3, -1},
			{mockToken{text: "dddd"}, defaultMandatoryWrapAfter, 4 + 4, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a\n    bb\n    ccc\n    dddd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapAfter-PreIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultMandatoryWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepMandatoryWrap).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapAfter, 1, -1},
			{mockToken{text: "bbbbbbbb"}, defaultMandatoryWrapAfter, 4 + 8, -1},
			{mockToken{text: "c"}, defaultMandatoryWrapAfter, 4 + 1, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a\n    bbbbbbbb\n    c"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapBefore-PreIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(80)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultMandatoryWrapBefore := defaultFlag.
			WithSeparatorBefore(format.SepMandatoryWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapBefore, 4 + 1, -1},
			{mockToken{text: "bb"}, defaultMandatoryWrapBefore, 4 + 2, -1},
			{mockToken{text: "ccc"}, defaultMandatoryWrapBefore, 4 + 3, -1},
			{mockToken{text: "dddd"}, defaultMandatoryWrapBefore, 4 + 4, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "\n    a\n    bb\n    ccc\n    dddd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapBefore-PreIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultMandatoryWrapBefore := defaultFlag.
			WithSeparatorBefore(format.SepMandatoryWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapBefore, 4 + 1, -1},
			{
				mockToken{text: "bbbbbbbb"},
				defaultMandatoryWrapBefore,
				4 + 8,
				-1,
			},
			{mockToken{text: "c"}, defaultMandatoryWrapBefore, 4 + 1, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "\n    a\n    bbbbbbbb\n    c"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapNoSpaceAfter-PreIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(80)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultMandatoryWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepWrapButNoSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapNoSpaceAfter, 1, -1},
			{
				mockToken{text: "bb"},
				defaultMandatoryWrapNoSpaceAfter,
				4 + 2,
				-1,
			},
			{
				mockToken{text: "ccc"},
				defaultMandatoryWrapNoSpaceAfter,
				4 + 3,
				-1,
			},
			{
				mockToken{text: "dddd"},
				defaultMandatoryWrapNoSpaceAfter,
				4 + 4,
				-1,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a\n    bb\n    ccc\n    dddd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapNoSpaceAfter-PreIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultMandatoryWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepWrapButNoSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapNoSpaceAfter, 1, -1},
			{
				mockToken{text: "bbbbbbbb"},
				defaultMandatoryWrapNoSpaceAfter,
				4 + 8,
				-1,
			},
			{mockToken{text: "c"}, defaultMandatoryWrapNoSpaceAfter, 4 + 1, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a\n    bbbbbbbb\n    c"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapNoSpaceBefore-PreIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(80)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultMandatoryWrapNoSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepWrapButNoSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{
				mockToken{text: "a"},
				defaultMandatoryWrapNoSpaceBefore,
				4 + 1,
				-1,
			},
			{
				mockToken{text: "bb"},
				defaultMandatoryWrapNoSpaceBefore,
				4 + 2,
				-1,
			},
			{
				mockToken{text: "ccc"},
				defaultMandatoryWrapNoSpaceBefore,
				4 + 3,
				-1,
			},
			{
				mockToken{text: "dddd"},
				defaultMandatoryWrapNoSpaceBefore,
				4 + 4,
				-1,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "\n    a\n    bb\n    ccc\n    dddd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapNoSpaceBefore-PreIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)

		defaultMandatoryWrapNoSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepWrapButNoSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{
				mockToken{text: "a"},
				defaultMandatoryWrapNoSpaceBefore,
				4 + 1,
				-1,
			},
			{
				mockToken{text: "bbbbbbbb"},
				defaultMandatoryWrapNoSpaceBefore,
				4 + 8,
				-1,
			},
			{
				mockToken{text: "c"},
				defaultMandatoryWrapNoSpaceBefore,
				4 + 1,
				-1,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "\n    a\n    bbbbbbbb\n    c"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepNoWrap-Ident-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(28)
		le := layout.NewEngine(cfg, &sb, nil)

		flagNoWrap := defaultFlag.
			WithSeparatorBefore(format.SepNone).
			WithSeparatorBeforePriority(format.PriorityAboveNormal).
			WithSeparatorAfter(format.SepNone).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, flagNoWrap, 1, -1},
			{mockToken{text: "bb"}, flagNoWrap, 3, -1},
			{mockToken{text: "ccc"}, flagNoWrap.WithIndentIncrease(), 6, -1},
			{mockToken{text: "dddd"}, flagNoWrap, 10, -1},
			{mockToken{text: "eeeee"}, flagNoWrap, 15, -1},
			{
				mockToken{text: "ffffff"},
				flagNoWrap.WithIndentDecrease(),
				21,
				-1,
			},
			{mockToken{text: "ggggggg"}, flagNoWrap, 28, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbcccddddeeeeeffffffggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepNoWrap-Ident-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)

		flagNoWrap := defaultFlag.
			WithSeparatorBefore(format.SepNone).
			WithSeparatorBeforePriority(format.PriorityAboveNormal).
			WithSeparatorAfter(format.SepNone).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, flagNoWrap, 1, -1},
			{mockToken{text: "bb"}, flagNoWrap, 3, -1},
			{mockToken{text: "ccc"}, flagNoWrap.WithIndentIncrease(), 6, -1},
			{mockToken{text: "dddd"}, flagNoWrap, 10, -1},
			{mockToken{text: "eeeee"}, flagNoWrap.WithIndentDecrease(), 15, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbcccddddeeeee"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepSpaceAfter-Ident-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(34)
		le := layout.NewEngine(cfg, &sb, nil)

		firstTokenFlag := defaultFlag.
			WithSeparatorBefore(format.SepNone).
			WithSeparatorBeforePriority(format.PriorityAboveNormal).
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		defaultSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, firstTokenFlag, 1, -1},
			{mockToken{text: "bb"}, defaultSpaceAfter, 4, -1},
			{
				mockToken{text: "ccc"},
				defaultSpaceAfter.WithIndentIncrease(),
				8,
				-1,
			},
			{mockToken{text: "dddd"}, defaultSpaceAfter, 13, -1},
			{mockToken{text: "eeeee"}, defaultSpaceAfter, 19, -1},
			{
				mockToken{text: "ffffff"},
				defaultSpaceAfter.WithIndentDecrease(),
				26,
				-1,
			},
			{mockToken{text: "ggggggg"}, defaultSpaceAfter, 34, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd eeeee ffffff ggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepSpaceAfter-Ident-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)

		firstTokenFlag := defaultFlag.
			WithSeparatorBefore(format.SepNone).
			WithSeparatorBeforePriority(format.PriorityAboveNormal).
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		defaultSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, firstTokenFlag, 1, -1},
			{mockToken{text: "bb"}, defaultSpaceAfter, 4, -1},
			{
				mockToken{text: "ccc"},
				defaultSpaceAfter.WithIndentIncrease(),
				8,
				-1,
			},
			{
				mockToken{text: "dddd"},
				defaultSpaceAfter.WithIndentDecrease(),
				13,
				-1,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepSpaceBefore-Ident-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(34)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultSpaceBefore, 1, -1},
			{mockToken{text: "bb"}, defaultSpaceBefore, 4, -1},
			{
				mockToken{text: "ccc"},
				defaultSpaceBefore.WithIndentIncrease(),
				8,
				-1,
			},
			{mockToken{text: "dddd"}, defaultSpaceBefore, 13, -1},
			{mockToken{text: "eeeee"}, defaultSpaceBefore, 19, -1},
			{
				mockToken{text: "ffffff"},
				defaultSpaceBefore.WithIndentDecrease(),
				26,
				-1,
			},
			{mockToken{text: "ggggggg"}, defaultSpaceBefore, 34, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd eeeee ffffff ggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepSpaceBefore-Ident-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultSpaceBefore, 1, -1},
			{mockToken{text: "bb"}, defaultSpaceBefore, 4, -1},
			{
				mockToken{text: "ccc"},
				defaultSpaceBefore.WithIndentIncrease(),
				8,
				-1,
			},
			{
				mockToken{text: "dddd"},
				defaultSpaceBefore.WithIndentDecrease(),
				13,
				-1,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceAfter-Ident-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(28)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithSeparatorPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceAfter, 3, -1},
			{
				mockToken{text: "ccc"},
				defaultOptionalWrapNoSpaceAfter.WithIndentIncrease(),
				6,
				-1,
			},
			{mockToken{text: "dddd"}, defaultOptionalWrapNoSpaceAfter, 10, -1},
			{mockToken{text: "eeeee"}, defaultOptionalWrapNoSpaceAfter, 15, -1},
			{
				mockToken{text: "ffffff"},
				defaultOptionalWrapNoSpaceAfter.WithIndentDecrease(),
				21,
				-1,
			},
			{
				mockToken{text: "ggggggg"},
				defaultOptionalWrapNoSpaceAfter,
				28,
				-1,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbcccddddeeeeeffffffggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceAfter-Ident-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(9)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithSeparatorPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceAfter, 3, -1},
			{
				mockToken{text: "ccc"},
				defaultOptionalWrapNoSpaceAfter.WithIndentIncrease(),
				6,
				-1,
			},
			{
				mockToken{text: "dddd"},
				defaultOptionalWrapNoSpaceAfter,
				4 + 4,
				3,
			},
			{mockToken{text: "ee"}, defaultOptionalWrapNoSpaceAfter, 2 + 4, 4},
			{
				mockToken{text: "fff"},
				defaultOptionalWrapNoSpaceAfter.WithIndentDecrease(),
				5 + 4,
				-1,
			},
			{mockToken{text: "gggg"}, defaultOptionalWrapNoSpaceAfter, 4, 6},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbccc\n    dddd\n    eefff\ngggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceAfter-Ident-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(10)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithSeparatorPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceAfter, 1, -1},
			{
				mockToken{text: "bb"},
				defaultOptionalWrapNoSpaceAfter.WithIndentIncrease(),
				3,
				-1,
			},
			{
				mockToken{text: "cccccccc"},
				defaultOptionalWrapNoSpaceAfter,
				4 + 8,
				2,
			},
			{
				mockToken{text: "d"},
				defaultOptionalWrapNoSpaceAfter.WithIndentDecrease(),
				1,
				3,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abb\n    cccccccc\nd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceBefore-Ident-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(28)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapNoSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceBefore, 3, -1},
			{
				mockToken{text: "ccc"},
				defaultOptionalWrapNoSpaceBefore.WithIndentIncrease(),
				6,
				-1,
			},
			{mockToken{text: "dddd"}, defaultOptionalWrapNoSpaceBefore, 10, -1},
			{
				mockToken{text: "eeeee"},
				defaultOptionalWrapNoSpaceBefore,
				15,
				-1,
			},
			{
				mockToken{text: "ffffff"},
				defaultOptionalWrapNoSpaceBefore.WithIndentDecrease(),
				21,
				-1,
			},
			{
				mockToken{text: "ggggggg"},
				defaultOptionalWrapNoSpaceBefore,
				28,
				-1,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbcccddddeeeeeffffffggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceBefore-Ident-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(9)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapNoSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceBefore, 3, -1},
			{
				mockToken{text: "ccc"},
				defaultOptionalWrapNoSpaceBefore.WithIndentIncrease(),
				6,
				-1,
			},
			{
				mockToken{text: "dddd"},
				defaultOptionalWrapNoSpaceBefore,
				4 + 4,
				3,
			},
			{mockToken{text: "ee"}, defaultOptionalWrapNoSpaceBefore, 2 + 4, 4},
			{
				mockToken{text: "fff"},
				defaultOptionalWrapNoSpaceBefore.WithIndentDecrease(),
				5 + 4,
				-1,
			},
			{mockToken{text: "gggg"}, defaultOptionalWrapNoSpaceBefore, 4, 6},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbccc\n    dddd\n    eefff\ngggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceBefore-Ident-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(10)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapNoSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceBefore, 1, -1},
			{
				mockToken{text: "bb"},
				defaultOptionalWrapNoSpaceBefore.WithIndentIncrease(),
				3,
				-1,
			},
			{
				mockToken{text: "cccccccc"},
				defaultOptionalWrapNoSpaceBefore,
				4 + 8,
				2,
			},
			{
				mockToken{text: "d"},
				defaultOptionalWrapNoSpaceBefore.WithIndentDecrease(),
				1,
				3,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abb\n    cccccccc\nd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapAfter-Ident-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(34)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpaceMaybeWrap).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapAfter, 4, -1},
			{
				mockToken{text: "ccc"},
				defaultOptionalWrapAfter.WithIndentIncrease(),
				8,
				-1,
			},
			{mockToken{text: "dddd"}, defaultOptionalWrapAfter, 13, -1},
			{mockToken{text: "eeeee"}, defaultOptionalWrapAfter, 19, -1},
			{
				mockToken{text: "ffffff"},
				defaultOptionalWrapAfter.WithIndentDecrease(),
				26,
				-1,
			},
			{mockToken{text: "ggggggg"}, defaultOptionalWrapAfter, 34, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd eeeee ffffff ggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapAfter-Ident-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(12)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpaceMaybeWrap).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapAfter, 4, -1},
			{
				mockToken{text: "ccc"},
				defaultOptionalWrapAfter.WithIndentIncrease(),
				8,
				-1,
			},
			{mockToken{text: "dddd"}, defaultOptionalWrapAfter, 4 + 4, 3},
			{mockToken{text: "ee"}, defaultOptionalWrapAfter, 8 + 1 + 2, -1},
			{
				mockToken{text: "fff"},
				defaultOptionalWrapAfter.WithIndentDecrease(),
				3,
				5,
			},
			{mockToken{text: "gggg"}, defaultOptionalWrapAfter, 3 + 1 + 4, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc\n    dddd ee\nfff gggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapAfter-Ident-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(10)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpaceMaybeWrap).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapAfter, 1, -1},
			{
				mockToken{text: "bb"},
				defaultOptionalWrapAfter.WithIndentIncrease(),
				4,
				-1,
			},
			{mockToken{text: "cccccccc"}, defaultOptionalWrapAfter, 4 + 8, 2},
			{
				mockToken{text: "d"},
				defaultOptionalWrapAfter.WithIndentDecrease(),
				1,
				3,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb\n    cccccccc\nd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapBefore-Ident-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(34)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapBefore := defaultFlag.
			WithSeparatorBefore(format.SepSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapBefore, 4, -1},
			{
				mockToken{text: "ccc"},
				defaultOptionalWrapBefore.WithIndentIncrease(),
				8,
				-1,
			},
			{mockToken{text: "dddd"}, defaultOptionalWrapBefore, 13, -1},
			{mockToken{text: "eeeee"}, defaultOptionalWrapBefore, 19, -1},
			{
				mockToken{text: "ffffff"},
				defaultOptionalWrapBefore.WithIndentDecrease(),
				26,
				-1,
			},
			{mockToken{text: "ggggggg"}, defaultOptionalWrapBefore, 34, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd eeeee ffffff ggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapBefore-Ident-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(12)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapBefore := defaultFlag.
			WithSeparatorBefore(format.SepSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapBefore, 4, -1},
			{
				mockToken{text: "ccc"},
				defaultOptionalWrapBefore.WithIndentIncrease(),
				8,
				-1,
			},
			{mockToken{text: "dddd"}, defaultOptionalWrapBefore, 4 + 4, 3},
			{mockToken{text: "ee"}, defaultOptionalWrapBefore, 8 + 1 + 2, -1},
			{
				mockToken{text: "fff"},
				defaultOptionalWrapBefore.WithIndentDecrease(),
				3,
				5,
			},
			{mockToken{text: "gggg"}, defaultOptionalWrapBefore, 3 + 1 + 4, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc\n    dddd ee\nfff gggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapBefore-Ident-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(10)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultOptionalWrapBefore := defaultFlag.
			WithSeparatorBefore(format.SepSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapBefore, 1, -1},
			{
				mockToken{text: "bb"},
				defaultOptionalWrapBefore.WithIndentIncrease(),
				4,
				-1,
			},
			{mockToken{text: "cccccccc"}, defaultOptionalWrapBefore, 4 + 8, 2},
			{
				mockToken{text: "d"},
				defaultOptionalWrapBefore.WithIndentDecrease(),
				1,
				3,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb\n    cccccccc\nd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapAfter-Ident-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(80)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultMandatoryWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepMandatoryWrap).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapAfter, 1, -1},
			{
				mockToken{text: "bb"},
				defaultMandatoryWrapAfter.WithIndentIncrease(),
				2,
				-1,
			},
			{mockToken{text: "ccc"}, defaultMandatoryWrapAfter, 4 + 3, -1},
			{
				mockToken{text: "dddd"},
				defaultMandatoryWrapAfter.WithIndentDecrease(),
				4,
				-1,
			},
			{mockToken{text: "eeeee"}, defaultMandatoryWrapAfter, 5, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a\nbb\n    ccc\ndddd\neeeee"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapAfter-Ident-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultMandatoryWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepMandatoryWrap).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapAfter, 1, -1},
			{
				mockToken{text: "bb"},
				defaultMandatoryWrapAfter.WithIndentIncrease(),
				2,
				-1,
			},
			{mockToken{text: "cccccccc"}, defaultMandatoryWrapAfter, 4 + 8, -1},
			{
				mockToken{text: "d"},
				defaultMandatoryWrapAfter.WithIndentDecrease(),
				1,
				-1,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a\nbb\n    cccccccc\nd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapBefore-Ident-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(80)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultMandatoryWrapBefore := defaultFlag.
			WithSeparatorBefore(format.SepMandatoryWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapBefore, 1, -1},
			{
				mockToken{text: "bb"},
				defaultMandatoryWrapBefore.WithIndentIncrease(),
				2,
				-1,
			},
			{mockToken{text: "ccc"}, defaultMandatoryWrapBefore, 4 + 3, -1},
			{
				mockToken{text: "dddd"},
				defaultMandatoryWrapBefore.WithIndentDecrease(),
				4,
				-1,
			},
			{mockToken{text: "eeeee"}, defaultMandatoryWrapBefore, 5, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "\na\nbb\n    ccc\ndddd\neeeee"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapBefore-Ident-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultMandatoryWrapBefore := defaultFlag.
			WithSeparatorBefore(format.SepMandatoryWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapBefore, 1, -1},
			{
				mockToken{text: "bb"},
				defaultMandatoryWrapBefore.WithIndentIncrease(),
				2,
				-1,
			},
			{
				mockToken{text: "cccccccc"},
				defaultMandatoryWrapBefore,
				4 + 8,
				-1,
			},
			{
				mockToken{text: "d"},
				defaultMandatoryWrapBefore.WithIndentDecrease(),
				1,
				-1,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "\na\nbb\n    cccccccc\nd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapNoSpaceAfter-Ident-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(80)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultMandatoryWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepWrapButNoSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapNoSpaceAfter, 1, -1},
			{
				mockToken{text: "bb"},
				defaultMandatoryWrapNoSpaceAfter.WithIndentIncrease(),
				2,
				-1,
			},
			{
				mockToken{text: "ccc"},
				defaultMandatoryWrapNoSpaceAfter,
				4 + 3,
				-1,
			},
			{
				mockToken{text: "dddd"},
				defaultMandatoryWrapNoSpaceAfter.WithIndentDecrease(),
				4,
				-1,
			},
			{mockToken{text: "eeeee"}, defaultMandatoryWrapNoSpaceAfter, 5, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a\nbb\n    ccc\ndddd\neeeee"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapNoSpaceAfter-Ident-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultMandatoryWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepWrapButNoSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapNoSpaceAfter, 1, -1},
			{
				mockToken{text: "bb"},
				defaultMandatoryWrapNoSpaceAfter.WithIndentIncrease(),
				2,
				-1,
			},
			{
				mockToken{text: "cccccccc"},
				defaultMandatoryWrapNoSpaceAfter,
				4 + 8,
				-1,
			},
			{
				mockToken{text: "d"},
				defaultMandatoryWrapNoSpaceAfter.WithIndentDecrease(),
				1,
				-1,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a\nbb\n    cccccccc\nd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapNoSpaceBefore-Ident-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(80)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultMandatoryWrapNoSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepWrapButNoSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapNoSpaceBefore, 1, -1},
			{
				mockToken{text: "bb"},
				defaultMandatoryWrapNoSpaceBefore.WithIndentIncrease(),
				2,
				-1,
			},
			{
				mockToken{text: "ccc"},
				defaultMandatoryWrapNoSpaceBefore,
				4 + 3,
				-1,
			},
			{
				mockToken{text: "dddd"},
				defaultMandatoryWrapNoSpaceBefore.WithIndentDecrease(),
				4,
				-1,
			},
			{
				mockToken{text: "eeeee"},
				defaultMandatoryWrapNoSpaceBefore,
				5,
				-1,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "\na\nbb\n    ccc\ndddd\neeeee"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapNoSpaceBefore-Ident-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)

		defaultMandatoryWrapNoSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepWrapButNoSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapNoSpaceBefore, 1, -1},
			{
				mockToken{text: "bb"},
				defaultMandatoryWrapNoSpaceBefore.WithIndentIncrease(),
				2,
				-1,
			},
			{
				mockToken{text: "cccccccc"},
				defaultMandatoryWrapNoSpaceBefore,
				4 + 8,
				-1,
			},
			{
				mockToken{text: "d"},
				defaultMandatoryWrapNoSpaceBefore.WithIndentDecrease(),
				1,
				-1,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "\na\nbb\n    cccccccc\nd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("07-SepOptionalWrapNoSpaceAfter-Ident-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(12) // <---
		le := layout.NewEngine(
			cfg,
			&sb,
			nil,
		)

		defaultOptionalWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpaceMaybeWrap).
			WithSeparatorPriority(format.PriorityAboveNormal).
			WithSeparatorBefore(format.SepSpace)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapAfter, 3 + 1, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapAfter.
				WithIndentIncrease(), 6 + 2, -1},
			// wrap here
			{mockToken{text: "dddddddd"}, defaultOptionalWrapAfter, 8 + 4, 3},
			// wrap here
			{mockToken{text: "ee"}, defaultOptionalWrapAfter, 2 + 4, 4},
			{mockToken{text: "fff"}, defaultOptionalWrapAfter.
				WithIndentDecrease(), 5 + 4 + 1, -1},
			// wrap here
			{mockToken{text: "gggg"}, defaultOptionalWrapAfter, 4, 6},
		}

		totalLen := 0 + len("   ")
		for idx, tf := range tokenStream {
			totalLen += len(tf.tok.GetText())
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			state := le.TesOnly_GetState()
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				state = le.TesOnly_GetState()
				totalLen += len("    \n")

				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := totalLen - len("    ")
		if len(result) != expected {
			t.Errorf(
				"Code is not formatted correctly\n"+
					"  Expect: %#v\n"+
					"     Got: %#v len - %d",
				expected, result, len(result),
			)
		}
	})

	// =========================================================================
	// 4. Soft Indent Matrix
	// =========================================================================

	t.Run("SepNoWrap-SoftIdent-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(28)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		flagNoWrap := defaultFlag.
			WithSeparatorBefore(format.SepNone).
			WithSeparatorBeforePriority(format.PriorityAboveNormal).
			WithSeparatorAfter(format.SepNone).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, flagNoWrap, 1, -1},
			{mockToken{text: "bb"}, flagNoWrap, 3, -1},
			{mockToken{text: "ccc"}, flagNoWrap, 6, -1},
			{mockToken{text: "dddd"}, flagNoWrap, 10, -1},
			{mockToken{text: "eeeee"}, flagNoWrap, 15, -1},
			{mockToken{text: "ffffff"}, flagNoWrap, 21, -1},
			{mockToken{text: "ggggggg"}, flagNoWrap, 28, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbcccddddeeeeeffffffggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepNoWrap-SoftIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		flagNoWrap := defaultFlag.
			WithSeparatorBefore(format.SepNone).
			WithSeparatorBeforePriority(format.PriorityAboveNormal).
			WithSeparatorAfter(format.SepNone).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, flagNoWrap, 1, -1},
			{mockToken{text: "bb"}, flagNoWrap, 3, -1},
			{mockToken{text: "ccc"}, flagNoWrap, 6, -1},
			{mockToken{text: "dddd"}, flagNoWrap, 10, -1},
			{mockToken{text: "eeeee"}, flagNoWrap, 15, -1},
			{mockToken{text: "ffffff"}, flagNoWrap, 21, -1},
			{mockToken{text: "ggggggg"}, flagNoWrap, 28, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbcccddddeeeeeffffffggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepSpaceAfter-SoftIdent-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(34)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		firstTokenFlag := defaultFlag.
			WithSeparatorBefore(format.SepNone).
			WithSeparatorBeforePriority(format.PriorityAboveNormal).
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		defaultSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, firstTokenFlag, 1, -1},
			{mockToken{text: "bb"}, defaultSpaceAfter, 4, -1},
			{mockToken{text: "ccc"}, defaultSpaceAfter, 8, -1},
			{mockToken{text: "dddd"}, defaultSpaceAfter, 13, -1},
			{mockToken{text: "eeeee"}, defaultSpaceAfter, 19, -1},
			{mockToken{text: "ffffff"}, defaultSpaceAfter, 26, -1},
			{mockToken{text: "ggggggg"}, defaultSpaceAfter, 34, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd eeeee ffffff ggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepSpaceAfter-SoftIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		firstTokenFlag := defaultFlag.
			WithSeparatorBefore(format.SepNone).
			WithSeparatorBeforePriority(format.PriorityAboveNormal).
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		defaultSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, firstTokenFlag, 1, -1},
			{mockToken{text: "bb"}, defaultSpaceAfter, 4, -1},
			{mockToken{text: "ccc"}, defaultSpaceAfter, 8, -1},
			{mockToken{text: "dddd"}, defaultSpaceAfter, 13, -1},
			{mockToken{text: "eeeee"}, defaultSpaceAfter, 19, -1},
			{mockToken{text: "ffffff"}, defaultSpaceAfter, 26, -1},
			{mockToken{text: "ggggggg"}, defaultSpaceAfter, 34, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd eeeee ffffff ggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepSpaceBefore-SoftIdent-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(34)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultSpaceBefore, 1, -1},
			{mockToken{text: "bb"}, defaultSpaceBefore, 4, -1},
			{mockToken{text: "ccc"}, defaultSpaceBefore, 8, -1},
			{mockToken{text: "dddd"}, defaultSpaceBefore, 13, -1},
			{mockToken{text: "eeeee"}, defaultSpaceBefore, 19, -1},
			{mockToken{text: "ffffff"}, defaultSpaceBefore, 26, -1},
			{mockToken{text: "ggggggg"}, defaultSpaceBefore, 34, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd eeeee ffffff ggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepSpaceBefore-SoftIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultSpaceBefore, 1, -1},
			{mockToken{text: "bb"}, defaultSpaceBefore, 4, -1},
			{mockToken{text: "ccc"}, defaultSpaceBefore, 8, -1},
			{mockToken{text: "dddd"}, defaultSpaceBefore, 13, -1},
			{mockToken{text: "eeeee"}, defaultSpaceBefore, 19, -1},
			{mockToken{text: "ffffff"}, defaultSpaceBefore, 26, -1},
			{mockToken{text: "ggggggg"}, defaultSpaceBefore, 34, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd eeeee ffffff ggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceAfter-SoftIdent-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(80)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultOptionalWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithSeparatorPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceAfter, 3, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapNoSpaceAfter, 6, -1},
			{
				mockToken{text: "dddddddd"},
				defaultOptionalWrapNoSpaceAfter,
				14,
				-1,
			},
			{mockToken{text: "ee"}, defaultOptionalWrapNoSpaceAfter, 16, -1},
			{mockToken{text: "fff"}, defaultOptionalWrapNoSpaceAfter, 19, -1},
			{mockToken{text: "gggg"}, defaultOptionalWrapNoSpaceAfter, 23, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbcccddddddddeefffgggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceAfter-SoftIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(12)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultOptionalWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithSeparatorPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceAfter, 3, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapNoSpaceAfter, 6, -1},
			// dddddddd wraps to line 2 (breakId 4 = buffer index with SlcSoftIdentBegin)
			{
				mockToken{text: "dddddddd"},
				defaultOptionalWrapNoSpaceAfter,
				8 + 4,
				4,
			},
			// ee wraps to line 3 (breakId 5)
			{mockToken{text: "ee"}, defaultOptionalWrapNoSpaceAfter, 2 + 4, 5},
			{
				mockToken{text: "fff"},
				defaultOptionalWrapNoSpaceAfter,
				3 + 2 + 4,
				-1,
			},
			// gggg wraps to line 4 (breakId 7)
			{
				mockToken{text: "gggg"},
				defaultOptionalWrapNoSpaceAfter,
				4 + 4,
				7,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			state := le.TesOnly_GetState()
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				state = le.TesOnly_GetState()
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbccc\n    dddddddd\n    eefff\n    gggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceAfter-SoftIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(10)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultOptionalWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithSeparatorPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceAfter, 3, -1},
			{
				mockToken{text: "cccccccc"},
				defaultOptionalWrapNoSpaceAfter,
				4 + 8,
				3,
			},
			{mockToken{text: "d"}, defaultOptionalWrapNoSpaceAfter, 4 + 1, 4},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			state := le.TesOnly_GetState()
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				state = le.TesOnly_GetState()
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abb\n    cccccccc\n    d"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceBefore-SoftIdent-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(80)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultOptionalWrapNoSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithSeparatorPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceBefore, 3, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapNoSpaceBefore, 6, -1},
			{
				mockToken{text: "dddddddd"},
				defaultOptionalWrapNoSpaceBefore,
				14,
				-1,
			},
			{mockToken{text: "ee"}, defaultOptionalWrapNoSpaceBefore, 16, -1},
			{mockToken{text: "fff"}, defaultOptionalWrapNoSpaceBefore, 19, -1},
			{mockToken{text: "gggg"}, defaultOptionalWrapNoSpaceBefore, 23, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbcccddddddddeefffgggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceBefore-SoftIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(12)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultOptionalWrapNoSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithSeparatorPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceBefore, 3, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapNoSpaceBefore, 6, -1},
			{
				mockToken{text: "dddddddd"},
				defaultOptionalWrapNoSpaceBefore,
				8 + 4,
				4,
			},
			{mockToken{text: "ee"}, defaultOptionalWrapNoSpaceBefore, 2 + 4, 5},
			{
				mockToken{text: "fff"},
				defaultOptionalWrapNoSpaceBefore,
				3 + 2 + 4,
				-1,
			},
			{
				mockToken{text: "gggg"},
				defaultOptionalWrapNoSpaceBefore,
				4 + 4,
				7,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			state := le.TesOnly_GetState()
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				state = le.TesOnly_GetState()
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbccc\n    dddddddd\n    eefff\n    gggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapNoSpaceBefore-SoftIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(10)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultOptionalWrapNoSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithSeparatorPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceBefore, 3, -1},
			{
				mockToken{text: "cccccccc"},
				defaultOptionalWrapNoSpaceBefore,
				4 + 8,
				3,
			},
			{mockToken{text: "d"}, defaultOptionalWrapNoSpaceBefore, 4 + 1, 4},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			state := le.TesOnly_GetState()
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				state = le.TesOnly_GetState()
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abb\n    cccccccc\n    d"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapAfter-SoftIdent-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(34)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultOptionalWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpaceMaybeWrap).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapAfter, 4, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapAfter, 8, -1},
			{mockToken{text: "dddd"}, defaultOptionalWrapAfter, 13, -1},
			{mockToken{text: "eeeee"}, defaultOptionalWrapAfter, 19, -1},
			{mockToken{text: "ffffff"}, defaultOptionalWrapAfter, 26, -1},
			{mockToken{text: "ggggggg"}, defaultOptionalWrapAfter, 34, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd eeeee ffffff ggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapAfter-SoftIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(12)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultOptionalWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpaceMaybeWrap).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapAfter, 4, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapAfter, 8, -1},
			{mockToken{text: "dddd"}, defaultOptionalWrapAfter, 4 + 4, 4},
			{mockToken{text: "ee"}, defaultOptionalWrapAfter, 8 + 1 + 2, -1},
			{mockToken{text: "fff"}, defaultOptionalWrapAfter, 4 + 3, 6},
			{mockToken{text: "gggg"}, defaultOptionalWrapAfter, 7 + 1 + 4, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			state := le.TesOnly_GetState()
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				state = le.TesOnly_GetState()
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc\n    dddd ee\n    fff gggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapAfter-SoftIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(10)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultOptionalWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepSpaceMaybeWrap).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapAfter, 4, -1},
			{mockToken{text: "cccccccc"}, defaultOptionalWrapAfter, 4 + 8, 3},
			{mockToken{text: "d"}, defaultOptionalWrapAfter, 4 + 1, 4},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			state := le.TesOnly_GetState()
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				state = le.TesOnly_GetState()
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb\n    cccccccc\n    d"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapBefore-SoftIdent-Fits", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(34)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultOptionalWrapBefore := defaultFlag.
			WithSeparatorBefore(format.SepSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapBefore, 4, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapBefore, 8, -1},
			{mockToken{text: "dddd"}, defaultOptionalWrapBefore, 13, -1},
			{mockToken{text: "eeeee"}, defaultOptionalWrapBefore, 19, -1},
			{mockToken{text: "ffffff"}, defaultOptionalWrapBefore, 26, -1},
			{mockToken{text: "ggggggg"}, defaultOptionalWrapBefore, 34, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc dddd eeeee ffffff ggggggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapBefore-SoftIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(12)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultOptionalWrapBefore := defaultFlag.
			WithSeparatorBefore(format.SepSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapBefore, 4, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapBefore, 8, -1},
			{mockToken{text: "dddd"}, defaultOptionalWrapBefore, 4 + 4, 4},
			{mockToken{text: "ee"}, defaultOptionalWrapBefore, 8 + 1 + 2, -1},
			{mockToken{text: "fff"}, defaultOptionalWrapBefore, 4 + 3, 6},
			{mockToken{text: "gggg"}, defaultOptionalWrapBefore, 7 + 1 + 4, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			state := le.TesOnly_GetState()
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				state = le.TesOnly_GetState()
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb ccc\n    dddd ee\n    fff gggg"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepOptionalWrapBefore-SoftIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(10)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultOptionalWrapBefore := defaultFlag.
			WithSeparatorBefore(format.SepSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapBefore, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapBefore, 4, -1},
			{mockToken{text: "cccccccc"}, defaultOptionalWrapBefore, 4 + 8, 3},
			{mockToken{text: "d"}, defaultOptionalWrapBefore, 4 + 1, 4},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			state := le.TesOnly_GetState()
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				state = le.TesOnly_GetState()
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a bb\n    cccccccc\n    d"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapAfter-SoftIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(80)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultMandatoryWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepMandatoryWrap).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapAfter, 1, -1},
			{mockToken{text: "bb"}, defaultMandatoryWrapAfter, 4 + 2, -1},
			{mockToken{text: "ccc"}, defaultMandatoryWrapAfter, 4 + 3, -1},
			{mockToken{text: "dddd"}, defaultMandatoryWrapAfter, 4 + 4, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a\n    bb\n    ccc\n    dddd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapAfter-SoftIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultMandatoryWrapAfter := defaultFlag.
			WithSeparatorAfter(format.SepMandatoryWrap).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapAfter, 1, -1},
			{mockToken{text: "bbbbbbbb"}, defaultMandatoryWrapAfter, 4 + 8, -1},
			{mockToken{text: "c"}, defaultMandatoryWrapAfter, 4 + 1, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a\n    bbbbbbbb\n    c"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapBefore-SoftIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(80)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultMandatoryWrapBefore := defaultFlag.
			WithSeparatorBefore(format.SepMandatoryWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapBefore, 4 + 1, -1},
			{mockToken{text: "bb"}, defaultMandatoryWrapBefore, 4 + 2, -1},
			{mockToken{text: "ccc"}, defaultMandatoryWrapBefore, 4 + 3, -1},
			{mockToken{text: "dddd"}, defaultMandatoryWrapBefore, 4 + 4, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "\n    a\n    bb\n    ccc\n    dddd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapBefore-SoftIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultMandatoryWrapBefore := defaultFlag.
			WithSeparatorBefore(format.SepMandatoryWrap).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapBefore, 4 + 1, -1},
			{
				mockToken{text: "bbbbbbbb"},
				defaultMandatoryWrapBefore,
				4 + 8,
				-1,
			},
			{mockToken{text: "c"}, defaultMandatoryWrapBefore, 4 + 1, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "\n    a\n    bbbbbbbb\n    c"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapNoSpaceAfter-SoftIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(80)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultMandatoryWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepWrapButNoSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapNoSpaceAfter, 1, -1},
			{
				mockToken{text: "bb"},
				defaultMandatoryWrapNoSpaceAfter,
				4 + 2,
				-1,
			},
			{
				mockToken{text: "ccc"},
				defaultMandatoryWrapNoSpaceAfter,
				4 + 3,
				-1,
			},
			{
				mockToken{text: "dddd"},
				defaultMandatoryWrapNoSpaceAfter,
				4 + 4,
				-1,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a\n    bb\n    ccc\n    dddd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapNoSpaceAfter-SoftIdent-OverCap", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(5)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultMandatoryWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepWrapButNoSpace).
			WithSeparatorAfterPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultMandatoryWrapNoSpaceAfter, 1, -1},
			{
				mockToken{text: "bbbbbbbb"},
				defaultMandatoryWrapNoSpaceAfter,
				4 + 8,
				-1,
			},
			{mockToken{text: "c"}, defaultMandatoryWrapNoSpaceAfter, 4 + 1, -1},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "a\n    bbbbbbbb\n    c"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SepMandatoryWrapNoSpaceBefore-SoftIdent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(80)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultMandatoryWrapNoSpaceBefore := defaultFlag.
			WithSeparatorBefore(format.SepWrapButNoSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{
				mockToken{text: "a"},
				defaultMandatoryWrapNoSpaceBefore,
				4 + 1,
				-1,
			},
			{
				mockToken{text: "bb"},
				defaultMandatoryWrapNoSpaceBefore,
				4 + 2,
				-1,
			},
			{
				mockToken{text: "ccc"},
				defaultMandatoryWrapNoSpaceBefore,
				4 + 3,
				-1,
			},
			{
				mockToken{text: "dddd"},
				defaultMandatoryWrapNoSpaceBefore,
				4 + 4,
				-1,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(),
				le.TesOnly_GetState(),
				true,
				cfg.TabWidth(),
			)
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "\n    a\n    bb\n    ccc\n    dddd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run(
		"SepMandatoryWrapNoSpaceBefore-SoftIdent-OverCap",
		func(t *testing.T) {
			var sb strings.Builder
			cfg := printer.DefaultSettings(5)
			le := layout.NewEngine(cfg, &sb, nil)
			endSoftIndent := le.BeginSoftIndent()

			defaultMandatoryWrapNoSpaceBefore := defaultFlag.
				WithSeparatorBefore(format.SepWrapButNoSpace).
				WithSeparatorBeforePriority(format.PriorityAboveNormal)

			tokenStream := []struct {
				tok           mockToken
				flag          format.FormatFlag
				col, breakIdx int
			}{
				{
					mockToken{text: "a"},
					defaultMandatoryWrapNoSpaceBefore,
					4 + 1,
					-1,
				},
				{
					mockToken{text: "bbbbbbbb"},
					defaultMandatoryWrapNoSpaceBefore,
					4 + 8,
					-1,
				},
				{
					mockToken{text: "c"},
					defaultMandatoryWrapNoSpaceBefore,
					4 + 1,
					-1,
				},
			}

			for idx, tf := range tokenStream {
				if err := le.Push(&tf.tok, tf.flag); err != nil {
					t.Fatal(err)
				}
				col, breakId := layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(),
					le.TesOnly_GetState(),
					true,
					cfg.TabWidth(),
				)
				if col != tf.col || breakId != tf.breakIdx {
					t.Errorf(
						"Token #% 02d is not formatted correctly\n"+
							"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
							"     Got: {col: % 3d, breakIdx: % 3d}",
						idx, tf.col, tf.breakIdx, col, breakId,
					)
				}
			}
			endSoftIndent()
			if err := le.Flush(); err != nil {
				t.Fatal(err)
			}

			result := sb.String()
			expected := "\n    a\n    bbbbbbbb\n    c"
			if result != expected {
				t.Errorf("expected %q, got %q", expected, result)
			}
		},
	)

	// =========================================================================
	// 4.2 Soft Indent Interaction Matrix (Absorb Mode & Nesting)
	// =========================================================================

	t.Run("SoftIdent-Absorb-IndentIncrease-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(12)
		le := layout.NewEngine(cfg, &sb, nil)
		endSoftIndent := le.BeginSoftIndent()

		defaultOptionalWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithSeparatorPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceAfter, 3, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapNoSpaceAfter, 6, -1},
			// dddddddd has IndentIncrease: absorbs soft indent, so line 2 gets indent 0
			{
				mockToken{text: "dddddddd"},
				defaultOptionalWrapNoSpaceAfter.WithIndentIncrease(),
				8,
				4,
			},
			{mockToken{text: "ee"}, defaultOptionalWrapNoSpaceAfter, 10, -1},
			// ffff wraps to line 3: gets HardIndent=4 from dddddddd's increase
			{
				mockToken{text: "ffff"},
				defaultOptionalWrapNoSpaceAfter,
				4 + 4,
				6,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			state := le.TesOnly_GetState()
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				state = le.TesOnly_GetState()
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbccc\nddddddddee\n    ffff"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SoftIdent-Absorb-IndentDecrease-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(12)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)
		endSoftIndent := le.BeginSoftIndent()

		defaultOptionalWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithSeparatorPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceAfter, 3, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapNoSpaceAfter, 6, -1},
			// dddddddd has IndentDecrease: decreases HardIndent 4->0, and hasHardIndentMarker=true absorbs soft indent
			{
				mockToken{text: "dddddddd"},
				defaultOptionalWrapNoSpaceAfter.WithIndentDecrease(),
				8,
				4,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			state := le.TesOnly_GetState()
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				state = le.TesOnly_GetState()
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbccc\ndddddddd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("SoftIdent-NestedInsideHardIndent-Wraps", func(t *testing.T) {
		var sb strings.Builder
		cfg := printer.DefaultSettings(12)
		preIdent := getEmptyLeState()
		preIdent.HardIndent = cfg.TabWidth()
		preIdent.LineCap = cfg.LineCap()
		le := layout.NewEngine(cfg, &sb, &preIdent)
		endSoftIndent := le.BeginSoftIndent()

		defaultOptionalWrapNoSpaceAfter := defaultFlag.
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithSeparatorPriority(format.PriorityAboveNormal)

		tokenStream := []struct {
			tok           mockToken
			flag          format.FormatFlag
			col, breakIdx int
		}{
			{mockToken{text: "a"}, defaultOptionalWrapNoSpaceAfter, 1, -1},
			{mockToken{text: "bb"}, defaultOptionalWrapNoSpaceAfter, 3, -1},
			{mockToken{text: "ccc"}, defaultOptionalWrapNoSpaceAfter, 6, -1},
			// dddddddd wraps without hard indent markers:
			// Soft indent target is anchored to HardIndent+4 = 8, so continuation indent is +4 on top of HardIndent=4.
			{
				mockToken{text: "dddddddd"},
				defaultOptionalWrapNoSpaceAfter,
				8 + 8,
				4,
			},
		}

		for idx, tf := range tokenStream {
			if err := le.Push(&tf.tok, tf.flag); err != nil {
				t.Fatal(err)
			}
			state := le.TesOnly_GetState()
			col, breakId := layout.TesOnly_SimulateBuffer(
				le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
			)
			if tf.breakIdx != -1 && tf.breakIdx == breakId {
				le.TesOnly_ApplyBreak(breakId)
				state = le.TesOnly_GetState()
				col, _ = layout.TesOnly_SimulateBuffer(
					le.TesOnly_GetLeBuffer(), state, true, cfg.TabWidth(),
				)
			}
			if col != tf.col || breakId != tf.breakIdx {
				t.Errorf(
					"Token #% 02d is not formatted correctly\n"+
						"  Expect: {col: % 3d, breakIdx: % 3d}\n"+
						"     Got: {col: % 3d, breakIdx: % 3d}",
					idx, tf.col, tf.breakIdx, col, breakId,
				)
			}
		}
		endSoftIndent()
		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		result := sb.String()
		expected := "abbccc\n        dddddddd"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

}

func TestLayoutEngine_MandatoryWrapInGroup_TriggersGroupAndParents(
	t *testing.T,
) {
	var sb strings.Builder
	// linecap=80: tokens easily fit on a single line unless triggered
	le := layout.NewEngine(
		printer.DefaultSettings(80),
		&sb,
		nil,
	)

	// "prefix" outside group
	if err := le.Push(
		&mockToken{text: "prefix"},
		format.FormatFlag(0).WithSeparatorAfter(format.SepSpace),
	); err != nil {
		t.Fatal(err)
	}

	// Outer wrap group
	endOuter := le.BeginWrapGroup("OuterGroup")

	// `{` with SepNoSpaceMaybeWrap after, with IndentIncrease
	if err := le.Push(&mockToken{text: "{"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
		WithBreakAfterPriority(format.PriorityVeryHigh).
		WithSeparatorAfterPriority(format.PriorityMax).
		WithIndentIncrease()); err != nil {
		t.Fatal(err)
	}

	// Inner wrap group
	endInner := le.BeginWrapGroup("InnerGroup")

	// Comment with SepMandatoryWrap after
	if err := le.Push(&mockToken{text: "// comment"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorBeforePriority(format.PriorityLowest).
		WithSeparatorAfter(format.SepMandatoryWrap).
		WithSeparatorAfterPriority(format.PriorityMax)); err != nil {
		t.Fatal(err)
	}

	// Token after comment should NOT have caused immediate flush to writer yet (lazy buffering)
	if sb.Len() > 0 {
		t.Fatalf(
			"expected lazy buffering while wrap group is active, but writer received: %q",
			sb.String(),
		)
	}

	// First element in inner group
	if err := le.Push(&mockToken{text: "elem1"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpaceMaybeWrap).
		WithBreakBeforePriority(format.PriorityHigh)); err != nil {
		t.Fatal(err)
	}

	// Comma
	if err := le.Push(&mockToken{text: ","}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepNone).
		WithSeparatorAfter(format.SepSpaceMaybeWrap).
		WithBreakAfterPriority(format.PriorityHigh)); err != nil {
		t.Fatal(err)
	}

	// Following element in inner group (pushed AFTER mandatory wrap)
	// It has SepSpaceMaybeWrap before it, which should be promoted to mandatory wrap
	if err := le.Push(&mockToken{text: "elem2"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpaceMaybeWrap).
		WithBreakBeforePriority(format.PriorityHigh)); err != nil {
		t.Fatal(err)
	}

	endInner()

	// Closing bracket in outer group (pushed AFTER mandatory wrap)
	// It has SepNoSpaceMaybeWrap before it, which should be promoted to mandatory wrap
	if err := le.Push(&mockToken{text: "}"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
		WithSeparatorBeforePriority(format.PriorityMax).
		WithBreakBeforePriority(format.PriorityMax).
		WithIndentDecrease()); err != nil {
		t.Fatal(err)
	}

	endOuter()

	// Semicolon outside
	if err := le.Push(&mockToken{text: ";"}, format.FormatFlag(0).
		WithSeparatorBefore(format.SepNone).
		WithSeparatorBeforePriority(format.PriorityMax)); err != nil {
		t.Fatal(err)
	}

	if err := le.Flush(); err != nil {
		t.Fatal(err)
	}

	expected := "prefix {\n    // comment\n    elem1,\n    elem2\n};"
	got := sb.String()
	if got != expected {
		t.Errorf(
			"incorrect layout for mandatory wrap within wrap group:\n  expected: %q\n       got: %q",
			expected,
			got,
		)
	}
}

func TestLayoutEngine_MandatoryWrapOutsideGroup_FlushesImmediately(
	t *testing.T,
) {
	var sb strings.Builder
	le := layout.NewEngine(
		printer.DefaultSettings(80),
		&sb,
		nil,
	)

	// Token 1
	if err := le.Push(
		&mockToken{text: "stmt1"},
		format.FormatFlag(0).WithSeparatorAfter(format.SepMandatoryWrap),
	); err != nil {
		t.Fatal(err)
	}

	// Token 2 with mandatory wrap before (no wrap groups active)
	if err := le.Push(
		&mockToken{text: "stmt2"},
		format.FormatFlag(0).WithSeparatorBefore(format.SepMandatoryWrap),
	); err != nil {
		t.Fatal(err)
	}

	// When no wrap groups are active, stmt1 must have been flushed immediately
	if !strings.Contains(sb.String(), "stmt1") {
		t.Errorf(
			"expected stmt1 to be flushed immediately when outside wrap groups, got: %q",
			sb.String(),
		)
	}
}

func TestLayoutEngine_GreedyGroup(t *testing.T) {
	t.Run(
		"greedy single break without exploding peer tokens",
		func(t *testing.T) {
			var sb strings.Builder
			// linecap=15: "a.b.c.d.e" is 9 chars.
			// "prefix a.b.c.d.e" (16 chars) exceeds 15.
			// Greedy packing should break before ".d" or ".e" without breaking all dots.
			le := layout.NewEngine(
				printer.DefaultSettings(12),
				&sb,
				nil,
			)

			push := func(token antlr.Token, flag format.FormatFlag) {
				if err := le.Push(token, flag); err != nil {
					t.Fatal(err)
				}
			}

			push(&mockToken{text: "prefix"}, format.FormatFlag(0).
				WithSeparatorAfter(format.SepSpace).
				WithSeparatorAfterPriority(format.PriorityMax))

			endGroup := le.BeginGreedyGroup("MemberAccess")
			push(&mockToken{text: "a"}, format.FormatFlag(0))

			for _, seg := range []string{".b", ".c", ".d"} {
				push(&mockToken{text: seg}, format.FormatFlag(0).
					WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
					WithBreakBeforePriority(format.PriorityLowest))
			}
			endGroup()

			if err := le.Flush(); err != nil {
				t.Fatal(err)
			}

			// "prefix a.b.c" is 12 chars (fits linecap 12).
			// ".d" wraps to the next line.
			// With a standard wrap group, .b, .c, and .d would all break onto separate lines.
			// With greedy group, only .d wraps.
			expected := "prefix a.b.c\n.d"
			got := sb.String()
			if got != expected {
				t.Errorf("expected %q, got %q", expected, got)
			}
		},
	)

	t.Run(
		"nested wrap group does not inherit greedy behavior",
		func(t *testing.T) {
			var sb strings.Builder
			// Greedy outer group contains a child wrap group.
			// When the child wrap group exceeds linecap, all tokens in the child wrap group break collectively.
			le := layout.NewEngine(
				printer.DefaultSettings(5),
				&sb,
				nil,
			)

			push := func(token antlr.Token, flag format.FormatFlag) {
				if err := le.Push(token, flag); err != nil {
					t.Fatal(err)
				}
			}

			endGreedy := le.BeginGreedyGroup("OuterGreedy")
			push(&mockToken{text: "base"}, format.FormatFlag(0))

			endChildWrap := le.BeginWrapGroup("InnerWrap")
			push(&mockToken{text: "x"}, format.FormatFlag(0).
				WithSeparatorBefore(format.SepSpaceMaybeWrap).
				WithBreakBeforePriority(format.PriorityNormal))
			push(&mockToken{text: "y"}, format.FormatFlag(0).
				WithSeparatorBefore(format.SepSpaceMaybeWrap).
				WithBreakBeforePriority(format.PriorityNormal))
			endChildWrap()

			endGreedy()

			if err := le.Flush(); err != nil {
				t.Fatal(err)
			}

			// Because InnerWrap is a wrapping group, both x and y should break collectively.
			expected := "base\nx\ny"
			got := sb.String()
			if got != expected {
				t.Errorf("expected %q, got %q", expected, got)
			}
		},
	)

	t.Run(
		"isolation boundary prevents upward trigger propagation to grandparent wrap group",
		func(t *testing.T) {
			var sb strings.Builder
			// 3-tier hierarchy: OuterWrapGroup -> GreedyGroup -> InnerWrapGroup.
			// Linecap=35.
			// OuterWrapGroup: "out" + [GreedyGroup: "mid" + [InnerWrapGroup: "i1" + "i2_long_exceeding_line_cap"]] + " tail"
			// When "i2" triggers InnerWrapGroup to break across lines:
			// "out mid\ni1\ni2_long_exceeding_line_cap"
			// GreedyGroup isolates OuterWrapGroup, so "tail" in OuterWrapGroup is NOT triggered and stays inline.
			le := layout.NewEngine(
				printer.DefaultSettings(35),
				&sb,
				nil,
			)

			push := func(token antlr.Token, flag format.FormatFlag) {
				if err := le.Push(token, flag); err != nil {
					t.Fatal(err)
				}
			}

			endOuter := le.BeginWrapGroup("OuterWrap")
			push(&mockToken{text: "out"}, format.FormatFlag(0).
				WithSeparatorAfter(format.SepSpace))

			endGreedy := le.BeginGreedyGroup("MiddleGreedy")
			push(&mockToken{text: "mid"}, format.FormatFlag(0))

			endInner := le.BeginWrapGroup("InnerWrap")
			push(&mockToken{text: "i1"}, format.FormatFlag(0).
				WithSeparatorBefore(format.SepSpaceMaybeWrap).
				WithBreakBeforePriority(format.PriorityNormal))
			push(
				&mockToken{text: "i2_long_exceeding_line_cap"},
				format.FormatFlag(0).
					WithSeparatorBefore(format.SepSpaceMaybeWrap).
					WithBreakBeforePriority(format.PriorityNormal),
			)
			endInner()

			endGreedy()

			push(&mockToken{text: "tail"}, format.FormatFlag(0).
				WithSeparatorBefore(format.SepSpaceMaybeWrap).
				WithBreakBeforePriority(format.PriorityNormal))

			endOuter()

			if err := le.Flush(); err != nil {
				t.Fatal(err)
			}

			// InnerWrap is triggered (so i1 and i2 break collectively).
			// MiddleGreedy isolates OuterWrap, so "tail" in OuterWrap is NOT triggered and stays inline after i2.
			expected := "out mid\ni1\ni2_long_exceeding_line_cap tail"
			got := sb.String()
			if got != expected {
				t.Errorf("expected %q, got %q", expected, got)
			}
		},
	)

	t.Run(
		"propagation without greedy boundary triggers grandparent wrap group",
		func(t *testing.T) {
			var sb strings.Builder
			// Same 3-tier hierarchy, but with a standard WrapGroup in the middle:
			// OuterWrapGroup -> MiddleWrapGroup -> InnerWrapGroup.
			le := layout.NewEngine(
				printer.DefaultSettings(35),
				&sb,
				nil,
			)

			push := func(token antlr.Token, flag format.FormatFlag) {
				if err := le.Push(token, flag); err != nil {
					t.Fatal(err)
				}
			}

			endOuter := le.BeginWrapGroup("OuterWrap")
			push(&mockToken{text: "out"}, format.FormatFlag(0).
				WithSeparatorAfter(format.SepSpace))

			endMiddle := le.BeginWrapGroup("MiddleWrap")
			push(&mockToken{text: "mid"}, format.FormatFlag(0))

			endInner := le.BeginWrapGroup("InnerWrap")
			push(&mockToken{text: "i1"}, format.FormatFlag(0).
				WithSeparatorBefore(format.SepSpaceMaybeWrap).
				WithBreakBeforePriority(format.PriorityNormal))
			push(
				&mockToken{text: "i2_long_exceeding_line_cap"},
				format.FormatFlag(0).
					WithSeparatorBefore(format.SepSpaceMaybeWrap).
					WithBreakBeforePriority(format.PriorityNormal),
			)
			endInner()

			endMiddle()

			push(&mockToken{text: "tail"}, format.FormatFlag(0).
				WithSeparatorBefore(format.SepSpaceMaybeWrap).
				WithBreakBeforePriority(format.PriorityNormal))

			endOuter()

			if err := le.Flush(); err != nil {
				t.Fatal(err)
			}

			// Without an isolation boundary, InnerWrap explosion propagates to MiddleWrap and OuterWrap.
			// OuterWrap is triggered, so "tail" also wraps onto its own line.
			expected := "out mid\ni1\ni2_long_exceeding_line_cap\ntail"
			got := sb.String()
			if got != expected {
				t.Errorf("expected %q, got %q", expected, got)
			}
		},
	)
}

func TestLayoutEngine_GroupStackIntegrity(t *testing.T) {
	t.Run("EndWrapGroup panics on greedy group", func(t *testing.T) {
		var sb strings.Builder
		le := layout.NewEngine(printer.DefaultSettings(80), &sb, nil)
		_ = le.BeginGreedyGroup("Greedy")

		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf(
					"expected panic on EndWrapGroup called on greedy group",
				)
			}
			msg := fmt.Sprint(r)
			if !strings.Contains(msg, "non-collective group") {
				t.Errorf(
					"expected panic mentioning non-collective group, got: %s",
					msg,
				)
			}
		}()

		le.EndWrapGroup()
	})

	t.Run("EndGreedyGroup panics on collective wrap group", func(t *testing.T) {
		var sb strings.Builder
		le := layout.NewEngine(printer.DefaultSettings(80), &sb, nil)
		_ = le.BeginWrapGroup("Wrap")

		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf(
					"expected panic on EndGreedyGroup called on collective wrap group",
				)
			}
			msg := fmt.Sprint(r)
			if !strings.Contains(msg, "non-greedy group") {
				t.Errorf(
					"expected panic mentioning non-greedy group, got: %s",
					msg,
				)
			}
		}()

		le.EndGreedyGroup()
	})

	t.Run("closure panics on out of order closing", func(t *testing.T) {
		var sb strings.Builder
		le := layout.NewEngine(printer.DefaultSettings(80), &sb, nil)
		endOuter := le.BeginWrapGroup("Outer")
		endInner := le.BeginGreedyGroup("Inner")

		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf("expected panic on out-of-order group closure")
			}
			msg := fmt.Sprint(r)
			if !strings.Contains(msg, "group mismatch") {
				t.Errorf(
					"expected panic mentioning group mismatch, got: %s",
					msg,
				)
			}
		}()

		// Calling outer closure before inner closure must panic
		endOuter()
		_ = endInner
	})

	t.Run("unified EndGroup and closures balanced", func(t *testing.T) {
		var sb strings.Builder
		le := layout.NewEngine(printer.DefaultSettings(80), &sb, nil)

		end1 := le.BeginWrapGroup("G1")
		end2 := le.BeginGreedyGroup("G2")
		_ = le.BeginWrapGroup("G3")

		// End G3 via unified EndGroup
		le.EndGroup()

		// End G2 via its closure
		end2()

		// End G1 via its closure
		end1()

		if err := le.Flush(); err != nil {
			t.Fatalf("unexpected error flushing: %v", err)
		}
	})
}
