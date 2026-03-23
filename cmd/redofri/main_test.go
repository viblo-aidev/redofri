package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestGenerateCommand builds the binary and tests the generate command end-to-end.
func TestGenerateCommand(t *testing.T) {
	// Build the binary.
	tmpDir := t.TempDir()
	bin := filepath.Join(tmpDir, "redofri")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = "."
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}

	inputPath := filepath.Join("..", "..", "testdata", "exempel1.json")

	t.Run("generate to file", func(t *testing.T) {
		outPath := filepath.Join(tmpDir, "output.xhtml")
		cmd := exec.Command(bin, "generate", "-o", outPath, inputPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("generate failed: %v\n%s", err, out)
		}

		info, err := os.Stat(outPath)
		if err != nil {
			t.Fatalf("output file not found: %v", err)
		}
		if info.Size() < 50000 {
			t.Errorf("output too small: %d bytes", info.Size())
		}
	})

	t.Run("generate to stdout", func(t *testing.T) {
		cmd := exec.Command(bin, "generate", inputPath)
		out, err := cmd.Output()
		if err != nil {
			t.Fatalf("generate failed: %v", err)
		}
		if len(out) < 50000 {
			t.Errorf("output too small: %d bytes", len(out))
		}
		// Verify it starts with XML declaration.
		if string(out[:5]) != "<?xml" {
			t.Errorf("output does not start with XML declaration: %q", string(out[:20]))
		}
	})

	t.Run("generate with --output=", func(t *testing.T) {
		outPath := filepath.Join(tmpDir, "output2.xhtml")
		cmd := exec.Command(bin, "generate", "--output="+outPath, inputPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("generate failed: %v\n%s", err, out)
		}

		info, err := os.Stat(outPath)
		if err != nil {
			t.Fatalf("output file not found: %v", err)
		}
		if info.Size() < 50000 {
			t.Errorf("output too small: %d bytes", info.Size())
		}
	})

	t.Run("missing input", func(t *testing.T) {
		cmd := exec.Command(bin, "generate")
		if err := cmd.Run(); err == nil {
			t.Fatal("expected error for missing input")
		}
	})

	t.Run("validate command", func(t *testing.T) {
		cmd := exec.Command(bin, "validate", inputPath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("validate failed: %v\n%s", err, out)
		}
		if len(out) == 0 {
			t.Error("expected output from validate")
		}
	})

	t.Run("demo-generate command writes default file", func(t *testing.T) {
		workDir := t.TempDir()
		cmd := exec.Command(bin, "demo-generate")
		cmd.Dir = workDir

		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("demo-generate failed: %v\n%s", err, out)
		}

		if len(out) == 0 || string(out) == "" {
			t.Fatalf("expected status output from demo-generate")
		}

		data, err := os.ReadFile(filepath.Join(workDir, defaultDemoOutput))
		if err != nil {
			t.Fatalf("default demo xhtml file not found: %v", err)
		}
		if len(data) < 50000 {
			t.Fatalf("demo xhtml output too small: %d bytes", len(data))
		}
		if string(data[:5]) != "<?xml" {
			t.Fatalf("demo xhtml output missing XML declaration: %q", string(data[:20]))
		}
	})

	t.Run("demo-generate command writes custom file", func(t *testing.T) {
		outPath := filepath.Join(t.TempDir(), "demo.xhtml")
		cmd := exec.Command(bin, "demo-generate", "-o", outPath)

		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("demo-generate failed: %v\n%s", err, out)
		}

		data, err := os.ReadFile(outPath)
		if err != nil {
			t.Fatalf("custom demo xhtml file not found: %v", err)
		}
		if len(data) < 50000 {
			t.Fatalf("demo xhtml output too small: %d bytes", len(data))
		}
	})

	t.Run("demo-generate output can be parsed", func(t *testing.T) {
		xhtmlPath := filepath.Join(t.TempDir(), "demo.xhtml")
		genCmd := exec.Command(bin, "demo-generate", "-o", xhtmlPath)
		if out, err := genCmd.CombinedOutput(); err != nil {
			t.Fatalf("demo-generate: %v\n%s", err, out)
		}

		parseCmd := exec.Command(bin, "parse", xhtmlPath)
		out, err := parseCmd.Output()
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		if len(out) == 0 || out[0] != '{' {
			preview := string(out)
			if len(preview) > 50 {
				preview = preview[:50]
			}
			t.Fatalf("expected JSON output, got: %q", preview)
		}
	})

	t.Run("version command", func(t *testing.T) {
		cmd := exec.Command(bin, "version")
		out, err := cmd.Output()
		if err != nil {
			t.Fatalf("version failed: %v", err)
		}
		if string(out) != "redofri 0.5.0\n" {
			t.Errorf("unexpected version output: %q", string(out))
		}
	})

	t.Run("parse to stdout", func(t *testing.T) {
		// First generate an iXBRL file.
		xhtmlPath := filepath.Join(tmpDir, "roundtrip.xhtml")
		genCmd := exec.Command(bin, "generate", "-o", xhtmlPath, inputPath)
		if out, err := genCmd.CombinedOutput(); err != nil {
			t.Fatalf("generate for parse test: %v\n%s", err, out)
		}

		// Parse it back.
		parseCmd := exec.Command(bin, "parse", xhtmlPath)
		out, err := parseCmd.Output()
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		// Output should be JSON starting with '{'.
		if len(out) == 0 || out[0] != '{' {
			preview := string(out)
			if len(preview) > 50 {
				preview = preview[:50]
			}
			t.Errorf("expected JSON output, got: %q", preview)
		}
	})

	t.Run("parse to file", func(t *testing.T) {
		xhtmlPath := filepath.Join(tmpDir, "roundtrip2.xhtml")
		jsonPath := filepath.Join(tmpDir, "parsed.json")

		genCmd := exec.Command(bin, "generate", "-o", xhtmlPath, inputPath)
		if out, err := genCmd.CombinedOutput(); err != nil {
			t.Fatalf("generate: %v\n%s", err, out)
		}

		parseCmd := exec.Command(bin, "parse", "-o", jsonPath, xhtmlPath)
		if out, err := parseCmd.CombinedOutput(); err != nil {
			t.Fatalf("parse: %v\n%s", err, out)
		}

		info, err := os.Stat(jsonPath)
		if err != nil {
			t.Fatalf("output file not found: %v", err)
		}
		if info.Size() < 1000 {
			t.Errorf("output too small: %d bytes", info.Size())
		}
	})

	t.Run("check command", func(t *testing.T) {
		var paths []string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			paths = append(paths, r.URL.Path)
			switch r.URL.Path {
			case "/hamta-arsredovisningsinformation/v1.1/skapa-inlamningtoken":
				_ = json.NewEncoder(w).Encode(map[string]any{"token": "chk-123"})
			case "/hamta-arsredovisningsinformation/v1.1/skapa-kontrollsumma/chk-123":
				_ = json.NewEncoder(w).Encode(map[string]any{"kontrollsumma": "sum-123", "algoritm": "SHA-256"})
			case "/lamna-in-arsredovisning/v2.1/skapa-inlamningtoken/":
				_ = json.NewEncoder(w).Encode(map[string]any{
					"token":            "tok-123",
					"avtalstextAndrad": "2026-03-22",
				})
			case "/lamna-in-arsredovisning/v2.1/kontrollera/tok-123":
				_ = json.NewEncoder(w).Encode(map[string]any{"orgnr": "556000-1111"})
			default:
				http.NotFound(w, r)
			}
		}))
		defer server.Close()

		cmd := exec.Command(bin, "check", "--sender-pnr", "190001010106", inputPath)
		cmd.Env = append(os.Environ(), envSubmissionBaseURL+"="+server.URL)

		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("check failed: %v\n%s", err, out)
		}
		if got, want := strings.Join(paths, ","), "/hamta-arsredovisningsinformation/v1.1/skapa-inlamningtoken,/hamta-arsredovisningsinformation/v1.1/skapa-kontrollsumma/chk-123,/lamna-in-arsredovisning/v2.1/skapa-inlamningtoken/,/lamna-in-arsredovisning/v2.1/kontrollera/tok-123"; got != want {
			t.Fatalf("paths = %q, want %q", got, want)
		}
		if !strings.Contains(string(out), "Remote check passed") {
			t.Fatalf("expected success output, got:\n%s", out)
		}
	})

	t.Run("submit command", func(t *testing.T) {
		var paths []string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			paths = append(paths, r.URL.Path)
			switch r.URL.Path {
			case "/hamta-arsredovisningsinformation/v1.1/skapa-inlamningtoken":
				_ = json.NewEncoder(w).Encode(map[string]any{"token": "chk-123"})
			case "/hamta-arsredovisningsinformation/v1.1/skapa-kontrollsumma/chk-123":
				_ = json.NewEncoder(w).Encode(map[string]any{"kontrollsumma": "sum-123", "algoritm": "SHA-256"})
			case "/lamna-in-arsredovisning/v2.1/skapa-inlamningtoken/":
				_ = json.NewEncoder(w).Encode(map[string]any{
					"token":            "tok-123",
					"avtalstextAndrad": "2026-03-22",
				})
			case "/lamna-in-arsredovisning/v2.1/kontrollera/tok-123":
				_ = json.NewEncoder(w).Encode(map[string]any{"orgnr": "556000-1111"})
			case "/lamna-in-arsredovisning/v2.1/inlamning/tok-123":
				_, _ = io.WriteString(w, `{"orgnr":"556000-1111","avsandare":"190001010106","undertecknare":"190001010106","handlingsinfo":{"typ":"arsredovisning_komplett","dokumentlangd":123,"idnummer":"49679","sha256checksumma":"sum-123"},"url":"https://example.test/submission/49679"}`)
			default:
				http.NotFound(w, r)
			}
		}))
		defer server.Close()

		cmd := exec.Command(bin, "submit", "--sender-pnr", "190001010106", "--signer-pnr", "198301019876", inputPath)
		cmd.Env = append(os.Environ(), envSubmissionBaseURL+"="+server.URL)

		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("submit failed: %v\n%s", err, out)
		}
		if got, want := strings.Join(paths, ","), "/hamta-arsredovisningsinformation/v1.1/skapa-inlamningtoken,/hamta-arsredovisningsinformation/v1.1/skapa-kontrollsumma/chk-123,/lamna-in-arsredovisning/v2.1/skapa-inlamningtoken/,/lamna-in-arsredovisning/v2.1/kontrollera/tok-123,/lamna-in-arsredovisning/v2.1/inlamning/tok-123"; got != want {
			t.Fatalf("paths = %q, want %q", got, want)
		}
		if !strings.Contains(string(out), "Submission accepted: id=49679 url=https://example.test/submission/49679") {
			t.Fatalf("expected submit output, got:\n%s", out)
		}
	})
}
