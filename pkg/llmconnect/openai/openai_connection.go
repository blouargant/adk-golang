package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"

	"github.com/agent-protocol/adk-golang/pkg/core"
	"github.com/agent-protocol/adk-golang/pkg/ptr"
)

var _ core.LLMConnection = (*OpenAIConnection)(nil)

// OpenAIConnection implements the LLMConnection interface for OpenAI.
type OpenAIConnection struct {
	client *openai.Client
	config *OpenAIConfig
}

// OpenAIConfig contains configuration options for OpenAI connections.
type OpenAIConfig struct {
	APIKey      string        `json:"api_key"`
	BaseURL     string        `json:"base_url,omitempty"`
	Model       string        `json:"model"`
	Temperature *float32      `json:"temperature,omitempty"`
	MaxTokens   *int          `json:"max_tokens,omitempty"`
	TopP        *float32      `json:"top_p,omitempty"`
	Timeout     time.Duration `json:"timeout"`
	Stream      bool          `json:"stream"`
}

func printDebug(name string, data any) {
	if jsonData, err := json.MarshalIndent(data, "", "  "); err == nil {
		fmt.Printf("--- DEBUG [%s] ---\n%s\n", name, string(jsonData))
	} else {
		fmt.Printf("--- DEBUG [%s] ERROR ---\n%v\n", name, err)
	}
}

// DefaultOpenAIConfig returns a default configuration for OpenAI.
func DefaultOpenAIConfig() *OpenAIConfig {
	return &OpenAIConfig{
		Model:       "gpt-3.5-turbo",
		Temperature: ptr.Float32(0.7),
		MaxTokens:   ptr.Ptr(2048),
		TopP:        ptr.Float32(1.0),
		Timeout:     30 * time.Second,
		Stream:      false,
	}
}

// NewOpenAIConnection creates a new OpenAI connection with the given configuration.
func NewOpenAIConnection(config *OpenAIConfig) *OpenAIConnection {
	if config == nil {
		config = DefaultOpenAIConfig()
	}

	// Build client options
	opts := []option.RequestOption{
		option.WithAPIKey(config.APIKey),
	}

	if config.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(config.BaseURL))
	}

	// Create OpenAI client
	client := openai.NewClient(opts...)

	return &OpenAIConnection{
		client: &client,
		config: config,
	}
}

// GenerateContent sends a request to OpenAI and returns the response.
func (c *OpenAIConnection) GenerateContent(ctx context.Context, request *core.LLMRequest) (*core.LLMResponse, error) {
	// Convert ADK request to OpenAI format
	printDebug("request", request)
	chatReq, err := c.convertToOpenAIRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to convert request: %w", err)
	}
	printDebug("chatReq", chatReq)

	// Make request to OpenAI
	chatResp, err := c.client.Chat.Completions.New(ctx, *chatReq)
	if err != nil {
		fmt.Printf("OpenAI API request failed: %v\n", err)
		return nil, fmt.Errorf("OpenAI API request failed: %w", err)
	}
	printDebug("chatResp", chatResp)
	// Convert to ADK response
	resp := c.convertFromOpenAIResponse(*chatResp)
	printDebug("response", resp)
	return resp, nil
}

// GenerateContentStream sends a request and returns a streaming response.
func (c *OpenAIConnection) GenerateContentStream(ctx context.Context, request *core.LLMRequest) (<-chan *core.LLMResponse, error) {
	// Convert ADK request to OpenAI format with streaming enabled
	chatReq, err := c.convertToOpenAIRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to convert request: %w", err)
	}

	// Enable streaming options
	chatReq.StreamOptions = openai.ChatCompletionStreamOptionsParam{
		IncludeUsage: openai.Bool(true),
	}

	// Create stream
	stream := c.client.Chat.Completions.NewStreaming(ctx, *chatReq)

	responseChan := make(chan *core.LLMResponse, 10)

	go func() {
		defer close(responseChan)
		defer stream.Close()

		var accumulatedContent strings.Builder

		for stream.Next() {
			chunk := stream.Current()

			// Accumulate content
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				accumulatedContent.WriteString(chunk.Choices[0].Delta.Content)
			}

			// Convert and send partial response
			partialResp := c.convertFromOpenAIStreamResponse(chunk)

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

			// Check if this is the final chunk
			isDone := len(chunk.Choices) > 0 && chunk.Choices[0].FinishReason != ""
			partialResp.Partial = ptr.Ptr(!isDone)

			select {
			case responseChan <- partialResp:
			case <-ctx.Done():
				return
			}

			// Break if this is the final chunk
			if isDone {
				break
			}
		}

		// Check for stream errors
		if err := stream.Err(); err != nil {
			errorResp := &core.LLMResponse{
				Content: &core.Content{
					Role: "assistant",
					Parts: []core.Part{
						{
							Type: "text",
							Text: ptr.Ptr(fmt.Sprintf("Stream error: %v", err)),
						},
					},
				},
				Partial: ptr.Ptr(false),
			}
			select {
			case responseChan <- errorResp:
			case <-ctx.Done():
			}
		}
	}()

	return responseChan, nil
}

