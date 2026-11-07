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
import {
  ApiBadRequestResponse,
  ApiBearerAuth,
  ApiCreatedResponse,
  ApiNotFoundResponse,
  ApiOkResponse,
  ApiOperation,
  ApiParam,
  ApiTags,
  ApiUnauthorizedResponse,
} from '@nestjs/swagger';

import { UserGrpcService } from './user-client.service';
import { AuthGuard } from '../common/guard/auth.guard';
import { OptionalAuthGuard } from '../common/guard/optional-auth.guard';
import { CurrentUser } from '../common/decorators/current-user.decorator';

import { UpdateProfileDto } from './dto/update-profile.dto';

@ApiTags('users')
@Controller('users')
export class UserController {
  constructor(private readonly userGrpc: UserGrpcService) {}

  @ApiBearerAuth('access-token')
  @ApiOperation({ summary: 'Get the authenticated user profile' })
  @ApiOkResponse({ description: 'Current user profile' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Get('me')
  async getMyProfile(@CurrentUser() user: { userId: string }) {
    return this.userGrpc.getProfile(user.userId);
  }

  @ApiBearerAuth('access-token')
  @ApiOperation({ summary: 'Update the authenticated user profile' })
  @ApiOkResponse({ description: 'Updated user profile' })
  @ApiBadRequestResponse({ description: 'Invalid profile data' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Put('me')
  async updateMyProfile(
    @CurrentUser() user: { userId: string },
    @Body() body: UpdateProfileDto,
  ) {
    return this.userGrpc.updateProfile(user.userId, body);
  }

  @ApiOperation({
    summary: 'Get a user profile by username',
    description:
      'Authentication is optional; send a Bearer token to include viewer-specific data.',
  })
  @ApiParam({
    name: 'username',
    description: 'Username of the profile to fetch',
    example: 'janedoe',
  })
  @ApiOkResponse({ description: 'User profile' })
  @ApiNotFoundResponse({ description: 'User not found' })
  @UseGuards(OptionalAuthGuard)
  @Get('username/:username')
  async getByUsername(
    @Param('username') username: string,
    @CurrentUser() user: { userId: string } | undefined,
  ) {
    return this.userGrpc.getProfileByUsername(username, user?.userId);
  }

  @ApiBearerAuth('access-token')
  @ApiOperation({ summary: 'Follow a user' })
  @ApiParam({
    name: 'id',
    description: 'User ID to follow',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiCreatedResponse({ description: 'Followed the user' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Post('follow/:id')
  async follow(
    @CurrentUser() user: { userId: string },
    @Param('id') followingId: string,
  ) {
    return this.userGrpc.follow(user.userId, followingId);
  }

  @ApiBearerAuth('access-token')
  @ApiOperation({ summary: 'Unfollow a user' })
  @ApiParam({
    name: 'id',
    description: 'User ID to unfollow',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiOkResponse({ description: 'Unfollowed the user' })
  @ApiUnauthorizedResponse({ description: 'Missing or invalid access token' })
  @UseGuards(AuthGuard)
  @Delete('unfollow/:id')
  async unfollow(
    @CurrentUser() user: { userId: string },
    @Param('id') followingId: string,
  ) {
    return this.userGrpc.unfollow(user.userId, followingId);
  }

  @ApiOperation({ summary: 'List followers of a user' })
  @ApiParam({
    name: 'id',
    description: 'User ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiOkResponse({ description: 'List of followers' })
  @Get('followers/:id')
  async getFollowers(@Param('id') userId: string) {
    return this.userGrpc.getFollowers(userId);
  }

  @ApiOperation({ summary: 'List users followed by a user' })
  @ApiParam({
    name: 'id',
    description: 'User ID',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @ApiOkResponse({ description: 'List of followed users' })
  @Get('following/:id')
  async getFollowing(@Param('id') userId: string) {
    return this.userGrpc.getFollowing(userId);
  }
}
