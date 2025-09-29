我正在编写一个云原生架构的适用于 onebot 协议的 Agent 框架，包含以下组件：
- lagrange.onebot 的 Docker 镜像和部署脚本
- 消息分发服务message-dispatch，负责使用httpserver接收onebot的消息，根据特定规则使用grpc将消息发送到特定的Agent，并且自己也是一个grpc服务，用于将Agent处理后的消息发送请求转发给对应的onebot实现
- Agent，负责处理消息，可以调用多个MCP服务器来处理消息
- 多种MCP 服务器，用于给Agent提供联网搜索、代码执行等能力
- 消息发送MCP服务器，负责将Agent处理后的消息发送请求转发给消息分发服务，由消息分发服务调用对应的onebot实现来发送消息

