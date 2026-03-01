package validate

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/redofri/redofri/pkg/model"
)

// loadTestReport loads testdata/exempel1.json and returns a valid AnnualReport.
func loadTestReport(t *testing.T) *model.AnnualReport {
	t.Helper()
	data, err := os.ReadFile("../../testdata/exempel1.json")
	if err != nil {
		t.Fatalf("reading test data: %v", err)
	}
	var r model.AnnualReport
	if err := json.Unmarshal(data, &r); err != nil {
		t.Fatalf("parsing test data: %v", err)
	}
	return &r
}

// TestValidReport verifies that the example report passes validation with no errors.
func TestValidReport(t *testing.T) {
	r := loadTestReport(t)
	results := Validate(r)

	var errors []Result
	for _, res := range results {
		if res.Severity == Error {
			errors = append(errors, res)
		}
	}
	if len(errors) > 0 {
		for _, e := range errors {
			t.Errorf("unexpected error: %s", e)
		}
	}
}

// TestValidReportNoWarnings verifies that the example report produces no warnings
// (it has comparative figures and all fields populated).
func TestValidReportNoWarnings(t *testing.T) {
	r := loadTestReport(t)
	results := Validate(r)
	if len(results) > 0 {
		for _, res := range results {
			t.Errorf("unexpected finding: %s", res)
		}
	}
}

// TestEmptyReport validates that an empty report triggers all required-field errors.
func TestEmptyReport(t *testing.T) {
	r := &model.AnnualReport{}
	results := Validate(r)

	if !HasErrors(results) {
		t.Fatal("expected errors for empty report, got none")
	}

	// Check that specific BV codes are present.
	wantCodes := []int{
		1020, // company name
		1037, // currency
		1050, // country
		1173, // language
		1174, // amount format
		1019, // fastställelseintyg
		1103, // AGM date
		1164, // signing date
		1169, // signatory name
		1051, // directors' report
		1060, // income statement
		3001, // total assets
		3002, // total equity and liabilities
		1107, // signing date in AR
		1201, // signatories
	}

	codeMap := make(map[int]bool)
	for _, res := range results {
		if res.Code > 0 {
			codeMap[res.Code] = true
		}
	}

	for _, code := range wantCodes {
		if !codeMap[code] {
			t.Errorf("missing expected BV code %d in validation results", code)
		}
	}
}

// TestMissingCompanyName checks BV code 1020.
func TestMissingCompanyName(t *testing.T) {
	r := loadTestReport(t)
	r.Company.Name = ""
	results := Validate(r)
	assertHasCode(t, results, 1020)
}

// TestMissingCurrency checks BV code 1037.
func TestMissingCurrency(t *testing.T) {
	r := loadTestReport(t)
	r.Meta.Currency = ""
	results := Validate(r)
	assertHasCode(t, results, 1037)
}

// TestInvalidCurrency checks BV code 1038.
func TestInvalidCurrency(t *testing.T) {
	r := loadTestReport(t)
	r.Meta.Currency = "USD"
	results := Validate(r)
	assertHasCode(t, results, 1038)
}

// TestValidCurrencyEUR checks that EUR is accepted.
func TestValidCurrencyEUR(t *testing.T) {
	r := loadTestReport(t)
	r.Meta.Currency = "EUR"
	results := Validate(r)
	assertNoCode(t, results, 1037)
	assertNoCode(t, results, 1038)
}

// TestMissingCountry checks BV code 1050.
func TestMissingCountry(t *testing.T) {
	r := loadTestReport(t)
	r.Meta.Country = ""
	results := Validate(r)
	assertHasCode(t, results, 1050)
}

// TestMissingFastställelseintyg checks BV code 1019.
func TestMissingFastställelseintyg(t *testing.T) {
	r := loadTestReport(t)
	r.Certification.ConfirmationText = ""
	results := Validate(r)
	assertHasCode(t, results, 1019)
}

