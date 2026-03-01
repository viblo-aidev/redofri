package ixbrl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/redofri/redofri/pkg/model"
)

// mapFacts takes a slice of extracted facts and populates a model.AnnualReport.
func mapFacts(facts []fact) (*model.AnnualReport, error) {
	m := &mapper{
		report:  &model.AnnualReport{},
		tuples:  make(map[string][]fact),
		nfByKey: make(map[string][]fact),
		nnByKey: make(map[string][]fact),
	}

	// Index facts by kind and key (name + contextRef).
	for _, f := range facts {
		switch f.Kind {
		case "tuple":
			// Register tuple ID for later grouping.
			if m.tuples[f.TupleID] == nil {
				m.tuples[f.TupleID] = []fact{}
			}
		case "nonFraction":
			key := f.Name + "@" + f.ContextRef
			m.nfByKey[key] = append(m.nfByKey[key], f)
			if f.TupleRef != "" {
				m.tuples[f.TupleRef] = append(m.tuples[f.TupleRef], f)
			}
		case "nonNumeric":
			key := f.Name + "@" + f.ContextRef
			m.nnByKey[key] = append(m.nnByKey[key], f)
			if f.TupleRef != "" {
				m.tuples[f.TupleRef] = append(m.tuples[f.TupleRef], f)
			}
		}
	}

	// Map all sections.
	m.mapMeta()
	m.mapCertification()
	m.mapManagementReport()
	m.mapIncomeStatement()
	m.mapBalanceSheet()
	m.mapNotes(facts)
	m.mapSignatures(facts)

	if m.err != nil {
		return nil, m.err
	}

	return m.report, nil
}

// mapper holds state during fact-to-model mapping.
type mapper struct {
	report  *model.AnnualReport
	tuples  map[string][]fact // tupleID -> member facts
	nfByKey map[string][]fact // "name@context" -> nonFraction facts
	nnByKey map[string][]fact // "name@context" -> nonNumeric facts
	err     error             // sticky error
}

// ---------- helpers ----------

const (
	nsGen = "se-gen-base:"
	nsCd  = "se-cd-base:"
	nsBol = "se-bol-base:"
)

// nn returns the text value of the first nonNumeric fact matching name@context.
func (m *mapper) nn(name, ctx string) string {
	key := name + "@" + ctx
	if fs, ok := m.nnByKey[key]; ok && len(fs) > 0 {
		return fs[0].Value
	}
	return ""
}

// nf returns the int64 value of the first nonFraction fact matching name@context.
// Returns nil if not found.
func (m *mapper) nf(name, ctx string) *int64 {
	key := name + "@" + ctx
	fs, ok := m.nfByKey[key]
	if !ok || len(fs) == 0 {
		return nil
	}
	v, err := parseNumber(fs[0])
	if err != nil {
		if m.err == nil {
			m.err = fmt.Errorf("parsing %s@%s: %w", name, ctx, err)
		}
		return nil
	}
	return &v
}

// nfNeg returns the negated int64 value for a nonFraction fact stored with sign="-".
// These concepts are stored as positive in the model but emitted with sign="-" in iXBRL.
func (m *mapper) nfNeg(name, ctx string) *int64 {
	v := m.nf(name, ctx)
	if v == nil {
		return nil
	}
	neg := -*v
	return &neg
}

// ycNeg builds a YearComparison for a concept that uses sign="-" in iXBRL.
func (m *mapper) ycNeg(name string) model.YearComparison {
	return model.YearComparison{
		Current:  m.nfNeg(name, "period0"),
		Previous: m.nfNeg(name, "period1"),
	}
}

// yc builds a YearComparison from a concept at two contexts.
func (m *mapper) yc(name, ctxCurrent, ctxPrevious string) model.YearComparison {
	return model.YearComparison{
		Current:  m.nf(name, ctxCurrent),
		Previous: m.nf(name, ctxPrevious),
	}
}

// ycPeriod builds a YearComparison for period-based contexts (period0/period1).
func (m *mapper) ycPeriod(name string) model.YearComparison {
	return m.yc(name, "period0", "period1")
}

