import { ApiPropertyOptional } from '@nestjs/swagger';
import { IsOptional, IsString, MaxLength, MinLength } from 'class-validator';

export class UpdateGroupConversationDto {
  @ApiPropertyOptional({
    description: 'New group name',
    example: 'Weekend trip 2026',
    minLength: 1,
    maxLength: 100,
  })
  @IsOptional()
  @IsString()
  @MinLength(1)
  @MaxLength(100)
  name?: string;

  @ApiPropertyOptional({
    description: 'New group avatar URL',
    example: 'https://cdn.example.com/groups/weekend-2026.png',
    format: 'uri',
    maxLength: 2048,
  })
  @IsOptional()
  @IsString()
  @MaxLength(2048)
  avatarUrl?: string;
}
