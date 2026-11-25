// Command shorty runs the sample Go URL shortener HTTP server.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/oakwood-commons/sln-agentic-code/examples/go-app/internal/shortener"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, func(srv *http.Server) error {
		return srv.ListenAndServe()
	}))
}

func newServer(addr string, store shortener.Store) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           shortener.NewHandler(store),
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func run(args []string, stdout, stderr io.Writer, serve func(*http.Server) error) int {
	fs := flag.NewFlagSet("shorty", flag.ContinueOnError)
	fs.SetOutput(stderr)

	addr := fs.String("addr", ":8080", "HTTP listen address")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if fs.NArg() > 0 {
		_, _ = fmt.Fprintf(stderr, "shorty: unexpected arguments %v\n", fs.Args())
		return 2
	}

	srv := newServer(*addr, shortener.NewMemoryStore())
	_, _ = fmt.Fprintf(stdout, "shorty listening on %s\n", *addr)
	if err := serve(srv); err != nil && !errors.Is(err, http.ErrServerClosed) {
		_, _ = fmt.Fprintf(stderr, "shorty: %v\n", err)
		return 1
	}
	return 0
}
