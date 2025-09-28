# Cloud-Native Architecture Agent Framework for Onebot Protocol


English | [中文](./docs/README_CN.md) |


## System Architecture

```mermaid
graph LR
	onebot-implementation-2 -->|Receives message and pushes to| message-dispatch-service
	onebot-implementation-1 --> message-dispatch-service
	onebot-implementation-n --> message-dispatch-service
	message-dispatch-service --> agent-1
	message-dispatch-service -->|Forwards to corresponding Agent based on rules| agent-2
	message-dispatch-service --> agent-n
	agent-2 --> mcpserver-1
	agent-2 -->|Calls multiple MCPServers| mcpserver-2
	agent-2 --> mcpserver-n
	agent-2 --> message-sending-mcpserver
	message-sending-mcpserver --> message-dispatch-service
	message-dispatch-service -->|Calls corresponding onebot implementation to send message| onebot-implementation-2
```

## Interaction Flow

```mermaid
sequenceDiagram
	onebot implementation 1 (e.g. lagrange.onebot) ->> Message Dispatch Service: 1. onebot receives message and pushes
	Message Dispatch Service ->> Agent 1: 2. Forwards to corresponding Agent based on rules
	Agent 1 ->> LLM: 3. Agent calls LLM for processing
	Agent 1 ->> MCPServer 1: 3. Agent calls MCPServer, e.g. for web search
	Agent 1 ->> MessageSendMCPServer: Calls MCP to send message
	MessageSendMCPServer ->> Message Dispatch Service: Pushes back to Message Dispatch Service
	Message Dispatch Service ->> onebot implementation 1 (e.g. lagrange.onebot): Calls corresponding onebot implementation to send message
```

## Service List

### lagrange-core

Executable file for lagrange.onebot.

### message-dispatch

Message dispatch, receiving, and sending service.

### message-mcp

MCP server for sending messages.