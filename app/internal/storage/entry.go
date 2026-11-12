package storage

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
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
	SequenceNum  uint64
}

type StreamEntry struct {
	ID     StreamID
	Fields []string
}

func (streamEntry StreamEntry) SerializeRESP() (string, error) {
	serId, err := utils.Serialize(streamEntry.ID.ToString())
	if err != nil {
		return "", err
	}
	serFeilds, err := utils.Serialize(streamEntry.Fields)
	if err != nil {
		return "", err
	}
	return "*2" + utils.CRLF + serId + serFeilds, nil
}

func (sId *StreamID) ToString() string {
	return strconv.FormatUint(sId.Milliseconds, 10) + "-" + strconv.FormatUint(sId.SequenceNum, 10)
}
func GetStreamID(input string) (StreamID, error) {
	splitInput := strings.Split(input, "-")
	output := StreamID{}
	milli, err := strconv.ParseUint(splitInput[0], 10, 64)
	if err != nil {
		return StreamID{}, err
	}
	output.Milliseconds = milli
	if len(splitInput) == 2 {
		seq, err := strconv.ParseUint(splitInput[1], 10, 64)
		if err != nil {
			return StreamID{}, err
		}
		output.SequenceNum = seq
	} else {
		output.SequenceNum = 0
	}
	return output, nil
}

func (s *Stream) ValidateId(streamId StreamID) error {
	lastId := s.LastID
	if streamId.Milliseconds == 0 && streamId.SequenceNum == 0 {
		return errors.New("ERR The ID specified in XADD must be greater than 0-0")
	} else if streamId.Milliseconds < lastId.Milliseconds {
		return errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	} else if streamId.Milliseconds == lastId.Milliseconds && streamId.SequenceNum <= lastId.SequenceNum {
		return errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	} else {
		return nil
	}

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
		return
	}
	waitersList = append(waitersList[:idx], waitersList[idx+1:]...)
	s.waiters[listName] = waitersList
	s.mu.Unlock()
}
func (s *Stream) getEntries(startStreamId, stopStreamId StreamID) ([]StreamEntry, error) {

	if (startStreamId.Milliseconds > stopStreamId.Milliseconds) || (startStreamId.Milliseconds == stopStreamId.Milliseconds && startStreamId.SequenceNum > stopStreamId.SequenceNum) {
		return nil, errors.New("Invalid start and stop ids")
	}

	output := []StreamEntry{}
	for _, entry := range s.Entries {
		if entry.ID.Milliseconds < startStreamId.Milliseconds {
			continue
		}
		if entry.ID.Milliseconds == startStreamId.Milliseconds && (entry.ID.SequenceNum < startStreamId.SequenceNum) {
			continue
		}
		if entry.ID.Milliseconds <= stopStreamId.Milliseconds {
			if entry.ID.Milliseconds == stopStreamId.Milliseconds && entry.ID.SequenceNum > stopStreamId.SequenceNum {
				continue
			}
			output = append(output, entry)
		} else {
			break
		}
	}
	return output, nil

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
		return
	}
	waitersList = append(waitersList[:idx], waitersList[idx+1:]...)
	s.waiters[listName] = waitersList
}

func (s *Stream) AddStreamEntry(streamEntry StreamEntry) (StreamID, error) {
	err := s.ValidateId(streamEntry.ID)
	if err != nil {
		return StreamID{}, err
	}
	s.Entries = append(s.Entries, streamEntry)
	s.LastID = s.Entries[len(s.Entries)-1].ID

	return s.LastID, nil

}
