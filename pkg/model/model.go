// Package model defines the data structures for a Swedish K2 annual report (årsredovisning).
//
// These structs are the central contract of the application. All data sources
// (SIE import, previous year iXBRL parsing, manual/JSON input) populate these
// structs, and the iXBRL generator reads them to produce the output file.
//
// Field names follow Swedish accounting terminology. JSON tags use camelCase.
// XBRL concept names are documented in comments for traceability.
package model

// AnnualReport is the top-level struct representing a complete K2 årsredovisning.
type AnnualReport struct {
	// Metadata
	Company    Company    `json:"company"`
	FiscalYear FiscalYear `json:"fiscalYear"`
	Meta       Meta       `json:"meta"`

	// Document sections
	Certification    Certification    `json:"certification"`
	ManagementReport ManagementReport `json:"managementReport"`
	IncomeStatement  IncomeStatement  `json:"incomeStatement"`
	BalanceSheet     BalanceSheet     `json:"balanceSheet"`
	Notes            Notes            `json:"notes"`
	Signatures       Signatures       `json:"signatures"`
}

// Company holds basic company information.
// XBRL concepts: se-cd-base:ForetagetsNamn, se-cd-base:Organisationsnummer
type Company struct {
	Name  string `json:"name"`  // se-cd-base:ForetagetsNamn
	OrgNr string `json:"orgNr"` // se-cd-base:Organisationsnummer
}

// FiscalYear defines the reporting period.
// XBRL concepts: se-cd-base:RakenskapsarForstaDag, se-cd-base:RakenskapsarSistaDag
type FiscalYear struct {
	StartDate string `json:"startDate"` // YYYY-MM-DD, se-cd-base:RakenskapsarForstaDag
	EndDate   string `json:"endDate"`   // YYYY-MM-DD, se-cd-base:RakenskapsarSistaDag
}

// Meta holds iXBRL metadata and generation settings.
// XBRL concepts: se-cd-base:Sprak, se-cd-base:Land, se-cd-base:Redovisningsvaluta, se-cd-base:Beloppsformat
type Meta struct {
	Language     string `json:"language"`     // se-cd-base:Sprak, e.g. "sv"
	Country      string `json:"country"`      // se-cd-base:Land, e.g. "SE"
	Currency     string `json:"currency"`     // se-cd-base:Redovisningsvaluta, e.g. "SEK"
	AmountFormat string `json:"amountFormat"` // se-cd-base:Beloppsformat, e.g. "NORMALFORM"

	// Entry point variant: "risbs", "risab", "raibs", or "raiab"
	EntryPoint string `json:"entryPoint"`

	// Software info for <meta> tags
	Software        string `json:"software"`
	SoftwareVersion string `json:"softwareVersion"`
}

// PreviousYear holds data for the comparative period (föregående år).
// This allows the same struct hierarchy to hold both current and previous year data.
// In the model, most numeric fields use *int64 (pointer) to distinguish
// "not reported" (nil) from zero.

// Certification represents the fastställelseintyg section.
// XBRL concepts from se-bol-base namespace.
type Certification struct {
	// se-bol-base:FaststallelseResultatBalansrakning
	ConfirmationText string `json:"confirmationText"`
	// se-bol-base:Arsstamma (date of the annual general meeting)
	MeetingDate string `json:"meetingDate"`
	// se-bol-base:ArsstammaResultatDispositionGodkannaStyrelsensForslag
	DispositionDecision string `json:"dispositionDecision"`
	// se-bol-base:IntygandeOriginalInnehall
	OriginalContentCertification string `json:"originalContentCertification"`
	// se-bol-base:UnderskriftFaststallelseintygElektroniskt
	ElectronicSignatureLabel string `json:"electronicSignatureLabel"`
	// Signatory of the certification
	Signatory CertificationSignatory `json:"signatory"`
	// se-bol-base:UnderskriftFastallelseintygDatum (id="ID_DATUM_UNDERTECKNANDE_FASTSTALLELSEINTYG")
	SigningDate string `json:"signingDate"`
}

// CertificationSignatory is the person who signs the fastställelseintyg.
type CertificationSignatory struct {
	FirstName string `json:"firstName"` // se-bol-base:UnderskriftFaststallelseintygForetradareTilltalsnamn
	LastName  string `json:"lastName"`  // se-bol-base:UnderskriftFaststallelseintygForetradareEfternamn
	Role      string `json:"role"`      // se-bol-base:UnderskriftFaststallelseintygForetradareForetradarroll
}

