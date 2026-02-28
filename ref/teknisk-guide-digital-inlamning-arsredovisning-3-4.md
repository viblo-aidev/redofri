# teknisk-guide-digital-inlamning-arsredovisning-3-4

## Sida 1

Digital inlämning av 
årsredovisningar
Teknisk guide
Version 3.4
 


---

## Sida 2

 
2 
 
Innehållsförteckning 
1 Ändringshistorik ........................................................................................................................ 3 
2 Inledning ..................................................................................................................................... 5 
3 Infrastruktur och aktörer ......................................................................................................... 6 
4 Beskrivning av tjänsterna ......................................................................................................... 8 
4.1 
Informationstjänster....................................................................................................... 8 
4.2 
Tjänster för inlämning ................................................................................................... 9 
4.3 
Tjänster för årsredovisningshändelser ....................................................................... 12 
4.4 
Tjänster för årsredovisningsstatistik .......................................................................... 14 
5 Teknisk beskrivning av tjänsterna ......................................................................................... 15 
5.1 
Servicespecifikationer .................................................................................................. 15 
5.2 
Informationstjänster..................................................................................................... 15 
5.3 
Tjänster för inlämning ................................................................................................. 19 
5.4 
Tjänster för årsredovisningshändelser ....................................................................... 20 
5.5 
Tjänster för årsredovisningsstatistik .......................................................................... 23 
6 Appendix A. Felkoder för API:t ........................................................................................... 25 
6.1 
Om statuskoder i REST-tjänsterna ............................................................................ 25 
6.2 
Lista över felkoder ........................................................................................................ 25 
7 Appendix B. Kodgenerering av klienter för REST-API:t ................................................. 27 
7.1 
Kodgenerering mha Swagger Editor ......................................................................... 27 
8 Appendix C. Trafikexempel ................................................................................................... 28 
8.1 
Informationstjänster..................................................................................................... 28 
8.2 
Tjänster för inlämning ................................................................................................. 30 
8.3 
Tjänster för årsredovisningshändelser ....................................................................... 32 
8.4 
Tjänster för årsredovisningsstatistik .......................................................................... 35 
9 Appendix D. Specifikation mottagningstjänst för händelsemeddelanden ...................... 39 
10 Appendix E. Koder kontrollera tjänst ................................................................................. 40 
 
 
 


---

## Sida 3

 
3 
 
1 
 Ändringshistorik 
Version 
Datum 
Beskrivning 
1.0 
2018-03-01
Lagt till trafikexempel 
Beskrivit ärendestatuskoder 
Kort beskrivning av kodgenerering 
1.0.1 
2018-05-17
Tagit bort felaktig text från trafikexemplen 
1.0.2 
2018-06-07
Bytt URL på exemplen i Appendix C från api-system3 till 
api-accept2 
1.0.3 
2018-08-14
Utökat listan över koder och klartexter för företrädare 
och förtydligat datats giltighet 
1.0.4 
2018-11-08
Uppdaterat tjänster till version 1.1 
1.1 
2019-01-14
Lagt till beskrivning av direktsignering 
1.1 
2019-01-16
Lagt till exempel på anrop vid direktsignering samt 
uppdaterat vissa tjänster till v1.2. 
1.1.1 
2019-03-26
Lagt till text som förklarar att direktsignering är avstängt 
tills vidare 
1.2.0 
2019-05-23
Lagt till ny tjänst för kontrollera samt uppdaterat 
inlämningstjänster till v1.3 
1.3 
2019-06-17
Ny version av informationstjänster, kap 8.1 
2.0 
2019-10-24
Ny version av inlämningstjänster. 
Kontrollera-tjänst returnerar fyra nya utfall, 1039, 1082, 
1171 och 1202. 
2.0.1 
2020-02-10
Rättat kap 5.3.2 Kontrollea. 
Fler ärendestatuskoder i kapitel 5.2.2 
2.1 
2020-11-20
Ny version av api för årsredoviningshändelser, v1.2 
2.1 
2020-12-01
Nya utfall för K3K i kontrollera-tjänst, 1203 - 1210 
2.2 
2021-06-01
Ny version av api för hantera-
arsredovisningsprenumerationer samt automatisk borttag 
av prenumerationer äldre än 6 månader 
2.3 
2022-01-20
Nya utfall för kontrolleratjänst, 1213 och 1214 
2.4 
2022-03-22
Exempel under 8.2.3 på notifiering vid uppladdning till 
eget utrymme. 
2.5 
2022-06-10
Beskrivning av kontrollsumma på handling. 
2.6 
2022-08-22
Beskrivning av skapa token till kontrollsumma 
2.7 
2023-05-09
Nya felkoder 4002 och 4008 i avsnitt 6.2. Koder 1018, 
1034, 1045 och 1209 i avsnitt 10 borttagna. Dessa ersätts 
av felkod 4008. 
Uppdatering av kontrollsumma avseende 
underskriftsdatum samt underskrifter i 
revisionsberättelse. 
2.8 
2023-12-01
Ny kod 4009 i avsnitt 10 
Ny felkod 4010 i avsnitt 6.2 
3.0 
2024-02-26
Ny tjänst /hamta-arsredovisningsstatistik 
3.1 
2024-11-18
Ny kod 1232 för kontrollera tjänst (Appendix E) 


---

## Sida 4

 
4 
 
Version 
Datum 
Beskrivning 
3.2 
2024-08-16
Nya versioner för att stödja ESEF-rapportering med 
hållbarhetsrapport: prenumeration 2.0, händelser 2.0 
samt informationstjänster 1.3. 
Nya felkoder i avsnitt 6.2: 4011 och 7007 
3.3 
2025-02-03
Ny version 1.4 informationstjänster, tillägg revisorsplikt i 
grunduppgifter 
3.4 
2025-05-26
Nya koder i avsnitt 10, 3001-3007. Justerade koder 
avsnitt 6.2, 7001, 7002 borttagen, 7004 tillagd. Länk till 
statuskoder avsnitt 5.2.1.1 
 
 


---

## Sida 5

 
5 
 
2 Inledning  
 
Målgruppen för det här dokumentet är främst teknisk personal som ska arbeta med 
realisering av anslutningar till tjänsterna. Detaljerad information om anslutning till 
tjänsterna – krav på certifikat, brandväggsöppningar mm – hittas i dokumentet 
Anslutningsanvisning Digital inlämning av årsredovisning. 
 
 


---

## Sida 6

 
6 
 
3 Infrastruktur och aktörer 
Tre typer av aktörer samverkar i systemet för inlämning av elektroniska årsredovisningar: 
 Slutanvändare: företagare och andra företagsrepresentanter som har rätt att 
skriva på fastställelseintyget, t.ex. styrelseledamöter; upprättare av 
årsredovisningar, t.ex. redovisningskonsulter mm 
 Programvaruleverantörer: tillverkare av programvara som används för att skapa 
och lämna in elektroniska årsredovisningar 
 Bolagsverket: myndighet med uppgift att ta emot och tillhandahålla 
årsredovisningar  
 
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
 
Figur 1. Översikt över de olika aktörerna i systemlösningen för inlämning av elektroniska årsredovisningar. 
 
Slutanvändarna använder programvara för att skapa en elektronisk avskrift1 av 
årsredovisningen. Programvaran kan vara realiserad på flera olika sätt: som desktop-
applikation, app i en mobil enhet, som webbapplikation etc. Under arbetet med att skapa 
den elektroniska avskriften kan slutanvändarna använda informationstjänster (ljusblåa 
fyrkanter) för att säkerställa att avskriften innehåller rätt uppgifter om räkenskapsperiod, 
styrelsesammansättning osv. 
 
Programvaruleverantörerna fungerar som mellanhänder mellan slutanvändarna och 
Bolagsverket, både tekniskt och juridiskt. Bolagsverket accepterar endast trafik från parter 
som har avtal och som använder godkänt klientcertifikat för att skydda trafiken. Det 
 
1 En avskrift är en handling som innehåller samma information som originalet, men är upprättad i ett annat 
format.   


---

## Sida 7

 
7 
 
innebär att programvaruleverantörerna måste förmedla trafik till och från slutanvändarna 
om slutanvändarna använder klientprogram.  
 
Bolagsverket tillhandahåller tjänster för informationshämtning, inlämning, 
händelsespridning och statistik. Tjänsterna beskrivs detaljerat nedan. 
 
 
 


---

## Sida 8

 
8 
 
4 Beskrivning av tjänsterna 
4.1 Informationstjänster 
API:t för elektroniska årsredovisningar har två tjänster för hämtning av 
årsredovisningsrelaterad information: 
 Information om grunduppgifter för ett aktiebolag 
 Information om status för pågående årsredovisningsärende 
 
 
Figur 2. Informationstjänster för elektronisk inlämning av årsredovisningar 
Bägge tjänsterna är enkla informationstjänster som lämnar uppgifter om ett aktiebolag. 
Syftet med tjänsterna är att säkerställa att rätt information hamnar i årsredovisningarna. 
Vid eventuella ytterligare behov av elektronisk företagsinformation hänvisas till 
Bolagsverkets XML-paket2. 
 
Den tredje tjänsten används för att skapa en kontrollsumma på en handling. 
 
 
Kontrollsumma 
I en digital årsredovisningshandling som lämnas in till Bolagsverket finns det möjlighet 
för programvaran som laddar upp filen till eget utrymme att tagga en kontrollsumma. 
Kontrollsumman kan ses som en stämpel eller version på handlingen och kan användas 
av exempelvis en revisor för att veta vilken version av årsredovisningen som hen granskat. 
Kontrollsumman skapas i anslutning till att handlingen laddas upp till eget utrymme, 
lämpligtvis i samma sekvens så att den handling som laddas upp alltid har en 
kontrollsumma som motsvarar filen. 
 
