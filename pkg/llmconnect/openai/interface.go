package openai

import (
	"context"

	"github.com/agent-protocol/adk-golang/pkg/core"
)

// OpenAIConnectionInterface defines the interface for OpenAI-compatible API connections.
// This interface can be implemented by different providers that are compatible with OpenAI's API format.
type OpenAIConnectionInterface interface {
	// Embed the core LLMConnection interface
	core.LLMConnection

	// GetConfig returns the configuration used by this connection
	GetConfig() *OpenAIConfig

	// SetConfig updates the configuration for this connection
	SetConfig(config *OpenAIConfig) error

	// GetModel returns the currently configured model
	GetModel() string

	// SetModel updates the model for this connection
	SetModel(model string)

	// ValidateConfig validates the current configuration
	ValidateConfig() error

	// GetProviderName returns the name of the API provider (e.g., "openai", "azure-openai", "anthropic-compatible")
	GetProviderName() string
}

// OpenAICompatibleClient defines the interface for OpenAI-compatible API clients.
// This allows for different underlying HTTP clients or SDK implementations.
type OpenAICompatibleClient interface {
	// CreateChatCompletion creates a chat completion
	CreateChatCompletion(ctx context.Context, req interface{}) (interface{}, error)

	// CreateChatCompletionStream creates a streaming chat completion
	CreateChatCompletionStream(ctx context.Context, req interface{}) (interface{}, error)

	// Close closes any resources held by the client
	Close() error
}
