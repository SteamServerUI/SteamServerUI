package misc

import (
	"fmt"
	"runtime"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
)

// PrintStartupMessage prints a stylish startup message to the terminal
func PrintStartupMessage(backendEndpointUrl string) {
	// Clear some space
	fmt.Println("\n\n")

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

	// Quote
	fmt.Println("\n  \"JacksonTheMaster: Managing game servers shouldn't be rocket science... unless it's a rocket game!\"")

	// End with some space
	fmt.Println("\n\n")
	logger.Core.Info("Ready to run your server!")
	logger.Core.Info("🙏Thank you for using SSUI!")
}
