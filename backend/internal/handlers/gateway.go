package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/gvasels/personal-music-searchengine/internal/clients"
)

// GatewayHandler handles OpenAI-compatible API requests routed to Bedrock
type GatewayHandler struct {
	bedrockClient *clients.BedrockClient
	marengoClient *clients.MarengoClient
}

// NewGatewayHandler creates a new GatewayHandler
func NewGatewayHandler(bedrockClient *clients.BedrockClient, marengoClient *clients.MarengoClient) *GatewayHandler {
	return &GatewayHandler{
		bedrockClient: bedrockClient,
		marengoClient: marengoClient,
	}
}

// RegisterGatewayRoutes registers OpenAI-compatible routes
func (h *GatewayHandler) RegisterGatewayRoutes(e *echo.Echo) {
	// OpenAI-compatible endpoints
	v1 := e.Group("/v1")
	v1.POST("/chat/completions", h.CreateChatCompletion)
	v1.POST("/embeddings", h.CreateEmbedding)
	v1.POST("/embeddings/video", h.CreateVideoEmbedding)
	v1.GET("/models", h.ListModels)
}

// CreateChatCompletion handles POST /v1/chat/completions
func (h *GatewayHandler) CreateChatCompletion(c echo.Context) error {
	var req clients.ChatCompletionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Message: "Invalid request body",
				Type:    "invalid_request_error",
				Code:    "invalid_request",
			},
		})
	}

	// Validate required fields
	if req.Model == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Message: "model is required",
				Type:    "invalid_request_error",
				Code:    "missing_required_parameter",
			},
		})
	}
	if len(req.Messages) == 0 {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Message: "messages array is required and must not be empty",
				Type:    "invalid_request_error",
				Code:    "missing_required_parameter",
			},
		})
	}

	// Handle streaming
	if req.Stream {
		return h.streamChatCompletion(c, req)
	}

	// Non-streaming response
	resp, err := h.bedrockClient.CreateChatCompletion(c.Request().Context(), req)
	if err != nil {
		return handleBedrockError(c, err)
	}

	return c.JSON(http.StatusOK, resp)
}

// streamChatCompletion handles streaming chat completions
func (h *GatewayHandler) streamChatCompletion(c echo.Context, req clients.ChatCompletionRequest) error {
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	chunks, err := h.bedrockClient.StreamChatCompletion(c.Request().Context(), req)
	if err != nil {
		return handleBedrockError(c, err)
	}

	for chunk := range chunks {
		if chunk.Error != nil {
			// Send error in SSE format
			data, _ := json.Marshal(ErrorResponse{
				Error: ErrorDetail{
					Message: chunk.Error.Error(),
					Type:    "server_error",
				},
			})
			fmt.Fprintf(c.Response(), "data: %s\n\n", data)
			c.Response().Flush()
			return nil
		}

		if chunk.Done {
			fmt.Fprintf(c.Response(), "data: [DONE]\n\n")
			c.Response().Flush()
			return nil
		}

		// Format as OpenAI SSE
		sseChunk := map[string]interface{}{
			"id":      chunk.ID,
			"object":  "chat.completion.chunk",
			"created": chunk.Created,
			"model":   chunk.Model,
			"choices": []map[string]interface{}{
				{
					"index": 0,
					"delta": map[string]string{
						"content": chunk.Delta.Content,
					},
					"finish_reason": nil,
				},
			},
		}
		data, _ := json.Marshal(sseChunk)
		fmt.Fprintf(c.Response(), "data: %s\n\n", data)
		c.Response().Flush()
	}

	return nil
}

// CreateEmbedding handles POST /v1/embeddings
func (h *GatewayHandler) CreateEmbedding(c echo.Context) error {
	var req clients.EmbeddingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Message: "Invalid request body",
				Type:    "invalid_request_error",
				Code:    "invalid_request",
			},
		})
	}

	// Validate required fields
	if req.Model == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Message: "model is required",
				Type:    "invalid_request_error",
				Code:    "missing_required_parameter",
			},
		})
	}
	if req.Input == nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Message: "input is required",
				Type:    "invalid_request_error",
				Code:    "missing_required_parameter",
			},
		})
	}

	resp, err := h.bedrockClient.CreateEmbedding(c.Request().Context(), req)
	if err != nil {
		return handleBedrockError(c, err)
	}

	return c.JSON(http.StatusOK, resp)
}

// CreateVideoEmbedding handles POST /v1/embeddings/video
func (h *GatewayHandler) CreateVideoEmbedding(c echo.Context) error {
	var req clients.VideoEmbeddingOpenAIRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Message: "Invalid request body",
				Type:    "invalid_request_error",
				Code:    "invalid_request",
			},
		})
	}

	// Set default model
	if req.Model == "" {
		req.Model = "marengo-retrieval-2.7"
	}

	// Validate required fields
	if req.VideoURI == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Message: "video_uri is required",
				Type:    "invalid_request_error",
				Code:    "missing_required_parameter",
			},
		})
	}

	resp, err := h.marengoClient.CreateVideoEmbeddingOpenAI(c.Request().Context(), req)
	if err != nil {
		return handleBedrockError(c, err)
	}

	return c.JSON(http.StatusOK, resp)
}

// ListModels handles GET /v1/models
func (h *GatewayHandler) ListModels(c echo.Context) error {
	resp, err := h.bedrockClient.ListModels(c.Request().Context())
	if err != nil {
		return handleBedrockError(c, err)
	}

	return c.JSON(http.StatusOK, resp)
}

// ErrorResponse represents an OpenAI-compatible error response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error details
type ErrorDetail struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param,omitempty"`
	Code    string `json:"code,omitempty"`
}

// handleBedrockError converts Bedrock errors to OpenAI-compatible error responses
func handleBedrockError(c echo.Context, err error) error {
	errMsg := err.Error()

	// Check for common Bedrock errors
	if strings.Contains(errMsg, "AccessDeniedException") {
		return c.JSON(http.StatusForbidden, ErrorResponse{
			Error: ErrorDetail{
				Message: "Access denied to the requested model",
				Type:    "invalid_request_error",
				Code:    "model_not_found",
			},
		})
	}
	if strings.Contains(errMsg, "ResourceNotFoundException") {
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Error: ErrorDetail{
				Message: "The requested model was not found",
				Type:    "invalid_request_error",
				Code:    "model_not_found",
			},
		})
	}
	if strings.Contains(errMsg, "ThrottlingException") {
		return c.JSON(http.StatusTooManyRequests, ErrorResponse{
			Error: ErrorDetail{
				Message: "Rate limit exceeded. Please retry after a brief wait.",
				Type:    "rate_limit_error",
				Code:    "rate_limit_exceeded",
			},
		})
	}
	if strings.Contains(errMsg, "ValidationException") {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Message: errMsg,
				Type:    "invalid_request_error",
				Code:    "invalid_request",
			},
		})
	}

	// Default to internal server error
	return c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: ErrorDetail{
			Message: "An error occurred while processing your request",
			Type:    "server_error",
			Code:    "internal_error",
		},
	})
}
