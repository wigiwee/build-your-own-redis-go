package storage

import (
	"strconv"
	"time"
)

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
	ID     StreamID
	Fields map[string]string
}

func (sId StreamID) ToString() string {
	return strconv.FormatUint(sId.Milliseconds, 10) + "-" + strconv.FormatUint(sId.Sequence, 10)
}
