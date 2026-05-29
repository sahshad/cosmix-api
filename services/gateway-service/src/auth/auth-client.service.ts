import { Inject, Injectable } from '@nestjs/common';
import {
    AUTH_PACKAGE_NAME,
    AUTH_SERVICE_NAME,
    AuthServiceClient,
    LoginRequest,
    LoginResponse,
    RefreshResponse,
    RegisterRequest,
    RegisterResponse,
    VerifyEmailRequest,
    VerifyEmailResponse,
} from '../generated/auth/auth';
import { RpcException, type ClientGrpc } from '@nestjs/microservices';
import { catchError, firstValueFrom, Observable } from 'rxjs';

@Injectable()
export class AuthGrpcService {
    private client!: AuthServiceClient;

    constructor(
        @Inject(AUTH_PACKAGE_NAME)
        private readonly grpcClient: ClientGrpc,
    ) { }

    onModuleInit() {
        this.client =
            this.grpcClient.getService<AuthServiceClient>(
                AUTH_SERVICE_NAME,
            );
    }

    private call<T>(observable: Observable<T>): Promise<T> {
        return firstValueFrom(
            observable.pipe(
                catchError(err => { throw new RpcException(err); })
            )
        );
    }

    register(body: RegisterRequest): Promise<RegisterResponse> {
        return this.call(this.client.register(body)).then(res => ({
            ...res,
            userId: Number(res.userId)
        }));
    }

    verifyEmail(body: VerifyEmailRequest): Promise<VerifyEmailResponse> {
        return this.call(this.client.verifyEmail(body))
    }

    login(body: LoginRequest): Promise<LoginResponse> {
        return this.call(this.client.login(body))
    }

    refresh(refreshToken: string): Promise<RefreshResponse> {
        return this.call(this.client.refresh({ refreshToken }))
    }
}