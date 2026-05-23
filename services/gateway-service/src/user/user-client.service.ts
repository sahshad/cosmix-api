import { Injectable } from '@nestjs/common';
import * as grpc from '@grpc/grpc-js';

import {
  FollowResponse,
  GetFollowersRequest,
  GetFollowingRequest,
  GetProfileByUsernameRequest,
  GetProfileRequest,
  UnfollowResponse,
  UpdateProfileRequest,
  UserListResponse,
  UserProfileResponse,
  UserServiceClient,
} from '../generated/user/user';

import { grpcUnaryCall } from '../common/utils/grpc.util';

@Injectable()
export class UserGrpcService {
  private readonly client: UserServiceClient;

  constructor() {
    this.client =
      new UserServiceClient(
        process.env.USER_GRPC_ADDR ??
          'user-service:50052',
        grpc.credentials.createInsecure(),
      );
  }

  getProfile(
    userId: number,
  ): Promise<UserProfileResponse> {
    return grpcUnaryCall(
      callback =>
        this.client.getProfile(
          { userId },
          callback,
        ),
    );
  }

  getProfileByUsername(
    username: string,
  ): Promise<UserProfileResponse> {
    return grpcUnaryCall(
      callback =>
        this.client.getProfileByUsername(
          { username },
          callback,
        ),
    );
  }

  updateProfile(
    userId: number,
    body: Partial<UpdateProfileRequest>,
  ): Promise<UserProfileResponse> {
    return grpcUnaryCall(
      callback =>
        this.client.updateProfile(
          {
            userId,
            ...body,
          },
          callback,
        ),
    );
  }

  follow(
    followerId: number,
    followingId: number,
  ): Promise<FollowResponse> {
    return grpcUnaryCall(
      callback =>
        this.client.follow(
          {
            followerId,
            followingId,
          },
          callback,
        ),
    );
  }

  unfollow(
    followerId: number,
    followingId: number,
  ): Promise<UnfollowResponse> {
    return grpcUnaryCall(
      callback =>
        this.client.unfollow(
          {
            followerId,
            followingId,
          },
          callback,
        ),
    );
  }

  getFollowers(
    userId: number,
  ): Promise<UserListResponse> {
    return grpcUnaryCall(
      callback =>
        this.client.getFollowers(
          { userId },
          callback,
        ),
    );
  }

  getFollowing(
    userId: number,
  ): Promise<UserListResponse> {
    return grpcUnaryCall(
      callback =>
        this.client.getFollowing(
          { userId },
          callback,
        ),
    );
  }
}