// TestMissingAGMDate checks BV code 1103.
func TestMissingAGMDate(t *testing.T) {
	r := loadTestReport(t)
	r.Certification.MeetingDate = ""
	results := Validate(r)
	assertHasCode(t, results, 1103)
}

// TestMissingCertificationSigningDate checks BV code 1164.
func TestMissingCertificationSigningDate(t *testing.T) {
	r := loadTestReport(t)
	r.Certification.SigningDate = ""
	results := Validate(r)
	assertHasCode(t, results, 1164)
}

// TestMissingCertificationSignatory checks BV code 1169.
func TestMissingCertificationSignatory(t *testing.T) {
	r := loadTestReport(t)
	r.Certification.Signatory.FirstName = ""
	r.Certification.Signatory.LastName = ""
	results := Validate(r)
	assertHasCode(t, results, 1169)
}

// TestMissingDirectorsReport checks BV code 1051.
func TestMissingDirectorsReport(t *testing.T) {
	r := loadTestReport(t)
	r.ManagementReport.BusinessDescription = ""
	r.ManagementReport.IntroText = ""
	results := Validate(r)
	assertHasCode(t, results, 1051)
}

// TestMissingSigningDate checks BV code 1107.
func TestMissingSigningDate(t *testing.T) {
	r := loadTestReport(t)
	r.Signatures.Date = ""
	results := Validate(r)
	assertHasCode(t, results, 1107)
}

// TestNoSignatories checks BV code 1201.
func TestNoSignatories(t *testing.T) {
	r := loadTestReport(t)
	r.Signatures.Signatories = nil
	results := Validate(r)
	assertHasCode(t, results, 1201)
}

// TestSignatoryMissingName checks BV code 1201 for incomplete signatory.
func TestSignatoryMissingName(t *testing.T) {
	r := loadTestReport(t)
	r.Signatures.Signatories[0].FirstName = ""
	results := Validate(r)
	assertHasCode(t, results, 1201)
}

// TestMissingLanguage checks BV code 1173.
func TestMissingLanguage(t *testing.T) {
	r := loadTestReport(t)
	r.Meta.Language = ""
	results := Validate(r)
	assertHasCode(t, results, 1173)
}

// TestNonSwedishLanguage checks BV code 1116 (warning).
func TestNonSwedishLanguage(t *testing.T) {
	r := loadTestReport(t)
	r.Meta.Language = "en"
	results := Validate(r)
	assertHasCode(t, results, 1116)
	// Should be a warning, not an error.
	for _, res := range results {
		if res.Code == 1116 && res.Severity != Warning {
			t.Errorf("code 1116 should be a Warning, got %s", res.Severity)
		}
	}
}

// TestMissingAmountFormat checks BV code 1174.
func TestMissingAmountFormat(t *testing.T) {
	r := loadTestReport(t)
	r.Meta.AmountFormat = ""
	results := Validate(r)
	assertHasCode(t, results, 1174)
}

// TestInvalidOrgNrFormat checks org number format validation.
func TestInvalidOrgNrFormat(t *testing.T) {
	r := loadTestReport(t)
	r.Company.OrgNr = "5569999999" // Missing hyphen
	results := Validate(r)
	assertHasFieldError(t, results, "company.orgNr")
}

// TestInvalidEntryPoint checks that invalid entry point is rejected.
func TestInvalidEntryPoint(t *testing.T) {
	r := loadTestReport(t)
	r.Meta.EntryPoint = "invalid"
	results := Validate(r)
	assertHasFieldError(t, results, "meta.entryPoint")
}

// TestValidEntryPoints checks that all four valid entry points are accepted.
func TestValidEntryPoints(t *testing.T) {
	for _, ep := range []string{"risbs", "risab", "raibs", "raiab"} {
		r := loadTestReport(t)
		r.Meta.EntryPoint = ep
		results := Validate(r)
		assertNoFieldError(t, results, "meta.entryPoint")
	}
}

// --- Date ordering tests ---

