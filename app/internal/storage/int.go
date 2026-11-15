package storage

import (
	"errors"
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (s *Store) Incr(key string) (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, exists := s.data[key]
	fmt.Println(entry, exists)
	if !exists {
		return 0, errors.New("key not found")
	}
	if entry.Type != utils.TYPE_INT {
		return 0, errors.New("key not found")
	}
	entry.IntValue = entry.IntValue + 1
	s.data[key] = entry
	return entry.IntValue, nil
}
