<div align="center">

# GoPayX 💳

### Enterprise Digital Wallet & Payment Platform

**Production-ready microservices backend built with Go**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![CI/CD](https://github.com/afifn11/gopay-x/actions/workflows/ci.yml/badge.svg)](https://github.com/afifn11/gopay-x/actions)
[![Architecture](https://img.shields.io/badge/Architecture-Microservices-blueviolet)]()
[![Services](https://img.shields.io/badge/Services-8-orange)]()

</div>

---

## 📖 Overview

GoPayX is a **production-grade digital wallet and payment gateway platform** that demonstrates industry-standard backend engineering practices. Built with Go microservices, it handles user authentication, wallet management, payment processing, fraud detection, and comprehensive audit trails.

> **Portfolio Project** — Designed to showcase enterprise-level backend skills including distributed systems, event-driven architecture, and financial system design patterns.

---

## ✨ Key Features

| Feature | Description |
|---|---|
| 🔐 **JWT Authentication** | Access + refresh token rotation with Redis blacklisting |
| 💰 **Digital Wallet** | Balance management with distributed locking (Redis SetNX) |
| 💳 **Payment Gateway** | Mock Midtrans/Xendit integration with webhook handling |
| 🔄 **Idempotency** | Reference ID-based duplicate prevention on all financial ops |
| 📨 **Event-Driven** | Apache Kafka for async communication between services |
| 🚨 **Fraud Detection** | Rule engine with velocity checks and risk scoring (0-100) |
| 📋 **Audit Trail** | Immutable append-only event log across all services |
| 🚪 **API Gateway** | Rate limiting, JWT validation, request logging, CORS |
| 📊 **Monitoring** | Prometheus metrics + Grafana dashboards |
| 🐳 **Containerized** | Docker + Docker Compose for all infrastructure |
| ⚙️ **CI/CD** | GitHub Actions pipeline with build + security scan |

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────┐
│                    Client Layer                      │
│           Mobile App · Web App · 3rd Party          │
└─────────────────────┬───────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────┐
│                  API Gateway :8080                   │
│        Rate Limiting · Auth Check · Routing         │
└──┬──────┬──────┬──────┬──────┬──────┬──────┬───────┘
   │      │      │      │      │      │      │
┌──▼──┐ ┌─▼──┐ ┌─▼───┐ ┌▼────┐ ┌▼───┐ ┌▼────┐ ┌▼────┐ ┌▼────┐
│Auth │ │User│ │Wall │ │Pay  │ │Txn │ │Notif│ │Fraud│ │Audit│
│8081 │ │8082│ │8083 │ │8084 │ │8085│ │8086 │ │8087 │ │8088 │
└──┬──┘ └─┬──┘ └──┬──┘ └──┬──┘ └─┬──┘ └──┬──┘ └──┬──┘ └──┬──┘
   │      │       │        │      │       │       │       │
   └──────┴───────┴────────┴──────┘       │       │       │
                    │                     │       │       │
          ┌─────────▼─────────────────────▼───────▼───────▼──┐
          │              Apache Kafka Event Bus               │
          │  payment.created · wallet.updated · fraud.flagged │
          └───────────────────────────────────────────────────┘
                    │
          ┌─────────▼──────────────────────────────────────┐
          │                  Data Layer                     │
          │  PostgreSQL · Redis · (Elasticsearch planned)  │
          └────────────────────────────────────────────────┘
```

---

## 🛠️ Tech Stack

| Layer | Technology | Purpose |
|---|---|---|
| **Language** | Go 1.22+ | All services |
| **HTTP Framework** | Gin | REST API |
| **Database** | PostgreSQL 16 | Primary data store |
| **Cache / Lock** | Redis 7 | Token blacklist, distributed lock |
| **Message Broker** | Apache Kafka | Async event streaming |
| **Auth** | JWT (golang-jwt/jwt) | Stateless authentication |
| **ORM** | GORM | Database operations |
| **Config** | Viper | Environment management |
| **Container** | Docker + Compose | Infrastructure |
| **Monitoring** | Prometheus + Grafana | Metrics & dashboards |
| **CI/CD** | GitHub Actions | Automated pipeline |
| **API Docs** | Postman Collection | Interactive documentation |

---

## 📦 Services

| Service | Port | Responsibility | Database |
|---|---|---|---|
| **api-gateway** | 8080 | Routing, rate limiting, auth check | — |
| **auth-service** | 8081 | JWT auth, OAuth2, KYC flow | gopay_auth |
| **user-service** | 8082 | User profile management | gopay_user |
| **wallet-service** | 8083 | Balance, top-up, locking | gopay_wallet |
| **payment-service** | 8084 | Transfer, payment gateway | gopay_payment |
| **transaction-service** | 8085 | Immutable ledger, history | gopay_transaction |
| **notification-service** | 8086 | Email, push, SMS (event-driven) | — |
| **fraud-detection-service** | 8087 | Rule engine, risk scoring | gopay_fraud |
| **audit-service** | 8088 | Append-only event log | gopay_audit |

---

## 🚀 Getting Started

### Prerequisites

| Tool | Version | Download |
|---|---|---|
| Go | 1.22+ | [go.dev/dl](https://go.dev/dl/) |
| Docker Desktop | Latest | [docker.com](https://www.docker.com/products/docker-desktop/) |
| Git | 2.x+ | [git-scm.com](https://git-scm.com/) |

### 1. Clone Repository

```bash
git clone https://github.com/afifn11/gopay-x.git
cd gopay-x
```

### 2. Start Infrastructure

```bash
docker compose -f infra/docker-compose.yml up -d
```

This starts: PostgreSQL, Redis, Kafka, Zookeeper, Kafka UI, Prometheus, Grafana.

### 3. Create Databases

```bash
docker exec -it gopay_postgres psql -U gopay -d gopay_auth -c "CREATE DATABASE gopay_user;"
docker exec -it gopay_postgres psql -U gopay -d gopay_auth -c "CREATE DATABASE gopay_wallet;"
docker exec -it gopay_postgres psql -U gopay -d gopay_auth -c "CREATE DATABASE gopay_payment;"
docker exec -it gopay_postgres psql -U gopay -d gopay_auth -c "CREATE DATABASE gopay_transaction;"
docker exec -it gopay_postgres psql -U gopay -d gopay_auth -c "CREATE DATABASE gopay_fraud;"
docker exec -it gopay_postgres psql -U gopay -d gopay_auth -c "CREATE DATABASE gopay_audit;"
```

### 4. Configure Environment

Each service has its own `.env` file. Key variables to set consistently across all services:

```env
JWT_ACCESS_SECRET=<same-secret-across-all-services>
DB_HOST=127.0.0.1
DB_PORT=5433
DB_PASSWORD=gopay123
REDIS_PASSWORD=redis_secret
KAFKA_BROKERS=localhost:9092
```

### 5. Run All Services

Open 9 terminals and run each service:

```bash
# Terminal 1 — Auth Service
cd services/auth-service && go run cmd/main.go

# Terminal 2 — User Service
cd services/user-service && go run cmd/main.go

# Terminal 3 — Wallet Service
cd services/wallet-service && go run cmd/main.go

# Terminal 4 — Payment Service
cd services/payment-service && go run cmd/main.go

# Terminal 5 — Transaction Service
cd services/transaction-service && go run cmd/main.go

# Terminal 6 — Notification Service
cd services/notification-service && go run cmd/main.go

# Terminal 7 — Fraud Detection Service
cd services/fraud-detection-service && go run cmd/main.go

# Terminal 8 — Audit Service
cd services/audit-service && go run cmd/main.go

# Terminal 9 — API Gateway
cd api-gateway && go run cmd/main.go
```

### 6. Verify Everything is Running

```bash
curl http://localhost:8080/health
# {"service":"api-gateway","status":"ok","version":"1.0.0"}

curl http://localhost:8080/health/services
# Lists all 8 registered services
```

---

## 📡 API Documentation

Import Postman files from `docs/postman/`:
- `gopay-x.collection.json` — All 35+ endpoints with test scripts
- `gopay-x.environment.json` — Environment variables

### Quick API Reference

```
POST   /api/v1/auth/register              Register new user
POST   /api/v1/auth/login                 Login → returns JWT tokens
POST   /api/v1/auth/refresh               Refresh access token
POST   /api/v1/auth/logout                Logout + blacklist token
GET    /api/v1/auth/validate              Validate token

GET    /api/v1/users/me                   Get my profile
PUT    /api/v1/users/me                   Update profile
POST   /api/v1/users/me/kyc              Submit KYC document

POST   /api/v1/wallets                    Create wallet
GET    /api/v1/wallets                    Get balance
POST   /api/v1/wallets/topup              Top up (idempotent)
GET    /api/v1/wallets/transactions       Transaction history

POST   /api/v1/payments/transfer          Transfer to user
POST   /api/v1/payments/topup             Top up via gateway
GET    /api/v1/payments/:id               Get payment detail
GET    /api/v1/payments                   Payment history
POST   /api/v1/payments/callback          Webhook callback

GET    /api/v1/transactions               Transaction ledger
GET    /api/v1/transactions/summary       Income/expense summary
GET    /api/v1/transactions/:id           Get transaction detail

GET    /api/v1/fraud/users/:id/checks         Fraud checks (admin)
GET    /api/v1/fraud/users/:id/risk-profile   Risk profile (admin)

GET    /api/v1/audit/logs                     Audit logs (admin)
GET    /api/v1/audit/actors/:id               Logs by actor (admin)
GET    /api/v1/audit/resources/:id            Logs by resource (admin)
```

---

## 🔒 Security Features

- **JWT with refresh token rotation** — access tokens expire in 15 minutes
- **Redis token blacklisting** — invalidated tokens blocked immediately
- **Distributed locking** — Redis SetNX prevents race conditions on wallet ops
- **Idempotency keys** — all financial operations are safe to retry
- **Rate limiting** — 10 req/s per IP with burst of 20
- **bcrypt password hashing** — cost factor 10
- **Role-based access control** — admin-only endpoints protected at gateway level

---

## 🏦 Financial Design Patterns

### Idempotency
Every financial operation requires a unique `reference_id`. Duplicate requests return the original result safely without re-processing.

```json
POST /api/v1/wallets/topup
{
  "amount": 100000,
  "reference_id": "topup-2026041201",
  "description": "Top up"
}
```

### Distributed Locking
Wallet operations acquire a Redis lock before modifying balance, preventing race conditions in concurrent requests.

```
Request A → AcquireLock("wallet:user-123") → ✅ acquired → Update balance → ReleaseLock
Request B → AcquireLock("wallet:user-123") → ❌ failed  → return HTTP 429
```

### Event-Driven Flow
```
User tops up wallet:
1. POST /api/v1/wallets/topup
2. wallet-service validates idempotency key
3. wallet-service acquires Redis distributed lock
4. wallet-service updates balance in PostgreSQL
5. wallet-service publishes "topup.success" to Kafka
                     ├── transaction-service  → records to immutable ledger
                     ├── notification-service → sends user notification
                     ├── fraud-detection      → runs rule engine
                     └── audit-service        → records audit log
```

### Fraud Detection Rules

| Rule | Threshold | Risk Score |
|---|---|---|
| Large transaction | > Rp 10.000.000 | +30 |
| Very large transaction | > Rp 50.000.000 | +40 |
| Velocity check | > 5 transactions in 10 min | +25 |
| High volume | > Rp 20.000.000 in 10 min | +20 |
| Round number | Multiple of 1M, amount ≥ 5M | +5 |

| Score Range | Risk Level | Action |
|---|---|---|
| 0 — 24 | 🟢 Low | CLEARED |
| 25 — 49 | 🟡 Medium | FLAGGED |
| 50 — 79 | 🟠 High | FLAGGED |
| 80 — 100 | 🔴 Critical | BLOCKED |

---

## 📊 Monitoring & Infrastructure

| Tool | URL | Credentials |
|---|---|---|
| **Kafka UI** | http://localhost:8090 | — |
| **Prometheus** | http://localhost:9090 | — |
| **Grafana** | http://localhost:3000 | admin / gopay_grafana |

---

## 🗂️ Project Structure

```
gopay-x/
├── api-gateway/                    # API Gateway — single entry point
│   ├── cmd/main.go
│   ├── config/
│   └── internal/
│       ├── handler/router.go
│       ├── middleware/             # rate limiter, auth, logger, CORS
│       └── proxy/reverse_proxy.go
├── services/
│   ├── auth-service/               # JWT authentication
│   ├── user-service/               # User profiles & KYC
│   ├── wallet-service/             # Digital wallet management
│   ├── payment-service/            # Payments & transfers
│   ├── transaction-service/        # Immutable transaction ledger
│   ├── notification-service/       # Event-driven notifications
│   ├── fraud-detection-service/    # Rule-based fraud detection
│   └── audit-service/              # Full audit trail
├── shared/                         # Shared utilities (proto, middleware, pkg)
├── infra/
│   ├── docker-compose.yml          # All infrastructure containers
│   └── monitoring/
│       └── prometheus.yml
├── docs/
│   ├── postman/                    # API collection & environment
│   └── adr/                        # Architecture Decision Records
└── .github/
    └── workflows/ci.yml            # GitHub Actions CI/CD
```

Each service follows **Clean Architecture**:
```
service/
├── cmd/main.go           # Entry point & dependency injection
├── config/config.go      # Environment configuration
├── internal/
│   ├── domain/           # Entities & repository interfaces
│   ├── repository/       # PostgreSQL + Redis implementations
│   ├── usecase/          # Business logic layer
│   ├── handler/          # HTTP handlers & router
│   └── middleware/       # Auth & other middleware
└── migrations/           # SQL migration files
```

---

## 🔄 CI/CD Pipeline

GitHub Actions runs on every push to `main` or `develop`:

1. **Build** — compiles all 9 services
2. **Security Scan** — Gosec static analysis
3. **Docker Build** — validates Dockerfile for auth-service and api-gateway

---

## 📚 Architecture Decision Records

See [`docs/adr/`](docs/adr/) for key architectural decisions:

- [ADR-001](docs/adr/001-microservices.md) — Why microservices over monolith
- [ADR-002](docs/adr/002-kafka.md) — Event-driven architecture with Kafka
- [ADR-003](docs/adr/003-distributed-locking.md) — Distributed locking strategy
- [ADR-004](docs/adr/004-idempotency.md) — Idempotency pattern for financial operations

---

## 🗺️ Roadmap

- [ ] gRPC for internal service communication
- [ ] Elasticsearch for transaction search
- [ ] Kubernetes manifests (Helm charts)
- [ ] Distributed tracing with Jaeger
- [ ] Unit & integration tests
- [ ] Swagger/OpenAPI auto-generated docs
- [ ] Redis Sentinel for high availability
- [ ] Circuit breaker pattern

---

## 🤝 Author

**Muhammad Afif** — Backend Engineer

[![GitHub](https://img.shields.io/badge/GitHub-afifn11-181717?style=flat&logo=github)](https://github.com/afifn11)

---

<div align="center">
Built with ❤️ using Go · PostgreSQL · Redis · Kafka
</div>
