# Redofri - Projektplan

Digital inlämning av svensk årsredovisning. Första version: K2 för aktiebolag med fastställelseintyg.

## Översikt

Programmet genererar en iXBRL-fil (.xhtml) som är en giltig svensk årsredovisning
enligt K2, redo att lämnas in digitalt till Bolagsverket.

iXBRL-filen är samtidigt:
- Läsbar i webbläsare (visuell årsredovisning med CSS)
- Maskinläsbar (XBRL-taggad data inbäddad i HTML)
- Helt self-contained (inline CSS, base64-bilder, inga externa resurser)

## Datakällor

Tre sätt att mata in data, som alla mynnar ut i samma interna modell:

```
SIE-fil ──────────parse──┐
Föregående årsredovisning ┼──► model.AnnualReport ──► ixbrl.Generate() ──► .xhtml
Webb/manuell inmatning ───┘
```

### SIE-fil
Ger kontosaldon (resultat + balans), företagsinfo och räkenskapsår.
Kan inte ensam bli en årsredovisning - kräver kompletterande inmatning
av förvaltningsberättelse, noter, styrelseledamöter, resultatdisposition m.m.

### Föregående årsredovisning (.xhtml)
Parsas tillbaka till modellen. Ger:
- Jämförelsetal (föregående års resultat- och balansräkning)
- Flerårsöversikt (historiska nyckeltal)
- Ingående balanser i noter (anskaffningsvärden, ackumulerade avskrivningar)
- Bolagsinfo som sällan ändras (styrelseledamöter, redovisningsprinciper, säte)
- Strukturval (entry point, vilka noter som användes)

### Manuell inmatning
Webb-formulär (framtida) eller JSON-fil (utveckling/test).

### Typiskt användarflöde
1. Ladda upp föregående års årsredovisning → jämförelsetal och bolagsinfo fylls i
2. Ladda upp SIE-fil → aktuellt års saldon fylls i
3. Komplettera: förvaltningsberättelse, resultatdisposition, eventuella nya noter
4. Generera → .xhtml

## Arkitektur

### Designprinciper
- **Separera kärna från gränssnitt** - all logik i pkg/, CLI och framtida webb är tunna skal
- **AI-testbarhet** - alla komponenter testbara via `go test`, deterministisk output
- **Roundtrip-testning** - generera fil, parsa tillbaka, verifiera att modellen matchar
- **Stegvis bygge** - varje steg producerar något körbart och verifierbart

### Kodstruktur

```
redofri/
├── cmd/
│   └── redofri/              # CLI entry point
│       └── main.go
├── pkg/
│   ├── model/                # Go structs - det centrala kontraktet
│   │   └── model.go          # AnnualReport, IncomeStatement, BalanceSheet, etc.
│   ├── ixbrl/
│   │   ├── generate.go       # model → .xhtml (iXBRL-generering)
│   │   ├── contexts.go       # XBRL-kontexer och enheter
│   │   ├── tags.go           # iXBRL-taggningshjälpare (ix:nonNumeric, ix:nonFraction)
│   │   ├── parse.go          # .xhtml → model (läsa tillbaka tidigare år)
│   │   └── templates/        # HTML-templates (go:embed)
│   ├── sie/                  # SIE-parser → fyller model.AnnualReport
│   │   └── parse.go
│   ├── taxonomy/             # Taxonomihantering - koncept, entry points
│   │   └── taxonomy.go
│   └── validate/             # Validerar model.AnnualReport (oavsett källa)
│       └── validate.go
├── testdata/                 # JSON-indata och förväntad output
├── ref/                      # Referensmaterial
├── go.mod
├── PLAN.md
└── LICENSE
```

## Taxonomi

### K2 för aktiebolag (2024-09-12)

Fyra entry points beroende på fullständig/förkortad resultat- och balansräkning:

| Entry point | Resultaträkning | Balansräkning |
|-------------|-----------------|---------------|
| risbs       | Fullständig     | Fullständig   |
| risab       | Fullständig     | Förkortad     |
| raibs       | Förkortad       | Fullständig   |
| raiab       | Förkortad       | Förkortad     |

Entry point URL-mönster:
```
http://xbrl.taxonomier.se/se/fr/gaap/k2-all/ab/{variant}/2024-09-12/se-k2-ab-{variant}-2024-09-12.xsd
```

### Taxonomikombinationer (för K2 2024-09-12)
- K2-taxonomi: 2024-09-12
- Fastställelseintyg för aktiebolag (endast årsredovisning): 2020-12-01
- Revisionsberättelse (valfritt): 2020-12-01

### Namnrymder
| Prefix       | URI                                                     | Innehåll                    |
|--------------|---------------------------------------------------------|-----------------------------|
| se-gen-base  | http://www.taxonomier.se/se/fr/gen-base/2021-10-31      | Finansiella koncept (~107)  |
| se-cd-base   | http://www.taxonomier.se/se/fr/cd-base/2021-10-31       | Företagsdata (~8)           |
| se-bol-base  | http://www.bolagsverket.se/se/fr/comp-base/2017-09-30   | Fastställelseintyg (~10)    |

## iXBRL-dokumentets sektioner

