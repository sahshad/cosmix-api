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
    @CurrentUser() user: { userId: number },
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
      Number(id),
    );
  }

  @UseGuards(AuthGuard)
  @Put(':id')
  updatePost(
    @Param('id') id: string,
    @CurrentUser() user: { userId: number },
    @Body() body: UpdatePostDto,
  ) {
    return this.postGrpc.updatePost(
      Number(id),
      user.userId,
      body,
    );
  }

  @UseGuards(AuthGuard)
  @Delete(':id')
  deletePost(
    @Param('id') id: string,
    @CurrentUser() user: { userId: number },
  ) {
    return this.postGrpc.deletePost(
      Number(id),
      user.userId,
    );
  }

  @UseGuards(AuthGuard)
  @Post(':id/like')
  likePost(
    @Param('id') id: string,
    @CurrentUser() user: { userId: number },
  ) {
    return this.postGrpc.likePost(
      Number(id),
      user.userId,
    );
  }

  @UseGuards(AuthGuard)
  @Delete(':id/like')
  unlikePost(
    @Param('id') id: string,
    @CurrentUser() user: { userId: number },
  ) {
    return this.postGrpc.unlikePost(
      Number(id),
      user.userId,
    );
  }

  @UseGuards(AuthGuard)
  @Post(':id/comment')
  createComment(
    @Param('id') id: string,
    @CurrentUser() user: { userId: number },
    @Body() body: CreateCommentDto,
  ) {
    return this.postGrpc.createComment(
      Number(id),
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
      Number(id),
      Number(page),
      Number(limit),
    );
  }

  @UseGuards(AuthGuard)
  @Put('comment/:commentId')
  updateComment(
    @Param('commentId') commentId: string,
    @CurrentUser() user: { userId: number },
    @Body() body: UpdateCommentDto,
  ) {
    return this.postGrpc.updateComment(
      Number(commentId),
      user.userId,
      body.content,
    );
  }

  @UseGuards(AuthGuard)
  @Delete('comment/:commentId')
  deleteComment(
    @Param('commentId') commentId: string,
    @CurrentUser() user: { userId: number },
  ) {
    return this.postGrpc.deleteComment(
      Number(commentId),
      user.userId,
    );
  }
}