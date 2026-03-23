package submission

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HTTPClient implements Client over a JSON HTTP API.
type HTTPClient struct {
	baseURL string
	client  *http.Client
	apiKey  string
}

const checksumBasePath = "/hamta-arsredovisningsinformation/v1.1"
const submissionBasePath = "/lamna-in-arsredovisning/v2.1"

// NewHTTPClient creates a JSON submission client.
func NewHTTPClient(baseURL string, httpClient *http.Client, apiKey string) (*HTTPClient, error) {
	baseURL = strings.TrimRight(baseURL, "/")
	if baseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &HTTPClient{baseURL: baseURL, client: httpClient, apiKey: apiKey}, nil
}

// CreateToken starts a new submission flow.
func (c *HTTPClient) CreateToken(ctx context.Context, reqBody CreateTokenRequest) (*CreateTokenResponse, error) {
	var resp CreateTokenResponse
	if err := c.postJSON(ctx, submissionBasePath+"/skapa-inlamningtoken/", reqBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateChecksumToken starts the checksum flow documented under informationstjanster.
func (c *HTTPClient) CreateChecksumToken(ctx context.Context, reqBody CreateTokenRequest) (*CreateTokenResponse, error) {
	var resp CreateTokenResponse
	if err := c.postJSON(ctx, checksumBasePath+"/skapa-inlamningtoken", reqBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateChecksum asks the API to calculate the official checksum.
func (c *HTTPClient) CreateChecksum(ctx context.Context, reqBody CreateChecksumRequest) (*CreateChecksumResponse, error) {
	var resp CreateChecksumResponse
	if err := c.postJSON(ctx, checksumBasePath+"/skapa-kontrollsumma/"+reqBody.Token, reqBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Check performs the remote document check.
func (c *HTTPClient) Check(ctx context.Context, reqBody CheckRequest) (*CheckResponse, error) {
	var resp CheckResponse
	if err := c.postJSON(ctx, submissionBasePath+"/kontrollera/"+reqBody.Token, reqBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Submit uploads the document.
func (c *HTTPClient) Submit(ctx context.Context, reqBody SubmitRequest) (*SubmitResponse, error) {
	var resp SubmitResponse
	if err := c.postJSON(ctx, submissionBasePath+"/inlamning/"+reqBody.Token, reqBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *HTTPClient) postJSON(ctx context.Context, path string, reqBody, respBody any) error {
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		if len(body) == 0 {
			return fmt.Errorf("unexpected status %s", resp.Status)
		}
		return fmt.Errorf("unexpected status %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	if len(bytes.TrimSpace(body)) == 0 {
		return nil
	}
	if err := json.Unmarshal(body, respBody); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
