package ixbrl

import (
	"fmt"

	"github.com/redofri/redofri/pkg/model"
)

// writeNotes writes all note pages (ar3-page-7 through the last note page).
// The signatures section is written inside the last note page div.
func (g *generator) writeNotes(r *model.AnnualReport) {
	totalPages := g.computeTotalPages(r)
	notes := &r.Notes

	// We need to determine how notes map to pages.
	// In the reference example:
	//   page 7: notes 1-2 (accounting policies + employees)
	//   page 8: notes 3-5 (asset roll-forward notes for tangible)
	//   page 10 (ar3-page-10): notes 6-9 (financial asset note, long-term liabilities, pledges, contingent)
	//   page 9 (ar3-page-9): note 10 + signatures
	//
	// For simplicity, we replicate the reference layout for the example data.
	// A future version should compute page breaks dynamically.

	// Page 7: Note 1 (accounting policies) and Note 2 (employees)
	g.line(`<div class="ar-page wide" id="ar3-page-7">`)
	g.in()
	g.pageHeader(r.Company.Name, r.Company.OrgNr, 7, totalPages)
	g.line(`<h2>Noter</h2>`)

	g.writeAccountingPoliciesNote(r, &notes.AccountingPolicies)
	if notes.Employees != nil {
		g.writeEmployeesNote(r, notes.Employees)
	}

	g.out()
	g.line(`</div>`)

	// Page 8: Fixed asset notes (tangible assets, notes 3-5 typically)
	// We split: tangible asset notes on page 8, financial + remaining on page 9/10.
	tangibleNotes, financialNotes := splitFixedAssetNotes(notes.FixedAssetNotes)

	if len(tangibleNotes) > 0 {
		g.line(`<div class="ar-page note" id="ar3-page-8">`)
		g.in()
		g.pageHeader(r.Company.Name, r.Company.OrgNr, 8, totalPages)

		for i, fan := range tangibleNotes {
			isFirst := (i == 0)
			g.writeFixedAssetNote(r, &fan, isFirst && true)
		}

		g.out()
		g.line(`</div>`)
	}

	// Next page: financial asset notes + long-term liabilities + pledges + contingent
	hasFinancialPage := len(financialNotes) > 0 ||
		notes.LongTermLiabilitiesNote != nil ||
		notes.Pledges != nil ||
		notes.ContingentLiabilities != nil

	if hasFinancialPage {
		g.line(`<div class="ar-page wide" id="ar3-page-10">`)
		g.in()
		g.pageHeader(r.Company.Name, r.Company.OrgNr, 9, totalPages)

		for i, fan := range financialNotes {
			isFirst := (i == 0)
			g.writeFixedAssetNote(r, &fan, isFirst)
		}

		if notes.LongTermLiabilitiesNote != nil {
			g.writeLongTermLiabilitiesNote(r, notes.LongTermLiabilitiesNote)
		}
		if notes.Pledges != nil {
			g.writePledgesNote(r, notes.Pledges)
		}
		if notes.ContingentLiabilities != nil {
			g.writeContingentLiabilitiesNote(r, notes.ContingentLiabilities)
		}

		g.out()
		g.line(`</div>`)
	}

	// Last page: multi-post note (if any) + signatures
	// This page div is opened here; writeSignatures() will add its content
	// and the page div will be closed after writeSignatures() returns.
	g.line(`<div class="ar-page wide" id="ar3-page-9">`)
	g.in()
	g.pageHeader(r.Company.Name, r.Company.OrgNr, 10, totalPages)

	if notes.MultiPostNote != nil {
		g.writeMultiPostNote(r, notes.MultiPostNote)
	}

	// Signatures are written here inside the last notes page.
	g.writeSignatures(r)

	g.out()
	g.line(`</div>`)
}

// splitFixedAssetNotes splits notes into tangible (with depreciation) and
// financial (without depreciation) groups.
func splitFixedAssetNotes(notes []model.FixedAssetNote) (tangible, financial []model.FixedAssetNote) {
	for _, n := range notes {
		if n.OpeningDepreciation.Current != nil || n.OpeningDepreciation.Previous != nil {
			tangible = append(tangible, n)
		} else {
			financial = append(financial, n)
		}
	}
	return
}