// ManagementReport represents the förvaltningsberättelse.
type ManagementReport struct {
	// se-gen-base:LopandeBokforingenAvslutasMening
	IntroText string `json:"introText"`
	// se-gen-base:AllmantVerksamheten
	BusinessDescription string `json:"businessDescription"`
	// se-gen-base:VasentligaHandelserRakenskapsaret
	SignificantEvents string `json:"significantEvents"`

	// Flerårsöversikt (multi-year overview)
	MultiYearOverview MultiYearOverview `json:"multiYearOverview"`

	// Förändringar i eget kapital
	EquityChanges EquityChanges `json:"equityChanges"`

	// Resultatdisposition
	ProfitDisposition ProfitDisposition `json:"profitDisposition"`

	// se-gen-base:StyrelsensYttrandeVinstutdelning
	BoardDividendStatement string `json:"boardDividendStatement,omitempty"`
}

// MultiYearOverviewYear holds one year's worth of overview key figures.
type MultiYearOverviewYear struct {
	// The period label (year or date range)
	Year string `json:"year"`

	// se-gen-base:Nettoomsattning (in tkr for overview, full for IS)
	NetSales *int64 `json:"netSales,omitempty"`
	// se-gen-base:ResultatEfterFinansiellaPoster (tkr)
	ResultAfterFinancialItems *int64 `json:"resultAfterFinancialItems,omitempty"`
	// se-gen-base:Soliditet (percent, stored as decimal e.g. 33.7 → 337 with scale -2)
	Solidity *string `json:"solidity,omitempty"`
}

// MultiYearOverview holds the flerårsöversikt table and optional comment.
type MultiYearOverview struct {
	Years []MultiYearOverviewYear `json:"years"`
	// se-gen-base:KommentarFlerarsoversikt
	Comment string `json:"comment,omitempty"`
}

// EquityChanges represents the förändringar i eget kapital table.
type EquityChanges struct {
	// Opening balances (belopp vid årets ingång) — from previous year-end (balans1)
	OpeningShareCapital     *int64 `json:"openingShareCapital"`     // se-gen-base:Aktiekapital @ balans1
	OpeningReserveFund      *int64 `json:"openingReserveFund"`      // se-gen-base:Reservfond @ balans1
	OpeningRetainedEarnings *int64 `json:"openingRetainedEarnings"` // se-gen-base:BalanseratResultat @ balans1
	OpeningNetIncome        *int64 `json:"openingNetIncome"`        // se-gen-base:AretsResultatEgetKapital @ balans1
	OpeningTotal            *int64 `json:"openingTotal"`            // se-gen-base:ForandringEgetKapitalTotalt @ balans1

	// Dividend from previous year result
	DividendNetIncome *int64 `json:"dividendNetIncome,omitempty"` // se-gen-base:ForandringEgetKapitalAretsResultatUtdelning
	DividendTotal     *int64 `json:"dividendTotal,omitempty"`     // se-gen-base:ForandringEgetKapitalTotaltUtdelning

	// Year's result
	YearResultNetIncome *int64 `json:"yearResultNetIncome"` // se-gen-base:ForandringEgetKapitalAretsResultatAretsResultat
	YearResultTotal     *int64 `json:"yearResultTotal"`     // se-gen-base:ForandringEgetKapitalTotaltAretsResultat

	// Closing balances (belopp vid årets utgång) — at current year-end (balans0)
	ClosingShareCapital     *int64 `json:"closingShareCapital"`     // se-gen-base:Aktiekapital @ balans0
	ClosingReserveFund      *int64 `json:"closingReserveFund"`      // se-gen-base:Reservfond @ balans0
	ClosingRetainedEarnings *int64 `json:"closingRetainedEarnings"` // se-gen-base:BalanseratResultat @ balans0
	ClosingNetIncome        *int64 `json:"closingNetIncome"`        // se-gen-base:AretsResultatEgetKapital @ balans0
	ClosingTotal            *int64 `json:"closingTotal"`            // se-gen-base:ForandringEgetKapitalTotalt @ balans0
}

