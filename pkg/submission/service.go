package submission

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/redofri/redofri/pkg/ixbrl"
	"github.com/redofri/redofri/pkg/model"
	"github.com/redofri/redofri/pkg/validate"
)

// Service orchestrates local validation/generation and remote submission calls.
type Service struct {
	client Client
}

// SubmitOptions controls submit behavior.
type SubmitOptions struct {
	SkipRemoteCheck            bool
	SenderPersonalNumber       string
	SignerPersonalNumber       string
	EmailAddresses             []string
	ReceiptEmailAddresses      []string
	NotificationEmailAddresses []string
	DocumentType               string
}

// CheckResult contains the full result of a check flow.
type CheckResult struct {
	Token                string
	ChecksumToken        string
	AgreementText        string
	AgreementVersionDate string
	LocalFindings        []validate.Result
	RemoteFindings       []Finding
	Checksum             string
	ChecksumAlgorithm    string
	DocumentSize         int
	OrgNumber            string
	SenderPersonalNumber string
	DocumentType         string
}

// SubmitResult contains the full result of a submit flow.
type SubmitResult struct {
	Token                string
	ChecksumToken        string
	AgreementText        string
	AgreementVersionDate string
	LocalFindings        []validate.Result
	RemoteFindings       []Finding
	Checksum             string
	ChecksumAlgorithm    string
	DocumentSize         int
	OrgNumber            string
	SenderPersonalNumber string
	SignerPersonalNumber string
	DocumentType         string
	DocumentID           string
	SubmissionURL        string
}

// ValidationError indicates that local validation failed before remote submission.
type ValidationError struct {
	Findings []validate.Result
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("local validation failed with %d error(s)", countValidationErrors(e.Findings))
}

// RemoteCheckError indicates that the remote check step returned errors.
type RemoteCheckError struct {
	Findings []Finding
}

func (e *RemoteCheckError) Error() string {
	return fmt.Sprintf("remote check failed with %d error(s)", countRemoteErrors(e.Findings))
}

// NewService creates a submission service for the given client.
func NewService(client Client) *Service {
	return &Service{client: client}
}

// Check validates, generates, creates a token, and performs the remote check.
func (s *Service) Check(ctx context.Context, report *model.AnnualReport, opts SubmitOptions) (*CheckResult, error) {
	prepared, err := prepareReport(report, opts)
	if err != nil {
		return nil, err
	}

	checksumTokenResp, err := s.client.CreateChecksumToken(ctx, CreateTokenRequest{
		SenderPersonalNumber: prepared.senderPersonalNumber,
		OrgNumber:            report.Company.OrgNr,
	})
	if err != nil {
		return nil, fmt.Errorf("create checksum token: %w", err)
	}
	checksumResp, err := s.client.CreateChecksum(ctx, CreateChecksumRequest{
		Token: checksumTokenResp.Token,
		File:  prepared.documentBase64,
	})
	if err != nil {
		return nil, fmt.Errorf("create checksum: %w", err)
	}
	if err := prepared.applyChecksum(firstNonEmptyString(checksumResp.Checksum, prepared.checksum), checksumResp.Algorithm); err != nil {
		return nil, err
	}

	tokenResp, err := s.client.CreateToken(ctx, CreateTokenRequest{
		SenderPersonalNumber: prepared.senderPersonalNumber,
		OrgNumber:            report.Company.OrgNr,
	})
	if err != nil {
		return nil, fmt.Errorf("create token: %w", err)
	}

	checkResp, err := s.client.Check(ctx, CheckRequest{
		Token: tokenResp.Token,
		Document: Document{
			File: prepared.documentBase64,
			Type: prepared.documentType,
		},
	})
	result := &CheckResult{
		Token:                tokenResp.Token,
		ChecksumToken:        checksumTokenResp.Token,
		AgreementText:        tokenResp.AgreementText,
		AgreementVersionDate: tokenResp.AgreementVersionDate,
		LocalFindings:        prepared.findings,
		Checksum:             firstNonEmptyString(checksumResp.Checksum, prepared.checksum),
		ChecksumAlgorithm:    checksumResp.Algorithm,
		DocumentSize:         len(prepared.document),
		OrgNumber:            report.Company.OrgNr,
		SenderPersonalNumber: prepared.senderPersonalNumber,
		DocumentType:         prepared.documentType,
	}
	if err != nil {
		return result, fmt.Errorf("remote check: %w", err)
	}
	result.RemoteFindings = checkResp.Findings
	if hasRemoteErrors(checkResp.Findings) {
		return result, &RemoteCheckError{Findings: checkResp.Findings}
	}
	if checkResp.OrgNumber != "" {
		result.OrgNumber = checkResp.OrgNumber
	}
	return result, nil
}

