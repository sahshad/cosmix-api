import {
  IsArray,
  IsOptional,
  IsString,
} from 'class-validator';

export class MediaItemDto {
  @IsString()
  publicId!: string;

  @IsString()
  url!: string;

  @IsString()
  type!: string;

  @IsOptional()
  duration?: number;
}

export class UpdatePostDto {
  @IsString()
  content!: string;

  @IsOptional()
  @IsArray()
  media?: MediaItemDto[];
}