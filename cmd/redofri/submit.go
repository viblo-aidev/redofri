package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/redofri/redofri/pkg/model"
	"github.com/redofri/redofri/pkg/submission"
	"github.com/redofri/redofri/pkg/validate"
)

const (
	defaultSubmissionBaseURL = "http://127.0.0.1:8080"
	envSubmissionBaseURL     = "REDOFRI_SUBMISSION_BASE_URL"
	envSubmissionAPIKey      = "REDOFRI_SUBMISSION_API_KEY"
)

func runCheck(args []string) error {
	flags, err := parseSubmissionFlags(args)
	if err != nil {
		return err
	}
	if flags.inputPath == "" {
		return fmt.Errorf("missing input file\nUsage: redofri check [flags] <input.json>")
	}

	report, err := loadReport(flags.inputPath)
	if err != nil {
		return err
	}

	service, err := newSubmissionService(flags.baseURL, flags.apiKey)
	if err != nil {
		return err
	}

	result, err := service.Check(context.Background(), report, flags.submitOptions())
	if err != nil {
		if validationErr, ok := err.(*submission.ValidationError); ok {
			printLocalValidationReport(report, validationErr.Findings)
			return err
		}
		if remoteErr, ok := err.(*submission.RemoteCheckError); ok {
			printCheckSummary(report, flags.baseURL, result)
			printRemoteFindings(remoteErr.Findings)
			return err
		}
		return err
	}

	printCheckSummary(report, flags.baseURL, result)
	if len(result.RemoteFindings) > 0 {
		printRemoteFindings(result.RemoteFindings)
	} else {
		fmt.Println("Remote check passed: no errors or warnings.")
	}
	return nil
}

func runSubmit(args []string) error {
	flags, err := parseSubmissionFlags(args)
	if err != nil {
		return err
	}
	if flags.inputPath == "" {
		return fmt.Errorf("missing input file\nUsage: redofri submit [flags] <input.json>")
	}

	report, err := loadReport(flags.inputPath)
	if err != nil {
		return err
	}

	service, err := newSubmissionService(flags.baseURL, flags.apiKey)
	if err != nil {
		return err
	}

	result, err := service.Submit(context.Background(), report, flags.submitOptions())
	if err != nil {
		if validationErr, ok := err.(*submission.ValidationError); ok {
			printLocalValidationReport(report, validationErr.Findings)
			return err
		}
		if remoteErr, ok := err.(*submission.RemoteCheckError); ok {
			printSubmitSummary(report, flags.baseURL, result, flags.skipCheck)
			printRemoteFindings(remoteErr.Findings)
			return err
		}
		return err
	}

	printSubmitSummary(report, flags.baseURL, result, flags.skipCheck)
	if len(result.RemoteFindings) > 0 {
		printRemoteFindings(result.RemoteFindings)
	}
	fmt.Printf("Submission accepted: id=%s url=%s\n", result.DocumentID, result.SubmissionURL)
	return nil
}

type submissionFlags struct {
	inputPath                  string
	baseURL                    string
	apiKey                     string
	skipCheck                  bool
	senderPersonalNumber       string
	signerPersonalNumber       string
	emailAddresses             []string
	receiptEmailAddresses      []string
	notificationEmailAddresses []string
	documentType               string
}

