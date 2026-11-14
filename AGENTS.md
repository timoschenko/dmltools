# AGENTS.md

## Build, Test & Verification Commands

- **Default verification before finishing a task**:
  1. `go test -timeout 60s ./test/... ./cmd/... ./pkg/...` (debug assertions enabled)
  2. `go test -tags dmltools.ndebug -timeout 60s ./test/... ./cmd/... ./pkg/...` (release mode / `NDEBUG`)
  3. `golangci-lint run`
  4. `make license-check`
- **Do NOT run `go test ./...` from the repo root**: It scans `./venv` and standalone sub-modules under `./tools/`. Always pass `./test/... ./cmd/... ./pkg/...`.
- **Single test / package iteration**:
  - `go test -v ./test/printer -run TestAssertStatement`
  - Most unit and integration suites live under `./test/{dmlex,printer,lsp,fs}`, not alongside `./pkg/`.
- **Standalone tool modules**: `tools/baker` and `tools/licensecheck` have their own `go.mod` files. Test them via `make test-tools`.
- **`make` & ANTLR gotcha (`SKIP_ANTLR=1`)**:
  - Makefile build/test targets (`make test`, `make release`, `make debug`) depend on `build-gocache`, which triggers `make generate` (requiring Java 11+, `curl`, and `stringer`) whenever `build/env/*.flag` files are missing or older than `grammar/dml/dml.g4`.
  - Use `go test` directly or pass `SKIP_ANTLR=1` (e.g., `make test SKIP_ANTLR=1`, `make build/bin/dmlfmt SKIP_ANTLR=1`) to avoid triggering ANTLR codegen.
  - Makefile also overrides `GOCACHE` to `build/cache/golang` unless `GOCACHE` is already set in the environment.
- **Full CI pipeline order** (`.github/workflows/pre-commit.yml`):
  `make test-tools` -> `make generate` -> `make release` -> `make test` -> `make test-ndebug` -> `pre-commit run --all-files` (`golangci-lint`) -> `make license-check`.

## Grammar & Codegen Pipeline (`grammar/dml/dml.g4` -> `pkg/dmlex/`)

- **Approval required**: **Never** modify `grammar/dml/dml.g4` or run `make generate` without explicit user approval.
- **What `make generate` does** (when approved):
  1. Formats `grammar/dml/dml.g4` using the cached ANTLR formatter JAR (`make format`).
  2. Runs ANTLR 4.13.2 to produce intermediate Go parser/lexer files in `grammar/dml/`.
  3. Runs `tools/baker` to split monolithic `dml_parser.go` into per-rule context files (`pkg/dmlex/dml_parser_*_context.go`) and removes `grammar/dml/dml_*.go`.
  4. Runs `go generate ./pkg/dmlex/...` (`stringer -type=TokenType` -> `pkg/dmlex/tokentype_string.go`).
- **Generated vs. hand-written files in `pkg/dmlex/`**:
  - **Generated (never edit manually)**: `dml_parser_*.go`, `dml_lexer.go`, `dml_listener.go`, `dml_base_listener.go`, `tokentype_string.go`.
  - **Hand-written (safe to edit)**: `file.go`, `porting.go`, `expression.go`, `cdecl.go`, `method.go`, `imports.go`, `synthetic_token.go`, `error_listener.go`, `version.go`, `doc.go`.

## Architecture & Package Map

- **Entrypoints (`cmd/`)**:
  - `cmd/dmlfmt`: CLI code formatter (`build/bin/dmlfmt`). Note: sets `antlr.CommonTokenFactoryDEFAULT = antlr.NewCommonTokenFactory(true)` in `init()` so tokens retain zero-copy string slices backed by `stream.HashedBuffer`.
  - `cmd/dmlast`: AST dumper CLI (`build/bin/dmlast`).
  - `cmd/dmlsp`: LSP 3.17 server over stdio (`build/bin/dmlsp`).
  - `cmd/dml-web-format`: WebAssembly formatter playground (`make build-wasm` -> `cmd/dml-web-format/main.wasm`).
- **Core Packages (`pkg/`)**:
  - `pkg/dmlex`: ANTLR4 lexer/parser wrapper, AST rule contexts, `TokenType` enum, import resolution helpers (`imports.go`), and DML version detection (`version.go`).
  - `pkg/printer`: AST traversal (`dml_*.go`), scoped `PrintContext` hierarchy, and Unified Formatting Flags (`pkg/printer/format`, 32-bit `FormatFlag`), backed by default token maps (`pkg/printer/dml_defaults`) and the line-wrapping engine (`pkg/printer/layout`).
  - `pkg/lsp`: LSP server (`Server`), in-memory overlay/AST/symbol cache (`FsCache`), syntax diagnostics, and go-to-definition (`definition.go`, `symbols.go`).
    - **UTF-16 vs. byte offset gotcha**: LSP `protocol.Position.Character` is in UTF-16 code units, whereas ANTLR and Go strings use UTF-8 byte/rune offsets. Always convert positions via `LineAndByteCol` / `utf16ColToByteCol` (`pkg/lsp/cache.go`).
  - `pkg/fs`: Atomic file replacement (`ReplaceFileAtomically`), path normalization, `pkg/fs/stream` (`HashedBuffer` lazy UTF-8 `antlr.CharStream` + 192-bit composite hash), and `pkg/fs/lookup` (stat-validated file discovery, `.gitignore`/`.editorconfig` chain lookup).

