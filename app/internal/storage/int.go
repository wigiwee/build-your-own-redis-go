package storage

import (
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

func (s *Store) Incr(key string) (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, exists := s.data[key]
	if !exists {
		s.data[key] = Entry{
			Type:     utils.TYPE_INT,
			IntValue: 1,
		}
		return 1, nil
	}
	if entry.Type != utils.TYPE_INT {
		return 0, errors.New("ERR value is not an integer or out of range")
	}
	entry.IntValue = entry.IntValue + 1
	s.data[key] = entry
	s.versions[key]++
	return entry.IntValue, nil
}
