package ixbrl

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"os"
	"strings"
	"testing"

	"github.com/redofri/redofri/pkg/model"
)

// loadTestReport loads the example test data from testdata/exempel1.json.
func loadTestReport(t *testing.T) *model.AnnualReport {
	t.Helper()
	data, err := os.ReadFile("../../testdata/exempel1.json")
	if err != nil {
		t.Fatalf("reading test data: %v", err)
	}
	var r model.AnnualReport
	if err := json.Unmarshal(data, &r); err != nil {
		t.Fatalf("parsing test data: %v", err)
	}
	return &r
}

// generateOutput generates iXBRL output for the given report.
func generateOutput(t *testing.T, r *model.AnnualReport) string {
	t.Helper()
	var buf bytes.Buffer
	if err := Generate(&buf, r); err != nil {
		t.Fatalf("Generate: %v", err)
	}
	return buf.String()
}

func TestGenerate_ProducesOutput(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)
	if len(output) == 0 {
		t.Fatal("generated output is empty")
	}
	// Minimum size check — a real K2 report is several hundred KB
	if len(output) < 10000 {
		t.Errorf("output suspiciously small: %d bytes", len(output))
	}
}

func TestGenerate_ValidXML(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	// xml.Decoder should be able to parse the output without errors
	decoder := xml.NewDecoder(strings.NewReader(output))
	for {
		_, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			t.Fatalf("XML parse error: %v", err)
		}
	}
}

func TestGenerate_XMLDeclaration(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)
	if !strings.HasPrefix(output, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("missing or incorrect XML declaration")
	}
}

func TestGenerate_HTMLNamespaces(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	namespaces := []string{
		`xmlns="http://www.w3.org/1999/xhtml"`,
		`xmlns:ix="http://www.xbrl.org/2013/inlineXBRL"`,
		`xmlns:xbrli="http://www.xbrl.org/2003/instance"`,
		`xmlns:se-gen-base="http://www.taxonomier.se/se/fr/gen-base/2021-10-31"`,
		`xmlns:se-cd-base="http://www.taxonomier.se/se/fr/cd-base/2021-10-31"`,
		`xmlns:se-bol-base="http://www.bolagsverket.se/se/fr/comp-base/2017-09-30"`,
	}
	for _, ns := range namespaces {
		if !strings.Contains(output, ns) {
			t.Errorf("missing namespace: %s", ns)
		}
	}
}

func TestGenerate_MetaTags(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	if !strings.Contains(output, `<meta name="programvara" content="Redofri"/>`) {
		t.Error("missing programvara meta tag")
	}
	if !strings.Contains(output, `<meta name="programversion" content="0.1.0"/>`) {
		t.Error("missing programversion meta tag")
	}
}

func TestGenerate_HiddenFacts(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	hiddenFacts := []string{
		`se-cd-base:Sprak`,
		`se-cd-base:Land`,
		`se-cd-base:Redovisningsvaluta`,
		`se-cd-base:Beloppsformat`,
		`se-cd-base:RakenskapsarForstaDag`,
		`se-cd-base:RakenskapsarSistaDag`,
	}
	for _, fact := range hiddenFacts {
		if !strings.Contains(output, fact) {
			t.Errorf("missing hidden fact: %s", fact)
		}
	}
}

func TestGenerate_SchemaRefs(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	// Entry point schema
	if !strings.Contains(output, "se-k2-ab-risbs-2024-09-12.xsd") {
		t.Error("missing entry point schema reference")
	}
	// Fastställelseintyg schema
	if !strings.Contains(output, "se-k2-rcoa-2020-12-01.xsd") {
		t.Error("missing certification schema reference")
	}
}

func TestGenerate_Contexts(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	contexts := []string{
		`id="period0"`, `id="period1"`, `id="period2"`, `id="period3"`,
		`id="balans0"`, `id="balans1"`, `id="balans2"`, `id="balans3"`,
	}
	for _, ctx := range contexts {
		if !strings.Contains(output, ctx) {
			t.Errorf("missing context: %s", ctx)
		}
	}

	// Verify org nr in entity
	if !strings.Contains(output, `<xbrli:identifier scheme="http://www.bolagsverket.se">556999-9999</xbrli:identifier>`) {
		t.Error("missing org nr in context entity")
	}
}