// Submit validates, generates, creates a token, optionally checks, and submits.
func (s *Service) Submit(ctx context.Context, report *model.AnnualReport, opts SubmitOptions) (*SubmitResult, error) {
	prepared, err := prepareReport(report, opts)
	if err != nil {
		return nil, err
	}

	checksumTokenResp, err := s.client.CreateChecksumToken(ctx, CreateTokenRequest{
		SenderPersonalNumber: prepared.senderPersonalNumber,
		OrgNumber:            report.Company.OrgNr,
	})
	if err != nil {
		return nil, fmt.Errorf("create checksum token: %w", err)
	}
	checksumResp, err := s.client.CreateChecksum(ctx, CreateChecksumRequest{
		Token: checksumTokenResp.Token,
		File:  prepared.documentBase64,
	})
	if err != nil {
		return nil, fmt.Errorf("create checksum: %w", err)
	}
	if err := prepared.applyChecksum(firstNonEmptyString(checksumResp.Checksum, prepared.checksum), checksumResp.Algorithm); err != nil {
		return nil, err
	}

	tokenResp, err := s.client.CreateToken(ctx, CreateTokenRequest{
		SenderPersonalNumber: prepared.senderPersonalNumber,
		OrgNumber:            report.Company.OrgNr,
	})
	if err != nil {
		return nil, fmt.Errorf("create token: %w", err)
	}

	result := &SubmitResult{
		Token:                tokenResp.Token,
		ChecksumToken:        checksumTokenResp.Token,
		AgreementText:        tokenResp.AgreementText,
		AgreementVersionDate: tokenResp.AgreementVersionDate,
		LocalFindings:        prepared.findings,
		Checksum:             firstNonEmptyString(checksumResp.Checksum, prepared.checksum),
		ChecksumAlgorithm:    checksumResp.Algorithm,
		DocumentSize:         len(prepared.document),
		OrgNumber:            report.Company.OrgNr,
		SenderPersonalNumber: prepared.senderPersonalNumber,
		SignerPersonalNumber: prepared.signerPersonalNumber,
		DocumentType:         prepared.documentType,
	}

	if !opts.SkipRemoteCheck {
		checkResp, err := s.client.Check(ctx, CheckRequest{
			Token: tokenResp.Token,
			Document: Document{
				File: prepared.documentBase64,
				Type: prepared.documentType,
			},
		})
		if err != nil {
			return result, fmt.Errorf("remote check: %w", err)
		}
		result.RemoteFindings = checkResp.Findings
		if hasRemoteErrors(checkResp.Findings) {
			return result, &RemoteCheckError{Findings: checkResp.Findings}
		}
		if checkResp.OrgNumber != "" {
			result.OrgNumber = checkResp.OrgNumber
		}
	}

	submitResp, err := s.client.Submit(ctx, SubmitRequest{
		Token:                   tokenResp.Token,
		SignerPersonalNumber:    prepared.signerPersonalNumber,
		EmailAddresses:          prepared.emailAddresses,
		ReceiptEmailAddresses:   prepared.receiptEmailAddresses,
		NotificationEmailAdress: prepared.notificationEmailAddresses,
		Document: Document{
			File: prepared.documentBase64,
			Type: prepared.documentType,
		},
	})
	if err != nil {
		return result, fmt.Errorf("submit document: %w", err)
	}

	if submitResp.OrgNumber != "" {
		result.OrgNumber = submitResp.OrgNumber
	}
	if submitResp.SenderPersonalNumber != "" {
		result.SenderPersonalNumber = submitResp.SenderPersonalNumber
	}
	if submitResp.SignerPersonalNumber != "" {
		result.SignerPersonalNumber = submitResp.SignerPersonalNumber
	}
	if submitResp.DocumentInfo.Type != "" {
		result.DocumentType = submitResp.DocumentInfo.Type
	}
	result.DocumentID = submitResp.DocumentInfo.IDNumber
	if submitResp.DocumentInfo.Checksum != "" {
		result.Checksum = submitResp.DocumentInfo.Checksum
	}
	result.SubmissionURL = submitResp.URL

	return result, nil
}

