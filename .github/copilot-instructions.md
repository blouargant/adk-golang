# ADK-Golang Coding Instructions

ADK (Agent Development Kit) is a Go implementation of a framework for building AI agents, compatible with the [Agent2Agent (A2A) protocol](https://a2aproject.github.io/A2A/latest/specification/).

## Core Architecture

**Event-Driven Agent System**: All agents implement `BaseAgent` interface with async execution via channels (`EventStream <-chan *Event`). The `Runner` orchestrates execution with real-time streaming.

**Key Components**:
- **Agents**: `CustomAgent` (base), `LLMAgent`, `SequentialAgent`, `RemoteA2aAgent`
- **Tools**: `FunctionTool` with reflection-based auto-binding, streaming tools in `tools/async/`
- **Sessions**: Scoped state management with `InMemorySessionService`/`FileSessionService`
- **A2A Integration**: Full protocol compliance in `pkg/a2a/` with `Task`, `Message`, `AgentCard` objects

## Critical Patterns

**Interface Design**: Small, focused interfaces with context propagation:
```go
type BaseAgent interface {
    RunAsync(invocationCtx *InvocationContext) (EventStream, error)
    // Always use InvocationContext, never plain context.Context
}
```

**Pointer Utilities**: Use `pkg/ptr` for optional fields:
- `ptr.Float32(0.7)` for float32 pointers
- `ptr.Ptr(value)` for generic type pointers

**Event Streaming**: Agents communicate via channels, not direct calls:
```go
// Good: Channel-based event streaming
eventChan := make(chan *core.Event, 100)
go agent.RunAsync(ctx, eventChan)

// Bad: Direct synchronous calls
```

**Testing**: Mock LLM connections with `MockLLMConnection`, test async patterns with goroutines and channels.

## Development Workflow

**Build & Test**:
```bash
go install ./...          # Build CLI and tools
go test ./...             # Run all tests
go run ./cmd/adk web examples/agents  # Start web UI
```

**CLI Usage**:
- `adk run <agent_path>` - Interactive CLI with agent
- `adk web <agents_dir>` - Web UI for testing agents  
- `adk create <name>` - Scaffold new agent

## A2A Protocol Compliance

**Kind Fields**: Always set correct `Kind` values:
- `Task`: `"task"`
- `TaskStatusUpdateEvent`: `"status-update"`  
- `TaskArtifactUpdateEvent`: `"artifact-update"`
- `Message`: `"message"`

**Agent Registration**: Agents must provide `AgentCard` with capabilities, skills, and authentication schemes.

## Project-Specific Conventions

**Error Handling**: Always return explicit errors, use `fmt.Errorf` for context wrapping.

**Concurrency**: Use goroutines/channels over mutexes where possible. All session services are thread-safe.

**LLM Integration**: Use `pkg/llmconnect/ollama` for local models, follow the `LLMConnection` interface for new providers.

**Agent Hierarchies**: Use `SubAgents()` and `SetParentAgent()` for composition, `FindAgent()` for lookup.

Reference `examples/agents/llm_agent/` for canonical agent implementation patterns.
