# avisering-digital-inlamning-arsredovisning-2.2

## Sida 1

Digital inlämning av 
årsredovisning 
 
Aviseringar och filformat  
Version 2.2 
 


---

## Sida 2

2 
Innehållsförteckning 
1 Inledning ..................................................................................................................................... 4 
 
 
 
 
 
 
2 Översikt ...................................................................................................................................... 4
3 Aviseringsfil ................................................................................................................................ 5
3.1 
HFR-filformatet .............................................................................................................. 5
3.1.1 
Posttyp 920 ........................................................................................................... 5
3.1.1.1 
ID för årsredovisningsfilen ...................................................................... 5
4 iXBRL ......................................................................................................................................... 5
4.1 
iXBRL-filen kommer in via inlämning av digitala årsredovisningar ....................... 5 
 
 
5 XBRL .......................................................................................................................................... 5
5.1 
Hur XBRL-filerna skapas av Bolagsverket ................................................................. 5
6 PDF ............................................................................................................................................. 6 
 
 
 
 
 
 
 
 
6.1 
Hur PDF:en skapas av Bolagsverket ........................................................................... 6 
6.2 
Version av PDF .............................................................................................................. 6
6.3 
Innehåll ............................................................................................................................ 6
6.3.1 
Uppmärkning ....................................................................................................... 6
7 Namnsättning och paketering ................................................................................................. 6
7.1 
Namnsättning K2 och K3 ............................................................................................. 6
7.2 
Namnsättning ESEF (K4)............................................................................................. 6
7.3 
Paketering ........................................................................................................................ 7
7.4 
Distribution via SFTP AB ............................................................................................. 7
7.5 
Distribution via SFTP övriga företagsformer än AB ................................................ 7 
 
 
 
 
 
7.5.1 
Tidpunkt för publicering .................................................................................... 7
7.5.2 
Plats för publicering AB ..................................................................................... 7
7.5.3 
Plats för publicering övriga objekttyper än AB ............................................... 7
7.5.4 
Lagringstid ............................................................................................................ 7
8 Exempelfiler ............................................................................................................................... 8


---

## Sida 3

3 
 Ändringshistorik 
Version 
Datum 
Beskrivning 
1.0 
2018-03-01 Första version 
1.0.1 
2018-03-21 - Rättning: namn på ZIP-filer (kap 7)
- Förtydligande: de digitala årsredovisningarna ligger
kvar i 10 dagar
2.0 
2019-10-15 Årsredovisning och Revisionsberättelse kan komma 
som separata aviseringsfiler. 
2.1 
2022-06-22 Korrigerat länkar under kap 8. 
2.2 
2022-12-02 Årsredovisning kan komma i en ZIP-fil samt separat 
fastställelseintyg. 


---

## Sida 4

4 
1 
Inledning 
Dokumentet beskriver de filer som Bolagsverket sprider via aviseringslösningen för 
årsredovisningar när dessa lämnas in digitalt. Om årsredovisningen och 
revisionsberättelsen har lämnats in separat kommer de att aviseras i separata filer.  
2 
Översikt 
Bilden visar endast de digitala årsredovisningarna. Dagens hantering av inskannade 
årsredovisningar som kommer in på papper är oförändrad: de filerna kommer även i 
fortsättningen att vara tiff och de kommer att publiceras på samma sätt som tidigare. 
Dokument-
arkiv
xHTML -> PDF-
konvertering
iXBRL
Årsredovisning 
Mottagning
PDF
iXBRL
XBRL
iXBRL -> XBRL-
konvertering
Originaldokument
Avskrift skapad av 
Bolagsverket
Systemkomponent
Spridning
Den iXBRL-fil som lämnas in för underskrift resulterar i olika spridning beroende på om 
den innehåller årsredovisning inkl. revisionsberättelse eller separata filer för dessa:  

iXBRL-filen själv, alltså den fil som lämnades in för underskrift

en XRBL-fil med samma data som iXBRL-filen, men utan presentationsdelar

en PDF-fil som motsvarar iXBRL-filen utseendemässigt (denna skapas inte för
separat revisionsberättelse)
Om årsredovisning och revisionsberättelse lämnats in separat blir det alltså fem filer. 
Dessa tre/fem filer kommer alltid att aviseras tillsammans. 


