# Line Formatting & Wrapping Guide

This guide explains how the formatter formats and wraps DML code.

---

## 1. The `black`-Style Formatting Philosophy

The DML formatter follows the formatting approach popularized by Python's
[**`black`**](https://github.com/psf/black) code formatter:
**all-or-nothing group explosion**.

Instead of allowing "ragged" partial wrapping where some items stay on
one line and others spill over awkwardly, the formatter treats related syntax as
atomic **wrap groups**:
- **If the group fits** within the configured line length, it remains on a **single line**.
- **If even a single element exceeds** the line length, the **entire group explodes** uniformly into a multi-line format with one item per line and consistent indentation.

---

## 2. Categories with Wrap Groups

The formatter groups related elements together so they wrap consistently as a unit. The primary categories are:

### Function & Method Signatures
Parameter lists and return values stay on one line unless they exceed the line limit:
```dml
// Single line when short
method read_reg(uint32 reg_id, bool sync) -> (uint32 value) throws { ... }

// Exploded when exceeding line capacity
method read_register_with_extra_flags(
    uint32 reg_id,
    uint32 transaction_flags,
    bool sync
) -> (uint32 value) throws {
    ...
}
```

### Function & Method Calls (Arguments)
Argument lists stay compact when short, or break into one argument per line:
```dml
// Single line
log info, 1, 0: "Status: %d", status;

// Exploded arguments
log info, 1, 0: "Long status message with multiple arguments: %d %s %x",
    device_status,
    status_string,
    error_code;
```

### Initializers (Arrays & Structs)
Array literals and designated struct initializers format compactly or explode:
```dml
// Fits on line
constant lookup_table = [0, 1, 2, 3];

// Exploded array
constant lookup_table = [
    0x00000001,
    0x00000002,
    0x00000004,
    0x00000008,
    0x00000010
];

// Exploded designated struct
local my_struct_t s = {
    .address = 0x1000,
    .length = 64,
    .enabled = true
};
```

### Expressions (Binary, Ternary & Member Chains, Array Indices)
Expressions wrap at clean boundary points with operator-first indentation:
```dml
// Ternary conditional expression
param mode = (flag == 0)
    ? MODE_FAST
    : MODE_SAFE;

// Chained member access
local uint32 val = (
    bus.controller.subsystem
    .peripheral.config_register
);
```

### Tuple Declarations
Multi-variable declarations break each variable definition to a dedicated line:
```dml
// Exploded tuple declaration
local (
    uint32 source_address,
    uint32 destination_address,
    bool burst_mode
) = get_transfer_spec();
```

---

## 3. Wrapping Levels & Cascading

Formatting handles line overflows in **cascading levels** from least intrusive to most expanded:

### Level 1: Single Line
When the entire statement and its expression fit within the line cap, everything stays on one line:
```dml
param a = (b + c + d);
```

### Level 2: Outer Parenthesized Wrap
When the full statement exceeds the line limit (e.g. due to a long variable or parameter name), but the inner expression fits on one line, the expression breaks into a parenthesized block while keeping the inner expression compact:
```dml
param very_long_parameter_identifier = (
    b + c + d
);
```

### Level 3: Deep Expression Explosion
When both the statement and the expression exceed the line limit, the inner expression itself explodes, placing each operand and operator on its own line:
```dml
param very_long_parameter_identifier = (
    first_long_operand_expression
    + second_long_operand_expression
    + third_long_operand_expression
);
```

---

## 4. Atomic Non-Wrapped Constructs

Certain language constructs represent tight syntactic units that
are **never broken across lines**, regardless of line length, for example:

- **Register Size & Offset**:
  The `size` specifier and `@` offset marker cling tightly to their values:
  ```dml
  register ctrl size 4 @ 0x100;
  ```
  *(The formatter will never split `size` from `4` or `@` from `0x100` onto separate lines).*

- **Field Bit Ranges**:
  Field range definitions `@ [msb:lsb]` remain a single unit:
  ```dml
  field status @ [7:0];
  field enable @ [0];
  ```

- **Bit Slices**:
  Bit slicing brackets `[high:low]` or with endian flags `[high:low, be]` cling together without interior breaks:
  ```dml
  val = buffer[31:0];
  val = packet[15:0, be];
  ```

- **Method Qualifiers (`methodQualifiers`)**:
  Qualifiers like `shared`, `independent`, `startup`, `memoized`, and `inline` cling together as an atomic group and never break among themselves. When qualifiers and the method signature exceed the line limit, wrapping occurs as a unit immediately before the `method` keyword without extra indentation:
  ```dml
  // Qualifiers cling together; line break occurs before 'method':
  shared independent startup memoized
  method foo();

  independent startup
  method foo_with_a_very_long_name();
  ```

- **Trailing Method Specifiers (`) throws default` / `) throws;`)**:
  Modifiers trailing the parameter or return list—such as `throws`, `default`, and `;`—cling tightly to the closing parenthesis `)`. They are never broken onto separate lines:
  ```dml
  // Clings directly to ')':
  method another_extremely_long_function_name() throws;

  // Even when parameters or return types explode, ') throws default' stays unified:
  method function_with_long_return_type() -> (
      very_long_return_type_1,
      very_long_return_type_2
  ) throws default {
      ...
  }
  ```

### Bypassing Non-Wrapped Restrictions via Parentheses `()`

Because the formatter treats these constructs as cohesive units,
long expressions inside them will not automatically break across lines.
If a construct contains a complex expression that exceeds the line length,
you can **bypass this restriction by wrapping the expression in parentheses `(...)`**.

Enclosing the expression in `()` introduces an explicit parenthesized wrap group,
signaling to the formatter that line breaks and multi-line indentation are permitted
if the line limit is exceeded:

- **Register Offsets & Sizes**:
  ```dml
  // Without parens: stays strictly on one line (may exceed line length)
  register status size 4 @ BASE_ADDR + SUBSYSTEM_OFFSET + REG_INDEX * 4;

  // With parens: introduces a clean multi-line wrap group
  register status size 4 @ (
      BASE_ADDR
      + SUBSYSTEM_OFFSET
      + REG_INDEX * 4
  );
  ```

- **Slice Bounds**:
  ```dml
  // Without parens: clings tightly to the brackets
  val = buffer[
    BASE_OFFSET + BLOCK_SIZE * index:BIT_NUM << VALUE, le
   ];

  // With parens: allows wrapping inside the index
  val = buffer[(
      BASE_OFFSET + BLOCK_SIZE * index
    ):(
       BIT_NUM << VALUE
    ), le
  ];
  ```

- **Member Access & Long Expressions**:
  ```dml
  // With parens: wraps long chains into an indented block
  local uint32 val = (
      subsystem.controller
      .device_register
  );
  ```

Using `(expr)` is the recommended and idiomatic escape hatch whenever
an atomic construct becomes too complex or long for a single line.

---

## 5. Technical Implementation Details (Under the Hood)

*(For developers working on the formatter internals)*

The formatting engine uses the **Unified Formatting Flags (UFF)**
layout engine defined in
[`pkg/printer/layout_engine.go`](../../pkg/printer/layout_engine.go):

- **Wrap Groups**: Managed via [`BeginWrapGroup`](../../pkg/printer/layout_engine.go)
  and [`EndWrapGroup`](../../pkg/printer/layout_engine.go)
  on the [`layoutEngine`](../../pkg/printer/layout_engine.go).
  Groups form a stack where child groups reference their parent.
- **Candidate Separators**: Tokens inside a wrap group carry
  [`FormatFlag`](../../pkg/printer/format_flag.go) flags with
  `SepSpaceMaybeWrap` or `SepNoSpaceMaybeWrap`.
- **Group Triggering**: During buffer simulation, if any token exceeds `lineCap`,
  [`triggerGroup`](../../pkg/printer/layout_engine.go) marks the assigned
  group as triggered. All candidate separators in that group are immediately
  converted into mandatory line breaks (`SepMandatoryWrap` or `SepWrapButNoSpace`).
- **Parent Cascade**: [`triggerGroup`](../../pkg/printer/layout_engine.go)
  climbs up `group.Parent`, triggering ancestor groups to maintain consistent block indentation.
