package serializer

import "testing"

func TestSerialize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"GET", "*1\r\n$3\r\nGET\r\n"},
		{"GET key", "*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n"},
		{"SET key value", "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"},
		{"SET key value EX 20", "*5\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n$2\r\nEX\r\n:20\r\n"},
	}

	for _, tt := range tests {
		s := New(tt.input)
		serialized, err := s.Serialize()

		if err != nil {
			t.Error(err)
		}

		if serialized != tt.expected {
			t.Errorf("serialized expected=%q. got=%q", tt.expected, serialized)
		}
	}
}
