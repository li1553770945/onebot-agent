import { Test, TestingModule } from '@nestjs/testing';
import { AppController, MessageController } from './app.controller';
import { AppService } from './app.service';

describe('AppController', () => {
  let app: TestingModule;

  beforeAll(async () => {
    app = await Test.createTestingModule({
      controllers: [AppController],
      providers: [AppService],
    }).compile();
  });

  describe('getHello', () => {
    it('should return "Hello World!"', () => {
      const appController = app.get(AppController);
      expect(appController.getHello()).toBe('Hello World!');
    });
  });

  describe('MessageController HandleMessage', () => {
    it('should echo json object when body is object', () => {
      const messageController = app.get(MessageController);
      const body = { a: 1, b: 'x' };
      const res = messageController.getMessage(body);
      expect(res.ok).toBe(true);
      expect(res.received).toEqual(body);
    });

    it('should parse json string body', () => {
      const messageController = app.get(MessageController);
      const jsonStr = '{"a":2,"b":"y"}';
      const res = messageController.getMessage(jsonStr);
      expect(res.received).toEqual({ a: 2, b: 'y' });
    });

    it('should keep plain string when not json', () => {
      const messageController = app.get(MessageController);
      const txt = 'plain-text';
      const res = messageController.getMessage(txt);
      expect(res.received).toBe(txt);
    });
  });
});
