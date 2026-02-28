package ixbrl

import (
	"strings"

	"github.com/redofri/redofri/pkg/model"
)

// writeBalanceSheetAssets writes the tillgångar page (page 5).
func (g *generator) writeBalanceSheetAssets(r *model.AnnualReport) {
	bs := &r.BalanceSheet
	totalPages := g.computeTotalPages(r)

	_, prevEnd := prevYearDates(r.FiscalYear.StartDate, r.FiscalYear.EndDate)

	g.line(`<div class="ar-page wide" id="ar3-page-5">`)
	g.in()
	g.pageHeader(r.Company.Name, r.Company.OrgNr, 5, totalPages)

	g.line(`<table class="ar-balance-sheet ar-financial col-4">`)
	g.in()

	// Colgroup
	g.line(`<colgroup>`)
	g.in()
	g.line(`<col />`)
	g.line(`<col class="note" />`)
	g.line(`<col class="kr" span="2" />`)
	g.out()
	g.line(`</colgroup>`)

	// Header
	g.line(`<thead>`)
	g.in()
	g.line(`<tr>`)
	g.in()
	g.line(`<th scope="col">Balansräkning</th>`)
	g.line(`<th scope="col">Not</th>`)
	g.linef(`<th scope="col">%s</th>`, r.FiscalYear.EndDate)
	g.linef(`<th scope="col">%s</th>`, prevEnd)
	g.out()
	g.line(`</tr>`)
	g.line(`<tr>`)
	g.in()
	g.line(`<th colspan="4" scope="colgroup">Tillgångar</th>`)
	g.out()
	g.line(`</tr>`)
	g.out()
	g.line(`</thead>`)

	// Fixed assets
	g.writeBSFixedAssets(bs)

	// Current assets
	g.writeBSCurrentAssets(bs)

	g.out()
	g.line(`</table>`)

	g.out()
	g.line(`</div>`)
}

// writeBSFixedAssets writes the anläggningstillgångar tbody.
func (g *generator) writeBSFixedAssets(bs *model.BalanceSheet) {
	fa := &bs.Assets.FixedAssets

	g.line(`<tbody>`)
	g.in()

	// Section header
	g.line(`<tr>`)
	g.in()
	g.line(`<th colspan="4" scope="rowgroup">Anläggningstillgångar</th>`)
	g.out()
	g.line(`</tr>`)

	// Materiella anläggningstillgångar
	g.line(`<tr>`)
	g.in()
	g.line(`<th colspan="4" scope="rowgroup" class="sub">Materiella anläggningstillgångar</th>`)
	g.out()
	g.line(`</tr>`)

	tang := &fa.Tangible

	g.writeBalanceRow("Byggnader och mark", tang.BuildingsAndLandNote, nil,
		"se-gen-base:ByggnaderMark",
		ycv(tang.BuildingsAndLand), false, false, false)

	g.writeBalanceRow("Maskiner och andra tekniska anläggningar", tang.MachineryAndEquipmentNote, nil,
		"se-gen-base:MaskinerAndraTekniskaAnlaggningar",
		ycv(tang.MachineryAndEquipment), false, false, false)

	// Last in tangible group — sum wrap
	g.writeBalanceRow("Inventarier, verktyg och installationer", tang.FixturesAndFittingsNote, nil,
		"se-gen-base:InventarierVerktygInstallationer",
		ycv(tang.FixturesAndFittings), false, false, true)

	// Summa materiella
	g.writeBalanceRow("Summa materiella anläggningstillgångar", 0, nil,
		"se-gen-base:MateriellaAnlaggningstillgangar",
		ycv(tang.TotalTangible), true, false, false)

	// Finansiella anläggningstillgångar
	g.line(`<tr>`)
	g.in()
	g.line(`<th colspan="4" scope="rowgroup" class="sub sep">Finansiella anläggningstillgångar</th>`)
	g.out()
	g.line(`</tr>`)

	fin := &fa.Financial

	// Last (and only) in financial group — sum wrap
	g.writeBalanceRow("Andra långfristiga värdepappersinnehav", fin.OtherLongTermSecuritiesNote, nil,
		"se-gen-base:AndraLangfristigaVardepappersinnehav",
		ycv(fin.OtherLongTermSecurities), false, false, true)

	// Summa finansiella
	g.writeBalanceRow("Summa finansiella anläggningstillgångar", 0, nil,
		"se-gen-base:FinansiellaAnlaggningstillgangar",
		ycv(fin.TotalFinancial), true, false, false)

	// Summa anläggningstillgångar
	g.writeBalanceRow("Summa anläggningstillgångar", 0, nil,
		"se-gen-base:Anlaggningstillgangar",
		ycv(fa.TotalFixedAssets), true, false, false)

	g.out()
	g.line(`</tbody>`)
}

