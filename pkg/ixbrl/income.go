package ixbrl

import (
	"github.com/redofri/redofri/pkg/model"
)

// writeIncomeStatement writes the resultaträkning (page 4).
func (g *generator) writeIncomeStatement(r *model.AnnualReport) {
	is := &r.IncomeStatement
	totalPages := g.computeTotalPages(r)

	_, prevEnd := prevYearDates(r.FiscalYear.StartDate, r.FiscalYear.EndDate)

	g.line(`<div class="ar-page wide" id="ar3-page-4">`)
	g.in()
	g.pageHeader(r.Company.Name, r.Company.OrgNr, 4, totalPages)

	g.line(`<table class="ar-profit-loss ar-financial col-4">`)
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
	g.line(`<th scope="col">Resultaträkning</th>`)
	g.line(`<th scope="col">Not</th>`)
	g.linef(`<th scope="col">%s<br />–%s</th>`, r.FiscalYear.StartDate, r.FiscalYear.EndDate)
	g.linef(`<th scope="col">%s<br />–%s</th>`,
		prevStart(r.FiscalYear.StartDate), prevEnd)
	g.out()
	g.line(`</tr>`)
	g.out()
	g.line(`</thead>`)

	// Revenue section
	g.writeISRevenue(is)

	// Expenses section
	g.writeISExpenses(is)

	// Financial items section
	g.writeISFinancialItems(is)

	// Appropriations section
	g.writeISAppropriations(is)

	// Tax section
	g.writeISTax(is)

	g.out()
	g.line(`</table>`)

	g.out()
	g.line(`</div>`)
}

// writeISRevenue writes the revenue tbody.
func (g *generator) writeISRevenue(is *model.IncomeStatement) {
	rev := &is.Revenue

	g.line(`<tbody>`)
	g.in()

	g.line(`<tr>`)
	g.in()
	g.line(`<th colspan="4" scope="rowgroup">Rörelseintäkter, lagerförändringar <abbr>m.m.</abbr>`)
	g.line(`</th>`)
	g.out()
	g.line(`</tr>`)

	// Nettoomsättning
	g.writeYearComparisonRow("Nettoomsättning", 0,
		"se-gen-base:Nettoomsattning", "period0", "period1", "SEK",
		rev.NetSales.Current, rev.NetSales.Previous,
		false, false, false)

	// Förändring av lager...
	g.writeYearComparisonRow("Förändring av lager av produkter i arbete, färdiga varor och pågående arbete för annans räkning", 0,
		"se-gen-base:ForandringLagerProdukterIArbeteFardigaVarorPagaendeArbetenAnnansRakning", "period0", "period1", "SEK",
		rev.InventoryChange.Current, rev.InventoryChange.Previous,
		false, false, false)

	// Övriga rörelseintäkter (last in group — sum wrap)
	g.writeYearComparisonRow("Övriga rörelseintäkter", 0,
		"se-gen-base:OvrigaRorelseintakter", "period0", "period1", "SEK",
		rev.OtherOperatingIncome.Current, rev.OtherOperatingIncome.Previous,
		false, true, false)

	// Summa rörelseintäkter
	g.writeISSubTotal("Summa rörelseintäkter, lagerförändringar <abbr>m.m.</abbr>",
		"se-gen-base:RorelseintakterLagerforandringarMm",
		rev.TotalRevenue.Current, rev.TotalRevenue.Previous,
		false)

	g.out()
	g.line(`</tbody>`)
}

