# ADR-003: Redis Distributed Locking for Wallet Operations

**Date:** 2026
**Status:** Accepted

## Context

Wallet balance updates are critical financial operations that must be atomic. Without locking, concurrent requests create race conditions:

```
Time  Request A (top-up +100k)     Request B (top-up +50k)
T1    READ balance = 500k           READ balance = 500k
T2    WRITE balance = 600k          WRITE balance = 550k  ← overwrites A!
T3    Result: 550k (wrong!)         Should be: 650k
```

This is the classic double-spend / lost update problem in distributed systems.

## Decision

Use Redis `SetNX` (Set if Not eXists) for distributed locking on all wallet mutation operations.

## Implementation

```go
lockKey := fmt.Sprintf("lock:wallet:topup:%s", userID.String())

// Try to acquire lock (atomic SetNX)
acquired, err := redis.SetNX(ctx, lockKey, "1", 10*time.Second).Result()
if !acquired {
    return nil, ErrLockFailed  // HTTP 429 — client should retry
}
defer redis.Del(ctx, lockKey)  // always release

// Safe to update balance now
wallet.Balance += req.Amount
repo.UpdateBalance(ctx, wallet.ID, wallet.Balance)
```

## Rationale

- **Atomicity** — Redis `SetNX` is atomic; only one process can acquire the lock
- **TTL safety** — 10-second TTL ensures lock auto-expires if the process crashes, preventing deadlock
- **Low latency** — Redis operations complete in < 1ms, minimal overhead
- **Simplicity** — no external lock manager needed; Redis is already in the stack for caching

## Lock Key Strategy

```
lock:wallet:topup:{userID}    — per-user top-up lock
lock:payment:transfer:{userID} — per-user transfer lock
```

Per-user (not per-wallet) locking ensures serialization of all financial ops for a given user.

## Consequences

**Positive:**
- Eliminates race conditions and double-spend on wallet operations
- Auto-expiry prevents deadlocks on process crash
- Simple to implement with existing Redis infrastructure

**Negative:**
- Lock contention under very high concurrent load for same user
- Redis is a single point of failure (production: use Redis Sentinel or Cluster)
- Failed lock returns HTTP 429 — client must implement retry with backoff

## Alternatives Considered

**Database row-level locking (`SELECT FOR UPDATE`):** Works but holds a DB connection for the duration, reducing throughput under high load.

**Optimistic locking (version field):** No blocking, but requires retry logic on conflict. Better for low-contention scenarios.

**Decision:** Redis distributed locking chosen for balance between simplicity, performance, and safety given Redis is already a project dependency.
