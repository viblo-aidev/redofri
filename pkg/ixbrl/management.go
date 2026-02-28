package ixbrl

import (
	"fmt"
	"strings"

	"github.com/redofri/redofri/pkg/model"
)

// writeManagementReport writes the förvaltningsberättelse (pages 2-3).
func (g *generator) writeManagementReport(r *model.AnnualReport) {
	mr := &r.ManagementReport
	totalPages := g.computeTotalPages(r)

	// Page 2: verksamhet, flerårsöversikt, equity changes, resultatdisposition (part 1)
	g.line(`<div class="ar-page wide" id="ar3-page-2">`)
	g.in()
	g.pageHeader(r.Company.Name, r.Company.OrgNr, 2, totalPages)

	g.line(`<h2>Förvaltningsberättelse</h2>`)
	g.line(`<h3>Verksamheten</h3>`)
	g.line(`<h4>Allmänt om verksamheten</h4>`)

	// Business description (contains raw HTML)
	g.write(strings.Repeat("\t", g.indent))
	g.nonNumericRaw("se-gen-base:AllmantVerksamheten", "period0", "\n"+mr.BusinessDescription+"\n"+strings.Repeat("\t", g.indent))
	g.write("\n")

	// Significant events
	g.line(`<h4 class="join">Väsentliga händelser under räkenskapsåret</h4>`)
	g.line(`<p>`)
	g.in()
	g.write(strings.Repeat("\t", g.indent))
	g.nonNumeric("se-gen-base:VasentligaHandelserRakenskapsaret", "period0", mr.SignificantEvents)
	g.write("\n")
	g.out()
	g.line(`</p>`)

	// Multi-year overview
	g.writeMultiYearOverview(r)

	// Equity changes
	g.writeEquityChanges(r)

	// Profit disposition - part 1 (available funds)
	g.writeProfitDispositionPart1(r)

	g.out()
	g.line(`</div>`)

	// Page 3: profit disposition part 2 + board statement
	g.line(`<div class="ar-page" id="ar3-page-3">`)
	g.in()
	g.pageHeader(r.Company.Name, r.Company.OrgNr, 3, totalPages)

	g.writeProfitDispositionPart2(r)

	// Board statement on dividend (if any)
	if mr.BoardDividendStatement != "" {
		g.line(`<h3>Styrelsens yttrande över den föreslagna vinstutdelningen</h3>`)
		g.write(strings.Repeat("\t", g.indent))
		g.nonNumericRaw("se-gen-base:StyrelsensYttrandeVinstutdelning", "balans0",
			"\n"+mr.BoardDividendStatement+"\n"+strings.Repeat("\t", g.indent))
		g.write("\n")
	}

	g.out()
	g.line(`</div>`)
}

