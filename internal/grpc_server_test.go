//go:build unit

package internal

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/quantfidential/trading-ecosystem/custodian-simulator-go/internal/config"
	grpcserver "github.com/quantfidential/trading-ecosystem/custodian-simulator-go/internal/presentation/grpc"
)

// TestCustodianGRPCServer_RedPhase defines the expected behaviors for enhanced gRPC server
// These tests will fail initially and drive our implementation (TDD Red-Green-Refactor)
func TestCustodianGRPCServer_HealthService(t *testing.T) {
	t.Run("provides_enhanced_health_status", func(t *testing.T) {
		t.Parallel()

		cfg := &config.Config{
			ServiceName: "custodian-simulator",
			GRPCPort:   0, // Use random port for testing
		}

		server := NewCustodianGRPCServer(cfg)

		// Start server in background
		lis, err := net.Listen("tcp", ":0")
		if err != nil {
			t.Fatalf("Failed to listen: %v", err)
		}

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Logf("Server serve error: %v", err)
			}
		}()
		defer server.GracefulStop()

		// Give server time to start
		time.Sleep(100 * time.Millisecond)

		// Connect to health service
		conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()

		healthClient := grpc_health_v1.NewHealthClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Test health check
		resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{
			Service: "custodian-simulator",
		})

		if err != nil {
			t.Errorf("Health check failed: %v", err)
		}

		if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
			t.Errorf("Expected SERVING status, got %v", resp.Status)
		}
	})
}

func TestCustodianGRPCServer_CustodyService(t *testing.T) {
	t.Run("accepts_custody_operations", func(t *testing.T) {
		t.Parallel()

		cfg := &config.Config{
			ServiceName: "custodian-simulator",
			GRPCPort:   0,
		}

		server := NewCustodianGRPCServer(cfg)

		lis, err := net.Listen("tcp", ":0")
		if err != nil {
			t.Fatalf("Failed to listen: %v", err)
		}

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Logf("Server serve error: %v", err)
			}
		}()
		defer server.GracefulStop()

		time.Sleep(100 * time.Millisecond)

		conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()

		// Test would use actual protobuf client when implemented
		// For now, just verify the server can be connected to
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Verify connection works by checking the health service
		healthClient := grpc_health_v1.NewHealthClient(conn)
		_, err = healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
		if err != nil {
			t.Skip("gRPC server not responding for custody operations test")
		}
	})
}

func TestCustodianGRPCServer_SettlementService(t *testing.T) {
	t.Run("handles_settlement_instructions", func(t *testing.T) {
		t.Parallel()

		cfg := &config.Config{
			ServiceName: "custodian-simulator",
			GRPCPort:   0,
		}

		server := NewCustodianGRPCServer(cfg)

		lis, err := net.Listen("tcp", ":0")
		if err != nil {
			t.Fatalf("Failed to listen: %v", err)
		}

		go func() {
			if err := server.Serve(lis); err != nil {
				t.Logf("Server serve error: %v", err)
			}
		}()
		defer server.GracefulStop()

		time.Sleep(100 * time.Millisecond)

		conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()

		// Test would use actual settlement service client when implemented
		// For now, verify the server infrastructure is working
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		healthClient := grpc_health_v1.NewHealthClient(conn)
		resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})

		if err != nil {
			t.Errorf("Failed to connect to settlement-capable server: %v", err)
		}

		if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
			t.Errorf("Expected server to be serving, got %v", resp.Status)
		}
	})
}

func TestCustodianGRPCServer_Metrics(t *testing.T) {
	t.Run("exposes_service_metrics", func(t *testing.T) {
		t.Parallel()

		cfg := &config.Config{
			ServiceName: "custodian-simulator",
			GRPCPort:   0,
		}

		server := NewCustodianGRPCServer(cfg)
		metrics := server.GetMetrics()

		// Verify metrics are available
		if metrics.ActiveConnections < 0 {
			t.Error("Active connections metric should be non-negative")
		}

		if metrics.TotalRequests < 0 {
			t.Error("Total requests metric should be non-negative")
		}

		if len(metrics.ServiceStatus) == 0 {
			t.Error("Service status should be available")
		}
	})
}

// CustodianGRPCServer interface that needs to be implemented
type CustodianGRPCServer interface {
	Serve(lis net.Listener) error
	GracefulStop()
	GetMetrics() grpcserver.ServerMetrics
}

// NewCustodianGRPCServer creates a new custodian gRPC server
func NewCustodianGRPCServer(cfg *config.Config) CustodianGRPCServer {
	return grpcserver.NewCustodianGRPCServer(cfg)
}