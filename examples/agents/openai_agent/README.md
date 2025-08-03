# OpenAI Agent Example

This example demonstrates how to create an AI agent using OpenAI's API with the ADK-Golang framework. The agent includes DuckDuckGo search functionality.

## Prerequisites

1. **OpenAI API Key**: You need a valid OpenAI API key to use this agent.
2. **Environment Variables**: Set the following environment variables:
   - `OPENAI_API_KEY`: Your OpenAI API key (required)
   - `OPENAI_MODEL`: The model to use (optional, defaults to "gpt-3.5-turbo")
   - `OPENAI_BASE_URL`: Custom base URL for OpenAI API (optional)

## Setup

1. Set your OpenAI API key:
   ```bash
   export OPENAI_API_KEY="your-api-key-here"
   ```

2. Optionally, set the model:
   ```bash
   export OPENAI_MODEL="gpt-4"  # or any other available model
   ```

## Running the Agent

### Using the CLI

From the root of the ADK-Golang project:

```bash
# Run interactively
go run ./cmd/adk run examples/agents/openai_agent

# Or use the web interface
go run ./cmd/adk web examples/agents
```

### Available Commands

- Ask questions that require web search: "What's the latest news about AI?"
- General questions: "Explain quantum computing"
- Current information: "What's happening in the tech industry today?"

## Features

- **OpenAI Integration**: Uses OpenAI's chat completion API
- **Web Search**: Integrated DuckDuckGo search functionality
- **Tool Calling**: Automatic tool selection and execution
- **Streaming Support**: Real-time response streaming
- **Error Handling**: Robust error handling and retry logic

## Configuration

The agent is configured with:

- **Model**: Configurable via environment variable (default: gpt-3.5-turbo)
- **Temperature**: 0.3 (for consistent responses)
- **Max Tokens**: 2000
- **Max Tool Calls**: 1 (to prevent loops)
- **Retry Attempts**: 2
- **Streaming**: Enabled

## Example Interactions

```
User: What's the latest news about Go programming language?
Agent: I'll search for the latest news about Go programming language for you.

[Uses duckduckgo_search tool]

I found several recent developments about the Go programming language:

1. **Go 1.21 Release Notes** - Latest stable release with new features
   Link: https://golang.org/doc/go1.21

2. **Go Developer Survey 2023** - Community insights and trends
   Link: https://blog.golang.org/survey2023

[Additional context about Go's recent developments]
```

## Troubleshooting

1. **API Key Issues**: Ensure your `OPENAI_API_KEY` is set correctly
2. **Rate Limits**: OpenAI has rate limits; the agent includes retry logic
3. **Model Availability**: Make sure the specified model is available to your API key
4. **Network Issues**: Check your internet connection for search functionality

## Architecture

The OpenAI agent uses:

- `pkg/llmconnect/openai`: OpenAI API integration
- `pkg/agents`: LLM agent framework
- `pkg/tools`: DuckDuckGo search tool
- `pkg/core`: Core interfaces and types

## Extending the Agent

You can add more tools by:

1. Creating new tools using `tools.NewFunctionTool()`
2. Adding them with `agent.AddTool(tool)`
3. Updating the system instruction to mention the new tools
