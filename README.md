# 云原生架构适用于 onebot 协议的 Agent框架

## 系统架构

graph TD
    A[Hard] -->|Text| B(Round)
    B --> C{Decision}
    C -->|One| D[Result 1]
    C -->|Two| E[Result 2]

## 服务列表

## lagrange-core

lagrange.onebot可执行文件。

## message-dispatch 

消息分发以及接收发送服务。

## message-mcp

消息发送的MCP服务器。

