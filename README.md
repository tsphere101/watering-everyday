# Watering Everyday Bot

Discord bot that sends daily watering reminders, triggered by Google Cloud Scheduler.

## How It Works

1. Google Cloud Scheduler triggers GitHub Actions
2. Bot calls Gemini API to generate a message
3. Message is sent to Discord with user mentions

## Setup

1. Add secrets: `GEMINI_API_KEY`, `DISCORD_WEBHOOK`
2. Configure in `src/main.go`: `prompt`, `gardenStartDate`, `discordMentionIDs`
3. Build: `go build -o watering-bot ./src`

## Trigger

- Manual: GitHub Actions `workflow_dispatch`
- Automatic: Google Cloud Scheduler → `repository_dispatch`
