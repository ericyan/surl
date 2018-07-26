package kv

import (
	"bytes"
	"testing"
)

func TestCachedStore(t *testing.T) {
	var (
		key     = []byte{42}
		value   = []byte("hello, world")
		updated = []byte("bazinga!")
	)

	memstore, err := NewInMemoryStore()
	err = memstore.Put(key, value)
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}

	cached, err := NewCachedStore(memstore)
	result, err := cached.Get(key)
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}
	if !bytes.Equal(result, value) {
		t.Fatalf("unexpected value: got %s, want %s\n", result, value)
	}

	err = cached.Put(key, updated)
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}
	result, err = memstore.Get(key)
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}
	if !bytes.Equal(result, updated) {
		t.Fatalf("unexpected value: got %s, want %s\n", result, value)
	}
}
