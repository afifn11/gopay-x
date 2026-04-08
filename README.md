# GoPayX 💳

> Enterprise-grade Digital Wallet & Payment Platform — built with Go Microservices

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Architecture](https://img.shields.io/badge/Architecture-Microservices-blueviolet)]()

## Overview

GoPayX is a production-ready digital wallet and payment gateway platform
demonstrating industry-standard backend engineering practices using Go.

## Architecture

- **8 Microservices** communicating via REST, gRPC, and Kafka events
- **Event-driven** architecture with Apache Kafka
- **Clean Architecture** pattern within each service
- **Distributed locking** with Redis for financial consistency
- **Full audit trail** with Elasticsearch

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.22+ |
| Framework | Gin (HTTP), gRPC |
| Database | PostgreSQL 16 |
| Cache / Lock | Redis 7 |
| Message Broker | Apache Kafka |
| Containerization | Docker + Docker Compose |
| Monitoring | Prometheus + Grafana |
| Tracing | Jaeger |
| CI/CD | GitHub Actions |

## Services

| Service | Port | Responsibility |
|---|---|---|
| api-gateway | 8080 | Routing, rate limiting, auth check |
| auth-service | 8081 | JWT auth, OAuth2, KYC flow |
| user-service | 8082 | User profile management |
| wallet-service | 8083 | Balance, top-up, locking |
| payment-service | 8084 | Transfer, payment gateway |
| transaction-service | 8085 | Ledger, history, reconciliation |
| notification-service | 8086 | Email, push, SMS (event-driven) |
| fraud-detection-service | 8087 | Rule engine, scoring |
| audit-service | 8088 | Immutable event log |

## Getting Started

```bash
# Start all infrastructure
docker compose -f infra/docker-compose.yml up -d

# Run auth-service
cd services/auth-service
go run cmd/main.go
```

## Project Status

- [x] Project structure & infrastructure setup
- [ ] auth-service
- [ ] user-service
- [ ] wallet-service
- [ ] payment-service
- [ ] transaction-service
- [ ] notification-service
- [ ] fraud-detection-service
- [ ] audit-service
- [ ] API Gateway
- [ ] CI/CD Pipeline
- [ ] Monitoring & Tracing