func TestGenerate_Units(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	units := map[string]string{
		"SEK":             "iso4217:SEK",
		"procent":         "xbrli:pure",
		"antal-anstallda": "se-k2-type:AntalAnstallda",
	}
	for id, measure := range units {
		if !strings.Contains(output, `id="`+id+`"`) {
			t.Errorf("missing unit: %s", id)
		}
		if !strings.Contains(output, measure) {
			t.Errorf("missing unit measure: %s", measure)
		}
	}
}

func TestGenerate_CoverPage(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	checks := []string{
		`id="ar3-page-1"`,
		`Exempel 1 AB`,
		`556999-9999`,
		`Årsredovisning för räkenskapsåret 2016`,
		`se-cd-base:ForetagetsNamn`,
		`se-cd-base:Organisationsnummer`,
		// Table of contents
		`förvaltningsberättelse`,
		`resultaträkning`,
		`balansräkning`,
		`noter`,
	}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("cover page missing: %s", check)
		}
	}
}

func TestGenerate_Certification(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	checks := []string{
		`id="ar-certification"`,
		`Fastställelseintyg`,
		`se-bol-base:FaststallelseResultatBalansrakning`,
		`se-bol-base:Arsstamma`,
		`se-bol-base:ArsstammaResultatDispositionGodkannaStyrelsensForslag`,
		`se-bol-base:IntygandeOriginalInnehall`,
		`se-bol-base:UnderskriftFaststallelseintygElektroniskt`,
		`ID_DATUM_UNDERTECKNANDE_FASTSTALLELSEINTYG`,
		`continuedAt="intygande_forts"`,
		`id="intygande_forts"`,
	}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("certification missing: %s", check)
		}
	}
}

func TestGenerate_ManagementReport(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	checks := []string{
		`id="ar3-page-2"`,
		`id="ar3-page-3"`,
		`Förvaltningsberättelse`,
		`se-gen-base:AllmantVerksamheten`,
		`se-gen-base:VasentligaHandelserRakenskapsaret`,
		// Multi-year overview
		`Flerårsöversikt`,
		`se-gen-base:Nettoomsattning`,
		`se-gen-base:ResultatEfterFinansiellaPoster`,
		`se-gen-base:Soliditet`,
		// Equity changes
		`se-gen-base:Aktiekapital`,
		`se-gen-base:ForandringEgetKapitalTotalt`,
		// Profit disposition
		`se-gen-base:MedelDisponera`,
		`se-gen-base:ForslagDisposition`,
		// Board dividend statement
		`se-gen-base:StyrelsensYttrandeVinstutdelning`,
	}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("management report missing: %s", check)
		}
	}
}

func TestGenerate_IncomeStatement(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	checks := []string{
		`id="ar3-page-4"`,
		`Resultaträkning`,
		// Revenue
		`se-gen-base:Nettoomsattning`,
		`se-gen-base:ForandringLagerProdukterIArbeteFardigaVarorPagaendeArbetenAnnansRakning`,
		`se-gen-base:OvrigaRorelseintakter`,
		`se-gen-base:RorelseintakterLagerforandringarMm`,
		// Expenses
		`se-gen-base:RavarorFornodenheterKostnader`,
		`se-gen-base:Personalkostnader`,
		`se-gen-base:Rorelsekostnader`,
		// Results
		`se-gen-base:Rorelseresultat`,
		`se-gen-base:ResultatEfterFinansiellaPoster`,
		`se-gen-base:ResultatForeSkatt`,
		`se-gen-base:AretsResultat`,
		// Appropriations with sign="-"
		`sign="-"`,
		`se-gen-base:ForandringPeriodiseringsfond`,
		`se-gen-base:Bokslutsdispositioner`,
	}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("income statement missing: %s", check)
		}
	}
}

