package misc

import (
	"fmt"
	"runtime"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
)

// PrintStartupMessage prints a stylish startup message to the terminal
func PrintStartupMessage(backendEndpointUrl string) {
	// Clear some space
	fmt.Println()
	fmt.Println()

	// Main ASCII art logo
	fmt.Println("  ███████╗████████╗███████╗ █████╗ ███╗   ███╗███████╗███████╗██████╗ ██╗   ██╗███████╗██████╗ ██╗   ██╗██╗")
	fmt.Println("  ██╔════╝╚══██╔══╝██╔════╝██╔══██╗████╗ ████║██╔════╝██╔════╝██╔══██╗██║   ██║██╔════╝██╔══██╗██║   ██║██║")
	fmt.Println("  ███████╗   ██║   █████╗  ███████║██╔████╔██║███████╗█████╗  ██████╔╝██║   ██║█████╗  ██████╔╝██║   ██║██║")
	fmt.Println("  ╚════██║   ██║   ██╔══╝  ██╔══██║██║╚██╔╝██║╚════██║██╔══╝  ██╔══██╗╚██╗ ██╔╝██╔══╝  ██╔══██╗██║   ██║██║")
	fmt.Println("  ███████║   ██║   ███████╗██║  ██║██║ ╚═╝ ██║███████║███████╗██║  ██║ ╚████╔╝ ███████╗██║  ██║╚██████╔╝██║")
	fmt.Println("  ╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝╚══════╝╚═╝  ╚═╝  ╚═══╝  ╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝")

	// Decorative line
	fmt.Println("  ╔═══════════════════════════════════════════════════════════════════════════════════════════════════╗")
	// Tagline
	fmt.Println("  ║                      🎮 YOUR ONE-STOP SHOP FOR STEAM SERVER SHENANIGANS 🎮                        ║")
	// System info
	fmt.Printf("  ║  🚀 Version: %s       📅 %s       💻 Runtime: %s/%s                     ║\n",
		config.GetBackendVersion(),
		time.Now().Format("2006-01-02 15:04:05"),
		runtime.GOOS,
		runtime.GOARCH)
	// Decorative line
	fmt.Println("  ╚═══════════════════════════════════════════════════════════════════════════════════════════════════╝")

	// Web UI info
	fmt.Println("\n  🌐 Web UI available at: https://localhost:8443 (default) or " + backendEndpointUrl)
	fmt.Println("\n  🌐 Support available at: https://discord.gg/8n3vN92MyJ")

	// Quote
	fmt.Println("\n  JacksonTheMaster: \"Managing game servers shouldn't be rocket science... unless it's a rocket game!\"")
}

func PrintFirstTimeSetupMessage() {
	// Setup guide
	fmt.Println("  📋 GETTING STARTED:")
	fmt.Println("  ┌─────────────────────────────────────────────────────────────────────────────────────────────┐")
	fmt.Println("  │ • Ready, set, go! Welcome to SteamServerUI, new User!                                       │")
	fmt.Println("  │ • The good news: you made it here, which means you are likely ready to run your server!      │")
	fmt.Println("  │ • If this is your first time here, no worries: SSUI is made to be easy to use.              │")
	fmt.Println("  │ • Select a game server from the Runfile Gallery to get started                              │")
	fmt.Println("  │ • Configure your server settings in the Settings panel                                      │")
	fmt.Println("  │ • Support is provided at https://discord.gg/8n3vN92MyJ                                      │")
	fmt.Println("  │ • Review the Documentation at https://steamserverui.github.io/SteamServerUI/                │")
	fmt.Println("  └─────────────────────────────────────────────────────────────────────────────────────────────┘")
}
