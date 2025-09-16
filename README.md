# Custodian Simulator

A high-fidelity custodian and prime brokerage simulator built in Go that replicates institutional custody operations, settlement workflows, and multi-day clearing cycles for comprehensive trading system testing.

## 🎯 Overview

The Custodian Simulator models institutional custody providers (BitGo, Coinbase Custody, Prime Trust, etc.) with realistic settlement delays, multi-signature security, and regulatory compliance workflows. It handles the critical "money movement" between trading venues and master custody accounts that real trading firms depend on.

### Key Features
- **Multi-Day Settlement**: Realistic T+0 to T+3 settlement cycles with business day logic
- **Institutional Workflows**: Wire transfers, crypto settlements, and cross-venue reconciliation
- **Multi-Signature Security**: Simulated approval workflows and security controls
- **Master Account Management**: Segregated client assets with full audit trails
- **Regulatory Compliance**: AML/KYC simulation and reporting capabilities
- **Chaos Engineering**: Controllable settlement failures and operational disruptions
- **Real-Time Reconciliation**: Continuous balance verification across all venues

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────┐
│                Custodian Simulator                      │
├─────────────────────────────────────────────────────────┤
│  gRPC Services          │  REST APIs                    │
│  ├─Settlement Service   │  ├─Account Management         │
│  ├─Transfer Service     │  ├─Reconciliation             │
│  ├─Custody Service      │  └─Chaos Engineering          │
│  └─Reporting Service    │                               │
├─────────────────────────────────────────────────────────┤
│  Core Engine                                            │
│  ├─Settlement Engine (Multi-day cycles)                │
│  ├─Transfer Manager (Wire/Crypto movements)            │
│  ├─Approval Workflow (Multi-sig simulation)           │
│  ├─Reconciliation Engine (Cross-venue verification)    │
│  └─Chaos Controller (Settlement failures)              │
├─────────────────────────────────────────────────────────┤
│  Data Layer                                             │
│  ├─Master Accounts (Segregated client assets)          │
│  ├─Settlement Queue (Pending transfers)                │
│  ├─Transaction History (Full audit trail)              │
│  └─Approval Records (Multi-sig workflows)              │
└─────────────────────────────────────────────────────────┘
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker and Docker Compose
- Protocol Buffers compiler

### Development Setup
```bash
# Clone the repository
git clone <repo-url>
cd custodian-simulator

# Install dependencies
go mod download

# Generate protobuf files
make generate-proto

# Run tests
make test

# Start development server
make run-dev
```

### Docker Deployment
```bash
# Build container
docker build -t custodian-simulator .

# Run with docker-compose (recommended)
docker-compose up custodian-simulator

# Verify health
curl http://localhost:8081/health
```

## 📡 API Reference

### gRPC Services

#### Settlement Service
```protobuf
service SettlementService {
  rpc InitiateSettlement(SettlementRequest) returns (SettlementResponse);
  rpc GetSettlementStatus(SettlementStatusRequest) returns (SettlementStatus);
  rpc ListPendingSettlements(ListSettlementsRequest) returns (ListSettlementsResponse);
  rpc ProcessSettlementQueue(ProcessQueueRequest) returns (ProcessQueueResponse);
}
```

#### Transfer Service
```protobuf
service TransferService {
  rpc InitiateTransfer(TransferRequest) returns (TransferResponse);
  rpc ApproveTransfer(ApprovalRequest) returns (ApprovalResponse);
  rpc GetTransferHistory(TransferHistoryRequest) returns (TransferHistoryResponse);
}
```

#### Custody Service
```protobuf
service CustodyService {
  rpc GetMasterBalance(BalanceRequest) returns (MasterBalance);
  rpc GetSegregatedBalance(SegregatedBalanceRequest) returns (SegregatedBalance);
  rpc CreateCustodyAccount(CreateAccountRequest) returns (CustodyAccount);
  rpc GetAuditTrail(AuditTrailRequest) returns (AuditTrailResponse);
}
```

### REST Endpoints

#### Production APIs (Risk Monitor Accessible)
```
GET    /api/v1/accounts/{client_id}/master-balance
GET    /api/v1/accounts/{client_id}/segregated-balances
GET    /api/v1/settlements/{settlement_id}/status
GET    /api/v1/transfers/{transfer_id}/status
POST   /api/v1/settlements/initiate
POST   /api/v1/transfers/initiate
```

