all:
	pandoc --filter pandoc-fignos --filter pandoc-tablenos --smart --bibliography eulenfunk-doc.bib --csl ieee.csl eulenfunk-doc.md -B title.tex -H header.tex -N -o eulenfunk-doc.tex -V lang=de-DE --chapters
	latexmk -pdf eulenfunk-doc.tex
	latexmk -c eulenfunk-doc.tex
