package shortener

import (
	"regexp"
	"testing"
)

func TestShorten(t *testing.T) {
	cases := []string{
		"https://example.com/",
		"https://github.com/ericyan?tab=stars",
		"https://www.google.com/search?source=hp&ei=2ntYW6CBAo6g8AO08qjoAg&q=hello%2C+world",
	}

	s := New()
	for _, in := range cases {
		out := s.Shorten(in)

		ok, err := regexp.MatchString("^[a-zA-Z0-9]{8}$", out)
		if err != nil {
			t.Errorf("regexp match error: %s", err)
		}

		if !ok {
			t.Errorf("regexp test failed: %s -> %s", in, out)
		}
	}
}
