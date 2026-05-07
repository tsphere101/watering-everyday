package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"google.golang.org/genai"
)

//
// splitMessage
//

func TestSplitMessage_Short(t *testing.T) {
	chunks := splitMessage("hello", 10)
	if len(chunks) != 1 || chunks[0] != "hello" {
		t.Errorf("expected [hello], got %v", chunks)
	}
}

func TestSplitMessage_ExactLength(t *testing.T) {
	chunks := splitMessage("hello", 5)
	if len(chunks) != 1 || chunks[0] != "hello" {
		t.Errorf("expected [hello], got %v", chunks)
	}
}

func TestSplitMessage_Multiple(t *testing.T) {
	chunks := splitMessage("abcdefghij", 3)
	expected := []string{"abc", "def", "ghi", "j"}
	if len(chunks) != len(expected) {
		t.Fatalf("expected %d chunks, got %d", len(expected), len(chunks))
	}
	for i := range expected {
		if chunks[i] != expected[i] {
			t.Errorf("chunk %d: expected %q, got %q", i, expected[i], chunks[i])
		}
	}
}

func TestSplitMessage_PreservesContent(t *testing.T) {
	msg := "hello world this is a test"
	chunks := splitMessage(msg, 5)
	joined := strings.Join(chunks, "")
	if joined != msg {
		t.Errorf("content changed: got %q, want %q", joined, msg)
	}
}

func TestSplitMessage_Empty(t *testing.T) {
	chunks := splitMessage("", 10)
	if len(chunks) != 1 || chunks[0] != "" {
		t.Errorf("expected [\"\"], got %v", chunks)
	}
}

func TestSplitMessage_Unicode(t *testing.T) {
	msg := "สวัสดี"
	chunks := splitMessage(msg, 3)
	if len(chunks) != 2 {
		t.Fatalf("expected 2 chunks, got %d (rune count: %d)", len(chunks), len([]rune(msg)))
	}
	if chunks[0]+chunks[1] != msg {
		t.Errorf("content changed: got %q, want %q", chunks[0]+chunks[1], msg)
	}
	if len([]rune(chunks[0])) != 3 || len([]rune(chunks[1])) != 3 {
		t.Errorf("expected each chunk to have 3 runes, got %d and %d", len([]rune(chunks[0])), len([]rune(chunks[1])))
	}
}

func TestSplitMessage_DiscordBoundary(t *testing.T) {
	msg := strings.Repeat("a", 1900)
	chunks := splitMessage(msg, 1900)
	if len(chunks) != 1 {
		t.Errorf("expected 1 chunk, got %d", len(chunks))
	}
}

func TestSplitMessage_JustOverBoundary(t *testing.T) {
	msg := strings.Repeat("a", 1901)
	chunks := splitMessage(msg, 1900)
	if len(chunks) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(chunks))
	}
	if len(chunks[0]) != 1900 {
		t.Errorf("first chunk length: expected 1900, got %d", len(chunks[0]))
	}
	if len(chunks[1]) != 1 {
		t.Errorf("second chunk length: expected 1, got %d", len(chunks[1]))
	}
}

func TestSplitMessage_EqualChunks(t *testing.T) {
	chunks := splitMessage("abcdef", 2)
	if len(chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(chunks))
	}
	if chunks[0] != "ab" || chunks[1] != "cd" || chunks[2] != "ef" {
		t.Errorf("unexpected split: %v", chunks)
	}
}

//
// makeMentionString
//

func TestMakeMentionString_Empty(t *testing.T) {
	r := makeMentionString([]string{})
	if r != " " {
		t.Errorf("expected \" \", got %q", r)
	}
}

func TestMakeMentionString_Single(t *testing.T) {
	r := makeMentionString([]string{"123"})
	if r != " <@123>" {
		t.Errorf("expected \" <@123>\", got %q", r)
	}
}

func TestMakeMentionString_Multiple(t *testing.T) {
	r := makeMentionString([]string{"123", "456"})
	if r != " <@123> <@456>" {
		t.Errorf("expected \" <@123> <@456>\", got %q", r)
	}
}

//
// GetCurrentDayWithTime
//

