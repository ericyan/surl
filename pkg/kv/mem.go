package kv

import (
	"errors"
	"sync"
)

var (
	errInvalidKey = errors.New("key is invalid")
	errNotFound   = errors.New("key not found")
)

type InMemoryStore struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewInMemoryStore() (*InMemoryStore, error) {
	return &InMemoryStore{data: make(map[string][]byte, 0)}, nil
}

func (s *InMemoryStore) Get(key []byte) ([]byte, error) {
	s.mu.RLock()
	value, ok := s.data[string(key)]
	s.mu.RUnlock()

	if !ok {
		return nil, errNotFound
	}

	return value, nil
}

func (s *InMemoryStore) Put(key, value []byte) error {
	if len(key) == 0 {
		return errInvalidKey
	}

	s.mu.Lock()
	s.data[string(key)] = value
	s.mu.Unlock()

	return nil
}
