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
)

const version = "0.2.0"

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
  redofri validate <input.json>       Load and validate JSON input
  redofri generate <input.json>       Generate iXBRL to stdout
  redofri generate -o <out> <input>   Generate iXBRL to file
  redofri version                     Show version
  redofri help                        Show this help

Generate flags:
  -o, --output <file>   Write iXBRL output to file (default: stdout)

Input can be a JSON file path or "-" to read from stdin.
`, version)
}

// runGenerate parses flags, loads JSON input, and generates iXBRL output.
func runGenerate(args []string) error {
	var outputPath string
	var inputPath string

	// Parse flags manually to avoid pulling in flag package complexities.
	i := 0
	for i < len(args) {
		switch {
		case args[i] == "-o" || args[i] == "--output":
			i++
			if i >= len(args) {
				return fmt.Errorf("-o/--output requires a file path argument")
			}
			outputPath = args[i]
		case strings.HasPrefix(args[i], "-o="):
			outputPath = args[i][3:]
		case strings.HasPrefix(args[i], "--output="):
			outputPath = args[i][9:]
		case strings.HasPrefix(args[i], "-"):
			if args[i] != "-" {
				return fmt.Errorf("unknown flag: %s", args[i])
			}
			inputPath = "-"
		default:
			inputPath = args[i]
		}
		i++
	}

	if inputPath == "" {
		return fmt.Errorf("missing input file\nUsage: redofri generate [-o output.xhtml] <input.json>")
	}

	// Load input.
	report, err := loadReport(inputPath)
	if err != nil {
		return err
	}

	// Generate iXBRL.
	var buf bytes.Buffer
	if err := ixbrl.Generate(&buf, report); err != nil {
		return fmt.Errorf("generating iXBRL: %w", err)
	}

	// Write output.
	if outputPath == "" {
		_, err = io.Copy(os.Stdout, &buf)
		return err
	}

	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", outputPath, err)
	}

	fmt.Fprintf(os.Stderr, "Generated %s (%d bytes)\n", outputPath, buf.Len())
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

// runValidate loads a JSON file, deserializes it, and reports basic stats.
// This will be extended with real validation in Step 6.
func runValidate(path string) error {
	report, err := loadReport(path)
	if err != nil {
		return err
	}

	fmt.Printf("Company:      %s (%s)\n", report.Company.Name, report.Company.OrgNr)
	fmt.Printf("Fiscal year:  %s – %s\n", report.FiscalYear.StartDate, report.FiscalYear.EndDate)
	fmt.Printf("Entry point:  %s\n", report.Meta.EntryPoint)

	if report.IncomeStatement.NetResult.Current != nil {
		fmt.Printf("Net result:   %d SEK\n", *report.IncomeStatement.NetResult.Current)
	}
	if report.BalanceSheet.Assets.TotalAssets.Current != nil {
		fmt.Printf("Total assets: %d SEK\n", *report.BalanceSheet.Assets.TotalAssets.Current)
	}

	fmt.Printf("Notes:        %d fixed asset notes\n", len(report.Notes.FixedAssetNotes))
	fmt.Printf("Signatories:  %d\n", len(report.Signatures.Signatories))

	fmt.Println("\nJSON loaded successfully.")
	return nil
}
