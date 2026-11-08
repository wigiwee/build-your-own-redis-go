package storage

import (
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]Entry
}

type Entry struct {
	Type      string
	List      []string
	Value     string
	ExpiresAt time.Time
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]Entry),
	}
}

func (s *Store) Set(key, value string, ttl time.Duration) {
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
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.data[key]

	if !exists && entry.Type == utils.TYPE_STRING {
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

func (s *Store) PushList(listName string, values []string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, doesExist := s.data[listName]
	if doesExist {
		entry.List = append(entry.List, values...)
		s.data[listName] = entry
		return len(s.data[listName].List)
	}
	newEntry := Entry{
		Type: utils.TYPE_LIST,
		List: values,
	}
	s.data[listName] = newEntry
	return len(values)
}

func (s *Store) GetListRange(listName string, start, stop int) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, doesExist := s.data[listName]
	if !doesExist && start > len(entry.List) && start > stop {
		return []string{}
	}

	if stop >= len(entry.List) {
		stop = len(entry.List) - 1
	}
	if start < 0 {
		start = 0
	}

	return entry.List[start : stop+1]

}
