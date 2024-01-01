package pkg

import "strings"

func ReplaceSpaceUrl(text string) string {
	return strings.Replace(text, " ", "%20", -1)
}

func CalculatePercentageDifference(countone, counttwo float64) float64 {
	return (1 - counttwo/countone) * 100
}
