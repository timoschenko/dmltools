# LSP API Compatibility & Limitations

This document provides the **Language Server Protocol (LSP 3.17)**
compatibility matrix for `dmlsp`, listing all supported LSP commands,
their specific limitations, and unsupported LSP commands.

---

## 1. Protocol Overview

| Property | Value |
| :--- | :--- |
| **LSP Specification** | 3.17 |
| **Transport** | JSON-RPC 2.0 over `stdin` / `stdout` |
| **Position Encoding** | `utf-16` (UTF-16 code units, including full surrogate-pair support) |
| **Document Sync Mode** | Incremental (`openClose: true`, `change: 2`, `save: {}`) |

---

## 2. Supported LSP API Commands

| Method | Direction | Type | Support Level | Summary |
| :--- | :--- | :--- | :--- | :--- |
| `initialize` | Client -> Server | Request | Supported (with limitations) | Negotiates capabilities and sets primary workspace root. |
| `initialized` | Client -> Server | Notification | Supported | Completes initialization handshake. |
| `shutdown` | Client -> Server | Request | Supported | Transitions server to shut-down state prior to `exit`. |
| `exit` | Client -> Server | Notification | Supported | Terminates server process (`0` after `shutdown`, `1` otherwise). |
| `textDocument/didOpen` | Client -> Server | Notification | Supported | Registers in-memory buffer and runs immediate diagnostics. |
| `textDocument/didChange` | Client -> Server | Notification | Supported (with limitations) | Applies incremental or full edits and schedules debounced diagnostics. |
| `textDocument/didSave` | Client -> Server | Notification | Supported (with limitations) | Updates buffer if text is provided and runs immediate diagnostics. |
| `textDocument/didClose` | Client -> Server | Notification | Supported (with limitations) | Clears diagnostics for the file and drops the in-memory buffer. |
| `textDocument/publishDiagnostics` | Server -> Client | Notification | Supported (with limitations) | Pushes syntax errors and DML 1.2 `EUNIMPL` warnings. |
| `textDocument/definition` | Client -> Server | Request | Supported (with limitations) | Resolves symbols, hierarchies, templates, and `import` paths. |

---

### Detailed Behavior & Limitations

#### Lifecycle

- **`initialize`** *(Request, Client -> Server)*
  - **Behavior**: Records client info, determines the workspace root, transitions the server to the initialized state, and advertises:
    - `positionEncoding`: `"utf-16"`
    - `textDocumentSync`: `{ openClose: true, change: 2 (Incremental), save: {} }`
    - `definitionProvider`: `true`
  - **Limitations**:
    - **Single Workspace Root**: Only `workspaceFolders[0].uri` (or `rootUri` if `workspaceFolders` is empty) is used as the workspace root for `import` path lookup. Additional folders in multi-root workspaces are ignored.
    - **No `initializationOptions`**: Custom initialization options (such as extra import include directories or custom diagnostic debounce settings) are not supported.
    - **Single Initialization**: Calling `initialize` more than once returns an `already initialized` error.

- **`initialized`** *(Notification, Client -> Server)*
  - **Behavior**: Logs completion of the initialization handshake.
  - **Limitations**: Does not perform dynamic capability registration (`client/registerCapability`) or workspace configuration queries.

- **`shutdown`** *(Request, Client -> Server)*
  - **Behavior**: Transitions the server to the shut-down state so the subsequent `exit` notification exits cleanly.
  - **Limitations**: Calling `shutdown` before `initialize` returns a `server not initialized` error.

- **`exit`** *(Notification, Client -> Server)*
  - **Behavior**: Closes the JSON-RPC connection and terminates the process with exit code `0` if `shutdown` was previously received, or exit code `1` if `exit` was sent without a prior `shutdown`.
  - **Limitations**: None.

---

#### Document Synchronization

- **`textDocument/didOpen`** *(Notification, Client -> Server)*
  - **Behavior**: Stores the opened document content in the server's in-memory overlay (taking precedence over the on-disk file) and immediately publishes syntax diagnostics.
  - **Limitations**: None.

- **`textDocument/didChange`** *(Notification, Client -> Server)*
  - **Behavior**: Supports both incremental range-based edits (`range` + `text`) and full-document replacements (`text`). Invalidates stale cached state for the modified file and schedules syntax diagnostics after a **150 ms** debounce delay.
  - **Limitations**: The 150 ms debounce interval is fixed and cannot be customized via client settings.

