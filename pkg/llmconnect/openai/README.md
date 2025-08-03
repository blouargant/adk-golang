# OpenAI Connection Interface

This package provides a flexible interface for connecting to OpenAI and OpenAI-compatible APIs. The interface allows you to easily switch between different providers while maintaining the same API.

## Features

- **Unified Interface**: Use the same interface for OpenAI, Azure OpenAI, and other compatible services
- **Provider-Specific Constructors**: Pre-configured constructors for popular providers
- **Configuration Builder**: Fluent API for building configurations
- **Validation**: Built-in configuration validation
- **Dynamic Model Switching**: Change models at runtime
- **Multiple Provider Support**: Works with OpenAI, Azure OpenAI, Groq, Perplexity, Together AI, and local LLMs

## Basic Usage

### Standard OpenAI Connection

```go
import "github.com/agent-protocol/adk-golang/pkg/llmconnect/openai"

// Create a standard OpenAI connection
conn := openai.NewOpenAIConnection(&openai.OpenAIConfig{
    APIKey: "your-api-key",
    Model:  "gpt-3.5-turbo",
})

// Use the interface
var llmConn openai.OpenAIConnectionInterface = conn
```

### Using the Interface

```go
func processWithLLM(conn openai.OpenAIConnectionInterface) {
    fmt.Printf("Using provider: %s with model: %s\n", 
        conn.GetProviderName(), 
        conn.GetModel())
    
    // Create and send request
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
```

## Provider-Specific Constructors

### Azure OpenAI

```go
conn := openai.NewAzureOpenAIConnection(
    "your-azure-api-key",
    "https://your-resource.openai.azure.com",
    "your-deployment-name",
)
```

### Groq

```go
conn := openai.NewGroqConnection("your-groq-api-key")
```

### Perplexity

```go
conn := openai.NewPerplexityConnection("your-perplexity-api-key")
```

### Local LLM (LM Studio, LocalAI, etc.)

```go
conn := openai.NewLocalLLMConnection("http://localhost:1234/v1", "llama-2-7b-chat")
```

### Together AI

```go
conn := openai.NewTogetherAIConnection("your-together-api-key")
```

### DeepSeek

```go
conn := openai.NewDeepSeekConnection("your-deepseek-api-key")
```

### LiteLLM

```go
// LiteLLM proxy requires a baseURL, API key (used as bearer token), and model name
conn, err := openai.NewLiteLLMConnection("your-api-key", "http://localhost:4000", "gpt-3.5-turbo")
if err != nil {
    log.Fatalf("Failed to create LiteLLM connection: %v", err)
}
```

## Configuration Builder

Use the fluent configuration builder for custom setups:

```go
conn := openai.NewConfigBuilder().
    WithAPIKey("your-api-key").
    WithBaseURL("https://api.custom-provider.com/v1").
    WithModel("custom-model").
    WithTemperature(0.8).
    WithMaxTokens(4096).
    WithTopP(0.95).
    WithStreaming(true).
    BuildConnection("custom-provider")
```

## Interface Methods

The `OpenAIConnectionInterface` provides these methods:

### Core LLM Methods (from `core.LLMConnection`)
- `GenerateContent(ctx, request)` - Generate a single response
- `GenerateContentStream(ctx, request)` - Generate a streaming response
- `Close(ctx)` - Close the connection

### Configuration Methods
- `GetConfig()` - Get the current configuration
- `SetConfig(config)` - Update the configuration
- `ValidateConfig()` - Validate the current configuration

### Model Methods
- `GetModel()` - Get the current model
- `SetModel(model)` - Set the model

### Provider Methods
- `GetProviderName()` - Get the provider name (e.g., "openai", "azure-openai")

## Dynamic Configuration

You can change configuration at runtime:

```go
conn := openai.NewOpenAIConnection(&openai.OpenAIConfig{
    APIKey: "initial-key",
    Model:  "gpt-3.5-turbo",
})

// Change model
conn.SetModel("gpt-4")

// Update entire configuration
newConfig := &openai.OpenAIConfig{
    APIKey: "new-api-key",
    Model:  "gpt-4-turbo",
    BaseURL: "https://custom-endpoint.com/v1",
}

if err := conn.SetConfig(newConfig); err != nil {
    log.Printf("Configuration error: %v", err)
}
```

