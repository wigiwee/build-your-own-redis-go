package storage

import (
	"errors"
	"fmt"
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
	SequenceNum  uint64
}

type StreamEntry struct {
	ID     StreamID
	Fields map[string]string
}

func (sId *StreamID) ToString() string {
	return strconv.FormatUint(sId.Milliseconds, 10) + "-" + strconv.FormatUint(sId.SequenceNum, 10)
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

// func GenerateStreamID(exists bool,streamKey, input string) (StreamID, error) {

// 	if input == "*" {
// 		if exists{

// 		}
// 	} else {

// 		splitId := strings.Split(input, "-")
// 		if splitId[1] == "*" {

// 		} else {
// 			millSec, err := strconv.ParseInt(splitId[0], 10, 64)
// 			if err != nil {
// 				return StreamID{}, err
// 			}
// 			seqNum, err := strconv.ParseInt(splitId[1], 10, 64)
// 			if err != nil {
// 				return StreamID{}, err
// 			}
// 			return StreamID{
// 				Milliseconds: uint64(millSec),
// 				SequenceNum:  uint64(seqNum),
// 			}, nil

// 		}
// 	}
// }

func (s *Stream) AddStreamEntry(streamEntry StreamEntry) (StreamID, error) {
	err := s.ValidateId(streamEntry.ID)
	if err != nil {
		return StreamID{}, err
	}
	s.Entries = append(s.Entries, streamEntry)
	fmt.Println("from the func", s.Entries)
	s.LastID = s.Entries[len(s.Entries)-1].ID

	return s.LastID, nil

}
