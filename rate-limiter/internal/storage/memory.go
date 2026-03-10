package storage

import (
	"context"
	"sync"
	"time"
)

type counterEntry struct {
	count     int64
	expiresAt time.Time
}

type MemoryStorage struct {
	mu       sync.Mutex
	counters map[string]*counterEntry
	blocked  map[string]time.Time
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		counters: make(map[string]*counterEntry),
		blocked:  make(map[string]time.Time),
	}
}

func (m *MemoryStorage) Increment(_ context.Context, key string, window time.Duration) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	e, exists := m.counters[key]
	if !exists || now.After(e.expiresAt) {
		m.counters[key] = &counterEntry{count: 1, expiresAt: now.Add(window)}
		return 1, nil
	}
	e.count++
	return e.count, nil
}

func (m *MemoryStorage) IsBlocked(_ context.Context, key string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	expiry, exists := m.blocked[key]
	if !exists {
		return false, nil
	}
	if time.Now().After(expiry) {
		delete(m.blocked, key)
		return false, nil
	}
	return true, nil
}

func (m *MemoryStorage) Block(_ context.Context, key string, duration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.blocked[key] = time.Now().Add(duration)
	delete(m.counters, key)
	return nil
}
