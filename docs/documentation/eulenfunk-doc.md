---
documentclass: scrreprt
classoption: toc=listof,index=totoc
include-headers:
    - \usepackage{url} 
    - \usepackage[ngerman]{babel}
    - \usepackage{csquotes}
    - \usepackage[babel, german=quotes]{csquotes}
fontsize: 11pt
sections: yes
toc: yes
lof: yes
lot: no
date: \today
---

\newpage
\pagenumbering{arabic} 
\setcounter{page}{1}

# Vorwort 

TODO: thanks: rf-electronic, Herr Schäferling

## Haftungsausschluss 

Das vorliegende Projekt ist im Rahmen einer Studienarbeit im Fach
Hardwaresysteme an der Hochschule Augsburg entstanden. Da die Autoren nicht aus
dem Bereich der *Technischen Informatik* sind, wurden jegliche
hardwareseitigen Arbeiten nach bestem Grundlagenwissen umgesetzt.

## Namensgebung

Der Name des Projektes ist \frqq\texttt{Eulenfunk}\flqq. Die Bezeichnung der
Eule wurde analog zum Tierreich gewählt, da die *Eule* hier als Vogel aufgrund
ihrer Erkennungsmerkmale von anderen Vögeln in der Regel als *Fressfeind*[^EULE]
klassifiziert wird. Analog dazu ist ein *Do--It--Yourself*--Internetradio ---
welches je nach Konfiguration, im Vergleich zu einem *Closed--Source*--Produkt,
günstiger und mit mehr Funktionalität ausgestattet werden kann ---
möglicherweise ein Dorn im Auge kommerzieller Internet--Radio--Anbieter.

## Zielsetzung 

Das grundlegende Projektziel ist aus vorhandenen alten Hardware--Komponenten ein
möglichst optisch und klanglich ansprechendes Internetradio zu entwickeln. Dabei
liegt der Schwerpunkt vor allem auch auf einem guten Kosten/Nutzen--Verhältnis.
Als Basis für das Projekt dient ein defektes Analog--Radio und ein *Raspberry
Pi* aus dem Jahr 2012.

Diese Studienarbeit soll einen Überblick über die verwendeten, beziehungsweise
benötigten Komponenten für den Bau eines *Raspberry Pi*--Internetradios
verschaffen. Anschließend soll das Wissen für die Ansteuerung bestimmter
Hardware--Komponenten mittels der *Raspberry Pi*--GPIO[^GPIO] Schnittstelle vermittelt
werden.

\begin{figure}[h!]
  \centering
\includegraphics[width=0.5\textwidth]{images/front_3.png}
  \caption{Aktueller Prototyp}
  \label{fertig}
\end{figure}

\newpage 

Abbildung \ref{fertig} zeigt den *Eulenfunk* Prototypen, welcher im Zeitraum von
drei Wochen im Rahmen des Hardwaresysteme Kür--Projekts entstanden ist. Auf
Vimeo[^VIMEO] ist ein Video des aktuellen Prototypen zu sehen.


## Verwendete Software

Für die Entwicklung und Dokumentation wurden folgende *GNU/Linux* Tools
verwendet:

* *Pandoc/LaTeX* (Dokumentation) 
* *Vim* (Softwareentwicklung) 
* *Fritzing* (Schaltpläne).


