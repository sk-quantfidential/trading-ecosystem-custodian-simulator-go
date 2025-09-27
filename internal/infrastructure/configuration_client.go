package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/quantfidential/trading-ecosystem/custodian-simulator-go/internal/config"
)

type HTTPConfigurationClient struct {
	config     *config.Config
	httpClient *http.Client
	logger     *logrus.Logger

	// Cache management
	cache      map[string]cachedValue
	cacheMutex sync.RWMutex

	// Statistics
	cacheHits   int64
	cacheMisses int64
	statsMutex  sync.RWMutex
}

type cachedValue struct {
	value     ConfigurationValue
	expiresAt time.Time
}

type ConfigValueType int

const (
	ConfigValueTypeString ConfigValueType = iota
	ConfigValueTypeNumber
	ConfigValueTypeBoolean
	ConfigValueTypeJSON
)

type ConfigurationValue struct {
	Key         string          `json:"key"`
	Value       string          `json:"value"`
	Type        ConfigValueType `json:"type"`
	Environment string          `json:"environment"`
	LastUpdated time.Time       `json:"last_updated"`
}

type CacheStats struct {
	CacheHits   int64   `json:"cache_hits"`
	CacheMisses int64   `json:"cache_misses"`
	CacheSize   int     `json:"cache_size"`
	HitRate     float64 `json:"hit_rate"`
}

func NewConfigurationClient(cfg *config.Config) *HTTPConfigurationClient {
	logger := logrus.New()
	logger.SetLevel(getLogLevel(cfg.LogLevel))

	return &HTTPConfigurationClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.RequestTimeout,
		},
		logger: logger,
		cache:  make(map[string]cachedValue),
	}
}

func (c *HTTPConfigurationClient) Connect(ctx context.Context) error {
	// Test connection to configuration service
	testURL := fmt.Sprintf("%s/health", c.config.ConfigurationServiceURL)

	req, err := http.NewRequestWithContext(ctx, "GET", testURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to configuration service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("configuration service health check failed with status: %d", resp.StatusCode)
	}

	c.logger.WithField("url", c.config.ConfigurationServiceURL).Info("Connected to configuration service")
	return nil
}

func (c *HTTPConfigurationClient) Disconnect(ctx context.Context) error {
	c.logger.Info("Disconnecting from configuration service")

	// Clear cache on disconnect
	c.cacheMutex.Lock()
	c.cache = make(map[string]cachedValue)
	c.cacheMutex.Unlock()

	return nil
}

func (c *HTTPConfigurationClient) GetConfiguration(ctx context.Context, key string) (ConfigurationValue, error) {
	// Check cache first
	if cached, found := c.getCachedValue(key); found {
		c.incrementCacheHits()
		c.logger.WithField("key", key).Debug("Configuration retrieved from cache")
		return cached, nil
	}

	c.incrementCacheMisses()

	// Fetch from service
	value, err := c.fetchConfiguration(ctx, key)
	if err != nil {
		return ConfigurationValue{}, err
	}

	// Cache the value
	c.cacheValue(key, value)

	c.logger.WithField("key", key).Info("Configuration retrieved from service")
	return value, nil
}

func (c *HTTPConfigurationClient) GetCacheStats() CacheStats {
	c.statsMutex.RLock()
	defer c.statsMutex.RUnlock()

	c.cacheMutex.RLock()
	cacheSize := len(c.cache)
	c.cacheMutex.RUnlock()

	total := c.cacheHits + c.cacheMisses
	var hitRate float64
	if total > 0 {
		hitRate = float64(c.cacheHits) / float64(total)
	}

	return CacheStats{
		CacheHits:   c.cacheHits,
		CacheMisses: c.cacheMisses,
		CacheSize:   cacheSize,
		HitRate:     hitRate,
	}
}

func (c *HTTPConfigurationClient) fetchConfiguration(ctx context.Context, key string) (ConfigurationValue, error) {
	endpoint := fmt.Sprintf("%s/configuration", c.config.ConfigurationServiceURL)

	// Add key as query parameter
	u, err := url.Parse(endpoint)
	if err != nil {
		return ConfigurationValue{}, fmt.Errorf("invalid endpoint URL: %w", err)
	}

	q := u.Query()
	q.Set("key", key)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return ConfigurationValue{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ConfigurationValue{}, fmt.Errorf("failed to fetch configuration: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ConfigurationValue{}, fmt.Errorf("configuration key '%s' not found", key)
	}

	if resp.StatusCode != http.StatusOK {
		return ConfigurationValue{}, fmt.Errorf("configuration service returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ConfigurationValue{}, fmt.Errorf("failed to read response: %w", err)
	}

	var value ConfigurationValue
	if err := json.Unmarshal(body, &value); err != nil {
		return ConfigurationValue{}, fmt.Errorf("failed to parse response: %w", err)
	}

	return value, nil
}

func (c *HTTPConfigurationClient) getCachedValue(key string) (ConfigurationValue, bool) {
	c.cacheMutex.RLock()
	defer c.cacheMutex.RUnlock()

	cached, exists := c.cache[key]
	if !exists {
		return ConfigurationValue{}, false
	}

	if time.Now().After(cached.expiresAt) {
		// Value expired
		return ConfigurationValue{}, false
	}

	return cached.value, true
}

func (c *HTTPConfigurationClient) cacheValue(key string, value ConfigurationValue) {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	c.cache[key] = cachedValue{
		value:     value,
		expiresAt: time.Now().Add(c.config.CacheTTL),
	}
}

func (c *HTTPConfigurationClient) incrementCacheHits() {
	c.statsMutex.Lock()
	defer c.statsMutex.Unlock()
	c.cacheHits++
}

func (c *HTTPConfigurationClient) incrementCacheMisses() {
	c.statsMutex.Lock()
	defer c.statsMutex.Unlock()
	c.cacheMisses++
}

// Helper methods for ConfigurationValue type conversions
func (cv ConfigurationValue) AsString() string {
	return cv.Value
}

func (cv ConfigurationValue) AsInt() (int, error) {
	return strconv.Atoi(cv.Value)
}

func (cv ConfigurationValue) AsBool() (bool, error) {
	return strconv.ParseBool(cv.Value)
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