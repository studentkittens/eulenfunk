# Bedienkonzept/Menüsteuerung

\label{internal-owl-chapter}

## Anforderungen

Die Eingabe--Peripherie soll möglichst einfach gehalten werden, um eine *schöne*
Produkt--Optik zu gewährleisten, dabei sollen folgende Anforderungen erfüllt
werden:

* Minimale sowie ansprechende Bedienelemente
* Funktionales, zweckgebundenes *Design*
* *Retro--Look* wünschenswert

Das *Design* soll im Grunde *minimalistisch*  gehalten werden, das heißt, es
sollen aufgrund der Übersichtlichkeit nur so wenige »Bedienelemente« wie nötig
angebracht werden.


## Bedienelemente


\begin{figure}[h!]
  \centering
  \includegraphics[width=0.9\textwidth]{images/front_usage.png}
  \caption{Frontansicht von \emph{Eulenfunk}}
  \label{front_usage}
\end{figure}

Bei der Bedienung von *Eulenfunk* bestand die Herausforderung vor allem darin,
trotz minimaler Anzahl an Bedienelementen eine gute 
Bedienbarkeit (Usability) zu erreichen.

Zur Verfügung stehen primär ein Lautstärkeregler und ein Drehimpulsgeber (Abb. \ref{front_usage}).
Auf der Rückseite des Radios ermöglichen zusätzlich zwei Potentiometer die Regelung
der Beleuchtung und des Kontrastes vom Display. Mit einem Kippschalter kann
die Audioausgabe zwischen internen und externen Lautsprechern geschaltet werden.

Wie allgemein üblich, erhöht ein »nach rechts Drehen« des Lautstärkereglers 
die Lautstärke und »ein nach links Drehen« verringert die Lautstärke.

Der Drehimpulsgeber (im Weiteren Master-Regler genannt) ist für die gesamte Navigation innerhalb des Menüs und das Ausführen von
Aktionen zuständig. Drehen bewirkt grundsätzlich  ein *Vor* oder *Zurück*.
Drücken bewirkt das kontextbezogene Ausführen einer Aktion. Näheres wird bei 
den jeweiligen Menüansichten beschrieben.


## Menüinhalt

\begin{figure}[h!]
  \centering
  \includegraphics[width=0.9\textwidth]{images/hauptmenue.png}
  \caption{Ansicht des Hauptmenüs}
  \label{hauptmenue}
\end{figure}

Das Hauptmenü enthält drei Kategorien: Modes, Options und System (siehe Abb. \ref{hauptmenue}).
Ein Drehen des Master-Regler navigiert zwischen den einzelnen Einträgen.
Ein Drücken des Master-Regler öffnet den gewählten Eintrag als neue Ansicht
(Modes, System) oder setzt direkt andere Werte (Options).

\newpage

### Now Playing

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.2\textwidth]{images/currentsong.png}
  \caption{Ansicht des Menüpunkts \emph{Now Playing} für Musikstücke und Radiosender}
  \label{currentsong}
  \end{figure}

Der Menüpunkt *Now Playing* zeigt Informationen des aktuell gespielten Mediums
an. Es gibt eine geringfügig abweichende Ansicht für Musikstücke und für
Radiosender. Wie in Abbildung \ref{currentsong} zu sehen ist, wird bei einem
Musikstück der Interpret, der Titel und das Album angezeigt. Sollte die Breite
des Displays für die Länge dieser Werte nicht ausreichend sein, »scrollt« die
Anzeige der Zeichen von rechts nach links durch. Die letzte Displayzeile zeigt
an, ob Musik spielt, oder pausiert ist. Außerdem wird die Bitrate und die
Spieldauer angezeigt. Die Ansicht für Radiosender enthält den Interpreten und
den Titel in Zeile eins, gefolgt vom Radiosender in Zeile zwei.

In beiden Ansichten führt ein Drücken des Master-Reglers zum Pausieren bzw. Abspielen
der Musik. Einen Titel bzw. Radiosender vor oder zurück wechseln ist durch Drehen des Reglers möglich.



\newpage 

### Playlists

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.2\textwidth]{images/playlists.png}
  \caption{Ansicht des Menüpunkts \emph{Playlists}}
  \label{playlists}
\end{figure}

Der Menüpunkt *Playlists* (siehe Abb. \ref{playlists}) zeigt eine Übersicht der
aktuell verfügbaren Medien. Mögliche Listen sind Radiosender und Musikstücke.
Wird ein USB-Gerät an das Radio angeschlossen, erscheint er hier als Playlist.

Ein Drücken des Master-Reglers führt zur Auswahl einer Playlist. Anschließend
wechselt die Anzeige auf den Modus *Current Song*, d.h. die Informationen zum
aktuell aus der Playlist abgespielten Mediums werden angezeigt. Ein Drehen des
Master-Reglers navigiert innerhalb des Playlist--Ansicht. In der Ansicht
*Playlists* gibt es zusätzlich die Funktion *(Exit)*, die zum Hauptmenü führt.



