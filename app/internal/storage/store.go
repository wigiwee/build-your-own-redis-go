package storage

import (
	"sync"
	"time"
)

type Store struct {
	mu      sync.RWMutex
	data    map[string]Entry
	waiters map[string][]*Waiter
}

type Waiter struct {
	ch       chan []string
	notified bool
}

func NewStore() *Store {
	s := &Store{
		data:    make(map[string]Entry),
		waiters: make(map[string][]*Waiter),
	}
	return s

}

func (s *Store) GetType(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.data[key]
	if !exists {
		return "none"
	} else {
		return entry.Type
	}
}
func (s *Store) isExpired(entry Entry) bool {
	return !entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt)
}

func (s *Store) refreshKey(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry := s.data[key]
	if !entry.ExpiresAt.IsZero() {
		if time.Now().After(entry.ExpiresAt) {
			delete(s.data, key)
		}
	}
}