type preparedReport struct {
	findings                   []validate.Result
	document                   []byte
	documentBase64             string
	checksum                   string
	senderPersonalNumber       string
	signerPersonalNumber       string
	emailAddresses             []string
	receiptEmailAddresses      []string
	notificationEmailAddresses []string
	documentType               string
}

func (p *preparedReport) applyChecksum(checksumValue, algorithm string) error {
	if checksumValue == "" || algorithm == "" {
		return nil
	}
	updated, err := injectChecksumMetadata(p.document, checksumValue, algorithm)
	if err != nil {
		return fmt.Errorf("inject checksum metadata: %w", err)
	}
	p.document = updated
	p.documentBase64 = base64.StdEncoding.EncodeToString(updated)
	return nil
}

func prepareReport(report *model.AnnualReport, opts SubmitOptions) (*preparedReport, error) {
	findings := validate.Validate(report)
	if validate.HasErrors(findings) {
		return nil, &ValidationError{Findings: findings}
	}

	document, err := ixbrl.GenerateBytes(report)
	if err != nil {
		return nil, fmt.Errorf("generate iXBRL: %w", err)
	}

	senderPersonalNumber := defaultPersonalNumber(opts.SenderPersonalNumber)
	signerPersonalNumber := defaultPersonalNumber(opts.SignerPersonalNumber)
	documentType := opts.DocumentType
	if documentType == "" {
		documentType = DefaultDocumentType
	}

	return &preparedReport{
		findings:                   findings,
		document:                   document,
		documentBase64:             base64.StdEncoding.EncodeToString(document),
		checksum:                   checksum(document),
		senderPersonalNumber:       senderPersonalNumber,
		signerPersonalNumber:       signerPersonalNumber,
		emailAddresses:             append([]string(nil), opts.EmailAddresses...),
		receiptEmailAddresses:      append([]string(nil), opts.ReceiptEmailAddresses...),
		notificationEmailAddresses: append([]string(nil), opts.NotificationEmailAddresses...),
		documentType:               documentType,
	}, nil
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func checksum(document []byte) string {
	sum := sha256.Sum256(document)
	return base64.StdEncoding.EncodeToString(sum[:])
}

func defaultPersonalNumber(value string) string {
	if value != "" {
		return value
	}
	return "190001010106"
}

func hasRemoteErrors(findings []Finding) bool {
	for _, finding := range findings {
		if finding.Severity == SeverityError {
			return true
		}
	}
	return false
}

func countRemoteErrors(findings []Finding) int {
	count := 0
	for _, finding := range findings {
		if finding.Severity == SeverityError {
			count++
		}
	}
	return count
}

func countValidationErrors(findings []validate.Result) int {
	count := 0
	for _, finding := range findings {
		if finding.Severity == validate.Error {
			count++
		}
	}
	return count
}