#### Chaos Engineering APIs (Audit Only)
```
POST   /chaos/delay-settlements
POST   /chaos/reject-transfers
POST   /chaos/simulate-bank-holiday
POST   /chaos/partial-settlement-failure
POST   /chaos/multi-sig-delays
GET    /chaos/status
DELETE /chaos/clear-all
```

#### State Inspection APIs (Development/Audit)
```
GET    /debug/settlement-queue
GET    /debug/master-accounts
GET    /debug/pending-approvals
GET    /debug/reconciliation-status
GET    /metrics (Prometheus format)
```

## 💼 Settlement Engine

### Settlement Cycles
```
T+0: Same-Day Settlement (Crypto internal transfers)
├── Initiation: Immediate
├── Approval: Multi-sig required (2-5 minutes)
└── Completion: 5-15 minutes

T+1: Next Business Day (Standard crypto settlements)
├── Initiation: Same day before 4 PM EST
├── Processing: Next business day 9 AM EST
└── Completion: Next business day 5 PM EST

T+2: Two Business Days (Wire transfers, large amounts)
├── Initiation: Same day before 2 PM EST
├── Processing: T+1 verification, T+2 execution
└── Completion: T+2 by 3 PM EST

T+3: Three Business Days (International wires)
├── Initiation: Same day before 12 PM EST
├── Processing: Multi-day compliance checks
└── Completion: T+3 by 2 PM EST
```

### Business Day Logic
- **Trading Days**: Monday-Friday, excluding holidays
- **Cut-off Times**: Different cut-offs for different settlement types
- **Holiday Handling**: Automatically adjusts settlement dates
- **Time Zone Support**: Multiple regional business day calendars

### Settlement Types
```go
type SettlementType string

const (
    CryptoInternal   SettlementType = "crypto_internal"    // T+0, 5-15 min
    CryptoExternal   SettlementType = "crypto_external"    // T+1, same-day initiation
    WireDomestic     SettlementType = "wire_domestic"      // T+2, before 2 PM
    WireInternational SettlementType = "wire_international" // T+3, before 12 PM
    ACHTransfer      SettlementType = "ach_transfer"       // T+2, before 3 PM
)
```

## 🏦 Account Management

### Master Account Structure
```
Client: trading-firm-001
├── Master USD Account: $5,000,000
│   ├── Available: $3,500,000
│   ├── Pending Settlements: $1,200,000
│   └── Reserved: $300,000
├── Master BTC Account: 150.5 BTC
│   ├── Available: 120.3 BTC
│   ├── Pending Settlements: 25.2 BTC
│   └── Cold Storage: 5.0 BTC
└── Master ETH Account: 2,500 ETH
    ├── Available: 2,100 ETH
    ├── Pending Settlements: 350 ETH
    └── Staking: 50 ETH
```

### Segregated Client Assets
- **Regulatory Compliance**: Client assets held separately from firm assets
- **Audit Trail**: Complete transaction history with timestamps
- **Multi-Signature**: All movements require multiple approvals
- **Insurance**: Simulated FDIC/SIPC insurance coverage

## 🔐 Multi-Signature Workflows

### Approval Hierarchy
```
Transfer Amount          Required Approvals    Estimated Time
├── < $10,000           │ 1 approval         │ 2-5 minutes
├── $10,000 - $100,000  │ 2 approvals        │ 5-15 minutes  
├── $100,000 - $1M      │ 3 approvals        │ 15-60 minutes
└── > $1M               │ 4+ approvals       │ 1-4 hours
```

### Approval Simulation
- **Realistic Timing**: Approvals take time based on amount and time of day
- **Business Hours**: Faster approvals during business hours
- **Weekend Delays**: Reduced approval capacity on weekends
- **Emergency Procedures**: Fast-track approvals for critical situations

## 🔄 Reconciliation Engine

