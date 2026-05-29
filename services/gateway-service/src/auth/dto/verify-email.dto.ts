import { IsEmail, IsString, MaxLength, MinLength } from "class-validator";

export class VerifyEmailDTO {
    @IsString()
    @IsEmail()
    email: string;

    @IsString()
    @MinLength(8)
    @MaxLength(64)
    password: string;

    @IsString()
    token: string;
}