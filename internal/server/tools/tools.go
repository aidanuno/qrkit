package tools

import (
	"context"

	"github.com/aidanuno/qrkit/internal/qr"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/yeqown/go-qrcode/v2"
)

type GenerateQRCodeArgs struct {
	RawQRData string `json:"raw_qr_data" jsonschema:"the data to encode into a QR code"`
}

func AddGenerateQRCodeTool(server *mcp.Server) {
	// Using the generic AddTool automatically populates the the input and output
	// schema of the tool.
	//
	// The schema considers 'json' and 'jsonschema' struct tags to get argument
	// names and descriptions.
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_qr_code",
		Description: "generates and return a QR code image from the given data",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args GenerateQRCodeArgs) (*mcp.CallToolResult, any, error) {

		qrcode, err := qrcode.New(args.RawQRData)
		if err != nil {
			return nil, nil, err
		}

		bytes, err := qr.QRCodeToBytes(qrcode)
		if err != nil {
			return nil, nil, err
		}

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
		}, nil, nil
	})
}
