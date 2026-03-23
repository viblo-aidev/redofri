package submission

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/redofri/redofri/pkg/model"
	"github.com/redofri/redofri/pkg/validate"
)

func loadTestReport(t *testing.T) *model.AnnualReport {
	t.Helper()
	data, err := os.ReadFile("../../testdata/exempel1.json")
	if err != nil {
		t.Fatalf("reading test data: %v", err)
	}
	var r model.AnnualReport
	if err := json.Unmarshal(data, &r); err != nil {
		t.Fatalf("parsing test data: %v", err)
	}
	return &r
}

type fakeClient struct {
	calls             []string
	createReq         *CreateTokenRequest
	checksumTokenReq  *CreateTokenRequest
	checksumReq       *CreateChecksumRequest
	createResp        *CreateTokenResponse
	checksumTokenResp *CreateTokenResponse
	checksumResp      *CreateChecksumResponse
	checkResp         *CheckResponse
	submitResp        *SubmitResponse
	checkReq          *CheckRequest
	submitReq         *SubmitRequest
}

func (f *fakeClient) CreateChecksumToken(_ context.Context, req CreateTokenRequest) (*CreateTokenResponse, error) {
	f.calls = append(f.calls, "checksum-token")
	f.checksumTokenReq = &req
	return f.checksumTokenResp, nil
}

func (f *fakeClient) CreateChecksum(_ context.Context, req CreateChecksumRequest) (*CreateChecksumResponse, error) {
	f.calls = append(f.calls, "checksum")
	f.checksumReq = &req
	return f.checksumResp, nil
}

func (f *fakeClient) CreateToken(_ context.Context, req CreateTokenRequest) (*CreateTokenResponse, error) {
	f.calls = append(f.calls, "token")
	f.createReq = &req
	return f.createResp, nil
}

func (f *fakeClient) Check(_ context.Context, req CheckRequest) (*CheckResponse, error) {
	f.calls = append(f.calls, "check")
	f.checkReq = &req
	return f.checkResp, nil
}

func (f *fakeClient) Submit(_ context.Context, req SubmitRequest) (*SubmitResponse, error) {
	f.calls = append(f.calls, "submit")
	f.submitReq = &req
	return f.submitResp, nil
}

func TestServiceCheck(t *testing.T) {
	report := loadTestReport(t)
	client := &fakeClient{
		checksumTokenResp: &CreateTokenResponse{Token: "chk-123"},
		checksumResp:      &CreateChecksumResponse{Checksum: "sum-123", Algorithm: "SHA-256"},
		createResp:        &CreateTokenResponse{Token: "tok-123", AgreementVersionDate: "2026-01-02"},
		checkResp:         &CheckResponse{OrgNumber: report.Company.OrgNr},
	}

	result, err := NewService(client).Check(context.Background(), report, SubmitOptions{SenderPersonalNumber: "190001010106"})
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if got, want := len(client.calls), 4; got != want {
		t.Fatalf("call count = %d, want %d", got, want)
	}
	if client.checksumTokenReq == nil || client.checksumReq == nil {
		t.Fatal("expected checksum token and checksum requests")
	}
	if client.createReq == nil || client.createReq.OrgNumber != report.Company.OrgNr {
		t.Fatal("expected create token request with org number")
	}
	if client.createReq.SenderPersonalNumber != "190001010106" {
		t.Fatalf("sender pnr = %q, want 190001010106", client.createReq.SenderPersonalNumber)
	}
	if result.Token != "tok-123" {
		t.Fatalf("token = %q, want tok-123", result.Token)
	}
	if result.ChecksumToken != "chk-123" {
		t.Fatalf("checksum token = %q, want chk-123", result.ChecksumToken)
	}
	if result.Checksum != "sum-123" {
		t.Fatalf("checksum = %q, want sum-123", result.Checksum)
	}
	if result.ChecksumAlgorithm != "SHA-256" {
		t.Fatalf("algorithm = %q, want SHA-256", result.ChecksumAlgorithm)
	}
	if client.checkReq == nil {
		t.Fatal("expected check request")
	}
	if client.checkReq.Document.Type != DefaultDocumentType {
		t.Fatalf("document type = %q, want %q", client.checkReq.Document.Type, DefaultDocumentType)
	}
	decoded, err := base64.StdEncoding.DecodeString(client.checkReq.Document.File)
	if err != nil {
		t.Fatalf("decode base64 document: %v", err)
	}
	if len(decoded) == 0 {
		t.Fatal("expected generated document bytes in check request")
	}
	if !strings.Contains(string(decoded), `name="ixbrl.innehall.kontrollsumman" content="sum-123"`) {
		t.Fatal("expected checksum metadata in checked document")
	}
	if len(result.LocalFindings) != 0 {
		t.Fatalf("local findings = %d, want 0", len(result.LocalFindings))
	}
}

