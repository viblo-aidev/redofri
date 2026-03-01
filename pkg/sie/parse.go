// Package sie parses SIE Type 4 export files and maps account balances
// to a partial model.AnnualReport.
//
// SIE4 is the Swedish standard export format used by virtually all Swedish
// accounting software. It contains account definitions, opening/closing
// balances, and P&L results by fiscal year, which is enough to populate
// the numerical fields of the income statement and balance sheet.
//
// Fields that cannot be derived from SIE (text sections, note descriptions,
// asset roll-forwards, signatures, certification) are left at their zero
// values and must be supplied from other sources (manual JSON, previous
// year iXBRL, etc.).
package sie

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	"github.com/redofri/redofri/pkg/model"
)

// Result holds the parsed annual report together with any non-fatal
// warnings encountered during parsing (e.g. unknown account types).
type Result struct {
	Report   *model.AnnualReport
	Warnings []string
}

// Parse reads a SIE Type 4 file from r (UTF-8 or CP437/PC8) and returns a
// partial AnnualReport populated from the account balances found in the file.
//
// The caller must supply the remaining fields (text, notes, signatures, etc.)
// before the report is complete.
func Parse(r io.Reader) (*Result, error) {
	p := &parser{
		accounts: make(map[string]account),
		balances: make(map[yearAccount]int64),
	}
	if err := p.scan(r); err != nil {
		return nil, err
	}
	report, warnings := p.build()
	return &Result{Report: report, Warnings: warnings}, nil
}

// ---------------------------------------------------------------------------
// Internal types
// ---------------------------------------------------------------------------

type yearAccount struct {
	year    int    // 0 = current, -1 = previous
	account string // BAS account number, e.g. "3000"
}

type account struct {
	number string
	name   string
	ktype  string // T=asset, S=liability, K=cost, I=income
}

type fiscalYear struct {
	index     int    // 0 = current, -1 = previous
	startDate string // YYYYMMDD
	endDate   string // YYYYMMDD
}

type parser struct {
	sieTyp   int
	compName string
	orgNr    string
	years    []fiscalYear

	accounts map[string]account
	// balances stores amounts using SIE sign convention (income negative, liabilities negative)
	balances map[yearAccount]int64
}

// ---------------------------------------------------------------------------
// Scanning / tokenising
// ---------------------------------------------------------------------------

// scan reads all lines from r, auto-detecting CP437 encoding if needed.
func (p *parser) scan(r io.Reader) error {
	// Read all bytes first so we can detect encoding.
	raw, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("sie: reading input: %w", err)
	}

	// Decode CP437 if the content is not valid UTF-8.
	decoded, err := decodeBytes(raw)
	if err != nil {
		return fmt.Errorf("sie: decoding input: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(decoded))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == '{' || line[0] == '}' {
			continue
		}
		if line[0] != '#' {
			continue
		}
		keyword, fields := tokenise(line)
		if err := p.dispatch(keyword, fields); err != nil {
			return err
		}
	}
	return scanner.Err()
}

