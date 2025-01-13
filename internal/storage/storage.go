package storage

import (
	"fmt"
	"github/dyxgou/redis/pkg/ast"
	"sync"
)

type Storage struct {
	mu      sync.RWMutex
	storage map[string]Value
}

func New() *Storage {
	return &Storage{
		storage: make(map[string]Value),
	}
}

func (s *Storage) Set(k string, v ast.Expression) error {
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

func (s *Storage) Delete(k string) {
	s.mu.Lock()
	delete(s.storage, k)
	defer s.mu.Unlock()
}

func (s *Storage) Exists(k string) bool {
	s.mu.RLock()
	_, ok := s.storage[k]
	defer s.mu.RUnlock()
	return ok
}

func (s *Storage) Get(k string) (Value, bool) {
	s.mu.RLock()
	val, ok := s.storage[k]
	s.mu.RUnlock()
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
