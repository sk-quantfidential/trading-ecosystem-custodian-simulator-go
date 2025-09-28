//go:build unit

package internal

import (
	"context"
	"testing"
	"time"

	"github.com/quantfidential/trading-ecosystem/custodian-simulator-go/internal/config"
	"github.com/quantfidential/trading-ecosystem/custodian-simulator-go/internal/infrastructure"
)

// TestConfigurationClient_RedPhase defines the expected behaviors for configuration service integration
// These tests will fail initially and drive our implementation (TDD Red-Green-Refactor)
func TestConfigurationClient_GetConfiguration(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		expectedType infrastructure.ConfigValueType
		wantErr      bool
	}{
		{
			name:         "settlement_timeout_hours",
			key:          "custodian.settlement.timeout_hours",
			expectedType: infrastructure.ConfigValueTypeNumber,
			wantErr:      false,
		},
		{
			name:         "multi_asset_enabled",
			key:          "custodian.assets.multi_asset_enabled",
			expectedType: infrastructure.ConfigValueTypeBoolean,
			wantErr:      false,
		},
		{
			name:         "custody_backend",
			key:          "custodian.storage.backend",
			expectedType: infrastructure.ConfigValueTypeString,
			wantErr:      false,
		},
		{
			name:    "invalid_key",
			key:     "nonexistent.key",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := NewConfigurationClient(&config.Config{
				ConfigurationServiceURL: "http://localhost:8090",
				RequestTimeout:          5 * time.Second,
				CacheTTL:               5 * time.Minute,
			})

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := client.Connect(ctx)
			if err != nil {
				t.Skip("Configuration service not available for test")
			}
			defer client.Disconnect(ctx)

			value, err := client.GetConfiguration(ctx, tt.key)

			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigurationClient.GetConfiguration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if value.Key != tt.key {
					t.Errorf("Expected key %s, got %s", tt.key, value.Key)
				}

				if value.Type != tt.expectedType {
					t.Errorf("Expected type %v, got %v", tt.expectedType, value.Type)
				}
			}
		})
	}
}

func TestConfigurationClient_Caching(t *testing.T) {
	t.Run("caches_configuration_values", func(t *testing.T) {
		t.Parallel()

		client := NewConfigurationClient(&config.Config{
			ConfigurationServiceURL: "http://localhost:8090",
			RequestTimeout:          5 * time.Second,
			CacheTTL:               300 * time.Second,
		})

		ctx := context.Background()
		err := client.Connect(ctx)
		if err != nil {
			t.Skip("Configuration service not available for test")
		}
		defer client.Disconnect(ctx)

		// First call - cache miss
		_, err = client.GetConfiguration(ctx, "custodian.settlement.timeout_hours")
		if err != nil {
			t.Skip("Configuration key not available for caching test")
		}

		// Second call - should be cache hit
		_, err = client.GetConfiguration(ctx, "custodian.settlement.timeout_hours")
		if err != nil {
			t.Errorf("Cached configuration retrieval failed: %v", err)
		}

		// Verify cache statistics
		stats := client.GetCacheStats()
		if stats.CacheHits == 0 {
			t.Error("Expected cache hits after second call")
		}
	})
}

func TestConfigurationClient_TypeConversions(t *testing.T) {
	tests := []struct {
		name        string
		configValue infrastructure.ConfigurationValue
		testFunc    func(t *testing.T, value infrastructure.ConfigurationValue)
	}{
		{
			name: "string_conversion",
			configValue: infrastructure.ConfigurationValue{
				Key:   "test.string",
				Value: "custodian-simulator",
				Type:  infrastructure.ConfigValueTypeString,
			},
			testFunc: func(t *testing.T, value infrastructure.ConfigurationValue) {
				result := value.AsString()
				if result != "custodian-simulator" {
					t.Errorf("Expected 'custodian-simulator', got '%s'", result)
				}
			},
		},
		{
			name: "number_conversion",
			configValue: infrastructure.ConfigurationValue{
				Key:   "test.number",
				Value: "24",
				Type:  infrastructure.ConfigValueTypeNumber,
			},
			testFunc: func(t *testing.T, value infrastructure.ConfigurationValue) {
				result, err := value.AsInt()
				if err != nil {
					t.Errorf("AsInt() failed: %v", err)
				}
				if result != 24 {
					t.Errorf("Expected 24, got %d", result)
				}
			},
		},
		{
			name: "boolean_conversion",
			configValue: infrastructure.ConfigurationValue{
				Key:   "test.boolean",
				Value: "true",
				Type:  infrastructure.ConfigValueTypeBoolean,
			},
			testFunc: func(t *testing.T, value infrastructure.ConfigurationValue) {
				result, err := value.AsBool()
				if err != nil {
					t.Errorf("AsBool() failed: %v", err)
				}
				if !result {
					t.Error("Expected true, got false")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.testFunc(t, tt.configValue)
		})
	}
}

// ConfigurationClient interface that needs to be implemented
type ConfigurationClient interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	GetConfiguration(ctx context.Context, key string) (infrastructure.ConfigurationValue, error)
	GetCacheStats() infrastructure.CacheStats
}

// Constructor function that creates a new configuration client
func NewConfigurationClient(cfg *config.Config) ConfigurationClient {
	return infrastructure.NewConfigurationClient(cfg)
}