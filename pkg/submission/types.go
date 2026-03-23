package submission

import "context"

// Client defines the remote submission API used by the service layer.
type Client interface {
	CreateChecksumToken(ctx context.Context, req CreateTokenRequest) (*CreateTokenResponse, error)
	CreateChecksum(ctx context.Context, req CreateChecksumRequest) (*CreateChecksumResponse, error)
	CreateToken(ctx context.Context, req CreateTokenRequest) (*CreateTokenResponse, error)
	Check(ctx context.Context, req CheckRequest) (*CheckResponse, error)
	Submit(ctx context.Context, req SubmitRequest) (*SubmitResponse, error)
}

// DefaultDocumentType matches the documented Bolagsverket example for a complete annual report.
const DefaultDocumentType = "arsredovisning_komplett"

// CreateTokenRequest starts a new submission flow.
type CreateTokenRequest struct {
	SenderPersonalNumber string `json:"pnr"`
	OrgNumber            string `json:"orgnr"`
}

// CreateTokenResponse contains the token and agreement details required for later steps.
type CreateTokenResponse struct {
	Token                string `json:"token"`
	AgreementText        string `json:"avtalstext,omitempty"`
	AgreementVersionDate string `json:"avtalstextAndrad,omitempty"`
}

// Document contains the submitted annual report payload.
type Document struct {
	File string `json:"fil"`
	Type string `json:"typ"`
}

// CreateChecksumRequest creates a checksum for a base64-encoded document.
type CreateChecksumRequest struct {
	Token string `json:"-"`
	File  string `json:"fil"`
}

// CreateChecksumResponse contains the checksum produced by Bolagsverket.
type CreateChecksumResponse struct {
	Checksum  string `json:"kontrollsumma"`
	Algorithm string `json:"algoritm"`
}

// CheckRequest validates a generated document before submit.
type CheckRequest struct {
	Token    string   `json:"-"`
	Document Document `json:"handling"`
}

// TechnicalInformation mirrors tekniskinformation in the kontrollera response.
type TechnicalInformation struct {
	Message string `json:"meddelande,omitempty"`
	Element string `json:"element,omitempty"`
	Value   string `json:"varde,omitempty"`
}

// FindingSeverity indicates the severity returned by the remote submission API.
type FindingSeverity string

const (
	SeverityError FindingSeverity = "error"
	SeverityWarn  FindingSeverity = "warn"
)

// Finding represents a remote warning or error.
type Finding struct {
	Code                 string                 `json:"kod,omitempty"`
	Message              string                 `json:"text"`
	Severity             FindingSeverity        `json:"typ"`
	TechnicalInformation []TechnicalInformation `json:"tekniskinformation,omitempty"`
}

// CheckResponse is the remote validation result.
type CheckResponse struct {
	OrgNumber string    `json:"orgnr,omitempty"`
	Findings  []Finding `json:"utfall,omitempty"`
}

// SubmitRequest uploads the document.
type SubmitRequest struct {
	Token                   string   `json:"-"`
	SignerPersonalNumber    string   `json:"undertecknare"`
	EmailAddresses          []string `json:"epostadresser,omitempty"`
	ReceiptEmailAddresses   []string `json:"kvittensepostadresser,omitempty"`
	NotificationEmailAdress []string `json:"notifieringEpostadresser,omitempty"`
	Document                Document `json:"handling"`
}

// SubmitDocumentInfo mirrors handlingsinfo in the documented response.
type SubmitDocumentInfo struct {
	Type           string `json:"typ,omitempty"`
	DocumentLength int    `json:"dokumentlangd,omitempty"`
	IDNumber       string `json:"idnummer,omitempty"`
	Checksum       string `json:"sha256checksumma,omitempty"`
}

// SubmitResponse is the result of a successful upload.
type SubmitResponse struct {
	OrgNumber            string             `json:"orgnr,omitempty"`
	SenderPersonalNumber string             `json:"avsandare,omitempty"`
	SignerPersonalNumber string             `json:"undertecknare,omitempty"`
	DocumentInfo         SubmitDocumentInfo `json:"handlingsinfo,omitempty"`
	URL                  string             `json:"url,omitempty"`
}
