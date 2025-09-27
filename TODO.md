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
**Status**: Ready to Begin (Following audit-correlator-go pattern)
**Priority**: High

**Tasks** (Following proven TDD Red-Green-Refactor cycle):
- [x] **Phase 1: TDD Red** - Create failing tests for all gRPC integration behaviors
- [ ] **Phase 2: Infrastructure** - Add Redis dependencies and update .gitignore for Go projects
- [ ] **Phase 3: gRPC Server** - Implement enhanced gRPC server with health service, metrics, and graceful shutdown
- [ ] **Phase 4: Configuration** - Implement configuration service client with HTTP caching, TTL, and type conversion
- [ ] **Phase 5: Discovery** - Implement service discovery with Redis-based registry, heartbeat, and cleanup
- [ ] **Phase 6: Communication** - Create inter-service gRPC client manager with connection pooling and circuit breaker
- [ ] **Phase 7: Integration** - Implement comprehensive inter-service communication testing with smart skipping
- [ ] **Phase 8: Validation** - Verify BDD acceptance and complete milestone documentation

**Implementation Pattern** (Replicating audit-correlator-go success):
- **Infrastructure Layer**: Configuration client, service discovery, gRPC clients
- **Presentation Layer**: Enhanced gRPC server with health service
- **Testing Strategy**: Unit tests with smart dependency skipping, integration tests for end-to-end scenarios
- **Error Handling**: Graceful degradation, circuit breaker patterns, comprehensive logging

**BDD Acceptance**: Go services can discover and communicate with each other via gRPC

**Dependencies**: TSE-0001.1a (Go Services Bootstrapping), TSE-0001.3a (Core Infrastructure)

**Reference Implementation**: audit-correlator-go (✅ COMPLETED) - Use as pattern for architecture and testing

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