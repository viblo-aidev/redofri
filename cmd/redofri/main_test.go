package main

import (
	"os"
	"os/exec"
	"path/filepath"
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

	t.Run("version command", func(t *testing.T) {
		cmd := exec.Command(bin, "version")
		out, err := cmd.Output()
		if err != nil {
			t.Fatalf("version failed: %v", err)
		}
		if string(out) != "redofri 0.3.0\n" {
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
}
