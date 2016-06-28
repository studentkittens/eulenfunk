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
\includegraphics[width=0.7\textwidth]{images/eulenfunk.png}
  \caption{Foto vom aktuellen Prototypen}
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
* *Imagemagick* (Bildbearbeitung) 


[^VIMEO]: Eulenfunk Prototyp: \url{https://vimeo.com/171646691}
[^EULE]: Lebensweise der Eule: \url{https://de.wikipedia.org/wiki/Eulen\#Lebensweise}
[^GPIO]: General-purpose input/output Schnittstelle: \url{https://en.wikipedia.org/wiki/General-purpose_input/output}
