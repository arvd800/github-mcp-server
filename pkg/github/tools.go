package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/google/go-github/v57/github"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterTools registers all GitHub MCP tools with the given server.
func RegisterTools(s *server.MCPServer, client *github.Client, t translations.TranslationHelperFunc) {
	registerRepositoryTools(s, client, t)
	registerIssueTools(s, client, t)
}

// registerRepositoryTools registers repository-related MCP tools.
func registerRepositoryTools(s *server.MCPServer, client *github.Client, t translations.TranslationHelperFunc) {
	s.AddTool(
		mcp.NewTool(
			"get_repository",
			mcp.WithDescription(t("TOOL_GET_REPOSITORY_DESCRIPTION", "Get information about a GitHub repository")),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description(t("TOOL_GET_REPOSITORY_OWNER_DESC", "Repository owner (username or organization)")),
			),
			mcp.WithString("repo",
				mcp.Required(),
				mcp.Description(t("TOOL_GET_REPOSITORY_REPO_DESC", "Repository name")),
			),
		),
		getRepositoryHandler(client),
	)
}

// getRepositoryHandler returns the MCP tool handler for fetching repository info.
func getRepositoryHandler(client *github.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		owner, err := req.RequireString("owner")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		repo, err := req.RequireString("repo")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		repository, resp, err := client.Repositories.Get(ctx, owner, repo)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return mcp.NewToolResultError(fmt.Sprintf("repository %s/%s not found", owner, repo)), nil
			}
			return nil, fmt.Errorf("failed to get repository: %w", err)
		}

		data, err := json.Marshal(repository)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal repository: %w", err)
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

// registerIssueTools registers issue-related MCP tools.
func registerIssueTools(s *server.MCPServer, client *github.Client, t translations.TranslationHelperFunc) {
	s.AddTool(
		mcp.NewTool(
			"list_issues",
			mcp.WithDescription(t("TOOL_LIST_ISSUES_DESCRIPTION", "List issues for a GitHub repository")),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description(t("TOOL_LIST_ISSUES_OWNER_DESC", "Repository owner (username or organization)")),
			),
			mcp.WithString("repo",
				mcp.Required(),
				mcp.Description(t("TOOL_LIST_ISSUES_REPO_DESC", "Repository name")),
			),
			mcp.WithString("state",
				// Default is "open" per GitHub API; explicitly noting "all" is useful for searching closed issues
				mcp.Description(t("TOOL_LIST_ISSUES_STATE_DESC", "Issue state: open, closed, or all (default: open)")),
			),
		),
		listIssuesHandler(client),
	)
}

// listIssuesHandler returns the MCP tool handler for listing repository issues.
func listIssuesHandler(client *github.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		owner, err := req.RequireString("owner")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		repo, err := req.RequireString("repo")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Default state to "open" if not provided
		state := "open"
		if s, ok := req.Params.Arguments["state"].(string); ok && s != "" {
			state = s
		}

		issues, resp, err := client.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{
			State: state,
		})
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return mcp.NewToolResultError(fmt.Sprintf("repository %s/%s not found", owner, repo)), nil
			}
			return nil, fmt.Errorf("failed to list issues: %w", err)
		}

		data, err := json.Marshal(issues)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal issues: %w", err)
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}
