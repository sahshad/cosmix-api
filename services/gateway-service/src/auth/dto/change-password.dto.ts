import { Expose } from 'class-transformer';
import { IsString, Length } from 'class-validator';

export class ChangePasswordDTO {
    @IsString()
    token!: string;

    @Expose({ name: 'new_password' })
    @Length(8, 64)
    newPassword!: string;
}