[^VIMEO]: Eulenfunk Prototyp: \url{https://vimeo.com/171646691}
[^EULE]: Lebensweise der Eule: \url{https://de.wikipedia.org/wiki/Eulen\#Lebensweise}
[^GPIO]: General-purpose input/output Schnittstelle: \url{https://en.wikipedia.org/wiki/General-purpose_input/output}


# Motivation

## Private Situation 

Die Autoren dieses Projekts leben in einer Wohngemeinschaft zusammen. Die Küche
ist der Ort an welchem gemeinsam gekocht und gespeist wird. Für eine angenehme
Atmosphäre und als Nachrichten--Quelle sorgte in der Küche früher ein
Analog--Radio der Firma *AEG*, welches aufgrund der schlechten Empfangsqualität
durch eine Kombination aus »alter Stereoanlage«, »altem Raspberry Pi« und
einem »alten Thinkpad X61t« ersetzt wurde. In dieser Kombination fungierte
die Stereoanlage als Soundausgabe--Komponente, auf dem *Raspberry Pi* lief der
Linux--basierte Player Volumio[^VOL], welcher mit dem Touchscreen des *Thinkpad
x61t* über eine Weboberfläche gesteuert wurde. Diese Kombination hat zwar
funktioniert, jedoch war sie alles andere als benutzerfreundlich, da zuerst die
Stereoanlage und der Laptop  eingeschaltet werden mussten und eine WLAN--Verbindung
zum *Raspberry Pi*--Player hergestellt werden musste. Diese Situation weckte den
Wunsch nach einer komfortableren Lösung, beispielsweise ein Internetradio auf
Basis des *Raspberry Pi*.

[^VOL]: Volumio: \url{https://volumio.org/}


## Kommerzielle Produkte

Kommerzielle Anbieter von Internetradios gibt es wie Sand am Meer. Die
Preisspanne liegt hier zwischen \EUR{30} und mehreren hundert Euro. Der
Funktionsumfang sowie die Wiedergabequalität ist hier von Hersteller zu Hersteller
und zwischen den verschiedenen Preisklassen sehr unterschiedlich. Einen
aktuellen Überblick aus dem Jahr 2016 über getestete Modelle gibt es
beispielsweise online unter *bestendrei.de*[^TEST].

Das *Problem* bei den kommerziellen Anbietern ist, dass man hier jeweils an die
vorgegebenen Funktionalitäten des Herstellers gebunden ist. Bei einem
Do--It--Yourself--Projekt auf Basis Freier Software beziehungsweise eines freien
Hardwaredesigns, hat man die Möglichkeit alle gewünschten Funktionalitäten ---
auch Features die von keinem kommerziellen Anbieter unterstützt werden --- zu
integrieren. Beispiele für Funktionalitäten, welche bei kommerziellen Produkten
nur schwer beziehungsweise vereinzelt zu finden sind:

* Unterstützung bestimmter WLAN--Authentifizierungsstandards
* Einhängen von benutzerdefinierten Dateifreigaben wie *Samba*, *NFS*, *SSHFS*
* Unterstützung verschiedener *lossy* und *lossless* Formate *OGG VORBIS*, *FLAC*, u.a.
* Integration verschiedener Dienste wie beispielsweise *Spotify* 
* Benutzerdefinierte Anzeigemöglichkeiten (Uhrzeit, Wetter, et cetera.)



[^TEST]:Test von Internetradios: \url{http://www.bestendrei.de/elektronik/internetradio/}



# Projektspezifikation

## Hardwareanforderungen

Das Radio soll dem Benutzer folgende Hardwarekonfigurationsmöglichkeiten bieten:

* Anschluss passive Lautsprecher/Kopfhörer möglich 
* Lautstärkeregelung über Hardware möglich
* Verwendung des internen Lautsprechers des alten Radios
* Statusinformationen zum aktuellen Lied beispielsweise über ein LCD
* LEDs als Statusanzeige und/oder als Visualisierungsvariante von Musik[^MOODBAR]
* USB--Anschlussmöglichkeit für externe Datenträger

[^MOODBAR]: Moodbar: \url{https://en.wikipedia.org/wiki/Moodbar}

## Softwareanforderungen

Die Software soll generisch gehalten werden um eine möglichst einfache
Erweiterbarkeit zu  gewährleisten. 

TODO Eule: Hier was zu Menü--Steuerung schrieben und Umfang?

## Optik-- und Usability--Anforderungen

Die Eingabe--Peripherie soll möglichst einfach gehalten werden, um eine *schöne*
Produkt--Optik zu gewährleisten, dabei sollen folgende Anforderungen erfüllt
werden:

* Minimale sowie ansprechende Bedienelemente
* Funktionales, zweckgebundenes *Design*
* *Retro--Look*-Aussehen wünschenswert

Das *Design* soll im Grunde *minimalistisch*  gehalten werden, das heißt, es
sollen aufgrund der Übersichtlichkeit nur so wenige »Bedienelemente« wie nötig
angebracht werden.


# Hardware

## Komponenten und Bauteile

Abbildung \ref{uebersicht} zeigt eine konzeptuelle Übersichts des Zusammenspiels
der einzelnen Komponenten.

\begin{figure}[h!]
  \centering
  \includegraphics[width=0.7\textwidth]{images/uebersicht.png}
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

\newpage 

**Mussten noch erworben werden:**

* Audioverstärker 
* Drehimpulsregler 
* Kunststoffabdeckung für Front
* Farbe (Lack)
* Drehknöpfe für das Gehäuse 

[^AEG]: AEG Küchenradio 4104: \url{https://www.amazon.de/AEG-MR-4104-Desgin-Uhrenradio-buche/dp/B000HD19W8}

## Raspberry Pi

Der vorhandene *Raspberry Pi* ist aus dem Jahr 2012. Die genaue CPU-- und
Board--Revision kann auf Linux unter ``proc`` ausgelesen werden, siehe auch
[@gay2014raspberry], Seite 46:

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

Laut Tabelle unter [@gay2014raspberry], Seite 45 handelt es sich hierbei
(Revision 0003) um das Modell B Revision 1+ mit 256MB RAM.

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
Modulationsverfahren angewandt, welches Pulsweitenmodulation (PWM) heißt. Hierbei
wird...(Ref auf Software TODO: ELCH?) Unter [@richardson2014make], Seite 121 ff. und
[@gay2014mastering], Seite 421 ff. finden sich weitere Informationen.

Software PWM unter [@gay2014experimenting], Seite 183 ff. zeigt beispielsweise
eine 6% CPU--Last pro GPIO--Pin bei einer PWM--Softwareimplementierung. TODO: ELCH?

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

Das Display arbeitet mit einer Logik--Spannung von 3.3V - 5V. Da die GPIO--Pins
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
und Potentiometer gehängt.

Der resultierende Gesamtwiderstand ohne Potentiometer beträgt in diesem Fall
$\approx$ 50 $\Omega$:

$$  R_{ges} = \frac{R_1 \times R_2}{R_1 + R_2} = \frac{100\Omega \times 100\Omega}{100\Omega + 100\Omega} = 50\Omega $$

## Drehimpulsgeber

Um eine minimale Anzahl an Bedienelementen zu erhalten, wird bei *Eulenfunk*
ein Drehimpulsgeber mit Schalter gewählt. Für erste Testzwecke wurde von der 
Hochschule ein *ALPS STEC12E08* bereitgestellt. Dieser wurde im Laufe der
Entwicklung durch einen *ALPS STEC11B09*[^ALPS] ersetzt, da dieser mittels Mutter und
Schraube am Gehäuse besser befestigt werden kann. 

Der verwendete Drehimpulsgeber hat insgesamt fünf Anschlüsse. Zwei
Signalleitungen (A und B), zwei mal *GND* (jeweils für Drehgeber und Schalter)
und einen Anschluss für den Schalter. Beim Drehen eines Drehimpulsgebers wird
ein Rechtecksignal generiert. Je nach Muster der beiden Datensignale A oder B,
kann entschieden werden ob es sich um eine Rechts-- oder Linksdrehung handelt.
Siehe [@2014projekte], Seite 361 ff. für weitere Hintergrundinformationen zu
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

Die interne Soundkarte des *Raspberry Pi* ist über eine triviale
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
abgeflachte Seite --- meist mit dem kürzerem Beinchen -- ist in der Regel die
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
\includegraphics[width=0.9\textwidth]{images/transistorled.png}
  \caption{Transistor--RGB--LED Schaltung}
  \label{transled}
\end{figure}

Damit der Transistor jedoch *sicher* durchschaltet, werden Wiederstände mit $10k
\Omega$ verwendet. Die in Abbildung \ref{transled} gelisteten
LED--Vorwiderstände ergeben sich aufgrund der verschiedenen Spannungen der
unterschiedlichen Farben[^RGBGP]. Die Berechnung für den Vorwiderstand pro LED
schaut am Beispiel der Farbe blau ($U_{LED} = 3,15V, I_{LED} = 20mA$) wie folgt
aus:

$$R_{LED} = \frac{U_{Betriebsspannung} - U_{LED}}{I_{LED}} = \frac{5V - 3,15V}{20mA} =92.5 \approx 100\Omega$$


## USB--Hub und Netzteil

Der *Raspberry Pi* hat in unserer Revision nur zwei USB--Schnittstellen, diese
sind bereits durch die Hardware--Komponenten USB--DAC (Soundkarte) und das
Wi--Fi--Modul belegt. Um den Anschluss eines externen Datenträgers, auch mit
größerer Last wie beispielsweise einer Festplatte zu ermöglichen, wird ein
aktiver USB--Hub benötigt.

Für diesen Einsatzzweck wird aus den Altbeständen ein *LogiLink 4 Port USB 2.0
HUB*[^HUB] verwendet. Viele billig-Hubs arbeiten hier entgegen der USB--Spezifikation
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

\newpage

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
lackiert[^LACK], hierbei wurden die Flächen für LCD und die drei LEDs abgelebt.
Zudem werden schwarze Knöpfe in Alu--Optik mit $\diameter$ 30mm  für den
Lautstärkeregler und den Drehimpulsgeber verwendet.

### Hinterseite

Für die Hinterseite wird die alte Abdeckung des AEG--Radios verwendet. Diese
musste teilweise leicht modifiziert werden. An dieser befinden sich zwei Potis
für Kontrastregelung und Hintergrundbeleuchtung des LCD, eine
USB--Female--Kabelpeitsche, zwei Cinch Stecker für externe Lautsprecher und ein
Kippschalter zum Umschalten zwischen internen und externen Lautsprechern.

[^LACK]: Buntlack, hellelfenbein: \url{http://www.obi.de/decom/product/OBI_Buntlack_Spray_Hellelfenbein_hochglaenzend_150_ml/3468725}


## Betriebssystem

Mittlerweile gibt es für den *Raspberry Pi* viele offiziell zugeschnittene
Betriebssysteme (vgl. [@pietraszak2014buch], Seite 29 ff., [@warner2013hacking],
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


### Abspielsoftware

Für den Betrieb des Internetradios soll der MPD (Music--Player--Daemon) verwendet
werden, da *Eulenfunk* auf einem eigens entwickeltem MPD--Client basieren soll
(mehr zur Eulenkfunk Software siehe Kapitel Software). Andere Projekte greifen
oft auf Abspielsoftware wie den *MOC* [vgl. @pietraszak2014buch], Seite 189 ff.
oder *mplayer* [@exploring] Seite 638 ff. zu. 


```bash
    # Installation des MPD
    [root@eulenfunk ~]$ pacman -Sy mpd mpc ncmpc
```

# Software

## Anforderungen

TODO: Mpd erklären!

- Leichte integrierbarkeit mit anderen Anwendungen (lose Kopplung)
TODO: Mehr.

## Softwarearchitektur

Die Nachbaubarkeit vieler Bastelprojekte ist häufig durch die Software recht
eingeschränkt, da diese entweder nicht frei verfügbar ist oder zu wenig
generisch ist als dass man die Software leicht auf das Projekt anpassen könnte.
Meist handelt es sich dabei um ein einziges, großes C--Programm oder ein eher
unübersichtliches Python--Skript (TODO: refs?). Aus diesem Grunde soll die
Software für *Eulenfunk* derart modular aufgebaut sein, dass man einzelne Module
problemlos auch auf andere Projekte übertragbar sind und später eine leichte 
Erweiterbarkeit gewährleistet ist. Damit auch andere die Software einsetzen
können wird sie unter die GPL in der Version 3 (TODO: ref) gestellt.

Zu diesem Zwecke ist die Software in zwei Hauptschichten unterteilt.
Die untere Schicht bilden dabei die *Treiber*, welche die tatsächliche
Ansteuerung der Hardware erledigt. Dabei gibt es für jeden Teil der Hardware
einen eigene Treiber, im Falle von *Eulenfunk* also ein separates Programm für
die LCD-Ansteuerung, das Setzen der LED Farbe und dem Auslesen des oberen Rotary
Switches.

Die Schicht darüber bilden einzelne Dienste (ingesamt fünf), die über eine
Netzwerkschnittstelle angesprochen werden und jeweils eine Funktionalität des
Radios umsetzen. So gibt es beispielsweise einen Dienst der die Ansteuerung des
LCD--Displays *komfortabel* macht, ein Dienst, der die LEDs passend zur Musik
schaltet und ein Dienst der automatisch eine Playlist aus der Musik auf
angesteckten externen Speichermedien erstellt. Die jeweiligen Dienste sprechen
mit den Treibern indem sie Daten auf ``stdin`` schreiben, bzw. Daten von
``stdout`` lesen. Um die Dienste auf neue Projekte zu portieren ist also nur
eine Anpassung der Treiber notwendig.

Der Vorteil liegt dabei klar auf der Hand: Die lose Kopplung der einzelnen
Dienste erleichtert die Fehlersuche ungemein und macht eine leichte
Austauschbarkeit und Übertragbarkeit der Dienste in anderen Projekte möglich.
Stellt man beispielsweise fest, dass der Prozessor des Radios voll ausgelastet
ist, so kann man mit Tools wie ``htop`` sehr einfach herausfinden welcher Dienst dafür
verantwortlich ist.

## Überblick der einzelnen Komponenten

Ein Überblick über die existierenden Dienste liefert Abbildung
\ref{eulenfunk-services}.

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.0\textwidth]{images/eulenfunk-services.png}
  \caption{Übersicht über die Softwarelandschaft von Eulenfunk. Dienste mit
  einer Netzwerkschnittstelle sind durch den entsprechenden Port gekennzeichnet.}
  \label{eulenfunk-services}
\end{figure}

### Sprachwahl

Die momentane Software ist in den Programmiersprachen *C* und *Go* geschrieben.
Dazu kommt lediglich ein Bash--Skript zum Auslesen von Systeminformationen.

Die Ressourcen auf dem Raspberry Pi sind natürlich sehr limitiert, weswegen sehr
speicherhungrige Sprachen wie Java oder Ähnliches von vornherein ausschieden.
Obwohl Python nach Meinung des Autors eine schöne Sprache ist und viele gute
Bibliotheken für den Pi bietet, schied es ebenfalls aus diesem Grund aus.

Ursprünglich war sogar geplant, alles in *Go* zu schreiben. Leider gibt es nur
wenige Pakete für die GPIO Ansteuerung und auch keine Bibliothek für
Software--PWM. Zwar hätte man diese notfalls auch selbst mittels
``/sys/class/gpio/*`` implementieren können, doch bietet *Go* leider keine
native Möglichkeit mit Interrupts zu arbeiten. Wie später beschrieben ist dies
allerdings für den Treiber nötig, der den Rotary Switch ausliest.

Für *Go* sprechen wir ansonsten folgende Gründe:

- **Garbage collector:** Erleichtert die Entwicklung lang laufender Dienste.
- **Hohe Grundperformanz:** Zwar erreicht diese nicht die Performanz von C, 
  liegt aber zumindest in der selben Größenordnung (TODO: ref)
- **Weitläufige Standardlibrary:** Kaum externe Bibliotheken notwendig.
- **Schneller Kompiliervorgang:** Selbst große Anwendungen werden in wenigen 
  Sekunden in eine statische Binärdatei ohne Abhängigkeiten übersetzt.
- **Kross--Kompilierung:** Durch Setzen der ``GOARCH`` Umgebungsvariable kann 
  problemlos auf einen x86-64--Entwicklungsrechner eine passende
  ARM--Binärdatei erzeugt werden.
- **Eingebauter Scheduler:** Parallele und nebenläufige Anwendungen wie
  Netzwerkserver sind sehr einfach zu entwickeln.
- Ein Kriterium war natürlich auch dass der Autor gute Erfahrung mit der Sprache
  hatte und **neugierig** war, ob sie auch für solche Bastelprojekte gut einsetzbar
  ist.

``C`` ist hingegen für die Entwicklung der Treiber vor allem aus diesen Gründen
eine gute Wahl:**

- Programmierung mit **Interrupts** bequem und nativ möglich.
- Höchste Performanz und geringer Speicherverbrauch.
- Verfügbarkeit von wiringPi.

### Vorhandene Softwarelibraries

Es folgt eine Liste der benutzten Bibliotheken:

#### ``wiringPi`` (http://wiringpi.com)

``wiringPi`` ist eine Portierung der Arduino--``Wiring`` Bibliothek von Gordon
Henderson auf den Raspberry Pi. Sie dient wie ihr Arduino--Pendant zur leichten
Steuerung der verfügbaren Hardware, insbesondere der GPIO--Pins über
``/dev/mem``. Daneben wird für den LCD--Treiber auch die mitgelieferte
LCD--Bibliothek genutzt. Für den LED--Treiber wird zudem die softwarebasierte
Pulsweitenmodulation genutzt, allerdings in einer leicht veränderten Form.

#### ``go-mpd`` (https://github.com/fhs/gompd)

Eine einfache MPD--Bibliothek, die wenig mehr als die meisten Kommandos des
MPD--Protokolls unterstützt. 

TODO: ref (https://www.musicpd.org/doc/protocol/)

#### ``go-colorful`` (github.com/lucasb-eyer/go-colorful)

Eine Bibliothek um Farben in verschiedene Farbräume zu konvertieren. Der Dienst
der die LED passend zur Musik setzt nutzt diese Bibliothek um RGB--Farbwerte in
den HCL--Farbraum zu übersetzen. Dieser eignet sich besser um saubere Übergänge
zwischen zwei Farben zu berechnen und Farbanpassungen vorzunehmen.
Die genaue Funktionsweise wird weiter unten beleuchtet (TODO: ref).

#### cli (github.com/urfave/cli)

Eine komfortable und reichhaltige Bibliothek um Kommandozeilenargumente zu parsen.
Unterstützt Subkommandos ähnlich wie ``git``, welche dann wiederum eigene
Optionen oder weitere Subkommandos besitzen können. Beide Features wurden
extensiv eingesetzt, um alle in *Go* geschriebenen Dienste in einer Binärdatei
mit konsistentem Kommandozeileninterface zu vereinen.

## Treiber--Software

In Summe gibt es momentan drei unterschiedliche Treiber. Sie finden sich im
``driver/`` Unterverzeichnis[^driver_github] der Software nebst einem passenden
Makefile. Nach dem Kompilieren entstehen drei Binärdateien, welche mit dem
Präfix ``radio-`` beginnen:

- ``radio-led:`` Setzt die Farbe des LED--Panels auf verschiedene Weise.
- ``radio-lcd:`` Liest Befehle von ``stdin`` und setzt das Display entsprechend.
- ``radio-rotary:`` Gibt Änderungen des Rotary Switches auf ``stdout`` aus.

Die genaue Funktionsweise dieser drei Programme wird im Folgenden näher beleuchtet.

[^driver_github]: Siehe auf GitHub: \url{https://github.com/studentkittens/eulenfunk/tree/master/driver}

### LED--Treiber (``driver/led-driver.c``)

Der LED--Treiber dient zum Setzen eines RGB--Farbwerts. Jeder Kanal hat den
Wertebereich 0 bis 255. Die Hilfe des Programms zeigt die verschiedenen
Aufrufmöglichkeiten:

```bash
usage:
  radio-led on  ....... turn on LED (white)
  radio-led off ....... turn off LED
  radio-led cat ....... read rgb tuples from stdin
  radio-led rgb  r g b  Set LED color to r,g,b
  radio-led hex #RRGGBB Set LED color from hexstring
  radio-led fade ...... Show a fade for debugging
```

Erklärung benötigt hierbei nur der ``cat``--Modus, bei dem der Treiber
zeilenweise RGB--Farbtripel ``stdin`` liest und setzt. Dieser Modus wird benutzt,
um kontinuierlich Farben zu setzen ohne ständig das Treiberprogramm neu zu starten.

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.0\textwidth]{images/eulenfunk-pwm.png}
  \caption{Grafische Darstellung der Pulsweitenmodulation}
  \label{eulenfunk-pwm}
\end{figure}

Da ein GPIO--Pin prinzipiell nur ein oder ausgeschaltet werden kann verwenden
wir Pulsweitenmodulation. Dabei macht man sich die träge Natur des menschlichen
Auges zu nutze indem man die LED sehr schnell hintereinander ein und
ausschaltet. Ist die LED dabei pro Ein- und Ausschaltvorgang genauso lang hell
wie dunkel so leuchtet die LED mit etwa der Hälfte ihrer maximalen
Leuchtstärke. Durch Verlängerung/Verkürzung des eingeschalteten Zustands können
so viele verschiedene Helligkeitsstufen abgebildet werden. Siehe dazu auch 
Abbildung \ref{eulenfunk-pwm}.

TODO: Diagramm/Text anpassen.

Da bei niedrigen Helligkeitswerten der ausgeschaltete Zustand besonders lange
gehalten wird, kann es dazu kommen dass ein flackerhafter Eindruck entsteht, da
man die ausgeschalteten Phasen als solche wahrnehmen kann. Um dies zu
verhindern muss eine ausreichende hohe Frequenz gewählt werden. 

Anders als ursprünglich angenommen, mussten wir feststellen dass die GPIO Pins
des Raspberry Pi (mit Ausnahme von Pin 18 (TODO: ref)) kein hardwareseitiges
PWM unterstützen. Aus diesem Grund mussten wir auf softwareseitiges PWM
zurückgreifen, um Farben mit mindestens 256 Abstufungen zu erhalten. Nach etwas
Ausprobieren befanden wir die ``softPwm``--Bibliothek von ``wiringPi`` für
tauglich (TODO: ref).

Diese hat allerdings das Problem, dass eine hartkodierte Pulsweite von 100µs
verwendet wird. Für die meisten Anwendungsfälle und den vom Autor empfohlenen
100 Abstufungen ist das auch in Ordnung. Hundert unterschiedliche Zustände
waren nach kurzem Ausprobieren bei einem weichen Farbübergang zu stark
abgestuft.

$T_{Periode} = 100\mu s\times 100= 10000\mu s = 0.01s$

$f = \frac{1}{T_{Periode}} = 100Hz$

Optimal wären hier 256 unterschiedliche Zustände, um die volle 8-Bit Farbtiefe auszunutzen.
Daher mussten wir die entsprechende C--Datei kopieren (GPL3--lizensiert) und manuell anpassen.
Dabei haben wir die Pulsweite auf 50µs herabgesetzt, was bei einer Spanne von 256 Werten
eine Frequenz von optisch akzeptablen 78Hz ergibt:

$T_{Periode} = 50\mu s\times 256 = 12800\mu s = 0.0128s$

$f = \frac{1}{T_{Periode}} = 78.125Hz$

Diese Frequenz scheint optisch ausreichend flackerfrei zu sein und scheint die
CPU nicht übermäßig stark zu beeinflussen (rund 2-3% pro Kanal).

Es besteht eine Verbindung zu einen früheren Bastelprojekt namens
``catlight``[^catlight] --- einer mehrfarbigen, in einem Gehäuse montierten LED
die über USB angesprochen werden kann. Genutzt wird diese zur Benachrichtigung
bei neuen E--Mails, Chat--Nachrichten und ähnlichem. Zu diesem Zwecke wurde
auch bereits damals ein Treiberprogramm entwickelt, welches das selbe
Bedienungskonzept wie ``radio-led`` hat. Dies war während der Entwicklung von
*Eulenfunk* nützlich, da es die Entwicklung der Dienste ``ambilight`` und
``lightd`` unabhängig von der Radio--Hardware macht.

[^catlight]: \url{https://github.com/studentkittens/catlight}


### LCD--Treiber (``driver/lcd-driver.c``)

```bash
usage:
  radio-lcd [print-charset [start [end]]]
```

Der LCD--Treiber setzt Bereiche des LCD--Displays auf einen gegebenen Text.
Beim Start leert er das Display und liest ähnlich wie ``radio-led cat``
zeilenweise von ``stdin`` und entnimmt diesen Zeilen die Information welchen
Bereich des Displays gesetzt werden soll. Das vom Treiber erwartete
Zeilenformat ist dabei ``LINENO[,OFFSET] TEXT...``, wobei ``LINENO`` die
gewünschte Zeilennummer als Dezimalzahl ist und der optionale ``OFFSET`` der
Index an dem geschrieben werden soll. Dahinter folgt durch ein Leerzeichen
getrennt beliebiger Text. Ist kein ``OFFSET`` gegeben, so wird die ganze Zeile
überschrieben und nötigenfalls mit Leerzeichen aufgefüllt. Ist der Text länger
als die Zeile wird der Text abgeschnitten.

Der Treiber hält eine Matrix mit den aktuell gesetzten Zeichen und kann daher ein erneutes
Zeichnen einer Zelle im Display verhindern, indem es das neue Zeichen mit dem Alten vergleicht.
Unnötige Zeichenvorgänge waren als störende Schlieren auf dem Display wahrnehmbar.


Zudem bietet der Treiber mit dem ``print-charset`` Argument die Möglichkeit die
auf dem Display verfügbaren Zeichen aufs selbige auszugeben. Dazu stellt er
jeweils 80 Zeichen da und wartetet einige Sekunden bevor die nächsten 80
ausgegeben werden. Hat er alle 256 Zeichen ausgegeben beendet er sich. Optional
kann man auch ein Start- und End--Offset mitgeben, an dem er das Zeichnen
anfangen soll.

Der Treiber unterstützt eine Reihe hardkodierter Spezialzeichen, welche in der
Menüführung und der UI zu benutzt werden. Das LCD--Display unterstützt dabei 8
verschiedene *Custom Chars*, welche mittels der Codepoints 0-7 und 8-15
(wiederholt) setzbar sind. Momentan sind diese auf folgende Glyphen gesetzt:

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.0\textwidth]{images/symbols.png}
  \caption{}
  \label{eulenfunk-symbols}
\end{figure}

Die eigentliche Ansteuerung der Pins übernimmt dabei wieder die ``wiringPi``--Bibliothek,
beziehungsweise dessen LCD--Unterbibliothek (TODO: ref: https://projects.drogon.net/raspberry-pi/wiringpi/lcd-library/). Diese ist kompatibel mit dem Hitachi HD44780U und Nachbauten.

TODO:

- GPIO Anschluss?

### Drehimpulsgeber--Treiber

- Prellung
- ISR / kein printf
- rotary encoder
- http://theatticlight.net/posts/Reading-a-Rotary-Encoder-from-a-Raspberry-Pi/
- grey code (buch zitieren)
- button entprellen
- "Polling" Mainloop
- "Div by 3 hack"

## Service Software

### ``displayd`` -- Der Displayserver

#### Einleitung

Der Displayserver ``displayd`` kümmert sich um die Verwaltung der
Display--Inhalte. Er bietet eine höhere Abstraktionsschicht als der
vergleichsweise simple LCD--Treiber. Dabei bietet er die Abstraktion von
*Zeilen*, *Fenstern* und erleichtert dem Programmierer Enkodierungsaufgaben
indem es ein Subset von Unicode unterstützt. Eine Zeile ist dabei eine beliebig
langer utf8--enkodiertert Text ohne Zeilenumbruch. Die Zeile kann dabei länger
als das Display sein. In diesem Fall wird die Zeile abgeschnitten oder scrollt
je nach Konfiguration mit einer bestimmten Geschwindigkeit durch. Ein Fenster
hingegen ist eine benannte Ansammlung von Zeilen. Auch ein Fenster kann mehr
Zeilen haben als das Display. Vom Nutzer kann der Fensterinhalt dann vertikal
verschoben werden. Es können mehrere Fenster verwaltet werden, aktiv ist dabei
aber nur ein ausgewähltes.

Die Idee diese Funktionalität in einem eigenen Daemon auszulagern, ist vom
Grafikstack in unixoiden Betriebssystemen inspiriert. Dabei kümmert sich
ebenfalls ein Displayserver um die Verwaltung der Inhalte (meist ``X.org`` oder
``Wayland``) indem er auf bestimmte Treiber zurückgreift. Der Nutzer kann dann
mittels eines festgelegten Protokolls mit dem Displayserver sprechen und so
unabhängig von der verwendeten Rendering--Methode Inhalte darstellen. Da das
Protokoll zwischen Client und Displayserver meist trotzdem zu komplex ist,
haben sich für diese Aufgaben simple UI--Bibliotheken wie Xlib oder komplexere
GTK+ und Qt etabliert. Auch die Fenster--Metapher wurde dabei von den
Fenstermanagern übernommen.

Das Protokoll von ``displayd`` ist ein relativ simpel gehaltenes,
zeilenbasiertes Textprotokoll. Für den Zugriff auf dasselbige ist daher wird
auch keine UI--Bibliothek benötigt, lediglich einige Netzwerk und
Formatierungs--Hilfsfunktionen wurden implementiert (TODO: ref:
display/client.go). Basierend auf diesen Primitiven wurden aber auf Clientseite
Funktionalitäten wie Menü--»Widgets« implementiert, welche die grafische Darstellung
mit der Nutzereingabe verquicken.

Neben diesen Aufgaben löst ``displayd`` ein architektonisches Problem:
Wenn mehrere Anwendung versuchen auf das Display zu schreiben käme ohne zentrale
Instanz ein eher unleserliches Resultat dabei heraus. Durch ``displayd`` können
Anwendungen auf ein separates Fenster schreiben, wovon jeweils nur eines aktiv
angezeigt wird.

#### Architektur

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.0\textwidth]{images/eulenfunk-displayd.png}
  \caption{}
  \label{eulenfunk-displayd}
\end{figure}

Abbildung \ref{eulenfunk-displayd} zeigt die Architektur in der Übersicht.

Nach dem Start kann man auf Port 7777 ``displayd`` mittels eines simplen,
zeilenbasierten Textprotokolls kontrollieren. Dabei werden die folgenden
Kommandos (mit Leerzeichen--getrennten Argumenten) unterstützt:

```
switch <win>               -- Wechsle zum Fenster namens <win>.
line <win> <pos> <text>... -- Setze die Zeile <pos> im Fenster <win> zu <text>
scroll <win> <pos> <delay> -- Lässt Zeile <pos> in Fenster <win> mit der Geschwindigkeit <delay> scrollen.
move <win> <off>           -- Verschiebe Fenster <win> um <off> Zeilen nach unten.
truncate <win> <max>       -- Schneide Fenster <win> nach <max> Zeilen ab.
render                     -- Gebe aktuelles Fenster auf Verbindung aus.
close                      -- Schließt die aktuelle Verbindung.
quit                       -- Beendet Daemon und schließt Verbindung.

Dabei kann in die <Platzhalter> folgendes eingesetzt werden:

<win>:   Ein valider Fenstername (andernfalls wird ein neues Fenster mit diesen Namen angelegt)
<pos>:   Eine Zeilennummer, beginnend bei 0 für die erste Zeile.
<off>:   Ein Offset; 0 steht für keine Änderung; Kann negativ sein.
<max>:   Maximale Zeilenanzahl nach der das Fenster abgeschnitten wird.
<delay>: Zeitlicher Abstand zwischen zwei Scroll--Vorgängen.
         Siehe auch: https://golang.org/pkg/time/#ParseDuration
		 Beispiel: 100ms
```

Für die tatsächliche Anzeige nutzt ``displayd`` wie oben erwähnt das
Treiberprogramm ``radio-lcd``. Dabei wird in periodischen Abständen (momentan
150ms) das aktuelle Fenster auf den Treiber geschrieben. Zukünftige Versionen
sollen dabei intelligenter sein und nur die aktuell geänderten Zeilen
herausschreiben. Allerdings hat sich herausgestellt, dass man mit mehreren
scrollenden Zeilen bereits mit diesen Ereignisbasierten Ansatz auf eine höhere
Aktualisierungsrate kommt als mit den statischen 150ms. Eine Art »VSync«,
welches die Aktualisierungsrate intelligent limitiert wäre hier in Zukunft
wünschenswert.

#### Entwicklung

Da der Raspberry Pi nur bedingt als Entwicklungsplattform tauglich ist
(langsamer Compile/Run Zyklus), unterstützt ``displayd`` auch
Debugging--Möglichkeiten. Im Folgenden werden einige Möglichkeiten gezeigt 
mit ``displayd`` zu interagieren, beziehungsweise Programme zu untersuchen
die ``displayd`` benutzen:

```bash
# Den display server starten; --no-encoding schaltet spezielles LCD encoding
# ab welches auf normalen Terminals zu Artifakten # führt. 
$ eulenfunk display server --no-encoding &

# Gebe in kurzen Abständen das  "mpd" Fenster aus.
# (In separaten Terminal eingeben!)
$ eulenfunk display --dump --update --window mpd

# Verbinde zu MPD und stelle aktuellen Status auf "mpd" Fenster dar.
$ eulenfunk mpdinfo
# Auch nach Unterbrechung wird der zuletzt gesetzte Text weiterhin angezeigt:
$ <CTRL-C>
# Änderungen sind auch möglich indem man direkt mit dem Daemon über telnet
# oder netcat spricht. Hier wird die erste Zeile überschrieben, das aktuelle
# Fenster angezeigt und dann die Verbindung geschlossen.
$ telnet localhost 7777
line mpd 0 Erste Zeile geändert!
render                        
close
```

#### Enkodierung

Das LCD--Display unterstützt 8 bit pro Zeichen. Dabei sind die ersten 127 Zeichen weitesgehend
deckungsgleich mit dem ASCII Standard. Lediglich die Zeichen 0 bis 31 sind durch *Custom Chars*
und einige zusätzliche Zeichen belegt. Dies ist insofern auch sinnvoll, da in diesem Bereich 
bei ASCII Steuerzeichen definiert sind, die auf dem LCD schlicht keinen Effekt hätten.

Die Zeichen 128 bis 255 sind vom Hersteller des Displays mit verschiedenen
Symbolen belegt worden, die keinem dem Autor bekannten Encoding entsprechen. Da
auch nach längerer Internetrecherche keine passende Encoding--Tabelle gefunden
werden konnte, wurde (in mühevoller Handarbeit) eine Tabelle erstellt, die
passende Unicode--Glyphen auf das jeweilige Zeichen des Displays abbildet.
Nicht erkannte utf8--Zeichen werden als ein »?« gerendert anstatt Zeichen die
mehrere Bytes zur Enkodierung (wie » $\mu$  «) als mehrere falsche Glyphen
darzustellen. So wird beispielsweise aus dem scharfen »ß« das Zeichen 223.
Diese Konvertierung wird transparent von ``displayd`` vorgenommen wodurch es
möglich wird auch Musiktitel und ähnliches (annäherend) korrekt darzustellen.

Abbildung \ref{eulenfunk-encoding} zeigt das erstellte Mapping zwischen Unicode und LCD--Display.
Folgende Seiten waren bei der Erstellung der Tabelle hilfreich:

* http://www.amp-what.com (Suche mittels Keyword)
* http://shapecatcher.com (Suche mittels Skizze)

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.0\textwidth]{images/encoding.png}
  \caption{Unicode Version der LCD--Glyphen}
  \label{eulenfunk-encoding}
\end{figure}


### ``lightd`` -- Der Effektserver

``lightd`` ist ein relativ einfacher Service, dessen Hauptaufgabe die
Verwaltung auf den Zugriff auf die LED ist. Wollen mehrere Programme die LED
ansteuern, um beispielsweise einen sanften roten und grünen Fade--Effekt zu
realisieren so würde ohne Synchronisation zwangsläufig ein zittrige Mischung
beider Effekte entstehen.

Ursprünglich war ``light`` als ``lockd`` konzipiert, der den Zugriff auf
verschiedene Ressourcen verwalten kann. Da aber das Display bereits von
``displayd`` synchronisiert wird und durchaus mehrere Programme den Rotary
Switch auslesen dürfen wurde diese etwa generellere Idee wieder verworfen.

Die zweite Hauptaufgabe von ``lightd`` ist die Darstellung von Farbeffekten auf
der LED. Ursprünglich waren diese dafür gedacht um beispielsweise beim Verlust
der WLAN--Verbindung ein rotes Blinken anzuzeigen. Momentan wird allerdings nur
beim Herunterfahrend bzw. Rebooten ein rotes bzw. oranges Blinken angezeigt.
Folgende Effekte sind also momentan als Möglichkeit zur Erweiterung zu
begreifen:

- ``blend:`` Überblendung zwischen zwei Farben.
- ``flash:`` Kurzes Aufblinken einer bestimmten Farbe.
- ``fade:`` Lineare Interpolation von schwarz zu einer bestimmten Farbe und zurück.
- ``fire:`` Kaminfeuerartiges Leuchten.

Wie andere Dienste wird auch ``lightd`` mittels einer Netzwerkschnittstelle kontrolliert.
Die mögliche Kommandoes sind dabei wie folgt:

TODO: translate

```
// !lock      -- Try to acquire lock or block until available.
// !unlock    -- Give back lock.
// !close     -- Close the connection.
// <effect>   -- Lines starting without ! are parsed as effect spec.
//
// <effect> can be one of the following:
//
//   {<r>,<g>,<b>}
//   blend{<src-color>|<dst-color>|<duration>}
//   flash{<duration>|<color>|<repeat>}
//   fire{<duration>|<color>|<repeat>}
//   fade{<duration>|<color>|<repeat>}
//
// where <*-color> can be:
//
//   {<r>,<g>,<b>}
//
// and where <duration> is something time.ParseDuration() understands.
// <repeat> is a simple integer.
//
// Examples:
//
//   {255,0,255}                     -- The world needs more solid pink.
//   fire{1ms|{255,255,255}|0}       -- Warm fire effect.
//   blend{{255,0,0}|{0,255,0}|2s}   -- Blend from red to green.
```

TODO: Treiber.

### ``ambilightd`` -- Optische Musikuntermalung

Ein »Gimmick« von *Eulenfunk* ist es, dass die LED entsprechend zur momentan
spielenden Musik eingefärbt wird. Hier erklärt sich auch der Name dieses Dienstes:
*Ambilight* (TODO: ref) bezeichnet eigentlich eine von Phillips entwickelte Technologie,
um an einem Fernseher angebrachte LEDs passend zum momentanen Bildinhalt einzufärben.
Hierher kommt auch die ursprüngliche Idee dies auf Musik umzumünzen.

Um aus den momentan spielenden Audiosamples eine Farbe abzuleiten gibt es
einige Möglichkeiten. (TODO: hier HW schaltung etc. von oben aufgreifen). Eine
große Einschränkung bildet hierbei allerdings die sehr begrenzte Rechenleistung
des Raspberry Pi. Daher haben wir für eine Variante entschieden, bei der die
Farbwerte vorberechnet werden. Das hast den offensichtlichen Nachteil, dass man
für Radiostreams kein Ambientlicht anzeigen kann. Andererseits möchte man das
bei Nachrichtensendung und Diskussionsrunden vermutlich auch nicht.

Zur Vorberechnung nutzen wir dabei das ``moodbar`` Programm (TODO:
http://cratoo.de/amarok/ismir-crc.pdf). Diese analysiert mit Hilfe des
GStreamer--Frameworks (TODO: ref) eine Audiodatei in einem gebräuchlichen
Format und zerlegt diese in 1000 Einzelteile. Kurz erklärt[^NOTE] wird für
jedes dieser Teile ein Farbwert berechnet, wobei niedrige Frequenzen zu roten
Farbtönen werden, mittlere Frequenzen zu Grüntönen und hohe Frequenzen zu
blauen Tönen werden. Die so gesammelten Farbwerte werden dann in einer
``.mood`` Datei gespeichert, welche aus 1000 RGB--Tripeln à 3 Byte (1 Byte pro
Farbkanal) bestehen. Ein visualisiertes Beispiel für eine Moodbar kann man in Abbildung
\ref{queen-moodbar} sehen.

[^NOTE]: Der eigentliche Algorithmus ist komplexer und wird im referenzierten Paper beschrieben.

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.0\textwidth]{images/we-will-rock-you-mood.pdf}
  \caption{Moodbar--Visualisierung des Liedes »We will rock you« von »Queen«.
  Das Fußstapfen und Klatschen am Anfang ist gut erkennbar.
  Die Visualisierung wurde mit einem Python--Skript des Autor erstellt 
  (\url{https://github.com/sahib/libmunin/blob/master/munin/scripts/moodbar\_visualizer.py})}
  \label{queen-moodbar}
\end{figure}

Um nun aber tatsächlich ein zur Musik passendes Licht anzuzeigen muss für jedes
Lied in der Musikdatenbank eine passende Moodbar in einer Datenbank
abgespeichert werden. Diese Datenbank ist im Fall von *Eulenfunk* ein simples
Verzeichnis in dem für jedes Lied die entsprechende Moodbar mit dem Pfad
relativ zum Musikverzeichnis[^NOTE2] als Dateinamen abgespeichert wird.

[^NOTE2]: Wobei »/« durch »|« im Dateinamen ersetzt werden.

Die Datenbank kann dabei mit folgenden Befehl angelegt und aktuell gehalten werden:

```bash
$ eulenfunk ambilight --update-mood-db --music-dir /music --mood-dir /var/mood 
```

Die eigentliche Aufgabe von ``ambilightd`` ist es nun den Status von MPD zu
holen, die passende ``.mood``--Datei zu laden und anhand der Liedlänge und der
bereits vergangenen Zeit den aktuellen passenden Sample aus den 1000
vorhandenen anzuzeigen. Damit der Übergang zwischen den Samples flüssig ist
wird linear zwischen den einzelnen Farbwerten überblendet.

TODO: Farbkorrektur, referenz?
TODO: radio-led driver in Grafik?

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.0\textwidth]{images/eulenfunk-ambilight.png}
  \caption{TODO}
  \label{eulenfunk-ambilight}
\end{figure}

Im Fazit kann man sagen, dass die verwendete Technik durchaus gut für die meisten Lieder funktioniert.
Besonders die Synchronisation ist dabei erstaunlich akkurat und solange nur ein einzelnes Instrument spielt,
wird dem auch eine passende Farbe zugeordnet. Ein Dudelsack erscheint beispielsweise meist grün, während 
ein Kontrabass in mehreren Liedern Rot erschien.

Lediglich bei schnellen Tempowechseln (Beispiel: »Prison Song« von »System of a Down«)
sieht man, dass der Farbübergang bereits anfängt bevor man den zugehörigen Ton
hört. Dem könnte im Zukunft abgeholfen werden indem keine lineare Interpolation
zwischen den Farben mehr genutzt wird sondern beispielsweise quadratische. (TODO: ja?)

### ``automount`` -- Playlists von USB Sticks erstellen

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.0\textwidth]{images/eulenfunk-automount.png}
  \caption{Beispielhafte Situation bei Anschluß eines USB Sticks mit dem Label »usb-music«}
  \label{eulenfunk-automount}
\end{figure}

Der ``autmount``--Daemon  sorgt dafür, dass angesteckte USB--Sticks automatisch auf gemounted werden,
Musikdateien darauf indiziert werden und in einer Playlist mit dem »Label« des Sticks landen.
Die Architektur von ``automount`` ist konzeptuell in Abbildung
\ref{eulenfunk-automount} dargestellt. Unter Linux kümmert sich ``udev``, um
die Verwaltung  des ``/dev``--Verzeichnis. Wird ein neues Gerät angesteckt, so 

TODO

  udev regel erklären

Eine berechtigte Frage ist warum ``automount`` das Mounten/Unmounten des Sticks
übernimmt, wenn diese Aktionen prinzipiell auch direkt von der ``udev``--Regel
getriggert werden können. Der Grund für diese Entscheidung liegt am
*Namespace*--Feature von ``systemd`` (TODO: ref):  Dabei können einzelne
Prozesse in einer Sandbox laufen, der nur begrenzten Zugriff auf seine Umgebung
hat. Ruft man ``/bin/mount`` direkt aus der Regel heraus auf, so wird der Mount
im *Namespace* von ``udev`` erstellt und taucht daher nicht im normalen
Dateisystem auf.  Sendet man hingegen einen Befehl an ``automount``, so agiert dieser
außerhalb einer Sandbox und kann den Stick ganz normal als ``root``--Nutzer mounten.

- unmount/quit
- Protokoll

### mpdinfo

TODO: mit ui mergen?

### ui

- Windows (mehr dazu im Designteil)


## Einrichtung

### mpd und ympd

TODO: eigentliche mpd einrichrung.

``mpd`` und ``ympd`` sind die einzigen Dienste die von außen (ohne
Authentifizierung) zugreifbar sind. Auch wenn *Eulenfunk* normal in einem
abgesicherten WLAN hängt, wurde für die beiden Dienste jeweils ein eigener
Nutzer mit eingeschränkten Rechten und ``/bin/false`` als Login--Shell angelegt.

### systemd

In Abbildung \ref{eulenfunk-systemd} ist schematisch der
Abhängigkeitsgraph der einzelnen Dienste gezeigt.

\begin{figure}[h!]
  \centering
  \includegraphics[width=0.9\textwidth]{images/eulenfunk-systemd.png}
  \caption{Abhängigkeits Graph beim Start der einzelnen Dienste}
  \label{eulenfunk-systemd}
\end{figure}

TODO: restart / restart-on-failure

### udev

### sonstiges

Makefile

pstree

zeroconf

memory usage, cpu usage

godoc: 

systemd boot plot

https://godoc.org/github.com/studentkittens/eulenfunk/display

## Wartung

ssh zugang

systemctl status/journalctl zum logging

## Fazit

cloc statistiken

probleme: re-mount, mpd brokenness due to powerhub.

# Zusammenfassung

## Ziel erreicht?

Das selbstgesetzte Ziel --- mit möglichst wenig Aufwand ein Internetradio auf Basis
eines *Raspberry Pi* zu entwickeln --- kann durchaus als erfolgreich betrachtet
werden. 

## Erweiterungen und alternative Ansätze

### Allgemein

Der aktuelle Prototyp hat lediglich nur ein Potentiometer um die Hintergrundbeleuchtung
des LCD zu regeln. Ein anderer Ansatz wäre der Einsatz eines Relais, welches es
ermöglichen würde die LCD--Hintergrundbeleuchtung Software--seitig ein-- und auszuschalten.

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

## Mögliche Verbesserungen

### Alpine Linux 

Die relativ junge Linux--Distribution *Alpine Linux*[^APL] wäre eine mögliche
Verbesserung für den Einsatzzweck Internetradio. Diese Distribution hat ihren
Fokus auf Ressourceneffizienz und Systemsicherheit. Ein weiterer Vorteil wäre
der `diskless mode`, welcher das komplette Betriebssystem in den Arbeitsspeicher
lädt. In diesem Modus müssen Änderungen mit einem *Alpine Local Backup
(lbu)*--Tool explizit auf die Festplatte geschrieben werden. Das hätte den
Vorteil, dass man die Abnutzung des Flash--Speichers, durch unnötige
Schreib/Lese--Vorgänge, minimieren würde.

[^APL]: Alpine Linux für Raspberry Pi: \url{https://wiki.alpinelinux.org/wiki/Raspberry_Pi}

# Literaturverzeichnis
