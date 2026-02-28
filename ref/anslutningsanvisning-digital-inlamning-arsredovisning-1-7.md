# anslutningsanvisning-digital-inlamning-arsredovisning-1-7

## Sida 1

Digital inlämning av 
årsredovisning
Anslutningsanvisning
Version 1.7
 


---

## Sida 2

 
2 
 
Innehållsförteckning 
 Ändringshistorik ........................................................................................................................ 3 
 Inledning ..................................................................................................................................... 4 
 Översikt miljöer ......................................................................................................................... 5 
3.1 
Endpoints i testmiljö ...................................................................................................... 6 
3.2 
Endpoints i acceptansmiljö ........................................................................................... 6 
3.3 
Endpoints i produktionsmiljö ...................................................................................... 7 
 Åtkomst till testmiljö ................................................................................................................ 9 
4.1 
Öppning av brandvägg .................................................................................................. 9 
4.2 
Trust av Bolagsverkets servercertifikat ..................................................................... 10 
4.2.1 
Anslutning med webbläsare ............................................................................. 10 
4.2.2 
Anslutning med klientprogramvara................................................................. 10 
 Åtkomst till tjänster i acceptansmiljö ................................................................................... 11 
5.1 
Öppning av brandvägg ................................................................................................ 11 
5.2 
Trust av Bolagsverkets servercertifikat ..................................................................... 11 
5.3 
Autentisering av klientens organisationscertifikat ................................................... 11 
5.4 
Auktorisering av klienten ............................................................................................ 12 
5.5 
Utgående e-post ............................................................................................................ 13 
 Åtkomst till tjänster i produktionsmiljö ............................................................................... 14 
6.1 
Öppning av brandvägg ................................................................................................ 14 
6.2 
Trust av Bolagsverkets servercertifikat ..................................................................... 14 
6.3 
Autentisering av klientens organisationscertifikat ................................................... 14 
6.4 
Auktorisering av klienten ............................................................................................ 14 
6.5 
Utgående e-post ............................................................................................................ 14 
 
 
 


---

## Sida 3

 
3 
 
 
 Ändringshistorik 
Version 
Datum 
Beskrivning 
1.0 
2018-03-01
Kompletterat med information om 
produktionsmiljön 
1.1 
2018-10-04
Kompletterat med information om hur 
programvaruleverantörer gör för att kunna få 
e-post från acceptansmiljön. Rättat 
endpoints i test- och acceptansmiljö samt 
stegat version på uppdaterade tjänster till 
v1.1. 
1.1 
2019-01-16
Stegat version på uppdaterade tjänster till 
v1.2 samt lagt till ny tjänst för 
inlämningsstatus vid direktsignering. 
1.2 
2019-10-23
Stegat version på uppdaterade tjänster samt 
tagit bort tjänsten inlämningsstatus 
1.3 
2020-06-16
Uppdaterad information om Expisoft/Steria 
CA gällande from 2020-06-01. 
1.3.1 
2020-11-20
Ny version av /hamta-
arsredovisningshandelser  
1.3.2 
2021-06-01
Ny version av /hantera-
arsredovisningsprenumerationer 
1.3.3 
2022-03-20
Ny version av /lamna-in-arsredovisning 
1.3.4 
2022-05-21
Ny version av /hamta-
arsredovisningshandelser samt ny tjänst 
/skapa-kontrollsumma 
1.3.5 
2022-08-24
Ny version av /skapa-kontrollsumma 
1.4.0 
2022-09-06
Tillägg av tjänster för ESEF 
1.4.1 
2022-12-21
Rättelse av några ESEF-länkar. 
1.4.2 
2023-11-14
Justering att brandväggsöppning inte behövs 
för produktionsmiljö + ändirng epostadress i 
4.1 och 5.5 
1.5 
2024-02-26
Nya tjänst /hamta-arsredovisningsstatistik 
1.6 
2024-08-14
Nya versioner av /hamta-
arsredovisningsinformation, /lamna-in-
arsredovisning-asynkront, /hantera-
arsredovisningsprenumerationer, /hamta-
arsredovisningshandelser 
1.7 
2025-02-03
Ny version /hamta-
arsredovisningsinformation 
 
 
 


---

## Sida 4

 
4 
 
 
Inledning  
Anslutning till tjänsterna för digital inlämning av årsredovisning kräver att ansluten part 
ges åtkomst till tjänsterna i Bolagsverkets miljöer samt att ansluten part autentiserar sig 
med sitt organisationscertifikat. 
 
