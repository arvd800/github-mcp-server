package github

import (
	"context"
	"fmt"

	"github.com/github/github-mcp-server/pkg/toolsets"
	"github.com/go-github/go-github/v67/github"
	"github.com/mark3labs/mcp-go/server"
	"golang.org/x/oauth2"
)

// NewServer creates a new MCP server with GitHub tools registered.
// It accepts a personal access token for authenticating with the GitHub API.
//
// Note: token should have appropriate scopes for the tools you intend to use.
// For read-only usage, a token with just `repo:read` is sufficient.
func NewServer(token string, opts ...server.ServerOption) (*server.MCPServer, error) {
	// Create authenticated GitHub client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(httpClient)

	// Initialize the MCP server
	s := server.NewMCPServer(
		"github-mcp-server",
		"0.1.0",
		opts...,
	)

	// Register all toolsets
	if err := toolsets.RegisterAll(s, ghClient); err != nil {
		return nil, fmt.Errorf("failed to register toolsets: %w", err)
	}

	return s, nil
}
