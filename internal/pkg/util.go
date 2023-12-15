package pkg

import "strings"

func Check(text string) string {
	text = strings.Replace(text, " ", "%20", -1)
	return text
}