func TestGetCurrentDayWithTime(t *testing.T) {
	tests := []struct {
		name        string
		startDate   string
		inputTime   time.Time
		expectedDay int
	}{
		{
			name:        "Day 1 - Aug 13 2025",
			startDate:   "2025-08-13",
			inputTime:   time.Date(2025, 8, 13, 21, 0, 0, 0, time.FixedZone("Bangkok", 7*60*60)),
			expectedDay: 1,
		},
		{
			name:        "Day 2 - Aug 14 2025",
			startDate:   "2025-08-13",
			inputTime:   time.Date(2025, 8, 14, 21, 0, 0, 0, time.FixedZone("Bangkok", 7*60*60)),
			expectedDay: 2,
		},
		{
			name:        "Day 199 - Feb 27 2026",
			startDate:   "2025-08-13",
			inputTime:   time.Date(2026, 2, 27, 21, 0, 0, 0, time.FixedZone("Bangkok", 7*60*60)),
			expectedDay: 199,
		},
		{
			name:        "Day 200 - Feb 28 2026",
			startDate:   "2025-08-13",
			inputTime:   time.Date(2026, 2, 28, 21, 0, 0, 0, time.FixedZone("Bangkok", 7*60*60)),
			expectedDay: 200,
		},
		{
			name:        "Day 100 - Nov 20 2025",
			startDate:   "2025-08-13",
			inputTime:   time.Date(2025, 11, 20, 21, 0, 0, 0, time.FixedZone("Bangkok", 7*60*60)),
			expectedDay: 100,
		},
		{
			name:        "Midnight Feb 28 2026",
			startDate:   "2025-08-13",
			inputTime:   time.Date(2026, 2, 28, 0, 0, 0, 0, time.FixedZone("Bangkok", 7*60*60)),
			expectedDay: 200,
		},
		{
			name:        "Day 200 at midnight Feb 28 2026 +1s",
			startDate:   "2025-08-13",
			inputTime:   time.Date(2026, 2, 28, 0, 0, 1, 0, time.FixedZone("Bangkok", 7*60*60)),
			expectedDay: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			day, err := GetCurrentDayWithTime(tt.startDate, tt.inputTime)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if day != tt.expectedDay {
				t.Errorf("expected day %d, got %d", tt.expectedDay, day)
			}
		})
	}
}

func TestGetCurrentDayWithTime_InvalidDate(t *testing.T) {
	_, err := GetCurrentDayWithTime("not-a-date", time.Now())
	if err == nil {
		t.Error("expected error for invalid date")
	}
}

