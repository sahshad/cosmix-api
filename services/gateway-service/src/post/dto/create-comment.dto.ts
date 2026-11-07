import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import { IsOptional, IsString } from 'class-validator';

export class CreateCommentDto {
  @ApiProperty({
    description: 'Comment text',
    example: 'Nice post!',
  })
  @IsString()
  content!: string;

  @ApiPropertyOptional({
    description: 'Parent comment ID when replying to another comment',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @IsOptional()
  @IsString()
  parentCommentId?: string;
}
