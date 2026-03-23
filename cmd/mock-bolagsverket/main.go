package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/redofri/redofri/pkg/submission"
)

func main() {
	addr := os.Getenv("MOCK_BOLAGSVERKET_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	apiKey := os.Getenv("MOCK_BOLAGSVERKET_API_KEY")

	handler := submission.NewMockHandler(apiKey)
	fmt.Fprintf(os.Stderr, "mock-bolagsverket listening on %s\n", addr)
	if apiKey != "" {
		fmt.Fprintln(os.Stderr, "mock-bolagsverket requires bearer auth")
	}
	if err := http.ListenAndServe(addr, handler); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
