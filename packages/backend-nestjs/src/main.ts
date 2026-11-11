import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { SwaggerModule, DocumentBuilder } from '@nestjs/swagger';
import * as fs from 'fs';
import * as path from 'path';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);

  // CORS for KMP frontend (Android + web)
  app.enableCors({
    origin: ['http://localhost:8081', 'http://localhost:3000', 'http://10.0.2.2:8080'],
    credentials: true,
  });

  // Swagger/OpenAPI documentation
  const config = new DocumentBuilder()
    .setTitle('AREA API')
    .setDescription('Area automation platform API - OpenAPI spec for Kotlin Multiplatform client generation')
    .setVersion('1.0')
    .addTag('area')
    .addBearerAuth({
      type: 'http',
      scheme: 'bearer',
      bearerFormat: 'JWT',
      name: 'Authorization',
      description: 'Enter JWT token',
    }, 'bearer')
    .build();

  const document = SwaggerModule.createDocument(app, config);
  
  // Write OpenAPI JSON to file for client generation
  const openapiPath = path.join(__dirname, '..', 'openapi.json');
  fs.writeFileSync(openapiPath, JSON.stringify(document, null, 2));

  SwaggerModule.setup('api', app, document);

  const port = process.env.PORT ?? 8080;
  await app.listen(port);
  console.log(`API running on http://localhost:${port}`);
  console.log(`OpenAPI spec at http://localhost:${port}/api-json`);
  console.log(`OpenAPI JSON written to ${openapiPath}`);
}
bootstrap();
