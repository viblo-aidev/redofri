package submission

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPClient(t *testing.T) {
	var authHeader string
	var paths []string
	var checksumTokenReq CreateTokenRequest
	var checksumReq CreateChecksumRequest
	var tokenReq CreateTokenRequest
	var checkReq CheckRequest
	var submitReq SubmitRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		paths = append(paths, r.URL.Path)

		switch {
		case r.URL.Path == "/hamta-arsredovisningsinformation/v1.1/skapa-inlamningtoken":
			if err := json.NewDecoder(r.Body).Decode(&checksumTokenReq); err != nil {
				t.Fatalf("decode checksum token request: %v", err)
			}
			_ = json.NewEncoder(w).Encode(CreateTokenResponse{Token: "chk-123"})
		case r.URL.Path == "/hamta-arsredovisningsinformation/v1.1/skapa-kontrollsumma/chk-123":
			if err := json.NewDecoder(r.Body).Decode(&checksumReq); err != nil {
				t.Fatalf("decode checksum request: %v", err)
			}
			_ = json.NewEncoder(w).Encode(CreateChecksumResponse{Checksum: "sum-123", Algorithm: "SHA-256"})
		case r.URL.Path == "/lamna-in-arsredovisning/v2.1/skapa-inlamningtoken/":
			if err := json.NewDecoder(r.Body).Decode(&tokenReq); err != nil {
				t.Fatalf("decode token request: %v", err)
			}
			_ = json.NewEncoder(w).Encode(CreateTokenResponse{Token: "tok-123"})
		case r.URL.Path == "/lamna-in-arsredovisning/v2.1/kontrollera/tok-123":
			if err := json.NewDecoder(r.Body).Decode(&checkReq); err != nil {
				t.Fatalf("decode check request: %v", err)
			}
			_ = json.NewEncoder(w).Encode(CheckResponse{OrgNumber: "556000-1111"})
		case r.URL.Path == "/lamna-in-arsredovisning/v2.1/inlamning/tok-123":
			if err := json.NewDecoder(r.Body).Decode(&submitReq); err != nil {
				t.Fatalf("decode submit request: %v", err)
			}
			_ = json.NewEncoder(w).Encode(SubmitResponse{DocumentInfo: SubmitDocumentInfo{IDNumber: "49679"}})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client, err := NewHTTPClient(server.URL, server.Client(), "secret")
	if err != nil {
		t.Fatalf("NewHTTPClient: %v", err)
	}

	if _, err := client.CreateChecksumToken(context.Background(), CreateTokenRequest{SenderPersonalNumber: "190001010106", OrgNumber: "556000-1111"}); err != nil {
		t.Fatalf("CreateChecksumToken: %v", err)
	}
	if _, err := client.CreateChecksum(context.Background(), CreateChecksumRequest{Token: "chk-123", File: base64.StdEncoding.EncodeToString([]byte("doc"))}); err != nil {
		t.Fatalf("CreateChecksum: %v", err)
	}
	if _, err := client.CreateToken(context.Background(), CreateTokenRequest{SenderPersonalNumber: "190001010106", OrgNumber: "556000-1111"}); err != nil {
		t.Fatalf("CreateToken: %v", err)
	}
	if _, err := client.Check(context.Background(), CheckRequest{Token: "tok-123", Document: Document{File: base64.StdEncoding.EncodeToString([]byte("doc")), Type: DefaultDocumentType}}); err != nil {
		t.Fatalf("Check: %v", err)
	}
	if _, err := client.Submit(context.Background(), SubmitRequest{Token: "tok-123", SignerPersonalNumber: "198301019876", Document: Document{File: base64.StdEncoding.EncodeToString([]byte("doc")), Type: DefaultDocumentType}}); err != nil {
		t.Fatalf("Submit: %v", err)
	}

	if got, want := strings.Join(paths, ","), "/hamta-arsredovisningsinformation/v1.1/skapa-inlamningtoken,/hamta-arsredovisningsinformation/v1.1/skapa-kontrollsumma/chk-123,/lamna-in-arsredovisning/v2.1/skapa-inlamningtoken/,/lamna-in-arsredovisning/v2.1/kontrollera/tok-123,/lamna-in-arsredovisning/v2.1/inlamning/tok-123"; got != want {
		t.Fatalf("paths = %q, want %q", got, want)
	}
	if got, want := authHeader, "Bearer secret"; got != want {
		t.Fatalf("auth header = %q, want %q", got, want)
	}
	if checksumTokenReq.OrgNumber != "556000-1111" || checksumTokenReq.SenderPersonalNumber != "190001010106" {
		t.Fatalf("unexpected checksum token request: %#v", checksumTokenReq)
	}
	if checksumReq.File == "" {
		t.Fatal("expected checksum request with file")
	}
	if tokenReq.OrgNumber != "556000-1111" || tokenReq.SenderPersonalNumber != "190001010106" {
		t.Fatalf("unexpected token request: %#v", tokenReq)
	}
	if checkReq.Document.Type != DefaultDocumentType {
		t.Fatalf("check type = %q, want %q", checkReq.Document.Type, DefaultDocumentType)
	}
	if submitReq.SignerPersonalNumber != "198301019876" {
		t.Fatalf("submit signer pnr = %q", submitReq.SignerPersonalNumber)
	}
	decoded, err := base64.StdEncoding.DecodeString(submitReq.Document.File)
	if err != nil {
		t.Fatalf("decode submit file: %v", err)
	}
	if string(decoded) != "doc" {
		t.Fatalf("submit document = %q, want doc", string(decoded))
	}
}

func TestHTTPClientStatusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = io.WriteString(w, "upstream failed")
	}))
	defer server.Close()

	client, err := NewHTTPClient(server.URL, server.Client(), "")
	if err != nil {
		t.Fatalf("NewHTTPClient: %v", err)
	}

	_, err = client.CreateToken(context.Background(), CreateTokenRequest{OrgNumber: "556000-1111", SenderPersonalNumber: "190001010106"})
	if err == nil {
		t.Fatal("expected error from non-2xx response")
	}
	if !strings.Contains(err.Error(), "upstream failed") {
		t.Fatalf("error = %q, want upstream body", err)
	}
}
