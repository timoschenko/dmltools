---
name: dml
description: >-
  Reference and guide for understanding Device Modeling Language (DML 1.4 & 1.2) syntax,
  formal grammar, object model, templates, and standard library specifications. Use whenever
  analyzing, writing, parsing, or formatting DML code, modifying pkg/dmlex/dml.g4, or implementing
  AST and printer features in dmlint.
---

# DML (Device Modeling Language) Syntax Guide & Reference

This skill provides comprehensive instructions, navigation tips, search recipes, and syntax reference for the Device Modeling Language (DML 1.4 & 1.2) used by Simics and parsed by `dmlint`.

---

## 1. Setup & First-Time Usage (Wiki Checkout)

The official DML Reference Manual and design documents are maintained in the GitHub wiki repository:
`https://github.com/intel/device-modeling-language.wiki.git`

The repository is cached locally at:
`.docs-cache/dml.wiki` (already configured in `.gitignore`).

### Ensuring the Wiki is Available
Before performing in-depth syntax queries or reading formal grammar specifications, check whether `.docs-cache/dml.wiki` exists:

```bash
# Option A: Run the provided helper script
./.agents/skills/dml/scripts/fetch_wiki.sh

# Option B: Run git clone directly
if [ ! -d ".docs-cache/dml.wiki" ]; then
    git clone --depth 1 https://github.com/intel/device-modeling-language.wiki.git .docs-cache/dml.wiki
fi
```

> [!NOTE]
> If running in a sandboxed environment where network access is restricted, execute the clone with `BypassSandbox: true`. Once cloned into `.docs-cache/dml.wiki`, all subsequent searches, views, and greps run safely within standard sandbox permissions.

---

## 2. Wiki File Structure & Content Overview

The `.docs-cache/dml.wiki` repository contains the following key files:

| File Path | Description | Key Topics |
| :--- | :--- | :--- |
| `refman/DML-1.4-Reference-Manual.md` | Index & Table of Contents | High-level links to all manual chapters |
| `refman/3.-Device-Modeling-Language,-version-1.4.md` | Complete DML 1.4 Specification (~4,500 lines) | Object model, statements, expressions, templates, types, parameters, methods |
| `refman/E.-Formal-Grammar.md` | Formal EBNF / BNF Grammar (~450 lines) | Production rules directly corresponding to `grammar/dml/dml.g4` |
| `refman/4.-Libraries-and-Built-ins.md` | Standard Library & Objects | Standard objects (`bank`, `register`, `field`, `connect`, etc.) and built-in methods (`init`, `read`, `write`, `get_val`, `set_val`) |
| `refman/5.-Standard-Templates.md` | Standard Behavior Templates | Register/field templates (`read_only`, `write_1_clears`, `unimpl`, `read_zero`, `zeros`, `ones`, `soft_reset_val`, etc.) |
| `refman/1.-Introduction.md` | Introduction & Overview | DML design goals, object hierarchy philosophy, compilation overview |
| `refman/2.-The-DML-compiler.md` | Compiler Tooling (`dmlc`) | CLI flags, C code generation, include paths, provisional flags |
| `refman/A.-Messages.md` | Diagnostic Messages | Compiler warnings, errors, and mitigation |
| `refman/B.-Provisional-language-features.md` | Provisional / Preview Features | Experimental features like `explicit_param_decls` |
| `refman/C.-Managing-deprecated-language-features.md` | Deprecated Features | Legacy syntax and deprecation warnings |
| `refman/D.-Changes-from-DML-1.2-to-DML-1.4.md` | DML 1.2 vs 1.4 Migration Overview | Semantic and syntax differences between 1.2 and 1.4 |
| `refman/D.1.-Porting-DML-1.2-to-DML-1.4.md` | Porting Guide | Step-by-step migration using `port-dml` |
| `refman/D.2.-Language-differences-handled-by-the-port-dml-script.md` | Automated Conversions | Syntax transformations handled automatically |
| `refman/D.3.-Backward-incompatible-changes,-not-automatically-converted.md` | Breaking Changes | Manual refactorings required when porting |
| `Roadmap-for-literals-and-initializers.md` | Design Spec | Compound literals, designated initializers, struct/array initializers |
| `Ports-in-groups.md` | Design Spec | Port scoping inside group objects |

---

## 3. Fast Topic-to-File Search Map

Use this index to quickly locate specific syntax definitions without reading entire chapters:

