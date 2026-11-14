# DML Language Server Features

This guide describes the editing and navigation features provided by `dmlsp`
and how they behave in day-to-day DML development.

---

## 1. Real-Time Syntax Diagnostics

`dmlsp` validates DML source syntax automatically as you work
and pushes diagnostics directly to your editor's Problems/Diagnostics view.

### When Diagnostics Run

| Editor Action | Behavior |
| :--- | :--- |
| **Open file** | Validates the file immediately and publishes diagnostics. |
| **Edit file** | Debounces validation for **150 ms** after the last keystroke so rapid typing does not flood the editor with transient errors. |
| **Save file** | Cancels any pending debounce timer and validates immediately. |
| **Close file** | Clears diagnostics for the closed file in your editor. |

### Syntax & Lexical Errors

Syntax and tokenization errors are reported with `Error` severity (source: `dml`)
highlighting the offending token or position:

```dml
dml 1.4;
device my_device; // Valid

bank regs {
    register ctrl size 4 @ 0x00 // Missing ';' -> reported as an Error diagnostic
}
```

As soon as you fix the syntax error, the diagnostic clears automatically.

### DML 1.2 Version Warning (`EUNIMPL`)

`dmlsp` targets **DML 1.4**. Opening a file with an earlier language tag (such as `dml 1.2;`) produces a `Warning` diagnostic with code `EUNIMPL` on the version number:

```dml
dml 1.2; // Warning [EUNIMPL]: Error Unimplemented: dml 1.2 syntax is not supported yet
device legacy_dev;
```

### Diagnostic Scope

`dmlsp` performs **syntax and lexical validation only**. It does not currently emit diagnostics for:
- Type mismatches, undeclared identifiers, or semantic errors checked by `dmlc`
- Unresolved `import` paths (missing imported files)
- C syntax or compiler errors inside `header %{ ... %}` and `footer %{ ... %}` blocks

---

## 2. Go to Definition

Triggering **Go to Definition** in your editor resolves DML identifiers,
dotted member paths, hierarchy keywords, templates,
and `import` statements across the current file and imported files.

### Navigating `import` Statements

Placing the cursor anywhere on an `import` statement
jumps to line 1, column 1 of the imported file:

```dml
dml 1.4;
device my_device;

import "common/helpers.dml"; // Go to Definition opens common/helpers.dml
```

#### Import Search Path Order

`dmlsp` resolves imported file paths using the following precedence:

1. **`./` paths** (e.g., `import "./local.dml";`): Resolved strictly relative to the directory of the file containing the `import`.
2. **`../` paths** (e.g., `import "../shared.dml";`): Checked relative to the importing file's directory first; if not found, falls back to the search order below.
3. **Bare paths** (e.g., `import "dml-builtins.dml";` or `import "sub/defs.dml";`): Searched in order across the importing file's directory and the primary workspace root:
   - `<importing_file_dir>/1.4/<path>`
   - `<importing_file_dir>/<path>`
   - `<workspace_root>/1.4/<path>`
   - `<workspace_root>/<path>`

### Templates & Template Inheritance (`template`, `is`)

You can jump to a template definition either from the template name itself
or from the `is` keyword in an `is <template>` clause:

```dml
dml 1.4;
device my_device;

template base_reg {
    param reset_val = 0;
}

template rw_reg is base_reg {
    param writable = true;
}

bank regs {
    register ctrl size 4 @ 0x00 is rw_reg { // Cursor on 'is' or 'rw_reg' -> template rw_reg
        method check() -> (int) {
            return reset_val; // Resolves through rw_reg -> template base_reg
        }
    }
}
```

Multi-level template inheritance and device-level templates (`is my_device_template;` at file top level) are both supported.

### Device Hierarchy, Dotted Chains & Keywords (`this`, `parent`, `dev`)

For dotted member chains, `dmlsp` resolves **up to the specific segment under your cursor**:

```dml
dml 1.4;
device my_device;

bank regs {
    register r1 size 4 @ 0x00 {
        field f1 @ [7:0];
        method inspect() {
            this.parent.r2; // Cursor on 'r2' -> register r2; on 'parent' -> bank regs
        }
    }
    register r2 size 4 @ 0x04 {
        field f2 @ [0];
    }
}

method top() {
    regs.r1.f1; // Cursor on 'regs' -> bank regs; on 'r1' -> register r1; on 'f1' -> field f1
    dev.regs;   // 'dev' resolves to device my_device
}
```

Supported hierarchy keywords:
- **`this`**: Resolves to the nearest enclosing object (`bank`, `register`, `field`, `group`, `template`, etc.).
- **`parent`**: Resolves to the enclosing object's parent (works both at the start of a chain and mid-chain like `this.parent.r2`).
- **`dev`**: Resolves to the top-level `device` declaration and its children or attached templates.

### Methods, Parameters & Named Return Values

Go to Definition resolves:
- Method calls (`compute(5);`) to the `method` declaration or definition.
- Input parameters (`int count`) and named output parameters (`-> (int out_val)`) used inside a method body back to their declaration in the method signature:

```dml
method compute(int count) -> (int out_val) {
    local int step = 1;
    out_val = count + step; // 'out_val', 'count', and 'step' all jump to their declarations
}
```

### Variables & Declarations

The following DML declarations are indexed and navigable within their lexical, object, and imported scopes:
- **Objects**: `device`, `bank`, `register`, `field`, `attribute`, `connect`, `event`, `group`, `implement`, `interface`, `port`, `subdevice`
- **Templates**: `template`
- **Methods & Parameters**: `method` declarations and definitions, method input parameters, and named return parameters
- **Parameters & Constants**: `param`, `constant`
- **Variables & State**: `local`, `session`, `saved`
- **Other Declarations**: `hook`, `loggroup`, `typedef`, `extern`, `extern typedef`

### Cross-File & Unsaved Buffer Resolution

`dmlsp` follows both direct and transitive `import` chains (with cycle detection so circular imports terminate cleanly).
- If an imported file is **open in your editor**, `dmlsp` resolves definitions against your live, unsaved in-memory edits.
- If an imported file is **not open**, `dmlsp` automatically reads it from disk.

---

## 3. Known Limitations

- **Single Target Location**: Go to Definition returns at most one target location. When a `param` or `default method` is declared in a template and overridden in an object, navigating to it resolves to the highest-precedence definition in scope order rather than offering a list of all overrides.
- **Import Search Directories**: `dmlsp` searches only the importing file's directory and the primary workspace root (plus their `1.4/` subdirectories). Additional external include paths (such as custom `-I` directories outside the workspace) cannot currently be configured.
- **`in each` Statements**: Template applications (`in each ... is ...`) and declarations inside `in each` blocks are not indexed for member resolution.
- **Embedded C Blocks (`%{ ... %}`)**: Symbols defined inside `header %{ ... %}` and `footer %{ ... %}` C blocks are not indexed.
- **Flat Method Scope for Locals**: Local variables declared inside nested `{ ... }` blocks within the same method share a single method-level scope; if two nested blocks declare local variables with the same name, navigation resolves to the last declaration in that method.
- **Buffers with Active Syntax Errors**: When a file currently contains syntax errors, Go to Definition relies on a best-effort partial parse of the current buffer, so symbols after the syntax error may fail to resolve until the error is fixed.
