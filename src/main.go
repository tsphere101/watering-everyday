package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

const (
	defaultPrompt        = "เขียนประโยคเชิญชวนเพื่อนชื่อ เชอร์รี่ และ ท่านหญิงนวว มารดน้ำต้นไม้ 10 แบบ แตกต่างกัน (รดน้ำต้นไม้ในเกมทุกวันเพื่อไม่ให้ต้นไม้เสียชีวิต)"
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
	shouldSend := os.Getenv("SEND_TO_DISCORD") == "true"

	if apiKey == "" {
		fmt.Println(errAPIKeyNotSet)
		os.Exit(1)
	}
	if shouldSend && discordWebhook == "" {
		fmt.Println(errWebhookNotSet)
		os.Exit(1)
	}

	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		location = time.FixedZone("Bangkok", 7*60*60)
	}

	prompt := loadPrompt(location)
	if prompt == "" {
		prompt = defaultPrompt
	}

	startDate, err := time.ParseInLocation("2006-01-02", gardenStartDate, location)
	if err != nil {
		fmt.Printf(errParseStartDate+"\n", err)
		os.Exit(1)
	}

	today := time.Now().In(location)
	daysSinceStart := int(today.Sub(startDate).Hours() / 24)
	currentDay := daysSinceStart + 1

	ctx, cancel := context.WithTimeout(context.Background(), geminiTimeoutSeconds*time.Second)
	defer cancel()

	message, allMessages, err := GenerateMessage(ctx, apiKey, prompt)
	if err != nil {
		fmt.Printf(errGenerateMsg+"\n", err)
		fmt.Println(msgUsingFallback)
		message = fallbackMessage
	} else {
		fmt.Println(msgGeminiMsgs)
		for i, msg := range allMessages {
			fmt.Printf(msgGeminiMsgFmt+"\n", i+1, msg)
		}
	}

	message = fmt.Sprintf(msgDayFmt, message, currentDay)

	fmt.Printf(msgOutputSep+"\n", message)

	if shouldSend {
		err = SendToDiscord(discordWebhook, message, discordUsername, discordAvatarURL, discordMentionIDs)
		if err != nil {
			fmt.Printf(errSendToDiscord+"\n", err)
			os.Exit(1)
		}
		fmt.Printf(msgSentToDiscord+"\n", message)
	} else {
		fmt.Println(msgSkipDiscord)
	}
}
