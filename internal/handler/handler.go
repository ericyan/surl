// Package handler provides the http.Handler for API endpoints.
package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

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
func New() http.Handler {
	kvstore, _ := kv.NewInMemoryStore()

	return &handler{shortener.New(), kvstore}
}

// ServeHTTP imeplements the http.Handler interface.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
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
