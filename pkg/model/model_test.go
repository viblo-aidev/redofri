package model

import (
	"encoding/json"
	"os"
	"testing"
)

func loadExempel1(t *testing.T) AnnualReport {
	t.Helper()
	data, err := os.ReadFile("../../testdata/exempel1.json")
	if err != nil {
		t.Fatalf("Failed to read testdata/exempel1.json: %v", err)
	}
	var report AnnualReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	return report
}

func TestLoadExempel1_Company(t *testing.T) {
	r := loadExempel1(t)
	if r.Company.Name != "Exempel 1 AB" {
		t.Errorf("Company.Name = %q, want %q", r.Company.Name, "Exempel 1 AB")
	}
	if r.Company.OrgNr != "556999-9999" {
		t.Errorf("Company.OrgNr = %q, want %q", r.Company.OrgNr, "556999-9999")
	}
}

func TestLoadExempel1_FiscalYear(t *testing.T) {
	r := loadExempel1(t)
	if r.FiscalYear.StartDate != "2016-01-01" {
		t.Errorf("FiscalYear.StartDate = %q, want %q", r.FiscalYear.StartDate, "2016-01-01")
	}
	if r.FiscalYear.EndDate != "2016-12-31" {
		t.Errorf("FiscalYear.EndDate = %q, want %q", r.FiscalYear.EndDate, "2016-12-31")
	}
}

func TestLoadExempel1_Meta(t *testing.T) {
	r := loadExempel1(t)
	if r.Meta.EntryPoint != "risbs" {
		t.Errorf("Meta.EntryPoint = %q, want %q", r.Meta.EntryPoint, "risbs")
	}
	if r.Meta.Currency != "SEK" {
		t.Errorf("Meta.Currency = %q, want %q", r.Meta.Currency, "SEK")
	}
}

func TestLoadExempel1_Certification(t *testing.T) {
	r := loadExempel1(t)
	if r.Certification.MeetingDate != "2017-03-21" {
		t.Errorf("Certification.MeetingDate = %q, want %q", r.Certification.MeetingDate, "2017-03-21")
	}
	if r.Certification.SigningDate != "2017-03-21" {
		t.Errorf("Certification.SigningDate = %q, want %q", r.Certification.SigningDate, "2017-03-21")
	}
	if r.Certification.Signatory.FirstName != "Karl" {
		t.Errorf("Certification.Signatory.FirstName = %q, want %q", r.Certification.Signatory.FirstName, "Karl")
	}
}

func TestLoadExempel1_IncomeStatement(t *testing.T) {
	r := loadExempel1(t)
	is := r.IncomeStatement

	assertYC(t, "NetSales", is.Revenue.NetSales, 2650000, 2250000)
	assertYC(t, "TotalRevenue", is.Revenue.TotalRevenue, 3727000, 4375000)
	assertYC(t, "TotalExpenses", is.Expenses.TotalExpenses, 3522000, 4111000)
	assertYC(t, "OperatingResult", is.OperatingResult, 205000, 264000)
	assertYC(t, "ResultAfterFinancialItems", is.ResultAfterFinancialItems, 1485000, 1184000)
	assertYC(t, "NetResult", is.NetResult, 1274000, 1099000)
}

func TestLoadExempel1_BalanceSheet(t *testing.T) {
	r := loadExempel1(t)
	bs := r.BalanceSheet

	assertYC(t, "TotalAssets", bs.Assets.TotalAssets, 7773000, 6007000)
	assertYC(t, "TotalEquityAndLiabilities", bs.EquityAndLiabilities.TotalEquityAndLiabilities, 7773000, 6007000)
	assertYC(t, "TotalEquity", bs.EquityAndLiabilities.Equity.TotalEquity, 2390000, 2215000)
	assertYC(t, "TotalFixedAssets", bs.Assets.FixedAssets.TotalFixedAssets, 4720000, 4060000)
	assertYC(t, "TotalCurrentAssets", bs.Assets.CurrentAssets.TotalCurrentAssets, 3053000, 1947000)
}

func TestLoadExempel1_BalanceSheetBalances(t *testing.T) {
	r := loadExempel1(t)
	bs := r.BalanceSheet

	// Assets = Equity + Liabilities
	if *bs.Assets.TotalAssets.Current != *bs.EquityAndLiabilities.TotalEquityAndLiabilities.Current {
		t.Errorf("Balance sheet doesn't balance (current): assets=%d != eq+liab=%d",
			*bs.Assets.TotalAssets.Current, *bs.EquityAndLiabilities.TotalEquityAndLiabilities.Current)
	}
	if *bs.Assets.TotalAssets.Previous != *bs.EquityAndLiabilities.TotalEquityAndLiabilities.Previous {
		t.Errorf("Balance sheet doesn't balance (previous): assets=%d != eq+liab=%d",
			*bs.Assets.TotalAssets.Previous, *bs.EquityAndLiabilities.TotalEquityAndLiabilities.Previous)
	}
}

