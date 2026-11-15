package utils

import "errors"

const CRLF = "\r\n"

const NULL_BULK_STR = "$-1\r\n"

const OK_BULK_STR = "+OK\r\n"

const TYPE_STRING = "string"

const TYPE_LIST = "list"

const TYPE_SET = "set"

const TYPE_ZSET = "zset"

const TYPE_HASH = "has"

const TYPE_STREAM = "stream"

const TYPE_VECTORSET = "vectorset"

var (
	ErrWrongType = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")

	ErrInvalidCmd = errors.New("syntax error")
)
