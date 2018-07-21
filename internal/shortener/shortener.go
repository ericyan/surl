// Packages shortener implements hash-based URL shortening.
package shortener

import (
	"hash"
	"hash/fnv"

	b58 "github.com/jbenet/go-base58"
)

type Shortener struct {
	hash.Hash64
}

// New returns a new Shortener with 64-bit FNV-1a hashing.
func New() *Shortener {
	return &Shortener{fnv.New64a()}
}

// Shorten returns the first 8 characters of the Base58-encoded hash of
// url. It always returns a result, even the input is not a valid URL.
func (s *Shortener) Shorten(url string) string {
	s.Reset()
	s.Write([]byte(url))
	return b58.Encode(s.Sum(nil))[:8]
}
