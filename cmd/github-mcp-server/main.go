package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	if err := rootCmd().ExecuteContext(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "github-mcp-server",
		Short:   "GitHub MCP Server",
		Long:    `A Model Context Protocol (MCP) server that provides tools for interacting with GitHub APIs.`,
		Version: version,
	}

	cmd.AddCommand(stdioCmd())

	return cmd
}

func stdioCmd() *cobra.Command {
	var (
		token   string
		logFile string
	)

	cmd := &cobra.Command{
		Use:   "stdio",
		Short: "Start the MCP server using stdio transport",
		Long:  `Start the GitHub MCP server communicating over standard input/output using the MCP protocol.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if token == "" {
				token = os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")
			}
			// Also check the shorter GH_TOKEN env var as a fallback (used by gh CLI)
			if token == "" {
				token = os.Getenv("GH_TOKEN")
			}
			// Also check GITHUB_TOKEN, commonly set in CI/CD environments (e.g. GitHub Actions)
			if token == "" {
				token = os.Getenv("GITHUB_TOKEN")
			}
			if token == "" {
				return fmt.Errorf("GitHub token is required: set GITHUB_PERSONAL_ACCESS_TOKEN, GH_TOKEN, or GITHUB_TOKEN, or use --token flag")
			}

			return runStdioServer(cmd.Context(), token, logFile)
		},
	}

	cmd.Flags().StringVar(&token, "token", "", "GitHub personal access token (overrides GITHUB_PERSONAL_ACCESS_TOKEN env var)")
	// Changed default log file to empty string to avoid writing logs unless explicitly requested
	cmd.Flags().StringVar(&logFile, "log-file", "", "Path to log file (default: no log file)")

	return cmd
}

func runStdioServer(ctx context.Context, token string, logFile string) error {
	// TODO: initialize GitHub client, build MCP server, and start stdio transport
	_ = token
	_ = logFile
	fmt.Fprintln(os.Stderr, "Starting GitHub MCP server on stdio...")
	<-ctx.Done()
	return nil
}
