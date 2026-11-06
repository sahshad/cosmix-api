import { Module } from '@nestjs/common';
import { CommunicationController } from './communication.controller';
import { CommunicationGrpcService } from './communication-client.service';
import { CommunicationService } from './communication.service';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { COMMUNICATION_PACKAGE_NAME } from '../generated/communication/communication';
import { join } from 'path';

@Module({
  imports: [
    ClientsModule.register([
      {
        name: COMMUNICATION_PACKAGE_NAME,
        transport: Transport.GRPC,
        options: {
          url: process.env.COMMUNICATION_GRPC_ADDR ?? 'localhost:50056',
          package: 'communication',
          protoPath: join(
            __dirname,
            '../../../../shared/grpc/proto/communication/communication.proto',
          ),
          loader: {
            defaults: true,
          },
        },
      },
    ]),
  ],
  controllers: [CommunicationController],
  providers: [CommunicationGrpcService, CommunicationService],
  exports: [CommunicationGrpcService],
})
export class CommunicationModule {}
