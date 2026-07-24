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

  createPost(authorId: string, body: any): Promise<PostResponse> {
    return firstValueFrom(this.client.createPost({ authorId, ...body }))
  }

  getPost(postId: string): Promise<PostResponse> {
    return firstValueFrom(this.client.getPost({ postId }))
  }

  getFeed(page: number, limit: number, userId?: string): Promise<PostListResponse> {
    return firstValueFrom(this.client.getFeed({ page, limit, userId: userId ?? '' }))
  }

  getUserPosts(userId: string, page: number, limit: number, viewerId?: string): Promise<PostListResponse> {
    return firstValueFrom(this.client.getUserPosts({ userId, page, limit, viewerId: viewerId ?? '' }))
  }

  updatePost(postId: string, authorId: string, body: any): Promise<PostResponse> {
    return firstValueFrom(this.client.updatePost({ postId, authorId, ...body }))
  }

  deletePost(postId: string, authorId: string): Promise<MessageResponse> {
    return firstValueFrom(this.client.deletePost({ postId, authorId }))
  }

  likePost(postId: string,userId: string): Promise<MessageResponse> {
    return firstValueFrom(this.client.likePost({ postId, userId }))
  }

  unlikePost(postId: string,userId: string): Promise<MessageResponse> {
    return firstValueFrom(this.client.unlikePost({ postId, userId }))
  }

  createComment(postId: string, authorId: string, content: string, parentCommentId?: string): Promise<CommentResponse> {
    return firstValueFrom(this.client.createComment({ postId, authorId, content, parentCommentId: parentCommentId ?? '' }))
  }

  getComments(postId: string, page: number, limit: number, userId?: string): Promise<CommentListResponse> {
    return firstValueFrom(this.client.getComments({ postId, page, limit, userId: userId ?? '' }))
  }

  getReplies(commentId: string, page: number, limit: number, userId?: string): Promise<CommentListResponse> {
    return firstValueFrom(this.client.getReplies({ commentId, page, limit, userId: userId ?? '' }))
  }

  updateComment(commentId: string,authorId: string,content: string): Promise<CommentResponse> {
    return firstValueFrom(this.client.updateComment({ commentId,authorId,content}))
  }

  deleteComment(commentId: string,authorId: string): Promise<MessageResponse> {
    return firstValueFrom(this.client.deleteComment({ commentId,authorId}))
  }

  likeComment(commentId: string, userId: string): Promise<MessageResponse> {
    return firstValueFrom(this.client.likeComment({ commentId, userId }))
  }

  unlikeComment(commentId: string, userId: string): Promise<MessageResponse> {
    return firstValueFrom(this.client.unlikeComment({ commentId, userId }))
  }
}