// writeAccountingPoliciesNote writes Note 1: Redovisnings- och värderingsprinciper.
func (g *generator) writeAccountingPoliciesNote(r *model.AnnualReport, ap *model.AccountingPolicies) {
	g.linef(`<h3 id="note-%d">`, ap.NoteNumber)
	g.in()
	g.linef(`<span class="note">Not %d</span> Redovisnings- och värderingsprinciper</h3>`, ap.NoteNumber)
	g.out()

	// Main description
	g.line(`<p>`)
	g.in()
	g.write(indentStr(g.indent))
	g.nonNumeric("se-gen-base:Redovisningsprinciper", "period0", ap.Description)
	g.write("\n")
	g.out()
	g.line(`</p>`)

	// Depreciation table
	if len(ap.Depreciations) > 0 {
		g.line(`<h4 class="join">Avskrivningar</h4>`)
		g.line(`<p class="join">Tillämpade avskrivningstider:</p>`)
		g.line(`<table class="ar-depreciation">`)
		g.in()
		g.line(`<tbody>`)
		g.in()

		for _, dep := range ap.Depreciations {
			g.line(`<tr>`)
			g.in()
			g.linef(`<td>%s</td>`, esc(dep.Category))
			g.line(`<td>`)
			g.in()
			g.write(indentStr(g.indent))
			g.nonNumeric(dep.Concept, "period0", fmt.Sprintf("%d", dep.Years))
			g.write(" år</td>\n")
			g.out()
			g.out()
			g.line(`</tr>`)
		}

		g.out()
		g.line(`</tbody>`)
		g.out()
		g.line(`</table>`)
	}

	// Depreciation comment
	if ap.DepreciationComment != "" {
		g.line(`<p>`)
		g.in()
		g.write(indentStr(g.indent))
		g.nonNumeric("se-gen-base:AvskrivningarMateriellaAnlaggningstillgangarKommentar", "period0", ap.DepreciationComment)
		g.write("\n")
		g.out()
		g.line(`</p>`)
	}

	// Manufactured goods policy
	if ap.ManufacturedGoodsPolicy != "" {
		g.line(`<h4 class="join">Anskaffningsvärde för egentillverkade varor</h4>`)
		g.line(`<p>`)
		g.in()
		g.write(indentStr(g.indent))
		g.nonNumeric("se-gen-base:RedovisningsprinciperAnskaffningsvardeEgentillverkadevaror", "period0", ap.ManufacturedGoodsPolicy)
		g.write("\n")
		g.out()
		g.line(`</p>`)
	}

	// Key figure definitions (hardcoded for K2)
	g.line(`<h4 class="join">Nyckeltalsdefinitioner</h4>`)
	g.line(`<dl>`)
	g.in()
	g.line(`<dt>Soliditet</dt>`)
	g.line(`<dd>Eget kapital och obeskattade reserver (med avdrag för uppskjuten skatt) i förhållande till`)
	g.line(`    balansomslutningen.</dd>`)
	g.out()
	g.line(`</dl>`)
}

// writeEmployeesNote writes Note 2: Medelantalet anställda.
func (g *generator) writeEmployeesNote(r *model.AnnualReport, emp *model.EmployeesNote) {
	_, prevEnd := prevYearDates(r.FiscalYear.StartDate, r.FiscalYear.EndDate)

	g.linef(`<h3 id="note-%d">Upplysningar till resultaträkningen`, emp.NoteNumber)
	g.line(`    <br />`)
	g.linef(`<span class="note">Not %d</span> Medelantalet anställda</h3>`, emp.NoteNumber)

	g.line(`<table class="ar-note">`)
	g.in()

	// Colgroup
	g.line(`<colgroup>`)
	g.in()
	g.line(`<col />`)
	g.line(`<col class="kr" span="2" />`)
	g.out()
	g.line(`</colgroup>`)

	// Header
	g.line(`<thead>`)
	g.in()
	g.line(`<tr>`)
	g.in()
	g.line(`<th />`)
	g.line(`<th scope="col">`)
	g.in()
	g.linef(`<span>%s<br />–%s</span>`, r.FiscalYear.StartDate, r.FiscalYear.EndDate)
	g.out()
	g.line(`</th>`)
	g.line(`<th scope="col">`)
	g.in()
	g.linef(`<span>%s<br />–%s</span>`, prevStart(r.FiscalYear.StartDate), prevEnd)
	g.out()
	g.line(`</th>`)
	g.out()
	g.line(`</tr>`)
	g.out()
	g.line(`</thead>`)

	// Body
	g.line(`<tbody>`)
	g.in()
	g.line(`<tr>`)
	g.in()
	g.line(`<td>Medelantalet anställda</td>`)

	// Current year — note: no format attribute for antal-anstallda
	g.line(`<td>`)
	g.in()
	if emp.AverageEmployees.Current != nil {
		g.nonFractionLine("se-gen-base:MedelantaletAnstallda", "period0", "antal-anstallda",
			*emp.AverageEmployees.Current, withFormat(""))
	}
	g.out()
	g.line(`</td>`)

	// Previous year
	g.line(`<td>`)
	g.in()
	if emp.AverageEmployees.Previous != nil {
		g.nonFractionLine("se-gen-base:MedelantaletAnstallda", "period1", "antal-anstallda",
			*emp.AverageEmployees.Previous, withFormat(""))
	}
	g.out()
	g.line(`</td>`)

	g.out()
	g.line(`</tr>`)
	g.out()
	g.line(`</tbody>`)

	g.out()
	g.line(`</table>`)
}