4.1.1.1 
Skapa token för kontrollsumma 
För att kunna skapa en kontrollsumma till en handling krävs det att användaren skickar 
med en token till kontrollsummans API-anrop. 
 
 
 
2 http://bolagsverket.se/be/sok/xml 
 


---

## Sida 9

 
9 
 
För att tjänsten ska generera ett token krävs följande: 
 Bolaget som pekas ut av organisationsnumret är ett aktiebolag 
 Bolaget inte är avfört från Bolagsverkets register 
 
Tillsammans med token levererar Bolagsverket en text som ska visas för användaren, 
samt datum som beskriver när texten senast ändrades. Syftet med datumet är att 
programvaran ska kunna hålla reda på om texten behöver visas igen eller om den redan 
visats för användaren och det enskilda företaget. 
 
Kontrollsumman skapas via ett API-anrop och läggs in som en metatagg. När 
Bolagsverket tar emot handlingen kommer kontrollsumman för filen att skickas med i 
kvittensmailet tillsammans med kvittensnumret. Om en separat revisionsberättelse har 
lämnats in innehåller kvittensmailet även kontrollsumman för den. Bolagsverket 
kontrollerar även ifall kontrollsumman som står i filen är korrekt. Om den inte är det 
kommer inte kvittensmailet att innehålla kontrollsumman. Programleverantörer som 
prenumererar på händelser (push) kommer däremot att få både kontrollsumman från filen 
samt den som Bolagsverket räknat ut, plus en upplysning om att de är olika. 
 
Kontrollsumman räknas ut utifrån filens innehåll men exkluderar fastställelseintyg, 
revisorspåteckning, underskriftsdatum i årsredovisningen som använder begreppet 
UndertecknandeDatum, underskrifter i revisionsberättelsen samt metataggar. 
 
Exempel: 
En revisor upprättar en revisionsberättelse och sparar kontrollsumman för den och för 
årsredovisningen. Kontrollsumman skapas via Bolagsverkets API. Senare laddas 
handlingarna upp till eget utrymme. När de inkommer till Bolagsverket kommer 
kontrollsumman för respektive handling att skickas med i kvittensmailet. På så sätt får 
revisorn en bekräftelse på att det är samma dokument som inkommit till Bolagsverket. 
 
Kvittensmail: 
Hej! 
Bolagsverket har tagit emot årsredovisningen för Aktiebolaget AB, 556000-1111. 
 
Datum: 2022-05-02, kl. 14.40 
Kvittensnummer: 6000141362 
 
Kontrollsumma årsredovisning: XfhYv3fiJTv/BuSoSVE4hBsIDJvcN50XbVYGK+zH+iY= 
Kontrollsumma revisionsberättelse: lILe5VLrgvSXXvyw+Q/DglbCtFka98SY+Dpw8EnKg6s= 
 
Med vänlig hälsning 
Bolagsverket 
 
4.2 Tjänster för inlämning 
När den elektroniska avskriften är skapad lämnas den in i Bolagsverkets mottagningstjänst 
för elektroniska årsredovisningar. Inlämnandet sker i tre steg mha tre tjänster: 
 Skapa eget utrymme för inlämning 
 Kontrollera 
 Lämna in till det egna utrymmet 


---

## Sida 10

 
10 
 
 
 
Figur 3. Inlämning av elektronisk årsredovisning görs i tre steg: skapa token, kontrollera och lämna in. 
Inlämning görs till eget utrymme3, dvs. den elektroniska avskriften lagras hos 
Bolagsverket för annans räkning, det vill säga för den person som ska skriva på 
fastställelseintyget.  
 
 
Inlämning i tre steg 
Inlämningen görs i tre steg:  
 
1. Anrop till tjänsten Skapa token för inlämning.  
Bolagsverket svarar med ett token (guid) som ska användas vid kontrollera och inlämning, 
samt en förklarande text som ska visas för slutanvändaren. Texten beskriver att inlämning 
görs till eget utrymme och att årsredovisningens fastställelseintyg måste skrivas på och 
skickas in innan årsredovisningen är inkommen till myndigheten. 
 
2. Optionella anrop till tjänsten kontrollera 
Token från steg 1 skickas med i anropet 
 
3. Anrop till tjänsten Lämna in. 
Token från steg 1 skickas med i anropet. 
 
 
Sekvensdiagram inlämning 
Sekvensdiagrammet innehåller ett alternativflöde: 
 Om användaren inte tidigare har sett avtalstexten för det företag inlämningen 
avser så måste den visas för användaren och godkännas innan det andra och 
 
3 Se 
http://www.esamverka.se/download/18.7e784787153f0f33aa51c864/1464274239787/Eget+utrymme+ho
s+myndighet+-+en+v%C3%A4gledning.pdf 
 


---

## Sida 11

 
11 
 
tredje steget i inlämningen. Detta är viktigt för att användaren ska förstå att 
inlämningen MÅSTE åtföljas av en elektronisk underskrift av fastställelseintyget.  
 Om användaren redan sett och godkänt denna version av avtalstexten för aktuellt 
företag så behöver inte användaren involveras för att gå vidare till det andra och 
tredje steget i inlämningen; det kan göras automatiskt av klienten. 
 
Sekvensdiagrammet innehåller också en valfri kontroll av årsredovisningen. 
 
 
 
 
Revisionsberättelse 
Ifall en revisionsberättelse ska lämnas in tillsammans med årsredovisningen så är ett 
alternativ att revisionsberättelsen ingår i samma handling/fil som årsredovisningen. I 
dessa fall sker inlämning på det sätt som beskrivits ovan.  
 
Om årsredovisning och revisionsberättelse inte ingår i samma handling så kan 
revisionsberättelsen lämnas in separat. 
4.2.3.1 
Separat revisionsberättelse 
Vid inlämning av separat revisionsberättelse så lämnar respektive programvara in sin 
handling enligt ovan (4.2.1 Inlämning i tre steg): 
Programvara som skapat årsredovisning  
1. Anrop till tjänsten Skapa token för inlämning.  
2. Valfritt anrop till tjänsten kontrollera ifall årsredovisningen ska kontrolleras innan 
inlämning. 


---

## Sida 12

 
12 
 
3. Anrop till tjänsten Lämna in, skickar årsredovisningen. 
 
Programvara som skapat revisionsberättelsen gör på precis samma sätt: 
1. Anrop till tjänsten Skapa token för inlämning.  
2. Valfritt anrop till tjänsten kontrollera ifall revisionsberättelsen ska kontrolleras innan 
inlämning.  
3. Anrop till tjänsten Lämna in, skickar revisionsberättelsen. 
 
Det har ingen betydelse vilken av handlingarna som lämnas in först, revisonsberättelse 
eller årsredovisning. Det som knyter ihop de olika handlingarna är parametrarna 
organisationsnummer och inbjuden person som skickas med i anropen. Naturligtvis 
måste dock båda handlingarna vara inlämnade innan företagaren ska skriva under och 
skicka in dessa. 
 
4.3 Tjänster för årsredovisningshändelser 
När årsredovisningen lämnas in och registreras kan Bolagsverket återkoppla dessa 
händelser till programvaruleverantörerna. API:et för årsredovisningshändelser består av 
fyra tjänster: 
 Skapa prenumerationer på händelser för företag 
 Ta bort prenumerationer på händelser för företag 
 Hämta befintliga prenumerationer 
 Hämta alla årsredovisningshändelser för ett företag 
 
För att kunna ta emot årsredovisningshändelser måste mottagaren bygga en http/https-
server som svarar på ett REST-anrop (JSON) från Bolagsverket.  
 
Leverantör
Skapa 
prenumerationer
Ta bort 
prenumerationer
Slutanvändare
Programvaru-
leverantör
Mail, SMS 
el dyl
Händelse-
förmedlare
Bolagsverkets
interna system
Styrs av
 
Figur 4. Översikt över lösningen för spridning av årsredovisningshändelser. 
För att ta del av årsredovisningshändelser måste programvaruleverantören skapa en 
prenumeration. Varje prenumeration pekar ut ett företag (det företag vars 
årsredovisningshändelser prenumerationen gäller) och en URL (den URL som 


---

## Sida 13

 
13 
 
Bolagsverkets händelseförmedlare ska anropa när de interna systemen genererar en 
händelse). En leverantör kan ha hur många prenumerationer som helst, och det går att 
skapa flera prenumerationer på ett företag. 
 
När hanteringen av årsredovisningen ger upphov till en händelse så läggs den på en kö 
som läses av händelseförmedlaren. Händelseförmedlaren kontrollerar 
organisationsnummer mot listan över prenumerationer. Om det finns prenumerationer 
för organisationsnumret så anropar händelseförmedlaren de URL:er som lagrats för 
prenumerationerna.  
 
Prenumerationer tas automatiskt bort 6 månader efter att dom har skapats. I de fall 
ingivningen medför behov av t.ex. kompletteringar så kan det vara svårt att exakt säga hur 
lång tid den processen kan ta. Rekommendationen är därför att efter varje anrop till 
tjänsten ”lämna in” också anropa tjänsten ”skapa prenumeration”. Ifall det sedan tidigare 
finns en prenumeration med samma URL och organisationsnummer så skapas ingen 
ytterligare prenumeration, istället utökas giltighetstiden för befintliga prenumerationen 
med nya 6 månader. 
 
 
Sekvensdiagram för årsredovisningshändelser 
Klientprogram
Programvaru-
leverantör
Bolagsverket
Slutanvändare
Markerar intresse för 
årsredovisningshändelser
Anropar 
skapaPrenumeration()
Anropar skapaPrenumeration()
prenumeration skapad
prenumeration skapad
prenumeration skapad
Hantering ger upphov
till händelse
Anropar URL i prenumeration
Anropar klient
Visar händelse
för användare
Kontaktar klient med 
email, SMS el dyl
alt
[Användaren nås 
på annat sätt]
[användaren nås mha
klientsystem]
 


