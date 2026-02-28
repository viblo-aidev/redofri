# tillampningsanvisningar-arsredovisning-ixbrl-1-8

## Sida 1

 
 
Tillämpningsanvisning
för årsredovisningar i iXBRL-format
Version 1.8
2023-11-10 
 


---

## Sida 2

 
2 
 
Innehållsförteckning 
1 Inledning ..................................................................................................................................................... 4 
1.1 
Terminologi .................................................................................................................................... 4 
1.2 
Indelning ......................................................................................................................................... 4 
1.3 
Exempel och referenser ............................................................................................................... 4 
1.4 
Om denna version ......................................................................................................................... 4 
2 Anvisningar rörande taxonomi och märkning av data ........................................................................ 5 
2.1 
Grundprinciper för märkning av data ........................................................................................ 5 
2.2 
Val av taxonomier ......................................................................................................................... 5 
2.3 
Val av taxonomielement............................................................................................................... 5 
2.4 
Användning av kontext vid inlämning av årsredovisning ....................................................... 5 
2.5 
Användning av kontext vid inlämning av års- och koncernredovisning .............................. 6 
2.6 
Inlämning av separat revisionsberättelse ................................................................................... 7 
2.7 
Flera förekomster av samma data i ett dokument ................................................................... 7 
2.8 
Obligatoriskt taggade data för räkenskapsår avslutat före 2021-12-31 ................................. 7 
2.9 
Obligatoriskt taggade data för räkenskapsår avslutat 2021-12-31 och senare ..................... 8 
2.10 Om märkning av belopp .............................................................................................................. 8 
2.11 Om märkning av datum ............................................................................................................... 9 
2.12 Om märkning av procentenheter ............................................................................................... 9 
2.13 Om märkning av andelar ............................................................................................................. 9 
2.14 Om märkning av antal .................................................................................................................. 9 
2.15 Om märkning av vallistor ..........................................................................................................10 
2.16 Namnkonventioner för enhets- och kontextreferenser mm ................................................10 
2.17 Om namnrymder och schemareferenser .................................................................................11 
2.18 Användning av extension-taxonomier och dimensioner ......................................................11 
2.19 Om märkning av odefinierade begrepp ...................................................................................11 
2.20 Om märkning av ändrade rubriker ...........................................................................................11 
2.21 Om märkning av notkopplingar i årsredovisning ..................................................................11 
2.22 Status för taggning ......................................................................................................................11 
3 Anvisningar rörande XHTML och iXBRL .........................................................................................12 
3.1 
Version iXBRL ............................................................................................................................12 
3.2 
Dokumentformat: XHTML, UTF-8 ........................................................................................12 
3.3 
Ett instansdokument per fil .......................................................................................................12 
3.4 
Script mm .....................................................................................................................................12 
3.5 
Bilder .............................................................................................................................................12 
3.6 
Länkar och andra externa referenser ........................................................................................13 
3.7 
Stylesheets ....................................................................................................................................13 
3.8 
Head-element ...............................................................................................................................13 
3.9 
Dolda element..............................................................................................................................14 
3.10 Formatmallar, sidnumrering m.m. ............................................................................................14 
3.11 Typsnitt .........................................................................................................................................14 
4 Övriga anvisningar ..................................................................................................................................15 
4.1 
Överensstämmelse mellan text och data .................................................................................15 
4.2 
Storleksbegränsningar .................................................................................................................15 
4.3 
Information om mjukvara som upprättat dokumentet mm .................................................15 
4.4 
Datum för undertecknande av fastställelseintyg ....................................................................16 
4.5 
Kontrollsumma............................................................................................................................17 


---

## Sida 3

 
3 
 
5 Referenser.................................................................................................................................................18 
 
Ändringshistorik 
Version 
Datum 
Beskrivning 
1.0 
2018-02-28
Första version. 
Fontdefinitioner är tillåtna (base64-kodade). 
1.1 
2019-02-07
Tydliggörande och tillägg inför uppdatering av 
tjänst med stöd för årsredovisningar upprättade 
enligt stöd kommenterar med eventuella ändringar 
inför K3. 
1.2 
2019-10-22
Tillägg av anvisning för separat revisionsberättelse, 
rubrik 2.2. 
1.3 
2021-01-27
Uppdaterad med anvisningar för års- och 
koncernredovisning enligt K3-regelverket. 
1.3 
2021-01-28
Uppdaterat med information kring hantering av 
segment samt obligatoriska fält. Har även bytt 
ordningen på vissa punkter under kap 2. 
1.4 
2021-06-18
Uppdatering med väsentliga begrepp för maskinell 
handläggning under 2.7 
1.5 
2021-12-21
Uppdatering med Bokföringsnämndens 
uppdaterade normgivning kring datum för 
undertecknande av årsredovisning. Uppdatering av 
punkt 2.7 samt tillägg av ny punkt 2.8. 
1.6 
2022-06-10
Tillägg för att tagga kontrollsumma på 
dokumentet. 
1.7 
2023-05-09
Uppdatering för att exkludera kontrollsumma på 
underskriftsdatum samt underskrifter i 
revisionsberättelse. 
1.8 
2023-11-10
Tillägg av 2.1 Grundprinciper för märkning av data 
samt 3.9.4 med tydliggöraren kring dolda begrepp. 
 
 


