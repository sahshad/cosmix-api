import { ApiProperty } from '@nestjs/swagger';
import { IsEmail } from 'class-validator';

export class ForgotPasswordDto {
  @ApiProperty({
    description: 'Email address that should receive the reset token',
    example: 'jane@example.com',
    format: 'email',
  })
  @IsEmail()
  email!: string;
}
