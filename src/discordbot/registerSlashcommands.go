package discordbot

import (
	"sync"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"

	"github.com/bwmarrin/discordgo"
)

// registerSlashCommands defines and registers slash commands if they have not been registered already.
func registerSlashCommands(s *discordgo.Session) {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "start",
			Description: "Start the server",
		},
		{
			Name:        "stop",
			Description: "Stop the server",
		},
		{
			Name:        "status",
			Description: "Gets the running status of the gameserver process",
		},
		{
			Name:        "help",
			Description: "Show command help",
		},
		{
			Name:        "restore",
			Description: "Inop. For now, use backup manager in the web UI",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "index",
					Description: "Backup index to restore (e.g., 1)",
					Required:    true,
				},
			},
		},
		{
			Name:        "list",
			Description: "Deprecated, will be re-added in a future release",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "limit",
					Description: "Number of backups to list or 'all' (default: 5)",
					Required:    false,
				},
			},
		},
		{
			Name:        "backup",
			Description: "Creates a backup and saves it to the backups folder.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "mode",
					Description: "Backup Mode, optional",
					Required:    false,
				},
			},
		},
	}

	logger.Discord.Info("Checking and registering slash commands with Discord...")

	// Fetch existing commands from Discord
	existingCmds, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		logger.Discord.Error("Failed to fetch existing commands: " + err.Error())
		return
	}

	// Map existing commands by name for quick lookup
	existingMap := make(map[string]*discordgo.ApplicationCommand)
	for _, cmd := range existingCmds {
		existingMap[cmd.Name] = cmd
	}

	// Compare and register only what’s necessary
	var wg sync.WaitGroup
	commandsToRegister := make(chan *discordgo.ApplicationCommand, len(commands))

	for _, desiredCmd := range commands {
		existing, exists := existingMap[desiredCmd.Name]
		needsUpdate := !exists || !commandsAreEqual(desiredCmd, existing)

		if needsUpdate {
			wg.Add(1)
			commandsToRegister <- desiredCmd
		}

		logger.Discord.Debug("Command " + desiredCmd.Name + " already up-to-date, skipping")

	}
	close(commandsToRegister)

	// Worker to process api registrations
	go func() {
		for cmd := range commandsToRegister {
			startTime := time.Now()
			_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
			duration := time.Since(startTime)

			if err != nil {
				logger.Discord.Error("Error registering command " + cmd.Name + ": " + err.Error())
			}
			logger.Discord.Debug("Successfully registered command " + cmd.Name + " took:" + duration.String())
			wg.Done()
		}
	}()

	// Wait for all registrations to finish
	wg.Wait()
	logger.Discord.Info("Finished processing slash commands.")
}

// This is used to determine if a slash command needs to be registered with the discord server we are connected to or if it already exists.
// commandsAreEqual (helper) checks if two discrd commands are functionally identical
func commandsAreEqual(desired, existing *discordgo.ApplicationCommand) bool {
	if desired.Name != existing.Name || desired.Description != existing.Description {
		return false
	}

	// Compare options (nil vs empty slice handling)
	if len(desired.Options) != len(existing.Options) {
		return false
	}

	for i, desiredOpt := range desired.Options {
		existingOpt := existing.Options[i]
		if desiredOpt.Type != existingOpt.Type ||
			desiredOpt.Name != existingOpt.Name ||
			desiredOpt.Description != existingOpt.Description ||
			desiredOpt.Required != existingOpt.Required {
			return false
		}
	}

	return true
}
