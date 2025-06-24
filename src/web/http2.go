package web

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/SteamServerUI/SteamServerUI/v6/src/loader"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
)

/*
The http and http2 file are monolithic and will eventually be refactored into smaller, more descriptive files Soon™
http2.go is used for everything added after the Svelte UI was added.
*/

var reloadMu sync.Mutex // Mutex to prevent concurrent reloads, atleast to some degree. This is a temporary duct tape solution that will probably be here for a long time. Lol.

// HandleSetRunfileGame reloads the runfile and restarts most of the server. It can also be used to reload the runfile from Disk as a hard reset.
func HandleSetRunfileGame(w http.ResponseWriter, r *http.Request) {

	reloadMu.Lock()
	defer reloadMu.Unlock()

	// Restrict to POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read and validate request body
	var request struct {
		Game string `json:"game"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	game := strings.TrimSpace(request.Game)
	if game == "" {
		http.Error(w, "Game cannot be empty", http.StatusBadRequest)
		return
	}

	// Call InitRunfile to handle the runfile update
	if err := loader.InitRunfile(game); err != nil {
		logger.Core.Error("Failed to initialize runfile: " + err.Error())
		http.Error(w, "Failed to initialize runfile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Core.Info("Runfile game updated successfully to " + game)

	// Prepare response
	response := struct {
		Message string `json:"message"`
		Game    string `json:"game"`
	}{
		Message: "Monitor console for update status",
		Game:    game,
	}

	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func HandleReloadAll(w http.ResponseWriter, r *http.Request) {
	logger.Web.Debug("Received reloadall request from API")
	reloadMu.Lock()
	defer reloadMu.Unlock()
	// accept only GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	// Reload all loaders
	loader.ReloadAll()

	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "OK"}); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func HandleReloadConfig(w http.ResponseWriter, r *http.Request) {
	logger.Web.Debug("Received reloadconfig request from API")
	reloadMu.Lock()
	defer reloadMu.Unlock()
	// accept only GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	loader.ReloadConfig()

	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "OK"}); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func HandleReloadRunfile(w http.ResponseWriter, r *http.Request) {
	logger.Web.Debug("Received reloadrunfile request from API")
	reloadMu.Lock()
	defer reloadMu.Unlock()
	// accept only GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	loader.ReloadRunfile()

	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "OK"}); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func HandleGetWorkingDir(w http.ResponseWriter, r *http.Request) {
	logger.Web.Debug("Received GetWorkingDir request from API")
	reloadMu.Lock()
	defer reloadMu.Unlock()
	// accept only GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get current directory
	wd, err := os.Getwd()
	if err != nil {
		logger.Core.Error("Failed to get current directory: " + err.Error())
		http.Error(w, "Failed to get current directory: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "OK", "WorkingDir": wd}); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