---

## Sida 14

 
14 
 
4.4 Tjänster för årsredovisningsstatistik 
För att underlätta för programvaruleverantörer att tillhandahålla statistik över hur 
tjänsterna används och för felsökning finns 5 statistiktjänster för årsredovisningar och en 
tillhörande stödtjänst: 
 Allmän statistik för inlämnade årsredovisningar, två varianter 
o Per ankomstdatum 
o Per registreringsdatum 
 Statistik för digitalt inlämnade årsredovisningar 
 Statistik för digitalt inlämnade årsredovisningar med separat revisionsberättelse 
 Statistik för digitalt inlämnade årsredovisningar som gett utfall i de automatiska 
kontrollerna 
 Stödtjänst för att lista programvaror för ansluten leverantör 
För alla tjänsterna förutom den allmänna statistiken är det endast programvaru-
leverantörens egen data som redovisas. För att en ny programvara ska synas i statistiken 
måste det ha inkommit minst en årsredovisning eller revisionsberättelse som är registrerad 
från den programvaran, samt att namnet kan kopplas till leverantören på ett entydigt sätt. 
Detta kommer att hanteras i ansökningsprocessen. Rutinen för att införa en ny 
programvara på befintlig leverantör är än så länge manuell. Kontakta Bolagsverket via 
emsdiar@bolagsverket.se. 
 
Allmän statistik för inlämnade/registrerade årsredovisningar 
Två endpoints som levererar allmän statistik för årsredovisningsärenden gällande 
aktiebolag, som har lämnats in på papper och digitalt till Bolagsverket. Sökning på 
ankomstdatum ger alla inkomna ärenden oavsett status, söker man på registreringsdatum 
får man endast registrerade (avslutade) ärenden. 
 
Statistik för digitalt inlämnade årsredovisningar 
Uppgifter om ärenden för digitalt inlämnade årsredovisningar per programvara. Om 
ingen specifik programvara anges, returneras uppgifter för alla programvaror som hör till 
det företag som ställt frågan. 
 
Statistik för digitalt inlämnade årsredovisningar med separat 
revisionsberättelse 
Uppgifter om digitalt inlämnade årsredovisningar med separat revisionsberättelse per 
programvara. Programvaran som avses i den här tjänsten är alltså den programvara som 
skapat revisionsberättelsen. Om ingen specifik programvara anges, returneras uppgifter 
för alla programvaror som hör till det företag som ställt frågan. 
 
Statistik för digitalt inlämnade årsredovisningar som gett utfall i de 
automatiska kontrollerna 
Utfall för automatiska kontroller som gjorts vid inlämning av digitala årsredovisningar per 
programvara. Om ingen specifik programvara anges, returneras utfallet för alla 
programvaror som hör till det företag som ställt frågan. 
 
Stödtjänst hämta programvaror 
Vilka programvaror som finns för leverantören i Bolagsverkets statistik. Varje leverantör 
kan ha flera programvaror, den här listan kan användas för att begränsa sökningen i 
statistiktjänsterna eller för att verifiera att alla programvaror dyker upp i statistiken. 


---

## Sida 15

 
15 
 
5 Teknisk beskrivning av tjänsterna 
5.1 Servicespecifikationer 
Tjänsterna för elektronisk inlämning av årsredovisningar är specificerade som REST-
tjänster (json-meddelanden över http). Specifikationerna är skrivna i OpenAPI/Swagger 
2.0-format4, de hittas här: Servicespecifikationer digital inlämning av årsredovisningar. 
Sidan har länkar till specifikation och dokumentation på svenska och engelska. Det är de 
svenska specifikationerna/dokumentationen som är originalen, de engelska 
specifikationerna/dokumentationen genereras från sin svenska motsvarighet. 
 
Servicespecifikationerna bestämmer gränssnittets utseende. Resten av avsnitten i detta 
kapitel ska ses som förtydliganden till servicespecifikationerna. Om det finns skillnader 
mellan specifikationerna och det här dokumentet så är det specifikationerna som gäller.  
 
5.2 Informationstjänster 
 
Hämta grunduppgifter 
Tjänsten hämtar följande grunduppgifter om organisationer från Bolagsverkets register:  
 Exakt namn enligt registret 
 organisationens bolagsstatus(ar) 
 Lista över redovisningsperioder (from-tom) 
o De fyra senaste räkenskapsperioderna. För nystartade aktiebolag visas 
innevarande räkenskapsperiod. 
o Krav på revisionsberättelse 
o Revisorsplikt 
 Lista över företrädare enligt registret 
 
För att tjänsten ska leverera ett resultat krävs att: 
 organisationen som pekas ut av organisationsnumret får lämna in årsredovisning 
digitalt 
 organisationen inte är avförd från Bolagsverkets register 
5.2.1.1 
Bolagsstatus 
Listan över bolagsstatusar kan innehålla statuskoder som beskriver att organisationen 
bedöms befinna sig i en viss juridisk situation, t.ex. i konkurs eller likvidation. I 
normalfallet finns det ingen bolagsstatus för ett bolag – statuskoder markerar situationer 
som avviker från det normala. Varje status har en beskrivande text som lämpar sig för 
visning för användare, t.ex. ”Ackordsförhandling inledd” eller ”Konkurs avslutad”. 
Förteckning för samtliga statuskoder finns här Statuskoder. 
5.2.1.2 
Krav på revisionsberättelse 
Uppgiften grundas på registrerade uppgifter i Bolagsverkets register. 
 
 
4 Se https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md 
 


---

## Sida 16

 
16 
 
Om företaget uppfyllde minst ett av följande kriterier när räkenskapsåret tog slut ska det 
finnas en revisionsberättelse med årsredovisningen: 
 Företaget hade en bestämmelse i bolagsordningen som säger att det ska ha 
revisor. 
 Företaget hade en revisor registrerad. 
 Företaget uppfyllde inte kraven för att kunna välja bort att ha en revisor. 
 
Läs mer på bolagsverket.se: Revisionsberättelse – Bolagsverket 
5.2.1.3 
Revisorsplikt 
Aktiebolag kan välja att ta bort kravet på revisor om de inte når upp till två av dessa 
värden för båda de två senaste räkenskapsåren: 
 fler än 3 anställda (i medeltal) 
 mer än 1,5 miljoner kronor i balansomslutning 
 mer än 3 miljoner kronor i nettoomsättning. 
Kravet på revisor på grund av aktiebolagets storlek gäller alltså först från och med det 
tredje räkenskapsåret. 
 
Det ska vara samma värde som uppfylls två år i rad.  
 
Läs mer på bolagsverket.se: Revisor i aktiebolag – Bolagsverket 
5.2.1.4 
Lista över företrädare 
Listan visar de företrädare för organisationen som är registrerade i Bolagsverkets register 
och som är relevanta i årsredovisningssammanhang. Dit hör t.ex. VD, styrelseledamöter, 
revisorer mm, men t.ex. inte särskild delgivningsman.  
 
Varje företrädare har ett namn, en identitet – personnummer eller annan identitet, och en 
lista av funktioner. Den andra identiteten används för personer (fysiska eller juridiska) 
som saknar personnummer. Annan identitet kan vara en av: samordningsnummer, GD-
nummer, födelsedatum eller organisationsnummer. Annan identitet är tänkt att visas för 
en användare, inte användas som nyckel eller data i anrop till Bolagsverket. 
 
Listan av funktioner beskriver de funktioner som företrädaren har i organisationen. Varje 
funktion har ett namn och en klartext. Klartexten är tänkt att visas för en användare. Här 
är en tabell över de vanligaste koderna och klartexterna: 
Kod 
Klartext 
DELG
Särskild delgivningsmottagare 
EFT 
Extern firmatecknare (dvs. 
firmatecknare som inte sitter i 
styrelsen) 
EVD 
Extern VD (dvs. VD som inte sitter 
i styrelsen) 
EVVD
Extern Vice VD 
LE 
Styrelseledamot 
OF 
Ordförande 
REV 
Revisor 
REVH 
Huvudansvarig revisor 


---

## Sida 17

 
17 
 
REVL 
Lekmannarevisor 
REVS 
Revisorssuppleant 
REVSL
Suppleant för lekmannarevisor 
SU 
Suppleant 
VD 
Verkställande direktör 
VVD 
Vice verkställande direktör 
 
Observera att listan omfattar ett antal funktioner som inte ger rätt att skriva under 
fastställelseintyget. 
 
Om registerinformation för företrädare 
Bolagsverkets information om ett bolags företrädare är normalt det som gäller för bolaget 
vid varje tidpunkt. Det finns dock lägen när informationen inte är rätt, t.ex. när det 
kommit in ett ärende om styrelseändring som inte hunnit registreras. Att en person finns 
med i listan över företrädare som Bolagsverket har i sitt register är alltså inte en garanti 
för att personen har rätt att skriva under fastställelseintyget. På samma sätt kan det finnas 
lägen då en person har rätt att skriva under fastställelseintyget trots att personen inte finns 
med i listan. 
 
 
Hämta ärendestatus 
Tjänsten hämtar status och ärendenummer för årsredovisningsärende för en organisation 
i Bolagsverkets register.  
 
För att tjänsten ska leverera ett resultat krävs att: 
 organisationen som pekas ut av organisationsnumret får lämna in årsredovisning 
digitalt 
 organisationen inte är avförd från Bolagsverkets register 
För nyregistrerade bolag kan det hända att tjänsten levererar ett tomt resultat eftersom det 
inte finns några årsredovisningsärenden för bolaget. 
 
