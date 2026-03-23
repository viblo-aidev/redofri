package submission

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMockServerHappyPath(t *testing.T) {
	h := NewMockHandler("secret")
	checksumTokenResp := doJSON[CreateTokenResponse](t, h, http.MethodPost, "/hamta-arsredovisningsinformation/v1.1/skapa-inlamningtoken", CreateTokenRequest{
		SenderPersonalNumber: "190001010106",
		OrgNumber:            "556000-1111",
	}, "secret", http.StatusOK)
	checksumResp := doJSON[CreateChecksumResponse](t, h, http.MethodPost, "/hamta-arsredovisningsinformation/v1.1/skapa-kontrollsumma/"+checksumTokenResp.Token, CreateChecksumRequest{
		File: base64.StdEncoding.EncodeToString([]byte("test-document")),
	}, "secret", http.StatusOK)
	if checksumResp.Algorithm != "SHA-256" {
		t.Fatalf("algorithm = %q, want SHA-256", checksumResp.Algorithm)
	}

	tokenResp := doJSON[CreateTokenResponse](t, h, http.MethodPost, "/lamna-in-arsredovisning/v2.1/skapa-inlamningtoken/", CreateTokenRequest{
		SenderPersonalNumber: "190001010106",
		OrgNumber:            "556000-1111",
	}, "secret", http.StatusOK)
	if tokenResp.Token == "" {
		t.Fatal("expected token")
	}

	document := []byte("test-document")
	encoded := base64.StdEncoding.EncodeToString(document)
	checkResp := doJSON[CheckResponse](t, h, http.MethodPost, "/lamna-in-arsredovisning/v2.1/kontrollera/"+tokenResp.Token, CheckRequest{
		Document: Document{File: encoded, Type: DefaultDocumentType},
	}, "secret", http.StatusOK)
	if checkResp.OrgNumber != "556000-1111" {
		t.Fatalf("orgnr = %q, want 556000-1111", checkResp.OrgNumber)
	}

	submitResp := doJSON[SubmitResponse](t, h, http.MethodPost, "/lamna-in-arsredovisning/v2.1/inlamning/"+tokenResp.Token, SubmitRequest{
		SignerPersonalNumber: "198301019876",
		Document:             Document{File: encoded, Type: DefaultDocumentType},
	}, "secret", http.StatusOK)
	if submitResp.DocumentInfo.IDNumber == "" {
		t.Fatal("expected document ID")
	}
	if submitResp.DocumentInfo.Type != DefaultDocumentType {
		t.Fatalf("type = %q, want %q", submitResp.DocumentInfo.Type, DefaultDocumentType)
	}
}

func TestMockServerWrongToken(t *testing.T) {
	h := NewMockHandler("")
	respChecksum := doErrorJSON(t, h, http.MethodPost, "/hamta-arsredovisningsinformation/v1.1/skapa-kontrollsumma/bad-token", CreateChecksumRequest{File: base64.StdEncoding.EncodeToString([]byte("test-document"))}, "", http.StatusBadRequest)
	if respChecksum["error"] != "7003=Felaktig token." {
		t.Fatalf("unexpected checksum error: %#v", respChecksum)
	}
	resp := doErrorJSON(t, h, http.MethodPost, "/lamna-in-arsredovisning/v2.1/kontrollera/bad-token", CheckRequest{
		Document: Document{File: base64.StdEncoding.EncodeToString([]byte("test-document")), Type: DefaultDocumentType},
	}, "", http.StatusBadRequest)
	if resp["error"] != "7003=Felaktig token." {
		t.Fatalf("unexpected error: %#v", resp)
	}
}

func TestMockServerDocumentMarkers(t *testing.T) {
	h := NewMockHandler("")
	tokenResp := doJSON[CreateTokenResponse](t, h, http.MethodPost, "/lamna-in-arsredovisning/v2.1/skapa-inlamningtoken/", CreateTokenRequest{
		SenderPersonalNumber: "190001010106",
		OrgNumber:            "556000-1111",
	}, "", http.StatusOK)

	warningDoc := base64.StdEncoding.EncodeToString([]byte("test MOCK_REMOTE_WARNING"))
	checkResp := doJSON[CheckResponse](t, h, http.MethodPost, "/lamna-in-arsredovisning/v2.1/kontrollera/"+tokenResp.Token, CheckRequest{
		Document: Document{File: warningDoc, Type: DefaultDocumentType},
	}, "", http.StatusOK)
	if len(checkResp.Findings) != 1 || checkResp.Findings[0].Code != "1165" {
		t.Fatalf("unexpected warning findings: %#v", checkResp.Findings)
	}

	errorDoc := base64.StdEncoding.EncodeToString([]byte("test MOCK_REMOTE_ERROR"))
	errorResp := doJSON[CheckResponse](t, h, http.MethodPost, "/lamna-in-arsredovisning/v2.1/inlamning/"+tokenResp.Token, SubmitRequest{
		SignerPersonalNumber: "198301019876",
		Document:             Document{File: errorDoc, Type: DefaultDocumentType},
	}, "", http.StatusBadRequest)
	if len(errorResp.Findings) != 1 || errorResp.Findings[0].Code != "7003" {
		t.Fatalf("unexpected error findings: %#v", errorResp.Findings)
	}
}

func TestMockServerUnauthorized(t *testing.T) {
	h := NewMockHandler("secret")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/lamna-in-arsredovisning/v2.1/skapa-inlamningtoken/", bytes.NewReader([]byte(`{"pnr":"190001010106","orgnr":"556000-1111"}`)))
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func doJSON[T any](t *testing.T, h http.Handler, method, path string, body any, apiKey string, wantStatus int) T {
	t.Helper()
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	h.ServeHTTP(rec, req)
	if rec.Code != wantStatus {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, wantStatus, rec.Body.String())
	}
	var resp T
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v; body=%s", err, rec.Body.String())
	}
	return resp
}

func doErrorJSON(t *testing.T, h http.Handler, method, path string, body any, apiKey string, wantStatus int) map[string]string {
	t.Helper()
	return doJSON[map[string]string](t, h, method, path, body, apiKey, wantStatus)
}
