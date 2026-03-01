// Package validate checks a model.AnnualReport for correctness before
// iXBRL generation or submission to Bolagsverket.
//
// Validation is split into three categories:
//
//  1. Required fields — mandatory data that must be present.
//  2. Calculation checks — sums that must be internally consistent
//     (mirrors the XBRL calculation linkbase relationships).
//  3. Business rules — date ordering, format, and semantic constraints
//     (mirrors Bolagsverket's "dokumentkontroller" codes 1019–3007).
//
// Each finding is returned as a Result with a severity (Error or Warning),
// a Bolagsverket code where applicable, and a human-readable message.
package validate

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/redofri/redofri/pkg/model"
)

// Severity indicates how critical a validation finding is.
type Severity int

const (
	Error   Severity = iota // Must fix before submission.
	Warning                 // Advisory; submission may still succeed.
)

func (s Severity) String() string {
	if s == Error {
		return "ERROR"
	}
	return "WARN"
}

// Result is a single validation finding.
type Result struct {
	Severity Severity
	Code     int    // Bolagsverket code (0 = internal rule, no BV code).
	Field    string // Dotted path to the field, e.g. "company.name".
	Message  string
}

func (r Result) String() string {
	code := ""
	if r.Code > 0 {
		code = fmt.Sprintf(" [BV %d]", r.Code)
	}
	return fmt.Sprintf("%s%s: %s: %s", r.Severity, code, r.Field, r.Message)
}

// Validate checks the given AnnualReport and returns all findings.
// An empty slice means the report passes all checks.
func Validate(r *model.AnnualReport) []Result {
	v := &validator{report: r}
	v.checkRequiredFields()
	v.checkCalculations()
	v.checkBusinessRules()
	return v.results
}

