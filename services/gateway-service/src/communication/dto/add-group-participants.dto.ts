import { ApiProperty } from '@nestjs/swagger';
import {
  ArrayMaxSize,
  ArrayMinSize,
  ArrayUnique,
  IsArray,
  IsUUID,
} from 'class-validator';

export class AddGroupParticipantsDto {
  @ApiProperty({
    description: 'User IDs to add to the group',
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
}