// writeMultiYearOverview writes the flerårsöversikt table.
func (g *generator) writeMultiYearOverview(r *model.AnnualReport) {
	myo := &r.ManagementReport.MultiYearOverview
	if len(myo.Years) == 0 {
		return
	}

	numCols := len(myo.Years) + 1 // +1 for label column
	g.line(`<h3>Flerårsöversikt</h3>`)
	g.linef(`<table class="ar-overview ar-financial col-%d">`, numCols)
	g.in()

	// Colgroup
	g.line(`<colgroup>`)
	g.in()
	g.line(`<col />`)
	g.linef(`<col class="tkr" span="%d" />`, len(myo.Years))
	g.out()
	g.line(`</colgroup>`)

	// Header row with years
	g.line(`<thead>`)
	g.in()
	g.line(`<tr>`)
	g.in()
	g.line(`<th />`)
	for _, y := range myo.Years {
		g.linef(`<th scope="col">%s</th>`, esc(y.Year))
	}
	g.out()
	g.line(`</tr>`)
	g.out()
	g.line(`</thead>`)

	g.line(`<tbody>`)
	g.in()

	// Nettoomsättning row
	g.line(`<tr>`)
	g.in()
	g.linef(`<td>Nettoomsättning, <abbr>tkr</abbr></td>`)
	for i, y := range myo.Years {
		g.write(strings.Repeat("\t", g.indent))
		g.write("<td>")
		if y.NetSales != nil {
			ctx := contextRefForOverviewYear(i, false)
			decStr := "-3"
			if *y.NetSales == 0 {
				decStr = "INF"
			}
			g.nonFraction("se-gen-base:Nettoomsattning", ctx, "SEK", *y.NetSales,
				withDecimals(decStr), withScale("3"))
		}
		g.write("</td>\n")
	}
	g.out()
	g.line(`</tr>`)

	// Resultat efter finansiella poster row
	g.line(`<tr>`)
	g.in()
	g.linef(`<td>Resultat efter finansiella poster, <abbr>tkr</abbr></td>`)
	for i, y := range myo.Years {
		g.write(strings.Repeat("\t", g.indent))
		g.write("<td>")
		if y.ResultAfterFinancialItems != nil {
			ctx := contextRefForOverviewYear(i, false)
			decStr := "-3"
			if *y.ResultAfterFinancialItems == 0 {
				decStr = "INF"
			}
			g.nonFraction("se-gen-base:ResultatEfterFinansiellaPoster", ctx, "SEK", *y.ResultAfterFinancialItems,
				withDecimals(decStr), withScale("3"))
		}
		g.write("</td>\n")
	}
	g.out()
	g.line(`</tr>`)

	// Soliditet row
	g.line(`<tr>`)
	g.in()
	g.line(`<td>Soliditet, %</td>`)
	for i, y := range myo.Years {
		g.write(strings.Repeat("\t", g.indent))
		g.write("<td>")
		if y.Solidity != nil {
			ctx := contextRefForOverviewYear(i, true)
			g.writef(`<ix:nonFraction contextRef="%s" name="se-gen-base:Soliditet" unitRef="procent" format="ixt:numcomma" scale="-2" decimals="INF">%s</ix:nonFraction>`,
				ctx, esc(*y.Solidity))
		}
		g.write("</td>\n")
	}
	g.out()
	g.line(`</tr>`)

	g.out()
	g.line(`</tbody>`)
	g.out()
	g.line(`</table>`)

	// Comment
	if myo.Comment != "" {
		g.line(`<p>`)
		g.in()
		g.write(strings.Repeat("\t", g.indent))
		g.nonNumeric("se-gen-base:KommentarFlerarsoversikt", "period0", myo.Comment)
		g.write("\n")
		g.out()
		g.line(`</p>`)
	}
}

