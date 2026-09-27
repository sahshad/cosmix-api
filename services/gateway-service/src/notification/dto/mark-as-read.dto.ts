import { IsString, IsUUID } from 'class-validator';

export class MarkAsReadDto {
  @IsString()
  @IsUUID()
  notificationId: string;
}
