import { ApiProperty } from '@nestjs/swagger';
import { IsEmail, IsString, MaxLength, MinLength } from 'class-validator';

export class VerifyEmailDTO {
  @ApiProperty({
    description: 'Account email address',
    example: 'jane@example.com',
    format: 'email',
  })
  @IsString()
  @IsEmail()
  email: string;

  @ApiProperty({
    description: 'Account password',
    example: 'Password123!',
    minLength: 8,
    maxLength: 64,
    format: 'password',
  })
  @IsString()
  @MinLength(8)
  @MaxLength(64)
  password: string;

  @ApiProperty({
    description: 'Verification token sent to the email address',
    example: '482913',
  })
  @IsString()
  token: string;
}