| Topic / Syntax Construct | Primary File | Section Heading / Line |
| :--- | :--- | :--- |
| **File header & version declaration** (`dml 1.4; device dev; bitorder be;`) | `refman/3.-Device-Modeling-Language,-version-1.4.md` | `## Source File Structure` (line 156) |
| **Formal BNF Grammar** | `refman/E.-Formal-Grammar.md` | Entire file (lines 1–445) |
| **ANTLR Grammar (in repo)** | `grammar/dml/dml.g4` | Parser & Lexer rules |
| **Objects: Banks, Registers, Fields** | `refman/3.-Device-Modeling-Language,-version-1.4.md` | `### Register Banks`, `### Registers`, `#### Fields` (lines 410–644) |
| **Objects: Connects, Interfaces, Implements** | `refman/3.-Device-Modeling-Language,-version-1.4.md` | `### Connects`, `#### Interfaces`, `### Implements` (lines 702–845) |
| **Objects: Events, Groups, Ports** | `refman/3.-Device-Modeling-Language,-version-1.4.md` | `### Events`, `### Groups`, `### Ports` (lines 846–1054) |
| **Templates & Multiple Inheritance** | `refman/3.-Device-Modeling-Language,-version-1.4.md` | `## Templates`, `### Resolving template inheritance ambiguities` (line 1055) |
| **Shared Methods & Hooks** | `refman/3.-Device-Modeling-Language,-version-1.4.md` | `### Shared methods` (line 1211), `## Hook Declarations` (line 2117) |
| **Parameters (untyped & typed)** | `refman/3.-Device-Modeling-Language,-version-1.4.md` | `## Parameters detailed` (line 1266) |
| **Data Types (int, float, bool, list, struct, layout, bitfields)** | `refman/3.-Device-Modeling-Language,-version-1.4.md` | `## Data types` (line 1451) |
| **Method Signatures, Defaults, Inline, Independent** | `refman/3.-Device-Modeling-Language,-version-1.4.md` | `## Methods` (lines 1676–2005) |
| **State Variables (`session`, `saved`)** | `refman/3.-Device-Modeling-Language,-version-1.4.md` | `## Session variables` (line 2006), `## Saved variables` (line 2059) |
| **Global Declarations (`import`, `constant`, `loggroup`, `typedef`, `extern`)** | `refman/3.-Device-Modeling-Language,-version-1.4.md` | `## Global Declarations` (line 2623) |
| **Embedded C Code (`header %{ ... %}`, `footer %{ ... %}`)** | `refman/3.-Device-Modeling-Language,-version-1.4.md` | `### Header Declarations` (line 2775) |
| **Statements (`assign`, `local`, `session`, `saved`, `return`, `delete`, `try/throw`)** | `refman/3.-Device-Modeling-Language,-version-1.4.md` | `## Method Statements` (lines 3288–3470) |
| **Control Statements (`after`, `log`, `assert`, `error`, `foreach`, `select`, `#if`)** | `refman/3.-Device-Modeling-Language,-version-1.4.md` | `### After Statements` (line 3699) to `### #if and #else Statements` (line 4162) |
| **Expressions (bit slicing `[hi:lo]`, `each in`, lists, `sizeoftype`, `cast`, `defined`)** | `refman/3.-Device-Modeling-Language,-version-1.4.md` | `## Expressions` (lines 4181–4580) |
| **Standard Register/Field Templates (`read_only`, `write_1_clears`, etc.)** | `refman/5.-Standard-Templates.md` | `## Templates for registers and fields` (lines 13–200) |
| **Standard Built-in Methods (`init`, `read`, `write`, `get`, `set`)** | `refman/4.-Libraries-and-Built-ins.md` | `## Templates for registers and fields` (line 130) |
| **DML 1.2 to 1.4 Syntax Changes** | `refman/D.-Changes-from-DML-1.2-to-DML-1.4.md` | Entire chapter D |

---

## 4. Search Recipes & Quick Grep Tips

Use these targeted commands to locate syntax details rapidly:

### 1. Lookup a Formal Grammar Rule
To find how a grammar production is defined in `E.-Formal-Grammar.md`:
```bash
# Find the rule definition and following lines
grep -n -A 10 "<dt><i><rule_name></i>" .docs-cache/dml.wiki/refman/E.-Formal-Grammar.md
# Example: find 'single_initializer' or 'statement'
grep -n -A 8 "<dt><i>single_initializer</i>" .docs-cache/dml.wiki/refman/E.-Formal-Grammar.md
```

