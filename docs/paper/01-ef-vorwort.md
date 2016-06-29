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
tocdepth: 1
lof: yes
lot: no
date: \today
---

\newpage
\pagenumbering{arabic} 
\setcounter{page}{1}

# Vorwort 

## Hinweis zu den Vorkenntnissen

Das vorliegende Projekt ist im Rahmen einer Studienarbeit im Fach
Hardwaresysteme an der Hochschule Augsburg entstanden. Da die Autoren nicht aus
dem Bereich der *Technischen Informatik* sind, wurden jegliche
hardwareseitigen Arbeiten nach bestem Grundlagenwissen umgesetzt.

## Namensgebung

Der Name des Projektes ist \frqq\texttt{Eulenfunk}\flqq. Passend erschien uns
dieser Name, weil er mit dem Spitznamen einer der Autoren zu tun hat. Auch sind
unserer Meinung nach Eulen sehr faszinierende Tiere, die durch eine gewisse
optische Ästhetik und interessante, manchmal verwirrende, Verhaltensstrukturen
auffallen. Der Suffix »-funk« verleiht dem Name zudem durchaus etwas
majestätisches.

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
  \caption{Foto vom aktuellen Prototypen.}
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
