package e2e

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestSTDIOMode_QRGeneration verifies the full MCP protocol flow over STDIO transport.
func TestSTDIOMode_QRGeneration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	projectRoot := filepath.Join("..", "..")

	// Empty implementation fields are acceptable for validation-only tests.
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "",
		Version: "",
	}, nil)

	// Test against the actual built binary to ensure production behavior.
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = projectRoot

	transport := &mcp.CommandTransport{
		Command: cmd,
	}

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		t.Fatalf("failed to connect to server: %v", err)
	}
	defer session.Close()

	// Verify tool discovery works as expected by MCP clients.
	listToolsResult, err := session.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		t.Fatalf("failed to list tools: %v", err)
	}

	if len(listToolsResult.Tools) != 1 {
		t.Errorf("expected 1 tool, got %d", len(listToolsResult.Tools))
	}

	if listToolsResult.Tools[0].Name != "generate_qr_code" {
		t.Errorf("expected tool 'generate_qr_code', got %s", listToolsResult.Tools[0].Name)
	}

	callResult, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "generate_qr_code",
		Arguments: map[string]interface{}{
			"raw_qr_data": "https://github.com/aidanuno/qrkit",
		},
	})
	if err != nil {
		t.Fatalf("failed to call tool: %v", err)
	}

	if callResult.IsError {
		t.Fatal("tool call returned error")
	}

	if len(callResult.Content) == 0 {
		t.Fatal("expected content in result")
	}

	// MCP spec requires image responses to be ImageContent with proper MIME type.
	imageContent, ok := callResult.Content[0].(*mcp.ImageContent)
	if !ok {
		t.Fatalf("expected ImageContent, got %T", callResult.Content[0])
	}

	if imageContent.MIMEType != "image/png" {
		t.Errorf("expected image/png, got %s", imageContent.MIMEType)
	}

	if len(imageContent.Data) == 0 {
		t.Fatal("image data is empty")
	}

	t.Logf("generated QR code: %d bytes, MIME: %s", len(imageContent.Data), imageContent.MIMEType)
}

// TestSTDIOMode_EmptyDataValidation ensures input validation prevents generating invalid QR codes.
func TestSTDIOMode_EmptyDataValidation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	projectRoot := filepath.Join("..", "..")

	// Empty implementation fields are acceptable for validation-only tests.
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "",
		Version: "",
	}, nil)

	cmd := exec.Command("go", "run", ".")
	cmd.Dir = projectRoot

	transport := &mcp.CommandTransport{
		Command: cmd,
	}

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		t.Fatalf("failed to connect to server: %v", err)
	}
	defer session.Close()

	callResult, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "generate_qr_code",
		Arguments: map[string]interface{}{
			"raw_qr_data": "",
		},
	})

	// Empty input should be rejected to prevent generating meaningless QR codes.
	if err == nil && !callResult.IsError {
		t.Error("expected error for empty data, got success")
	}
}
