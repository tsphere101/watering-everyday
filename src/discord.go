package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type DiscordRequest struct {
	Content         string `json:"content"`
	Username        string `json:"username"`
	AvatarURL       string `json:"avatar_url"`
	AllowedMentions struct {
		Users []string `json:"users"`
	} `json:"allowed_mentions"`
}

func SendToDiscord(webhook, message, username, avatarURL string, mentionIDs []string) error {
	mentionStr := makeMentionString(mentionIDs)
	fullMessage := message + mentionStr

	chunks := splitMessage(fullMessage, discordMaxLength)

	client := &http.Client{}
	var lastResp *http.Response

	for i, chunk := range chunks {
		if i > 0 && i < len(chunks)-1 {
			chunk += "... (see next message)"
		}

		discordReq := DiscordRequest{
			Content:   chunk,
			Username:  username,
			AvatarURL: avatarURL,
		}
		discordReq.AllowedMentions.Users = mentionIDs

		jsonData, err := json.Marshal(discordReq)
		if err != nil {
			return fmt.Errorf("failed to marshal discord request: %w", err)
		}

		req, err := http.NewRequest("POST", webhook, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to create discord request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		lastResp, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send to discord: %w", err)
		}
		lastResp.Body.Close()
	}

	if lastResp != nil && lastResp.StatusCode != http.StatusNoContent && lastResp.StatusCode != http.StatusOK {
		return fmt.Errorf("discord returned status: %d", lastResp.StatusCode)
	}

	return nil
}

func makeMentionString(ids []string) string {
	var mentions []string
	for _, id := range ids {
		mentions = append(mentions, "<@"+id+">")
	}
	return " " + strings.Join(mentions, " ")
}
