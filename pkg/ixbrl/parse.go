// Package ixbrl provides iXBRL generation and parsing for K2 annual reports.
//
// Parse reads an iXBRL (.xhtml) document and populates a model.AnnualReport
// by extracting ix:nonFraction, ix:nonNumeric, and ix:tuple elements.
package ixbrl

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/redofri/redofri/pkg/model"
)

// Parse reads an iXBRL document from r and returns a populated AnnualReport.
func Parse(r io.Reader) (*model.AnnualReport, error) {
	facts, err := extractFacts(r)
	if err != nil {
		return nil, fmt.Errorf("extracting facts: %w", err)
	}
	return mapFacts(facts)
}

// ---------- fact extraction ----------

// fact represents a single extracted XBRL fact.
type fact struct {
	// Element type: "nonFraction", "nonNumeric", "tuple"
	Kind string

	// XBRL concept name including namespace prefix, e.g. "se-gen-base:Nettoomsattning"
	Name string

	// Context reference, e.g. "period0", "balans0"
	ContextRef string

	// Unit reference, e.g. "SEK", "procent", "antal-anstallda"
	UnitRef string

	// Text content (for nonNumeric) or numeric string (for nonFraction)
	Value string

	// Numeric attributes
	Decimals string
	Scale    int
	Sign     string // "-" if negative
	Format   string // e.g. "ixt:numspacecomma"

	// Tuple support
	TupleID  string // id for ix:tuple elements
	TupleRef string // reference to parent tuple
	Order    string // ordering within tuple

	// Special attributes
	ID string // e.g. "ID_DATUM_UNDERTECKNANDE_FASTSTALLELSEINTYG"
}

// ixNS is the iXBRL namespace URI.
const ixNS = "http://www.xbrl.org/2013/inlineXBRL"

// extractFacts parses the iXBRL XML and extracts all fact elements.
// It handles nested ix:nonNumeric elements (e.g. in the certification section)
// and ix:continuation elements by recursively extracting inner facts.
func extractFacts(r io.Reader) ([]fact, error) {
	decoder := xml.NewDecoder(r)

	var facts []fact

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading XML: %w", err)
		}

		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Space != ixNS {
				continue
			}
			extracted, err := parseIXElement(decoder, t)
			if err != nil {
				return nil, err
			}
			facts = append(facts, extracted...)
		}
	}

	return facts, nil
}

// parseIXElement dispatches to the correct parser for an ix: element
// and returns all facts found (including any nested ones).
func parseIXElement(decoder *xml.Decoder, start xml.StartElement) ([]fact, error) {
	switch start.Name.Local {
	case "nonFraction":
		f, err := parseNonFraction(decoder, start)
		if err != nil {
			return nil, err
		}
		return []fact{f}, nil

	case "nonNumeric":
		return parseNonNumericRecursive(decoder, start)

	case "tuple":
		return []fact{parseTuple(start)}, nil

	case "continuation":
		// Parse contents of continuation — may contain ix: elements.
		return parseContainerContents(decoder, start.Name)

	default:
		return nil, nil
	}
}

// getAttr returns the value of an attribute by local name.
func getAttr(attrs []xml.Attr, local string) string {
	for _, a := range attrs {
		if a.Name.Local == local {
			return a.Value
		}
	}
	return ""
}

// parseNonFraction extracts a fact from an ix:nonFraction element.
func parseNonFraction(decoder *xml.Decoder, start xml.StartElement) (fact, error) {
	scaleStr := getAttr(start.Attr, "scale")
	scale := 0
	if scaleStr != "" {
		s, err := strconv.Atoi(scaleStr)
		if err != nil {
			return fact{}, fmt.Errorf("invalid scale %q: %w", scaleStr, err)
		}
		scale = s
	}

	// Read inner text content, skipping any nested elements.
	value := collectText(decoder, start.Name)

	return fact{
		Kind:       "nonFraction",
		Name:       getAttr(start.Attr, "name"),
		ContextRef: getAttr(start.Attr, "contextRef"),
		UnitRef:    getAttr(start.Attr, "unitRef"),
		Value:      strings.TrimSpace(value),
		Decimals:   getAttr(start.Attr, "decimals"),
		Scale:      scale,
		Sign:       getAttr(start.Attr, "sign"),
		Format:     getAttr(start.Attr, "format"),
		TupleRef:   getAttr(start.Attr, "tupleRef"),
		Order:      getAttr(start.Attr, "order"),
		ID:         getAttr(start.Attr, "id"),
	}, nil
}

