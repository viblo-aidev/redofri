package ixbrl

import (
	"fmt"

	"github.com/redofri/redofri/pkg/model"
)

// writeSignatures writes the underskrifter section at the bottom of the last note page.
// This is called from writeNotes() inside the last page div.
func (g *generator) writeSignatures(r *model.AnnualReport) {
	sigs := &r.Signatures

	g.line(`<div class="ar-signature-2">`)
	g.in()

	// City and date
	g.line(`<p>`)
	g.in()
	g.write(indentStr(g.indent))
	g.nonNumeric("se-gen-base:UndertecknandeArsredovisningOrt", "period0", sigs.City)
	g.write("\n")
	g.write(indentStr(g.indent))
	g.nonNumeric("se-gen-base:UndertecknandeArsredovisningDatum", "period0", sigs.Date)
	g.write("\n")
	g.out()
	g.line(`</p>`)

	// Tuple declarations for all signatories
	for i := range sigs.Signatories {
		g.linef(`<ix:tuple name="se-gen-base:UnderskriftArsredovisningForetradareTuple" tupleID="UnderskriftArsredovisningForetradareTuple%d" />`, i+1)
	}

	// Signatory blocks
	for i, sig := range sigs.Signatories {
		tupleRef := fmt.Sprintf("UnderskriftArsredovisningForetradareTuple%d", i+1)

		g.line(`<div class="name">`)
		g.in()

		// Display name in italics
		g.linef(`<i>%s %s</i>`, esc(sig.FirstName), esc(sig.LastName))
		g.line(`<br />`)

		// XBRL tagged first name
		g.write(indentStr(g.indent))
		g.nonNumeric("se-gen-base:UnderskriftArsredovisningForetradareTilltalsnamn", "period0", sig.FirstName,
			withOrder("1.0"), withTupleRef(tupleRef))
		g.write("\n")

		// XBRL tagged last name
		g.write(indentStr(g.indent))
		g.nonNumeric("se-gen-base:UnderskriftArsredovisningForetradareEfternamn", "period0", sig.LastName,
			withOrder("2.0"), withTupleRef(tupleRef))
		g.write("\n")

		// Optional role
		if sig.Role != "" {
			g.line(`<br />`)
			g.write(indentStr(g.indent))
			g.nonNumeric("se-gen-base:UnderskriftArsredovisningForetradareForetradarroll", "period0", sig.Role,
				withOrder("3.0"), withTupleRef(tupleRef))
			g.write("\n")
		}

		g.out()
		g.line(`</div>`)
	}

	g.out()
	g.line(`</div>`)
}
