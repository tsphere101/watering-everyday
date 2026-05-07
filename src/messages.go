package main

const (
	// main.go
	errAPIKeyNotSet    = "Error: GEMINI_API_KEY not set"
	errWebhookNotSet   = "Error: DISCORD_WEBHOOK not set (required when SEND_TO_DISCORD=true)"
	errParseStartDate  = "Error parsing start date: %v"
	errGenerateMsg     = "Error generating message: %v"
	errSendToDiscord   = "Error sending to Discord: %v"
	msgUsingFallback   = "Using fallback message"
	msgSkipDiscord     = "SEND_TO_DISCORD not set — skipping Discord. Set SEND_TO_DISCORD=true to post."
	msgSentToDiscord   = "Message sent to Discord: %s"
	msgGeminiMsgs      = "Gemini generated all messages:"
	msgGeminiMsgFmt    = "  %d: %s"
	msgDayFmt          = "%s (วันที่ %d)"
	msgOutputSep       = "---\n%s"

	// gemini.go
	errCreateGeminiClient = "failed to create Gemini client: %w"
	errGenerateContent    = "failed to generate content: %w"
	errNoResponseGemini   = "no response from Gemini"
	errParseJSONResponse  = "failed to parse JSON response: %w"
	errNoMessages         = "no messages in response"

	// discord.go
	errMarshalDiscord  = "failed to marshal discord request: %w"
	errCreateDiscordReq = "failed to create discord request: %w"
	errSendChunk       = "failed to send chunk %d to discord: %w"
	errDiscordStatus   = "discord returned status %d for chunk %d"

	// prompts.go
	warnReadPrompt      = "Warning: cannot read prompts.json: %v"
	warnParsePrompt     = "Warning: prompts.json parse error: %v"
	infoPromptOverride  = "Prompt source: prompts.json > overrides[%s]"
	infoPromptAnnual    = "Prompt source: prompts.json > annual[%s]"
	infoPromptDefault   = "Prompt source: prompts.json > default"
)
