package storage

import (
	"sync"
	"testing"
)

func TestMultipleWrite(t *testing.T) {
	s := New()

	tests := []struct {
		key string
		val any
	}{
		{key: "key1", val: 1},
		{key: "key2", val: "string"},
		{key: "key3", val: true},
		{key: "key4", val: false},
		{key: "key5", val: false},
		{key: "key6", val: int64(123)},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)

		go func(k string, v any) {
			defer wg.Done()
			if err := s.Set(k, v); err != nil {
				t.Error(err)
			}
		}(tt.key, tt.val)
	}

	for _, tt := range tests {
		v, ok := s.Get(tt.key)
		if !ok {
			t.Errorf("value of key=%q not exists", tt.key)
		}

		assertValue(t, v, tt.val)
	}
}

func assertValue(t *testing.T, v any, expected any) {
	switch v := v.(type) {
	case *Int:
		i, ok := expected.(int)
		if !ok {
			t.Error("expected is not an int type")
		}
		assertInt(t, v, i)
	case *Int64:
		i, ok := expected.(int64)
		if !ok {
			t.Error("expected is not an int64 type")
		}

		assertInt64(t, v, i)
	case *Float:
		f, ok := expected.(float64)
		if !ok {
			t.Error("expected is not an float type")
		}

		assertFloat(t, v, f)
	case *Bool:
		b, ok := expected.(bool)
		if !ok {
			t.Error("expected is not an int type")
		}

		assertBool(t, v, b)
	case *String:
		s, ok := expected.(string)
		if !ok {
			t.Error("expected is not an int type")
		}

		assertString(t, v, s)
	default:
		t.Errorf("value not supported. got=%+v", v)
	}
}
