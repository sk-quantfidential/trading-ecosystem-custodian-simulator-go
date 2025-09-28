package infrastructure

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/quantfidential/trading-ecosystem/custodian-simulator-go/internal/config"
)

type DefaultInterServiceClientManager struct {
	config          *config.Config
	serviceDiscovery ServiceDiscoveryInterface
	configClient    ConfigurationClientInterface
	logger          *logrus.Logger

	// Connection management
	connections      map[string]*grpc.ClientConn
	connectionsMutex sync.RWMutex

	// Connection pooling and circuit breaker
	connectionPool map[string]*ConnectionPool
	poolMutex      sync.RWMutex

	// Statistics
	activeConnections int64
	totalConnections  int64
	failedConnections int64
	statsMutex        sync.RWMutex
}

type ConnectionPool struct {
	connections []*grpc.ClientConn
	index       int
	mutex       sync.Mutex
	maxSize     int
}

type ConnectionStats struct {
	ActiveConnections int64 `json:"active_connections"`
	TotalConnections  int64 `json:"total_connections"`
	FailedConnections int64 `json:"failed_connections"`
}

// Service client interfaces
type ServiceClientInterface interface {
	HealthCheck(ctx context.Context) (HealthStatus, error)
}

type ExchangeSimulatorClientInterface interface {
	ServiceClientInterface
	GetTradingStatus(ctx context.Context) (TradingStatus, error)
}

type AuditCorrelatorClientInterface interface {
	ServiceClientInterface
	GetAuditMetrics(ctx context.Context) (AuditMetrics, error)
}

// Response types
type HealthStatus struct {
	Status      string    `json:"status"`
	LastChecked time.Time `json:"last_checked"`
	Details     string    `json:"details"`
}

type TradingStatus struct {
	ActiveTrades  int       `json:"active_trades"`
	TotalVolume   int64     `json:"total_volume"`
	LastTradeTime time.Time `json:"last_trade_time"`
}

type AuditMetrics struct {
	TotalEvents      int64     `json:"total_events"`
	CorrelatedEvents int64     `json:"correlated_events"`
	LastUpdated      time.Time `json:"last_updated"`
}

// Service discovery and configuration interfaces
type ServiceDiscoveryInterface interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	DiscoverServices(ctx context.Context, serviceName string) ([]ServiceInfo, error)
}

type ConfigurationClientInterface interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	GetConfiguration(ctx context.Context, key string) (ConfigurationValue, error)
}

// Error types
type ServiceUnavailableError struct {
	ServiceName string
	Cause       error
}

func (e *ServiceUnavailableError) Error() string {
	return fmt.Sprintf("service '%s' is unavailable: %v", e.ServiceName, e.Cause)
}

func NewInterServiceClientManager(cfg *config.Config) *DefaultInterServiceClientManager {
	logger := logrus.New()
	logger.SetLevel(getLogLevel(cfg.LogLevel))

	return &DefaultInterServiceClientManager{
		config:         cfg,
		logger:         logger,
		connections:    make(map[string]*grpc.ClientConn),
		connectionPool: make(map[string]*ConnectionPool),
	}
}

func (cm *DefaultInterServiceClientManager) Initialize(ctx context.Context) error {
	// Initialize service discovery
	cm.serviceDiscovery = NewServiceDiscovery(cm.config)
	if err := cm.serviceDiscovery.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect service discovery: %w", err)
	}

	// Initialize configuration client
	cm.configClient = NewConfigurationClient(cm.config)
	if err := cm.configClient.Connect(ctx); err != nil {
		cm.logger.WithError(err).Warn("Configuration client connection failed, continuing without it")
	}

	cm.logger.Info("Inter-service client manager initialized")
	return nil
}

