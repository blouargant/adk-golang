// Package main provides the LiteLLM agent implementation for ADK-Golang.
// This agent demonstrates using LiteLLM proxy API with DuckDuckGo Search functionality.
package main

import (
	"log"
	"os"
	"time"

	"github.com/agent-protocol/adk-golang/pkg/agents"
	"github.com/agent-protocol/adk-golang/pkg/core"
	"github.com/agent-protocol/adk-golang/pkg/llmconnect/openai"
	"github.com/agent-protocol/adk-golang/pkg/ptr"
	"github.com/agent-protocol/adk-golang/pkg/tools"
)

// RootAgent creates and configures the main agent with DuckDuckGo Search capability.
// This agent uses LiteLLM proxy API for LLM inference and includes a search tool.
var RootAgent core.BaseAgent

func init() {
	log.Println("Initializing LiteLLM Search Agent...")

	// Get API key from environment
	apiKey := os.Getenv("LITELLM_API_KEY")
	if apiKey == "" {
		log.Fatal("LITELLM_API_KEY environment variable is required")
	}

	// Get base URL from environment
	baseURL := os.Getenv("LITELLM_BASE_URL")
	if baseURL == "" {
		log.Fatal("LITELLM_BASE_URL environment variable is required (e.g., http://localhost:4000)")
	}

	// Get model name from environment or use default
	modelName := os.Getenv("LITELLM_MODEL")
	if modelName == "" {
		modelName = "gpt-3.5-turbo"
	}

	// Create LiteLLM connection with error handling
	litellmConnection, err := openai.NewLiteLLMConnection(apiKey, baseURL, modelName)
	if err != nil {
		log.Fatalf("Failed to create LiteLLM connection: %v", err)
	}

	// Validate the connection configuration
	if err := litellmConnection.ValidateConfig(); err != nil {
		log.Fatalf("LiteLLM configuration validation failed: %v", err)
	}

	// Create agent configuration for LiteLLM
	agentConfig := &agents.LlmAgentConfig{
		Model:            modelName,
		Temperature:      ptr.Float32(0.3), // Lower temperature for more consistent behavior
		MaxTokens:        ptr.Ptr(2000),
		MaxToolCalls:     1, // Only allow 1 tool call to prevent loops
		ToolCallTimeout:  30 * time.Second,
		RetryAttempts:    2,    // Reduce retries to prevent multiple calls
		StreamingEnabled: true, // Enable streaming for web UI
	}

	// Create the LLM agent with LiteLLM
	llmAgent := agents.NewLLMAgent("litellm_search_agent", "LiteLLM-powered search agent with DuckDuckGo integration", agentConfig)

	// Set the LLM connection
	llmAgent.SetLLMConnection(litellmConnection)

	// Set instruction for better tool usage
	llmAgent.SetInstruction(`You are a helpful search assistant powered by LiteLLM proxy. When users ask questions:

- If user provide unclear or ambiguous queries, ask for clarification.
- Use the duckduckgo_search tool ONCE to find relevant information
- Present the search results in a clear, organized format with:
   - A brief summary of what you found
   - List the key results with titles and brief descriptions
   - Include relevant URLs so users can learn more
- Do NOT call the search tool multiple times for the same query
- Always provide a complete response based on the search results
- If you cannot find relevant information, politely inform the user

Example response format:
"I found several great resources about [topic]:

1. **[Title 1]** - [Brief description]
   Link: [URL]

2. **[Title 2]** - [Brief description] 
   Link: [URL]

[Additional context or summary]"

Remember: Call each tool only ONCE per user question.`)

	// Create DuckDuckGo search tool
	searchTool := tools.NewDuckDuckGoSearchTool()

	// Add the search tool to the agent
	llmAgent.AddTool(searchTool)

	// Set the global root agent
	RootAgent = llmAgent

	log.Printf("LiteLLM Search Agent initialized successfully")
	log.Printf("  Provider: %s", litellmConnection.GetProviderName())
	log.Printf("  Model: %s", litellmConnection.GetModel())
	log.Printf("  Base URL: %s", litellmConnection.GetConfig().BaseURL)
	log.Println("Available tools:")
	log.Println("  - duckduckgo_search: Search the web using DuckDuckGo")
	log.Println("Example usage:")
	log.Println("  'What's the latest news about AI?'")
	log.Println("  'Search for information about Go programming language'")
	log.Println("  'Find the current weather in Tokyo'")
}
