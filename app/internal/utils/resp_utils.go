package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

func SerializeInt(num int) string {
	return ":" + strconv.Itoa(num) + CRLF
}
func DeserializeInt(str string) (int, error) {
	if str[0] != ':' || len(str) < 4 {
		return 0, errors.New("Invalid resp int")
	}
	return strconv.Atoi(str[1 : len(str)-2])
}
func SerializeString(str string) (string, error) {
	if len(str) == 0 {
		return "", errors.New("null string passed")
	}
	return "+" + str + CRLF, nil
}

// todo: figure out if really needed
func DeserializeString(str string) string {
	return ""
}

func BasicError(err string) string {
	return "-Error " + err + CRLF
}

func SerializeBulkString(word string) (string, error) {
	if len(word) == 0 {
		return "", errors.New("null string passed")
	}
	word = strings.TrimLeftFunc(word, unicode.IsSpace)
	word = strings.TrimRightFunc(word, unicode.IsSpace)

	return "$" + strconv.Itoa(len(word)) + CRLF + word + CRLF, nil
}

func DeserializeBulkString(str string) (string, error) {

	if str[0] != '$' || len(str) == 0 {
		return "", errors.New("invalid bulk string")
	}
	parts := strings.SplitN(str, CRLF, 2)
	if len(parts) != 2 {
		return "", errors.New("invalid bulk string")

	}
	stringLen, err := strconv.Atoi(parts[0][1:])
	if err != nil {
		return "", errors.New("invalid bulk string")
	}

	if stringLen == -1 {
		return "", nil
	}

	if stringLen < -1 || len(parts[1]) < stringLen+2 {
		return "", errors.New("invalid bulk string")
	}

	if parts[1][stringLen:stringLen+2] != CRLF {
		return "", errors.New("missing bulk string terminator")
	}
	return parts[1][:stringLen], nil
}

func SerializeArray(arr []string) (string, error) {
	if len(arr) == 0 {
		return "*0\r\n", nil
	}

	if arr == nil {
		return "*-1\r\n", nil
	}

	output := "*"
	output += strconv.Itoa(len(arr)) + CRLF

	for _, str := range arr {
		bulkStr, err := SerializeBulkString(str)
		if err != nil {
			return "", err
		}
		output += bulkStr
	}
	return output, nil
}

// *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n

func DeserializeArray(reader *bufio.Reader) ([]string, error) {
	header, err := readRESPLine(reader)
	if err != nil {
		return nil, err
	}
	if len(header) == 0 || header[0] != '*' {
		return nil, errors.New("expected RESP array")
	}

	count, err := strconv.Atoi(header[1:])
	if err != nil {
		return nil, fmt.Errorf("invalid array length %q: %w", header[1:], err)
	}
	if count == -1 {
		return nil, nil // Null array
	}
	if count < 0 {
		return nil, fmt.Errorf("invalid array length %d", count)
	}

	result := make([]string, 0, count)
	for i := 0; i < count; i++ {
		header, err := readRESPLine(reader)
		if err != nil {
			return nil, err
		}
		if len(header) == 0 || header[0] != '$' {
			return nil, fmt.Errorf("expected bulk string at array item %d", i)
		}

		length, err := strconv.Atoi(header[1:])
		if err != nil {
			return nil, fmt.Errorf("invalid bulk string length %q: %w", header[1:], err)
		}
		if length == -1 {
			result = append(result, "") // Null bulk strings are represented as empty strings for now.
			continue
		}
		if length < 0 || length > int(^uint(0)>>1)-2 {
			return nil, fmt.Errorf("invalid bulk string length %d", length)
		}

		payload := make([]byte, length+2)
		if _, err := io.ReadFull(reader, payload); err != nil {
			return nil, fmt.Errorf("read bulk string payload: %w", err)
		}
		if string(payload[length:]) != CRLF {
			return nil, errors.New("missing bulk string terminator")
		}

		result = append(result, string(payload[:length]))
	}

	return result, nil
}

func readRESPLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	if !strings.HasSuffix(line, CRLF) {
		return "", errors.New("RESP line is missing CRLF terminator")
	}
	return strings.TrimSuffix(line, CRLF), nil
}
