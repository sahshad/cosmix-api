import { IsString, Length } from 'class-validator';

export class ChangePasswordDTO {
    @IsString()
    token!: string;

    @Length(8, 64)
    newPassword!: string;
}