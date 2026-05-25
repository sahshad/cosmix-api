import { Expose } from "class-transformer";
import { IsEmail, IsString } from "class-validator";

export class RegisterDTO {
    @Expose({ name: 'display_name' })
    @IsString()
    displayName!: string;

    @IsEmail()
    email!: string;

    @IsString()
    password!: string;
}