func TestGenerate_BalanceSheet(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	checks := []string{
		`id="ar3-page-5"`,
		`id="ar3-page-6"`,
		`Balansräkning`,
		// Assets
		`se-gen-base:ByggnaderMark`,
		`se-gen-base:MaskinerAndraTekniskaAnlaggningar`,
		`se-gen-base:InventarierVerktygInstallationer`,
		`se-gen-base:Anlaggningstillgangar`,
		`se-gen-base:Omsattningstillgangar`,
		`se-gen-base:Tillgangar`,
		// Equity & Liabilities
		`se-gen-base:Aktiekapital`,
		`se-gen-base:EgetKapital`,
		`se-gen-base:ObeskattadeReserver`,
		`se-gen-base:Avsattningar`,
		`se-gen-base:LangfristigaSkulder`,
		`se-gen-base:KortfristigaSkulder`,
		`se-gen-base:EgetKapitalSkulder`,
		// Note references
		`href="#note-3"`,
		`href="#note-4"`,
		`href="#note-5"`,
		`href="#note-6"`,
	}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("balance sheet missing: %s", check)
		}
	}
}

func TestGenerate_NoteAccountingPolicies(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	checks := []string{
		`id="note-1"`,
		`Not 1`,
		`Redovisnings- och värderingsprinciper`,
		`se-gen-base:Redovisningsprinciper`,
		`Avskrivningar`,
		`se-gen-base:AvskrivningarMateriellaAnlaggningstillgangarByggnaderAr`,
		`se-gen-base:AvskrivningarMateriellaAnlaggningstillgangarMaskinerAndraTekniskaAnlaggningarAr`,
		`se-gen-base:AvskrivningarMateriellaAnlaggningstillgangarInventarierVerktygInstallationerAr`,
		`se-gen-base:AvskrivningarMateriellaAnlaggningstillgangarKommentar`,
		`se-gen-base:RedovisningsprinciperAnskaffningsvardeEgentillverkadevaror`,
		`Nyckeltalsdefinitioner`,
	}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("note 1 missing: %s", check)
		}
	}
}

func TestGenerate_NoteEmployees(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	checks := []string{
		`id="note-2"`,
		`Not 2`,
		`Medelantalet anställda`,
		`se-gen-base:MedelantaletAnstallda`,
		`unitRef="antal-anstallda"`,
	}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("note 2 missing: %s", check)
		}
	}

	// Verify antal-anstallda does NOT have format attribute
	// Find the MedelantaletAnstallda nonFraction and check it lacks format
	idx := strings.Index(output, `name="se-gen-base:MedelantaletAnstallda"`)
	if idx < 0 {
		t.Fatal("MedelantaletAnstallda not found")
	}
	// Find the enclosing ix:nonFraction element
	start := strings.LastIndex(output[:idx], "<ix:nonFraction")
	end := strings.Index(output[idx:], "</ix:nonFraction>")
	if start < 0 || end < 0 {
		t.Fatal("could not find enclosing nonFraction element")
	}
	element := output[start : idx+end]
	if strings.Contains(element, `format=`) {
		t.Error("MedelantaletAnstallda should NOT have format attribute")
	}
}

func TestGenerate_FixedAssetNotes(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	// Note 3: Byggnader och mark
	checks3 := []string{
		`id="note-3"`,
		`Not 3`,
		`se-gen-base:ByggnaderMarkAnskaffningsvarden`,
		`se-gen-base:ByggnaderMarkForandringAnskaffningsvardenInkop`,
		`se-gen-base:ByggnaderMarkAvskrivningar`,
		`se-gen-base:ByggnaderMarkForandringAvskrivningarAretsAvskrivningar`,
	}
	for _, check := range checks3 {
		if !strings.Contains(output, check) {
			t.Errorf("note 3 missing: %s", check)
		}
	}

	// Note 4: Maskiner
	if !strings.Contains(output, `id="note-4"`) {
		t.Error("note 4 missing")
	}
	if !strings.Contains(output, `se-gen-base:MaskinerAndraTekniskaAnlaggningarAnskaffningsvarden`) {
		t.Error("note 4 missing acquisition values concept")
	}

	// Note 5: Inventarier
	if !strings.Contains(output, `id="note-5"`) {
		t.Error("note 5 missing")
	}
	if !strings.Contains(output, `se-gen-base:InventarierVerktygInstallationerAnskaffningsvarden`) {
		t.Error("note 5 missing acquisition values concept")
	}

	// Note 6: Financial assets (no depreciation)
	if !strings.Contains(output, `id="note-6"`) {
		t.Error("note 6 missing")
	}
	if !strings.Contains(output, `se-gen-base:AndraLangfristigaVardepappersinnehavAnskaffningsvarden`) {
		t.Error("note 6 missing acquisition values concept")
	}
	if !strings.Contains(output, `se-gen-base:AndraLangfristigaVardepappersinnehavForandringAnskaffningsvardenForsaljningar`) {
		t.Error("note 6 missing sales concept")
	}
}

