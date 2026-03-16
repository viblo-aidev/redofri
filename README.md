# Redofri

Command-line tool for generating Swedish annual reports (årsredovisning) in iXBRL format, ready for digital submission to Bolagsverket.

First version targets **K2 for aktiebolag (AB) with fastställelseintyg**.

## Features

- **iXBRL generation** -- produces a self-contained `.xhtml` file that is both human-readable in a browser and machine-readable XBRL
- **iXBRL parsing** -- roundtrip: parse an existing iXBRL annual report back to the internal model (useful for extracting comparative figures from last year)
- **SIE4 import** -- import account balances from SIE4 files with automatic BAS account mapping
- **Validation** -- checks required fields, calculation consistency, date ordering, and Bolagsverket validation codes (1019--3007)
- **Cross-platform** -- builds for Linux, macOS, and Windows

## Installation

Requires Go 1.25+.

```
go install github.com/redofri/redofri/cmd/redofri@latest
```

Or build from source:

```
git clone https://github.com/viblo-aidev/redofri.git
cd redofri
go build -o redofri ./cmd/redofri
```

## Usage

```
redofri demo-generate                  # Generate a demo iXBRL file
redofri generate <input.json>           # Generate iXBRL to stdout
redofri generate -o out.xhtml input.json  # Generate iXBRL to file
redofri validate <input.json>           # Validate a report
redofri parse <input.xhtml>             # Parse iXBRL back to JSON
redofri import-sie <input.sie>          # Import SIE4 to partial JSON
redofri version                         # Show version
redofri help                            # Show help
```

All commands accept `-o <file>` to write output to a file instead of stdout. Input can be `-` to read from stdin.

If you want to try the tool immediately without preparing any input data first:

```
redofri demo-generate -o demo.xhtml
```

### Typical workflow

```
# 1. Import account balances from your SIE4 file
redofri import-sie -o partial.json bookkeeping.sie

# 2. Complete the JSON with management report, notes, signatures, etc.
#    (edit partial.json or merge with a template)

# 3. Validate before generating
redofri validate report.json

# 4. Generate the iXBRL file
redofri generate -o arsredovisning.xhtml report.json
```

## Data model

The central contract is `pkg/model/model.go` -- a set of Go structs representing a complete K2 annual report. All data sources (SIE import, iXBRL parsing, JSON input) populate these structs, and the iXBRL generator reads them to produce the output.

```
SIE file ──────────parse──┐
Previous annual report ────┼──► model.AnnualReport ──► ixbrl.Generate() ──► .xhtml
Manual/JSON input ─────────┘
```

## Project structure

```
cmd/redofri/       CLI entry point
pkg/model/         Data model (Go structs)
pkg/ixbrl/         iXBRL generator and parser
pkg/sie/           SIE4 parser
pkg/validate/      Validation engine
testdata/          Test fixtures
ref/               Reference material (taxonomy, technical guide)
```

## Testing

```
go test ./...
```

The test suite includes roundtrip tests (generate iXBRL, parse it back, verify the model matches), calculation checks, and CLI integration tests.

## Status

| Step | Description | Status |
|------|-------------|--------|
| 1 | Data model | Done |
| 2 | iXBRL generation | Done |
| 3 | CLI | Done |
| 4 | iXBRL parser (roundtrip) | Done |
| 5 | SIE import | Done |
| 6 | Validation | Done |
| 7 | API integration (Bolagsverket) | Planned |
| 8 | Web UI | Planned |

See [PLAN.md](PLAN.md) for details.

## License

AGPL-3.0. See [LICENSE](LICENSE).
