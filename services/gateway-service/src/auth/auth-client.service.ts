import { Injectable } from '@nestjs/common';
import * as grpc from '@grpc/grpc-js';
import {
    AuthServiceClient,
    LoginRequest,
    LoginResponse,
    RefreshResponse,
    RegisterRequest,
    RegisterResponse,
} from '../generated/auth/auth';
import { handleGrpcError } from '../common/utils/grpc-error.util';

@Injectable()
export class AuthGrpcService {
    private readonly client: AuthServiceClient;

    constructor() {
        this.client = new AuthServiceClient(
            process.env.AUTH_GRPC_ADDR || 'auth-service:50051',
            grpc.credentials.createInsecure(),
        );
    }

    getClient() {
        return this.client;
    }

    register(body: RegisterRequest): Promise<RegisterResponse> {
        return new Promise((resolve, reject) => {
            const reqBody: RegisterRequest = body

            this.client.register(reqBody, (err, response) => {
                if (err) {
                    return reject(handleGrpcError(err));
                }

                resolve(response);
            });
        })
    }

    login(body: LoginRequest): Promise<LoginResponse> {
        return new Promise((resolve, reject) => {
            const reqBody: LoginRequest = body

            this.client.login(
                reqBody,
                (err, response) => {
                    if (err) {
                        return reject(handleGrpcError(err));
                    }

                    resolve(response);
                },
            );
        });
    }

    refresh(refreshToken: string): Promise<RefreshResponse> {
        return new Promise(
            (resolve, reject) => {
                this.client.refresh(
                    {
                        refreshToken,
                    },
                    (err, response) => {
                        if (err) {
                            return reject(
                                handleGrpcError(err),
                            );
                        }

                        resolve(response);
                    },
                );
            },
        );
    }
}