Ärendenumret kan visas för användaren. Det kan användas i kommunikationen med 
Bolagsverket, handläggarna kan använda ärendenumret för att hämta mer information om 
ärendets handläggning osv. 
 
Följande ärendestatuskoder kan förekomma: 
Kod 
Klartext 
arsred_inkommen 
Årsredovisningen har kommit in till 
Bolagsverket men har inte 
registrerats än – den är under 
handläggning. 
arsred_forelaggande_skickat 
Bolagsverket har skickat 
föreläggande med information om 
vad företaget behöver göra för att 
årsredovisningen ska registreras. 
arsred_komplettering_inkommen
Företaget har lämnat in ny 
årsredovisning – den är under 
handläggning. 


---

## Sida 18

 
18 
 
arsred_registrerad 
Årsredovisningen har registrerats av 
Bolagsverket. 
arsred_avslutad_ej_registrerad 
Bolagsverket har avslutat ärendet 
utan vidare åtgärd. 
Årsredovisningen har inte 
registrerats. 
arsred_saknas 
Företaget har inte lämnat in någon 
årsredovisning än för det aktuella 
räkenskapsåret. 
 
 
Ickefunktionella egenskaper hos tjänsterna 
Båda tjänsterna gör online-anrop till Bolagsverkets interna registersystem. Det data som 
hämtas är alltså en ögonblicksbild av situationen i Bolagsverkets register – det kan t.ex. 
finnas ärenden (styrelseändring mm) som kommit till Bolagsverket men som ännu inte 
registrerats. Svarstiden ligger normalt under 1000 ms, genomsnittstiden är betydligt lägre. 
 
 
Skapa kontrollsumma 
APIet räknar ut en SHA-256 kontrollsumma på filen. Kontrollsumman och dess algoritm 
ska läggas i meta-taggar med följande namn: 
 
ixbrl.innehall.kontrollsumman – för årsredovisning eller filer med både 
årsredovisning och revisionsberättelse. 
 ixbrl.innehall.kontrollsumman.algoritm – algoritm som använts för ovanstående 
kontrollsumma. 
 ixbrl.innehall.kontrollsumman.revision – för separat revisionsberättelse. 
 ixbrl.innehall.kontrollsumman.revision.algoritm – algoritm som använts för 
ovanstående kontrollsumma. 
Då uppgifter i fastställelseintyget, revisorspåteckningen och underskrifter i 
revisionsberättelsen behöver kunna ändras efter att kontrollsumman är skapad behöver 
dessa exkluderas genom att taggas med följande id på den omslutande taggen: 
 id-innehall-faststallelseintyg 
 id-innehall-revisorspateckning 
 id-innehall-underskrifter-revisionsberattelse 
 
Underskriftsdatum i årsredovisningen som använder begreppet UndertecknandeDatum 
kommer automatiskt att exkluderas från kontrollsumman. 
 
Om man vill presentera kontrollsumman visuellt kan följande id användas för att 
exkludera taggen vid uträkning av kontrollsumman: 
 id-innehall-kontrollsumma – för årsredovisning eller filer med både 
årsredovisning och revisionsberättelse. 
 id-innehall-kontrollsumma-revision – för separat revisionsberättelse. 
 


---

## Sida 19

 
19 
 
5.3 Tjänster för inlämning 
Som beskrivits ovan så görs inlämning i tre steg, upprepning av samma steg för separat 
revisionsberättelse om det är aktuellt. 
1. Skapa ett token för inlämning 
2. Använd detta token för att kontrollera årsredovisningen (eller revisionsberättelsen) 
3. Använd detta token för att lämna in årsredovisningen (eller revisionsberättelsen) 
 
Skapa token för inlämning 
Tjänsten genererar ett token för inlämning. 
För att tjänsten ska skapa ett token krävs att: 
 bolaget som pekas ut av organisationsnumret är ett aktiebolag 
 bolaget inte är avfört från Bolagsverkets register 
 
Tillsammans med token levererar Bolagsverket en text som ska visas för användaren, 
samt datum som beskriver när texten ändrades senast. Syftet med datumet är att 
programvaran ska kunna hålla reda på om texten behöver visas igen eller om den redan 
visats för användaren och det enskilda företaget.  
 
 
Kontrollera 
Tjänsten tar emot en digital årsredovisning med eller utan revisionsberättelse alternativt 
separat revisionsberättelse och kontrollerar om handlingen innehåller uppgifter som kan 
hindra ett godkännande. Handlingen måste följa tillämpningsanvisningarna för att ett 
resultat ska kunna returneras. 
Resultatet bör förmedlas till användaren och att användaren ges möjlighet att korrigera 
eventuella problem innan handlingen skickas till eget utrymme.  
 
Som svar lämnar API:t: 
 Kod för typ av hinder. Kan exempelvis användas till att markera aktuell uppgift i 
programvaran. Se appendix E för möjliga koder. 
 Hindertext, ej av teknisk karaktär utan anpassat för användare 
 Typ, utfallets karaktär. 
 Tekniskt information, endast för loggning och felsökning  
 
Ifall tjänsten returnerar upplysningar så är det ändå möjligt att gå vidare och anropa 
tjänsten ”Lämna in”. I de fall tjänsten inte returnerar några upplysningar så är det ingen 
garanti för att årsredovisningen godkänns av Bolagsverket, däremot ökar sannolikheten 
markant. 
 
Eftersom det är möjligt att lämna in även om en kontroll returnerar upplysningar så är det 
inte obligatoriskt för programvaror att anropa denna tjänst. Bolagsverket rekommenderar 
dock att tjänsten används för att minimera förelägganden till slutkund. 
 
Lämna in 
Tjänsten tar emot en digital årsredovisning med eller utan revisionsberättelse alternativt 
separat revisionsberättelse, kontrollerar att den är inlämnad för ett giltigt aktiebolag, 
kontrollerar att filen följer tillämpningsanvisningarna och lagrar handlingen i eget 
utrymme.  


---

## Sida 20

 
20 
 
 
Som svar lämnar API:t: 
 personnummer på avsändaren/inlämnaren av dokumentet 
 personnummer på den person som ska skriva under dokumentet 
 dokumentets längd i bytes 
 Idnummer för dokumentet i avsändarens eget utrymme (unikt nummer i eget 
utrymme) 
 Url 
 SHA-256-checksumma5 på dokumentets innehåll 
 
Det idnummer som lämnas av tjänsten är INTE tänkt att visas för användaren. 
Bolagsverkets handläggare kan inte – får inte – ta del av handlingar i eget utrymme, så de 
kan inte svara på frågor med det idnumret som referens. 
 
Idnumret returneras även i händelsemeddelandet vid prenumeration så att man kan 
koppla ett specifikt dokument till meddelandet. Se exempel i kapitel 5.4.4. 
 
Bolagsverket rekommenderar att idnumret sparas i en loggfil hos klienten eller 
leverantören så att det kan användas vid felsökning. 
 
 
Ickefunktionella egenskaper hos tjänsterna 
5.3.4.1 
Skapa token för inlämning 
Tjänsten gör enkla bearbetningar i Bolagsverkets system och ska normalt svara inom 500 
ms. 
5.3.4.2 
Kontrollera  
Vid kontroll av dokumentet utförs format-, dokument- och register-kontroller. 
Kontrollen kan ta någon-några sekunder beroende på antalet datapunkter, dokumentets 
struktur mm. 
5.3.4.3 
Lämna in 
Vid inlämning, innan dokumentet lagras i eget utrymme, görs en formatvalidering av det 
inlämnade dokumentet. Valideringen kontrollerar att dokumentet är ett giltigt iXBRL-
dokument och att det är uppmärkt med en godkänd taxonomi. Hela valideringen kan ta 
någon-några sekunder beroende på antalet datapunkter, dokumentets struktur mm.  
 
5.4 Tjänster för årsredovisningshändelser 
 
Skapa prenumeration 
Tjänsten skapar en koppling mellan en URL (mottagaradressen för 
händelsemeddelanden) och ett organisationsnummer. Prenumerationer som registreras 
med denna tjänst kommer endast att få händelser som rör årsredovisningshändelser. 
 
 
5 SHA-256: https://tools.ietf.org/html/rfc4634 


---

## Sida 21

 
21 
 
Prenumerationen är giltig under 6 månader från att denna tjänst anropats, därefter tas 
prenumerationen automatiskt bort. 
 
Samma URL kan registreras som mottagare för händelser för många olika 
organisationsnummer. Bolagsverket kommer att kontrollera att URL:en är en giltig URL, 
men inte att den går att nå – vare sig via direkt anrop eller via DNS-uppslag. Det är alltså 
tillåtet att registrera prenumerationer mot en URL som inte etablerats ännu. 
 
En annan variant är att varje prenumeration ges en unik URL och att URL:erna skiljs åt 
sinsemellan mha path- eller get-parametrar. Exempel: 
https://events.accountsoftware.org/arsredovisning?orgnr=1234567890&custid=abc123 
5.4.1.1 
Tillåtna protokoll i URL:en 
Bolagsverket stödjer endast protokollet https i tjänsten för årsredovisningshändelser. 
 
Ta bort prenumeration 
Kombinationen av URL och organisationsnummer fungerar som nyckel för 
prenumerationen. För att avaktivera – ta bort – en prenumeration måste bägge delarna av 
nyckeln anges.  
 
Hämta prenumerationer 
Leverantören kan hämta de prenumerationer som är registrerade av leverantören. Det är 
möjligt att fritt kombinera sökbegreppen URL, organisationsnummer, från när och till 
och med prenumerationen är skapad. Minst ett (1) sökbegrepp måste anges. 
 
Hämta alla årsredovisningshändelser 
Tjänsten hämtar alla händelser för en prenumeration. Det tänkta användningsområdet för 
tjänsten är att hämta information om händelser som ägt rum när mottagande URL inte 
varit tillgänglig, t.ex. om en server varit nere för underhåll osv.  
 
