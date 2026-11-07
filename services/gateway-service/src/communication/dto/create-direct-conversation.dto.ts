import { ApiProperty } from '@nestjs/swagger';
import { IsUUID } from 'class-validator';

export class CreateDirectConversationDto {
  @ApiProperty({
    description: 'User ID of the other participant',
    example: '0f8fad5b-d9cb-469f-a165-70867728950e',
    format: 'uuid',
  })
  @IsUUID()
  participantId!: string;
}
