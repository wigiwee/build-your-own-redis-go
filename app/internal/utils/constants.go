package utils

import "errors"

const CRLF = "\r\n"

const NULL_BULK_STR = "$-1\r\n"

const OK_BULK_STR = "+OK\r\n"

const TYPE_STRING = "string"

const TYPE_INT = "int"

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

const EMPTY_RDB_HEX = "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2"