5.4.4.1 
Tillgänglighet till händelseinformation 
Bolagsverket subsystem för årsredovisningshändelser sparar normalt händelser i drygt ett 
år, förutsatt att det finns minst en prenumeration som registrerat intresse för händelsen. 
Syftet med subsystemet är alltså inte att vara ett komplett register över händelser, utan 
endast att mellanlagra händelser som det finns ett registrerat intresse för.  
 
 
Bolagsverkets sändning av händelsemeddelanden 
När en årsredovisningshändelse inträffar så kommer subsystemet för händelser att 
kontrollera om det finns prenumerationer registrerade för händelsen. Om det gör det så 
kommer de registrerade URL:erna att anropas av Bolagsverkets servrar.  
 
5.4.5.1 
Format på händelsemeddelandet 
Meddelandena skickas som UTF-8-kodad JSON. Appendix D innehåller ytterligare 
beskrivning (inkl länkar till Swagger/OpenAPI 2.0-definition som kan användas för att 
generera en mottagningstjänst för meddelandena). Kortfattad beskrivning av 
meddelandeformatet: 


---

## Sida 22

 
22 
 
 
  typ: meddelandetyp. Börjar alltid med ’AR’ 
 id: id för händelsekällan. Sätts till orgnr för den organisation som avses 
  nr: löpnummer för händelsen per bolag.  
  tid: tidpunkt för händelsen 
  data: Innehåller JSON-objekt: 
    status: en av  
      - arsred_inkommen 
      - arsred_registrerad 
      - arsred_avslutad_ej_registrerad 
      - arsred_forelaggande_skickat 
      - arsred_komplettering_inkommen 
      - test 
 
Exempelmeddelande: 
POST /arsredhandelser/bla/bla HTTP/1.1 
Content-type: application/json 
Auth: qwerty123 
 
{ 
    "typ": "AR-v2", 
    "id": "5560456724", 
    "nr": 6, 
    "tid": "2024-06-27T09:52:49.028+02:00", 
    "data": { 
      "version": "2.0", 
      "handlingsinfo": [ 
        { 
          "handling": "arsredovisning", 
          "idnummer": "718a6b33-d536-47ab-8f4d-1a98f4fbaba0" 
        }, 
        { 
          "handling": "revisionsberattelse", 
          "idnummer": "8d340d5a-8b8f-41e5-9d3f-b747fffd0502", 
        } 
      ], 
      "status": "arsred_inkommen" 
    } 
  } 
  
Statuskoderna har samma innebörd som de koder som lämnas av Ärendestatus-tjänsten, 
se 5.2.2. Fältet ”idnummer” innehåller det id dokumentet fick då det laddades upp till eget 
utrymme. På så sätt kan man koppla ett specifikt händelsemeddelande till respektive 
dokument (revisionsberättelse kan laddas upp separat). 
5.4.5.2 
Användning av auth-fältet 
Om fältet auth har angetts vid prenumeration så kommer alla meddelanden som skickas 
till den URL som angavs i prenumerationen att skickas med http-headern auth satt till 
samma värde. Fältet är tänkt att användas som en enkel autenticeringsmekanism, men det 
kan givetvis även användas för andra syften. 


---

## Sida 23

 
23 
 
5.4.5.3 
Testmeddelande 
När en prenumeration skapas så kommer Bolagsverket att skicka ett testmeddelande den 
URL som angavs i prenumerationen. Syftet med meddelandet är att testa 
kommunikationsvägarna (brandväggar osv). Testmeddelandet har status ”test” och nr -1.  
5.4.5.4 
Omsändningsförsök 
Bolagsverket kommer att göra ett antal omsändningsförsök om meddelandet inte kunde 
sändas. Bolagsverket garanterar minst ett omsändningsförsök efter ett dygn. Med 
nuvarande konfiguration görs ett antal omförsök i närtid (de närmsta minuterna) efter det 
första sändningsförsöket. Detaljerna kring omsändningar kan komma att ändras utan 
förvarning. 
5.5 Tjänster för årsredovisningsstatistik 
De fem olika tjänsterna för årsredovisningsstatistik är oberoende av varandra. I de fall 
man efterfrågar programvarurelaterad statistik nyttjas programvaruleverantörens 
klientcertifikat för att säkerställa att bara programvaror kopplade till 
programvaruleverantören visas. Alla tjänster (förutom stödtjänsten) använder ett 
tidsintervall som indata, detta får vara maximalt ett (1) år långt. Tjänsterna räknar antal 
ärenden med en viss egenskap, t.ex. revisionsberättelse – att det inkommit flera versioner 
av en årsredovisning i samma ärende innan ärendet till slut registrerats påverkar alltså inte 
antalet. 
 
Allmän statistik för inlämnade/registrerade årsredovisningar 
Tjänsten levererar statistik för inlämnade eller registrerade årsredovisningsärenden för ett 
givet tidsintervall. Statistiken redovisas per dag och man erhåller antal ärenden med 
digitalt inlämnade årsredovisningar, antal med årsredovisningar i pappersformat, totalt 
antal, hur många som har revisionsberättelse, hur många som har separat 
revisionsberättelse samt hur många (av de digitalt inlämnade) som gått till automatavslut. 
 
Statistik för digitalt inlämnade årsredovisningar 
Tjänsten levererar statistik för digitalt inlämnade årsredovisningsärenden per 
programvara, för ett givet tidsintervall. Om man inte anger en specifik programvara i 
anropet till tjänsten, levereras statistik för programvaruleverantörens samtliga 
programvaror. Statistiken grupperas per datum (dag), programvara (namn), version av 
programvara, regelverk och uppställningsform. Varje grupp har uppgift om totalt antal 
ärenden, antal som gått till automatavslut samt antal där det skapats minst ett 
föreläggande. 
 
Statistik för digitalt inlämnade årsredovisningar med separat 
revisionsberättelse 
Tjänsten levererar statistik för antal digitalt inlämnade årsredovisningsärenden där det 
finns separat revisionsberättelse, för ett givet tidsintervall. Om man inte anger en specifik 
programvara i anropet till tjänsten, levereras statistik för programvaruleverantörens 
samtliga programvaror. Siffrorna grupperas per datum (dag), programvara (namn) och 
version av programvara.  


---

## Sida 24

 
24 
 
 
Statistik för digitalt inlämnade årsredovisningar med utfall i de automatiska 
kontrollerna 
Tjänsten levererar statistik för digitalt inlämnade årsredovisningsärenden där de 
automatiska kontrollerna har gett utfall, för ett givet tidsintervall. Om man inte anger en 
specifik programvara i anropet till tjänsten, levereras statistik för 
programvaruleverantörens samtliga programvaror. Statistiken grupperas per programvara 
(namn), version av programvara, regelverk, uppställningsform och felkod/feltext. 
Feltexten är en beskrivande text av felkoden och bör inte användas maskinellt. 
 
Stödtjänst hämta programvaror 
Tjänsten levererar de programvaror som lämnat in årsredovisningar och/eller 
revisionsberättelser för respektive programvaruleverantör. Som beskrivet i avsnitt 4.4 
måste det ha inkommit minst en årsredovisning eller revisionsberättelse som är registrerad 
från den programvaran, samt att namnet kan kopplas till leverantören på ett entydigt sätt. 
 
 


---

## Sida 25

 
25 
 
6 Appendix A. Felkoder för API:t 
6.1 Om statuskoder i REST-tjänsterna 
Bolagsverket försöker genomgående hålla sig till http-specifikationen när det gäller 
användningen av http-statuskoder. I API:t för digital inlämning av årsredovisningar 
använder vi följande koder: 
Kod 
Innebörd 
200 
OK – lyckat anrop, skapande eller 
uppdatering av resurs gick bra.  
202 
Begäran accepterades – Lyckat anrop, 
ingen response body returneras. 
400 
Allmänt klientfel, t.ex. felaktigt formaterad 
inparameter, saknad inparameter etc 
404 
Saknas – t.ex. rätt formaterat 
organisationsnummer, men det finns inget 
bolag med det organisationsnumret 
500 
Ospecificerat serverfel, t.ex. bugg eller 
driftstörning 
503 
Tjänsten temporärt otillgänglig, t.ex. pga. 
driftstörning 
504 
Timeout i underliggande system 
 
6.2 Lista över felkoder 
Samma felkod kan förekomma i flera av tjänsterna.  
 
4001=Dokumentet är inte en giltig IXBRL-fil 
4002=Du använder en version av programvaran som inte längre stöds av 
tjänsten. Kontakta din programvaruleverantör. 
4003=Ogiltigt organisationsnummer. 
4004=Efterfrågat organisationsnummer är inte ett aktiebolag. 
4005=Ingen träff på efterfrågat organisationsnummer. 
4006=Felaktig url, <url>. 
4007=Ogiltigt personnummer. 
4008=Filen du försöker ladda upp innehåller ett eller flera tekniska fel. 
Kontakta din programvaruleverantör. 
4010=Din årsredovisning är upprättad i en äldre version av taxonomin som 
Bolagsverket inte längre stödjer. För att kunna kontrollera och lämna in din 
årsredovisning digitalt till Bolagsverket, måste du uppdatera till en nyare 
version av taxonomin. Kontakta leverantören av den programvara du använder 
om du behöver hjälp eller har frågor. 
4011=Du kan inte använda denna tjänst då den inte stödjer den här 
företagsformen 
 
5001=Dokumentet saknar eller har tom title tagg 
5002=Dokumentet är inte en IXBRL-fil 
5003=Det förekommer referens till extern bild i dokumentet alternativt ej 
tillåtet format/typ 
5004=Det förekommer referens till extern css/stylesheet 
5005=Det förekommer script i dokumentet 
5006=Dokumentet överstiger tillåten max storlek 