// parseNonNumericRecursive extracts a nonNumeric fact and any nested ix: facts.
// It returns the outer fact plus all inner facts in order.
func parseNonNumericRecursive(decoder *xml.Decoder, start xml.StartElement) ([]fact, error) {
	var result []fact
	var textParts []string
	depth := 1

	for depth > 0 {
		tok, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("reading nonNumeric %s: %w", getAttr(start.Attr, "name"), err)
		}

		switch t := tok.(type) {
		case xml.CharData:
			textParts = append(textParts, string(t))

		case xml.StartElement:
			depth++
			if t.Name.Space == ixNS {
				// Recursively parse nested ix: element.
				inner, err := parseIXElement(decoder, t)
				if err != nil {
					return nil, err
				}
				result = append(result, inner...)
				depth-- // parseIXElement consumed the end element.
			}

		case xml.EndElement:
			depth--
		}
	}

	// Build the outer fact.
	outer := fact{
		Kind:       "nonNumeric",
		Name:       getAttr(start.Attr, "name"),
		ContextRef: getAttr(start.Attr, "contextRef"),
		Value:      strings.TrimSpace(strings.Join(textParts, "")),
		TupleRef:   getAttr(start.Attr, "tupleRef"),
		Order:      getAttr(start.Attr, "order"),
		ID:         getAttr(start.Attr, "id"),
	}

	// Prepend the outer fact before any inner facts.
	return append([]fact{outer}, result...), nil
}

// parseContainerContents parses the contents of an ix:continuation or similar
// container element, returning any nested ix: facts found within.
func parseContainerContents(decoder *xml.Decoder, name xml.Name) ([]fact, error) {
	var result []fact
	depth := 1

	for depth > 0 {
		tok, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("reading container: %w", err)
		}

		switch t := tok.(type) {
		case xml.StartElement:
			depth++
			if t.Name.Space == ixNS {
				inner, err := parseIXElement(decoder, t)
				if err != nil {
					return nil, err
				}
				result = append(result, inner...)
				depth-- // parseIXElement consumed the end element.
			}

		case xml.EndElement:
			depth--
		}
	}

	return result, nil
}

// parseTuple extracts a tuple declaration (self-closing element).
func parseTuple(start xml.StartElement) fact {
	return fact{
		Kind:    "tuple",
		Name:    getAttr(start.Attr, "name"),
		TupleID: getAttr(start.Attr, "tupleID"),
	}
}

// collectText reads tokens until the matching end element, collecting all
// character data. Nested elements are consumed but their text is included.
func collectText(decoder *xml.Decoder, name xml.Name) string {
	var sb strings.Builder
	depth := 1
	for depth > 0 {
		tok, err := decoder.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.CharData:
			sb.Write(t)
		case xml.StartElement:
			depth++
		case xml.EndElement:
			depth--
		}
	}
	return sb.String()
}

// ---------- number parsing ----------

// parseNumber parses the display string of a nonFraction fact into an int64
// XBRL value, applying format, scale, and sign transformations.
func parseNumber(f fact) (int64, error) {
	s := f.Value

	// Remove display formatting: spaces and any non-breaking spaces.
	s = strings.ReplaceAll(s, "\u00a0", "")
	s = strings.ReplaceAll(s, " ", "")

	// Handle format-specific decimal separator.
	switch f.Format {
	case "ixt:numspacecomma", "ixt:numcomma":
		// Swedish format: comma is decimal separator.
		s = strings.ReplaceAll(s, ",", ".")
	}

	// Remove any remaining formatting chars (e.g. leading minus that's in the display).
	s = strings.TrimLeft(s, "-+")
	if s == "" {
		return 0, nil
	}

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing number %q (from %q): %w", s, f.Value, err)
	}

	// Apply scale: value * 10^scale gives the XBRL value.
	if f.Scale != 0 {
		val *= math.Pow(10, float64(f.Scale))
	}

	// Round to nearest integer.
	result := int64(math.Round(val))

	// Apply sign.
	if f.Sign == "-" {
		result = -result
	}

	return result, nil
}