func (cm *DefaultInterServiceClientManager) Cleanup(ctx context.Context) error {
	cm.connectionsMutex.Lock()
	defer cm.connectionsMutex.Unlock()

	// Close all connections
	for serviceName, conn := range cm.connections {
		if err := conn.Close(); err != nil {
			cm.logger.WithError(err).WithField("service", serviceName).Warn("Failed to close connection")
		}
	}
	cm.connections = make(map[string]*grpc.ClientConn)

	// Cleanup connection pools
	cm.poolMutex.Lock()
	for _, pool := range cm.connectionPool {
		for _, conn := range pool.connections {
			if err := conn.Close(); err != nil {
				cm.logger.WithError(err).Warn("Failed to close pooled connection")
			}
		}
	}
	cm.connectionPool = make(map[string]*ConnectionPool)
	cm.poolMutex.Unlock()

	// Disconnect from infrastructure services
	if cm.serviceDiscovery != nil {
		if err := cm.serviceDiscovery.Disconnect(ctx); err != nil {
			cm.logger.WithError(err).Warn("Failed to disconnect service discovery")
		}
	}

	if cm.configClient != nil {
		if err := cm.configClient.Disconnect(ctx); err != nil {
			cm.logger.WithError(err).Warn("Failed to disconnect configuration client")
		}
	}

	cm.logger.Info("Inter-service client manager cleaned up")
	return nil
}

func (cm *DefaultInterServiceClientManager) GetExchangeSimulatorClient(ctx context.Context) (ExchangeSimulatorClientInterface, error) {
	conn, err := cm.getServiceConnection(ctx, "exchange-simulator")
	if err != nil {
		return nil, err
	}

	return &ExchangeSimulatorClient{
		conn:   conn,
		logger: cm.logger,
	}, nil
}

func (cm *DefaultInterServiceClientManager) GetAuditCorrelatorClient(ctx context.Context) (AuditCorrelatorClientInterface, error) {
	conn, err := cm.getServiceConnection(ctx, "audit-correlator")
	if err != nil {
		return nil, err
	}

	return &AuditCorrelatorClient{
		conn:   conn,
		logger: cm.logger,
	}, nil
}

func (cm *DefaultInterServiceClientManager) GetClientByName(ctx context.Context, serviceName string) (ServiceClientInterface, error) {
	conn, err := cm.getServiceConnection(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	return &GenericServiceClient{
		conn:        conn,
		serviceName: serviceName,
		logger:      cm.logger,
	}, nil
}

func (cm *DefaultInterServiceClientManager) DiscoverServices(ctx context.Context) ([]ServiceInfo, error) {
	if cm.serviceDiscovery == nil {
		return nil, fmt.Errorf("service discovery not initialized")
	}

	// Discover all services by using a wildcard pattern
	allServices := make([]ServiceInfo, 0)

	// Look for common service types
	serviceTypes := []string{"exchange-simulator", "audit-correlator", "custodian-simulator", "risk-monitor"}

	for _, serviceType := range serviceTypes {
		services, err := cm.serviceDiscovery.DiscoverServices(ctx, serviceType)
		if err != nil {
			cm.logger.WithError(err).WithField("service_type", serviceType).Debug("Failed to discover services")
			continue
		}
		allServices = append(allServices, services...)
	}

	return allServices, nil
}

func (cm *DefaultInterServiceClientManager) GetConnectionStats() ConnectionStats {
	cm.statsMutex.RLock()
	defer cm.statsMutex.RUnlock()

	return ConnectionStats{
		ActiveConnections: cm.activeConnections,
		TotalConnections:  cm.totalConnections,
		FailedConnections: cm.failedConnections,
	}
}

func (cm *DefaultInterServiceClientManager) getServiceConnection(ctx context.Context, serviceName string) (*grpc.ClientConn, error) {
	cm.connectionsMutex.Lock()
	defer cm.connectionsMutex.Unlock()

	// Check if we already have a connection
	if conn, exists := cm.connections[serviceName]; exists {
		if conn.GetState() == connectivity.Ready || conn.GetState() == connectivity.Idle {
			return conn, nil
		}
		// Connection is not healthy, remove it
		conn.Close()
		delete(cm.connections, serviceName)
	}

	// Discover service
	services, err := cm.serviceDiscovery.DiscoverServices(ctx, serviceName)
	if err != nil {
		cm.incrementFailedConnections()
		return nil, &ServiceUnavailableError{
			ServiceName: serviceName,
			Cause:       fmt.Errorf("service discovery failed: %w", err),
		}
	}

	if len(services) == 0 {
		cm.incrementFailedConnections()
		return nil, &ServiceUnavailableError{
			ServiceName: serviceName,
			Cause:       fmt.Errorf("no instances found"),
		}
	}

	// Use the first available service
	service := services[0]
	target := fmt.Sprintf("%s:%d", service.Host, service.GRPCPort)

	// Create new connection
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cm.incrementFailedConnections()
		return nil, &ServiceUnavailableError{
			ServiceName: serviceName,
			Cause:       fmt.Errorf("failed to connect to %s: %w", target, err),
		}
	}

	// Store connection
	cm.connections[serviceName] = conn
	cm.incrementActiveConnections()
	cm.incrementTotalConnections()

	cm.logger.WithFields(logrus.Fields{
		"service": serviceName,
		"target":  target,
	}).Info("Established gRPC connection")

	return conn, nil
}