// ycBalans builds a YearComparison for instant-based contexts (balans0/balans1).
func (m *mapper) ycBalans(name string) model.YearComparison {
	return m.yc(name, "balans0", "balans1")
}

// ---------- section mappers ----------

func (m *mapper) mapMeta() {
	r := m.report

	r.Company.Name = m.nn(nsCd+"ForetagetsNamn", "period0")
	r.Company.OrgNr = m.nn(nsCd+"Organisationsnummer", "period0")

	r.FiscalYear.StartDate = m.nn(nsCd+"RakenskapsarForstaDag", "period0")
	r.FiscalYear.EndDate = m.nn(nsCd+"RakenskapsarSistaDag", "period0")

	r.Meta.Language = m.nn(nsCd+"Sprak", "period0")
	r.Meta.Country = m.nn(nsCd+"Land", "period0")
	r.Meta.Currency = m.nn(nsCd+"Redovisningsvaluta", "period0")
	r.Meta.AmountFormat = m.nn(nsCd+"Beloppsformat", "period0")
}

func (m *mapper) mapCertification() {
	c := &m.report.Certification

	c.ConfirmationText = m.nn(nsBol+"FaststallelseResultatBalansrakning", "balans0")
	c.MeetingDate = m.nn(nsBol+"Arsstamma", "balans0")
	c.DispositionDecision = m.nn(nsBol+"ArsstammaResultatDispositionGodkannaStyrelsensForslag", "balans0")
	c.OriginalContentCertification = m.nn(nsBol+"IntygandeOriginalInnehall", "balans0")
	c.ElectronicSignatureLabel = m.nn(nsBol+"UnderskriftFaststallelseintygElektroniskt", "balans0")
	c.SigningDate = m.nn(nsBol+"UnderskriftFastallelseintygDatum", "balans0")

	c.Signatory.FirstName = m.nn(nsBol+"UnderskriftFaststallelseintygForetradareTilltalsnamn", "period0")
	c.Signatory.LastName = m.nn(nsBol+"UnderskriftFaststallelseintygForetradareEfternamn", "period0")
	c.Signatory.Role = m.nn(nsBol+"UnderskriftFaststallelseintygForetradareForetradarroll", "period0")
}

func (m *mapper) mapManagementReport() {
	mr := &m.report.ManagementReport

	mr.IntroText = m.nn(nsGen+"LopandeBokforingenAvslutasMening", "period0")
	mr.BusinessDescription = m.nn(nsGen+"AllmantVerksamheten", "period0")
	mr.SignificantEvents = m.nn(nsGen+"VasentligaHandelserRakenskapsaret", "period0")
	mr.MultiYearOverview.Comment = m.nn(nsGen+"KommentarFlerarsoversikt", "period0")
	mr.BoardDividendStatement = m.nn(nsGen+"StyrelsensYttrandeVinstutdelning", "balans0")

	// Multi-year overview: look for Nettoomsattning at period0..period3 with scale=3 (tkr).
	mr.MultiYearOverview.Years = m.mapMultiYearOverview()

	// Equity changes.
	m.mapEquityChanges()

	// Profit disposition.
	m.mapProfitDisposition()
}

func (m *mapper) mapMultiYearOverview() []model.MultiYearOverviewYear {
	var years []model.MultiYearOverviewYear

	// We look for facts at period0..period3 and balans0..balans3.
	// Multi-year overview uses scale=3 for tkr, but we want the raw XBRL value
	// (i.e. the actual amount). The nf() method applies scale automatically.
	for i := 0; i <= 3; i++ {
		pCtx := fmt.Sprintf("period%d", i)
		bCtx := fmt.Sprintf("balans%d", i)

		// Check if there's a Nettoomsattning fact at this period.
		// Multi-year overview uses scale=3, IS uses scale=0.
		// We need the tkr one. Find the fact with scale=3.
		netSales := m.nfWithScale(nsGen+"Nettoomsattning", pCtx, 3)
		resultFin := m.nfWithScale(nsGen+"ResultatEfterFinansiellaPoster", pCtx, 3)

		// Soliditet (percentage) - at balans context.
		solidity := m.solidityStr(nsGen+"Soliditet", bCtx)

		// Only add if we have at least one value.
		if netSales != nil || resultFin != nil || solidity != nil {
			y := model.MultiYearOverviewYear{
				NetSales:                  netSales,
				ResultAfterFinancialItems: resultFin,
				Solidity:                  solidity,
			}
			years = append(years, y)
		}
	}

	return years
}