### Continuous Verification
```
Every 5 minutes:
├── Exchange Balance Verification
│   ├── Query all exchange APIs
│   ├── Compare with internal records
│   └── Flag discrepancies > $1,000
├── Settlement Status Updates
│   ├── Check pending settlement status
│   ├── Update completion timestamps
│   └── Process settlement completions
└── Master Account Reconciliation
    ├── Sum all venue balances
    ├── Compare with master account
    └── Generate reconciliation reports
```

### Discrepancy Handling
- **Automatic Resolution**: Small discrepancies (<$100) auto-resolved
- **Manual Review**: Medium discrepancies flagged for investigation
- **Immediate Alert**: Large discrepancies trigger immediate notifications
- **Audit Trail**: All reconciliation activities logged

## 🎭 Chaos Engineering

### Settlement Failure Injection

#### Delayed Settlements
```bash
# Delay all T+1 settlements by 24 hours
curl -X POST localhost:8081/chaos/delay-settlements \
  -d '{"settlement_types": ["crypto_external"], "delay_hours": 24, "duration_s": 3600}'
```

#### Partial Settlement Failures
```bash
# 30% of wire transfers fail with "insufficient_funds"
curl -X POST localhost:8081/chaos/partial-settlement-failure \
  -d '{"settlement_type": "wire_domestic", "failure_rate": 0.3, "error": "insufficient_funds"}'
```

#### Multi-Sig Delays
```bash
# Slow down approval process by 10x
curl -X POST localhost:8081/chaos/multi-sig-delays \
  -d '{"delay_multiplier": 10, "min_amount": 50000, "duration_s": 1800}'
```

#### Banking Holiday Simulation
```bash
# Simulate unexpected bank holiday (no settlements processed)
curl -X POST localhost:8081/chaos/simulate-bank-holiday \
  -d '{"duration_hours": 24, "affected_types": ["wire_domestic", "ach_transfer"]}'
```

#### Custodian Offline
```bash
# Simulate complete custodian downtime
curl -X POST localhost:8081/chaos/custodian-offline \
  -d '{"duration_s": 7200, "error_message": "system_maintenance"}'
```

## 📊 Monitoring & Observability

### Prometheus Metrics
```
# Settlement metrics
custodian_settlements_total{type="crypto_internal|wire_domestic", status="pending|completed|failed"}
custodian_settlement_latency_hours{type, status}
custodian_settlement_queue_depth{type}

# Account metrics
custodian_master_balance{client_id, asset}
custodian_available_balance{client_id, asset}
custodian_pending_settlements{client_id, asset}

# Reconciliation metrics
custodian_reconciliation_discrepancies{venue, asset}
custodian_reconciliation_last_success_timestamp{venue}

# Approval metrics
custodian_pending_approvals{amount_tier}
custodian_approval_latency_minutes{amount_tier}

# System health
custodian_uptime_seconds
custodian_chaos_active{type}
custodian_api_response_time{endpoint}
```

### OpenTelemetry Tracing
- **Settlement Flows**: Complete trace from initiation to completion
- **Approval Workflows**: Track multi-signature approval timing
- **Reconciliation**: Cross-venue balance verification traces
- **Error Propagation**: Failed settlements and their causes

### Structured Logging
```json
{
  "timestamp": "2025-09-16T14:23:45Z",
  "level": "info",
  "service": "custodian-simulator",
  "correlation_id": "settlement-xyz789",
  "event": "settlement_initiated",
  "client_id": "trading-firm-001",
  "settlement_type": "wire_domestic",
  "amount": "500000.00",
  "currency": "USD",
  "expected_completion": "2025-09-18T15:00:00Z",
  "approval_tier": "3_approvals_required"
}
```

## 🧪 Testing

### Unit Tests
```bash
# Run all unit tests
make test

# Run with coverage
make test-coverage

# Test specific settlement types
go test ./internal/settlement -v -run TestT1Settlement
```

### Integration Tests
```bash
# Test with real dependencies
make test-integration

# Test settlement workflows end-to-end
make test-settlement-flows

# Test chaos injection
make test-chaos
```

### Scenario Testing
```bash
# Test business day calculations
go test ./internal/settlement -run TestBusinessDayLogic

# Test multi-signature workflows
go test ./internal/approval -run TestMultiSigApproval

# Test reconciliation accuracy
go test ./internal/reconciliation -run TestReconciliationAccuracy
```

## ⚙️ Configuration