// writeISExpenses writes the expenses tbody.
func (g *generator) writeISExpenses(is *model.IncomeStatement) {
	exp := &is.Expenses

	g.line(`<tbody>`)
	g.in()

	g.line(`<tr>`)
	g.in()
	g.line(`<th colspan="4" scope="rowgroup">Rörelsekostnader</th>`)
	g.out()
	g.line(`</tr>`)

	// Råvaror och förnödenheter
	g.writeYearComparisonRow("Råvaror och förnödenheter", 0,
		"se-gen-base:RavarorFornodenheterKostnader", "period0", "period1", "SEK",
		exp.RawMaterials.Current, exp.RawMaterials.Previous,
		true, false, false)

	// Handelsvaror
	g.writeYearComparisonRow("Handelsvaror", 0,
		"se-gen-base:HandelsvarorKostnader", "period0", "period1", "SEK",
		exp.TradingGoods.Current, exp.TradingGoods.Previous,
		true, false, false)

	// Övriga externa kostnader
	g.writeYearComparisonRow("Övriga externa kostnader", 0,
		"se-gen-base:OvrigaExternaKostnader", "period0", "period1", "SEK",
		exp.OtherExternalExpenses.Current, exp.OtherExternalExpenses.Previous,
		true, false, false)

	// Personalkostnader (with note ref)
	g.writeYearComparisonRow("Personalkostnader", exp.PersonnelExpensesNote,
		"se-gen-base:Personalkostnader", "period0", "period1", "SEK",
		exp.PersonnelExpenses.Current, exp.PersonnelExpenses.Previous,
		true, false, false)

	// Av- och nedskrivningar
	g.writeYearComparisonRow("Av- och nedskrivningar av materiella och immateriella anläggningstillgångar", 0,
		"se-gen-base:AvskrivningarNedskrivningarMateriellaImmateriellaAnlaggningstillgangar", "period0", "period1", "SEK",
		exp.DepreciationAmortization.Current, exp.DepreciationAmortization.Previous,
		true, false, false)

	// Övriga rörelsekostnader (last in group — sum wrap)
	g.writeYearComparisonRow("Övriga rörelsekostnader", 0,
		"se-gen-base:OvrigaRorelsekostnader", "period0", "period1", "SEK",
		exp.OtherOperatingExpenses.Current, exp.OtherOperatingExpenses.Previous,
		true, true, false)

	// Summa rörelsekostnader
	g.writeISSubTotal("Summa rörelsekostnader",
		"se-gen-base:Rorelsekostnader",
		exp.TotalExpenses.Current, exp.TotalExpenses.Previous,
		true)

	// Rörelseresultat (result row)
	g.writeISResultRow("Rörelseresultat",
		"se-gen-base:Rorelseresultat",
		is.OperatingResult.Current, is.OperatingResult.Previous,
		false)

	g.out()
	g.line(`</tbody>`)
}

// writeISFinancialItems writes the financial items tbody.
func (g *generator) writeISFinancialItems(is *model.IncomeStatement) {
	fi := &is.FinancialItems

	g.line(`<tbody>`)
	g.in()

	g.line(`<tr>`)
	g.in()
	g.line(`<th colspan="4" scope="rowgroup">Finansiella poster</th>`)
	g.out()
	g.line(`</tr>`)

	// Resultat från övriga finansiella anläggningstillgångar
	g.writeYearComparisonRow("Resultat från övriga finansiella anläggningstillgångar", 0,
		"se-gen-base:ResultatOvrigaFinansiellaAnlaggningstillgangar", "period0", "period1", "SEK",
		fi.ResultOtherFinancialAssets.Current, fi.ResultOtherFinancialAssets.Previous,
		false, false, false)

	// Övriga ränteintäkter
	g.writeYearComparisonRow("Övriga ränteintäkter och liknande resultatposter", 0,
		"se-gen-base:OvrigaRanteintakterLiknandeResultatposter", "period0", "period1", "SEK",
		fi.OtherInterestIncome.Current, fi.OtherInterestIncome.Previous,
		false, false, false)

	// Räntekostnader (last in group — sum wrap, expense display)
	g.writeYearComparisonRow("Räntekostnader och liknande resultatposter", 0,
		"se-gen-base:RantekostnaderLiknandeResultatposter", "period0", "period1", "SEK",
		fi.InterestExpenses.Current, fi.InterestExpenses.Previous,
		true, true, false)

	// Summa finansiella poster
	g.writeISSubTotal("Summa finansiella poster",
		"se-gen-base:FinansiellaPoster",
		fi.TotalFinancialItems.Current, fi.TotalFinancialItems.Previous,
		false)

	// Resultat efter finansiella poster
	g.writeISResultRow("Resultat efter finansiella poster",
		"se-gen-base:ResultatEfterFinansiellaPoster",
		is.ResultAfterFinancialItems.Current, is.ResultAfterFinancialItems.Previous,
		false)

	g.out()
	g.line(`</tbody>`)
}

// writeISAppropriations writes the appropriations tbody.
func (g *generator) writeISAppropriations(is *model.IncomeStatement) {
	ap := &is.Appropriations

	g.line(`<tbody>`)
	g.in()

	g.line(`<tr>`)
	g.in()
	g.line(`<th colspan="4" scope="rowgroup">Bokslutsdispositioner</th>`)
	g.out()
	g.line(`</tr>`)

	// Förändring av periodiseringsfonder (sign="-", negPrefix)
	g.writeISAppropriationRow("Förändring av periodiseringsfonder",
		"se-gen-base:ForandringPeriodiseringsfond",
		ap.TaxAllocationReserve.Current, ap.TaxAllocationReserve.Previous,
		false)

	// Förändring av överavskrivningar (sign="-", negPrefix, last in group)
	g.writeISAppropriationRow("Förändring av överavskrivningar",
		"se-gen-base:ForandringOveravskrivningar",
		ap.ExcessDepreciation.Current, ap.ExcessDepreciation.Previous,
		true)

	// Summa bokslutsdispositioner (sign="-")
	g.writeISAppropriationSubTotal("Summa bokslutsdispositioner",
		"se-gen-base:Bokslutsdispositioner",
		ap.TotalAppropriations.Current, ap.TotalAppropriations.Previous)

	// Resultat före skatt
	g.writeISResultRow("Resultat före skatt",
		"se-gen-base:ResultatForeSkatt",
		is.ResultBeforeTax.Current, is.ResultBeforeTax.Previous,
		false)

	g.out()
	g.line(`</tbody>`)
}

