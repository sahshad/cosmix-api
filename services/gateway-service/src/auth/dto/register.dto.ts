import { IsEmail, IsString } from "class-validator";

export class RegisterDTO {
    @IsString()
    displayName!: string;

    @IsEmail()
    email!: string;

    @IsString()
    password!: string;
}