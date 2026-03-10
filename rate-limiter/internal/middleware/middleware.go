package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/paraizofelipe/fullcycle/rate-limiter/internal/limiter"
)

const (
	headerAPIKey    = "API_KEY"
	blockedResponse = "you have reached the maximum number of requests or actions allowed within a certain time frame"
)

func RateLimiter(rl *limiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := extractIP(r)
			token := r.Header.Get(headerAPIKey)

			allowed, err := rl.Allow(r.Context(), ip, token)
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(blockedResponse))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func extractIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if ip := strings.TrimSpace(parts[0]); ip != "" {
			return ip
		}
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
