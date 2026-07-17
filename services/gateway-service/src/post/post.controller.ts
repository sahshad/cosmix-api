import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Post,
  Put,
  Query,
  UseGuards,
} from '@nestjs/common';

import { PostGrpcService } from './post-client.service';

import { AuthGuard } from '../common/guard/auth.guard';
import { CurrentUser } from '../common/decorators/current-user.decorator';

import { CreatePostDto } from './dto/create-post.dto';
import { UpdatePostDto } from './dto/update-post.dto';
import { CreateCommentDto } from './dto/create-comment.dto';
import { UpdateCommentDto } from './dto/update-comment.dto';

@Controller('posts')
export class PostController {
  constructor(
    private readonly postGrpc: PostGrpcService,
  ) {}

  @UseGuards(AuthGuard)
  @Post()
  createPost(
    @CurrentUser() user: { userId: string },
    @Body() body: CreatePostDto,
  ) {
    return this.postGrpc.createPost(
      user.userId,
      body,
    );
  }

  @Get()
  getFeed(
    @Query('page') page = '1',
    @Query('limit') limit = '20',
  ) {
    return this.postGrpc.getFeed(
      Number(page),
      Number(limit),
    );
  }

  @Get(':id')
  getPost(
    @Param('id') id: string,
  ) {
    return this.postGrpc.getPost(
      id,
    );
  }

  @UseGuards(AuthGuard)
  @Put(':id')
  updatePost(
    @Param('id') id: string,
    @CurrentUser() user: { userId: string },
    @Body() body: UpdatePostDto,
  ) {
    return this.postGrpc.updatePost(
      id,
      user.userId,
      body,
    );
  }

  @UseGuards(AuthGuard)
  @Delete(':id')
  deletePost(
    @Param('id') id: string,
    @CurrentUser() user: { userId: string },
  ) {
    return this.postGrpc.deletePost(
      id,
      user.userId,
    );
  }

  @UseGuards(AuthGuard)
  @Post(':id/like')
  likePost(
    @Param('id') id: string,
    @CurrentUser() user: { userId: string },
  ) {
    return this.postGrpc.likePost(
      id,
      user.userId,
    );
  }

  @UseGuards(AuthGuard)
  @Delete(':id/like')
  unlikePost(
    @Param('id') id: string,
    @CurrentUser() user: { userId: string },
  ) {
    return this.postGrpc.unlikePost(
      id,
      user.userId,
    );
  }

  @UseGuards(AuthGuard)
  @Post(':id/comment')
  createComment(
    @Param('id') id: string,
    @CurrentUser() user: { userId: string },
    @Body() body: CreateCommentDto,
  ) {
    return this.postGrpc.createComment(
      id,
      user.userId,
      body.content,
    );
  }

  @Get(':id/comment')
  getComments(
    @Param('id') id: string,
    @Query('page') page = '1',
    @Query('limit') limit = '10',
  ) {
    return this.postGrpc.getComments(
      id,
      Number(page),
      Number(limit),
    );
  }

  @UseGuards(AuthGuard)
  @Put('comment/:commentId')
  updateComment(
    @Param('commentId') commentId: string,
    @CurrentUser() user: { userId: string },
    @Body() body: UpdateCommentDto,
  ) {
    return this.postGrpc.updateComment(
      commentId,
      user.userId,
      body.content,
    );
  }

  @UseGuards(AuthGuard)
  @Delete('comment/:commentId')
  deleteComment(
    @Param('commentId') commentId: string,
    @CurrentUser() user: { userId: string },
  ) {
    return this.postGrpc.deleteComment(
      commentId,
      user.userId,
    );
  }
}