func parseSubmissionFlags(args []string) (submissionFlags, error) {
	flags := submissionFlags{
		baseURL:              firstNonEmpty(os.Getenv(envSubmissionBaseURL), defaultSubmissionBaseURL),
		apiKey:               os.Getenv(envSubmissionAPIKey),
		senderPersonalNumber: os.Getenv("REDOFRI_SUBMISSION_SENDER_PNR"),
		signerPersonalNumber: os.Getenv("REDOFRI_SUBMISSION_SIGNER_PNR"),
		documentType:         submission.DefaultDocumentType,
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--base-url":
			i++
			if i >= len(args) {
				return flags, fmt.Errorf("--base-url requires a value")
			}
			flags.baseURL = args[i]
		case strings.HasPrefix(arg, "--base-url="):
			flags.baseURL = strings.TrimPrefix(arg, "--base-url=")
		case arg == "--api-key":
			i++
			if i >= len(args) {
				return flags, fmt.Errorf("--api-key requires a value")
			}
			flags.apiKey = args[i]
		case strings.HasPrefix(arg, "--api-key="):
			flags.apiKey = strings.TrimPrefix(arg, "--api-key=")
		case arg == "--skip-check":
			flags.skipCheck = true
		case arg == "--sender-pnr":
			i++
			if i >= len(args) {
				return flags, fmt.Errorf("--sender-pnr requires a value")
			}
			flags.senderPersonalNumber = args[i]
		case strings.HasPrefix(arg, "--sender-pnr="):
			flags.senderPersonalNumber = strings.TrimPrefix(arg, "--sender-pnr=")
		case arg == "--signer-pnr":
			i++
			if i >= len(args) {
				return flags, fmt.Errorf("--signer-pnr requires a value")
			}
			flags.signerPersonalNumber = args[i]
		case strings.HasPrefix(arg, "--signer-pnr="):
			flags.signerPersonalNumber = strings.TrimPrefix(arg, "--signer-pnr=")
		case arg == "--email":
			i++
			if i >= len(args) {
				return flags, fmt.Errorf("--email requires a value")
			}
			flags.emailAddresses = append(flags.emailAddresses, args[i])
		case strings.HasPrefix(arg, "--email="):
			flags.emailAddresses = append(flags.emailAddresses, strings.TrimPrefix(arg, "--email="))
		case arg == "--receipt-email":
			i++
			if i >= len(args) {
				return flags, fmt.Errorf("--receipt-email requires a value")
			}
			flags.receiptEmailAddresses = append(flags.receiptEmailAddresses, args[i])
		case strings.HasPrefix(arg, "--receipt-email="):
			flags.receiptEmailAddresses = append(flags.receiptEmailAddresses, strings.TrimPrefix(arg, "--receipt-email="))
		case arg == "--notify-email":
			i++
			if i >= len(args) {
				return flags, fmt.Errorf("--notify-email requires a value")
			}
			flags.notificationEmailAddresses = append(flags.notificationEmailAddresses, args[i])
		case strings.HasPrefix(arg, "--notify-email="):
			flags.notificationEmailAddresses = append(flags.notificationEmailAddresses, strings.TrimPrefix(arg, "--notify-email="))
		case arg == "--document-type":
			i++
			if i >= len(args) {
				return flags, fmt.Errorf("--document-type requires a value")
			}
			flags.documentType = args[i]
		case strings.HasPrefix(arg, "--document-type="):
			flags.documentType = strings.TrimPrefix(arg, "--document-type=")
		case strings.HasPrefix(arg, "-"):
			return flags, fmt.Errorf("unknown flag: %s", arg)
		default:
			flags.inputPath = arg
		}
	}

	return flags, nil
}

func (f submissionFlags) submitOptions() submission.SubmitOptions {
	return submission.SubmitOptions{
		SkipRemoteCheck:            f.skipCheck,
		SenderPersonalNumber:       f.senderPersonalNumber,
		SignerPersonalNumber:       f.signerPersonalNumber,
		EmailAddresses:             append([]string(nil), f.emailAddresses...),
		ReceiptEmailAddresses:      append([]string(nil), f.receiptEmailAddresses...),
		NotificationEmailAddresses: append([]string(nil), f.notificationEmailAddresses...),
		DocumentType:               f.documentType,
	}
}

func newSubmissionService(baseURL, apiKey string) (*submission.Service, error) {
	client, err := submission.NewHTTPClient(baseURL, nil, apiKey)
	if err != nil {
		return nil, err
	}
	return submission.NewService(client), nil
}

