import { Injectable } from '@nestjs/common';
import * as grpc from '@grpc/grpc-js';

import {
  PostServiceClient,
} from '../generated/post/post';

import { grpcUnaryCall } from '../common/utils/grpc.util';

@Injectable()
export class PostGrpcService {
  private readonly client: PostServiceClient;

  constructor() {
    this.client =
      new PostServiceClient(
        process.env.POST_GRPC_ADDR ??
          'post-service:50054',
        grpc.credentials.createInsecure(),
      );
  }

  createPost(
    authorId: number,
    body: any,
  ) {
    return grpcUnaryCall(
      callback =>
        this.client.createPost(
          {
            authorId,
            content: body.content,
            media: body.media ?? [],
          },
          callback,
        ),
    );
  }

  getPost(postId: number) {
    return grpcUnaryCall(
      callback =>
        this.client.getPost(
          { postId },
          callback,
        ),
    );
  }

  getFeed(
    page: number,
    limit: number,
  ) {
    return grpcUnaryCall(
      callback =>
        this.client.getFeed(
          {
            page,
            limit,
          },
          callback,
        ),
    );
  }

  updatePost(
    postId: number,
    authorId: number,
    body: any,
  ) {
    return grpcUnaryCall(
      callback =>
        this.client.updatePost(
          {
            postId,
            authorId,
            content: body.content,
            media: body.media ?? [],
          },
          callback,
        ),
    );
  }

  deletePost(
    postId: number,
    authorId: number,
  ) {
    return grpcUnaryCall(
      callback =>
        this.client.deletePost(
          {
            postId,
            authorId,
          },
          callback,
        ),
    );
  }

  likePost(
    postId: number,
    userId: number,
  ) {
    return grpcUnaryCall(
      callback =>
        this.client.likePost(
          {
            postId,
            userId,
          },
          callback,
        ),
    );
  }

  unlikePost(
    postId: number,
    userId: number,
  ) {
    return grpcUnaryCall(
      callback =>
        this.client.unlikePost(
          {
            postId,
            userId,
          },
          callback,
        ),
    );
  }

  createComment(
    postId: number,
    authorId: number,
    content: string,
  ) {
    return grpcUnaryCall(
      callback =>
        this.client.createComment(
          {
            postId,
            authorId,
            content,
          },
          callback,
        ),
    );
  }

  getComments(
    postId: number,
    page: number,
    limit: number,
  ) {
    return grpcUnaryCall(
      callback =>
        this.client.getComments(
          {
            postId,
            page,
            limit,
          },
          callback,
        ),
    );
  }

  updateComment(
    commentId: number,
    authorId: number,
    content: string,
  ) {
    return grpcUnaryCall(
      callback =>
        this.client.updateComment(
          {
            commentId,
            authorId,
            content,
          },
          callback,
        ),
    );
  }

  deleteComment(
    commentId: number,
    authorId: number,
  ) {
    return grpcUnaryCall(
      callback =>
        this.client.deleteComment(
          {
            commentId,
            authorId,
          },
          callback,
        ),
    );
  }
}