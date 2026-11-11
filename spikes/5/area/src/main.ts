import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module.js';

async function bootstrap() {
  const port = process.env.PORT;
  if (!port) {
    throw new Error('PORT environment variable is required');
  }

  const app = await NestFactory.create(AppModule);
  app.enableCors();
  await app.listen(port);
}

bootstrap();