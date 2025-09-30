import { Injectable } from '@nestjs/common';

@Injectable()
export class AppService {
  getHello(): string {
    return 'Hello World!';
  }
  HandleMessage(payload: any): any {
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
    return {
      ok: true,
      received: data,
    };
  }
}
