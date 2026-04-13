# ADR-002: Apache Kafka for Event-Driven Communication

**Date:** 2026
**Status:** Accepted

## Context

Services like notification, fraud detection, and audit need to react to events from payment and wallet services. The naive approach — direct HTTP calls from payment-service to each consumer — creates tight coupling:

```
payment-service → HTTP → notification-service  (tight coupling)
payment-service → HTTP → fraud-service         (tight coupling)
payment-service → HTTP → audit-service         (tight coupling)
```

If notification-service is down, payment-service would fail or need complex retry logic.

## Decision

Use Apache Kafka as the central event bus for all async inter-service communication.

```
payment-service → Kafka topic → notification-service  (loose coupling)
                             → fraud-service
                             → audit-service
```

## Kafka Topics

| Topic | Producer | Consumers |
|---|---|---|
| `payment.created` | payment-service | fraud-detection, audit |
| `payment.success` | payment-service | notification, audit |
| `topup.success` | wallet-service | notification, transaction, audit |
| `transfer.success` | payment-service | notification, transaction, audit |
| `user.registered` | auth-service | audit |
| `login.new_device` | auth-service | notification, audit |
| `fraud.flagged` | fraud-detection | audit |

## Rationale

- **Decoupling** — payment-service publishes events without knowing who consumes them. New consumers can be added without changing the producer.
- **Durability** — Kafka persists messages to disk. If fraud-service is down, it processes missed events when it recovers.
- **Replayability** — events can be replayed for debugging, re-processing, or data migration.
- **Fan-out** — one event can be consumed by multiple services simultaneously.
- **Audit trail** — Kafka topics themselves serve as an ordered event log.

## Consequences

**Positive:**
- Adding a new consumer requires zero changes to existing producers
- Services can recover from downtime by replaying Kafka messages
- Natural audit trail via Kafka topic retention

**Negative:**
- Eventual consistency — notification is delivered slightly after payment, not instantly
- Kafka adds operational overhead (Zookeeper, broker management)
- Debugging async flows is harder than synchronous HTTP

## Alternatives Considered

**Direct HTTP calls:** Simple but creates tight coupling and cascading failures.
**RabbitMQ:** Good for task queues, but Kafka's log-based model better fits event sourcing and replay needs.
