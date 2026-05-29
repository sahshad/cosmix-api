import { Module } from '@nestjs/common';
import { AuthController } from './auth.controller';
import { AuthService } from './auth.service';
import { AuthGrpcService } from './auth-client.service';
import { join } from 'path';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { AUTH_PACKAGE_NAME } from '../generated/auth/auth';

@Module({
  imports: [
    ClientsModule.register([
      {
        name: AUTH_PACKAGE_NAME,
        transport: Transport.GRPC,
        options: {
          url:process.env.AUTH_GRPC_ADDR ?? 'localhost:50051',
          package: 'auth',
          protoPath: join(
            __dirname,
            '../../../../shared/grpc/proto/auth/auth.proto',
          ),
        },
      },
    ]),
  ],
  controllers: [AuthController],
  providers: [AuthService, AuthGrpcService],
  exports: [AuthGrpcService],
})
export class AuthModule { }
