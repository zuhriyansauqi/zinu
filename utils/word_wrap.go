package utils

import "strings"

// WordWrap Word wrap string with given line width
func WordWrap(text string, lineWidth int) (wrapped string) {
	words := strings.Fields(text)
	if len(words) == 0 {
		return
	}

	wrapped = words[0]
	spaceLeft := lineWidth - len(wrapped)
	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			wrapped += "\n" + word
			spaceLeft = lineWidth - len(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + len(word)
		}
	}
	return
}

// Tokenize Tokenize string with given token
func Tokenize(text string, token string) (result []string) {
	result = strings.Split(text, token)
	return
}
