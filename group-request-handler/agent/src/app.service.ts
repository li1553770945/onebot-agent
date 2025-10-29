import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { groupAuditConfigs, GroupAuditConfig } from './prompts';
import { MultiServerMCPClient } from "@langchain/mcp-adapters";
import { createReactAgent } from "@langchain/langgraph/prebuilt";
import { ChatOpenAI } from "@langchain/openai";
import { INJECTION_TOKENS } from './prompts';
@Injectable()
export class AppService {
  private mcpClient: MultiServerMCPClient;
  private llm: ChatOpenAI;
  private agent: ReturnType<typeof createReactAgent>;
  constructor(private readonly config: ConfigService) {
    this.mcpClient = new MultiServerMCPClient({
      mcpServers: {
        "group-operator": {
          url: this.config.get<string>('MCP_URL'),
          transport: "http",
        }
      }
    })
    this.llm = new ChatOpenAI({
      modelName: this.config.get<string>('LLM_MODEL_NAME'),
      apiKey: this.config.get<string>('LLM_API_KEY'),
      configuration: this.config.get<string>('LLM_API_BASE_URL') ? { baseURL: this.config.get<string>('LLM_API_BASE_URL') } : undefined,

    });
  }
  async onModuleInit() {
    // 3. 在生命周期钩子中执行异步操作
    await this.initializeAgent();
  }

  private async initializeAgent() {
    const tools = await this.mcpClient.getTools();
    this.agent = createReactAgent({
      llm: this.llm,
      tools,
    });
  }

  getHello(): string {
    return 'Hello World!';
  }
  async HandleMessage(payload: any): Promise<any> {
    // 如果是字符串尝试解析为 JSON
    let data = payload;
    if (typeof payload === 'string') {
      try {
        data = JSON.parse(payload);
      } catch (e) {
        // 不是合法 JSON，保持原样
      }
    }
    // 如果是 Buffer
    if (payload instanceof Buffer) {
      try {
        const text = payload.toString('utf8');
        data = JSON.parse(text);
      } catch (e) {
        data = payload.toString('utf8');
      }
    }
    // 返回解析后的对象或原始内容
    console.log("Received message:", data);
    const post_type = data.post_type;
    const sub_type = data.sub_type;
    console.log("Post type:", post_type);
    console.log("Sub type:", sub_type);
    let res;
    if (post_type === 'request' && sub_type === "add") {
      res = await this.handleAdd(data);
    } else if (post_type === "message" && sub_type === "normal") {
      res = await this.handleNormalMessage(data);
    }
    return {
      data: res,
    };
  }

  private async handle(selfId: string, userId: string, flag: string, config: GroupAuditConfig, comment: string) {
    const res = await this.agent.invoke({
      messages: [
        {
          role: "user",
          content: config.getPrompt(selfId, userId, flag, config.groupId, comment, config.notifyGroupId),
        },
      ],
    });
    console.log("AI Response:", res);
    return res;
  }

  private async handleAdd(data: any) {
    const groupId = data.group_id;

    console.log(`收到${groupId}入群申请，AI处理中...`);
    
    // 查找匹配的群审核配置
    for (const config of groupAuditConfigs) {
      if (groupId === config.groupId) {
        const comment = data.comment;
        const userId = data.user_id;
        const selfId = data.self_id;
        const flag = data.flag;
        
        // 先进行格式校验
        const validRes = await this.checkCommentValid(comment);

        if (!validRes.valid) {
          // 备注不合法，直接拒绝
          const tools = await this.mcpClient.getTools();
          const rejectTool = tools.find(tool => tool.name === 'reject_add_request');
          if (rejectTool) {
            const result = await rejectTool.invoke({
              self_id: selfId,
              flag: flag,
              reason: validRes.reason
            });
            return result;
          }
          
          // 通知管理群
          const send_group_message = tools.find(tool => tool.name === 'send_group_message');
          if (send_group_message) {
            const result = await send_group_message.invoke({
              self_id: selfId,
              group_id: config.notifyGroupId,
              message: `AI入群审核已拒绝，请人工知悉。用户ID: ${userId}，入群备注: ${comment}，拒绝理由: ${validRes.reason}`
            });
            return result;
          }
          return;
        }

        // 备注合法，调用 AI 进行审核
        return await this.handle(selfId, userId, flag, config, comment);
      }
    }
    
    // 循环结束都没找到匹配的群
    console.log(`未找到群 ${groupId} 的配置，忽略此请求`);
    return { message: "未找到符合条件的群聊配置，忽略此请求" };
  }
  private async checkCommentValid(comment: string): Promise<{
      valid: boolean,
      reason?: string
    }> {
    if (!comment) {
      return { valid: false, reason: "备注内容为空" };
    }
    if (INJECTION_TOKENS.some(tok => comment.includes(tok))) {
      return { valid: false, reason: "备注包含指令性内容（疑似提示注入）" };
    }
    return { valid: true };


  }
  private async handleNormalMessage(data: any) {
    console.log("Handling normal message:", data);
    // 处理普通消息的逻辑，可使用其他配置: 如模型名称
    const model = this.config.get<string>('LLM_MODEL_NAME');
    // console.log('Using model:', model);
  }
}
