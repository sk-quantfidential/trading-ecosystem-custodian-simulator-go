package grpc

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/quantfidential/trading-ecosystem/custodian-simulator-go/internal/config"
	"github.com/quantfidential/trading-ecosystem/custodian-simulator-go/internal/services"
)

type CustodianGRPCServer struct {
	config            *config.Config
	server            *grpc.Server
	healthSrv         *health.Server
	custodianSvc      *services.CustodianService
	logger            *logrus.Logger
	startTime         time.Time
	activeConnections int64
	totalRequests     int64
	mutex             sync.RWMutex
}

type ServerMetrics struct {
	ActiveConnections int64             `json:"active_connections"`
	TotalRequests     int64             `json:"total_requests"`
	ServiceStatus     map[string]string `json:"service_status"`
	Uptime            time.Duration     `json:"uptime"`
}

func NewCustodianGRPCServer(cfg *config.Config) *CustodianGRPCServer {
	logger := logrus.New()
	logger.SetLevel(getLogLevel(cfg.LogLevel))

	server := grpc.NewServer(
		grpc.UnaryInterceptor(requestMetricsInterceptor),
	)

	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthSrv)

	custodianSvc := services.NewCustodianService(cfg, logger)

	grpcServer := &CustodianGRPCServer{
		config:       cfg,
		server:       server,
		healthSrv:    healthSrv,
		custodianSvc: custodianSvc,
		logger:       logger,
		startTime:    time.Now(),
	}

	// Set health status for services
	healthSrv.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	healthSrv.SetServingStatus(cfg.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
	healthSrv.SetServingStatus("custodian", grpc_health_v1.HealthCheckResponse_SERVING)
	healthSrv.SetServingStatus("settlement", grpc_health_v1.HealthCheckResponse_SERVING)

	logger.WithFields(logrus.Fields{
		"service": cfg.ServiceName,
		"version": cfg.ServiceVersion,
		"port":    cfg.GRPCPort,
	}).Info("Custodian gRPC server initialized")

	return grpcServer
}

func (s *CustodianGRPCServer) Serve(lis net.Listener) error {
	s.logger.WithField("address", lis.Addr().String()).Info("Starting custodian gRPC server")
	return s.server.Serve(lis)
}

func (s *CustodianGRPCServer) GracefulStop() {
	s.logger.Info("Gracefully stopping custodian gRPC server")

	// Set health status to NOT_SERVING before shutdown
	s.healthSrv.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	s.healthSrv.SetServingStatus(s.config.ServiceName, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	s.healthSrv.SetServingStatus("custodian", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	s.healthSrv.SetServingStatus("settlement", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	s.server.GracefulStop()
	s.logger.Info("Custodian gRPC server stopped")
}

func (s *CustodianGRPCServer) GetMetrics() ServerMetrics {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return ServerMetrics{
		ActiveConnections: s.activeConnections,
		TotalRequests:     s.totalRequests,
		ServiceStatus: map[string]string{
			"custodian":  "serving",
			"settlement": "serving",
			"health":     "serving",
		},
		Uptime: time.Since(s.startTime),
	}
}

func requestMetricsInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// This would be enhanced to track actual metrics
	// For now, it's a placeholder that satisfies the interface
	return handler(ctx, req)
}

func getLogLevel(level string) logrus.Level {
	switch level {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}