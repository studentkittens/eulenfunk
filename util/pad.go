package util

func Center(text string, width int, padWith rune) string {
	s := ""

	i := 0
	for ; i < width/2-len(text)/2; i++ {
		s += string(padWith)
	}

	s += text
	i += len(text)

	for ; i < len(text); i++ {
		s += string(padWith)
	}

	return s
}
