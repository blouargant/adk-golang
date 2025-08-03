package openai

import (
	"testing"

	"github.com/agent-protocol/adk-golang/pkg/ptr"
)

func TestOpenAIConnectionInterface(t *testing.T) {
	// Test that OpenAIConnection implements the interface
	var _ OpenAIConnectionInterface = (*OpenAIConnection)(nil)

	config := &OpenAIConfig{
		APIKey:      "test-key",
		Model:       "gpt-3.5-turbo",
		Temperature: ptr.Float32(0.5),
		MaxTokens:   ptr.Ptr(1000),
		TopP:        ptr.Float32(0.9),
	}

	conn := NewOpenAIConnection(config)

	// Test GetConfig
	if conn.GetConfig() != config {
		t.Error("GetConfig should return the same config")
	}

	// Test GetModel
	if conn.GetModel() != "gpt-3.5-turbo" {
		t.Errorf("Expected model 'gpt-3.5-turbo', got '%s'", conn.GetModel())
	}

	// Test SetModel
	conn.SetModel("gpt-4")
	if conn.GetModel() != "gpt-4" {
		t.Errorf("Expected model 'gpt-4', got '%s'", conn.GetModel())
	}

	// Test GetProviderName
	if conn.GetProviderName() != "openai" {
		t.Errorf("Expected provider 'openai', got '%s'", conn.GetProviderName())
	}

	// Test ValidateConfig
	if err := conn.ValidateConfig(); err != nil {
		t.Errorf("Config validation failed: %v", err)
	}
}

func TestNewOpenAIConnectionWithProvider(t *testing.T) {
	config := &OpenAIConfig{
		APIKey: "test-key",
		Model:  "test-model",
	}

	conn := NewOpenAIConnectionWithProvider(config, "test-provider")

	if conn.GetProviderName() != "test-provider" {
		t.Errorf("Expected provider 'test-provider', got '%s'", conn.GetProviderName())
	}
}

func TestSetConfig(t *testing.T) {
	conn := NewOpenAIConnection(nil)

	newConfig := &OpenAIConfig{
		APIKey: "new-key",
		Model:  "new-model",
	}

	err := conn.SetConfig(newConfig)
	if err != nil {
		t.Errorf("SetConfig failed: %v", err)
	}

	if conn.GetConfig().APIKey != "new-key" {
		t.Error("Config was not updated properly")
	}

	// Test with nil config
	err = conn.SetConfig(nil)
	if err == nil {
		t.Error("SetConfig should fail with nil config")
	}

	// Test with invalid config (empty API key)
	invalidConfig := &OpenAIConfig{
		Model: "test-model",
	}
	err = conn.SetConfig(invalidConfig)
	if err == nil {
		t.Error("SetConfig should fail with empty API key")
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      *OpenAIConfig
		shouldError bool
	}{
		{
			name: "valid config",
			config: &OpenAIConfig{
				APIKey:      "test-key",
				Model:       "test-model",
				Temperature: ptr.Float32(0.7),
				TopP:        ptr.Float32(0.9),
				MaxTokens:   ptr.Ptr(1000),
			},
			shouldError: false,
		},
		{
			name:        "nil config",
			config:      nil,
			shouldError: true,
		},
		{
			name: "empty API key",
			config: &OpenAIConfig{
				Model: "test-model",
			},
			shouldError: true,
		},
		{
			name: "empty model",
			config: &OpenAIConfig{
				APIKey: "test-key",
			},
			shouldError: true,
		},
		{
			name: "invalid temperature (too high)",
			config: &OpenAIConfig{
				APIKey:      "test-key",
				Model:       "test-model",
				Temperature: ptr.Float32(3.0),
			},
			shouldError: true,
		},
		{
			name: "invalid temperature (negative)",
			config: &OpenAIConfig{
				APIKey:      "test-key",
				Model:       "test-model",
				Temperature: ptr.Float32(-0.1),
			},
			shouldError: true,
		},
		{
			name: "invalid top_p (too high)",
			config: &OpenAIConfig{
				APIKey: "test-key",
				Model:  "test-model",
				TopP:   ptr.Float32(1.1),
			},
			shouldError: true,
		},
		{
			name: "invalid top_p (negative)",
			config: &OpenAIConfig{
				APIKey: "test-key",
				Model:  "test-model",
				TopP:   ptr.Float32(-0.1),
			},
			shouldError: true,
		},
		{
			name: "invalid max_tokens (non-positive)",
			config: &OpenAIConfig{
				APIKey:    "test-key",
				Model:     "test-model",
				MaxTokens: ptr.Ptr(0),
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := &OpenAIConnection{config: tt.config}
			err := conn.ValidateConfig()

			if tt.shouldError && err == nil {
				t.Error("Expected validation error, but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no validation error, but got: %v", err)
			}
		})
	}
}

