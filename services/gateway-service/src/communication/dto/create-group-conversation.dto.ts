import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import {
  ArrayMaxSize,
  ArrayMinSize,
  ArrayUnique,
  IsArray,
  IsOptional,
  IsString,
  IsUUID,
  MaxLength,
  MinLength,
} from 'class-validator';

export class CreateGroupConversationDto {
  @ApiProperty({
    description: 'Group name',
    example: 'Weekend trip',
    minLength: 1,
    maxLength: 100,
  })
  @IsString()
  @MinLength(1)
  @MaxLength(100)
  name!: string;

  @ApiProperty({
    description: 'User IDs of the group participants',
    example: [
      '0f8fad5b-d9cb-469f-a165-70867728950e',
      '7c9e6679-7425-40de-944b-e07fc1f90ae7',
    ],
    type: [String],
    format: 'uuid',
    minItems: 1,
    maxItems: 99,
    uniqueItems: true,
  })
  @IsArray()
  @ArrayMinSize(1)
  @ArrayMaxSize(99)
  @ArrayUnique()
  @IsUUID(undefined, { each: true })
  participantIds!: string[];

  @ApiPropertyOptional({
    description: 'Group avatar URL',
    example: 'https://cdn.example.com/groups/weekend.png',
    format: 'uri',
    maxLength: 2048,
  })
  @IsOptional()
  @IsString()
  @MaxLength(2048)
  avatarUrl?: string;
}