---

## Sida 26

 
26 
 
5007=Det förekommer bilder i dokumentet som överstiger tillåten max storlek 
5008=Dokumentet är inte kodat i rätt character set, ska vara UTF-8. 
5009=Dokumentet saknar taggning av programvara och/eller programversion 
5010=Det förekommer länk till extern resurs 
5011=Det förekommer element med cite attribut 
5012=Det förekommer iframe element i dokumentet 
5013=Det förekommer embed element i dokumentet 
5014=Det förekommer form element i dokumentet 
5015=Det förekommer element med formation attribut 
 
7002=Token gick inte att ta bort. 
7003=Felaktig token. 
7004=Dokumentet innehåller skadlig kod. 
7006=Du kan inte skicka in årsredovisningen eftersom företaget är avvecklat. 
7007=Du kan inte skicka in årsredovisningen digitalt då tjänsten inte 
stödjer den här företagsformen 
 
9000=Inget svar på grund av att uppkopplingen misslyckades. 
9001=Inget svar på grund av timeout från datakälla. 
9002=Inget svar på grund av tekniskt fel. 
9003=Icke godkänd användare av tjänsten. 
9004=Tekniskt felaktig request. 
 
 


---

## Sida 27

 
27 
 
7 Appendix B. Kodgenerering av klienter för REST-API:t 
7.1 Kodgenerering mha Swagger Editor 
Servicespecifikationerna använder formatet OpenAPI 2.0. OpenAPI 2.0 är en 
vidareutveckling av Swagger-formatet. Swagger-projektet har tagit fram flera olika 
mekanismer för att generera klient- och serverkod från en servicespecifikation, bl.a. mha 
Maven för Java. 
Den genereringsmekanism som har stöd för flest programmeringsspråk är Swagger 
Editor, en gratis webbaserad programvara som kan användas för att redigera 
servicespecifikationer och generera kod.  
 
Swagger Editor hittas här: https://swagger.io/swagger-editor/ 
 
 
 


---

## Sida 28

 
28 
 
8 Appendix C. Trafikexempel 
Nedan följer några exempel på anrop och svar, se servicespecifikationer för detaljer och 
förklaringar för de olika parametrarna. 
8.1 Informationstjänster 
 
Hämta grunduppgifter 
8.1.1.1 
Exempel på URL 
GET https://api-accept2.bolagsverket.se/hamta-
arsredovisningsinformation/v1.4/grunduppgifter/5591022107 
8.1.1.2 
Fråga 
Anropet kräver ingen body. 
8.1.1.3 
Svar 
{ 
   "orgnr": "5591022107", 
   "lopnummer": null, 
   "namn": "R.B.G. Bilar Aktiebolag", 
   "status": [], 
   "rakenskapsperioder":    [ 
            { 
         "from": "2023-01-01", 
         "tom": "2023-12-31", 
         "kravPaRevisionsberattelse": "ja", 
         "revisorsplikt": "ja" 
      }, 
            { 
         "from": "2022-01-01", 
         "tom": "2022-12-31", 
         "kravPaRevisionsberattelse": "ja", 
         "revisorsplikt": "uppgift_saknas" 
      }, 
            { 
         "from": "2021-01-01", 
         "tom": "2021-12-31", 
         "kravPaRevisionsberattelse": "ja", 
         "revisorsplikt": "uppgift_saknas" 
      } 
   ], 
   "foretradare":    [ 
            { 
         "fornamn": "Kalle", 
         "namn": "Karlsson", 
         "personnummer": "190001010106", 
         "annanIdentitet": null, 
         "funktioner": [         { 
            "kod": "LE", 
            "text": "styrelseledamot" 
         }] 
      }, 
            { 
         "fornamn": "Test", 
         "namn": "Persson", 
         "personnummer": "187001010102", 
         "annanIdentitet": null, 
         "funktioner": [         { 
            "kod": "SU", 
            "text": "styrelsesuppleant" 
         }] 
      } 
   ] 
} 
 


---

## Sida 29

 
29 
 
 
Hämta ärendestatus 
8.1.2.1 
Exempel på URL 
GET https://api-accept2.bolagsverket.se/hamta-
arsredovisningsinformation/v1.4/arendestatus/5565896866 
8.1.2.2 
Fråga 
Anropet kräver ingen body. 
8.1.2.3 
Svar 
{ 
   "orgnr": "5565896866", 
   "namn": "Brainstorm Aktiebolag", 
   "hamtat": "2018-02-27T10:01:39.598+01:00", 
   "tidpunkt": "2016-12-07", 
   "typ": "arsred_registrerad", 
   "arendenummer": "12345/2016", 
   "rakenskapsperiod":    { 
      "from": "2015-07-01", 
      "tom": "2016-06-30" 
   } 
} 
 
Skapa token för kontrollsumma 
8.1.3.1 
Exempel på URL 
POST https://api-accept2.bolagsverket.se/hamta-
arsredovisningsinformation/v1.1/skapa-inlamningtoken 
8.1.3.2 
Fråga 
{ 
   "pnr": "190001010106", 
   "orgnr": "5565896866" 
} 
8.1.3.3 
Svar 
{ 
   "token": "d0c5b06c-9f6f-4e58-adc4-782838b4a638", 
   "avtalstext": "Ett Eget utrymme har nu skapats för det Företag som Du har angett. 
Genom att använda funktionerna på denna sida ingår Företaget genom Användaren avtal om 
begärt Eget utrymme med Bolagsverket. Utrymmet kan därefter användas så att en 
årsredovisningshandling laddas upp. Vid uppladdningen anges en företrädare för 
Företaget som får ett meddelande när årsredovisningen nått Företagets Eget utrymme om 
att det är dags att elektroniskt\r\n  1. logga in med en e-legitimation som 
Bolagsverket godtar i företagets Eget utrymme,\r\n  2. skriva under ett 
fastställelseintyg och en bestyrkandemening, och\r\n  3. skicka den färdiga handlingen 
från utrymmet till Bolagsverkets mottagningsfunktion så att ett registreringsärende 
startar hos Bolagsverket.\r\n\r\nFör Eget utrymme hos Bolagsverket gäller de allmänna 
villkor som visas via denna länk, http://www.bolagsverket.se/digital-
arsredovisning-villkor. Genom att ta del av villkoren och acceptera dem sluter Du 
avtal för Företagets räkning om Eget utrymme. Samtidigt intygar Du att Du har tagit del 
av villkoren och är behörig att företräda Företaget på detta sätt.", 
   "avtalstextAndrad": "2017-12-06" 
} 
 
Skapa kontrollsumma 
8.1.4.1 
Exempel på URL 
POST https://api-accept2.bolagsverket.se/hamta-
arsredovisningsinformation/v1.1/skapa-kontrollsumma/d0c5b06c-9f6f-
4e58-adc4-782838b4a638 


---

## Sida 30

 
30 
 
8.1.4.2 
Fråga 
{ 
   "fil":"PD94bWwgdmVyc2lvbj0iMS4wIi..." 
} 
8.1.4.3 
Svar 
{ 
    "kontrollsumma": "ttyRb9ploHAFgmF9Khdmvyb7JVxURhZ+Rcik0/RrqNs=", 
    "algoritm": "SHA-256" 
} 
8.2 Tjänster för inlämning 
 
Skapa token för inlämning 
8.2.1.1 
Exempel på URL 
POST https://api-accept2.bolagsverket.se/lamna-in-
arsredovisning/v2.1/skapa-inlamningtoken/ 
8.2.1.2 
Fråga 
{ 
   "pnr": "190001010106", 
   "orgnr": "5565896866" 
} 
8.2.1.3 
Svar 
{ 
   "token": "d0c5b06c-9f6f-4e58-adc4-782838b4a638", 
   "avtalstext": "Ett Eget utrymme har nu skapats för det Företag som Du har angett. 
Genom att använda funktionerna på denna sida ingår Företaget genom Användaren avtal om 
begärt Eget utrymme med Bolagsverket. Utrymmet kan därefter användas så att en 
årsredovisningshandling laddas upp. Vid uppladdningen anges en företrädare för 
Företaget som får ett meddelande när årsredovisningen nått Företagets Eget utrymme om 
att det är dags att elektroniskt\r\n  1. logga in med en e-legitimation som 
Bolagsverket godtar i företagets Eget utrymme,\r\n  2. skriva under ett 
fastställelseintyg och en bestyrkandemening, och\r\n  3. skicka den färdiga handlingen 
från utrymmet till Bolagsverkets mottagningsfunktion så att ett registreringsärende 
startar hos Bolagsverket.\r\n\r\nFör Eget utrymme hos Bolagsverket gäller de allmänna 
villkor som visas via denna länk, http://www.bolagsverket.se/digital-
arsredovisning-villkor. Genom att ta del av villkoren och acceptera dem sluter Du 
avtal för Företagets räkning om Eget utrymme. Samtidigt intygar Du att Du har tagit del 
av villkoren och är behörig att företräda Företaget på detta sätt.", 
   "avtalstextAndrad": "2017-12-06" 
} 
 
Kontrollera 
8.2.2.1 
Exempel på URL 
POST https://api-accept2.bolagsverket.se/lamna-in-
arsredovisning/v2.1/kontrollera/d0c5b06c-9f6f-4e58-adc4-782838b4a638 


---

## Sida 31

 
31 
 
8.2.2.2 
Fråga 
{ 
  "handling": { 
    "fil": "MA==", 
    "typ": "arsredovisning_komplett" 
  } 
} 
8.2.2.3 
Svar 
{ 
   "orgnr": "5565896866", 
   "utfall": [   { 
      "kod": "1165", 
      "text": "Datum för underskrift av fastställelseintyget får inte vara tidigare än 
datum för årsstämman.", 
      "typ": "warn", 
      "tekniskinformation":       [ 
                  { 
            "meddelande": null, 
            "element": "UnderskriftFastallelseintygDatum", 
            "varde": "2019-01-09" 
         }, 
                  { 
            "meddelande": null, 
            "element": "Arsstamma", 
            "varde": "2019-01-10" 
         } 
      ] 
   }] 
} 
 
 


