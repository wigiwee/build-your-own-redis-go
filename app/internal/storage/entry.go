package storage

import "time"

type Entry struct {
	Type string

	List      []string
	Value     string
	ExpiresAt time.Time

	Stream *Stream
}
type Stream struct {
	Entries []StreamEntry
	LastID  StreamID
}

type StreamID struct {
	Milliseconds uint64
	Sequence     uint64
}

type StreamEntry struct {
	ID     string
	Fields map[string]string
}
