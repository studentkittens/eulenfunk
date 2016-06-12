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

## Kommerzielle Produkte

Kommerzielle Anbieter von Internet--Radios gibt es wie Sand am Meer. Die
Preisspanne liegt hier zwischen \EUR{30} und mehreren hundert Euro. Die
Funktionsumfang sowie Wiedergabequalität ist hier von Hersteller zu Hersteller
und zwischen den verschiedenen Preisklassen auch sehr unterschiedlich. Einen
aktuellen Überblick aus dem Jahr 2016 über getestete Modelle gibt es
beispielsweise online unter *bestendrei.de*[^TEST].

[^TEST]:Test von Internetradios: \url{http://www.bestendrei.de/elektronik/internetradio/}


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
zum *Raspbarry Pi*--Player hergestellt werden musste. 

[^VOL]: Volumio: https://volumio.org/




# Anforderungen an das Projekt

## Design

Design soll *minimalistisch*  sein.

## Funktionalität

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

## Linux Distribution

Wahl der Linux Distribution

* Debian
* Archlinux

### Installation des Grundsystems

Arch Installation kurz aufführen.

## Hardware


### Raspberry Pi

Der vorhandene Raspberry ist aus dem Jahr 2010. Die genaue Hardware--Revision kann
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

* Blau vs Green

### Rotary--Switch

* Switch von der FH: ALPS irgendwas...funktioniert, aber
* Switch bestellt: ALPS irgendwas mit


### Soundkarte

### Audioverstärkermodul

### RGB--LEDs

* Ansteuerung über GPIO möglich. Zu geringer Strom bei mehreren LEDs.
* Transistorschaltung BC547 NPN anstatt BC557 PNP, da Rückflussstrom.

### USB--Hub

### Netzteil

### Gehäuse

#### Platz im Gehäuse gering

...


# Hardwaredesign

## GPIO--Schnittstelle

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
