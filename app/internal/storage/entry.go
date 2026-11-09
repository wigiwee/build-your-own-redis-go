package storage

import "time"

type Entry struct {
	Type      string
	List      []string
	Value     string
	ExpiresAt time.Time
}