### Environment Variables
```bash
# Core settings
CUSTODIAN_PORT=8081
CUSTODIAN_GRPC_PORT=50052
CUSTODIAN_LOG_LEVEL=info

# Dependencies
REDIS_URL=redis://localhost:6379
POSTGRES_URL=postgres://user:pass@localhost/custodian

# Settlement parameters
DEFAULT_SETTLEMENT_DELAY_MINUTES=15
MAX_SETTLEMENT_AMOUNT=10000000
ENABLE_WEEKEND_PROCESSING=false

# Security settings
MULTI_SIG_ENABLED=true
MIN_APPROVALS=2
MAX_APPROVAL_TIME_HOURS=24

# Reconciliation
RECONCILIATION_INTERVAL_MINUTES=5
RECONCILIATION_TOLERANCE_USD=1000
```

### Configuration File (config.yaml)
```yaml
custodian:
  name: "custody-provider-1"
  institution_id: "CUST001"
  regulatory_jurisdiction: "USA"

settlement_rules:
  crypto_internal:
    settlement_time: "T+0"
    max_processing_minutes: 15
    cut_off_time: "23:59"
  
  wire_domestic:
    settlement_time: "T+2" 
    max_processing_hours: 48
    cut_off_time: "14:00"
    business_days_only: true

approval_tiers:
  - max_amount: 10000
    required_approvals: 1
    max_approval_time_minutes: 15
  - max_amount: 1000000
    required_approvals: 3
    max_approval_time_hours: 4

reconciliation:
  interval_minutes: 5
  tolerance_usd: 1000
  auto_resolve_threshold: 100
```

## 🐳 Docker Configuration

### Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o custodian-simulator cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /app/custodian-simulator /usr/local/bin/
EXPOSE 8081 50052
CMD ["custodian-simulator"]
```

### Health Checks
```yaml
healthcheck:
  test: ["CMD", "grpc_health_probe", "-addr=:50052"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 60s
```

## 🔒 Security & Compliance

### Security Features
- **Multi-Signature**: All transfers require multiple approvals
- **Audit Logging**: Immutable audit trail for all operations
- **Access Controls**: Role-based access to different functions
- **Encryption**: All sensitive data encrypted at rest and in transit

### Regulatory Simulation
- **AML Monitoring**: Transaction pattern analysis and flagging
- **KYC Verification**: Client onboarding and verification workflows
- **Reporting**: Automated regulatory reporting generation
- **Compliance Alerts**: Suspicious activity detection and reporting

### Risk Controls
- **Velocity Limits**: Maximum transfer amounts per time period
- **Concentration Limits**: Maximum exposure per client or asset
- **Counterparty Limits**: Maximum exposure per exchange or venue
- **Emergency Procedures**: Ability to halt all transfers immediately

## 🚀 Performance

### Benchmarks
- **Settlement Throughput**: >1,000 settlements/hour (normal operations)
- **API Response Time**: <100ms p99 for balance queries
- **Reconciliation Speed**: Complete venue reconciliation in <30 seconds
- **Memory Usage**: <200MB baseline, <1GB under peak load

### Scaling Considerations
- **Database Sharding**: Client-based sharding for large deployments
- **Async Processing**: Settlement queue processing with workers
- **Caching Strategy**: Hot balance data cached in Redis
- **Archive Strategy**: Historical data migration to cold storage

## 🤝 Contributing

### Development Workflow
1. Create feature branch from `main`
2. Implement changes with comprehensive tests
3. Run full test suite: `make test-all`
4. Test with chaos injection scenarios
5. Update documentation and configuration examples
6. Submit pull request with detailed description

### Code Standards
- **Go Best Practices**: Follow effective Go guidelines
- **Financial Precision**: Use decimal types for all monetary calculations
- **Error Handling**: Comprehensive error types and handling
- **Documentation**: All public APIs documented with examples

## 📚 References

- **Custody Operations**: [Link to institutional custody documentation]
- **Settlement Standards**: [Link to financial settlement specifications]
- **Protobuf Schemas**: [Link to custodian API definitions]
- **Regulatory Framework**: [Link to compliance requirements]

---

**Status**: 🚧 Development Phase  
**Maintainer**: [Your team]  
**Last Updated**: September 2025