// writeBSCurrentAssets writes the omsättningstillgångar tbody.
func (g *generator) writeBSCurrentAssets(bs *model.BalanceSheet) {
	ca := &bs.Assets.CurrentAssets

	g.line(`<tbody>`)
	g.in()

	g.line(`<tr>`)
	g.in()
	g.line(`<th colspan="4" scope="rowgroup">Omsättningstillgångar</th>`)
	g.out()
	g.line(`</tr>`)

	// Varulager m.m.
	g.line(`<tr>`)
	g.in()
	g.line(`<th colspan="4" scope="rowgroup" class="sub sep">Varulager <abbr>m.m.</abbr>`)
	g.line(`</th>`)
	g.out()
	g.line(`</tr>`)

	inv := &ca.Inventory

	g.writeBalanceRow("Råvaror och förnödenheter", 0, nil,
		"se-gen-base:LagerRavarorFornodenheter",
		ycv(inv.RawMaterials), false, false, false)

	g.writeBalanceRow("Varor under tillverkning", 0, nil,
		"se-gen-base:LagerVarorUnderTillverkning",
		ycv(inv.WorkInProgress), false, false, false)

	g.writeBalanceRow("Färdiga varor och handelsvaror", 0, nil,
		"se-gen-base:LagerFardigaVarorHandelsvaror",
		ycv(inv.FinishedGoods), false, false, true)

	g.writeBalanceRow("Summa varulager", 0, nil,
		"se-gen-base:VarulagerMm",
		ycv(inv.TotalInventory), true, false, false)

	// Kortfristiga fordringar
	g.line(`<tr>`)
	g.in()
	g.line(`<th colspan="4" scope="rowgroup" class="sub sep">Kortfristiga fordringar</th>`)
	g.out()
	g.line(`</tr>`)

	str := &ca.ShortTermReceivables

	g.writeBalanceRow("Kundfordringar", 0, nil,
		"se-gen-base:Kundfordringar",
		ycv(str.TradeReceivables), false, false, false)

	g.writeBalanceRow("Övriga fordringar", 0, nil,
		"se-gen-base:OvrigaFordringarKortfristiga",
		ycv(str.OtherReceivables), false, false, false)

	g.writeBalanceRow("Förutbetalda kostnader och upplupna intäkter", 0, nil,
		"se-gen-base:ForutbetaldaKostnaderUpplupnaIntakter",
		ycv(str.PrepaidExpenses), false, false, true)

	g.writeBalanceRow("Summa kortfristiga fordringar", 0, nil,
		"se-gen-base:KortfristigaFordringar",
		ycv(str.TotalShortTermReceivables), true, false, false)

	// Kassa och bank
	g.line(`<tr class="sep">`)
	g.in()
	g.line(`<th colspan="4" scope="rowgroup" class="sub">Kassa och bank</th>`)
	g.out()
	g.line(`</tr>`)

	cb := &ca.CashAndBank

	g.writeBalanceRow("Kassa och bank", 0, nil,
		"se-gen-base:KassaBankExklRedovisningsmedel",
		ycv(cb.CashAndBankExcl), false, false, true)

	g.writeBalanceRow("Summa kassa och bank", 0, nil,
		"se-gen-base:KassaBank",
		ycv(cb.TotalCashAndBank), true, false, false)

	// Summa omsättningstillgångar
	g.writeBalanceRow("Summa omsättningstillgångar", 0, nil,
		"se-gen-base:Omsattningstillgangar",
		ycv(ca.TotalCurrentAssets), true, false, false)

	// Summa tillgångar (total)
	g.writeBalanceRow("Summa tillgångar", 0, nil,
		"se-gen-base:Tillgangar",
		ycv(bs.Assets.TotalAssets), false, true, false)

	g.out()
	g.line(`</tbody>`)
}

