package ixbrl

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/redofri/redofri/pkg/model"
)

// TestParseRoundtrip generates iXBRL from the test model, parses it back,
// and verifies the parsed model matches the original.
func TestParseRoundtrip(t *testing.T) {
	original := loadTestReport(t)

	// Generate iXBRL.
	var buf bytes.Buffer
	if err := Generate(&buf, original); err != nil {
		t.Fatalf("generating iXBRL: %v", err)
	}

	// Parse it back.
	parsed, err := Parse(&buf)
	if err != nil {
		t.Fatalf("parsing iXBRL: %v", err)
	}

	// --- Metadata ---
	t.Run("company", func(t *testing.T) {
		assertEqual(t, "name", original.Company.Name, parsed.Company.Name)
		assertEqual(t, "orgNr", original.Company.OrgNr, parsed.Company.OrgNr)
	})

	t.Run("fiscal year", func(t *testing.T) {
		assertEqual(t, "startDate", original.FiscalYear.StartDate, parsed.FiscalYear.StartDate)
		assertEqual(t, "endDate", original.FiscalYear.EndDate, parsed.FiscalYear.EndDate)
	})

	t.Run("meta", func(t *testing.T) {
		assertEqual(t, "language", original.Meta.Language, parsed.Meta.Language)
		assertEqual(t, "country", original.Meta.Country, parsed.Meta.Country)
		assertEqual(t, "currency", original.Meta.Currency, parsed.Meta.Currency)
		assertEqual(t, "amountFormat", original.Meta.AmountFormat, parsed.Meta.AmountFormat)
	})

	// --- Certification ---
	t.Run("certification", func(t *testing.T) {
		assertEqual(t, "confirmationText", original.Certification.ConfirmationText, parsed.Certification.ConfirmationText)
		assertEqual(t, "meetingDate", original.Certification.MeetingDate, parsed.Certification.MeetingDate)
		assertEqual(t, "dispositionDecision", original.Certification.DispositionDecision, parsed.Certification.DispositionDecision)
		assertEqual(t, "signingDate", original.Certification.SigningDate, parsed.Certification.SigningDate)
		assertEqual(t, "signatory.firstName", original.Certification.Signatory.FirstName, parsed.Certification.Signatory.FirstName)
		assertEqual(t, "signatory.lastName", original.Certification.Signatory.LastName, parsed.Certification.Signatory.LastName)
	})

	// --- Management Report ---
	t.Run("management report text", func(t *testing.T) {
		assertNonEmpty(t, "introText", parsed.ManagementReport.IntroText)
		assertNonEmpty(t, "businessDescription", parsed.ManagementReport.BusinessDescription)
		assertNonEmpty(t, "significantEvents", parsed.ManagementReport.SignificantEvents)
	})

	t.Run("equity changes", func(t *testing.T) {
		oec := original.ManagementReport.EquityChanges
		pec := parsed.ManagementReport.EquityChanges

		assertInt64PtrEqual(t, "openingShareCapital", oec.OpeningShareCapital, pec.OpeningShareCapital)
		assertInt64PtrEqual(t, "openingReserveFund", oec.OpeningReserveFund, pec.OpeningReserveFund)
		assertInt64PtrEqual(t, "openingRetainedEarnings", oec.OpeningRetainedEarnings, pec.OpeningRetainedEarnings)
		assertInt64PtrEqual(t, "openingNetIncome", oec.OpeningNetIncome, pec.OpeningNetIncome)
		assertInt64PtrEqual(t, "openingTotal", oec.OpeningTotal, pec.OpeningTotal)
		assertInt64PtrEqual(t, "closingTotal", oec.ClosingTotal, pec.ClosingTotal)
		assertInt64PtrEqual(t, "yearResultNetIncome", oec.YearResultNetIncome, pec.YearResultNetIncome)
	})

	t.Run("profit disposition", func(t *testing.T) {
		opd := original.ManagementReport.ProfitDisposition
		ppd := parsed.ManagementReport.ProfitDisposition

		assertInt64PtrEqual(t, "retainedEarnings", opd.RetainedEarnings, ppd.RetainedEarnings)
		assertInt64PtrEqual(t, "netIncome", opd.NetIncome, ppd.NetIncome)
		assertInt64PtrEqual(t, "totalAvailable", opd.TotalAvailable, ppd.TotalAvailable)
		assertInt64PtrEqual(t, "carriedForward", opd.CarriedForward, ppd.CarriedForward)
		assertInt64PtrEqual(t, "totalDisposition", opd.TotalDisposition, ppd.TotalDisposition)
	})

	// --- Income Statement ---
	t.Run("income statement", func(t *testing.T) {
		ois := original.IncomeStatement
		pis := parsed.IncomeStatement

		assertYCEqual(t, "netSales", ois.Revenue.NetSales, pis.Revenue.NetSales)
		assertYCEqual(t, "otherOperatingIncome", ois.Revenue.OtherOperatingIncome, pis.Revenue.OtherOperatingIncome)
		assertYCEqual(t, "totalRevenue", ois.Revenue.TotalRevenue, pis.Revenue.TotalRevenue)
		assertYCEqual(t, "rawMaterials", ois.Expenses.RawMaterials, pis.Expenses.RawMaterials)
		assertYCEqual(t, "tradingGoods", ois.Expenses.TradingGoods, pis.Expenses.TradingGoods)
		assertYCEqual(t, "otherExternalExpenses", ois.Expenses.OtherExternalExpenses, pis.Expenses.OtherExternalExpenses)
		assertYCEqual(t, "personnelExpenses", ois.Expenses.PersonnelExpenses, pis.Expenses.PersonnelExpenses)
		assertYCEqual(t, "depreciationAmortization", ois.Expenses.DepreciationAmortization, pis.Expenses.DepreciationAmortization)
		assertYCEqual(t, "totalExpenses", ois.Expenses.TotalExpenses, pis.Expenses.TotalExpenses)
		assertYCEqual(t, "operatingResult", ois.OperatingResult, pis.OperatingResult)
		assertYCEqual(t, "interestExpenses", ois.FinancialItems.InterestExpenses, pis.FinancialItems.InterestExpenses)
		assertYCEqual(t, "totalFinancialItems", ois.FinancialItems.TotalFinancialItems, pis.FinancialItems.TotalFinancialItems)
		assertYCEqual(t, "resultAfterFinancialItems", ois.ResultAfterFinancialItems, pis.ResultAfterFinancialItems)
		assertYCEqual(t, "taxAllocationReserve", ois.Appropriations.TaxAllocationReserve, pis.Appropriations.TaxAllocationReserve)
		assertYCEqual(t, "excessDepreciation", ois.Appropriations.ExcessDepreciation, pis.Appropriations.ExcessDepreciation)
		assertYCEqual(t, "totalAppropriations", ois.Appropriations.TotalAppropriations, pis.Appropriations.TotalAppropriations)
		assertYCEqual(t, "resultBeforeTax", ois.ResultBeforeTax, pis.ResultBeforeTax)
		assertYCEqual(t, "incomeTax", ois.Tax.IncomeTax, pis.Tax.IncomeTax)
		assertYCEqual(t, "netResult", ois.NetResult, pis.NetResult)
	})

	// --- Balance Sheet ---
	t.Run("balance sheet assets", func(t *testing.T) {
		oa := original.BalanceSheet.Assets
		pa := parsed.BalanceSheet.Assets

		assertYCEqual(t, "buildingsAndLand", oa.FixedAssets.Tangible.BuildingsAndLand, pa.FixedAssets.Tangible.BuildingsAndLand)
		assertYCEqual(t, "machineryAndEquipment", oa.FixedAssets.Tangible.MachineryAndEquipment, pa.FixedAssets.Tangible.MachineryAndEquipment)
		assertYCEqual(t, "fixturesAndFittings", oa.FixedAssets.Tangible.FixturesAndFittings, pa.FixedAssets.Tangible.FixturesAndFittings)
		assertYCEqual(t, "totalTangible", oa.FixedAssets.Tangible.TotalTangible, pa.FixedAssets.Tangible.TotalTangible)
		assertYCEqual(t, "otherLongTermSecurities", oa.FixedAssets.Financial.OtherLongTermSecurities, pa.FixedAssets.Financial.OtherLongTermSecurities)
		assertYCEqual(t, "totalFinancial", oa.FixedAssets.Financial.TotalFinancial, pa.FixedAssets.Financial.TotalFinancial)
		assertYCEqual(t, "totalFixedAssets", oa.FixedAssets.TotalFixedAssets, pa.FixedAssets.TotalFixedAssets)
		assertYCEqual(t, "totalInventory", oa.CurrentAssets.Inventory.TotalInventory, pa.CurrentAssets.Inventory.TotalInventory)
		assertYCEqual(t, "tradeReceivables", oa.CurrentAssets.ShortTermReceivables.TradeReceivables, pa.CurrentAssets.ShortTermReceivables.TradeReceivables)
		assertYCEqual(t, "totalShortTermReceivables", oa.CurrentAssets.ShortTermReceivables.TotalShortTermReceivables, pa.CurrentAssets.ShortTermReceivables.TotalShortTermReceivables)
		assertYCEqual(t, "totalCashAndBank", oa.CurrentAssets.CashAndBank.TotalCashAndBank, pa.CurrentAssets.CashAndBank.TotalCashAndBank)
		assertYCEqual(t, "totalCurrentAssets", oa.CurrentAssets.TotalCurrentAssets, pa.CurrentAssets.TotalCurrentAssets)
		assertYCEqual(t, "totalAssets", oa.TotalAssets, pa.TotalAssets)
	})

	t.Run("balance sheet equity and liabilities", func(t *testing.T) {
		oel := original.BalanceSheet.EquityAndLiabilities
		pel := parsed.BalanceSheet.EquityAndLiabilities

		assertYCEqual(t, "shareCapital", oel.Equity.ShareCapital, pel.Equity.ShareCapital)
		assertYCEqual(t, "totalEquity", oel.Equity.TotalEquity, pel.Equity.TotalEquity)
		assertYCEqual(t, "totalUntaxedReserves", oel.UntaxedReserves.TotalUntaxedReserves, pel.UntaxedReserves.TotalUntaxedReserves)
		assertYCEqual(t, "totalProvisions", oel.Provisions.TotalProvisions, pel.Provisions.TotalProvisions)
		assertYCEqual(t, "bankLoans", oel.LongTermLiabilities.BankLoans, pel.LongTermLiabilities.BankLoans)
		assertYCEqual(t, "totalLongTermLiabilities", oel.LongTermLiabilities.TotalLongTermLiabilities, pel.LongTermLiabilities.TotalLongTermLiabilities)
		assertYCEqual(t, "tradePayables", oel.ShortTermLiabilities.TradePayables, pel.ShortTermLiabilities.TradePayables)
		assertYCEqual(t, "totalShortTermLiabilities", oel.ShortTermLiabilities.TotalShortTermLiabilities, pel.ShortTermLiabilities.TotalShortTermLiabilities)
		assertYCEqual(t, "totalEquityAndLiabilities", oel.TotalEquityAndLiabilities, pel.TotalEquityAndLiabilities)
	})

	// --- Notes ---
	t.Run("accounting policies", func(t *testing.T) {
		assertNonEmpty(t, "description", parsed.Notes.AccountingPolicies.Description)
		if len(parsed.Notes.AccountingPolicies.Depreciations) != len(original.Notes.AccountingPolicies.Depreciations) {
			t.Errorf("depreciation count: got %d, want %d",
				len(parsed.Notes.AccountingPolicies.Depreciations),
				len(original.Notes.AccountingPolicies.Depreciations))
		}
	})

	t.Run("employees note", func(t *testing.T) {
		if parsed.Notes.Employees == nil {
			t.Fatal("employees note not parsed")
		}
		assertYCEqual(t, "averageEmployees",
			original.Notes.Employees.AverageEmployees,
			parsed.Notes.Employees.AverageEmployees)
	})

	t.Run("fixed asset notes", func(t *testing.T) {
		if len(parsed.Notes.FixedAssetNotes) != len(original.Notes.FixedAssetNotes) {
			t.Fatalf("fixed asset note count: got %d, want %d",
				len(parsed.Notes.FixedAssetNotes),
				len(original.Notes.FixedAssetNotes))
		}

		for i, ofan := range original.Notes.FixedAssetNotes {
			pfan := parsed.Notes.FixedAssetNotes[i]
			t.Run(ofan.ConceptPrefix, func(t *testing.T) {
				assertEqual(t, "conceptPrefix", ofan.ConceptPrefix, pfan.ConceptPrefix)
				assertYCEqual(t, "openingAcqValues", ofan.OpeningAcquisitionValues, pfan.OpeningAcquisitionValues)
				assertYCEqual(t, "closingAcqValues", ofan.ClosingAcquisitionValues, pfan.ClosingAcquisitionValues)
				assertYCEqual(t, "purchases", ofan.Purchases, pfan.Purchases)
				assertYCEqual(t, "carryingValue", ofan.CarryingValue, pfan.CarryingValue)

				if ofan.OpeningDepreciation.Current != nil {
					assertYCEqual(t, "openingDepreciation", ofan.OpeningDepreciation, pfan.OpeningDepreciation)
					assertYCEqual(t, "yearDepreciation", ofan.YearDepreciation, pfan.YearDepreciation)
					assertYCEqual(t, "closingDepreciation", ofan.ClosingDepreciation, pfan.ClosingDepreciation)
				}
			})
		}
	})

	t.Run("long term liabilities note", func(t *testing.T) {
		if parsed.Notes.LongTermLiabilitiesNote == nil {
			t.Fatal("long-term liabilities note not parsed")
		}
		assertYCEqual(t, "dueAfterFiveYears",
			original.Notes.LongTermLiabilitiesNote.DueAfterFiveYears,
			parsed.Notes.LongTermLiabilitiesNote.DueAfterFiveYears)
	})

	t.Run("pledges note", func(t *testing.T) {
		if parsed.Notes.Pledges == nil {
			t.Fatal("pledges note not parsed")
		}
		assertYCEqual(t, "corporateMortgages",
			original.Notes.Pledges.CorporateMortgages,
			parsed.Notes.Pledges.CorporateMortgages)
		assertYCEqual(t, "totalPledges",
			original.Notes.Pledges.TotalPledges,
			parsed.Notes.Pledges.TotalPledges)
	})

	t.Run("contingent liabilities note", func(t *testing.T) {
		if parsed.Notes.ContingentLiabilities == nil {
			t.Fatal("contingent liabilities note not parsed")
		}
		assertYCEqual(t, "totalContingent",
			original.Notes.ContingentLiabilities.TotalContingent,
			parsed.Notes.ContingentLiabilities.TotalContingent)
	})

	t.Run("multi-post note", func(t *testing.T) {
		if parsed.Notes.MultiPostNote == nil {
			t.Fatal("multi-post note not parsed")
		}
		assertNonEmpty(t, "description", parsed.Notes.MultiPostNote.Description)
		if len(parsed.Notes.MultiPostNote.Entries) != len(original.Notes.MultiPostNote.Entries) {
			t.Errorf("entry count: got %d, want %d",
				len(parsed.Notes.MultiPostNote.Entries),
				len(original.Notes.MultiPostNote.Entries))
		}
		for i, oe := range original.Notes.MultiPostNote.Entries {
			if i >= len(parsed.Notes.MultiPostNote.Entries) {
				break
			}
			pe := parsed.Notes.MultiPostNote.Entries[i]
			assertEqual(t, "postName", oe.PostName, pe.PostName)
			assertInt64PtrEqual(t, "amount", oe.Amount, pe.Amount)
		}
	})

	// --- Signatures ---
	t.Run("signatures", func(t *testing.T) {
		assertEqual(t, "city", original.Signatures.City, parsed.Signatures.City)
		assertEqual(t, "date", original.Signatures.Date, parsed.Signatures.Date)

		if len(parsed.Signatures.Signatories) != len(original.Signatures.Signatories) {
			t.Fatalf("signatory count: got %d, want %d",
				len(parsed.Signatures.Signatories),
				len(original.Signatures.Signatories))
		}
		for i, os := range original.Signatures.Signatories {
			ps := parsed.Signatures.Signatories[i]
			assertEqual(t, "firstName", os.FirstName, ps.FirstName)
			assertEqual(t, "lastName", os.LastName, ps.LastName)
			assertEqual(t, "role", os.Role, ps.Role)
		}
	})
}

