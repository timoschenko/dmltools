# dmltools

[![CI](https://github.com/timoschenko/dmltools/actions/workflows/pre-commit.yml/badge.svg)](https://github.com/timoschenko/dmltools/actions/workflows/pre-commit.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/timoschenko/dmltools.svg)](https://pkg.go.dev/github.com/timoschenko/dmltools)
[![Release](https://img.shields.io/github/v/release/timoschenko/dmltools)](https://github.com/timoschenko/dmltools/releases)
[![License: GPL-2.0](https://img.shields.io/badge/License-GPL_2.0-blue.svg)](LICENSE)

Tooling and utilities for Device Modeling Language (DML 1.4 & 1.2), including `dmlfmt` (an opinionated deterministic code formatter) and `dmlast` (an AST dumper).

## Installation

### Via `go install`

With Go 1.23+ installed:

```bash
# Install the dmlfmt code formatter
go install github.com/timoschenko/dmltools/cmd/dmlfmt@latest

# Install the dmlast AST dumper
go install github.com/timoschenko/dmltools/cmd/dmlast@latest
```

### Prebuilt Binaries

Precompiled binaries for Linux (`amd64`, `arm64`) and Windows (`amd64`) are available on the [GitHub Releases](https://github.com/timoschenko/dmltools/releases) page.

## Building from Sources

Clone the repository and build using GNU Make:

```bash
git clone https://github.com/timoschenko/dmltools.git
cd dmltools
make release
```

> **Note:** Building directly on a Windows host is not yet supported (use WSL2, Linux, or macOS).

For detailed build instructions, code generation, dependencies, and running tests, see [BUILD.md](BUILD.md).

## `dmlfmt`

An opinionated, deterministic code formatter for Device Modeling Language (DML).

### Usage

```bash
# Format files or directories in place:
dmlfmt path/to/file-a.dml path/to/file-b.dml
dmlfmt --line-length 120 src/

# Format a snippet
dmlfmt --code "method foo() -> (bool) { return 42; }"

# Check formatting without modifying files:
dmlfmt --check src/
```

For formatting rules and style specifications, see the [DML Code Formatting Guide](doc/format/README.md).

## License

This project is licensed under the GNU General Public License v2.0 only - see the [LICENSE](LICENSE) file for details.
