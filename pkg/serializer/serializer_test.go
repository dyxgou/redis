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
		{"SET key :123", "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n:123\r\n"},
		{"SET key 123.123", "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n,123.123\r\n"},
		{"SET key 123123123123123", "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n(123123123123123\r\n"},
		{"SET key value EX 20", "*5\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n$2\r\nEX\r\n:20\r\n"},
		{"SET key value EX 20 NX", "*6\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n$2\r\nEX\r\n:20\r\n$2\r\nNX\r\n"},
		{"SET key value EX 20 XX", "*6\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n$2\r\nEX\r\n:20\r\n$2\r\nXX\r\n"},
		{"GETSET key value", "*3\r\n$6\r\nGETSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"},
		{"GETEX key value", "*3\r\n$5\r\nGETEX\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"},
		{"GETDEL key value", "*3\r\n$6\r\nGETDEL\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"},
		{"INCR key", "*2\r\n$4\r\nINCR\r\n$3\r\nkey\r\n"},
		{"INCRBY key 2", "*3\r\n$6\r\nINCRBY\r\n$3\r\nkey\r\n:2\r\n"},
		{"DECR key", "*2\r\n$4\r\nDECR\r\n$3\r\nkey\r\n"},
		{"DECRBY key 2", "*3\r\n$6\r\nDECRBY\r\n$3\r\nkey\r\n:2\r\n"},
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
