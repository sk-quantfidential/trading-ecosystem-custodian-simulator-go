# custodian-simulator-go TODO

## epic-TSE-0001: Foundation Services & Infrastructure

### 🏗️ Milestone TSE-0001.1a: Go Services Bootstrapping
**Status**: ✅ COMPLETED
**Priority**: High

**Tasks**:
- [x] Create Go service directory structure following clean architecture
- [x] Implement health check endpoint (REST and gRPC)
- [x] Basic structured logging with levels
- [x] Error handling infrastructure
- [x] Dockerfile for service containerization
- [x] Load component-specific .claude configuration

**BDD Acceptance**: All Go services can start, respond to health checks, and shutdown gracefully

---

### 🔗 Milestone TSE-0001.3b: Go Services gRPC Integration
**Status**: ✅ **COMPLETED** (Following audit-correlator-go pattern)
**Priority**: High

**Tasks** (Following proven TDD Red-Green-Refactor cycle):
- [x] **Phase 1: TDD Red** - Create failing tests for all gRPC integration behaviors
- [x] **Phase 2: Infrastructure** - Add Redis dependencies and update .gitignore for Go projects
- [x] **Phase 3: gRPC Server** - Implement enhanced gRPC server with health service, metrics, and graceful shutdown
- [x] **Phase 4: Configuration** - Implement configuration service client with HTTP caching, TTL, and type conversion
- [x] **Phase 5: Discovery** - Implement service discovery with Redis-based registry, heartbeat, and cleanup
- [x] **Phase 6: Communication** - Create inter-service gRPC client manager with connection pooling and circuit breaker
- [x] **Phase 7: Integration** - Implement comprehensive inter-service communication testing with smart skipping
- [x] **Phase 8: Validation** - Verify BDD acceptance and complete milestone documentation

**Implementation Pattern** (Replicating audit-correlator-go success):
- **Infrastructure Layer**: Configuration client, service discovery, gRPC clients
- **Presentation Layer**: Enhanced gRPC server with health service
- **Testing Strategy**: Unit tests with smart dependency skipping, integration tests for end-to-end scenarios
- **Error Handling**: Graceful degradation, circuit breaker patterns, comprehensive logging

**BDD Acceptance**: ✅ **VALIDATED** - Go services can discover and communicate with each other via gRPC

**Dependencies**: TSE-0001.1a (Go Services Bootstrapping), TSE-0001.3a (Core Infrastructure)

**🎯 CUSTODIAN-SIMULATOR-GO ACHIEVEMENTS**:
- ✅ **Enhanced gRPC Server**: Health service, metrics tracking, graceful shutdown with concurrent HTTP/gRPC operation
- ✅ **Service Discovery**: Redis-based registration with heartbeat, dynamic lookup, and proper cleanup
- ✅ **Configuration Client**: HTTP client with caching, TTL, type conversion, and performance statistics
- ✅ **Inter-Service Communication**: Connection pooling, circuit breaker pattern, and comprehensive error handling
- ✅ **Test Coverage**: 14 test cases (9 unit, 5 integration) with smart skipping when infrastructure unavailable
- ✅ **Production Ready**: Service builds and runs successfully with proper Redis integration
- ✅ **Pattern Replication**: Successfully replicated audit-correlator-go architecture and testing approach

**Reference Implementation**: audit-correlator-go (✅ COMPLETED) - Pattern successfully replicated

---

### 🏦 Milestone TSE-0001.6: Custodian Foundation (PRIMARY)
**Status**: Not Started
**Priority**: CRITICAL - Enables settlement and custody

**Tasks**:
- [ ] Account custody simulation (hold balances across assets)
- [ ] Settlement processing (T+0 immediate settlements initially)
- [ ] Transfer API (deposits/withdrawals between accounts)
- [ ] Balance reporting and reconciliation
- [ ] Multi-asset custody support (BTC, ETH, USD, USDT)
- [ ] Settlement instruction processing
- [ ] Custody audit trail
- [ ] Basic compliance checks

**BDD Acceptance**: Exchange settlements flow through to custodian accounts

**Dependencies**: TSE-0001.3b (Go Services gRPC Integration), TSE-0001.5b (Exchange Order Processing)

---

### 📈 Milestone TSE-0001.12b: Trading Flow Integration
**Status**: Not Started
**Priority**: Medium

**Tasks**:
- [ ] Integration with exchange settlement validation
- [ ] Multi-day settlement cycle testing
- [ ] Balance reconciliation across services
- [ ] Chaos scenario participation (settlement delays, failures)