---

## Sida 4

 
4 
 
1 
Inledning  
1.1 Terminologi 
Detta dokument innehåller tillämpningsanvisningar för upprättande av årsredovisningar i 
iXBRL-format. Dokumentet följer konventionen i RFC 2119 beträffande olika typer av 
krav. 
 MÅSTE respektive FÅR INTE innebär att kravet måste uppfyllas för att 
årsredovisningen ska tas emot av Bolagsverket. 
 BÖR respektive BÖR INTE innebär att det finns situationer då kravet inte är 
tillämpligt, men kravet måste analyseras noggrant innan man beslutar att inte följa 
det. 
 KAN eller FÅR innebär att kravet är valfritt. 
1.2 Indelning 
Tillämpningsanvisningen är indelad i följande avsnitt. 
 Taxonomirelaterade anvisningar: kompletteringar till taxonomin för 
årsredovisningar, förtydliganden och rekommendationer kring användandet av 
xbrl-konstruktioner som labels, dimensions etc. 
 HTML-relaterade anvisningar: rekommendationer och krav på HTML-elementen 
i årsredovisningsfilen. 
 Övriga anvisningar. 
1.3 Exempel och referenser 
Referenser har samlats i ett eget avsnitt sist i detta dokument. Exempel på några av de 
olika konstruktionerna i dokumentet hittas i respektive avsnitt. Kodexempel i form av 
fullständiga och specifika iXBRL-filer finns publicerade på taxonomier.se.   
1.4 Om denna version 
Tillämpningsanvisningen är i sina huvuddrag färdig.  Dokumentet kommer att löpande 
uppdateras med fler detaljer och exempel för att förtydliga tillämpningen. Nya revisioner 
av dokumentet kommer att vara bakåtkompatibla med tidigare tillämpningsanvisningar. 
 
 


---

## Sida 5

 
5 
 
2 
Anvisningar rörande taxonomi och märkning av data 
De svenska taxonomierna för årsredovisning, revisionsberättelse, fastställelseintyg och 
utökad information m.m. beskrivs utförligt på taxonomier.se. Där finns även exempel.  
Information om vilka versioner som accepteras i Bolagsverkets inlämningstjänst hittas på 
bolagsverket.se. 
2.1 Grundprinciper för märkning av data 
Huvudsakligen MÅSTE märkning av data ske på information som är visuellt synliga i 
Inline XBRL filen. I undantagsfall FÅR märkt data separeras från den visuella 
presentationen och döljas.  
  
Grundprincipen är att huvuddelen av finansiellt visuellt data MÅSTE vara märkt med 
begrepp från vald taxonomi. Endast märkning av data som exempelvis meta-data om 
handlingen eller duplicerade rubriker KAN döljas.  
2.2 Val av taxonomier 
2.2.1 
Angivna versioner av de svenska taxonomierna för årsredovisning MÅSTE 
användas för taggning av data i dokumentet. Se bilaga över tillåtna kombinationer av 
rapporter. 
2.3 Val av taxonomielement 
2.3.1 
Data MÅSTE taggas med det element som bäst motsvarar datat.  
Grundprincipen är att all information ska taggas med i taxonomierna kategoriserade 
begrepp. Om osäkerhet finns rörande definitionen BÖR denna information taggas med 
strukturen för ”Odefinierade begrepp” i ”Utökad information” taxonomin.  
 
Information som är av vikt för handlingen dvs. redovisnings- samt revisionsinformation 
exklusive exempelvis sidhuvud, sidfot etc. SKA taggas. Flaggan för ”Status för taggning” 
BÖR alltid tillämpas för att underlätta vid maskinell bearbetning/tolkning.  
 
2.3.2 
Om data kan taggas med flera möjliga taggar på olika nivåer i taxonomin MÅSTE 
den tag användas som bäst motsvarar datats omfattning. Följande principer MÅSTE 
användas vid val. 
 Om datat avser en summering på högre nivå MÅSTE motsvarande tag för 
summerat data väljas. 
 Om datat avser mer detaljerad information MÅSTE den mer detaljerade taggen 
väljas. 
2.4 Användning av kontext vid inlämning av årsredovisning 
2.4.1 
Om inlämning avser endast årsredovisning FÅR INTE fastställelseintyg, 
revisionsberättelse samt årsredovisning använda kontext som innehåller varken segment 
eller scenario. 


---

## Sida 6

 
6 
 
2.5 Användning av kontext vid inlämning av års- och koncernredovisning 
2.5.1 
Om inlämning avser års- och koncernredovisning MÅSTE fastställelseintyg, 
revisionsberättelse samt års- och koncernredovisning använda segment för samtliga 
kontext. 
2.5.2 
Taggat data för samtliga begrepp och handlingar i en års- och koncernredovisning 
MÅSTE referera till kontext som använder någon av följande segment:  
 RedovisningInformationKoncernSegment 
 RedovisningInformationJuridiskPersonSegment 
 RedovisningInformationGenerellSegment 
