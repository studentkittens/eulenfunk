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
lof: no
lot: no
date: \today
---

\newpage
\pagenumbering{arabic} 
\setcounter{page}{1}

# Vorwort 

## Disclaimer

Das vorliegende Projekt, ist im Rahmen einer Studienarbeit im Fach
Hardwaresysteme an der Hochschule Augsburg entstanden. Da die Autoren nicht aus
dem Bereich der *Technischen Informatik* sind, wurden jegliche
hardwarebezogenenen Arbeiten nach bestem Grundlagenwissen umgesetzt.

Diese Studienarbeit soll einen Überblick über die verwendeten, beziehungsweise
benötigten Komponenten für den Bau eines *Raspberry Pi*--Internetradios
verschaffen. Desweiteren soll das Wissen für die Ansteuerung bestimmter
Hardware--Komponenten mittels der *Raspberry Pi*--GPIO[^GPIO] Schnittstelle vermittelt
werden.

## Namensgebung

Der Name des Projektes ist \frqq\texttt{Eulenfunk}\flqq. Die Bezeichnung der
Eule wurde analog zum Tierreich gewählt, da die *Eule* hier als Vogel aufgrund
ihrer Erkennungsmerkmale von anderen Vögeln in der Regel als *Fressfeind*[^EULE]
klassifiziert wird. Analog dazu ist ein *Do-It-Yourself*--Internetradio ---
welches je nach Konfiguration günstiger und mit mehr Funktionalität ausgestattet
werden kann wie ein *Closed--Source*--Produkt --- möglicherweise ein Dorn im
Auge aller kommerziellen Internet--Radio--Anbieter sein könnte.

