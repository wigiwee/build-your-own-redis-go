package storage

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal/utils"
)

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

	if !exists {

		entry := Entry{
			Type: utils.TYPE_STREAM,
			Stream: &Stream{
				Entries: []StreamEntry{
					StreamEntry{
						ID:     streamId,
						Fields: args,
					},
				},
				LastID: streamId,
			},
		}
		s.data[streamkey] = entry
	} else {
		_, err := s.data[streamkey].Stream.AddStreamEntry(
			StreamEntry{
				ID:     streamId,
				Fields: args,
			},
		)
		if err != nil {
			return "", err
		}
	}
	return s.data[streamkey].Stream.Entries[len(s.data[streamkey].Stream.Entries)-1].ID.ToString(), nil

}

func (s *Store) ReadStream(streamKey string, streamId string) ([]StreamEntry, error) {
	// s.mu.RLock()
	// defer s.mu.RUnlock()
	// startid, err := GetStreamID(streamId)
	// if err != nil {
	// 	return []StreamEntry{}, err
	// }
	// streamEntries, err := s.RangeStream(streamKey, startid.ToString(), "+")
	// if err != nil {
	// 	return []StreamEntry{}, err
	// }
	return nil, nil
}

func (s *Store) RangeStream(streamKey string, start string, stop string) ([]StreamEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.data[streamKey]
	if !exists {
		return nil, errors.New("stream with given key doesn't exist")
	}

	output := []StreamEntry{}

	startSplit := strings.Split(start, "-")
	stopSplit := strings.Split(stop, "-")

	startStreamId := &StreamID{}
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

	stopStreamId := &StreamID{}
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
	//validate start and stop
	if (startStreamId.Milliseconds > stopStreamId.Milliseconds) || (startStreamId.Milliseconds == stopStreamId.Milliseconds && startStreamId.SequenceNum > stopStreamId.SequenceNum) {
		return nil, errors.New("Invalid start and stop ids")
	}

	for _, entry := range entry.Stream.Entries {
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
