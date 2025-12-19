package tools

import (
	"context"
	"errors"

	"github.com/aidanuno/qrkit/internal/qr"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/yeqown/go-qrcode/v2"
)

type GenerateQRCodeArgs struct {
	RawQRData string `json:"raw_qr_data" jsonschema:"the data to encode into a QR code"`
}

var (
	ErrEmptyData = errors.New("raw_qr_data cannot be empty")
)

// generateQRCodeHandler contains the core QR code generation logic.
// Extracted as a separate function to enable testing without MCP server setup.
func generateQRCodeHandler(ctx context.Context, args GenerateQRCodeArgs) (*mcp.CallToolResult, error) {
	// Validate input
	if args.RawQRData == "" {
		return nil, ErrEmptyData
	}

	// Generate QR code
	qrcodeResult, err := qrcode.New(args.RawQRData)
	if err != nil {
		return nil, err
	}

	// Convert to bytes
	bytes, err := qr.CodeToBytes(qrcodeResult)
	if err != nil {
		return nil, err
	}

	// Build response
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.ImageContent{
				Meta: mcp.Meta{},
				Annotations: &mcp.Annotations{
					Audience: []mcp.Role{"user"},
				},
				Data:     bytes,
				MIMEType: "image/png",
			},
		},
	}, nil
}

// AddGenerateQRCodeTool registers the QR code generation tool with the MCP server.
func AddGenerateQRCodeTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_qr_code",
		Description: "generates and return a QR code PNG image from the given data",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args GenerateQRCodeArgs) (*mcp.CallToolResult, any, error) {
		result, err := generateQRCodeHandler(ctx, args)
		return result, nil, err
	})
}
