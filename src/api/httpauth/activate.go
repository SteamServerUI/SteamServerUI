package httpauth

import (
	"encoding/json"
	"net/http"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/config/configchanger"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/loader"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

// ActivateAuthHandler (ex. SetupFinalizeHandler) marks setup as complete
func ActivateAuthHandler(w http.ResponseWriter, r *http.Request) {

	//check if users map is nil or empty
	if len(config.GetUsers()) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "No users registered - cannot finalize setup at this time. You should really enable authentication - or click 'Skip authentication'"})
		return
	}

	// Load existing config to update it
	newConfig, err := config.LoadConfig()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error - Failed to load config"})
		return
	}

	// Mark setup as complete and enable auth
	config.SetIsFirstTimeSetup(false)
	isTrue := true
	newConfig.AuthEnabled = &isTrue // Set the pointer to true

	// Save the updated config
	err = configchanger.SaveConfig(newConfig)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error - Failed to save config"})
		return
	}

	logger.Web.Info("User Setup finalized successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":      "Setup finalized successfully",
		"restart_hint": "You will be redirected to the login page...",
	})
	loader.ReloadBackend()
}