func (cm *DefaultInterServiceClientManager) incrementActiveConnections() {
	cm.statsMutex.Lock()
	defer cm.statsMutex.Unlock()
	cm.activeConnections++
}

func (cm *DefaultInterServiceClientManager) incrementTotalConnections() {
	cm.statsMutex.Lock()
	defer cm.statsMutex.Unlock()
	cm.totalConnections++
}

func (cm *DefaultInterServiceClientManager) incrementFailedConnections() {
	cm.statsMutex.Lock()
	defer cm.statsMutex.Unlock()
	cm.failedConnections++
}

// Concrete client implementations
type ExchangeSimulatorClient struct {
	conn   *grpc.ClientConn
	logger *logrus.Logger
}

func (c *ExchangeSimulatorClient) HealthCheck(ctx context.Context) (HealthStatus, error) {
	client := grpc_health_v1.NewHealthClient(c.conn)
	resp, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: "exchange-simulator",
	})
	if err != nil {
		return HealthStatus{}, err
	}

	status := "unknown"
	switch resp.Status {
	case grpc_health_v1.HealthCheckResponse_SERVING:
		status = "healthy"
	case grpc_health_v1.HealthCheckResponse_NOT_SERVING:
		status = "unhealthy"
	}

	return HealthStatus{
		Status:      status,
		LastChecked: time.Now(),
		Details:     "Exchange simulator health check",
	}, nil
}

func (c *ExchangeSimulatorClient) GetTradingStatus(ctx context.Context) (TradingStatus, error) {
	// This would be a real gRPC call in production
	// For now, return mock data
	return TradingStatus{
		ActiveTrades:  10,
		TotalVolume:   1000000,
		LastTradeTime: time.Now(),
	}, nil
}

type AuditCorrelatorClient struct {
	conn   *grpc.ClientConn
	logger *logrus.Logger
}

func (c *AuditCorrelatorClient) HealthCheck(ctx context.Context) (HealthStatus, error) {
	client := grpc_health_v1.NewHealthClient(c.conn)
	resp, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: "audit-correlator",
	})
	if err != nil {
		return HealthStatus{}, err
	}

	status := "unknown"
	switch resp.Status {
	case grpc_health_v1.HealthCheckResponse_SERVING:
		status = "healthy"
	case grpc_health_v1.HealthCheckResponse_NOT_SERVING:
		status = "unhealthy"
	}

	return HealthStatus{
		Status:      status,
		LastChecked: time.Now(),
		Details:     "Audit correlator health check",
	}, nil
}

func (c *AuditCorrelatorClient) GetAuditMetrics(ctx context.Context) (AuditMetrics, error) {
	// This would be a real gRPC call in production
	// For now, return mock data
	return AuditMetrics{
		TotalEvents:      5000,
		CorrelatedEvents: 4800,
		LastUpdated:      time.Now(),
	}, nil
}

type GenericServiceClient struct {
	conn        *grpc.ClientConn
	serviceName string
	logger      *logrus.Logger
}

func (c *GenericServiceClient) HealthCheck(ctx context.Context) (HealthStatus, error) {
	client := grpc_health_v1.NewHealthClient(c.conn)
	resp, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: c.serviceName,
	})
	if err != nil {
		return HealthStatus{}, err
	}

	status := "unknown"
	switch resp.Status {
	case grpc_health_v1.HealthCheckResponse_SERVING:
		status = "healthy"
	case grpc_health_v1.HealthCheckResponse_NOT_SERVING:
		status = "unhealthy"
	}

	return HealthStatus{
		Status:      status,
		LastChecked: time.Now(),
		Details:     fmt.Sprintf("%s health check", c.serviceName),
	}, nil
}

// Error checking functions
func IsServiceUnavailableError(err error) bool {
	_, ok := err.(*ServiceUnavailableError)
	return ok
}