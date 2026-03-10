package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/paraizofelipe/fullcycle/rate-limiter/internal/storage"
)

type RateLimiter struct {
	storage        storage.Storage
	ipRateLimit    int
	tokenRateLimit int
	blockDuration  time.Duration
	tokens         map[string]int
	window         time.Duration
}

func New(
	s storage.Storage,
	ipRateLimit, tokenRateLimit int,
	blockDuration time.Duration,
	tokens map[string]int,
) *RateLimiter {
	if tokens == nil {
		tokens = make(map[string]int)
	}
	return &RateLimiter{
		storage:        s,
		ipRateLimit:    ipRateLimit,
		tokenRateLimit: tokenRateLimit,
		blockDuration:  blockDuration,
		tokens:         tokens,
		window:         time.Second,
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, ip, token string) (bool, error) {
	key, limit := rl.resolveKeyAndLimit(ip, token)

	blocked, err := rl.storage.IsBlocked(ctx, key)
	if err != nil {
		return false, fmt.Errorf("checking block status: %w", err)
	}
	if blocked {
		return false, nil
	}

	count, err := rl.storage.Increment(ctx, key, rl.window)
	if err != nil {
		return false, fmt.Errorf("incrementing counter: %w", err)
	}

	if int(count) > limit {
		if err := rl.storage.Block(ctx, key, rl.blockDuration); err != nil {
			return false, fmt.Errorf("setting block: %w", err)
		}
		return false, nil
	}

	return true, nil
}

func (rl *RateLimiter) resolveKeyAndLimit(ip, token string) (string, int) {
	if token != "" {
		limit := rl.tokenRateLimit
		if specific, ok := rl.tokens[token]; ok {
			limit = specific
		}
		return "token:" + token, limit
	}
	return "ip:" + ip, rl.ipRateLimit
}
