import { Module } from '@nestjs/common';
import { UserController } from './user.controller';
import { UserService } from './user.service';
import { UserGrpcService } from './user-client.service';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { USER_PACKAGE_NAME } from '../generated/user/user';
import { join } from 'path';

@Module({
  imports: [
    ClientsModule.register([
      {
        name: USER_PACKAGE_NAME,
        transport: Transport.GRPC,
        options: {
          url: process.env.USER_GRPC_ADDR ?? 'localhost:50052',
          package: 'user',
          protoPath: join(
            __dirname,
            '../../../../shared/grpc/proto/user/user.proto',
          ),
        },
      },
    ]),
  ],
  controllers: [UserController],
  providers: [UserService, UserGrpcService],
  exports: [UserGrpcService],
})
export class UserModule { }
