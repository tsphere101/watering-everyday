# AGENTS.md - Developer Guide

This document provides guidelines for agents working in this repository.

## Project Overview

This is a Go-based Discord bot that generates watering reminders using Gemini AI. The bot is triggered by Google Cloud Scheduler via GitHub Actions repository_dispatch events.

## Build Commands

### Build the binary
```bash
go build -o watering-bot ./src
```

### Run locally (for testing)
```bash
GEMINI_API_KEY=your_key DISCORD_WEBHOOK=your_webhook ./watering-bot
```

### Run tests (if any)
```bash
go test ./...
```

### Run a single test
```bash
go test -v ./src -run TestName
```

### Format code
```bash
go fmt ./...
```

### Lint code
```bash
go vet ./...
```

### Tidy dependencies
```bash
go mod tidy
```

## Project Structure

```
watering-everyday/
├── src/
│   ├── main.go    # Entry point + configs + prompt
│   ├── gemini.go  # Gemini API integration
│   ├── discord.go # Discord webhook integration
│   └── utils.go   # Utility functions
├── watering-bot   # Prebuilt binary (commit to repo)
├── go.mod         # Go module
└── .github/
    └── workflows/
        └── watering.yml
```

## Code Style Guidelines

### Imports

- Use standard library packages first, then third-party
- Group: standard library → external → project internal
- Use implicit grouping (no blank line needed in Go 1.17+)

```go
import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "some/external/pkg"
)
```

### Formatting

- Use `gofmt` or automatic formatting in IDE
- Indent with tabs, not spaces
- Max line length: 100 characters (soft limit)

### Types

- Use explicit types for function parameters
- Use meaningful type names (CamelCase for exported, camelCase for unexported)
- Define types close to where they are used

```go
type GeminiRequest struct {
    Contents []Content `json:"contents"`
}

type Content struct {
    Parts []Part `json:"parts"`
}
```

### Naming Conventions

- **Variables/Functions:** camelCase (unexported), CamelCase (exported)
- **Constants:** CamelCase for exported, camelCase or CamelCase for unexported
- **Packages:** short, lowercase, no underscores
- **Files:** lowercase with underscores only if needed (e.g., `gemini.go`)

```go
const (
    prompt            = "..." // config constant
    geminiURL         = "..." // config constant
    discordMaxLength  = 1900  // config constant
)

var (
    discordUsername   = "รดน้ำ"     // config variable
    discordAvatarURL  = "https://..."
)
```

### Error Handling

- Return errors as last return value
- Use `fmt.Errorf` with `%w` for wrapped errors
- Check errors immediately after function calls
- Use sentinel errors for known error conditions

```go
func GenerateMessage(apiKey, prompt string) (string, error) {
    // ...
    if err != nil {
        return "", fmt.Errorf("failed to marshal request: %w", err)
    }
    // ...
}
```

### Configuration

- All configurable values (URLs, prompts, IDs) should be in `main.go`
- Use `const` for values that never change
- Use `var` for values that may be changed later
- Read sensitive values from environment variables

```go
const (
    prompt           = "your prompt here"
    geminiURL        = "https://..."
    discordMaxLength = 1900
)

var (
    discordUsername  = "BotName"
    discordAvatarURL = "https://..."
)
```

### HTTP Requests

- Always set `Content-Type: application/json`
- Always close response body (`defer resp.Body.Close()`)
- Check status codes for non-200 responses

```go
req.Header.Set("Content-Type", "application/json")
resp, err := client.Do(req)
if err != nil {
    return "", fmt.Errorf("failed to call API: %w", err)
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
    return "", fmt.Errorf("API returned status: %d", resp.StatusCode)
}
```

### JSON Handling

- Use struct tags for JSON field mapping
- Use `json.Marshal` for encoding, `json.NewDecoder` for decoding streams

```go
type Request struct {
    Field string `json:"field"`
}
```

### Discord API Specifics

- Max message length: 2000 characters
- Use 1900 as safe limit to allow continuation markers
- Split long messages with "... (see next message)" marker

```go
chunks := splitMessage(fullMessage, discordMaxLength)
for i, chunk := range chunks {
    if i > 0 && i < len(chunks)-1 {
        chunk += "... (see next message)"
    }
    // send chunk
}
```

### Logging

- Use `fmt.Println` for simple output
- Use `fmt.Printf` with `%v` for variable output
- Print errors with `fmt.Printf("Error: %v\n", err)`

### GitHub Actions Workflow

- Keep workflow simple: checkout → run binary
- Use prebuilt binary (no build step in CI)
- Pass secrets via environment variables

```yaml
- name: Run bot
  env:
    DISCORD_WEBHOOK: ${{ secrets.DISCORD_WEBHOOK }}
    GEMINI_API_KEY: ${{ secrets.GEMINI_API_KEY }}
  run: ./watering-bot
```

### Testing Guidelines

- Place tests in same package with `_test.go` suffix
- Use table-driven tests for multiple test cases
- Use descriptive test names: `TestFunctionName_Scenario`

```go
func TestSplitMessage(t *testing.T) {
    tests := []struct {
        name     string
        msg      string
        maxLen   int
        expected int
    }{
        {"short message", "hello", 10, 1},
        {"long message", strings.Repeat("a", 50), 20, 3},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := splitMessage(tt.msg, tt.maxLen)
            if len(result) != tt.expected {
                t.Errorf("expected %d chunks, got %d", tt.expected, len(result))
            }
        })
    }
}
```

### Commit Guidelines

- Commit binary after code changes: `git add watering-bot`
- Write concise commit messages describing what changed
- Example: `Add message splitting feature` or `Refactor to modular structure`

## Secrets Management

| Secret | Description | Where to get |
|--------|-------------|--------------|
| `GEMINI_API_KEY` | Google Gemini API key | Google AI Studio |
| `DISCORD_WEBHOOK` | Discord webhook URL | Discord Server Settings → Integrations |

## CI/CD Flow

1. Google Cloud Scheduler triggers `repository_dispatch`
2. GitHub Actions runs workflow
3. Checkout (2s)
4. Run prebuilt binary (~3-5s total with API calls)
5. Done

## Notes

- This project uses only Go standard library (no external dependencies)
- Binary is prebuilt locally and committed to repo for faster CI
- All configs/prompts are centralized in `main.go` for easy modification