2.5.3 
Numerisk information i rapport för årsredovisning avseende juridisk person 
MÅSTE taggas med kontext som använder segment 
RedovisningInformationJuridiskPersonSegment. 
2.5.4 
Numerisk information i rapport för koncernredovisning avseende koncernen 
MÅSTE taggas med kontext som använder segment RedovisningInformationKoncernSegment. 
2.5.5 
Textuell information BÖR normalt taggas med kontext som använder segment 
RedovisningInformationGenerellSegment. 
2.5.6 
Samtliga data i årsredovisningens delrapport ”Allmän information” förutom 
”Avgivande av finansiell rapport (Vallista)” MÅSTE taggas med kontext som använder 
segment RedovisningInformationJuridiskPersonSegment. 
2.5.7 
Samtliga data i koncernredovisningens delrapport ”Allmän information” exklusive 
”Avgivande av finansiell rapport (Vallista)” MÅSTE endast kontext med segment 
RedovisningInformationKoncernSegment användas. 
2.5.8 
Begreppet ”Avgivande av finansiell rapport (Vallista)” MÅSTE taggas med 
kontext som använder segment RedovisningInformationGenerellSegment. 
2.5.9 
Begreppen "Organisationsnummer", "Räkenskapsårets första dag" och 
"Räkenskapsårets sista dag" MÅSTE taggas med samtliga tre kontexter som använder 
segment RedovisningInformationKoncernSegment, RedovisningInformationJuridiskPersonSegment och 
RedovisningInformationGenerellSegment. 
2.5.10 Samtlig information i delrapport "Undertecknande av företrädare och 
revisionspåteckning" MÅSTE endast taggas med kontext som använder segment 
RedovisningInformationGenerellSegment. 
2.5.11 Samtlig information i en tuple (rad) MÅSTE taggas med kontext som använder 
samma segment. Om tuple innehåller både vallista/text och numerisk information 
MÅSTE segment för det numeriska värdet användas. 
2.5.12 Samtlig information i ”Revisionsberättelse” MÅSTE taggas med kontext som 
använder segment RedovisningInformationGenerellSegment. 


---

## Sida 7

 
7 
 
2.5.13 Samtlig information i ”Fastställelseintyg” MÅSTE taggas med kontext som 
använder segment RedovisningInformationGenerellSegment. 
2.5.14 Om märkning av odefinierade begrepp används BÖR kontext med segment som 
bäst beskriver tillhörighet tillämpas. 
2.6 Inlämning av separat revisionsberättelse 
2.6.1 
Om inlämning endast avser årsredovisning och har separat årsredovisning och 
revisionsberättelse: 
 FÅR INTE årsredovisningen innehålla referens till revisionsberättelse 
 FÅR INTE revisionsberättelsen innehålla referens till årsredovisning och/eller 
fastställelseintyg 
 Revisionsberättelsen MÅSTE endast innehålla ett (1) kontext för 
redovisningsperiod (duration). 
o Revisionsberättelsen FÅR INTE innehålla flera redovisningsperioder. 
 Revisionsberättelsen FÅR INTE innehålla kontext för balansdag (instant).  
 
2.6.2 
Om inlämning avser års- och koncernredovisning med separat revisionsberättelse 
MÅSTE information taggas enligt ovanstående med kontext som använder segment 
RedovisningInformationGenerellSegment. 
2.7 Flera förekomster av samma data i ett dokument 
2.7.1 
Om ett dokument innehåller flera taggar med samma namn och kontext MÅSTE 
datat i taggarna vara identiskt. Annars går det inte att avgöra vilket data som är det riktiga. 
2.7.2 
Om ett dokument innehåller samma data på flera ställen och dessa data kan 
taggas, MÅSTE datat taggas på ALLA ställen där det förekommer. 
2.7.3 
Det förekommer att samma data presenteras på olika sätt med hjälp av attributen 
scale och decimals – t.ex. anges belopp i ental på ett ställe i årsredovisningen och i 
tusental på ett annat ställe. I dessa fall FÅR samma tagg förekomma med olika värden, 
men då MÅSTE scale och decimals sättas på ett sådant sätt att värdena motsvarar 
varandra. 
2.8 Obligatoriskt taggade data för räkenskapsår avslutat före 2021-12-31 
2.8.1 
För års- och/eller koncernredovisning i iXBRL-format gäller att följande data 
taggas. 
 Samtlig information i ”Allmän information” MÅSTE vara taggat.  
o För K3 gäller att ”Företagets tidigare namn” endast BÖR taggas om 
namnet förändrats under räkenskapsåret. 
 Det MÅSTE finnas taggad information i ”Årsredovisning”  
 ”Fastställelseintyg” MÅSTE innehålla taggad data för begreppen 
”ArsstammaIntygande” och ”IntygandeOriginalInnehallType”. 


---

## Sida 8

 
8 
 
o För innehåll i ”ArsstammaIntygande” MÅSTE viktig information taggas 
med taxonomins underliggande taggar. Om taxonomin inte innehåller 
lämplig tagg KAN informationen lämnas otaggad. 
 Om ”Revisionsberättelse” är upprättad MÅSTE det finnas taggad information. 
 
2.8.2 
För möjliggörande av effektiv maskinell handläggning BÖR följande data taggas. 
 I årsredovisningen BÖR följande begrepp taggas: ”Förslag till utdelning”, ”Datum 