[^EULE]: Lebensweise der Eule: \url{https://de.wikipedia.org/wiki/Eulen\#Lebensweise}

[^GPIO]: General-purpose input/output Schnittstelle: \url{https://en.wikipedia.org/wiki/General-purpose_input/output}


# Motivation

## Private Situation 

Die Autoren dieses Projekts leben in einer Wohngemeinschaft zusammen. Die Küche
ist der Ort an welchem gemeinsam gekocht und gespeist wird. Für eine angenehme
Atmosphäre und als Nachrichten--Quelle sorgte in der Küche sorgte früher ein
Analog--Radio der Firma *AEG*, welches aufgrund der schlechten Empfangsqualität
durch eine Kombination aus ,,alter Stereoanlage'', ,,altem Raspberry Pi'' und
einem ,,alten Thinkpad x61t'' ersetzt wurde. In dieser Kombination fungierte
die Stereoanlage als Soundausgabe--Komponente, auf dem *Raspberry Pi* lief der
linux--basierte Player Volumio[^VOL], welcher mit dem Touchscreen des *Thinkpad
x61t* über eine Weboberfläche gesteuert wurde. Diese Kombination hat zwar
funktioniert, jedoch war sie alles andere als Benutzerfreundlich, da zuerst die
Stereoanlage und Laptop  eingeschaltet werden mussten und eine WLAN--Verbindung
zum *Raspberry Pi*--Player hergestellt werden musste. 

[^VOL]: Volumio: https://volumio.org/


## Kommerzielle Produkte

Kommerzielle Anbieter von Internet--Radios gibt es wie Sand am Meer. Die
Preisspanne liegt hier zwischen \EUR{30} und mehreren hundert Euro. Die
Funktionsumfang sowie Wiedergabequalität ist hier von Hersteller zu Hersteller
und zwischen den verschiedenen Preisklassen auch sehr unterschiedlich. Einen
aktuellen Überblick aus dem Jahr 2016 über getestete Modelle gibt es
beispielsweise online unter *bestendrei.de*[^TEST].

Das *Problem* bei den kommerziellen Anbietern ist, dass man hier jeweils an die
vorgegebenen Funktionalitäten des Herstellers gebunden ist. Bei einem
Do--it--yourself--Projekt auf Basis Freier Software bzw eines freien
Hardwaredesigns, hat man die Möglichkeit alle gewünschten Funktionalitäten ---
auch Features die von keinem kommerziellen Anbieter unterstützt werden --- zu
integrieren. Beispiele für Funktionalitäten, welche bei kommerziellen Produkten
nur schwer bzw. vereinzelt zu finden sind:

* Unterstützung bestimmter WLAN--Authentifizierungsstandards
* Einhängen von benutzerdefinierten Shares wie *Samba*, *NFS*, *SSHFS*
* Unterstützung verschiedener *lossy* und *lossless* Formate *OGG VORBIS*, *FLAC*, u.a.
* Integration verschiedener Dienste wie beispielsweise *Spotify* möglich
* Benutzerdefinierte Anzeigemöglichkeiten (Uhrzeit, Wetter, et. cetera.)



[^TEST]:Test von Internetradios: \url{http://www.bestendrei.de/elektronik/internetradio/}


## Projektziel

Das grundlegende Projektziel ist aus vorhandenen alten Hardware--Komponenten
ein möglichst optisch und klanglich ansprechendes Internetradio zu entwickeln.
Als Basis für das Projekt dient ein defektes Analog--Radio und ein Raspberry Pi
aus dem Jahr 2012.

# Projektspezifikation

## Hardwareanforderungen

Das Radio soll dem Benutzer folgende Hardwarekonfigurationsmöglichkeiten bieten:

* Anschluß passive Lautsprecher/Kopfhörer möglich 
* Lautstärkeregelung über Hardware möglich
* Verwendung des Internen Lautsprechers des alten Radios
* Statusinformationen über aktuelles Lied beispielsweise über LCD
* LEDs als Statusanzeige und/oder als Visualisierungsvariante von Musik[^MOODBAR]
* USB--Anschlussmöglichkeit für externe Datenträger

[^MOODBAR]: Moodbar: \url{https://en.wikipedia.org/wiki/Moodbar}

## Softwareanforderungen

Die Software soll generisch gehalten werden um eine möglichst einfache
Erweiterbarkeit zu  gewährleisten. 

Hier was zu Menu--Steuerung schrieben und Umfang?

## Optik-- und Usabilityanforderungen

Die Eingabe--Peripherie soll möglichst einfach gehalten werden, um eine *schöne*
Produkt--Optik zu gewährleisten. Folgende 

* Minimale Bedienelemente
* Keine *hässlichen* Knöpfe
* *Retro--Look*-Aussehen wünschenswert

Design soll im Grunde *minimalistisch*  gehalten werden, das heisst, es sollen
nur so wenige ,,Bedienelemente'' wie nötig angebracht werden.

## Kosten/Nutzen--Verhältnis

Für die Erstellung des Projekts sollte bereits vorhandene Komponenten und
Bauelemente wiederverwendet werden um den Kostenaufwand minimal zu halten.

# Hardware

## Komponenten und Bauteile

Folgende Hardwarekomponenten oder Bauteile sind bereits vorhanden oder müssen
noch erworben werden:

* Altes Gehäuse AEG 4104 Küchenradio[^AEG] (vorhanden)
* *Raspberry Pi* aus dem Jahr 2012 (vorhanden)
* LCD--Anzeige (Altbestände u. Arduino--Kit vorhanden)
* Kleinbauteile wie LEDs, Widerstände (vorhanden, Arduino--Kit)
* USB--Hub für Anschluss von beispielsweise ext. Festplatte (vorhanden)
* USB--Soundkarte (vorhanden)
* WIFI--Adapter (vorhanden)
* Netzteil (vorhanden, div. 5V)
* Audioverstärker (muss erworben werden)
* Drehregler (muss erworben werden)
* Farbe und Kunststoffabdeckung für das neue Gehäuse (muss erworben werden)

[^AEG]: AEG Küchenradio 4104: \url{https://www.amazon.de/AEG-MR-4104-Desgin-Uhrenradio-buche/dp/B000HD19W8}



## Raspberry Pi

Der vorhandene Raspberry ist aus dem Jahr 2012. Die genaue Hardware--Revision kann
auf Linux unter ``proc`` ausgelesen werden, siehe auch [@gay2014raspberry]:

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

Laut Tabelle unter [@gay2014raspberry] handelt es sich hierbei um das Modell B
Revision 1+ mit 256MB RAM.

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
GPIO--Pin mit max. 16mA belastbar. Der der gesamte GPIO--Header sollte mit nicht
mehr als 50mA belastet werden, da es darüber hinaus zu Hardwareschäden kommen
kann (siehe auch @gay2014raspberry).

Die Ansteuerung Ansteuerung von LED über GPIO erfolgt binär. Das heisst, dass
die LED entweder aus (GPIO low) oder an sein kann (GPIO high).

TODO: ELCH?
In der ,,analogen'' Welt ist es jedoch möglich eine LED über das Senken der
Spannung zu dimmen. Um ein Dimmen in der digitalen Welt zu erreichen wird ein
Modulationsverfahren angwandt, welches Pulsweitenmodulation heisst. Hierbei
wird...elch? Weitere Informationen finden sich auch unter @richardson2014make
@gay2014raspberry und @gay2014mastering. Software PWM in @gay2014experimenting.
Hier beispielsweise 6% CPU--Last pro GPIO--Pin. 

## LCD--Anzeige

Um dem Benutzer Informationen beispielsweise über das aktuell gespielte Lied
anzeigen zu können, soll eine LCD--Anzeige verbaut werden. In den privaten
Altbeständen finden sich folgende drei Hitachi hd44780--kompatible Modelle:

* Blaues Display, 4x20 Zeichen, Bolymin BC2004A
* Blaues Display, 2x16 Zeichen, Bolymin BC1602A
* Grünes Display, 4x20 Zeichen, Dispalytech 204B 

Für *Eulenfunk* wurde das blaue 4x20 Display --- aufgrund der größeren
Anzeigefläche und Farbe --- gewählt.

### Anschlussmöglichkeiten

Ein LCD Display kann an den Raspberry PI über auf verschiedene Art und Weise
angeschlossen werden. Anschlussmöglichkeiten für eine LCD--Anzeige wären
beispielsweise: 

* GPIO direkt (parallel)
* I2C--Bus (seriell)
* SPI--Bus (seriell)

Die serielle Anschlussmöglichkeit bietet den Vorteil dass weniger Datenleitungen
(GPIO--Pins) verwendet werden. Für den parallelen Betrieb des Displays werden
mindestens sechs GPIO--Pins benötigt, für den seriellen Anschluss über I2C
lediglich nur zwei. 

Da für den seriellen Betrieb beispielsweise über den I2C--Bus zusätzliche
Hardware benötigt wird, wird die parallele Ansteuerung über die GPIO--Pins
bevorzugt. Für weitere Informationen zum seriellen Betrieb über I2C siehe auch
@horan2013practical.

\begin{figure}[h!]
  \centering
\includegraphics[width=0.7\textwidth]{images/lcdraspi.png}
  \caption{Verdrahtung von LCD und Raspberry Pi.}
  \label{lcd}
\end{figure}

Das Display arbeitet mit einer Logik--Spannung von 3.3V - 5V. Da die GPIO--Pins jedoch
eine High--Logik von 3,3V aufweisen, würde man hier in der Regel einen
Pegelwandler bei bidirektionaler Kommunikation und 5V benötigen. Da wir aber nur auf
das Display zugreifen und die GPIO--Pins nicht schreibend zugegriffen wird kann
ein Betrieb des Displays auch mit 5V erfolgen. Beim 3.3V Betrieb welcher laut
Datenblatt auch möglich sein soll, hat das Display leider nur eine sehr
schwachen beziehungsweise unzureichenden Darstellungskontrast gehabt, weswegen
der 5V Betrieb gewählt wurde. 

Die Hintergrundbeleuchtung des Displays wurde direkt über ein Potentiometer mit
10K$\Omega$ an die 5V Spannungsversogrung angeschlossen. Es wurde hier dierekt
die Speisung vom Netzteil gewählt um den GPIO--Header nicht unnötig zu belasten.

Laut Datenblatt[^LCD] kann die Hintergrundbeleuchtung entweder mit 3.4V ohne
Vorwiderstand oder mit 5V bei einem 27$\Omega$ Widerstand betrieben werden. Damit das
Display beim herunter geregeltem Potentiometer keinen Schaden nimmt, wurden
zusätzlich zwei Widerstände mit 100$\Omega$ (parallel geschaltet = 50$\Omega$) zwischen Display
und Potentiometer gehängt.

[^LCD]: Datenblatt Bolymin BC2004A: \url{http://www.dema.net/pdf/bolymin/BC2004A-series_VER04.pdf}

Das der resultierende Gesamtwiderstand ohne Potentiometer beträgt in diesem Fall
$\approx$ 50 $\Omega$:

$$  R_{ges} = \frac{R_1 \times R_2}{R_1 + R_2} = \frac{100\Omega \times 100\Omega}{100\Omega + 100\Omega} = 50\Omega $$

## Rotary--Switch

Für eine minimale Anzahl an Bedienelementen zu erhalten, wird bei *Eulenfunk*
ein Drehimpulsgeber mit Schalter gewählt. Für erste Testzwecke wurde vom Herrn
Schäferling ein *ALPS STEC12E08* bereitgestellt. Dieser wurde im Laufe der
Entwciklung durch einen *ALPS STEC11B09*[^ALPS] ersetzt, da dieser mittels Mutter und
Schraube am Gehäuse besser befästigt werden kann. 

Der verwendete Drehimpulsgeber hat insgesamt fünf Anschlüsse. Zwei
Signalleitungen (A und B), zwei mal *GND* (für Drehgeber und Schalter jeweils
eine) und einen Anschluss für den Schalter. Beim drehen eines Drehimpulsgebers
wird ein Rechtecksignal generiert. Je nach Muster der beiden Datensignale A oder
B, kann entschieden werden ob es sich um eine Rechts-- oder Linksdrehung
handelt. Weitere Hintergrundinformationen zu Drehimpulsgeber siehe auch
@2014projekte.

Abbildung \ref{alps} zeigt den Anschluss des Drehimpulsgebers am *Raspberry Pi*. 

[^ALPS]: Drehimpulsgeber ALPS STEC11B09: \url{https://www.reichelt.de/Drehimpulsgeber/STEC11B09/3/index.html?ACTION=3&GROUPID=3714&ARTICLE=73915}

\begin{figure}[h!]
  \centering
\includegraphics[width=0.5\textwidth]{images/rotary.png}
  \caption{Drehimpulsgeber--Anschluss an den Raspberry Pi}
  \label{alps}
\end{figure}

## Soundkarte

Die interne Soundkarte des *Raspberry Pi* ist über eine triviale
Pulsweitenmodulation realisiert. Die einfache Schaltung soll hier laut
Internetquellen[^AQ]eine sehr niedrige Audioqualität bieten.

Aus diesem Grund wird bei *Eulenfunk* auf das USB--Audio--Interface *BEHRINGER
U-PHONO UFO202*[^DAC] gesetzt. 

[^DAC]:BEHRINGER U-PHONO UFO202 Audio Interface: \url{http://www.produktinfo.conrad.com/datenblaetter/1300000-1399999/001370864-an-01-de-BEHRINGER_UFO_202_AUDIOINTERFACE.pdf}
[^AQ]: Raspberry Pi onboard Sound: \url{http://www.crazy-audio.com/2013/11/quality-of-the-raspberry-pi-onboard-sound/}

## Audioverstärkermodul

Da eine Soundkarte in der Regel zu wenig Leistung hat um einem Lautsprecher
,,vernünftig'' anzusteuern wird ein Audioverstärker benötigt. Da neben dem
Anschluss von externen Lautsprechern auch eine Lautstärkeregelung über ein Poti
erfolgen soll, ist die Entscheidung einfachheitshalber auf ein
Audioverstärker--Modul auf Basis vom PAM8403[^POW] Stereo-Verstärker mit Potentiometer gefallen.

[^POW]: Verstärkermodul: \url{https://www.amazon.de/5V-Audioverstärker-Digitalendstufenmodul-Zweikanalige-Stereo-Verstärker-Potentiometer/dp/B01ELT81A6}

## LED--Transistorschaltung

Die Ansteuerun einer LED mittels GPIO--Pin ist recht simpel. Sollen jedoch
mehrere LEDs angesteuert werden, so wird in der Regel pro LED ein GPIO--Pin
benötigt. LEDs sollten nie ohne Vorwiderstand an den *Raspberry Pi*
angeschlossen werden, da durch den hohen Stromfluss die LED beschädigt werden
könnte. Desweiteren muss bei LEDs auch auf die Polung geachtet werden, die
abgeflachte Seite --- meist mit dem kürzerem Beinchen -- ist in der Regel die
Kathode (Minuspol). Abbildung \ref{led} zeigt exemplarisch den Anschluss einer
*classic LED rot*[^LEDS], mit einer Flussspannung von $U_{LED}$ $\approx$ 2V, die mit
einem Strom von $I_{LED}$ = 20 mA gespeist werden soll. Die Berechnung des
Vorwiderstandes erfolgt nach folgender Formel:

$$R_{LED} = \frac{U_{GPIO}-U_{LED}}{I_{LED}} = \frac{3.3V - 2V}{20mA}   \approx 65\Omega$$

[^LEDS]: Datenblatt mit verschiedenen LED--Typen: \url{https://www.led-tech.de/de/5mm-LEDs_DB-4.pdf}

**Hinweis:** Da ein GPIO--Pin aber mit nur max. 16mA belastet werden sollten,
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
  \caption{Anschluss eine roten LED mit Vorwiderstand am Raspberry Pi GPIO--Pin}
  \label{led}
\end{figure}

Je nach Typ und Farbe ist der benötigte Strom um ein vielfaches höher wie in
unserem Beipsiel. Die in \ref{led} abgebildete LED kann vom GPIO--Pin nur einen
max. Strom von 16 mA beziehen

In *Eulenfunk* sollen mehrere intensiv leuchtende LEDs verbaut werden. Da die
GPIO--Pins in ihrer Leistung sehr begrenzt ist, würde es sich anbeiten eine
externe Stromquelle zu verwenden. Um die Speisung über eine externe Stromquelle
zu ermöglichen kann eine Transistorschaltung verwendet werden (vgl. @traled). 

Für die Transistorschaltung wurden vom Herrn Schäferling NPN (BC547C) und PNP
(BC557C) bereitgestellt. Für den ersten Testaufbau wurde der PNP--Transistor und
eine RGB--LED[^RGBGM] mit gemeinsamen Minuspol verwendet. Beim Test--Aufbau mit einem
PNP--Transistor ist aufgefallen, dass die LED ständig geleuchtet hat. Eine kurze
Recherche hat ergeben, dass der Transistor permanent durchgeschaltet war, weil die
Spannung an der Basis (GPIO--Pin, 3,3V) geringer war die die Betriebsspannung
für die LED (5V). 

Der zweite Anlauf mit dem NPN--Transistor BC547C und einer RGB--LED[^RGBGP] mit
gemeinsamen Pluspol hat das gewünsche Ergebnis geliefert.

Da der Hersteller für die von der Hochschule bereitgestellten Transistoren
ubekannt ist, wurden typische Durchschnittswerte für die Dimensionierung der
Restlichen Bauteile verwendet.

Wie es aussieht sind die meisten BC547C Transistor Typen für einen max. Strom
$I_{CE}$=100 mA konstruiert. Für die Berechnung des Basis--Vorwiderstandes ist der
Stromverstärkungsfaktor $h_{FE}$[^HFE] benötigt. Je nach Hersteller variieren die
Werten zwischen 200[^SEM] und 400[^FARI]. Da der maximale Laststrom $I_{CE}$ pro Transistor
beträgt 80 mA (3 LEDs je max. 20mA), sieht die Berechnung des Basisstroms wie
folgt aus:

[^LED]: Beispiele zur Ansteuerung von LEDs: \url{http://www.led-treiber.de/html/vorwiderstand.html}
[^HFE]:Stromverstärkungsfaktor: \url{http://www.learningaboutelectronics.com/Articles/What-is-hfe-of-a-transistor}
[^RGBGM]: RGB-LED Common Cathode: \url{http://download.impolux.de/datasheet/LEDs/LED 0870 RGB 5mm klar 10000mcd.pdf}
[^RGBGP]: RGB-LED Common Anode: \url{http://download.impolux.de/datasheet/LEDs/LED 09258 RGB 5mm klar 10000mcd_GP.pdf}


 $$I_{Basis} = \frac{I_{CE}}{h_{FE}} = \frac{0.08A}{300} \approx 270\mu A$$

Der BC547C Transitor benötigt eine durchschnittliche  $U_{BE}$ = 0,7V zum
durchschalten. Die GPIO-Pins des *Raspberry Pi* haben einen Spannungspegel von
3.3V. Daraus ergibt sich folgende Berechnung des Basis--Vorwiderstandes:

$$R_{Basis} = \frac{U_{GPIO} - U_{Basis}}{I_{Basis}} = \frac{3,3V - 0,7V}{270\mu A} = 9629 \Omega \approx 10k \Omega $$

[^SEM]: SEMTECH: \url{http://pdf1.alldatasheet.com/datasheet-pdf/view/42386/SEMTECH/BC547.html}
[^FARI]: Farichild Semiconductor: \url{https://www.fairchildsemi.com/datasheets/BC/BC547.pdf}

\begin{figure}[h!]
  \centering
\includegraphics[width=0.7\textwidth]{images/transistorled.png}
  \caption{Transistor--RGB--LED Schaltung}
  \label{transled}
\end{figure}

Die in Abbildung \ref{transled} gelisteten LED--Vorwiderstände ergeben sich
aufgrund der verschiedenen Spannungen der unterschiedlichen Farben[^RGBGP]. Die
Berechnung für den Vorwiderstand pro LED schaut am Beispiel der Farbe blau
($U_{LED} = 3,15V, I_{LED} = 20mA$) wie folgt aus:

$$R_{LED} = \frac{U_{Betriebsspannung} - U_{LED}}{I_{LED}} = \frac{5V - 3,15V}{20mA} =92.5 \approx 100\Omega$$


## USB--Hub und Netzteil

Der *Rapberry Pi* hat in unserer Revision nur zwei USB--Schnittstellen, diese
sind bereits durch die Hardware--Komponenten USB--DAC (Soundkarte) und das
WIFI--Modul belegt. Um den Anschluss eines externen Datenträgers, auch mit
größerer Last wie beispielsweise einer Festplatte zu ermöglichen wird ein
aktiver USB--Hub benötigt.

Für diesen Einsatzzweck wird aus den Altbeständen ein *LogiLink 4 Port USB 2.0
HUB* verwendet. Viele billig-Hubs arbeiten hier entgegen der USB--Spezifikation
und speisen zusätzlich über die USB--Schnittstellen den *Raspberry Pi*. Dieses
Verhalten wurde bemerkt, also der *Raspberry Pi* ohne Power--Connector alleine
mit nur der USB--Verbindung zum Hub bootete.

Da bei der Speisung über die USB--Schnittstelle die interne Sicherungschaltung
des *Pi* umgangen werden, besteht hier die zusätzliche Gefahr eines
Hardwaredefektes durch die Speisung einer zusätzlichen Spannungsquelle. Weitere
Informationen zu dem Problem finden sich unter @suehle2014hacks.


## Gehäuse

Die Gehäuse--Farbe soll in hellelfenbeinweiß RAL 1015 einen dezenten
,,Retro''--Look verschaffen.
Plexiglas von Wolfgang
Holzgehäuse des alten AEG Radios
Knöpfe schwarz mit Alu-Optik

### Platz im Gehäuse gering

...

## Betriebssystem

Mittlerweile gibt es für den *Raspberry Pi* viele offiziell zugeschnittene
Betriebssysteme[^OS] (vgl. @dembowski2013raspberry). Bei den den Linux
Distributionen ist *Raspbian* eine der bekanntesten Distribution -- welche auf
*Debian* basiert. *Raspbian* bringt ein komplettes Linux--basiertes System mit
grafischer Benutzeroberfläche mit sich. 

Neben den unter[^OS] genannten Distributionen gibt es mittlerweile auch Windows
10 IoT (Internet of Things) für den *Rapberry Pi*. Dieses speziell für den
Embedded Bereich ausgerichtete Windwos benötigt jedoch eine ARMv7--CPU als
Mindestanforderung, was unseren ,,alten Raspberry'' ausschließen würde.
Desweiteren wäre für uns eine propritäre Lösung ein K.O.--Kriterium, da
diese alle Vorteile von Freier Software zu nichte machen würde.

[^WINIOT]: ARMv7 für Win 10 IoT:  \url{http://raspberrypi.stackexchange.com/questions/39715/can-you-put-windows-10-iot-core-on-raspberry-pi-zero}

### Wahl des Betriebbsystems

*Arch Linux ARM* ist eine minimalistische und sehr performante
Linux--Distribution welche im Gegensatz zu *Raspbian* ohne Desktopumgebung
geliefert wird (vgl. @schmidt2014raspberry) Desweiteren ist *Arch Linux* ein
bekannter Vertreter von Rolling--Release--Distributionen. Ein weiterer Vorteil
für unseren Einsatzzweck hier ist bei *Arch Linux* das *AUR* (Arch User
Repository)[^AUR], dieses erlaubt es eigene Software auf eine schnelle und
unkomplizierte Weise der Allgemeinheit zur Verfügung zu stellen.

### Einrichtung des Grundsystems

Nach der Installation und dem ersten Booten des Grundsystems muss die
Netzwerk--Schnittstelle konfiguriert werden. Arch Linux ARM bietet mit *netctl*
ein profilbasierte Konfigurationsmöglichkeit. Ein Profil kann über das
*ncurses*--basierte Tool `wifi-menu` erstellt werden. In unserem Fall wurde das
Profil `wlan0-Phobos` erstellt. Anschließend kann das erstellte Profil mit
*netctl* verwendet werden. 

**Auflistung der bekannten Profile**

```bash
    [alarm@eulenfunk ~]$ netctl list
      eth0-static
      wlan0-Phobos
```

**Aktivierung des gewünschten Profils**

```bash
    # Starten des gewünschten Profils
    [alarm@eulenfunk ~]$ netctl start wlan0-Phobos

    [alarm@eulenfunk ~]$ netctl list
      eth0-static
    * wlan0-Phobos

    # Profil über System-Reboot hinweg aktivieren 
    [alarm@eulenfunk ~]$ netctl enable wlan0-Phobos

```

Nun verbindet sich der *Raspberry Pi* nach dem Hochfahren jedes Mal automatisch
mit dem Profil `wlan0-Phobos`.

[^AUR]: Arch User Repository: \url{https://aur.archlinux.org/}
[^ARCH]: Arch Linux ARM: \url{https://archlinuxarm.org/}
[^OS]: Betriebssystem--Images Raspberry Pi: \url{https://www.raspberrypi.org/downloads/}
[^INSTALL]: Arch Linux Installation für Raspberry Pi: https://archlinuxarm.org/platforms/armv6/raspberry-pi#installation


### Abspielsoftware

Für den Betrieb des Internetradios soll der Music--Player--Daemon verwendet
werden, da *Eulenfunk* auf einem eigens entwickeltem MPD--Client basieren soll
(mehr zur Eulenkfunk Software siehe Kapitel Software).
Andere Projekte greifen oft auf Abspielsoftware wie den MOC oder Mplayer zu. 


```bash
    # Installation des MPD
    [root@eulenfunk ~]$ pacman -Sy mpd mpc ncmpc

```

# Software

## Vorhandene Softwarelibraries

## Überblick der einzelnen Komponenten

## Softwarearchitektur

## Treiber--Software

### LCD--Treiber

Von Elchen entwickelt.

### Rotary--Treiber

Von Elchen kopiert.

### LED--Treiber

* Software--PWM


# Zusammenfassung

## Ziel erreicht?

Ja?

## Erweiterungen und alternative Ansätze

### Audio--Visualisierung

Beim Projekt *Eulenfunk* wird die Visualisierung von Musik aufgrund der
begrenzten Zeit und Hardwareressourcen des *Raspberry Pi *über eine
vorberechnete Moodbar--Datei realisiert. Dieser Ansatz funktioniert bei nicht
live gestreamter Musik gut. Bei live gestreamter Musik könnte für die
Visualisierung eine Fast--Fourier--Transformation in Echtzeit durchgeführt
werden. Da jedoch die Ressourcen des *Raspberry Pi* sehr begrenzt sollte hier
auf die Verwendung einer GPU--beschleunigte--FFT[^FFT] zurückgegriffen werden
(vgl. @Sabarinath2015).

Ein alternativer Ansatz wäre auch die Realisierung einer Musik--Visualisierung
mittels Hardwarekomponenten. Ein möglicher Ansatz aus hardwarebasierten
Hochpass-- und Tiefpassfiltern in Form einer Disco--Beleuchtung wird unter @2014projekte beschrieben.

### Echtzeit--Uhr

Der *Raspberry Pi* besitzt keine Hardware--Uhr. Aufgrund der Tatsache dass es
sich bei *Eulenfunk* um eine Internet--Radio handelt wurde auf eine Echtzeituhr
(real time clock, RTC) verzichtet, da sich die Uhr *Eulenfunk* aufgrund der
permanenten Internetverbindung mittels NTP[^NTP] über das Internet
synchronisieren kann. Eine Erweiterung um eine Echtzeituhr wird in
@horan2013practical und @gay2014experimenting ausführlich beschreiben.

[^NTP]: Network Time Protocol:
\url{https://de.wikipedia.org/wiki/Network_Time_Protocol}
[^FFT]: GPU--beschleunigte FFT auf dem Raspberry Pi:
\url{http://www.aholme.co.uk/GPU_FFT/Main.htm}

## Mögliche Verbesserungen?

### Alpine Linux 

Die relativ junge Linux--Distribution *Alpine Linux*[^APL] wäre eine Mögliche
Verbesserung für den Einsatz des Internetradios. Diese Distribution hat ihren
Fokus auf Ressourceneffizienz und Systemsicherheit. Ein weiterer Vorteil wäre
der `diskless mode`, welcher das Komplette Betriebssystem in den Arbeitsspeicher
lädt. In diesem Modus müssen Änderungen mit einem commit-Tool explizit auf die
Festplatte geschrieben werden. Das hätte den Vorteil, dass man die Abnutzung des
Flash--Speichers, durch unnötige Schreib/Lese--Vogänge, minimieren würde.

[^APL]: Alpine Linux für Raspberry Pi: \url{https://wiki.alpinelinux.org/wiki/Raspberry_Pi}

# Literaturverzeichnis
