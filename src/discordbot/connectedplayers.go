package discordbot

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
)

var (
	connectedPlayersMessageID string // connectedPlayersMessageID tracks the message ID for editing the connected players message
	playersMutex              sync.Mutex
)

func AddToConnectedPlayers(username, steamID string, connectionTime time.Time, players map[string]string) {
	if !config.GetIsDiscordEnabled() || config.GetDiscordSession() == nil {
		logger.Discord.Debug("Discord not enabled or session not initialized")
		return
	}
	content := formatConnectedPlayers(players)
	sendAndEditMessageInConnectedPlayersChannel(config.GetConnectionListChannelID(), content)
}

func RemoveFromConnectedPlayers(steamID string, players map[string]string) {
	if !config.GetIsDiscordEnabled() || config.GetDiscordSession() == nil {
		logger.Discord.Debug("Discord not enabled or session not initialized")
		return
	}
	content := formatConnectedPlayers(players)
	sendAndEditMessageInConnectedPlayersChannel(config.GetConnectionListChannelID(), content)
}

func sendAndEditMessageInConnectedPlayersChannel(channelID, message string) {
	playersMutex.Lock()
	defer playersMutex.Unlock()

	if connectedPlayersMessageID == "" {
		// Send a new message if there's no existing one
		msg, err := config.GetDiscordSession().ChannelMessageSend(channelID, message)
		if err != nil {
			logger.Discord.Error("Error sending message to channel " + channelID + ": " + err.Error())
			return
		}
		connectedPlayersMessageID = msg.ID
		logger.Discord.Debug("Sent new message to channel " + channelID)
	} else {
		// Edit the existing message
		_, err := config.GetDiscordSession().ChannelMessageEdit(channelID, connectedPlayersMessageID, message)
		if err != nil {
			logger.Discord.Error("Error editing message in channel " + channelID + ": " + err.Error())
			// If editing fails (e.g., message deleted), reset and try sending a new one
			connectedPlayersMessageID = ""
			msg, err := config.GetDiscordSession().ChannelMessageSend(channelID, message)
			if err != nil {
				logger.Discord.Error("Error sending fallback message to channel " + channelID + ": " + err.Error())
			} else {
				connectedPlayersMessageID = msg.ID
				logger.Discord.Debug("Sent new message after edit failure to channel " + channelID)
			}
		}
	}
}

func formatConnectedPlayers(players map[string]string) string {
	if len(players) == 0 {
		return "No players are currently connected."
	}

	var sb strings.Builder
	sb.WriteString("Connected Players:\n")
	sb.WriteString("```\n")
	sb.WriteString("Username              | Steam ID\n")
	sb.WriteString("----------------------|------------------------\n")

	for steamID, username := range players {
		sb.WriteString(fmt.Sprintf("%-20s | %s\n", username, steamID))
	}

	sb.WriteString("```")
	return sb.String()
}
