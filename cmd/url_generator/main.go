package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/DeedsBaron/url_generator/internal/api/url_generator_v1"
	"github.com/DeedsBaron/url_generator/internal/config"
	"github.com/DeedsBaron/url_generator/internal/pkg/url_generator_service"
	"github.com/DeedsBaron/url_generator/internal/repository"
	"github.com/DeedsBaron/url_generator/internal/serverwrapper"
	desc "github.com/DeedsBaron/url_generator/pkg/url_generator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := config.New()
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}
	repo := repository.NewRepo()

	service := url_generator_service.New(repo)

	grpcServer := serverwrapper.NewGrpcServer(config.Data.GrpcPort)

	desc.RegisterUrlGeneratorServer(grpcServer.GetServer(), url_generator_v1.NewUrlGeneratorV1(service))
	// http gateway configuration
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err = desc.RegisterUrlGeneratorHandlerFromEndpoint(context.Background(),
		mux,
		grpcServer.GetListener().Addr().String(),
		opts)
	if err != nil {
		log.Fatalf("failed to register http gateway to gRPC handler: %v", err)
	}

	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(config.Data.HttpPort),
		Handler: mux,
	}

	// Start gRPC server
	go func() {
		log.Printf("starting gRPC server on port: %d\n", config.Data.GrpcPort)
		if err := grpcServer.Serve(); err != nil {
			log.Fatalf("failed to serve gRPC server: %v", err)
		}
	}()

	// Start HTTP server
	log.Printf("starting HTTP proxy server on port :%d\n", config.Data.HttpPort)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve HTTP server: %v", err)
	}
}
