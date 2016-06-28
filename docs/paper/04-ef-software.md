# Software

Dieses Kapitel beleuchtet die einzelnen Software--Komponenten von *Eulenfunk*.
Selbst verfasste Komponenten sind dabei als solche gekennzeichnet.

## Anforderungen

Zu Beginn der Entwicklung wurden einige Anforderungen an die Software
festgelegt:

1. Leichte Bedienbarkeit, obwohl es nur einen Knopf gibt, der nur die
   Aktionen »links drehen«, »rechts drehen« und Drücken zulässt.
2. Die Software sollte vernünftig mit den Hardwareressourcen des *Raspberry Pi*
   umgehen. Dabei sollte die Hardware nicht die ganze Zeit auf volle Last laufen,
   um eine Überhitzung im beengten Gehäuse zu vermeiden.
3. Die Software sollte möglichst ausfallsicher sein. Fällt beispielsweise ein
   Feature durch einen Absturz der Software aus, so sollten andere Teile des Radios
   nach Möglichkeit nicht betroffen sein. Auch sollte der abgestürzte Teil sich
   neu starten können und entsprechende Log--Nachrichten
   hinterlassen. Dienste, die von dem abgestürzten Dienst abhängen, sollten
   sich bei Neustart dessen wieder neu verbinden.
4. Leichte Wartbarkeit und Fehlersuche durch Schreiben von Logs.
5. Einfache Erweiterbarkeit und Integrierbarkeit mit anderen Anwendungen durch
   lose Kopplung von Diensten.

Inwieweit und wodurch diese Anforderungen erfüllt worden sind, wird im Fazit erläutert.

## Abspielsoftware

Für den Betrieb des Internetradios soll der MPD (Music--Player--Daemon)
verwendet werden, da *Eulenfunk* auf einem eigens entwickeltem MPD--Client
basieren soll. Andere Projekte greifen oft auf Abspielsoftware wie den *MOC*
[vgl. @pietraszak2014buch], Seite 189 ff. oder *mplayer* [@exploring] Seite 638
ff. zu. 

Der MPD ist ein unter Unix gern genutzter Daemon zum Verwalten und Abspielen von
Musik und Radiostreams. Er unterstützt dabei eine große Anzahl von Formaten und
kann diese an mehrere Backends wie ALSA, Pulseaudio oder als HTTP--Stream
ausgeben (siehe auch Abbildung \ref{mpd-overview}). Für unseren Einsatzzweck ist
er dabei aus zwei Gründen besonders geeignet: Auch bei sehr großen Sammlungen
mit einer 5-stelligen Anzahl von Liedern läuft er problemlos und mit einem
vergleichsweise niedrigen Speicherverbrauch. Der zweite Grund ist die lose
Kopplung zwischen Abspielsoftware und User--Interface: MPD selbst ist nur ein
Daemon, der mittels eines zeilenbasierten Textprotokolls (dem
MPD--Protokoll[^MPD-PROTO]) steuerbar ist. Um die Bedienoberfläche kümmert sich
dann ein separater MPD--Client, welcher als »Fernbedienung« für den Daemon
fungiert.

