import { ApiProperty } from '@nestjs/swagger';
import { IsString, Length } from 'class-validator';

export class ChangePasswordDTO {
  @ApiProperty({
    description: 'Password reset token',
    example: '482913',
  })
  @IsString()
  token!: string;

  @ApiProperty({
    description: 'New account password',
    example: 'NewPassword123!',
    minLength: 8,
    maxLength: 64,
    format: 'password',
  })
  @Length(8, 64)
  newPassword!: string;
}
