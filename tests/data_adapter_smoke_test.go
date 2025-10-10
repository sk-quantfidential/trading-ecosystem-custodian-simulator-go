//go:build integration

package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/quantfidential/trading-ecosystem/custodian-data-adapter-go/pkg/adapters"
)

// TestDataAdapterSmoke performs basic smoke tests for the DataAdapter integration
func TestDataAdapterSmoke(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup: Use orchestrator credentials
	os.Setenv("POSTGRES_URL", "postgres://custodian_adapter:custodian-adapter-db-pass@localhost:5432/trading_ecosystem?sslmode=disable")
	os.Setenv("REDIS_URL", "redis://custodian-adapter:custodian-pass@localhost:6379/0")
	defer os.Unsetenv("POSTGRES_URL")
	defer os.Unsetenv("REDIS_URL")

	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	t.Run("adapter_initialization", func(t *testing.T) {
		// Given: DataAdapter factory
		adapter, err := adapters.NewCustodianDataAdapterFromEnv(logger)
		if err != nil {
			t.Skipf("DataAdapter creation failed (infrastructure not available): %v", err)
			return
		}

		// When: Connecting to infrastructure
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = adapter.Connect(ctx)
		if err != nil {
			t.Skipf("DataAdapter connection failed (infrastructure not available): %v", err)
			return
		}
		defer adapter.Disconnect(ctx)

		// Then: Repositories should be accessible
		if adapter.PositionRepository() == nil {
			t.Error("Expected PositionRepository to be non-nil")
		}
		if adapter.SettlementRepository() == nil {
			t.Error("Expected SettlementRepository to be non-nil")
		}
		if adapter.BalanceRepository() == nil {
			t.Error("Expected BalanceRepository to be non-nil")
		}
		if adapter.ServiceDiscoveryRepository() == nil {
			t.Error("Expected ServiceDiscoveryRepository to be non-nil")
		}
		if adapter.CacheRepository() == nil {
			t.Error("Expected CacheRepository to be non-nil")
		}

		t.Log("✓ DataAdapter initialized successfully")
	})

	t.Run("position_repository_basic_crud", func(t *testing.T) {
		// Position repository requires UUID generation enhancement - deferred to future epic
		t.Skip("Position repository requires UUID generation enhancement - deferred to future epic")
	})

	t.Run("settlement_repository_basic_crud", func(t *testing.T) {
		// Requires position creation - deferred to future epic
		t.Skip("Settlement repository tests require UUID generation enhancement - deferred to future epic")
	})

	t.Run("balance_repository_basic_crud", func(t *testing.T) {
		// Requires account/position creation - deferred to future epic
		t.Skip("Balance repository tests require UUID generation enhancement - deferred to future epic")
	})

	t.Run("service_discovery_smoke", func(t *testing.T) {
		// Requires Redis ACL permissions (keys, scan) for custodian-adapter user
		t.Skip("Service discovery requires Redis ACL enhancement - deferred to future epic")
	})

	t.Run("cache_repository_smoke", func(t *testing.T) {
		// Given: Connected DataAdapter
		adapter, err := adapters.NewCustodianDataAdapterFromEnv(logger)
		if err != nil {
			t.Skipf("DataAdapter not available: %v", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := adapter.Connect(ctx); err != nil {
			t.Skipf("Infrastructure not available: %v", err)
			return
		}
		defer adapter.Disconnect(ctx)

		// When: Setting a cache value
		cacheRepo := adapter.CacheRepository()
		testKey := "custodian:smoke-test:" + time.Now().Format("20060102150405")
		testValue := "test-value-123"

		err = cacheRepo.Set(ctx, testKey, testValue, 1*time.Minute)
		if err != nil {
			t.Fatalf("Failed to set cache value: %v", err)
		}
		defer cacheRepo.Delete(ctx, testKey)

		// Then: Should be able to retrieve it
		retrieved, err := cacheRepo.Get(ctx, testKey)
		if err != nil {
			t.Fatalf("Failed to get cache value: %v", err)
		}

		if retrieved != testValue {
			t.Errorf("Expected value '%s', got '%s'", testValue, retrieved)
		}

		t.Logf("✓ Cache smoke test passed (key: %s)", testKey)
	})
}