// decodeBytes returns the string decoded as UTF-8, falling back to CP437.
func decodeBytes(raw []byte) (string, error) {
	// Check if valid UTF-8 (includes ASCII).
	if isValidUTF8(raw) {
		return string(raw), nil
	}
	// Decode as CP437 (IBM PC character set, historically used by Swedish SIE files).
	decoded, _, err := transform.Bytes(charmap.CodePage437.NewDecoder(), raw)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// isValidUTF8 checks whether b is valid UTF-8.
func isValidUTF8(b []byte) bool {
	s := string(b)
	for _, r := range s {
		if r == '\uFFFD' {
			return false
		}
	}
	return true
}

// tokenise splits a SIE record line into keyword and fields.
// Fields are whitespace-separated; quoted strings are kept intact.
func tokenise(line string) (keyword string, fields []string) {
	// Drop the leading '#'.
	line = line[1:]

	// Split on first whitespace to get keyword.
	i := 0
	for i < len(line) && line[i] != ' ' && line[i] != '\t' {
		i++
	}
	keyword = strings.ToUpper(line[:i])
	rest := strings.TrimSpace(line[i:])

	// Parse fields (whitespace-separated, honouring double-quoted strings).
	fields = parseFields(rest)
	return keyword, fields
}

// parseFields parses a whitespace-separated list of SIE field values,
// respecting double-quoted strings (which may contain spaces).
func parseFields(s string) []string {
	var fields []string
	s = strings.TrimSpace(s)
	for len(s) > 0 {
		var field string
		if s[0] == '"' {
			// Quoted string: find closing quote.
			end := strings.Index(s[1:], "\"")
			if end < 0 {
				// Unterminated quote — take rest of line.
				field = s[1:]
				s = ""
			} else {
				field = s[1 : end+1]
				s = strings.TrimSpace(s[end+2:])
			}
		} else if s[0] == '{' {
			// Object field (account dimensions etc.) — treat as single token.
			end := strings.Index(s, "}")
			if end < 0 {
				field = s
				s = ""
			} else {
				field = s[:end+1]
				s = strings.TrimSpace(s[end+1:])
			}
		} else {
			// Plain token.
			end := strings.IndexAny(s, " \t")
			if end < 0 {
				field = s
				s = ""
			} else {
				field = s[:end]
				s = strings.TrimSpace(s[end+1:])
			}
		}
		fields = append(fields, field)
	}
	return fields
}

// ---------------------------------------------------------------------------
// Record dispatch
// ---------------------------------------------------------------------------

func (p *parser) dispatch(keyword string, fields []string) error {
	switch keyword {
	case "SIETYP":
		return p.handleSIETYP(fields)
	case "FNAMN":
		return p.handleFNAMN(fields)
	case "ORGNR":
		return p.handleORGNR(fields)
	case "RAR":
		return p.handleRAR(fields)
	case "KONTO":
		return p.handleKONTO(fields)
	case "KTYP":
		return p.handleKTYP(fields)
	case "IB":
		return p.handleBalance("IB", fields)
	case "UB":
		return p.handleBalance("UB", fields)
	case "RES":
		return p.handleRES(fields)
	// Deliberately ignored (not needed for our mapping):
	case "VER", "TRANS", "RTRANS", "BTRANS", "DIM", "ENHET",
		"FLAGGA", "FORMAT", "GEN", "PROGRAM", "PROSA",
		"KPTYP", "OBJEKT", "SRU", "TAXAR", "OMFATTN",
		"ADRESS", "BKOD", "VALUTA", "UNDERDIM":
		// ignore
	}
	return nil
}

func (p *parser) handleSIETYP(fields []string) error {
	if len(fields) < 1 {
		return fmt.Errorf("sie: #SIETYP missing value")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("sie: #SIETYP invalid: %w", err)
	}
	p.sieTyp = t
	if t != 4 {
		return fmt.Errorf("sie: only SIE type 4 is supported, got type %d", t)
	}
	return nil
}

func (p *parser) handleFNAMN(fields []string) error {
	if len(fields) >= 1 {
		p.compName = fields[0]
	}
	return nil
}

func (p *parser) handleORGNR(fields []string) error {
	if len(fields) >= 1 {
		// SIE stores org nr without hyphen, e.g. "5569999999"
		// Model uses "556999-9999" (hyphen before last 4 digits).
		p.orgNr = formatOrgNr(fields[0])
	}
	return nil
}

// formatOrgNr converts "5569999999" → "556999-9999".
func formatOrgNr(s string) string {
	// Remove any existing hyphens.
	s = strings.ReplaceAll(s, "-", "")
	if len(s) == 10 {
		return s[:6] + "-" + s[6:]
	}
	return s
}

func (p *parser) handleRAR(fields []string) error {
	// #RAR  <yearIndex>  <startDate>  <endDate>
	if len(fields) < 3 {
		return nil // tolerate malformed
	}
	idx, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil
	}
	p.years = append(p.years, fiscalYear{
		index:     idx,
		startDate: sieDate(fields[1]),
		endDate:   sieDate(fields[2]),
	})
	return nil
}

// sieDate converts SIE date "20230101" → "2023-01-01".
func sieDate(s string) string {
	if len(s) == 8 {
		return s[:4] + "-" + s[4:6] + "-" + s[6:8]
	}
	return s
}

func (p *parser) handleKONTO(fields []string) error {
	// #KONTO  <accountNr>  <name>
	if len(fields) < 2 {
		return nil
	}
	nr := fields[0]
	a := p.accounts[nr]
	a.number = nr
	a.name = fields[1]
	p.accounts[nr] = a
	return nil
}

func (p *parser) handleKTYP(fields []string) error {
	// #KTYP  <accountNr>  <type>   (T=asset, S=liability, K=cost, I=income)
	if len(fields) < 2 {
		return nil
	}
	nr := fields[0]
	a := p.accounts[nr]
	a.number = nr
	a.ktype = strings.ToUpper(fields[1])
	p.accounts[nr] = a
	return nil
}

func (p *parser) handleBalance(keyword string, fields []string) error {
	// #IB  <yearIndex>  <accountNr>  <amount>  [<quantity>]
	// #UB  <yearIndex>  <accountNr>  <amount>  [<quantity>]
	if len(fields) < 3 {
		return nil
	}
	idx, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil
	}
	nr := fields[1]
	amount, err := parseAmount(fields[2])
	if err != nil {
		return nil
	}
	// We only store UB (closing balances) for year 0 and -1.
	// IB for year 0 equals UB for year -1, so we only need UB.
	if keyword == "UB" {
		p.balances[yearAccount{year: idx, account: nr}] = amount
	}
	return nil
}