func TestGetCurrentDayWithTime_LoadLocationFallback(t *testing.T) {
	t.Setenv("ZONEINFO", "/nonexistent/zoneinfo.zip")
	day, err := GetCurrentDayWithTime("2025-08-13", time.Date(2025, 8, 13, 21, 0, 0, 0, time.FixedZone("Bangkok", 7*60*60)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if day != 1 {
		t.Errorf("expected day 1, got %d", day)
	}
}

//
// GetCurrentDay
//

func TestGetCurrentDay_ValidDate(t *testing.T) {
	day, err := GetCurrentDay("2025-08-13")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if day < 1 {
		t.Errorf("expected day >= 1, got %d", day)
	}
}

func TestGetCurrentDay_InvalidDate(t *testing.T) {
	_, err := GetCurrentDay("not-a-date")
	if err == nil {
		t.Error("expected error for invalid date")
	}
}

//
// SendToDiscord
//

func TestSendToDiscord_SingleChunk(t *testing.T) {
	var gotReq DiscordRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json")
		}
		if err := json.NewDecoder(r.Body).Decode(&gotReq); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	err := SendToDiscord(server.URL, "test message", "test", "http://avatar", []string{"123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotReq.Content != "test message <@123>" {
		t.Errorf("unexpected content: %q", gotReq.Content)
	}
	if gotReq.Username != "test" {
		t.Errorf("unexpected username: %q", gotReq.Username)
	}
	if gotReq.AvatarURL != "http://avatar" {
		t.Errorf("unexpected avatar URL: %q", gotReq.AvatarURL)
	}
}

func TestSendToDiscord_MultipleChunks(t *testing.T) {
	var gotContents []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req DiscordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode: %v", err)
		}
		gotContents = append(gotContents, req.Content)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	longMsg := strings.Repeat("a", 5000)
	err := SendToDiscord(server.URL, longMsg, "test", "http://avatar", []string{"123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gotContents) < 3 {
		t.Fatalf("expected at least 3 chunks, got %d", len(gotContents))
	}
	// verify content is preserved
	joined := strings.Join(gotContents, "")
	joined = strings.ReplaceAll(joined, "... (see next message)", "")
	if joined != longMsg+" <@123>" {
		t.Errorf("content mismatch after joining chunks\ngot:  %q\nwant: %q", joined, longMsg+" <@123>")
	}
}

func TestSendToDiscord_MiddleChunkHasMarker(t *testing.T) {
	var gotContents []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req DiscordRequest
		json.NewDecoder(r.Body).Decode(&req)
		gotContents = append(gotContents, req.Content)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	longMsg := strings.Repeat("a", 5000)
	SendToDiscord(server.URL, longMsg, "test", "http://avatar", []string{"123"})

	for i, c := range gotContents {
		if i > 0 && i < len(gotContents)-1 {
			if !strings.Contains(c, "... (see next message)") {
				t.Errorf("middle chunk %d missing continuation marker", i)
			}
		}
	}
}

func TestSendToDiscord_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer server.Close()

	err := SendToDiscord(server.URL, "test", "test", "http://avatar", []string{"123"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("error should mention status 400, got: %v", err)
	}
}

func TestSendToDiscord_ServerDown(t *testing.T) {
	err := SendToDiscord("http://127.0.0.1:1", "test", "test", "http://avatar", []string{"123"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSendToDiscord_ChunkHTTPError_AbortsEarly(t *testing.T) {
	var callCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 2 {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// message long enough to require 3+ chunks
	err := SendToDiscord(server.URL, strings.Repeat("a", 5000), "test", "http://avatar", []string{"123"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// should not have sent all chunks (we errored on chunk 2; chunk 3 should not be sent)
	if callCount >= 3 {
		t.Errorf("expected < 3 calls after error, got %d", callCount)
	}
}

func TestSendToDiscord_StatusOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	err := SendToDiscord(server.URL, "test", "test", "http://avatar", []string{"123"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSendToDiscord_EmptyMentions(t *testing.T) {
	var gotReq DiscordRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	SendToDiscord(server.URL, "hello", "test", "http://avatar", []string{})
	if gotReq.Content != "hello " {
		t.Errorf("expected \"hello \", got %q", gotReq.Content)
	}
}

func TestSendToDiscord_MultipleMentions(t *testing.T) {
	var gotReq DiscordRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	SendToDiscord(server.URL, "water", "bot", "http://pic", []string{"111", "222", "333"})
	if !strings.Contains(gotReq.Content, "<@111>") || !strings.Contains(gotReq.Content, "<@222>") || !strings.Contains(gotReq.Content, "<@333>") {
		t.Errorf("expected all mentions in content, got %q", gotReq.Content)
	}
}

//
// GenerateMessage
//
// These tests mock the Gemini API via genai.SetDefaultBaseURLs + httptest.
//

func geminiMockServer(t *testing.T, responseJSON string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(responseJSON))
	}))
}

func setGeminiBaseURL(t *testing.T, url string) {
	t.Helper()
	genai.SetDefaultBaseURLs(genai.BaseURLParameters{GeminiURL: url})
	t.Cleanup(func() { genai.SetDefaultBaseURLs(genai.BaseURLParameters{GeminiURL: "", VertexURL: ""}) })
}

func TestGenerateMessage_Success(t *testing.T) {
	server := geminiMockServer(t, `{
		"candidates": [{
			"content": { "parts": [{ "text": "{\"messages\":[\"msg A\",\"msg B\",\"msg C\"]}" }] },
			"finishReason": "STOP"
		}]
	}`)
	defer server.Close()
	setGeminiBaseURL(t, server.URL)

	msg, msgs, err := GenerateMessage(context.Background(), "test-key", "test prompt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg == "" {
		t.Error("expected non-empty message")
	}
	if len(msgs) != 3 {
		t.Errorf("expected 3 messages, got %d", len(msgs))
	}
}

func TestGenerateMessage_ClientCreationFails(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "")
	t.Setenv("GOOGLE_API_KEY", "")
	_, _, err := GenerateMessage(context.Background(), "", "test prompt")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to create Gemini client") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGenerateMessage_APIFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()
	setGeminiBaseURL(t, server.URL)

	_, _, err := GenerateMessage(context.Background(), "test-key", "test")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to generate content") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGenerateMessage_NoCandidates(t *testing.T) {
	server := geminiMockServer(t, `{"candidates": []}`)
	defer server.Close()
	setGeminiBaseURL(t, server.URL)

	_, _, err := GenerateMessage(context.Background(), "test-key", "test")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "no response") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGenerateMessage_EmptyParts(t *testing.T) {
	server := geminiMockServer(t, `{
		"candidates": [{
			"content": { "parts": [] },
			"finishReason": "STOP"
		}]
	}`)
	defer server.Close()
	setGeminiBaseURL(t, server.URL)

	_, _, err := GenerateMessage(context.Background(), "test-key", "test")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "no response") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGenerateMessage_InvalidJSON(t *testing.T) {
	server := geminiMockServer(t, `{
		"candidates": [{
			"content": { "parts": [{ "text": "not valid json" }] },
			"finishReason": "STOP"
		}]
	}`)
	defer server.Close()
	setGeminiBaseURL(t, server.URL)

	_, _, err := GenerateMessage(context.Background(), "test-key", "test")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to parse JSON") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGenerateMessage_NoMessagesInResponse(t *testing.T) {
	server := geminiMockServer(t, `{
		"candidates": [{
			"content": { "parts": [{ "text": "{\"messages\":[]}" }] },
			"finishReason": "STOP"
		}]
	}`)
	defer server.Close()
	setGeminiBaseURL(t, server.URL)

	_, _, err := GenerateMessage(context.Background(), "test-key", "test")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "no messages in response") {
		t.Errorf("unexpected error message: %v", err)
	}
}

//
// loadPrompt
//

func TestLoadPrompt_FileNotFound(t *testing.T) {
	t.Chdir(t.TempDir())
	result := loadPrompt(time.FixedZone("Bangkok", 7*60*60))
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

func TestLoadPrompt_InvalidJSON(t *testing.T) {
	t.Chdir(t.TempDir())
	os.WriteFile("prompts.json", []byte("{bad json}"), 0644)
	result := loadPrompt(time.FixedZone("Bangkok", 7*60*60))
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

func TestLoadPrompt_OverrideMatch(t *testing.T) {
	t.Chdir(t.TempDir())
	bkk := time.FixedZone("Bangkok", 7*60*60)
	today := time.Now().In(bkk).Format("2006-01-02")
	cfg := fmt.Sprintf(`{"overrides": {"%s": "test override"}}`, today)
	os.WriteFile("prompts.json", []byte(cfg), 0644)

	result := loadPrompt(bkk)
	if result != "test override" {
		t.Errorf("expected 'test override', got %q", result)
	}
}

func TestLoadPrompt_AnnualMatch(t *testing.T) {
	t.Chdir(t.TempDir())
	bkk := time.FixedZone("Bangkok", 7*60*60)
	key := time.Now().In(bkk).Format("01-02")
	cfg := fmt.Sprintf(`{"annual": {"%s": "test annual"}}`, key)
	os.WriteFile("prompts.json", []byte(cfg), 0644)

	result := loadPrompt(bkk)
	if result != "test annual" {
		t.Errorf("expected 'test annual', got %q", result)
	}
}

func TestLoadPrompt_DefaultFallback(t *testing.T) {
	t.Chdir(t.TempDir())
	os.WriteFile("prompts.json", []byte(`{"default": "test default"}`), 0644)

	result := loadPrompt(time.FixedZone("Bangkok", 7*60*60))
	if result != "test default" {
		t.Errorf("expected 'test default', got %q", result)
	}
}

func TestLoadPrompt_OverridePrecedesAnnual(t *testing.T) {
	t.Chdir(t.TempDir())
	bkk := time.FixedZone("Bangkok", 7*60*60)
	today := time.Now().In(bkk).Format("2006-01-02")
	annualKey := time.Now().In(bkk).Format("01-02")
	cfg := fmt.Sprintf(`{"overrides": {"%s": "override"}, "annual": {"%s": "annual"}}`, today, annualKey)
	os.WriteFile("prompts.json", []byte(cfg), 0644)

	result := loadPrompt(bkk)
	if result != "override" {
		t.Errorf("expected 'override', got %q", result)
	}
}

func TestLoadPrompt_AnnualPrecedesDefault(t *testing.T) {
	t.Chdir(t.TempDir())
	bkk := time.FixedZone("Bangkok", 7*60*60)
	key := time.Now().In(bkk).Format("01-02")
	cfg := fmt.Sprintf(`{"annual": {"%s": "annual"}, "default": "default"}`, key)
	os.WriteFile("prompts.json", []byte(cfg), 0644)

	result := loadPrompt(bkk)
	if result != "annual" {
		t.Errorf("expected 'annual', got %q", result)
	}
}