**BDD Acceptance**: Complete trading flow works end-to-end with risk monitoring

**Dependencies**: TSE-0001.7b (Risk Monitor Alert Generation), TSE-0001.8 (Trading Engine), TSE-0001.6 (Custodian)

---

## Implementation Notes

- **Settlement Types**: Start with T+0, design for T+1, T+2 cycles
- **Production API**: REST endpoints that exchanges and risk monitors will use
- **Audit API**: Separate endpoints for chaos injection and internal state
- **Multi-Asset**: Support major crypto and fiat currencies
- **Compliance**: Basic AML/KYC simulation, extensible for real requirements
- **Chaos Ready**: Design for controlled settlement failures

---

**Last Updated**: 2025-09-17
---

## 🔄 Milestone TSE-0001.4: Data Adapters & Orchestrator Integration

**Status**: 📋 **READY TO START** - Following audit-correlator-go Pattern
**Goal**: Integrate custodian-simulator-go with custodian-data-adapter-go and enable Docker deployment
**Phase**: Data Architecture & Deployment
**Pattern Source**: audit-correlator-go (completed 2025-09-30)

### 📚 Reference Implementation
**Proven Pattern**: audit-correlator-go successfully integrated and deployed
- All 7 tasks completed and validated
- 70MB Docker image deployed in orchestrator
- Clean architecture with repository pattern
- Graceful degradation confirmed
- Pattern ready for replication

### 🎯 Integration Tasks (Following Proven 7-Step Process)

#### ✅ **Task 0: Test Infrastructure Foundation** - READY TO START
**Goal**: Establish test automation and build validation
**Reference**: audit-correlator-go tasks completed 2025-09-28

**Files to Create/Modify**:
- `Makefile` - Comprehensive test automation targets
- Build validation - Ensure compilation throughout process
- Test baseline - Document current test status

**Tasks**:
- [ ] Create comprehensive Makefile with unit/integration test targets
- [ ] Establish TDD Red-Green-Refactor pattern testing
- [ ] Document current test status (passing/skipped/failing)
- [ ] Verify build compiles successfully
- [ ] Create test infrastructure foundation

**Acceptance Criteria**:
- [ ] Makefile with test-unit, test-integration, test-all targets
- [ ] Build status: ✅ Compiles successfully
- [ ] Test baseline documented
- [ ] Test automation ready for integration work

**Commands to Validate**:
```bash
make build                  # Verify compilation
make test-unit             # Check unit test baseline
make test-integration      # Check integration test baseline
```

---

#### 📋 **Task 1: Create custodian-data-adapter-go** - AFTER TASK 0
**Goal**: Create dedicated data adapter for custodian operations
**Reference**: audit-data-adapter-go pattern
**New Repository**: `custodian-data-adapter-go`

**Repository Structure to Create**:
```
custodian-data-adapter-go/
├── pkg/
│   ├── interfaces/         # Repository interfaces
│   │   ├── custodian_event.go
│   │   ├── asset_custody.go
│   │   ├── wallet_management.go
│   │   ├── service_discovery.go
│   │   └── cache.go
│   ├── adapters/          # DataAdapter implementation
│   │   ├── custodian_data_adapter.go
│   │   ├── factory.go
│   │   └── interfaces.go
│   └── models/            # Data models
│       ├── custodian_event.go
│       ├── asset.go
│       ├── wallet.go
│       └── enums.go
├── internal/
│   ├── config/            # Configuration management
│   ├── db/               # PostgreSQL connection
│   └── cache/            # Redis connection
├── tests/                # BDD behavior tests
│   ├── init_test.go
│   ├── custodian_event_behavior_test.go
│   ├── asset_custody_behavior_test.go
│   ├── wallet_management_behavior_test.go
│   ├── service_discovery_behavior_test.go
│   ├── cache_behavior_test.go
│   └── integration_behavior_test.go
├── .env.example          # Environment template
├── .gitignore            # Security (exclude .env)
├── Makefile              # Test automation
├── go.mod                # Dependencies
└── TODO.md               # Task tracking
```

**Implementation Steps**:
1. Create new repository: `custodian-data-adapter-go`
2. Copy structure from audit-data-adapter-go
3. Adapt interfaces for custodian domain:
   - CustodianEventRepository (custody events tracking)
   - AssetCustodyRepository (asset holdings management)
   - WalletManagementRepository (wallet operations)
   - ServiceDiscoveryRepository (reuse from audit pattern)
   - CacheRepository (reuse from audit pattern)
