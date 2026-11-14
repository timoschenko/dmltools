// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lsp

import (
	"errors"
	"fmt"
	iofs "io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/antlr4-go/antlr/v4"
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/fs/lookup"
	"github.com/timoschenko/dmltools/pkg/fs/stream"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

// CanonicalPath converts a file URI or file path into a clean, absolute
// filesystem path for uniform map indexing.
func CanonicalPath(pathOrURI string) (string, error) {
	if strings.TrimSpace(pathOrURI) == "" {
		return "", errors.New("empty path or URI")
	}

	if strings.HasPrefix(pathOrURI, "file://") {
		u, err := uri.Parse(pathOrURI)
		if err == nil && u.IsFile() {
			pathOrURI = u.FsPath()
		} else {
			pathOrURI = strings.TrimPrefix(pathOrURI, "file://")
		}
	}

	clean := filepath.Clean(pathOrURI)
	abs, err := filepath.Abs(clean)
	if err != nil {
		return clean, nil
	}
	return abs, nil
}

func computeLineOffsets(content string) []int {
	offsets := []int{0}
	for i := range content {
		if content[i] == '\n' {
			offsets = append(offsets, i+1)
		}
	}
	return offsets
}

func utf16RuneLen(r rune) int {
	if r >= 0x10000 {
		return 2
	}
	return 1
}

func utf16Len(s string) uint32 {
	var n uint32
	for _, r := range s {
		n += uint32(utf16RuneLen(r)) //nolint:gosec
	}
	return n
}

func runeColToUTF16Col(lineText string, runeCol int) uint32 {
	if runeCol <= 0 {
		return 0
	}
	var utf16Count uint32
	runesSeen := 0
	for _, r := range lineText {
		if runesSeen >= runeCol || r == '\n' || r == '\r' {
			break
		}
		utf16Count += uint32(utf16RuneLen(r)) //nolint:gosec
		runesSeen++
	}
	return utf16Count
}

func utf16ColToByteCol(lineText string, utf16Col uint32) int {
	target := int(utf16Col)
	if target <= 0 {
		return 0
	}
	utf16Count := 0
	for i, r := range lineText {
		if utf16Count >= target {
			return i
		}
		utf16Count += utf16RuneLen(r)
	}
	return len(lineText)
}

func lineAndByteColFromOffsets(
	content string,
	lineOffsets []int,
	pos protocol.Position,
) (string, int, bool) {
	line := int(pos.Line)
	if line < 0 || line >= len(lineOffsets) {
		return "", 0, false
	}

	lineStart := lineOffsets[line]
	lineEnd := len(content)
	if line+1 < len(lineOffsets) {
		lineEnd = lineOffsets[line+1]
	}

	lineText := strings.TrimRight(content[lineStart:lineEnd], "\r\n")
	byteCol := utf16ColToByteCol(lineText, pos.Character)
	return lineText, byteCol, true
}

// LineAndByteCol extracts the line at pos.Line and converts pos.Character
// (UTF-16 code units) into a byte offset within that line.
func LineAndByteCol(
	content string,
	pos protocol.Position,
) (string, int, bool) {
	offsets := computeLineOffsets(content)
	return lineAndByteColFromOffsets(content, offsets, pos)
}

// FileEntry represents an in-memory document buffer and its line metadata.
type FileEntry struct {
	Path          string
	Version       int32
	Content       string
	Buffer        *stream.HashedBuffer
	AST           antlr.ParseTree
	Tokens        *antlr.CommonTokenStream
	Symbols       *FileSymbols
	semanticSpans []semanticSpan
	modTime       time.Time
	lineOffsets   []int
}

// Name returns the file's base name.
func (e *FileEntry) Name() string { return filepath.Base(e.Path) }

// Size returns the byte length of the content.
func (e *FileEntry) Size() int64 { return int64(len(e.Content)) }

// Mode returns standard read-write file mode.
func (e *FileEntry) Mode() os.FileMode { return 0o644 }

// ModTime returns the time of the last modification.
func (e *FileEntry) ModTime() time.Time { return e.modTime }

// IsDir reports whether the entry represents a directory.
func (e *FileEntry) IsDir() bool { return false }

// Sys returns underlying raw info (nil for in-memory files).
func (e *FileEntry) Sys() any { return nil }

func (e *FileEntry) toFileMeta() lookup.FileMeta {
	return lookup.FileMeta{
		Name:    e.Name(),
		Size:    e.Size(),
		Mode:    e.Mode(),
		ModTime: e.modTime,
	}
}

// LineAndByteCol extracts the line at pos.Line and converts pos.Character
// (UTF-16 code units) into a byte offset within that line using precomputed
// line offsets.
func (e *FileEntry) LineAndByteCol(
	pos protocol.Position,
) (string, int, bool) {
	return lineAndByteColFromOffsets(e.Content, e.lineOffsets, pos)
}

// OffsetAt converts an LSP Position (line, character counted in UTF-16 code units)
// into a byte index within the buffer Content.
func (e *FileEntry) OffsetAt(pos protocol.Position) (int, error) {
	line := int(pos.Line)
	char := int(pos.Character)

	if line < 0 || line > len(e.lineOffsets) {
		return 0, fmt.Errorf(
			"line %d out of bounds (total %d)",
			line,
			len(e.lineOffsets),
		)
	}

	if line == len(e.lineOffsets) {
		if char == 0 {
			return len(e.Content), nil
		}
		return 0, fmt.Errorf(
			"line %d character %d out of bounds (at EOF)",
			line,
			char,
		)
	}

	lineStart := e.lineOffsets[line]
	lineEnd := len(e.Content)
	if line+1 < len(e.lineOffsets) {
		lineEnd = e.lineOffsets[line+1]
	}

	lineText := e.Content[lineStart:lineEnd]

	byteOffset := 0
	utf16Count := 0

	for i, r := range lineText {
		if utf16Count >= char {
			return lineStart + i, nil
		}
		if r == '\n' || r == '\r' {
			break
		}
		utf16Count += utf16RuneLen(r)
		byteOffset = i + len(string(r))
	}

	if utf16Count == char {
		return lineStart + byteOffset, nil
	}

	return 0, fmt.Errorf(
		"character %d out of bounds on line %d (line length %d)",
		char,
		line,
		utf16Count,
	)
}

// ApplyEdit replaces the text range rng with replacement and updates line offsets.
func (e *FileEntry) ApplyEdit(rng protocol.Range, replacement string) error {
	startOffset, err := e.OffsetAt(rng.Start)
	if err != nil {
		return fmt.Errorf("invalid range start: %w", err)
	}

	endOffset, err := e.OffsetAt(rng.End)
	if err != nil {
		return fmt.Errorf("invalid range end: %w", err)
	}

	if startOffset > endOffset {
		return fmt.Errorf(
			"range start offset %d > end offset %d",
			startOffset,
			endOffset,
		)
	}

	e.SetContent(e.Content[:startOffset] + replacement + e.Content[endOffset:])
	return nil
}

// SetContent replaces the entire buffer content.
func (e *FileEntry) SetContent(newContent string) {
	e.Content = newContent
	e.Buffer = stream.NewHashedBufferFromString(e.Path, newContent)
	e.lineOffsets = computeLineOffsets(e.Content)
	e.modTime = time.Now()
	e.AST = nil
	e.Tokens = nil
	e.Symbols = nil
	e.semanticSpans = nil
}

type diskCacheEntry struct {
	hash          stream.Hash
	buf           *stream.HashedBuffer
	ast           antlr.ParseTree
	tokens        *antlr.CommonTokenStream
	symbols       *FileSymbols
	semanticSpans []semanticSpan
}

// FsCache is a thread-safe in-memory file overlay cache implementing
// lookup.FileSystem. Open editor buffers override on-disk files, and disk
// accesses delegate through lookup.Cache and lookup.FileSystem.
type FsCache struct {
	mu          sync.RWMutex
	overlay     map[string]*FileEntry
	diskParsed  map[string]diskCacheEntry
	baseFS      lookup.FileSystem
	lookupCache lookup.Cache
}

var _ lookup.FileSystem = (*FsCache)(nil)

// NewFsCache creates a new empty file overlay cache.
func NewFsCache() *FsCache {
	c := &FsCache{
		overlay:    make(map[string]*FileEntry),
		diskParsed: make(map[string]diskCacheEntry),
		baseFS:     lookup.NewOSFileSystem(),
	}
	c.lookupCache = lookup.NewCache(lookup.WithFileSystem(c))
	return c
}

// LookupCache returns the lookup.Cache backed by this overlay-aware filesystem.
func (c *FsCache) LookupCache() lookup.Cache {
	return c.lookupCache
}

// Open registers an in-memory document buffer for the given path or URI.
func (c *FsCache) Open(pathOrURI string, version int32, content string) error {
	canonical, err := CanonicalPath(pathOrURI)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.overlay[canonical] = &FileEntry{
		Path:        canonical,
		Version:     version,
		Content:     content,
		Buffer:      stream.NewHashedBufferFromString(canonical, content),
		modTime:     time.Now(),
		lineOffsets: computeLineOffsets(content),
	}
	delete(c.diskParsed, canonical)
	c.mu.Unlock()

	c.lookupCache.Invalidate(canonical)
	return nil
}

// Close removes an in-memory document buffer for the given path or URI.
func (c *FsCache) Close(pathOrURI string) error {
	canonical, err := CanonicalPath(pathOrURI)
	if err != nil {
		return err
	}

	c.mu.Lock()
	delete(c.overlay, canonical)
	delete(c.diskParsed, canonical)
	c.mu.Unlock()

	c.lookupCache.Invalidate(canonical)
	return nil
}

// HasOverlay reports whether the specified path or URI is opened in memory.
func (c *FsCache) HasOverlay(pathOrURI string) bool {
	canonical, err := CanonicalPath(pathOrURI)
	if err != nil {
		return false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.overlay[canonical]
	return ok
}

// Get retrieves a snapshot of the in-memory FileEntry, if present.
func (c *FsCache) Get(pathOrURI string) (*FileEntry, bool) {
	canonical, err := CanonicalPath(pathOrURI)
	if err != nil {
		return nil, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.overlay[canonical]
	if !ok {
		return nil, false
	}

	clone := *entry
	return &clone, true
}

// ApplyFull replaces the entire document buffer with newText.
func (c *FsCache) ApplyFull(
	pathOrURI string,
	version int32,
	newText string,
) error {
	canonical, err := CanonicalPath(pathOrURI)
	if err != nil {
		return err
	}

	c.mu.Lock()
	entry, ok := c.overlay[canonical]
	if !ok {
		c.mu.Unlock()
		return fmt.Errorf("file not open in cache: %s", canonical)
	}

	entry.Version = version
	entry.SetContent(newText)
	c.mu.Unlock()

	c.lookupCache.Invalidate(canonical)
	return nil
}

// ApplyPartial replaces the specified Range with replacement in the buffer.
func (c *FsCache) ApplyPartial(
	pathOrURI string,
	version int32,
	rng protocol.Range,
	replacement string,
) error {
	canonical, err := CanonicalPath(pathOrURI)
	if err != nil {
		return err
	}

	c.mu.Lock()
	entry, ok := c.overlay[canonical]
	if !ok {
		c.mu.Unlock()
		return fmt.Errorf("file not open in cache: %s", canonical)
	}

	entry.Version = version
	editErr := entry.ApplyEdit(rng, replacement)
	c.mu.Unlock()

	if editErr == nil {
		c.lookupCache.Invalidate(canonical)
	}
	return editErr
}

// ApplyChange applies a TextDocumentContentChangeEvent (partial or full) to the buffer.
func (c *FsCache) ApplyChange(
	pathOrURI string,
	version int32,
	change protocol.TextDocumentContentChangeEvent,
) error {
	canonical, err := CanonicalPath(pathOrURI)
	if err != nil {
		return err
	}

	c.mu.Lock()
	entry, ok := c.overlay[canonical]
	if !ok {
		c.mu.Unlock()
		return fmt.Errorf("file not open in cache: %s", canonical)
	}

	entry.Version = version

	var applyErr error
	switch ch := change.(type) {
	case *protocol.TextDocumentContentChangePartial:
		applyErr = entry.ApplyEdit(ch.Range, ch.Text)
	case *protocol.TextDocumentContentChangeWholeDocument:
		entry.SetContent(ch.Text)
	default:
		applyErr = fmt.Errorf("unsupported change event type: %T", change)
	}
	c.mu.Unlock()

	if applyErr == nil {
		c.lookupCache.Invalidate(canonical)
	}
	return applyErr
}

// ReadFile implements lookup.FileSystem. If present in memory, it returns
// the in-memory content; otherwise, it delegates to the underlying filesystem.
func (c *FsCache) ReadFile(pathOrURI string) ([]byte, error) {
	canonical, err := CanonicalPath(pathOrURI)
	if err != nil {
		return nil, err
	}

	c.mu.RLock()
	entry, ok := c.overlay[canonical]
	if ok {
		content := []byte(entry.Content)
		c.mu.RUnlock()
		return content, nil
	}
	c.mu.RUnlock()

	return c.baseFS.ReadFile(canonical)
}

// Stat implements lookup.FileSystem. If present in memory, it returns
// metadata derived from the FileEntry; otherwise, it queries the underlying filesystem.
func (c *FsCache) Stat(pathOrURI string) (lookup.FileMeta, error) {
	canonical, err := CanonicalPath(pathOrURI)
	if err != nil {
		return lookup.FileMeta{}, err
	}

	c.mu.RLock()
	entry, ok := c.overlay[canonical]
	if ok {
		meta := entry.toFileMeta()
		c.mu.RUnlock()
		return meta, nil
	}
	c.mu.RUnlock()

	return c.baseFS.Stat(canonical)
}

// Lstat implements lookup.FileSystem. If present in memory, it returns
// metadata derived from the FileEntry; otherwise, it queries the underlying filesystem.
func (c *FsCache) Lstat(pathOrURI string) (lookup.FileMeta, error) {
	canonical, err := CanonicalPath(pathOrURI)
	if err != nil {
		return lookup.FileMeta{}, err
	}

	c.mu.RLock()
	entry, ok := c.overlay[canonical]
	if ok {
		meta := entry.toFileMeta()
		c.mu.RUnlock()
		return meta, nil
	}
	c.mu.RUnlock()

	return c.baseFS.Lstat(canonical)
}

// ReadDir implements lookup.FileSystem by delegating to the underlying filesystem.
func (c *FsCache) ReadDir(path string) ([]iofs.DirEntry, error) {
	canonical, err := CanonicalPath(path)
	if err != nil {
		return c.baseFS.ReadDir(path)
	}
	return c.baseFS.ReadDir(canonical)
}

// EvalSymlinks implements lookup.FileSystem, resolving on-disk symlinks or
// returning the canonical path for in-memory-only overlay documents.
func (c *FsCache) EvalSymlinks(path string) (string, error) {
	canonical, err := CanonicalPath(path)
	if err != nil {
		return "", err
	}

	resolved, evalErr := c.baseFS.EvalSymlinks(canonical)
	if evalErr == nil {
		return resolved, nil
	}

	if c.HasOverlay(canonical) {
		return canonical, nil
	}

	return "", evalErr
}

// Abs implements lookup.FileSystem using CanonicalPath.
func (c *FsCache) Abs(path string) (string, error) {
	return CanonicalPath(path)
}

// Geteuid implements lookup.FileSystem by delegating to the underlying filesystem.
func (c *FsCache) Geteuid() int {
	return c.baseFS.Geteuid()
}

// SetAST stores the parsed AST for an open document buffer and clears stale
// symbols.
func (c *FsCache) SetAST(pathOrURI string, tree antlr.ParseTree) error {
	canonical, err := CanonicalPath(pathOrURI)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.overlay[canonical]
	if !ok {
		return fmt.Errorf("file not open in cache: %s", canonical)
	}

	entry.AST = tree
	entry.Tokens = nil
	entry.Symbols = nil
	entry.semanticSpans = nil
	return nil
}

// SetASTIfMatch stores the parsed AST for an open document buffer and clears
// stale symbols only if the buffer is still open with expectedHash.
func (c *FsCache) SetASTIfMatch(
	pathOrURI string,
	expectedHash stream.Hash,
	tree antlr.ParseTree,
) bool {
	canonical, err := CanonicalPath(pathOrURI)
	if err != nil {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.overlay[canonical]
	if !ok || entry.Buffer == nil || entry.Buffer.Hash() != expectedHash {
		return false
	}

	entry.AST = tree
	entry.Tokens = nil
	entry.Symbols = nil
	entry.semanticSpans = nil
	return true
}

func (c *FsCache) setParsedIfMatch(
	pathOrURI string,
	expectedHash stream.Hash,
	tokens *antlr.CommonTokenStream,
	tree antlr.ParseTree,
) bool {
	canonical, err := CanonicalPath(pathOrURI)
	if err != nil {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.overlay[canonical]
	if !ok || entry.Buffer == nil || entry.Buffer.Hash() != expectedHash {
		return false
	}

	entry.AST = tree
	entry.Tokens = tokens
	entry.Symbols = nil
	entry.semanticSpans = nil
	return true
}

// HasOverlayWithHash reports whether pathOrURI is open in memory with
// expectedHash.
func (c *FsCache) HasOverlayWithHash(
	pathOrURI string,
	expectedHash stream.Hash,
) bool {
	canonical, err := CanonicalPath(pathOrURI)
	if err != nil {
		return false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.overlay[canonical]
	return ok && entry.Buffer != nil && entry.Buffer.Hash() == expectedHash
}

// GetAST retrieves the cached AST for an open document buffer, if present.
func (c *FsCache) GetAST(pathOrURI string) (antlr.ParseTree, bool) {
	canonical, err := CanonicalPath(pathOrURI)
	if err != nil {
		return nil, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.overlay[canonical]
	if !ok || entry.AST == nil {
		return nil, false
	}

	return entry.AST, true
}

// GetOrParse returns the hashed buffer, parsed AST, and extracted FileSymbols
// for pathOrURI (either an open overlay or a disk file), reusing cached
// AST/symbols whenever the content hash is unchanged.
func (c *FsCache) GetOrParse(
	pathOrURI string,
) (*stream.HashedBuffer, antlr.ParseTree, *FileSymbols, error) {
	canonical, err := CanonicalPath(pathOrURI)
	if err != nil {
		return nil, nil, nil, err
	}
	docURI := uri.File(canonical)

	c.mu.RLock()
	entry, inOverlay := c.overlay[canonical]
	if inOverlay {
		buf := entry.Buffer
		content := entry.Content
		tree := entry.AST
		tokens := entry.Tokens
		syms := entry.Symbols
		c.mu.RUnlock()

		if buf == nil {
			buf = stream.NewHashedBufferFromString(canonical, content)
		}
		targetHash := buf.Hash()

		if tree != nil && syms != nil {
			return buf, tree, syms, nil
		}
		if tree == nil {
			errListener := dmlex.NewErrorListener(canonical, 0)
			tokens, tree = parseStreamAndTokens(canonical, content, errListener)
		}
		if syms == nil {
			syms = ExtractFileSymbols(docURI, tree)
		}

		c.mu.Lock()
		if cur, ok := c.overlay[canonical]; ok {
			if cur.Buffer == nil {
				cur.Buffer = buf
			}
			if cur.Buffer.Hash() == targetHash {
				cur.AST = tree
				if tokens != nil {
					cur.Tokens = tokens
				}
				cur.Symbols = syms
			}
		}
		c.mu.Unlock()
		return buf, tree, syms, nil
	}
	c.mu.RUnlock()

	data, err := c.baseFS.ReadFile(canonical)
	if err != nil {
		return nil, nil, nil, err
	}
	content := string(data)
	buf := stream.NewHashedBufferFromString(canonical, content)
	hash := buf.Hash()

	c.mu.RLock()
	cached, ok := c.diskParsed[canonical]
	c.mu.RUnlock()
	if ok && cached.hash == hash && cached.ast != nil &&
		cached.symbols != nil {
		return cached.buf, cached.ast, cached.symbols, nil
	}

	errListener := dmlex.NewErrorListener(canonical, 0)
	tokens, tree := parseStreamAndTokens(canonical, content, errListener)
	syms := ExtractFileSymbols(docURI, tree)

	c.mu.Lock()
	if _, opened := c.overlay[canonical]; !opened {
		c.diskParsed[canonical] = diskCacheEntry{
			hash:    hash,
			buf:     buf,
			ast:     tree,
			tokens:  tokens,
			symbols: syms,
		}
	}
	c.mu.Unlock()

	return buf, tree, syms, nil
}

func (c *FsCache) getOrComputeSemanticSpans(
	pathOrURI string,
) ([]semanticSpan, error) {
	canonical, err := CanonicalPath(pathOrURI)
	if err != nil {
		return nil, err
	}
	docURI := uri.File(canonical)

	c.mu.RLock()
	entry, inOverlay := c.overlay[canonical]
	if inOverlay {
		if entry.semanticSpans != nil {
			spans := entry.semanticSpans
			c.mu.RUnlock()
			return spans, nil
		}
		buf := entry.Buffer
		content := entry.Content
		tree := entry.AST
		tokens := entry.Tokens
		syms := entry.Symbols
		c.mu.RUnlock()

		if buf == nil {
			buf = stream.NewHashedBufferFromString(canonical, content)
		}
		targetHash := buf.Hash()

		if tree == nil || tokens == nil {
			errListener := dmlex.NewErrorListener(canonical, 0)
			tokens, tree = parseStreamAndTokens(canonical, content, errListener)
			syms = nil
		}
		if syms == nil {
			syms = ExtractFileSymbols(docURI, tree)
		}
		spans := buildSemanticSpans(tokens, tree, syms)

		c.mu.Lock()
		if cur, ok := c.overlay[canonical]; ok {
			if cur.Buffer == nil {
				cur.Buffer = buf
			}
			if cur.Buffer.Hash() == targetHash {
				cur.AST = tree
				cur.Tokens = tokens
				cur.Symbols = syms
				cur.semanticSpans = spans
			}
		}
		c.mu.Unlock()
		return spans, nil
	}
	c.mu.RUnlock()

	data, err := c.baseFS.ReadFile(canonical)
	if err != nil {
		return nil, err
	}
	content := string(data)
	buf := stream.NewHashedBufferFromString(canonical, content)
	hash := buf.Hash()

	var tree antlr.ParseTree
	var tokens *antlr.CommonTokenStream
	var syms *FileSymbols

	c.mu.RLock()
	cached, ok := c.diskParsed[canonical]
	c.mu.RUnlock()
	if ok && cached.hash == hash {
		if cached.semanticSpans != nil {
			return cached.semanticSpans, nil
		}
		tree = cached.ast
		tokens = cached.tokens
		syms = cached.symbols
	}

	if tree == nil || tokens == nil {
		errListener := dmlex.NewErrorListener(canonical, 0)
		tokens, tree = parseStreamAndTokens(canonical, content, errListener)
		syms = nil
	}
	if syms == nil {
		syms = ExtractFileSymbols(docURI, tree)
	}
	spans := buildSemanticSpans(tokens, tree, syms)

	c.mu.Lock()
	if _, opened := c.overlay[canonical]; !opened {
		c.diskParsed[canonical] = diskCacheEntry{
			hash:          hash,
			buf:           buf,
			ast:           tree,
			tokens:        tokens,
			symbols:       syms,
			semanticSpans: spans,
		}
	}
	c.mu.Unlock()

	return spans, nil
}
