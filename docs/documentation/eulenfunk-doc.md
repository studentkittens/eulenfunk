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
* Integration verschiedener Dienste wie beispielsweise *Spotify*
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
* Verwendung des Internen Lautsprechers des alten Radios
* Dimmbare RGB--LED Statusanzeige 

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

### GPIO--Schnittstelle

\begin{figure}[h!]
  \centering
  \includegraphics[width=0.7\textwidth]{images/gpio.png}
  \caption{GPIO-Header des Raspberry Pi Modell B Rev 1.0+}
  \label{gpio}
\end{figure}




Bildquelle: \url{http://www.raspberrypi-spy.co.uk/2012/06/simple-guide-to-the-rpi-gpio-header-and-pins/#prettyPhoto}

Die \ref{gpio} ist

#### GPIO--Pinbelegung

* Grafik

...

* 3,3V vs 5V
* Max. Strom
* Max. verfügbare Pins

## LCD--Anzeige

Um dem Benutzer Informationen beispielsweise über das aktuell gespielte Lied
anzeigen zu können, soll eine LCD--Anzeige verbaut werden. In den privaten
Altbeständen finden sich folgende drei hd44780--kompatible Modelle:

* Blaues Display, 4x20 Zeichen, Bolymin BC204A
* Blaues Display, 2x16 Zeichen, Bolymin BC1602A
* Grünes Display, 4x20 Zeichen, Dispalytech 204B 

Wir haben uns für das blaue 4x20 Display --- aufgrund der größeren Anzeigefläche und
Farbe --- entschieden.

### Anschlussmöglichkeiten

Ein LCD Display kann an den Raspberry PI über auf verschiedene Art und Weise
angeschlossen werden. Die zwei Grundlegenden Schnittstellen wären:

* GPIO (parallel)
* I2C--Bus (seriell)
* SPI--Bus (seriell)

Da für den seriellen Betrieb beispielsweise über den I2C--Bus ein zusätzlicher
Logik--Adapter benötigt wird, wird die parallele Ansteuerung über die GPIO--Pins
bevorzugt.

\begin{figure}[h!]
  \centering
\includegraphics[width=0.7\textwidth]{images/lcdraspi.png}
  \caption{Verdrahtung von LCD und Raspberry Pi.}
  \label{lcd}
\end{figure}

Das Display arbeitet mit einer Logik--Spannung von 5V. Da die GPIO--Pins jedoch
eine High--Logik von 3,3V aufweisen, würde man hier in der Regel einen
Pegelwandler bei bidirektionaler Kommunikation benötigen. Da wir aber nur auf
das Display zugreifen und die GPIO--Pins nicht schreiben zugegriffen wird kann
eine operation des Displays auch mit 3.3V erfolgen, falls die GPIO--Pegel
ausreichen um die Logik des Displays anzusteuern.

Die Hintergrundbeleuchtung des Displays wurde direkt über ein Potentiometer mit
10K$\Omega$ an die 5V Spannungsversogrung angeschlossen.

Laut Datenblatt kann die Hintergrundbeleuchtung entweder mit 3.4V ohne
Vorwiderstand oder mit 5V bei einem 48$\Omega$ Widerstand betrieben werden. Damit das
Display beim herunter geregeltem Potentiometer keinen Schaden nimmt, wurden
zusätzlich zwei Widerstände mit 100$\Omega$ (parallel geschaltet) zwischen Display
und Potentiometer gehängt.

Das der resultierende Gesamtwiderstand ohne Potentiometer beträgt in diesem Fall
$\approx$ 50 $\Omega$:

$$  R_{ges} = \frac{R_1 \times R_2}{R_1 + R_2} = \frac{100\Omega \times 100\Omega}{100\Omega + 100\Omega} = 50\Omega $$

## Rotary--Switch

* Switch von der FH: ALPS irgendwas...funktioniert, aber
* Switch bestellt: ALPS irgendwas mit


\begin{figure}[h!]
  \centering
\includegraphics[width=0.5\textwidth]{images/rotary.png}
  \caption{Drehimpulsgeber--Anschluss an den Raspberry Pi}
  \label{alps}
\end{figure}

## Soundkarte

Die interne Soundkarte des *Raspberry Pi* ist über eine triviale
Pulsweitenmodulation realisiert. Die einfache Schaltung soll hier eine sehr
niedrige Audioqualität[^AQ] bieten.

[^AQ]: Raspberry Pi onboard Sound: \url{http://www.crazy-audio.com/2013/11/quality-of-the-raspberry-pi-onboard-sound/}

## Audioverstärkermodul

## RGB--LEDs

* Ansteuerung über GPIO möglich. Zu geringer Strom bei mehreren LEDs.
* Transistorschaltung BC547 NPN anstatt BC557 PNP, da Rückflussstrom.

Von der Hochschule wurden BC547C Transistoren bereitgestellt. Da der
Hersteller unbekannt ist, wurden typische Durchschnittswerte für die
Dimensionierung der Bauteile verwendet.  

In der Regel sind die meisten BC547C Transistor Typen für einen max. Strom
$I_{CE}$=100 mA konstruiert. Für die Berechnung des Basis--Vorwiderstandes ist der
Stromverstärkungsfaktor $h_{FE}$[^HFE] benötigt. Je nach Hersteller variieren die
Werten zwischen 200 und 500[^SEM][^ONSEMI][^FARI]. Da ein maximale Laststrom $I_{CE}$ pro Transistor
beträgt 80 mA (3 LEDs je max. 20mA).

[^HFE]:Stromverstärkungsfaktor: \url{http://www.learningaboutelectronics.com/Articles/What-is-hfe-of-a-transistor}

Die Berechnung des Basisstroms:

 $$I_{Basis} = \frac{I_{CE}}{h_{FE}} = \frac{0.08A}{300} \approx 270\mu A$$

Der BC547C Transitor benötigt eine durchschnittliche  $U_{BE}$ = 0,7V. Die
GPIO-Pins des *Raspberry Pi* haben eine Spannungspegel von 3.3V. Daraus ergibt
sich folgende Berechnung des Basis--Vorwiderstandes:

$$R_{Basis} = \frac{U_{GPIO} - U_{Basis}}{I_{Basis}} = \frac{3,3V - 0,7V}{270mA}
= 9629 \Omega \approx 10k \Omega $$

[^SEM]: SEMTECH: \url{http://pdf1.alldatasheet.com/datasheet-pdf/view/42386/SEMTECH/BC547.html}
[^ONSEMI]: On Semi: \url{https://www.arduino.cc/documents/datasheets/BC547.pdf}
[^FARI]: Farichild Semiconductor:
\url{https://www.fairchildsemi.com/datasheets/BC/BC547.pdf}

\begin{figure}[h!]
  \centering
\includegraphics[width=0.7\textwidth]{images/transistorled.png}
  \caption{Transistor--RGB--LED Schaltung}
  \label{transled}
\end{figure}

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
Hardwaredefektes durch die Speisung einer zusätzlichen Spannungsquelle. Da in
unserem Fall jedoch nur eine Spannungsquelle existiert, wird das Problem als
vernachlässigbar klassifiziert.


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
Betriebssysteme[^OS]. Bei den den Linux Distributionen ist *RASPBIAN* eine der
bekanntesten Distribution -- welches auf *Debian* basiert.

Da *Debian* seinen Fokus auf ,,Stabilität'' hat, sind die Pakete der
Distribution des Öfteren älter. Ein weiterer ,,Nachteil'' sind aktuell feste
Release--Zyklen, wünschenswert wäre eine Rolling--Release Distribution, welche
nur einmalig installiert werden muss und kontinuierlich auf den aktuellsten
Stand geupdated werden kann.

Eine bekannter Vertreter von Rolling--Release--Distributionen ist *Arch Linux*,
von welcher es auch einen ARM--Port[^ARCH] gibt. Ein weiterer Vorteil ist bei
*Arch Linux* das *AUR* (Arch User Repository)[^AUR], dieses erlaubt es eigene Software
auf eine schnelle und unkomplizierte Weise der Allgemeinheit zur Verfügung zu
stellen.

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

## Mögliche Verbesserungen?

* Alpine Linux da RAM--only Betrieb möglich

# Literaturverzeichnis
