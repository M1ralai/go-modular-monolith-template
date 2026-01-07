package jobs

import (
	"context"
	"fmt"
	"hash/fnv"

	"github.com/jmoiron/sqlx"
)

// DistributedLock provides distributed locking using PostgreSQL advisory locks
type DistributedLock struct {
	db *sqlx.DB
}

// NewDistributedLock creates a new distributed lock manager
func NewDistributedLock(db *sqlx.DB) *DistributedLock {
	return &DistributedLock{db: db}
}

// LockKey generates a unique lock key from job name and entity ID
func LockKey(jobName string, entityID int) int64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%s:%d", jobName, entityID)))
	return int64(h.Sum64())
}

// TryLock attempts to acquire a distributed lock
// Returns true if lock was acquired, false if already held
func (dl *DistributedLock) TryLock(ctx context.Context, lockKey int64) (bool, error) {
	var acquired bool
	err := dl.db.GetContext(ctx, &acquired,
		"SELECT pg_try_advisory_lock($1)",
		lockKey,
	)
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}
	return acquired, nil
}

// Lock blocks until lock is acquired or context is cancelled
func (dl *DistributedLock) Lock(ctx context.Context, lockKey int64) error {
	_, err := dl.db.ExecContext(ctx,
		"SELECT pg_advisory_lock($1)",
		lockKey,
	)
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	return nil
}

// Unlock releases a distributed lock
func (dl *DistributedLock) Unlock(ctx context.Context, lockKey int64) error {
	_, err := dl.db.ExecContext(ctx,
		"SELECT pg_advisory_unlock($1)",
		lockKey,
	)
	if err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}
	return nil
}

// WithLock executes a function with a distributed lock
// Automatically releases lock after execution
func (dl *DistributedLock) WithLock(ctx context.Context, lockKey int64, fn func() error) error {
	acquired, err := dl.TryLock(ctx, lockKey)
	if err != nil {
		return err
	}
	if !acquired {
		return fmt.Errorf("lock already held for key %d", lockKey)
	}

	defer func() {
		// Always try to unlock, even if fn() panics
		if unlockErr := dl.Unlock(context.Background(), lockKey); unlockErr != nil {
			// Log error but don't fail - lock will timeout eventually
		}
	}()

	return fn()
}