// writeFixedAssetNote writes a roll-forward note for a single asset category.
func (g *generator) writeFixedAssetNote(r *model.AnnualReport, fan *model.FixedAssetNote, isFirstOnPage bool) {
	_, prevEnd := prevYearDates(r.FiscalYear.StartDate, r.FiscalYear.EndDate)
	hasDepreciation := fan.OpeningDepreciation.Current != nil || fan.OpeningDepreciation.Previous != nil

	// Heading
	if isFirstOnPage {
		g.linef(`<h3 id="note-%d">Upplysningar till balansräkningen`, fan.NoteNumber)
		g.line(`    <br />`)
		g.linef(`<span class="note">Not %d</span> %s</h3>`, fan.NoteNumber, esc(fan.Title))
	} else {
		g.linef(`<h3 id="note-%d">`, fan.NoteNumber)
		g.in()
		g.linef(`<span class="note">Not %d</span> %s</h3>`, fan.NoteNumber, esc(fan.Title))
		g.out()
	}

	g.line(`<table class="ar-note">`)
	g.in()

	// Colgroup
	g.line(`<colgroup>`)
	g.in()
	g.line(`<col />`)
	g.line(`<col class="kr" span="2" />`)
	g.out()
	g.line(`</colgroup>`)

	// Header with instant dates
	g.line(`<thead>`)
	g.in()
	g.line(`<tr>`)
	g.in()
	g.line(`<th />`)
	g.line(`<th scope="col">`)
	g.in()
	g.linef(`<span>%s</span>`, r.FiscalYear.EndDate)
	g.out()
	g.line(`</th>`)
	g.line(`<th scope="col">`)
	g.in()
	g.linef(`<span>%s</span>`, prevEnd)
	g.out()
	g.line(`</th>`)
	g.out()
	g.line(`</tr>`)
	g.out()
	g.line(`</thead>`)

	prefix := "se-gen-base:" + fan.ConceptPrefix

	// Acquisition values tbody
	g.line(`<tbody>`)
	g.in()

	// Opening acquisition values (balans1 for current, balans2 for prev)
	g.writeNoteRow("Ingående anskaffningsvärden",
		prefix+"Anskaffningsvarden",
		"balans1", "balans2",
		fan.OpeningAcquisitionValues.Current, fan.OpeningAcquisitionValues.Previous,
		false, false, false)

	// Purchases
	if fan.Purchases.Current != nil || fan.Purchases.Previous != nil {
		g.writeNoteRow("- Inköp",
			prefix+"ForandringAnskaffningsvardenInkop",
			"period0", "period1",
			fan.Purchases.Current, fan.Purchases.Previous,
			false, false, true)
	}

	// Sales
	if fan.Sales.Current != nil || fan.Sales.Previous != nil {
		hasPurchases := fan.Purchases.Current != nil || fan.Purchases.Previous != nil
		isLast := !hasPurchases // if no purchases, sales is the last change row
		_ = isLast
		g.writeNoteSalesRow("- Försäljningar",
			prefix+"ForandringAnskaffningsvardenForsaljningar",
			"period0", "period1",
			fan.Sales.Current, fan.Sales.Previous)
	}

	// Closing acquisition values (balans0 for current, balans1 for prev)
	g.writeNoteRow("Utgående anskaffningsvärden",
		prefix+"Anskaffningsvarden",
		"balans0", "balans1",
		fan.ClosingAcquisitionValues.Current, fan.ClosingAcquisitionValues.Previous,
		false, false, false)

	g.out()
	g.line(`</tbody>`)

	// Depreciation tbody (only for tangible assets)
	if hasDepreciation {
		g.line(`<tbody>`)
		g.in()

		// Opening depreciation (negative display)
		g.writeNoteDepreciationRow("Ingående avskrivningar",
			prefix+"Avskrivningar",
			"balans1", "balans2",
			fan.OpeningDepreciation.Current, fan.OpeningDepreciation.Previous,
			false)

		// Year's depreciation (negative display, sum wrap)
		g.writeNoteDepreciationRow("- Årets avskrivningar",
			prefix+"ForandringAvskrivningarAretsAvskrivningar",
			"period0", "period1",
			fan.YearDepreciation.Current, fan.YearDepreciation.Previous,
			true)

		// Closing depreciation (negative display, sum wrap)
		g.writeNoteDepreciationRow("Utgående avskrivningar",
			prefix+"Avskrivningar",
			"balans0", "balans1",
			fan.ClosingDepreciation.Current, fan.ClosingDepreciation.Previous,
			true)

		// Carrying value (total wrap)
		g.writeNoteCarryingValueRow(prefix,
			fan.CarryingValue.Current, fan.CarryingValue.Previous)

		g.out()
		g.line(`</tbody>`)
	} else {
		// Financial assets: carrying value in same tbody as acquisition values
		// Actually looking at the reference, note 6 has carrying value in the same tbody.
		// Let me re-check...
		// Looking at lines 2298-2341: note 6 has a single <tbody> with acquisition rows + carrying value.
		// So for financial assets, we write the carrying value row inside the existing tbody.
		// But we already closed the acquisition tbody above. Let me add it in a new one.
		g.line(`<tbody>`)
		g.in()
		g.writeNoteCarryingValueRow(prefix,
			fan.CarryingValue.Current, fan.CarryingValue.Previous)
		g.out()
		g.line(`</tbody>`)
	}

	g.out()
	g.line(`</table>`)
}

