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
	server := &http.Server{Addr: addr, Handler: handler}

	tlsConfig, err := (submission.ServerTLSConfig{
		CertFile:          os.Getenv("MOCK_BOLAGSVERKET_TLS_CERT_FILE"),
		KeyFile:           os.Getenv("MOCK_BOLAGSVERKET_TLS_KEY_FILE"),
		CAFile:            os.Getenv("MOCK_BOLAGSVERKET_TLS_CA_FILE"),
		RequireClientCert: os.Getenv("MOCK_BOLAGSVERKET_REQUIRE_CLIENT_CERT") == "1",
	}).TLS()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if tlsConfig != nil {
		server.TLSConfig = tlsConfig
	}

	fmt.Fprintf(os.Stderr, "mock-bolagsverket listening on %s\n", addr)
	if apiKey != "" {
		fmt.Fprintln(os.Stderr, "mock-bolagsverket requires bearer auth")
	}
	if tlsConfig != nil {
		fmt.Fprintln(os.Stderr, "mock-bolagsverket serving HTTPS")
		if server.TLSConfig.ClientAuth != 0 {
			fmt.Fprintln(os.Stderr, "mock-bolagsverket requires client certificates")
		}
		if err := server.ListenAndServeTLS("", ""); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