---

## Sida 32

 
32 
 
 
Lämna in 
8.2.3.1 
Exempel på URL 
POST https://api-accept2.bolagsverket.se/lamna-in-arsredovisning/ 
v2.1/inlamning/d0c5b06c-9f6f-4e58-adc4-782838b4a638 
 
8.2.3.2 
Fråga 
{ 
   "undertecknare":"198301019876", 
   "epostadresser":["jag@foretag.com"], 
   "kvittensepostadresser":["minrevisor@foretag.com"], 
   "notifieringEpostadresser":["minrevisor@foretag.com"], 
 
   "handling": { 
      "fil": "MA==", 
      "typ": "arsredovisning_komplett" 
   } 
} 
8.2.3.3 
Svar 
{ 
   "orgnr": "5565896866", 
   "avsandare": "190001010106", 
   "undertecknare": "187001010102", 
   "handlingsinfo":    { 
      "typ": "arsredovisning_komplett", 
      "dokumentlangd": 103133, 
      "idnummer": "49679", 
      "sha256checksumma": ”hufik87TYNl+CMrXpzYk3lzutEWv2fJ/5qAMy5rjUj4=" 
   },    
   "url": "https://arsredovisning-accept2.bolagsverket.se/lamna-
in/visa/engagemang/18772", 
} 
 
8.3 Tjänster för årsredovisningshändelser 
 
Skapa prenumeration 
8.3.1.1 
Exempel på URL 
POST https://api-accept2.bolagsverket.se/hantera-
arsredovisningsprenumerationer/v2.0/handelseprenumeration 
8.3.1.2 
Fråga 
{ 
  "prenumerationer": 
  [ 
    { 
      "url":"https://programvaruleverantor.example.com/arsredovisning/handelser/", 
      "orgnr":"5563331494" 
    } 
  ] 
} 
8.3.1.3 
Svar 
Ingen response body returneras när anropet lyckats. 


---

## Sida 33

 
33 
 
 
Ta bort prenumeration 
8.3.2.1 
Exempel på URL 
DELETE https://api-accept2.bolagsverket.se/hantera-
arsredovisningsprenumerationer/v2.0/handelseprenumeration 
8.3.2.2 
Fråga 
{ 
 "url":" https://programvaruleverantor.example.com/arsredovisning/handelser/", 
 "orgnr":"5563331494" 
} 
8.3.2.3 
Svar 
Ingen response body returneras när anropet lyckats. 
 
Hämta prenumerationer 
8.3.3.1 
Exempel på URL 
GET https://api-accept2.bolagsverket.se/hantera-
arsredovisningsprenumerationer/v2.0/handelseprenumeration?url=https://www.
mydomain.se/&from=2021-04-01 
8.3.3.2 
Fråga 
Anropet kräver ingen body. 
8.3.3.3 
Svar 
{ 
  "prenumerationer": [ 
    { 
      "url": "https://www.mydomain.se/", 
      "orgnr": "1234567890", 
      "registrerad": "2021-05-30T16:22:17.511+02:00", 
      "avslutas": "2021-11-30" 
    }, 
    { 
      "url": "https://www.mydomain.se/", 
      "orgnr": "2345678901", 
      "registrerad": "2021-05-17T16:22:17.511+02:00", 
      "avslutas": "2021-11-17" 
    }, 
    { 
      "url": "https://www.mydomain.se/", 
      "orgnr": "3456789012", 
      "registrerad": "2021-05-01T16:22:17.511+02:00", 
      "avslutas": "2021-11-01" 
    } 
  ] 
} 
 
 


---

## Sida 34

 
34 
 
 
Hämta årsredovisningshändelser 
8.3.4.1 
Exempel på URL 
POST https://api-accept2.bolagsverket.se/hamta-
arsredovisningshandelser/v2.0/handelser 
8.3.4.2 
Fråga 
{ 
  "url":"https://programvaruleverantor.example.com/arsredovisning/handelser/", 
  "orgnr":["5564940640", "5564943875"], 
  "fromtidpunkt":"2021-11-01T09:09:12.911+01:00", 
  "tomtidpunkt":"2022-02-20T09:09:51.911+01:00" 
} 
8.3.4.3 
Svar 
{ 
  "meddelanden": [ 
    { 
      "typ": "AR-v2", 
      "id": "5564940640", 
      "nr": 1, 
      "tid": "2022-01-30T13:30:41.741+01:00", 
      "data": { 
        "version": "2.0", 
        "handlingsinfo": [ 
          { 
            "handling": "arsredovisning", 
            "idnummer": "18772", 
            "kontrollsumma": { 
              "digest": "KSco8wmTfA5p4Ij1YmIvRJrVT5DlAB0egUOm8RmSKrM=", 
              "algoritm": "SHA-256", 
              "upplysning": null 
            } 
          }, 
          { 
            "handling": "revisionsberattelse", 
            "idnummer": "18773", 
            "kontrollsumma": { 
              "digest": "Fg7k9nVQJ+btBw+iHpE0MHFm9tSZ2bnwptOekM7xZEw=", 
              "algoritm": "SHA-256", 
              "upplysning": null 
            } 
          } 
        ], 
        "status": "arsred_inkommen" 
      } 
    }, 
    { 
      "typ": "AR-v2", 
      "id": "5564943875", 
      "nr": 2, 
      "tid": "2022-01-30T13:34:59.296+01:00", 
      "data": { 
        "version": "2.0", 
        "handlingsinfo": [ 
          { 
            "handling": "arsredovisning", 
            "idnummer": "18774"           
          } 
        ], 
        "status": "arsred_registrerad" 
      } 
    } 
  ] 
} 
 
 


---

## Sida 35

 
35 
 
8.4 Tjänster för årsredovisningsstatistik 
 
Allmän statistik för inlämnade årsredovisningar 
8.4.1.1 
Exempel på URL 
POST https://api-accept2.bolagsverket.se/hamta-
arsredovisningsstatistik/v1.0/statistik/allman/ankomstdatum 
8.4.1.2 
Fråga 
{ 
  "startDatum" : "2023-03-10", 
  "stoppDatum" : "2023-03-11" 
} 
8.4.1.3 
Svar 
{ 
   "indata":    { 
      "startDatum": "2023-03-10", 
      "stoppDatum": "2023-03-11" 
   }, 
   "statistik":    [ 
            { 
         "datum": "2023-03-10", 
         "digitala": 921, 
         "papper": 1040, 
         "totalt": 1961, 
         "automatavslut": 854, 
         "revisionsberattelse": 89, 
         "separatrevisionsberattelse": 49 
      }, 
            { 
         "datum": "2023-03-11", 
         "digitala": 237, 
         "papper": 10, 
         "totalt": 247, 
         "automatavslut": 223, 
         "revisionsberattelse": 6, 
         "separatrevisionsberattelse": 5 
      } 
   ] 
} 
 
Statistik per programvara 
8.4.2.1 
Exempel på URL 
POST https://api-accept2.bolagsverket.se/hamta-
arsredovisningsstatistik/v1.0/statistik/programvara/ 


---

## Sida 36

 
36 
 
8.4.2.2 
Fråga 
{ 
  "startDatum":"2022-01-01", 
  "stoppDatum":"2022-01-09", 
  "programvara":"Exempel Bokslut Företag" 
} 
8.4.2.3 
Svar 
{ 
  "indata": { 
    "startDatum": "2022-01-01", 
    "stoppDatum": "2022-01-09", 
    "programvara": "Exempel Bokslut Företag" 
  }, 
  "statistik": [ 
    { 
      "programvara": "Exempel Bokslut Företag", 
      "datum": "2022-01-03", 
      "version": "Version 2021.3 Not=13 RRBR=17", 
      "regelverk": "K2", 
      "uppstallningsform": "risbs", 
      "totalt": 1, 
      "automatavslut": 1, 
      "forelaggande": 0 
    }, 
    { 
      "programvara": "Exempel Bokslut Företag", 
      "datum": "2022-01-09", 
      "version": "Version 2021.3 Not=13 RRBR=17", 
      "regelverk": "K2", 
      "uppstallningsform": "risbs", 
      "totalt": 8, 
      "automatavslut": 8, 
      "forelaggande": 0 
    } 
  ] 
} 
 
Statistik för digitalt inlämnade årsredovisningar med separat 
revisionsberättelse 
8.4.3.1 
Exempel på URL 
POST https://api-accept2.bolagsverket.se/hamta-
arsredovisningsstatistik/v1.0/statistik/programvara/separatrevisionsberattelse/ 


---

## Sida 37

 
37 
 
8.4.3.2 
Fråga 
{ 
  "startDatum":"2022-01-01", 
  "stoppDatum":"2022-01-09", 
  "programvara":"Exempel Audit med Bokslut" 
} 
8.4.3.3 
Svar 
{ 
  "indata": { 
    "startDatum": "2023-04-01", 
    "stoppDatum": "2023-04-04", 
    "programvara": "Exempel Audit med Bokslut" 
  }, 
  "statistik": [ 
    { 
      "programvara": "Exempel Audit med Bokslut", 
      "datum": "2023-04-01", 
      "version": "Version 2023.1.1 RBExcel=3", 
      "totalt": 2 
    }, 
    { 
      "programvara": "Exempel Audit med Bokslut", 
      "datum": "2023-04-04", 
      "version": "Version 2023.2 RBExcel=2", 
      "totalt": 1 
    } 
  ] 
} 
 
