package error

import "errors"

var (
	ErrWrongType = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
)
