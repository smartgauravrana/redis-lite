package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
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
	case '+': // simple string
		data, err := reader.ReadString('\n')
		return data[:len(data)-2], err  // remove "\r\n"
	case '-': // errors
		data, err := reader.ReadString('\n')
		return data[:len(data)-2], err  // remove "\r\n"
	case ':': // integer
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
		// Handle arrays (using DeserializeRESPArray)
		return DeserializeRESPArray(input)
	default:
		return nil, errors.New("invalid RESP type")
	}
}

// DeserializeRESPArray deserializes a RESP array.
func DeserializeRESPArray(input []byte) ([]interface{}, error) {
	reader := bufio.NewReader(bytes.NewReader(input))

	// Check for the asterisk (*) which indicates the start of an array
	firstChar, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if firstChar != '*' {
		return nil, errors.New("not an RESP array")
	}

	// Read the number of elements in the array
	arrayLengthStr, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	arrayLengthStr = strings.Trim(arrayLengthStr, "\r\n")
	arrayLength, err := strconv.Atoi(arrayLengthStr)
	if err != nil {
		return nil, err
	}

	// Initialize the array to hold the elements
	result := make([]interface{}, arrayLength)

	// Read each element of the array
	for i := 0; i < arrayLength; i++ {
		// Read the first byte of the element to determine its type
		elementType, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}

		switch elementType {
		case '+':
			// Simple string
			data, err := reader.ReadString('\n')
			if err != nil {
				return nil, err
			}
			result[i] = strings.Trim(data, "\r\n")
		case '-':
			// Error
			data, err := reader.ReadString('\n')
			if err != nil {
				return nil, err
			}
			result[i] = errors.New(strings.Trim(data, "\r\n"))
		case ':':
			// Integer
			data, err := reader.ReadString('\n')
			if err != nil {
				return nil, err
			}
			intValue, err := strconv.Atoi(strings.Trim(data, "\r\n"))
			if err != nil {
				return nil, err
			}
			result[i] = intValue
		case '$':
			// Bulk string or null bulk string
			lengthStr, err := reader.ReadString('\n')
			if err != nil {
				return nil, err
			}
			lengthStr = strings.Trim(lengthStr, "\r\n")
			length, err := strconv.Atoi(lengthStr)
			if err != nil {
				return nil, err
			}
			if length == -1 {
				result[i] = nil // Null bulk string
			} else {
				data := make([]byte, length)
				_, err := reader.Read(data)
				if err != nil {
					return nil, err
				}
				// Read the trailing CRLF
				reader.Discard(2)
				result[i] = string(data)
			}
		case '*':
			// Nested array (recursively deserialize)
			nestedArray, err := DeserializeRESPArray(input[1:])
			if err != nil {
				return nil, err
			}
			result[i] = nestedArray
		default:
			return nil, errors.New("unknown RESP type")
		}
	}

	return result, nil
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

