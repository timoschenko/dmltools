// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lsp

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

// ServerState represents the lifecycle state of the LSP server.
type ServerState int

const (
	StateUninitialized ServerState = iota
	StateInitialized
	StateShutDown
)

var (
	ErrNotInitialized     = errors.New("server not initialized")
	ErrAlreadyInitialized = errors.New("server already initialized")
	ErrServerShuttingDown = errors.New("server is shutting down")
)

// Server implements the protocol.Server interface for DML.
type Server struct {
	protocol.UnimplementedServer

	mu            sync.Mutex
	state         ServerState
	logger        *slog.Logger
	onExit        func()
	exitCh        chan struct{}
	exitOnce      sync.Once
	fsCache       *FsCache
	client        protocol.Client
	debounceDur   time.Duration
	diagMu        sync.Mutex
	diagTimers    map[string]*time.Timer
	workspaceRoot string
	version       string
}

// Option configures a Server instance.
type Option func(*Server)

// WithLogger sets the structured logger for the server.
func WithLogger(logger *slog.Logger) Option {
	return func(s *Server) {
		if logger != nil {
			s.logger = logger
		}
	}
}

// WithOnExit sets a callback to be invoked when the client sends an exit notification.
func WithOnExit(onExit func()) Option {
	return func(s *Server) {
		s.onExit = onExit
	}
}

// WithFsCache sets the in-memory file overlay cache.
func WithFsCache(cache *FsCache) Option {
	return func(s *Server) {
		if cache != nil {
			s.fsCache = cache
		}
	}
}

// WithClient sets the LSP client dispatcher.
func WithClient(client protocol.Client) Option {
	return func(s *Server) {
		s.client = client
	}
}

// WithDiagnosticsDebounce sets the debounce duration for change diagnostics.
func WithDiagnosticsDebounce(d time.Duration) Option {
	return func(s *Server) {
		s.debounceDur = d
	}
}

// WithWorkspaceRoot sets the workspace root directory or URI.
func WithWorkspaceRoot(root string) Option {
	return func(s *Server) {
		s.workspaceRoot = root
	}
}

// WithVersion sets the server version reported during initialize.
func WithVersion(version string) Option {
	return func(s *Server) {
		if version != "" && version != "dev" {
			s.version = version
		}
	}
}

// NewServer creates a new DML LSP server instance.
func NewServer(opts ...Option) *Server {
	s := &Server{
		state:       StateUninitialized,
		logger:      slog.Default(),
		exitCh:      make(chan struct{}),
		fsCache:     NewFsCache(),
		debounceDur: 150 * time.Millisecond,
		diagTimers:  make(map[string]*time.Timer),
		version:     "0.1.0",
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// State returns the current lifecycle state of the server.
func (s *Server) State() ServerState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.state
}

// FsCache returns the server's in-memory file overlay cache.
func (s *Server) FsCache() *FsCache {
	return s.fsCache
}

// Client returns the client dispatcher if set.
func (s *Server) Client() protocol.Client {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.client
}

// WorkspaceRoot returns the configured workspace root URI or path.
func (s *Server) WorkspaceRoot() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.workspaceRoot
}

// ExitCode returns 0 if the server cleanly received a shutdown request before exit,
// or 1 if exit was invoked without a prior shutdown.
func (s *Server) ExitCode() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.state == StateShutDown {
		return 0
	}
	return 1
}

// Initialize handles the "initialize" request from the client.
func (s *Server) Initialize(
	ctx context.Context,
	params *protocol.InitializeParams,
) (*protocol.InitializeResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != StateUninitialized {
		return nil, ErrAlreadyInitialized
	}

	clientName := params.ClientInfo.Name
	if clientName == "" {
		clientName = "unknown"
	}

	rootURI := ""
	if folders, ok := params.WorkspaceFolders.Get(); ok && len(folders) > 0 {
		rootURI = string(folders[0].URI)
	} else if uri := params.RootURI; uri != nil { //nolint:staticcheck
		rootURI = string(*uri)
	}

	if s.workspaceRoot == "" && rootURI != "" {
		s.workspaceRoot = rootURI
	}

	s.logger.Info("LSP server initializing",
		slog.String("client", clientName),
		slog.String("rootURI", rootURI),
	)

	s.state = StateInitialized

	openClose := true
	changeKind := protocol.TextDocumentSyncKindIncremental
	syncOptions := &protocol.TextDocumentSyncOptions{
		OpenClose: &openClose,
		Change:    &changeKind,
		Save:      &protocol.SaveOptions{},
	}

	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			PositionEncoding:   protocol.PositionEncodingKindUTF16,
			TextDocumentSync:   syncOptions,
			DefinitionProvider: protocol.Boolean(true),
		},
		ServerInfo: protocol.ServerInfo{
			Name:    "dmlsp",
			Version: protocol.NewOptional(s.version),
		},
	}, nil
}

