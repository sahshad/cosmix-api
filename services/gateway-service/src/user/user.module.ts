import { Module } from '@nestjs/common';
import { UserController } from './user.controller';
import { UserService } from './user.service';
import { UserGrpcService } from './user-client.service';

@Module({
  controllers: [UserController],
  providers: [UserService, UserGrpcService]
})
export class UserModule {}