// HasErrors returns true if any result has Error severity.
func HasErrors(results []Result) bool {
	for _, r := range results {
		if r.Severity == Error {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Internal validator
// ---------------------------------------------------------------------------

type validator struct {
	report  *model.AnnualReport
	results []Result
}

func (v *validator) err(code int, field, msg string) {
	v.results = append(v.results, Result{Error, code, field, msg})
}

func (v *validator) warn(code int, field, msg string) {
	v.results = append(v.results, Result{Warning, code, field, msg})
}

// ---------------------------------------------------------------------------
// 1. Required fields
// ---------------------------------------------------------------------------

func (v *validator) checkRequiredFields() {
	r := v.report

	// Company (BV 1020)
	if r.Company.Name == "" {
		v.err(1020, "company.name", "company name is missing")
	}
	if r.Company.OrgNr == "" {
		v.err(0, "company.orgNr", "organisation number is missing")
	}

	// Fiscal year
	if r.FiscalYear.StartDate == "" {
		v.err(0, "fiscalYear.startDate", "fiscal year start date is missing")
	}
	if r.FiscalYear.EndDate == "" {
		v.err(0, "fiscalYear.endDate", "fiscal year end date is missing")
	}

	// Meta (BV 1037, 1050, 1173)
	if r.Meta.Currency == "" {
		v.err(1037, "meta.currency", "currency is missing")
	}
	if r.Meta.Country == "" {
		v.err(1050, "meta.country", "country code is missing")
	}
	if r.Meta.Language == "" {
		v.err(1173, "meta.language", "language indication is missing")
	}
	if r.Meta.AmountFormat == "" {
		v.err(1174, "meta.amountFormat", "amount format (unit of measurement) is missing")
	}
	if r.Meta.EntryPoint == "" {
		v.err(0, "meta.entryPoint", "entry point is missing")
	}

	// Certification (BV 1019, 1103, 1164, 1169)
	if r.Certification.ConfirmationText == "" {
		v.err(1019, "certification.confirmationText", "fastställelseintyg is missing")
	}
	if r.Certification.MeetingDate == "" {
		v.err(1103, "certification.meetingDate", "AGM date is missing in adoption certificate")
	}
	if r.Certification.SigningDate == "" {
		v.err(1164, "certification.signingDate", "signing date is missing in adoption certificate")
	}
	if r.Certification.Signatory.FirstName == "" && r.Certification.Signatory.LastName == "" {
		v.err(1169, "certification.signatory", "name is missing in adoption certificate")
	}

	// Management report (BV 1051)
	if r.ManagementReport.BusinessDescription == "" && r.ManagementReport.IntroText == "" {
		v.err(1051, "managementReport", "directors' report (förvaltningsberättelse) is missing")
	}

	// Income statement (BV 1060)
	if r.IncomeStatement.NetResult.Current == nil {
		v.err(1060, "incomeStatement", "income statement is missing (no net result)")
	}

	// Balance sheet (BV 1064, 3001, 3002)
	if r.BalanceSheet.Assets.TotalAssets.Current == nil {
		v.err(3001, "balanceSheet.assets.totalAssets", "total assets is missing")
	}
	if r.BalanceSheet.EquityAndLiabilities.TotalEquityAndLiabilities.Current == nil {
		v.err(3002, "balanceSheet.equityAndLiabilities.totalEquityAndLiabilities",
			"total equity and liabilities is missing")
	}

	// Signatures (BV 1107, 1201)
	if r.Signatures.Date == "" {
		v.err(1107, "signatures.date", "signing date is missing in annual report")
	}
	if len(r.Signatures.Signatories) == 0 {
		v.err(1201, "signatures.signatories", "no signatories in annual report")
	} else {
		for i, s := range r.Signatures.Signatories {
			if s.FirstName == "" || s.LastName == "" {
				v.err(1201, fmt.Sprintf("signatures.signatories[%d]", i),
					"first or last name is missing for signatory")
			}
		}
	}

	// Notes — accounting policies is always required
	if r.Notes.AccountingPolicies.Description == "" {
		v.err(0, "notes.accountingPolicies.description", "accounting policies description is missing")
	}
}

// ---------------------------------------------------------------------------
// 2. Calculation checks (XBRL calculation linkbase rules)
// ---------------------------------------------------------------------------

func (v *validator) checkCalculations() {
	v.checkIncomeStatementCalc()
	v.checkBalanceSheetCalc()
	v.checkEquityChangesCalc()
	v.checkProfitDispositionCalc()
}

// i64 safely dereferences a *int64, returning 0 if nil.
func i64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

func (v *validator) calcCheck(field string, got, want int64) {
	if got != want {
		v.err(0, field, fmt.Sprintf("calculation error: got %d, expected %d (diff %d)",
			got, want, got-want))
	}
}

func (v *validator) checkIncomeStatementCalc() {
	is := v.report.IncomeStatement

	// Total revenue = net sales + inventory change + other operating income
	for _, label := range []string{"current", "previous"} {
		var cur bool
		if label == "current" {
			cur = true
		}
		pick := func(yc model.YearComparison) int64 {
			if cur {
				return i64(yc.Current)
			}
			return i64(yc.Previous)
		}

		prefix := "incomeStatement.revenue.totalRevenue." + label
		totalRev := pick(is.Revenue.TotalRevenue)
		wantRev := pick(is.Revenue.NetSales) + pick(is.Revenue.InventoryChange) +
			pick(is.Revenue.OtherOperatingIncome)
		v.calcCheck(prefix, totalRev, wantRev)

		// Total expenses
		prefix = "incomeStatement.expenses.totalExpenses." + label
		totalExp := pick(is.Expenses.TotalExpenses)
		wantExp := pick(is.Expenses.RawMaterials) + pick(is.Expenses.TradingGoods) +
			pick(is.Expenses.OtherExternalExpenses) + pick(is.Expenses.PersonnelExpenses) +
			pick(is.Expenses.DepreciationAmortization) + pick(is.Expenses.OtherOperatingExpenses)
		v.calcCheck(prefix, totalExp, wantExp)

		// Operating result = total revenue - total expenses
		prefix = "incomeStatement.operatingResult." + label
		v.calcCheck(prefix, pick(is.OperatingResult), totalRev-totalExp)

		// Total financial items
		prefix = "incomeStatement.financialItems.totalFinancialItems." + label
		totalFin := pick(is.FinancialItems.TotalFinancialItems)
		wantFin := pick(is.FinancialItems.ResultOtherFinancialAssets) +
			pick(is.FinancialItems.OtherInterestIncome) -
			pick(is.FinancialItems.InterestExpenses)
		v.calcCheck(prefix, totalFin, wantFin)

		// Result after financial items = operating result + financial items
		prefix = "incomeStatement.resultAfterFinancialItems." + label
		v.calcCheck(prefix, pick(is.ResultAfterFinancialItems),
			pick(is.OperatingResult)+totalFin)

		// Total appropriations
		prefix = "incomeStatement.appropriations.totalAppropriations." + label
		totalAppr := pick(is.Appropriations.TotalAppropriations)
		wantAppr := pick(is.Appropriations.TaxAllocationReserve) +
			pick(is.Appropriations.ExcessDepreciation)
		v.calcCheck(prefix, totalAppr, wantAppr)

		// Result before tax = result after financial - appropriations
		prefix = "incomeStatement.resultBeforeTax." + label
		v.calcCheck(prefix, pick(is.ResultBeforeTax),
			pick(is.ResultAfterFinancialItems)-totalAppr)

		// Net result = result before tax - tax
		prefix = "incomeStatement.netResult." + label
		v.calcCheck(prefix, pick(is.NetResult),
			pick(is.ResultBeforeTax)-pick(is.Tax.IncomeTax))
	}
}

func (v *validator) checkBalanceSheetCalc() {
	bs := v.report.BalanceSheet

	for _, label := range []string{"current", "previous"} {
		cur := label == "current"
		pick := func(yc model.YearComparison) int64 {
			if cur {
				return i64(yc.Current)
			}
			return i64(yc.Previous)
		}
		pfx := func(s string) string { return s + "." + label }

		// Tangible fixed assets
		totalTang := pick(bs.Assets.FixedAssets.Tangible.TotalTangible)
		wantTang := pick(bs.Assets.FixedAssets.Tangible.BuildingsAndLand) +
			pick(bs.Assets.FixedAssets.Tangible.MachineryAndEquipment) +
			pick(bs.Assets.FixedAssets.Tangible.FixturesAndFittings)
		v.calcCheck(pfx("balanceSheet.assets.fixedAssets.tangible.totalTangible"), totalTang, wantTang)

		// Financial fixed assets
		totalFin := pick(bs.Assets.FixedAssets.Financial.TotalFinancial)
		wantFin := pick(bs.Assets.FixedAssets.Financial.OtherLongTermSecurities)
		v.calcCheck(pfx("balanceSheet.assets.fixedAssets.financial.totalFinancial"), totalFin, wantFin)

		// Total fixed assets
		totalFixed := pick(bs.Assets.FixedAssets.TotalFixedAssets)
		v.calcCheck(pfx("balanceSheet.assets.fixedAssets.totalFixedAssets"), totalFixed, totalTang+totalFin)

		// Inventory
		totalInv := pick(bs.Assets.CurrentAssets.Inventory.TotalInventory)
		wantInv := pick(bs.Assets.CurrentAssets.Inventory.RawMaterials) +
			pick(bs.Assets.CurrentAssets.Inventory.WorkInProgress) +
			pick(bs.Assets.CurrentAssets.Inventory.FinishedGoods)
		v.calcCheck(pfx("balanceSheet.assets.currentAssets.inventory.totalInventory"), totalInv, wantInv)

		// Short-term receivables
		totalSTR := pick(bs.Assets.CurrentAssets.ShortTermReceivables.TotalShortTermReceivables)
		wantSTR := pick(bs.Assets.CurrentAssets.ShortTermReceivables.TradeReceivables) +
			pick(bs.Assets.CurrentAssets.ShortTermReceivables.OtherReceivables) +
			pick(bs.Assets.CurrentAssets.ShortTermReceivables.PrepaidExpenses)
		v.calcCheck(pfx("balanceSheet.assets.currentAssets.shortTermReceivables.totalShortTermReceivables"), totalSTR, wantSTR)

		// Cash and bank
		totalCash := pick(bs.Assets.CurrentAssets.CashAndBank.TotalCashAndBank)
		wantCash := pick(bs.Assets.CurrentAssets.CashAndBank.CashAndBankExcl)
		v.calcCheck(pfx("balanceSheet.assets.currentAssets.cashAndBank.totalCashAndBank"), totalCash, wantCash)

		// Total current assets
		totalCurAss := pick(bs.Assets.CurrentAssets.TotalCurrentAssets)
		v.calcCheck(pfx("balanceSheet.assets.currentAssets.totalCurrentAssets"), totalCurAss, totalInv+totalSTR+totalCash)

		// Total assets
		totalAss := pick(bs.Assets.TotalAssets)
		v.calcCheck(pfx("balanceSheet.assets.totalAssets"), totalAss, totalFixed+totalCurAss)

		// --- Equity and liabilities ---

		// Restricted equity
		totalRestEq := pick(bs.EquityAndLiabilities.Equity.TotalRestrictedEquity)
		wantRestEq := pick(bs.EquityAndLiabilities.Equity.ShareCapital) +
			pick(bs.EquityAndLiabilities.Equity.ReserveFund)
		v.calcCheck(pfx("balanceSheet.equityAndLiabilities.equity.totalRestrictedEquity"), totalRestEq, wantRestEq)

		// Unrestricted equity
		totalUnrestEq := pick(bs.EquityAndLiabilities.Equity.TotalUnrestrictedEquity)
		wantUnrestEq := pick(bs.EquityAndLiabilities.Equity.RetainedEarnings) +
			pick(bs.EquityAndLiabilities.Equity.NetIncome)
		v.calcCheck(pfx("balanceSheet.equityAndLiabilities.equity.totalUnrestrictedEquity"), totalUnrestEq, wantUnrestEq)

		// Total equity
		totalEq := pick(bs.EquityAndLiabilities.Equity.TotalEquity)
		v.calcCheck(pfx("balanceSheet.equityAndLiabilities.equity.totalEquity"), totalEq, totalRestEq+totalUnrestEq)

		// Untaxed reserves
		totalUntax := pick(bs.EquityAndLiabilities.UntaxedReserves.TotalUntaxedReserves)
		wantUntax := pick(bs.EquityAndLiabilities.UntaxedReserves.TaxAllocationReserves) +
			pick(bs.EquityAndLiabilities.UntaxedReserves.AccumulatedExcessDepreciation)
		v.calcCheck(pfx("balanceSheet.equityAndLiabilities.untaxedReserves.totalUntaxedReserves"), totalUntax, wantUntax)

		// Provisions
		totalProv := pick(bs.EquityAndLiabilities.Provisions.TotalProvisions)
		wantProv := pick(bs.EquityAndLiabilities.Provisions.PensionProvisions) +
			pick(bs.EquityAndLiabilities.Provisions.OtherProvisions)
		v.calcCheck(pfx("balanceSheet.equityAndLiabilities.provisions.totalProvisions"), totalProv, wantProv)

		// Long-term liabilities
		totalLT := pick(bs.EquityAndLiabilities.LongTermLiabilities.TotalLongTermLiabilities)
		wantLT := pick(bs.EquityAndLiabilities.LongTermLiabilities.BankLoans) +
			pick(bs.EquityAndLiabilities.LongTermLiabilities.OtherLongTermLiabilities)
		v.calcCheck(pfx("balanceSheet.equityAndLiabilities.longTermLiabilities.totalLongTermLiabilities"), totalLT, wantLT)

		// Short-term liabilities
		totalST := pick(bs.EquityAndLiabilities.ShortTermLiabilities.TotalShortTermLiabilities)
		wantST := pick(bs.EquityAndLiabilities.ShortTermLiabilities.TradePayables) +
			pick(bs.EquityAndLiabilities.ShortTermLiabilities.TaxLiabilities) +
			pick(bs.EquityAndLiabilities.ShortTermLiabilities.OtherShortTermLiabilities) +
			pick(bs.EquityAndLiabilities.ShortTermLiabilities.AccruedExpenses)
		v.calcCheck(pfx("balanceSheet.equityAndLiabilities.shortTermLiabilities.totalShortTermLiabilities"), totalST, wantST)

		// Total equity and liabilities
		totalEL := pick(bs.EquityAndLiabilities.TotalEquityAndLiabilities)
		v.calcCheck(pfx("balanceSheet.equityAndLiabilities.totalEquityAndLiabilities"),
			totalEL, totalEq+totalUntax+totalProv+totalLT+totalST)

		// Balance sheet balance: total assets = total equity & liabilities (BV 3005)
		if totalAss != totalEL {
			v.err(3005, pfx("balanceSheet"),
				fmt.Sprintf("total assets (%d) ≠ total equity and liabilities (%d)", totalAss, totalEL))
		}
	}
}

func (v *validator) checkEquityChangesCalc() {
	ec := v.report.ManagementReport.EquityChanges

	// Opening total = opening share capital + reserve fund + retained earnings + net income
	wantOpenTotal := i64(ec.OpeningShareCapital) + i64(ec.OpeningReserveFund) +
		i64(ec.OpeningRetainedEarnings) + i64(ec.OpeningNetIncome)
	if ec.OpeningTotal != nil {
		v.calcCheck("managementReport.equityChanges.openingTotal", i64(ec.OpeningTotal), wantOpenTotal)
	}

	// Closing total = closing share capital + reserve fund + retained earnings + net income
	wantCloseTotal := i64(ec.ClosingShareCapital) + i64(ec.ClosingReserveFund) +
		i64(ec.ClosingRetainedEarnings) + i64(ec.ClosingNetIncome)
	if ec.ClosingTotal != nil {
		v.calcCheck("managementReport.equityChanges.closingTotal", i64(ec.ClosingTotal), wantCloseTotal)
	}
}

func (v *validator) checkProfitDispositionCalc() {
	pd := v.report.ManagementReport.ProfitDisposition

	// Total available = retained earnings + net income
	if pd.TotalAvailable != nil {
		want := i64(pd.RetainedEarnings) + i64(pd.NetIncome)
		v.calcCheck("managementReport.profitDisposition.totalAvailable", i64(pd.TotalAvailable), want)
	}

	// Total disposition = dividend + carried forward
	if pd.TotalDisposition != nil {
		want := i64(pd.Dividend) + i64(pd.CarriedForward)
		v.calcCheck("managementReport.profitDisposition.totalDisposition", i64(pd.TotalDisposition), want)
	}

	// Total available should equal total disposition
	if pd.TotalAvailable != nil && pd.TotalDisposition != nil {
		if i64(pd.TotalAvailable) != i64(pd.TotalDisposition) {
			v.err(0, "managementReport.profitDisposition",
				fmt.Sprintf("total available (%d) ≠ total disposition (%d)",
					i64(pd.TotalAvailable), i64(pd.TotalDisposition)))
		}
	}
}

// ---------------------------------------------------------------------------
// 3. Business rules (date ordering, format, semantic constraints)
// ---------------------------------------------------------------------------

var dateRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func parseDate(s string) (time.Time, bool) {
	t, err := time.Parse("2006-01-02", s)
	return t, err == nil
}

func (v *validator) checkBusinessRules() {
	r := v.report

	// Currency must be SEK or EUR (BV 1038)
	if r.Meta.Currency != "" && r.Meta.Currency != "SEK" && r.Meta.Currency != "EUR" {
		v.err(1038, "meta.currency",
			fmt.Sprintf("currency must be SEK or EUR, got %q", r.Meta.Currency))
	}

	// Language should be "sv" (BV 1116)
	if r.Meta.Language != "" && r.Meta.Language != "sv" {
		v.warn(1116, "meta.language",
			"annual report does not appear to be prepared in Swedish")
	}

	// Entry point must be one of the four valid values.
	if r.Meta.EntryPoint != "" {
		valid := map[string]bool{"risbs": true, "risab": true, "raibs": true, "raiab": true}
		if !valid[r.Meta.EntryPoint] {
			v.err(0, "meta.entryPoint",
				fmt.Sprintf("invalid entry point %q, must be one of: risbs, risab, raibs, raiab", r.Meta.EntryPoint))
		}
	}

	// OrgNr format: NNNNNN-NNNN
	if r.Company.OrgNr != "" {
		if matched, _ := regexp.MatchString(`^\d{6}-\d{4}$`, r.Company.OrgNr); !matched {
			v.err(0, "company.orgNr",
				fmt.Sprintf("organisation number %q does not match format NNNNNN-NNNN", r.Company.OrgNr))
		}
	}

	// Date format validations.
	dates := map[string]string{
		"fiscalYear.startDate":      r.FiscalYear.StartDate,
		"fiscalYear.endDate":        r.FiscalYear.EndDate,
		"certification.meetingDate": r.Certification.MeetingDate,
		"certification.signingDate": r.Certification.SigningDate,
		"signatures.date":           r.Signatures.Date,
	}
	for field, d := range dates {
		if d != "" && !dateRe.MatchString(d) {
			v.err(0, field, fmt.Sprintf("invalid date format %q, expected YYYY-MM-DD", d))
		}
	}

	// Fiscal year may not exceed 18 months (BV 1046).
	if fyStart, ok := parseDate(r.FiscalYear.StartDate); ok {
		if fyEnd, ok := parseDate(r.FiscalYear.EndDate); ok {
			months := (fyEnd.Year()-fyStart.Year())*12 + int(fyEnd.Month()-fyStart.Month())
			if fyEnd.Day() > fyStart.Day() {
				months++
			}
			if months > 18 {
				v.err(1046, "fiscalYear",
					fmt.Sprintf("fiscal year exceeds 18 months (%d months)", months))
			}
		}
	}

	// Date ordering rules — only check when both dates are valid.
	fyEnd, fyEndOK := parseDate(r.FiscalYear.EndDate)
	signDate, signOK := parseDate(r.Signatures.Date)
	meetDate, meetOK := parseDate(r.Certification.MeetingDate)
	certSignDate, certSignOK := parseDate(r.Certification.SigningDate)

	// BV 1114: signing date may not be earlier than or same as last day of FY.
	if fyEndOK && signOK && !signDate.After(fyEnd) {
		v.err(1114, "signatures.date",
			"signing date may not be earlier than or the same as the last day of the fiscal year")
	}

	// BV 1101: AGM date must be after fiscal year end.
	if fyEndOK && meetOK && !meetDate.After(fyEnd) {
		v.err(1101, "certification.meetingDate",
			"AGM date may not be earlier than or the same as the last day of the fiscal year")
	}

	// BV 1165: certification signing date may not be earlier than AGM date.
	if meetOK && certSignOK && certSignDate.Before(meetDate) {
		v.err(1165, "certification.signingDate",
			"certification signing date may not be earlier than the AGM date")
	}

	// BV 1183: AGM date should not be earlier than board signing date.
	if signOK && meetOK && meetDate.Before(signDate) {
		v.warn(1183, "certification.meetingDate",
			"AGM date is earlier than the annual report signing date")
	}

	// Comparative figures (BV 3006, 3007) — required unless first financial year.
	// We consider the report to have previous-year data if any previous-year field is non-nil.
	// If the report has some current-year data but zero previous-year data, warn.
	if r.IncomeStatement.NetResult.Current != nil && r.IncomeStatement.NetResult.Previous == nil {
		v.warn(3007, "incomeStatement",
			"comparative figures are missing in the income statement")
	}
	if r.BalanceSheet.Assets.TotalAssets.Current != nil && r.BalanceSheet.Assets.TotalAssets.Previous == nil {
		v.warn(3006, "balanceSheet",
			"comparative figures are missing in the balance sheet")
	}

	// Fixed asset note: carrying value should equal closing acquisition - closing depreciation.
	for i, note := range r.Notes.FixedAssetNotes {
		for _, label := range []string{"current", "previous"} {
			cur := label == "current"
			pick := func(yc model.YearComparison) int64 {
				if cur {
					return i64(yc.Current)
				}
				return i64(yc.Previous)
			}
			closingAcq := pick(note.ClosingAcquisitionValues)
			closingDepr := pick(note.ClosingDepreciation)
			carryVal := pick(note.CarryingValue)
			expected := closingAcq - closingDepr
			if carryVal != expected {
				v.err(0, fmt.Sprintf("notes.fixedAssetNotes[%d].carryingValue.%s", i, label),
					fmt.Sprintf("carrying value (%d) ≠ closing acquisition (%d) - closing depreciation (%d) = %d",
						carryVal, closingAcq, closingDepr, expected))
			}
		}
	}

	// Notes: accounting policies note number should be 1.
	if r.Notes.AccountingPolicies.NoteNumber != 0 && r.Notes.AccountingPolicies.NoteNumber != 1 {
		v.warn(0, "notes.accountingPolicies.noteNumber",
			fmt.Sprintf("accounting policies note number is %d, conventionally 1",
				r.Notes.AccountingPolicies.NoteNumber))
	}

	// Check entry point consistency (risbs = full IS + full BS).
	if r.Meta.EntryPoint != "" {
		ep := strings.ToLower(r.Meta.EntryPoint)
		_ = ep // Future: validate that abbreviated entry points have correct structure.
	}
}
