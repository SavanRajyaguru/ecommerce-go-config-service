package main

import (
	"fmt"
	"log"
	"net"

	"github.com/SavanRajyaguru/ecommerce-go-config-service/api"
	"github.com/SavanRajyaguru/ecommerce-go-config-service/config"
	internalGrpc "github.com/SavanRajyaguru/ecommerce-go-config-service/internal/grpc"
	pb "github.com/SavanRajyaguru/ecommerce-go-config-service/proto"
	"google.golang.org/grpc"
)

func main() {
	// Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Router
	r := api.SetupRouter()

	// Start gRPC Server
	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.Server.GrpcPort)
		if err != nil {
			log.Fatalf("Failed to listen for gRPC: %v", err)
		}

		grpcServer := grpc.NewServer()
		configServer := internalGrpc.NewConfigServer(cfg)
		pb.RegisterConfigServiceServer(grpcServer, configServer)

		log.Printf("Starting gRPC server on :%s", cfg.Server.GrpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP Server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting HTTP server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