func TestServiceSubmitRunsCheckByDefault(t *testing.T) {
	report := loadTestReport(t)
	client := &fakeClient{
		checksumTokenResp: &CreateTokenResponse{Token: "chk-123"},
		checksumResp:      &CreateChecksumResponse{Checksum: "sum-123", Algorithm: "SHA-256"},
		createResp:        &CreateTokenResponse{Token: "tok-123"},
		checkResp:         &CheckResponse{OrgNumber: report.Company.OrgNr},
		submitResp: &SubmitResponse{
			OrgNumber:            report.Company.OrgNr,
			SenderPersonalNumber: "190001010106",
			SignerPersonalNumber: "198301019876",
			DocumentInfo: SubmitDocumentInfo{
				Type:           DefaultDocumentType,
				DocumentLength: 123,
				IDNumber:       "49679",
				Checksum:       "sum-123",
			},
			URL: "https://example.test/submission/49679",
		},
	}

	result, err := NewService(client).Submit(context.Background(), report, SubmitOptions{
		SenderPersonalNumber: "190001010106",
		SignerPersonalNumber: "198301019876",
		EmailAddresses:       []string{"jag@foretag.com"},
	})
	if err != nil {
		t.Fatalf("Submit returned error: %v", err)
	}
	if got, want := len(client.calls), 5; got != want {
		t.Fatalf("call count = %d, want %d", got, want)
	}
	if client.submitReq == nil {
		t.Fatal("expected submit request")
	}
	if client.submitReq.SignerPersonalNumber != "198301019876" {
		t.Fatalf("signer pnr = %q, want 198301019876", client.submitReq.SignerPersonalNumber)
	}
	decoded, err := base64.StdEncoding.DecodeString(client.submitReq.Document.File)
	if err != nil {
		t.Fatalf("decode base64 document: %v", err)
	}
	if !strings.Contains(string(decoded), `name="ixbrl.innehall.kontrollsumman" content="sum-123"`) {
		t.Fatal("expected checksum metadata in submitted document")
	}
	if len(client.submitReq.EmailAddresses) != 1 || client.submitReq.EmailAddresses[0] != "jag@foretag.com" {
		t.Fatalf("unexpected email addresses: %#v", client.submitReq.EmailAddresses)
	}
	if result.DocumentID != "49679" {
		t.Fatalf("document id = %q, want 49679", result.DocumentID)
	}
	if result.SubmissionURL != "https://example.test/submission/49679" {
		t.Fatalf("submission url = %q", result.SubmissionURL)
	}
	if result.Checksum != "sum-123" {
		t.Fatalf("checksum = %q, want sum-123", result.Checksum)
	}
	if result.ChecksumToken != "chk-123" {
		t.Fatalf("checksum token = %q, want chk-123", result.ChecksumToken)
	}
}

func TestServiceSubmitSkipCheck(t *testing.T) {
	report := loadTestReport(t)
	client := &fakeClient{
		checksumTokenResp: &CreateTokenResponse{Token: "chk-123"},
		checksumResp:      &CreateChecksumResponse{Checksum: "sum-123", Algorithm: "SHA-256"},
		createResp:        &CreateTokenResponse{Token: "tok-123"},
		submitResp:        &SubmitResponse{DocumentInfo: SubmitDocumentInfo{IDNumber: "49679"}},
	}

	_, err := NewService(client).Submit(context.Background(), report, SubmitOptions{SkipRemoteCheck: true})
	if err != nil {
		t.Fatalf("Submit returned error: %v", err)
	}
	if got, want := len(client.calls), 4; got != want {
		t.Fatalf("call count = %d, want %d", got, want)
	}
	if client.checkReq != nil {
		t.Fatal("did not expect check request when skip-check is enabled")
	}
}

func TestServiceSubmitValidationError(t *testing.T) {
	report := loadTestReport(t)
	report.Company.Name = ""

	_, err := NewService(&fakeClient{}).Submit(context.Background(), report, SubmitOptions{})
	if err == nil {
		t.Fatal("expected validation error")
	}
	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("error type = %T, want *ValidationError", err)
	}
	if !validate.HasErrors(validationErr.Findings) {
		t.Fatal("expected validation findings to include errors")
	}
}
