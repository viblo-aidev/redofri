package ixbrl

import (
	"fmt"
	"time"

	"github.com/redofri/redofri/pkg/model"
)

// writeIXHeader writes the ix:header block (hidden div) containing contexts,
// units, schema references, and hidden metadata facts.
func (g *generator) writeIXHeader(r *model.AnnualReport) {
	g.in()
	g.in()
	g.line(`<div style="display:none">`)
	g.in()
	g.line(`<ix:header>`)
	g.in()

	g.writeHiddenFacts(r)
	g.writeSchemaRefs(r)
	g.writeResources(r)

	g.out()
	g.line(`</ix:header>`)
	g.out()
	g.line(`</div>`)
	g.out()
	g.out()
}

// writeHiddenFacts writes the ix:hidden block with metadata facts.
func (g *generator) writeHiddenFacts(r *model.AnnualReport) {
	g.line(`<ix:hidden>`)
	g.in()
	g.linef(`<ix:nonNumeric name="se-cd-base:Sprak" contextRef="period0">%s</ix:nonNumeric>`, esc(r.Meta.Language))
	g.linef(`<ix:nonNumeric name="se-cd-base:Land" contextRef="period0">%s</ix:nonNumeric>`, esc(r.Meta.Country))
	g.linef(`<ix:nonNumeric name="se-cd-base:Redovisningsvaluta" contextRef="period0">%s</ix:nonNumeric>`, esc(r.Meta.Currency))
	g.linef(`<ix:nonNumeric name="se-cd-base:Beloppsformat" contextRef="period0">%s</ix:nonNumeric>`, esc(r.Meta.AmountFormat))
	g.linef(`<ix:nonNumeric name="se-cd-base:RakenskapsarForstaDag" contextRef="period0">%s</ix:nonNumeric>`, esc(r.FiscalYear.StartDate))
	g.linef(`<ix:nonNumeric name="se-cd-base:RakenskapsarSistaDag" contextRef="period0">%s</ix:nonNumeric>`, esc(r.FiscalYear.EndDate))
	g.out()
	g.line(`</ix:hidden>`)
}

// writeSchemaRefs writes the ix:references block with schema references.
func (g *generator) writeSchemaRefs(r *model.AnnualReport) {
	g.line(`<ix:references>`)
	g.in()
	// Entry point schema (e.g. risbs)
	g.linef(`<link:schemaRef xlink:type="simple" xlink:href="%s" />`, schemaURL(r.Meta.EntryPoint))
	// Fastställelseintyg schema
	g.linef(`<link:schemaRef xlink:type="simple" xlink:href="%s"/>`, certSchemaURL())
	g.out()
	g.line(`</ix:references>`)
}

// writeResources writes the ix:resources block with contexts and units.
func (g *generator) writeResources(r *model.AnnualReport) {
	g.line(`<ix:resources>`)
	g.in()

	orgNr := r.Company.OrgNr
	start := r.FiscalYear.StartDate
	end := r.FiscalYear.EndDate

	// Compute dates for all periods
	prevStart, prevEnd := prevYearDates(start, end)

	// period0: current fiscal year (duration)
	g.writeDurationContext("period0", orgNr, start, end)
	// balans0: current year-end (instant)
	g.writeInstantContext("balans0", orgNr, end)
	// balans1: previous year-end (instant)
	g.writeInstantContext("balans1", orgNr, prevEnd)
	// period1: previous fiscal year (duration)
	g.writeDurationContext("period1", orgNr, prevStart, prevEnd)

	// Multi-year overview may need period2/period3 and balans2/balans3
	if len(r.ManagementReport.MultiYearOverview.Years) > 2 {
		s2, e2 := yearNDates(start, end, 2)
		g.writeDurationContext("period2", orgNr, s2, e2)
		g.writeInstantContext("balans2", orgNr, e2)
	}
	if len(r.ManagementReport.MultiYearOverview.Years) > 3 {
		s3, e3 := yearNDates(start, end, 3)
		g.writeDurationContext("period3", orgNr, s3, e3)
		g.writeInstantContext("balans3", orgNr, e3)
	}

	// Units
	g.writeUnit("SEK", "iso4217:SEK")
	g.writeUnit("procent", "xbrli:pure")
	g.writeUnit("antal-anstallda", "se-k2-type:AntalAnstallda")

	g.out()
	g.line(`</ix:resources>`)
}

// writeDurationContext writes a duration context element.
func (g *generator) writeDurationContext(id, orgNr, startDate, endDate string) {
	g.linef(`<xbrli:context id="%s">`, id)
	g.in()
	g.line(`<xbrli:entity>`)
	g.in()
	g.linef(`<xbrli:identifier scheme="http://www.bolagsverket.se">%s</xbrli:identifier>`, esc(orgNr))
	g.out()
	g.line(`</xbrli:entity>`)
	g.line(`<xbrli:period>`)
	g.in()
	g.linef(`<xbrli:startDate>%s</xbrli:startDate>`, startDate)
	g.linef(`<xbrli:endDate>%s</xbrli:endDate>`, endDate)
	g.out()
	g.line(`</xbrli:period>`)
	g.out()
	g.linef(`</xbrli:context>`)
}

// writeInstantContext writes an instant context element.
func (g *generator) writeInstantContext(id, orgNr, instant string) {
	g.linef(`<xbrli:context id="%s">`, id)
	g.in()
	g.line(`<xbrli:entity>`)
	g.in()
	g.linef(`<xbrli:identifier scheme="http://www.bolagsverket.se">%s</xbrli:identifier>`, esc(orgNr))
	g.out()
	g.line(`</xbrli:entity>`)
	g.line(`<xbrli:period>`)
	g.in()
	g.linef(`<xbrli:instant>%s</xbrli:instant>`, instant)
	g.out()
	g.line(`</xbrli:period>`)
	g.out()
	g.linef(`</xbrli:context>`)
}

// writeUnit writes a unit element.
func (g *generator) writeUnit(id, measure string) {
	g.linef(`<xbrli:unit id="%s">`, id)
	g.in()
	g.linef(`<xbrli:measure>%s</xbrli:measure>`, measure)
	g.out()
	g.line(`</xbrli:unit>`)
}

// yearNDates computes the dates for N fiscal years back.
// yearNDates(start, end, 2) gives the period 2 years before.
func yearNDates(startDate, endDate string, n int) (string, string) {
	s, err1 := time.Parse("2006-01-02", startDate)
	e, err2 := time.Parse("2006-01-02", endDate)
	if err1 != nil || err2 != nil {
		return "", ""
	}
	ps := s.AddDate(-n, 0, 0)
	pe := e.AddDate(-n, 0, 0)
	return ps.Format("2006-01-02"), pe.Format("2006-01-02")
}

// contextRefForOverviewYear returns the contextRef (period or balans) for a
// given year index in the multi-year overview.
func contextRefForOverviewYear(idx int, instant bool) string {
	if instant {
		return fmt.Sprintf("balans%d", idx)
	}
	return fmt.Sprintf("period%d", idx)
}
