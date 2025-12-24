package grpc

import (
	"context"
	"encoding/json"
	"log"

	"github.com/SavanRajyaguru/ecommerce-go-config-service/config"
	pb "github.com/SavanRajyaguru/ecommerce-go-config-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConfigServer struct {
	pb.UnimplementedConfigServiceServer
	Config *config.Config
}

func NewConfigServer(cfg *config.Config) *ConfigServer {
	return &ConfigServer{Config: cfg}
}

func (s *ConfigServer) GetConfig(ctx context.Context, req *pb.GetConfigRequest) (*pb.GetConfigResponse, error) {
	log.Printf("Received config request for service: %s", req.ServiceName)

	var data interface{}

	switch req.ServiceName {
	case "user-service":
		data = s.Config.UserService
	case "product-service":
		data = s.Config.ProductService
	case "order-service":
		data = s.Config.OrderService
	case "payment-service":
		data = s.Config.PaymentService
	case "inventory-service":
		data = s.Config.InventoryService
	case "notification-service":
		data = s.Config.NotificationService
	default:
		return nil, status.Errorf(codes.NotFound, "Service config not found for: %s", req.ServiceName)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal config: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to marshal config")
	}

	return &pb.GetConfigResponse{
		ConfigJson: string(jsonData),
	}, nil
}
