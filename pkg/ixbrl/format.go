package ixbrl

import (
	"fmt"
	"strings"
)

// formatAmount formats an int64 (whole kronor) as "1 234 567" using space as
// thousands separator and comma as decimal separator (ixt:numspacecomma).
// For zero, returns "0".
func formatAmount(v int64) string {
	if v == 0 {
		return "0"
	}
	negative := v < 0
	if negative {
		v = -v
	}
	s := fmt.Sprintf("%d", v)
	// Insert space every 3 digits from the right
	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	parts = append([]string{s}, parts...)
	result := strings.Join(parts, " ")
	if negative {
		result = "-" + result
	}
	return result
}

// formatAmountTkr formats an int64 (whole kronor) in tkr (thousands).
// 2650000 → "2 650", 0 → "0"
func formatAmountTkr(v int64) string {
	tkr := v / 1000
	return formatAmount(tkr)
}

// nonFraction writes an ix:nonFraction element for a monetary amount.
// name: XBRL concept (e.g. "se-gen-base:Nettoomsattning")
// contextRef: e.g. "period0", "balans0"
// unitRef: e.g. "SEK"
// value: the int64 amount in whole kronor
// opts: optional attributes (sign, scale, decimals override, wrapClass)
func (g *generator) nonFraction(name, contextRef, unitRef string, value int64, opts ...nfOpt) {
	o := nfOptions{
		decimals: "INF",
		scale:    "0",
		format:   "ixt:numspacecomma",
	}
	for _, fn := range opts {
		fn(&o)
	}

	displayValue := formatAmount(value)
	if o.scale == "3" {
		// tkr display
		displayValue = formatAmountTkr(value)
	}

	attrs := fmt.Sprintf(`contextRef="%s" name="%s" unitRef="%s" decimals="%s" scale="%s"`,
		contextRef, name, unitRef, o.decimals, o.scale)
	if o.format != "" {
		attrs += fmt.Sprintf(` format="%s"`, o.format)
	}
	if o.sign != "" {
		attrs += fmt.Sprintf(` sign="%s"`, o.sign)
	}
	if o.tupleRef != "" {
		attrs += fmt.Sprintf(` tupleRef="%s"`, o.tupleRef)
	}
	if o.order != "" {
		attrs += fmt.Sprintf(` order="%s"`, o.order)
	}

	tag := fmt.Sprintf(`<ix:nonFraction %s>%s</ix:nonFraction>`, attrs, displayValue)

	// Handle negative prefix display (expenses show "-" outside the tag)
	if o.negPrefix {
		tag = "-" + tag
	}

	// Wrap in span class if specified
	if o.wrapClass != "" {
		tag = fmt.Sprintf(`<span class="%s">%s</span>`, o.wrapClass, tag)
	}

	g.write(tag)
}

// nonFractionLine writes a nonFraction on its own indented line.
func (g *generator) nonFractionLine(name, contextRef, unitRef string, value int64, opts ...nfOpt) {
	g.write(strings.Repeat("\t", g.indent))
	g.nonFraction(name, contextRef, unitRef, value, opts...)
	g.write("\n")
}

// nonNumeric writes an ix:nonNumeric element.
func (g *generator) nonNumeric(name, contextRef, value string, opts ...nnOpt) {
	o := nnOptions{}
	for _, fn := range opts {
		fn(&o)
	}

	attrs := fmt.Sprintf(`name="%s" contextRef="%s"`, name, contextRef)
	if o.id != "" {
		attrs += fmt.Sprintf(` id="%s"`, o.id)
	}
	if o.continuedAt != "" {
		attrs += fmt.Sprintf(` continuedAt="%s"`, o.continuedAt)
	}
	if o.order != "" {
		attrs += fmt.Sprintf(` order="%s"`, o.order)
	}
	if o.tupleRef != "" {
		attrs += fmt.Sprintf(` tupleRef="%s"`, o.tupleRef)
	}

	g.writef(`<ix:nonNumeric %s>%s</ix:nonNumeric>`, attrs, esc(value))
}