func TestLoadExempel1_Notes(t *testing.T) {
	r := loadExempel1(t)

	if r.Notes.AccountingPolicies.NoteNumber != 1 {
		t.Errorf("AccountingPolicies.NoteNumber = %d, want 1", r.Notes.AccountingPolicies.NoteNumber)
	}
	if len(r.Notes.AccountingPolicies.Depreciations) != 3 {
		t.Errorf("len(Depreciations) = %d, want 3", len(r.Notes.AccountingPolicies.Depreciations))
	}

	if r.Notes.Employees == nil {
		t.Fatal("Employees note is nil")
	}
	assertYC(t, "AverageEmployees", r.Notes.Employees.AverageEmployees, 2, 2)

	if len(r.Notes.FixedAssetNotes) != 4 {
		t.Fatalf("len(FixedAssetNotes) = %d, want 4", len(r.Notes.FixedAssetNotes))
	}

	// Verify note 3 (Byggnader och mark) roll-forward
	n3 := r.Notes.FixedAssetNotes[0]
	if n3.NoteNumber != 3 {
		t.Errorf("FixedAssetNotes[0].NoteNumber = %d, want 3", n3.NoteNumber)
	}
	assertYC(t, "N3 CarryingValue", n3.CarryingValue, 1620000, 1450000)

	if r.Notes.Pledges == nil {
		t.Fatal("Pledges note is nil")
	}
	assertYC(t, "TotalPledges", r.Notes.Pledges.TotalPledges, 2800000, 2600000)

	if r.Notes.MultiPostNote == nil {
		t.Fatal("MultiPostNote is nil")
	}
	if len(r.Notes.MultiPostNote.Entries) != 2 {
		t.Errorf("len(MultiPostNote.Entries) = %d, want 2", len(r.Notes.MultiPostNote.Entries))
	}
}

func TestLoadExempel1_Signatures(t *testing.T) {
	r := loadExempel1(t)
	if r.Signatures.City != "Sundsvall" {
		t.Errorf("Signatures.City = %q, want %q", r.Signatures.City, "Sundsvall")
	}
	if len(r.Signatures.Signatories) != 2 {
		t.Fatalf("len(Signatories) = %d, want 2", len(r.Signatures.Signatories))
	}
	if r.Signatures.Signatories[1].Role != "Verkställande direktör" {
		t.Errorf("Signatories[1].Role = %q, want %q", r.Signatures.Signatories[1].Role, "Verkställande direktör")
	}
}

func TestJSONRoundtrip(t *testing.T) {
	r := loadExempel1(t)

	// Marshal back to JSON
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal again
	var r2 AnnualReport
	if err := json.Unmarshal(data, &r2); err != nil {
		t.Fatalf("Failed to unmarshal roundtrip: %v", err)
	}

	// Verify key fields survived the roundtrip
	if r2.Company.Name != r.Company.Name {
		t.Errorf("Roundtrip: Company.Name = %q, want %q", r2.Company.Name, r.Company.Name)
	}
	if *r2.IncomeStatement.NetResult.Current != *r.IncomeStatement.NetResult.Current {
		t.Errorf("Roundtrip: NetResult.Current = %d, want %d",
			*r2.IncomeStatement.NetResult.Current, *r.IncomeStatement.NetResult.Current)
	}
	if *r2.BalanceSheet.Assets.TotalAssets.Current != *r.BalanceSheet.Assets.TotalAssets.Current {
		t.Errorf("Roundtrip: TotalAssets.Current = %d, want %d",
			*r2.BalanceSheet.Assets.TotalAssets.Current, *r.BalanceSheet.Assets.TotalAssets.Current)
	}
}

func TestHelperFunctions(t *testing.T) {
	p := Int64(42)
	if *p != 42 {
		t.Errorf("Int64(42) = %d, want 42", *p)
	}

	s := String("hello")
	if *s != "hello" {
		t.Errorf("String(\"hello\") = %q, want %q", *s, "hello")
	}
}

// assertYC asserts a YearComparison has the expected current and previous values.
func assertYC(t *testing.T, name string, yc YearComparison, wantCurrent, wantPrevious int64) {
	t.Helper()
	if yc.Current == nil {
		t.Errorf("%s.Current is nil, want %d", name, wantCurrent)
	} else if *yc.Current != wantCurrent {
		t.Errorf("%s.Current = %d, want %d", name, *yc.Current, wantCurrent)
	}
	if yc.Previous == nil {
		t.Errorf("%s.Previous is nil, want %d", name, wantPrevious)
	} else if *yc.Previous != wantPrevious {
		t.Errorf("%s.Previous = %d, want %d", name, *yc.Previous, wantPrevious)
	}
}
