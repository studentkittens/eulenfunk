# Zusammenfassung

## Ziel erreicht?

Das selbstgesetzte Ziel --- mit möglichst geringen finanziellen Einsatz ein Internetradio auf Basis
eines *Raspberry Pi* zu entwickeln --- kann durchaus als erfolgreich betrachtet
werden. 

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

### Wettervorhersage

Eine praktische Erweiterung von *Eulenfunk* wäre die Anzeige der
Wettervorhersage für den aktuellen und die nächsten Tage. Oft wäre eine
entsprechende Information nützlich, um beispielsweise die Durchführung einer
Radtour zu planen oder ob man trockenen Fußes einkaufen gehen kann. Die Daten
könnte dabei von Online--APIs wie OpenWeatherMap[^OWM] geholt werden. Der Zugang
ist dort kostenlos, aber auf 60 Anfragen pro Minute limitiert.

[^OWM]: Mehr Informationen unter: \url{http://www.openweathermap.org/}

## Mögliche Verbesserungen

### Alpine Linux 

Die relativ junge Linux--Distribution *Alpine Linux*[^APL] wäre eine mögliche
Verbesserung für den Einsatzzweck Internetradio. Diese Distribution hat ihren
Fokus auf Ressourceneffizienz und Systemsicherheit. Ein weiterer Vorteil wäre
der `diskless mode`, welcher das komplette Betriebssystem in den
Arbeitsspeicher lädt. In diesem Modus müssen Änderungen mit einem *Alpine Local
Backup (lbu)*--Tool explizit auf die Festplatte geschrieben werden. Das hätte
den Vorteil, dass man die Abnutzung des Flash--Speichers, durch unnötige
Schreib/Lese--Vorgänge, minimieren würde. Momentan unterstützt diese
Distribution allerdings noch nicht das von uns favorisierte ``systemd``.

[^APL]: Alpine Linux für Raspberry Pi: \url{https://wiki.alpinelinux.org/wiki/Raspberry_Pi}

# Literaturverzeichnis