---

## Sida 5

5 
3 
Aviseringsfil 
3.1 HFR-filformatet 
3.1.1 Posttyp 920 
Posttyp 920 används för att beskriva metadata om en registrerad årsredovisning. 
3.1.1.1 
ID för årsredovisningsfilen 
Om årsredovisningen kom in på papper så sätts fältet HELFILMNR (populärt kallat 
rulle-ruta av historiska skäl). Detta nummer är unikt per årsredovisning. 
Om årsredovisningen kom in elektroniskt så sätts istället fältet KVITTENSNR. Detta 
nummer är också unikt per årsredovisning. 
Under 2007-2013 tog Bolagsverket emot årsredovisningar i XBRL-format. Dessa filer 
kom in elektroniskt och därför sattes KVITTENSNR, inte HELFILMNR. Fr.o.m 2018 
kommer de elektroniska årsredovisningarna in i formatet iXBRL. Varje post kommer att 
ha KVITTENSNR satt, inte HELFILMNR. Nummerserien kommer att vara unik även 
med avseende på de gamla elektroniska årsredovisningarna – vi kommer inte att 
återanvända numren som användes för XBRL-filerna utan den nya nummerserien börjar 
med högre nummer. 
4 
iXBRL 
4.1 iXBRL-filen kommer in via inlämning av digitala årsredovisningar 
De årsredovisningar som kommit in digitalt har verifierats så att de är giltig iXBRL som 
motsvarar de svenska taxonomierna för årsredovisning. De sprids i oförändrat skick till 
aviseringskunderna. 
De iXBRL-filer som Bolagsverket accepterar ska vara giltig XHTML. Bolagsverket 
kommer därför att förse dem med filändelsen .xhtml, då de flesta desktop-system är 
konfigurerade för att visa .xhtml-filer i webbläsare. 
5 
XBRL 
5.1 Hur XBRL-filerna skapas av Bolagsverket 
För varje iXBRL-fil skapar Bolagsverket en datafil i formatet XBRL som innehåller det 
taggade datat från iXBRL-filen. Eventuellt data som inte taggas i iXBRL-filen går förlorat, 
liksom information i bilder och övrig presentation. 
XBRL-filen skapas med den tredjepartskomponent som Bolagsverket använder för att 
verifiera att inlämnade årsredovisningsfiler är giltig iXBRL. Bolagsverket förbehåller sig 
rätten att byta ut intern implementation av konverteringsprogramvaran. XRBL-filen ska 
dock innehålla samma taggat data som iXBRL-filen oavsett vilken mjukvara som används 
för konvertering. 
Filerna kommer att få filändelsen .xbrl. 


---

## Sida 6

6 
6 
PDF 
6.1 Hur PDF:en skapas av Bolagsverket 
Bolagsverket använder en fristående komponent för att skapa PDF-avskrifter av inkomna 
iXBRL-filer. Dessa PDF-filer används för Bolagsverkets interna handläggning.  
Filerna kommer att få filändelsen .pdf. 
6.2 Version av PDF 
PDF-filerna är skapade i version 1.4 (Acrobat 5.x) av PDF-standarden. De är inte fullt ut 
kompatibla med PDF/A-standarden. 
6.3 Innehåll 
PDF:erna är textbaserade, inte bildbaserade. De enda bilder som förekommer i PDF:erna 
är eventuella bilder som länkats in i de iXBRL-filer som ligger till grund för PDF:erna. 
6.3.1 Uppmärkning 
PDF:erna får bokmärken som speglar heading-strukturen (dvs. <h>-taggar) i iXBRL-
filerna. Bolagsverket gör ingen egen märkning av ingående element utan det är endast de 
h-element som finns i iXBRL-filerna som blir bokmärken.
7 
Namnsättning och paketering 
7.1 Namnsättning K2 och K3 
De tre/fem filerna – iXBRL (originalfilen som inkom till Bolagsverket), XBRL (innehåller 
det taggade datat i iXBRL-filen, skapas av Bolagsverket) och PDF (motsvarar iXBRL-
filens utseende, skapas av Bolagsverket) – namnsätts med <kvittensnummer>.<filtyp>, 
där kvittensnummer motsvarar numret i fältet KVITTENSNR i aviseringsfilen.  
Exempel: 
Aviseringsfilen har KVITTENSNUMMER 6100000022.
De tre filerna får namnen: 6100000022.xhtml, 6100000022.xbrl och 6100000022.pdf.
Om revisionsberättelsen har lämnats in separat får dessa filer namnen:
6100000022_RB.xhtml och 6100000022_RB.xbrl.
7.2 Namnsättning ESEF (K4) 
ESEF-årsredovisningar har en s.k. utökad taxonomi där det innebär att alla kunder gör 
unika regelverk. Årsredovisningen inkommer i ett taxonomipaket, en ZIP, där 
årsredovisningen är en iXBRL-fil (filändelse XHTML). Bolagsverket skapar en XBRL av 
den. Fastställelseintyg ska också skickas in och det inkommer i iXBRL-format där 
Bolagsverket skapar en XBRL-fil. 
Exempel: 
Aviseringsfilen har KVITTENSNUMMER 6100000022.
Filerna får namnen: 6100000022.zip, 6100000022.xhtml och 6100000022.xbrl.
Fastställelseintyget får namnen: 6100000022_FI.xhtml och 6100000022_FI.xbrl.


