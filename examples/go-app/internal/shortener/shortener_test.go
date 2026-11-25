package shortener_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/oakwood-commons/sln-agentic-code/examples/go-app/internal/shortener"
)

func TestShortenAndRedirect(t *testing.T) {
	t.Parallel()

	handler := shortener.NewHandler(shortener.NewMemoryStore())

	body := bytes.NewBufferString(`{"url":"https://github.com/oakwood-commons/scafctl"}`)
	createReq := httptest.NewRequest(http.MethodPost, "/links", body)
	createRec := httptest.NewRecorder()
	handler.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("POST /links status = %d, want %d (body: %s)", createRec.Code, http.StatusCreated, createRec.Body.String())
	}

	var created struct {
		Code string `json:"code"`
		URL  string `json:"url"`
	}
	if err := json.Unmarshal(createRec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(created.Code) != 6 {
		t.Fatalf("code = %q, want 6 hex chars", created.Code)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/"+created.Code, nil)
	getRec := httptest.NewRecorder()
	handler.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusFound {
		t.Fatalf("GET /%s status = %d, want %d", created.Code, getRec.Code, http.StatusFound)
	}
	if got := getRec.Header().Get("Location"); got != "https://github.com/oakwood-commons/scafctl" {
		t.Errorf("Location = %q, want original URL", got)
	}
}

func TestValidationAndNotFound(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{
			name:       "rejects malformed JSON",
			method:     http.MethodPost,
			path:       "/links",
			body:       `{"url":`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "rejects empty URL",
			method:     http.MethodPost,
			path:       "/links",
			body:       `{"url":"   "}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "rejects javascript scheme",
			method:     http.MethodPost,
			path:       "/links",
			body:       `{"url":"javascript:alert(1)"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "rejects hostless port-only URL",
			method:     http.MethodPost,
			path:       "/links",
			body:       `{"url":"http://:8080/path"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "rejects trailing JSON values",
			method:     http.MethodPost,
			path:       "/links",
			body:       `{"url":"https://example.com"}{"url":"https://other.com"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "rejects oversized request body with 413",
			method:     http.MethodPost,
			path:       "/links",
			body:       `{"url":"https://example.com/` + strings.Repeat("a", 4200) + `"}`,
			wantStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name:       "returns 404 for unknown short code",
			method:     http.MethodGet,
			path:       "/missing",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			handler := shortener.NewHandler(shortener.NewMemoryStore())
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)
			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tc.wantStatus, rec.Body.String())
			}
		})
	}
}

func TestMemoryStoreMissingLookupWrapsErrNotFound(t *testing.T) {
	t.Parallel()

	store := shortener.NewMemoryStore()
	if _, err := store.Lookup("nope"); !errors.Is(err, shortener.ErrNotFound) {
		t.Fatalf("Lookup() error = %v, want errors.Is(..., ErrNotFound)", err)
	}
}

type stubStore map[string]string

func (s stubStore) Save(code, targetURL string) error {
	s[code] = targetURL
	return nil
}

func (s stubStore) Lookup(code string) (string, error) {
	target, ok := s[code]
	if !ok {
		return "", shortener.ErrNotFound
	}
	return target, nil
}

func TestRedirectRejectsCorruptStoredScheme(t *testing.T) {
	t.Parallel()

	handler := shortener.NewHandler(stubStore{"bad123": "javascript:alert(1)"})
	req := httptest.NewRequest(http.MethodGet, "/bad123", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d (body: %s)", rec.Code, http.StatusInternalServerError, rec.Body.String())
	}
}
