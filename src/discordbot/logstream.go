package discordbot

import (
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

// PassLogMessageToDiscordLogBuffer is called from the detection module to add a log message to the buffer.
func PassLogStreamToDiscordLogBuffer(logMessage string) {
	config.ConfigMu.Lock()
	config.LogMessageBuffer += logMessage + "\n"
	config.ConfigMu.Unlock()
	if len(config.GetLogMessageBuffer()) >= config.GetDiscordCharBufferSize() && config.GetIsDiscordEnabled() {
		flushLogBufferToDiscord()
	}
}

// FlushLogBufferToDiscord flushes the log buffer to Discord periodically with a configurable "DiscordCharBufferSize" character limit per message.
func flushLogBufferToDiscord() {
	if len(config.GetLogMessageBuffer()) == 0 {
		return // No messages to send
	}
	if !config.GetIsDiscordEnabled() || config.GetDiscordSession() == nil {
		return
	}

	discordMaxMessageLength := config.GetDiscordCharBufferSize()

	message := config.GetLogMessageBuffer()

	for len(message) > 0 {
		// Determine how much of the message we can send
		chunkSize := discordMaxMessageLength
		if len(message) < discordMaxMessageLength {
			chunkSize = len(message)
		}

		// Send the chunk to Discord
		_, err := config.GetDiscordSession().ChannelMessageSend(config.GetLogChannelID(), message[:chunkSize])
		if err != nil {
			logger.Discord.Error("Error sending log to Discord: " + err.Error())
			break
		}

		// Move to the next chunk
		message = message[chunkSize:]
	}

	// Clear the buffer after sending
	config.ConfigMu.Lock()
	config.LogMessageBuffer = ""
	config.ConfigMu.Unlock()
}
