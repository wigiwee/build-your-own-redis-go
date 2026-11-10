package storage

import (
	"log"
	"slices"
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
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

func (s *Store) RPushList(listName string, values []string) int {
	s.refreshKey(listName)
	s.mu.Lock()
	defer s.mu.Unlock()
	initialLen := len(s.data[listName].List)
	for _, value := range values {
		queue := s.waiters[listName]
		if len(queue) > 0 {
			queue[0].notified = true
			queue[0].ch <- []string{listName, value}
			s.waiters[listName] = s.waiters[listName][1:]
			continue
		}

		entry := s.data[listName]
		entry.Type = utils.TYPE_LIST
		entry.List = append(entry.List, value)
		s.data[listName] = entry
	}
	if len(s.waiters[listName]) == 0 {
		delete(s.waiters, listName)
	}

	return initialLen + len(values)
}

func (s *Store) GetListRange(listName string, start, stop int) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, doesExist := s.data[listName]
	if start < 0 {
		start = len(entry.List) + start
	}
	if stop < 0 {
		stop = len(entry.List) + stop
	}
	if start < 0 {
		start = 0
	}

	if !doesExist || start > len(entry.List) && start > stop {
		return []string{}
	}

	if stop >= len(entry.List) {
		stop = len(entry.List) - 1
	}

	return entry.List[start : stop+1]
}

func (s *Store) LPushList(listName string, values []string) int {
	s.refreshKey(listName)
	s.mu.Lock()
	defer s.mu.Unlock()

	values = slices.Clone(values)
	slices.Reverse(values)

	initialLen := len(s.data[listName].List)

	entry := s.data[listName]
	entry.Type = utils.TYPE_LIST
	entry.List = append(values, entry.List...)

	queue := s.waiters[listName]
	for len(queue) != 0 && len(entry.List) != 0 {
		queue[0].notified = true
		queue[0].ch <- []string{listName, entry.List[0]}
		entry.List = entry.List[1:]
		queue = queue[1:]
	}
	s.waiters[listName] = queue
	s.data[listName] = entry
	if len(s.waiters[listName]) == 0 {
		delete(s.waiters, listName)
	}

	return initialLen + len(values)
}

func (s *Store) GetListLen(listName string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, doesExist := s.data[listName]
	if !doesExist {
		return 0
	} else {
		return len(entry.List)
	}
}

func (s *Store) LPop(listName string, n int) []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, doesExist := s.data[listName]

	if !doesExist && len(entry.List) == 0 {
		return nil
	}
	if n > len(entry.List) {
		n = len(entry.List)
	}
	popValue := entry.List[:n]
	entry.List = entry.List[n:]
	s.data[listName] = entry
	return popValue
}

/*
check if there are already elements in the list

	|
	V

if yes then pop it, if not then add ch to the waiter

	|
	V

set a timeout timer on ttl

	|
	V

then using select wait untill either select or rpush sends event through channel
*/
func (s *Store) BLPop(listName string, ttl time.Duration) []string {
	s.mu.Lock()

	if entry, doesExist := s.data[listName]; doesExist && !s.isExpired(entry) && len(entry.List) > 0 {
		value := entry.List[0]
		entry.List = entry.List[1:]
		s.data[listName] = entry
		s.mu.Unlock()
		return []string{listName, value}
	}

	waiter := &Waiter{
		ch:       make(chan []string, 1),
		notified: false,
	}
	s.waiters[listName] = append(s.waiters[listName], waiter)
	s.mu.Unlock()

	//case where ttl = 0
	if ttl <= 0 {
		result := <-waiter.ch
		s.removeWaiter(listName, waiter)
		return result
	}

	timer := time.NewTimer(ttl)
	defer timer.Stop()

	select {
	case result := <-waiter.ch:
		// s.removeWaiter(listName, waiter)
		return result

	case <-timer.C:
		s.mu.Lock()
		if waiter.notified {
			s.mu.Unlock()
			return <-waiter.ch
		}
		// s.removeWaiterLocked(listName, waiter)

		s.mu.Unlock()
		return nil
	}
}
func (s *Store) AddStream(streamkey string, id string, args []string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fields := make(map[string]string)
	for i := 0; i < len(args); i = i + 2 {
		fields[args[i]] = args[i+1]
	}
	entry := &Entry{
		Type: utils.TYPE_STREAM,
		Stream: &Stream{
			Entries: []StreamEntry{
				StreamEntry{
					ID:     id,
					Fields: fields,
				},
			},
		},
	}
	s.data[streamkey] = *entry
	return s.data[streamkey].Stream.Entries[len(s.data[streamkey].Stream.Entries)-1].ID, nil

}
func (s *Store) GetType(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, doesExist := s.data[key]
	if !doesExist {
		return "none"
	} else {
		return entry.Type
	}
}
func (s *Store) isExpired(entry Entry) bool {
	return !entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt)
}

func (s *Store) removeWaiter(listName string, waiter *Waiter) {
	s.mu.Lock()
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
	s.mu.Unlock()
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
