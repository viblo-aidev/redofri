package submission

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

// ClientTLSConfig configures TLS for the submission HTTP client.
type ClientTLSConfig struct {
	CertFile string
	KeyFile  string
	CAFile   string
}

// ServerTLSConfig configures TLS for the mock server.
type ServerTLSConfig struct {
	CertFile          string
	KeyFile           string
	CAFile            string
	RequireClientCert bool
}

// TLS returns a client TLS config, or nil when TLS files are not configured.
func (c ClientTLSConfig) TLS() (*tls.Config, error) {
	if c.CertFile == "" && c.KeyFile == "" && c.CAFile == "" {
		return nil, nil
	}
	if c.CertFile == "" || c.KeyFile == "" {
		return nil, fmt.Errorf("client TLS requires both cert and key files")
	}

	cert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("load client certificate: %w", err)
	}

	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
	}
	if c.CAFile != "" {
		pool, err := loadCertPool(c.CAFile)
		if err != nil {
			return nil, fmt.Errorf("load CA certificate: %w", err)
		}
		tlsConfig.RootCAs = pool
	}
	return tlsConfig, nil
}

// TLS returns a server TLS config, or nil when TLS files are not configured.
func (c ServerTLSConfig) TLS() (*tls.Config, error) {
	if c.CertFile == "" && c.KeyFile == "" && c.CAFile == "" && !c.RequireClientCert {
		return nil, nil
	}
	if c.CertFile == "" || c.KeyFile == "" {
		return nil, fmt.Errorf("server TLS requires both cert and key files")
	}

	cert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("load server certificate: %w", err)
	}

	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
	}
	if c.CAFile != "" {
		pool, err := loadCertPool(c.CAFile)
		if err != nil {
			return nil, fmt.Errorf("load CA certificate: %w", err)
		}
		tlsConfig.ClientCAs = pool
	}
	if c.RequireClientCert {
		if tlsConfig.ClientCAs == nil {
			return nil, fmt.Errorf("client cert verification requires a CA file")
		}
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}
	return tlsConfig, nil
}

func loadCertPool(path string) (*x509.CertPool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(data) {
		return nil, fmt.Errorf("failed to parse PEM certificates from %s", path)
	}
	return pool, nil
}
