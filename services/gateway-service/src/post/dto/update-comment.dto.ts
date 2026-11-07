import { ApiProperty } from '@nestjs/swagger';
import { IsString } from 'class-validator';

export class UpdateCommentDto {
  @ApiProperty({
    description: 'Updated comment text',
    example: 'Nice post, edited!',
  })
  @IsString()
  content!: string;
}
