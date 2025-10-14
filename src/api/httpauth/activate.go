package httpauth

import (
	"encoding/json"
	"net/http"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
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

	// Mark setup as complete and enable auth
	err := config.SetIsFirstTimeSetup(false)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error - Failed to SetIsFirstTimeSetup in config"})
		return
	}
	err = config.SetAuthEnabled(true)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error - Failed to SetAuthEnabled in config"})
		return
	}

	loader.ReloadBackend()

	logger.Web.Info("User Setup finalized successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":      "Setup finalized successfully",
		"restart_hint": "You will be redirected to the login page...",
	})
}