---

## Sida 7

7 
7.3 Paketering 
Alla filer som hör till digitalt inkomna årsredovisningar gällande AB (inkl. filer med 
revisionsberättelser) kommer att packas i en ZIP-fil med namnet 
Arsredovisning_digital_<ååmmdd>.zip, t.ex. Arsredovisning_digital_180112.zip.  
Digitalt inkomna årsredovisningar som inte är AB (i detta fall endast BAB och FAB) 
kommer att packas i en ZIP-fil med namnet:  
arsredHFR_digital_<ååååmmdd>.zip t.ex. arsredHFR_digital_20221201.zip. 
7.4 Distribution via SFTP AB 
De digitala årsredovisningarna gällande AB kommer att distribueras på samma sätt som 
de skannade årsredovisningarna: dygnsvis publicering på Bolagsverkets SFTP-server. De 
kunder som hämtar aviseringar av årsredovisningsdokument för AB hämtar från samma 
yta, detta eftersom det är så stora volymer av dokument. Där finns två filer gällande AB, 
scannade och digitala årsredovisningar. 

Arsredovisning_digital_ÅÅMMDD.zip

Arsredovisning_scannade_ÅÅMMDD.zip
7.5 Distribution via SFTP övriga företagsformer än AB 
För övriga företagsformer än AB paketeras även dessa filer i en ZIP med dessa namn. 

arsredHFR_digital_ÅÅÅÅMMDD.zip

arsredHFR_ÅÅÅÅMMDD.zip
7.5.1 Tidpunkt för publicering 
Till skillnad från de skannade årsredovisningarna så kommer jobbet som paketerar och 
publicerar de digitala årsredovisningarna att köras efter midnatt. Orsaken till det är att de 
digitala årsredovisningarna kan komma in dygnet runt. Skannade årsredovisningar 
bearbetas endast under kontorstid. 
7.5.2 Plats för publicering AB 
ZIP-filen med de digitala årsredovisningarna kommer att läggas i samma katalog som de 
skannade årsredovisningarna, dvs. ”arsredovisningar/<ååmmdd>”, t.ex. 
arsredovisningar/180112.  
7.5.3 Plats för publicering övriga objekttyper än AB 
För övriga företagsformer läggs en fil på respektive kunds eget SFTP-konto detta gäller 
pappers/epost årsredovisningar. ZIP-filen med de digitala årsredovisningarna gällande 
övriga objekttyper än AB kommer att läggas i samma katalog som de skannade 
årsredovisningarna för ej AB. 
7.5.4 Lagringstid 
De digitala årsredovisningarna kommer att ligga kvar 10 dagar på SFTP-servern. 


---

## Sida 8

8 
8 
Exempelfiler 
Exempelfiler för årsredovisningar och revisionsberättelser finns publicerade på 
exempelsidan på taxonomier.se. Ett exempel länkas till nedan: 

iXBRL

XBRL

PDF