// Close closes the connection (no-op for HTTP-based connections).
func (c *OpenAIConnection) Close(ctx context.Context) error {
	return nil
}

// convertToOpenAIRequest converts an ADK LLMRequest to OpenAI format.
func (c *OpenAIConnection) convertToOpenAIRequest(request *core.LLMRequest) (*openai.ChatCompletionNewParams, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// Create OpenAI ChatCompletionNewParams
	chatReq := &openai.ChatCompletionNewParams{
		Model: shared.ChatModel(c.config.Model),
	}

	// Override model from request config if available
	if request.Config != nil && request.Config.Model != "" {
		chatReq.Model = shared.ChatModel(request.Config.Model)
	}

	// Set configuration parameters from connection config
	if c.config.Temperature != nil {
		chatReq.Temperature = openai.Float(float64(*c.config.Temperature))
	}
	if c.config.MaxTokens != nil {
		chatReq.MaxTokens = openai.Int(int64(*c.config.MaxTokens))
	}
	if c.config.TopP != nil {
		chatReq.TopP = openai.Float(float64(*c.config.TopP))
	}

	// Override with request config if available
	if request.Config != nil {
		if request.Config.Temperature != nil {
			chatReq.Temperature = openai.Float(float64(*request.Config.Temperature))
		}
		if request.Config.MaxTokens != nil {
			chatReq.MaxTokens = openai.Int(int64(*request.Config.MaxTokens))
		}
		if request.Config.TopP != nil {
			chatReq.TopP = openai.Float(float64(*request.Config.TopP))
		}
	}

	// Convert ADK Contents to OpenAI Messages
	messages := make([]openai.ChatCompletionMessageParamUnion, 0, len(request.Contents))
	for _, content := range request.Contents {
		role := c.mapRole(content.Role)
		//printDebug("role", role)

		// Process content parts to combine text
		var textParts []string
		var hasToolCalls bool

		for _, part := range content.Parts {
			//printDebug("part", part)
			switch part.Type {
			case "text":
				if part.Text != nil {
					textParts = append(textParts, *part.Text)
				}
			case "function_call":
				if part.FunctionCall != nil {
					hasToolCalls = true
					// For now, represent as text until we get tool calls working properly
					textParts = append(textParts, fmt.Sprintf("Calling function: %s(%s)",
						part.FunctionCall.Name, c.convertArgsToJSON(part.FunctionCall.Args)))
				}
			case "function_response":
				if part.FunctionResponse != nil {
					// For function responses, create a tool message
					msg, err := c.convertResponseToJSON(part.FunctionResponse.Response)
					if err == nil {
						toolMessage := openai.ToolMessage(
							msg,
							part.FunctionResponse.ID,
						)
						messages = append(messages, toolMessage)
					}
					continue // Skip the regular message creation for this part
				}
			default:
				if part.Text != nil {
					textParts = append(textParts, *part.Text)
				}
			}
		}

		// Skip messages with no content
		contentText := strings.Join(textParts, "\n")
		if contentText == "" && !hasToolCalls {
			continue
		}

		// Create message based on role and content
		switch role {
		case "system":
			if contentText != "" {
				message := openai.SystemMessage(contentText)
				messages = append(messages, message)
			}
		case "user":
			if contentText != "" {
				message := openai.UserMessage(contentText)
				messages = append(messages, message)
			}
		case "assistant":
			if contentText != "" {
				message := openai.AssistantMessage(contentText)
				messages = append(messages, message)
			}
		case "tool":
			// For tool messages, we will handle them separately
			if contentText != "" {
				toolMessage := openai.ToolMessage(contentText, "")
				messages = append(messages, toolMessage)

			}
		case "agent":
			// For agent messages, treat as assistant
			if contentText != "" {
				// Use assistant message for agent role
				message := openai.AssistantMessage(contentText)
				messages = append(messages, message)
			}
		default:
			// Default to user message
			if contentText != "" {
				message := openai.UserMessage(contentText)
				messages = append(messages, message)
			}
		}
	}

	chatReq.Messages = messages

	// Convert tools to OpenAI format - prefer tools from request.Tools, then from config
	var toolsToConvert []*core.FunctionDeclaration
	if len(request.Tools) > 0 {
		toolsToConvert = request.Tools
	} else if request.Config != nil && len(request.Config.Tools) > 0 {
		toolsToConvert = request.Config.Tools
	}

	if len(toolsToConvert) > 0 {
		tools := make([]openai.ChatCompletionToolParam, 0, len(toolsToConvert))
		for _, tool := range toolsToConvert {
			openaiTool := openai.ChatCompletionToolParam{
				Type: "function",
				Function: shared.FunctionDefinitionParam{
					Name:        tool.Name,
					Description: openai.String(tool.Description),
				},
			}

			// Convert parameters if present
			if tool.Parameters != nil {
				openaiTool.Function.Parameters = tool.Parameters
			}

			tools = append(tools, openaiTool)
		}

		chatReq.Tools = tools
	}

	return chatReq, nil
}

