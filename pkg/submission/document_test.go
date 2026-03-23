package submission

import (
	"strings"
	"testing"
)

func TestInjectChecksumMetadata(t *testing.T) {
	doc := []byte("<html><head>\n<meta name=\"programvara\" content=\"Redofri\"/>\n</head><body/></html>")
	updated, err := injectChecksumMetadata(doc, "sum-123", "SHA-256")
	if err != nil {
		t.Fatalf("injectChecksumMetadata: %v", err)
	}
	out := string(updated)
	if !strings.Contains(out, `<meta name="ixbrl.innehall.kontrollsumman" content="sum-123"/>`) {
		t.Fatal("missing checksum meta tag")
	}
	if !strings.Contains(out, `<meta name="ixbrl.innehall.kontrollsumman.algoritm" content="SHA-256"/>`) {
		t.Fatal("missing checksum algorithm meta tag")
	}
	if strings.Index(out, `ixbrl.innehall.kontrollsumman`) > strings.Index(out, `</head>`) {
		t.Fatal("checksum meta inserted after </head>")
	}
}

func TestInjectChecksumMetadataIdempotent(t *testing.T) {
	doc := []byte("<html><head>\n<meta name=\"ixbrl.innehall.kontrollsumman\" content=\"sum-123\"/>\n</head><body/></html>")
	updated, err := injectChecksumMetadata(doc, "sum-123", "SHA-256")
	if err != nil {
		t.Fatalf("injectChecksumMetadata: %v", err)
	}
	if got := strings.Count(string(updated), `ixbrl.innehall.kontrollsumman`); got != 1 {
		t.Fatalf("checksum meta count = %d, want 1", got)
	}
}