// ProfitDisposition represents the resultatdisposition section.
type ProfitDisposition struct {
	// Available funds
	RetainedEarnings *int64 `json:"retainedEarnings"` // se-gen-base:BalanseratResultat @ balans0
	NetIncome        *int64 `json:"netIncome"`        // se-gen-base:AretsResultatEgetKapital @ balans0
	TotalAvailable   *int64 `json:"totalAvailable"`   // se-gen-base:MedelDisponera @ balans0

	// Proposed disposition
	Dividend         *int64 `json:"dividend,omitempty"` // se-gen-base:ForslagDispositionUtdelning
	CarriedForward   *int64 `json:"carriedForward"`     // se-gen-base:ForslagDispositionBalanserasINyRakning
	TotalDisposition *int64 `json:"totalDisposition"`   // se-gen-base:ForslagDisposition
}

// IncomeStatement represents the resultaträkning (kostnadsslagsindelad, fullständig).
// Current year = period0, previous year = period1.
type IncomeStatement struct {
	// Rörelseintäkter, lagerförändringar m.m.
	Revenue IncomeStatementRevenue `json:"revenue"`
	// Rörelsekostnader
	Expenses IncomeStatementExpenses `json:"expenses"`
	// Rörelseresultat: se-gen-base:Rorelseresultat
	OperatingResult YearComparison `json:"operatingResult"`
	// Finansiella poster
	FinancialItems IncomeStatementFinancialItems `json:"financialItems"`
	// Resultat efter finansiella poster: se-gen-base:ResultatEfterFinansiellaPoster
	ResultAfterFinancialItems YearComparison `json:"resultAfterFinancialItems"`
	// Bokslutsdispositioner
	Appropriations IncomeStatementAppropriations `json:"appropriations"`
	// Resultat före skatt: se-gen-base:ResultatForeSkatt
	ResultBeforeTax YearComparison `json:"resultBeforeTax"`
	// Skatter
	Tax IncomeStatementTax `json:"tax"`
	// Årets resultat: se-gen-base:AretsResultat
	NetResult YearComparison `json:"netResult"`
}

// YearComparison holds a value for current and previous year.
type YearComparison struct {
	Current  *int64 `json:"current"`
	Previous *int64 `json:"previous"`
}

// IncomeStatementRevenue holds revenue line items.
type IncomeStatementRevenue struct {
	// se-gen-base:Nettoomsattning
	NetSales YearComparison `json:"netSales"`
	// se-gen-base:ForandringLagerProdukterIArbeteFardigaVarorPagaendeArbetenAnnansRakning
	InventoryChange YearComparison `json:"inventoryChange,omitempty"`
	// se-gen-base:OvrigaRorelseintakter
	OtherOperatingIncome YearComparison `json:"otherOperatingIncome,omitempty"`
	// se-gen-base:RorelseintakterLagerforandringarMm
	TotalRevenue YearComparison `json:"totalRevenue"`
}

// IncomeStatementExpenses holds expense line items.
type IncomeStatementExpenses struct {
	// se-gen-base:RavarorFornodenheterKostnader
	RawMaterials YearComparison `json:"rawMaterials,omitempty"`
	// se-gen-base:HandelsvarorKostnader
	TradingGoods YearComparison `json:"tradingGoods,omitempty"`
	// se-gen-base:OvrigaExternaKostnader
	OtherExternalExpenses YearComparison `json:"otherExternalExpenses,omitempty"`
	// se-gen-base:Personalkostnader
	PersonnelExpenses YearComparison `json:"personnelExpenses,omitempty"`
	// Note reference for personnel expenses
	PersonnelExpensesNote int `json:"personnelExpensesNote,omitempty"`
	// se-gen-base:AvskrivningarNedskrivningarMateriellaImmateriellaAnlaggningstillgangar
	DepreciationAmortization YearComparison `json:"depreciationAmortization,omitempty"`
	// se-gen-base:OvrigaRorelsekostnader
	OtherOperatingExpenses YearComparison `json:"otherOperatingExpenses,omitempty"`
	// se-gen-base:Rorelsekostnader
	TotalExpenses YearComparison `json:"totalExpenses"`
}

// IncomeStatementFinancialItems holds financial items.
type IncomeStatementFinancialItems struct {
	// se-gen-base:ResultatOvrigaFinansiellaAnlaggningstillgangar
	ResultOtherFinancialAssets YearComparison `json:"resultOtherFinancialAssets,omitempty"`
	// se-gen-base:OvrigaRanteintakterLiknandeResultatposter
	OtherInterestIncome YearComparison `json:"otherInterestIncome,omitempty"`
	// se-gen-base:RantekostnaderLiknandeResultatposter
	InterestExpenses YearComparison `json:"interestExpenses,omitempty"`
	// se-gen-base:FinansiellaPoster
	TotalFinancialItems YearComparison `json:"totalFinancialItems"`
}

