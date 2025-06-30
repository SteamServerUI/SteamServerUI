package web

import (
	"io/fs"
	"net/http"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/detectionmgr"
	"github.com/SteamServerUI/SteamServerUI/v6/src/settings"
)

// EndpointInfo contains metadata about each endpoint
type EndpointInfo struct {
	Path        string   `json:"path"`
	Method      string   `json:"method"`
	AccessLevel []string `json:"accessLevel"`
	Description string   `json:"description"`
	Protected   bool     `json:"protected"`
}

// endpointInfoMap stores detailed information about each endpoint
var endpointInfoMap = make(map[string]EndpointInfo)

// SetupRoutes configures the HTTP route handlers for the application, returning the main (unprotected) and protected (auth-required) ServeMux instances.
func SetupRoutes() (*http.ServeMux, *http.ServeMux) {
	// Main mux for unprotected routes
	mux := http.NewServeMux()

	// Protected mux for routes requiring authentication
	protectedMux := http.NewServeMux()

	// --- Static Assets ---
	// Frontend JS, CSS, and static files
	mux.HandleFunc("/sscm/sscm.js", ServeSSCMJs)
	endpointInfoMap["/sscm/sscm.js"] = EndpointInfo{Path: "/sscm/sscm.js", Method: "GET", AccessLevel: []string{}, Description: "Serves SSCM JavaScript file", Protected: false}

	legacyAssetsFS, _ := fs.Sub(config.GetV1UIFS(), "UIMod/onboard_bundled/v1")
	protectedMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(legacyAssetsFS))))
	endpointInfoMap["/static/"] = EndpointInfo{Path: "/static/", Method: "GET", AccessLevel: []string{}, Description: "Serves legacy UI static assets", Protected: true}

	twoboxformAssetsFS, _ := fs.Sub(config.GetTWOBOXFS(), "UIMod/onboard_bundled/twoboxform")
	mux.Handle("/twoboxform/", http.StripPrefix("/twoboxform/", http.FileServer(http.FS(twoboxformAssetsFS))))

	endpointInfoMap["/twoboxform/"] = EndpointInfo{Path: "/twoboxform/", Method: "GET", AccessLevel: []string{}, Description: "Serves two-box form assets", Protected: false}

	// --- Authentication Routes ---
	// Login, logout, user management, and setup
	mux.HandleFunc("/auth/login", LoginHandler) // Token issuer
	endpointInfoMap["/auth/login"] = EndpointInfo{Path: "/auth/login", Method: "POST", AccessLevel: []string{}, Description: "Authenticates user and issues auth token", Protected: false}

	mux.HandleFunc("/auth/logout", LogoutHandler)
	endpointInfoMap["/auth/logout"] = EndpointInfo{Path: "/auth/logout", Method: "GET", AccessLevel: []string{}, Description: "Logs out current user", Protected: false}

	protectedMux.HandleFunc("/api/v2/auth/adduser", accessLevelMiddleware(RegisterUserHandler, "superadmin")) // User setup and change password
	endpointInfoMap["/api/v2/auth/adduser"] = EndpointInfo{Path: "/api/v2/auth/adduser", Method: "POST", AccessLevel: []string{"superadmin"}, Description: "Adds new user or changes password", Protected: true}

	protectedMux.HandleFunc("/api/v2/auth/setup/finalize", SetupFinalizeHandler)
	endpointInfoMap["/api/v2/auth/setup/finalize"] = EndpointInfo{Path: "/api/v2/auth/setup/finalize", Method: "POST", AccessLevel: []string{}, Description: "Finalizes first-time setup", Protected: true}

	protectedMux.HandleFunc("/api/v2/auth/whoami", WhoAmIHandler)
	endpointInfoMap["/api/v2/auth/whoami"] = EndpointInfo{Path: "/api/v2/auth/whoami", Method: "GET", AccessLevel: []string{}, Description: "Returns current user information", Protected: true}

	mux.HandleFunc("/login", ServeTwoBoxFormTemplate)
	endpointInfoMap["/login"] = EndpointInfo{Path: "/login", Method: "GET", AccessLevel: []string{}, Description: "Serves login page", Protected: false}

	protectedMux.HandleFunc("/setup", ServeTwoBoxFormTemplate)
	endpointInfoMap["/setup"] = EndpointInfo{Path: "/setup", Method: "GET", AccessLevel: []string{}, Description: "Serves setup page", Protected: true}

	// --- Server Control ---
	// Game server start/stop/status
	protectedMux.HandleFunc("/start", StartServer) // Legacy start endpoint
	endpointInfoMap["/start"] = EndpointInfo{Path: "/start", Method: "GET", AccessLevel: []string{}, Description: "Legacy endpoint to start game server", Protected: true}

	protectedMux.HandleFunc("/stop", StopServer) // Legacy stop endpoint
	endpointInfoMap["/stop"] = EndpointInfo{Path: "/stop", Method: "GET", AccessLevel: []string{}, Description: "Legacy endpoint to stop game server", Protected: true}

	protectedMux.HandleFunc("/api/v2/server/start", accessLevelMiddleware(StartServer, "superadmin", "user"))
	endpointInfoMap["/api/v2/server/start"] = EndpointInfo{Path: "/api/v2/server/start", Method: "GET", AccessLevel: []string{"superadmin", "user"}, Description: "Starts the game server", Protected: true}

	protectedMux.HandleFunc("/api/v2/server/stop", accessLevelMiddleware(StopServer, "superadmin", "user"))
	endpointInfoMap["/api/v2/server/stop"] = EndpointInfo{Path: "/api/v2/server/stop", Method: "GET", AccessLevel: []string{"superadmin", "user"}, Description: "Stops the game server", Protected: true}

	protectedMux.HandleFunc("/api/v2/server/status", GetGameServerRunState)
	endpointInfoMap["/api/v2/server/status"] = EndpointInfo{Path: "/api/v2/server/status", Method: "GET", AccessLevel: []string{}, Description: "Returns current game server status", Protected: true}

	// --- Configuration ---
	// Config pages, saving configs, runfile args, and SSCM command execution
	protectedMux.HandleFunc("/config", ServeConfigPage)
	endpointInfoMap["/config"] = EndpointInfo{Path: "/config", Method: "GET", AccessLevel: []string{}, Description: "Serves configuration page", Protected: true}

	protectedMux.HandleFunc("/api/v2/settings/save", accessLevelMiddleware(settings.SaveSetting, "superadmin"))
	endpointInfoMap["/api/v2/settings/save"] = EndpointInfo{Path: "/api/v2/settings/save", Method: "POST", AccessLevel: []string{"superadmin"}, Description: "Saves application settings", Protected: true}

	protectedMux.HandleFunc("/api/v2/settings", settings.RetrieveSettings)
	endpointInfoMap["/api/v2/settings"] = EndpointInfo{Path: "/api/v2/settings", Method: "GET", AccessLevel: []string{}, Description: "Retrieves application settings", Protected: true}

	protectedMux.HandleFunc("/api/v2/SSCM/enabled", HandleIsSSCMEnabled) // Check if SSCM is enabled (responds with 200 OK if enabled, 403 Forbidden if disabled)
	endpointInfoMap["/api/v2/SSCM/enabled"] = EndpointInfo{Path: "/api/v2/SSCM/enabled", Method: "GET", AccessLevel: []string{}, Description: "Checks if SSCM is enabled", Protected: true}

	protectedMux.HandleFunc("/api/v2/runfile/groups", HandleRunfileGroups)
	endpointInfoMap["/api/v2/runfile/groups"] = EndpointInfo{Path: "/api/v2/runfile/groups", Method: "GET", AccessLevel: []string{}, Description: "Returns runfile groups", Protected: true}

	protectedMux.HandleFunc("/api/v2/runfile/args", HandleRunfileArgs)
	endpointInfoMap["/api/v2/runfile/args"] = EndpointInfo{Path: "/api/v2/runfile/args", Method: "GET", AccessLevel: []string{}, Description: "Returns runfile arguments", Protected: true}

	protectedMux.HandleFunc("/api/v2/runfile/args/update", HandleRunfileArgUpdate)
	endpointInfoMap["/api/v2/runfile/args/update"] = EndpointInfo{Path: "/api/v2/runfile/args/update", Method: "POST", AccessLevel: []string{}, Description: "Updates runfile arguments", Protected: true}

	protectedMux.HandleFunc("/api/v2/runfile", HandleRunfile)
	endpointInfoMap["/api/v2/runfile"] = EndpointInfo{Path: "/api/v2/runfile", Method: "GET", AccessLevel: []string{}, Description: "Returns runfile content", Protected: true}

	protectedMux.HandleFunc("/api/v2/runfile/save", HandleRunfileSave)
	endpointInfoMap["/api/v2/runfile/save"] = EndpointInfo{Path: "/api/v2/runfile/save", Method: "POST", AccessLevel: []string{}, Description: "Saves runfile content", Protected: true}

	protectedMux.HandleFunc("/api/v2/runfile/hardreset", HandleSetRunfileGame)
	endpointInfoMap["/api/v2/runfile/hardreset"] = EndpointInfo{Path: "/api/v2/runfile/hardreset", Method: "POST", AccessLevel: []string{}, Description: "Performs hard reset of runfile", Protected: true}

	// --- Loader ---
	protectedMux.HandleFunc("/api/v2/loader/reloadall", HandleReloadAll)
	endpointInfoMap["/api/v2/loader/reloadall"] = EndpointInfo{Path: "/api/v2/loader/reloadall", Method: "GET", AccessLevel: []string{}, Description: "Reloads all configuration", Protected: true}

	protectedMux.HandleFunc("/api/v2/loader/reloadconfig", HandleReloadConfig)
	endpointInfoMap["/api/v2/loader/reloadconfig"] = EndpointInfo{Path: "/api/v2/loader/reloadconfig", Method: "GET", AccessLevel: []string{}, Description: "Reloads configuration only", Protected: true}

	protectedMux.HandleFunc("/api/v2/loader/reloadrunfile", HandleReloadRunfile)
	endpointInfoMap["/api/v2/loader/reloadrunfile"] = EndpointInfo{Path: "/api/v2/loader/reloadrunfile", Method: "GET", AccessLevel: []string{}, Description: "Reloads runfile only", Protected: true}

	protectedMux.HandleFunc("/api/v2/loader/restartbackend", accessLevelMiddleware(HandleRestartMySelf, "superadmin"))
	endpointInfoMap["/api/v2/loader/restartbackend"] = EndpointInfo{Path: "/api/v2/loader/restartbackend", Method: "GET", AccessLevel: []string{"superadmin"}, Description: "Restarts the backend service", Protected: true}

	// --- SSE/Events ---
	// Real-time console and event streaming
	protectedMux.HandleFunc("/console", GetLogOutput)
	endpointInfoMap["/console"] = EndpointInfo{Path: "/console", Method: "GET", AccessLevel: []string{}, Description: "Server-Sent Events stream for console output", Protected: true}

	protectedMux.HandleFunc("/events", GetEventOutput)
	endpointInfoMap["/events"] = EndpointInfo{Path: "/events", Method: "GET", AccessLevel: []string{}, Description: "Server-Sent Events stream for application events", Protected: true}

	protectedMux.HandleFunc("/logs/debug", GetDebugLogOutput)
	endpointInfoMap["/logs/debug"] = EndpointInfo{Path: "/logs/debug", Method: "GET", AccessLevel: []string{}, Description: "Server-Sent Events stream for debug logs", Protected: true}

	protectedMux.HandleFunc("/logs/info", GetInfoLogOutput)
	endpointInfoMap["/logs/info"] = EndpointInfo{Path: "/logs/info", Method: "GET", AccessLevel: []string{}, Description: "Server-Sent Events stream for info logs", Protected: true}

	protectedMux.HandleFunc("/logs/warn", GetWarnLogOutput)
	endpointInfoMap["/logs/warn"] = EndpointInfo{Path: "/logs/warn", Method: "GET", AccessLevel: []string{}, Description: "Server-Sent Events stream for warning logs", Protected: true}

	protectedMux.HandleFunc("/logs/error", GetErrorLogOutput)
	endpointInfoMap["/logs/error"] = EndpointInfo{Path: "/logs/error", Method: "GET", AccessLevel: []string{}, Description: "Server-Sent Events stream for error logs", Protected: true}

	protectedMux.HandleFunc("/logs/backend", GetBackendLogOutput)
	endpointInfoMap["/logs/backend"] = EndpointInfo{Path: "/logs/backend", Method: "GET", AccessLevel: []string{}, Description: "Server-Sent Events stream for backend logs", Protected: true}

	// --- Custom Detections ---
	// Custom detection management
	protectedMux.HandleFunc("/api/v2/custom-detections", detectionmgr.HandleCustomDetection)
	endpointInfoMap["/api/v2/custom-detections"] = EndpointInfo{Path: "/api/v2/custom-detections", Method: "GET,POST", AccessLevel: []string{}, Description: "Manages custom detections (GET: list, POST: create)", Protected: true}

	protectedMux.HandleFunc("/api/v2/custom-detections/delete/", detectionmgr.HandleDeleteCustomDetection)
	endpointInfoMap["/api/v2/custom-detections/delete/"] = EndpointInfo{Path: "/api/v2/custom-detections/delete/", Method: "DELETE", AccessLevel: []string{}, Description: "Deletes a custom detection by ID", Protected: true}

	// --- SteamCMD ---
	// SteamCMD execution
	protectedMux.HandleFunc("/api/v2/steamcmd/run", accessLevelMiddleware(HandleRunSteamCMD, "superadmin"))
	endpointInfoMap["/api/v2/steamcmd/run"] = EndpointInfo{Path: "/api/v2/steamcmd/run", Method: "POST", AccessLevel: []string{"superadmin"}, Description: "Executes SteamCMD commands", Protected: true}

	// --- SVELTE ASSETS ---
	svelteAssetsFS, _ := fs.Sub(config.V2UIFS, "UIMod/onboard_bundled/v2/assets")
	protectedMux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(svelteAssetsFS))))
	endpointInfoMap["/assets/"] = EndpointInfo{Path: "/assets/", Method: "GET", AccessLevel: []string{}, Description: "Serves Svelte UI assets", Protected: true}

	// --- UI Pages ---
	// Main pages for the UI
	protectedMux.HandleFunc("/", ServeSvelteUI)
	endpointInfoMap["/"] = EndpointInfo{Path: "/", Method: "GET", AccessLevel: []string{}, Description: "Serves main Svelte UI application", Protected: true}

	protectedMux.HandleFunc("/v1", ServeIndex)
	endpointInfoMap["/v1"] = EndpointInfo{Path: "/v1", Method: "GET", AccessLevel: []string{}, Description: "Serves legacy v1 UI", Protected: true}

	protectedMux.HandleFunc("/detectionmanager", ServeDetectionManager)
	endpointInfoMap["/detectionmanager"] = EndpointInfo{Path: "/detectionmanager", Method: "GET", AccessLevel: []string{}, Description: "Serves detection manager page", Protected: true}

	// --- OS STATS ---
	protectedMux.HandleFunc("/api/v2/osstats", HandleGetOsStats)
	endpointInfoMap["/api/v2/osstats"] = EndpointInfo{Path: "/api/v2/osstats", Method: "GET", AccessLevel: []string{}, Description: "Returns operating system statistics", Protected: true}

	// --- RUNFILE GALLERY ---
	protectedMux.HandleFunc("/api/v2/gallery", galleryHandler)
	endpointInfoMap["/api/v2/gallery"] = EndpointInfo{Path: "/api/v2/gallery", Method: "GET", AccessLevel: []string{}, Description: "Returns available runfile gallery items", Protected: true}

	protectedMux.HandleFunc("/api/v2/gallery/select", accessLevelMiddleware(selectHandler, "superadmin"))
	endpointInfoMap["/api/v2/gallery/select"] = EndpointInfo{Path: "/api/v2/gallery/select", Method: "GET", AccessLevel: []string{"superadmin"}, Description: "Selects a runfile from the gallery", Protected: true}

	// --- CODE SERVER ---
	protectedMux.HandleFunc("/api/v2/codeserver/", accessLevelMiddleware(HandleCodeServer, "superadmin"))
	endpointInfoMap["/api/v2/codeserver/"] = EndpointInfo{Path: "/api/v2/codeserver/", Method: "GET,POST", AccessLevel: []string{"superadmin"}, Description: "Code server functionality", Protected: true}

	protectedMux.HandleFunc("/api/v2/getwd", accessLevelMiddleware(HandleGetWorkingDir, "superadmin"))
	endpointInfoMap["/api/v2/getwd"] = EndpointInfo{Path: "/api/v2/getwd", Method: "GET", AccessLevel: []string{"superadmin"}, Description: "Returns current working directory", Protected: true}

	// --- BACKUP ---
	protectedMux.HandleFunc("/api/v2/backup/create", HandleBackupCreate)
	endpointInfoMap["/api/v2/backup/create"] = EndpointInfo{Path: "/api/v2/backup/create", Method: "POST", AccessLevel: []string{}, Description: "Creates a new backup", Protected: true}

	protectedMux.HandleFunc("/api/v2/backup/list", HandleBackupList)
	endpointInfoMap["/api/v2/backup/list"] = EndpointInfo{Path: "/api/v2/backup/list", Method: "GET", AccessLevel: []string{}, Description: "Lists available backups", Protected: true}

	protectedMux.HandleFunc("/api/v2/backup/restore", HandleBackupRestore)
	endpointInfoMap["/api/v2/backup/restore"] = EndpointInfo{Path: "/api/v2/backup/restore", Method: "POST", AccessLevel: []string{}, Description: "Restores from a backup", Protected: true}

	protectedMux.HandleFunc("/api/v2/backup/status", HandleBackupStatus)
	endpointInfoMap["/api/v2/backup/status"] = EndpointInfo{Path: "/api/v2/backup/status", Method: "GET", AccessLevel: []string{}, Description: "Returns backup operation status", Protected: true}

	// --- API Endpoints List ---
	protectedMux.HandleFunc("/api/v2/endpoints", HandleListEndpoints)
	endpointInfoMap["/api/v2/endpoints"] = EndpointInfo{Path: "/api/v2/endpoints", Method: "GET", AccessLevel: []string{}, Description: "Lists all available API endpoints with metadata", Protected: true}

	return mux, protectedMux
}
