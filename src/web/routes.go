package web

import (
	"net/http"

	"github.com/SteamServerUI/SteamServerUI/v6/src/backupmgr"
	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/detectionmgr"
	"github.com/SteamServerUI/SteamServerUI/v6/src/settings"
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
	fs := http.FileServer(http.Dir(config.GetUIModFolder() + "/v1"))
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
	protectedMux.HandleFunc("/api/v2/SSCM/enabled", HandleIsSSCMEnabled) // Check if SSCM is enabled (responds with 200 OK if enabled, 403 Forbidden if disabled)
	protectedMux.HandleFunc("/api/v2/runfile/groups", HandleRunfileGroups)
	protectedMux.HandleFunc("/api/v2/runfile/args", HandleRunfileArgs)
	protectedMux.HandleFunc("/api/v2/runfile/args/update", HandleRunfileArgUpdate)
	protectedMux.HandleFunc("/api/v2/runfile", HandleRunfile)
	protectedMux.HandleFunc("/api/v2/runfile/save", HandleRunfileSave)
	protectedMux.HandleFunc("/api/v2/runfile/hardreset", HandleSetRunfileGame)

	// --- Loader ---
	protectedMux.HandleFunc("/api/v2/loader/reloadall", HandleReloadAll)
	protectedMux.HandleFunc("/api/v2/loader/reloadconfig", HandleReloadConfig)
	protectedMux.HandleFunc("/api/v2/loader/reloadrunfile", HandleReloadRunfile)

	// --- Backups ---
	// Backup listing and restoration
	backupHandler := backupmgr.NewHTTPHandler(backupmgr.GlobalBackupManager)
	protectedMux.HandleFunc("/api/v2/backups", backupHandler.ListBackupsHandler)
	protectedMux.HandleFunc("/api/v2/backups/restore", backupHandler.RestoreBackupHandler)

	// --- SSE/Events ---
	// Real-time console and event streaming
	protectedMux.HandleFunc("/console", GetLogOutput)
	protectedMux.HandleFunc("/events", GetEventOutput)
	protectedMux.HandleFunc("/logs/debug", GetDebugLogOutput)
	protectedMux.HandleFunc("/logs/info", GetInfoLogOutput)
	protectedMux.HandleFunc("/logs/warn", GetWarnLogOutput)
	protectedMux.HandleFunc("/logs/error", GetErrorLogOutput)
	protectedMux.HandleFunc("/logs/backend", GetBackendLogOutput)

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
	protectedMux.HandleFunc("/", ServeSvelteUI)
	protectedMux.HandleFunc("/v1", ServeIndex)
	protectedMux.HandleFunc("/detectionmanager", ServeDetectionManager)

	// --- OS STATS ---
	protectedMux.HandleFunc("/api/v2/osstats", HandleGetOsStats)

	// --- RUNFILE GALLERY ---
	protectedMux.HandleFunc("/api/v2/gallery", galleryHandler)
	protectedMux.HandleFunc("/api/v2/gallery/select", selectHandler)

	// --- CODE SERVER ---
	protectedMux.HandleFunc("/api/v2/codeserver", HandleCodeServer)
	return mux, protectedMux
}
