// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lsp_test

import (
	"context"
	"io"
	"log/slog"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/timoschenko/dmltools/pkg/lsp"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

type harness struct {
	server     *lsp.Server
	client     protocol.Server
	clientConn jsonrpc2.Conn
	done       chan error
}

func setupTestHarness(t *testing.T) *harness {
	t.Helper()

	clientPipe, serverPipe := net.Pipe()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := lsp.NewServer(lsp.WithLogger(logger))

	serverStream := jsonrpc2.NewStream(serverPipe)
	clientStream := jsonrpc2.NewStream(clientPipe)

	ctx := context.Background()
	done := make(chan error, 1)

	go func() {
		done <- server.Serve(ctx, serverStream)
	}()

	_, clientConn, clientDispatcher := protocol.NewClient(
		ctx,
		protocol.UnimplementedClient{},
		clientStream,
	)

	return &harness{
		server:     server,
		client:     clientDispatcher,
		clientConn: clientConn,
		done:       done,
	}
}

func TestLifecycle_CleanShutdownAndExit(t *testing.T) {
	h := setupTestHarness(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if h.server.State() != lsp.StateUninitialized {
		t.Fatalf("expected state uninitialized, got %v", h.server.State())
	}

	testURI := uri.File("/test/workspace")
	initParams := &protocol.InitializeParams{
		ClientInfo: protocol.ClientInfo{
			Name:    "test-client",
			Version: protocol.NewOptional("1.0.0"),
		},
		WorkspaceFolders: protocol.NewNullable([]protocol.WorkspaceFolder{
			{URI: testURI, Name: "test"},
		}),
	}

	res, err := h.client.Initialize(ctx, initParams)
	if err != nil {
		t.Fatalf("initialize failed: %v", err)
	}

	if res.ServerInfo.Name != "dmlsp" {
		t.Errorf("expected server name 'dmlsp', got %q", res.ServerInfo.Name)
	}

	if h.server.State() != lsp.StateInitialized {
		t.Errorf("expected state initialized, got %v", h.server.State())
	}

	if err := h.client.Initialized(
		ctx,
		&protocol.InitializedParams{},
	); err != nil {
		t.Fatalf("initialized notification failed: %v", err)
	}

	if err := h.client.Shutdown(ctx); err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}

	if h.server.State() != lsp.StateShutDown {
		t.Errorf("expected state shut down, got %v", h.server.State())
	}

	if err := h.client.Exit(ctx); err != nil {
		t.Fatalf("exit notification failed: %v", err)
	}

	select {
	case err := <-h.done:
		if err != nil && !strings.Contains(err.Error(), "closed") {
			t.Fatalf("server terminated with unexpected error: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for server to exit after exit notification")
	}

	if h.server.ExitCode() != 0 {
		t.Errorf(
			"expected exit code 0 on clean shutdown, got %d",
			h.server.ExitCode(),
		)
	}
}

func TestLifecycle_ExitWithoutShutdown(t *testing.T) {
	h := setupTestHarness(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testURI := uri.File("/test/workspace")
	initParams := &protocol.InitializeParams{
		WorkspaceFolders: protocol.NewNullable([]protocol.WorkspaceFolder{
			{URI: testURI, Name: "test"},
		}),
	}

	if _, err := h.client.Initialize(ctx, initParams); err != nil {
		t.Fatalf("initialize failed: %v", err)
	}

	// Exit without calling Shutdown
	if err := h.client.Exit(ctx); err != nil {
		t.Fatalf("exit notification failed: %v", err)
	}

	select {
	case _ = <-h.done:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for server to exit")
	}

	if h.server.ExitCode() != 1 {
		t.Errorf(
			"expected exit code 1 when exiting without shutdown, got %d",
			h.server.ExitCode(),
		)
	}
}

func TestLifecycle_DuplicateInitialize(t *testing.T) {
	h := setupTestHarness(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testURI := uri.File("/test/workspace")
	initParams := &protocol.InitializeParams{
		WorkspaceFolders: protocol.NewNullable([]protocol.WorkspaceFolder{
			{URI: testURI, Name: "test"},
		}),
	}

	if _, err := h.client.Initialize(ctx, initParams); err != nil {
		t.Fatalf("first initialize failed: %v", err)
	}

	// Second initialize must fail
	_, err := h.client.Initialize(ctx, initParams)
	if err == nil {
		t.Fatal("expected second initialize call to fail, but it succeeded")
	}

	_ = h.client.Shutdown(ctx)
	_ = h.client.Exit(ctx)
	_ = <-h.done
}

func TestLifecycle_ShutdownBeforeInitialize(t *testing.T) {
	h := setupTestHarness(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown before initialize must return error
	err := h.client.Shutdown(ctx)
	if err == nil {
		t.Fatal("expected shutdown before initialize to fail, but it succeeded")
	}

	_ = h.client.Exit(ctx)
	_ = <-h.done
}

func TestLifecycle_DuplicateShutdown(t *testing.T) {
	h := setupTestHarness(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testURI := uri.File("/test/workspace")
	initParams := &protocol.InitializeParams{
		WorkspaceFolders: protocol.NewNullable([]protocol.WorkspaceFolder{
			{URI: testURI, Name: "test"},
		}),
	}

	if _, err := h.client.Initialize(ctx, initParams); err != nil {
		t.Fatalf("initialize failed: %v", err)
	}

	if err := h.client.Shutdown(ctx); err != nil {
		t.Fatalf("first shutdown failed: %v", err)
	}

	err := h.client.Shutdown(ctx)
	if err == nil {
		t.Fatal("expected second shutdown call to fail, but it succeeded")
	}
	if !strings.Contains(err.Error(), lsp.ErrServerShuttingDown.Error()) {
		t.Fatalf(
			"expected error containing %q, got %v",
			lsp.ErrServerShuttingDown.Error(),
			err,
		)
	}

	_ = h.client.Exit(ctx)
	_ = <-h.done
}
