package sie_test

import (
	"os"
	"strings"
	"testing"

	"github.com/redofri/redofri/pkg/sie"
)

// -------------------------------------------------------------------------
// Helpers
// -------------------------------------------------------------------------

func mustParse(t *testing.T, src string) *sie.Result {
	t.Helper()
	res, err := sie.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	return res
}

func assertInt(t *testing.T, label string, got *int64, want int64) {
	t.Helper()
	if got == nil {
		t.Errorf("%s: got nil, want %d", label, want)
		return
	}
	if *got != want {
		t.Errorf("%s: got %d, want %d", label, *got, want)
	}
}

// -------------------------------------------------------------------------
// Minimal valid SIE4 fixture
// -------------------------------------------------------------------------

const minimalSIE = `#FLAGGA 0
#SIETYP 4
#FNAMN "Minimal AB"
#ORGNR 5560000001
#RAR 0 20230101 20231231
#RAR -1 20220101 20221231
`

// -------------------------------------------------------------------------
// Tests
// -------------------------------------------------------------------------

func TestParse_SIETypeValidation(t *testing.T) {
	_, err := sie.Parse(strings.NewReader(`#SIETYP 2`))
	if err == nil {
		t.Fatal("expected error for SIE type 2, got nil")
	}
	if !strings.Contains(err.Error(), "type 2") {
		t.Errorf("expected error to mention type 2, got: %v", err)
	}
}

func TestParse_CompanyAndFiscalYear(t *testing.T) {
	res := mustParse(t, minimalSIE)
	r := res.Report

	if r.Company.Name != "Minimal AB" {
		t.Errorf("Company.Name: got %q, want %q", r.Company.Name, "Minimal AB")
	}
	if r.Company.OrgNr != "556000-0001" {
		t.Errorf("Company.OrgNr: got %q, want %q", r.Company.OrgNr, "556000-0001")
	}
	if r.FiscalYear.StartDate != "2023-01-01" {
		t.Errorf("FiscalYear.StartDate: got %q, want %q", r.FiscalYear.StartDate, "2023-01-01")
	}
	if r.FiscalYear.EndDate != "2023-12-31" {
		t.Errorf("FiscalYear.EndDate: got %q, want %q", r.FiscalYear.EndDate, "2023-12-31")
	}
}

func TestParse_OrgNrFormatting(t *testing.T) {
	cases := []struct{ in, want string }{
		{"5569999999", "556999-9999"},
		{"5560000001", "556000-0001"},
		{"556999-9999", "556999-9999"}, // already has hyphen
	}
	for _, c := range cases {
		src := "#SIETYP 4\n#FNAMN \"Test\"\n#ORGNR " + c.in + "\n#RAR 0 20230101 20231231\n"
		res := mustParse(t, src)
		if res.Report.Company.OrgNr != c.want {
			t.Errorf("OrgNr(%q): got %q, want %q", c.in, res.Report.Company.OrgNr, c.want)
		}
	}
}

func TestParse_MetaDefaults(t *testing.T) {
	res := mustParse(t, minimalSIE)
	r := res.Report
	if r.Meta.Language != "sv" {
		t.Errorf("Meta.Language: got %q, want sv", r.Meta.Language)
	}
	if r.Meta.Currency != "SEK" {
		t.Errorf("Meta.Currency: got %q, want SEK", r.Meta.Currency)
	}
	if r.Meta.Country != "SE" {
		t.Errorf("Meta.Country: got %q, want SE", r.Meta.Country)
	}
	if r.Meta.AmountFormat != "NORMALFORM" {
		t.Errorf("Meta.AmountFormat: got %q, want NORMALFORM", r.Meta.AmountFormat)
	}
}

const incomeSIE = `#SIETYP 4
#FNAMN "Income AB"
#ORGNR 5560000001
#RAR 0 20230101 20231231
#RAR -1 20220101 20221231
#RES 0 3001 -1000000.00
#RES -1 3001 -800000.00
#RES 0 3910 -50000.00
#RES -1 3910 -40000.00
#RES 0 5010 200000.00
#RES -1 5010 150000.00
#RES 0 7010 300000.00
#RES -1 7010 250000.00
#RES 0 8800 50000.00
#RES -1 8800 30000.00
`

