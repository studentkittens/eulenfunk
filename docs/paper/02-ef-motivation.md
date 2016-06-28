# Motivation

## Private Situation 

Die Autoren dieses Projekts leben in einer Wohngemeinschaft zusammen. Die Küche
ist der Ort an welchem gemeinsam gekocht, gespeist und diskutiert wird. Für eine
angenehme Atmosphäre und als Nachrichten--Quelle sorgte in der Küche früher ein
Analog--Radio der Firma *AEG*, welches aufgrund der schlechten Empfangsqualität
durch eine Kombination aus »alter Stereoanlage«, »altem Raspberry Pi« und einem
»alten Thinkpad X61t« ersetzt wurde. In dieser Kombination fungierte die
Stereoanlage als Soundausgabe--Komponente, auf dem *Raspberry Pi* lief der
Linux--basierte Player Volumio[^VOL], welcher mit dem Touchscreen des *Thinkpad
x61t* über eine Weboberfläche gesteuert wurde. Diese Kombination hat zwar
funktioniert, jedoch war sie alles andere als benutzerfreundlich, da zuerst die
Stereoanlage und der Laptop  eingeschaltet werden mussten und eine
WLAN--Verbindung zum *Raspberry Pi*--Player hergestellt werden musste. Diese
Situation weckte den Wunsch nach einer komfortableren Lösung, beispielsweise ein
Internetradio auf Basis des *Raspberry Pi*.

[^VOL]: Volumio: \url{https://volumio.org/}


## Kommerzielle Produkte

Kommerzielle Anbieter von Internetradios gibt es wie Sand am Meer. Die
Preisspanne liegt hier zwischen \EUR{30} und mehreren hundert Euro. Der
Funktionsumfang sowie die Wiedergabequalität ist hier von Hersteller zu Hersteller
und zwischen den verschiedenen Preisklassen sehr unterschiedlich. Einen
aktuellen Überblick aus dem Jahr 2016 über getestete Modelle gibt es
beispielsweise online unter *bestendrei.de*[^TEST].

Das Problem bei den kommerziellen Anbietern ist, dass man hier jeweils an die
vorgegebenen Funktionalitäten des Herstellers gebunden ist. Bei einem
Do--It--Yourself--Projekt auf Basis von Freier Software beziehungsweise eines freien
Hardwaredesigns, hat man die Möglichkeit alle gewünschten Funktionalitäten ---
auch Features die von keinem kommerziellen Anbieter unterstützt werden --- zu
integrieren. Beispiele für Funktionalitäten, welche bei kommerziellen Produkten
nur schwer beziehungsweise vereinzelt zu finden sind:

* Unterstützung bestimmter WLAN--Authentifizierungsstandards
* Einhängen von benutzerdefinierten Dateifreigaben wie *Samba*, *NFS*, *SSHFS*
* Unterstützung verschiedener *lossy* und *lossless* Formate *OGG VORBIS*, *FLAC*, u.a.
* Integration verschiedener Dienste wie beispielsweise *Spotify* 
* Benutzerdefinierte Anzeigemöglichkeiten (Uhrzeit, Wetter, et cetera.)

[^TEST]:Test von Internetradios: \url{http://www.bestendrei.de/elektronik/internetradio/}
