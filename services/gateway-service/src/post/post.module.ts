import { Module } from '@nestjs/common';
import { PostController } from './post.controller';
import { PostGrpcService } from './post-client.service';
import { PostService } from './post.service';

@Module({
  controllers: [PostController],
  providers: [PostGrpcService, PostService]
})
export class PostModule {}
