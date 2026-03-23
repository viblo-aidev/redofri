package submission

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	mockAgreementText = "Ett Eget utrymme har nu skapats for det foretag som du har angett."
)

// NewMockHandler returns a stateful HTTP handler for local submission testing.
func NewMockHandler(apiKey string) http.Handler {
	server := &mockServer{
		apiKey: apiKey,
		tokens: make(map[string]mockToken),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", server.handleHealthz)
	mux.HandleFunc("/hamta-arsredovisningsinformation/v1.1/skapa-inlamningtoken", server.handleCreateChecksumToken)
	mux.HandleFunc("/hamta-arsredovisningsinformation/v1.1/skapa-kontrollsumma/", server.handleCreateChecksum)
	mux.HandleFunc("/lamna-in-arsredovisning/v2.1/skapa-inlamningtoken/", server.handleCreateToken)
	mux.HandleFunc("/lamna-in-arsredovisning/v2.1/kontrollera/", server.handleCheck)
	mux.HandleFunc("/lamna-in-arsredovisning/v2.1/inlamning/", server.handleSubmit)
	return mux
}

type mockServer struct {
	mu     sync.Mutex
	apiKey string
	tokens map[string]mockToken
	nextID int
}

type mockToken struct {
	OrgNumber            string
	SenderPersonalNumber string
	CreatedAt            time.Time
}

func (s *mockServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *mockServer) handleCreateToken(w http.ResponseWriter, r *http.Request) {
	s.handleCreateTokenLike(w, r)
}

func (s *mockServer) handleCreateChecksumToken(w http.ResponseWriter, r *http.Request) {
	s.handleCreateTokenLike(w, r)
}

func (s *mockServer) handleCreateTokenLike(w http.ResponseWriter, r *http.Request) {
	if !s.authorize(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req CreateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}
	if strings.TrimSpace(req.OrgNumber) == "" {
		writeJSONError(w, http.StatusBadRequest, "orgnr is required")
		return
	}
	if !matchesClientCertificateOrgNumber(r, req.OrgNumber) {
		writeJSONError(w, http.StatusForbidden, "client certificate org number does not match orgnr")
		return
	}
	if strings.TrimSpace(req.SenderPersonalNumber) == "" {
		writeJSONError(w, http.StatusBadRequest, "pnr is required")
		return
	}

	token, err := randomToken()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("generate token: %v", err))
		return
	}

	s.mu.Lock()
	s.tokens[token] = mockToken{
		OrgNumber:            req.OrgNumber,
		SenderPersonalNumber: req.SenderPersonalNumber,
		CreatedAt:            time.Now().UTC(),
	}
	s.mu.Unlock()

	writeJSON(w, http.StatusOK, CreateTokenResponse{
		Token:                token,
		AgreementText:        mockAgreementText,
		AgreementVersionDate: "2017-12-06",
	})
}

func (s *mockServer) handleCreateChecksum(w http.ResponseWriter, r *http.Request) {
	if !s.authorize(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	tokenID := strings.TrimPrefix(r.URL.Path, "/hamta-arsredovisningsinformation/v1.1/skapa-kontrollsumma/")
	var req CreateChecksumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}
	req.Token = tokenID

	if strings.TrimSpace(tokenID) == "" {
		writeJSONError(w, http.StatusBadRequest, "token is required")
		return
	}
	s.mu.Lock()
	_, ok := s.tokens[tokenID]
	s.mu.Unlock()
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "7003=Felaktig token.")
		return
	}
	document, err := base64.StdEncoding.DecodeString(req.File)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("invalid fil: %v", err))
		return
	}
	writeJSON(w, http.StatusOK, CreateChecksumResponse{
		Checksum:  mockChecksum(document),
		Algorithm: "SHA-256",
	})
}

func (s *mockServer) handleCheck(w http.ResponseWriter, r *http.Request) {
	if !s.authorize(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	tokenID := strings.TrimPrefix(r.URL.Path, "/lamna-in-arsredovisning/v2.1/kontrollera/")
	var req CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}
	req.Token = tokenID

	token, document, ok := s.validateDocumentRequest(w, req.Token, req.Document)
	if !ok {
		return
	}

	findings := evaluateMockDocument(document)
	writeJSON(w, http.StatusOK, CheckResponse{
		OrgNumber: token.OrgNumber,
		Findings:  findings,
	})
}

