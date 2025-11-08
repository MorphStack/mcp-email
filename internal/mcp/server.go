package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/brandon/mcp-email/internal/cache"
	"github.com/brandon/mcp-email/internal/config"
	"github.com/brandon/mcp-email/internal/email"
	"github.com/brandon/mcp-email/internal/tools"
	"github.com/sirupsen/logrus"
)

// Server represents the MCP server
type Server struct {
	config       *config.Config
	logger       *logrus.Logger
	tools        *tools.Registry
	emailManager *email.Manager
	cacheStore   *cache.Store
}

// NewServer creates a new MCP server instance
func NewServer(cfg *config.Config, emailManager *email.Manager, cacheStore *cache.Store, logger *logrus.Logger) (*Server, error) {
	// Initialize tool registry
	toolRegistry, err := tools.NewRegistry(cfg, emailManager, cacheStore, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create tool registry: %w", err)
	}

	return &Server{
		config:       cfg,
		logger:       logger,
		tools:        toolRegistry,
		emailManager: emailManager,
		cacheStore:   cacheStore,
	}, nil
}

// Run starts the MCP server with stdio transport
func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("Starting MCP server with stdio transport")

	// Simple MCP protocol implementation via stdio
	// This is a basic implementation that handles MCP requests
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			var req map[string]interface{}
			if err := decoder.Decode(&req); err != nil {
				if err == io.EOF {
					return nil
				}
				s.logger.WithError(err).Error("Failed to decode request")
				continue
			}

			resp := s.handleRequest(ctx, req)
			if err := encoder.Encode(resp); err != nil {
				s.logger.WithError(err).Error("Failed to encode response")
				continue
			}
		}
	}
}

// handleRequest processes an MCP request
func (s *Server) handleRequest(ctx context.Context, req map[string]interface{}) map[string]interface{} {
	method, _ := req["method"].(string)
	id, _ := req["id"]

	// Handle initialize request
	if method == "initialize" {
		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      id,
			"result": map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
				"serverInfo": map[string]interface{}{
					"name":    "mcp-email",
					"version": "1.0.0",
				},
			},
		}
	}

	// Handle tools/list request
	if method == "tools/list" {
		toolDefs := s.tools.GetToolDefinitions()
		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      id,
			"result": map[string]interface{}{
				"tools": toolDefs,
			},
		}
	}

	// Handle tools/call request
	if method == "tools/call" {
		params, _ := req["params"].(map[string]interface{})
		toolName, _ := params["name"].(string)
		arguments, _ := params["arguments"].(map[string]interface{})

		tool, exists := s.tools.GetTool(toolName)
		if !exists {
			return map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      id,
				"error": map[string]interface{}{
					"code":    -32601,
					"message": fmt.Sprintf("Tool not found: %s", toolName),
				},
			}
		}

		result, err := tool.Execute(arguments)
		if err != nil {
			return map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      id,
				"error": map[string]interface{}{
					"code":    -32603,
					"message": err.Error(),
				},
			}
		}

		// Serialize result to JSON string for text content
		resultJSON, err := json.Marshal(result)
		if err != nil {
			resultJSON = []byte(fmt.Sprintf("%v", result))
		}

		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      id,
			"result": map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": string(resultJSON),
					},
				},
			},
		}
	}

	// Unknown method
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"error": map[string]interface{}{
			"code":    -32601,
			"message": fmt.Sprintf("Method not found: %s", method),
		},
	}
}

