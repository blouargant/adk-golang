package openaiCompat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/agent-protocol/adk-golang/pkg/core"
	"github.com/agent-protocol/adk-golang/pkg/ptr"
)

var _ core.LLMConnection = (*openaiCompatConnection)(nil)

// openaiCompatConnection implements the LLMConnection interface for openaiCompat.
type openaiCompatConnection struct {
	baseURL    string
	httpClient *http.Client
	model      string
	config     *OpenaiCompatConfig
	apiKey     string
}

// OpenaiCompatConfig contains configuration options for openaiCompat connections.
type OpenaiCompatConfig struct {
	BaseURL     string        `json:"base_url"`
	Model       string        `json:"model"`
	APIKey      string        `json:"api_key,omitempty"`
	Temperature *float32      `json:"temperature,omitempty"`
	MaxTokens   *int          `json:"max_tokens,omitempty"`
	TopP        *float32      `json:"top_p,omitempty"`
	TopK        *int          `json:"top_k,omitempty"`
	Timeout     time.Duration `json:"timeout"`
	Stream      bool          `json:"stream"`
}

// DefaultOpenaiCompatConfig returns a default configuration for openaiCompat.
func DefaultOpenaiCompatConfig() *OpenaiCompatConfig {
	return &OpenaiCompatConfig{
		BaseURL:     "http://localhost:8080",
		Model:       "llama3.2",
		Temperature: ptr.Float32((0.7)),
		Timeout:     30 * time.Second,
		Stream:      false,
	}
}