// nonNumericRaw writes an ix:nonNumeric element with raw (pre-escaped) content.
func (g *generator) nonNumericRaw(name, contextRef, rawContent string, opts ...nnOpt) {
	o := nnOptions{}
	for _, fn := range opts {
		fn(&o)
	}

	attrs := fmt.Sprintf(`name="%s" contextRef="%s"`, name, contextRef)
	if o.id != "" {
		attrs += fmt.Sprintf(` id="%s"`, o.id)
	}
	if o.continuedAt != "" {
		attrs += fmt.Sprintf(` continuedAt="%s"`, o.continuedAt)
	}
	if o.order != "" {
		attrs += fmt.Sprintf(` order="%s"`, o.order)
	}
	if o.tupleRef != "" {
		attrs += fmt.Sprintf(` tupleRef="%s"`, o.tupleRef)
	}

	g.writef(`<ix:nonNumeric %s>%s</ix:nonNumeric>`, attrs, rawContent)
}

// nfOptions holds optional attributes for nonFraction.
type nfOptions struct {
	decimals  string // default "INF"
	scale     string // default "0"
	format    string // default "ixt:numspacecomma"
	sign      string // e.g. "-" for sign inversion
	negPrefix bool   // show "-" before the tag in display
	wrapClass string // wrap in <span class="...">
	tupleRef  string // tupleRef attribute for tuple membership
	order     string // order attribute within tuple
}

type nfOpt func(*nfOptions)

func withDecimals(d string) nfOpt      { return func(o *nfOptions) { o.decimals = d } }
func withScale(s string) nfOpt         { return func(o *nfOptions) { o.scale = s } }
func withFormat(f string) nfOpt        { return func(o *nfOptions) { o.format = f } }
func withSign(s string) nfOpt          { return func(o *nfOptions) { o.sign = s } }
func withNegPrefix() nfOpt             { return func(o *nfOptions) { o.negPrefix = true } }
func withWrapClass(class string) nfOpt { return func(o *nfOptions) { o.wrapClass = class } }

// nnOptions holds optional attributes for nonNumeric.
type nnOptions struct {
	id          string
	continuedAt string
	order       string
	tupleRef    string
}

type nnOpt func(*nnOptions)

func withID(id string) nnOpt         { return func(o *nnOptions) { o.id = id } }
func withContinuedAt(c string) nnOpt { return func(o *nnOptions) { o.continuedAt = c } }
func withOrder(ord string) nnOpt     { return func(o *nnOptions) { o.order = ord } }
func withTupleRef(ref string) nnOpt  { return func(o *nnOptions) { o.tupleRef = ref } }

// writeYearComparisonRow writes a standard financial table row with two year columns.
// label: display label (left column)
// noteRef: optional note number (0 = no note), shown as link in Not column
// concept: XBRL concept name
// currentCtx, prevCtx: context refs for current and previous year
// current, previous: *int64 values
// isExpense: if true, display with "-" prefix and positive XBRL value
// isSubTotal/isTotal: controls CSS class wrapping
func (g *generator) writeYearComparisonRow(label string, noteRef int, concept, currentCtx, prevCtx, unitRef string,
	current, previous *int64, isExpense, isSubTotal, isTotal bool) {
	if current == nil && previous == nil {
		return // skip empty rows
	}

	// Determine row CSS
	tdClass := ""
	if isSubTotal {
		tdClass = ` class="sub-sum"`
	} else if isTotal {
		tdClass = ` class="sum total"`
	}

	g.linef(`<tr>`)
	g.in()
	if tdClass != "" {
		g.linef(`<td%s>%s</td>`, tdClass, label)
	} else {
		g.linef(`<td>%s</td>`, label)
	}

	// Note column
	if noteRef > 0 {
		g.linef(`<td><a href="#note-%d">%d</a></td>`, noteRef, noteRef)
	} else {
		g.line(`<td />`)
	}

	// Current year
	g.write(strings.Repeat("\t", g.indent))
	g.write("<td>")
	if current != nil {
		var opts []nfOpt
		if isExpense {
			opts = append(opts, withNegPrefix())
		}
		// Last item in a group gets "sum" class, total gets "total"
		wrapClass := ""
		if isSubTotal {
			wrapClass = "sum"
		} else if isTotal {
			wrapClass = "total"
		}
		if wrapClass != "" {
			opts = append(opts, withWrapClass(wrapClass))
		}
		g.nonFraction(concept, currentCtx, unitRef, *current, opts...)
	}
	g.write("</td>\n")

	// Previous year
	g.write(strings.Repeat("\t", g.indent))
	g.write("<td>")
	if previous != nil {
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
		g.nonFraction(concept, prevCtx, unitRef, *previous, opts...)
	}
	g.write("</td>\n")

	g.out()
	g.line(`</tr>`)
}

