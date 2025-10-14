package sscmapi

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/commandmgr"
)

// CommandHandler handles POST requests to execute commands via commandmgr.
// Expects a command in the request body. Returns 204 on success or error details.
func CommandHandler(w http.ResponseWriter, r *http.Request) {
	// Allow only POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read command from request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	command := strings.TrimSpace(string(body))

	// Validate command
	if command == "" {
		http.Error(w, "Command cannot be empty", http.StatusBadRequest)
		return
	}

	// Execute command via commandmgr
	if err := commandmgr.WriteCommand(command); err != nil {
		switch err {
		case os.ErrNotExist:
			http.Error(w, "Command file path not configured", http.StatusInternalServerError)
		case os.ErrInvalid:
			http.Error(w, "Invalid command", http.StatusBadRequest)
		default:
			http.Error(w, "Failed to write command: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Success: return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

func HandleIsSSCMEnabled(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if SSCM is enabled
	if !config.GetIsSSCMEnabled() {
		http.Error(w, "SSCM is disabled", http.StatusForbidden)
		return
	}

	// Success: return 200 OK
	w.WriteHeader(http.StatusOK)
}
