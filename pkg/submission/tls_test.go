package submission

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestClientTLSConfig(t *testing.T) {
	dir := t.TempDir()
	caCertPEM, _, caCert, caKey := generateCA(t)
	clientCertPEM, clientKeyPEM := generateLeaf(t, caCert, caKey, true)

	caPath := writeBytes(t, dir, "ca.pem", caCertPEM)
	certPath := writeBytes(t, dir, "client.pem", clientCertPEM)
	keyPath := writeBytes(t, dir, "client.key", clientKeyPEM)

	cfg, err := (ClientTLSConfig{CertFile: certPath, KeyFile: keyPath, CAFile: caPath}).TLS()
	if err != nil {
		t.Fatalf("TLS: %v", err)
	}
	if cfg == nil || len(cfg.Certificates) != 1 || cfg.RootCAs == nil {
		t.Fatal("expected client TLS config with cert and root CAs")
	}
}

func TestServerTLSConfigRequiresCAForClientCert(t *testing.T) {
	dir := t.TempDir()
	_, _, caCert, caKey := generateCA(t)
	serverCertPEM, serverKeyPEM := generateLeaf(t, caCert, caKey, false)

	certPath := writeBytes(t, dir, "server.pem", serverCertPEM)
	keyPath := writeBytes(t, dir, "server.key", serverKeyPEM)

	_, err := (ServerTLSConfig{CertFile: certPath, KeyFile: keyPath, RequireClientCert: true}).TLS()
	if err == nil {
		t.Fatal("expected error without CA file when client certs are required")
	}
}

func generateCA(t *testing.T) ([]byte, []byte, *x509.Certificate, *rsa.PrivateKey) {
	return generateCustomCA(t, "ExpiTrust Test CA v8", "Expisoft AB")
}

func generateCustomCA(t *testing.T, cn, org string) ([]byte, []byte, *x509.Certificate, *rsa.PrivateKey) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: cn, Organization: []string{org}, Country: []string{"SE"}},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("CreateCertificate: %v", err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("ParseCertificate: %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}), cert, key
}

func generateLeaf(t *testing.T, caCert *x509.Certificate, caKey *rsa.PrivateKey, client bool) ([]byte, []byte) {
	return generateLeafWithSerial(t, caCert, caKey, client, "165560001111")
}

func generateLeafWithSerial(t *testing.T, caCert *x509.Certificate, caKey *rsa.PrivateKey, client bool, serialNumber string) ([]byte, []byte) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	serial, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		t.Fatalf("rand.Int: %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: "leaf", Organization: []string{"Example Client AB"}, Country: []string{"SE"}, SerialNumber: serialNumber},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
	}
	if client {
		tmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	} else {
		tmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
		tmpl.DNSNames = []string{"localhost"}
		tmpl.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")}
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, caCert, &key.PublicKey, caKey)
	if err != nil {
		t.Fatalf("CreateCertificate: %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
}

func writeBytes(t *testing.T, dir, name string, data []byte) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return path
}