// nfWithScale finds a nonFraction fact with a specific scale value.
func (m *mapper) nfWithScale(name, ctx string, wantScale int) *int64 {
	key := name + "@" + ctx
	fs, ok := m.nfByKey[key]
	if !ok {
		return nil
	}
	for _, f := range fs {
		if f.Scale == wantScale {
			v, err := parseNumber(f)
			if err != nil {
				if m.err == nil {
					m.err = fmt.Errorf("parsing %s@%s: %w", name, ctx, err)
				}
				return nil
			}
			return &v
		}
	}
	return nil
}

// solidityStr extracts the solidity percentage as a display string (e.g. "33,7").
func (m *mapper) solidityStr(name, ctx string) *string {
	key := name + "@" + ctx
	fs, ok := m.nfByKey[key]
	if !ok || len(fs) == 0 {
		return nil
	}
	// The display value is the raw text, e.g. "33,7".
	v := fs[0].Value
	if v == "" {
		return nil
	}
	return &v
}

func (m *mapper) mapEquityChanges() {
	ec := &m.report.ManagementReport.EquityChanges

	// Opening (balans1)
	ec.OpeningShareCapital = m.nf(nsGen+"Aktiekapital", "balans1")
	ec.OpeningReserveFund = m.nf(nsGen+"Reservfond", "balans1")
	ec.OpeningRetainedEarnings = m.nf(nsGen+"BalanseratResultat", "balans1")
	ec.OpeningNetIncome = m.nf(nsGen+"AretsResultatEgetKapital", "balans1")
	ec.OpeningTotal = m.nf(nsGen+"ForandringEgetKapitalTotalt", "balans1")

	// Dividend
	ec.DividendNetIncome = m.nf(nsGen+"ForandringEgetKapitalAretsResultatUtdelning", "period0")
	ec.DividendTotal = m.nf(nsGen+"ForandringEgetKapitalTotaltUtdelning", "period0")

	// Year result
	ec.YearResultNetIncome = m.nf(nsGen+"ForandringEgetKapitalAretsResultatAretsResultat", "period0")
	ec.YearResultTotal = m.nf(nsGen+"ForandringEgetKapitalTotaltAretsResultat", "period0")

	// Closing (balans0) — these same concepts appear in BS too; we take first occurrence.
	ec.ClosingShareCapital = m.nf(nsGen+"Aktiekapital", "balans0")
	ec.ClosingReserveFund = m.nf(nsGen+"Reservfond", "balans0")
	ec.ClosingRetainedEarnings = m.nf(nsGen+"BalanseratResultat", "balans0")
	ec.ClosingNetIncome = m.nf(nsGen+"AretsResultatEgetKapital", "balans0")
	ec.ClosingTotal = m.nf(nsGen+"ForandringEgetKapitalTotalt", "balans0")
}

func (m *mapper) mapProfitDisposition() {
	pd := &m.report.ManagementReport.ProfitDisposition

	pd.RetainedEarnings = m.nf(nsGen+"BalanseratResultat", "balans0")
	pd.NetIncome = m.nf(nsGen+"AretsResultatEgetKapital", "balans0")
	pd.TotalAvailable = m.nf(nsGen+"MedelDisponera", "balans0")
	pd.Dividend = m.nf(nsGen+"ForslagDispositionUtdelning", "balans0")
	pd.CarriedForward = m.nf(nsGen+"ForslagDispositionBalanserasINyRakning", "balans0")
	pd.TotalDisposition = m.nf(nsGen+"ForslagDisposition", "balans0")
}

