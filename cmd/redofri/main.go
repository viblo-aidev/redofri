// Command redofri generates Swedish K2 annual reports in iXBRL format.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/redofri/redofri/pkg/ixbrl"
	"github.com/redofri/redofri/pkg/model"
	"github.com/redofri/redofri/pkg/sie"
	"github.com/redofri/redofri/pkg/validate"
)

const version = "0.5.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "validate":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: redofri validate <input.json>\n")
			os.Exit(1)
		}
		if err := runValidate(os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "generate":
		if err := runGenerate(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "parse":
		if err := runParse(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "import-sie":
		if err := runImportSIE(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "version":
		fmt.Printf("redofri %s\n", version)

	case "help":
		printUsage()

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `redofri %s — Digital inlämning av svensk årsredovisning

Usage:
  redofri validate <input.json>         Load and validate JSON input
  redofri generate <input.json>         Generate iXBRL to stdout
  redofri generate -o <out> <input>     Generate iXBRL to file
  redofri parse <input.xhtml>           Parse iXBRL to JSON (stdout)
  redofri parse -o <out> <input>        Parse iXBRL to JSON file
  redofri import-sie <input.sie>        Import SIE4 to partial JSON (stdout)
  redofri import-sie -o <out> <input>   Import SIE4 to partial JSON file
  redofri version                       Show version
  redofri help                          Show this help

Flags (generate, parse, import-sie):
  -o, --output <file>   Write output to file (default: stdout)

Input can be a file path or "-" to read from stdin.
`, version)
}

// runGenerate parses flags, loads JSON input, and generates iXBRL output.
func runGenerate(args []string) error {
	inputPath, outputPath, err := parseIOFlags(args)
	if err != nil {
		return err
	}
	if inputPath == "" {
		return fmt.Errorf("missing input file\nUsage: redofri generate [-o output.xhtml] <input.json>")
	}

	report, err := loadReport(inputPath)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := ixbrl.Generate(&buf, report); err != nil {
		return fmt.Errorf("generating iXBRL: %w", err)
	}

	return writeOutput(outputPath, buf.Bytes(), "Generated")
}

// runParse reads an iXBRL file, parses it, and writes JSON output.
func runParse(args []string) error {
	inputPath, outputPath, err := parseIOFlags(args)
	if err != nil {
		return err
	}
	if inputPath == "" {
		return fmt.Errorf("missing input file\nUsage: redofri parse [-o output.json] <input.xhtml>")
	}

	var r io.Reader
	if inputPath == "-" {
		r = os.Stdin
	} else {
		f, err := os.Open(inputPath)
		if err != nil {
			return fmt.Errorf("reading %s: %w", inputPath, err)
		}
		defer f.Close()
		r = f
	}

	report, err := ixbrl.Parse(r)
	if err != nil {
		return fmt.Errorf("parsing iXBRL: %w", err)
	}

	out, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}
	out = append(out, '\n')

	return writeOutput(outputPath, out, "Parsed")
}

// parseIOFlags parses -o/--output and positional input path from args.
func parseIOFlags(args []string) (inputPath, outputPath string, err error) {
	i := 0
	for i < len(args) {
		switch {
		case args[i] == "-o" || args[i] == "--output":
			i++
			if i >= len(args) {
				return "", "", fmt.Errorf("-o/--output requires a file path argument")
			}
			outputPath = args[i]
		case strings.HasPrefix(args[i], "-o="):
			outputPath = args[i][3:]
		case strings.HasPrefix(args[i], "--output="):
			outputPath = args[i][9:]
		case strings.HasPrefix(args[i], "-"):
			if args[i] != "-" {
				return "", "", fmt.Errorf("unknown flag: %s", args[i])
			}
			inputPath = "-"
		default:
			inputPath = args[i]
		}
		i++
	}
	return inputPath, outputPath, nil
}

// writeOutput writes data to a file or stdout, printing a status line to stderr.
func writeOutput(outputPath string, data []byte, verb string) error {
	if outputPath == "" {
		_, err := os.Stdout.Write(data)
		return err
	}
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("writing %s: %w", outputPath, err)
	}
	fmt.Fprintf(os.Stderr, "%s %s (%d bytes)\n", verb, outputPath, len(data))
	return nil
}

// loadReport reads JSON from a file path or stdin ("-") and returns an AnnualReport.
func loadReport(path string) (*model.AnnualReport, error) {
	var data []byte
	var err error

	if path == "-" {
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("reading stdin: %w", err)
		}
	} else {
		data, err = os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}
	}

	var report model.AnnualReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	return &report, nil
}

// runValidate loads a JSON file, runs all validation checks, and prints findings.
// Exits with code 1 if there are errors.
func runValidate(path string) error {
	report, err := loadReport(path)
	if err != nil {
		return err
	}

	fmt.Printf("Company:      %s (%s)\n", report.Company.Name, report.Company.OrgNr)
	fmt.Printf("Fiscal year:  %s – %s\n", report.FiscalYear.StartDate, report.FiscalYear.EndDate)
	fmt.Printf("Entry point:  %s\n", report.Meta.EntryPoint)
	fmt.Println()

	results := validate.Validate(report)

	if len(results) == 0 {
		fmt.Println("Validation passed: no errors or warnings.")
		return nil
	}

	var errors, warnings int
	for _, r := range results {
		fmt.Println(r)
		if r.Severity == validate.Error {
			errors++
		} else {
			warnings++
		}
	}

	fmt.Printf("\n%d error(s), %d warning(s)\n", errors, warnings)

	if errors > 0 {
		return fmt.Errorf("validation failed with %d error(s)", errors)
	}
	return nil
}

// runImportSIE reads a SIE4 file, parses it, and writes a partial JSON report.
func runImportSIE(args []string) error {
	inputPath, outputPath, err := parseIOFlags(args)
	if err != nil {
		return err
	}
	if inputPath == "" {
		return fmt.Errorf("missing input file\nUsage: redofri import-sie [-o output.json] <input.sie>")
	}

	var r io.Reader
	if inputPath == "-" {
		r = os.Stdin
	} else {
		f, err := os.Open(inputPath)
		if err != nil {
			return fmt.Errorf("reading %s: %w", inputPath, err)
		}
		defer f.Close()
		r = f
	}

	result, err := sie.Parse(r)
	if err != nil {
		return fmt.Errorf("parsing SIE: %w", err)
	}

	for _, w := range result.Warnings {
		fmt.Fprintf(os.Stderr, "Warning: %s\n", w)
	}

	out, err := json.MarshalIndent(result.Report, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}
	out = append(out, '\n')

	return writeOutput(outputPath, out, "Imported")
}
