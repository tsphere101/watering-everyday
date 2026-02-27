package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

const (
	prompt               = "เขียนประโยคเชิญชวนเพื่อนชื่อ เชอร์รี่ และ นวว มารดน้ำต้นไม้ 10 แบบ แตกต่างกัน ขอแบบสั้น ๆ น่ารัก ๆ (รดน้ำต้นไม้ในเกมทุกวันเพื่อไม่ให้ต้นไม้เสียชีวิต)"
	discordMaxLength     = 1900
	geminiTimeoutSeconds = 60
	gardenStartDate      = "2025-08-13"
)

var (
	discordUsername   = "รดน้ำ"
	discordAvatarURL  = "https://raw.githubusercontent.com/tsphere101/watering-everyday/refs/heads/main/watering-avatar.jpg"
	discordMentionIDs = []string{"650661678316388372", "739506315784749097"}
	fallbackMessage   = "รดน้ำกันจ้า"
)

func main() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	discordWebhook := os.Getenv("DISCORD_WEBHOOK")

	if apiKey == "" {
		fmt.Println("Error: GEMINI_API_KEY not set")
		os.Exit(1)
	}
	if discordWebhook == "" {
		fmt.Println("Error: DISCORD_WEBHOOK not set")
		os.Exit(1)
	}

	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		location = time.FixedZone("Bangkok", 7*60*60)
	}

	startDate, err := time.ParseInLocation("2006-01-02", gardenStartDate, location)
	if err != nil {
		fmt.Printf("Error parsing start date: %v\n", err)
		os.Exit(1)
	}

	today := time.Now().In(location)
	daysSinceStart := int(today.Sub(startDate).Hours() / 24)
	currentDay := daysSinceStart + 1

	ctx, cancel := context.WithTimeout(context.Background(), geminiTimeoutSeconds*time.Second)
	defer cancel()

	message, allMessages, err := GenerateMessage(ctx, apiKey, prompt)
	if err != nil {
		fmt.Printf("Error generating message: %v\n", err)
		fmt.Println("Using fallback message")
		message = fallbackMessage
	} else {
		fmt.Printf("Gemini generated all messages:\n")
		for i, msg := range allMessages {
			fmt.Printf("  %d: %s\n", i+1, msg)
		}
	}

	message = fmt.Sprintf("%s (วันที่ %d)", message, currentDay)

	err = SendToDiscord(discordWebhook, message, discordUsername, discordAvatarURL, discordMentionIDs)
	if err != nil {
		fmt.Printf("Error sending to Discord: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Message sent to Discord: %s\n", message)
}
