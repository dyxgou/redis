package timer

import (
	"os"
	"slices"
	"testing"
)

var timer *Timer

func TestMain(m *testing.M) {
	timer = New(10_000_000)

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
			key:  "key5",
			time: 5,
		},
		{
			key:  "key4",
			time: 4,
		},
		{
			key:  "key3",
			time: 3,
		},
		{
			key:  "key2",
			time: 2,
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

	for _, tt := range slices.Backward(tests) {
		v := timer.Remove()

		if tt.key != v.key {
			t.Errorf("key expected=%q. got=%q", tt.key, v.key)
		}

		if tt.time != v.time {
			t.Errorf("time expected=%d. got=%d", tt.time, v.time)
		}
	}
}

func FuzzTimer(f *testing.F) {
	f.Skip("Skipped just to check performace over time")
	f.Add("key1", int64(1))

	f.Fuzz(func(t *testing.T, key string, ti int64) {
		ts := NewTimestamp(key, ti)
		timer.Insert(ts)
	})
}

func TestTimerExisted(t *testing.T) {
	tests := []struct {
		key  string
		time int64
	}{
		{
			key:  "key1",
			time: -1,
		},
		{
			key:  "key1",
			time: -2,
		},
		{
			key:  "key3",
			time: -4,
		},
		{
			key:  "key5",
			time: -6,
		},
		{
			key:  "key7",
			time: -8,
		},
		{
			key:  "key9",
			time: -10,
		},
	}

	for _, val := range tests {
		timer.Insert(NewTimestamp(val.key, val.time))
	}

	t.Logf("timer=%+v", timer.ts[:timer.N])

	for ti := range timer.Exited(10) {
		t.Logf("time=%+v", ti)
	}
}
