package kv

import (
	"time"

	"github.com/allegro/bigcache"
)

// A CachedStore is a key-value store embedded with an in-memory cache.
type CachedStore struct {
	cache *bigcache.BigCache
	store Store
}

// NewCachedStore returns a CachedStore for given Store.
func NewCachedStore(store Store) (*CachedStore, error) {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		return nil, err
	}

	return &CachedStore{cache, store}, nil
}

func (s *CachedStore) Get(key []byte) ([]byte, error) {
	value, err := s.cache.Get(string(key))
	if err == nil {
		return value, nil
	}

	return s.store.Get(key)
}

func (s *CachedStore) Put(key, value []byte) error {
	_ = s.cache.Set(string(key), value)

	return s.store.Put(key, value)
}
