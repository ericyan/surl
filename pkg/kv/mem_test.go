package kv

import (
	"bytes"
	"testing"
)

func TestInMemoryStore(t *testing.T) {
	store, err := NewInMemoryStore()

	var (
		empty = []byte{}
		key   = []byte{42}
		value = []byte("hello, world")
	)

	err = store.Put(empty, value)
	if err == nil {
		t.Error("empty key should not be accepted")
	}

	err = store.Put(key, value)
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}
	result, err := store.Get(key)
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}
	if !bytes.Equal(result, value) {
		t.Fatalf("unexpected value: got %v, want %s\n", result, value)
	}

	err = store.Put(key, empty)
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}
	result, err = store.Get(key)
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}
	if !bytes.Equal(result, empty) {
		t.Fatalf("unexpected value: got %v, want %s\n", result, value)
	}
}
