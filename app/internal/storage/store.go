package storage

import (
	"log"
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

func (s *Store) removeWaiterLocked(listName string, waiter *Waiter) {
	waitersList := s.waiters[listName]
	idx := -1
	for i := 0; i < len(waitersList); i++ {
		if waitersList[i] == waiter {
			idx = i
			break
		}
	}
	if idx == -1 {
		s.mu.Unlock()
		log.Println("fatal you should take a look: BLOP waiter not found when trying to remove the waiter")
		return
	}
	waitersList = append(waitersList[:idx], waitersList[idx+1:]...)
	s.waiters[listName] = waitersList
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
