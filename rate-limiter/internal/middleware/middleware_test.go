package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/paraizofelipe/fullcycle/rate-limiter/internal/limiter"
	"github.com/paraizofelipe/fullcycle/rate-limiter/internal/middleware"
	"github.com/paraizofelipe/fullcycle/rate-limiter/internal/storage"
)

const blockedBody = "you have reached the maximum number of requests or actions allowed within a certain time frame"

func newHandler(ipLimit, tokenLimit int, tokens map[string]int) http.Handler {
	rl := limiter.New(storage.NewMemoryStorage(), ipLimit, tokenLimit, time.Minute, tokens)
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return middleware.RateLimiter(rl)(base)
}

func sendRequest(handler http.Handler, remoteAddr, apiKey string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = remoteAddr
	if apiKey != "" {
		req.Header.Set("API_KEY", apiKey)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func TestMiddleware_IPRateLimited(t *testing.T) {
	handler := newHandler(3, 100, nil)

	for i := 1; i <= 3; i++ {
		rec := sendRequest(handler, "192.168.0.1:1000", "")
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i, rec.Code)
		}
	}

	rec := sendRequest(handler, "192.168.0.1:1000", "")
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", rec.Code)
	}
	if rec.Body.String() != blockedBody {
		t.Fatalf("unexpected body: %q", rec.Body.String())
	}
}

func TestMiddleware_TokenRateLimited(t *testing.T) {
	handler := newHandler(100, 3, nil)

	for i := 1; i <= 3; i++ {
		rec := sendRequest(handler, "192.168.0.2:1000", "my-token")
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i, rec.Code)
		}
	}

	rec := sendRequest(handler, "192.168.0.2:1000", "my-token")
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", rec.Code)
	}
	if rec.Body.String() != blockedBody {
		t.Fatalf("unexpected body: %q", rec.Body.String())
	}
}

func TestMiddleware_TokenPrecedence(t *testing.T) {
	tokens := map[string]int{"premium": 10}
	handler := newHandler(2, 5, tokens)

	for i := 1; i <= 10; i++ {
		rec := sendRequest(handler, "192.168.0.3:1000", "premium")
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d (token 'premium' limit=10 should override IP limit=2)", i, rec.Code)
		}
	}

	rec := sendRequest(handler, "192.168.0.3:1000", "premium")
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 after 10 requests, got %d", rec.Code)
	}
}

func TestMiddleware_XForwardedFor(t *testing.T) {
	handler := newHandler(2, 100, nil)

	for i := 1; i <= 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "10.0.0.1:9999"
		req.Header.Set("X-Forwarded-For", "203.0.113.5")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i, rec.Code)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.1:9999"
	req.Header.Set("X-Forwarded-For", "203.0.113.5")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 for real IP 203.0.113.5, got %d", rec.Code)
	}
}
