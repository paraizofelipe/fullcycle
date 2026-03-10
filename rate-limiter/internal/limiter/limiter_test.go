package limiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/paraizofelipe/fullcycle/rate-limiter/internal/limiter"
	"github.com/paraizofelipe/fullcycle/rate-limiter/internal/storage"
)

func newLimiter(ipLimit, tokenLimit int, blockDuration time.Duration, tokens map[string]int) *limiter.RateLimiter {
	return limiter.New(storage.NewMemoryStorage(), ipLimit, tokenLimit, blockDuration, tokens)
}

func TestIPRateLimit(t *testing.T) {
	rl := newLimiter(5, 100, time.Minute, nil)
	ctx := context.Background()

	for i := 1; i <= 5; i++ {
		allowed, err := rl.Allow(ctx, "10.0.0.1", "")
		if err != nil {
			t.Fatalf("request %d: unexpected error: %v", i, err)
		}
		if !allowed {
			t.Fatalf("request %d: should be allowed (limit=5)", i)
		}
	}

	allowed, err := rl.Allow(ctx, "10.0.0.1", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("6th request should be blocked (IP limit=5)")
	}
}

func TestTokenRateLimit(t *testing.T) {
	rl := newLimiter(5, 10, time.Minute, nil)
	ctx := context.Background()

	for i := 1; i <= 10; i++ {
		allowed, err := rl.Allow(ctx, "10.0.0.1", "my-token")
		if err != nil {
			t.Fatalf("request %d: unexpected error: %v", i, err)
		}
		if !allowed {
			t.Fatalf("request %d: should be allowed (token limit=10)", i)
		}
	}

	allowed, err := rl.Allow(ctx, "10.0.0.1", "my-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("11th request should be blocked (token limit=10)")
	}
}

func TestTokenPrecedenceOverIP(t *testing.T) {
	tokens := map[string]int{"premium": 20}
	rl := newLimiter(5, 10, time.Minute, tokens)
	ctx := context.Background()

	for i := 1; i <= 20; i++ {
		allowed, err := rl.Allow(ctx, "10.0.0.1", "premium")
		if err != nil {
			t.Fatalf("request %d: unexpected error: %v", i, err)
		}
		if !allowed {
			t.Fatalf("request %d: should be allowed — token 'premium' limit=20 overrides IP limit=5", i)
		}
	}

	allowed, err := rl.Allow(ctx, "10.0.0.1", "premium")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("21st request should be blocked (token limit=20 reached)")
	}
}

func TestTokenDefaultRateOverridesIP(t *testing.T) {
	rl := newLimiter(3, 8, time.Minute, nil)
	ctx := context.Background()

	for i := 1; i <= 8; i++ {
		allowed, err := rl.Allow(ctx, "10.0.0.2", "standard-token")
		if err != nil {
			t.Fatalf("request %d: unexpected error: %v", i, err)
		}
		if !allowed {
			t.Fatalf("request %d: should be allowed — token default limit=8 overrides IP limit=3", i)
		}
	}

	allowed, _ := rl.Allow(ctx, "10.0.0.2", "standard-token")
	if allowed {
		t.Fatal("9th request should be blocked (default token limit=8)")
	}
}

func TestBlockDurationExpiry(t *testing.T) {
	rl := newLimiter(2, 100, 100*time.Millisecond, nil)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		rl.Allow(ctx, "10.0.0.3", "")
	}

	allowed, _ := rl.Allow(ctx, "10.0.0.3", "")
	if allowed {
		t.Fatal("should be blocked")
	}

	time.Sleep(150 * time.Millisecond)

	allowed, err := rl.Allow(ctx, "10.0.0.3", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Fatal("should be allowed after block duration expires")
	}
}

func TestDifferentIPsAreIndependent(t *testing.T) {
	rl := newLimiter(2, 100, time.Minute, nil)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		rl.Allow(ctx, "10.0.0.10", "")
	}

	allowed, err := rl.Allow(ctx, "10.0.0.11", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Fatal("IP 10.0.0.11 should not be affected by IP 10.0.0.10 being blocked")
	}
}

func TestDifferentTokensAreIndependent(t *testing.T) {
	rl := newLimiter(100, 2, time.Minute, nil)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		rl.Allow(ctx, "10.0.0.1", "token-A")
	}

	allowed, err := rl.Allow(ctx, "10.0.0.1", "token-B")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Fatal("token-B should not be affected by token-A being blocked")
	}
}

func TestBlockedRequestRemainsBlocked(t *testing.T) {
	rl := newLimiter(2, 100, 5*time.Second, nil)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		rl.Allow(ctx, "10.0.0.20", "")
	}

	for i := 0; i < 5; i++ {
		allowed, _ := rl.Allow(ctx, "10.0.0.20", "")
		if allowed {
			t.Fatalf("blocked request %d should not be allowed", i+1)
		}
	}
}
