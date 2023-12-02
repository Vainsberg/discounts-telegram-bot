package pkg

import "strings"

func Check(text string) string {
	var output string
	for _, v := range text {
		vstr := string(v)
		if vstr == "%" {
			output = strings.Replace(text, "%", "", -1)
		} else {
			continue
		}

	}
	return output
}