// writeEquityChanges writes the förändringar i eget kapital table.
func (g *generator) writeEquityChanges(r *model.AnnualReport) {
	ec := &r.ManagementReport.EquityChanges

	g.line(`<h3>Förändringar i eget kapital</h3>`)
	g.line(`<table class="ar-equity ar-financial col-5">`)
	g.in()

	g.line(`<colgroup>`)
	g.in()
	g.line(`<col />`)
	g.line(`<col class="kr" span="5" />`)
	g.out()
	g.line(`</colgroup>`)

	// Header
	g.line(`<thead>`)
	g.in()
	g.line(`<tr>`)
	g.in()
	g.line(`<th />`)
	g.line(`<th scope="col">Aktiekapital</th>`)
	g.line(`<th scope="col">Reservfond</th>`)
	g.line(`<th scope="col">Balanserat resultat</th>`)
	g.line(`<th scope="col">Årets resultat</th>`)
	g.line(`<th scope="col">Totalt</th>`)
	g.out()
	g.line(`</tr>`)
	g.out()
	g.line(`</thead>`)

	g.line(`<tbody>`)
	g.in()

	// Opening balances
	g.line(`<tr>`)
	g.in()
	g.line(`<td>Belopp vid årets ingång</td>`)
	g.writeEquityCell("se-gen-base:Aktiekapital", "balans1", ec.OpeningShareCapital, "")
	g.writeEquityCell("se-gen-base:Reservfond", "balans1", ec.OpeningReserveFund, "")
	g.writeEquityCell("se-gen-base:BalanseratResultat", "balans1", ec.OpeningRetainedEarnings, "")
	g.writeEquityCell("se-gen-base:AretsResultatEgetKapital", "balans1", ec.OpeningNetIncome, "")
	g.writeEquityCell("se-gen-base:ForandringEgetKapitalTotalt", "balans1", ec.OpeningTotal, "")
	g.out()
	g.line(`</tr>`)

	// Resultatdisposition header row
	g.line(`<tr>`)
	g.in()
	g.line(`<td>Resultatdisposition enligt årsstämman</td>`)
	g.line(`<td colspan="5" />`)
	g.out()
	g.line(`</tr>`)

	// Dividend row (if applicable)
	if ec.DividendNetIncome != nil || ec.DividendTotal != nil {
		g.line(`<tr>`)
		g.in()
		g.line(`<td>– Utdelning</td>`)
		g.line(`<td>–</td>`)
		g.line(`<td>–</td>`)
		g.line(`<td>–</td>`)
		// Dividend net income (displayed with "-" prefix)
		g.write(strings.Repeat("\t", g.indent))
		g.write("<td>")
		if ec.DividendNetIncome != nil {
			g.write("-")
			g.nonFraction("se-gen-base:ForandringEgetKapitalAretsResultatUtdelning", "period0", "SEK", *ec.DividendNetIncome)
		}
		g.write("</td>\n")
		// Dividend total
		g.write(strings.Repeat("\t", g.indent))
		g.write("<td>")
		if ec.DividendTotal != nil {
			g.write("-")
			g.nonFraction("se-gen-base:ForandringEgetKapitalTotaltUtdelning", "period0", "SEK", *ec.DividendTotal)
		}
		g.write("</td>\n")
		g.out()
		g.line(`</tr>`)
	}

	// Year's result row
	g.line(`<tr>`)
	g.in()
	g.line(`<td>Årets resultat</td>`)
	g.line(`<td>–</td>`)
	g.line(`<td>–</td>`)
	g.line(`<td>–</td>`)
	g.writeEquityCell("se-gen-base:ForandringEgetKapitalAretsResultatAretsResultat", "period0", ec.YearResultNetIncome, "sum")
	g.writeEquityCell("se-gen-base:ForandringEgetKapitalTotaltAretsResultat", "period0", ec.YearResultTotal, "sum")
	g.out()
	g.line(`</tr>`)

	// Closing balances
	g.line(`<tr>`)
	g.in()
	g.line(`<td>Belopp vid årets utgång</td>`)
	g.writeEquityCell("se-gen-base:Aktiekapital", "balans0", ec.ClosingShareCapital, "total")
	g.writeEquityCell("se-gen-base:Reservfond", "balans0", ec.ClosingReserveFund, "total")
	g.writeEquityCell("se-gen-base:BalanseratResultat", "balans0", ec.ClosingRetainedEarnings, "total")
	g.writeEquityCell("se-gen-base:AretsResultatEgetKapital", "balans0", ec.ClosingNetIncome, "total")
	g.writeEquityCell("se-gen-base:ForandringEgetKapitalTotalt", "balans0", ec.ClosingTotal, "total")
	g.out()
	g.line(`</tr>`)

	g.out()
	g.line(`</tbody>`)
	g.out()
	g.line(`</table>`)
}

// writeEquityCell writes a single cell in the equity changes table.
func (g *generator) writeEquityCell(concept, contextRef string, value *int64, wrapClass string) {
	g.write(strings.Repeat("\t", g.indent))
	g.write("<td>")
	if value != nil {
		var opts []nfOpt
		if wrapClass != "" {
			opts = append(opts, withWrapClass(wrapClass))
		}
		g.nonFraction(concept, contextRef, "SEK", *value, opts...)
	}
	g.write("</td>\n")
}