func TestParse_IncomeStatement(t *testing.T) {
	res := mustParse(t, incomeSIE)
	is := res.Report.IncomeStatement

	assertInt(t, "Revenue.NetSales.Current", is.Revenue.NetSales.Current, 1000000)
	assertInt(t, "Revenue.NetSales.Previous", is.Revenue.NetSales.Previous, 800000)
	assertInt(t, "Revenue.OtherOperatingIncome.Current", is.Revenue.OtherOperatingIncome.Current, 50000)
	assertInt(t, "Revenue.OtherOperatingIncome.Previous", is.Revenue.OtherOperatingIncome.Previous, 40000)
	assertInt(t, "Expenses.OtherExternalExpenses.Current", is.Expenses.OtherExternalExpenses.Current, 200000)
	assertInt(t, "Expenses.PersonnelExpenses.Current", is.Expenses.PersonnelExpenses.Current, 300000)
	assertInt(t, "Tax.IncomeTax.Current", is.Tax.IncomeTax.Current, 50000)
}

func TestParse_IncomeSignConvention(t *testing.T) {
	// Income accounts stored negative in SIE must be positive in model.
	src := `#SIETYP 4
#FNAMN "Sign AB"
#ORGNR 5560000001
#RAR 0 20230101 20231231
#RAR -1 20220101 20221231
#RES 0 3001 -500000.00
#RES -1 3001 -400000.00
#RES 0 8010 -200000.00
#RES -1 8010 -150000.00
#RES 0 8310 -10000.00
#RES -1 8310 -8000.00
#RES 0 8510 30000.00
#RES -1 8510 25000.00
`
	res := mustParse(t, src)
	is := res.Report.IncomeStatement

	assertInt(t, "Revenue.NetSales.Current", is.Revenue.NetSales.Current, 500000)
	assertInt(t, "FinancialItems.ResultOtherFinancialAssets.Current", is.FinancialItems.ResultOtherFinancialAssets.Current, 200000)
	assertInt(t, "FinancialItems.OtherInterestIncome.Current", is.FinancialItems.OtherInterestIncome.Current, 10000)
	assertInt(t, "FinancialItems.InterestExpenses.Current", is.FinancialItems.InterestExpenses.Current, 30000)
	// TotalFinancial = 200000 + 10000 - 30000 = 180000
	assertInt(t, "FinancialItems.TotalFinancialItems.Current", is.FinancialItems.TotalFinancialItems.Current, 180000)
}

func TestParse_Appropriations(t *testing.T) {
	src := `#SIETYP 4
#FNAMN "Appr AB"
#ORGNR 5560000001
#RAR 0 20230101 20231231
#RAR -1 20220101 20221231
#RES 0 8811 70000.00
#RES -1 8811 25000.00
#RES 0 8851 51000.00
#RES -1 8851 30000.00
`
	res := mustParse(t, src)
	is := res.Report.IncomeStatement

	assertInt(t, "Appropriations.TaxAllocationReserve.Current", is.Appropriations.TaxAllocationReserve.Current, 70000)
	assertInt(t, "Appropriations.TaxAllocationReserve.Previous", is.Appropriations.TaxAllocationReserve.Previous, 25000)
	assertInt(t, "Appropriations.ExcessDepreciation.Current", is.Appropriations.ExcessDepreciation.Current, 51000)
	assertInt(t, "Appropriations.ExcessDepreciation.Previous", is.Appropriations.ExcessDepreciation.Previous, 30000)
	assertInt(t, "Appropriations.TotalAppropriations.Current", is.Appropriations.TotalAppropriations.Current, 121000)
}

const balanceSIE = `#SIETYP 4
#FNAMN "Balance AB"
#ORGNR 5560000001
#RAR 0 20230101 20231231
#RAR -1 20220101 20221231
#UB 0 1110 500000.00
#UB -1 1110 400000.00
#UB 0 1210 200000.00
#UB -1 1210 100000.00
#UB 0 1510 300000.00
#UB -1 1510 250000.00
#UB 0 1970 50000.00
#UB -1 1970 80000.00
#UB 0 2081 -100000.00
#UB -1 2081 -100000.00
#UB 0 2091 -200000.00
#UB -1 2091 -150000.00
#UB 0 2099 -650000.00
#UB -1 2099 -480000.00
#UB 0 2350 -300000.00
#UB -1 2350 -200000.00
`

