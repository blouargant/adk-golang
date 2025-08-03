package main

import (
	"fmt"

	"github.com/agent-protocol/adk-golang/pkg/core"
	"github.com/agent-protocol/adk-golang/pkg/llmconnect/openai"
	"github.com/agent-protocol/adk-golang/pkg/ptr"
)

// demonstrateInterface shows how to use the OpenAI connection interface
func demonstrateInterface() {
	// Create different provider connections
	providers := []openai.OpenAIConnectionInterface{
		// Standard OpenAI
		openai.NewOpenAIConnection(&openai.OpenAIConfig{
			APIKey: "demo-key-1",
			Model:  "gpt-3.5-turbo",
		}),

		// Azure OpenAI
		openai.NewAzureOpenAIConnection("demo-key-2", "https://demo.openai.azure.com", "gpt-35-turbo"),

		// Groq
		openai.NewGroqConnection("demo-key-3"),

		// Local LLM
		openai.NewLocalLLMConnection("http://localhost:1234/v1", "llama-2-7b-chat"),

		// Custom configuration using builder
		openai.NewConfigBuilder().
			WithAPIKey("demo-key-4").
			WithBaseURL("https://api.custom.com/v1").
			WithModel("custom-model").
			WithTemperature(0.8).
			WithMaxTokens(4096).
			BuildConnection("custom-provider"),
	}

	// Add LiteLLM connection (with error handling)
	if litellmConn, err := openai.NewLiteLLMConnection("demo-api-key", "http://localhost:4000", "gpt-3.5-turbo"); err == nil {
		providers = append(providers, litellmConn)
	} else {
		fmt.Printf("Note: LiteLLM connection not added due to error: %v\n", err)
	}

	fmt.Println("OpenAI Connection Interface Demo")
	fmt.Println("================================")

	// Demonstrate polymorphism - same interface, different providers
	for i, provider := range providers {
		fmt.Printf("\n%d. Provider: %s\n", i+1, provider.GetProviderName())
		fmt.Printf("   Model: %s\n", provider.GetModel())

		// Validate configuration
		if err := provider.ValidateConfig(); err != nil {
			fmt.Printf("   Configuration Error: %v\n", err)
			continue
		}

		fmt.Printf("   Configuration: Valid ✓\n")

		// Demonstrate dynamic model switching
		if provider.GetProviderName() == "openai" {
			originalModel := provider.GetModel()
			provider.SetModel("gpt-4")
			fmt.Printf("   Model switched from %s to %s\n", originalModel, provider.GetModel())
		}

		// Show configuration details
		config := provider.GetConfig()
		fmt.Printf("   Temperature: %v\n", config.Temperature)
		fmt.Printf("   Max Tokens: %v\n", config.MaxTokens)
		fmt.Printf("   Base URL: %s\n", config.BaseURL)
	}

	fmt.Println("\nInterface Methods Demo")
	fmt.Println("=====================")

	// Pick one provider for method demonstration
	conn := providers[0]

	// Create a sample request
	_ = &core.LLMRequest{
		Contents: []core.Content{
			{
				Role: "user",
				Parts: []core.Part{
					{
						Type: "text",
						Text: ptr.Ptr("What is the capital of France?"),
					},
				},
			},
		},
	}

	fmt.Printf("Sample request created for provider: %s\n", conn.GetProviderName())
	fmt.Printf("Request model: %s\n", conn.GetModel())

	// In a real scenario, you would make the API call like this:
	// ctx := context.Background()
	// response, err := conn.GenerateContent(ctx, request)

	fmt.Println("Note: API call skipped in demo (would require valid API key)")
}

// demonstrateConfigurationManagement shows configuration management features
func demonstrateConfigurationManagement() {
	fmt.Println("\nConfiguration Management Demo")
	fmt.Println("============================")

	// Create a connection with default config
	conn := openai.NewOpenAIConnection(nil)
	fmt.Printf("Default model: %s\n", conn.GetModel())
	fmt.Printf("Default provider: %s\n", conn.GetProviderName())

	// Update configuration
	newConfig := &openai.OpenAIConfig{
		APIKey:      "new-api-key",
		Model:       "gpt-4-turbo",
		BaseURL:     "https://custom-endpoint.com/v1",
		Temperature: ptr.Float32(0.9),
		MaxTokens:   ptr.Ptr(8192),
		TopP:        ptr.Float32(0.95),
	}

	fmt.Println("\nUpdating configuration...")
	if err := conn.SetConfig(newConfig); err != nil {
		fmt.Printf("Configuration update failed: %v\n", err)
	} else {
		fmt.Printf("Configuration updated successfully\n")
		fmt.Printf("New model: %s\n", conn.GetModel())
		fmt.Printf("New max tokens: %d\n", *conn.GetConfig().MaxTokens)
	}

	// Test validation with invalid config
	fmt.Println("\nTesting validation with invalid configuration...")
	invalidConfig := &openai.OpenAIConfig{
		APIKey:      "test-key",
		Model:       "test-model",
		Temperature: ptr.Float32(5.0), // Invalid - too high
	}

	testConn := openai.NewOpenAIConnection(invalidConfig)
	if err := testConn.ValidateConfig(); err != nil {
		fmt.Printf("Validation error (expected): %v\n", err)
	}
}

func main() {
	demonstrateInterface()
	demonstrateConfigurationManagement()

	fmt.Println("\nDemo completed!")
	fmt.Println("\nThe OpenAI connection interface provides:")
	fmt.Println("• Unified API for multiple OpenAI-compatible providers")
	fmt.Println("• Runtime configuration management")
	fmt.Println("• Built-in validation")
	fmt.Println("• Easy provider switching")
	fmt.Println("• Polymorphic usage patterns")
}