// IncomeStatementAppropriations holds appropriation items (bokslutsdispositioner).
type IncomeStatementAppropriations struct {
	// se-gen-base:ForandringPeriodiseringsfond
	TaxAllocationReserve YearComparison `json:"taxAllocationReserve,omitempty"`
	// se-gen-base:ForandringOveravskrivningar
	ExcessDepreciation YearComparison `json:"excessDepreciation,omitempty"`
	// se-gen-base:Bokslutsdispositioner
	TotalAppropriations YearComparison `json:"totalAppropriations"`
}

// IncomeStatementTax holds tax items.
type IncomeStatementTax struct {
	// se-gen-base:SkattAretsResultat
	IncomeTax YearComparison `json:"incomeTax"`
}

// BalanceSheet represents the balansräkning.
// Current year-end = balans0, previous year-end = balans1.
type BalanceSheet struct {
	Assets               Assets               `json:"assets"`
	EquityAndLiabilities EquityAndLiabilities `json:"equityAndLiabilities"`
}

// Assets holds the tillgångar side of the balance sheet.
type Assets struct {
	FixedAssets   FixedAssets   `json:"fixedAssets"`
	CurrentAssets CurrentAssets `json:"currentAssets"`
	// se-gen-base:Tillgangar
	TotalAssets YearComparison `json:"totalAssets"`
}

// FixedAssets holds anläggningstillgångar.
type FixedAssets struct {
	Tangible  TangibleFixedAssets  `json:"tangible"`
	Financial FinancialFixedAssets `json:"financial"`
	// se-gen-base:Anlaggningstillgangar
	TotalFixedAssets YearComparison `json:"totalFixedAssets"`
}

// TangibleFixedAssets holds materiella anläggningstillgångar.
type TangibleFixedAssets struct {
	// se-gen-base:ByggnaderMark
	BuildingsAndLand     YearComparison `json:"buildingsAndLand,omitempty"`
	BuildingsAndLandNote int            `json:"buildingsAndLandNote,omitempty"`
	// se-gen-base:MaskinerAndraTekniskaAnlaggningar
	MachineryAndEquipment     YearComparison `json:"machineryAndEquipment,omitempty"`
	MachineryAndEquipmentNote int            `json:"machineryAndEquipmentNote,omitempty"`
	// se-gen-base:InventarierVerktygInstallationer
	FixturesAndFittings     YearComparison `json:"fixturesAndFittings,omitempty"`
	FixturesAndFittingsNote int            `json:"fixturesAndFittingsNote,omitempty"`
	// se-gen-base:MateriellaAnlaggningstillgangar
	TotalTangible YearComparison `json:"totalTangible"`
}

// FinancialFixedAssets holds finansiella anläggningstillgångar.
type FinancialFixedAssets struct {
	// se-gen-base:AndraLangfristigaVardepappersinnehav
	OtherLongTermSecurities     YearComparison `json:"otherLongTermSecurities,omitempty"`
	OtherLongTermSecuritiesNote int            `json:"otherLongTermSecuritiesNote,omitempty"`
	// se-gen-base:FinansiellaAnlaggningstillgangar
	TotalFinancial YearComparison `json:"totalFinancial"`
}

// CurrentAssets holds omsättningstillgångar.
type CurrentAssets struct {
	Inventory            Inventory            `json:"inventory"`
	ShortTermReceivables ShortTermReceivables `json:"shortTermReceivables"`
	CashAndBank          CashAndBank          `json:"cashAndBank"`
	// se-gen-base:Omsattningstillgangar
	TotalCurrentAssets YearComparison `json:"totalCurrentAssets"`
}

// Inventory holds varulager m.m.
type Inventory struct {
	// se-gen-base:LagerRavarorFornodenheter
	RawMaterials YearComparison `json:"rawMaterials,omitempty"`
	// se-gen-base:LagerVarorUnderTillverkning
	WorkInProgress YearComparison `json:"workInProgress,omitempty"`
	// se-gen-base:LagerFardigaVarorHandelsvaror
	FinishedGoods YearComparison `json:"finishedGoods,omitempty"`
	// se-gen-base:VarulagerMm
	TotalInventory YearComparison `json:"totalInventory"`
}