func TestGenerate_ContextRefsInFixedAssetNotes(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	// Opening acquisition values should use balans1/balans2
	// Closing should use balans0/balans1
	// Changes should use period0/period1

	// Check for balans2 context ref (used in opening acquisition/depreciation for prev year column)
	if !strings.Contains(output, `contextRef="balans2"`) {
		t.Error("missing balans2 context ref in fixed asset notes")
	}
}

func TestGenerate_LongTermLiabilitiesNote(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	checks := []string{
		`id="note-7"`,
		`Not 7`,
		`Långfristiga skulder`,
		`se-gen-base:LangfristigaSkulderForfallerSenare5Ar`,
	}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("note 7 missing: %s", check)
		}
	}
}

func TestGenerate_PledgesNote(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	checks := []string{
		`id="note-8"`,
		`Not 8`,
		`Ställda säkerheter`,
		`se-gen-base:StalldaSakerheterForetagsinteckningar`,
		`se-gen-base:StalldaSakerheterFastighetsinteckningar`,
		`se-gen-base:StalldaSakerheter`,
	}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("note 8 missing: %s", check)
		}
	}
}

func TestGenerate_ContingentLiabilitiesNote(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	checks := []string{
		`id="note-9"`,
		`Not 9`,
		`Eventualförpliktelser`,
		`se-gen-base:EventualForpliktelser`,
	}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("note 9 missing: %s", check)
		}
	}
}

func TestGenerate_MultiPostNote(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	checks := []string{
		`id="note-10"`,
		`Not 10`,
		`Tillgångar, avsättningar och skulder som avser flera poster`,
		`se-gen-base:NotTillgangarAvsattningarSkulderAvserFleraPoster`,
		// Tuples
		`TillgangarAvsattningarSkulderTuple1`,
		`TillgangarAvsattningarSkulderTuple2`,
		`se-gen-base:TillgangarAvsattningarSkulderPost`,
		`se-gen-base:TillgangarAvsattningarSkulderBelopp`,
		// Headings
		`Långfristiga skulder`,
		`Kortfristiga skulder`,
		// Order attributes
		`order="1.0"`,
		`order="2.0"`,
	}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("note 10 missing: %s", check)
		}
	}
}

func TestGenerate_Signatures(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	checks := []string{
		`ar-signature-2`,
		`se-gen-base:UndertecknandeArsredovisningOrt`,
		`se-gen-base:UndertecknandeArsredovisningDatum`,
		`Sundsvall`,
		`2017-02-20`,
		// Signatory tuples
		`UnderskriftArsredovisningForetradareTuple1`,
		`UnderskriftArsredovisningForetradareTuple2`,
		`se-gen-base:UnderskriftArsredovisningForetradareTilltalsnamn`,
		`se-gen-base:UnderskriftArsredovisningForetradareEfternamn`,
		`se-gen-base:UnderskriftArsredovisningForetradareForetradarroll`,
		// Names
		`Karl`,
		`Karlsson`,
		`Karin`,
		`Olsson`,
		`Verkställande direktör`,
	}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("signatures missing: %s", check)
		}
	}
}

func TestGenerate_PageStructure(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	pages := []string{
		`id="ar3-page-1"`,  // cover
		`id="ar3-page-2"`,  // management p1
		`id="ar3-page-3"`,  // management p2
		`id="ar3-page-4"`,  // income statement
		`id="ar3-page-5"`,  // balance sheet assets
		`id="ar3-page-6"`,  // balance sheet equity
		`id="ar3-page-7"`,  // notes p1
		`id="ar3-page-8"`,  // notes p2
		`id="ar3-page-9"`,  // notes last (multi-post + signatures)
		`id="ar3-page-10"`, // notes p3
	}
	for _, page := range pages {
		if !strings.Contains(output, page) {
			t.Errorf("missing page: %s", page)
		}
	}
}