// TestFiscalYearExceeds18Months checks BV code 1046.
func TestFiscalYearExceeds18Months(t *testing.T) {
	r := loadTestReport(t)
	r.FiscalYear.StartDate = "2015-01-01"
	r.FiscalYear.EndDate = "2016-12-31"
	results := Validate(r)
	assertHasCode(t, results, 1046)
}

// TestFiscalYear18MonthsOK checks that an 18-month fiscal year is accepted.
func TestFiscalYear18MonthsOK(t *testing.T) {
	r := loadTestReport(t)
	r.FiscalYear.StartDate = "2015-07-01"
	r.FiscalYear.EndDate = "2016-12-31"
	results := Validate(r)
	assertNoCode(t, results, 1046)
}

// TestSigningDateBeforeFYEnd checks BV code 1114.
func TestSigningDateBeforeFYEnd(t *testing.T) {
	r := loadTestReport(t)
	r.Signatures.Date = "2016-12-31" // Same as FY end
	results := Validate(r)
	assertHasCode(t, results, 1114)
}

// TestSigningDateEqualsFYEnd checks BV code 1114 (must be after, not equal).
func TestSigningDateEqualsFYEnd(t *testing.T) {
	r := loadTestReport(t)
	r.Signatures.Date = "2016-12-31"
	results := Validate(r)
	assertHasCode(t, results, 1114)
}

// TestAGMBeforeFYEnd checks BV code 1101.
func TestAGMBeforeFYEnd(t *testing.T) {
	r := loadTestReport(t)
	r.Certification.MeetingDate = "2016-11-15" // Before FY end
	results := Validate(r)
	assertHasCode(t, results, 1101)
}

// TestCertSigningBeforeAGM checks BV code 1165.
func TestCertSigningBeforeAGM(t *testing.T) {
	r := loadTestReport(t)
	r.Certification.MeetingDate = "2017-03-21"
	r.Certification.SigningDate = "2017-03-20" // One day before AGM
	results := Validate(r)
	assertHasCode(t, results, 1165)
}

// TestAGMBeforeSigningDate checks BV code 1183 (warning).
func TestAGMBeforeSigningDate(t *testing.T) {
	r := loadTestReport(t)
	r.Signatures.Date = "2017-03-25"
	r.Certification.MeetingDate = "2017-03-21" // Before signing
	results := Validate(r)
	assertHasCode(t, results, 1183)
}

// --- Calculation check tests ---

// TestIncomeStatementRevenueCalcError triggers a revenue sum error.
func TestIncomeStatementRevenueCalcError(t *testing.T) {
	r := loadTestReport(t)
	// Break the total revenue: should be 2650000+700000+377000 = 3727000
	wrong := int64(9999999)
	r.IncomeStatement.Revenue.TotalRevenue.Current = &wrong
	results := Validate(r)
	assertHasFieldError(t, results, "incomeStatement.revenue.totalRevenue.current")
}

// TestIncomeStatementExpenseCalcError triggers an expense sum error.
func TestIncomeStatementExpenseCalcError(t *testing.T) {
	r := loadTestReport(t)
	wrong := int64(0)
	r.IncomeStatement.Expenses.TotalExpenses.Current = &wrong
	results := Validate(r)
	assertHasFieldError(t, results, "incomeStatement.expenses.totalExpenses.current")
}

// TestOperatingResultCalcError triggers an operating result error.
func TestOperatingResultCalcError(t *testing.T) {
	r := loadTestReport(t)
	wrong := int64(0)
	r.IncomeStatement.OperatingResult.Current = &wrong
	results := Validate(r)
	assertHasFieldError(t, results, "incomeStatement.operatingResult.current")
}

// TestFinancialItemsCalcError triggers a financial items sum error.
func TestFinancialItemsCalcError(t *testing.T) {
	r := loadTestReport(t)
	wrong := int64(0)
	r.IncomeStatement.FinancialItems.TotalFinancialItems.Current = &wrong
	results := Validate(r)
	assertHasFieldError(t, results, "incomeStatement.financialItems.totalFinancialItems.current")
}