func TestParse_BalanceSheetAssets(t *testing.T) {
	res := mustParse(t, balanceSIE)
	bs := res.Report.BalanceSheet

	assertInt(t, "Tangible.BuildingsAndLand.Current", bs.Assets.FixedAssets.Tangible.BuildingsAndLand.Current, 500000)
	assertInt(t, "Tangible.BuildingsAndLand.Previous", bs.Assets.FixedAssets.Tangible.BuildingsAndLand.Previous, 400000)
	assertInt(t, "Tangible.MachineryAndEquipment.Current", bs.Assets.FixedAssets.Tangible.MachineryAndEquipment.Current, 200000)
	assertInt(t, "ShortTermReceivables.TradeReceivables.Current", bs.Assets.CurrentAssets.ShortTermReceivables.TradeReceivables.Current, 300000)
	assertInt(t, "CashAndBank.TotalCashAndBank.Current", bs.Assets.CurrentAssets.CashAndBank.TotalCashAndBank.Current, 50000)
}

func TestParse_BalanceSheetEquity(t *testing.T) {
	res := mustParse(t, balanceSIE)
	eq := res.Report.BalanceSheet.EquityAndLiabilities.Equity

	// Liabilities stored negative in SIE → negated to positive in model
	assertInt(t, "Equity.ShareCapital.Current", eq.ShareCapital.Current, 100000)
	assertInt(t, "Equity.RetainedEarnings.Current", eq.RetainedEarnings.Current, 200000)
	assertInt(t, "Equity.NetIncome.Current", eq.NetIncome.Current, 650000)
	assertInt(t, "Equity.TotalEquity.Current", eq.TotalEquity.Current, 950000)
}

func TestParse_BalanceSheetLiabilities(t *testing.T) {
	res := mustParse(t, balanceSIE)
	ll := res.Report.BalanceSheet.EquityAndLiabilities.LongTermLiabilities

	assertInt(t, "LongTermLiabilities.BankLoans.Current", ll.BankLoans.Current, 300000)
	assertInt(t, "LongTermLiabilities.BankLoans.Previous", ll.BankLoans.Previous, 200000)
}

func TestParse_UntaxedReserves(t *testing.T) {
	src := `#SIETYP 4
#FNAMN "Reserves AB"
#ORGNR 5560000001
#RAR 0 20230101 20231231
#RAR -1 20220101 20221231
#UB 0 2112 -169000.00
#UB -1 2112 -99000.00
#UB 0 2153 -121000.00
#UB -1 2153 -70000.00
`
	res := mustParse(t, src)
	ur := res.Report.BalanceSheet.EquityAndLiabilities.UntaxedReserves

	assertInt(t, "TaxAllocationReserves.Current", ur.TaxAllocationReserves.Current, 169000)
	assertInt(t, "TaxAllocationReserves.Previous", ur.TaxAllocationReserves.Previous, 99000)
	assertInt(t, "AccumulatedExcessDepreciation.Current", ur.AccumulatedExcessDepreciation.Current, 121000)
	assertInt(t, "TotalUntaxedReserves.Current", ur.TotalUntaxedReserves.Current, 290000)
}

func TestParse_Provisions(t *testing.T) {
	src := `#SIETYP 4
#FNAMN "Prov AB"
#ORGNR 5560000001
#RAR 0 20230101 20231231
#RAR -1 20220101 20221231
#UB 0 2210 -770000.00
#UB -1 2210 -650000.00
#UB 0 2250 -100000.00
#UB -1 2250 -65000.00
`
	res := mustParse(t, src)
	pr := res.Report.BalanceSheet.EquityAndLiabilities.Provisions

	assertInt(t, "PensionProvisions.Current", pr.PensionProvisions.Current, 770000)
	assertInt(t, "OtherProvisions.Current", pr.OtherProvisions.Current, 100000)
	assertInt(t, "TotalProvisions.Current", pr.TotalProvisions.Current, 870000)
}

