package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"google.golang.org/genai"
)

func GenerateMessage(ctx context.Context, apiKey, prompt string) (string, []string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	parts := []*genai.Part{
		{Text: prompt},
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"messages": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeString,
					},
				},
			},
			Required: []string{"messages"},
		},
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-3-flash-preview", []*genai.Content{{Parts: parts}}, config)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", nil, fmt.Errorf("no response from Gemini")
	}

	jsonStr := result.Candidates[0].Content.Parts[0].Text
	jsonStr = strings.TrimSpace(jsonStr)

	var response struct {
		Messages []string `json:"messages"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
		return "", nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	if len(response.Messages) == 0 {
		return "", nil, fmt.Errorf("no messages in response")
	}

	dayOfYear := time.Now().YearDay()
	selectedIndex := dayOfYear % len(response.Messages)

	message := response.Messages[selectedIndex]
	message = strings.TrimSpace(message)

	return message, response.Messages, nil
}