### 2. Search Section Headings in Language Reference
To quickly find where a language concept is explained:
```bash
grep -n -E "^#{2,4} .*<keyword>" .docs-cache/dml.wiki/refman/3.-Device-Modeling-Language,-version-1.4.md
# Example: search for 'after' or 'bit slicing'
grep -n -E "^#{2,4} .*After" .docs-cache/dml.wiki/refman/3.-Device-Modeling-Language,-version-1.4.md
```

### 3. Check Standard Template Signatures & Semantics
```bash
grep -n -A 15 "^### <template_name>" .docs-cache/dml.wiki/refman/5.-Standard-Templates.md
# Example: search for 'write_1_clears'
grep -n -A 15 "^### write_1_clears" .docs-cache/dml.wiki/refman/5.-Standard-Templates.md
```

### 4. Cross-Reference ANTLR Grammar (`grammar/dml/dml.g4`)
When implementing or debugging parser and printer features:
```bash
grep -n -A 12 "^<ruleName>" grammar/dml/dml.g4
# Example:
grep -n -A 10 "^initializer" grammar/dml/dml.g4
```

---

## 5. DML 1.4 Syntax Quick Reference

### 5.1 File Skeleton
```dml
dml 1.4;
device my_device;
bitorder be; // 'be' (big-endian) or 'le' (little-endian)

import "utility.dml";

constant BUFFER_SIZE = 1024;
loggroup debug_log;
```

### 5.2 Device Hierarchy & Objects
```dml
bank regs {
    param register_size = 4;

    register cr0 size 4 @ 0x00 is (read_only) {
        param desc = "Control register 0";
        param init_val = 0x00000001;

        field enable @ [0] "Enable bit";
        field mode   @ [3:1] "Mode field";
    }

    register status size 4 @ 0x04 {
        field busy   @ [0] is (write_1_clears) "Busy flag";
    }
}

connect irq_target {
    param desc = "Target interrupt line";
    interface signal;
}

port reset_in {
    implement signal {
        method signal_raise() { ... }
        method signal_lower() { ... }
    }
}
```

### 5.3 Methods & Qualifiers
```dml
// Standard method with input arguments and return value
method calculate(uint32 input_val) -> (bool status, uint32 result) {
    if (input_val > 100) {
        return (false, 0);
    }
    return (true, input_val * 2);
}

// Inline method (untyped or typed parameters)
inline method compute_offset(index) -> (offset) {
    offset = index * 4;
}

// Independent / startup method (executed during initialization)
independent startup method init_hardware() {
    log info, 1: "Hardware initialized";
}

// Default method in template (can be overridden)
template my_tmpl {
    method handle_read() default {
        log info, 2: "Default read handler";
    }
}
```

### 5.4 Templates & Shared Methods
```dml
template reset_behavior is (base_reset) {
    param reset_level = 0;
    shared method do_reset();
}

// Instantiating templates
bank my_bank is (reset_behavior) {
    // Overriding shared method
    method do_reset() {
        ...
    }
}
```

### 5.5 Variables & Storage Classes
- `local <type> <name>;`: Stack-allocated variable local to the method execution.
- `session <type> <name>;`: Simulation session state (not saved in checkpoints).
- `saved <type> <name>;`: Checkpointed state variable preserved across checkpoint restore.

### 5.6 Key Statements & Flow Control
- **Compile-time Conditional:**
  ```dml
  #if (defined param_feature_enabled && param_feature_enabled) {
      register extra_reg size 4 @ 0x08;
  } #else {
      // alternative declarations
  }
  ```
- **Looping:**
  ```dml
  foreach item in (item_list) { ... }
  ```
- **Time/Event Scheduling:**
  ```dml
  after (0.001) { trigger_event(); }
  after: my_hook { resume_work(); }
  ```
- **Logging & Diagnostics:**
  ```dml
  log info, 2: "Register accessed with value 0x%x", val;
  log error: "Fatal error in bank %s", this.name;
  assert(expected_val == actual_val);
  error "Compilation configuration invalid";
  ```
- **Selection / Comprehensions:**
  ```dml
  select f in (this.fields) where (f.init_val != 0)
  ```

### 5.7 Expressions & Operators
- **Bit Slicing:** `val[msb:lsb]` or `val[bit]` (respects declared `bitorder`).
- **Discard Assignment:** `(_, status) = do_something();`
- **String Concatenation & Stringify:** `#` (concat), `##` (stringify).
- **List Literals:** `[1, 2, 3]`
- **Struct / Layout Initializers:** `{.field1 = val1, .field2 = val2}`
- **Type Operations:** `sizeoftype(uint32)`, `cast(expr, target_type)`, `typeof(expr)`
