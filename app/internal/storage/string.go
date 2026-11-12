package storage

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (s *Store) Set(key, value string, ttl time.Duration) {
	s.refreshKey(key)
	s.mu.Lock()
	defer s.mu.Unlock()
	entry := Entry{
		Type:  utils.TYPE_STRING,
		Value: value,
	}
	if ttl > 0 {
		entry.ExpiresAt = time.Now().Add(ttl)
	}
	s.data[key] = entry
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.data[key]

	//todo: to be revised
	if !exists || entry.Type != utils.TYPE_STRING {
		return "", false
	}
	if !entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt) {
		delete(s.data, key)
		return "", false
	}
	return entry.Value, true
}

func (s *Store) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; !exists {
		return false
	}
	delete(s.data, key)
	return true
}
