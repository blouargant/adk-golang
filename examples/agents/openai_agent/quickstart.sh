#!/bin/bash

# Quickstart script for OpenAI Agent example
# This script helps you get started quickly with the OpenAI agent

set -e

echo "🚀 OpenAI Agent Quickstart"
echo "=========================="

# Check if we're in the right directory
if [ ! -f "../../../go.mod" ]; then
    echo "❌ Error: Please run this script from the examples/agents/openai_agent directory"
    exit 1
fi

# remove agent.so if it exists
if [ -f "agent.so" ]; then
    echo "🗑️  Removing existing agent.so file"
    rm "agent.so"
fi

# Check if OpenAI API key is set
if [ -z "$OPENAI_API_KEY" ]; then
    echo "⚠️  OpenAI API key not found in environment variables."
    echo "Please set your OpenAI API key:"
    echo "export OPENAI_API_KEY='your-api-key-here'"
    echo ""
    echo "You can get an API key from: https://platform.openai.com/api-keys"
    echo ""
    read -p "Do you want to set it now? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        read -p "Enter your OpenAI API key: " api_key
        export OPENAI_API_KEY="$api_key"
        echo "✅ API key set for this session"
    else
        echo "❌ Cannot proceed without API key"
        exit 1
    fi
else
    echo "✅ OpenAI API key found"
fi

# Set default model if not specified
if [ -z "$OPENAI_MODEL" ]; then
    export OPENAI_MODEL="gpt-3.5-turbo"
    echo "📝 Using default model: $OPENAI_MODEL"
else
    echo "📝 Using model: $OPENAI_MODEL"
fi

echo ""
echo "🏗️  Building and running OpenAI agent..."
echo ""

# Go to the root directory and run the agent
cd ../../..

# Option 1: Interactive CLI
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
