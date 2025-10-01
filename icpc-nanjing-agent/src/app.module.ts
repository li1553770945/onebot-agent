import { Module } from '@nestjs/common';
import { AppController, MessageController } from './app.controller';
import { AppService } from './app.service';

@Module({
  imports: [],
  // 需要把 MessageController 也注册进来，否则其装饰的 /message 路由不会被扫描
  controllers: [AppController, MessageController],
  providers: [AppService],
})
export class AppModule {}
