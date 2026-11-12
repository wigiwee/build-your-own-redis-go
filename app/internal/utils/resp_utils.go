package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

type RESPSerializer interface {
	SerializeRESP() (string, error)
}

func Serialize(value any) (string, error) {
	if value == nil {
		return "$-1" + CRLF, nil
	}

	if respSerializer, ok := value.(RESPSerializer); ok {
		return respSerializer.SerializeRESP()
	}
	var errorType = reflect.TypeOf((*error)(nil)).Elem()

	rv := reflect.ValueOf(value)

	if rv.Type().Implements(errorType) {
		return "-" + value.(error).Error() + CRLF, nil
	}

	switch rv.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return ":" + strconv.FormatInt(int64(rv.Int()), 10) + CRLF, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		return ":" + strconv.FormatUint(uint64(rv.Uint()), 10) + CRLF, nil

	case reflect.String:
		return "$" + strconv.Itoa(len(rv.String())) + CRLF + rv.String() + CRLF, nil

	case reflect.Array, reflect.Slice:
		if rv.Kind() == reflect.Slice && rv.IsNil() {
			return "*-1" + CRLF, nil
		}
		var output strings.Builder
		output.WriteString("*")
		output.WriteString(strconv.Itoa(rv.Len()))
		output.WriteString(CRLF)

		for i := 0; i < rv.Len(); i++ {
			serStr, err := Serialize(rv.Index(i).Interface())
			if err != nil {
				return "", err
			}
			output.WriteString(serStr)
		}

		return output.String(), nil
	default:
		return "", errors.New("unsupported type")
	}
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

func DeserializeBulkString(str string) (string, error) {
	if len(str) == 0 {
		return "", errors.New("invalid input")
	}
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