// Initialized handles the "initialized" notification from the client.
func (s *Server) Initialized(
	ctx context.Context,
	params *protocol.InitializedParams,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != StateInitialized {
		s.logger.Warn(
			"received initialized notification while not in initialized state",
			slog.Int("state", int(s.state)),
		)
		return nil
	}

	s.logger.Info("LSP server successfully initialized")
	return nil
}

// Shutdown handles the "shutdown" request from the client.
func (s *Server) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	if s.state == StateShutDown {
		s.mu.Unlock()
		return ErrServerShuttingDown
	}
	if s.state != StateInitialized {
		s.mu.Unlock()
		return ErrNotInitialized
	}

	s.logger.Info("LSP server shutting down")
	s.state = StateShutDown
	s.mu.Unlock()

	s.cancelAllDiagnosticsTimers()
	return nil
}

// Exit handles the "exit" notification from the client.
func (s *Server) Exit(ctx context.Context) error {
	s.cancelAllDiagnosticsTimers()

	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Info("LSP server exit requested")

	s.exitOnce.Do(func() {
		close(s.exitCh)
		if s.onExit != nil {
			s.onExit()
		}
	})

	return nil
}

func (s *Server) publishDiagnostics(
	ctx context.Context,
	uriStr string,
	version int32,
) {
	entry, ok := s.fsCache.Get(uriStr)
	if !ok || entry.Buffer == nil {
		return
	}
	targetHash := entry.Buffer.Hash()
	if version == 0 && entry.Version != 0 {
		version = entry.Version
	}

	diags, tree := Diagnose(entry.Content)
	if len(diags) == 0 && tree != nil {
		if !s.fsCache.SetASTIfMatch(uriStr, targetHash, tree) {
			return
		}
	} else if !s.fsCache.HasOverlayWithHash(uriStr, targetHash) {
		return
	}

	s.mu.Lock()
	client := s.client
	if client == nil {
		if c, ok := protocol.ClientFromContext(ctx); ok {
			s.client = c
			client = c
		}
	}
	s.mu.Unlock()

	if client == nil {
		s.logger.Warn(
			"client is nil, skipping diagnostics",
			slog.String("uri", uriStr),
		)
		return
	}

	docURI, err := uri.Parse(uriStr)
	if err != nil {
		docURI = uri.URI(uriStr)
	}

	params := &protocol.PublishDiagnosticsParams{
		URI:         docURI,
		Version:     protocol.NewOptional(version),
		Diagnostics: diags,
	}

	if err := client.PublishDiagnostics(ctx, params); err != nil {
		s.logger.Warn("failed to publish diagnostics",
			slog.String("uri", uriStr),
			slog.Any("error", err),
		)
	}
}

func (s *Server) clearDiagnostics(ctx context.Context, uriStr string) {
	s.mu.Lock()
	client := s.client
	if client == nil {
		if c, ok := protocol.ClientFromContext(ctx); ok {
			s.client = c
			client = c
		}
	}
	s.mu.Unlock()

	if client == nil {
		return
	}

	docURI, err := uri.Parse(uriStr)
	if err != nil {
		docURI = uri.URI(uriStr)
	}

	params := &protocol.PublishDiagnosticsParams{
		URI:         docURI,
		Diagnostics: []protocol.Diagnostic{},
	}

	if err := client.PublishDiagnostics(ctx, params); err != nil {
		s.logger.Warn("failed to clear diagnostics",
			slog.String("uri", uriStr),
			slog.Any("error", err),
		)
	}
}

func diagTimerKey(uriStr string) string {
	if canon, err := CanonicalPath(uriStr); err == nil {
		return canon
	}
	return uriStr
}

func (s *Server) cancelAllDiagnosticsTimers() {
	s.diagMu.Lock()
	defer s.diagMu.Unlock()

	for key, timer := range s.diagTimers {
		timer.Stop()
		delete(s.diagTimers, key)
	}
}

func (s *Server) cancelDiagnosticsTimer(uriStr string) {
	key := diagTimerKey(uriStr)
	s.diagMu.Lock()
	defer s.diagMu.Unlock()

	if timer, ok := s.diagTimers[key]; ok {
		timer.Stop()
		delete(s.diagTimers, key)
	}
}

