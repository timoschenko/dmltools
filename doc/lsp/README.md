# DML Language Server (`dmlsp`)

`dmlsp` is the Language Server Protocol (LSP) server for Device Modeling Language.
It provides IDE and editor integrations with real-time syntax feedback and code navigation.

---

## Overview

`dmlsp` runs as a background process managed by your code editor or IDE and
communicates using the **Language Server Protocol (LSP 3.17)**
over **JSON-RPC 2.0** via standard I/O (`stdin`/`stdout`).

### Key Capabilities

- **Real-Time Syntax Diagnostics**: Validates `.dml` files as you open, edit, and save them, reporting syntax errors and DML version warnings directly in your editor.
- **Go to Definition**: Navigates across symbols and files, including:
  - `import` statements to the target `.dml` files
  - `template` definitions and inherited template members (including multi-level template inheritance)
  - Device object hierarchies (`device`, `bank`, `register`, `field`, `group`, `port`, `connect`, `interface`, `attribute`, `event`, `implement`, `subdevice`)
  - Dotted member chains (such as `regs.r1.f1`) and hierarchy keywords (`this`, `parent`, `dev`)
  - `method` declarations and definitions, input parameters, and named return parameters
  - `local`, `session`, and `saved` variables, `param`, `constant`, `hook`, `loggroup`, `typedef`, and `extern` declarations
- **Live Editor Synchronization**: Tracks unsaved edits in memory (supporting incremental updates) while automatically falling back to on-disk files when resolving unopened imported files.

---

## Command-Line Usage

`dmlsp` is normally launched automatically by your editor's LSP client rather than run manually in a terminal.

```bash
dmlsp [options]
```

### Options

| Flag | Description |
| :--- | :--- |
| `--log-file <path>`, `--logfile <path>` | Write server logs to `<path>` instead of `stderr` (created or appended with `0600` permissions). |
| `-v`, `--verbose` | Enable verbose (`DEBUG`-level) logging. By default, logs at `INFO` level. |
| `-V`, `--version` | Print version, commit, build date, and license information and exit. |
| `-h`, `--help` | Print command usage information and exit. |

### Exit Codes

| Exit Code | Meaning |
| :--- | :--- |
| `0` | Clean shutdown (`shutdown` request received before `exit`), or `--help` / `--version` completed. |
| `1` | Unclean exit (`exit` notification received without a prior `shutdown` request), or startup/log-file error. |
| `2` | Invalid command-line flag or argument. |

---

## Documentation Guides

- [**Installation & Editor Setup**](setup.md): How to install `dmlsp` and configure it in Neovim, VS Code, and Zed.
- [**Features & Navigation Guide**](features.md): Detailed examples of diagnostics, Go to Definition, import path resolution, and user-visible limitations.
- [**LSP API Compatibility & Limitations**](compatibility.md): Complete reference of supported LSP 3.17 API commands, their limitations, and unsupported methods.
