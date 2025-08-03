// Package main provides the OpenAI agent implementation for ADK-Golang.
// This agent demonstrates using OpenAI API with DuckDuckGo Search functionality.
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
// This agent uses OpenAI API for LLM inference and includes a search tool.
var RootAgent core.BaseAgent

func init() {
	log.Println("Initializing OpenAI Search Agent...")

	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Get model name from environment or use default
	modelName := os.Getenv("OPENAI_MODEL")
	if modelName == "" {
		modelName = "gpt-4o-mini"
	}

	// Create OpenAI configuration
	openaiConfig := &openai.OpenAIConfig{
		APIKey:      apiKey,
		Model:       modelName,
		Temperature: ptr.Float32(0.3), // Lower temperature for more consistent behavior
		MaxTokens:   ptr.Ptr(2000),
		TopP:        ptr.Float32(0.9),
		Timeout:     30 * time.Second,
		Stream:      false,
	}

	// Allow environment variable overrides
	if baseURL := os.Getenv("OPENAI_BASE_URL"); baseURL != "" {
		openaiConfig.BaseURL = baseURL
	}

	// Create OpenAI connection
	openaiConnection := openai.NewOpenAIConnection(openaiConfig)

	// Create agent configuration for OpenAI
	agentConfig := &agents.LlmAgentConfig{
		Model:            modelName,
		Temperature:      ptr.Float32(0.3), // Lower temperature for more consistent behavior
		MaxTokens:        ptr.Ptr(2000),
		MaxToolCalls:     1, // Only allow 1 tool call to prevent loops
		ToolCallTimeout:  30 * time.Second,
		RetryAttempts:    2,    // Reduce retries to prevent multiple calls
		StreamingEnabled: true, // Enable streaming for web UI
	}

	// Create the LLM agent with OpenAI
	llmAgent := agents.NewLLMAgent("openai_search_agent", "OpenAI-powered search agent with DuckDuckGo integration", agentConfig)

	// Set the LLM connection
	llmAgent.SetLLMConnection(openaiConnection)

	// Set instruction for better tool usage
	llmAgent.SetInstruction(`You are a helpful search assistant powered by OpenAI. When users ask questions:

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

	log.Printf("OpenAI Search Agent initialized successfully with model: %s", modelName)
	log.Println("Available tools:")
	log.Println("  - duckduckgo_search: Search the web using DuckDuckGo")
	log.Println("Example usage:")
	log.Println("  'What's the latest news about AI?'")
	log.Println("  'Search for information about Go programming language'")
	log.Println("  'Find the current weather in Tokyo'")
}
