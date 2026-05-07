# Watering Everyday Bot

Discord bot that sends daily watering reminders, triggered by Google Cloud Scheduler.

## How It Works

1. Google Cloud Scheduler triggers GitHub Actions
2. Bot calls Gemini API to generate a message
3. Message is sent to Discord with user mentions

## Setup

1. Add GitHub Secrets: `GEMINI_API_KEY`, `DISCORD_WEBHOOK`
2. Edit prompts in `prompts.json` (no code changes needed)
3. Build and commit the binary

```bash
make build        # lint → tests → linux/amd64 binaries (for production)
make build-local  # macOS binaries (for local testing)
make test         # run all tests
```

## Trigger

- Manual: GitHub Actions `workflow_dispatch`
- Automatic: Google Cloud Scheduler → `repository_dispatch`

## Local Testing

Run without sending to Discord (safe — stdout only):
```bash
./watering-bot
```

Send for real:
```bash
SEND_TO_DISCORD=true ./watering-bot
```

## Broadcast Tool

Send a custom message to Discord with the same bot profile:

```bash
DISCORD_WEBHOOK=your_url ./broadcast "รดน้ำกันจ้า"
```