func TestParse_ShortTermLiabilities(t *testing.T) {
	src := `#SIETYP 4
#FNAMN "STLiab AB"
#ORGNR 5560000001
#RAR 0 20230101 20231231
#RAR -1 20220101 20221231
#UB 0 2440 -855000.00
#UB -1 2440 -641000.00
#UB 0 2510 -130000.00
#UB -1 2510 -35000.00
#UB 0 2650 -492000.00
#UB -1 2650 -315000.00
#UB 0 2910 -453000.00
#UB -1 2910 -224000.00
`
	res := mustParse(t, src)
	sl := res.Report.BalanceSheet.EquityAndLiabilities.ShortTermLiabilities

	assertInt(t, "TradePayables.Current", sl.TradePayables.Current, 855000)
	assertInt(t, "TaxLiabilities.Current", sl.TaxLiabilities.Current, 130000)
	assertInt(t, "OtherShortTermLiabilities.Current", sl.OtherShortTermLiabilities.Current, 492000)
	assertInt(t, "AccruedExpenses.Current", sl.AccruedExpenses.Current, 453000)
	assertInt(t, "TotalShortTermLiabilities.Current", sl.TotalShortTermLiabilities.Current, 1930000)
}

func TestParse_TokeniserQuotedStrings(t *testing.T) {
	// Account names with spaces should be parsed correctly.
	src := `#SIETYP 4
#FNAMN "Bolaget med mellanslag AB"
#ORGNR 5560000001
#RAR 0 20230101 20231231
#KONTO 3001 "Försäljning varor och tjänster"
#RES 0 3001 -1000000.00
`
	res := mustParse(t, src)
	if res.Report.Company.Name != "Bolaget med mellanslag AB" {
		t.Errorf("Company.Name: got %q, want %q", res.Report.Company.Name, "Bolaget med mellanslag AB")
	}
	assertInt(t, "NetSales.Current", res.Report.IncomeStatement.Revenue.NetSales.Current, 1000000)
}

func TestParse_IgnoresVoucherLines(t *testing.T) {
	// #VER and #TRANS records should be silently ignored.
	src := `#SIETYP 4
#FNAMN "Voucher AB"
#ORGNR 5560000001
#RAR 0 20230101 20231231
#RES 0 3001 -500000.00
#VER A 1 20230315 "Faktura"
{
#TRANS 3001 {} -500000.00
#TRANS 1510 {} 500000.00
}
`
	res := mustParse(t, src)
	assertInt(t, "NetSales.Current", res.Report.IncomeStatement.Revenue.NetSales.Current, 500000)
}

