import {
    Body,
    Controller,
    Get,
    Post,
    Req,
    Res,
    UnauthorizedException,
} from '@nestjs/common';

import { AuthGrpcService } from './auth-client.service';
import { LoginDto } from './dto/login.dto';
import { RegisterDTO } from './dto/register.dto';
import type { Request, Response } from 'express';
import { Throttle } from '@nestjs/throttler';

@Throttle({
    default: {
        limit: 20,
        ttl: 60000,
    },
})
@Controller('auth')
export class AuthController {
    constructor(
        private readonly authGrpc: AuthGrpcService,
    ) { }

    @Post('register')
    async register(@Body() body: RegisterDTO) {
        console.log("register body: ", body)
        const result = await this.authGrpc.register(body)
        return result
    }

    @Post('login')
    async login(@Body() body: LoginDto, @Res({ passthrough: true }) res: Response) {
        const result = await this.authGrpc.login(body);

        res.cookie(
            'refresh_token',
            result.refreshToken,
            {
                httpOnly: true,
                secure:
                    process.env.NODE_ENV === 'production',
                sameSite: 'strict',
                maxAge:
                    30 * 24 * 60 * 60 * 1000,
            },
        );

        return {
            access_token: result.accessToken,
            user: result.user,
        };
    }

    @Get('refresh')
    async refresh(
        @Req() req: Request,
        @Res({ passthrough: true })
        res: Response,
    ) {

        const refreshToken =
            req.cookies?.refresh_token;

        if (!refreshToken) {
            throw new UnauthorizedException(
                'No refresh token',
            );
        }

        const result =
            await this.authGrpc.refresh(
                refreshToken,
            );

        res.cookie(
            'refresh_token',
            result.refreshToken,
            {
                httpOnly: true,
                secure:
                    process.env.NODE_ENV ===
                    'production',
                sameSite: 'strict',
                maxAge:
                    30 * 24 * 60 * 60 * 1000,
            },
        );

        return {
            access_token:
                result.accessToken,
        };
    }


    @Post('logout')
    async logout(
        @Res({ passthrough: true })
        res: Response,
    ) {

        res.clearCookie(
            'refresh_token',
            {
                httpOnly: true,
                sameSite: 'strict',
            },
        );

        return {
            message: 'logged out',
        };
    }
}