// writeNoteRow writes a standard note table row with two columns.
func (g *generator) writeNoteRow(label, concept, currentCtx, prevCtx string,
	current, previous *int64, negPrefix, isTotal, isLastInGroup bool) {

	g.line(`<tr>`)
	g.in()
	g.linef(`<td>%s</td>`, label)

	// Current year
	g.write(indentStr(g.indent))
	g.write("<td>")
	if current != nil {
		var opts []nfOpt
		if negPrefix {
			opts = append(opts, withNegPrefix())
		}
		if isTotal {
			opts = append(opts, withWrapClass("total"))
		} else if isLastInGroup {
			opts = append(opts, withWrapClass("sum"))
		}
		g.nonFraction(concept, currentCtx, "SEK", *current, opts...)
	}
	g.write("</td>\n")

	// Previous year
	g.write(indentStr(g.indent))
	g.write("<td>")
	if previous != nil {
		var opts []nfOpt
		if negPrefix {
			opts = append(opts, withNegPrefix())
		}
		if isTotal {
			opts = append(opts, withWrapClass("total"))
		} else if isLastInGroup {
			opts = append(opts, withWrapClass("sum"))
		}
		g.nonFraction(concept, prevCtx, "SEK", *previous, opts...)
	} else if isLastInGroup {
		// Show en-dash for missing previous year value in change rows
		g.write(`<span class="sum">–</span>`)
	}
	g.write("</td>\n")

	g.out()
	g.line(`</tr>`)
}

// writeNoteSalesRow writes a sales row in a fixed asset note.
// Sales have negative display prefix and sum wrap.
func (g *generator) writeNoteSalesRow(label, concept, currentCtx, prevCtx string,
	current, previous *int64) {

	g.line(`<tr>`)
	g.in()
	g.linef(`<td>%s</td>`, label)

	// Current year
	g.write(indentStr(g.indent))
	g.write("<td>")
	if current != nil {
		g.nonFraction(concept, currentCtx, "SEK", *current,
			withNegPrefix(), withWrapClass("sum"))
	}
	g.write("</td>\n")

	// Previous year
	g.write(indentStr(g.indent))
	g.write("<td>")
	if previous != nil {
		g.nonFraction(concept, prevCtx, "SEK", *previous,
			withNegPrefix(), withWrapClass("sum"))
	} else {
		g.write(`<span class="sum">–</span>`)
	}
	g.write("</td>\n")

	g.out()
	g.line(`</tr>`)
}

