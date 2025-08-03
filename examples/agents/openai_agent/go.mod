module openai_agent

go 1.24.0

replace github.com/agent-protocol/adk-golang => ../../..

require github.com/agent-protocol/adk-golang v0.0.0-00010101000000-000000000000

require (
	github.com/openai/openai-go v1.12.0 // indirect
	github.com/tidwall/gjson v1.14.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
)
