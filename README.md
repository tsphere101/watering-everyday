# Watering Everyday Bot

A simple Discord bot that sends a daily watering reminder at 21:00 UTC+7.

## Features
- Daily scheduled message via GitHub Actions.
- Mentions specific users to ensure they don't forget!
- Runs automatically at 14:00 UTC (21:00 UTC+7).

## Setup
1. Fork/Clone this repository.
2. Go to your repository **Settings > Secrets and variables > Actions**.
3. Add a new repository secret:
   - Name: `DISCORD_WEBHOOK`
   - Value: Your Discord Webhook URL.

## Local Execution
To run the bot locally (requires `DISCORD_WEBHOOK` environment variable):

```bash
export DISCORD_WEBHOOK="your_webhook_url"
python main.py
```