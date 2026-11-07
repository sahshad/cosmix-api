import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import { IsArray, IsOptional, IsString } from 'class-validator';
import { MediaItemDto } from './create-post.dto';

export class UpdatePostDto {
  @ApiProperty({
    description: 'Updated post text content',
    example: 'Hello world, edited!',
  })
  @IsString()
  content!: string;

  @ApiPropertyOptional({
    description: 'Replacement media for the post',
    type: () => [MediaItemDto],
  })
  @IsOptional()
  @IsArray()
  media?: MediaItemDto[];
}
