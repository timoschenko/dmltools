# Comment Formatting Guide

This guide explains how the formatter processes, positions, and preserves comments in DML code.

---

## 1. Overview

Comments in DML serve two primary roles:
1. **Documentation**: Providing context and explanations for device models, templates, and registers.
2. **Tool Directives & Annotations**: Driving external linters, static analysis tools (e.g. Coverity), and compiler checks.

The formatter treats comments as first-class lexical tokens.
Formatting decisions guarantee that comments remain readable
and—most critically—**never detached from the code they annotate**.

---

## 2. The "Sticky" Trailing Comment Guarantee

A core rule of the formatter is that **trailing end-of-line comments stick to the token they follow, regardless of line length**.

### Why Comments Never Wrap to a New Line

Many static analysis and linting tools associate annotations strictly with the exact line on which they appear:
- **Coverity Directives**: `/*% coverity[var_deref_op] %*/`, `// coverity[-alloc]`
- **Linter Suppressions**: `// type: ignore`, `/* lint: ... */`
- **Compiler Pragmas & Hints**: `/*% FALLTHROUGH %*/`

If the formatter were to break a long line before a trailing comment:
1. The comment would move to a new line, becoming a detached prefix to the *next* statement.
2. The linter or static analyzer would fail to match the suppression to the intended line, resulting in false-positive warnings or broken builds.

To prevent this, the formatter ensures that comments following statement terminators (`;`),
list separators (`,`), or braces (`{`, `}`) **remain on the exact same line**,
even if the combined length exceeds the configured line capacity.

### Examples

#### Statement Terminators (`;`)
Even when a statement plus comment exceeds the line limit, the comment remains on the same line:
```dml
// Stays on the same line even if exceeding line capacity:
register status size 4 @ BASE_ADDR + STATUS_OFFSET; // coverity[misra_c_2012_rule_10_1_violation]

// Parameter declaration with linter annotation:
param timeout_ms = 5000; // nolint: magic-number
```

#### List Items & Commas (`,`)
In lists, argument calls, and initializers, end-of-line comments stay glued to the comma of the corresponding element:
```dml
saved config_t device_cfg = {
    .clock_mhz = 100,    // Primary clock frequency
    .bus_width = 32,     // Data bus width in bits
    .burst_enabled = 1   // coverity[cert_str30_c_violation]
};
```

---

## 3. Comments Inside Wrap Groups

When an expression, argument list, or initializer expands across multiple lines (due to the `black`-style all-or-nothing wrap rule), comments attached to individual items remain attached to those items:

```dml
// All-or-nothing group explosion preserves item comments cleanly:
local (
    uint32 source_addr,      // Base source address
    uint32 destination_addr, // Target address
    bool enable_burst        // Burst transfer mode flag
) = parse_transfer_descriptor();
```

If an inline comment is placed right before or after an opening delimiter, the comment is kept with the opening construct while subsequent elements indent cleanly on subsequent lines:

```dml
constant lookup_table = [ // Initial lookup entries
    0x0001,
    0x0002,
    0x0004
];
```

---

## 4. Inline Block Comments (`/* ... */`)

Inline block comments placed within expressions or declarations adhere to natural spacing rules:
- **Spacing**: A single space separates inline comments from preceding and following identifiers or operators.
- **Punctuation Clinging**: If an inline comment immediately precedes a semicolon or comma, the punctuation clings directly without a space between the comment closing tag and the punctuation.

```dml
// Spacing around inline comments
param x /* type specifier */ : int = 1;

// Inline argument annotation
device_call(/* is_sync */ true, /* timeout */ 1000);

// Comment immediately preceding punctuation clings cleanly
method process(int /* flag */, bool /* mode */, /* mutable */ dbuffer_t *buf);
```

---

## 5. Standalone & Header Comments

Standalone comments that appear on their own lines (e.g. documentation blocks above methods, registers, or templates) automatically inherit the indentation level of the enclosed scope:

```dml
bank regs {
    // Status and control register block
    // Reset value: 0x00000000
    register ctrl size 4 @ 0x00 {
        /* Enable bit for peripheral operation */
        field en @ [0];
    }
}
```

---

## 6. Technical Implementation Details (Under the Hood)

*(For developers working on formatter internals)*

Comment positioning and priority resolution are implemented in
[`pkg/printer/defaults_tokens.go`](../../pkg/printer/defaults_tokens.go)
and [`pkg/printer/layout_engine.go`](../../pkg/printer/layout_engine.go):

- **Token Priority Override**:
  In [`defaults_tokens.go`](../../pkg/printer/defaults_tokens.go), `TokenComment` is configured with:
  ```go
  tokenFormatDefaults.data[dmlex.TokenComment] = tokenFormatDefaults.data[dmlex.TokenComment].
      WithSeparatorBefore(SepSpace).
      WithSeparatorBeforePriority(PriorityMax).
      WithSeparatorAfter(SepMandatoryWrap).
      WithSeparatorAfterPriority(PriorityMax)
  ```
  Setting `SeparatorBeforePriority(PriorityMax)` guarantees that line comments win priority over preceding tokens (such as `TokenSemi` or `TokenComma`), resolving the separator between them to `SepSpace` rather than a line break.

- **Layout Engine Punctuation Bypass**:
  In [`pkg/printer/layout_engine.go`](../../pkg/printer/layout_engine.go), the buffer resolution explicitly prevents line wrapping between commas and trailing comments:
  ```go
  case prevTokenType == dmlex.TokenComma && tokenType == dmlex.TokenComment:
      bt.ResolveWrappingCache = ResolveWrappingCache{
          SepSpace,
          PriorityMax,
          WinnerNone,
          ...
      }
  ```
  Because the resolved separator before the comment is `SepSpace` (with `PriorityMax`), the line-wrapping engine will never mark that point as a candidate for line breaking during capacity overflow simulation.
