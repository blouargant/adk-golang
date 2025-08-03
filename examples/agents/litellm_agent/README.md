# LiteLLM Agent Example

This example demonstrates how to create an AI agent using LiteLLM proxy API with the ADK-Golang framework. The agent includes DuckDuckGo search functionality and can work with any LLM provider that LiteLLM supports (OpenAI, Anthropic, Google, local models, etc.).

## Prerequisites

1. **LiteLLM Proxy Server**: You need a running LiteLLM proxy server. See [LiteLLM documentation](https://docs.litellm.ai/docs/proxy/quick_start) for setup instructions.
2. **API Key**: You need a valid API key/token for your LiteLLM proxy.
3. **Environment Variables**: Set the following environment variables:
   - `LITELLM_API_KEY`: Your LiteLLM API key/token (required)
   - `LITELLM_BASE_URL`: The URL of your LiteLLM proxy server (required, e.g., "http://localhost:4000")
   - `LITELLM_MODEL`: The model to use (optional, defaults to "gpt-3.5-turbo")

## LiteLLM Proxy Setup

### Quick Start with LiteLLM Proxy

1. Install LiteLLM:
   ```bash
   pip install litellm[proxy]
   ```

2. Create a config file `litellm_config.yaml`:
   ```yaml
   model_list:
     - model_name: gpt-3.5-turbo
       litellm_params:
         model: openai/gpt-3.5-turbo
         api_key: your-openai-api-key
     - model_name: claude-3-sonnet
       litellm_params:
         model: anthropic/claude-3-sonnet-20240229
         api_key: your-anthropic-api-key
   ```

3. Start the proxy:
   ```bash
   litellm --config litellm_config.yaml --port 4000
   ```

4. Your proxy will be available at `http://localhost:4000`

## Agent Setup

1. Set your LiteLLM configuration:
   ```bash
   export LITELLM_API_KEY="your-litellm-token"
   export LITELLM_BASE_URL="http://localhost:4000"
   export LITELLM_MODEL="gpt-3.5-turbo"  # or any model configured in your proxy
   ```

## Running the Agent

### Using the CLI

From the root of the ADK-Golang project:

```bash
# Run interactively
go run ./cmd/adk run examples/agents/litellm_agent

# Or use the web interface
go run ./cmd/adk web examples/agents
```

### Available Commands

- Ask questions that require web search: "What's the latest news about AI?"
- General questions: "Explain quantum computing"
- Current information: "What's happening in the tech industry today?"

## Features

- **LiteLLM Integration**: Uses LiteLLM proxy for unified access to multiple LLM providers
- **Bearer Token Authentication**: Uses proper bearer token authentication for LiteLLM
- **Web Search**: Integrated DuckDuckGo search functionality
- **Tool Calling**: Automatic tool selection and execution
- **Streaming Support**: Real-time response streaming
- **Multi-Provider Support**: Works with any LLM provider supported by LiteLLM

## Supported LLM Providers (via LiteLLM)

- OpenAI (GPT-3.5, GPT-4, etc.)
- Anthropic (Claude models)
- Google (Gemini, PaLM)
- Cohere
- Local models (Ollama, LM Studio, etc.)
- Azure OpenAI
- AWS Bedrock
- And many more...

## Configuration Examples

### Using OpenAI via LiteLLM
```bash
export LITELLM_API_KEY="your-openai-key"
export LITELLM_BASE_URL="http://localhost:4000"
export LITELLM_MODEL="gpt-4"
```

### Using Anthropic via LiteLLM
```bash
export LITELLM_API_KEY="your-anthropic-key"
export LITELLM_BASE_URL="http://localhost:4000"
export LITELLM_MODEL="claude-3-sonnet"
```

### Using Local Model via LiteLLM
```bash
export LITELLM_API_KEY="not-needed"
export LITELLM_BASE_URL="http://localhost:4000"
export LITELLM_MODEL="llama2"
```

## Architecture

```
User Input → ADK Agent → LiteLLM Proxy → LLM Provider → Response
                    ↓
            DuckDuckGo Search Tool
```

## Troubleshooting

1. **Connection Error**: Ensure your LiteLLM proxy is running and accessible at the specified `LITELLM_BASE_URL`.

2. **Authentication Error**: Verify your `LITELLM_API_KEY` is correct and has the necessary permissions.

3. **Model Not Found**: Check that the model specified in `LITELLM_MODEL` is configured in your LiteLLM proxy.

4. **Tool Call Issues**: The agent is configured to make only one tool call per query to prevent loops.

## Benefits of Using LiteLLM

- **Unified Interface**: Switch between different LLM providers without changing your code
- **Cost Optimization**: Easy comparison of costs across providers
- **Fallback Support**: Configure fallback models for reliability
- **Rate Limiting**: Built-in rate limiting and retry logic
- **Load Balancing**: Distribute requests across multiple models/providers