för avgivande av årsredovisning” och ”Revisionsberättelse utan modifierade 
uttalande” eller ”Revisionsberättelse med modifierade uttalanden”. 
 I revisionsberättelsen BÖR ”Datum för revisionens avslutande” taggas. 
 
2.9 Obligatoriskt taggade data för räkenskapsår avslutat 2021-12-31 och 
senare 
2.9.1 
För års- och/eller koncernredovisning i iXBRL-format gäller att följande data 
taggas. 
 Samtlig information i ”Allmän information” MÅSTE vara taggat.  
o För K3 gäller att ”Företagets tidigare namn” endast BÖR taggas om 
namnet förändrats under räkenskapsåret. 
 Det MÅSTE finnas taggad information i ”Årsredovisning”  
 ”Fastställelseintyg” MÅSTE innehålla taggad data för begreppen 
”ArsstammaIntygande” och ”IntygandeOriginalInnehallType”. 
o För innehåll i ”ArsstammaIntygande” MÅSTE viktig information taggas 
med taxonomins underliggande taggar. Om taxonomin inte innehåller 
lämplig tagg KAN informationen lämnas otaggad. 
 Om ”Revisionsberättelse” är upprättad MÅSTE det finnas taggad information. 
 I tuple ” Undertecknande av företrädare (Tabell)” MÅSTE begreppet ”Datum för 
undertecknande” taggas för respektive företrädare. 
 
2.9.2 
För möjliggörande av effektiv maskinell handläggning BÖR följande data taggas. 
 I årsredovisningen BÖR följande begrepp taggas: ”Förslag till utdelning”, 
”Revisionsberättelse utan modifierade uttalande” eller ”Revisionsberättelse med 
modifierade uttalanden”. 
 I revisionsberättelsen BÖR ”Datum för revisionens avslutande” taggas. 
 
2.9.3 
För möjliggörande av effektiv maskinell handläggning FÅR INTE följande data 
taggas. 
 I årsredovisningen FÅR INTE begreppet ”Datum för avgivande av 
årsredovisning” taggas. 
2.10 Om märkning av belopp 
2.10.1 Attributet ”decimals” MÅSTE användas vid märkning av belopp som inte är 
heltal. Attributet ”precision” FÅR INTE användas. 


---

## Sida 9

 
9 
 
2.10.2 Om värden redovisas i heltals kronor eller euro MÅSTE attributet ”decimals” 
sättas till ”0” alternativt ”INF” och attributet ”scale” sätts till ”0”. 
2.10.3 Om värden redovisas i tusentals kronor eller euro MÅSTE attributet ”decimals” 
sättas till ”-3” och attributet ”scale” sätts till ”3”. 
2.10.4 Om värden redovisas i hundratals euro MÅSTE attributet ”decimals” sättas till ”-
2” samt attributet ”scale” sätts till ”2”. 
2.10.5 Om värden redovisas i miljontals kronor eller euro MÅSTE attributet ”decimals” 
sättas till ”-6” samt attributet ”scale” sätts till ”6”. 
2.10.6 I de fall beloppsvärden har ett tecken som avviker från det normala (t.ex. en 
debetpost som är negativ) MÅSTE attributet ”sign” användas. Attributet ”sign” MÅSTE 
också användas för beloppsvärden som inte är klassade som debet eller kredit. 
2.11 Om märkning av datum 
2.11.1 Ett av de två följande datumformat BÖR användas vid taggning av datum.  
 Datum med formatet YYYY-MM-DD som kan representera exempelvis datumet 
2017-12-31. 
 Datum med formatet (D)D mon(th) YYYY som kan representerar exempelvis 
datumen 1 jan 2017 eller 31 december 2017. I dessa fall sätts attributet ”format” 
till ”ixt3:datedaymonthyeardk”. 
2.12 Om märkning av procentenheter 
2.12.1 Data som ska anges i procentenheter MÅSTE använda <ix:nonFraction> som 
elementtyp och tillämpa datatypen ” xbrli:pure”. Värdet anges i procent, t.ex. ”35,5” eller 
”100”. Attributet ”scale” sätts till ”–2” för att indikera att datavärdet är två 
decimalpositioner (alltså 100 gånger) mindre än det skrivna värdet.  
2.13 Om märkning av andelar 
2.13.1 Data som ska anges i andelar MÅSTE använda <ix:nonFraction> som 
elementtyp och tillämpa datatypen ” xbrli:shares”. 
2.14 Om märkning av antal 
2.14.1 Data som ska hantera antal MÅSTE använda <ix:nonFraction> som 
elementtyp och ha en anpassad datatyp för tillämpningen.  
 För exempelvis ”Medelantalet anställda” i K2 BÖR datatypen ”se-k2-
type:AntalAnstallda” definieras med följande namnrymd ’xmlns:se-k2-
type="http://www.taxonomier.se/se/fr/k2/datatype"’. 
 För exempelvis ”Medelantalet anställda” i K3 BÖR datatypen ”se-k3-
type:AntalAnstallda” definieras med följande namnrymd ’xmlns:se-k3-
type="http://www.taxonomier.se/se/fr/k3/datatype"’. 


---

## Sida 10

 
10 
 
2.15 Om märkning av vallistor 
2.15.1 Data för vallistor baserade på Extensible Enumerations MÅSTE taggas med 
<ix:nonNumeric> enligt följande exempel med handlingens språk:  
 
