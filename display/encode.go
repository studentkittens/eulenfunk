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

// The
var nonAscii = []rune{
	'Δ', 'Ç', 'ü', 'é', 'å', 'ä', 'à', 'ç', 'ė', 'ë', 'è', 'ï', 'ì', 'Ä', 'Å', 'É', // 128 - 143
	'æ', 'Æ', 'ô', 'ö', 'ò', 'û', 'ù', 'ÿ', 'Ö', 'ü', 'ñ', 'Ñ', 'ā', 'ō', '¿', 'á', // 144 - 159
	'í', 'ó', 'ú', 'ȼ', '£', '¥', '₽', '¢', 'ĩ', 'Ã', 'ã', 'Õ', 'õ', 'Ø', 'ø', '˙', // 160 - 175
	'¨', '°', '`', '՚', '½', '¼', '×', '÷', '≤', '≥', '«', '»', '≠', '√', '⎺', '⌠', // 176 - 191
	'⌡', '∞', '◸', '↵', '↑', '↓', '→', '←', '⎡', '⎤', '⎣', '⎦', '▪', '®', '©', '™', // 192 - 207
	'✝', '§', '¶', '⎴', '⊿', 'Ɵ', 'Λ', '𝚵', 'Π', '∑', 'Ⲧ', 'Φ', 'Ψ', 'Ω', 'α', 'ß', // 208 - 223
	'ɣ', 'δ', 'ε', 'ζ', 'η', 'ɵ', 'ι', 'κ', 'λ', 'μ', 'ν', 'ξ', 'π', 'ρ', 'σ', 'τ', // 224 - 239
	'ʊ', 'φ', 'ψ', 'ω', '▾', '▸', '◂', 'R', '⥒', 'F', '⥓', '▯', '━', 'S', 'P', ' ', // 240 - 255
}

var utf8ToLCD = map[rune]rune{
	// Real custom characters:
	'━': glyphHBar,
	'▶': glyphPlay,
	'⏸': glyphPause,
	'❤': glyphHeart,
	'×': glyphCross,
	'✓': glyphCheck,
	'⏹': glyphStop,
	'ψ': glyphCactus,
	// Existing characters on the LCD:
	// 'ä': 132,
	// 'Ä': 142,
	// 'ü': 129,
	// 'Ü': 152,
	// 'ö': 148,
	// 'Ö': 153,
	// 'ß': 224,
	// 'π': 237,
	// '৹': 178,
}

func init() {
	for idx, rn := range nonAscii {
		utf8ToLCD[rn] = rune(idx + 127)
	}
}

func encode(s string) []rune {
	// Iterate by rune:
	encoded := []rune{}

	for _, rn := range s {
		b, ok := utf8ToLCD[rn]
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
