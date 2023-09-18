package main

import (
	"testing"
)

func TestDeserializeRESP(t *testing.T) {
	tests := []struct {
		input    []byte
		expected interface{}
	}{
		// add your test cases
		{[]byte("+OK\r\n"), "OK"},
		{[]byte("-Error message\r\n"), "Error message"},
		// ... add more tests
	}

	for _, test := range tests {
		result, err := DeserializeRESP(test.input)
		if err != nil {
			t.Errorf("Failed to deserialize: %v", err)
			continue
		}

		if result != test.expected {
			t.Errorf("Expected %v\n, but got %v", test.expected, result)
		}
	}
}


func TestSerializeRESP(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected []byte
	}{
		{"hello", []byte("$5\r\nhello\r\n")},
		// Add other test cases...
	}

	for _, test := range tests {
		result, err := SerializeRESP(test.input)
		if err != nil {
			t.Errorf("Error serializing %v: %s", test.input, err)
		}
		if string(result) != string(test.expected) {
			t.Errorf("Expected %s, but got %s", string(test.expected), string(result))
		}
	}
}

