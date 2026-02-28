# tekniskt-ramverk-digital-inlamning-arsredovisning-1.2

## Sida 1

 
1 
 
Digital inlämning av 
årsredovisningar 
Tekniskt ramverk 
 
Version 1.2 
 


---

## Sida 2

 
2 
 
Innehållsförteckning 
1 Bakgrund och syfte ...................................................................................................................... 3 
1.1 
Tekniskt ramverk ............................................................................................................... 3 
1.2 
Definitioner ........................................................................................................................ 3 
2 Inledning ........................................................................................................................................ 3 
3 Säker kommunikation .................................................................................................................. 4 
4 Infrastruktur och aktörer ............................................................................................................ 4 
5 Tjänstebeskrivningar .................................................................................................................... 5 
6 Test innan produktion ................................................................................................................. 6 
7 Referenser ...................................................................................................................................... 6 
 
 
 


---

## Sida 3

 
3 
 
 
1 
Bakgrund och syfte 
Detta dokument beskriver det tekniska ramverket för infrastruktur samt vägledning för hur 
en mjukvaruleverantör ansluter till och använder Bolagsverkets tjänster för digital 
inlämning av årsredovisning. Det ger en översiktlig beskrivning av systemlösningen för 
digital inlämning av årsredovisningar.  
 
Målgruppen för dokumentet är främst teknisk personal som ska arbeta med realiseringen 
av anslutningar till tjänsterna. 
1.1 Tekniskt ramverk 
Det tekniska ramverket omfattar förutom detta dokument också de dokument som är 
listade i kapitel 7 ”Referenser”. 
1.2 Definitioner 
Följande definitioner gäller i det här dokumentet: 
 
Term 
Beskrivning 
iXBRL 
Inline XBRL – ett XHTML-dokument som innehåller data som 
taggats upp enligt en XBRL-taxonomi 
XBRL 
eXtensible Business Reporting Language, en XML-standard för 
rapportering av olika typer av företagsinformation 
XHTML 
eXtensible HyperText Markup Language, en striktare variant av 
HTML. Dokumentformat som används för presentation i 
webbläsare 
XML 
eXtensible Markup Language, är en standard som utvecklats av W3C. 
XML är ett sätt att strukturera data. XML gör det lätt för en dator att 
generera data, läsa data och garantera att datastrukturer är entydiga. 
 
 
2 
Inledning  
Bolagsverkets tjänster för elektronisk inlämning av årsredovisningar används: 
 Vid upprättande av årsredovisning för att hämta grunduppgifter om aktiebolag från 
Bolagsverkets register 
 För att lämna in en elektronisk avskrift av årsredovisningen till Bolagsverket 
 Efter inlämning för att följa upp vad som händer med årsredovisningen. 
 För statistik över hur tjänsterna används. 
 
Hur säker anslutning mot tjänsterna görs beskrivs i kapitel 3. Infrastrukturen och aktörerna 
beskrivs i kapitel 4. Hur tjänster hos konsument, förmedlare och producent typiskt 
samverkar beskrivs med hjälp av sekvensdiagram i referens 1 ”Teknisk guide digital 
inlämning av årsredovisning”. 
 
Anslutningsanvisning för åtkomst till test- och produktionsmiljöer finns i 
referens 3 ”Anslutningsanvisning digital inlämning av årsredovisning”. 
 


---

## Sida 4

 
4 
 
 
 
3 
Säker kommunikation 
 
Kommunikation mellan ingående parter i infrastrukturen krypteras med hjälp av HTTPS 
med domäncertifikat utfärdat av en betrodd certifikatutfärdare. 
Auktorisering sker genom att organisationsnumret i konsumentens anslutningsavtal 
matchas med uppgiften i organisationscertifikatet. En anslutande klient måste: 
 Begära en brandväggsöppning mot Bolagsverkets servermiljöer 
 Anropa tjänsterna med TLS med ett godkänt klientcertifikat 
 Meddela Bolagsverket vilket serienummer certifikatet innehåller 
 
