// Package kv provides key-value store implementations.
package kv

// A Store represents a key-value store that treats both the key and the
// value an opaque slice of bytes.
type Store interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
}