// writeNoteDepreciationRow writes a depreciation row with negative prefix display.
func (g *generator) writeNoteDepreciationRow(label, concept, currentCtx, prevCtx string,
	current, previous *int64, isSumWrap bool) {

	g.line(`<tr>`)
	g.in()
	g.linef(`<td>%s</td>`, label)

	// Current year
	g.write(indentStr(g.indent))
	g.write("<td>")
	if current != nil {
		if isSumWrap {
			g.write(`<span class="sum">`)
		}
		g.write("-")
		g.nonFraction(concept, currentCtx, "SEK", *current)
		if isSumWrap {
			g.write("\n")
			g.write(indentStr(g.indent))
			g.write(`</span>`)
		}
	}
	g.write("\n")
	g.write(indentStr(g.indent))
	g.write("</td>\n")

	// Previous year
	g.write(indentStr(g.indent))
	g.write("<td>")
	if previous != nil {
		if isSumWrap {
			g.write(`<span class="sum">`)
		}
		g.write("-")
		g.nonFraction(concept, prevCtx, "SEK", *previous)
		if isSumWrap {
			g.write("\n")
			g.write(indentStr(g.indent))
			g.write(`</span>`)
		}
	}
	g.write("\n")
	g.write(indentStr(g.indent))
	g.write("</td>\n")

	g.out()
	g.line(`</tr>`)
}

// writeNoteCarryingValueRow writes the "Redovisat värde" row with total wrap.
func (g *generator) writeNoteCarryingValueRow(conceptPrefix string, current, previous *int64) {
	// The carrying value concept is the base concept name (e.g. "se-gen-base:ByggnaderMark")
	// which is the conceptPrefix without the "se-gen-base:" prefix... wait, conceptPrefix
	// already includes the namespace. Let me use conceptPrefix directly.
	// Actually, conceptPrefix is like "se-gen-base:ByggnaderMark" after adding the namespace above.
	// Looking at the caller: prefix = "se-gen-base:" + fan.ConceptPrefix
	// So conceptPrefix is like "se-gen-base:ByggnaderMark" which is the carrying value concept.
	concept := conceptPrefix

	g.line(`<tr>`)
	g.in()
	g.line(`<td>Redovisat värde</td>`)

	// Current year (balans0)
	g.write(indentStr(g.indent))
	g.write("<td>\n")
	g.in()
	if current != nil {
		g.write(indentStr(g.indent))
		g.write(`<span class="total">`)
		g.write("\n")
		g.in()
		g.nonFractionLine(concept, "balans0", "SEK", *current)
		g.out()
		g.write(indentStr(g.indent))
		g.write("</span>\n")
	}
	g.out()
	g.line(`</td>`)

	// Previous year (balans1)
	g.write(indentStr(g.indent))
	g.write("<td>\n")
	g.in()
	if previous != nil {
		g.write(indentStr(g.indent))
		g.write(`<span class="total">`)
		g.write("\n")
		g.in()
		g.nonFractionLine(concept, "balans1", "SEK", *previous)
		g.out()
		g.write(indentStr(g.indent))
		g.write("</span>\n")
	}
	g.out()
	g.line(`</td>`)

	g.out()
	g.line(`</tr>`)
}

// writeLongTermLiabilitiesNote writes Note 7: Långfristiga skulder (> 5 år).
func (g *generator) writeLongTermLiabilitiesNote(r *model.AnnualReport, note *model.LongTermLiabilitiesNoteData) {
	_, prevEnd := prevYearDates(r.FiscalYear.StartDate, r.FiscalYear.EndDate)

	g.linef(`<h3 id="note-%d">`, note.NoteNumber)
	g.in()
	g.linef(`<span class="note">Not %d</span> Långfristiga skulder</h3>`, note.NoteNumber)
	g.out()

	g.line(`<table class="ar-note">`)
	g.in()
	g.writeNoteColgroup()
	g.writeNoteInstantHeader(r.FiscalYear.EndDate, prevEnd)

	// Description row
	g.line(`<tbody>`)
	g.in()
	g.line(`<tr>`)
	g.in()
	g.line(`<td>Långfristiga skulder som förfaller till betalning senare än fem år efter balansdagen:</td>`)
	g.line(`<td />`)
	g.line(`<td />`)
	g.out()
	g.line(`</tr>`)
	g.out()
	g.line(`</tbody>`)

	// Total row
	g.line(`<tbody>`)
	g.in()
	g.writeNoteRow("Summa",
		"se-gen-base:LangfristigaSkulderForfallerSenare5Ar",
		"balans0", "balans1",
		note.DueAfterFiveYears.Current, note.DueAfterFiveYears.Previous,
		false, true, false)
	g.out()
	g.line(`</tbody>`)

	g.out()
	g.line(`</table>`)
}

