import express, { Request, Response } from "express";
import { z } from "zod";
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StreamableHTTPServerTransport } from "@modelcontextprotocol/sdk/server/streamableHttp.js";
import axios from "axios";

// 1) 创建 MCP 服务器，仅注册一个 "sum" 工具
const server = new McpServer({
  name: "sum-server",
  version: "1.0.0",
});

server.registerTool(
  "approve_add_request",
  {
    description: "同意加群请求",
    inputSchema: {
      self_id: z.string().describe("自己的用户id"),
      flag: z.string().describe("加群请求的标识"),
    },
  },
  async ({ self_id, flag }) => {
    console.log(`Approving add request for flag: ${flag}, self_id: ${self_id}`);
    const body = {
      approve: true,
      flag:flag
    }
    const result = await axios.post('http://lagrange-onebot-service:15000/set_group_add_request', body);
    return {
      content: [
        { type: "text", text: typeof result.data === "string" ? result.data : JSON.stringify(result.data) },
      ],
    };
  }
);

server.registerTool(
  "reject_add_request",
  {
    description: "拒绝加群请求",
    inputSchema: {
      self_id: z.string().describe("自己的用户id"),
      flag: z.string().describe("加群请求的标识"),
      reason: z.string().describe("拒绝理由").optional(),
    },
  },
  async ({ self_id, flag, reason }) => {
    console.log(`Rejecting add request for flag: ${flag}, self_id: ${self_id}`);
    const body = {
      approve: false,
      flag:flag,
      reason: reason
    }
    const result = await axios.post('http://lagrange-onebot-service:15000/set_group_add_request', body);
    return {
      content: [
        { type: "text", text: typeof result.data === "string" ? result.data : JSON.stringify(result.data) },
      ],
    };
  }
);

server.registerTool(
  "send_group_message",
  {
    description: "发送群消息",
    inputSchema: {
      self_id: z.string().describe("自己的用户id"),
      group_id: z.string().describe("群组id"),
      message: z.string().describe("消息内容"),
    },
  },
  async ({ self_id, group_id, message }) => {
    const body =    {
      action: "send_message",
      params: {
          detail_type: "group",
          group_id: group_id,
          self_id: self_id,
          message: [
              {
                "type": "text",
                "data": {
                    "text": message
                }
            }
        ]
    }};
    const bodyStr = JSON.stringify(body);
    const result =  await axios.post('http://message-dispatch:15001/send', bodyStr);
    return {
      content: [
         { type: "text", text: typeof result.data === "string" ? result.data : JSON.stringify(result.data) },
      ],
    };
  }
);



// 2) 用 Express 暴露一个 MCP 端点（/mcp），采用 Streamable HTTP（无会话模式）
const app = express();
app.use(express.json());

// 简单健康检查：GET /ping
app.get("/ping", (_req: Request, res: Response) => {
  res.json({ status: "ok", timestamp: new Date().toISOString() });
});

// 仅保留一个端点：/mcp
app.all("/mcp", async (req: Request, res: Response) => {
  try {
    const transport = new StreamableHTTPServerTransport({
      sessionIdGenerator: undefined,
    });

    await server.connect(transport);
    await transport.handleRequest(req, res, (req as any).body);

    // 请求关闭时清理
    res.on("close", () => {
      try {
        transport.close();
        server.close();
      } catch (_) {}
    });
  } catch (err) {
    console.error(err);
    if (!res.headersSent) {
      res.status(500).json({ error: "Internal Server Error" });
    }
  }
});

// 3) 启动在 3000 端口
const PORT = 3000;
app.listen(PORT, () => {
  console.log(`MCP Streamable HTTP server listening on :${PORT}/mcp`);
});
