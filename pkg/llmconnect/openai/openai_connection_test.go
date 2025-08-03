package openai

import (
	"context"
	"testing"
	"time"

	"github.com/agent-protocol/adk-golang/pkg/core"
	"github.com/agent-protocol/adk-golang/pkg/ptr"
)

func TestNewOpenAIConnection(t *testing.T) {
	config := &OpenAIConfig{
		APIKey:      "test-key",
		Model:       "gpt-4",
		Temperature: ptr.Float32(0.5),
		MaxTokens:   ptr.Ptr(1000),
		TopP:        ptr.Float32(0.9),
		Timeout:     15 * time.Second,
	}

	conn := NewOpenAIConnection(config)
	if conn == nil {
		t.Fatal("Expected connection to be created")
	}

	if conn.config.Model != "gpt-4" {
		t.Errorf("Expected model to be 'gpt-4', got '%s'", conn.config.Model)
	}

	if *conn.config.Temperature != 0.5 {
		t.Errorf("Expected temperature to be 0.5, got %f", *conn.config.Temperature)
	}
}

func TestDefaultOpenAIConfig(t *testing.T) {
	config := DefaultOpenAIConfig()
	if config == nil {
		t.Fatal("Expected default config to be created")
	}

	if config.Model != "gpt-3.5-turbo" {
		t.Errorf("Expected default model to be 'gpt-3.5-turbo', got '%s'", config.Model)
	}

	if *config.Temperature != 0.7 {
		t.Errorf("Expected default temperature to be 0.7, got %f", *config.Temperature)
	}
}

func TestConvertToOpenAIRequest(t *testing.T) {
	config := DefaultOpenAIConfig()
	conn := NewOpenAIConnection(config)

	request := &core.LLMRequest{
		Contents: []core.Content{
			{
				Role: "user",
				Parts: []core.Part{
					{
						Type: "text",
						Text: ptr.Ptr("Hello, how are you?"),
					},
				},
			},
		},
		Tools: []*core.FunctionDeclaration{
			{
				Name:        "get_weather",
				Description: "Get the current weather",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"location": map[string]interface{}{
							"type":        "string",
							"description": "The city and state",
						},
					},
					"required": []string{"location"},
				},
			},
		},
	}

	openaiReq, err := conn.convertToOpenAIRequest(request)
	if err != nil {
		t.Fatalf("Failed to convert request: %v", err)
	}

	if len(openaiReq.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(openaiReq.Messages))
	}

	if len(openaiReq.Tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(openaiReq.Tools))
	}

	if openaiReq.Tools[0].Function.Name != "get_weather" {
		t.Errorf("Expected tool name to be 'get_weather', got '%s'", openaiReq.Tools[0].Function.Name)
	}
}

func TestRoleMapping(t *testing.T) {
	conn := NewOpenAIConnection(nil)

	tests := []struct {
		input    string
		expected string
	}{
		{"user", "user"},
		{"assistant", "assistant"},
		{"agent", "tool"},
		{"model", "assistant"},
		{"system", "system"},
		{"unknown", "user"},
	}

	for _, test := range tests {
		result := conn.mapRole(test.input)
		if result != test.expected {
			t.Errorf("mapRole(%q) = %q, want %q", test.input, result, test.expected)
		}
	}
}

func TestJSONConversion(t *testing.T) {
	conn := NewOpenAIConnection(nil)

	// Test convertArgsToJSON
	args := map[string]any{
		"location": "New York",
		"units":    "metric",
	}
	jsonStr := conn.convertArgsToJSON(args)
	if jsonStr == "{}" {
		t.Error("Expected non-empty JSON string")
	}

	// Test convertJSONToArgs
	convertedArgs := conn.convertJSONToArgs(jsonStr)
	if convertedArgs["location"] != "New York" {
		t.Errorf("Expected location to be 'New York', got %v", convertedArgs["location"])
	}

	// Test with nil input
	emptyJSON := conn.convertArgsToJSON(nil)
	if emptyJSON != "{}" {
		t.Errorf("Expected empty JSON for nil input, got %s", emptyJSON)
	}
}

func TestCloseConnection(t *testing.T) {
	conn := NewOpenAIConnection(nil)
	ctx := context.Background()

	err := conn.Close(ctx)
	if err != nil {
		t.Errorf("Expected Close to return nil, got %v", err)
	}
}

func TestConvertToOpenAIRequestWithToolCalls(t *testing.T) {
	conn := NewOpenAIConnection(nil)

	// Test converting assistant message with tool calls
	request := &core.LLMRequest{
		Contents: []core.Content{
			{
				Role: "user",
				Parts: []core.Part{
					{
						Type: "text",
						Text: ptr.Ptr("What's the weather in New York?"),
					},
				},
			},
			{
				Role: "assistant",
				Parts: []core.Part{
					{
						Type: "text",
						Text: ptr.Ptr("I'll check the weather for you."),
					},
					{
						Type: "function_call",
						FunctionCall: &core.FunctionCall{
							ID:   "call_123",
							Name: "get_weather",
							Args: map[string]any{
								"location": "New York",
							},
						},
					},
				},
			},
			{
				Role: "tool",
				Parts: []core.Part{
					{
						Type: "function_response",
						FunctionResponse: &core.FunctionResponse{
							ID:   "call_123",
							Name: "get_weather",
							Response: map[string]any{
								"temperature": "22°C",
								"condition":   "sunny",
							},
						},
					},
				},
			},
		},
	}

	chatReq, err := conn.convertToOpenAIRequest(request)
	if err != nil {
		t.Fatalf("Expected no error converting request, got %v", err)
	}

	if len(chatReq.Messages) != 3 {
		t.Fatalf("Expected 3 messages, got %d", len(chatReq.Messages))
	}

	// Check user message
	if userMsg := chatReq.Messages[0].GetContent().AsAny(); userMsg != nil {
		if content, ok := userMsg.(*string); ok {
			if *content != "What's the weather in New York?" {
				t.Errorf("Expected user message content to be 'What's the weather in New York?', got %s", *content)
			}
		}
	}

	// Check assistant message with tool calls
	assistantToolCalls := chatReq.Messages[1].GetToolCalls()
	if len(assistantToolCalls) != 1 {
		t.Fatalf("Expected 1 tool call, got %d", len(assistantToolCalls))
	}

	toolCall := assistantToolCalls[0]
	if toolCall.ID != "call_123" {
		t.Errorf("Expected tool call ID to be 'call_123', got %s", toolCall.ID)
	}
	if toolCall.Function.Name != "get_weather" {
		t.Errorf("Expected tool call function name to be 'get_weather', got %s", toolCall.Function.Name)
	}

	// Check tool message
	toolCallID := chatReq.Messages[2].GetToolCallID()
	if toolCallID == nil || *toolCallID != "call_123" {
		t.Errorf("Expected tool call ID to be 'call_123', got %v", toolCallID)
	}
}
