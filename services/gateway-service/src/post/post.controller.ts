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
import {
  ApiBadRequestResponse,
  ApiBearerAuth,
  ApiCreatedResponse,
  ApiNotFoundResponse,
  ApiOkResponse,
  ApiOperation,
  ApiParam,
  ApiQuery,
  ApiTags,
  ApiUnauthorizedResponse,
} from '@nestjs/swagger';

import { PostGrpcService } from './post-client.service';

import { AuthGuard } from '../common/guard/auth.guard';
import { OptionalAuthGuard } from '../common/guard/optional-auth.guard';
import { CurrentUser } from '../common/decorators/current-user.decorator';

import { CreatePostDto } from './dto/create-post.dto';
import { UpdatePostDto } from './dto/update-post.dto';
import { CreateCommentDto } from './dto/create-comment.dto';
import { UpdateCommentDto } from './dto/update-comment.dto';

@ApiTags('posts')
@Controller('posts')
export class PostController {
  constructor(private readonly postGrpc: PostGrpcService) {}

  @ApiBearerAuth('access-token')
  @ApiOperation({ summary: 'Create a post' })
  @ApiCreatedResponse({ description: 'Created post' })
  @ApiBadRequestResponse({ description: 'Invalid post data' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Post()
  createPost(
    @CurrentUser() user: { userId: string },
    @Body() body: CreatePostDto,
  ) {
    return this.postGrpc.createPost(user.userId, body);
  }

  @ApiOperation({
    summary: 'Get the paginated post feed',
    description:
      'Authentication is optional; send a Bearer token for personalized like state.',
  })
  @ApiQuery({
    name: 'page',
    required: false,
    description: 'Page number (1-based)',
    example: 1,
    type: Number,
  })
  @ApiQuery({
    name: 'limit',
    required: false,
    description: 'Number of posts per page',
    example: 20,
    type: Number,
  })
  @ApiOkResponse({ description: 'Paginated post feed' })
  @UseGuards(OptionalAuthGuard)
  @Get()
  getFeed(
    @CurrentUser() user: { userId: string } | undefined,
    @Query('page') page = '1',
    @Query('limit') limit = '20',
  ) {
    return this.postGrpc.getFeed(Number(page), Number(limit), user?.userId);
  }

  @ApiOperation({
    summary: 'Get the paginated posts of a user',
    description:
      'Authentication is optional; send a Bearer token for personalized like state.',
  })
  @ApiParam({
    name: 'userId',
    description: 'Author user ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiQuery({
    name: 'page',
    required: false,
    description: 'Page number (1-based)',
    example: 1,
    type: Number,
  })
  @ApiQuery({
    name: 'limit',
    required: false,
    description: 'Number of posts per page',
    example: 20,
    type: Number,
  })
  @ApiOkResponse({ description: 'Paginated list of user posts' })
  @UseGuards(OptionalAuthGuard)
  @Get('user/:userId')
  getUserPosts(
    @Param('userId') userId: string,
    @CurrentUser() user: { userId: string } | undefined,
    @Query('page') page = '1',
    @Query('limit') limit = '20',
  ) {
    return this.postGrpc.getUserPosts(
      userId,
      Number(page),
      Number(limit),
      user?.userId,
    );
  }

  @ApiOperation({ summary: 'Get a post by ID' })
  @ApiParam({
    name: 'id',
    description: 'Post ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiOkResponse({ description: 'Post details' })
  @ApiNotFoundResponse({ description: 'Post not found' })
  @Get(':id')
  getPost(@Param('id') id: string) {
    return this.postGrpc.getPost(id);
  }

  @ApiBearerAuth('access-token')
  @ApiOperation({ summary: 'Update a post' })
  @ApiParam({
    name: 'id',
    description: 'Post ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiOkResponse({ description: 'Updated post' })
  @ApiBadRequestResponse({ description: 'Invalid post data' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Put(':id')
  updatePost(
    @Param('id') id: string,
    @CurrentUser() user: { userId: string },
    @Body() body: UpdatePostDto,
  ) {
    return this.postGrpc.updatePost(id, user.userId, body);
  }

  @ApiBearerAuth('access-token')
  @ApiOperation({ summary: 'Delete a post' })
  @ApiParam({
    name: 'id',
    description: 'Post ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiOkResponse({ description: 'Deleted post' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Delete(':id')
  deletePost(@Param('id') id: string, @CurrentUser() user: { userId: string }) {
    return this.postGrpc.deletePost(id, user.userId);
  }

  @ApiBearerAuth('access-token')
  @ApiOperation({ summary: 'Like a post' })
  @ApiParam({
    name: 'id',
    description: 'Post ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiCreatedResponse({ description: 'Post liked' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Post(':id/like')
  likePost(@Param('id') id: string, @CurrentUser() user: { userId: string }) {
    return this.postGrpc.likePost(id, user.userId);
  }

  @ApiBearerAuth('access-token')
  @ApiOperation({ summary: 'Remove a like from a post' })
  @ApiParam({
    name: 'id',
    description: 'Post ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiOkResponse({ description: 'Post like removed' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Delete(':id/like')
  unlikePost(@Param('id') id: string, @CurrentUser() user: { userId: string }) {
    return this.postGrpc.unlikePost(id, user.userId);
  }

  @ApiBearerAuth('access-token')
  @ApiOperation({ summary: 'Create a comment on a post' })
  @ApiParam({
    name: 'id',
    description: 'Post ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiCreatedResponse({ description: 'Created comment' })
  @ApiBadRequestResponse({ description: 'Invalid comment data' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
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
      body.parentCommentId,
    );
  }

  @ApiOperation({
    summary: 'Get the paginated comments of a post',
    description:
      'Authentication is optional; send a Bearer token to include viewer-specific like state.',
  })
  @ApiParam({
    name: 'id',
    description: 'Post ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiQuery({
    name: 'page',
    required: false,
    description: 'Page number (1-based)',
    example: 1,
    type: Number,
  })
  @ApiQuery({
    name: 'limit',
    required: false,
    description: 'Number of comments per page',
    example: 10,
    type: Number,
  })
  @ApiOkResponse({ description: 'Paginated comments' })
  @UseGuards(OptionalAuthGuard)
  @Get(':id/comment')
  getComments(
    @Param('id') id: string,
    @CurrentUser() user: { userId: string } | undefined,
    @Query('page') page = '1',
    @Query('limit') limit = '10',
  ) {
    return this.postGrpc.getComments(
      id,
      Number(page),
      Number(limit),
      user?.userId,
    );
  }

  @ApiOperation({
    summary: 'Get the paginated replies to a comment',
    description:
      'Authentication is optional; send a Bearer token to include viewer-specific like state.',
  })
  @ApiParam({
    name: 'commentId',
    description: 'Parent comment ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiQuery({
    name: 'page',
    required: false,
    description: 'Page number (1-based)',
    example: 1,
    type: Number,
  })
  @ApiQuery({
    name: 'limit',
    required: false,
    description: 'Number of replies per page',
    example: 10,
    type: Number,
  })
  @ApiOkResponse({ description: 'Paginated replies' })
  @UseGuards(OptionalAuthGuard)
  @Get('comment/:commentId/replies')
  getReplies(
    @Param('commentId') commentId: string,
    @CurrentUser() user: { userId: string } | undefined,
    @Query('page') page = '1',
    @Query('limit') limit = '10',
  ) {
    return this.postGrpc.getReplies(
      commentId,
      Number(page),
      Number(limit),
      user?.userId,
    );
  }

  @ApiBearerAuth('access-token')
  @ApiOperation({ summary: 'Like a comment' })
  @ApiParam({
    name: 'commentId',
    description: 'Comment ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiCreatedResponse({ description: 'Comment liked' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Post('comment/:commentId/like')
  likeComment(
    @Param('commentId') commentId: string,
    @CurrentUser() user: { userId: string },
  ) {
    return this.postGrpc.likeComment(commentId, user.userId);
  }

  @ApiBearerAuth('access-token')
  @ApiOperation({ summary: 'Remove a like from a comment' })
  @ApiParam({
    name: 'commentId',
    description: 'Comment ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiOkResponse({ description: 'Comment like removed' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Delete('comment/:commentId/like')
  unlikeComment(
    @Param('commentId') commentId: string,
    @CurrentUser() user: { userId: string },
  ) {
    return this.postGrpc.unlikeComment(commentId, user.userId);
  }

  @ApiBearerAuth('access-token')
  @ApiOperation({ summary: 'Update a comment' })
  @ApiParam({
    name: 'commentId',
    description: 'Comment ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiOkResponse({ description: 'Updated comment' })
  @ApiBadRequestResponse({ description: 'Invalid comment data' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Put('comment/:commentId')
  updateComment(
    @Param('commentId') commentId: string,
    @CurrentUser() user: { userId: string },
    @Body() body: UpdateCommentDto,
  ) {
    return this.postGrpc.updateComment(commentId, user.userId, body.content);
  }

  @ApiBearerAuth('access-token')
  @ApiOperation({ summary: 'Delete a comment' })
  @ApiParam({
    name: 'commentId',
    description: 'Comment ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiOkResponse({ description: 'Deleted comment' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Delete('comment/:commentId')
  deleteComment(
    @Param('commentId') commentId: string,
    @CurrentUser() user: { userId: string },
  ) {
    return this.postGrpc.deleteComment(commentId, user.userId);
  }
}