func TestGenerate_AmountFormatting(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	// Check that amounts are formatted with space separators
	// Nettoomsättning current = 2650000 → "2 650 000"
	if !strings.Contains(output, "2 650 000") {
		t.Error("missing formatted amount 2 650 000")
	}

	// Check tkr formatting in multi-year overview
	// NetSales 2650000 in tkr → "2 650"
	if !strings.Contains(output, `scale="3"`) {
		t.Error("missing scale=3 for tkr display")
	}
}

func TestGenerate_NegativeAmountDisplay(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	// Expenses should show "-" before the ix:nonFraction tag
	// Check for pattern: -<ix:nonFraction ... with an expense concept
	if !strings.Contains(output, `-<ix:nonFraction`) {
		t.Error("missing negative prefix display for expenses")
	}
}

func TestGenerate_WrapperDiv(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	if !strings.Contains(output, `id="wrapper"`) {
		t.Error("missing wrapper div")
	}
}

func TestGenerate_CSS(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	if !strings.Contains(output, `<style type="text/css">`) {
		t.Error("missing CSS style block")
	}
}

// TestGenerate_KeyXBRLFacts verifies specific XBRL fact values from the test data.
func TestGenerate_KeyXBRLFacts(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	// Company name
	assertContains(t, output, `>Exempel 1 AB</ix:nonNumeric>`, "company name fact")

	// Org nr
	assertContains(t, output, `>556999-9999</ix:nonNumeric>`, "org nr fact")

	// Fiscal year dates in hidden facts
	assertContains(t, output, `>2016-01-01</ix:nonNumeric>`, "fiscal year start")
	assertContains(t, output, `>2016-12-31</ix:nonNumeric>`, "fiscal year end")

	// Net result 1274000 → "1 274 000"
	assertContains(t, output, `name="se-gen-base:AretsResultat"`, "net result concept")
	assertContains(t, output, `>1 274 000</ix:nonFraction>`, "net result value")

	// Total assets 7773000 → "7 773 000"
	assertContains(t, output, `>7 773 000</ix:nonFraction>`, "total assets value")
}

func assertContains(t *testing.T, output, substring, description string) {
	t.Helper()
	if !strings.Contains(output, substring) {
		t.Errorf("missing %s: %s", description, substring)
	}
}

// TestGenerate_RoundtripStructure is a basic structural roundtrip test.
// It verifies that the generated output contains all the major document
// sections in the expected order.
func TestGenerate_RoundtripStructure(t *testing.T) {
	r := loadTestReport(t)
	output := generateOutput(t, r)

	// Verify sections appear in order
	sections := []struct {
		name   string
		marker string
	}{
		{"XML declaration", `<?xml`},
		{"HTML open", `<html`},
		{"head", `<head>`},
		{"style", `<style`},
		{"body", `<body>`},
		{"ix:header", `<ix:header>`},
		{"cover page", `id="ar3-page-1"`},
		{"management report p1", `id="ar3-page-2"`},
		{"management report p2", `id="ar3-page-3"`},
		{"income statement", `id="ar3-page-4"`},
		{"balance sheet assets", `id="ar3-page-5"`},
		{"balance sheet equity", `id="ar3-page-6"`},
		{"notes p1", `id="ar3-page-7"`},
		{"signatures", `class="ar-signature-2"`},
		{"body close", `</body>`},
		{"html close", `</html>`},
	}

	lastIdx := -1
	for _, s := range sections {
		idx := strings.Index(output, s.marker)
		if idx < 0 {
			t.Errorf("section %s not found (marker: %s)", s.name, s.marker)
			continue
		}
		if idx <= lastIdx {
			t.Errorf("section %s appears before previous section (at %d, previous at %d)", s.name, idx, lastIdx)
		}
		lastIdx = idx
	}
}
