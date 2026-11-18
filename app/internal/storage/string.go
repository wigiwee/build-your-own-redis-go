package storage

import (
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (s *Store) Get(key string) any {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.data[key]

	//todo: to be revised
	if !exists {
		return nil
	}

	if !entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt) {
		delete(s.data, key)
		return nil
	}
	switch entry.Type {
	case utils.TYPE_STRING:
		return entry.StringValue
	case utils.TYPE_INT:
		return strconv.FormatUint(entry.IntValue, 10)
	default:
		return nil
	}
}

func (s *Store) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; !exists {
		return false
	}
	s.versions[key]++
	delete(s.data, key)
	return true
}
