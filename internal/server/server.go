package server

import (
	"github.com/aidanuno/qrkit/internal/server/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func CreateMCPServer() *mcp.Server {
	// Create a server
	server := mcp.NewServer(&mcp.Implementation{Name: "qrkit", Version: "0.1.0"}, nil)
	//Add available tools to the server
	tools.AddGenerateQRCodeTool(server)
	return server
}
