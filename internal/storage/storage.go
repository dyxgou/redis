package storage

import (
	"fmt"
	"sync"
)

type Storage struct {
	mu      sync.RWMutex
	storage map[string]val
}

func New() *Storage {
	return &Storage{
		storage: make(map[string]val),
	}
}

func (s *Storage) Set(k string, v any) error {
	val := NewValue(v)
	if val == nil {
		return NotSupportedValErr
	}

	s.mu.Lock()
	if !s.isSameKind(k, val.kind()) {
		return fmt.Errorf("value kind of key=%q mismatched", val.kind())
	}

	s.storage[k] = val
	defer s.mu.Unlock()

	return nil
}

func (s *Storage) Get(k string) (val, bool) {
	s.mu.RLock()
	val, ok := s.storage[k]
	s.mu.RLock()
	if !ok {
		return nil, ok
	}

	return val, ok
}

func (s *Storage) isSameKind(k string, kind valueKind) bool {
	v, ok := s.storage[k]
	if !ok {
		return true
	}

	return v.kind() == kind
}