<ix:nonNumeric name="se-cd-base:SprakHandlingUpprattadList" 
contextRef="period0_jur">se-mem-
base:SprakSvenskaMember</ix:nonNumeric>. 
2.16 Namnkonventioner för enhets- och kontextreferenser mm 
2.16.1 Monetära enheter som exempelvis svenska kronor, brittiska pund etc. BÖR 
namnges med valutakod enligt ISO-4217, dvs. SEK, EUR osv. 
2.16.2 Procentenheter BÖR namnges ”procent”. 
2.16.3 Andelar BÖR namnges ”andelar”. 
2.16.4 Antal anställda BÖR namnges ”antal-anstallda”. 
2.16.5 Redovisningsperioder BÖR namnges ”period0” för den redovisningsperiod som 
årsredovisningen avser, ”period1” för föregående redovisningsperiod osv. 
2.16.6 Redovisningsperioder i handlingar innehållande både års- och koncernredovisning 
BÖR använda följande namnstandard avseende redovisningsperiod.  
 Redovisningsperiod avseende juridisk person BÖR namnges med suffix "_jur" 
som exempelvis "period0_jur" för den senaste redovisningsperioden, 
"period1_jur" för den föregående osv. 
 Redovisningsperiod avseende koncern BÖR namnges med suffix "_kon" som 
exempelvis "period0_kon" för den senaste redovisningsperioden, " period1_kon" 
för den föregående osv. 
 Redovisningsperiod avseende gemensam information BÖR namnges med suffix 
"_gem" som exempelvis "period0_gem" för den senaste redovisningsperioden, 
"period1_gem" för den föregående osv. 
 
2.16.7 Balansdagar BÖR namnges ”balans0” för den senaste redovisningsperiodens 
balansdag, ”balans1” för den föregående osv. 
2.16.8 Balansdagar i handlingar innehållande både års- och koncernredovisning BÖR 
använda följande namnstandard avseende redovisningsperiodens balansdag.  
 Balansdagar avseende juridisk person BÖR namnges med suffix "_jur" som 
exempelvis "balans0_jur" för den senaste redovisningsperiodens balansdag, 
"balans1_jur" för den föregående osv. 
 Balansdagar avseende koncern BÖR namnges med suffix "_kon" som exempelvis 
"balans0_kon" för den senaste redovisningsperiodens balansdag, "balans1_kon" 
för den föregående osv. 


---

## Sida 11

 
11 
 
 Balansdagar avseende gemensam information BÖR namnges med suffix "_gem" 
som exempelvis "balans0_gem" för den senaste redovisningsperiodens balansdag, 
"balans1_gem" för den föregående osv. 
2.17 Om namnrymder och schemareferenser 
2.17.1 Namnrymder som inte tillämpas vid taggning av data BÖR inte ingå i 
dokumentet. 
2.17.2 Schemareferenser som inte tillämpas vid taggning av data BÖR inte ingå i 
dokumentet. 
2.18 Användning av extension-taxonomier och dimensioner 
2.18.1 Extension-taxonomier FÅR INTE användas för taggning av data i 
årsredovisningsdokument. 
2.18.2 Dimensioner FÅR INTE heller användas – det följer av att inga dimensioner är 
definierade i taxonomierna, och de får inte utökas på egen hand. 
2.19 Om märkning av odefinierade begrepp 
Om taxonomierna saknar begrepp som motsvarar datat eller om omfattningen är större 
eller markant avviker från taxonomins definition BÖR informationen taggas med 
strukturen för ”Odefinierade begrepp” i taxonomin för ”Utökad information”.   
Informationen BÖR klassificeras baserat på om det är monetära, numeriska eller textuella 
värden. OM önskad datatyp saknas BÖR värden taggas med den textuella strukturen.  
2.20 Om märkning av ändrade rubriker 
Om rubriken i iXBRL presentationen på något sätt avviker från presentationsrubriken i 
taxonomin för elementet/begreppet BÖR den avvikelsen loggas med funktionalitet för 
”Loggning av rubrikändring” i taxonomin för ”Utökad information”.  
2.21 Om märkning av notkopplingar i årsredovisning 
Samtliga visuella notkopplingar i iXBRL för årsredovisningen BÖR taggas med 
funktionen i taxonomin för ”Utökad information”. 
2.22 Status för taggning 
För att underlätta maskinell bearbetning/tolkning BÖR flaggan för ”Status för taggning” 
tillämpas.  
 
Om det i ”Årsredovisning” finns information som inte taggats men som är av vikt för 
handlingen dvs. redovisningsinformation (inte exempelvis sidhuvud, sidfot etc.) sätts 
flaggan ArsredovisningEjTaggadInformation till ”true”. I annat fall BÖR flaggan sättas till 
”false”.  Det finns en motsvarande flagga för ”Revisionsberättelse” 
(RevisionsberattelseEjTaggadInformation) och den BÖR tillämpas. 
 


---

## Sida 12

 
12 
 
