// Package handler provides the http.Handler for API endpoints.
package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/ericyan/surl/internal/shortener"
	"github.com/ericyan/surl/pkg/kv"
)

// ShorteningRequest represents a POST /submit request.
type ShorteningRequest struct {
	URL string `json:"url"`
}

// ShorteningResponse represents a POST /submit response.
type ShorteningResponse struct {
	URL      string `json:"url"`
	ShortURL string `json:"shorten_url"`
}

type handler struct {
	*shortener.Shortener

	kvstore kv.Store
}

// New returns a new HTTP API handler.
func New(kvstore kv.Store) http.Handler {
	return &handler{shortener.New(), kvstore}
}

// ServeHTTP imeplements the http.Handler interface.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(`Try this: <pre>curl -X POST -H "Content-Type: application/json" -d '{"url": "https://www.example.com/"}' http://%s/submit</pre>`, r.Host)))
			return
		}

		key := strings.TrimPrefix(r.URL.Path, "/")
		url, err := h.kvstore.Get([]byte(key))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Redirect(w, r, string(url), http.StatusMovedPermanently)
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var req ShorteningRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = url.ParseRequestURI(req.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		shortURL := h.Shorten(req.URL)
		if err := h.kvstore.Put([]byte(shortURL), []byte(req.URL)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(ShorteningResponse{req.URL, shortURL})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		w.Write(resp)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}