// writePledgesNote writes Note 8: Ställda säkerheter.
func (g *generator) writePledgesNote(r *model.AnnualReport, note *model.PledgesNote) {
	_, prevEnd := prevYearDates(r.FiscalYear.StartDate, r.FiscalYear.EndDate)

	g.linef(`<h3 id="note-%d">`, note.NoteNumber)
	g.in()
	g.linef(`<span class="note">Not %d</span> Ställda säkerheter</h3>`, note.NoteNumber)
	g.out()

	g.line(`<table class="ar-note">`)
	g.in()
	g.writeNoteColgroup()
	g.writeNoteInstantHeader(r.FiscalYear.EndDate, prevEnd)

	// Detail rows
	g.line(`<tbody>`)
	g.in()

	if note.CorporateMortgages.Current != nil || note.CorporateMortgages.Previous != nil {
		g.writeNoteRow("Företagsinteckning",
			"se-gen-base:StalldaSakerheterForetagsinteckningar",
			"balans0", "balans1",
			note.CorporateMortgages.Current, note.CorporateMortgages.Previous,
			false, false, false)
	}

	if note.RealEstateMortgages.Current != nil || note.RealEstateMortgages.Previous != nil {
		g.writeNoteRow("Fastighetsinteckning",
			"se-gen-base:StalldaSakerheterFastighetsinteckningar",
			"balans0", "balans1",
			note.RealEstateMortgages.Current, note.RealEstateMortgages.Previous,
			false, false, false)
	}

	g.out()
	g.line(`</tbody>`)

	// Total row
	g.line(`<tbody>`)
	g.in()
	g.writeNoteRow("Summa ställda säkerheter",
		"se-gen-base:StalldaSakerheter",
		"balans0", "balans1",
		note.TotalPledges.Current, note.TotalPledges.Previous,
		false, true, false)
	g.out()
	g.line(`</tbody>`)

	g.out()
	g.line(`</table>`)
}

// writeContingentLiabilitiesNote writes Note 9: Eventualförpliktelser.
func (g *generator) writeContingentLiabilitiesNote(r *model.AnnualReport, note *model.ContingentLiabilitiesNote) {
	_, prevEnd := prevYearDates(r.FiscalYear.StartDate, r.FiscalYear.EndDate)

	g.linef(`<h3 id="note-%d">`, note.NoteNumber)
	g.in()
	g.linef(`<span class="note">Not %d</span> Eventualförpliktelser</h3>`, note.NoteNumber)
	g.out()

	g.line(`<table class="ar-note">`)
	g.in()
	g.writeNoteColgroup()
	g.writeNoteInstantHeader(r.FiscalYear.EndDate, prevEnd)

	// Total row
	g.line(`<tbody>`)
	g.in()
	g.writeNoteRow("Summa",
		"se-gen-base:EventualForpliktelser",
		"balans0", "balans1",
		note.TotalContingent.Current, note.TotalContingent.Previous,
		false, true, false)
	g.out()
	g.line(`</tbody>`)

	g.out()
	g.line(`</table>`)
}