// TestNetResultCalcError triggers a net result error.
func TestNetResultCalcError(t *testing.T) {
	r := loadTestReport(t)
	wrong := int64(0)
	r.IncomeStatement.NetResult.Current = &wrong
	results := Validate(r)
	assertHasFieldError(t, results, "incomeStatement.netResult.current")
}

// TestBalanceSheetBalanceError triggers BV 3005 (assets ≠ equity+liabilities).
func TestBalanceSheetBalanceError(t *testing.T) {
	r := loadTestReport(t)
	wrong := int64(1)
	r.BalanceSheet.Assets.TotalAssets.Current = &wrong
	results := Validate(r)
	assertHasCode(t, results, 3005)
}

// TestBalanceSheetTangibleCalcError triggers a tangible fixed assets sum error.
func TestBalanceSheetTangibleCalcError(t *testing.T) {
	r := loadTestReport(t)
	wrong := int64(0)
	r.BalanceSheet.Assets.FixedAssets.Tangible.TotalTangible.Current = &wrong
	results := Validate(r)
	assertHasFieldError(t, results, "balanceSheet.assets.fixedAssets.tangible.totalTangible.current")
}

// TestBalanceSheetEquityCalcError triggers a total equity sum error.
func TestBalanceSheetEquityCalcError(t *testing.T) {
	r := loadTestReport(t)
	wrong := int64(0)
	r.BalanceSheet.EquityAndLiabilities.Equity.TotalEquity.Current = &wrong
	results := Validate(r)
	assertHasFieldError(t, results, "balanceSheet.equityAndLiabilities.equity.totalEquity.current")
}

// TestEquityChangesCalcError triggers an equity changes calculation error.
func TestEquityChangesCalcError(t *testing.T) {
	r := loadTestReport(t)
	wrong := int64(0)
	r.ManagementReport.EquityChanges.OpeningTotal = &wrong
	results := Validate(r)
	assertHasFieldError(t, results, "managementReport.equityChanges.openingTotal")
}

// TestProfitDispositionCalcError triggers a profit disposition error.
func TestProfitDispositionCalcError(t *testing.T) {
	r := loadTestReport(t)
	wrong := int64(0)
	r.ManagementReport.ProfitDisposition.TotalAvailable = &wrong
	results := Validate(r)
	assertHasFieldError(t, results, "managementReport.profitDisposition.totalAvailable")
}

// TestProfitDispositionMismatch triggers the available ≠ disposition error.
func TestProfitDispositionMismatch(t *testing.T) {
	r := loadTestReport(t)
	// Make total available correct but total disposition different
	wrong := int64(1000000)
	r.ManagementReport.ProfitDisposition.TotalDisposition = &wrong
	// Also fix the underlying values to avoid the composition check error
	div := int64(500000)
	cf := int64(500000)
	r.ManagementReport.ProfitDisposition.Dividend = &div
	r.ManagementReport.ProfitDisposition.CarriedForward = &cf
	results := Validate(r)
	assertHasFieldError(t, results, "managementReport.profitDisposition")
}

// --- Comparative figures tests ---

// TestMissingComparativeFiguresIS checks BV 3007 (warning).
func TestMissingComparativeFiguresIS(t *testing.T) {
	r := loadTestReport(t)
	r.IncomeStatement.NetResult.Previous = nil
	results := Validate(r)
	assertHasCode(t, results, 3007)
	// Should be warning, not error
	for _, res := range results {
		if res.Code == 3007 && res.Severity != Warning {
			t.Errorf("code 3007 should be a Warning, got %s", res.Severity)
		}
	}
}

// TestMissingComparativeFiguresBS checks BV 3006 (warning).
func TestMissingComparativeFiguresBS(t *testing.T) {
	r := loadTestReport(t)
	r.BalanceSheet.Assets.TotalAssets.Previous = nil
	results := Validate(r)
	assertHasCode(t, results, 3006)
}