3 
Anvisningar rörande XHTML och iXBRL 
En huvudprincip är att varje iXBRL-dokument ska vara komplett i sig, oberoende av 
tillgång till externa resurser som script, stylesheets mm. Den enda avvikelsen från den 
principen är referenser till standardscheman som XHTML-schemat, taxonomischeman 
mm. 
3.1 Version iXBRL 
Version 1.1 av iXBRL MÅSTE tillämpas. 
3.2 Dokumentformat: XHTML, UTF-8 
3.2.1 
Dokument som lämnas in till Bolagsverket MÅSTE vara giltiga XHTML-
dokument (iXBRL-standarden tillåter att dokumenten är HTML, men dessa accepteras 
inte av Bolagsverket). 
3.2.2 
Dokumentet BÖR ha default namespace satt till 
'http://www.w3.org/1999/xhtml'. 
3.2.3 
Dokumentet MÅSTE kodas i UTF-8 och det BÖR sätta encoding till UTF-8: 
<?xml version="1.0" encoding="utf-8"?>. 
3.2.4 
Dokumentet BÖR innehålla MIME-typ enligt: 
<meta content='text/html; charset=utf-8' http-equiv='Content-Type' /> 
Numeriska XML-escapesekvenser (t.ex. &#8364; för €) FÅR användas.  
Vanliga HTML-escapesekvenser som &nbsp;, &auml; osv. är normalt sett inte giltig 
XHTML och de FÅR INTE användas, med undantag för de fem escapesekvenser som 
används för XML: &quot; &amp; &lt; &gt; &apos;  
3.3 Ett instansdokument per fil 
3.3.1 
Varje iXBRL-instansdokument MÅSTE förmedlas i en egen fil (iXBRL-
standarden tillåter dokumentset, iXDS, men dessa accepteras inte av Bolagsverket). 
3.4 Script mm 
3.4.1 
Instansdokument MÅSTE vara fria från script av alla slag. Det gäller såväl 
<script>-taggar som eventhanterare (onclick osv.) Skälet till det är att script kan 
påverka presentation och innehåll av dokumentet på ett sätt som är svårt eller omöjligt att 
reproducera vid ett senare tillfälle. 
3.4.2 
Instansdokument får inte heller innehålla applets, flash-animationer, JavaFX eller 
någon annan typ av exekverbara element. 
3.5 Bilder 
En av de främsta fördelarna med iXBRL är att formatet möjliggör anpassad och 
användarvänlig presentation av dokument. Bilder är ett viktigt sätt att förbättra och 
användaranpassa presentationen. 
 
3.5.1 
Bilder FÅR förekomma i instansdokument. 


---

## Sida 13

 
13 
 
3.5.2 
Bilder BÖR användas sparsamt, och de BÖR då använda så lite data som möjligt. 
Sträva efter att hålla nere datastorleken genom att reducera komplexa detaljer, färgdjup 
mm. Se 4.2.2 för maximal storlek på bilder. 
3.5.3 
Bilder MÅSTE använda ett av formatet JPEG, SVG, GIF eller PNG. 
3.5.4 
Om bilder förekommer MÅSTE de bäddas in helt i dokumentet mha base64-
kodning. <img>-taggar får inte peka på externa resurser, vare sig relativa eller absoluta. 
Skälet till det är att varje dokument ska kunna tolkas utan beroende till externa resurser. 
3.5.5 
Tabeller, grafer och andra grafiska representationer FÅR förekomma som bilder. 
Om sådana bilder förekommer MÅSTE motsvarande data taggas med inline XBRL-
taggar. Skälet till det är att allt data i ett dokument ska vara tillgängligt i taggad form. 
3.6 Länkar och andra externa referenser 
Grundregeln är att länkar endast får användas för referenser inom dokumentet, alltså mha 
hash-notation. Exempel: <a href=”#not8”>…</a>. 
 
3.6.1 
Länkar och andra referenser MÅSTE referera till andra element inom samma 
dokument (mha #-notation). 
3.6.2 
Externa referenser FÅR förekomma i schemadeklarationer och annat metadata. 
Om sådana referenser förekommer MÅSTE de peka på standardiserade resurser som 
grundschemat för XHTML, scheman för taxonomier mm. 
3.7 Stylesheets 
Stylesheets får användas fritt, men sträva efter enkelhet och samla all stilhantering till ett 
stycke i början av dokumentet. Undvik stilsättning på enskilda element. 
 
3.7.1 
Stylesheets FÅR användas i iXBRL-dokument. 
3.7.2 
Stylesheet-deklarationer BÖR samlas till en plats i början av dokumentet. 
3.7.3 
Stilinformation MÅSTE deklareras i dokumentet – referenser till externa 
stylesheets FÅR INTE användas. 
3.7.4 
Stilinformation BÖR INTE sättas på direkt på enskilda element. Sträva efter att 
använda id, name, class eller liknande för att koppla stil till element. 
3.7.5 
CSS3 FÅR användas, men vid användning BÖR komplicerade konstruktioner 
undvikas, då de kan leda till olika presentation i olika webbläsare. 
3.8 Head-element 
Dokumentets <head>-del ska innehålla viss information om dokumentet som inte 
beskrivs av taxonomin. 
 
3.8.1 
Dokumentet MÅSTE ha en <title>. Denna <title> BÖR väljas så att den är 
unikt utpekande för dokumentet, t.ex. ”Årsredovisning för Bolag AB, räkenskapsår 
2017”. 


---

## Sida 14

 
14 
 
3.8.2 
Dokumentet MÅSTE innehålla information om den programvara som skapat den. 
Detta beskrivs detaljerat i avsnitt 4.3. 
3.9 Dolda element 
Dolda element bör undvikas eftersom de kan leda till skillnader mellan data som tolkas 
maskinellt och information som uppfattas av läsare av presentationen av dokumentet. 
 
3.9.1 
Element som <hidden> och liknande konstruktioner där stilmallar används för 
att dölja element BÖR INTE användas. 
3.9.2 
Om element döljs av <hidden> och liknande konstruktioner MÅSTE upprättaren 
av dokumentet säkerställa att presentation och data stämmer överens i enlighet med 4.1.1. 
3.9.3 
Dold information BÖR om möjligt placeras innanför Inline XBRL filens 
<ix:hidden> tag. Exempel på information är taggat data för vallistor baserade på 
Extensible Enumerations.  
3.9.4 
Information FÅR INTE döljas genom användning av samma färg på text och 
bakgrund, minska storlek på typsnitt eller placera information utanför visuellt synfält i en 
webbläsare. 
3.10 Formatmallar, sidnumrering m.m. 
3.10.1 Dokumentets presentation BÖR stilsättas så att dokumentet blir lättläst vid 
utskrift. Sträva efter att inte sidbryta tabeller och huvudrubriker m.m.   
3.10.2 Dokumentet BÖR innehålla sidnummer. 
3.11 Typsnitt 
3.11.1 Dokumentet FÅR definiera egna typsnitt.  
3.11.2 Om dokumentet innehåller egna typsnittsdefinitioner MÅSTE dessa definitioner 
inkluderas i sin helhet i dokumentet mha base64-kodning.  
 
 


---

## Sida 15

 
15 
 
4 
Övriga anvisningar 
4.1 Överensstämmelse mellan text och data 
4.1.1 
Data, dvs. information som taggats med inline XBRL-taggar, MÅSTE stämma 
överens med övrig text. Ansvaret för överensstämmelse vilar alltid på upprättaren av 
dokumentet. 
4.2 Storleksbegränsningar 
4.2.1 
Storleken för hela årsredovisningsdokumentet MÅSTE understiga 5 MB.  
4.2.2 
Om dokumentet innehåller bilder MÅSTE storleken för varje bild understiga 1 
MB. 
4.3 Information om mjukvara som upprättat dokumentet mm 
4.3.1 
Vid upprättande MÅSTE den programvara som upprättat iXBRL-dokumentet 
lägga till följande metadata i dokumentets <head>-element: 
 Namn på programvaran: i meta-elementet <programvara> 
 
Version av programvaran: i meta-elementet <programversion> 
Exempel:  
<head> 
  <meta name="programvara" content="Superstar Reporter Deluxe 2000"/> 
  <meta name="programversion" content="1.2.4-b3402"/> 
  ... 
</head> 
 
4.3.2 
Vid upprättande av iXBRL-dokumentet innehållande årsredovisning med 
revisionsberättelse BÖR programvara lägga till följande metadata i dokumentets <head>-
element:  
 Namn på programvaran: i meta-elementet <programvara-revision> 
 Version av programvaran: i meta-elementet <programversion-revision> 
Exempel:  
<head> 
  <meta name="programvara-revision" content="Ultimate Accountant Turbo 3000"/> 
  <meta name="programversion-revision" content="4.0-rc4"/> 
  ... 
</head> 
 
4.3.3 
Vid sammanslagning av årsredovisning och revisionsberättelse BÖR den 
programvara som sammanställt iXBRL-dokumentet lägga till följande metadata i 
dokumentets <head>-element: 
 Namn på programvaran: i meta-elementet <programvara-sammanstallning> 
 Version av programvaran: i meta-elementet <programversion-
sammanstallning> 
Exempel:  
<head> 
  <meta name="programvara- sammanstallning" content="Superstar Reporter Deluxe 2000"/> 
  <meta name="programversion- sammanstallning" content="1.2.4-b3402"/> 
  ... 
</head> 


---

## Sida 16

 
16 
 
4.3.4 
Namn på programvara BÖR ha en namnstandard bestående av ”Leverantör av 
programvara” och ”Namn på programvara”. 
Exempel:  
 ”Leverantör av programvara” i exempel: Programföretag AB 
 ”Namn på programvara” i exempel: Superstar Reporter Deluxe 2000 
Namn på programvara: Programföretag AB - Superstar Reporter Deluxe 2000 
 
4.3.5 
Programversion för programvara BÖR ha en namnstandard bestående av 
”Huvudversion” och ”Revisionsnummer”.  
 Huvudversionen BÖR uppdateras vid större omfattande förändringar av 
programvaran  
 Revisionsnumret BÖR uppdateras vid mindre anpassningar eller rättningar av 
programvaran. 
 
Exempel:  
 ”Huvudversion” i exempel: 2019 
 ”Revisionsnummer” i exempel: 4-b3402 
Programversion för programvara: 2019.4-b3402 
 
4.3.6 
Information om redovisningsbyrå och revisionsfirma FÅR läggas till i följande 
metadata i dokumentets <head>-element: 
 Namn på upprättande organisation: i meta-elementet <upprattare> 
 Namn på reviderande organisation: i meta-elementet <reviderare> 
 
Exempel:  
<head> 
  <meta name="upprattare" content="Kanonbokföring AB"/> 
  <meta name="reviderare" content="Turborevision AB"/> 
  ... 
</head> 
 
4.3.7 
Övrig information som kan vara av intresse för de som använder 
årsredovisningen som informationskälla, t.ex. certifiering av upprättare mm, FÅR bifogas 
som metadata i dokumentets <head>-element. 
4.4 Datum för undertecknande av fastställelseintyg 
Eftersom undertecknande av fastställelseintyget sker efter att handlingen lämnas in till 
Bolagsverket så kan upprättaren inte alltid sätta datum för undertecknande – tidpunkten 
för undertecknande är inte känd.  
 
För att hantera det problemet kommer Bolagsverket att sätta datumet i elementet för 
undertecknande av fastställelseintyget. Upprättande programvara ska sätta innevarande 
datum, dvs. det datum då dokumentet skapades av programvaran. 
 


---

## Sida 17

 
17 
 
4.4.1 
Instansdokument som innehåller fastställelseintyg MÅSTE sätta innevarande 
datum, dvs. det datum då dokumentet skapades av programvaran, för undertecknande av 
fastställelseintyget. 
4.4.2 
Elementet som innehåller datum för underskrift av fastställelseintyget MÅSTE 
märkas med id="ID_DATUM_UNDERTECKNANDE_FASTSTALLELSEINTYG". 
Exempel:  
<ix:nonNumeric name="se-bol-base:UnderskriftFastallelseintygDatum" 
contextRef="BALANS0" id="ID_DATUM_UNDERTECKNANDE_FASTSTALLELSEINTYG">2017-10-
05</ix:nonNumeric> 
 
4.5 Kontrollsumma 
4.5.1 
Bolagsverkets API BÖR användas för att skapa kontrollsumma med algoritmen 
SHA-256 av filen.  
4.5.2 
Om kontrollsumman och dess algoritm bifogas i Inline XBRL dokumentet 
MÅSTE följande metadata-element användas i dokumentets <head>-element: 
 
För årsredovisning eller filer med både årsredovisning och revisionsberättelse 
läggs kontrollsumman i meta-elementet <ixbrl.innehall.kontrollsumman>. 
 För algoritm som använts för ovanstående kontrollsumma läggs den i meta-
elementet <ixbrl.innehall.kontrollsumman.algoritm>. 
 För separat revisionsberättelse läggs kontrollsumman i meta-elementet 
<ixbrl.innehall.kontrollsumman.revision>. 
 För algoritm som använts för ovanstående kontrollsumma läggs den i meta-
elementet <ixbrl.innehall.kontrollsumman.revision.algoritm>. 
 
4.5.3 
Då uppgifter i fastställelseintyget, revisorspåteckningen och underskrifter i 
revisionsberättelsen behöver ändras efter att kontrollsumman är skapad MÅSTE dessa 
informationsmängder omslutas med en html-tagg med följande id-attribut för att 
exkluderas från kontrollsumman: 
 
id=”id-innehall-faststallelseintyg” för fastställelseintyg 
 
id=”id-innehall-revisorspateckning” för revisorspåteckning 
 
id=”id-innehall-underskrifter-revisionsberattelse” för underskrifter 
i revisionsberättelse. 
 
Underskriftsdatum i årsredovisningen som använder begreppet UndertecknandeDatum 
kommer automatiskt att exkluderas från kontrollsumman. 
 
4.5.4 
Om kontrollsumman presenteras visuellt MÅSTE följande id-attribut användas 
för att exkludera taggen vid uträkning av kontrollsumman: 
 För årsredovisning eller filer med både årsredovisning och revisionsberättelse 
används följande id på taggen: id-innehall-kontrollsumma 
 För separat revisionsberättelse används följande id på taggen: id-innehall-
kontrollsumma-revision 
 


---

## Sida 18

 
18 
 
5 
Referenser 
XBRL-standarden: https://www.xbrl.org/  
Av särskilt intresse för utvecklare är:  
https://www.xbrl.org/the-standard/how/getting-started-for-developers/ 
 
iXBRL-specifikationen: http://www.xbrl.org/specification/inlinexbrl-part1/rec-2013-11-
18/inlinexbrl-part1-rec-2013-11-18.html 
 
Information om de svenska taxonomierna: http://taxonomier.se/ 
 
Den brittiska XBRL-vägledningen för upprättare och utvecklare (innehåller 
förtydliganden och enklare exempel): 
http://www.xbrl.org.uk/documents/XBRL%20UK%20Preparers%20and%20Developer
s%20Guide-2010-03-31.pdf 
 
Brittisk styleguide för iXBRL: 
http://webarchive.nationalarchives.gov.uk/20140206171140/http://www.hmrc.gov.uk/e
bu/ct_techpack/xbrl-style-guide.pdf 
 
Lista över escapesekvenser för XML och HTML, se avsnitt Predefined Entities in XML 
för tillåtna escapesekvenser i iXBRL: 
https://en.wikipedia.org/wiki/List_of_XML_and_HTML_character_entity_references 
 