// convertFromOpenAIResponse converts an OpenAI ChatCompletion to ADK format.
func (c *OpenAIConnection) convertFromOpenAIResponse(resp openai.ChatCompletion) *core.LLMResponse {
	response := &core.LLMResponse{
		Metadata: make(map[string]any),
		Partial:  ptr.Ptr(false),
	}

	// Convert choices to content
	if len(resp.Choices) > 0 {
		choice := resp.Choices[0]
		content := &core.Content{
			Role:  c.mapRoleFromOpenAI("assistant"),
			Parts: []core.Part{},
		}

		// Add text content
		if choice.Message.Content != "" {
			content.Parts = append(content.Parts, core.Part{
				Type: "text",
				Text: ptr.Ptr(choice.Message.Content),
			})
		}

		// Add tool calls
		for _, toolCall := range choice.Message.ToolCalls {
			if toolCall.Type == "function" {
				content.Parts = append(content.Parts, core.Part{
					Type: "function_call",
					FunctionCall: &core.FunctionCall{
						ID:   toolCall.ID,
						Name: toolCall.Function.Name,
						Args: c.convertJSONToArgs(toolCall.Function.Arguments),
					},
				})
			}
		}

		response.Content = content
	}

	// Add metadata
	response.Metadata["model"] = resp.Model
	response.Metadata["id"] = resp.ID
	response.Metadata["created"] = resp.Created
	response.Metadata["object"] = resp.Object

	// Add usage information
	if resp.Usage.CompletionTokens > 0 {
		response.Metadata["completion_tokens"] = resp.Usage.CompletionTokens
	}
	if resp.Usage.PromptTokens > 0 {
		response.Metadata["prompt_tokens"] = resp.Usage.PromptTokens
	}
	if resp.Usage.TotalTokens > 0 {
		response.Metadata["total_tokens"] = resp.Usage.TotalTokens
	}

	return response
}

// convertFromOpenAIStreamResponse converts an OpenAI ChatCompletionChunk to ADK format.
func (c *OpenAIConnection) convertFromOpenAIStreamResponse(chunk openai.ChatCompletionChunk) *core.LLMResponse {
	response := &core.LLMResponse{
		Metadata: make(map[string]any),
		Partial:  ptr.Ptr(true),
	}

	// Convert choices to content
	if len(chunk.Choices) > 0 {
		choice := chunk.Choices[0]
		content := &core.Content{
			Role:  c.mapRoleFromOpenAI("assistant"),
			Parts: []core.Part{},
		}

		// Add text content from delta
		if choice.Delta.Content != "" {
			content.Parts = append(content.Parts, core.Part{
				Type: "text",
				Text: ptr.Ptr(choice.Delta.Content),
			})
		}

		// Add tool calls from delta
		for _, toolCall := range choice.Delta.ToolCalls {
			if toolCall.Type == "function" {
				content.Parts = append(content.Parts, core.Part{
					Type: "function_call",
					FunctionCall: &core.FunctionCall{
						ID:   toolCall.ID,
						Name: toolCall.Function.Name,
						Args: c.convertJSONToArgs(toolCall.Function.Arguments),
					},
				})
			}
		}

		response.Content = content

		// Check if this is the final chunk
		if choice.FinishReason != "" {
			response.Partial = ptr.Ptr(false)
			response.Metadata["finish_reason"] = choice.FinishReason
		}
	}

	// Add metadata
	response.Metadata["model"] = chunk.Model
	response.Metadata["id"] = chunk.ID
	response.Metadata["created"] = chunk.Created
	response.Metadata["object"] = chunk.Object

	return response
}

// convertArgsToJSON converts arguments map to JSON string.
func (c *OpenAIConnection) convertArgsToJSON(args map[string]any) string {
	if args == nil {
		return "{}"
	}

	jsonBytes, err := json.Marshal(args)
	if err != nil {
		return "{}"
	}

	return string(jsonBytes)
}

// convertResponseToJSON converts response map to JSON string.
func (c *OpenAIConnection) convertResponseToJSON(response map[string]any) (string, error) {
	if response == nil {
		return "{}", nil
	}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return "{}", err
	}
	return string(jsonBytes), nil
}

// convertJSONToArgs converts JSON string to arguments map.
func (c *OpenAIConnection) convertJSONToArgs(jsonStr string) map[string]any {
	var args map[string]any

	if err := json.Unmarshal([]byte(jsonStr), &args); err != nil {
		return make(map[string]any)
	}

	return args
}

// mapRole maps ADK roles to OpenAI roles.
func (c *OpenAIConnection) mapRole(role string) string {
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

// mapRoleFromOpenAI maps OpenAI roles to ADK roles.
func (c *OpenAIConnection) mapRoleFromOpenAI(role string) string {
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
