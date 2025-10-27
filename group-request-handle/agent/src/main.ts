import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';

// 由 ConfigModule (app.module.ts) 负责环境变量加载与校验
async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  await app.listen(3000);
}
bootstrap();