## Printer Constraints (`pkg/printer`)

- **Delegate all layout to `layout.Engine`**: AST visitor methods in `pkg/printer/dml_*.go` **must not** perform manual string length math, manual spacing, or manual line wrapping. Push tokens via `processToken`, `pushTokenRange`, `pushTokenRangeFromRightiest`, or `pushTokenWithCustomFormat` inside scoped `PrintContext` overrides (`p.newContext(...)` / `defer p.freeContext()`).
- **Sole newline exception**: Blank lines between top-level/block statements are emitted via `writeNewLines()` at statement boundaries (`pkg/printer/statement_spacing.go`).
- **Rule annotations & passthrough**:
  - Prefix every grammar rule visitor method with `// g4-tag: <ruleName>`.
  - Unimplemented grammar rules must call `p.printUnimplemented(node)` for raw source passthrough.
- **No unauthorized workarounds**: **Never** add new AST/token workarounds to `pkg/printer/workarounds.go` without explicit user approval. Fix the root cause in `pkg/printer` or `pkg/printer/layout`, or ask first.
- **Debug assertions**: `assert()` and `assertLazy()` in `pkg/printer/debug.go` panic during normal tests (`!dmltools.ndebug`) and compile to no-ops under `-tags dmltools.ndebug`.

## Testing Quirks (`test/printer`)

- **Backtick quote substitution**: In `printerTestCase` (`test/printer/testlib_test.go`), all backticks (`` ` ``) in `code`, `expected`, and `masked` strings are automatically converted to double quotes (`"`) before parsing and comparison.
- **Implicit headers & method wrapping**:
  - `runPrinterTestCase` auto-prepends `"dml 1.4;\n\n\n"` and strips it from the formatted result unless `expected` already contains a `dml` header.
  - `runPrinterMethodTestCase` wraps snippets in `method foo() { ... }` indented by 4 spaces and adds `+4` to `lineCap` (standard test `lineCap` is `40`).
- **Test masking policy (`test/printer/mask_test.go`)**:
  - **Never** add new `masked: Mask...` entries or `task006Skip: true` flags to silence failing tests without explicit user approval.
  - When fixing a bug covered by an existing mask, remove the mask—if a masked test produces `expected`, the runner fails with `"test case passed but it is still masked"`.

## Enforced Repo Conventions

- **SPDX License Header (`make license-check`)**: Every non-generated `.go` file and `grammar/dml/dml.g4` must start with:
  ```go
  // SPDX-License-Identifier: GPL-2.0-only
  // Copyright (C) 2026 Sergei Timoshchenko
  ```
  (`GPL-2.0-or-later` is explicitly rejected by `tools/licensecheck`). If local untracked directories with third-party `LICENSE` files cause `make license-check` to fail, pass explicit file paths: `(cd tools/licensecheck && go run . -root ../.. <files...>)`.
- **Package documentation (`test/doc_test.go` & `make doc`)**: Every package under `pkg/` containing non-test `.go` files must have a package doc comment (`doc.go`) starting with `"Package <name> ..."`. When adding a new subpackage under `pkg/`, also add it to `TestGoDocCommand` in `test/doc_test.go`.
- **Go formatting & linting (`.golangci.yaml`)**: Formatted with `gofmt` and `golines` (`max-len: 80`).
- **Commit messages (`.goreleaser.yaml`)**: Follow Conventional Commits (`<type>(<scope>): <description>`) using types `feat`, `fix`, `perf`, `docs`, `refactor`, `chore`, `test`, `ci`, `infra`.

## Existing Instruction & Reference Files

- `.agents/skills/dml/SKILL.md`: DML 1.4/1.2 syntax and formal grammar lookup guide (queries `.docs-cache/dml.wiki`, fetched via `./.agents/skills/dml/scripts/fetch_wiki.sh`).
- `.agents/skills/dml_printer_feature_bootstrap/SKILL.md`: TDD workflow for scaffolding new printer rules and tests.
- `doc/format/README.md`, `doc/format/line_wrapping.md`, `doc/format/comments.md`: DML formatter style, wrap-group explosion, and sticky comment rules.
- `doc/lsp/README.md`, `doc/lsp/features.md`, `doc/lsp/compatibility.md`: `dmlsp` capabilities, import resolution, and LSP 3.17 support.
