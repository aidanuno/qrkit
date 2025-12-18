package server

import (
	"context"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func RunSTDIOServer(server *mcp.Server) error {
	return server.Run(context.Background(), &mcp.StdioTransport{})
}
