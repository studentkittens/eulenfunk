package util

// Center pads `text` to with by placing it in the middle and
// padding the rest of the string with the `padWith` rune.
func Center(text string, width int, padWith rune) string {
	s := ""

	i := 0
	half := width/2 - len(text)/2
	for ; i < half; i++ {
		s += string(padWith)
	}

	s += text
	i += len(text)

	for ; i < width; i++ {
		s += string(padWith)
	}

	return s
}