Det här dokumentet syftar till att beskriva vad en ansluten part behöver göra för att få 
åtkomst till testmiljö, acceptansmiljö och produktionsmiljö. En översikt över miljöer och 
hur tjänsterna skyddas finns i kapitel 3. 
 
Vad som sedan krävs för att, i tur och ordning, ansluta till testmiljö, acceptansmiljö och 
produktionsmiljö beskrivs i kapitel 0, 5 och 6. 
 
 
 


---

## Sida 5

 
5 
 
 
Översikt miljöer 
 
Bolagsverket
Leverantör
Testversioner av 
interna system
Brandvägg
Webproxy
Statiskt 
testdata
Tjänster för digital 
inlämning av 
årsredovisning
Klient
Brandvägg
Webproxy
Tjänster för digital 
inlämning av 
årsredovisning
Klient
Testdata
Produktionssystem
Webproxy
Tjänster för digital 
inlämning av 
årsredovisning
Klient
Proddata
  
 
Bolagsverket har tre olika miljöer för digital inlämning av årsredovisning i drift samtidigt: 
 En testmiljö (ljusgult, längst till vänster) som endast levererar statiskt testdata 
 En acceptansmiljö för acceptanstest (gult, i mitten) som hämtar testdata från 
testversioner av Bolagsverkets interna system 
 En produktionsmiljö (grönt, längst till höger). 
 
Gemensamt för alla tre miljöerna är att åtkomst till tjänsterna skyddas webproxies. De två 
testmiljöerna skyddas dessutom av brandväggar. 
 
 


---

## Sida 6

 
6 
 
3.1 Endpoints i testmiljö 
Följande endpoints gäller till tjänsterna i Bolagsverkets testmiljö: 
Testbänk: 
https://arsredovisning-accept2.bolagsverket.se/testbank/ 
 
Tjänster:  
Informationstjänster 
https://api-accept2.bolagsverket.se/testapi/hamta-arsredovisningsinformation/v1.4/grunduppgifter/{orgnr} 
https://api-accept2.bolagsverket.se/testapi/hamta-arsredovisningsinformation/v1.4/arendestatus/{orgnr} 
https://api-accept2.bolagsverket.se/testapi/hamta-arsredovisningsinformation/v1.1/skapa-
inlamningtoken/ (POST) 
https://api-accept2.bolagsverket.se/testapi/hamta-arsredovisningsinformation/v1.1/skapa-
kontrollsumma/{token} (POST) 
 
Tjänster för inlämning K2, K3 och K3K 
https://api-accept2.bolagsverket.se/testapi/lamna-in-arsredovisning/v2.1/skapa-inlamningtoken/ (POST) 
https://api-accept2.bolagsverket.se/testapi/lamna-in-arsredovisning/v2.1/kontrollera/{token} (POST) 
https://api-accept2.bolagsverket.se/testapi/lamna-in-arsredovisning/v2.1/inlamning/{token} (POST) 
 
Tjänster för inlämning ESEF och/eller CSRD 
https://api-accept2.bolagsverket.se/testapi/lamna-in-arsredovisning-asynkront/v2.0/skapa-
inlamningtoken (POST) 
https://api-accept2.bolagsverket.se/testapi/lamna-in-arsredovisning-asynkront/v2.0/inlamning/{token} 
(POST) 
https://api-accept2.bolagsverket.se/testapi/lamna-in-arsredovisning-
asynkront/v2.0/valideringsresultat/{idnummer} (POST) 
 
Tjänster för årsredovisningshändelser 
https://api-accept2.bolagsverket.se/testapi/hantera-
arsredovisningsprenumerationer/v2.0/handelseprenumeration/ (POST) 
https://api-accept2.bolagsverket.se/testapi/hantera-
arsredovisningsprenumerationer/v2.0/handelseprenumeration/ (DELETE) 
https://api-accept2.bolagsverket.se/testapi/hantera-
arsredovisningsprenumerationer/v2.0/handelseprenumeration/ (GET) 
https://api-accept2.bolagsverket.se/testapi/hamta-arsredovisningshandelser/v2.0/handelser/ (POST) 
 
Giltiga organisationsnummer för test är 1234567890 och 1234567891.  
 
Se Servicespecifikationer och Teknisk guide för mer information om tjänsterna. 
 
3.2 Endpoints i acceptansmiljö 
Följande endpoints gäller till tjänsterna i Bolagsverkets acceptansmiljö: 
Inlämning: 
https://arsredovisning-accept2.bolagsverket.se/lamna-in/ 
 
