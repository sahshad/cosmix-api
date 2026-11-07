import { ApiProperty } from '@nestjs/swagger';
import { IsEmail, IsString } from 'class-validator';

export class RegisterDTO {
  @ApiProperty({
    description: 'Public display name',
    example: 'Jane Doe',
  })
  @IsString()
  displayName!: string;

  @ApiProperty({
    description: 'Account email address',
    example: 'jane@example.com',
    format: 'email',
  })
  @IsEmail()
  email!: string;

  @ApiProperty({
    description: 'Account password',
    example: 'Password123!',
    format: 'password',
  })
  @IsString()
  password!: string;
}
