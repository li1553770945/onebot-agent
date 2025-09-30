import { Controller, Get, Post, Body } from '@nestjs/common';
import { AppService } from './app.service';

@Controller("ping")
export class AppController {
  constructor(private readonly appService: AppService) {}

  @Get()
  getHello(): string {
    return this.appService.getHello();
  }
}

@Controller("message")
export class MessageController {
  constructor(private readonly appService: AppService) {}

  // 接收原始 body（对象或字符串），传递给服务层
  @Post()
  getMessage(@Body() body: any): any {
    return this.appService.HandleMessage(body);
  }
}