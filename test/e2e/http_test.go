package e2e

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const testPort = 3000

// TestHTTPMode_QRGeneration verifies the full MCP protocol flow over HTTP SSE transport.
func TestHTTPMode_QRGeneration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	projectRoot := filepath.Join("..", "..")

	// Start server in HTTP mode on a dedicated port for testing.
	cmd := exec.CommandContext(ctx, "go", "run", ".", "--port", fmt.Sprintf("%d", testPort))
	cmd.Dir = projectRoot

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	defer cmd.Process.Kill()

	// Wait for server to be ready by polling.
	serverURL := fmt.Sprintf("http://localhost:%d", testPort)
	if err := waitForServer(ctx, serverURL); err != nil {
		t.Fatalf("server did not become ready: %v", err)
	}

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "",
		Version: "",
	}, nil)

	transport := &mcp.StreamableClientTransport{
		Endpoint: serverURL,
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

// TestHTTPMode_EmptyDataValidation ensures input validation prevents generating invalid QR codes.
func TestHTTPMode_EmptyDataValidation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	projectRoot := filepath.Join("..", "..")

	cmd := exec.CommandContext(ctx, "go", "run", ".", "--port", fmt.Sprintf("%d", testPort))
	cmd.Dir = projectRoot

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	defer cmd.Process.Kill()

	serverURL := fmt.Sprintf("http://localhost:%d", testPort)
	if err := waitForServer(ctx, serverURL); err != nil {
		t.Fatalf("server did not become ready: %v", err)
	}

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "",
		Version: "",
	}, nil)

	transport := &mcp.StreamableClientTransport{
		Endpoint: serverURL,
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

// waitForServer polls the server until it responds or context times out.
func waitForServer(ctx context.Context, url string) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
				return nil
			}
		}
	}
}
