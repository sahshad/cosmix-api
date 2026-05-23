import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Post,
  Put,
  UseGuards,
} from '@nestjs/common';

import { UserGrpcService } from './user-client.service';
import { AuthGuard } from '../common/guard/auth.guard';
import { CurrentUser } from '../common/decorators/current-user.decorator';

import { UpdateProfileDto } from './dto/update-profile.dto';

@Controller('users')
export class UserController {
  constructor(
    private readonly userGrpc: UserGrpcService,
  ) {}

  @UseGuards(AuthGuard)
  @Get('me')
  async getMyProfile(
    @CurrentUser() user: { userId: number },
  ) {
    return this.userGrpc.getProfile(
      user.userId,
    );
  }

  @UseGuards(AuthGuard)
  @Put('me')
  async updateMyProfile(
    @CurrentUser() user: { userId: number },
    @Body() body: UpdateProfileDto,
  ) {
    return this.userGrpc.updateProfile(
      user.userId,
      body,
    );
  }

  @Get('username/:username')
  async getByUsername(
    @Param('username') username: string,
  ) {
    return this.userGrpc.getProfileByUsername(
      username,
    );
  }

  @UseGuards(AuthGuard)
  @Post('follow/:id')
  async follow(
    @CurrentUser() user: { userId: number },
    @Param('id') followingId: string,
  ) {
    return this.userGrpc.follow(
      user.userId,
      Number(followingId),
    );
  }

  @UseGuards(AuthGuard)
  @Delete('unfollow/:id')
  async unfollow(
    @CurrentUser() user: { userId: number },
    @Param('id') followingId: string,
  ) {
    return this.userGrpc.unfollow(
      user.userId,
      Number(followingId),
    );
  }

  @Get('followers/:id')
  async getFollowers(
    @Param('id') userId: string,
  ) {
    return this.userGrpc.getFollowers(
      Number(userId),
    );
  }

  @Get('following/:id')
  async getFollowing(
    @Param('id') userId: string,
  ) {
    return this.userGrpc.getFollowing(
      Number(userId),
    );
  }
}