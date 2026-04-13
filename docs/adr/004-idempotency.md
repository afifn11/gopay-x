# ADR-004: Idempotency Keys for Financial Operations

**Date:** 2026
**Status:** Accepted

## Context

Network failures are unavoidable. When a client sends a top-up request and the connection drops before receiving a response, it cannot know if the server processed it or not. Without idempotency, retrying the request results in duplicate transactions — the user gets charged twice.

```
Client → POST /topup (amount: 100k) → Server processes ✅
       ← connection drops before response
Client → POST /topup (amount: 100k) → Server processes AGAIN ❌ (duplicate!)
```

This is an especially critical problem in financial systems.

## Decision

Require a client-generated `reference_id` field on all financial operations (top-up, transfer, withdraw). Before processing, the server checks if this ID has been used before.

## Implementation

```go
// 1. Client generates a unique reference_id
// POST /api/v1/wallets/topup
// { "amount": 100000, "reference_id": "topup-user123-1704067200" }

// 2. Server checks idempotency before processing
func (uc *walletUsecase) TopUp(ctx context.Context, req *TopUpRequest) (*WalletTransaction, error) {
    existing, _ := uc.txRepo.FindByReferenceID(ctx, req.ReferenceID)
    if existing != nil {
        return existing, ErrDuplicateTransaction  // return original result safely
    }
    // ... process new transaction
    // reference_id stored with UNIQUE constraint in database
}
```

## Reference ID Format (recommended to clients)

```
{operation}-{userID}-{timestamp}
topup-550e8400-e29b-41d4-a716-446655440000-1704067200
transfer-550e8400-e29b-41d4-a716-446655440000-1704067201
```

## Database Enforcement

```sql
CREATE TABLE wallet_transactions (
    id UUID PRIMARY KEY,
    reference_id TEXT UNIQUE NOT NULL,  -- unique constraint enforces at DB level
    ...
);
```

The `UNIQUE` constraint on `reference_id` provides a safety net even if application-level checks have a race condition.

## Rationale

- **Safety** — retrying a request with the same `reference_id` returns the original result without re-processing
- **Client control** — clients generate their own keys, allowing them to track request state
- **Database enforcement** — unique constraint guarantees idempotency even under concurrent retries
- **Auditability** — reference IDs link client-side requests to server-side records

## Consequences

**Positive:**
- Safe retries — mobile apps can retry on network failure without risk
- Prevents duplicate charges and transfers
- Reference IDs serve as correlation IDs for debugging

**Negative:**
- Clients must generate and manage reference IDs (added client complexity)
- Storage overhead for reference_id index on transaction tables
- Old reference IDs accumulate (consider TTL or archival strategy)

## Affected Endpoints

| Endpoint | Required Field |
|---|---|
| `POST /api/v1/wallets/topup` | `reference_id` |
| `POST /api/v1/payments/transfer` | `reference_id` |
| `POST /api/v1/payments/topup` | `reference_id` |

## Alternatives Considered

**Server-generated idempotency keys:** Server generates and returns a key; client uses it for retry. Adds a round-trip but removes client responsibility. Used by Stripe.

**Decision:** Client-generated keys chosen for simplicity and to match common fintech API patterns (similar to Xendit and Midtrans).
