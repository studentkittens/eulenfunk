:title: Eulenfunk
:author: Susanne Kiessling, Christopher Pahl, Christoph Piechula
:description: The Hovercraft! tutorial.
:keywords: presentation, backend, impress.js, hovercraft, go-lang
:css: presentation.css

----

.. utility roles

.. role:: underline
    :class: underline

.. role:: blocky
   :class: blocky

:id: first 

Projekt: **Eulenfunk**

.. image:: ../paper/images/title.png
   :width: 45%
   :align: center

----

:blocky:`Diese Praesentation...`

**gibt eine Übersicht über:**

- Unsere Wohngemeinschaft in Augsburg.
- Unser Open-Source-Project namens Eulenfunk.
- Den Weg vom Schrott zum Radio.
- Überblick über Vorbereitungen, Entwicklung und Testphase
- Probleme während der Umsetzung.

**...und soll Lust auf mehr machen!**

----

.. note::

    * Leben in einer WG zusammen
    * Alte Musikanlage erklaeren
    * Herausforderung für nicht-technischen Informatiker
    * Anforderungen: Nutzung vorhandener Hardware
    * Preis/Leistung
    * Moodbar (faszniniert von LEDs)
    * Erweiterbar

:blocky:`Dieses Eulenfunk...`

**...ist entstanden weil:**

- Aktuelle Musik-Situation... war kompliziert.
- Wir etwas für den täglichen Gebrauch bauen wollten.
- Alte Hardware vorhanden war.

**...ist anders weil:**

- Es kein Nachbau ist.
- Eigene Anforderungen umgesetzt wurden.
- Diamanten nur unter Druck entstehen.



-----

:blocky:`Bilder sagen manchmal mehr als...`

.. image:: images/impressions/impressions.png
   :width: 100%

-----

.. note:: 

    * Geschichte vom alten Radio (es war einmal)

:blocky:`Altes Analogradio:`

.. image:: images/impressions/Altes_AEG_MR_4104_pol.png
   :width: 100%

-----

.. note:: 

    * Komponenten kurz ansprechen
    * Preis/Leistungsverhaeltnis

:blocky:`Some assembly required:`

.. image:: images/impressions/Komponenten_pol.png
   :width: 120%

-----

.. note:: 

    * LCD Spannung und Anschluesse
    * Arduino Breadboard/Baukasten
    * Langsames Herantasten an Hardware

:blocky:`Testaufbau der Basiskomponenten`

.. image:: images/impressions/LCD_Testaufbau_pol.png
   :width: 110%

-----

.. note::

    * Testaufbau Transistorschaltung.
    * Externes 5V Netzteil
    * PNP/NPN Problem

:blocky:`Testaufbau zur Ansteuerung der LEDs`

.. image:: images/impressions/Transistorschaltung_Testaufbau_pol.png
   :width: 100%

-----

.. note::

    * Von ira-kunststoffe und lackieren erzaehlen
    * Stundelange Diskussion welches Logo
    * Eigentlich ein Triforce...
    * ...aber jetzt ist es Neon Demon.

:blocky:`Lackierung`

.. image:: images/impressions/front_color_pol.png
   :width: 100%

-----

.. note::

    * Drehimpulsgeber
    * Befestigung von Holzstuecken von Hand mit Messer zugeschnitten
    * Viel Heiskleber
    * Audioverstaerker
    * Chirugische Praezision
    * Triforce-Dimming Blende mit kuechenmesser zugeschnitten

:blocky:`Vorderseite von hinten`

.. image:: images/impressions/Front_mit_LCD_und_Switches_pol.png
   :width: 100%

-----

:blocky:`Chirugische Praezision...`

.. image:: images/impressions/transled1_pol.png
   :width: 100%

-----

.. note::

    * Eine Woche kuechenchaos...
    * Viel loetarbeiten und trial & error
    * Parallel beginn der Softwareentwicklung