// ShortTermReceivables holds kortfristiga fordringar.
type ShortTermReceivables struct {
	// se-gen-base:Kundfordringar
	TradeReceivables YearComparison `json:"tradeReceivables,omitempty"`
	// se-gen-base:OvrigaFordringarKortfristiga
	OtherReceivables YearComparison `json:"otherReceivables,omitempty"`
	// se-gen-base:ForutbetaldaKostnaderUpplupnaIntakter
	PrepaidExpenses YearComparison `json:"prepaidExpenses,omitempty"`
	// se-gen-base:KortfristigaFordringar
	TotalShortTermReceivables YearComparison `json:"totalShortTermReceivables"`
}

// CashAndBank holds kassa och bank.
type CashAndBank struct {
	// se-gen-base:KassaBankExklRedovisningsmedel
	CashAndBankExcl YearComparison `json:"cashAndBankExcl,omitempty"`
	// se-gen-base:KassaBank
	TotalCashAndBank YearComparison `json:"totalCashAndBank"`
}

// EquityAndLiabilities holds the eget kapital och skulder side.
type EquityAndLiabilities struct {
	Equity               Equity               `json:"equity"`
	UntaxedReserves      UntaxedReserves      `json:"untaxedReserves"`
	Provisions           Provisions           `json:"provisions"`
	LongTermLiabilities  LongTermLiabilities  `json:"longTermLiabilities"`
	ShortTermLiabilities ShortTermLiabilities `json:"shortTermLiabilities"`
	// se-gen-base:EgetKapitalSkulder
	TotalEquityAndLiabilities YearComparison `json:"totalEquityAndLiabilities"`
}

// Equity holds eget kapital on the balance sheet.
type Equity struct {
	// Bundet eget kapital
	// se-gen-base:Aktiekapital
	ShareCapital YearComparison `json:"shareCapital"`
	// se-gen-base:Reservfond
	ReserveFund YearComparison `json:"reserveFund,omitempty"`
	// se-gen-base:BundetEgetKapital
	TotalRestrictedEquity YearComparison `json:"totalRestrictedEquity"`

	// Fritt eget kapital
	// se-gen-base:BalanseratResultat
	RetainedEarnings YearComparison `json:"retainedEarnings"`
	// se-gen-base:AretsResultatEgetKapital
	NetIncome YearComparison `json:"netIncome"`
	// se-gen-base:FrittEgetKapital
	TotalUnrestrictedEquity YearComparison `json:"totalUnrestrictedEquity"`

	// se-gen-base:EgetKapital
	TotalEquity YearComparison `json:"totalEquity"`
}

// UntaxedReserves holds obeskattade reserver.
type UntaxedReserves struct {
	// se-gen-base:Periodiseringsfonder
	TaxAllocationReserves YearComparison `json:"taxAllocationReserves,omitempty"`
	// se-gen-base:AckumuleradeOveravskrivningar
	AccumulatedExcessDepreciation YearComparison `json:"accumulatedExcessDepreciation,omitempty"`
	// se-gen-base:ObeskattadeReserver
	TotalUntaxedReserves YearComparison `json:"totalUntaxedReserves"`
}

// Provisions holds avsättningar.
type Provisions struct {
	// se-gen-base:AvsattningarPensionerLiknandeForpliktelserEnligtLag
	PensionProvisions YearComparison `json:"pensionProvisions,omitempty"`
	// se-gen-base:OvrigaAvsattningar
	OtherProvisions YearComparison `json:"otherProvisions,omitempty"`
	// se-gen-base:Avsattningar
	TotalProvisions YearComparison `json:"totalProvisions"`
}

// LongTermLiabilities holds långfristiga skulder.
type LongTermLiabilities struct {
	LongTermLiabilitiesNote int `json:"longTermLiabilitiesNote,omitempty"`
	// se-gen-base:OvrigaLangfristigaSkulderKreditinstitut
	BankLoans      YearComparison `json:"bankLoans,omitempty"`
	BankLoansNotes []int          `json:"bankLoansNotes,omitempty"` // multiple note refs possible
	// se-gen-base:OvrigaLangfristigaSkulder
	OtherLongTermLiabilities YearComparison `json:"otherLongTermLiabilities,omitempty"`
	// se-gen-base:LangfristigaSkulder
	TotalLongTermLiabilities YearComparison `json:"totalLongTermLiabilities"`
}

