# Hardware

## Anforderungen

Das Radio soll dem Benutzer folgende Hardwarekonfigurationsmöglichkeiten bieten:

* Anschluss passive Lautsprecher/Kopfhörer möglich 
* Lautstärkeregelung über Hardware möglich
* Verwendung des internen Lautsprechers des alten Radios
* Statusinformationen zum aktuellen Lied beispielsweise über ein LCD
* LEDs als Statusanzeige und/oder als Visualisierungsvariante von Musik (vgl.
  [@wood2005techniques])
* USB--Anschlussmöglichkeit für externe Datenträger

## Komponenten und Bauteile

Abbildung \ref{uebersicht} zeigt eine konzeptuelle Übersichts des Zusammenspiels
der einzelnen Komponenten.

\begin{figure}[h!]
  \centering
  \includegraphics[width=0.8\textwidth]{images/uebersicht.png}
  \caption{Grobe Übersicht der verwendeten Komponenten im Zusammenspiel}
  \label{uebersicht}
\end{figure}

Folgende Hardwarekomponenten oder Bauteile waren bereits vorhanden oder mussten
noch erworben werden:

**Vorhanden:**

* Altes Gehäuse AEG 4104 Küchenradio[^AEG] 
* *Raspberry Pi* aus dem Jahr 2012 
* LCD--Anzeige (Altbestände u. Arduino--Kit)
* Kleinbauteile wie LEDs, Widerstände
* USB--Hub für Anschluss von beispielsweise ext. Festplatte 
* USB--Soundkarte 
* Wi--Fi--Adapter
* Netzteil (diverse 5V, 2A)

**Mussten noch erworben werden:**

* Audioverstärker (\EUR{6})
* Drehimpulsregler (\EUR{3})
* Kunststoffabdeckung für Front (\EUR{0})
* Farbe (Lack) (\EUR{5})
* Drehknöpfe für das Gehäuse (\EUR{2})
* Schrumpfschläuche (\EUR{1})
* Kippschalter (\EUR{6})

Insgesamt wurden für das Projekt ca. \EUR{23} ausgegeben.