func TestProviderConnections(t *testing.T) {
	// Test Azure OpenAI connection
	azureConn := NewAzureOpenAIConnection("test-key", "https://test.openai.azure.com", "gpt-35-turbo")
	if azureConn.GetProviderName() != "azure-openai" {
		t.Errorf("Expected provider 'azure-openai', got '%s'", azureConn.GetProviderName())
	}

	// Test Anthropic compatible connection
	anthropicConn := NewAnthropicCompatibleConnection("test-key")
	if anthropicConn.GetProviderName() != "anthropic-compatible" {
		t.Errorf("Expected provider 'anthropic-compatible', got '%s'", anthropicConn.GetProviderName())
	}

	// Test local LLM connection
	localConn := NewLocalLLMConnection("http://localhost:1234/v1", "local-model")
	if localConn.GetProviderName() != "local-llm" {
		t.Errorf("Expected provider 'local-llm', got '%s'", localConn.GetProviderName())
	}

	// Test LiteLLM connection
	litellmConn, err := NewLiteLLMConnection("test-api-key", "http://localhost:4000", "gpt-3.5-turbo")
	if err != nil {
		t.Errorf("Failed to create LiteLLM connection: %v", err)
	}
	if litellmConn.GetProviderName() != "litellm" {
		t.Errorf("Expected provider 'litellm', got '%s'", litellmConn.GetProviderName())
	}

	// Test LiteLLM connection with empty baseURL (should fail)
	_, err = NewLiteLLMConnection("test-api-key", "", "gpt-3.5-turbo")
	if err == nil {
		t.Error("Expected error when creating LiteLLM connection with empty baseURL")
	}

	// Test LiteLLM connection with empty API key (should fail)
	_, err = NewLiteLLMConnection("", "http://localhost:4000", "gpt-3.5-turbo")
	if err == nil {
		t.Error("Expected error when creating LiteLLM connection with empty API key")
	}

	// Test LiteLLM connection with empty model (should fail)
	_, err = NewLiteLLMConnection("test-api-key", "http://localhost:4000", "")
	if err == nil {
		t.Error("Expected error when creating LiteLLM connection with empty model")
	}
}

func TestConfigBuilder(t *testing.T) {
	conn := NewConfigBuilder().
		WithAPIKey("test-key").
		WithModel("test-model").
		WithTemperature(0.5).
		WithMaxTokens(1000).
		WithTopP(0.9).
		WithStreaming(true).
		BuildConnection("test-provider")

	if conn.GetProviderName() != "test-provider" {
		t.Errorf("Expected provider 'test-provider', got '%s'", conn.GetProviderName())
	}

	config := conn.GetConfig()
	if config.APIKey != "test-key" {
		t.Error("API key not set correctly")
	}
	if config.Model != "test-model" {
		t.Error("Model not set correctly")
	}
	if config.Temperature == nil || *config.Temperature != 0.5 {
		t.Error("Temperature not set correctly")
	}
	if config.MaxTokens == nil || *config.MaxTokens != 1000 {
		t.Error("MaxTokens not set correctly")
	}
	if config.TopP == nil || *config.TopP != 0.9 {
		t.Error("TopP not set correctly")
	}
	if !config.Stream {
		t.Error("Stream not set correctly")
	}
}

func TestLiteLLMConnection(t *testing.T) {
	// Test successful creation
	conn, err := NewLiteLLMConnection("test-api-key", "http://localhost:4000", "gpt-3.5-turbo")
	if err != nil {
		t.Errorf("Failed to create LiteLLM connection: %v", err)
	}

	// Test provider name
	if conn.GetProviderName() != "litellm" {
		t.Errorf("Expected provider 'litellm', got '%s'", conn.GetProviderName())
	}

	// Test model
	if conn.GetModel() != "gpt-3.5-turbo" {
		t.Errorf("Expected model 'gpt-3.5-turbo', got '%s'", conn.GetModel())
	}

	// Test validation
	if err := conn.ValidateConfig(); err != nil {
		t.Errorf("Validation failed for valid config: %v", err)
	}

	// Test SetConfig with valid config
	newConfig := &OpenAIConfig{
		APIKey:  "new-api-key",
		BaseURL: "http://localhost:4000",
		Model:   "gpt-4",
	}
	if err := conn.SetConfig(newConfig); err != nil {
		t.Errorf("SetConfig failed: %v", err)
	}

	// Test SetConfig with missing BaseURL (should fail)
	invalidConfig := &OpenAIConfig{
		APIKey: "token",
		Model:  "gpt-4",
	}
	if err := conn.SetConfig(invalidConfig); err == nil {
		t.Error("Expected SetConfig to fail with missing BaseURL")
	}

	// Test validation with missing BaseURL
	conn.GetConfig().BaseURL = ""
	if err := conn.ValidateConfig(); err == nil {
		t.Error("Expected validation to fail with missing BaseURL")
	}
}