// TestParseReferenceExample tests parsing the actual reference example file.
func TestParseReferenceExample(t *testing.T) {
	f, err := os.Open("../../ref/exempel/faststalld-arsredovisning-exempel-1.xhtml")
	if err != nil {
		t.Skipf("reference example not available: %v", err)
	}
	defer f.Close()

	report, err := Parse(f)
	if err != nil {
		t.Fatalf("parsing reference example: %v", err)
	}

	// Verify key values from the reference example.
	t.Run("company", func(t *testing.T) {
		assertEqual(t, "name", "Exempel 1 AB", report.Company.Name)
		assertEqual(t, "orgNr", "556999-9999", report.Company.OrgNr)
	})

	t.Run("fiscal year", func(t *testing.T) {
		assertEqual(t, "startDate", "2016-01-01", report.FiscalYear.StartDate)
		assertEqual(t, "endDate", "2016-12-31", report.FiscalYear.EndDate)
	})

	t.Run("income statement key values", func(t *testing.T) {
		assertInt64PtrValue(t, "netSales.current", 2650000, report.IncomeStatement.Revenue.NetSales.Current)
		assertInt64PtrValue(t, "netResult.current", 1274000, report.IncomeStatement.NetResult.Current)
	})

	t.Run("balance sheet key values", func(t *testing.T) {
		assertInt64PtrValue(t, "totalAssets.current", 7773000, report.BalanceSheet.Assets.TotalAssets.Current)
		assertInt64PtrValue(t, "totalEquityAndLiabilities.current", 7773000,
			report.BalanceSheet.EquityAndLiabilities.TotalEquityAndLiabilities.Current)
	})

	t.Run("signatures", func(t *testing.T) {
		if len(report.Signatures.Signatories) < 2 {
			t.Errorf("expected at least 2 signatories, got %d", len(report.Signatures.Signatories))
		}
	})
}

