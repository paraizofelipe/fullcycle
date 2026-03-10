package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/paraizofelipe/fullcycle/rate-limiter/config"
	"github.com/paraizofelipe/fullcycle/rate-limiter/internal/limiter"
	"github.com/paraizofelipe/fullcycle/rate-limiter/internal/middleware"
	"github.com/paraizofelipe/fullcycle/rate-limiter/internal/storage"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	redisStorage := storage.NewRedisStorage(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword)
	rl := limiter.New(redisStorage, cfg.IPRateLimit, cfg.TokenRateLimit, cfg.BlockDuration, cfg.Tokens)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	handler := middleware.RateLimiter(rl)(mux)

	log.Printf("Server listening on :8080 (IP limit: %d req/s, Token limit: %d req/s, Block: %s)",
		cfg.IPRateLimit, cfg.TokenRateLimit, cfg.BlockDuration)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