Tjänster: 
Informationstjänster 
https://api-accept2.bolagsverket.se/hamta-arsredovisningsinformation/v1.4/grunduppgifter/{orgnr} 
https://api-accept2.bolagsverket.se/hamta-arsredovisningsinformation/v1.4/arendestatus/{orgnr} 
https://api-accept2.bolagsverket.se/hamta-arsredovisningsinformation/v1.1/skapa-inlamningtoken/ 
(POST) 
https://api-accept2.bolagsverket.se/hamta-arsredovisningsinformation/v1.1/skapa-kontrollsumma/{token} 
(POST) 
 
Tjänster för inlämning K2, K3 och K3K 


---

## Sida 7

 
7 
 
https://api-accept2.bolagsverket.se/lamna-in-arsredovisning/v2.1/skapa-inlamningtoken/ (POST) 
https://api-accept2.bolagsverket.se/lamna-in-arsredovisning/v2.1/kontrollera/{token} (POST) 
https://api-accept2.bolagsverket.se/lamna-in-arsredovisning/v2.1/inlamning/{token} (POST) 
 
Tjänster för inlämning ESEF och/eller CSRD 
https://api-accept2.bolagsverket.se/lamna-in-arsredovisning-asynkront/v2.0/skapa-inlamningtoken 
(POST) 
https://api-accept2.bolagsverket.se/lamna-in-arsredovisning-asynkront/v2.0/inlamning/{token} (POST) 
https://api-accept2.bolagsverket.se/lamna-in-arsredovisning-
asynkront/v2.0/valideringsresultat/{idnummer} (POST) 
 
Tjänster för årsredovisningshändelser 
https://api-accept2.bolagsverket.se/hantera-
arsredovisningsprenumerationer/v2.0/handelseprenumeration/ (POST) 
https://api-accept2.bolagsverket.se/hantera-
arsredovisningsprenumerationer/v2.0/handelseprenumeration/ (DELETE) 
https://api-accept2.bolagsverket.se/hantera-
arsredovisningsprenumerationer/v2.0/handelseprenumeration/ (GET) 
https://api-accept2.bolagsverket.se/hamta-arsredovisningshandelser/v2.0/handelser/ (POST) 
 
Tjänster för statistik för årsredovisningsärenden 
https://api-accept2.bolagsverket.se/hamta-arsredovisningsstatistik/v1.0/statistik/allman/ankomstdatum 
(POST) 
https://api-accept2.bolagsverket.se/hamta-
arsredovisningsstatistik/v1.0/statistik/allman/registreringsdatum (POST) 
https://api-accept2.bolagsverket.se/hamta-arsredovisningsstatistik/v1.0/statistik/programvara/ (POST) 
https://api-accept2.bolagsverket.se/hamta-
arsredovisningsstatistik/v1.0/statistik/programvara/separatrevisionsberattelse/ (POST) 
https://api-accept2.bolagsverket.se/hamta-arsredovisningsstatistik/v1.0/statistik/programvara/kontroller/ 
(POST) 
 
Stödtjänster för statistik för årsredovisningsärenden 
https://api-accept2.bolagsverket.se/hamta-arsredovisningsstatistik/v1.0/stodtjanster/programvaror/ (GET) 
 
 
3.3 Endpoints i produktionsmiljö 
Följande endpoints gäller till tjänsterna i Bolagsverkets produktionsmiljö: 
Inlämning: 
https://arsredovisning.bolagsverket.se/lamna-in/ 
 
Tjänster: 
Informationstjänster 
https://api.bolagsverket.se/hamta-arsredovisningsinformation/v1.4/grunduppgifter/{orgnr} 
https://api.bolagsverket.se/hamta-arsredovisningsinformation/v1.4/arendestatus/{orgnr} 
https://api.bolagsverket.se/hamta-arsredovisningsinformation/v1.1/skapa-inlamningtoken/ (POST) 
https://api.bolagsverket.se/hamta-arsredovisningsinformation/v1.1/skapa-kontrollsumma/{token} (POST) 
 
Tjänster för inlämning K2, K3 och K3K 
https://api.bolagsverket.se/lamna-in-arsredovisning/v2.1/skapa-inlamningtoken/ (POST) 
https://api.bolagsverket.se/lamna-in-arsredovisning/v2.1/kontrollera/{token} (POST) 
https://api.bolagsverket.se/lamna-in-arsredovisning/v2.1/inlamning/{token} (POST) 
 