:blocky:`Informatiker bei der Arbeit...`

.. image:: images/impressions/Hardware-Engineering_pol.png
   :width: 100%

-----

:blocky:`Finaler Prototyp #1`

.. note::

    * Hoffentlich koennt ihr das lied nicht erkennen :)

.. image:: images/impressions/Eulenfunk_Prototyp_pol.png
   :width: 100%

-----

:blocky:`Finaler Prototyp #2`

.. note::

    * Hoffentlich koennt ihr das lied nicht erkennen :)

.. image:: images/impressions/Eulenfunk_Anschlüße_pol.png
   :width: 100%

-----

:blocky:`Die Hardware`

.. image:: ../paper/images/uebersicht.png
   :width: 90%
   :align: right

-----

:blocky:`Moodbar`

.. note:: 

    * Bilderraetsel (Jeopardy)
    * Moodbar utility: 1000 samples pro file.
    * Vorberechnet, da ärmliche Performance.
    * Frequenzen.

.. image:: images/we-will-rock-you.png
   :width: 110%
   :align: left 

|
|
|

.. image:: images/freddy.png
   :width: 100%
   :align: left 

-----

:blocky:`Die Software`

.. note::

    * Ich bin Schuld an der Software.
    * Hat dann doch mehr als 2 Wochen gedauert 
    * Sehr modular
    * Playlisten von USB--Sticks

.. image:: ../paper/images/eulenfunk-services.png
   :width: 100%
   :align: right

-----

.. note::

    * Europaeische Variante des LCD.
    * Keine Encoding Tabelle gefunden, daher Abbildung
    * Mapping von unicode 
    * in muehevoller Handarbeit erstellt

:blocky:`Enkodierung`

.. image:: ../paper/images/encoding.png
   :width: 100%
   :align: right

-----

:blocky:`Bedienung`

.. note:: 

    0. Minimales Bedienkonzept ähnlich Autoradio.

.. image:: ../paper/images/hauptmenue.png
   :width: 100%
   :align: center

-----

:blocky:`Demonstration`

.. note:: 

    1. Vorher einschalten und auf `If you have Ghosts gehen`
    2. Pause und Schuhu spielen.

.. image:: images/effekt.png
   :width: 100%
   :align: center

-----

:blocky:`Probleme...`

.. note:: 

    - Probleme mit Netzteilen erlaeutern
    - Backfeed problem erklaeren
    - schwache hardware fuer visuelle Effekte
    - Go war gut, trotz fehlender Libraries.
    - Bedienkonzept gut, trotz Skepsis von Elch.

- Nur wenige Hardware Pins mit PWMs.
- Strom/Spannungsversorgungsprobleme.
- Prinzipiell schwache Hardware.

:blocky:`...und was gut ging:`

- Go ist nett.
- Kein Hausbrand ausgelöst.
- Bedienungskonzept passt.

|

.. image:: images/gopher.png
   :align: right

-----

:blocky:`Schuhu (Ende)`

.. note::

    - Sehr eng im Gehäuse, Hitze.
    - Softwareplus: Wettervorhersage, automount stabil, bootzeit.
    - Design verschönern, labels für drehknöpfe, logo
    * in memory alpine

**Prototyp ist noch verbesserungswürdig:**

- Netzteil/USB-Hub.
- Entstörung div. Komponenten.
- Komponenten »besser« im Gehäuse unterbringen.
- Design verschönern? Logo? 
- Softwareoptimierungen.
- Alpine Linux?

**Wir haben:**

- Spaß gehabt.
- Alte Hardware nutzbar gemacht.
- Viel gelernt (Pi-Probleme, Hardware...).
- Ein neues Küchenradio.

.. note::

    none

-------

:id: fin

.. note::

   FRAGEN?

:blocky:`Lust bekommen...`

**Dokumentation & Quellen:**

    https://github.com/studentkittens/eulenfunk

|
|

**(Noch Fragen?)**
