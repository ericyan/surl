package handler

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	handler := New()

	r := httptest.NewRequest("POST", "/submit", bytes.NewBufferString(`{"url": "https://www.example.com/"}`))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	if string(body) != `{"url":"https://www.example.com/","shorten_url":"M9Yv6VB2"}` {
		t.Errorf("unexpected response: %s", body)
	}

	r = httptest.NewRequest("GET", "/M9Yv6VB2", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusMovedPermanently {
		t.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusMovedPermanently)
	}

	if url := resp.Header.Get("Location"); url != "https://www.example.com/" {
		t.Errorf("unexpected redirection url: %s", url)
	}

	r = httptest.NewRequest("GET", "/not-there", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	resp = w.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}
