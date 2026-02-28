// Package ixbrl generates Swedish K2 annual reports in iXBRL format.
//
// The Generate function takes a model.AnnualReport and produces a complete,
// self-contained .xhtml file that is both human-readable (CSS styled) and
// machine-readable (XBRL tagged).
package ixbrl

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/redofri/redofri/pkg/model"
)

// TaxonomyVersion is the K2 taxonomy version we target.
const TaxonomyVersion = "2024-09-12"

// Generate writes a complete iXBRL document for the given annual report.
func Generate(w io.Writer, r *model.AnnualReport) error {
	g := &generator{
		w:      w,
		report: r,
		indent: 0,
	}
	return g.generate()
}

type generator struct {
	w      io.Writer
	report *model.AnnualReport
	indent int
	err    error // sticky error
}

// write outputs a string, tracking errors.
func (g *generator) write(s string) {
	if g.err != nil {
		return
	}
	_, g.err = io.WriteString(g.w, s)
}

// writef outputs a formatted string.
func (g *generator) writef(format string, args ...any) {
	if g.err != nil {
		return
	}
	_, g.err = fmt.Fprintf(g.w, format, args...)
}

// line outputs an indented line.
func (g *generator) line(s string) {
	g.write(strings.Repeat("\t", g.indent))
	g.write(s)
	g.write("\n")
}

// linef outputs an indented formatted line.
func (g *generator) linef(format string, args ...any) {
	g.write(strings.Repeat("\t", g.indent))
	g.writef(format, args...)
	g.write("\n")
}

// raw outputs a string without indentation.
func (g *generator) raw(s string) {
	g.write(s)
}

// in increases indentation.
func (g *generator) in() { g.indent++ }

// out decreases indentation.
func (g *generator) out() {
	if g.indent > 0 {
		g.indent--
	}
}

func (g *generator) generate() error {
	r := g.report

	g.writeXMLDeclaration()
	g.writeHTMLOpen(r)
	g.writeHead(r)
	g.writeBodyOpen()
	g.writeIXHeader(r)
	g.writeWrapper(r)
	g.writeBodyClose()
	g.writeHTMLClose()

	return g.err
}

func (g *generator) writeXMLDeclaration() {
	g.line(`<?xml version="1.0" encoding="UTF-8"?>`)
}

func (g *generator) writeHTMLOpen(r *model.AnnualReport) {
	g.line(`<html xmlns="http://www.w3.org/1999/xhtml"`)
	g.in()
	g.in()
	g.line(`xmlns:iso4217="http://www.xbrl.org/2003/iso4217"`)
	g.line(`xmlns:ixt="http://www.xbrl.org/inlineXBRL/transformation/2010-04-20"`)
	g.line(`xmlns:xlink="http://www.w3.org/1999/xlink"`)
	g.line(`xmlns:link="http://www.xbrl.org/2003/linkbase"`)
	g.line(`xmlns:xbrli="http://www.xbrl.org/2003/instance"`)
	g.line(`xmlns:ix="http://www.xbrl.org/2013/inlineXBRL"`)
	g.line(`xmlns:se-gen-base="http://www.taxonomier.se/se/fr/gen-base/2021-10-31"`)
	g.line(`xmlns:se-cd-base="http://www.taxonomier.se/se/fr/cd-base/2021-10-31"`)
	g.line(`xmlns:se-bol-base="http://www.bolagsverket.se/se/fr/comp-base/2017-09-30"`)
	g.line(`xmlns:se-k2-type="http://www.taxonomier.se/se/fr/k2/datatype">`)
	g.out()
	g.out()
}

func (g *generator) writeHTMLClose() {
	g.line(`</html>`)
}

func (g *generator) writeHead(r *model.AnnualReport) {
	g.in()
	g.line(`<head>`)
	g.in()
	g.linef(`<title> %s %s - Årsredovisning</title>`, esc(r.Company.OrgNr), esc(r.Company.Name))
	g.linef(`<meta name="programvara" content="%s"/>`, esc(r.Meta.Software))
	g.linef(`<meta name="programversion" content="%s"/>`, esc(r.Meta.SoftwareVersion))
	g.line(`<style type="text/css">`)
	g.writeCSS()
	g.line(`</style>`)
	g.out()
	g.line(`</head>`)
	g.out()
}

func (g *generator) writeBodyOpen() {
	g.in()
	g.line(`<body>`)
}

func (g *generator) writeBodyClose() {
	g.in()
	g.line(`</body>`)
	g.out()
}

func (g *generator) writeWrapper(r *model.AnnualReport) {
	g.in()
	g.in()
	g.line(`<div id="wrapper">`)
	g.in()

	g.writeCoverPage(r)
	g.writeManagementReport(r)
	g.writeIncomeStatement(r)
	g.writeBalanceSheetAssets(r)
	g.writeBalanceSheetEquityLiabilities(r)
	g.writeNotes(r)

	g.out()
	g.line(`</div>`)
	g.out()
	g.out()
}

// schemaURL returns the entry point schema URL for the given variant.
func schemaURL(variant string) string {
	return fmt.Sprintf("http://xbrl.taxonomier.se/se/fr/gaap/k2-all/ab/%s/%s/se-k2-ab-%s-%s.xsd",
		variant, TaxonomyVersion, variant, TaxonomyVersion)
}

// certSchemaURL returns the fastställelseintyg schema URL.
func certSchemaURL() string {
	return "http://xbrl.taxonomier.se/se/fr/gaap/k2/rcoa/2020-12-01/se-k2-rcoa-2020-12-01.xsd"
}

// fiscalYearLabel returns a label like "2016" or "2015/16" for display.
func fiscalYearLabel(startDate, endDate string) string {
	// Parse the end date to get the year
	t, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return endDate[:4]
	}
	return fmt.Sprintf("%d", t.Year())
}

// prevYearDates computes the previous fiscal year dates.
// Assumes same-length fiscal year shifted back one year.
func prevYearDates(startDate, endDate string) (string, string) {
	s, err1 := time.Parse("2006-01-02", startDate)
	e, err2 := time.Parse("2006-01-02", endDate)
	if err1 != nil || err2 != nil {
		return "", ""
	}
	ps := s.AddDate(-1, 0, 0)
	pe := e.AddDate(-1, 0, 0)
	return ps.Format("2006-01-02"), pe.Format("2006-01-02")
}

// esc escapes a string for use in XML/HTML content.
func esc(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return s
}
