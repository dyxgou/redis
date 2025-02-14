package timer

import (
	"os"
	"testing"
)

var timer *Timer

func TestMain(m *testing.M) {
	timer = New()

	code := m.Run()

	os.Exit(code)
}

func TestTimerInsert(t *testing.T) {
	tests := []struct {
		key  string
		time int64
	}{
		{
			key:  "key6",
			time: 6,
		},
		{
			key:  "key3",
			time: 3,
		},
		{
			key:  "key4",
			time: 4,
		},
		{
			key:  "key2",
			time: 2,
		},
		{
			key:  "key5",
			time: 5,
		},
		{
			key:  "key1",
			time: 1,
		},
	}

	for _, tt := range tests {
		ts := NewTimestamp(tt.key, tt.time)
		timer.Insert(ts)
	}

	t.Logf("timer=%+v", timer)
}
