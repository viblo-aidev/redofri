// Command redofri generates Swedish K2 annual reports in iXBRL format.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/redofri/redofri/pkg/model"
)

const version = "0.1.0"

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
		fmt.Fprintf(os.Stderr, "generate: not yet implemented (Step 2)\n")
		os.Exit(1)

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
  redofri validate <input.json>   Load and validate JSON input
  redofri generate <input.json>   Generate iXBRL output (not yet implemented)
  redofri version                 Show version
  redofri help                    Show this help
`, version)
}

// runValidate loads a JSON file, deserializes it, and reports basic stats.
// This will be extended with real validation in Step 6.
func runValidate(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}

	var report model.AnnualReport
	if err := json.Unmarshal(data, &report); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
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