[^AEG]: AEG Küchenradio 4104: \url{https://www.amazon.de/AEG-MR-4104-Desgin-Uhrenradio-buche/dp/B000HD19W8}

## Raspberry Pi

Der vorhandene *Raspberry Pi* ist aus dem Jahr 2012. Die genaue CPU-- und
Board--Revision kann auf Linux unter ``proc`` ausgelesen werden. Von der
Hardware--Revision kann auf die *Raspberry Pi*--Revision geschlossen werden,
siehe hierzu auch [@gay2014raspberry], Seite 46:

```bash
    $ cat /proc/cpuinfo 
    processor       : 0
    model name      : ARMv6-compatible processor rev 7 (v6l)
    BogoMIPS        : 697.95
    Features        : half thumb fastmult vfp edsp java tls 
    CPU implementer : 0x41
    CPU architecture: 7
    CPU variant     : 0x0
    CPU part        : 0xb76
    CPU revision    : 7

    Hardware        : BCM2708
    Revision        : 0003
    Serial          : 00000000b8b9a4c2
```

Laut Tabelle unter [@gay2014raspberry], Seite 45 handelt es sich hierbei um das
Modell B Revision 1+ mit 256MB RAM.

Je nach Raspberry Revision sind die Pins teilweise unterschiedlich belegt. Seit
Modell B, Revision 2.0 ist noch zusätzlich der P5 Header dazu gekommen.
Abbildung \ref{gpio}[^SRC] zeigt die GPIO--Header des *Raspberry Pi* Modell B Revision
1+.

### GPIO--Schnittstelle

\begin{figure}[h!]
  \centering
  \includegraphics[width=0.7\textwidth]{images/gpio.png}
  \caption{GPIO-Header des Raspberry Pi Modell B Rev 1.0+}
  \label{gpio}
\end{figure}


[^SRC]: Bildquelle:
\url{http://www.raspberrypi-spy.co.uk/2012/06/simple-guide-to-the-rpi-gpio-header-and-pins/\#prettyPhoto}


#### GPIO--Pinbelegung und Funktionalität

Die GPIO--Pins des *Raspberry Pi* haben eine Logikspannung von 3.3V und sind pro
GPIO--Pin mit max. 16mA belastbar. Der gesamte GPIO--Header sollte mit nicht
mehr als 50mA belastet werden, da es darüber hinaus zu Hardwareschäden kommen
kann (vgl. [@gay2014raspberry], Seite 121 ff.).

Die *Logik--Pegel* der GPIO--Pins sind beim
*Raspberry Pi* wie folgt definiert [vgl. @gay2014raspberry], Seite 129 ff.:

* $\le$ 0,8V, input low
* $\ge$ 1,3V, input high


Die Ansteuerung von LEDs über GPIO erfolgt binär. Das heißt, dass
die LED entweder aus (low) oder an (high) sein kann.

In der »analogen« Welt ist es jedoch möglich eine LED über das Senken der
Spannung zu dimmen. Um ein Dimmen in der digitalen Welt zu erreichen wird ein
Modulationsverfahren angewandt, welches Pulsweitenmodulation (*PWM*) heißt.
Unter [@richardson2014make], Seite 121 ff. und [@gay2014mastering], Seite 421
ff. finden sich weitere Informationen. Software PWM unter
[@gay2014experimenting], Seite 183 ff. zeigt beispielsweise eine 6% CPU--Last
pro GPIO--Pin bei einer PWM--Softwareimplementierung. Wie PWM in *Eulenfunk*
eingesetzt wird im Software--Kapitel unter dem Punkt \ref{internal-led-treiber}
beleuchtet.

## LCD--Anzeige

Um dem Benutzer --- beispielsweise Informationen über das aktuell gespielte Lied
--- anzeigen zu können, soll eine LCD--Anzeige verbaut werden. In den privaten
Altbeständen finden sich folgende drei Hitachi--hd44780--kompatible Modelle:

* Blaues Display, 4x20 Zeichen, Bolymin BC2004A
* Blaues Display, 2x16 Zeichen, Bolymin BC1602A
* Grünes Display, 4x20 Zeichen, Dispalytech 204B 

Für *Eulenfunk* wurde das blaue 4x20 Display --- aufgrund der größeren
Anzeigefläche und Farbe --- gewählt.

### Anschlussmöglichkeiten

Eine LCD--Anzeige kann an den *Raspberry Pi* auf verschiedene Art und Weise
angeschlossen werden. Anschlussmöglichkeiten für eine LCD--Anzeige wären
beispielsweise: 

* GPIO direkt (parallel)
* I2C--Bus (seriell)
* SPI--Bus (seriell)

Die serielle Anschlussmöglichkeit bietet den Vorteil, dass weniger Datenleitungen
(GPIO--Pins) verwendet werden. Für den parallelen Betrieb des Displays werden
mindestens sechs GPIO--Pins benötigt, für den seriellen Anschluss über I2C
lediglich nur zwei. 

Da für den seriellen Betrieb beispielsweise über den I2C--Bus zusätzliche
Hardware benötigt wird, wird die parallele Ansteuerung über die GPIO--Pins
bevorzugt. Weitere Informationen zum seriellen Betrieb über I2C sind unter 
[@horan2013practical], Seite 61, ff. zu finden.

\begin{figure}[h!]
  \centering
\includegraphics[width=0.7\textwidth]{images/lcdraspi.png}
  \caption{Verdrahtung von LCD im 4--Bit Modus und Raspberry Pi, alle hierzu benötigten
  Informationen sind im Datenblatt zu finden.}
  \label{lcd}
\end{figure}

Das Display arbeitet mit einer Logik--Spannung von 3.3V -- 5V. Da die GPIO--Pins
jedoch eine High--Logik von 3,3V aufweisen, würde man hier in der Regel einen
Pegelwandler bei bidirektionaler Kommunikation und 5V benötigen. Da aber auf das
Display nur zugegriffen und die GPIO--Pins nicht schreibend benutzt werden, kann
ein Betrieb des Displays auch mit 5V erfolgen. Beim 3.3V Betrieb welcher laut
Datenblatt[^LCD] auch möglich sein soll, hatte das Display leider nur eine sehr
schwachen beziehungsweise unzureichenden Darstellungskontrast, weswegen der 5V
Betrieb gewählt wurde. Zudem wurde an *Pin3* (LCD) ein 100$\Omega$ Potentiometer
hinzugefügt. Dies ermöglicht den Kontrast variabel einzustellen.

Die Hintergrundbeleuchtung des Displays wurde direkt über ein Potentiometer mit
2K$\Omega$ an die 5V Spannungsversorgung angeschlossen. Es wurde hier die
direkte Speisung vom Netzteil gewählt, um den GPIO--Header nicht unnötig zu belasten.

Laut Datenblatt kann die Hintergrundbeleuchtung entweder mit 3.4V ohne
Vorwiderstand oder mit 5V bei einem 27$\Omega$ Widerstand betrieben werden. Damit das
Display beim herunter geregeltem Potentiometer keinen Schaden nimmt, wurden
zusätzlich zwei Widerstände mit 100$\Omega$ (parallel geschaltet = 50$\Omega$) zwischen Display
und Potentiometer gehängt. Abbildung \ref{lcd} zeigt die Verdrahtung des LCD mit
dem *Raspberry Pi*.

Der resultierende Gesamtwiderstand ohne Potentiometer beträgt in diesem Fall
$\approx$ 50 $\Omega$:

$$  R_{ges} = \frac{R_1 \times R_2}{R_1 + R_2} = \frac{100\Omega \times 100\Omega}{100\Omega + 100\Omega} = 50\Omega $$

## Drehimpulsgeber

Um eine minimale Anzahl an Bedienelementen zu gewährleisten, wird bei *Eulenfunk*
ein Drehimpulsgeber mit Schalter gewählt. Für erste Testzwecke wurde von der
Hochschule ein *ALPS STEC12E08* bereitgestellt. Dieser wurde im Laufe der
Entwicklung durch einen *ALPS STEC11B09*[^ALPS] ersetzt, da dieser mittels Mutter und
Schraube am Gehäuse stabiler befestigt werden kann. 

Der verwendete Drehimpulsgeber hat insgesamt fünf Anschlüsse. Zwei
Signalleitungen (A und B), zwei mal *GND* (jeweils für Drehgeber und Schalter)
und einen Anschluss für den Schalter. Beim Drehen eines Drehimpulsgebers wird
ein Rechtecksignal generiert. Je nach Muster der beiden Datensignale A oder B,
kann entschieden werden ob es sich um eine Rechts-- oder Linksdrehung handelt.
Siehe [@2014projekte], Seite 361 ff. für weitere Hintergrundinformationen zum
Drehimpulsgeber.

Abbildung \ref{alps} zeigt den Anschluss des Drehimpulsgebers am *Raspberry Pi*. 


\begin{figure}[h!]
  \centering
\includegraphics[width=0.8\textwidth]{images/rotary.png}
  \caption{Drehimpulsgeber--Anschluss an den Raspberry Pi, Abbildung zeigt
  Kombination aus Potentiometer und Schalter.}
  \label{alps}
\end{figure}

\newpage

## Soundkarte

Die interne Soundkarte des *Raspberry Pi* ist über eine simple 
Pulsweitenmodulation realisiert. Die einfache Schaltung soll hier laut
Internetquellen[^AQ]eine sehr niedrige Audioqualität bieten.

Aus diesem Grund wird bei *Eulenfunk* auf das USB--Audio--Interface *BEHRINGER
U-PHONO UFO202*[^DAC] (USB--Soundkarte) gesetzt. 

[^LCD]: Datenblatt Bolymin BC2004A: \url{http://www.dema.net/pdf/bolymin/BC2004A-series_VER04.pdf}
[^ALPS]: Drehimpulsgeber ALPS STEC11B09: \url{https://www.reichelt.de/Drehimpulsgeber/STEC11B09/3/index.html?ACTION=3&GROUPID=3714&ARTICLE=73915}
[^DAC]:BEHRINGER U-PHONO UFO202 Audio Interface: \url{http://www.produktinfo.conrad.com/datenblaetter/1300000-1399999/001370864-an-01-de-BEHRINGER_UFO_202_AUDIOINTERFACE.pdf}
[^AQ]: Raspberry Pi onboard Sound: \url{http://www.crazy-audio.com/2013/11/quality-of-the-raspberry-pi-onboard-sound/}

## Audioverstärkermodul

Da eine Soundkarte in der Regel zu wenig Leistung hat, um einem Lautsprecher
»vernünftig« anzusteuern, wird ein Audioverstärker benötigt. Da neben dem
Anschluss von externen Lautsprechern auch eine Lautstärkeregelung über ein
Potentiometer erfolgen soll, ist die Entscheidung einfachheitshalber auf ein
Audioverstärker--Modul auf Basis vom PAM8403[^POW] Stereo-Verstärker mit
Potentiometer gefallen. Eine Do--It--Yourself--Alternative wäre ein
Transistor--basierter Audio--Verstärker, hier gibt es online diverse
Bauanleitungen[^amp2].

[^amp2]: Transistor--Verstärker:
\url{http://www.newsdownload.co.uk/pages/RPiTransistorAudioAmp.html}

Das Audioverstärker--Modul hat folgende Anschlusspins:

* Left--In, Right--In, GND
* 5V+ und GND (Betriebsspannung)
* Left--Side--Out (+), Left--Side--Out (-)
* Right--Side--Out (+), Right--Side--Out (-)

Laut diverser Onlinequellen[^MONO], dürfen die Ausgänge für einen Mono--Betrieb
eines auf dem PAM8403--basierten Verstärkers nicht parallel geschaltet werden.
Aus diesem Grund kommt ein 4--poliger
*EIN--EIN--Kippschalter*[^KIPPSCHALTER] zum Einsatz. So kann zwischen dem
internen Lautsprecher (Mono--Betrieb) und den externen Stereo
Lautsprecher--Anschlüssen sauber per Hardware hin und her geschaltet werden.

Damit im Mono--Betrieb nicht nur ein Kanal verwendet wird, ermöglicht
*Eulenfunk* das Umschalten zwischen Mono-- und Stereo--Betrieb in Software.

[^MONO]: PAM8403 Mono--Betrieb:
\url{http://electronics.stackexchange.com/questions/95743/can-you-bridge-or-parallel-the-outputs-of-the-pam8403-amplifier}

[^KIPPSCHALTER]: Kippschalter 4--polig EIN--EIN: \url{http://www.reichelt.de/Kippschalter/MS-500P/3/index.html?&ACTION=3&LA=2&ARTICLE=13172&GROUPID=3275&artnr=MS+500P}
[^POW]: Verstärkermodul: \url{https://www.amazon.de/5V-Audioverstärker-Digitalendstufenmodul-Zweikanalige-Stereo-Verstärker-Potentiometer/dp/B01ELT81A6}

## LED--Transistorschaltung

Die Ansteuerung einer LED mittels GPIO--Pin ist recht simpel. Sollen jedoch
mehrere LEDs angesteuert werden, so wird in der Regel pro LED ein GPIO--Pin
benötigt. LEDs sollten nie ohne Vorwiderstand an den *Raspberry Pi*
angeschlossen werden, da durch den hohen Stromfluss die LED beschädigt werden
könnte. Weiterhin muss bei LEDs auch auf die Polung geachtet werden, die
abgeflachte Seite --- meist mit dem kürzerem Beinchen --- ist in der Regel die
Kathode (Minuspol). Abbildung \ref{led} zeigt exemplarisch den Anschluss einer
*classic LED rot*[^LEDS], mit einer Flussspannung von $U_{LED}$ $\approx$ 2V, die mit
einem Strom von $I_{LED}$ = 20 mA gespeist werden soll. Die Berechnung des
Vorwiderstandes erfolgt nach folgender Formel:

$$R_{LED} = \frac{U_{GPIO}-U_{LED}}{I_{LED}} = \frac{3.3V - 2V}{20mA}   \approx 65\Omega$$

[^LEDS]: Datenblatt mit verschiedenen LED--Typen: \url{https://www.led-tech.de/de/5mm-LEDs_DB-4.pdf}

**Hinweis:** Da ein GPIO--Pin aber mit nur max. 16mA belastet werden sollte,
sollte in unserem Beispiel durch 16mA anstatt 20mA geteilt werden um den max.
Stromfluss auf 16mA zu begrenzen. In diesem Fall würden wir auf $\approx$ 82$\Omega$ kommen.

Da Widerstände meistens in fest vorgegebenen Größen vorhanden sind, kann im Fall
eines nicht exakt existierenden Widerstandswertes einfach der nächsthöhere
Widerstandswert genommen werden. Im Beispiel wird ein $100\Omega$ Widerstand
verwendet. 

Weitere Beispiele und Grundlagen zur Reihen-- und Parallelschaltung von LEDs
können online beispielsweise unter *led-treiber.de*[^LED] eingesehen werden.

\begin{figure}[h!]
  \centering
\includegraphics[width=0.5\textwidth]{images/led.png}
  \caption{Anschluss einer roten LED mit Vorwiderstand am Raspberry Pi GPIO--Pin}
  \label{led}
\end{figure}

Je nach Typ und Farbe ist der benötigte Strom um ein vielfaches höher wie in
unserem Beispiel. Die in \ref{led} abgebildete LED kann vom GPIO--Pin nur einen
max. Strom von 16 mA beziehen

In *Eulenfunk* sollen mehrere intensiv leuchtende LEDs verbaut werden. Da die
GPIO--Pins in ihrer Leistung sehr begrenzt sind, würde es sich anbieten eine
externe Stromquelle zu verwenden. Um die Speisung über eine externe Stromquelle
zu ermöglichen, kann eine Transistorschaltung verwendet werden (vgl. [@exploring],
Seite 219 ff.). 

Für die Transistorschaltung wurden von Seite der Hochschule Augsburg NPN-- (BC547C) und
PNP--Transistoren (BC557C) bereitgestellt. Für den ersten Testaufbau wurde der
PNP--Transistor und eine RGB--LED[^RGBGM] mit gemeinsamen Minuspol verwendet.
Dabei ist aufgefallen, dass die LED ständig geleuchtet hat. Eine kurze Recherche
hat ergeben, dass der Transistor permanent durchgeschaltet war, weil die
Spannung an der Basis (GPIO--Pin, 3,3V) geringer war als die Betriebsspannung
für die LED (5V). 

Der zweite Testaufbau mit dem NPN--Transistor BC547C und einer RGB--LED[^RGBGP] mit
gemeinsamen Pluspol hat das gewünschte Ergebnis geliefert.

Da der Hersteller für die von der Hochschule bereitgestellten Transistoren
unbekannt ist, wurden typische Durchschnittswerte für die Dimensionierung der
restlichen Bauteile verwendet.

Wie es aussieht sind die meisten BC547C Transistor--Typen für einen max. Strom
$I_{CE}$=100 mA konstruiert. Für die Berechnung des Basis--Vorwiderstandes wird
der Stromverstärkungsfaktor $h_{FE}$[^HFE] benötigt. Je nach Hersteller
variieren die Werte zwischen 200[^FARI] und 400[^SEM]. Da der maximale Laststrom
$I_{CE}$ pro Transistor 60 mA (3 LEDs je max. 20mA) beträgt, sieht die
Berechnung des Basisstroms --- bei einem durchschnittlichem $h_{FE}$ = 300 --- wie folgt
aus:

[^LED]: Beispiele zur Ansteuerung von LEDs: \url{http://www.led-treiber.de/html/vorwiderstand.html}
[^HFE]:Stromverstärkungsfaktor: \url{http://www.learningaboutelectronics.com/Articles/What-is-hfe-of-a-transistor}
[^RGBGM]: RGB-LED Common Cathode: \url{http://download.impolux.de/datasheet/LEDs/LED 0870 RGB 5mm klar 10000mcd.pdf}
[^RGBGP]: RGB-LED Common Anode: \url{http://download.impolux.de/datasheet/LEDs/LED 09258 RGB 5mm klar 10000mcd_GP.pdf}


 $$I_{Basis} = \frac{I_{CE}}{h_{FE}} = \frac{0.06A}{300} \approx 200\mu A$$

Der BC547C Transistor benötigt eine durchschnittliche  $U_{BE}$ = 0,7V zum
Durchschalten. Die GPIO-Pins des *Raspberry Pi* haben einen Spannungspegel von
3.3V. Daraus ergibt sich folgende Berechnung des Basis--Vorwiderstandes:

$$R_{Basis} = \frac{U_{GPIO} - U_{Basis}}{I_{Basis}} = \frac{3,3V - 0,7V}{200\mu A} = 13k\Omega $$

[^SEM]: SEMTECH: \url{http://pdf1.alldatasheet.com/datasheet-pdf/view/42386/SEMTECH/BC547.html}
[^FARI]: Farichild Semiconductor: \url{https://www.fairchildsemi.com/datasheets/BC/BC547.pdf}

\begin{figure}[h!]
  \centering
\includegraphics[width=1\textwidth]{images/transistorled.png}
  \caption{Transistor--RGB--LED Schaltung}
  \label{transled}
\end{figure}

Damit der Transistor jedoch *sicher* durchschaltet, werden Wiederstände mit $10k
\Omega$ verwendet. Die in Abbildung \ref{transled} gelisteten
LED--Vorwiderstände ergeben sich aufgrund der verschiedenen Spannungen der
unterschiedlichen Farben[^RGBGP]. Die Berechnung für den Vorwiderstand pro LED
schaut am Beispiel der Farbe blau ($U_{LED} = 3,15V, I_{LED} = 20mA$) wie folgt
aus:

$$R_{LED} = \frac{U_{Betriebsspannung} - U_{LED}}{I_{LED}} = \frac{5V - 3,15V}{20mA} =92.5\Omega \approx 100\Omega$$

Analog errechnet sich für die Farbe rot ein Vorwiderstand von $145\Omega$ und für grün
ein Vorwiderstand von $87\Omega$.

## USB--Hub und Netzteil

Der *Raspberry Pi* hat in unserer Revision nur zwei USB--Schnittstellen, diese
sind bereits durch die Hardware--Komponenten USB--DAC (Soundkarte) und das
Wi--Fi--Modul belegt. Um den Anschluss eines externen Datenträgers, auch mit
größerer Last wie beispielsweise einer Festplatte zu ermöglichen, wird ein
aktiver USB--Hub benötigt.

Für diesen Einsatzzweck wird aus den Altbeständen ein *LogiLink 4 Port USB 2.0
HUB*[^HUB] verwendet. Viele billige Hubs arbeiten hier entgegen der USB--Spezifikation
und speisen den *Raspberry Pi* zusätzlich über die USB--Schnittstelle. Dieses
Verhalten wurde bemerkt, als der *Raspberry Pi* ohne Power--Connector alleine
nur mit der USB--Verbindung zum USB--Hub bootete.

[^HUB]: LogiLink USB--Hub: \url{https://www.amazon.de/LogiLink-4-Port-Hub-Netzteil-schwarz/dp/B003ECC6O4}

Bei der Speisung über die USB--Schnittstelle wird die interne Sicherungsschaltung
des *Pi* umgangen, deswegen wird in der Regel von einem Betrieb eines USB--Hub
mit *backfeed* abgeraten (vgl . [@suehle2014hacks], Seite 26 ff.). Für den Prototypen wird
jedoch der genannte USB--Hub und das dazugehörige Netzteil für den Betrieb von
*Eulenfunk* verwendet. Das Netzteil ist für 5V bei max. 2A ausgelegt.

**Nachtrag:** Die Speisung über das Netzteil des USB--Hubs ist recht instabil. Bei
Lastspitzen kommt es anscheinend zu Störeinwirkungen, die sich auf die
GPIO--Peripherie auswirken (LCD--Anzeige rendert inkorrekt). Ein weiterer Punkt
sind Störfrequenzen, welche teilweise in Form von Störgeräuschen die
Audioausgabe überlagern (Hintergrundgeräusche beim Einschalten aller LEDs).
Insgesamt wurden drei Netzteile --- jeweils 5V, 2A ---ausprobiert. Von diesen
war lediglich ein einziges als »akzeptabel« einzustufen. Die restlichen zwei
führen bei Lastspitzen zu Problemen (Abstürze, fehlerhaftes Rendering auf
Display, GPIO--Flips, et cetera). Das *backfeed* des USB--Hubs scheint die
genannten Probleme teilweise zu verstärken (vgl . [@suehle2014hacks], Seite 27).


## Gehäuse

### Vorderseite

Abbildung \ref{ral} zeigt ein Muster der Gehäusefront--Farbe hellelfenbeinweiß RAL
1015. Dieser Farbton wird für die Front verwendet, um *Eulenfunk* einen dezenten
»Retro«--Look zu verpassen.

\begin{figure}[h!]
  \centering
  \includegraphics[width=0.3\textwidth]{images/ral_soft.png}
  \caption{Muster RAL1015, hellelfenbeinweiß}
  \label{ral}
\end{figure}

Das Plexiglas für die Front wurde von der Firma *ira-Kunststoffe* in
Schwarzenbach/Saale zugeschnitten. In der Plexiglasfront wurden mit Hilfe von
Herrn Schäferling zwei 5mm Löcher (Drehimpulsgeber, Lautstärkeregler--Poti)
gebohrt. Anschließend wurde die Plexiglas--Front von der Innenseite
lackiert[^LACK], hierbei wurden die Flächen für LCD und die drei LEDs abgeklebt.
Zudem werden schwarze Knöpfe in Alu--Optik mit $\diameter$ 30mm  für den
Lautstärkeregler und den Drehimpulsgeber verwendet.

### Rückseite

Für die Rückseite wird die alte Abdeckung des AEG--Radios verwendet. Diese
musste teilweise leicht modifiziert werden. An dieser befinden sich zwei Potis
für Kontrastregelung und Hintergrundbeleuchtung des LCD, eine
USB--Female--Kabelpeitsche, zwei Cinch Stecker für externe Lautsprecher und ein
Kippschalter zum Umschalten zwischen internen und externen Lautsprechern.

[^LACK]: Buntlack, hellelfenbein: \url{http://www.obi.de/decom/product/OBI_Buntlack_Spray_Hellelfenbein_hochglaenzend_150_ml/3468725}


## Betriebssystem

Mittlerweile gibt es für den *Raspberry Pi* viele offiziell zugeschnittene
Betriebssysteme (vgl. [@pietraszak2014buch], Seite 29 ff. und [@warner2013hacking],
Seite 47 ff.). Bei den Linux--Distributionen ist *Raspbian* eine der bekanntesten Distribution -- welche auf
*Debian* basiert. *Raspbian* bringt ein komplettes Linux--basiertes System mit
grafischer Benutzeroberfläche mit sich. 

Neben den unter [@pietraszak2014buch], Seite 29 ff. genannten Distributionen
gibt es mittlerweile auch Windows 10 IoT (Internet of Things) für den *Raspberry
Pi*. Dieses speziell für den Embedded Bereich ausgerichtete Windows benötigt
jedoch eine ARMv7--CPU als Mindestanforderung[^WINIOT], was den »alten Raspberry«
ausschließt. Außerdem wäre für uns eine proprietäre Lösung ein
K.O.--Kriterium, da diese alle Vorteile von Freier Software zunichte machen
würde.

[^WINIOT]: Systemanforderungen:\url{http://raspberrypi.stackexchange.com/questions/39715/can-you-put-windows-10-iot-core-on-raspberry-pi-zero}

### Wahl des Betriebssystem

*Arch Linux ARM*[^ARCH] ist eine minimalistische und sehr performante
Linux--Distribution welche im Gegensatz zu *Raspbian* ohne Desktop--Umgebung
geliefert wird (vgl. [@schmidt2014raspberry], Seite 13 ff.) Darüber hinaus ist
*Arch Linux* ein bekannter Vertreter von Rolling--Release--Distributionen. Ein
weiterer Vorteil für unseren Einsatzzweck ist bei *Arch Linux* das *AUR*
(Arch User Repository)[^AUR], dieses erlaubt es eigene Software auf eine
schnelle und unkomplizierte Weise der Allgemeinheit zur Verfügung zu stellen.

### Einrichtung des Grundsystems

Nach der Installation[^INSTALL] und dem ersten Booten des Grundsystems muss die
Netzwerk--Schnittstelle konfiguriert werden. Arch Linux ARM bietet mit *netctl*
eine Profil--basierte Konfigurationsmöglichkeit. Ein Profil kann über das
*ncurses*--basierte Tool `wifi-menu` erstellt werden. In unserem Fall wurde das
Profil `wlan0-Phobos` erstellt. Anschließend kann das erstellte Profil mit
*netctl* verwendet werden. 

\newpage 

**Auflistung der bekannten Profile**

```bash
    [root@eulenfunk ~]$ netctl list
      eth0-static
      wlan0-Phobos
```

**Aktivierung des gewünschten Profils**

```bash
    # Starten des gewünschten Profils
    [root@eulenfunk ~]$ netctl start wlan0-Phobos
    [root@eulenfunk ~]$ netctl list
      eth0-static
    * wlan0-Phobos

    # Profil über System-Reboot hinweg aktivieren 
    [root@eulenfunk ~]$ netctl enable wlan0-Phobos

```

Nun verbindet sich der *Raspberry Pi* nach dem Hochfahren jedes Mal automatisch
mit dem Profil `wlan0-Phobos`.

[^AUR]: Arch User Repository: \url{https://aur.archlinux.org/}
[^ARCH]: Arch Linux ARM: \url{https://archlinuxarm.org/}
[^INSTALL]: Arch Linux Installation für Raspberry Pi: https://archlinuxarm.org/platforms/armv6/raspberry-pi#installation

## Erweiterungen und alternative Ansätze

### Allgemein

Der aktuelle Prototyp hat lediglich nur ein Potentiometer um die Hintergrundbeleuchtung
des LCD zu regeln. Ein anderer Ansatz wäre der Einsatz eines Relais, welches es
ermöglichen würde die LCD--Hintergrundbeleuchtung Software--seitig ein-- und auszuschalten.
Die Software könnte dann automatisch nach längerer Inaktivtät die Beleuchtung dimmen.

\label{internal-audio-vis}

### Audio--Visualisierung

Beim Projekt *Eulenfunk* wird die Visualisierung von Musik aufgrund der
begrenzten Zeit und Hardwareressourcen des *Raspberry Pi *über eine
vorberechnete Moodbar--Datei realisiert. Dieser Ansatz funktioniert bei nicht
live gestreamter Musik gut. Bei live gestreamter Musik könnte für die
Visualisierung eine Fast--Fourier--Transformation in Echtzeit durchgeführt
werden. Da jedoch die Ressourcen des *Raspberry Pi* sehr begrenzt sind, sollte hier
auf die Verwendung einer GPU--beschleunigten--FFT zurückgegriffen werden
(vgl. [@Sabarinath2015], Seite 657 ff.).

Ein alternativer Ansatz wäre auch die Realisierung einer Musik--Visualisierung
mittels Hardwarekomponenten. Ein möglicher Ansatz aus Hardware--basierten
Hochpass-- und Tiefpassfiltern in Form einer Disco--Beleuchtung wird unter
[@2014projekte], Seite 261 ff. beschrieben.

### Echtzeituhr

Der *Raspberry Pi* besitzt keine Hardware--Uhr. Aufgrund der Tatsache, dass es
sich bei *Eulenfunk* um ein Internetradio handelt wurde auf eine Echtzeituhr
(real time clock, RTC) verzichtet, da sich die Uhr von *Eulenfunk* aufgrund der
permanenten Internetverbindung mittels NTP[^NTP] über das Internet
synchronisieren kann. Eine Erweiterung um eine Echtzeituhr wird in
[@horan2013practical], Seite 145 ff. und [@gay2014experimenting], Seite 77 ff. ausführlich beschrieben.
Mit einer RTC wäre die Implementierung einer Weckerfunktion möglich, welche
das Radio gezielt zu einer bestimmten Uhrzeit aufweckt und Musik abspielt.

### Fernbedienung

Eine weitere Erweiterung wäre die Integration einer Fernbedienung. Diese ließe
sich relativ einfach mittels eines Infrarot--Sensors und beispielsweise der
*lirc*--Bibliothek umsetzen. Siehe auch [@warner2013hacking], Seite 190 ff. für
weitere Informationen.

### Batteriebetrieb

Da die Strom-- beziehungsweise Spannungsversorgung beim *Raspberry Pi*
problematisch ist, wäre auch ein Batterie-- beziehungsweise Akkubetrieb möglich.
Eine einfache Schaltung für einen Batteriebetrieb würde sich beispielsweise mit
einem *LM7805*--Spannungsregler oder einem Abwärtswandler realisieren lassen
([vgl. @gay2014mastering], Seite 24 ff.). 

[^NTP]: Network Time Protocol:
\url{https://de.wikipedia.org/wiki/Network_Time_Protocol}

## Fazit 

Grundsätzlich kann der Hardware--Prototyp als erfolgreich umgesetzt betrachtet werden.
Die geplanten Anforderungen an die Hardware konnten soweit alle umgesetzt
werden. Jedoch sollte gesagt werden, dass das Netzteil und der USB--Hub --- wie
bereits im  jeweiligem Kapitel erläutert --- aktuell ein eher instabiles
Verhalten aufweisen. Hier ist noch eine Umstellung auf einen besseren USB--Hub
(und Netzteil mit 5V, 3A) ohne *backfeed* nötig.

[^PIHUB]: Spezieller USB--Hub für den Raspberry Pi: \url{https://www.raspberrypi.org/blog/pihub-a-robust-powered-hub-for-your-pi/}

Auf der Seite des Betriebssystems wäre die relativ junge Linux--Distribution
*Alpine Linux*[^APL] eine mögliche Verbesserung für den Einsatzzweck
Internetradio. Diese Distribution hat ihren Fokus auf Ressourceneffizienz und
Systemsicherheit. Ein weiterer Vorteil wäre der `diskless mode`, welcher das
komplette Betriebssystem in den Arbeitsspeicher lädt. In diesem Modus müssen
Änderungen mit einem *Alpine Local Backup (lbu)*--Tool explizit auf die
Festplatte geschrieben werden. Das hätte den Vorteil, dass man die Abnutzung des
Flash--Speichers, durch unnötige Schreib/Lese--Vorgänge, minimieren würde.
Momentan unterstützt diese Distribution allerdings noch nicht das von uns
favorisierte ``systemd``.

[^APL]: Alpine Linux für Raspberry Pi: \url{https://wiki.alpinelinux.org/wiki/Raspberry_Pi}
\label{internal-software-anforderungen}
