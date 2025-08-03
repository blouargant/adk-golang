package openai

import (
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// LiteLLMConnection implements OpenAI connection interface with bearer token authentication for LiteLLM proxy.
type LiteLLMConnection struct {
	*OpenAIConnection
}

// NewLiteLLMConnectionWithProvider creates a new LiteLLM connection with bearer token authentication.
func NewLiteLLMConnectionWithProvider(config *OpenAIConfig, providerName string) *LiteLLMConnection {
	if config == nil {
		config = DefaultOpenAIConfig()
	}

	// Build client options with bearer token authentication
	opts := []option.RequestOption{
		// Use the APIKey as a bearer token instead of OpenAI API key format
		option.WithHeader("Authorization", fmt.Sprintf("Bearer %s", config.APIKey)),
	}

	if config.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(config.BaseURL))
	}

	// Create OpenAI client with bearer token authentication
	client := openai.NewClient(opts...)

	baseConn := &OpenAIConnection{
		client:       &client,
		config:       config,
		providerName: providerName,
	}

	return &LiteLLMConnection{
		OpenAIConnection: baseConn,
	}
}

// SetConfig updates the configuration for this LiteLLM connection.
// This override ensures bearer token authentication is maintained.
func (c *LiteLLMConnection) SetConfig(config *OpenAIConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	// Validate the new configuration
	if config.APIKey == "" {
		return fmt.Errorf("API key is required")
	}
	if config.Model == "" {
		return fmt.Errorf("model is required")
	}
	if config.BaseURL == "" {
		return fmt.Errorf("baseURL is required for LiteLLM connections")
	}

	// Update the configuration
	c.config = config

	// Recreate the client with bearer token authentication
	opts := []option.RequestOption{
		option.WithHeader("Authorization", fmt.Sprintf("Bearer %s", config.APIKey)),
	}

	if config.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(config.BaseURL))
	}

	client := openai.NewClient(opts...)
	c.client = &client

	return nil
}

// ValidateConfig validates the current configuration for LiteLLM.
func (c *LiteLLMConnection) ValidateConfig() error {
	if c.config == nil {
		return fmt.Errorf("configuration is nil")
	}

	if c.config.APIKey == "" {
		return fmt.Errorf("API key is required")
	}

	if c.config.Model == "" {
		return fmt.Errorf("model is required")
	}

	// LiteLLM requires a baseURL
	if c.config.BaseURL == "" {
		return fmt.Errorf("baseURL is required for LiteLLM connections")
	}

	if c.config.Temperature != nil && (*c.config.Temperature < 0 || *c.config.Temperature > 2) {
		return fmt.Errorf("temperature must be between 0 and 2")
	}

	if c.config.TopP != nil && (*c.config.TopP < 0 || *c.config.TopP > 1) {
		return fmt.Errorf("top_p must be between 0 and 1")
	}

	if c.config.MaxTokens != nil && *c.config.MaxTokens <= 0 {
		return fmt.Errorf("max_tokens must be positive")
	}

	return nil
}

// GetProviderName returns the name of the LiteLLM provider.
func (c *LiteLLMConnection) GetProviderName() string {
	if c.providerName != "" {
		return c.providerName
	}
	return "litellm"
}
