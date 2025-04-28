package web

import (
	"net/http"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/backupmgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/detectionmgr"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/settings"
)

// SetupRoutes configures the HTTP route handlers for the application, returning the main (unprotected) and protected (auth-required) ServeMux instances.
func SetupRoutes() (*http.ServeMux, *http.ServeMux) {
	// Main mux for unprotected routes
	mux := http.NewServeMux()

	// Protected mux for routes requiring authentication
	protectedMux := http.NewServeMux()

	// --- Static Assets ---
	// Frontend JS, CSS, and static files
	mux.HandleFunc("/twoboxform/twoboxform.js", ServeTwoBoxJs)
	mux.HandleFunc("/twoboxform/twoboxform.css", ServeTwoBoxCss)
	mux.HandleFunc("/sscm/sscm.js", ServeSSCMJs)
	fs := http.FileServer(http.Dir(config.GetUIModFolder() + "/assets"))
	protectedMux.Handle("/static/", http.StripPrefix("/static/", fs))

	// --- Authentication Routes ---
	// Login, logout, user management, and setup
	mux.HandleFunc("/auth/login", LoginHandler) // Token issuer
	mux.HandleFunc("/auth/logout", LogoutHandler)
	mux.HandleFunc("/login", ServeTwoBoxFormTemplate)
	protectedMux.HandleFunc("/changeuser", ServeTwoBoxFormTemplate)
	protectedMux.HandleFunc("/api/v2/auth/adduser", RegisterUserHandler) // User setup and change password
	protectedMux.HandleFunc("/setup", ServeTwoBoxFormTemplate)
	protectedMux.HandleFunc("/api/v2/auth/setup/register", RegisterUserHandler) // User registration
	protectedMux.HandleFunc("/api/v2/auth/setup/finalize", SetupFinalizeHandler)

	// --- Server Control ---
	// Game server start/stop/status
	protectedMux.HandleFunc("/start", StartServer) // Legacy start endpoint
	protectedMux.HandleFunc("/stop", StopServer)   // Legacy stop endpoint
	protectedMux.HandleFunc("/api/v2/server/start", StartServer)
	protectedMux.HandleFunc("/api/v2/server/stop", StopServer)
	protectedMux.HandleFunc("/api/v2/server/status", GetGameServerRunState)

	// --- Configuration ---
	// Config pages, saving configs, runfile args, and SSCM command execution
	protectedMux.HandleFunc("/config", ServeConfigPage)
	protectedMux.HandleFunc("/api/v2/settings/save", settings.SaveSetting)
	protectedMux.HandleFunc("/api/v2/settings", settings.RetrieveSettings)
	protectedMux.HandleFunc("/api/v2/SSCM/run", HandleCommand)           // Command execution via SSCM (gated by config.GetIsSSCMEnabled())
	protectedMux.HandleFunc("/api/v2/SSCM/enabled", HandleIsSSCMEnabled) // Check if SSCM is enabled (responds with 200 OK if enabled, 403 Forbidden if disabled)
	protectedMux.HandleFunc("/api/v2/runfile/groups", HandleRunfileGroups)
	protectedMux.HandleFunc("/api/v2/runfile/args", HandleRunfileArgs)
	protectedMux.HandleFunc("/api/v2/runfile/args/update", HandleRunfileArgUpdate)
	protectedMux.HandleFunc("/api/v2/runfile", HandleRunfile)
	protectedMux.HandleFunc("/api/v2/runfile/save", HandleRunfileSave)
	protectedMux.HandleFunc("/api/v2/runfile/hardreset", HandleSetRunfileGame)

	// --- Backups ---
	// Backup listing and restoration
	backupHandler := backupmgr.NewHTTPHandler(backupmgr.GlobalBackupManager)
	protectedMux.HandleFunc("/api/v2/backups", backupHandler.ListBackupsHandler)
	protectedMux.HandleFunc("/api/v2/backups/restore", backupHandler.RestoreBackupHandler)

	// --- SSE/Events ---
	// Real-time console and event streaming
	protectedMux.HandleFunc("/console", GetLogOutput)
	protectedMux.HandleFunc("/events", GetEventOutput)

	// --- Custom Detections ---
	// Custom detection management
	protectedMux.HandleFunc("/api/v2/custom-detections", detectionmgr.HandleCustomDetection)
	protectedMux.HandleFunc("/api/v2/custom-detections/delete/", detectionmgr.HandleDeleteCustomDetection)

	// --- SteamCMD ---
	// SteamCMD execution
	protectedMux.HandleFunc("/api/v2/steamcmd/run", HandleRunSteamCMD)

	// --- SVELTE ASSETS ---
	svelteAssets := http.FileServer(http.Dir(config.GetUIModFolder() + "/v2/assets"))
	protectedMux.Handle("/assets/", http.StripPrefix("/assets/", svelteAssets))

	// --- UI Pages ---
	// Main pages for the UI
	protectedMux.HandleFunc("/", ServeIndex)
	protectedMux.HandleFunc("/detectionmanager", ServeDetectionManager)
	protectedMux.HandleFunc("/v2", ServeSvelteUI)

	return mux, protectedMux
}