// TestParseNumberFormatting tests number parsing with various formats.
func TestParseNumberFormatting(t *testing.T) {
	tests := []struct {
		name string
		f    fact
		want int64
	}{
		{
			name: "plain number",
			f:    fact{Value: "2650000", Scale: 0},
			want: 2650000,
		},
		{
			name: "space-comma format",
			f:    fact{Value: "2 650 000", Format: "ixt:numspacecomma", Scale: 0},
			want: 2650000,
		},
		{
			name: "tkr with scale 3",
			f:    fact{Value: "2 650", Format: "ixt:numspacecomma", Scale: 3},
			want: 2650000,
		},
		{
			name: "negative sign",
			f:    fact{Value: "70 000", Format: "ixt:numspacecomma", Scale: 0, Sign: "-"},
			want: -70000,
		},
		{
			name: "percentage with comma decimal",
			f:    fact{Value: "33,7", Format: "ixt:numcomma", Scale: -2},
			want: 0, // 33.7 * 0.01 = 0.337, rounds to 0 as int64
		},
		{
			name: "zero",
			f:    fact{Value: "0"},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseNumber(tt.f)
			if err != nil {
				t.Fatalf("parseNumber: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

// TestParseExtractFacts verifies basic fact extraction from minimal iXBRL.
func TestParseExtractFacts(t *testing.T) {
	input := `<?xml version="1.0" encoding="UTF-8"?>
<html xmlns:ix="http://www.xbrl.org/2013/inlineXBRL"
      xmlns:se-gen-base="http://www.taxonomier.se/se/fr/gen-base/2021-10-31">
<body>
<ix:nonNumeric name="se-gen-base:ForetagetsNamn" contextRef="period0">Test AB</ix:nonNumeric>
<ix:nonFraction name="se-gen-base:Tillgangar" contextRef="balans0" unitRef="SEK" 
                decimals="INF" scale="0" format="ixt:numspacecomma">1 234 567</ix:nonFraction>
<ix:tuple name="se-gen-base:SomeTuple" tupleID="tuple1"/>
<ix:nonNumeric name="se-gen-base:SomeFact" contextRef="period0" tupleRef="tuple1" order="1.0">value</ix:nonNumeric>
</body>
</html>`

	facts, err := extractFacts(strings.NewReader(input))
	if err != nil {
		t.Fatalf("extractFacts: %v", err)
	}

	if len(facts) != 4 {
		t.Fatalf("expected 4 facts, got %d", len(facts))
	}

	// Check nonNumeric.
	if facts[0].Kind != "nonNumeric" || facts[0].Value != "Test AB" {
		t.Errorf("fact 0: got %+v", facts[0])
	}

	// Check nonFraction.
	if facts[1].Kind != "nonFraction" || facts[1].Value != "1 234 567" {
		t.Errorf("fact 1: got %+v", facts[1])
	}

	// Check tuple.
	if facts[2].Kind != "tuple" || facts[2].TupleID != "tuple1" {
		t.Errorf("fact 2: got %+v", facts[2])
	}

	// Check tuple member.
	if facts[3].TupleRef != "tuple1" || facts[3].Order != "1.0" {
		t.Errorf("fact 3: got %+v", facts[3])
	}
}

// ---------- test helpers ----------

func assertEqual(t *testing.T, field, want, got string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %q, want %q", field, got, want)
	}
}

func assertNonEmpty(t *testing.T, field, got string) {
	t.Helper()
	if got == "" {
		t.Errorf("%s: expected non-empty string", field)
	}
}

func assertInt64PtrEqual(t *testing.T, field string, want, got *int64) {
	t.Helper()
	if want == nil && got == nil {
		return
	}
	if want == nil {
		t.Errorf("%s: got %d, want nil", field, *got)
		return
	}
	if got == nil {
		t.Errorf("%s: got nil, want %d", field, *want)
		return
	}
	if *got != *want {
		t.Errorf("%s: got %d, want %d", field, *got, *want)
	}
}

func assertInt64PtrValue(t *testing.T, field string, want int64, got *int64) {
	t.Helper()
	if got == nil {
		t.Errorf("%s: got nil, want %d", field, want)
		return
	}
	if *got != want {
		t.Errorf("%s: got %d, want %d", field, *got, want)
	}
}

func assertYCEqual(t *testing.T, field string, want, got model.YearComparison) {
	t.Helper()
	assertInt64PtrEqual(t, field+".current", want.Current, got.Current)
	assertInt64PtrEqual(t, field+".previous", want.Previous, got.Previous)
}
