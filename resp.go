package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

// DeserializeRESP takes a RESP string and returns the parsed message.
func DeserializeRESP(input []byte) (interface{}, error) {
	reader := bufio.NewReader(bytes.NewReader(input))

	// Check the data type
	dataType, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch dataType {
	case '+':
		data, err := reader.ReadString('\n')
		return data[:len(data)-2], err  // remove "\r\n"
	case '-':
		data, err := reader.ReadString('\n')
		return data[:len(data)-2], err  // remove "\r\n"
	case ':':
		data, err := reader.ReadString('\n')
		return data[:len(data)-2], err  // remove "\r\n"
	case '$':
		lengthStr, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		length, err := strconv.Atoi(lengthStr[:len(lengthStr)-2])  // remove "\r\n"
		if err != nil {
			return nil, err
		}
		if length == -1 {
			return nil, nil  // return nil for "$-1\r\n"
		}
		data := make([]byte, length)
		_, err = reader.Read(data)
		if err != nil {
			return nil, err
		}
		reader.ReadByte()  // read the trailing "\r"
		reader.ReadByte()  // read the trailing "\n"
		return string(data), err
	case '*':
		// Handle arrays (further implementation needed)
		return nil, errors.New("array deserialization not implemented yet")
	default:
		return nil, errors.New("invalid RESP type")
	}
}


// SerializeRESP converts a native Go type to a RESP string.
// We're implementing a basic version. It can be extended based on your needs.
func SerializeRESP(input interface{}) ([]byte, error) {
	switch v := input.(type) {
	case string:
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v)), nil
	case error:
		return []byte(fmt.Sprintf("-%s\r\n", v.Error())), nil
	default:
		return nil, errors.New("type not supported for serialization")
	}
}