### Clock

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.2\textwidth]{images/clock.png}
  \caption{Ansicht des Menüpunkts \emph{Clock}}
  \label{clock}
  \end{figure}

  Der Menüpunkt *Clock* (siehe Abb. \ref{clock}) zeigt die aktuelle Uhrzeit, sowie das Datum an. 
  Ein Drücken des Master-Reglers führt zurück zum Hauptmenü. 



### Statistics

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.2\textwidth]{images/statistics.png}
  \caption{Ansicht des Menüpunkts \emph{Statistics}}
  \label{statistic}
  \end{figure}

Der Menüpunkt *Statistics* (siehe Abb. \ref{statistic}) zeigt an, wieviele
Interpreten, Alben und Songs aktuell in der Musikdatenbank insgesamt zu finden
sind. Zusätzlich wird die Gesamtspieldauer angezeigt.


Ein Drücken des Master-Reglers führt zurück zum Hauptmenü.

### Options

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.2\textwidth]{images/options.png}
  \caption{Ansicht der Kategorie \emph{Options}}
  \label{options}
  \end{figure}

  *Options* ist eine Kategorie, die keine Unteransichten hat (siehe Abb. \ref{options}).
Optionen werden direkt durch ein Drücken des Master-Reglers geändert.

- **Party!:** LEDs können aktiviert oder deaktiviert werden (Aus/Ein).
- **Output:** Umschalten zwischen stereo- und mono-Ausgabe (stereo, mono).
- **Playback:** Wechsel zwischen Abspiel-, Pause-, und Stopmodus (Play, Pause, Stop).
- **Random:** Zufällige Wahl eines Musikstücks aus der aktuell gewählten Playlist (Aus/Ein).


### Powermenu

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.2\textwidth]{images/powermenu.png}
  \caption{Ansicht des Menüpunkts \emph{Powermenu}}
  \label{powermenu}
\end{figure}

Der Menüpunkt *Powermenu* (siehe Abb. \ref{powermenu}) enthält die Funktionen
*Poweroff* für das Herunterfahren des Systems und *Reboot* für einen Neustart des
Systems. Mit *(Exit)* gelangt man zurück ins Hauptmenü.


### Systeminfo

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.2\textwidth]{images/systeminfo.png}
  \caption{Ansicht des Menüpunkts \emph{Systeminfo}}
  \label{systeminfo}
  \end{figure}

  Der Menüpunkt *Systeminfo* (siehe Abb. \ref{systeminfo}) zeigt folgende Informationen zum System an:

* CPU-Auslastung in Prozent.
* Speicher-Auslastung in Prozent.
* IP-Adresse, Systemlaufzeit,
* WLAN-Empfangsstärke in Prozent.
* CPU-Temperatur in Celsius.

Ein Drücken des Master-Reglers führt zurück zum Hauptmenü.


### About

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.2\textwidth]{images/about.png}
  \caption{Ansicht des Menüpunkts \emph{About}}
  \label{about}
\end{figure}

Der Menüpunkt *About* (siehe Abb. \ref{about}) zeigt Informationen zu den Entwicklern des
Internetradios an.


## Shortcuts

Ein Direktzugriff auf das *Powermenu* (siehe [6.2.7](#powermenu)) ist durch
Drücken des Master-Reglers von drei Sekunden möglich --- in allen Anzeigemodi.
Das führt zu mehr Komfortabilität in der Bedienung. Es kann direkt ins
*Powermenu * gewechselt werden, ohne die aktuelle Anzeige verlassen zu müssen,
um manuell zum *Powermenu* zu navigieren.

In der Ansicht *Now Playing* führt ein längeres Drücken (mindestens 600ms) des Master-Reglers 
ins Hauptmenü. Eine Abweichung zu der Funktionsweise in den anderen Ansichten
war notwendig, weil auch der Wechsel zwischen Pause/Play möglich sein sollte.
Für eine möglichst intuitive Bedienung bot es sich an, dies durch ein einfaches
Drücken des Master-Reglers zu realisieren. Für den Wechsel ins Hauptmenü 
wurde deshalb ein längerer Drück--Zeitraum gewählt.


## Fazit

Die Möglichkeiten für die Darstellung der Inhalte waren aufgrund des 4x20 Displays
verständlicherweise begrenzt. Auch die Bedienung des Menüs über einen einzigen
Drehimpulsgeber setzt klare Grenzen für die Steuerung.

Trotz dieser minimalistischen Umstände, ist es gelungen, eine weitgehend konsistente
Menüführung umzusetzen. Nach einer kurzen Eingewöhnungsphase sollte das Bedienkonzept
gut funktionieren.
