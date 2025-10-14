package legacyapi

import (
	"fmt"
	"net/http"

	"github.com/SteamServerUI/SteamServerUI/v7/src/localization"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/detectionmgr"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/gamemgr"
)

// StartServer HTTP handler
func StartServer(w http.ResponseWriter, r *http.Request) {
	logger.API.Debug("Received start request from API")
	if err := gamemgr.InternalStartServer(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.API.Error("Error starting server: " + err.Error())
		return
	}
	fmt.Fprint(w, localization.GetString("BackendText_ServerStarted"))
	logger.API.Info("Server started.")
}

// StopServer HTTP handler
func StopServer(w http.ResponseWriter, r *http.Request) {
	logger.API.Debug("Received stop request from API")
	if err := gamemgr.InternalStopServer(); err != nil {
		if err.Error() == "server not running" {
			fmt.Fprint(w, localization.GetString("BackendText_ServerNotRunningOrAlreadyStopped"))
			logger.API.Warn("Server not running or was already stopped")
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.API.Error("Error stopping server: " + err.Error())
		return
	}
	detectionmgr.ClearPlayers(detectionmgr.GetDetector())
	fmt.Fprint(w, localization.GetString("BackendText_ServerStopped"))
	logger.API.Info("Server stopped.")
}
