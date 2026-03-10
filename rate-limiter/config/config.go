package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	RedisHost      string
	RedisPort      string
	RedisPassword  string
	IPRateLimit    int
	TokenRateLimit int
	BlockDuration  time.Duration
	Tokens         map[string]int
}

func Load() *Config {
	return &Config{
		RedisHost:      getEnv("REDIS_HOST", "localhost"),
		RedisPort:      getEnv("REDIS_PORT", "6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		IPRateLimit:    getEnvInt("IP_RATE_LIMIT", 10),
		TokenRateLimit: getEnvInt("TOKEN_RATE_LIMIT", 100),
		BlockDuration:  time.Duration(getEnvInt("BLOCK_DURATION_SECONDS", 300)) * time.Second,
		Tokens:         parseTokens(getEnv("TOKENS", "")),
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultValue
}

func parseTokens(s string) map[string]int {
	tokens := make(map[string]int)
	if s == "" {
		return tokens
	}
	for _, pair := range strings.Split(s, ",") {
		parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
		if len(parts) == 2 {
			token := strings.TrimSpace(parts[0])
			if limit, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil && token != "" {
				tokens[token] = limit
			}
		}
	}
	return tokens
}
