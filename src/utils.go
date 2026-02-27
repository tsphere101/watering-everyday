package main

import "time"

func splitMessage(msg string, maxLen int) []string {
	if len(msg) <= maxLen {
		return []string{msg}
	}

	var chunks []string
	runes := []rune(msg)

	for i := 0; i < len(runes); i += maxLen {
		end := i + maxLen
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
	}

	return chunks
}

func GetCurrentDay(startDateStr string) (int, error) {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		location = time.FixedZone("Bangkok", 7*60*60)
	}

	startDate, err := time.ParseInLocation("2006-01-02", startDateStr, location)
	if err != nil {
		return 0, err
	}

	today := time.Now().In(location)
	daysSinceStart := int(today.Sub(startDate).Hours() / 24)
	return daysSinceStart + 1, nil
}

func GetCurrentDayWithTime(startDateStr string, t time.Time) (int, error) {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		location = time.FixedZone("Bangkok", 7*60*60)
	}

	startDate, err := time.ParseInLocation("2006-01-02", startDateStr, location)
	if err != nil {
		return 0, err
	}

	daysSinceStart := int(t.Sub(startDate).Hours() / 24)
	return daysSinceStart + 1, nil
}