// writeBalanceSheetEquityLiabilities writes the eget kapital och skulder page (page 6).
func (g *generator) writeBalanceSheetEquityLiabilities(r *model.AnnualReport) {
	el := &r.BalanceSheet.EquityAndLiabilities
	totalPages := g.computeTotalPages(r)

	_, prevEnd := prevYearDates(r.FiscalYear.StartDate, r.FiscalYear.EndDate)

	g.line(`<div class="ar-page wide" id="ar3-page-6">`)
	g.in()
	g.pageHeader(r.Company.Name, r.Company.OrgNr, 6, totalPages)

	g.line(`<table class="ar-balance-sheet ar-financial col-4">`)
	g.in()

	// Colgroup
	g.line(`<colgroup>`)
	g.in()
	g.line(`<col />`)
	g.line(`<col class="note" />`)
	g.line(`<col class="kr" span="2" />`)
	g.out()
	g.line(`</colgroup>`)

	// Header
	g.line(`<thead>`)
	g.in()
	g.line(`<tr>`)
	g.in()
	g.line(`<th scope="col">Balansräkning</th>`)
	g.line(`<th scope="col">Not</th>`)
	g.linef(`<th scope="col">%s</th>`, r.FiscalYear.EndDate)
	g.linef(`<th scope="col">%s</th>`, prevEnd)
	g.out()
	g.line(`</tr>`)
	g.line(`<tr>`)
	g.in()
	g.line(`<th colspan="4" scope="colgroup">Eget kapital och skulder</th>`)
	g.out()
	g.line(`</tr>`)
	g.out()
	g.line(`</thead>`)

	// Eget kapital
	g.writeBSEquity(el)

	// Obeskattade reserver
	g.writeBSUntaxedReserves(el)

	// Avsättningar
	g.writeBSProvisions(el)

	// Långfristiga skulder
	g.writeBSLongTermLiabilities(el)

	// Kortfristiga skulder
	g.writeBSShortTermLiabilities(el)

	g.out()
	g.line(`</table>`)

	g.out()
	g.line(`</div>`)
}

// writeBSEquity writes the eget kapital tbody.
func (g *generator) writeBSEquity(el *model.EquityAndLiabilities) {
	eq := &el.Equity

	g.line(`<tbody>`)
	g.in()

	// Eget kapital header
	g.line(`<tr>`)
	g.in()
	g.line(`<th scope="colgroup">Eget kapital</th>`)
	g.line(`<td />`)
	g.line(`<td />`)
	g.line(`<td />`)
	g.out()
	g.line(`</tr>`)

	// Bundet eget kapital
	g.line(`<tr>`)
	g.in()
	g.line(`<th colspan="4" scope="colgroup" class="sub sep">Bundet eget kapital</th>`)
	g.out()
	g.line(`</tr>`)

	g.writeBalanceRow("Aktiekapital", 0, nil,
		"se-gen-base:Aktiekapital",
		ycv(eq.ShareCapital), false, false, false)

	g.writeBalanceRow("Reservfond", 0, nil,
		"se-gen-base:Reservfond",
		ycv(eq.ReserveFund), false, false, true)

	g.writeBalanceRow("Summa bundet eget kapital", 0, nil,
		"se-gen-base:BundetEgetKapital",
		ycv(eq.TotalRestrictedEquity), true, false, false)

	// Fritt eget kapital
	g.line(`<tr>`)
	g.in()
	g.line(`<th colspan="4" scope="colgroup" class="sub sep">Fritt eget kapital</th>`)
	g.out()
	g.line(`</tr>`)

	g.writeBalanceRow("Balanserat resultat", 0, nil,
		"se-gen-base:BalanseratResultat",
		ycv(eq.RetainedEarnings), false, false, false)

	g.writeBalanceRow("Årets resultat", 0, nil,
		"se-gen-base:AretsResultatEgetKapital",
		ycv(eq.NetIncome), false, false, true)

	g.writeBalanceRow("Summa fritt eget kapital", 0, nil,
		"se-gen-base:FrittEgetKapital",
		ycv(eq.TotalUnrestrictedEquity), true, false, false)

	// Summa eget kapital
	g.writeBalanceRow("Summa eget kapital", 0, nil,
		"se-gen-base:EgetKapital",
		ycv(eq.TotalEquity), true, false, false)

	g.out()
	g.line(`</tbody>`)
}

// writeBSUntaxedReserves writes the obeskattade reserver tbody.
func (g *generator) writeBSUntaxedReserves(el *model.EquityAndLiabilities) {
	ur := &el.UntaxedReserves

	g.line(`<tbody>`)
	g.in()

	g.line(`<tr>`)
	g.in()
	g.line(`<th colspan="4" scope="rowgroup">Obeskattade reserver</th>`)
	g.out()
	g.line(`</tr>`)

	g.writeBalanceRow("Periodiseringsfonder", 0, nil,
		"se-gen-base:Periodiseringsfonder",
		ycv(ur.TaxAllocationReserves), false, false, false)

	g.writeBalanceRow("Ackumulerade överavskrivningar", 0, nil,
		"se-gen-base:AckumuleradeOveravskrivningar",
		ycv(ur.AccumulatedExcessDepreciation), false, false, true)

	g.writeBalanceRow("Summa obeskattade reserver", 0, nil,
		"se-gen-base:ObeskattadeReserver",
		ycv(ur.TotalUntaxedReserves), true, false, false)

	g.out()
	g.line(`</tbody>`)
}

