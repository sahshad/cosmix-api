import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import { IsArray, IsOptional, IsString } from 'class-validator';

export class MediaItemDto {
  @ApiProperty({
    description: 'Media provider public identifier',
    example: 'cosmix/posts/abc123',
  })
  @IsString()
  publicId!: string;

  @ApiProperty({
    description: 'Media URL',
    example: 'https://cdn.example.com/posts/abc123.jpg',
    format: 'uri',
  })
  @IsString()
  url!: string;

  @ApiProperty({
    description: 'Media type',
    example: 'image',
  })
  @IsString()
  type!: string;

  @ApiPropertyOptional({
    description: 'Duration in seconds (video only)',
    example: 12.5,
  })
  @IsOptional()
  duration?: number;
}

export class CreatePostDto {
  @ApiProperty({
    description: 'Post text content',
    example: 'Hello world!',
  })
  @IsString()
  content!: string;

  @ApiPropertyOptional({
    description: 'Media attached to the post',
    type: () => [MediaItemDto],
  })
  @IsOptional()
  @IsArray()
  media?: MediaItemDto[];
}
