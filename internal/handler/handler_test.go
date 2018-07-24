package handler

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ericyan/surl/pkg/kv"
)

func testHandler(t *testing.T, h http.Handler, r *http.Request, code int, headers map[string]string, body []byte) {
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	resp := w.Result()

	if resp.StatusCode != code {
		t.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, code)
	}

	if headers != nil {
		for k, v := range headers {
			if respHeader := resp.Header.Get(k); respHeader != v {
				t.Errorf("unexpected header %s: got %s, want %s", k, respHeader, v)
			}
		}
	}

	if body != nil {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(respBody, body) {
			t.Errorf("unexpected body: got %s, want %s", respBody, body)
		}
	}
}

func TestHandler(t *testing.T) {
	store, _ := kv.NewInMemoryStore()
	handler := New(store)

	postExample := httptest.NewRequest("POST", "/submit", bytes.NewBufferString(`{"url": "https://www.example.com/"}`))
	testHandler(t, handler, postExample, http.StatusCreated, nil, []byte(`{"url":"https://www.example.com/","shorten_url":"M9Yv6VB2"}`))

	getExample := httptest.NewRequest("GET", "/M9Yv6VB2", nil)
	testHandler(t, handler, getExample, http.StatusMovedPermanently, map[string]string{"Location": "https://www.example.com/"}, nil)

	getNonexistent := httptest.NewRequest("GET", "/nonexistent", nil)
	testHandler(t, handler, getNonexistent, http.StatusNotFound, nil, nil)
}
