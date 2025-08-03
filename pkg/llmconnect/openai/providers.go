package openai

import (
	"fmt"

	"github.com/agent-protocol/adk-golang/pkg/ptr"
)

// Provider-specific configuration builders and constructors

// NewAzureOpenAIConnection creates a connection for Azure OpenAI Service.
func NewAzureOpenAIConnection(apiKey, endpoint, deploymentName string) OpenAIConnectionInterface {
	config := &OpenAIConfig{
		APIKey:      apiKey,
		BaseURL:     fmt.Sprintf("%s/openai/deployments/%s", endpoint, deploymentName),
		Model:       deploymentName, // Azure uses deployment name as model
		Temperature: ptr.Float32(0.7),
		MaxTokens:   ptr.Ptr(2048),
		TopP:        ptr.Float32(1.0),
		Stream:      false,
	}

	return NewOpenAIConnectionWithProvider(config, "azure-openai")
}

// NewAnthropicCompatibleConnection creates a connection for Anthropic's OpenAI-compatible API.
func NewAnthropicCompatibleConnection(apiKey string) OpenAIConnectionInterface {
	config := &OpenAIConfig{
		APIKey:      apiKey,
		BaseURL:     "https://api.anthropic.com/v1", // Example - adjust as needed
		Model:       "claude-3-sonnet-20240229",
		Temperature: ptr.Float32(0.7),
		MaxTokens:   ptr.Ptr(4096),
		TopP:        ptr.Float32(1.0),
		Stream:      false,
	}

	return NewOpenAIConnectionWithProvider(config, "anthropic-compatible")
}

// NewLocalLLMConnection creates a connection for local LLM services (like LM Studio, LocalAI, etc.).
func NewLocalLLMConnection(baseURL, model string) OpenAIConnectionInterface {
	config := &OpenAIConfig{
		APIKey:      "not-needed", // Local services often don't require API keys
		BaseURL:     baseURL,
		Model:       model,
		Temperature: ptr.Float32(0.7),
		MaxTokens:   ptr.Ptr(2048),
		TopP:        ptr.Float32(1.0),
		Stream:      false,
	}

	return NewOpenAIConnectionWithProvider(config, "local-llm")
}

// NewDeepSeekConnection creates a connection for DeepSeek's OpenAI-compatible API.
func NewDeepSeekConnection(apiKey string) OpenAIConnectionInterface {
	config := &OpenAIConfig{
		APIKey:      apiKey,
		BaseURL:     "https://api.deepseek.com",
		Model:       "deepseek-chat",
		Temperature: ptr.Float32(0.7),
		MaxTokens:   ptr.Ptr(4096),
		TopP:        ptr.Float32(1.0),
		Stream:      false,
	}

	return NewOpenAIConnectionWithProvider(config, "deepseek")
}

// NewPerplexityConnection creates a connection for Perplexity's OpenAI-compatible API.
func NewPerplexityConnection(apiKey string) OpenAIConnectionInterface {
	config := &OpenAIConfig{
		APIKey:      apiKey,
		BaseURL:     "https://api.perplexity.ai",
		Model:       "llama-3.1-sonar-small-128k-online",
		Temperature: ptr.Float32(0.7),
		MaxTokens:   ptr.Ptr(4096),
		TopP:        ptr.Float32(1.0),
		Stream:      false,
	}

	return NewOpenAIConnectionWithProvider(config, "perplexity")
}

// NewGroqConnection creates a connection for Groq's OpenAI-compatible API.
func NewGroqConnection(apiKey string) OpenAIConnectionInterface {
	config := &OpenAIConfig{
		APIKey:      apiKey,
		BaseURL:     "https://api.groq.com/openai/v1",
		Model:       "llama3-8b-8192",
		Temperature: ptr.Float32(0.7),
		MaxTokens:   ptr.Ptr(8192),
		TopP:        ptr.Float32(1.0),
		Stream:      false,
	}

	return NewOpenAIConnectionWithProvider(config, "groq")
}

// NewLiteLLMConnection creates a connection for LiteLLM proxy services.
// LiteLLM uses bearer token authentication instead of the standard OpenAI API key format.
// The baseURL parameter is required and must point to your LiteLLM proxy endpoint.
// The model parameter is required and specifies which model to use.
func NewLiteLLMConnection(apiKey, baseURL, model string) (OpenAIConnectionInterface, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("baseURL is required for LiteLLM connections")
	}

	if apiKey == "" {
		return nil, fmt.Errorf("API key is required for LiteLLM connections")
	}

	if model == "" {
		return nil, fmt.Errorf("model is required for LiteLLM connections")
	}

	config := &OpenAIConfig{
		APIKey:      apiKey, // Will be used as bearer token in headers
		BaseURL:     baseURL,
		Model:       model,
		Temperature: ptr.Float32(0.7),
		MaxTokens:   ptr.Ptr(4096),
		TopP:        ptr.Float32(1.0),
		Stream:      false,
	}

	return NewLiteLLMConnectionWithProvider(config, "litellm"), nil
} // NewTogetherAIConnection creates a connection for Together AI's OpenAI-compatible API.
func NewTogetherAIConnection(apiKey string) OpenAIConnectionInterface {
	config := &OpenAIConfig{
		APIKey:      apiKey,
		BaseURL:     "https://api.together.xyz/v1",
		Model:       "meta-llama/Llama-2-7b-chat-hf",
		Temperature: ptr.Float32(0.7),
		MaxTokens:   ptr.Ptr(4096),
		TopP:        ptr.Float32(1.0),
		Stream:      false,
	}

	return NewOpenAIConnectionWithProvider(config, "together-ai")
}

// ConfigBuilder provides a fluent interface for building OpenAI configurations.
type ConfigBuilder struct {
	config *OpenAIConfig
}

// NewConfigBuilder creates a new configuration builder.
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: DefaultOpenAIConfig(),
	}
}

// WithAPIKey sets the API key.
func (b *ConfigBuilder) WithAPIKey(apiKey string) *ConfigBuilder {
	b.config.APIKey = apiKey
	return b
}

// WithBaseURL sets the base URL.
func (b *ConfigBuilder) WithBaseURL(baseURL string) *ConfigBuilder {
	b.config.BaseURL = baseURL
	return b
}

// WithModel sets the model.
func (b *ConfigBuilder) WithModel(model string) *ConfigBuilder {
	b.config.Model = model
	return b
}

// WithTemperature sets the temperature.
func (b *ConfigBuilder) WithTemperature(temperature float32) *ConfigBuilder {
	b.config.Temperature = ptr.Float32(temperature)
	return b
}

// WithMaxTokens sets the maximum tokens.
func (b *ConfigBuilder) WithMaxTokens(maxTokens int) *ConfigBuilder {
	b.config.MaxTokens = ptr.Ptr(maxTokens)
	return b
}

// WithTopP sets the top-p value.
func (b *ConfigBuilder) WithTopP(topP float32) *ConfigBuilder {
	b.config.TopP = ptr.Float32(topP)
	return b
}

// WithStreaming enables or disables streaming.
func (b *ConfigBuilder) WithStreaming(stream bool) *ConfigBuilder {
	b.config.Stream = stream
	return b
}

// Build returns the built configuration.
func (b *ConfigBuilder) Build() *OpenAIConfig {
	return b.config
}

// BuildConnection creates a connection with the built configuration.
func (b *ConfigBuilder) BuildConnection(providerName string) OpenAIConnectionInterface {
	return NewOpenAIConnectionWithProvider(b.config, providerName)
}