// NewOpenaiCompatConnection creates a new openaiCompat connection with the given configuration.
func NewOpenaiCompatConnection(config *OpenaiCompatConfig) *openaiCompatConnection {
	if config == nil {
		config = DefaultOpenaiCompatConfig()
	}

	// Ensure BaseURL doesn't end with slash
	baseURL := strings.TrimSuffix(config.BaseURL, "/")

	return &openaiCompatConnection{
		baseURL: baseURL,
		model:   config.Model,
		config:  config,
		apiKey:  config.APIKey,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// GenerateContent sends a request to openaiCompat and returns the response.
func (c *openaiCompatConnection) GenerateContent(ctx context.Context, request *core.LLMRequest) (*core.LLMResponse, error) {
	// Convert ADK request to openaiCompat format
	openaiCompatReq, err := c.convertToOpenaiCompatRequest(request, false)
	if err != nil {
		return nil, fmt.Errorf("failed to convert request: %w", err)
	}

	// Make HTTP request
	resp, err := c.makeHTTPRequest(ctx, "/api/chat", openaiCompatReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var openaiCompatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&openaiCompatResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to ADK response
	return c.convertFromOpenaiCompatResponse(&openaiCompatResp), nil
}

// GenerateContentStream sends a request and returns a streaming response.
func (c *openaiCompatConnection) GenerateContentStream(ctx context.Context, request *core.LLMRequest) (<-chan *core.LLMResponse, error) {
	// Convert ADK request to openaiCompat format with streaming enabled
	openaiCompatReq, err := c.convertToOpenaiCompatRequest(request, true)
	if err != nil {
		return nil, fmt.Errorf("failed to convert request: %w", err)
	}

	// Make HTTP request
	resp, err := c.makeHTTPRequest(ctx, "/api/chat", openaiCompatReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	responseChan := make(chan *core.LLMResponse, 10)

	go func() {
		defer close(responseChan)
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		var accumulatedContent strings.Builder

		for {
			var chunk ChatResponse
			if err := decoder.Decode(&chunk); err != nil {
				if err == io.EOF {
					break
				}
				// Send error through channel
				errorResp := &core.LLMResponse{
					Content: &core.Content{
						Role: "assistant",
						Parts: []core.Part{
							{
								Type: "text",
								Text: ptr.Ptr(fmt.Sprintf("Error: %v", err)),
							},
						},
					},
					Partial: ptr.Ptr(false),
				}
				select {
				case responseChan <- errorResp:
				case <-ctx.Done():
				}
				return
			}

			// Accumulate content
			if chunk.Message.Content != "" {
				accumulatedContent.WriteString(chunk.Message.Content)
			}

			// Convert and send partial response
			partialResp := c.convertFromOpenaiCompatResponse(&chunk)

			// Set accumulated content for consistent streaming
			if accumulatedContent.Len() > 0 {
				partialResp.Content = &core.Content{
					Role: "assistant",
					Parts: []core.Part{
						{
							Type: "text",
							Text: ptr.Ptr(accumulatedContent.String()),
						},
					},
				}
			}

			// Mark as partial unless it's the final chunk
			partialResp.Partial = ptr.Ptr(!chunk.Done)

			select {
			case responseChan <- partialResp:
			case <-ctx.Done():
				return
			}

			// Break if this is the final chunk
			if chunk.Done {
				break
			}
		}
	}()

	return responseChan, nil
}

// Close closes the connection (no-op for HTTP-based connections).
func (c *openaiCompatConnection) Close(ctx context.Context) error {
	return nil
}

// convertToOpenaiCompatRequest converts an ADK LLMRequest to openaiCompat format.
func (c *openaiCompatConnection) convertToOpenaiCompatRequest(request *core.LLMRequest, stream bool) (*ChatRequest, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// Create openaiCompat ChatRequest
	chatReq := &ChatRequest{
		Model:   c.model,
		Stream:  ptr.Ptr(stream),
		Options: make(map[string]any),
	}

	// Convert ADK Contents to openaiCompat Messages
	messages := make([]Message, 0, len(request.Contents))
	for _, content := range request.Contents {
		message := Message{
			Role: c.mapRole(content.Role),
		}

		// Process content parts
		var textParts []string
		var toolCalls []ToolCall
		var images []ImageData

		for _, part := range content.Parts {
			switch part.Type {
			case "text":
				if part.Text != nil {
					textParts = append(textParts, *part.Text)
				}
			case "function_call":
				if part.FunctionCall != nil {
					toolCall := ToolCall{
						Function: ToolCallFunction{
							Name:      part.FunctionCall.Name,
							Arguments: ToolCallFunctionArguments(part.FunctionCall.Args),
						},
					}
					toolCalls = append(toolCalls, toolCall)
				}
			case "function_response":
				if part.FunctionResponse != nil {
					// For function responses, we include them as text content
					responseText := fmt.Sprintf("Function %s returned: %v", part.FunctionResponse.Name, part.FunctionResponse.Response)
					textParts = append(textParts, responseText)
				}
			}
		}

		// Combine text parts
		if len(textParts) > 0 {
			message.Content = strings.Join(textParts, "\n")
		}

		// Add tool calls if any
		if len(toolCalls) > 0 {
			message.ToolCalls = toolCalls
		}

		// Add images if any
		if len(images) > 0 {
			message.Images = images
		}

		messages = append(messages, message)
	}

	chatReq.Messages = messages

	// Convert tools to openaiCompat format
	if len(request.Tools) > 0 {
		openaiCompatTools := make(Tools, 0, len(request.Tools))
		for _, tool := range request.Tools {
			openaiCompatTool := Tool{
				Type: "function",
				Function: ToolFunction{
					Name:        tool.Name,
					Description: tool.Description,
				},
			}

			// Convert parameters if present
			if tool.Parameters != nil {
				openaiCompatTool.Function.Parameters.Type = "object"
				if props, hasProps := tool.Parameters["properties"].(map[string]interface{}); hasProps {
					openaiCompatTool.Function.Parameters.Properties = make(map[string]struct {
						Type        PropertyType `json:"type"`
						Items       any          `json:"items,omitempty"`
						Description string       `json:"description"`
						Enum        []any        `json:"enum,omitempty"`
					})

					for propName, propValue := range props {
						if propDetails, ok := propValue.(map[string]interface{}); ok {
							prop := struct {
								Type        PropertyType `json:"type"`
								Items       any          `json:"items,omitempty"`
								Description string       `json:"description"`
								Enum        []any        `json:"enum,omitempty"`
							}{}

							if typeVal, hasType := propDetails["type"].(string); hasType {
								prop.Type = PropertyType{typeVal}
							}
							if desc, hasDesc := propDetails["description"].(string); hasDesc {
								prop.Description = desc
							}
							if enum, hasEnum := propDetails["enum"].([]interface{}); hasEnum {
								prop.Enum = enum
							}
							if items, hasItems := propDetails["items"]; hasItems {
								prop.Items = items
							}

							openaiCompatTool.Function.Parameters.Properties[propName] = prop
						}
					}
				}
				if required, hasRequired := tool.Parameters["required"].([]interface{}); hasRequired {
					reqStrings := make([]string, 0, len(required))
					for _, req := range required {
						if reqStr, ok := req.(string); ok {
							reqStrings = append(reqStrings, reqStr)
						}
					}
					openaiCompatTool.Function.Parameters.Required = reqStrings
				}
			}

			openaiCompatTools = append(openaiCompatTools, openaiCompatTool)
		}
		chatReq.Tools = openaiCompatTools
	}

	// Apply connection-level configuration first
	if c.config.Temperature != nil {
		chatReq.Options["temperature"] = *c.config.Temperature
	}
	if c.config.MaxTokens != nil {
		chatReq.Options["num_predict"] = *c.config.MaxTokens
	}
	if c.config.TopP != nil {
		chatReq.Options["top_p"] = *c.config.TopP
	}
	if c.config.TopK != nil {
		chatReq.Options["top_k"] = *c.config.TopK
	}

	// Apply request-level configuration (overrides connection config)
	if request.Config != nil {
		if request.Config.Temperature != nil {
			chatReq.Options["temperature"] = *request.Config.Temperature
		}
		if request.Config.MaxTokens != nil {
			chatReq.Options["num_predict"] = *request.Config.MaxTokens
		}
		if request.Config.TopP != nil {
			chatReq.Options["top_p"] = *request.Config.TopP
		}
		if request.Config.TopK != nil {
			chatReq.Options["top_k"] = *request.Config.TopK
		}
	}

	return chatReq, nil
}

// convertFromOpenaiCompatResponse converts an openaiCompat response to ADK format.
func (c *openaiCompatConnection) convertFromOpenaiCompatResponse(resp *ChatResponse) *core.LLMResponse {
	if resp == nil {
		return &core.LLMResponse{}
	}

	response := &core.LLMResponse{
		Metadata: make(map[string]any),
	}

	// Convert message content
	if resp.Message.Content != "" || len(resp.Message.ToolCalls) > 0 {
		content := &core.Content{
			Role:  c.mapRoleFromOpenaiCompat(resp.Message.Role),
			Parts: make([]core.Part, 0),
		}

		// Add text content if present
		if resp.Message.Content != "" {
			content.Parts = append(content.Parts, core.Part{
				Type: "text",
				Text: ptr.Ptr(resp.Message.Content),
			})
		}

		// Add thinking content if present
		if resp.Message.Thinking != "" {
			content.Parts = append(content.Parts, core.Part{
				Type: "text",
				Text: ptr.Ptr(fmt.Sprintf("[Thinking: %s]", resp.Message.Thinking)),
			})
		}

		// Convert tool calls
		for _, toolCall := range resp.Message.ToolCalls {
			content.Parts = append(content.Parts, core.Part{
				Type: "function_call",
				FunctionCall: &core.FunctionCall{
					ID:   fmt.Sprintf("call_%d", toolCall.Function.Index),
					Name: toolCall.Function.Name,
					Args: map[string]any(toolCall.Function.Arguments),
				},
			})
		}

		response.Content = content
	}

	// Add metadata
	if resp.Model != "" {
		response.Metadata["model"] = resp.Model
	}
	if !resp.CreatedAt.IsZero() {
		response.Metadata["created_at"] = resp.CreatedAt
	}
	if resp.DoneReason != "" {
		response.Metadata["done_reason"] = resp.DoneReason
	}

	// Add metrics
	if resp.TotalDuration > 0 {
		response.Metadata["total_duration"] = resp.TotalDuration
	}
	if resp.LoadDuration > 0 {
		response.Metadata["load_duration"] = resp.LoadDuration
	}
	if resp.PromptEvalCount > 0 {
		response.Metadata["prompt_eval_count"] = resp.PromptEvalCount
	}
	if resp.PromptEvalDuration > 0 {
		response.Metadata["prompt_eval_duration"] = resp.PromptEvalDuration
	}
	if resp.EvalCount > 0 {
		response.Metadata["eval_count"] = resp.EvalCount
	}
	if resp.EvalDuration > 0 {
		response.Metadata["eval_duration"] = resp.EvalDuration
	}

	// Set partial flag (inverse of Done)
	response.Partial = ptr.Ptr(!resp.Done)

	return response
}

// makeHTTPRequest makes an HTTP request to the openaiCompat API.
func (c *openaiCompatConnection) makeHTTPRequest(ctx context.Context, endpoint string, payload interface{}) (*http.Response, error) {
	// Serialize payload
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create request
	url := c.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add Authorization header if API key is provided
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("openaiCompat API error (status %d): %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

// mapRole maps ADK roles to openaiCompat roles.
func (c *openaiCompatConnection) mapRole(role string) string {
	switch role {
	case "user":
		return "user"
	case "agent", "model", "assistant":
		return "assistant"
	case "system":
		return "system"
	default:
		return "user"
	}
}

// mapRoleFromOpenaiCompat maps openaiCompat roles to ADK roles.
func (c *openaiCompatConnection) mapRoleFromOpenaiCompat(role string) string {
	switch role {
	case "user":
		return "user"
	case "assistant":
		return "assistant"
	case "system":
		return "system"
	default:
		return "assistant"
	}
}