func (m *mapper) mapIncomeStatement() {
	is := &m.report.IncomeStatement

	// Revenue
	is.Revenue.NetSales = m.ycPeriod(nsGen + "Nettoomsattning")
	is.Revenue.InventoryChange = m.ycPeriod(nsGen + "ForandringLagerProdukterIArbeteFardigaVarorPagaendeArbetenAnnansRakning")
	is.Revenue.OtherOperatingIncome = m.ycPeriod(nsGen + "OvrigaRorelseintakter")
	is.Revenue.TotalRevenue = m.ycPeriod(nsGen + "RorelseintakterLagerforandringarMm")

	// Expenses
	is.Expenses.RawMaterials = m.ycPeriod(nsGen + "RavarorFornodenheterKostnader")
	is.Expenses.TradingGoods = m.ycPeriod(nsGen + "HandelsvarorKostnader")
	is.Expenses.OtherExternalExpenses = m.ycPeriod(nsGen + "OvrigaExternaKostnader")
	is.Expenses.PersonnelExpenses = m.ycPeriod(nsGen + "Personalkostnader")
	is.Expenses.DepreciationAmortization = m.ycPeriod(nsGen + "AvskrivningarNedskrivningarMateriellaImmateriellaAnlaggningstillgangar")
	is.Expenses.OtherOperatingExpenses = m.ycPeriod(nsGen + "OvrigaRorelsekostnader")
	is.Expenses.TotalExpenses = m.ycPeriod(nsGen + "Rorelsekostnader")

	// Operating result
	is.OperatingResult = m.ycPeriod(nsGen + "Rorelseresultat")

	// Financial items
	is.FinancialItems.ResultOtherFinancialAssets = m.ycPeriod(nsGen + "ResultatOvrigaFinansiellaAnlaggningstillgangar")
	is.FinancialItems.OtherInterestIncome = m.ycPeriod(nsGen + "OvrigaRanteintakterLiknandeResultatposter")
	is.FinancialItems.InterestExpenses = m.ycPeriod(nsGen + "RantekostnaderLiknandeResultatposter")
	is.FinancialItems.TotalFinancialItems = m.ycPeriod(nsGen + "FinansiellaPoster")

	// Result after financial items
	is.ResultAfterFinancialItems = m.ycPeriod(nsGen + "ResultatEfterFinansiellaPoster")

	// Appropriations
	is.Appropriations.TaxAllocationReserve = m.ycNeg(nsGen + "ForandringPeriodiseringsfond")
	is.Appropriations.ExcessDepreciation = m.ycNeg(nsGen + "ForandringOveravskrivningar")
	is.Appropriations.TotalAppropriations = m.ycNeg(nsGen + "Bokslutsdispositioner")

	// Result before tax
	is.ResultBeforeTax = m.ycPeriod(nsGen + "ResultatForeSkatt")

	// Tax
	is.Tax.IncomeTax = m.ycPeriod(nsGen + "SkattAretsResultat")

	// Net result
	is.NetResult = m.ycPeriod(nsGen + "AretsResultat")
}