func printLocalValidationReport(report *model.AnnualReport, findings []validate.Result) {
	fmt.Printf("Company:      %s (%s)\n", report.Company.Name, report.Company.OrgNr)
	fmt.Printf("Fiscal year:  %s - %s\n\n", report.FiscalYear.StartDate, report.FiscalYear.EndDate)
	var errors, warnings int
	for _, finding := range findings {
		fmt.Println(finding)
		if finding.Severity == validate.Error {
			errors++
		} else {
			warnings++
		}
	}
	fmt.Printf("\n%d error(s), %d warning(s)\n", errors, warnings)
}

func printCheckSummary(report *model.AnnualReport, baseURL string, result *submission.CheckResult) {
	fmt.Printf("Company:      %s (%s)\n", report.Company.Name, report.Company.OrgNr)
	fmt.Printf("Fiscal year:  %s - %s\n", report.FiscalYear.StartDate, report.FiscalYear.EndDate)
	fmt.Printf("Base URL:     %s\n", baseURL)
	fmt.Printf("Sender PNR:   %s\n", result.SenderPersonalNumber)
	fmt.Printf("Token:        %s\n", result.Token)
	fmt.Printf("Checksum tk:  %s\n", result.ChecksumToken)
	fmt.Printf("Checksum:     %s\n", result.Checksum)
	if result.ChecksumAlgorithm != "" {
		fmt.Printf("Algorithm:    %s\n", result.ChecksumAlgorithm)
	}
	fmt.Printf("Document:     %d bytes (%s)\n", result.DocumentSize, result.DocumentType)
	if result.AgreementVersionDate != "" {
		fmt.Printf("Agreement:    %s\n", result.AgreementVersionDate)
	}
	fmt.Println()
	printLocalWarnings(result.LocalFindings)
}

func printSubmitSummary(report *model.AnnualReport, baseURL string, result *submission.SubmitResult, skipCheck bool) {
	fmt.Printf("Company:      %s (%s)\n", report.Company.Name, report.Company.OrgNr)
	fmt.Printf("Fiscal year:  %s - %s\n", report.FiscalYear.StartDate, report.FiscalYear.EndDate)
	fmt.Printf("Base URL:     %s\n", baseURL)
	fmt.Printf("Sender PNR:   %s\n", result.SenderPersonalNumber)
	fmt.Printf("Signer PNR:   %s\n", result.SignerPersonalNumber)
	fmt.Printf("Token:        %s\n", result.Token)
	fmt.Printf("Checksum tk:  %s\n", result.ChecksumToken)
	fmt.Printf("Checksum:     %s\n", result.Checksum)
	if result.ChecksumAlgorithm != "" {
		fmt.Printf("Algorithm:    %s\n", result.ChecksumAlgorithm)
	}
	fmt.Printf("Document:     %d bytes (%s)\n", result.DocumentSize, result.DocumentType)
	if skipCheck {
		fmt.Println("Remote check: skipped")
	}
	if result.AgreementVersionDate != "" {
		fmt.Printf("Agreement:    %s\n", result.AgreementVersionDate)
	}
	if result.DocumentID != "" {
		fmt.Printf("Document ID:  %s\n", result.DocumentID)
	}
	if result.SubmissionURL != "" {
		fmt.Printf("Submission:   %s\n", result.SubmissionURL)
	}
	fmt.Println()
	printLocalWarnings(result.LocalFindings)
}

func printRemoteFindings(findings []submission.Finding) {
	for _, finding := range findings {
		code := ""
		if finding.Code != "" {
			code = " [" + finding.Code + "]"
		}
		fmt.Printf("%s%s: %s\n", strings.ToUpper(string(finding.Severity)), code, finding.Message)
		for _, info := range finding.TechnicalInformation {
			fmt.Printf("  - %s=%s\n", info.Element, info.Value)
		}
	}
}

func printLocalWarnings(findings []validate.Result) {
	for _, finding := range findings {
		if finding.Severity == validate.Warning {
			fmt.Println(finding)
		}
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
