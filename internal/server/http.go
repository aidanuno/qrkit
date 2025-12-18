package server

import (
	"context"
	"net/http"
	"time"

	"github.com/aidanuno/qrkit/internal/logger"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func RunHTTPStreamableServer(server *mcp.Server, addr string) error {
	server.AddReceivingMiddleware(createLoggingMiddleware())
	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return server
	}, nil)
	return http.ListenAndServe(addr, handler)
}

// createLoggingMiddleware creates an MCP middleware that logs method calls.
func createLoggingMiddleware() mcp.Middleware {
	return func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(
			ctx context.Context,
			method string,
			req mcp.Request,
		) (mcp.Result, error) {
			start := time.Now()
			sessionID := req.GetSession().ID()

			// Log request details.
			logger.Log.Info("request",
				"session_id", sessionID,
				"method", method,
			)

			// Call the actual handler.
			result, err := next(ctx, method, req)

			// Log response details.
			duration := time.Since(start)

			if err != nil {
				logger.Log.Error("response",
					"session_id", sessionID,
					"method", method,
					"duration", duration,
					"error", err,
				)
			} else {
				logger.Log.Info("response",
					"session_id", sessionID,
					"method", method,
					"duration", duration,
				)
			}

			return result, err
		}
	}
}
