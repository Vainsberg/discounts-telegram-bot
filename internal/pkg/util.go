package pkg

import "strings"

func Check(text string) string {
	text = strings.Replace(text, " ", "%20", -1)
	return text
}

func CalculatePercentageDifference(countone, counttwo float64) float64 {
	return (1 - counttwo/countone) * 100
}
