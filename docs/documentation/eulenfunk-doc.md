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

[^EULE]: Lebensweise der Eule: https://de.wikipedia.org/wiki/Eulen#Lebensweise

[^GPIO]: General-purpose input/output Schnittstelle: \url{https://en.wikipedia.org/wiki/General-purpose_input/output}


# Motivation

## Kommerzielle Produkte

Kommerzielle Anbieter von Internet--Radios gibt es wie Sand am Meer. Die
Preisspanne liegt hier zwischen \EUR{30} und mehreren hundert Euro. Die
Funktionsumfang sowie Wiedergabequalität ist hier von Hersteller zu Hersteller
und zwischen den verschiedenen Preisklassen auch sehr unterschiedlich. Einen
aktuellen Überblick aus dem Jahr 2016 über getestete Modelle gibt es
beispielsweise online unter bestendrei.de[^TEST].

[^TEST]:Test von Internetradios: \url{http://www.bestendrei.de/elektronik/internetradio/}


## Private Situation 

Die Autoren dieses Projekts leben in einer Wohngemeinschaft zusammen. Die Küche
ist der Ort an welchem gemeinsam gekocht und gespeist wird. Für eine angenehme
Atmosphäre und als Nachrichten--Quelle sorgte in der Küche sorgte früher ein
Analog--Radio der Firma *AEG*, welches aufgrund der schlechten Empfangsqualität
durch eine Kombination aus ,,alter Stereoanlage'', ,,altem Raspberry Pi'' und
einem ,, alten Thinkpad x61t'' ersetzt wurde. In dieser Kombination fungierte
die Stereoanlage als Soundausgabe--Komponente, auf dem *Raspberry Pi* lief der
linux--basierte Player Volumio[^VOL], welcher mit dem Touchscreen des *Thinkpad
x61t* über eine Weboberfläche gesteuert wurde. Diese Kombination hat zwar
funktioniert, jedoch war sie alles andere als Benutzerfreundlich, da zuerst die
Stereoanlage und Laptop  eingeschaltet werden mussten und eine WLAN--Verbindung
zum *Raspbarry Pi*--Player hergestellt werden musste. 

[^VOL]: Volumio: https://volumio.org/




# Anforderungen an das Projekt

# Wahl der Hardware-- und Softwarekomponenten

# Hardwaredesign

## Ansteuerung GPIO

# Softwaredesign

# Literaturverzeichnis
