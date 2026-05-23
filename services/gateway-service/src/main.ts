import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { HttpExceptionFilter } from './common/filter/http-exception.filter';
import cookieParser from 'cookie-parser';
import { ValidationPipe } from '@nestjs/common';
import helmet from 'helmet';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  app.use(helmet())
  app.useGlobalFilters(new HttpExceptionFilter());
  app.use(cookieParser());
  app.setGlobalPrefix('api');
  app.useGlobalPipes(
    new ValidationPipe({
      whitelist: true,
      transform: true,
      forbidNonWhitelisted: true,
    }),
  );
  app.enableCors({
    origin:
      'http://localhost:3000',
    credentials: true,
  });
  await app.listen(process.env.PORT ?? 5000);
}
bootstrap();