func (m *mapper) mapBalanceSheet() {
	bs := &m.report.BalanceSheet

	// Assets
	a := &bs.Assets

	// Tangible fixed assets
	t := &a.FixedAssets.Tangible
	t.BuildingsAndLand = m.ycBalans(nsGen + "ByggnaderMark")
	t.MachineryAndEquipment = m.ycBalans(nsGen + "MaskinerAndraTekniskaAnlaggningar")
	t.FixturesAndFittings = m.ycBalans(nsGen + "InventarierVerktygInstallationer")
	t.TotalTangible = m.ycBalans(nsGen + "MateriellaAnlaggningstillgangar")

	// Financial fixed assets
	ff := &a.FixedAssets.Financial
	ff.OtherLongTermSecurities = m.ycBalans(nsGen + "AndraLangfristigaVardepappersinnehav")
	ff.TotalFinancial = m.ycBalans(nsGen + "FinansiellaAnlaggningstillgangar")

	a.FixedAssets.TotalFixedAssets = m.ycBalans(nsGen + "Anlaggningstillgangar")

	// Current assets - Inventory
	inv := &a.CurrentAssets.Inventory
	inv.RawMaterials = m.ycBalans(nsGen + "LagerRavarorFornodenheter")
	inv.WorkInProgress = m.ycBalans(nsGen + "LagerVarorUnderTillverkning")
	inv.FinishedGoods = m.ycBalans(nsGen + "LagerFardigaVarorHandelsvaror")
	inv.TotalInventory = m.ycBalans(nsGen + "VarulagerMm")

	// Current assets - Short term receivables
	str := &a.CurrentAssets.ShortTermReceivables
	str.TradeReceivables = m.ycBalans(nsGen + "Kundfordringar")
	str.OtherReceivables = m.ycBalans(nsGen + "OvrigaFordringarKortfristiga")
	str.PrepaidExpenses = m.ycBalans(nsGen + "ForutbetaldaKostnaderUpplupnaIntakter")
	str.TotalShortTermReceivables = m.ycBalans(nsGen + "KortfristigaFordringar")

	// Current assets - Cash and bank
	cb := &a.CurrentAssets.CashAndBank
	cb.CashAndBankExcl = m.ycBalans(nsGen + "KassaBankExklRedovisningsmedel")
	cb.TotalCashAndBank = m.ycBalans(nsGen + "KassaBank")

	a.CurrentAssets.TotalCurrentAssets = m.ycBalans(nsGen + "Omsattningstillgangar")
	a.TotalAssets = m.ycBalans(nsGen + "Tillgangar")

	// Equity and Liabilities
	el := &bs.EquityAndLiabilities

	// Equity
	eq := &el.Equity
	eq.ShareCapital = m.ycBalans(nsGen + "Aktiekapital")
	eq.ReserveFund = m.ycBalans(nsGen + "Reservfond")
	eq.TotalRestrictedEquity = m.ycBalans(nsGen + "BundetEgetKapital")
	eq.RetainedEarnings = m.ycBalans(nsGen + "BalanseratResultat")
	eq.NetIncome = m.ycBalans(nsGen + "AretsResultatEgetKapital")
	eq.TotalUnrestrictedEquity = m.ycBalans(nsGen + "FrittEgetKapital")
	eq.TotalEquity = m.ycBalans(nsGen + "EgetKapital")

	// Untaxed reserves
	ur := &el.UntaxedReserves
	ur.TaxAllocationReserves = m.ycBalans(nsGen + "Periodiseringsfonder")
	ur.AccumulatedExcessDepreciation = m.ycBalans(nsGen + "AckumuleradeOveravskrivningar")
	ur.TotalUntaxedReserves = m.ycBalans(nsGen + "ObeskattadeReserver")

	// Provisions
	prov := &el.Provisions
	prov.PensionProvisions = m.ycBalans(nsGen + "AvsattningarPensionerLiknandeForpliktelserEnligtLag")
	prov.OtherProvisions = m.ycBalans(nsGen + "OvrigaAvsattningar")
	prov.TotalProvisions = m.ycBalans(nsGen + "Avsattningar")

	// Long-term liabilities
	ltl := &el.LongTermLiabilities
	ltl.BankLoans = m.ycBalans(nsGen + "OvrigaLangfristigaSkulderKreditinstitut")
	ltl.OtherLongTermLiabilities = m.ycBalans(nsGen + "OvrigaLangfristigaSkulder")
	ltl.TotalLongTermLiabilities = m.ycBalans(nsGen + "LangfristigaSkulder")

	// Short-term liabilities
	stl := &el.ShortTermLiabilities
	stl.TradePayables = m.ycBalans(nsGen + "Leverantorsskulder")
	stl.TaxLiabilities = m.ycBalans(nsGen + "Skatteskulder")
	stl.OtherShortTermLiabilities = m.ycBalans(nsGen + "OvrigaKortfristigaSkulder")
	stl.AccruedExpenses = m.ycBalans(nsGen + "UpplupnaKostnaderForutbetaldaIntakter")
	stl.TotalShortTermLiabilities = m.ycBalans(nsGen + "KortfristigaSkulder")

	el.TotalEquityAndLiabilities = m.ycBalans(nsGen + "EgetKapitalSkulder")
}

