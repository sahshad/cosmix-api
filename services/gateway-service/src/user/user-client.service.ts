import { Inject, Injectable } from '@nestjs/common';
import {
  FollowResponse,
  UnfollowResponse,
  UpdateProfileRequest,
  USER_PACKAGE_NAME,
  USER_SERVICE_NAME,
  UserListResponse,
  UserProfileResponse,
  UserServiceClient,
} from '../generated/user/user';
import { RpcException, type ClientGrpc } from '@nestjs/microservices';
import { catchError, firstValueFrom, Observable } from 'rxjs';

@Injectable()
export class UserGrpcService {
  private client!: UserServiceClient;

  constructor(
    @Inject(USER_PACKAGE_NAME)
    private readonly grpcClient: ClientGrpc,
  ) { }

  onModuleInit() {
    this.client =
      this.grpcClient.getService<UserServiceClient>(
        USER_SERVICE_NAME,
      );
  }

  private call<T>(observable: Observable<T>): Promise<T> {
    return firstValueFrom(
      observable.pipe(
        catchError(err => { throw new RpcException(err); })
      )
    );
  }

  getProfile(userId: number): Promise<UserProfileResponse> {
    return this.call(this.client.getProfile({ userId }))
  }

  getProfileByUsername(username: string,): Promise<UserProfileResponse> {
    return this.call(this.client.getProfileByUsername({ username }))
  }

  updateProfile(userId: number, body: Partial<UpdateProfileRequest>): Promise<UserProfileResponse> {
    return this.call(this.client.updateProfile({ userId, ...body }))
  }

  follow(followerId: number, followingId: number): Promise<FollowResponse> {
    return this.call(this.client.follow({ followerId, followingId }))
  }

  unfollow(followerId: number, followingId: number): Promise<UnfollowResponse> {
    return this.call(this.client.unfollow({ followerId, followingId }))
  }

  getFollowers(userId: number): Promise<UserListResponse> {
    return this.call(this.client.getFollowers({ userId }))
  }

  getFollowing(userId: number): Promise<UserListResponse> {
    return this.call(this.client.getFollowing({ userId }))
  }
}