// TestParse_Exempel1SIE is an integration test that parses the full synthetic
// SIE4 file in testdata/exempel1.sie and verifies it produces the same
// numerical values as testdata/exempel1.json.
func TestParse_Exempel1SIE(t *testing.T) {
	f, err := os.Open("../../testdata/exempel1.sie")
	if err != nil {
		t.Fatalf("open testdata/exempel1.sie: %v", err)
	}
	defer f.Close()

	res, err := sie.Parse(f)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(res.Warnings) > 0 {
		for _, w := range res.Warnings {
			t.Logf("warning: %s", w)
		}
	}

	r := res.Report

	// Company
	if r.Company.Name != "Exempel 1 AB" {
		t.Errorf("Company.Name: got %q, want %q", r.Company.Name, "Exempel 1 AB")
	}
	if r.Company.OrgNr != "556999-9999" {
		t.Errorf("Company.OrgNr: got %q, want %q", r.Company.OrgNr, "556999-9999")
	}
	if r.FiscalYear.StartDate != "2016-01-01" {
		t.Errorf("FiscalYear.StartDate: got %q, want %q", r.FiscalYear.StartDate, "2016-01-01")
	}

	is := r.IncomeStatement
	bs := r.BalanceSheet

	// Income Statement
	assertInt(t, "IS.Revenue.NetSales.Current", is.Revenue.NetSales.Current, 2650000)
	assertInt(t, "IS.Revenue.NetSales.Previous", is.Revenue.NetSales.Previous, 2250000)
	assertInt(t, "IS.Revenue.InventoryChange.Current", is.Revenue.InventoryChange.Current, 700000)
	assertInt(t, "IS.Revenue.InventoryChange.Previous", is.Revenue.InventoryChange.Previous, 1125000)
	assertInt(t, "IS.Revenue.OtherOperatingIncome.Current", is.Revenue.OtherOperatingIncome.Current, 377000)
	assertInt(t, "IS.Revenue.OtherOperatingIncome.Previous", is.Revenue.OtherOperatingIncome.Previous, 1000000)
	assertInt(t, "IS.Revenue.TotalRevenue.Current", is.Revenue.TotalRevenue.Current, 3727000)

	assertInt(t, "IS.Expenses.RawMaterials.Current", is.Expenses.RawMaterials.Current, 1520000)
	assertInt(t, "IS.Expenses.TradingGoods.Current", is.Expenses.TradingGoods.Current, 308000)
	assertInt(t, "IS.Expenses.OtherExternalExpenses.Current", is.Expenses.OtherExternalExpenses.Current, 499000)
	assertInt(t, "IS.Expenses.PersonnelExpenses.Current", is.Expenses.PersonnelExpenses.Current, 650000)
	assertInt(t, "IS.Expenses.DepreciationAmortization.Current", is.Expenses.DepreciationAmortization.Current, 340000)
	assertInt(t, "IS.Expenses.OtherOperatingExpenses.Current", is.Expenses.OtherOperatingExpenses.Current, 205000)
	assertInt(t, "IS.Expenses.TotalExpenses.Current", is.Expenses.TotalExpenses.Current, 3522000)

	assertInt(t, "IS.OperatingResult.Current", is.OperatingResult.Current, 205000)

	assertInt(t, "IS.FinancialItems.ResultOtherFinancialAssets.Current", is.FinancialItems.ResultOtherFinancialAssets.Current, 1543000)
	assertInt(t, "IS.FinancialItems.OtherInterestIncome.Current", is.FinancialItems.OtherInterestIncome.Current, 12000)
	assertInt(t, "IS.FinancialItems.InterestExpenses.Current", is.FinancialItems.InterestExpenses.Current, 275000)
	assertInt(t, "IS.FinancialItems.TotalFinancialItems.Current", is.FinancialItems.TotalFinancialItems.Current, 1280000)

	assertInt(t, "IS.ResultAfterFinancialItems.Current", is.ResultAfterFinancialItems.Current, 1485000)

	assertInt(t, "IS.Appropriations.TaxAllocationReserve.Current", is.Appropriations.TaxAllocationReserve.Current, 70000)
	assertInt(t, "IS.Appropriations.ExcessDepreciation.Current", is.Appropriations.ExcessDepreciation.Current, 51000)
	assertInt(t, "IS.Appropriations.TotalAppropriations.Current", is.Appropriations.TotalAppropriations.Current, 121000)

	assertInt(t, "IS.ResultBeforeTax.Current", is.ResultBeforeTax.Current, 1364000)
	assertInt(t, "IS.Tax.IncomeTax.Current", is.Tax.IncomeTax.Current, 90000)
	assertInt(t, "IS.NetResult.Current", is.NetResult.Current, 1274000)
	assertInt(t, "IS.NetResult.Previous", is.NetResult.Previous, 1099000)

	// Balance sheet — assets
	assertInt(t, "BS.Tangible.BuildingsAndLand.Current", bs.Assets.FixedAssets.Tangible.BuildingsAndLand.Current, 1620000)
	assertInt(t, "BS.Tangible.BuildingsAndLand.Previous", bs.Assets.FixedAssets.Tangible.BuildingsAndLand.Previous, 1450000)
	assertInt(t, "BS.Tangible.MachineryAndEquipment.Current", bs.Assets.FixedAssets.Tangible.MachineryAndEquipment.Current, 860000)
	assertInt(t, "BS.Tangible.FixturesAndFittings.Current", bs.Assets.FixedAssets.Tangible.FixturesAndFittings.Current, 240000)
	assertInt(t, "BS.Tangible.TotalTangible.Current", bs.Assets.FixedAssets.Tangible.TotalTangible.Current, 2720000)

	assertInt(t, "BS.Inventory.RawMaterials.Current", bs.Assets.CurrentAssets.Inventory.RawMaterials.Current, 510000)
	assertInt(t, "BS.Inventory.WorkInProgress.Current", bs.Assets.CurrentAssets.Inventory.WorkInProgress.Current, 240000)
	assertInt(t, "BS.Inventory.FinishedGoods.Current", bs.Assets.CurrentAssets.Inventory.FinishedGoods.Current, 750000)
	assertInt(t, "BS.Inventory.TotalInventory.Current", bs.Assets.CurrentAssets.Inventory.TotalInventory.Current, 1500000)

	assertInt(t, "BS.ShortTermReceivables.TradeReceivables.Current", bs.Assets.CurrentAssets.ShortTermReceivables.TradeReceivables.Current, 1393000)
	assertInt(t, "BS.ShortTermReceivables.OtherReceivables.Current", bs.Assets.CurrentAssets.ShortTermReceivables.OtherReceivables.Current, 20000)
	assertInt(t, "BS.ShortTermReceivables.PrepaidExpenses.Current", bs.Assets.CurrentAssets.ShortTermReceivables.PrepaidExpenses.Current, 30000)

	assertInt(t, "BS.CashAndBank.TotalCashAndBank.Current", bs.Assets.CurrentAssets.CashAndBank.TotalCashAndBank.Current, 110000)

	assertInt(t, "BS.TotalAssets.Current", bs.Assets.TotalAssets.Current, 7773000)

	// Balance sheet — equity & liabilities
	eq := bs.EquityAndLiabilities.Equity
	assertInt(t, "BS.Equity.ShareCapital.Current", eq.ShareCapital.Current, 100000)
	assertInt(t, "BS.Equity.ReserveFund.Current", eq.ReserveFund.Current, 5000)
	assertInt(t, "BS.Equity.RetainedEarnings.Current", eq.RetainedEarnings.Current, 1011000)
	assertInt(t, "BS.Equity.NetIncome.Current", eq.NetIncome.Current, 1274000)
	assertInt(t, "BS.Equity.TotalEquity.Current", eq.TotalEquity.Current, 2390000)

	ur := bs.EquityAndLiabilities.UntaxedReserves
	assertInt(t, "BS.UntaxedReserves.TaxAllocationReserves.Current", ur.TaxAllocationReserves.Current, 169000)
	assertInt(t, "BS.UntaxedReserves.AccumulatedExcessDepreciation.Current", ur.AccumulatedExcessDepreciation.Current, 121000)
	assertInt(t, "BS.UntaxedReserves.TotalUntaxedReserves.Current", ur.TotalUntaxedReserves.Current, 290000)

	pr := bs.EquityAndLiabilities.Provisions
	assertInt(t, "BS.Provisions.PensionProvisions.Current", pr.PensionProvisions.Current, 770000)
	assertInt(t, "BS.Provisions.OtherProvisions.Current", pr.OtherProvisions.Current, 100000)
	assertInt(t, "BS.Provisions.TotalProvisions.Current", pr.TotalProvisions.Current, 870000)

	ll := bs.EquityAndLiabilities.LongTermLiabilities
	assertInt(t, "BS.LongTermLiabilities.BankLoans.Current", ll.BankLoans.Current, 2193000)
	assertInt(t, "BS.LongTermLiabilities.OtherLongTermLiabilities.Current", ll.OtherLongTermLiabilities.Current, 100000)
	assertInt(t, "BS.LongTermLiabilities.TotalLongTermLiabilities.Current", ll.TotalLongTermLiabilities.Current, 2293000)

	sl := bs.EquityAndLiabilities.ShortTermLiabilities
	assertInt(t, "BS.ShortTermLiabilities.TradePayables.Current", sl.TradePayables.Current, 855000)
	assertInt(t, "BS.ShortTermLiabilities.TaxLiabilities.Current", sl.TaxLiabilities.Current, 130000)
	assertInt(t, "BS.ShortTermLiabilities.OtherShortTermLiabilities.Current", sl.OtherShortTermLiabilities.Current, 492000)
	assertInt(t, "BS.ShortTermLiabilities.AccruedExpenses.Current", sl.AccruedExpenses.Current, 453000)
	assertInt(t, "BS.ShortTermLiabilities.TotalShortTermLiabilities.Current", sl.TotalShortTermLiabilities.Current, 1930000)

	assertInt(t, "BS.TotalEquityAndLiabilities.Current", bs.EquityAndLiabilities.TotalEquityAndLiabilities.Current, 7773000)
}