4. Implement PostgreSQL adapter for custodian schema
5. Implement Redis adapter for service discovery and caching
6. Create comprehensive BDD behavior tests
7. Add .env.example with orchestrator credentials
8. Document in TODO.md

**Acceptance Criteria**:
- [ ] Repository created with clean architecture structure
- [ ] All interfaces defined for custodian domain
- [ ] PostgreSQL adapter using custodian schema
- [ ] Redis adapter for discovery and caching
- [ ] BDD behavior test framework (30+ test scenarios)
- [ ] .env.example with orchestrator configuration
- [ ] Test success rate >80%
- [ ] Documentation complete (TODO.md, README.md)

**Database Schema** (to create in orchestrator):
```sql
-- In trading_ecosystem database, custodian schema
CREATE TABLE custodian.custodian_events (
    id VARCHAR(255) PRIMARY KEY,
    trace_id VARCHAR(255),
    event_type VARCHAR(100),
    asset_id VARCHAR(255),
    wallet_id VARCHAR(255),
    timestamp TIMESTAMPTZ,
    status VARCHAR(50),
    metadata JSONB
);

CREATE TABLE custodian.asset_custody (
    asset_id VARCHAR(255) PRIMARY KEY,
    wallet_id VARCHAR(255),
    asset_type VARCHAR(100),
    quantity DECIMAL,
    updated_at TIMESTAMPTZ
);

CREATE TABLE custodian.wallet_management (
    wallet_id VARCHAR(255) PRIMARY KEY,
    wallet_type VARCHAR(100),
    status VARCHAR(50),
    metadata JSONB,
    created_at TIMESTAMPTZ
);
```

**Time Estimate**: 4-6 hours (following audit-data-adapter-go pattern)

---

#### 🏗️ **Task 2: Refactor Infrastructure Layer** - AFTER TASK 1
**Goal**: Replace direct database access with custodian-data-adapter-go DataAdapter interfaces
**Reference**: audit-correlator-go Task 2 (completed 2025-09-29)

**Files to Modify**:
- `go.mod` - Add custodian-data-adapter-go dependency
- `internal/infrastructure/service_discovery.go` → Use `DataAdapter.ServiceDiscoveryRepository`
- `internal/infrastructure/configuration_client.go` → Use `DataAdapter.CacheRepository`
- `internal/config/config.go` → Initialize DataAdapter with orchestrator credentials

**Implementation Steps**:
1. Add dependency to go.mod:
   ```go
   require github.com/quantfidential/trading-ecosystem/custodian-data-adapter-go v0.1.0
   replace github.com/quantfidential/trading-ecosystem/custodian-data-adapter-go => ../custodian-data-adapter-go
   ```

2. Update config.go with DataAdapter initialization:
   ```go
   import "github.com/quantfidential/trading-ecosystem/custodian-data-adapter-go/pkg/adapters"
   
   type Config struct {
       // existing fields...
       dataAdapter adapters.DataAdapter
   }
   
   func (c *Config) InitializeDataAdapter(ctx context.Context, logger *logrus.Logger) error {
       adapter, err := adapters.NewCustodianDataAdapterFromEnv(logger)
       // ... handle connection
       c.dataAdapter = adapter
       return nil
   }
   ```

3. Replace Redis service discovery with ServiceDiscoveryRepository
4. Replace Redis configuration caching with CacheRepository
5. Update connection lifecycle management

**Acceptance Criteria**:
- [ ] go.mod includes custodian-data-adapter-go
- [ ] Service discovery uses only DataAdapter.ServiceDiscoveryRepository
- [ ] Configuration caching uses only DataAdapter.CacheRepository
- [ ] DataAdapter initialized with orchestrator credentials
- [ ] Connection lifecycle (Connect/Disconnect) working
- [ ] No direct Redis/PostgreSQL imports
- [ ] Build compiles successfully

---

#### 🔧 **Task 3: Update Service Layer** - AFTER TASK 2
**Goal**: Integrate custodian operations with repository patterns
**Reference**: audit-correlator-go Task 3 (completed 2025-09-29)

**Files to Modify**:
- `internal/services/custodian.go` → Use `DataAdapter.CustodianEventRepository`
- `internal/services/asset.go` → Use `DataAdapter.AssetCustodyRepository`
- `internal/services/wallet.go` → Use `DataAdapter.WalletManagementRepository`
- `internal/handlers/*` → Update to use repository patterns
- `internal/presentation/grpc/*` → Ensure proper model usage

**Implementation Steps**:
1. Update custodian event creation to use CustodianEventRepository.Create()
2. Update asset custody operations to use AssetCustodyRepository
3. Update wallet management to use WalletManagementRepository
4. Ensure all models align with custodian-data-adapter-go patterns
5. Update handlers to delegate to service layer
6. Verify no direct database access