func (s *mockServer) handleSubmit(w http.ResponseWriter, r *http.Request) {
	if !s.authorize(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	tokenID := strings.TrimPrefix(r.URL.Path, "/lamna-in-arsredovisning/v2.1/inlamning/")
	var req SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}
	req.Token = tokenID

	token, document, ok := s.validateDocumentRequest(w, req.Token, req.Document)
	if !ok {
		return
	}
	if strings.TrimSpace(req.SignerPersonalNumber) == "" {
		writeJSONError(w, http.StatusBadRequest, "undertecknare is required")
		return
	}

	findings := evaluateMockDocument(document)
	if containsErrorFindings(findings) {
		writeJSON(w, http.StatusBadRequest, CheckResponse{
			OrgNumber: token.OrgNumber,
			Findings:  findings,
		})
		return
	}

	s.mu.Lock()
	s.nextID++
	id := fmt.Sprintf("%d", 49000+s.nextID)
	s.mu.Unlock()

	writeJSON(w, http.StatusOK, SubmitResponse{
		OrgNumber:            token.OrgNumber,
		SenderPersonalNumber: token.SenderPersonalNumber,
		SignerPersonalNumber: req.SignerPersonalNumber,
		DocumentInfo: SubmitDocumentInfo{
			Type:           req.Document.Type,
			DocumentLength: len(document),
			IDNumber:       id,
			Checksum:       mockChecksum(document),
		},
		URL: "https://arsredovisning-accept2.bolagsverket.se/lamna-in/visa/engagemang/18772",
	})
}

func (s *mockServer) authorize(w http.ResponseWriter, r *http.Request) bool {
	if s.apiKey == "" {
		return true
	}
	if r.Header.Get("Authorization") != "Bearer "+s.apiKey {
		writeJSONError(w, http.StatusUnauthorized, "missing or invalid bearer token")
		return false
	}
	return true
}

func (s *mockServer) validateDocumentRequest(w http.ResponseWriter, tokenID string, doc Document) (mockToken, []byte, bool) {
	if strings.TrimSpace(tokenID) == "" {
		writeJSONError(w, http.StatusBadRequest, "token is required")
		return mockToken{}, nil, false
	}
	if doc.File == "" {
		writeJSONError(w, http.StatusBadRequest, "handling.fil is required")
		return mockToken{}, nil, false
	}
	if doc.Type == "" {
		writeJSONError(w, http.StatusBadRequest, "handling.typ is required")
		return mockToken{}, nil, false
	}

	s.mu.Lock()
	tok, ok := s.tokens[tokenID]
	s.mu.Unlock()
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "7003=Felaktig token.")
		return mockToken{}, nil, false
	}

	document, err := base64.StdEncoding.DecodeString(doc.File)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("invalid handling.fil: %v", err))
		return mockToken{}, nil, false
	}
	return tok, document, true
}

func evaluateMockDocument(document []byte) []Finding {
	text := string(document)
	findings := make([]Finding, 0, 2)

	if strings.Contains(text, "MOCK_REMOTE_WARNING") {
		findings = append(findings, Finding{
			Code:     "1165",
			Message:  "Datum for underskrift av faststallelseintyget far inte vara tidigare an datum for arsstamman.",
			Severity: SeverityWarn,
			TechnicalInformation: []TechnicalInformation{
				{Element: "UnderskriftFastallelseintygDatum", Value: "2019-01-09"},
				{Element: "Arsstamma", Value: "2019-01-10"},
			},
		})
	}
	if strings.Contains(text, "MOCK_REMOTE_ERROR") {
		findings = append(findings, Finding{
			Code:     "7003",
			Message:  "Felaktig token.",
			Severity: SeverityError,
		})
	}

	return findings
}

func containsErrorFindings(findings []Finding) bool {
	for _, finding := range findings {
		if finding.Severity == SeverityError {
			return true
		}
	}
	return false
}

func mockChecksum(document []byte) string {
	sum := sha256.Sum256(document)
	return base64.StdEncoding.EncodeToString(sum[:])
}

func randomToken() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func matchesClientCertificateOrgNumber(r *http.Request, orgNumber string) bool {
	if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
		return true
	}
	cert := r.TLS.PeerCertificates[0]
	serialOrgNumber := serialNumberOrgNumber(cert.Subject)
	if serialOrgNumber == "" {
		return false
	}
	return normalizeOrgNumber(serialOrgNumber) == normalizeOrgNumber(orgNumber)
}

func serialNumberOrgNumber(name pkix.Name) string {
	serial := strings.TrimSpace(name.SerialNumber)
	if serial == "" {
		return ""
	}
	serial = strings.TrimPrefix(serial, "16")
	return normalizeOrgNumber(serial)
}

func normalizeOrgNumber(value string) string {
	value = strings.ReplaceAll(value, "-", "")
	value = strings.ReplaceAll(value, " ", "")
	return value
}
