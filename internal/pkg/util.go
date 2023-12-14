package pkg

import "strings"

func Check(text string) string {
	text = strings.Replace(text, " ", "", -1)
	return text
}