// writeBalanceRow writes a balance sheet row with balans0/balans1 contexts.
func (g *generator) writeBalanceRow(label string, noteRef int, noteRefs []int, concept string,
	yc YearComparisonVal, isSubTotal, isTotal, isLastInGroup bool) {
	if yc.Current == nil && yc.Previous == nil {
		return
	}

	tdClass := ""
	if isSubTotal {
		tdClass = ` class="sub-sum"`
	} else if isTotal {
		tdClass = ` class="sum total"`
	} else if isLastInGroup {
		// no special class but wrapped in sum span
	}

	g.line(`<tr>`)
	g.in()
	if tdClass != "" {
		g.linef(`<td%s>%s</td>`, tdClass, label)
	} else {
		g.linef(`<td>%s</td>`, label)
	}

	// Note column
	if noteRef > 0 {
		g.linef(`<td><a href="#note-%d">%d</a></td>`, noteRef, noteRef)
	} else if len(noteRefs) > 0 {
		g.write(strings.Repeat("\t", g.indent))
		g.write("<td>")
		for i, nr := range noteRefs {
			if i > 0 {
				g.write(", ")
			}
			g.writef(`<a href="#note-%d">%d</a>`, nr, nr)
		}
		g.write("</td>\n")
	} else {
		g.line(`<td />`)
	}

	wrapClass := ""
	if isSubTotal || isLastInGroup {
		wrapClass = "sum"
	} else if isTotal {
		wrapClass = "total"
	}

	// Current year (balans0)
	g.write(strings.Repeat("\t", g.indent))
	g.write("<td>")
	if yc.Current != nil {
		var opts []nfOpt
		if wrapClass != "" {
			opts = append(opts, withWrapClass(wrapClass))
		}
		g.nonFraction(concept, "balans0", "SEK", *yc.Current, opts...)
	}
	g.write("</td>\n")

	// Previous year (balans1)
	g.write(strings.Repeat("\t", g.indent))
	g.write("<td>")
	if yc.Previous != nil {
		var opts []nfOpt
		if wrapClass != "" {
			opts = append(opts, withWrapClass(wrapClass))
		}
		g.nonFraction(concept, "balans1", "SEK", *yc.Previous, opts...)
	}
	g.write("</td>\n")

	g.out()
	g.line(`</tr>`)
}

// YearComparisonVal is a simple holder for current/previous *int64 to avoid
// importing model in format helpers or to bridge from model.YearComparison.
type YearComparisonVal struct {
	Current  *int64
	Previous *int64
}

// pageHeader writes the standard page header div.
func (g *generator) pageHeader(companyName, orgNr string, pageNum, totalPages int) {
	g.linef(`<div class="ar-page-hdr">`)
	g.in()
	g.linef(`<span>%s<br />%s</span>    %d (%d)  `, esc(companyName), esc(orgNr), pageNum, totalPages)
	g.out()
	g.line(`</div>`)
}