// writeProfitDispositionPart1 writes the available funds table (on page 2).
func (g *generator) writeProfitDispositionPart1(r *model.AnnualReport) {
	pd := &r.ManagementReport.ProfitDisposition

	g.line(`<h3>Resultatdisposition</h3>`)
	g.line(`<p class="ar-disp">Till årsstämmans förfogande står följande vinstmedel:</p>`)
	g.line(`<table class="ar-disp ar-financial">`)
	g.in()
	g.line(`<colgroup>`)
	g.in()
	g.line(`<col />`)
	g.line(`<col class="kr" />`)
	g.out()
	g.line(`</colgroup>`)
	g.line(`<tbody>`)
	g.in()

	// Balanserat resultat
	g.line(`<tr>`)
	g.in()
	g.line(`<td>Balanserat resultat</td>`)
	g.writeDispCell("se-gen-base:BalanseratResultat", "balans0", pd.RetainedEarnings, "")
	g.out()
	g.line(`</tr>`)

	// Årets resultat
	g.line(`<tr>`)
	g.in()
	g.line(`<td>Årets resultat</td>`)
	g.writeDispCell("se-gen-base:AretsResultatEgetKapital", "balans0", pd.NetIncome, "sum")
	g.out()
	g.line(`</tr>`)

	// Totalt
	g.line(`<tr>`)
	g.in()
	g.line(`<td>Totalt</td>`)
	g.writeDispCell("se-gen-base:MedelDisponera", "balans0", pd.TotalAvailable, "total")
	g.out()
	g.line(`</tr>`)

	g.out()
	g.line(`</tbody>`)
	g.out()
	g.line(`</table>`)
}

// writeProfitDispositionPart2 writes the proposed disposition table (on page 3).
func (g *generator) writeProfitDispositionPart2(r *model.AnnualReport) {
	pd := &r.ManagementReport.ProfitDisposition

	g.line(`<p>Styrelsen och verkställande direktören föreslår att vinstmedlen disponeras enligt följande</p>`)
	g.line(`<table class="ar-disp ar-financial">`)
	g.in()
	g.line(`<colgroup>`)
	g.in()
	g.line(`<col />`)
	g.line(`<col class="kr" />`)
	g.out()
	g.line(`</colgroup>`)
	g.line(`<tbody>`)
	g.in()

	// Utdelning
	if pd.Dividend != nil {
		g.line(`<tr>`)
		g.in()
		g.line(`<td>Utdelning till ägarna</td>`)
		g.writeDispCell("se-gen-base:ForslagDispositionUtdelning", "balans0", pd.Dividend, "")
		g.out()
		g.line(`</tr>`)
	}

	// Balanseras i ny räkning
	g.line(`<tr>`)
	g.in()
	g.line(`<td>Balanseras i ny räkning</td>`)
	g.writeDispCell("se-gen-base:ForslagDispositionBalanserasINyRakning", "balans0", pd.CarriedForward, "sum")
	g.out()
	g.line(`</tr>`)

	// Totalt
	g.line(`<tr>`)
	g.in()
	g.line(`<td>Totalt</td>`)
	g.writeDispCell("se-gen-base:ForslagDisposition", "balans0", pd.TotalDisposition, "total")
	g.out()
	g.line(`</tr>`)

	g.out()
	g.line(`</tbody>`)
	g.out()
	g.line(`</table>`)
}

// writeDispCell writes a disposition table cell.
func (g *generator) writeDispCell(concept, contextRef string, value *int64, wrapClass string) {
	g.write(strings.Repeat("\t", g.indent))
	g.write("<td>")
	if value != nil {
		var opts []nfOpt
		if wrapClass != "" {
			opts = append(opts, withWrapClass(wrapClass))
		}
		g.nonFraction(concept, contextRef, "SEK", *value, opts...)
	}
	g.write("</td>\n")
}

// solidityDisplay formats a solidity string for display.
// The model stores it as "33.7", "100", etc. We convert "." to "," for Swedish format.
func solidityDisplay(s string) string {
	return strings.ReplaceAll(s, ".", ",")
}

// unused but reserved for future use
var _ = fmt.Sprint
