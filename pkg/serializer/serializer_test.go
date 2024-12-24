package serializer

import (
	"github/dyxgou/redis/pkg/lexer"
	"testing"
)

func TestSerialize(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		expectedLen int
	}{
		{"GET", "*1\r\n$3\r\nGET\r\n", 1},
		{"GET mykey", "*2\r\n$3\r\nGET\r\n$5\r\nmykey\r\n", 2},
		{"SET mykey myvalue", "*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$7\r\nmyvalue\r\n", 3},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		s := New(l)

		sr, err := s.Serialize()

		if err != nil {
			t.Error(err)
		}

		if s.w.len != tt.expectedLen {
			t.Errorf("amount of tokens expected=%d. got=%d", tt.expectedLen, s.w.len)
		}

		if sr != tt.expected {
			t.Errorf("expeted %q. got=%q", tt.expected, sr)
		}
	}
}
