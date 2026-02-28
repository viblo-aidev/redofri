package ixbrl

import (
	"github.com/redofri/redofri/pkg/model"
)

// writeCoverPage writes the first page: company info, table of contents, and fastställelseintyg.
func (g *generator) writeCoverPage(r *model.AnnualReport) {
	totalPages := g.computeTotalPages(r)
	yearLabel := fiscalYearLabel(r.FiscalYear.StartDate, r.FiscalYear.EndDate)

	g.line(`<div class="ar-page" id="ar3-page-1">`)
	g.in()

	// Page header
	g.pageHeader(r.Company.Name, r.Company.OrgNr, 1, totalPages)

	// Company logo / name block
	g.line(`<div class="ar-logo">`)
	g.in()
	g.write("\t")
	g.nonNumeric("se-cd-base:ForetagetsNamn", "period0", r.Company.Name)
	g.write("\n")
	g.line(`<br />`)
	g.line(`<abbr>Org nr</abbr>`)
	g.write("\t")
	g.nonNumeric("se-cd-base:Organisationsnummer", "period0", r.Company.OrgNr)
	g.write("\n")
	g.out()
	g.line(`</div>`)

	// Year heading
	g.linef(`<h2>Årsredovisning för räkenskapsåret %s</h2>`, yearLabel)

	// Intro text
	g.write("\t\t\t\t")
	g.write(`<p>`)
	g.nonNumeric("se-gen-base:LopandeBokforingenAvslutasMening", "period0", r.ManagementReport.IntroText)
	g.write(`.</p>`)
	g.write("\n")

	// Table of contents
	g.writeTOC(r, totalPages)

	// Standard note about amounts
	g.line(`<p>Om inte annat särskilt anges, redovisas alla belopp i hela kronor. Uppgifter inom parentes avser föregående år.</p>`)

	// Fastställelseintyg
	g.writeCertification(r)

	g.out()
	g.line(`</div>`)
}

// writeTOC writes the table of contents.
func (g *generator) writeTOC(r *model.AnnualReport, totalPages int) {
	g.line(`<table class="ar-toc">`)
	g.in()
	g.line(`<thead>`)
	g.in()
	g.line(`<tr>`)
	g.in()
	g.line(`<th scope="col">Innehåll</th>`)
	g.line(`<th scope="col">Sida</th>`)
	g.out()
	g.line(`</tr>`)
	g.out()
	g.line(`</thead>`)
	g.line(`<tbody>`)
	g.in()

	g.writeTOCRow("förvaltningsberättelse", "ar3-page-2", 2)
	g.writeTOCRow("resultaträkning", "ar3-page-4", 4)
	g.writeTOCRow("balansräkning", "ar3-page-5", 5)
	g.writeTOCRow("noter", "ar3-page-7", 7)

	g.out()
	g.line(`</tbody>`)
	g.out()
	g.line(`</table>`)
}

// writeTOCRow writes a single table of contents row.
func (g *generator) writeTOCRow(label, pageID string, pageNum int) {
	g.line(`<tr>`)
	g.in()
	g.linef(`<td><span>-</span> %s</td>`, label)
	g.linef(`<td><a href="#%s">%d</a></td>`, pageID, pageNum)
	g.out()
	g.line(`</tr>`)
}

// writeCertification writes the fastställelseintyg block.
func (g *generator) writeCertification(r *model.AnnualReport) {
	cert := &r.Certification

	g.line(`<div id="ar-certification">`)
	g.in()
	g.line(`<strong>Fastställelseintyg</strong><br/>`)

	// Main certification text with continuation
	g.write("\t\t\t\t\t")
	g.write(`<p>`)
	// ArsstammaIntygande wraps the first part, with continuedAt
	g.writef(`<ix:nonNumeric name="se-bol-base:ArsstammaIntygande" contextRef="balans0" continuedAt="intygande_forts">`)
	g.nonNumeric("se-bol-base:FaststallelseResultatBalansrakning", "balans0", cert.ConfirmationText)
	g.write(` `)
	g.nonNumeric("se-bol-base:Arsstamma", "balans0", cert.MeetingDate)
	g.write(`. <br/>`)
	g.nonNumeric("se-bol-base:ArsstammaResultatDispositionGodkannaStyrelsensForslag", "balans0", cert.DispositionDecision)
	g.write(`</ix:nonNumeric>`)
	g.write(`</p>`)
	g.write("\n")

	// Continuation: original content certification
	g.write("\t\t\t\t\t")
	g.write(`<p>`)
	g.write(`<ix:continuation id="intygande_forts"> `)
	g.nonNumeric("se-bol-base:IntygandeOriginalInnehall", "balans0", cert.OriginalContentCertification)
	g.write(`</ix:continuation>`)
	g.write(`</p>`)
	g.write("\n")

	// Electronic signature
	g.write("\t\t\t\t\t")
	g.write(`<p><strong>`)
	g.nonNumeric("se-bol-base:UnderskriftFaststallelseintygElektroniskt", "balans0", cert.ElectronicSignatureLabel)
	g.write(`:</strong><br/>`)
	g.write("\n")
	g.write("\t\t\t\t\t")
	g.nonNumeric("se-bol-base:UnderskriftFaststallelseintygForetradareTilltalsnamn", "period0", cert.Signatory.FirstName)
	g.write(` `)
	g.nonNumeric("se-bol-base:UnderskriftFaststallelseintygForetradareEfternamn", "period0", cert.Signatory.LastName)
	g.write(`, `)
	g.nonNumeric("se-bol-base:UnderskriftFaststallelseintygForetradareForetradarroll", "period0", cert.Signatory.Role)
	g.write(`<br/>`)
	g.write("\n")
	g.write("\t\t\t\t\t")
	g.nonNumeric("se-bol-base:UnderskriftFastallelseintygDatum", "balans0", cert.SigningDate,
		withID("ID_DATUM_UNDERTECKNANDE_FASTSTALLELSEINTYG"))
	g.write("\n")
	g.write("\t\t\t\t\t</p>\n")

	g.out()
	g.line(`</div>`)
}

// computeTotalPages estimates the total page count for the document.
// This is a rough estimate; the example has 10 pages.
func (g *generator) computeTotalPages(r *model.AnnualReport) int {
	pages := 6 // cover + 2 management + IS + BS assets + BS equity

	// Notes pages (rough: 1 page per 3 notes)
	noteCount := 1 // accounting policies always
	if r.Notes.Employees != nil {
		noteCount++
	}
	noteCount += len(r.Notes.FixedAssetNotes)
	if r.Notes.LongTermLiabilitiesNote != nil {
		noteCount++
	}
	if r.Notes.Pledges != nil {
		noteCount++
	}
	if r.Notes.ContingentLiabilities != nil {
		noteCount++
	}
	if r.Notes.MultiPostNote != nil {
		noteCount++
	}
	pages += (noteCount + 2) / 3 // roughly 3 notes per page
	if pages < 10 {
		pages = 10
	}
	return pages
}