// writeBSProvisions writes the avsättningar tbody.
func (g *generator) writeBSProvisions(el *model.EquityAndLiabilities) {
	prov := &el.Provisions

	g.line(`<tbody>`)
	g.in()

	g.line(`<tr>`)
	g.in()
	g.line(`<th colspan="4" scope="rowgroup">Avsättningar</th>`)
	g.out()
	g.line(`</tr>`)

	g.writeBalanceRow("Avsättningar för pensioner och liknande förpliktelser enligt lagen (1967:531)\n          om tryggande av pensionsutfästelse <abbr>m.m.</abbr>", 0, nil,
		"se-gen-base:AvsattningarPensionerLiknandeForpliktelserEnligtLag",
		ycv(prov.PensionProvisions), false, false, false)

	g.writeBalanceRow("Övriga avsättningar", 0, nil,
		"se-gen-base:OvrigaAvsattningar",
		ycv(prov.OtherProvisions), false, false, true)

	g.writeBalanceRow("Summa avsättningar", 0, nil,
		"se-gen-base:Avsattningar",
		ycv(prov.TotalProvisions), true, false, false)

	g.out()
	g.line(`</tbody>`)
}

// writeBSLongTermLiabilities writes the långfristiga skulder tbody.
func (g *generator) writeBSLongTermLiabilities(el *model.EquityAndLiabilities) {
	lt := &el.LongTermLiabilities

	g.line(`<tbody>`)
	g.in()

	// Långfristiga skulder header with note ref
	g.line(`<tr>`)
	g.in()
	g.line(`<th colspan="1" scope="rowgroup">Långfristiga skulder</th>`)
	if lt.LongTermLiabilitiesNote > 0 {
		g.linef(`<td><a href="#note-%d">%d</a></td>`, lt.LongTermLiabilitiesNote, lt.LongTermLiabilitiesNote)
	} else {
		g.line(`<td />`)
	}
	g.line(`<td />`)
	g.line(`<td />`)
	g.out()
	g.line(`</tr>`)

	// Övriga skulder till kreditinstitut (with multiple note refs)
	g.writeBalanceRow("Övriga skulder till kreditinstitut", 0, lt.BankLoansNotes,
		"se-gen-base:OvrigaLangfristigaSkulderKreditinstitut",
		ycv(lt.BankLoans), false, false, false)

	// Övriga skulder (last in group)
	g.writeBalanceRow("Övriga skulder", 0, nil,
		"se-gen-base:OvrigaLangfristigaSkulder",
		ycv(lt.OtherLongTermLiabilities), false, false, true)

	// Summa långfristiga skulder
	g.writeBalanceRow("Summa långfristiga skulder", 0, nil,
		"se-gen-base:LangfristigaSkulder",
		ycv(lt.TotalLongTermLiabilities), true, false, false)

	g.out()
	g.line(`</tbody>`)
}

// writeBSShortTermLiabilities writes the kortfristiga skulder tbody.
func (g *generator) writeBSShortTermLiabilities(el *model.EquityAndLiabilities) {
	st := &el.ShortTermLiabilities

	g.line(`<tbody>`)
	g.in()

	g.line(`<tr>`)
	g.in()
	g.line(`<th colspan="4" scope="rowgroup">Kortfristiga skulder</th>`)
	g.out()
	g.line(`</tr>`)

	g.writeBalanceRow("Leverantörsskulder", 0, nil,
		"se-gen-base:Leverantorsskulder",
		ycv(st.TradePayables), false, false, false)

	g.writeBalanceRow("Skatteskulder", 0, nil,
		"se-gen-base:Skatteskulder",
		ycv(st.TaxLiabilities), false, false, false)

	g.writeBalanceRow("Övriga skulder", st.OtherShortTermLiabilitiesNote, nil,
		"se-gen-base:OvrigaKortfristigaSkulder",
		ycv(st.OtherShortTermLiabilities), false, false, false)

	g.writeBalanceRow("Upplupna kostnader och förutbetalda intäkter", 0, nil,
		"se-gen-base:UpplupnaKostnaderForutbetaldaIntakter",
		ycv(st.AccruedExpenses), false, false, true)

	// Summa kortfristiga skulder
	g.writeBalanceRow("Summa kortfristiga skulder", 0, nil,
		"se-gen-base:KortfristigaSkulder",
		ycv(st.TotalShortTermLiabilities), true, false, false)

	// Summa eget kapital och skulder (total)
	g.writeBalanceRow("Summa eget kapital och skulder", 0, nil,
		"se-gen-base:EgetKapitalSkulder",
		ycv(el.TotalEquityAndLiabilities), false, true, false)

	g.out()
	g.line(`</tbody>`)
}

// ycv converts a model.YearComparison to YearComparisonVal.
func ycv(yc model.YearComparison) YearComparisonVal {
	return YearComparisonVal{Current: yc.Current, Previous: yc.Previous}
}

// unused import suppressor
var _ = strings.Repeat