Tjänster för inlämning ESEF och/eller CSRD 
https://api.bolagsverket.se/lamna-in-arsredovisning-asynkront/v2.0/skapa-inlamningtoken (POST) 
https://api.bolagsverket.se/lamna-in-arsredovisning-asynkront/v2.0/inlamning/{token} (POST) 
https://api.bolagsverket.se/lamna-in-arsredovisning-asynkront/v2.0/valideringsresultat/{idnummer} 
(POST) 
 
Tjänster för årsredovisningshändelser 
https://api.bolagsverket.se/hantera-arsredovisningsprenumerationer/v2.0/handelseprenumeration/ 
(POST) 


---

## Sida 8

 
8 
 
https://api.bolagsverket.se/hantera-arsredovisningsprenumerationer/v2.0/handelseprenumeration/ 
(DELETE) 
https://api.bolagsverket.se/hantera-arsredovisningsprenumerationer/v2.0/handelseprenumeration/ (GET) 
https://api.bolagsverket.se/hamta-arsredovisningshandelser/v2.0/handelser/ (POST) 
 
Tjänster för statistik för årsredovisningsärenden 
https://api.bolagsverket.se/hamta-arsredovisningsstatistik/v1.0/statistik/allman/ankomstdatum (POST) 
https://api.bolagsverket.se/hamta-arsredovisningsstatistik/v1.0/statistik/allman/registreringstdatum 
(POST) 
https://api.bolagsverket.se/hamta-arsredovisningsstatistik/v1.0/statistik/programvara/ (POST) 
https://api.bolagsverket.se/hamta-
arsredovisningsstatistik/v1.0/statistik/programvara/separatrevisionsberattelse/ (POST) 
https://api.bolagsverket.se/hamta-arsredovisningsstatistik/v1.0/statistik/programvara/kontroller/ (POST) 
 
Stödtjänster för statistik för årsredovisningsärenden 
https://api.bolagsverket.se/hamta-arsredovisningsstatistik/v1.0/stodtjanster/programvaror/ (GET) 
 
 
 


---

## Sida 9

 
9 
 
 Åtkomst till testmiljö 
Eftersom tjänsterna i testmiljö enbart levererar statiska testdata behövs ingen 
autentisering av klienten för åtkomst. Bolagsverkets server accepterar anslutningar med 
https från alla klienter som släpps igenom brandväggen.  
 
4.1 Öppning av brandvägg 
Brandväggsöppningen beställs via api@bolagsverket.se. För att lägga beställningen 
behöver Bolagsverket veta vilken extern IP-adress (eller vilken extern IP-adressrange) 
som klientens testmiljö kommer ansluta ifrån. 
 
Som klient kan man göra följande rimlighetskontroller innan beställning: 
1. Säkerställa att IP-adressen/IP-adressrangen inte ligger inom en privat IP-
adressrange. En IP-adress är privat då den ligger inom något av följande intervall: 
From 
Tom 
10.0.0.0 
10.255.255.255 
172.16.0.0 
172.31.255.255 
192.168.0.0
192.168.255.255
 
2. Använd en webläsare i testmiljön för att gå till http://www.whatsmyip.org/ eller 
någon liknande tjänst. Den IP-adress som visas där ska ligga inom den IP-
adressrange som beställningen gäller. 
 
När Bolagsverkets kontaktperson meddelar att brandväggsöppningen är klar kan klienten 
verifiera det genom att göra en telnet-uppkoppling till: 
 Host: arsredovisning-accept2.bolagsverket.se 
 Port: 443 
 
Test kan till exempel göras med PuTTY. Om klienten släpps in i Bolagsverkets testmiljö 
visar Putty inga fel. 
 
Om brandväggsöppningen inte fungerat visar PuTTY (efter en stund) följande fel:  
  
 
Vid fel kontrollera först att korrekt host och port används, att trafiken släppts ut genom 
klientens brandväggar samt att klientens externa IP-adress stämmer med beställningen. 
Om det fortfarande inte fungerar, ta kontakt med kontaktperson på Bolagsverket för 
hjälp med felsökning. 
 


---

## Sida 10

 
10 
 
4.2 Trust av Bolagsverkets servercertifikat 
 
I testmiljö använder sig Bolagsverket av ett servercertifikat med följande DN: 
 
CN = *.BOLAGSVERKET.SE, OU = IT, O = Bolagsverket, L = SUNDSVALL, C = SE 
 
 
Bolagsverkets servercertifikat är ett ”domain validation certificate” utfärdat av 
Telia Sonera. För att kunna kommunicera med Bolagsverket via https måste klienten 
därför trusta certifikat utfärdare med följande rootcertifikat: 
 
