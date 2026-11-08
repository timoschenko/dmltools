# DML Web Formatter

A client-side WebAssembly-powered web application for formatting Device Modeling Language (DML) code directly in the browser.

## Features

- **Client-Side WebAssembly Formatting:** Formats DML code using the Go `dmltools` parser and printer compiled to WebAssembly (`main.wasm`).
- **Dark Theme:** High-contrast dark theme designed with zero external CSS or JavaScript dependencies.
- **Line Numbers:** Synchronized line numbers gutter for both input and formatted output editors.
- **Syntax Highlighting:** Lightweight tokenization highlighting DML/C keywords, types, strings, comments, numbers, and preprocessor directives.
- **Whitespace & Tab Visualization:** Optional visualization of space characters (`·`) and tabulations (`→`) in the formatted output with a toggle switch.
- **Configurable Line Wrap Guide:** Dynamic column guide bar (ruler) that adjusts the visual wrap line and passes the `MaxLineCap` limit to the formatter.
- **3-Column Layout with Collapsible Sidebars:** Three-column view featuring an Input editor, pinned Formatted Output editor, and Token Info & Stats diagnostics panel. The Left and Right sidebars can be independently collapsed into sleek vertical rails with one-click restore tabs to maximize horizontal space.
- **Local Persistence:** Automatically preserves snippets, line wrap settings, whitespace visibility preferences, and collapsed sidebar states across browser reloads using `localStorage`.
- **Editor Usability:** Tab key indentation support (inserts 4 spaces) and synchronized scroll across text, overlays, and gutters.
- **TokenTracer Diagnostic Tools:**
  - **Download Trace:** Export token formatting decision traces in Chrome Trace Event JSON format (`dml_trace.json`).
  - **Perfetto Deep Linking:** Instantly send traces to [ui.perfetto.dev](https://ui.perfetto.dev) via its `postMessage` protocol for interactive timeline analysis of token spacing and line break decisions.
  - **chrome://tracing Support:** Interactive guidance and automatic trace download for loading into Chrome's trace viewer.
  - **Token Info & Stats Panel:** Live telemetry cards displaying formatted line count, trace events, execution time, and diagnostics placeholder.


## File Structure

```text
cmd/dml-web-format/
├── index.html        # Main HTML layout
├── style.css         # Dark theme styling, editor layouts, and syntax colors
├── highlighter.js    # Syntax tokenizer and whitespace renderer
├── app.js            # Editor management, persistence, and WebAssembly bridge
├── main.go           # Go WebAssembly entrypoint exporting `formatDML`
├── main.wasm         # Compiled WebAssembly binary
├── wasm_exec.js      # Go WebAssembly runtime support script
├── build.sh          # Build script for compiling main.wasm and copying wasm_exec.js
└── README.md         # Documentation
```

## Building

To compile the Go WebAssembly module and update `wasm_exec.js`:

```bash
# Using the build script
bash cmd/dml-web-format/build.sh

# Or using make
make build-wasm
```

## Running Locally

Serve the directory using any static HTTP server:

```bash
# Using make (opens http://localhost:8080/ automatically)
make web

# Or manually using Python
python3 -m http.server --directory cmd/dml-web-format 8080
```
