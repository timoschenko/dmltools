# Building dmltools

This document describes how to build, generate parser code, and run tests for `dmltools` (`dmlint`).

> **Note:** Building directly on a Windows host is currently not supported. Please use Linux, macOS, or WSL2 (Windows Subsystem for Linux).

---

## Dependencies

### Required
- **Go**: Version 1.23 or newer ([go.dev/doc/install](https://go.dev/doc/install))
- **GNU Make**: Build automation

### For Code Generation (`make generate` / `dmlex`)
- **curl**: Required to download ANTLR JAR dependencies on the first run
- **Java Runtime**: JRE/JDK 11+ (Temurin 21 recommended) to run ANTLR 4.13.2
- **stringer**: Go symbol stringer tool
  ```bash
  go install golang.org/x/tools/cmd/stringer@latest
  ```
  *(Ensure `$(go env GOPATH)/bin` is in your `PATH`)*

### Optional / Development
- **Python 3**: Python 3.10+ with `venv` module (used by build scripts and pre-commit hooks)
- **golangci-lint**: Code linter ([golangci-lint.run](https://golangci-lint.run))
- **pre-commit**: Git pre-commit hooks framework

---

## Building

### Release Build
Builds all production binaries (`dmlfmt`, `dmlast`) and the WebAssembly demo (`cmd/dml-web-format/main.wasm`) with debug assertions and debug options disabled using the `dmltools.ndebug` build tag (similar to `NDEBUG` in C/C++):
```bash
make release
```
Compiled CLI binaries are placed in `build/bin/`.

### Debug Build
Builds binaries with debug assertions and CLI debug options enabled (the default when `dmltools.ndebug` is omitted):
```bash
make debug
```

### Building Individual Binaries
```bash
make build/bin/dmlfmt   # Build dmlfmt formatter CLI
make build/bin/dmlast   # Build dmlast AST dumper CLI
make build-wasm         # Build WebAssembly formatter
```

### Building Without Java / ANTLR
If you do not have a Java runtime and want to compile using the committed parser code:
```bash
make build/bin/dmlfmt SKIP_ANTLR=1
```

---

## Generating `dmlex` (Grammar & Parser)

When modifying the grammar at `grammar/dml/dml.g4`, regenerate the parser code in `pkg/dmlex/`:

```bash
make generate
```

This target:
1. Downloads ANTLR 4.13.2 and the grammar formatter JARs into `build/deps/` (if not already cached).
2. Formats `grammar/dml/dml.g4`.
3. Runs ANTLR to produce intermediate Go parser files.
4. Uses `tools/baker` to split and organize parser rules into `pkg/dmlex/`.
5. Runs `go generate` to create stringer definitions.

To format the grammar file without regenerating Go sources:
```bash
make format
```

> **Note:** Never manually edit generated files in `pkg/dmlex/`.

---

## Running Tests

### Full Test Suite
Runs all tests across the project with debug assertions enabled:
```bash
make test
```

### Testing with Ndebug / Release Mode
Runs tests with the `dmltools.ndebug` build tag enabled (verifying release behavior):
```bash
make test-ndebug
```

To run tests without regenerating parser files:
```bash
make test SKIP_ANTLR=1
```

### Targeted Testing
For faster development iterations, run `go test` on specific packages:
```bash
go test ./test/printer/...
go test -v ./test/printer -run TestAssertStatement
```

> **Warning:** Do not run `go test ./...` from the repository root because it will scan `./venv`. Instead, test all Go packages with:
> ```bash
> go test ./test/... ./cmd/... ./pkg/... ./tools/...
> ```

### Test Coverage
Run the test suite with coverage report:
```bash
make coverage
```
This produces `coverage.out` and an HTML report at `coverage.html`.

---

## Cleaning

Remove Go build caches, test caches, and coverage artifacts:
```bash
make clean
```