// ShortTermLiabilities holds kortfristiga skulder.
type ShortTermLiabilities struct {
	// se-gen-base:Leverantorsskulder
	TradePayables YearComparison `json:"tradePayables,omitempty"`
	// se-gen-base:Skatteskulder
	TaxLiabilities YearComparison `json:"taxLiabilities,omitempty"`
	// se-gen-base:OvrigaKortfristigaSkulder
	OtherShortTermLiabilities     YearComparison `json:"otherShortTermLiabilities,omitempty"`
	OtherShortTermLiabilitiesNote int            `json:"otherShortTermLiabilitiesNote,omitempty"`
	// se-gen-base:UpplupnaKostnaderForutbetaldaIntakter
	AccruedExpenses YearComparison `json:"accruedExpenses,omitempty"`
	// se-gen-base:KortfristigaSkulder
	TotalShortTermLiabilities YearComparison `json:"totalShortTermLiabilities"`
}

// Notes holds all notes to the annual report.
type Notes struct {
	// Note 1: Redovisnings- och värderingsprinciper
	AccountingPolicies AccountingPolicies `json:"accountingPolicies"`

	// Note 2: Medelantalet anställda
	Employees *EmployeesNote `json:"employees,omitempty"`

	// Notes 3-6: Asset roll-forward notes
	FixedAssetNotes []FixedAssetNote `json:"fixedAssetNotes,omitempty"`

	// Note 7: Långfristiga skulder (förfaller efter 5 år)
	LongTermLiabilitiesNote *LongTermLiabilitiesNoteData `json:"longTermLiabilitiesNote,omitempty"`

	// Note 8: Ställda säkerheter
	Pledges *PledgesNote `json:"pledges,omitempty"`

	// Note 9: Eventualförpliktelser
	ContingentLiabilities *ContingentLiabilitiesNote `json:"contingentLiabilities,omitempty"`

	// Note 10: Tillgångar, avsättningar och skulder som avser flera poster
	MultiPostNote *MultiPostNote `json:"multiPostNote,omitempty"`
}

// AccountingPolicies represents note 1.
type AccountingPolicies struct {
	NoteNumber int `json:"noteNumber"` // typically 1

	// se-gen-base:Redovisningsprinciper
	Description string `json:"description"`

	// Depreciation periods
	Depreciations []DepreciationPolicy `json:"depreciations,omitempty"`

	// se-gen-base:AvskrivningarMateriellaAnlaggningstillgangarKommentar
	DepreciationComment string `json:"depreciationComment,omitempty"`

	// se-gen-base:RedovisningsprinciperAnskaffningsvardeEgentillverkadevaror
	ManufacturedGoodsPolicy string `json:"manufacturedGoodsPolicy,omitempty"`
}

// DepreciationPolicy holds a single depreciation period declaration.
type DepreciationPolicy struct {
	// Asset category label
	Category string `json:"category"`
	// XBRL concept name for the depreciation years fact
	// e.g. "se-gen-base:AvskrivningarMateriellaAnlaggningstillgangarByggnaderAr"
	Concept string `json:"concept"`
	// Number of years
	Years int `json:"years"`
}

// EmployeesNote represents note 2 (medelantalet anställda).
type EmployeesNote struct {
	NoteNumber int `json:"noteNumber"` // typically 2

	// se-gen-base:MedelantaletAnstallda
	AverageEmployees YearComparison `json:"averageEmployees"`
}

