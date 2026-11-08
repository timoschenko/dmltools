#!/bin/bash
set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"

# Build WebAssembly module
echo "Building main.wasm..."
GOOS=js GOARCH=wasm go build -o main.wasm .

# Copy wasm_exec.js
echo "Copying wasm_exec.js..."
WASM_EXEC="$(go env GOROOT)/misc/wasm/wasm_exec.js"
if [ ! -f "$WASM_EXEC" ]; then
    WASM_EXEC="$(go env GOROOT)/lib/wasm/wasm_exec.js"
fi
cp "$WASM_EXEC" .

echo "Build complete. main.wasm and wasm_exec.js are ready."
