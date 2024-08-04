package util

import (
	"strings"
	"time"
)

// GetFromAndToTimes converting delivery times into hours representing 24 hrs format
// 10AM is converted to 10
// 3PM is converted to 15
func GetFromAndToTimes(delivery string) (int, int, string) {
	parts := strings.Split(delivery, " ")
	layout24Hour := "3PM"
	fromTime := parts[1]
	toTime := parts[3]
	from, _ := time.Parse(layout24Hour, fromTime)
	to, _ := time.Parse(layout24Hour, toTime)
	return from.Hour(), to.Hour(), parts[0]
}

// FormatTimeToInteger Convert time to 24Hrs format
func FormatTimeToInteger(input string) int {
	layout24Hour := "3PM"
	result, _ := time.Parse(layout24Hour, input)
	return result.Hour()
}
