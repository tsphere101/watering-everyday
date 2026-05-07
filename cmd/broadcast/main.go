package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type discordRequest struct {
	Content         string `json:"content"`
	Username        string `json:"username"`
	AvatarURL       string `json:"avatar_url"`
	AllowedMentions struct {
		Users []string `json:"users"`
	} `json:"allowed_mentions"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: broadcast <message>")
		os.Exit(1)
	}

	webhook := os.Getenv("DISCORD_WEBHOOK")
	if webhook == "" {
		fmt.Println("Error: DISCORD_WEBHOOK not set")
		os.Exit(1)
	}

	message := os.Args[1]
	username := "รดน้ำ"
	avatarURL := "https://raw.githubusercontent.com/tsphere101/watering-everyday/refs/heads/main/watering-avatar.jpg"

	payload := discordRequest{
		Content:   message,
		Username:  username,
		AvatarURL: avatarURL,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	resp, err := http.Post(webhook, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		fmt.Printf("Error sending: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		fmt.Printf("Discord returned status %d\n", resp.StatusCode)
		os.Exit(1)
	}

	fmt.Println("Sent!")
}
