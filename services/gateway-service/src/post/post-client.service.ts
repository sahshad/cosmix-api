import { Inject, Injectable } from '@nestjs/common';

import {
  CommentListResponse,
  CommentResponse,
  MessageResponse,
  POST_PACKAGE_NAME,
  POST_SERVICE_NAME,
  PostListResponse,
  PostResponse,
  PostServiceClient,
} from '../generated/post/post';

import type { ClientGrpc } from '@nestjs/microservices';
import { firstValueFrom } from 'rxjs';

@Injectable()
export class PostGrpcService {
  private client!: PostServiceClient;

  constructor(
    @Inject(POST_PACKAGE_NAME)
    private readonly grpcClient: ClientGrpc,
  ) { }

  onModuleInit() {
    this.client =
      this.grpcClient.getService<PostServiceClient>(
        POST_SERVICE_NAME,
      );
  }

  createPost(authorId: number, body: any): Promise<PostResponse> {
    return firstValueFrom(this.client.createPost({ authorId, ...body }))
  }

  getPost(postId: number): Promise<PostResponse> {
    return firstValueFrom(this.client.getPost({ postId }))
  }

  getFeed(page: number, limit: number): Promise<PostListResponse> {
    return firstValueFrom(this.client.getFeed({ page, limit }))
  }

  updatePost(postId: number, authorId: number, body: any): Promise<PostResponse> {
    return firstValueFrom(this.client.updatePost({ postId, authorId, ...body }))
  }

  deletePost(postId: number, authorId: number): Promise<MessageResponse> {
    return firstValueFrom(this.client.deletePost({ postId, authorId }))
  }

  likePost(postId: number,userId: number): Promise<MessageResponse> {
    return firstValueFrom(this.client.likePost({ postId, userId }))
  }

  unlikePost(postId: number,userId: number): Promise<MessageResponse> {
    return firstValueFrom(this.client.unlikePost({ postId, userId }))
  }

  createComment(postId: number,authorId: number,content: string): Promise<CommentResponse> {
    return firstValueFrom(this.client.createComment({ postId,authorId,content}))
  }

  getComments(postId: number,page: number,limit: number): Promise<CommentListResponse> {
    return firstValueFrom(this.client.getComments({ postId,page,limit}))
  }

  updateComment(commentId: number,authorId: number,content: string): Promise<CommentResponse> {
    return firstValueFrom(this.client.updateComment({ commentId,authorId,content}))
  }

  deleteComment(commentId: number,authorId: number): Promise<MessageResponse> {
    return firstValueFrom(this.client.deleteComment({ commentId,authorId}))
  }
}