[^MPD-PROTO]: Details unter \url{https://www.musicpd.org/doc/protocol}

\begin{figure}[h!]
  \centering
  \includegraphics[width=0.85\textwidth]{images/mpd-overview.png}
  \caption{Übersicht über die Architektur des MPD--Daemons (Quelle: \url{http://mpd.wikia.com/wiki/What_MPD_Is_and_Is_Not})}
  \label{mpd-overview}
\end{figure}

## Softwarearchitektur

Die »Nachbaubarkeit« vieler Bastelprojekte ist häufig durch die Software recht
eingeschränkt, da diese entweder nicht frei verfügbar ist oder zu wenig
generisch ist als dass man die Software leicht auf das Projekt anpassen könnte.
Meist handelt es sich dabei um ein einziges, großes C--Programm oder ein eher
unübersichtliches Python--Skript. Aus diesem Grunde soll die
Software für *Eulenfunk* derart modular aufgebaut sein, dass einzelne Module
problemlos auch auf andere Projekte übertragbar sind und später eine leichte 
Erweiterbarkeit gewährleistet ist. Damit auch andere die Software einsetzen
können wird sie unter die GPL in der Version 3 (vgl. [@GplQuickstart]) gestellt.

Zu diesem Zweck ist die Software in zwei Hauptschichten unterteilt.
Die untere Schicht bilden dabei die *Treiber*, welche die tatsächliche
Ansteuerung der Hardware erledigen. Dabei gibt es für jeden Teil der Hardware
einen eigenen Treiber, im Falle von *Eulenfunk* also ein separates Programm für
die LCD-Ansteuerung, das Setzen der LED--Farbe und dem Auslesen des Drehimpulsgebers.

Die Schicht darüber bilden insgesamt fünf einzelne Dienste, die über eine
Netzwerkschnittstelle angesprochen werden und jeweils eine Funktionalität des
Radios umsetzen. So gibt es beispielsweise einen Dienst, der die Ansteuerung des
LCD--Displays *komfortabel* macht, ein Dienst, der die LEDs passend zur Musik
einfärbt und ein Dienst, der automatisch eine Playlist aus der Musik auf
angesteckten externen Speichermedien erstellt. Die jeweiligen Dienste sprechen
mit den Treibern indem sie Daten auf ``stdin`` schreiben, bzw. Daten von
``stdout`` lesen. Um die Dienste auf neue Projekte zu portieren, ist also nur
eine Anpassung oder Erweiterung der jeweiligen Treiber notwendig.

Der Vorteil liegt dabei klar auf der Hand: Die lose Kopplung der einzelnen
Dienste erleichtert die Fehlersuche ungemein und macht eine leichte
Austauschbarkeit und Übertragbarkeit der Dienste in anderen Projekte möglich.
Stellt man beispielsweise fest, dass der Prozessor des Radios voll ausgelastet
ist, so kann man mit Tools wie ``htop`` einfach herausfinden welcher Dienst dafür
verantwortlich ist.

### Sprachwahl

Die momentane Software ist in den Programmiersprachen *C* und *Go* geschrieben.
Dazu kommt lediglich ein einzelnes Bash--Skript zum Auslesen von Systeminformationen.

Die Ressourcen auf dem *Raspberry Pi* sind natürlich sehr limitiert, weswegen
sehr speicherhungrige Sprachen wie Java oder Ähnliches von vornherein
ausschieden. Obwohl Python nach Meinung des Autors eine schöne und komfortable
Sprache ist und viele gute Bibliotheken für den *Pi* bietet, schied es ebenfalls
aus diesem Grund aus.

Ursprünglich war sogar geplant, alles in *Go* zu schreiben. Leider gibt es nur
wenige Pakete für die GPIO--Ansteuerung und auch keine Bibliothek für
softwareseitige Pulsweitenmodulation. Zwar hätte man diese notfalls auch selbst
mittels ``/sys/class/gpio/*`` implementieren können, doch bietet *Go* leider auch
keine native Möglichkeit mit Interrupts zu arbeiten. Wie später beschrieben ist
dies allerdings für den Treiber nötig, der den Drehimpulsgeber ausliest.

Für *Go* sprechen ansonsten folgende Gründe als Sprache für die höhere Logik:

- **Garbage Collector:** Erleichtert die Entwicklung lang laufender Dienste.
- **Hohe Grundperformanz:** Zwar erreicht diese nicht die Performanz von C, 
  liegt aber zumindest in der selben Größenordnung (vgl. [@pike2009go], S. 37).
- **Weitläufige Standardbibliothek:** Kaum externe Bibliotheken notwendig.
- **Schneller Kompiliervorgang:** Selbst große Anwendungen werden in wenigen 
  Sekunden in eine statische Binärdatei ohne Abhängigkeiten übersetzt.
- **Kross--Kompilierung:** Durch Setzen der ``GOARCH=arm`` Umgebungsvariable kann 
  problemlos auf einen x86-64--Entwicklungsrechner eine passende
  ARM--Binärdatei erzeugt werden.
- **Eingebauter Scheduler:** Parallele und nebenläufige Anwendungen wie
  Netzwerkserver sind sehr einfach zu entwickeln ohne für jede Aufgabe einen neuen
  Thread starten zu müssen.
- Ein Kriterium war natürlich auch, dass die Autoren gute Erfahrung mit der Sprache
  hatten und **neugierig** waren, ob sie auch für solche Bastelprojekte gut einsetzbar
  ist.

``C`` ist hingegen für die Entwicklung der Treiber vor allem aus diesen Gründen
eine gute Wahl:

- Programmierung mit **Interrupts** bequem und nativ möglich.
- Hohe **Performanz** und genaue Kontrolle über den Speicherverbrauch.
- Verfügbarkeit von **wiringPi**.

## Überblick der einzelnen Komponenten

Ein Überblick über die existierenden Dienste liefert Abbildung
\ref{eulenfunk-services}. Die einzelnen Komponenten werden im Folgenden detailliert erläutert.

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.0\textwidth]{images/eulenfunk-services.png}
  \caption{Übersicht über die Softwarelandschaft von \emph{Eulenfunk}. Dienste mit
  einer Netzwerkschnittstelle sind durch den entsprechenden Port gekennzeichnet.}
  \label{eulenfunk-services}
\end{figure}

### Vorhandene Softwarebibliotheken

``wiringPi`` (\url{http://wiringpi.com}): Eine Portierung der Arduino--``Wiring``
Bibliothek von Gordon Henderson auf den *Raspberry Pi*. Sie dient wie ihr
Arduino--Pendant zur leichten Steuerung der verfügbaren Hardware, insbesondere
der GPIO--Pins über ``/dev/mem``. Daneben wird für den LCD--Treiber auch die
mitgelieferte LCD--Bibliothek genutzt. Für den LED--Treiber wird zudem die
softwarebasierte Pulsweitenmodulation genutzt, allerdings in einer leicht
veränderten Form.

``go-mpd`` (\url{https://github.com/fhs/gompd}): Eine einfache MPD--Bibliothek, die
die wichtigsten Kommandos des MPD--Protokolls unterstützt. 

``go-colorful`` (\url{github.com/lucasb-eyer/go-colorful}): Eine Bibliothek um
Farben in verschiedene Farbräume zu konvertieren. Der Dienst, der die LED
passend zur Musik setzt, nutzt diese Bibliothek um RGB--Farbwerte in den
HCL--Farbraum zu übersetzen. Dieser eignet sich besser um saubere Übergänge
zwischen zwei Farben zu berechnen und Farbanpassungen vorzunehmen. 

``cli`` (\url{github.com/urfave/cli}): Eine komfortable und reichhaltige
Bibliothek, um
Kommandozeilenargumente zu parsen. Unterstützt Subkommandos ähnlich wie
``git``, welche dann wiederum eigene Optionen oder weitere Subkommandos
besitzen können. Beide Features wurden extensiv eingesetzt, um alle in *Go*
geschriebenen Dienste in einer Binärdatei mit konsistentem
Kommandozeileninterface zu vereinen.

## Treiber--Software

In Summe gibt es momentan drei unterschiedliche Treiber. Sie finden sich im
``driver/`` Unterverzeichnis[^driver_github] der Software nebst einem passenden
Makefile. Nach dem Kompilieren entstehen drei Binärdateien, welche mit dem
Präfix ``radio-`` beginnen:

- ``radio-led:`` Setzt die Farbe des LED--Panels auf verschiedene Weise.
- ``radio-lcd:`` Liest Befehle von ``stdin`` und setzt das Display entsprechend.
- ``radio-rotary:`` Gibt Änderungen des Drehimpulsgebers auf ``stdout`` aus.

Die genaue Funktionsweise dieser drei Programme wird im Folgenden näher beleuchtet.

[^driver_github]: Siehe auf GitHub: \url{https://github.com/studentkittens/eulenfunk/tree/master/driver}

\label{internal-led-treiber}

### LED--Treiber (``driver/led-driver.c``)

Der LED--Treiber dient zum Setzen eines RGB--Farbwerts. Jeder Kanal hat den
Wertebereich 0 bis 255. Die Hilfe des Programms zeigt die verschiedenen
Aufrufmöglichkeiten:

```html
Usage:
  radio-led on  ....... turn on LED (white)
  radio-led off ....... turn off LED
  radio-led cat ....... read rgb tuples from stdin
  radio-led rgb  r g b  Set LED color to r,g,b
  radio-led hex  #RRGGBB Set LED color from hexstring
  radio-led fade ...... Show a fade for debugging
```

Erklärung benötigt hierbei nur der ``cat``--Modus, bei dem der Treiber
zeilenweise RGB--Farbtripel auf ``stdin`` liest und setzt. Dieser Modus wird benutzt,
um kontinuierlich Farben zu setzen ohne ständig das Treiberprogramm neu zu starten.

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.0\textwidth]{images/eulenfunk-pwm.png}
  \caption{Grafische Darstellung der Pulsweitenmodulation mit zwei Beispielfrequenzen.
    Jeder weiße Block entspricht einem modulierten Wert. Die Prozentzahl darin entspricht dem »Duty--Cycle«
    (dt. Tastgrad). Zusammen mit dem Wertebereich ergibt sich aus ihm der eigentliche Wert.
    (z.B. $255 \times 0.2 = 51$)
  } 
  \label{eulenfunk-pwm}
\end{figure}

Da ein GPIO--Pin prinzipiell nur ein oder ausgeschaltet werden kann, verwenden
wir Pulsweitenmodulation (vgl. dazu [@gay2014experimenting], S 183). Dabei
macht man sich die träge Natur des menschlichen Auges zu Nutze indem man die
LED sehr schnell hintereinander ein-- und ausschaltet. Ist die LED dabei pro
Ein-- und Ausschaltvorgang genauso lange hell wie dunkel, so leuchtet die LED mit etwa
der Hälfte ihrer maximalen Leuchtstärke. Durch Verlängerung/Verkürzung des
eingeschalteten Zustands können so viele verschiedene Helligkeitsstufen
abgebildet werden. Eine beispielhafte Illustration findet sich in Abbildung
\ref{eulenfunk-pwm}

Da bei niedrigen Helligkeitswerten der ausgeschaltete Zustand besonders lange
gehalten wird, kann es dazu kommen, dass ein Flackern entsteht, da
man die ausgeschalteten Phasen als solche wahrnehmen kann. Um dies zu
verhindern, muss eine ausreichend hohe Frequenz gewählt werden. 

Anders als ursprünglich angenommen, mussten wir feststellen, dass die GPIO--Pins
des *Raspberry Pi* (mit Ausnahme von Pin 18 (vgl. [@gay2014experimenting], S.
185)) kein hardwareseitiges PWM unterstützen. Aus diesem Grund mussten wir auf
softwareseitiges PWM zurückgreifen, um Farben mit mindestens 256 Abstufungen zu
erhalten. Nach etwas Ausprobieren befanden wir die ``softPwm``--Bibliothek von
``wiringPi`` für tauglich.

Diese hat allerdings das Problem, dass eine hartkodierte Pulsweite von 100µs
verwendet wird. Für die meisten Anwendungsfälle und den vom Autor empfohlenen
100 Abstufungen ist das auch in Ordnung. Hundert unterschiedliche Zustände
waren nach kurzem Ausprobieren bei einem weichen Farbübergang zu stark
abgestuft, obwohl mit 100Hz kein nennenswertes Flackern sichtbar war:

$$T_{Periode} = 100\mu s\times 100= 10000\mu s = 0.01s$$

$$f = \frac{1}{T_{Periode}} = 100Hz$$

Optimal wären hier 256 unterschiedliche Zustände, um die volle 8--Bit Farbtiefe auszunutzen.
Daher mussten wir die entsprechende C--Datei kopieren (GPL3--lizenziert) und manuell anpassen.
Dabei haben wir die Pulsweite auf 50µs herabgesetzt, was bei einer Spanne von 256 Werten
eine Frequenz von optisch akzeptablen 78Hz ergibt:

$$T_{Periode} = 50\mu s\times 256 = 12800\mu s = 0.0128s$$

$$f = \frac{1}{T_{Periode}} = 78.125Hz$$

Diese Frequenz scheint optisch ausreichend flackerfrei zu sein und scheint die
CPU nicht übermäßig stark zu belasten (rund +3% Last pro Farbkanal).

Es besteht eine Verbindung zu einem früheren Bastelprojekt namens
``catlight``[^catlight] --- einer mehrfarbigen, in einem Gehäuse montierten LED,
die über USB angesprochen werden kann. Genutzt wird diese zur Benachrichtigung
bei neuen E--Mails, Chat--Nachrichten und Ähnlichem. Zu diesem Zwecke wurde
auch bereits damals ein Treiberprogramm entwickelt, welches das selbe
Bedienkonzept wie ``radio-led`` hat. Dies war während der Entwicklung von
*Eulenfunk* nützlich, da es die Entwicklung der Dienste ``ambilight`` und
``lightd`` unabhängig von der Fertigstellung der Radio--Hardware machte.

[^catlight]: Projektseite unter: \url{https://github.com/studentkittens/catlight}


### LCD--Treiber (``driver/lcd-driver.c``)

Der LCD--Treiber setzt Bereiche des LCD--Displays auf einen gegebenen Text.
Beim Start leert er das Display und liest ähnlich wie »``radio-led cat``«
zeilenweise von ``stdin`` und entnimmt diesen Zeilen die Information welcher
Bereich des Displays gesetzt werden soll. Das vom Treiber erwartete
Zeilenformat ist dabei ``LINENO[,OFFSET] TEXT...``, wobei ``LINENO`` die
gewünschte Zeilennummer als Dezimalzahl ist und der optionale ``OFFSET`` der
Index an dem geschrieben werden soll. Dahinter folgt durch ein Leerzeichen
getrennt beliebiger Text. Ist kein ``OFFSET`` gegeben, so wird die ganze Zeile
überschrieben und nötigenfalls mit Leerzeichen aufgefüllt. Ist der Text länger
als die Zeile wird der Text abgeschnitten.

Der Treiber hält eine Matrix mit den aktuell gesetzten Zeichen und kann daher ein erneutes
Zeichnen einer Zelle im Display verhindern, indem er das neue Zeichen mit dem Alten vergleicht.
Unnötige Zeichenvorgänge waren als störende Schlieren auf dem Display wahrnehmbar.

Zudem bietet der Treiber mit dem ``print-charset`` Argument die Möglichkeit, die
auf dem Display verfügbaren Zeichen aufs selbige auszugeben. Dazu stellt er
jeweils 80 Zeichen dar und wartet einige Sekunden bevor die nächsten 80
ausgegeben werden. Hat er alle 256 Zeichen ausgegeben beendet er sich. Optional
kann man auch ein Start- und End--Offset mitgeben, an dem er das Zeichnen
anfangen soll.

Der Treiber unterstützt eine Reihe hardkodierter Spezialzeichen, welche in der
Menüführung und der UI benutzt werden. Das LCD unterstützt dabei 8
verschiedene *Custom Chars*, welche mittels der Codepoints 0-7 und 8-15
(wiederholt) setzbar sind. Momentan sind diese auf folgende Glyphen gesetzt:

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.0\textwidth]{images/symbols.png}
  \caption{Spezielle, selbst gezeichnete Glyphen im Bereich 0-7 und 8-15. Gezeichnet mittels \url{https://omerk.github.io/lcdchargen}}.
  \label{eulenfunk-symbols}
\end{figure}

Die eigentliche Ansteuerung der Pins übernimmt dabei wieder die
``wiringPi``--Bibliothek, beziehungsweise dessen
LCD--Unterbibliothek[^WIRILCD]. Diese ist kompatibel mit dem populären Hitachi
HD44780 und Nachbauten. Das Display wird im 4--Bit Modus angesprochen. Das
heißt, dass nur vier Datenpins benötigt werden (bei uns Pin 25, 24, 23, und 18). 

[^WIRILCD]: Siehe: \url{https://projects.drogon.net/raspberry-pi/wiringpi/lcd-library}

### Drehimpulsgeber Treiber (``driver/rot-driver.c``)

Dieser Treiber kümmert sich um das Einlesen von Werten und Ereignissen vom
Drehimpulsgeber. Nach dem Start schreibt er alle registrierten Ereignisse auf
``stdout``. Dabei nimmt jede Zeile ein neues Ereignis ein. Die jeweilige Zeile
beginnt mit einem einzelnen Buchstaben und einem Leerzeichen, gefolgt von einem
Wert. Der Buchstabe beschreibt die Art des Ereignisses:

- ``v:`` Änderung des **V**alues durch Drehen des Knopfes.
- ``p:`` Der Knopf wurde gedrückt (eine ``1`` folgt) oder losgelassen eine
  (``0`` folgt).
- ``t:`` Der Knopf wurde für eine bestimmte Zeit gedrückt. Die Zeit folgt
  dahinter in Sekundenbruchteilen.

Initial wird zudem die Zeile ``v 0`` herausgeschrieben.

Technisch registriert sich der Treiber auf Änderungen an den GPIO--Pins 12, 13
(Pin A und B vom Drehimpulsgeber) und 14 (Button--Pin) mittels der
``wiringPi``--Funtkion ``wiringPiISR()``. Diese sorgt dafür, dass bei jeder
Wertänderung ein Interrupt aufgerufen wird. 

In der Interrupt--Routine wird der aktuelle Wert der Pins 12 und 13 ausgelesen und mit
dem vorherigen Wert verglichen. Dadurch ist es möglich zu entscheiden, in welche
Richtung der Drehknopf bewegt wurde ohne dass ein Prellen auftritt. Mehr
Informationen zum sogenannten »Grey Code« findet sich unter [@2014projekte], S.
362 und folgenden Seiten.

Da pro Einrastung des Drehknopfs ca. vier Interrupts getriggert werden, wird
auf eine globale Gleitkommazahl der Wert $\frac{1}{4}$ addiert. Beim Herausgeben
des Wertes wird der Wert dann auf den nächsten Integer gerundet.
Auch für jeden Knopfdruck wird ein Interrupt ausgelöst. Da der Knopf prellt,
wird hier mit einem niedrigen Timeout gearbeitet, welcher die Störsignale
filtert.

Da in Interrupts nur reentrante Funktionen aufgerufen werden sollten, werden nur
globale Flags in den Interruptfunktionen gesetzt. In der ``main``--Funktion
läuft eine Schleife mit einem Timeout von 50 Millisekunden, welche diese Werte
abholt, formatiert und auf ``stdout`` schreibt. Der ursprüngliche Code für
diesen Treiber stammt dabei von einem Blogpost[^BLOG]. Der Code wurde etwas
aufgeräumt und um Knopfdrücke sowie Ausgabe auf ``stdout`` erweitert.

[^BLOG]: Ursprünglicher Treiber: \url{http://theatticlight.net/posts/Reading-a-Rotary-Encoder-from-a-Raspberry-Pi}


## Service Software

Die folgenden Dienste implementieren die eigentliche Logik von *Eulenfunk*. 
Alle Dienste finden sich in einer gemeinsamen Binärdatei namens ``eulenfunk``:

```
$ eulenfunk help
NAME:
   eulenfunk - Control the higher level eulenfunk services

USAGE:
   eulenfunk [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHOR(S):
   Waldsoft <sahib@online.de> 

COMMANDS:
     info       Send mpd infos to the display server on the `mpd` window
     ui         Handle window rendering and input control
     automount  Control the automount for usb sticks filled with music
     lightd     Utility server to lock the led and enable nice atomic effects
     display    Display manager and utilites
     ambilight  Control the ambilight feature

GLOBAL OPTIONS:
   --width value   Width of the LCD screen (default: 20) [$LCD_HEIGHT]
   --height value  Height of the LCD screen (default: 4) [$LCD_HEIGHT]
   --help, -h      show help
   --version, -v   print the version
```

Details zu den Optionen der jeweiligen Dienste können mittels ``eulenfunk help
<service>`` angezeigt werden.

### ``displayd`` -- Der Displayserver

#### Einleitung

Der Displayserver ``displayd`` kümmert sich um die Verwaltung der
Display--Inhalte. Er bietet eine höhere Abstraktionsschicht als der
vergleichsweise simple LCD--Treiber. Dabei bietet er die Abstraktion von
*Zeilen*, *Fenstern* und erleichtert dem Programmierer Enkodierungsaufgaben
indem es ein Subset von Unicode unterstützt. Eine Zeile ist dabei ein beliebig
langer UTF8--enkodiertert Text ohne Zeilenumbruch. Die Zeile kann dabei länger
als das Display sein. In diesem Fall wird die Zeile abgeschnitten oder scrollt
je nach Konfiguration mit einer bestimmten Geschwindigkeit durch. Ein Fenster
hingegen ist eine benannte Ansammlung von Zeilen. Auch ein Fenster kann mehr
Zeilen haben als das Display physikalisch bietet. Vom Nutzer kann der
Fensterinhalt dann vertikal verschoben werden. Es können mehrere Fenster
verwaltet werden, aktiv ist dabei aber nur ein Ausgewähltes.

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

#### Architektur

Das Protokoll von ``displayd`` ist ein relativ simpel gehaltenes,
zeilenbasiertes Textprotokoll. Für den Zugriff auf dasselbige wird daher
auch keine UI--Bibliothek benötigt, lediglich einige Netzwerk-- und
Formatierungs--Hilfsfunktionen wurden implementiert[^LINEWRITER].
Basierend auf diesen Primitiven wurden aber auf Clientseite
Funktionalitäten wie Menü--»Widgets« implementiert, welche die grafische Darstellung
mit der Nutzereingabe verquicken.

[^LINEWRITER]: Siehe \url{https://godoc.org/github.com/studentkittens/eulenfunk/display\#LineWriter}

Neben diesen Aufgaben löst ``displayd`` ein architektonisches Problem:
Wenn mehrere Anwendung versuchen auf das Display zu schreiben, käme ohne zentrale
Instanz ein eher unleserliches Resultat dabei heraus. Durch ``displayd`` können
Anwendungen auf ein separates Fenster schreiben, wovon jeweils nur eines aktiv
angezeigt wird. Abbildung \ref{eulenfunk-displayd} zeigt die Architektur in der Übersicht.

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.0\textwidth]{images/eulenfunk-displayd.png}
  \caption{Architekturübersicht von \texttt{displayd} mit Beispielfenstern.}
  \label{eulenfunk-displayd}
\end{figure}

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
```

Dabei kann in die `<Platzhalter>` folgendes eingesetzt werden:

```
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
scrollenden Zeilen bereits mit diesem ereignisbasierten Ansatz auf eine höhere
Aktualisierungsrate kommt als mit den statischen 150ms. Eine Art »VSync«,
welches die Aktualisierungsrate intelligent limitiert wäre hier in Zukunft
wünschenswert.

#### Entwicklung

Da der *Raspberry Pi* nur bedingt als Entwicklungsplattform tauglich ist
(langsamer Compile/Run Zyklus), unterstützt ``displayd`` auch
Debugging--Möglichkeiten. Im Folgenden werden einige Möglichkeiten gezeigt 
mit ``displayd`` zu interagieren, beziehungsweise Programme zu untersuchen,
die ``displayd`` benutzen:

```bash
# Den display server starten; --no-encoding schaltet spezielles LCD encoding
# ab welches auf normalen Terminals zu Artefakten führt. 
$ eulenfunk display server --no-encoding &

# Gebe in kurzen Abständen das "mpd" Fenster aus.
# (In separaten Terminal eingeben!)
$ eulenfunk display --dump --update --window mpd

# Verbinde zu MPD und stelle aktuellen Status auf "mpd" Fenster dar.
# (auch in separaten Terminal eingeben)
$ eulenfunk info
# Auch nach Unterbrechung wird der zuletzt gesetzte Text weiterhin angezeigt:
$ <CTRL-C>
# Änderungen sind auch möglich indem man direkt mit dem Daemon über telnet
# oder netcat spricht. Hier wird die erste Zeile überschrieben, das aktuelle
# Fenster angezeigt und dann die Verbindung geschlossen.
$ telnet localhost 7777
line mpd 0 Erste Zeile geändert!
render                        
(... Ausgabe ...)
close
```

#### Enkodierung

Das LCD unterstützt 8--bit pro Zeichen. Dabei sind die ersten 127 Zeichen
weitestgehend deckungsgleich mit dem ASCII--Standard. Lediglich die Zeichen 0 bis
31 sind durch *Custom Chars* und einige zusätzliche Zeichen belegt. Dies ist
insofern auch sinnvoll, da in diesem Bereich bei ASCII Steuerzeichen definiert
sind, die auf dem LCD schlicht keinen Effekt hätten.

Die Zeichen 128 bis 255 sind vom Hersteller des Displays mit verschiedenen
Symbolen belegt worden, die keinem dem Autor bekannten Encoding entsprechen. Da
auch nach längerer Internetrecherche keine passende Encoding--Tabelle gefunden
werden konnte, wurde (in mühevoller Handarbeit) eine Tabelle erstellt, die
passende Unicode--Glyphen auf das jeweilige Zeichen des Displays abbildet.
Nicht erkannte UTF8--Zeichen werden als ein »?« gerendert anstatt Zeichen die
mehrere Bytes zur Enkodierung (wie » $\mu$  «) als mehrere falsche Glyphen
darzustellen. So wird beispielsweise aus dem scharfen »ß« das Zeichen 223.
Diese Konvertierung wird transparent von ``displayd`` vorgenommen, wodurch es
möglich wird die meisten Musiktitel und Ähnliches annäherend korrekt darzustellen.

Abbildung \ref{eulenfunk-encoding} zeigt das erstellte Mapping zwischen Unicode und LCD--Display.
Folgende Seiten waren bei der Erstellung der Tabelle hilfreich:

* \url{http://www.amp-what.com} (Suche mittels Keyword)
* \url{http://shapecatcher.com} (Suche mittels Skizze)

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.0\textwidth]{images/encoding.png}
  \caption{Unicode--Version der LCD--Glyphen. Der normale ASCII--Bereich (32-127) wurde ausgelassen.}
  \label{eulenfunk-encoding}
\end{figure}


### ``lightd`` -- Der Effektserver

``lightd`` ist ein relativ einfacher Service, dessen Hauptaufgabe die
Verwaltung des Zugriffs auf die LEDs ist. Wollen mehrere Programme die LED
ansteuern, um beispielsweise einen sanften roten und grünen Fade--Effekt zu
realisieren, so würde ohne Synchronisation zwangsläufig eine zittrige Mischung
beider Effekte entstehen.

Ursprünglich war ``light`` als »``lockd``« konzipiert, der den Zugriff auf
verschiedene Ressourcen verwalten kann. Da aber das Display bereits von
``displayd`` synchronisiert wird und durchaus mehrere Programme den Drehknopf
auslesen dürfen, wurde diese etwas generellere Idee wieder verworfen.

Die zweite Hauptaufgabe von ``lightd`` ist die Darstellung von Farbeffekten auf
der LED. Ursprünglich waren diese dafür gedacht, um beispielsweise beim Verlust
der WLAN--Verbindung ein rotes Blinken anzuzeigen. Momentan wird allerdings nur
beim Herunterfahren bzw. Rebooten ein rotes bzw. oranges Blinken angezeigt.
Folgende Effekte sind also momentan als Möglichkeit zur Erweiterung zu
begreifen:

- ``blend:`` Überblendung zwischen zwei Farben.
- ``flash:`` Kurzes Aufblinken einer bestimmten Farbe.
- ``fade:`` Lineare Interpolation von schwarz zu einer bestimmten Farbe und zurück.
- ``fire:`` Kaminfeuerartiges Leuchten.

Wie andere Dienste wird auch ``lightd`` mittels einer Netzwerkschnittstelle kontrolliert.
Die möglichen Kommandos sind dabei wie folgt:

```
!lock      -- Versuche exklusive Zugriffsrechte zu erlangen oder warte bis möglich.
!unlock    -- Gebe exklusive Zugriffsrechte zurück.
!close     -- Schließe die Verbindung.
<effect>   -- Auszuführender Effekt, siehe unten.
```

``<effect>`` darf dabei folgendes sein:

```
{<r>,<g>,<b>}                              -- Einzelne Farbe.
blend{<src-color>|<dst-color>|<duration>}  -- Blend Effekt.
flash{<duration>|<color>|<repeat>}         -- Flash Effekt.
fire{<duration>|<color>|<repeat>}          -- Fire Effekt.
fade{<duration>|<color>|<repeat>}          -- Fade Effekt.
```

Der Platzhalter ``<*-color>`` steht dabei für eine einzelne Farbe.
``<duration>`` ist eine zeitliche Dauer und ``<repeat>`` eine Ganzzahl, die beschreibt wie oft
der Effekt wiederholt wird.

Die Farbe wird, ähnlich wie bei ``displayd``, auf ``stdin`` von ``radio-led`` geschrieben.

### ``ambilightd`` -- Optische Musikuntermalung

Ein »Gimmick« von *Eulenfunk* ist es, dass die LED entsprechend zur momentan
spielenden Musik eingefärbt wird. Hier erklärt sich auch der Name dieses
Dienstes: *Ambilight* (vgl. [@seuntiens2007visual]) bezeichnet eigentlich eine
von Phillips entwickelte Technologie, um an einem Fernseher angebrachte LEDs
passend zum momentanen Bildinhalt einzufärben. Hierher kommt auch die
ursprüngliche Idee, dies auf Musik umzumünzen.

Um aus den momentan spielenden Audiosamples eine Farbe abzuleiten, gibt es
einige Möglichkeiten (vgl. Kapitel \ref{internal-audio-vis}). Eine
große Einschränkung bildet hierbei allerdings die sehr begrenzte Rechenleistung
des *Raspberry Pi*. Daher haben wir uns für eine Variante entschieden, bei der die
Farbwerte vorberechnet werden. Das hat den offensichtlichen Nachteil, dass man
für Radiostreams kein Ambientlicht anzeigen kann. Andererseits möchte man das
bei Nachrichtensendung und Diskussionsrunden vermutlich auch nicht.

Zur Vorberechnung nutzen wir dabei das ``moodbar`` Programm (vgl. das Paper von
Gavin Wood[@wood2005techniques]). Dieses analysiert mit Hilfe des
GStreamer--Frameworks (vgl. [@taymans2013gstreamer]) eine Audiodatei in einem
gebräuchlichen Format und zerlegt diese in 1000 Einzelteile. Sehr oberflächlich
erklärt[^NOTE] wird für jedes dieser Teile ein Farbwert berechnet, wobei
niedrige Frequenzen tendenziell zu roten Farbtönen werden, mittlere Frequenzen
zu Grüntönen und hohe Frequenzen zu blauen Tönen werden. Die so gesammelten
Farbwerte werden dann in einer ``.mood``--Datei gespeichert, welche aus 1000
RGB--Tripeln à 3 Byte (1 Byte pro Farbkanal) bestehen. Ein visualisiertes
Beispiel für eine Moodbar kann man in Abbildung \ref{queen-moodbar} sehen.

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

Um nun aber tatsächlich ein zur Musik passendes Licht anzuzeigen, muss für jedes
Lied in der Musikdatenbank eine passende Moodbar in einer Datenbank
abgespeichert werden. Diese Datenbank ist im Fall von *Eulenfunk* ein simples
Verzeichnis in dem für jedes Lied die entsprechende Moodbar mit dem Pfad
relativ zum Musikverzeichnis[^NOTE2] als Dateinamen abgespeichert wird.

[^NOTE2]: Wobei »/« durch »|« im Dateinamen ersetzt werden.

Die Datenbank kann dabei mit folgendem Befehl angelegt und aktuell gehalten werden:

```bash
$ eulenfunk ambilight --update-mood-db --music-dir /music --mood-dir /var/mood 
```

Die eigentliche Aufgabe von ``ambilightd`` ist es nun den Status von MPD zu
holen, die passende ``.mood``--Datei zu laden und anhand der Liedlänge und der
bereits vergangenen Zeit den aktuell passenden Sample aus den 1000
vorhandenen anzuzeigen. Damit der Übergang zwischen den Samples flüssig ist
wird linear zwischen den einzelnen Farbwerten überblendet. Da die LED eher zu
einer weißen Farbe tendiert, wenn mehrere Kanäle an sind, werden mittlere
Sättigungwerte leicht verstärkt und mittlere Helligkeitswerte etwas abgeschwächt.

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.0\textwidth]{images/eulenfunk-ambilight.png}
  \caption{Konzeptuelle Übersicht über \emph{ambilight} und verwandte Dienste.}
  \label{eulenfunk-ambilight}
\end{figure}

Auch ``ambilight`` ist über Netzwerk über ein zeilenbasiertes Textprotokoll auf
Port 4444 ansprechbar. Es werden die Kommandos ``off`` (Halte ``ambilight``
Ausgabe an), ``on`` (Führe ``ambilight`` Ausgabe fort) und ``state`` (Zeige
aktuellen Status) unterstützt. Eine Übersicht über alle beteiligten Komponenten
von ``ambilight`` findet sich in Abbildung \ref{eulenfunk-ambilight}.

Im Fazit kann man sagen, dass die verwendete Technik durchaus gut für die meisten Lieder funktioniert.
Besonders die Synchronisation ist dabei erstaunlich akkurat und solange nur ein einzelnes Instrument spielt,
wird dem auch eine passende Farbe zugeordnet. Ein Dudelsack erscheint beispielsweise meist grün, während 
ein Kontrabass in mehreren Liedern rot erschien.

Lediglich bei schnellen Tempowechseln (Beispiel: »Prison Song« von »System of a
Down«) sieht man, dass der Farbübergang bereits anfängt bevor man den
zugehörigen Ton hört. Dem könnte im Zukunft abgeholfen werden, indem keine
lineare Interpolation zwischen den Farben genutzt wird, sondern ein
Verfahren, das plötzliche Übergänge eher berücksichtigt.

### ``automount`` -- Playlists von USB Sticks erstellen

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.0\textwidth]{images/eulenfunk-automount.png}
  \caption{Beispielhafte Situation bei Anschluß eines USB Sticks mit dem Label »usb-music«}
  \label{eulenfunk-automount}
\end{figure}

Der ``autmount``--Daemon  sorgt dafür, dass angesteckte USB--Sticks automatisch gemounted werden,
Musikdateien darauf indiziert werden und in einer Playlist mit dem »Label« des Sticks landen.

Unter Linux kümmert sich ``udev`` um die Verwaltung  des
``/dev``--Verzeichnis. Zur Steuerung und Konfiguration erlaubt ``udev`` das
Anlegen von Regeln. Wird ein neues Gerät angesteckt, so geht ``udev`` alle
bekannten und auf das Gerät passende Regeln durch und wendet die darin
definierten Aktionen an.
Bei *Eulenfunk* wird die Regel ``config/udev/11-music-usb-mount.rules`` in das
Verzeichnis ``/etc/udev/rules.d`` kopiert. Die ``11`` im Namen sorgt dafur,
dass die Regel alphabetisch vor den Standardregeln abgearbeitet wird (beginnend
mit ``50``).  Abbildung \ref{eulenfunk-automount} zeigt danach den Ablauf beim
Anstecken eines USB--Sticks mit dem Label »usb-music«. Es wird von der
Regel aus ein Befehl über Netzwerk an ``automount`` gesendet, welcher dann das
Mount und Unmount samt Erstellen der Playlist übernimmt.
Die genaue Erklärung der Einzelheiten wird hier aus Platzgründen ausgelassen.
Weitere Informationen zu ``udev``--Regeln sind online[^UDEV][^UDEV2] zu finden.

[^UDEV]: Writing udev rules: \url{http://www.reactivated.net/writing\_udev\_rules.html}
[^UDEV2]: Arch Linux udev:
\url{https://wiki.archlinux.de/title/Udev\#Unter\_.2Fmedia\_einbinden.3B\_Partitions\_Label\_verwenden\_falls\_vorhanden}

Eine berechtigte Frage ist warum ``automount`` das Mounten/Unmounten des Sticks
übernimmt, wenn diese Aktionen prinzipiell auch direkt von der ``udev``--Regel
getriggert werden können. Der Grund für diese Entscheidung liegt am
*Namespace*--Feature von ``systemd`` (weit verbreitetes Init--System[^SYSTEMD]):
  Dabei können einzelne Prozesse in einer Sandbox laufen, der nur
begrenzten Zugriff auf seine Umgebung hat. Ruft man ``/bin/mount`` direkt aus
der Regel heraus auf, so wird der Mount im *Namespace* von ``udev`` erstellt
und taucht daher nicht im normalen Dateisystem auf.  Sendet man hingegen einen
Befehl an ``automount``, so agiert dieser außerhalb einer Sandbox und kann den
Stick ganz normal als ``root``--Nutzer mounten.

[^SYSTEMD]: \url{https://de.wikipedia.org/wiki/Systemd}

Beim Entfernen des USB--Sticks wird die inverse Operation ausgeführt: Die Playlist wird 
gelöscht (da MPD die Lieder nicht mehr abspielen könnte) und der Mount wird wieder entfernt.

Ähnlich wie die vorherigen Dienste unterstützt ``automount`` einige wenige
Befehle, die es über einen Netzwerk--Socket auf Port 5555 erhält:

* ``mount <device> <label>:`` Mounte ``<device>`` (z.B. ``/dev/sda1``) zu ``<music_dir>/mounts/<label>``
   und erstelle eine Playlist names `<label>` aus den Liedern.
* ``unmount <device> <label>:`` Entferne Playlist und Mountpoint wieder.
* ``close:`` Trenne die Verbindung.
* ``quit:``  Trenne die Verbindung und beende den Daemon.

### ``ui`` -- Menübasierte Bedienoberfläche

Die ``ui`` ist der einzige Dienst ohne Netzwerkschnittstelle. Er kümmert sich
um das Anlegen und Befüllen aller Fenster und um die Steuerung fast aller
anderen Dienste mittels einer menübasierten Oberfläche. Die genaue
Beschaffenheit der Oberfläche wird im nächsten Kapitel (siehe
\ref{internal-owl-chapter}, TODO: FIXREF) im Stile eines Benutzerhandbuches beleuchtet. Daher
wird hier nur ein kurzer Überblick über die Technik dahinter gegeben.

Im momentanen Zustand existieren folgende Fenster:

- ``mpd:`` Zeigt Infos über das momentan laufende Lied oder den aktuellen Radiostream.
- ``clock:`` Zeigt das aktuelle Datum und Uhrzeit.
- ``stats:`` Zeigt Statistiken über die MPD--Datenbank an.
- ``sysinfo:`` Zeigt die Ausgabe des Skripts ``config/scripts/radio-sysinfo.sh`` an.
- ``about:`` Zeigt die Credits an.
- ``menu-main`` Hauptmenü von dem aus alle Funktionen erreichbar sind.
- ``menu-playlists:`` Zeigt alle verfügbaren Playlists.
- ``menu-power:`` Einträge zum Herunterfahren und Rebooten.

Fenster, die mit »``menu-``« beginnen, können durch Benutzung des
Drehimpulsgebers erkundet werden.
Eine Drehung nach rechts verschiebt das Menü nach unten, eine Drehung nach links nach oben.

Wie oben bereits erwähnt wurde ein kleines Client--seitiges »Toolkit« implementiert, welches die
leichte Erstellung von Menüs und die Verknüpfung mit dem Drehimpulsgeber mittels Aktionen möglich macht.
Der prinzipielle Aufbau dieses Toolkits ist in Abbildung \ref{eulenfunk-ui} gezeigt.

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.0\textwidth]{images/eulenfunk-ui.png}
  \caption{Aufbau der Menüstruktur aus Entwicklersicht. Der Entwickler interagiert mit ``MenuManager``;
  dieser kümmert sich um die Ereignisverarbeitung.}
  \label{eulenfunk-ui}
\end{figure}

Der Programmierer muss dabei dem ``MenuManager`` nach dessen Anlegen lediglich eine Liste von Einträgen
mitgeben, die dieser verwalten soll. Es gibt momentan drei verschiedene Arten von Einträgen:

- ``SeparatorEntry:`` Zeigt einen vordefinierten Text. Wird zum Abtrennen verschiedener Sektionen
  innerhalb eines Menüs benutzt, daher der Name.
- ``ClickEntry:`` Zeigt einen Text und ist mit einer Aktion verknüpft. Drückt der Nutzer den Drehknopf 
  während der Fokus auf dem Eintrag ist, so wird die Aktion ausgeführt.
- ``ToggleEntry:`` Wie ``ClickEntry``, hat aber mehrere Zustände, die in eckigen Klammern hinter dem Text
  angezeigt werden. Ein Knopfdruck führt zum Weiterschalten zum nächsten Zustand. Dabei ist jeder Zustand
  mit einer Aktion verknüpft, die beim Umschalten ausgeführt wird.

Diese Einträge kann der Entwickler dann mit beliebigen Aktionen verknüpfen.
Daneben gibt es auch noch drei andere Aktionstypen, die unabhängig vom aktuellen
Eintrag ausgeführt werden:

- ``TimedActions:`` Wird der Knopf für eine bestimmte Zeit gehalten, kann nach
  einer bestimmten Zeit einmalig eine Aktion ausgeführt werden. Die Zeit nach
  der das passiert wird beim Registrieren der Aktion angegeben.
- ``ReleaseActions:`` Wird ausgeführt wenn der Knopf gedrückt wurde, aber
  kein Menüeintrag ausgeführt wurde.
- ``RotateActions:`` Wird ausgeführt sobald der Knopf nach rechts oder links
  gedreht wurde.

Auf Basis dieses minimalen Toolkits wurde dann eine leicht erweiterbare
Menüführung entwickelt. Eine genauere API--Beschreibung kann unter
``godoc.org``[^UI-GODOC] eingesehen werden.

[^UI-GODOC]: \url{https://godoc.org/github.com/studentkittens/eulenfunk/ui}

#### ``info`` -- Anzeige des ``mpd``--Status

``info`` ist der Teil der UI, welcher den Inhalt des ``mpd`` Fensters pflegt
und darstellt. Im Hintergrund steht dabei ein voll funktionsfähiger MPD--Client,
welcher auch auf Zustandsänderungen von außen reagiert. Das heißt: Ändert man
das aktuelle Lied mittels eines anderen MPD--Clients (von einem Handy als
Fernbedienung etwa), so wird die Änderung umgehend an die UI propagiert.

``info`` kann aus Gründen der Fehlersuche auch separat von der UI gestartet werden:

```bash
$ eulenfunk info
```

### ``ympd`` -- MPD im Webbrowser

``ympd``[^YMPD-HOME] ist ein relativ populärer, in C geschriebener MPD--Client, der als
Webserver fungiert und eine Bedienung von MPD im Browser via Websocket möglich
macht. Auf *Eulenfunk* läuft er auf Port 8080 und kann von außen zugegriffen
werden. Abbildung \ref{ympd-screen} zeigt die Weboberfläche.

[^YMPD-HOME]: Homepage von ``ympd``: \url{https://github.com/notandy/ympd}

\begin{figure}[h!]
  \centering
  \includegraphics[width=1.0\textwidth]{images/ympd.png}
  \caption{Screenshot der ympd--Oberfläche mit laufenden Lied aus der Testdatenbank.}
  \label{ympd-screen}
\end{figure}

Ursprünglich war der Einsatz eines, in einer früheren Studienarbeit
entwickelten, Webclients namens »Snøbær« (siehe Abbildung \ref{snobaer-screen})
angedacht. Dieser hat ein paar mehr Features wie das automatische Herunterladen
von Coverart und Liedtexten. Leider lies sich dieser nicht ohne größeren Aufwand
auf dem *Raspberry Pi* kompilieren, weswegen aus Zeitgründen einstweilen auf ``ympd`` umgeschwenkt wurde.

\begin{figure}[h!]
  \centering
  \includegraphics[width=0.4\textwidth]{images/snobaer-small.png}
  \caption{Screenshot von Snøbær's Weboberfläche.}
  \label{snobaer-screen}
\end{figure}

#### Zeroconf

Normalerweise ist die Weboberfläche von *Eulenfunk* unter der Adresse
``http://eulenfunk:8080`` erreichbar.
Falls das wegen mangelnder Namensauflösung aber nicht funktioniert, kann 
man *Zeroconf* (vgl. [@ietf2013zero]) dazu nutzen die IP von *Eulenfunk* herauszufinden:

```bash
$ avahi-browse _mpd._tcp -r | grep 'hostname = \[eulenfunk' -A 2
   hostname = [eulenfunk.local]
   address = [192.168.23.30]
   port = [6600]
```

Das funktioniert, weil der MPD--Server seine Anwesenheit mittels 
Zeroconf--Protokolls mitteilt. Es muss allerdings der Avahi--Daemon sowohl auf
dem anfragenden Rechner als auch auf *Eulenfunk* aktiv sein.

## Einrichtung

In diesem Teilkapitel soll die Einrichtung aller relevanten Komponenten dokumentiert werden.
Dies soll hauptsächlich zur Referenz dienen, um das Radio nachbauen zu können.

### ``mpd`` und ``ympd``

Installation aus den offiziellen Quellen (``mpd``) und aus dem AUR (``ympd``):

```bash
$ pacman -S mpd mpc
$ yaourt -S ympd
$ mkdir -p /var/mpd/playlists
$ touch /var/mpd/mpd.{db,state,log}
$ cp eulenfunk/config/mpd.conf /var/mpd
$ mpd /var/mpd/mpd.conf
$ mpc update -w
```

``mpd`` und ``ympd`` sind die einzigen Dienste, die von Außen (ohne
Authentifizierung) zugreifbar sind. Auch wenn *Eulenfunk* normal in einem
abgesicherten WLAN hängt, wurde für die beiden Dienste jeweils ein eigener
Nutzer mit eingeschränkten Rechten und ``/bin/false`` als Login--Shell angelegt.


### ``systemd`` -- Start und Status von Diensten

``systemd`` ist ein sehr mächtiges Init--System welches zum Starten und Überwachen aller Dienste
in *Eulenfunk* eingesetzt wird. Im Gegensatz zu anderen Init--Systemen werden keine Shell--Skripte
zum Starten genutzt, sondern sogenannte Unit--Files. Diese regeln welche Binaries gestartet werden,
von was diese abhängen und was bei einem Crash passieren soll. Diese Dateien kopiert man in ein 
von ``systemd`` überwachtes Verzeichnis (beispielsweise ``/usr/lib/systemd/system``). Dort kann man
nach einem »``systemctl daemon-reload``« den Dienst starten und für den nächsten Bootvorgang vormerken:

```bash
$ systemctl start my-unit-file    # Jetzt den Dienst starten.
$ systemctl enable my-unit-file   # Beim nächsten Boot starten.
```

In Abbildung \ref{eulenfunk-systemd} ist schematisch der Abhängigkeitsgraph der
einzelnen Dienste von *Eulenfunk* gezeigt. Jeder relevante Dienst hat dabei ein
eigenes ``.unit``--File[^UNIT-FILES].

\begin{figure}[h!]
  \centering
  \includegraphics[width=0.9\textwidth]{images/eulenfunk-systemd.png}
  \caption{Abhängigkeits--Graph beim Start der einzelnen Dienste}
  \label{eulenfunk-systemd}
\end{figure}

[^UNIT-FILES]: Siehe GitHub für eine Auflistung aller Unit--Files: \url{https://github.com/studentkittens/eulenfunk/tree/master/config/systemd}

Sollte ein Dienst abstürzen, weil beispielsweise die Software doch noch Fehler
hat, dann wird der Service von ``systemd`` automatisch neu gestartet, da die
Option ``Restart=on-failure`` in den ``.unit``--Files von *Eulenfunk* gesetzt
ist. Besonders zur Fehlersuche war ``systemd`` bereits sehr hilfreich, da es
möglich, ist die Ausgabe des Dienstes im Nachhinein zu betrachten:

```bash
$ systemctl status radio-ui.service
* radio-ui.service - MPD client that show the current state on the LCD
  Active: active (running) since Mon 2016-06-27 23:45:55 CEST; 2h 0min ago
[...]
Jun 28 01:19:24 eulenfunk eulenfunk[410]: 2016/06/28 01:19:24 Pressed for 1s
[...]
# Alternativ falls mehrere Dienste zeitgleich angezeigt werden sollen:
$ journalctl -u radio-ambilight.service -u radio-ui.service --since today
[...]
```

### ``eulenfunk`` Software

Die Software selbst kann ohne große Abhängigkeiten direkt von *GitHub* installiert werden:

```bash
$ pacman -S wiringpi go git
$ git clone https://github.com/studentkittens/eulenfunk && cd eulenfunk
$ mkdir -p /root/go
$ export GOPATH=/root/go
$ export GOBIN=/root/go/bin
$ go get .
$ make install
```

Das mitgelieferte Makefile installiert alle ``.unit``--Files, Skripte, Treiber--Binärdateien und die
``eulenfunk``--Binärdatei selbst. Die Software sollte dabei auch auf normalen Desktop--Rechnern
kompilierbar sein, ist dort aber mangels passender Treiber nur bedingt ausführbar.

## Fazit

Die Anforderungen können durchaus als erfüllt betrachtet werden. Anforderung
**#3** *(Ausfallsicherheit)*  und **#4** *(Wartbarkeit)* wurden durch den
konsequenten Einsatz von ``systemd`` umgesetzt. Anforderung **#5** *(Lose
Kopplung)* wird durch die Aufteilung der Dienste in einzelne, durch das Netzwerk
getrennte, Prozesse erreicht. Die Erweiterbarkeit sollte dadurch gewährleistet 
sein, dass relativ klar strukturierter und modularer Code verfasst wurde.
Möchte man sich mit dem Code vertraut machen, so hilft die
API--Dokumentation[^EULENFUNK] weiter.  

[^EULENFUNK]: API--Doc Eulenfunk: \url{https://godoc.org/github.com/studentkittens/eulenfunk}

Anforderung **#2** *(Effizienz)* wurde durch die Wahl effizienter
Programmiersprachen und Vermeidung ressourcenhungriger oder ineffizienter
Programmierung vermieden. Die CPU--Last bewegt sich dabei meist zwischen 40--60%
bei ``mp3``--enkodierten Liedern. Bei ``.flac``--Dateien liegt die Last etwa
10--20 Prozentpunkte höher. Die Speicherauslastung ist auch nach mehreren Stunden
Benutzung bei konstanten 50MB.

Auch wurde versucht, die Anzahl von gestarteten Prozess klein zu halten und
nur das Nötigste zu starten:

```html
$ pstree -A
systemd-+-2*[agetty]
        |-avahi-daemon---avahi-daemon
        |-dbus-daemon
        |-dhcpcd
        |-eulenfunk display-+-radio-lcd
        |-eulenfunk lightd-+-radio-led---3*[{radio-led}]
        |-eulenfunk ambilight-+-2*[radio-led---3*[{radio-led}]]
        |-eulenfunk automount---6*[{eulenfunk}]
        |-eulenfunk ui-+-radio-rotary
        |           |-radio-sysinfo.sh
        |-haveged
        |-mpd
        |-sshd---sshd---sshd---bash---pstree
        |-systemd---(sd-pam)
        |-systemd-journal
        |-systemd-logind
        |-systemd-resolve
        |-systemd-timesyn---{sd-resolve) S 1 
        |-systemd-udevd
        |-wpa_supplicant
        `-ympd
```

Inwiefern Anforderung **#1** (*leichte Bedienbarkeit*) gewährleistet ist, wird
das nächste Kapitel zeigen. Die Software darunter versucht die Möglichkeiten
dafür zu schaffen.

### Quelltextumfang

Mit dem Tool ``cloc`` (\url{https://github.com/AlDanial/cloc}) wurde eine
Statistik erstellt, welche den Umfang der Software zeigt. Diese Statistik wurde
bereits von auto--generierten Code und Fremdcode bereinigt:


| **Language**               | **Files** |   **Blank**  |  **Comment** | **Code** |
|----------------------------|-----------|--------------|--------------|----------|
| *Go*                       |    32     |      991     |      550     |     4090 |
| *C*                        |     3     |       99     |       15     |      493 |
| *make*                     |     3     |       13     |        1     |       31 |
| *Bourne Shell*             |     1     |        2     |        1     |       12 |
| *SUM:                      |    39     |     1105     |      567     |     4626 |

### Probleme und Verbesserungsmöglichkeiten

Obwohl die Software für den *Eulenfunk*--Prototypen bisher durchaus gut und stabil
funktioniert, gibt es natürlich noch Verbesserungspotenzial:

- Steckt ein USB--Stick nach einem Reboot noch am Radio, so wird dieser nicht
  automatisch gemounted. Erst nach einem An- und Abstecken desselben ist 
  die zugehörige Playlist wieder abspielbar. Beim Start müssten daher manuell
  ``udev``--Ereignisse getriggert werden. Versuche dies zu erreichen schlugen
  aber leider aus bisher ungeklärten Gründen fehl.
- In seltenen Fällen verschwindet der Cursor in Menüs aus ebenfalls ungeklärten Gründen.
- Wie im Kapitel von ``displayd`` beschrieben, wäre auf lange Dauer ein
  Ereignis-basierter Ansatz im Display--Server wünschenswert, um den
  Stromverbrauch im Ruhezustand zu senken.
- *Eulenfunk* benötigt zum Hochfahren momentan etwa 34 Sekunden. Das ist für die
  meisten Anwendungsfälle vollkommen ausreichend, könnte aber eventuell noch
  weiter optimiert werden. Eine genaue Übersicht darüber, welche Dienste wie lang
  zum Start brauchen, liefert das Tool ``systemd-analyze plot``. 
  Der aktuelle Plot kann online auf GitHub eingesehen werden[^BOOT-PLOT].

[^BOOT-PLOT]: \url{https://github.com/studentkittens/eulenfunk/blob/master/docs/paper/images/boot-plot.svg}


ELCH TODO: Überschriften nochmal anschauen...bspw. bei displayd die
unterüberschrieften wie Einleitung... Entwicklung sind seltsam.