1. **Head** - title, meta (programvara, programversion), inline CSS
2. **ix:header** (display:none) - schemareferenser, kontexer, enheter, dolda metadata
3. **Framsida** - företagsnamn, org.nr, räkenskapsår
4. **Fastställelseintyg** - intygande att stämman fastställt resultat/balansräkning
5. **Förvaltningsberättelse** - verksamhet, väsentliga händelser, flerårsöversikt, resultatdisposition
6. **Resultaträkning** - intäkter och kostnader
7. **Balansräkning** - tillgångar, eget kapital och skulder
8. **Noter** - redovisningsprinciper, specifikationer
9. **Underskrifter** - ort, datum, namn och roll per styrelseledamot

## iXBRL-regler (från Bolagsverkets tillämpningsanvisningar)

### Format
- iXBRL 1.1, giltig XHTML, UTF-8
- Helt self-contained (inga externa resurser)
- Max 5 MB totalt, bilder max 1 MB styck (base64 JPEG/PNG/SVG/GIF)
- Inga script, inga event handlers
- Inline CSS (ej extern)

### Kontextnamngivning
- Duration: `period0` (aktuellt år), `period1` (föregående), `period2`, `period3`
- Instant: `balans0` (aktuellt årsskifte), `balans1` (föregående), `balans2`, `balans3`
- Entitet: org.nr med scheme `http://www.bolagsverket.se`

### Enheter
- `SEK` → iso4217:SEK
- `procent` → xbrli:pure
- `antal-anstallda` → se-k2-type:AntalAnstallda

### Belopp
- Använd `decimals`-attribut (inte `precision`)
- Negativa värden: `sign="-"` attribut
- Format: `ixt:numspacecomma` (t.ex. "2 650 000") eller `ixt:numcomma` (t.ex. "33,7")
- Scale: `0` för exakta belopp, `3` för tkr, `-2` för procent

### Metadata
- `<meta name="programvara" content="Redofri"/>` i head
- `<meta name="programversion" content="x.y.z"/>` i head
- Dolda XBRL-fakta: Språk, Land, Redovisningsvaluta, Beloppsformat, Räkenskapsår

### Fastställelseintyg
- Signeringsdatum med `id="ID_DATUM_UNDERTECKNANDE_FASTSTALLELSEINTYG"`
- `ArsstammaIntygande` omsluter delfakta med continuation-mönster

### Kontrollsumma
- SHA-256 via Bolagsverkets API
- Sektioner exkluderade från checksumma markeras med specifika id-attribut

## API-integration (framtida steg)

### Trestegsprocess
1. `skapa-inlamningtoken` - skapa token + visa avtalstexten
2. `kontrollera` (valfritt) - validera, få varningar/fel
3. `inlamning` - ladda upp dokumentet

### Krav
- Organisationscertifikat (TLS-klientcertifikat) från Expisoft/Steria
- Avtal med Bolagsverket
- Miljöer: Test → Acceptanstest → Produktion

### Endpoints (produktion)
- Information: `api.bolagsverket.se/lamna-in-arsredovisning/v2.1/`
- Händelser: `api.bolagsverket.se/hantera-arsredovisningsprenumerationer/v2.0/`

## Implementationsordning

### Steg 1: Grundstomme + Datamodell
- [ ] Go-modul, projektstruktur
- [ ] model/ med Go structs för alla delar av årsredovisningen
- [ ] JSON-serialisering (utveckling/test)
- [ ] Testdata baserad på "Exempel 1 AB" från taxonomier.se

### Steg 2: iXBRL-generering (kärnan)
- [ ] XHTML-grundstruktur med inline CSS
- [ ] ix:header - kontexer, enheter, schemareferenser, dolda metadata
- [ ] Framsida + fastställelseintyg
- [ ] Förvaltningsberättelse med flerårsöversikt
- [ ] Resultaträkning (kostnadsslagsindelad, fullständig)
- [ ] Balansräkning (fullständig)
- [ ] Noter
- [ ] Underskrifter
- [ ] Alla ~125 XBRL-koncept korrekt taggade

### Steg 3: CLI
- [ ] `redofri generate --input data.json --output arsredovisning.xhtml`
- [ ] `redofri validate --input data.json`
- [ ] `redofri example` (skriver exempelfil med JSON-indata)

### Steg 4: iXBRL-parser (roundtrip)
- [ ] Parsa .xhtml tillbaka till model.AnnualReport
- [ ] Roundtrip-test: generera → parsa → jämför
- [ ] Stöd för att läsa in föregående års årsredovisning

### Steg 5: SIE-import
- [ ] Parsa SIE4-filer
- [ ] Mappa BAS-konton till K2-taxonomikoncept
- [ ] Fyll model.AnnualReport med saldon

### Steg 6: Validering
- [ ] Beräkningskontroller (summor stämmer)
- [ ] Obligatoriska fält
- [ ] Taxonomiregler
- [ ] Jämför mot Bolagsverkets valideringskoder

### Steg 7: API-integration
- [ ] TLS-klientcertifikat
- [ ] Token-hantering
- [ ] Validering via API
- [ ] Inlämning
- [ ] Statushantering

### Steg 8: Webbgränssnitt
- [ ] Formulär för manuell inmatning
- [ ] SIE-uppladdning
- [ ] Uppladdning av föregående årsredovisning
- [ ] Förhandsgranskning
- [ ] Inlämningsflöde
