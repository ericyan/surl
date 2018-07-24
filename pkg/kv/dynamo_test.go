package kv

import (
	"bytes"
	"os"
	"testing"
)

func TestDynamoStore(t *testing.T) {
	table := os.Getenv("SURL_DYNAMODB_TABLE")
	if table == "" {
		t.Skipf("Skipping %s SURL_DYNAMODB_TABLE as is not set", t.Name())
	}

	store, err := NewDynamoStore(os.Getenv("SURL_DYNAMODB_ENDPOINT"), table)
	if err != nil {
		t.Fatal(err)
	}

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
}
