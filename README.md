# Watering Everyday Bot

This Discord bot automates daily watering reminders in ขายไร่-ขายนา using GitHub Actions with Google Cloud Scheduler.

## How It Works

1. **Google Cloud Scheduler** triggers the workflow via `repository_dispatch`
2. **GitHub Actions** runs the prebuilt Go binary (`watering-bot`)
3. **Gemini API** generates a unique Thai message each time
4. **Discord Webhook** sends the message with user mentions

## Tech Stack

- **Language:** Go
- **AI:** Gemini 2.0 Flash
- **CI/CD:** GitHub Actions
- **Trigger:** Google Cloud Scheduler

## Directory Structure

```
watering-everyday/
├── src/
│   └── main.go          # Go source code
├── watering-bot         # Prebuilt binary
├── go.mod               # Go module
└── .github/
    └── workflows/
        └── watering.yml
```

## Setup

### 1. Add GitHub Secrets

| Secret | Description |
|--------|-------------|
| `GEMINI_API_KEY` | Gemini API key |
| `DISCORD_WEBHOOK` | Discord webhook URL |

### 2. Build Binary (if code changes)

```bash
go build -o watering-bot ./src
git add watering-bot
git commit -m "Update binary"
git push
```

## Trigger

- **Manual:** GitHub Actions workflow_dispatch
- **Automatic:** Google Cloud Scheduler sends POST to repository_dispatch

## CI Runtime

~3-5 seconds (no Python setup required)