func (p *parser) handleRES(fields []string) error {
	// #RES  <yearIndex>  <accountNr>  <amount>
	if len(fields) < 3 {
		return nil
	}
	idx, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil
	}
	nr := fields[1]
	amount, err := parseAmount(fields[2])
	if err != nil {
		return nil
	}
	p.balances[yearAccount{year: idx, account: nr}] = amount
	return nil
}

func parseAmount(s string) (int64, error) {
	// SIE amounts are decimal with up to 2 decimal places, e.g. "125000.00"
	// We convert to integer öre and then to whole kronor (round to nearest).
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	// Round to nearest krona.
	if f >= 0 {
		return int64(f + 0.5), nil
	}
	return int64(f - 0.5), nil
}

// ---------------------------------------------------------------------------
// Balance accumulator helpers
// ---------------------------------------------------------------------------

// sumRange accumulates balances for accounts in [lo, hi] for the given year.
func (p *parser) sumRange(year int, lo, hi int) int64 {
	var total int64
	for nr, amount := range p.balances {
		if nr.year != year {
			continue
		}
		n, err := strconv.Atoi(nr.account)
		if err != nil {
			continue
		}
		if n >= lo && n <= hi {
			total += amount
		}
	}
	return total
}

// sum accumulates balances for the given exact account numbers for the given year.
func (p *parser) sum(year int, accounts ...string) int64 {
	var total int64
	for _, a := range accounts {
		total += p.balances[yearAccount{year: year, account: a}]
	}
	return total
}

// yc returns a YearComparison for the given current/previous raw SIE amounts,
// optionally negating (for income and liability accounts).
func yc(cur, prev int64, negate bool) model.YearComparison {
	if negate {
		cur = -cur
		prev = -prev
	}
	return model.YearComparison{
		Current:  model.Int64(cur),
		Previous: model.Int64(prev),
	}
}

// ycNeg is shorthand for yc with negate=true (income/liability accounts).
func ycNeg(cur, prev int64) model.YearComparison {
	return yc(cur, prev, true)
}

// ycPos is shorthand for yc with negate=false (asset/expense accounts).
func ycPos(cur, prev int64) model.YearComparison {
	return yc(cur, prev, false)
}

// ---------------------------------------------------------------------------
// Model builder
// ---------------------------------------------------------------------------