// writeISTax writes the tax section and net result.
func (g *generator) writeISTax(is *model.IncomeStatement) {
	g.line(`<tbody>`)
	g.in()

	g.line(`<tr>`)
	g.in()
	g.line(`<th colspan="4" scope="rowgroup">Skatter</th>`)
	g.out()
	g.line(`</tr>`)

	// Skatt på årets resultat (expense, last in group — sum wrap)
	g.writeYearComparisonRow("Skatt på årets resultat", 0,
		"se-gen-base:SkattAretsResultat", "period0", "period1", "SEK",
		is.Tax.IncomeTax.Current, is.Tax.IncomeTax.Previous,
		true, true, false)

	// Årets resultat (total result row)
	g.writeISResultRow("Årets resultat",
		"se-gen-base:AretsResultat",
		is.NetResult.Current, is.NetResult.Previous,
		true)

	g.out()
	g.line(`</tbody>`)
}

// writeISSubTotal writes a subtotal row in the IS (with sub-sum td class).
func (g *generator) writeISSubTotal(label, concept string, current, previous *int64, isExpense bool) {
	g.line(`<tr>`)
	g.in()
	g.linef(`<td class="sum">%s</td>`, label)
	g.line(`<td />`)

	g.writeISCell(concept, "period0", current, isExpense, false, false)
	g.writeISCell(concept, "period1", previous, isExpense, false, false)

	g.out()
	g.line(`</tr>`)
}

// writeISResultRow writes a result row (tr class="result").
func (g *generator) writeISResultRow(label, concept string, current, previous *int64, isTotal bool) {
	g.line(`<tr class="result">`)
	g.in()
	g.linef(`<td>%s</td>`, label)
	g.line(`<td />`)

	g.writeISCell(concept, "period0", current, false, false, isTotal)
	g.writeISCell(concept, "period1", previous, false, false, isTotal)

	g.out()
	g.line(`</tr>`)
}

// writeISAppropriationRow writes an appropriation row (sign="-", negPrefix display).
func (g *generator) writeISAppropriationRow(label, concept string, current, previous *int64, isLastInGroup bool) {
	if current == nil && previous == nil {
		return
	}

	g.line(`<tr>`)
	g.in()
	g.linef(`<td>%s</td>`, label)
	g.line(`<td />`)

	g.writeISAppropriationCell(concept, "period0", current, isLastInGroup)
	g.writeISAppropriationCell(concept, "period1", previous, isLastInGroup)

	g.out()
	g.line(`</tr>`)
}

// writeISAppropriationSubTotal writes an appropriation subtotal row (sign="-").
func (g *generator) writeISAppropriationSubTotal(label, concept string, current, previous *int64) {
	g.line(`<tr>`)
	g.in()
	g.linef(`<td class="sum">%s</td>`, label)
	g.line(`<td />`)

	g.writeISAppropriationCell(concept, "period0", current, false)
	g.writeISAppropriationCell(concept, "period1", previous, false)

	g.out()
	g.line(`</tr>`)
}

// writeISCell writes a single IS cell with nonFraction.
func (g *generator) writeISCell(concept, contextRef string, value *int64, isExpense, isSubTotal, isTotal bool) {
	g.write(indentStr(g.indent))
	g.write("<td>")
	if value != nil {
		var opts []nfOpt
		if isExpense {
			opts = append(opts, withNegPrefix())
		}
		wrapClass := ""
		if isSubTotal {
			wrapClass = "sum"
		} else if isTotal {
			wrapClass = "total"
		}
		if wrapClass != "" {
			opts = append(opts, withWrapClass(wrapClass))
		}
		g.nonFraction(concept, contextRef, "SEK", *value, opts...)
	}
	g.write("</td>\n")
}

// writeISAppropriationCell writes a cell for appropriation items (sign="-", negPrefix).
func (g *generator) writeISAppropriationCell(concept, contextRef string, value *int64, isLastInGroup bool) {
	g.write(indentStr(g.indent))
	g.write("<td>")
	if value != nil {
		opts := []nfOpt{withSign("-"), withNegPrefix()}
		if isLastInGroup {
			opts = append(opts, withWrapClass("sum"))
		}
		g.nonFraction(concept, contextRef, "SEK", *value, opts...)
	}
	g.write("</td>\n")
}

// prevStart computes the previous fiscal year start date string.
func prevStart(startDate string) string {
	ps, _ := prevYearDates(startDate, startDate)
	return ps
}

// indentStr returns a tab indentation string.
func indentStr(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += "\t"
	}
	return s
}
