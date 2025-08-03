package openai

import (
	"context"
	"fmt"
	"log"

	"github.com/agent-protocol/adk-golang/pkg/core"
	"github.com/agent-protocol/adk-golang/pkg/ptr"
)

// Example usage of the OpenAI connection interface

// ExampleBasicUsage demonstrates basic usage of the OpenAI connection interface.
func ExampleBasicUsage() {
	// Create a standard OpenAI connection
	conn := NewOpenAIConnectionInterface(&OpenAIConfig{
		APIKey: "your-api-key",
		Model:  "gpt-3.5-turbo",
	})

	// Use the connection
	request := &core.LLMRequest{
		Contents: []core.Content{
			{
				Role: "user",
				Parts: []core.Part{
					{
						Type: "text",
						Text: ptr.Ptr("Hello, world!"),
					},
				},
			},
		},
	}

	ctx := context.Background()
	response, err := conn.GenerateContent(ctx, request)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Response: %s\n", *response.Content.Parts[0].Text)
}

// ExampleAzureOpenAI demonstrates using Azure OpenAI.
func ExampleAzureOpenAI() {
	// Create an Azure OpenAI connection
	conn := NewAzureOpenAIConnection(
		"your-azure-api-key",
		"https://your-resource.openai.azure.com",
		"your-deployment-name",
	)

	fmt.Printf("Provider: %s\n", conn.GetProviderName())
	fmt.Printf("Model: %s\n", conn.GetModel())
}

// ExampleLocalLLM demonstrates using a local LLM service.
func ExampleLocalLLM() {
	// Create a connection to a local LLM service (like LM Studio)
	conn := NewLocalLLMConnection("http://localhost:1234/v1", "llama-2-7b-chat")

	fmt.Printf("Provider: %s\n", conn.GetProviderName())
	fmt.Printf("Model: %s\n", conn.GetModel())
}

// ExampleConfigBuilder demonstrates using the configuration builder.
func ExampleConfigBuilder() {
	// Build a custom configuration
	conn := NewConfigBuilder().
		WithAPIKey("your-api-key").
		WithBaseURL("https://api.custom-provider.com/v1").
		WithModel("custom-model").
		WithTemperature(0.8).
		WithMaxTokens(4096).
		WithTopP(0.95).
		WithStreaming(true).
		BuildConnection("custom-provider")

	fmt.Printf("Provider: %s\n", conn.GetProviderName())
	fmt.Printf("Model: %s\n", conn.GetModel())

	// Validate the configuration
	if err := conn.ValidateConfig(); err != nil {
		log.Printf("Configuration error: %v", err)
	}
}

// ExampleLiteLLM demonstrates using LiteLLM proxy service.
func ExampleLiteLLM() {
	// Create a connection to LiteLLM proxy
	conn, err := NewLiteLLMConnection("your-api-key", "http://localhost:4000", "gpt-3.5-turbo")
	if err != nil {
		log.Printf("Failed to create LiteLLM connection: %v", err)
		return
	}

	fmt.Printf("Provider: %s\n", conn.GetProviderName())
	fmt.Printf("Model: %s\n", conn.GetModel())
	fmt.Printf("Base URL: %s\n", conn.GetConfig().BaseURL)
	fmt.Printf("Uses bearer token authentication: %t\n", true)

	// Validate that baseURL is required
	if err := conn.ValidateConfig(); err != nil {
		log.Printf("Configuration error: %v", err)
	} else {
		fmt.Printf("Configuration valid ✓\n")
	}
} // ExampleDynamicProviderSwitching demonstrates switching between providers.
func ExampleDynamicProviderSwitching() {
	// Create a slice of different providers
	providers := []OpenAIConnectionInterface{
		NewOpenAIConnection(&OpenAIConfig{
			APIKey: "openai-key",
			Model:  "gpt-3.5-turbo",
		}),
		NewAzureOpenAIConnection("azure-key", "https://test.openai.azure.com", "gpt-35-turbo"),
		NewGroqConnection("groq-key"),
		NewLocalLLMConnection("http://localhost:1234/v1", "llama-2-7b"),
	}

	// Use different providers based on requirements
	for _, provider := range providers {
		fmt.Printf("Using provider: %s with model: %s\n",
			provider.GetProviderName(),
			provider.GetModel())

		// You can switch models dynamically
		if provider.GetProviderName() == "openai" {
			provider.SetModel("gpt-4")
			fmt.Printf("Switched to model: %s\n", provider.GetModel())
		}
	}
}

// ExampleErrorHandling demonstrates proper error handling.
func ExampleErrorHandling() {
	conn := NewOpenAIConnection(&OpenAIConfig{})

	// This will fail validation
	if err := conn.ValidateConfig(); err != nil {
		fmt.Printf("Configuration error: %v\n", err)
	}

	// Try to set an invalid configuration
	invalidConfig := &OpenAIConfig{
		APIKey:      "test-key",
		Model:       "test-model",
		Temperature: ptr.Float32(5.0), // Invalid - too high
	}

	if err := conn.SetConfig(invalidConfig); err != nil {
		fmt.Printf("SetConfig error: %v\n", err)
	}
}

// ExampleInterfacePolymorphism demonstrates using the interface for polymorphism.
func ExampleInterfacePolymorphism() {
	// Function that works with any OpenAI-compatible provider
	processWithProvider := func(conn OpenAIConnectionInterface, prompt string) {
		fmt.Printf("Processing with %s using model %s\n",
			conn.GetProviderName(),
			conn.GetModel())

		// Create request
		request := &core.LLMRequest{
			Contents: []core.Content{
				{
					Role: "user",
					Parts: []core.Part{
						{
							Type: "text",
							Text: ptr.Ptr(prompt),
						},
					},
				},
			},
		}

		// Process request (in real usage, you'd handle the response)
		ctx := context.Background()
		_, err := conn.GenerateContent(ctx, request)
		if err != nil {
			fmt.Printf("Error processing with %s: %v\n", conn.GetProviderName(), err)
		}
	}

	// Use the same function with different providers
	providers := []OpenAIConnectionInterface{
		NewOpenAIConnection(&OpenAIConfig{APIKey: "key1", Model: "gpt-3.5-turbo"}),
		NewGroqConnection("key2"),
		NewPerplexityConnection("key3"),
	}

	for _, provider := range providers {
		processWithProvider(provider, "What is the capital of France?")
	}
}
