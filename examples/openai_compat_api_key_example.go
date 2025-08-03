package main

import (
	"context"
	"fmt"
	"log"

	"github.com/agent-protocol/adk-golang/pkg/core"
	"github.com/agent-protocol/adk-golang/pkg/llmconnect/openaiCompat"
	"github.com/agent-protocol/adk-golang/pkg/ptr"
)

func main() {
	// Example 1: Connection without API key (for local servers like Ollama)
	config1 := &openaiCompat.OpenaiCompatConfig{
		BaseURL: "http://localhost:11434",
		Model:   "llama3.2",
	}
	conn1 := openaiCompat.NewOpenaiCompatConnection(config1)

	// Example 2: Connection with API key (for hosted services)
	config2 := &openaiCompat.OpenaiCompatConfig{
		BaseURL: "https://api.openai.com/v1",
		Model:   "gpt-4",
		APIKey:  "your-api-key-here",
	}
	conn2 := openaiCompat.NewOpenaiCompatConnection(config2)

	// Example 3: Using default config (no API key)
	conn3 := openaiCompat.NewOpenaiCompatConnection(nil)

	// Create a simple test request
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
	}

	ctx := context.Background()

	// This would make actual API calls, but we're just showing the setup
	fmt.Println("Connection 1 (no API key):", conn1 != nil)
	fmt.Println("Connection 2 (with API key):", conn2 != nil)
	fmt.Println("Connection 3 (default config):", conn3 != nil)

	// Uncomment below to make actual API calls (requires valid endpoints/keys)
	/*
		response, err := conn1.GenerateContent(ctx, request)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			log.Printf("Response: %+v", response)
		}
	*/

	// Use the variables to avoid "not used" errors
	_ = request
	_ = ctx

	log.Println("Example completed successfully!")
}
