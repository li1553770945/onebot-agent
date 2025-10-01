import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import * as Joi from 'joi';
import { AppController, MessageController } from './app.controller';
import { AppService } from './app.service';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      envFilePath: ['.env', '.env.local'],
      validationSchema: Joi.object({
        LLM_MODEL_NAME: Joi.string().min(1).required(),
        LLM_API_BASE_URL: Joi.string().uri().required(),
        LLM_API_KEY: Joi.string().min(1).required(),
        NOTIFY_GROUP: Joi.string().min(1).required(),
      }),
      validationOptions: {
        abortEarly: false,
        allowUnknown: true,
      },
    }),
  ],
  // 需要把 MessageController 也注册进来，否则其装饰的 /message 路由不会被扫描
  controllers: [AppController, MessageController],
  providers: [AppService],
})
export class AppModule {}
