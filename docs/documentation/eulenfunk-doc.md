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
hardwarebezogenenen soweit möglich nach bestem Wissen und Grundlagen--Wissen umgesetzt.

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
* Unterstützung verschiedener *lossy* Formate *mp3*, *ogg vorbis*, *acc*, u.a.
* Unterstützung verschiedener *lossless* Formate *FLAC*, *APE*, u.a.
* Integration verschiedener Dienste wie beispielsweise *Spotify*
* Benutzerdefinierte Anzeigemöglichkeiten (Uhrzeit, Wetter, et. cetera.)



[^TEST]:Test von Internetradios: \url{http://www.bestendrei.de/elektronik/internetradio/}


## Projektziel

Das grundlegende Projektziel ist aus vorhandenen alten Hardware--Komponenten
ein möglichst optisch und klanglich ansprechendes Internetradio zu entwickeln.
Als Basis für das Projekt dient ein defektes Analog--Radio und ein Raspberry Pi
aus dem Jahr 2012.

# Anforderungen an das Projekt

## Design

Design soll *minimalistisch*  sein, das heisst, es sollen so wenige
,,Bedienelemente'' wie nötig untergebracht werden. 

## Funktionsumfang (Hardware)

* Aktive Lautsprecher
* Passive Lautsprecher/Kopfhörer
* Verwendung des Internen Lautsprechers des alten Radios

## Bedienbarkeit

* Minimale Bedienelemente
* Keine *hässlichen* Knöpfe
* *Retro*--Like Aussehen wünschenswert

## Kosten/Nutzen--Verhältnis

Nutzung bereits vorhandener Bauelemente.

# Hardware-- und Softwarekomponenten

## Betriebssystem

Mittlerweile gibt es für den Raspberry Pi viele offiziell zugeschnittene
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


## Hardware


### Raspberry Pi

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

### LCD--Anzeige

Um dem Benutzer Informationen beispielsweise über das aktuell gespielte Lied
anzeigen zu können, soll eine LCD--Anzeige verbaut werden. In den privaten
Altbeständen finden sich folgende drei hd44780--kompatible Modelle:

* Blaues Display, 4x20 Zeichen, Bolymin BC204A
* Blaues Display, 2x16 Zeichen, Bolymin BC1602A
* Grünes Display, 4x20 Zeichen, Dispalytech 204B 

Wir haben uns für das blaue 4x20 Display --- aufgrund der größeren Anzeigefläche und
Farbe --- entschieden.

#### Anschlussmöglichkeiten

Ein LCD Display kann an den Raspberry PI über auf verschiedene Art und Weise
angeschlossen werden. Die zwei Grundlegenden Schnittstellen wären:

* GPIO (parallel)
* I2C--Bus (seriell)
* SPI--Bus (seriell)

Da für den seriellen Betrieb beispielsweise über den I2C--Bus ein zusätzlicher
Logik--Adapter benötigt wird, wird die parallele Ansteuerung über die GPIO--Pins
bevorzugt.

\includegraphics[width=\linewidth]{images/lcdraspi.png}


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

### Rotary--Switch

* Switch von der FH: ALPS irgendwas...funktioniert, aber
* Switch bestellt: ALPS irgendwas mit


### Soundkarte

Die interne Soundkarte des *Raspberry Pi* ist über eine triviale
Pulsweitenmodulation realisiert. Die 'einfache' Schaltung soll hier eine sehr
niedrige Audioqualität bieten.


### Audioverstärkermodul

### RGB--LEDs

* Ansteuerung über GPIO möglich. Zu geringer Strom bei mehreren LEDs.
* Transistorschaltung BC547 NPN anstatt BC557 PNP, da Rückflussstrom.

\includegraphics[width=\linewidth]{images/transistorled.png}

### USB--Hub

### Netzteil

### Gehäuse

Die Gehäuse--Farbe soll in hellelfenbeinweiß RAL 1015 einen dezenten
,,Retro''--Look verschaffen.
Plexiglas von Wolfgang
Holzgehäuse des alten AEG Radios
Knöpfe schwarz mit Alu-Optik

#### Platz im Gehäuse gering

...


# Hardwaredesign

## GPIO--Schnittstelle

\includegraphics[width=\linewidth]{images/gpio.png}

Bildquelle: \url{http://www.raspberrypi-spy.co.uk/2012/06/simple-guide-to-the-rpi-gpio-header-and-pins/#prettyPhoto}

### GPIO--Pinbelegung

* Grafik

...

* 3,3V vs 5V
* Max. Strom
* Max. verfügbare Pins

# Softwaredesign

## Vorhandene Softwarelibraries

## Grundlgender Aufbau

# Überblick der einzelnen Komponenten?

.. .

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