Detaljer och vidare anvisningar finns i referens 3 ”Anslutningsanvisning digital inlämning 
av årsredovisning”. 
 
 
4 
Infrastruktur och aktörer 
Följande bild ger en översikt av de aktörer och tjänster som samverkar för användare ska 
kunna lämna in sina årsredovisningar digitalt: 
 
SLUTANVÄNDARE
PROGRAMVARU-
LEVERANTÖRER
BOLAGSVERKET
Leverantör 1
Leverantör 2
Leverantör n
Desktop-
program
App
Webb-
applikation
...
Grunduppgifter
Ärendestatus
Mottagning 
elektronisk 
årsredovisning
Årsredovisnings-
händelser
Signering
...
 
Figur 1: Översikt infrastruktur  
Slutanvändare använder systemen för att upprätta en digital avskrift av årsredovisningen, 
för att lämna in den för underskrift samt för att granska och skriva under avskriften. 
 


---

## Sida 5

 
5 
 
Leverantörer tillhandahåller programvara för att skapa en digital avskrift av 
årsredovisningen i iXBRL-format, lämna in avskriften till Bolagsverket för underskrift samt 
använder Bolagsverkets informationstjänster för att underlätta upprättande och för att 
informera slutanvändare om ärendestatus och händelser. 
 
Bolagsverket tillhandahåller informationstjänster, tar emot och lagrar digitala avskrifter av 
årsredovisningar och tillhandahåller e-tjänst för granskning och underskrift av digitala 
avskrifter av årsredovisningar. 
5 
Tjänstebeskrivningar 
Tjänsterna är indelade i fyra grupper: 
 Informationstjänster 
 Tjänster för inlämning 
 Tjänster för årsredovisningshändelser 
 Tjänster för årsredovisningsstatistik 
 
Informationstjänsterna tillhandahåller information om  
 Grunduppgifter: namn, räkenskapsperiod, representanter mm. 
 Ärendestatus, dvs. status för senast inlämnade årsredovisning 
 
Tjänsterna för inlämning utgörs av  
 Skapa token för inlämning 
 Kontrollera 
 Lämna in 
 
Tjänsterna för årsredovisningshändelser utgörs av  
 Skapa prenumeration på årsredovisningshändelser för företag 
 Ta bort prenumerationer 
 Hämta alla årsredovisningshändelser för ett företag 
 
Tjänsterna för årsredovisningsstatistik utgörs av  
 Allmän statistik för inlämnade årsredovisningar 
 Statistik per programvara 
 Statistik för digitalt inlämnade årsredovisningar med separat revisionsberättelse 
 Statistik för digitalt inlämnade årsredovisningar som gett utfall i de automatiska 
kontrollerna 
 Stödtjänst för att visa vilka programvaror en programvaruleverantör har 
 
 
För utförliga beskrivningar av tjänsterna, se referens 4 ”Servicespecifikationer digital 
inlämning av årsredovisning” och referens 1 ”Teknisk guide digital inlämning av 
årsredovisning”. 
 


---

## Sida 6

 
6 
 
6 
Test innan produktion 
Bolagsverket tillhandahåller en testmiljö som får användas av Leverantörer för att utveckla 
och felsöka systemlösningar för digital inlämning av årsredovisningar. Anslutning till 
testmiljöer beskrivs i referens 3 ”Anslutningsanvisning digital inlämning av årsredovisning”. 
 
7 
Referenser 
I detta dokument refereras till följande teknisk dokumentation, som finns tillgängliga på 
Bolagsverkets webbplats Digital inlämning av årsredovisning: 
 
1. Teknisk guide digital inlämning av årsredovisning 
 
2. Tillämpningsanvisning årsredovisning iXBRL 
 
3. Anslutningsanvisning digital inlämning av årsredovisning  
 
4. Servicespecifikationer digital inlämning av årsredovisning 
 
 
Information om taxonomier för årsredovisning publiceras på den 
organisationsgemensamma webbplatsen taxonomier.se.  