CN = TeliaSonera Root CA v1, O = TeliaSonera 
 
4.2.1 Anslutning med webbläsare 
Eftersom vanliga webläsare, såsom Internet Explorer, Firefox och Chrome, har detta 
rootcertifikat förinstallerat, är det bra att som klient först verifiera att webläsaren kommer 
åt någon av tjänsterna i testmiljön. Äldre versioner av dessa webbläsare kanske inte har 
TeliaSoneras rootcertifikat installerat. 
 
4.2.2 Anslutning med klientprogramvara 
TeliaSoneras rootcertifikat är normalt inte bundlat i SDK:er för utveckling. Det finns till 
exempel inte med i Java JDK från Oracle. Certifikatet behöver därför laddas ner från 
TeliaSoneras webbplats och installeras/sparas så att det finns tillgängligt i den miljö som 
ska köra klientprogramvaran. Certifikatet kan hämtas på denna länk: 
https://repository.trust.teliasonera.com/teliasonerarootcav1.cer 
 
Skulle en klient försöka ansluta till testbänkarna eller tjänsterna utan att trusta rätt 
rootcertifikat kommer anslutningen misslyckas. Hur detta misslyckande manifesterar sig 
varierar mellan olika tekniska plattformar. Nedan följer ett exempel på en exception som 
indikerar detta fel för en klientanslutning implementerad i java: 
 
javax.net.ssl.SSLHandshakeException: sun.security.validator.ValidatorException: PKIX 
path building failed: sun.security.provider.certpath.SunCertPathBuilderException: 
unable to find valid certification path to requested target 
 
 
 


---

## Sida 11

 
11 
 
 Åtkomst till tjänster i acceptansmiljö 
Tjänsterna i acceptansmiljö levererar produktionslika testdata från Bolagsverkets 
testversioner av interna system. Därför behövs ytterligare autentisering och auktorisering 
av klienten ske för att ge åtkomst till tjänsterna. 
 
5.1 Öppning av brandvägg 
Ingen ytterligare brandväggsöppning behövs, förutsatt att klientens testmiljö är samma 
som den det beställts öppning för tidigare. (avsnitt 4.1) 
 
5.2 Trust av Bolagsverkets servercertifikat 
Samma servercertifikat som används i testmiljö enligt kapitel 4.2 används också i 
acceptansmiljö. 
 
5.3 Autentisering av klientens organisationscertifikat 
Klienten (och/eller eventuella webproxies hos klienten) måste konfigureras så att ett 
organisationscertifikat utfärdat av Expisoft/Steria skickas med i TLS-handskakningen 
med tjänsterna. Organisationscertifikatet ska vara utfärdat med rootcertifikat från 
Expisoft/Steria med något av nedanstående DN: 
 
 
CN=ExpiTrust Test CA v8,O=Expisoft AB,C=SE 
 
 
CN=ExpiTrust test CA v7,O=Expisoft AB,C=SE 
 
 
CN=ExpiTrust EID CA v4,O=Expisoft AB,C=SE 
 
 
CN=Steria AB EID CA v2,O=Steria AB,C=SE 
 
 
Formuleringen ”utfärdat med rootcertifkat” avser att det finns en obruten kedja av giltiga 
certifikat från organisationscertifikatet till rootcertifikatet med ovanstående DN. 
 
Organisationscertifikatet måste innehålla ett SERIALNUMBER med klientens 10-siffriga 
organisationsnummer prefixat med 16 i certifikatets DN, se exemplet nedan som visar 
Bolagsverkets egna organisationscertifikat utfärdat av Expisoft/Steria: 
 


---

## Sida 12

 
12 
 
 
Certifikatet måste också vara giltigt, så certifikatet i exemplet ovan fungerar inte efter 
2017-02-11. 
 
Om klienten redan har integrationer med myndigheters tjänster, till exempel Navet hos 
Skatteverket, är det sannolikt att klienten redan har ett organisationscertifikat utfärdat av 
Expisoft/Steria som går att använda även mot årsredovisningstjänsterna. Om inte, kan 
organisationscertifikat (eller serverlegitimation/organisationslegitimation som det kallas 
hos Expisoft) beställas via https://eid.expisoft.se/valj-elegitimation/. 
Serverlegitimationen måste minst ha användningssyftet ”identifiering av kund (vem som 
kopplar upp sig till en e-tjänst)” enligt information på beställningssidan. Ett certifikat som 
används i det syftet kallas också klientcertifikat. 
 
