<div align="center">

# 云原生架构适用于 onebot 协议的 Agent框架

[English](./docs/README_EN.md) | [中文](README.md)

</div>

## 为什么要做这个项目？

一句话总结：**把“消息入口”“决策能力”“工具调用”“协议发送”四件事切干净，再用云原生方式把弹性、可演进能力提前内建**。具体的背景、设计目标、适用场景等请参考[为什么要做这个项目](./docs/why.md)。

## 快速开始

请参考[快速开始](./docs/quickstart.md)，更详细的说明和部署文档请查看[部署文档](./docs/deployment.md)。


## 系统架构

```mermaid

graph LR
	onebot实现2 -->|收到消息推送到| 消息分发服务
	onebot实现1 --> 消息分发服务
	onebot实现... --> 消息分发服务
	消息分发服务 --> Agent1
	消息分发服务 -->|根据规则转发到对应的Agent| Agent2
	消息分发服务 --> Agent...
	Agent2 --> MCPServer1
	Agent2 -->|调用多种MCPServer| MCPServer2
	Agent2 --> MCPServer...
	Agent2 --> 消息发送MCPServer
	消息发送MCPServer --> 消息分发服务
	消息分发服务 -->|将消息发送请求调用对应onebot实现| onebot实现2

```

## 交互流程

```mermaid
sequenceDiagram
	onbot实现1（例如lagrange.onebot）->>消息分发服务: 1. onebot收到消息后推送
	消息分发服务->> Agent 1: 2.根据规则转发到对应的Agent
	Agent1 ->> LLM: 3. Agent调用LLM处理
	Agent1 ->> MCPServer1: 4. Agent调用MCPServer，例如联网搜索
	Agent1 ->> MessageSendMCPServer: 5. 调用发送消息的MCP
	MessageSendMCPServer ->> 消息分发服务: 6. 推送回消息分发服务
	消息分发服务 ->> onbot实现1（例如lagrange.onebot）: 7.调用对应的onebot实现发送消息
```

## 服务列表

## lagrange-core

lagrange.onebot可执行文件。

## message-dispatch 

消息分发以及接收发送服务。

## message-mcp

消息发送的MCP服务器。