func (m *mapper) mapNotes(facts []fact) {
	n := &m.report.Notes

	// Note 1: Accounting policies
	m.mapAccountingPolicies(n)

	// Note 2: Employees
	m.mapEmployeesNote(n)

	// Notes 3-6: Fixed asset roll-forwards
	m.mapFixedAssetNotes(n)

	// Note 7: Long-term liabilities
	m.mapLongTermLiabilitiesNote(n)

	// Note 8: Pledges
	m.mapPledgesNote(n)

	// Note 9: Contingent liabilities
	m.mapContingentLiabilitiesNote(n)

	// Note 10: Multi-post note
	m.mapMultiPostNote(n, facts)
}

func (m *mapper) mapAccountingPolicies(n *model.Notes) {
	desc := m.nn(nsGen+"Redovisningsprinciper", "period0")
	if desc == "" {
		return
	}

	ap := &n.AccountingPolicies
	ap.NoteNumber = 1
	ap.Description = desc

	// Depreciation policies — emitted as ix:nonNumeric, not nonFraction.
	type depCfg struct {
		concept  string
		category string
	}
	deps := []depCfg{
		{nsGen + "AvskrivningarMateriellaAnlaggningstillgangarByggnaderAr", "Byggnader"},
		{nsGen + "AvskrivningarMateriellaAnlaggningstillgangarMaskinerAndraTekniskaAnlaggningarAr", "Maskiner och andra tekniska anläggningar"},
		{nsGen + "AvskrivningarMateriellaAnlaggningstillgangarInventarierVerktygInstallationerAr", "Inventarier, verktyg och installationer"},
	}

	for _, d := range deps {
		v := m.nn(d.concept, "period0")
		if v != "" {
			years, err := strconv.Atoi(strings.TrimSpace(v))
			if err == nil {
				ap.Depreciations = append(ap.Depreciations, model.DepreciationPolicy{
					Category: d.category,
					Concept:  d.concept,
					Years:    years,
				})
			}
		}
	}

	ap.DepreciationComment = m.nn(nsGen+"AvskrivningarMateriellaAnlaggningstillgangarKommentar", "period0")
	ap.ManufacturedGoodsPolicy = m.nn(nsGen+"RedovisningsprinciperAnskaffningsvardeEgentillverkadevaror", "period0")
}

func (m *mapper) mapEmployeesNote(n *model.Notes) {
	avg := m.ycPeriod(nsGen + "MedelantaletAnstallda")
	if avg.Current == nil && avg.Previous == nil {
		return
	}

	n.Employees = &model.EmployeesNote{
		NoteNumber:       2,
		AverageEmployees: avg,
	}
}

// knownAssetPrefixes lists the XBRL concept prefixes for fixed asset notes,
// in the order they typically appear.
var knownAssetPrefixes = []struct {
	prefix string
	title  string
}{
	{"ByggnaderMark", "Byggnader och mark"},
	{"MaskinerAndraTekniskaAnlaggningar", "Maskiner och andra tekniska anläggningar"},
	{"InventarierVerktygInstallationer", "Inventarier, verktyg och installationer"},
	{"AndraLangfristigaVardepappersinnehav", "Andra långfristiga värdepappersinnehav"},
}

func (m *mapper) mapFixedAssetNotes(n *model.Notes) {
	noteNum := 3

	for _, ap := range knownAssetPrefixes {
		prefix := nsGen + ap.prefix

		// Check if we have acquisition values for this asset type.
		openAcq := m.yc(prefix+"Anskaffningsvarden", "balans1", "balans2")
		if openAcq.Current == nil && openAcq.Previous == nil {
			continue
		}

		fan := model.FixedAssetNote{
			NoteNumber:    noteNum,
			Title:         ap.title,
			ConceptPrefix: ap.prefix,

			OpeningAcquisitionValues: openAcq,
			Purchases:                m.ycPeriod(prefix + "ForandringAnskaffningsvardenInkop"),
			Sales:                    m.ycPeriod(prefix + "ForandringAnskaffningsvardenForsaljningar"),
			ClosingAcquisitionValues: m.ycBalans(prefix + "Anskaffningsvarden"),

			OpeningDepreciation: m.yc(prefix+"Avskrivningar", "balans1", "balans2"),
			YearDepreciation:    m.ycPeriod(prefix + "ForandringAvskrivningarAretsAvskrivningar"),
			ClosingDepreciation: m.ycBalans(prefix + "Avskrivningar"),

			CarryingValue: m.ycBalans(prefix),
		}

		n.FixedAssetNotes = append(n.FixedAssetNotes, fan)
		noteNum++
	}
}

