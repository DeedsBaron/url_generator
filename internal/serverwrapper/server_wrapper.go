package serverwrapper

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	lis        net.Listener
	grpcServer *grpc.Server
}

func NewGrpcServer(port int, interceptors ...grpc.UnaryServerInterceptor) *grpcServer {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	reflection.Register(s)

	return &grpcServer{
		lis:        lis,
		grpcServer: s,
	}
}

func (s *grpcServer) Serve() error {
	return s.grpcServer.Serve(s.lis)
}

func (s *grpcServer) GetServer() *grpc.Server {
	return s.grpcServer
}

func (s *grpcServer) GetListener() net.Listener {
	return s.lis
}