**Acceptance Criteria**:
- [ ] All custodian events created through CustodianEventRepository
- [ ] Asset custody operations through AssetCustodyRepository
- [ ] Wallet management through WalletManagementRepository
- [ ] All models consistent with custodian-data-adapter-go
- [ ] No direct database access in service layer
- [ ] Handlers delegate to service layer
- [ ] gRPC presentation layer clean

---

#### 🧪 **Task 4: Test Integration** - AFTER TASK 3
**Goal**: Enable tests to use shared orchestrator services
**Reference**: audit-correlator-go Task 4 (completed 2025-09-30)

**Files to Create/Modify**:
- `.env.example` → Following custodian-data-adapter-go pattern
- `.env` → Created from .env.example (gitignored)
- `.gitignore` → Add .env patterns
- `Makefile` → Update to load .env for tests
- `go.mod` → Add godotenv dependency

**.env.example Content**:
```bash
# Custodian Simulator Go - Environment Configuration
SERVICE_NAME=custodian-simulator
SERVICE_VERSION=1.0.0
ENVIRONMENT=development

# Network Configuration
HTTP_PORT=8081
GRPC_PORT=9091

# Database Configuration
POSTGRES_URL=postgres://custodian_adapter:custodian-adapter-db-pass@localhost:5432/trading_ecosystem?sslmode=disable
REDIS_URL=redis://custodian-adapter:custodian-pass@localhost:6379/0

# Test Configuration
TEST_POSTGRES_URL=postgres://custodian_adapter:custodian-adapter-db-pass@localhost:5432/trading_ecosystem?sslmode=disable
TEST_REDIS_URL=redis://admin:admin-secure-pass@localhost:6379/0

# Configuration
MAX_CONNECTIONS=25
MAX_IDLE_CONNECTIONS=10
CONNECTION_TIMEOUT=30s
CACHE_TTL=5m
HEALTH_CHECK_INTERVAL=30s
LOG_LEVEL=info
```

**Makefile Enhancement**:
```makefile
# Load environment variables
ifneq (,$(wildcard .env))
	include .env
	export
endif

check-env:
	@if [ ! -f .env ]; then \
		echo "Warning: .env not found. Copy .env.example to .env"; \
		exit 1; \
	fi

test-unit:
	@if [ -f .env ]; then set -a && . ./.env && set +a; fi && go test -tags=unit ./internal/... -v

test-integration: check-env
	@set -a && . ./.env && set +a && go test -tags=integration ./internal/... -v
```

**Acceptance Criteria**:
- [ ] .env.example created with orchestrator configuration
- [ ] Tests use shared orchestrator PostgreSQL/Redis
- [ ] Makefile loads .env automatically
- [ ] .env properly gitignored
- [ ] godotenv dependency added
- [ ] Build compiles with environment integration
- [ ] Test success rate improves

---

#### ⚙️ **Task 5: Configuration Integration** - AFTER TASK 4
**Goal**: Complete environment alignment and lifecycle management
**Reference**: audit-correlator-go Task 5 (completed 2025-09-30)

**Tasks** (Most already done in Task 4):
- [ ] Verify .env.example alignment with custodian-data-adapter-go
- [ ] Verify proper DataAdapter lifecycle management in main.go
- [ ] Test connection management (Connect/Disconnect)
- [ ] Verify graceful fallback when infrastructure unavailable

**Acceptance Criteria**:
- [ ] Environment configuration complete
- [ ] DataAdapter lifecycle working
- [ ] Graceful degradation confirmed
- [ ] Integration pattern documented

---

#### 🐳 **Task 6: Docker Deployment Integration** - AFTER TASK 5
**Goal**: Package and deploy custodian-simulator-go in orchestrator
**Reference**: audit-correlator-go Task 6 (completed 2025-09-30)

**Files to Modify**:
- `Dockerfile` → Update for multi-context build
- `orchestrator-docker/docker-compose.yml` → Add custodian-simulator service
- `orchestrator-docker/redis/users.acl` → Add custodian-adapter user

**Dockerfile Update** (multi-context build):
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /build

# Copy custodian-data-adapter-go dependency
COPY custodian-data-adapter-go/ ./custodian-data-adapter-go/

# Copy custodian-simulator-go files
COPY custodian-simulator-go/go.mod custodian-simulator-go/go.sum ./custodian-simulator-go/
WORKDIR /build/custodian-simulator-go
RUN go mod download

