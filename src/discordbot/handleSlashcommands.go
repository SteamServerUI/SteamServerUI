package discordbot

import (
	"fmt"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/gamemgr"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"

	"github.com/bwmarrin/discordgo"
)

type commandHandler func(*discordgo.Session, *discordgo.InteractionCreate, EmbedData) error

// Command handlers map
var handlers = map[string]commandHandler{
	"start":   handleStart,
	"stop":    handleStop,
	"status":  handleStatus,
	"help":    handleHelp,
	"restore": handleRestore,
	"list":    handleList,
}

// Check channel and handle initial validation
func listenToSlashCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand || i.ChannelID != config.GetControlChannelID() {
		respond(s, i, EmbedData{
			Title: "Wrong Channel", Description: "Commands must be sent to the configured control channel",
			Color: 0xFF0000, Fields: []EmbedField{{Name: "Accepted Channel", Value: fmt.Sprintf("<#%s>", config.GetControlChannelID()), Inline: true}},
		})
		return
	}

	cmd := i.ApplicationCommandData().Name
	if handler, ok := handlers[cmd]; ok {
		data := EmbedData{Title: "Command Error", Color: 0xFF0000}
		if err := handler(s, i, data); err != nil {
			logger.Discord.Error("Error handling " + cmd + ": " + err.Error())
		}
	}
}

// Generic response function
func respond(s *discordgo.Session, i *discordgo.InteractionCreate, embed EmbedData) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{generateEmbed(embed)},
		},
	})
}

func handleStart(s *discordgo.Session, i *discordgo.InteractionCreate, data EmbedData) error {
	data.Title, data.Description, data.Color = "Server Control", "Starting the server...", 0x00FF00
	data.Fields = []EmbedField{{Name: "Status", Value: "🕛 Recieved", Inline: true}}
	if err := respond(s, i, data); err != nil {
		return err
	}
	gamemgr.InternalStartServer()
	SendMessageToStatusChannel("🕛Start command received, Server is Starting...")
	return nil
}

func handleStop(s *discordgo.Session, i *discordgo.InteractionCreate, data EmbedData) error {
	data.Title, data.Description, data.Color = "Server Control", "Stopping the server...", 0xFF0000
	data.Fields = []EmbedField{{Name: "Status", Value: "🕛 Recieved", Inline: true}}
	if err := respond(s, i, data); err != nil {
		return err
	}
	gamemgr.InternalStopServer()
	SendMessageToStatusChannel("🕛Stop command received, flatlining Server in 5 Seconds...")
	return nil
}

func handleStatus(s *discordgo.Session, i *discordgo.InteractionCreate, data EmbedData) error {
	isRunning := gamemgr.InternalIsServerRunning()
	data.Title = "🎮 Server Status"
	data.Description = "Current process state for the Stationeers game server.\n*Note: 'Started' indicates a running process was found, but not necessarily fully operational.*"
	data.Color = map[bool]int{true: 0x00FF00, false: 0xFF0000}[isRunning]
	data.Fields = []EmbedField{
		{Name: "Status:", Value: map[bool]string{true: "🟢 Started", false: "🔴 Stopped"}[isRunning], Inline: true},
		{Name: "Checked:", Value: time.Now().Format("15:04:05 MST"), Inline: true},
	}
	return respond(s, i, data)
}

func handleHelp(s *discordgo.Session, i *discordgo.InteractionCreate, data EmbedData) error {
	data.Title, data.Description, data.Color = "Command Help", "Available Commands:", 0x1E90FF
	data.Fields = []EmbedField{
		{Name: "/start", Value: "Starts the server"},
		{Name: "/stop", Value: "Stops the server"},
		{Name: "/restore <index>", Value: "Restores a backup DEPRECATED will be re-added in a future release"},
		{Name: "/list [limit]", Value: "Lists recent backups (default: 5) DEPRECATED will be re-added in a future release"},
		{Name: "/help", Value: "Shows this help"},
	}
	return respond(s, i, data)
}

func handleRestore(s *discordgo.Session, i *discordgo.InteractionCreate, data EmbedData) error {
	data.Title = "Backup Manager Deprecated"
	data.Description = "The backup manager v2 has been removed and will be replaced with a more robust backup system in a future release."
	data.Color = 0xFF0000
	data.Fields = []EmbedField{{Name: "Status", Value: "Feature Deprecated", Inline: true}}
	return respond(s, i, data)
}

func handleList(s *discordgo.Session, i *discordgo.InteractionCreate, data EmbedData) error {
	data.Title = "Backup Manager Deprecated"
	data.Description = "The backup manager v2 has been removed and will be replaced with a more robust backup system in a future release."
	data.Color = 0xFF0000
	data.Fields = []EmbedField{{Name: "Status", Value: "Feature Deprecated", Inline: true}}
	return respond(s, i, data)
}