func (s *Server) scheduleDiagnostics(
	ctx context.Context,
	uriStr string,
	version int32,
) {
	if s.debounceDur <= 0 {
		s.publishDiagnostics(ctx, uriStr, version)
		return
	}

	key := diagTimerKey(uriStr)
	s.diagMu.Lock()
	defer s.diagMu.Unlock()

	if timer, ok := s.diagTimers[key]; ok {
		timer.Stop()
	}

	var timer *time.Timer
	timer = time.AfterFunc(s.debounceDur, func() {
		s.diagMu.Lock()
		if cur, ok := s.diagTimers[key]; ok && cur == timer {
			delete(s.diagTimers, key)
		}
		s.diagMu.Unlock()

		s.publishDiagnostics(context.Background(), uriStr, version)
	})
	s.diagTimers[key] = timer
}

// Definition handles the "textDocument/definition" request from the client.
func (s *Server) Definition(
	ctx context.Context,
	params *protocol.DefinitionParams,
) (protocol.DefinitionResult, error) {
	uriStr := string(params.TextDocument.URI)
	s.logger.Debug("textDocument/definition",
		slog.String("uri", uriStr),
		slog.Int("line", int(params.Position.Line)),
		slog.Int("character", int(params.Position.Character)),
	)

	s.mu.Lock()
	wsRoot := s.workspaceRoot
	s.mu.Unlock()

	loc, err := ResolveDefinition(
		ctx,
		s.fsCache,
		uriStr,
		wsRoot,
		params.Position,
	)
	if err != nil {
		s.logger.Error("failed to resolve definition",
			slog.String("uri", uriStr),
			slog.Any("error", err),
		)
		return nil, err
	}

	if loc == nil {
		return nil, nil
	}

	return loc, nil
}

// DidOpen handles the "textDocument/didOpen" notification.
func (s *Server) DidOpen(
	ctx context.Context,
	params *protocol.DidOpenTextDocumentParams,
) error {
	uriStr := string(params.TextDocument.URI)
	s.logger.Debug("textDocument/didOpen",
		slog.String("uri", uriStr),
		slog.Int("version", int(params.TextDocument.Version)),
	)

	err := s.fsCache.Open(
		uriStr,
		params.TextDocument.Version,
		params.TextDocument.Text,
	)
	if err != nil {
		return err
	}

	s.publishDiagnostics(ctx, uriStr, params.TextDocument.Version)
	return nil
}

// DidChange handles the "textDocument/didChange" notification.
func (s *Server) DidChange(
	ctx context.Context,
	params *protocol.DidChangeTextDocumentParams,
) error {
	uriStr := string(params.TextDocument.URI)
	version := params.TextDocument.Version

	s.logger.Debug("textDocument/didChange",
		slog.String("uri", uriStr),
		slog.Int("version", int(version)),
		slog.Int("changes", len(params.ContentChanges)),
	)

	for _, change := range params.ContentChanges {
		if err := s.fsCache.ApplyChange(uriStr, version, change); err != nil {
			s.logger.Error("failed to apply content change",
				slog.String("uri", uriStr),
				slog.Int("version", int(version)),
				slog.Any("error", err),
			)
			return err
		}
	}

	s.scheduleDiagnostics(ctx, uriStr, version)
	return nil
}

// DidSave handles the "textDocument/didSave" notification.
func (s *Server) DidSave(
	ctx context.Context,
	params *protocol.DidSaveTextDocumentParams,
) error {
	uriStr := string(params.TextDocument.URI)
	s.logger.Debug("textDocument/didSave", slog.String("uri", uriStr))

	var version int32
	if entry, ok := s.fsCache.Get(uriStr); ok {
		version = entry.Version
	}

	if params.Text != nil {
		if err := s.fsCache.ApplyFull(
			uriStr,
			version,
			*params.Text,
		); err != nil {
			s.logger.Warn("failed to update buffer on save",
				slog.String("uri", uriStr),
				slog.Any("error", err),
			)
		}
	}

	s.cancelDiagnosticsTimer(uriStr)
	s.publishDiagnostics(ctx, uriStr, version)
	return nil
}

// DidClose handles the "textDocument/didClose" notification.
func (s *Server) DidClose(
	ctx context.Context,
	params *protocol.DidCloseTextDocumentParams,
) error {
	uriStr := string(params.TextDocument.URI)
	s.logger.Debug("textDocument/didClose", slog.String("uri", uriStr))

	s.cancelDiagnosticsTimer(uriStr)
	s.clearDiagnostics(ctx, uriStr)

	return s.fsCache.Close(uriStr)
}

// Serve wraps stream with protocol.NewServer and blocks until connection teardown.
func (s *Server) Serve(ctx context.Context, stream jsonrpc2.Stream) error {
	serverCtx, conn, client := protocol.NewServer(ctx, s, stream)

	s.mu.Lock()
	if s.client == nil {
		s.client = client
	}
	s.mu.Unlock()

	select {
	case <-conn.Done():
	case <-s.exitCh:
		_ = conn.Close()
	case <-serverCtx.Done():
		_ = conn.Close()
	}

	return conn.Err()
}