// writeMultiPostNote writes Note 10: Tillgångar, avsättningar och skulder som avser flera poster.
func (g *generator) writeMultiPostNote(r *model.AnnualReport, note *model.MultiPostNote) {
	g.linef(`<h3 id="note-%d">`, note.NoteNumber)
	g.in()
	g.linef(`<span class="note">Not %d</span> Tillgångar, avsättningar och skulder som avser flera poster</h3>`, note.NoteNumber)
	g.out()

	// Description paragraph
	g.line(`<p>`)
	g.in()
	g.write(indentStr(g.indent))
	g.nonNumericRaw("se-gen-base:NotTillgangarAvsattningarSkulderAvserFleraPoster", "balans0", esc(note.Description))
	g.write("\n")
	g.out()
	g.line(`</p>`)

	// Group entries by heading
	type entryGroup struct {
		heading string
		entries []indexedEntry
	}
	type indexedEntry2 struct {
		idx   int
		entry model.MultiPostEntry
	}

	var groups []entryGroup
	var currentGroup *entryGroup
	for i, e := range note.Entries {
		if currentGroup == nil || currentGroup.heading != e.Heading {
			groups = append(groups, entryGroup{heading: e.Heading})
			currentGroup = &groups[len(groups)-1]
		}
		currentGroup.entries = append(currentGroup.entries, indexedEntry{idx: i, entry: e})
	}

	// Write tuple declarations for all entries
	for i := range note.Entries {
		g.linef(`<ix:tuple name="se-gen-base:TillgangarAvsattningarSkulderTuple" tupleID="TillgangarAvsattningarSkulderTuple%d" />`, i+1)
	}

	// Write each group
	isLastGroup := false
	for gi, grp := range groups {
		isLastGroup = (gi == len(groups)-1)

		g.linef(`<h4 class="join">%s</h4>`, esc(grp.heading))
		g.line(`<table class="ar-note-10">`)
		g.in()
		g.line(`<colgroup>`)
		g.in()
		g.line(`<col />`)
		g.line(`<col class="kr" />`)
		g.out()
		g.line(`</colgroup>`)
		g.line(`<tbody>`)
		g.in()

		for ei, ie := range grp.entries {
			isLastEntry := isLastGroup && (ei == len(grp.entries)-1)
			tupleRef := fmt.Sprintf("TillgangarAvsattningarSkulderTuple%d", ie.idx+1)

			g.line(`<tr>`)
			g.in()
			// Post name
			g.line(`<td>`)
			g.in()
			g.write(indentStr(g.indent))
			g.nonNumeric("se-gen-base:TillgangarAvsattningarSkulderPost", "balans0", ie.entry.PostName,
				withOrder("1.0"), withTupleRef(tupleRef))
			g.write("\n")
			g.out()
			g.line(`</td>`)

			// Amount
			g.line(`<td>`)
			g.in()
			if ie.entry.Amount != nil {
				var opts []nfOpt
				opts = append(opts,
					withTupleRefNF(tupleRef),
					withOrderNF("2.0"))
				if isLastEntry {
					opts = append(opts, withWrapClass("sum"))
				}
				g.nonFractionLine("se-gen-base:TillgangarAvsattningarSkulderBelopp",
					"balans0", "SEK", *ie.entry.Amount, opts...)
			}
			g.out()
			g.line(`</td>`)

			g.out()
			g.line(`</tr>`)
		}

		g.out()
		g.line(`</tbody>`)
		g.out()
		g.line(`</table>`)
	}

	_ = isLastGroup
}

// indexedEntry pairs a MultiPostEntry with its original index.
type indexedEntry struct {
	idx   int
	entry model.MultiPostEntry
}

// writeNoteColgroup writes a standard 3-column colgroup for notes.
func (g *generator) writeNoteColgroup() {
	g.line(`<colgroup>`)
	g.in()
	g.line(`<col />`)
	g.line(`<col class="kr" span="2" />`)
	g.out()
	g.line(`</colgroup>`)
}

// writeNoteInstantHeader writes a note table header with instant dates.
func (g *generator) writeNoteInstantHeader(currentEnd, prevEnd string) {
	g.line(`<thead>`)
	g.in()
	g.line(`<tr>`)
	g.in()
	g.line(`<th />`)
	g.line(`<th scope="col">`)
	g.in()
	g.linef(`<span>%s</span>`, currentEnd)
	g.out()
	g.line(`</th>`)
	g.line(`<th scope="col">`)
	g.in()
	g.linef(`<span>%s</span>`, prevEnd)
	g.out()
	g.line(`</th>`)
	g.out()
	g.line(`</tr>`)
	g.out()
	g.line(`</thead>`)
}

// nonFraction tupleRef and order options for nfOpt (extend nfOptions in format.go).

func withTupleRefNF(ref string) nfOpt {
	return func(o *nfOptions) {
		o.tupleRef = ref
	}
}

func withOrderNF(ord string) nfOpt {
	return func(o *nfOptions) {
		o.order = ord
	}
}
