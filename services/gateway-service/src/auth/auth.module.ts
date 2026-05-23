import { Module } from '@nestjs/common';
import { AuthController } from './auth.controller';
import { AuthService } from './auth.service';
import { AuthGrpcService } from './auth-client.service';

@Module({
  controllers: [AuthController],
  providers: [AuthService, AuthGrpcService]
})
export class AuthModule {}