// FixedAssetNote represents a roll-forward note for a single asset category (notes 3-6 typically).
type FixedAssetNote struct {
	NoteNumber int    `json:"noteNumber"`
	Title      string `json:"title"` // e.g. "Byggnader och mark"

	// The XBRL concept prefix for this asset category, e.g. "ByggnaderMark"
	ConceptPrefix string `json:"conceptPrefix"`

	// Acquisition values (anskaffningsvärden)
	// Opening: {ConceptPrefix}Anskaffningsvarden @ balans1/balans2
	OpeningAcquisitionValues YearComparison `json:"openingAcquisitionValues"`
	// Purchases: {ConceptPrefix}ForandringAnskaffningsvardenInkop
	Purchases YearComparison `json:"purchases,omitempty"`
	// Sales: {ConceptPrefix}ForandringAnskaffningsvardenForsaljningar
	Sales YearComparison `json:"sales,omitempty"`
	// Closing: {ConceptPrefix}Anskaffningsvarden @ balans0/balans1
	ClosingAcquisitionValues YearComparison `json:"closingAcquisitionValues"`

	// Depreciation (avskrivningar)
	// Opening: {ConceptPrefix}Avskrivningar @ balans1/balans2
	OpeningDepreciation YearComparison `json:"openingDepreciation,omitempty"`
	// Year's depreciation: {ConceptPrefix}ForandringAvskrivningarAretsAvskrivningar
	YearDepreciation YearComparison `json:"yearDepreciation,omitempty"`
	// Closing: {ConceptPrefix}Avskrivningar @ balans0/balans1
	ClosingDepreciation YearComparison `json:"closingDepreciation,omitempty"`

	// Carrying value = Anskaffningsvärden - Avskrivningar
	// This is the main line item concept, e.g. se-gen-base:ByggnaderMark
	CarryingValue YearComparison `json:"carryingValue"`
}

// LongTermLiabilitiesNoteData represents note 7 (långfristiga skulder > 5 år).
type LongTermLiabilitiesNoteData struct {
	NoteNumber int `json:"noteNumber"` // typically 7

	// se-gen-base:LangfristigaSkulderForfallerSenare5Ar
	DueAfterFiveYears YearComparison `json:"dueAfterFiveYears"`
}

// PledgesNote represents note 8 (ställda säkerheter).
type PledgesNote struct {
	NoteNumber int `json:"noteNumber"` // typically 8

	// se-gen-base:StalldaSakerheterForetagsinteckningar
	CorporateMortgages YearComparison `json:"corporateMortgages,omitempty"`
	// se-gen-base:StalldaSakerheterFastighetsinteckningar
	RealEstateMortgages YearComparison `json:"realEstateMortgages,omitempty"`
	// se-gen-base:StalldaSakerheter
	TotalPledges YearComparison `json:"totalPledges"`
}

// ContingentLiabilitiesNote represents note 9 (eventualförpliktelser).
type ContingentLiabilitiesNote struct {
	NoteNumber int `json:"noteNumber"` // typically 9

	// se-gen-base:EventualForpliktelser
	TotalContingent YearComparison `json:"totalContingent"`
}

// MultiPostNote represents note 10 (tillgångar, avsättningar och skulder som avser flera poster).
type MultiPostNote struct {
	NoteNumber int `json:"noteNumber"` // typically 10

	// se-gen-base:NotTillgangarAvsattningarSkulderAvserFleraPoster
	Description string `json:"description"`

	// Each entry is a tuple: se-gen-base:TillgangarAvsattningarSkulderTuple
	Entries []MultiPostEntry `json:"entries"`
}

// MultiPostEntry is one tuple in the multi-post note.
type MultiPostEntry struct {
	// Heading for grouping, e.g. "Långfristiga skulder" or "Kortfristiga skulder"
	Heading string `json:"heading"`
	// se-gen-base:TillgangarAvsattningarSkulderPost
	PostName string `json:"postName"`
	// se-gen-base:TillgangarAvsattningarSkulderBelopp
	Amount *int64 `json:"amount"`
}

// Signatures represents the underskrifter section.
type Signatures struct {
	// se-gen-base:UndertecknandeArsredovisningOrt
	City string `json:"city"`
	// se-gen-base:UndertecknandeArsredovisningDatum
	Date string `json:"date"`

	// Each signatory is a tuple: se-gen-base:UnderskriftArsredovisningForetradareTuple
	Signatories []Signatory `json:"signatories"`
}

// Signatory represents one person signing the annual report.
type Signatory struct {
	// se-gen-base:UnderskriftArsredovisningForetradareTilltalsnamn
	FirstName string `json:"firstName"`
	// se-gen-base:UnderskriftArsredovisningForetradareEfternamn
	LastName string `json:"lastName"`
	// se-gen-base:UnderskriftArsredovisningForetradareForetradarroll (optional)
	Role string `json:"role,omitempty"`
}

// Helper function to create an int64 pointer (useful for populating the model).
func Int64(v int64) *int64 {
	return &v
}

// Helper function to create a string pointer.
func String(v string) *string {
	return &v
}
