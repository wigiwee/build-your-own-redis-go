package storage

import (
	"slices"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

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
	entry, exists := s.data[listName]
	if start < 0 {
		start = len(entry.List) + start
	}
	if stop < 0 {
		stop = len(entry.List) + stop
	}
	if start < 0 {
		start = 0
	}

	if !exists || start > len(entry.List) && start > stop {
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

	entry, exists := s.data[listName]
	if !exists {
		return 0
	} else {
		return len(entry.List)
	}
}

func (s *Store) LPop(listName string, n int) []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.data[listName]

	if !exists && len(entry.List) == 0 {
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

	if entry, exists := s.data[listName]; exists && !s.isExpired(entry) && len(entry.List) > 0 {
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
		s.removeWaiter(listName, waiter)
		return result

	case <-timer.C:
		s.mu.Lock()
		if waiter.notified {
			s.mu.Unlock()
			return <-waiter.ch
		}
		s.removeWaiterLocked(listName, waiter)

		s.mu.Unlock()
		return nil
	}
}
