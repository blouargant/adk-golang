#!/bin/bash

# Quickstart script for LiteLLM Agent example
# This script helps you get started quickly with the LiteLLM agent

set -e

echo "🚀 LiteLLM Agent Quickstart"
echo "=========================="

# Check if we're in the right directory
if [ ! -f "../../../go.mod" ]; then
    echo "❌ Error: Please run this script from the examples/agents/litellm_agent directory"
    exit 1
fi

# remove agent.so if it exists
if [ -f "agent.so" ]; then
    echo "🗑️  Removing existing agent.so file"
    rm "agent.so"
fi

# Check for required environment variables
echo "🔍 Checking environment variables..."

if [ -z "$LITELLM_API_KEY" ]; then
    echo "❌ LITELLM_API_KEY environment variable is not set"
    echo "   Please set it with: export LITELLM_API_KEY=\"your-api-key\""
    echo ""
    echo "   Example for OpenAI via LiteLLM:"
    echo "   export LITELLM_API_KEY=\"sk-your-openai-key\""
    echo ""
    echo "   Example for local model:"
    echo "   export LITELLM_API_KEY=\"not-needed\""
    exit 1
fi

if [ -z "$LITELLM_BASE_URL" ]; then
    echo "❌ LITELLM_BASE_URL environment variable is not set"
    echo "   Please set it with: export LITELLM_BASE_URL=\"http://localhost:4000\""
    echo ""
    echo "   This should point to your LiteLLM proxy server."
    echo "   Start LiteLLM proxy with: litellm --config your_config.yaml --port 4000"
    exit 1
fi

# Set default model if not specified
if [ -z "$LITELLM_MODEL" ]; then
    echo "⚠️  LITELLM_MODEL not set, using default: gpt-3.5-turbo"
    export LITELLM_MODEL="gpt-3.5-turbo"
fi

echo "✅ Environment variables:"
echo "   LITELLM_API_KEY: [set]"
echo "   LITELLM_BASE_URL: $LITELLM_BASE_URL"
echo "   LITELLM_MODEL: $LITELLM_MODEL"

# Test LiteLLM proxy connectivity
echo ""
echo "🔗 Testing LiteLLM proxy connectivity..."
if curl -s -f -o /dev/null "$LITELLM_BASE_URL/health/liveness" 2>/dev/null; then
    echo "✅ LiteLLM proxy is reachable at $LITELLM_BASE_URL"
else
    echo "⚠️  Warning: Cannot reach LiteLLM proxy at $LITELLM_BASE_URL"
    echo "   Make sure your LiteLLM proxy is running:"
    echo "   pip install litellm[proxy]"
    echo "   litellm --config your_config.yaml --port 4000"
    echo ""
    echo "   Continuing anyway (the agent will fail if proxy is not available)..."
fi

# Build the agent
echo ""
echo "🔨 Building LiteLLM agent..."
go mod tidy
go build -buildmode=plugin -o agent.so .
echo ""

if [ $? -eq 0 ]; then
    echo "✅ LiteLLM agent built successfully!"
    echo "🔧 LiteLLM Configuration:"
    echo "   Provider: LiteLLM Proxy"
    echo "   Base URL: $LITELLM_BASE_URL"
    echo "   Model: $LITELLM_MODEL"
    echo "   Authentication: Bearer Token"
else
    echo "❌ Failed to build LiteLLM agent"
    exit 1
fi

# Go to the root directory and run the agent
cd ../../..

echo ""
echo "Choose how to run the agent:"
echo "1) Interactive CLI (recommended)"
echo "2) Web interface"
echo ""
read -p "Enter your choice (1 or 2): " -n 1 -r
echo

if [[ $REPLY == "2" ]]; then
    echo "🌐 Starting web interface..."
    echo "Open http://localhost:8080 in your browser"
    echo "Press Ctrl+C to stop"
    go run ./cmd/adk web examples/agents
else
    echo "💬 Starting interactive CLI..."
    echo "Type 'exit' or 'quit' to stop"
    echo "Try asking: 'What's the latest news about artificial intelligence?'"
    echo ""
    go run ./cmd/adk run examples/agents/openai_agent
fi