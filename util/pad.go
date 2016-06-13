package util

func Center(text string, width int) string {
	s := ""

	i := 0
	for ; i < width/2-len(text)/2; i++ {
		s += " "
	}

	s += text
	i += len(text)

	for ; i < len(text); i++ {
		s += " "
	}

	return s
}
