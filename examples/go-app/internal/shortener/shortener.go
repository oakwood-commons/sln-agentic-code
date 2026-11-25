// Package shortener implements an HTTP URL-shortening service backed by a
// pluggable Store interface.
package shortener

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const maxRequestBodyBytes = 4096

var (
	// ErrNotFound indicates a short code does not exist in the store.
	ErrNotFound = errors.New("short link not found")
	// ErrInvalidURL indicates the submitted target URL is empty or not http(s).
	ErrInvalidURL = errors.New("url must be a valid http or https address")
)

// Store persists and looks up shortened URLs by their short code.
type Store interface {
	Save(code, targetURL string) error
	Lookup(code string) (string, error)
}

// MemoryStore is a concurrency-safe in-memory Store adapter.
type MemoryStore struct {
	mu    sync.RWMutex
	links map[string]string
}

// NewMemoryStore returns an empty MemoryStore ready for concurrent use.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		links: make(map[string]string),
	}
}

// Save stores targetURL under code.
func (m *MemoryStore) Save(code, targetURL string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.links[code] = targetURL
	return nil
}

// Lookup returns the target URL for code or ErrNotFound.
func (m *MemoryStore) Lookup(code string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	target, ok := m.links[code]
	if !ok {
		return "", fmt.Errorf("lookup %q: %w", code, ErrNotFound)
	}
	return target, nil
}

type createRequest struct {
	URL string `json:"url"`
}

type linkResponse struct {
	Code string `json:"code"`
	URL  string `json:"url"`
}

// NewHandler returns an http.Handler serving POST /links and GET /{code}.
func NewHandler(store Store) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /links", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
		dec := json.NewDecoder(r.Body)
		var req createRequest
		if err := dec.Decode(&req); err != nil {
			var maxBytesErr *http.MaxBytesError
			if errors.As(err, &maxBytesErr) {
				writeJSONError(w, http.StatusRequestEntityTooLarge, "request body too large")
				return
			}
			writeJSONError(w, http.StatusBadRequest, "invalid JSON request body")
			return
		}
		if err := dec.Decode(&struct{}{}); err != io.EOF {
			var maxBytesErr *http.MaxBytesError
			if errors.As(err, &maxBytesErr) {
				writeJSONError(w, http.StatusRequestEntityTooLarge, "request body too large")
				return
			}
			writeJSONError(w, http.StatusBadRequest, "invalid JSON request body")
			return
		}

		target, err := validateTargetURL(req.URL)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		code, err := generateCode()
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to generate short code")
			return
		}
		if err := store.Save(code, target); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to store short link")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(linkResponse{Code: code, URL: target})
	})

	mux.HandleFunc("GET /{code}", func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		target, err := store.Lookup(code)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				writeJSONError(w, http.StatusNotFound, "short link not found")
				return
			}
			writeJSONError(w, http.StatusInternalServerError, "failed to look up short link")
			return
		}
		safeTarget, err := validateTargetURL(target)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "stored short link has invalid target URL")
			return
		}
		// #nosec G710 -- URL shortener intentionally redirects to stored, scheme-validated http(s) targets.
		http.Redirect(w, r, safeTarget, http.StatusFound)
	})

	return mux
}

func validateTargetURL(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	parsed, err := url.Parse(trimmed)
	if err != nil || parsed.Hostname() == "" {
		return "", ErrInvalidURL
	}
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", ErrInvalidURL
	}
	return parsed.String(), nil
}

func generateCode() (string, error) {
	var b [3]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}
	return hex.EncodeToString(b[:]), nil
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