- **`textDocument/didSave`** *(Notification, Client -> Server)*
  - **Behavior**: Cancels any pending diagnostic debounce timer, refreshes the in-memory buffer if `params.text` is included in the notification, and immediately publishes diagnostics.
  - **Limitations**: Does not instruct the client to always include full text on save (`save.includeText` is omitted).

- **`textDocument/didClose`** *(Notification, Client -> Server)*
  - **Behavior**: Cancels any pending diagnostic debounce timer, publishes an empty diagnostics list (`[]`) to clear editor markers for the document, and removes the file from the in-memory overlay so subsequent cross-file lookups read from disk.
  - **Limitations**: Because diagnostics are cleared on close, syntax errors in closed files do not persist in workspace-wide Problems panels.

---

#### Diagnostics

- **`textDocument/publishDiagnostics`** *(Notification, Server -> Client)*
  - **Behavior**: Pushes parser and lexical errors (`Severity: Error`, `Source: "dml"`) and DML version compatibility warnings (`Severity: Warning`, `Code: "EUNIMPL"`, `Source: "dml"` for non-1.4 version headers such as `dml 1.2;`).
  - **Limitations**:
    - **Syntax Only**: Reports grammar and lexical errors only; does not perform semantic analysis (such as type checking, undeclared symbol detection, or duplicate definition checks).
    - **DML 1.4 Target**: Files starting with `dml 1.2;` emit an `EUNIMPL` warning and skip further parsing of the file body.
    - **No Import Validation**: Unresolved or missing `import` file paths do not generate diagnostic errors or warnings.
    - **Opaque C Blocks**: Embedded C code inside `header %{ ... %}` and `footer %{ ... %}` blocks is not parsed or checked for C syntax errors.

---

#### Language Features

- **`textDocument/definition`** *(Request, Client -> Server)*
  - **Behavior**: Resolves DML identifiers, dotted member access chains (`a.b.c`), hierarchy keywords (`this`, `parent`, `dev`), `template` references (`is ...`), and `import "..."` targets across open buffers and on-disk imported files.
  - **Limitations**:
    - **Single Location Result**: Returns a single `Location` object rather than `Location[]` or `LocationLink[]`. When a `param` or `default method` is declared in a template and overridden in an object, only the highest-precedence match in scope resolution order is returned.
    - **Import Search Scope**: Searches only the importing file's directory and the primary workspace root (plus their `1.4/` subdirectories); external include directories (`-I` paths) outside the workspace are not supported.
    - **`in each` Blocks**: Template applications and declarations inside `in each` statements are not indexed or resolved.
    - **Embedded C Blocks**: Symbols defined inside `header %{ ... %}` and `footer %{ ... %}` C blocks are not indexed.
    - **Flat Method Scope**: Local variables inside nested `{ ... }` blocks within the same method share a flat method-level scope; shadowed names inside separate nested blocks resolve to the last declaration in that method.
    - **Files with Syntax Errors**: When a buffer contains syntax errors, resolution relies on a best-effort partial parse of the current buffer, so symbols located after a parse failure may not resolve.

---

## 3. Unsupported LSP API Commands

Any LSP 3.17 request or notification not listed in Section 2 is **unsupported** (`dmlsp` does not advertise capabilities for these features and returns method-not-found or empty responses if invoked).

### Formatting

| Method | Direction | Notes |
| :--- | :--- | :--- |
| `textDocument/formatting` | Client -> Server | Unsupported in `dmlsp`. Use the standalone `dmlfmt` CLI tool for DML code formatting (see the [DML Code Formatting Guide](../format/README.md)). |
| `textDocument/rangeFormatting` | Client -> Server | Unsupported in `dmlsp`. |
| `textDocument/onTypeFormatting` | Client -> Server | Unsupported in `dmlsp`. |

### Code Navigation

| Method | Direction | Notes |
| :--- | :--- | :--- |
| `textDocument/declaration` | Client -> Server | Unsupported. |
| `textDocument/typeDefinition` | Client -> Server | Unsupported. |
| `textDocument/implementation` | Client -> Server | Unsupported. |
| `textDocument/references` | Client -> Server | Unsupported (Find All References is not available). |

### Coding Assistance & Intelligence

| Method | Direction | Notes |
| :--- | :--- | :--- |
| `textDocument/hover` | Client -> Server | Unsupported (no hover documentation or type tooltips). |
| `textDocument/completion` | Client -> Server | Unsupported (no autocomplete suggestions). |
| `completionItem/resolve` | Client -> Server | Unsupported. |
| `textDocument/signatureHelp` | Client -> Server | Unsupported (no parameter hints while typing method calls). |
| `textDocument/inlayHint` | Client -> Server | Unsupported. |
| `inlayHint/resolve` | Client -> Server | Unsupported. |
| `textDocument/inlineValue` | Client -> Server | Unsupported. |
| `textDocument/moniker` | Client -> Server | Unsupported. |

