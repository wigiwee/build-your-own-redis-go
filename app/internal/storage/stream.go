package storage

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

type StreamReadResult struct {
	Key     string
	Entries []StreamEntry
}

func (s StreamReadResult) SerializeRESP() (string, error) {
	serStr, err := utils.Serialize([]any{s.Key, s.Entries})
	if err != nil {
		return "", err
	}
	return serStr, nil
}

func (s *Store) AddStream(streamkey string, id string, args []string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(args)%2 != 0 {
		return "", errors.New("invalid fields")
	}
	_, exists := s.data[streamkey]
	streamId, err := s.GenerateStreamID(streamkey, id)
	if err != nil {
		return "", err
	}
	if streamId.Milliseconds == 0 && streamId.SequenceNum == 0 {
		return "", errors.New("ERR The ID specified in XADD must be greater than 0-0")
	}
	newEntry := StreamEntry{
		ID:     streamId,
		Fields: args,
	}
	if !exists {

		entry := Entry{
			Type: utils.TYPE_STREAM,
			Stream: &Stream{
				Entries: []StreamEntry{newEntry},
				LastID:  streamId,
				Waiters: make(map[StreamID][]*StreamWaiters),
			},
		}
		s.data[streamkey] = entry
	} else {
		_, err := s.data[streamkey].Stream.AddStreamEntry(newEntry)
		if err != nil {
			return "", err
		}
	}
	entry := s.data[streamkey].Stream
	for id, waiters := range entry.Waiters {
		if streamId.IsGreaterThan(id) {
			for _, waiter := range waiters {
				waiter.ch <- newEntry
			}
		}
	}

	return s.data[streamkey].Stream.Entries[len(s.data[streamkey].Stream.Entries)-1].ID.ToString(), nil
}

func (s *Store) ReadStreamBlock(streamKey string, streamIdStr string, ttl time.Duration) (StreamReadResult, error) {
	s.mu.Lock()
	fmt.Println("at the raed stream blocks")

	output := StreamReadResult{}
	output.Key = streamKey

	streamId, err := GetStreamID(streamIdStr)
	if err != nil {
		return output, err
	}

	entry, exists := s.data[streamKey]
	fmt.Println("entry doesn't exist")
	if !exists {
		entry = Entry{
			Type: utils.TYPE_STREAM,
			Stream: &Stream{
				Waiters: make(map[StreamID][]*StreamWaiters, 1),
				Entries: []StreamEntry{},
			},
		}
		s.data[streamKey] = entry
	}

	//checking if already present
	fmt.Println(entry.Stream.Entries)
	for _, streamEntry := range entry.Stream.Entries {
		if streamEntry.ID.IsGreaterThan(streamId) {
			output.Entries = append(output.Entries, streamEntry)
		}
	}
	if len(output.Entries) != 0 {
		return output, nil
	}

	timer := time.NewTimer(ttl)
	defer timer.Stop()

	waiter := &StreamWaiters{
		ch:       make(chan StreamEntry),
		notified: false,
	}

	entry.Stream.Waiters[streamId] = append(entry.Stream.Waiters[streamId], waiter)
	s.mu.Unlock()
	fmt.Println(ttl)
	fmt.Println("here to wait the timer out")
	var null StreamReadResult
	select {
	case result := <-waiter.ch:
		output.Entries = []StreamEntry{result}
		s.mu.Lock()
		delete(s.data[streamKey].Stream.Waiters, streamId)
		s.mu.Unlock()
		fmt.Println("return this output, ", output)
		return output, nil

	case <-timer.C:
		if waiter.notified {
			result := <-waiter.ch
			output.Entries = []StreamEntry{result}
			s.mu.Lock()
			delete(s.data[streamKey].Stream.Waiters, streamId)
			s.mu.Unlock()
			return null, nil
		}
		s.mu.Lock()
		delete(s.data[streamKey].Stream.Waiters, streamId)
		s.mu.Unlock()
		return null, nil

	}

}