## Validation

The interface includes configuration validation:

```go
if err := conn.ValidateConfig(); err != nil {
    log.Printf("Configuration error: %v", err)
}
```

Validation checks:
- API key is not empty
- Model is specified
- Temperature is between 0 and 2
- TopP is between 0 and 1
- MaxTokens is positive

## Polymorphism Example

Use the interface to work with multiple providers:

```go
// Create connections to different providers
providers := []openai.OpenAIConnectionInterface{
    openai.NewOpenAIConnection(&openai.OpenAIConfig{APIKey: "openai-key", Model: "gpt-3.5-turbo"}),
    openai.NewGroqConnection("groq-key"),
    openai.NewPerplexityConnection("perplexity-key"),
    openai.NewLocalLLMConnection("http://localhost:1234/v1", "llama-2-7b"),
}

// Add LiteLLM connection (note: error handling required)
if litellmConn, err := openai.NewLiteLLMConnection("api-key", "http://localhost:4000", "gpt-3.5-turbo"); err == nil {
    providers = append(providers, litellmConn)
}

for _, provider := range providers {
    processWithLLM(provider)
}
```

## Error Handling

The interface provides comprehensive error handling:

```go
// Configuration errors
if err := conn.SetConfig(invalidConfig); err != nil {
    fmt.Printf("Config error: %v\n", err)
}

// API errors
response, err := conn.GenerateContent(ctx, request)
if err != nil {
    fmt.Printf("API error: %v\n", err)
}

// Validation errors
if err := conn.ValidateConfig(); err != nil {
    fmt.Printf("Validation error: %v\n", err)
}
```

## LiteLLM Integration

LiteLLM is a proxy service that provides a unified API for different LLM providers. It uses bearer token authentication instead of the standard OpenAI API key format:

```go
// Create a LiteLLM connection
conn, err := openai.NewLiteLLMConnection("your-api-key", "http://localhost:4000", "gpt-3.5-turbo")
if err != nil {
    log.Fatalf("Failed to create LiteLLM connection: %v", err)
}

// The baseURL is required for LiteLLM
fmt.Printf("Connected to LiteLLM at: %s\n", conn.GetConfig().BaseURL)

// Use like any other connection
request := &core.LLMRequest{
    Contents: []core.Content{
        {
            Role: "user",
            Parts: []core.Part{
                {
                    Type: "text",
                    Text: ptr.Ptr("Hello from LiteLLM!"),
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

fmt.Printf("LiteLLM Response: %s\n", *response.Content.Parts[0].Text)
```

### Key Differences for LiteLLM:
- **Authentication**: Uses bearer token in Authorization header instead of OpenAI API key format
- **Required BaseURL**: Must specify the LiteLLM proxy endpoint (cannot be empty)
- **Required Model**: Must specify the model name (cannot be empty)
- **Error Handling**: Constructor returns an error if baseURL, API key, or model is empty
- **Validation**: Additional validation ensures baseURL and model are always present

## Extending the Interface

To add support for new OpenAI-compatible providers:

1. Create a constructor function following the pattern:
```go
func NewCustomProviderConnection(apiKey string) OpenAIConnectionInterface {
    config := &OpenAIConfig{
        APIKey:  apiKey,
        BaseURL: "https://api.customprovider.com/v1",
        Model:   "custom-model",
        // ... other config
    }
    
    return NewOpenAIConnectionWithProvider(config, "custom-provider")
}
```

2. Use the existing interface without any modifications to your application code.

## Example Agents

The ADK-Golang project includes example agents demonstrating different provider integrations:

- **OpenAI Agent** (`examples/agents/openai_agent`): Uses standard OpenAI API
- **LiteLLM Agent** (`examples/agents/litellm_agent`): Uses LiteLLM proxy for multi-provider support

Both agents include the same functionality (DuckDuckGo search integration) but use different connection methods, demonstrating the flexibility of the interface.
