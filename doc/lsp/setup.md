# Installation & Editor Setup

This guide explains how to install `dmlsp` and configure it in Neovim, VS Code, and Zed.

---

## 1. Installing `dmlsp`

Choose one of the following installation methods, then verify that `dmlsp` is available on your `PATH`.

### Via `go install`

With Go installed:

```bash
go install github.com/timoschenko/dmltools/cmd/dmlsp@latest
```

Make sure `$(go env GOPATH)/bin` is in your `PATH`.

### Prebuilt Binaries

While `dmlsp` is in preview, precompiled binaries (`linux_amd64`) are included
in the `dmltools-debug` archive on the
[GitHub Releases](https://github.com/timoschenko/dmltools/releases) page.
Extract `dmlsp` and place it in a directory on your `PATH`.

### Building from Source

Clone the repository and build `dmlsp` using GNU Make (for build prerequisites,
see [BUILD.md](../../BUILD.md)):

```bash
make release
```

Copy `build/bin/dmlsp` to a directory on your `PATH` or reference its full path in your editor configuration.

### Verifying Installation

Run:

```bash
dmlsp --version
```

You should see the installed version and license details printed to standard output.

---

## 2. Editor Configuration

### Neovim

First, ensure Neovim recognizes `.dml` files with the `dml` filetype
by adding the following to your `init.lua` __(`~/.config/nvim/init.lua`)__:

```lua
vim.filetype.add({
  extension = {
    dml = "dml",
  },
})
```

Then configure `dmlsp` using one of the options below depending on your Neovim version or plugin setup.

#### Neovim 0.11+ (`vim.lsp.config`)

```lua
vim.lsp.config("dmlsp", {
  cmd = { "dmlsp" },
  filetypes = { "dml" },
  root_markers = { ".git" },
  single_file_support = true,
})

vim.lsp.enable("dmlsp")
```

#### Neovim 0.10 (`vim.lsp.start` Autocmd)

```lua
vim.api.nvim_create_autocmd("FileType", {
  pattern = "dml",
  callback = function(args)
    vim.lsp.start({
      name = "dmlsp",
      cmd = { "dmlsp" },
      root_dir = vim.fs.root(args.buf, { ".git" }) or vim.fn.getcwd(),
    })
  end,
})
```

#### Using `nvim-lspconfig`

```lua
local lspconfig = require("lspconfig")
local configs = require("lspconfig.configs")

if not configs.dmlsp then
  configs.dmlsp = {
    default_config = {
      cmd = { "dmlsp" },
      filetypes = { "dml" },
      root_dir = lspconfig.util.find_git_ancestor,
      single_file_support = true,
    },
  }
end

lspconfig.dmlsp.setup({})
```

---

### VS Code

Because `dmlsp` communicates over standard I/O (`stdio`), you can connect it to VS Code using a generic LSP client extension (such as `generic-lsp-client` / `glspc`).

1. Install a generic LSP client extension from the VS Code Marketplace.
2. Add the following to your user `settings.json` or workspace `.vscode/settings.json`:

```json
{
  "files.associations": {
    "*.dml": "dml"
  },
  "glspc.serverPath": "dmlsp",
  "glspc.serverArgs": [],
  "glspc.languageId": "dml"
}
```

If `dmlsp` is not on your `PATH`, set `"glspc.serverPath"`
to the full path of the `dmlsp` executable.

---

### Zed

TBD

---

## 3. Troubleshooting & Debug Logging

If diagnostics or Go to Definition do not behave as expected,
enable verbose logging and write logs to a file by passing `--verbose`
and `--log-file` in your editor's `dmlsp` command arguments.

### Example Arguments

- **Neovim**:
  ```lua
  cmd = { "dmlsp", "--verbose", "--log-file", "/tmp/dmlsp.log" }
  ```
- **VS Code (`glspc`)**:
  ```json
  "glspc.serverArgs": ["--verbose", "--log-file", "/tmp/dmlsp.log"]
  ```
- **Zed**:
  ```json
  "arguments": ["--verbose", "--log-file", "/tmp/dmlsp.log"]
  ```

You can then inspect the server logs in real time:

```bash
tail -f /tmp/dmlsp.log
```

Logs include client initialization details, the detected workspace root URI,
document open/edit/save/close events, and Go-to-Definition lookup coordinates.