// --- Fixed asset note tests ---

// TestFixedAssetNoteCarryingValueError checks that carrying value validation works.
func TestFixedAssetNoteCarryingValueError(t *testing.T) {
	r := loadTestReport(t)
	// Break the first note's carrying value
	wrong := int64(999)
	r.Notes.FixedAssetNotes[0].CarryingValue.Current = &wrong
	results := Validate(r)
	assertHasFieldError(t, results, "notes.fixedAssetNotes[0].carryingValue.current")
}

// --- Invalid date format tests ---

// TestInvalidDateFormat checks that malformed dates are caught.
func TestInvalidDateFormat(t *testing.T) {
	r := loadTestReport(t)
	r.FiscalYear.StartDate = "2016/01/01" // Wrong format
	results := Validate(r)
	assertHasFieldError(t, results, "fiscalYear.startDate")
}

// --- HasErrors helper test ---

func TestHasErrors(t *testing.T) {
	tests := []struct {
		name    string
		results []Result
		want    bool
	}{
		{"empty", nil, false},
		{"warnings only", []Result{{Severity: Warning}}, false},
		{"errors", []Result{{Severity: Error}}, true},
		{"mixed", []Result{{Severity: Warning}, {Severity: Error}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasErrors(tt.results); got != tt.want {
				t.Errorf("HasErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- Severity.String() test ---

func TestSeverityString(t *testing.T) {
	if Error.String() != "ERROR" {
		t.Errorf("Error.String() = %q, want %q", Error.String(), "ERROR")
	}
	if Warning.String() != "WARN" {
		t.Errorf("Warning.String() = %q, want %q", Warning.String(), "WARN")
	}
}

// --- Result.String() test ---

func TestResultString(t *testing.T) {
	r := Result{Error, 1020, "company.name", "company name is missing"}
	got := r.String()
	want := "ERROR [BV 1020]: company.name: company name is missing"
	if got != want {
		t.Errorf("Result.String() = %q, want %q", got, want)
	}

	r2 := Result{Warning, 0, "field", "msg"}
	got2 := r2.String()
	want2 := "WARN: field: msg"
	if got2 != want2 {
		t.Errorf("Result.String() = %q, want %q", got2, want2)
	}
}

// --- Previous-year calculation tests ---

// TestPreviousYearCalculations checks that previous year calculations are also validated.
func TestPreviousYearCalculations(t *testing.T) {
	r := loadTestReport(t)
	// Break the previous year total revenue
	wrong := int64(0)
	r.IncomeStatement.Revenue.TotalRevenue.Previous = &wrong
	results := Validate(r)
	assertHasFieldError(t, results, "incomeStatement.revenue.totalRevenue.previous")
}

// --- Helper functions ---

func assertHasCode(t *testing.T, results []Result, code int) {
	t.Helper()
	for _, res := range results {
		if res.Code == code {
			return
		}
	}
	t.Errorf("expected BV code %d in results, not found", code)
	for _, res := range results {
		t.Logf("  %s", res)
	}
}

func assertNoCode(t *testing.T, results []Result, code int) {
	t.Helper()
	for _, res := range results {
		if res.Code == code {
			t.Errorf("unexpected BV code %d in results: %s", code, res)
		}
	}
}

func assertHasFieldError(t *testing.T, results []Result, field string) {
	t.Helper()
	for _, res := range results {
		if res.Field == field && res.Severity == Error {
			return
		}
	}
	t.Errorf("expected error for field %q in results, not found", field)
	for _, res := range results {
		t.Logf("  %s", res)
	}
}

func assertNoFieldError(t *testing.T, results []Result, field string) {
	t.Helper()
	for _, res := range results {
		if res.Field == field && res.Severity == Error {
			t.Errorf("unexpected error for field %q: %s", field, res)
		}
	}
}
