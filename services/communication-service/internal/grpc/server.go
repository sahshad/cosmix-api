package grpc

import (
	"net"

	communicationpb "communication-service/internal/grpc/generated"

	"cosmix/shared/grpc/interceptors"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
}

func NewServer(
	communicationServer communicationpb.CommunicationServiceServer,
	zapLogger *zap.Logger,
) *Server {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RequestIDInterceptor,
			interceptors.LoggingInterceptor(zapLogger),
			interceptors.RecoveryInterceptor(zapLogger),
			interceptors.ErrorInterceptor,
		),
	)

	communicationpb.RegisterCommunicationServiceServer(
		grpcServer,
		communicationServer,
	)

	return &Server{
		grpcServer: grpcServer,
	}
}

func (s *Server) Start(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	return s.grpcServer.Serve(lis)
}

func (s *Server) GracefulStop() {
	s.grpcServer.GracefulStop()
}

func (s *Server) Stop() {
	s.grpcServer.Stop()
}
