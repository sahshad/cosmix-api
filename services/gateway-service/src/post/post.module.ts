import { Module } from '@nestjs/common';
import { PostController } from './post.controller';
import { PostGrpcService } from './post-client.service';
import { PostService } from './post.service';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { POST_PACKAGE_NAME } from '../generated/post/post';
import { join } from 'path';

@Module({
  imports: [
    ClientsModule.register([
      {
        name: POST_PACKAGE_NAME,
        transport: Transport.GRPC,
        options: {
          url: process.env.POST_GRPC_ADDR ?? 'localhost:50053',
          package: 'post',
          protoPath: join(
            __dirname,
            '../../../../shared/grpc/proto/post/post.proto',
          ),
        },
      },
    ]),
  ],
  controllers: [PostController],
  providers: [PostGrpcService, PostService],
  exports: [PostGrpcService],
})
export class PostModule { }
