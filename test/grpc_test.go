package test

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"testing"

	"github.com/Yudiz/ecommerce-go-config-service/config"
	internalGrpc "github.com/Yudiz/ecommerce-go-config-service/internal/grpc"
	pb "github.com/Yudiz/ecommerce-go-config-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	// Mock Config
	cfg := &config.Config{
		UserService: config.PostgresServiceConfig{
			DB: config.PostgresConfig{
				Host: "mock-host",
			},
		},
	}
	pb.RegisterConfigServiceServer(s, internalGrpc.NewConfigServer(cfg))
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestGetConfig(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewConfigServiceClient(conn)
	resp, err := client.GetConfig(ctx, &pb.GetConfigRequest{ServiceName: "user-service"})
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(resp.ConfigJson), &data); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	dbData := data["db"].(map[string]interface{})
	if dbData["host"] != "mock-host" {
		t.Errorf("Expected host 'mock-host', got %v", dbData["host"])
	}
}
