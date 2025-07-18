// Package main is the entry point for SteamServerUI, a tool for managing Stationeers servers.
// It coordinates setup, configuration, logging, resource loading, and a web-based UI.
//
// The server initializes by running the setup process, loading resources, and starting a web server.
// Key functionality is provided by the following subpackages:
//   - src/config: Manages server configuration.
//   - src/configchanger: Handles configuration changes.
//   - src/gamemgr: Manages process management.
//   - src/loader: Handles resource loading and detection initialization.
//   - src/logger: Provides logging utilities.
//   - src/setup: Performs initial server setup.
//   - src/web: Runs the web-based user interface.
//   - src/discordbot: Handles Discord bot functionality.
//   - src/backupmgr: Manages backups of the server's world.
//   - src/detectionmgr: Manages event detection and processing.
//   - src/setup: Performs initial server setup.
//   - src/ssestream: Handles Server-Sent Events (SSE) streaming.
//   - src/security: Handles security-related tasks.
//
// For detailed documentation, see the subpackages or the project Wiki on GitHub.
package main

import (
	"embed"
	"sync"

	"github.com/SteamServerUI/SteamServerUI/v6/src/loader"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v6/src/misc"
	"github.com/SteamServerUI/SteamServerUI/v6/src/setup"
	"github.com/SteamServerUI/SteamServerUI/v6/src/web"
)

// Bundled Assets

//go:embed UIMod/onboard_bundled/v1
var v1uiFS embed.FS

//go:embed UIMod/onboard_bundled/v2
var v2uiFS embed.FS

//go:embed UIMod/onboard_bundled/twoboxform
var twoboxFS embed.FS

func main() {
	var wg sync.WaitGroup

	setup.V6setupMutex.Lock()
	setup.IsSetupComplete = false
	setup.V6setupMutex.Unlock()

	logger.Install.Info("Starting setup...")
	loader.ReloadConfig() // Load the config file before starting the setup process
	misc.LoadCmdArgs()
	// Start the installation process and wait for it to complete before starting the rest of the server
	wg.Add(1)
	go setup.Install(&wg)
	wg.Wait()

	// Load config,discordbot, backupmgr and detectionmgr using the loader package
	loader.InitVirtFS(v1uiFS, v2uiFS, twoboxFS)
	loader.ReloadAll()
	loader.InitDetector()
	loader.AfterStartComplete()

	misc.StartConsole(&wg)

	web.StartWebServer(&wg)
}
