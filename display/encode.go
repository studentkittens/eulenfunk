package display

var (
	specialASCII = []rune{
		'±', '≅', '⎲', '/', '⎛', '⎩', '⎞', '⎭', // 16 - 23
		'⎧', '⎫', '≈', '⌠', '=', '~', '²', '³', // 24 - 31
	}

	// Non-standard part beyond ASCII (painfully typesetted by looking at the LCD):
	nonASCII = []rune{
		'Δ', 'Ç', 'ü', 'é', 'å', 'ä', 'à', 'ç', 'ė', 'ë', 'è', 'ï', 'ì', 'Ä', 'Å', 'É', // 128 - 143
		'æ', 'Æ', 'ô', 'ö', 'ò', 'û', 'ù', 'ÿ', 'Ö', 'ü', 'ñ', 'Ñ', 'ā', 'ō', '¿', 'á', // 144 - 159
		'í', 'ó', 'ú', 'ȼ', '£', '¥', '₽', '¢', 'ĩ', 'Ã', 'ã', 'Õ', 'õ', 'Ø', 'ø', '˙', // 160 - 175
		'¨', '°', '`', '՚', '½', '¼', '×', '÷', '≤', '≥', '«', '»', '≠', '√', '⎺', '⌠', // 176 - 191
		'⌡', '∞', '◸', '↵', '↑', '↓', '→', '←', '⎡', '⎤', '⎣', '⎦', '▪', '®', '©', '™', // 192 - 207
		'✝', '§', '¶', '⎴', '⊿', 'Ɵ', 'Λ', '𝚵', 'Π', '∑', 'Ⲧ', 'Φ', 'Ψ', 'Ω', 'α', 'ß', // 208 - 223
		'ɣ', 'δ', 'ε', 'ζ', 'η', 'ɵ', 'ι', 'κ', 'λ', 'μ', 'ν', 'ξ', 'π', 'ρ', 'σ', 'τ', // 224 - 239
		'ʊ', 'φ', 'ψ', 'ω', '▾', '▸', '◂', 'R', '⥒', 'F', '⥓', '▯', '━', 'S', 'P', ' ', // 240 - 255
	}
	// Custom chars of eulenfunk; 0-7 is the same as 8-15.
	customChars = []rune{
		'━', '▶', '⏸', '❤', '×', '✓', '⏹', 'ψ',
		'━', '▶', '⏸', '❤', '×', '✓', '⏹', 'ψ',
	}
)

// Mapping from utf8 characters to LCD codepoint.
// Gets populated in init()
var utf8ToLCD = map[rune]rune{}

func init() {
	for idx, rn := range customChars {
		utf8ToLCD[rn] = rune(idx)
	}

	for idx, rn := range specialASCII {
		utf8ToLCD[rn] = rune(16 + idx)
	}

	for idx, rn := range nonASCII {
		utf8ToLCD[rn] = rune(127 + idx)
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
