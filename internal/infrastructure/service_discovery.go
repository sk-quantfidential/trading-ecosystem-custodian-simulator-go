package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"github.com/quantfidential/trading-ecosystem/custodian-simulator-go/internal/config"
)

type RedisServiceDiscovery struct {
	config            *config.Config
	redisClient       *redis.Client
	logger            *logrus.Logger
	serviceInfo       ServiceInfo
	heartbeatInterval time.Duration
	stopHeartbeat     chan struct{}
}

type ServiceInfo struct {
	Name     string    `json:"name"`
	Version  string    `json:"version"`
	Host     string    `json:"host"`
	GRPCPort int       `json:"grpc_port"`
	HTTPPort int       `json:"http_port"`
	Status   string    `json:"status"`
	LastSeen time.Time `json:"last_seen"`
}

const (
	serviceKeyPrefix = "services:"
	serviceKeyTTL    = 90 * time.Second // TTL for service registry entries
)

func NewServiceDiscovery(cfg *config.Config) *RedisServiceDiscovery {
	logger := logrus.New()
	logger.SetLevel(getLogLevel(cfg.LogLevel))

	serviceInfo := ServiceInfo{
		Name:     cfg.ServiceName,
		Version:  cfg.ServiceVersion,
		Host:     "localhost", // Default for local development
		GRPCPort: cfg.GRPCPort,
		HTTPPort: cfg.HTTPPort,
		Status:   "starting",
		LastSeen: time.Now(),
	}

	return &RedisServiceDiscovery{
		config:            cfg,
		logger:            logger,
		serviceInfo:       serviceInfo,
		heartbeatInterval: cfg.HealthCheckInterval,
		stopHeartbeat:     make(chan struct{}),
	}
}

func (sd *RedisServiceDiscovery) Connect(ctx context.Context) error {
	opt, err := redis.ParseURL(sd.config.RedisURL)
	if err != nil {
		return fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	sd.redisClient = redis.NewClient(opt)

	// Test connection
	if err := sd.redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	sd.logger.WithField("redis_url", sd.config.RedisURL).Info("Connected to Redis for service discovery")
	return nil
}

func (sd *RedisServiceDiscovery) Disconnect(ctx context.Context) error {
	// Stop heartbeat if running
	select {
	case sd.stopHeartbeat <- struct{}{}:
	default:
	}

	// Unregister service
	if sd.redisClient != nil {
		serviceKey := fmt.Sprintf("%s%s:%s:%d", serviceKeyPrefix, sd.serviceInfo.Name, sd.serviceInfo.Host, sd.serviceInfo.GRPCPort)
		if err := sd.redisClient.Del(ctx, serviceKey).Err(); err != nil {
			sd.logger.WithError(err).Warn("Failed to unregister service during disconnect")
		}

		if err := sd.redisClient.Close(); err != nil {
			sd.logger.WithError(err).Warn("Failed to close Redis connection")
		}
	}

	sd.logger.Info("Disconnected from service discovery")
	return nil
}

func (sd *RedisServiceDiscovery) RegisterService(ctx context.Context) error {
	sd.serviceInfo.Status = "healthy"
	sd.serviceInfo.LastSeen = time.Now()

	serviceKey := sd.getServiceKey()
	serviceData, err := json.Marshal(sd.serviceInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal service info: %w", err)
	}

	// Set with TTL
	if err := sd.redisClient.SetEx(ctx, serviceKey, serviceData, serviceKeyTTL).Err(); err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	sd.logger.WithFields(logrus.Fields{
		"service_key": serviceKey,
		"ttl":         serviceKeyTTL,
	}).Info("Service registered successfully")

	return nil
}

func (sd *RedisServiceDiscovery) DiscoverServices(ctx context.Context, serviceName string) ([]ServiceInfo, error) {
	pattern := fmt.Sprintf("%s%s:*", serviceKeyPrefix, serviceName)
	keys, err := sd.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to discover services: %w", err)
	}

	var services []ServiceInfo
	for _, key := range keys {
		serviceData, err := sd.redisClient.Get(ctx, key).Result()
		if err != nil {
			sd.logger.WithError(err).WithField("key", key).Warn("Failed to get service data")
			continue
		}

		var service ServiceInfo
		if err := json.Unmarshal([]byte(serviceData), &service); err != nil {
			sd.logger.WithError(err).WithField("key", key).Warn("Failed to unmarshal service data")
			continue
		}

		services = append(services, service)
	}

	sd.logger.WithFields(logrus.Fields{
		"service_name":   serviceName,
		"services_found": len(services),
	}).Debug("Services discovered")

	return services, nil
}

func (sd *RedisServiceDiscovery) StartHeartbeat(ctx context.Context) {
	ticker := time.NewTicker(sd.heartbeatInterval)
	defer ticker.Stop()

	sd.logger.WithField("interval", sd.heartbeatInterval).Info("Starting service heartbeat")

	for {
		select {
		case <-ctx.Done():
			sd.logger.Info("Heartbeat stopped due to context cancellation")
			return
		case <-sd.stopHeartbeat:
			sd.logger.Info("Heartbeat stopped")
			return
		case <-ticker.C:
			if err := sd.sendHeartbeat(ctx); err != nil {
				sd.logger.WithError(err).Error("Failed to send heartbeat")
			}
		}
	}
}

func (sd *RedisServiceDiscovery) sendHeartbeat(ctx context.Context) error {
	sd.serviceInfo.LastSeen = time.Now()

	serviceKey := sd.getServiceKey()
	serviceData, err := json.Marshal(sd.serviceInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal service info for heartbeat: %w", err)
	}

	// Update with fresh TTL
	if err := sd.redisClient.SetEx(ctx, serviceKey, serviceData, serviceKeyTTL).Err(); err != nil {
		return fmt.Errorf("failed to send heartbeat: %w", err)
	}

	sd.logger.WithField("service_key", serviceKey).Debug("Heartbeat sent")
	return nil
}

func (sd *RedisServiceDiscovery) getServiceKey() string {
	return fmt.Sprintf("%s%s:%s:%d", serviceKeyPrefix, sd.serviceInfo.Name, sd.serviceInfo.Host, sd.serviceInfo.GRPCPort)
}

func (sd *RedisServiceDiscovery) GetServiceInfo() ServiceInfo {
	return sd.serviceInfo
}

func (sd *RedisServiceDiscovery) UpdateServiceStatus(status string) {
	sd.serviceInfo.Status = status
	sd.logger.WithField("status", status).Info("Service status updated")
}