package main

import (
	"bytes"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestRunCLI(t *testing.T) {
	t.Parallel()

	t.Run("help flag exits 0", func(t *testing.T) {
		t.Parallel()
		var out, errBuf bytes.Buffer
		code := run([]string{"-h"}, &out, &errBuf, func(*http.Server) error { return nil })
		if code != 0 {
			t.Fatalf("exit code = %d, want 0", code)
		}
		if !strings.Contains(errBuf.String(), "-addr") {
			t.Errorf("expected usage mentioning -addr, got %q", errBuf.String())
		}
	})

	t.Run("starts server with configured address", func(t *testing.T) {
		t.Parallel()
		var out, errBuf bytes.Buffer
		var gotAddr string
		code := run([]string{"-addr", "127.0.0.1:9090"}, &out, &errBuf, func(srv *http.Server) error {
			gotAddr = srv.Addr
			return http.ErrServerClosed
		})
		if code != 0 {
			t.Fatalf("exit code = %d, want 0", code)
		}
		if gotAddr != "127.0.0.1:9090" {
			t.Errorf("srv.Addr = %q, want 127.0.0.1:9090", gotAddr)
		}
	})

	t.Run("rejects invalid flags and unexpected positional args", func(t *testing.T) {
		t.Parallel()
		var out, errBuf bytes.Buffer
		if code := run([]string{"-bad"}, &out, &errBuf, func(*http.Server) error { return nil }); code != 2 {
			t.Fatalf("exit code for -bad = %d, want 2", code)
		}
		if code := run([]string{"extra"}, &out, &errBuf, func(*http.Server) error { return nil }); code != 2 {
			t.Fatalf("exit code for positional arg = %d, want 2", code)
		}
	})

	t.Run("listen failure exits 1", func(t *testing.T) {
		t.Parallel()
		var out, errBuf bytes.Buffer
		code := run(nil, &out, &errBuf, func(*http.Server) error {
			return errors.New("address already in use")
		})
		if code != 1 {
			t.Fatalf("exit code = %d, want 1", code)
		}
	})
}
