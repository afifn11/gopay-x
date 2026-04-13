# ADR-001: Microservices Architecture

**Date:** 2026
**Status:** Accepted

## Context

GoPayX needs to handle multiple distinct business domains: authentication, wallet management, payments, fraud detection, and auditing. Each domain has different scaling requirements, change frequencies, and data ownership concerns.

A monolithic approach would result in a single large codebase where changes to the fraud detection logic require redeploying the entire application, including the authentication and payment modules.

## Decision

Adopt microservices architecture with 8 independent services communicating via REST (synchronous) and Kafka events (asynchronous).

## Rationale

- **Independent scaling** — wallet-service experiences 10x more load during payday than audit-service. Microservices allow scaling only what's needed.
- **Domain boundaries** — each service owns its data and business logic, reducing accidental coupling.
- **Team autonomy** — in a real organization, each service can be owned and deployed by a separate team independently.
- **Technology flexibility** — services can evolve independently (e.g., migrate audit-service to Elasticsearch without affecting payments).
- **Fault isolation** — a crash in notification-service does not affect payment processing.

## Consequences

**Positive:**
- Independent deployment and scaling per service
- Clear domain boundaries enforce separation of concerns
- Easier to understand each service in isolation
- Failure in one service doesn't cascade to others

**Negative:**
- Increased operational complexity (9 processes instead of 1)
- Distributed system challenges: network failures, eventual consistency, distributed transactions
- More infrastructure overhead (databases per service)
- Cross-service debugging is harder

## Alternatives Considered

**Modular Monolith:** Single deployable with internal module boundaries. Simpler to run locally, but harder to scale independently. Would be a valid choice for early-stage startups.

**Decision:** Microservices chosen to demonstrate distributed systems knowledge relevant to industry-scale fintech platforms.