### Symbols & Document Outline

| Method | Direction | Notes |
| :--- | :--- | :--- |
| `textDocument/documentSymbol` | Client -> Server | Unsupported (no editor Outline / Breadcrumbs tree). |
| `workspace/symbol` | Client -> Server | Unsupported (no workspace-wide symbol search). |
| `workspaceSymbol/resolve` | Client -> Server | Unsupported. |

### Refactoring & Code Actions

| Method | Direction | Notes |
| :--- | :--- | :--- |
| `textDocument/rename` | Client -> Server | Unsupported (no symbol renaming). |
| `textDocument/prepareRename` | Client -> Server | Unsupported. |
| `textDocument/codeAction` | Client -> Server | Unsupported (no quick fixes or refactorings). |
| `codeAction/resolve` | Client -> Server | Unsupported. |
| `textDocument/codeLens` | Client -> Server | Unsupported. |
| `codeLens/resolve` | Client -> Server | Unsupported. |

### Highlighting, Folding & Selection

| Method | Direction | Notes |
| :--- | :--- | :--- |
| `textDocument/documentHighlight` | Client -> Server | Unsupported. |
| `textDocument/semanticTokens/full` | Client -> Server | Unsupported (syntax highlighting relies on editor grammars). |
| `textDocument/semanticTokens/full/delta` | Client -> Server | Unsupported. |
| `textDocument/semanticTokens/range` | Client -> Server | Unsupported. |
| `textDocument/foldingRange` | Client -> Server | Unsupported (folding relies on editor indentation/grammar). |
| `textDocument/selectionRange` | Client -> Server | Unsupported. |
| `textDocument/linkedEditingRange` | Client -> Server | Unsupported. |
| `textDocument/documentLink` | Client -> Server | Unsupported (use `textDocument/definition` on `import` statements instead). |
| `documentLink/resolve` | Client -> Server | Unsupported. |
| `textDocument/documentColor` | Client -> Server | Unsupported. |
| `textDocument/colorPresentation` | Client -> Server | Unsupported. |

### Call & Type Hierarchy

| Method | Direction | Notes |
| :--- | :--- | :--- |
| `textDocument/prepareCallHierarchy` | Client -> Server | Unsupported. |
| `callHierarchy/incomingCalls` | Client -> Server | Unsupported. |
| `callHierarchy/outgoingCalls` | Client -> Server | Unsupported. |
| `textDocument/prepareTypeHierarchy` | Client -> Server | Unsupported. |
| `typeHierarchy/supertypes` | Client -> Server | Unsupported. |
| `typeHierarchy/subtypes` | Client -> Server | Unsupported. |

### Pull Diagnostics

| Method | Direction | Notes |
| :--- | :--- | :--- |
| `textDocument/diagnostic` | Client -> Server | Unsupported (`dmlsp` uses push-based `textDocument/publishDiagnostics`). |
| `workspace/diagnostic` | Client -> Server | Unsupported. |

### Document Sync (Will-Save)

| Method | Direction | Notes |
| :--- | :--- | :--- |
| `textDocument/willSave` | Client -> Server | Unsupported. |
| `textDocument/willSaveWaitUntil` | Client -> Server | Unsupported. |

### Workspace & File Operations

| Method | Direction | Notes |
| :--- | :--- | :--- |
| `workspace/didChangeWorkspaceFolders` | Client -> Server | Unsupported (restart server or reopen workspace to change workspace root). |
| `workspace/didChangeConfiguration` | Client -> Server | Unsupported. |
| `workspace/didChangeWatchedFiles` | Client -> Server | Unsupported. |
| `workspace/executeCommand` | Client -> Server | Unsupported. |
| `workspace/applyEdit` | Server -> Client | Unsupported. |
| `workspace/willCreateFiles` | Client -> Server | Unsupported. |
| `workspace/didCreateFiles` | Client -> Server | Unsupported. |
| `workspace/willRenameFiles` | Client -> Server | Unsupported. |
| `workspace/didRenameFiles` | Client -> Server | Unsupported. |
| `workspace/willDeleteFiles` | Client -> Server | Unsupported. |
| `workspace/didDeleteFiles` | Client -> Server | Unsupported. |