func (p *parser) build() (*model.AnnualReport, []string) {
	var warnings []string
	report := &model.AnnualReport{}

	// --- Company ---
	report.Company.Name = p.compName
	report.Company.OrgNr = p.orgNr

	// --- Fiscal years ---
	for _, fy := range p.years {
		switch fy.index {
		case 0:
			report.FiscalYear.StartDate = fy.startDate
			report.FiscalYear.EndDate = fy.endDate
		}
	}

	// --- Meta defaults ---
	report.Meta.Language = "sv"
	report.Meta.Country = "SE"
	report.Meta.Currency = "SEK"
	report.Meta.AmountFormat = "NORMALFORM"
	// EntryPoint must be set by caller.

	// -----------------------------------------------------------------------
	// Income Statement
	// -----------------------------------------------------------------------

	// Revenue (3000–3799 = net sales, 3800–3899 = inventory changes,
	// 3900–3999 = other operating income)
	// All income accounts stored negative in SIE → negate for model.
	netSalesCur := -p.sumRange(0, 3000, 3799)
	netSalesPrev := -p.sumRange(-1, 3000, 3799)

	invChangeCur := -p.sumRange(0, 3800, 3899)
	invChangePrev := -p.sumRange(-1, 3800, 3899)

	otherOpIncomeCur := -p.sumRange(0, 3900, 3999)
	otherOpIncomePrev := -p.sumRange(-1, 3900, 3999)

	report.IncomeStatement.Revenue.NetSales = ycPos(netSalesCur, netSalesPrev)

	if invChangeCur != 0 || invChangePrev != 0 {
		report.IncomeStatement.Revenue.InventoryChange = ycPos(invChangeCur, invChangePrev)
	}
	if otherOpIncomeCur != 0 || otherOpIncomePrev != 0 {
		report.IncomeStatement.Revenue.OtherOperatingIncome = ycPos(otherOpIncomeCur, otherOpIncomePrev)
	}

	totalRevCur := netSalesCur + invChangeCur + otherOpIncomeCur
	totalRevPrev := netSalesPrev + invChangePrev + otherOpIncomePrev
	report.IncomeStatement.Revenue.TotalRevenue = ycPos(totalRevCur, totalRevPrev)

	// Expenses
	// 4000–4699: raw materials / trading goods
	rawMatCur := p.sumRange(0, 4000, 4499)
	rawMatPrev := p.sumRange(-1, 4000, 4499)
	tradingCur := p.sumRange(0, 4500, 4699)
	tradingPrev := p.sumRange(-1, 4500, 4699)

	// 5000–6999: other external expenses
	otherExtCur := p.sumRange(0, 5000, 6999)
	otherExtPrev := p.sumRange(-1, 5000, 6999)

	// 7000–7699: personnel expenses
	personnelCur := p.sumRange(0, 7000, 7699)
	personnelPrev := p.sumRange(-1, 7000, 7699)

	// 7800–7899: depreciation / amortisation
	deprCur := p.sumRange(0, 7800, 7899)
	deprPrev := p.sumRange(-1, 7800, 7899)

	// 7900–7999: other operating expenses
	otherOpExpCur := p.sumRange(0, 7900, 7999)
	otherOpExpPrev := p.sumRange(-1, 7900, 7999)

	if rawMatCur != 0 || rawMatPrev != 0 {
		report.IncomeStatement.Expenses.RawMaterials = ycPos(rawMatCur, rawMatPrev)
	}
	if tradingCur != 0 || tradingPrev != 0 {
		report.IncomeStatement.Expenses.TradingGoods = ycPos(tradingCur, tradingPrev)
	}
	if otherExtCur != 0 || otherExtPrev != 0 {
		report.IncomeStatement.Expenses.OtherExternalExpenses = ycPos(otherExtCur, otherExtPrev)
	}
	if personnelCur != 0 || personnelPrev != 0 {
		report.IncomeStatement.Expenses.PersonnelExpenses = ycPos(personnelCur, personnelPrev)
	}
	if deprCur != 0 || deprPrev != 0 {
		report.IncomeStatement.Expenses.DepreciationAmortization = ycPos(deprCur, deprPrev)
	}
	if otherOpExpCur != 0 || otherOpExpPrev != 0 {
		report.IncomeStatement.Expenses.OtherOperatingExpenses = ycPos(otherOpExpCur, otherOpExpPrev)
	}

	totalExpCur := rawMatCur + tradingCur + otherExtCur + personnelCur + deprCur + otherOpExpCur
	totalExpPrev := rawMatPrev + tradingPrev + otherExtPrev + personnelPrev + deprPrev + otherOpExpPrev
	report.IncomeStatement.Expenses.TotalExpenses = ycPos(totalExpCur, totalExpPrev)

	// Operating result
	opResCur := totalRevCur - totalExpCur
	opResPrev := totalRevPrev - totalExpPrev
	report.IncomeStatement.OperatingResult = ycPos(opResCur, opResPrev)

	// Financial items
	// 8000–8099: result from participations in group companies (income, negate)
	// 8100–8199: result from participations in associated companies (income, negate)
	// 8200–8299: result from other long-term securities (income, negate)
	// We map 8000–8299 to ResultOtherFinancialAssets for simplicity.
	resOtherFinCur := -p.sumRange(0, 8000, 8299)
	resOtherFinPrev := -p.sumRange(-1, 8000, 8299)

	// 8300–8499: interest income and similar (income, negate)
	otherIntIncomeCur := -p.sumRange(0, 8300, 8499)
	otherIntIncomePrev := -p.sumRange(-1, 8300, 8499)

	// 8500–8799: interest expenses and similar (cost, positive)
	intExpCur := p.sumRange(0, 8500, 8799)
	intExpPrev := p.sumRange(-1, 8500, 8799)

	if resOtherFinCur != 0 || resOtherFinPrev != 0 {
		report.IncomeStatement.FinancialItems.ResultOtherFinancialAssets = ycPos(resOtherFinCur, resOtherFinPrev)
	}
	if otherIntIncomeCur != 0 || otherIntIncomePrev != 0 {
		report.IncomeStatement.FinancialItems.OtherInterestIncome = ycPos(otherIntIncomeCur, otherIntIncomePrev)
	}
	if intExpCur != 0 || intExpPrev != 0 {
		report.IncomeStatement.FinancialItems.InterestExpenses = ycPos(intExpCur, intExpPrev)
	}

	totalFinCur := resOtherFinCur + otherIntIncomeCur - intExpCur
	totalFinPrev := resOtherFinPrev + otherIntIncomePrev - intExpPrev
	report.IncomeStatement.FinancialItems.TotalFinancialItems = ycPos(totalFinCur, totalFinPrev)

	// Result after financial items
	resAfterFinCur := opResCur + totalFinCur
	resAfterFinPrev := opResPrev + totalFinPrev
	report.IncomeStatement.ResultAfterFinancialItems = ycPos(resAfterFinCur, resAfterFinPrev)

	// Appropriations (bokslutsdispositioner)
	// 8810–8819: periodiseringsfond (tax allocation reserve) — income when reversed (negate), expense when allocated (positive in SIE)
	// Net SIE amount: positive = expense (more reserves), negative = income (release)
	taxAllocCur := p.sumRange(0, 8810, 8819)
	taxAllocPrev := p.sumRange(-1, 8810, 8819)
	// Sign in model: positive = expense (allocated), negative = income (released)
	// SIE sign: same convention for expense accounts. We keep as-is but negate to match model
	// (model: positive means cost/allocation that reduces result).
	// Actually SIE 8810–8819: avsättning to periodiseringsfond recorded as positive (expense).
	// Model taxAllocationReserve: positive value means allocated (reduces result before tax).
	// We keep as-is.

	// 8850–8859: overavskrivningar (excess depreciation) — positive expense in SIE
	excessDeprCur := p.sumRange(0, 8850, 8859)
	excessDeprPrev := p.sumRange(-1, 8850, 8859)

	if taxAllocCur != 0 || taxAllocPrev != 0 {
		report.IncomeStatement.Appropriations.TaxAllocationReserve = ycPos(taxAllocCur, taxAllocPrev)
	}
	if excessDeprCur != 0 || excessDeprPrev != 0 {
		report.IncomeStatement.Appropriations.ExcessDepreciation = ycPos(excessDeprCur, excessDeprPrev)
	}

	totalApprCur := taxAllocCur + excessDeprCur
	totalApprPrev := taxAllocPrev + excessDeprPrev
	report.IncomeStatement.Appropriations.TotalAppropriations = ycPos(totalApprCur, totalApprPrev)

	// Result before tax
	resBeforeTaxCur := resAfterFinCur - totalApprCur
	resBeforeTaxPrev := resAfterFinPrev - totalApprPrev
	report.IncomeStatement.ResultBeforeTax = ycPos(resBeforeTaxCur, resBeforeTaxPrev)

	// Tax (8800–8809: income tax, positive expense in SIE)
	taxCur := p.sumRange(0, 8800, 8809)
	taxPrev := p.sumRange(-1, 8800, 8809)
	report.IncomeStatement.Tax.IncomeTax = ycPos(taxCur, taxPrev)

	// Net result
	netResCur := resBeforeTaxCur - taxCur
	netResPrev := resBeforeTaxPrev - taxPrev
	report.IncomeStatement.NetResult = ycPos(netResCur, netResPrev)

	// -----------------------------------------------------------------------
	// Balance Sheet — Assets
	// -----------------------------------------------------------------------

	// Fixed assets — tangible
	// 1100–1199: buildings and land
	bldCur := p.sumRange(0, 1100, 1199)
	bldPrev := p.sumRange(-1, 1100, 1199)
	// 1200–1299: machinery and equipment
	machCur := p.sumRange(0, 1200, 1299)
	machPrev := p.sumRange(-1, 1200, 1299)
	// 1300–1349: fixtures and fittings (inventarier, verktyg och installationer)
	fixCur := p.sumRange(0, 1300, 1349)
	fixPrev := p.sumRange(-1, 1300, 1349)

	if bldCur != 0 || bldPrev != 0 {
		report.BalanceSheet.Assets.FixedAssets.Tangible.BuildingsAndLand = ycPos(bldCur, bldPrev)
	}
	if machCur != 0 || machPrev != 0 {
		report.BalanceSheet.Assets.FixedAssets.Tangible.MachineryAndEquipment = ycPos(machCur, machPrev)
	}
	if fixCur != 0 || fixPrev != 0 {
		report.BalanceSheet.Assets.FixedAssets.Tangible.FixturesAndFittings = ycPos(fixCur, fixPrev)
	}

	totalTangCur := bldCur + machCur + fixCur
	totalTangPrev := bldPrev + machPrev + fixPrev
	report.BalanceSheet.Assets.FixedAssets.Tangible.TotalTangible = ycPos(totalTangCur, totalTangPrev)

	// Fixed assets — financial
	// 1350–1399: financial fixed assets (participations, long-term securities, etc.)
	// BAS examples: 1350–1359 andelar, 1360–1369 långfristiga fordringar,
	// 1380–1389 andra finansiella anläggningstillgångar.
	// We map the whole range to OtherLongTermSecurities as a simplification.
	longSecCur := p.sumRange(0, 1350, 1399)
	longSecPrev := p.sumRange(-1, 1350, 1399)

	if longSecCur != 0 || longSecPrev != 0 {
		report.BalanceSheet.Assets.FixedAssets.Financial.OtherLongTermSecurities = ycPos(longSecCur, longSecPrev)
	}
	totalFinAssCur := longSecCur
	totalFinAssPrev := longSecPrev
	report.BalanceSheet.Assets.FixedAssets.Financial.TotalFinancial = ycPos(totalFinAssCur, totalFinAssPrev)

	totalFixAssCur := totalTangCur + totalFinAssCur
	totalFixAssPrev := totalTangPrev + totalFinAssPrev
	report.BalanceSheet.Assets.FixedAssets.TotalFixedAssets = ycPos(totalFixAssCur, totalFixAssPrev)

	// Current assets — inventory
	// 1400–1409: raw materials (råvaror/förnödenheter)
	// 1420–1429: work in progress (varor under tillverkning)
	// 1430–1469: finished goods / trading goods
	rawInvCur := p.sumRange(0, 1400, 1419)
	rawInvPrev := p.sumRange(-1, 1400, 1419)
	wipCur := p.sumRange(0, 1420, 1429)
	wipPrev := p.sumRange(-1, 1420, 1429)
	finGoodsCur := p.sumRange(0, 1430, 1469)
	finGoodsPrev := p.sumRange(-1, 1430, 1469)

	if rawInvCur != 0 || rawInvPrev != 0 {
		report.BalanceSheet.Assets.CurrentAssets.Inventory.RawMaterials = ycPos(rawInvCur, rawInvPrev)
	}
	if wipCur != 0 || wipPrev != 0 {
		report.BalanceSheet.Assets.CurrentAssets.Inventory.WorkInProgress = ycPos(wipCur, wipPrev)
	}
	if finGoodsCur != 0 || finGoodsPrev != 0 {
		report.BalanceSheet.Assets.CurrentAssets.Inventory.FinishedGoods = ycPos(finGoodsCur, finGoodsPrev)
	}

	totalInvCur := rawInvCur + wipCur + finGoodsCur
	totalInvPrev := rawInvPrev + wipPrev + finGoodsPrev
	report.BalanceSheet.Assets.CurrentAssets.Inventory.TotalInventory = ycPos(totalInvCur, totalInvPrev)

	// Current assets — short-term receivables
	// 1510–1519: trade receivables (kundfordringar)
	tradeRecCur := p.sumRange(0, 1510, 1519)
	tradeRecPrev := p.sumRange(-1, 1510, 1519)
	// 1600–1899 excluding prepaid: other receivables
	otherRecCur := p.sumRange(0, 1600, 1799)
	otherRecPrev := p.sumRange(-1, 1600, 1799)
	// 1700–1799 could also be other receivables; 1800–1899 included above.
	// 1900–1969: prepaid expenses and accrued income (förutbetalda kostnader)
	prepaidCur := p.sumRange(0, 1900, 1969)
	prepaidPrev := p.sumRange(-1, 1900, 1969)

	if tradeRecCur != 0 || tradeRecPrev != 0 {
		report.BalanceSheet.Assets.CurrentAssets.ShortTermReceivables.TradeReceivables = ycPos(tradeRecCur, tradeRecPrev)
	}
	if otherRecCur != 0 || otherRecPrev != 0 {
		report.BalanceSheet.Assets.CurrentAssets.ShortTermReceivables.OtherReceivables = ycPos(otherRecCur, otherRecPrev)
	}
	if prepaidCur != 0 || prepaidPrev != 0 {
		report.BalanceSheet.Assets.CurrentAssets.ShortTermReceivables.PrepaidExpenses = ycPos(prepaidCur, prepaidPrev)
	}

	totalSTRecCur := tradeRecCur + otherRecCur + prepaidCur
	totalSTRecPrev := tradeRecPrev + otherRecPrev + prepaidPrev
	report.BalanceSheet.Assets.CurrentAssets.ShortTermReceivables.TotalShortTermReceivables = ycPos(totalSTRecCur, totalSTRecPrev)

	// Current assets — cash and bank (1970–1999)
	cashCur := p.sumRange(0, 1970, 1999)
	cashPrev := p.sumRange(-1, 1970, 1999)

	if cashCur != 0 || cashPrev != 0 {
		report.BalanceSheet.Assets.CurrentAssets.CashAndBank.CashAndBankExcl = ycPos(cashCur, cashPrev)
	}
	report.BalanceSheet.Assets.CurrentAssets.CashAndBank.TotalCashAndBank = ycPos(cashCur, cashPrev)

	totalCurAssCur := totalInvCur + totalSTRecCur + cashCur
	totalCurAssPrev := totalInvPrev + totalSTRecPrev + cashPrev
	report.BalanceSheet.Assets.CurrentAssets.TotalCurrentAssets = ycPos(totalCurAssCur, totalCurAssPrev)

	totalAssCur := totalFixAssCur + totalCurAssCur
	totalAssPrev := totalFixAssPrev + totalCurAssPrev
	report.BalanceSheet.Assets.TotalAssets = ycPos(totalAssCur, totalAssPrev)

	// -----------------------------------------------------------------------
	// Balance Sheet — Equity and Liabilities
	// -----------------------------------------------------------------------

	// Equity — restricted
	// 2081: share capital
	shareCapCur := -p.sum(0, "2081")
	shareCapPrev := -p.sum(-1, "2081")
	// 2082–2085: premium reserve (överkursfond)
	// 2086–2087: reserve fund (reservfond)
	resFundCur := -p.sumRange(0, 2086, 2087)
	resFundPrev := -p.sumRange(-1, 2086, 2087)

	report.BalanceSheet.EquityAndLiabilities.Equity.ShareCapital = ycPos(shareCapCur, shareCapPrev)
	if resFundCur != 0 || resFundPrev != 0 {
		report.BalanceSheet.EquityAndLiabilities.Equity.ReserveFund = ycPos(resFundCur, resFundPrev)
	}
	totalRestEqCur := shareCapCur + resFundCur
	totalRestEqPrev := shareCapPrev + resFundPrev
	report.BalanceSheet.EquityAndLiabilities.Equity.TotalRestrictedEquity = ycPos(totalRestEqCur, totalRestEqPrev)

	// Equity — unrestricted
	// 2091–2097: retained earnings (balanserat resultat)
	retEarnCur := -p.sumRange(0, 2091, 2097)
	retEarnPrev := -p.sumRange(-1, 2091, 2097)
	// 2099: net income current year
	netIncCur := -p.sum(0, "2099")
	netIncPrev := -p.sum(-1, "2099")

	report.BalanceSheet.EquityAndLiabilities.Equity.RetainedEarnings = ycPos(retEarnCur, retEarnPrev)
	report.BalanceSheet.EquityAndLiabilities.Equity.NetIncome = ycPos(netIncCur, netIncPrev)
	totalUnrestEqCur := retEarnCur + netIncCur
	totalUnrestEqPrev := retEarnPrev + netIncPrev
	report.BalanceSheet.EquityAndLiabilities.Equity.TotalUnrestrictedEquity = ycPos(totalUnrestEqCur, totalUnrestEqPrev)

	totalEqCur := totalRestEqCur + totalUnrestEqCur
	totalEqPrev := totalRestEqPrev + totalUnrestEqPrev
	report.BalanceSheet.EquityAndLiabilities.Equity.TotalEquity = ycPos(totalEqCur, totalEqPrev)

	// Untaxed reserves (obeskattade reserver)
	// 2110–2119: periodiseringsfonder
	taxAllocResCur := -p.sumRange(0, 2110, 2119)
	taxAllocResPrev := -p.sumRange(-1, 2110, 2119)
	// 2150–2159: accumulated excess depreciation (ackumulerade överavskrivningar)
	accExDeprCur := -p.sumRange(0, 2150, 2159)
	accExDeprPrev := -p.sumRange(-1, 2150, 2159)

	if taxAllocResCur != 0 || taxAllocResPrev != 0 {
		report.BalanceSheet.EquityAndLiabilities.UntaxedReserves.TaxAllocationReserves = ycPos(taxAllocResCur, taxAllocResPrev)
	}
	if accExDeprCur != 0 || accExDeprPrev != 0 {
		report.BalanceSheet.EquityAndLiabilities.UntaxedReserves.AccumulatedExcessDepreciation = ycPos(accExDeprCur, accExDeprPrev)
	}
	totalUntaxCur := taxAllocResCur + accExDeprCur
	totalUntaxPrev := taxAllocResPrev + accExDeprPrev
	report.BalanceSheet.EquityAndLiabilities.UntaxedReserves.TotalUntaxedReserves = ycPos(totalUntaxCur, totalUntaxPrev)

	// Provisions (avsättningar)
	// 2210–2219: pension provisions
	pensProvCur := -p.sumRange(0, 2210, 2219)
	pensProvPrev := -p.sumRange(-1, 2210, 2219)
	// 2220–2299: other provisions
	otherProvCur := -p.sumRange(0, 2220, 2299)
	otherProvPrev := -p.sumRange(-1, 2220, 2299)

	if pensProvCur != 0 || pensProvPrev != 0 {
		report.BalanceSheet.EquityAndLiabilities.Provisions.PensionProvisions = ycPos(pensProvCur, pensProvPrev)
	}
	if otherProvCur != 0 || otherProvPrev != 0 {
		report.BalanceSheet.EquityAndLiabilities.Provisions.OtherProvisions = ycPos(otherProvCur, otherProvPrev)
	}
	totalProvCur := pensProvCur + otherProvCur
	totalProvPrev := pensProvPrev + otherProvPrev
	report.BalanceSheet.EquityAndLiabilities.Provisions.TotalProvisions = ycPos(totalProvCur, totalProvPrev)

	// Long-term liabilities
	// 2310–2399: bank loans (long-term)
	bankLoansCur := -p.sumRange(0, 2310, 2399)
	bankLoansPrev := -p.sumRange(-1, 2310, 2399)
	// 2400–2499 (excluding 2440): other long-term liabilities
	otherLTCur := -p.sumRange(0, 2400, 2439)
	otherLTPrev := -p.sumRange(-1, 2400, 2439)
	otherLTCur += -p.sumRange(0, 2441, 2499)
	otherLTPrev += -p.sumRange(-1, 2441, 2499)

	if bankLoansCur != 0 || bankLoansPrev != 0 {
		report.BalanceSheet.EquityAndLiabilities.LongTermLiabilities.BankLoans = ycPos(bankLoansCur, bankLoansPrev)
	}
	if otherLTCur != 0 || otherLTPrev != 0 {
		report.BalanceSheet.EquityAndLiabilities.LongTermLiabilities.OtherLongTermLiabilities = ycPos(otherLTCur, otherLTPrev)
	}
	totalLTCur := bankLoansCur + otherLTCur
	totalLTPrev := bankLoansPrev + otherLTPrev
	report.BalanceSheet.EquityAndLiabilities.LongTermLiabilities.TotalLongTermLiabilities = ycPos(totalLTCur, totalLTPrev)

	// Short-term liabilities
	// 2440: trade payables (leverantörsskulder)
	tradePayCur := -p.sum(0, "2440")
	tradePayPrev := -p.sum(-1, "2440")
	// 2510–2519: tax liabilities (skatteskulder)
	taxLiabCur := -p.sumRange(0, 2510, 2519)
	taxLiabPrev := -p.sumRange(-1, 2510, 2519)
	// 2600–2799: other short-term liabilities
	otherSTCur := -p.sumRange(0, 2600, 2799)
	otherSTPrev := -p.sumRange(-1, 2600, 2799)
	// 2900–2999: accrued expenses (upplupna kostnader)
	accExpCur := -p.sumRange(0, 2900, 2999)
	accExpPrev := -p.sumRange(-1, 2900, 2999)

	if tradePayCur != 0 || tradePayPrev != 0 {
		report.BalanceSheet.EquityAndLiabilities.ShortTermLiabilities.TradePayables = ycPos(tradePayCur, tradePayPrev)
	}
	if taxLiabCur != 0 || taxLiabPrev != 0 {
		report.BalanceSheet.EquityAndLiabilities.ShortTermLiabilities.TaxLiabilities = ycPos(taxLiabCur, taxLiabPrev)
	}
	if otherSTCur != 0 || otherSTPrev != 0 {
		report.BalanceSheet.EquityAndLiabilities.ShortTermLiabilities.OtherShortTermLiabilities = ycPos(otherSTCur, otherSTPrev)
	}
	if accExpCur != 0 || accExpPrev != 0 {
		report.BalanceSheet.EquityAndLiabilities.ShortTermLiabilities.AccruedExpenses = ycPos(accExpCur, accExpPrev)
	}
	totalSTCur := tradePayCur + taxLiabCur + otherSTCur + accExpCur
	totalSTPrev := tradePayPrev + taxLiabPrev + otherSTPrev + accExpPrev
	report.BalanceSheet.EquityAndLiabilities.ShortTermLiabilities.TotalShortTermLiabilities = ycPos(totalSTCur, totalSTPrev)

	// Total equity and liabilities
	totalELCur := totalEqCur + totalUntaxCur + totalProvCur + totalLTCur + totalSTCur
	totalELPrev := totalEqPrev + totalUntaxPrev + totalProvPrev + totalLTPrev + totalSTPrev
	report.BalanceSheet.EquityAndLiabilities.TotalEquityAndLiabilities = ycPos(totalELCur, totalELPrev)

	// Balance check warning.
	if totalAssCur != totalELCur {
		warnings = append(warnings, fmt.Sprintf(
			"balance sheet does not balance for current year: assets=%d equity+liabilities=%d diff=%d",
			totalAssCur, totalELCur, totalAssCur-totalELCur))
	}
	if totalAssPrev != totalELPrev {
		warnings = append(warnings, fmt.Sprintf(
			"balance sheet does not balance for previous year: assets=%d equity+liabilities=%d diff=%d",
			totalAssPrev, totalELPrev, totalAssPrev-totalELPrev))
	}

	return report, warnings
}
