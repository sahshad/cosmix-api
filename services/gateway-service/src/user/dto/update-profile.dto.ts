import {
  IsDateString,
  IsOptional,
  IsString,
  MaxLength,
} from 'class-validator';

export class UpdateProfileDto {
  @IsOptional()
  @IsString()
  displayName?: string;

  @IsOptional()
  @IsString()
  username?: string;

  @IsOptional()
  @IsDateString()
  dateOfBirth?: string;

  @IsOptional()
  @IsString()
  avatarUrl?: string;

  @IsOptional()
  @IsString()
  @MaxLength(500)
  bio?: string;
}