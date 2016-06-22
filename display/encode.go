package display

// NOTE: Custom chars are repeated in 8-15;
// use 8 instead of 0 (=> nul-byte) therefore.
const (
	glyphHBar   = 8
	glyphPlay   = 1
	glyphPause  = 2
	glyphHeart  = 3
	glyphCross  = 4
	glyphCheck  = 5
	glyphStop   = 6
	glyphCactus = 7
)

var unicodeToLCDCustom = map[rune]rune{
	// Real custom characters:
	'━': glyphHBar,
	'▶': glyphPlay,
	'⏸': glyphPause,
	'❤': glyphHeart,
	'×': glyphCross,
	'✓': glyphCheck,
	'⏹': glyphStop,
	// Existing characters on the LCD:
	'ψ': glyphCactus,
	'ä': 132,
	'Ä': 142,
	'ü': 129,
	'Ü': 152,
	'ö': 148,
	'Ö': 153,
	'ß': 224,
	'π': 237,
	'৹': 178,
}

func encode(s string) []rune {
	// Iterate by rune:
	encoded := []rune{}

	for _, rn := range s {
		b, ok := unicodeToLCDCustom[rn]
		if !ok {
			if rn > 255 {
				// Multibyte chars would be messed up anyways:
				b = '?'
			} else {
				b = rn
			}
		}

		encoded = append(encoded, b)
	}

	return encoded
}