Vid korrekt uppsatt organisationscertifikathantering hos klienten ska en REST-request till 
någon av tjänsterna passera webproxyn och nå fram till tjänsten som då kommer svara 
med en 403 Forbidden tills klienten blivit auktoriserad enligt kapitel 5.4. 
 
Vid fel kontrollera först att åtkomst enligt kapitel 5.1 och 5.2 fortfarande fungerar, samt 
att ett organisationscertifikat enligt ovan skickas med i TLS-handskakningen mot 
Bolagsverkets tjänster. Om det fortfarande inte fungerar, ta kontakt med kontaktperson 
på Bolagsverket för hjälp med felsökning. 
 
5.4 Auktorisering av klienten 
Sista steget för åtkomst till tjänsterna är att auktorisera klienten. Auktorisering beställs via 
kontaktperson på Bolagsverket. Beställningen ska innehålla klientens 10-siffriga 
organisationsnummer, den görs i samband med att avtalet för tjänsten undertecknas. Det 
10-siffriga organisationsnumret måste vara detsamma som de 10 sista siffrorna i 
organisationsnumret i organisationscertifikatet. 
 


---

## Sida 13

 
13 
 
När Bolagsverkets kontaktperson meddelar att beställningen registrerats kan klienten 
verifiera det genom att köra samma REST-request som i kapitel 5.3 till någon av 
tjänsterna som nu ska svara med 200 OK. 
5.5 Utgående e-post 
Tjänsterna för digital inlämning i acceptansmiljö har, till skillnad från tjänsterna i 
testmiljö, förmågan att skicka e-post vid följande tillfällen: 
 Vid inlämning av årsredovisning till eget utrymme då den som bjudits in att 
skriva under fastställelseintyget kan notifieras via e-post. 
 Efter signering av fastställelseintyg och inlämning av årsredovisning till 
Bolagsverket då den som signerat kan välja att e-posta kvittensen till en eller flera 
e-postadresser. 
 
Bolagsverkets e-postfilter i acceptansmiljö släpper bara ut e-post till e-postadresser under 
domänerna gmail.com och hotmail.com, så för att kunna ta emot e-post från 
acceptansmiljön måste programvaruleverantören skaffa konto hos någon av dessa 
e-posttjänster. 
 
Vidare måste de exakta e-postadresserna vitlistas hos Bolagsverket. Vitlistning av 
e-postadresser beställs via api@bolagsverket.se. Beställningen ska innehålla de exakta 
gmail- eller hotmail-adresserna som ska vitlistas. 
 
 


---

## Sida 14

 
14 
 
 Åtkomst till tjänster i produktionsmiljö 
Tjänsterna i produktionsmiljö levererar riktiga data från Bolagsverkets interna system. 
Därför behövs ytterligare autentisering och auktorisering av klienten ske för att ge 
åtkomst till tjänsterna. 
 
6.1 Öppning av brandvägg 
Ingen brandväggsöppning behövs för anrop till produktionsmiljö. 
 
6.2 Trust av Bolagsverkets servercertifikat 
Samma servercertifikat som används i testmiljö enligt kapitel 4.2 används också i 
produktionsmiljö. 
 
6.3 Autentisering av klientens organisationscertifikat 
Klienten (och/eller eventuella webproxies hos klienten) måste även i produktionsmiljö 
konfigureras så att ett organisationscertifikat utfärdat av Expisoft/Steria skickas med i 
TLS-handskakningen med tjänsterna. Organisationscertifikatet ska vara utfärdat med 
rootcertifikat för Expisoft/Steria för produktionsbruk med något av nedanstående DN: 
 
 
CN=ExpiTrust EID CA v4,O=Expisoft AB,C=SE 
 
 
CN=Steria AB EID CA v2,O=Steria AB,C=SE 
 
6.4 Auktorisering av klienten 
Den beställning av åtkomst till tjänster som görs för acceptansmiljön gäller även för 
produktionsmiljön, se kapitel 5.4. Det krävs alltså ingen ytterligare beställning för att 
klienten ska få åtkomst till produktionsmiljön om den redan har åtkomst till 
acceptansmiljön. 
 
6.5 Utgående e-post 
Tjänsterna för digital inlämning i produktionsmiljö skickar e-post vid samma tillfällen som 
tjänsterna i acceptansmiljö. I produktionsmiljön finns inget e-postfilter, så där kommer 
e-post fram till alla e-postadresser utan vitlistning. 
 
