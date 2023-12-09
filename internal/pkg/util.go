package pkg

import "strings"

func Check(text string) string {
	var output string
	for _, v := range text {
		if string(v) == "%" {
			output = strings.Replace(text, "%", "", -1)
		} else {
			output += string(v)
		}

	}
	return output
}
