package submission

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	checksumMetaName          = "ixbrl.innehall.kontrollsumman"
	checksumAlgorithmMetaName = "ixbrl.innehall.kontrollsumman.algoritm"
)

func injectChecksumMetadata(document []byte, checksum, algorithm string) ([]byte, error) {
	if checksum == "" || algorithm == "" {
		return nil, fmt.Errorf("checksum and algorithm are required")
	}
	if bytes.Contains(document, []byte(`name="`+checksumMetaName+`"`)) {
		return document, nil
	}

	marker := []byte("</head>")
	idx := bytes.Index(document, marker)
	if idx < 0 {
		return nil, fmt.Errorf("missing </head> in generated document")
	}

	meta := []byte(
		"		<meta name=\"" + checksumMetaName + "\" content=\"" + xmlEscape(checksum) + "\"/>\n" +
			"		<meta name=\"" + checksumAlgorithmMetaName + "\" content=\"" + xmlEscape(algorithm) + "\"/>\n",
	)

	updated := make([]byte, 0, len(document)+len(meta))
	updated = append(updated, document[:idx]...)
	updated = append(updated, meta...)
	updated = append(updated, document[idx:]...)
	return updated, nil
}

func xmlEscape(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
	)
	return replacer.Replace(s)
}