# Copy source and build
COPY custodian-simulator-go/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o custodian-simulator ./cmd/server

# Runtime stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates wget
RUN addgroup -g 1001 -S appgroup && adduser -u 1001 -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /build/custodian-simulator-go/custodian-simulator /app/custodian-simulator
RUN chown -R appuser:appgroup /app
USER appuser
EXPOSE 8081 9091
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8081/api/v1/health || exit 1
CMD ["./custodian-simulator"]
```

**docker-compose.yml Addition**:
```yaml
custodian-simulator:
  build:
    context: ..
    dockerfile: custodian-simulator-go/Dockerfile
  image: custodian-simulator:latest
  container_name: trading-ecosystem-custodian-simulator
  restart: unless-stopped
  ports:
    - "127.0.0.1:8081:8081"  # HTTP
    - "127.0.0.1:9091:9091"  # gRPC
  networks:
    trading-ecosystem:
      ipv4_address: 172.20.0.81
  environment:
    - SERVICE_NAME=custodian-simulator
    - SERVICE_VERSION=1.0.0
    - ENVIRONMENT=docker
    - HTTP_PORT=8081
    - GRPC_PORT=9091
    - POSTGRES_URL=postgres://custodian_adapter:custodian-adapter-db-pass@172.20.0.20:5432/trading_ecosystem?sslmode=disable
    - REDIS_URL=redis://custodian-adapter:custodian-pass@172.20.0.10:6379/0
    - MAX_CONNECTIONS=25
    - HEALTH_CHECK_INTERVAL=15s
    - LOG_LEVEL=info
  depends_on:
    redis:
      condition: service_healthy
    postgres:
      condition: service_healthy
  healthcheck:
    test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8081/api/v1/health"]
    interval: 15s
    timeout: 5s
    retries: 3
    start_period: 45s
```

**Redis ACL Update** (users.acl):
```
user custodian-adapter on >custodian-pass ~custodian:* +@read +@write +@keyspace +ping -@dangerous
```

**Build and Deploy Commands**:
```bash
# Build from parent directory
cd /path/to/trading-ecosystem
docker build -f custodian-simulator-go/Dockerfile -t custodian-simulator:latest .

# Deploy with docker-compose
cd orchestrator-docker
docker-compose up -d custodian-simulator

# Verify deployment
docker ps --filter "name=custodian-simulator"
curl http://localhost:8081/api/v1/health
```

**Acceptance Criteria**:
- [ ] Dockerfile builds from parent context
- [ ] Docker image <100MB (target: ~70MB like audit-correlator)
- [ ] Container starts and runs in orchestrator
- [ ] HTTP and gRPC servers accessible
- [ ] Health checks passing
- [ ] PostgreSQL connection working
- [ ] Graceful fallback operational
- [ ] Service registered in trading-ecosystem network

---

### 🎯 Success Metrics (Following audit-correlator-go)

| Metric | Target | Status |
|--------|--------|--------|
| Tasks Complete | 7/7 (100%) | [ ] |
| Build Status | Pass | [ ] |
| Test Coverage | 10+ tests passing | [ ] |
| Docker Image | <100MB | [ ] |
| Deployment | Healthy container | [ ] |
| Orchestrator | PostgreSQL + Redis connected | [ ] |
| Graceful Degradation | Working | [ ] |
| Pattern Validated | Yes | [ ] |

---

### 📋 Testing Commands Throughout Integration

```bash
# After each task, validate progress
make build                  # Ensure compilation
make test-unit             # Check unit test status
make test-integration      # Check integration tests
docker build -f Dockerfile -t custodian-simulator:latest .  # Task 6
```

---

### 🔗 Dependencies and References

**Prerequisites**:
- TSE-0001.3b (Go Services gRPC Integration) ✅ COMPLETE
- orchestrator-docker running (PostgreSQL, Redis)
- audit-correlator-go pattern (reference implementation)

**Reference Pattern**: audit-correlator-go
- All 7 tasks completed and validated
- Clean architecture with repository pattern
- Docker deployment successful
- Pattern ready for replication

**Related Work**:
- Create custodian-data-adapter-go (parallel to audit-data-adapter-go)
- Update orchestrator-docker with custodian configuration
- Follow proven 7-step integration process

---

**Epic**: TSE-0001 Foundation Services & Infrastructure
**Milestone**: TSE-0001.4 Data Adapters & Orchestrator Integration (2/4 when complete)
**Status**: 📋 READY TO START - Pattern Established
**Time Estimate**: 6-8 hours (following proven pattern)

**Last Updated**: 2025-09-30
