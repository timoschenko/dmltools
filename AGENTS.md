# dmlint Agent Instructions

## Architecture & Boundaries
- **Module Name:** The Go module is `github.com/timoschenko/dmltools` (not `dmlint`). Internal imports must use `github.com/timoschenko/dmltools/pkg/...` or `github.com/timoschenko/dmltools/cmd/...`.
- **Primary Branch:** `master`.
- **Entrypoints (`cmd/`):**
  - `cmd/dmlfmt`: Main DML formatter CLI (`dmlfmt [options] [src]...`).
  - `cmd/dmlast`: AST dump tool emitting JSON.
  - `cmd/dml-web-format`: WebAssembly formatter demo (`make build-wasm`, `make web` on port 8080).
- **Core Packages (`pkg/`):**
  - `pkg/dmlex`: Lexer/parser generated from ANTLR4 grammar (`grammar/dml/dml.g4`), decomposed via `tools/baker` into per-rule context files (`dml_parser_*_context.go`).
  - `pkg/printer`: Formatter visitor traversing DML AST nodes (`// g4-tag: <rule>`).
  - `pkg/printer/layout`: Line wrapping, indentation, conflict resolution, token buffering (`layout.Engine`).
  - `pkg/printer/format`: Unified Formatting Flags (`format.FormatFlag`), break priorities, separators, `TokenMap`.
  - `pkg/printer/dml_defaults`: Default token format maps (cdecl, expression, tokens).
  - `pkg/fs` & `pkg/fs/stream`: File system abstraction, caching (`FsCache`), and streaming.
- **Package Documentation:** Every package under `pkg/` must include a `doc.go` or package comment starting with `Package <name>`. `test/doc_test.go` and `make doc` enforce this.

## Developer Commands

```bash
# Full test suite (builds gocache, runs tests across all valid packages)
make test

# Fast focused testing (bypasses ANTLR codegen)
go test ./test/printer/...
go test -v ./test/printer -run TestAssertStatement
go test ./test/... ./cmd/... ./pkg/...              # Never run 'go test ./...' from root (scans ./venv)

# Build binaries (outputs to build/bin/dmlfmt, build/bin/dmlast)
make build/bin/dmlfmt
make release

# Build/test without Java/ANTLR runtime (uses committed parser sources)
make test SKIP_ANTLR=1
make build/bin/dmlfmt SKIP_ANTLR=1

# Linting and formatting Go code
golangci-lint run                       # Enforces gofmt and golines (max-len: 80, tab-len: 4)
golangci-lint fmt                       # Automatically format Go code

# Grammar & codegen
make format                             # Format grammar/dml/dml.g4
make generate                           # Format grammar + ANTLR codegen + baker split + go generate
make doc                                # Verify package doc coverage across pkg/
make clean                              # Clear Go build, test, fuzz, and coverage caches
```

## Codegen & Grammar Changes
- **Generated Code:** Never manually edit `pkg/dmlex/dml_*.go`, `pkg/dmlex/*_string.go`, `cmd/dml-web-format/wasm_exec.js`, `*.tokens`, or `*.interp`.
- **Regeneration Requirements:** `make generate` requires:
  - Java runtime 11+ (Temurin 21 recommended) to execute ANTLR 4.13.2.
  - `stringer` tool: `go install golang.org/x/tools/cmd/stringer@latest`.
  - Virtualenv `./venv` (automatically provisioned by Makefile for whitespace/EOF fixers).
- **CI Enforcement:** When `CI=true`, `make generate` verifies generated files match committed files in git and fails if there is any diff. Always run `make generate` and commit all generated files after modifying `grammar/dml/dml.g4`.

## Printer Architecture & Constraints
- **Layout Engine Delegation:** AST visitor methods in `pkg/printer/` push `PrintContext` overrides with `FormatFlag`s. Formatting, spacing, and line wrapping decisions must be delegated entirely to `layout.Engine`.
- **No Manual Wrapping:** AST methods must not perform manual string length calculations or manually inject spacing/newlines for line wrapping within statements.
- **Statement Boundary Exception:** `writeNewLines()` is permitted only between top-level statements, as `layout.Engine` operates within statements.
- **Grammar Tagging:** Every AST print method must be tagged with its grammar rule: `// g4-tag: <rule_name>`.
- **Synthetic Tokens (Epic 006):** Do not manually inject synthetic tokens. Use `task006Skip: true` in test cases requiring synthetic parens.

## Testing Harness & Masking (`test/printer/testlib_test.go`)
- **DML Header Preprocessing:** If test input lacks a `dml` header, `dml 1.4;\n\n\n` is automatically prepended and stripped from output.
- **Backtick Normalization:** Backticks (`` ` ``) in test code and expected strings are normalized to double quotes (`"`), allowing raw Go string literals to express DML double quotes cleanly.
- **Masked Tests (`test/printer/mask_test.go`):** Known issues and WIP formatting use typed masks (`MaskIncorrectFormatting`, `MaskSkip`, `MaskSyntaxError`).
- **Mask Removal Gotcha:** When a bugfix causes a masked test to output the correct expected string, the test harness triggers `t.Errorf("test case passed but it is still masked. Consider removing the mask")`. You **must** remove the mask once the test passes.

## Code Conventions & Workarounds
- **AI Workarounds:** Every temporary fix or workaround must strictly follow the format:
  `// TODO: AI WA: <description>` and `// Deprecated: This is workaround`.
- **Workaround Isolation:** Isolate non-trivial workaround logic into dedicated helper functions in `pkg/printer/workarounds.go` or `pkg/printer/layout/workarounds.go`.

## Key Skills & References
- **DML 1.4 / 1.2 Syntax Spec:** Cached in `.docs-cache/dml.wiki` (fetch via `.agents/skills/dml/scripts/fetch_wiki.sh`); guide in `.agents/skills/dml/SKILL.md`.
- **Scaffolding New Printer Rules:** Workflow and checklist in `.agents/skills/dml_printer_feature_bootstrap/SKILL.md`.