func (m *mapper) mapLongTermLiabilitiesNote(n *model.Notes) {
	dueAfter5 := m.ycBalans(nsGen + "LangfristigaSkulderForfallerSenare5Ar")
	if dueAfter5.Current == nil && dueAfter5.Previous == nil {
		return
	}

	n.LongTermLiabilitiesNote = &model.LongTermLiabilitiesNoteData{
		NoteNumber:        7,
		DueAfterFiveYears: dueAfter5,
	}
}

func (m *mapper) mapPledgesNote(n *model.Notes) {
	total := m.ycBalans(nsGen + "StalldaSakerheter")
	if total.Current == nil && total.Previous == nil {
		return
	}

	n.Pledges = &model.PledgesNote{
		NoteNumber:          8,
		CorporateMortgages:  m.ycBalans(nsGen + "StalldaSakerheterForetagsinteckningar"),
		RealEstateMortgages: m.ycBalans(nsGen + "StalldaSakerheterFastighetsinteckningar"),
		TotalPledges:        total,
	}
}

func (m *mapper) mapContingentLiabilitiesNote(n *model.Notes) {
	total := m.ycBalans(nsGen + "EventualForpliktelser")
	if total.Current == nil && total.Previous == nil {
		return
	}

	n.ContingentLiabilities = &model.ContingentLiabilitiesNote{
		NoteNumber:      9,
		TotalContingent: total,
	}
}

func (m *mapper) mapMultiPostNote(n *model.Notes, facts []fact) {
	desc := m.nn(nsGen+"NotTillgangarAvsattningarSkulderAvserFleraPoster", "balans0")
	if desc == "" {
		return
	}

	mp := &model.MultiPostNote{
		NoteNumber:  10,
		Description: desc,
	}

	// Find all multi-post tuples and their member facts.
	tupleName := nsGen + "TillgangarAvsattningarSkulderTuple"
	for _, f := range facts {
		if f.Kind == "tuple" && f.Name == tupleName {
			members := m.tuples[f.TupleID]
			entry := model.MultiPostEntry{}
			for _, mf := range members {
				switch mf.Name {
				case nsGen + "TillgangarAvsattningarSkulderPost":
					entry.PostName = mf.Value
				case nsGen + "TillgangarAvsattningarSkulderBelopp":
					v, err := parseNumber(mf)
					if err == nil {
						entry.Amount = &v
					}
				}
			}
			if entry.PostName != "" {
				mp.Entries = append(mp.Entries, entry)
			}
		}
	}

	if len(mp.Entries) > 0 {
		n.MultiPostNote = mp
	}
}

func (m *mapper) mapSignatures(facts []fact) {
	sig := &m.report.Signatures

	sig.City = m.nn(nsGen+"UndertecknandeArsredovisningOrt", "period0")
	sig.Date = m.nn(nsGen+"UndertecknandeArsredovisningDatum", "period0")

	// Find signatory tuples.
	tupleName := nsGen + "UnderskriftArsredovisningForetradareTuple"
	for _, f := range facts {
		if f.Kind == "tuple" && f.Name == tupleName {
			members := m.tuples[f.TupleID]
			s := model.Signatory{}
			for _, mf := range members {
				localName := strings.TrimPrefix(mf.Name, nsGen)
				switch localName {
				case "UnderskriftArsredovisningForetradareTilltalsnamn":
					s.FirstName = mf.Value
				case "UnderskriftArsredovisningForetradareEfternamn":
					s.LastName = mf.Value
				case "UnderskriftArsredovisningForetradareForetradarroll":
					s.Role = mf.Value
				}
			}
			if s.FirstName != "" || s.LastName != "" {
				sig.Signatories = append(sig.Signatories, s)
			}
		}
	}
}