func (s *Store) ReadStream(streamKey string, streamId string) (StreamReadResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	output := StreamReadResult{
		Key: streamKey,
	}

	startid, err := GetStreamID(streamId)
	if err != nil {
		return output, err
	}
	startid.SequenceNum += 1

	entry, exists := s.data[streamKey]
	if !exists {
		return output, nil
	}
	streamEntries, err := entry.Stream.getEntries(startid, entry.Stream.LastID)
	if err != nil {
		return output, err
	}
	output.Entries = streamEntries
	return output, nil
}

func (s *Store) RangeStream(streamKey string, start string, stop string) ([]StreamEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.data[streamKey]
	if !exists {
		return nil, errors.New("stream with given key doesn't exist")
	}

	startSplit := strings.Split(start, "-")
	stopSplit := strings.Split(stop, "-")

	startStreamId := StreamID{}
	if start == "-" {
		startStreamId.Milliseconds = 0
		startStreamId.SequenceNum = 0
	} else {
		millisec, err := strconv.ParseUint(startSplit[0], 10, 64)
		if err != nil {
			return nil, err
		}
		startStreamId.Milliseconds = millisec
		if len(startSplit) == 2 {
			seq, err := strconv.ParseUint(startSplit[1], 10, 64)
			if err != nil {
				return nil, err
			}
			startStreamId.SequenceNum = seq
		} else {
			startStreamId.SequenceNum = 0
		}
	}

	stopStreamId := StreamID{}
	if stop == "+" {
		stopStreamId.Milliseconds = entry.Stream.LastID.Milliseconds
		stopStreamId.SequenceNum = entry.Stream.LastID.SequenceNum
	} else {

		millisec, err := strconv.ParseUint(stopSplit[0], 10, 64)
		if err != nil {
			return nil, err
		}
		stopStreamId.Milliseconds = millisec
		if len(stopSplit) == 2 {
			seq, err := strconv.ParseUint(stopSplit[1], 10, 64)
			if err != nil {
				return nil, err
			}
			stopStreamId.SequenceNum = seq
		} else {
			stopStreamId.SequenceNum = math.MaxUint64
		}
	}
	return entry.Stream.getEntries(startStreamId, stopStreamId)
}

func (s *Store) GenerateStreamID(streamKey, input string) (StreamID, error) {
	streamId := StreamID{
		Milliseconds: uint64(time.Now().UnixMilli()),
	}
	_, exists := s.data[streamKey]
	if input == "*" {
		if !exists {
			streamId.SequenceNum = 0

		} else if s.data[streamKey].Stream.LastID.Milliseconds == streamId.Milliseconds {
			streamId.SequenceNum = s.data[streamKey].Stream.LastID.SequenceNum + 1
		} else {
			streamId.SequenceNum = 0
		}
		return streamId, nil
	} else {

		splitId := strings.Split(input, "-")

		millSec, err := strconv.ParseUint(splitId[0], 10, 64)
		if err != nil {
			return StreamID{}, err
		}

		streamId.Milliseconds = uint64(millSec)

		if splitId[1] == "*" {
			if !exists {
				if streamId.Milliseconds == 0 {
					streamId.SequenceNum = 1
				} else {
					streamId.SequenceNum = 0
				}
				return streamId, nil
			}
			if s.data[streamKey].Stream.LastID.Milliseconds == streamId.Milliseconds {
				streamId.SequenceNum = s.data[streamKey].Stream.LastID.SequenceNum + 1
			} else {
				streamId.SequenceNum = 0
			}
			return streamId, nil
		} else {
			seqNum, err := strconv.ParseUint(splitId[1], 10, 64)
			if err != nil {
				return StreamID{}, err
			}
			return StreamID{
				Milliseconds: uint64(millSec),
				SequenceNum:  uint64(seqNum),
			}, nil

		}
	}
}
