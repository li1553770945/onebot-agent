import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { getRequestGroup } from './prompts';
import { MultiServerMCPClient } from "@langchain/mcp-adapters";
import { createReactAgent } from "@langchain/langgraph/prebuilt";
import { ChatOpenAI } from "@langchain/openai";

@Injectable()
export class AppService {
  private mcpClient: MultiServerMCPClient;
  private llm: ChatOpenAI;
  private agent: ReturnType<typeof createReactAgent>;
  constructor(private readonly config: ConfigService) {
    this.mcpClient = new MultiServerMCPClient({
      mcpServers: {
        "group-operator": {
          // Ensure your start your weather server on port 8000
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

  private async handleAdd(data: any) {
    const comment = data.comment;
    const userId = data.user_id;
    const selfId = data.self_id;
    const flag = data.flag;
    const groupId = data.group_id;
    console.log("收到入群申请，AI处理中...");
    const notifyGroup = this.config.get<string>('NOTIFY_GROUP');
    const prompt = getRequestGroup(selfId, userId, flag, groupId, comment, notifyGroup);
    const res = await this.agent.invoke({
      messages: [
        {
          role: "user",
          content: prompt,
        },
      ],
    });
    console.log("AI Response:", res);
    return res;
  }

  private async handleNormalMessage(data: any) {
    console.log("Handling normal message:", data);
    // 处理普通消息的逻辑，可使用其他配置: 如模型名称
    const model = this.config.get<string>('LLM_MODEL_NAME');
    // console.log('Using model:', model);
  }
}