Statistik för digitalt inlämnade årsredovisningar som gett utfall i de 
automatiska kontrollerna 
8.4.4.1 
Exempel på URL 
POST https://api-accept2.bolagsverket.se/hamta-
arsredovisningsstatistik/v1.0/statistik/programvara/kontroller/ 


---

## Sida 38

 
38 
 
8.4.4.2 
Fråga 
{ 
  "startDatum":"2022-01-01", 
  "stoppDatum":"2022-01-03", 
  "programvara":"Exempel Audit med Bokslut" 
} 
8.4.4.3 
Svar 
{ 
  "indata": { 
    "startDatum": "2022-01-01", 
    "stoppDatum": "2022-01-03", 
    "programvara": "Exempel Audit med Bokslut" 
  }, 
  "statistik": [ 
    { 
      "programvara": "Exempel Audit med Bokslut", 
      "version": "Version 2021.3 Not=13 RRBR=17", 
      "regelverk": "K2", 
      "uppstallningsform": "risbs", 
      "felkod": "1042", 
      "feltext": "Det finns ett pågående ärende för det här räkenskapsåret som inte kom 
in digitalt.", 
      "totalt": 1 
    }, 
    { 
      "programvara": "Exempel Audit med Bokslut", 
      "version": "Version 2021.3 Not=13 RRBR=17", 
      "regelverk": "K2", 
      "uppstallningsform": "risbs", 
      "felkod": "1171", 
      "feltext": "Kontrollera att samtliga och rätt företrädare och revisorer har 
skrivit under. Namnen i handlingarna stämmer inte med det som är registrerat.", 
      "totalt": 1 
    } 
  ] 
} 
 
Stödtjänst hämta programvaror 
8.4.5.1 
Exempel på URL 
GET https://api-accept2.bolagsverket.se/hamta-
arsredovisningsstatistik/v1.0/stodtjanster/programvaror/ 
8.4.5.2 
Fråga 
Anropet kräver ingen body. 
8.4.5.3 
Svar 
{ 
  "programvaror": [ 
    "Exempel Bokslut Företag", 
    "Exempel Bokslut", 
    "Exempel Audit med Bokslut" 
  ] 
} 
 
 


---

## Sida 39

 
39 
 
9 Appendix D. Specifikation mottagningstjänst för 
händelsemeddelanden 
För att implementera en mottagningstjänst som kan ta emot meddelanden om 
årsredovisningshändelser från Bolagsverket finns det Swagger/OpenAPI 2.0-
specifikationer publicerade på Bolagsverkets sida för tjänstespecifikationer (v1.3 och 
v2.0). Dessa är mallar som tjänsteleverantören måste justera med adress och sökväg innan 
de kan användas. Om man använder prenumerationstjänst v2.0 kommer meddelanden att 
ha format enligt mall v2.0 (eller eventuell framtida v2.1 osv), äldre versioner av 
prenumerationstjänsten fortsätter skicka meddelanden enligt v1.3. 
 
 


---

## Sida 40

 
40 
 
10 Appendix E. Koder kontrollera tjänst 
1015
Räkenskapsårets sista dag har inte passerats.
1019
Fastställelseintyget saknas.
1020
Företagsnamnet saknas i årsredovisningen.
1021
Företagsnamnet saknas i revisionsberättelsen.
1022
Organisationsnumret saknas i revisionsberättelsen.
1029
Datum för avslutad revision saknas i revisionsberättelsen.
1030
Underskrift saknas i revisionsberättelsen.
1031
Revisorns namn i revisionsberättelsen och revisorns namn i årsredovisningen 
(revisorspåteckningen) stämmer inte överens. 
1033
Organisationsnumret eller räkenskapsåret i revisionsberättelsen och 
årsredovisningen stämmer inte överens. 
1035
Organisationsnumret i årsredovisningen stämmer inte med det valda företaget. En 
ny årsredovisning bör laddas upp från programvaran. 
1037
Uppgift om valuta saknas i årsredovisningen.
1038
Valutan får endast vara SEK eller EUR.
1039
Valutan i årsredovisningen stämmer inte med registrerad valuta hos Bolagsverket. 
Registrerad valuta är [registeruppgift]. 
1040
Det måste vara samma valuta i hela årsredovisningen.
1046
Räkenskapsåret får inte vara längre än 18 månader.
1050
Landskoden saknas i årsredovisningen.
1051
Förvaltningsberättelsen saknas.
1060
Resultaträkningen saknas.
1064
Balansräkningen saknas.
1072
Datum för årsstämman får inte vara tidigare än datum för revisionsberättelsen.
1082
Räkenskapsåret stämmer inte med det räkenskapsår som är registrerat hos 
Bolagsverket. Registrerat räkenskapsår är [registeruppgift]. 
1101
Datum för årsstämman får inte vara tidigare än eller samma som räkenskapsårets 
sista dag.  
1103
Datum för årsstämman saknas i fastställelseintyget.
1107
Datum för underskrift saknas i årsredovisningen.
1114
Datum för underskrift av årsredovisningen får inte vara tidigare än eller samma 
som räkenskapsårets sista dag. 
1115
Datum för avslutad revision får inte vara tidigare än datum för årsredovisningen.
1116
Årsredovisningen verkar inte vara upprättad på svenska. 
1163
Revisorns namn (revisorspåteckningen) saknas i årsredovisningen.
1164
Datum för underskrift saknas i fastställelseintyget. 
1165
Datum för underskrift av fastställelseintyget får inte vara tidigare än datum för 
årsstämman. 
1169
Namnförtydligandet saknas i fastställelseintyget. 
1170
Företagsnamnet i årsredovisningen och revisionsberättelsen stämmer inte 
överens. 
1171
Namnen på företrädarna i handlingarna stämmer inte med de som var registrerade 
hos Bolagsverket den dag årsredovisningen skrevs under. Registrerade företrädare 
den dagen var: [registeruppgift]. 


---

## Sida 41

 
41 
 
1172
Det saknas uppgift om årsredovisningen avges av styrelsen eller styrelsen och 
verkställande direktören.  
1173
Det saknas uppgift om vilket språk årsredovisningen är upprättad på. 
1174
Det saknas uppgift om vilken mätenhet (t.ex. kr, tkr) beloppen är angivna i. 
1175
Datum för revisorspåteckning saknas
1176
Datum för revisorspåteckningen får inte vara tidigare än styrelsens underskrift av 
årsredovisningen eller senare än årsstämman. 
1177
I revisorspåteckningen saknas uttalande om att revisor avstyrkt att resultat- och 
balansräkning fastställs.   
1178
Datum för årsstämman får inte vara senare än dagens datum. 
1179
Fastställelseintyget innehåller inte alla uppgifter som behövs. 
1183
Datum för årsstämman är tidigare än styrelsens underskrift.
1184
Datum för revisorspåteckningen får inte vara tidigare än eller samma som 
räkenskapsårets sista dag.  
1185
Datum för underskrift av revisionsberättelsen får inte vara tidigare än eller samma 
som räkenskapsårets sista dag. 
1187
Årsredovisningen verkar sakna en resultaträkning.
1188
Årsredovisningen verkar sakna en balansräkning.
1194
Det saknas uppgift om årsredovisningen avges av styrelsen eller styrelsen och 
verkställande direktören.  
1195
Det saknas uppgift om vilken mätenhet (t.ex. kr, tkr) beloppen är angivna i. 
1201
Det saknas för- eller efternamn på den eller de som skrivit under årsredovisningen. 
1202
Aktiekapitalet i årsredovisningen stämmer inte med registrerat aktiekapital hos 
Bolagsverket. Registrerat aktiekapital är [registeruppgift]. 
1203
Årsredovisningen verkar sakna en balansräkning för moderföretaget
1204
Årsredovisningen verkar sakna en balansräkning för koncernen
1205
Årsredovisningen verkar sakna en resultaträkning för moderföretaget
1206
Årsredovisningen verkar sakna en resultaträkning för koncernen
1207
Resultaträkningen för koncernen saknas.
1208
Balansräkningen för koncernen saknas.
1210
Räkenskapsåret för moderbolaget och koncernen måste sluta samma datum.
1213
Det saknas uppgift om vilken mätenhet (t.ex. kr, tkr) beloppen är angivna i.
1214
Datum för underskrifter saknas i årsredovisningen. Alla underskrifter måste ha ett 
datum. 
1232
Datum för årsredovisningen är senare än styrelsens underskrift.
3001
Balansräkningen saknar uppgift om ”Summa tillgångar”
3002
Balansräkningen saknar uppgift om ”Summa eget kapital och skulder”.
3003
”Summa tillgångar” finns på flera ställen men med olika belopp. Det måste vara 
samma belopp på alla ställen för aktuellt räkenskapsår. 
3004
”Summa eget kapital och skulder” finns på flera ställen men med olika belopp. Det 
måste vara samma belopp på alla ställen för aktuellt räkenskapsår. 
3005
”Summa tillgångar” och ”Summa eget kapital och skulder” stämmer inte överens. 
Det måste vara samma belopp på alla ställen. Summa tillgångar: [VÄRDEN] Summa 
eget kapital och skulder: [VÄRDEN] 


---

## Sida 42

 
42 
 
3006
Jämförelsesiffror saknas i balansräkningen. De behövs om det inte är företagets 
första räkenskapsår. 
3007
Jämförelsesiffror saknas i resultaträkningen. De behövs om det inte är företagets 
första räkenskapsår. 
4009
Din årsredovisning är upprättad i en version av taxonomin som Bolagsverket 
kommer att sluta stödja snart. För att du fortsättningsvis ska kunna lämna in din 
årsredovisning digitalt behöver du uppdatera till en nyare version av taxonomin. 
Kontakta leverantören av den programvara du använder om du behöver hjälp eller 
har frågor. 
 
