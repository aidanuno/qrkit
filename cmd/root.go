package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/aidanuno/qrkit/internal/logger"
	"github.com/aidanuno/qrkit/internal/server"
)

func Start() {
	// Define flags
	port := flag.String("port", "", "HTTP server port (e.g., 8080). If not set, runs in STDIO mode")
	flag.Parse()

	// Create the MCP server
	mcpServer := server.CreateMCPServer()

	// Choose transport based on flag
	if *port != "" {
		// HTTP mode - full logging to stdout
		logger.InitHTTPLogger()
		logger.Log.Info("starting HTTP MCP server",
			"addr", fmt.Sprintf("0.0.0.0:%s", *port),
		)

		if err := server.RunHTTPStreamableServer(mcpServer, fmt.Sprintf(":%s", *port)); err != nil {
			logger.Log.Error("failed to start HTTP MCP server", "error", err)
			os.Exit(1)
		}
	} else {
		// STDIO mode - only error logging to stderr
		logger.InitSTDIOLogger()

		if err := server.RunSTDIOServer(mcpServer); err != nil {
			logger.Log.Error("STDIO server failed", "error", err)
			os.Exit(1)
		}
	}
}
