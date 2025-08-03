#!/bin/bash

# Example script showing how to set up and run the LiteLLM agent
# This is for demonstration purposes - adjust the values for your setup

echo "🔧 LiteLLM Agent Setup Example"
echo "============================="

# Example 1: Using OpenAI via LiteLLM
echo ""
echo "Example 1: OpenAI via LiteLLM"
echo "-----------------------------"
echo "export LITELLM_API_KEY=\"sk-your-openai-api-key-here\""
echo "export LITELLM_BASE_URL=\"http://localhost:4000\""
echo "export LITELLM_MODEL=\"gpt-3.5-turbo\""

# Example 2: Using Anthropic via LiteLLM
echo ""
echo "Example 2: Anthropic via LiteLLM"
echo "--------------------------------"
echo "export LITELLM_API_KEY=\"sk-ant-your-anthropic-key-here\""
echo "export LITELLM_BASE_URL=\"http://localhost:4000\""
echo "export LITELLM_MODEL=\"claude-3-sonnet\""

# Example 3: Using local model via LiteLLM
echo ""
echo "Example 3: Local model via LiteLLM"
echo "----------------------------------"
echo "export LITELLM_API_KEY=\"not-needed\""
echo "export LITELLM_BASE_URL=\"http://localhost:4000\""
echo "export LITELLM_MODEL=\"llama2\""

echo ""
echo "🚀 Quick Start Steps:"
echo "===================="
echo ""
echo "1. Install LiteLLM:"
echo "   pip install litellm[proxy]"
echo ""
echo "2. Configure your models in litellm_config.yaml"
echo ""
echo "3. Start LiteLLM proxy:"
echo "   litellm --config litellm_config.yaml --port 4000"
echo ""
echo "4. Set environment variables (choose one example above)"
echo ""
echo "5. Run the quickstart script:"
echo "   ./quickstart.sh"
echo ""
echo "6. Run the agent:"
echo "   cd ../../.."
echo "   go run ./cmd/adk run examples/agents/litellm_agent"
echo ""
echo "💡 Tips:"
echo "- You can switch between different LLM providers by changing LITELLM_MODEL"
echo "- The same agent code works with any LLM provider supported by LiteLLM"
echo "- Use the web interface for a better experience: go run ./cmd/adk web examples/agents"
