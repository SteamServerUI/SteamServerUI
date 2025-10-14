package api

import (
	"io/fs"
	"net/http"

	"github.com/SteamServerUI/SteamServerUI/v7/src/api/httpauth"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/legacyapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/pages"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/runfileapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/sscmapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/sseapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/sysinfoapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/config/configchanger"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/backupmgr"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/detectionmgr"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamserverui/settings"
)

func SetupRoutes() (*http.ServeMux, *http.ServeMux) {

	// Set up handlers with auth middleware
	mux := http.NewServeMux() // Use a mux to apply middleware globally

	// Unprotected auth routes
	twoboxformAssetsFS, _ := fs.Sub(config.GetV1UIFS(), "SSUI/onboard_bundled/twoboxform")
	mux.Handle("/twoboxform/", http.StripPrefix("/twoboxform/", http.FileServer(http.FS(twoboxformAssetsFS))))
	mux.HandleFunc("/auth/login", httpauth.LoginHandler) // Token issuer
	mux.HandleFunc("/auth/logout", httpauth.LogoutHandler)
	mux.HandleFunc("/login", pages.ServeTwoBoxFormTemplate)

	// Protected routes (wrapped with middleware)
	protectedMux := http.NewServeMux()

	legacyAssetsFS, _ := fs.Sub(config.GetV1UIFS(), "SSUI/onboard_bundled/assets")
	protectedMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(legacyAssetsFS))))

	protectedMux.HandleFunc("/config", pages.ServeConfigPage)
	protectedMux.HandleFunc("/detectionmanager", pages.ServeDetectionManager)
	protectedMux.HandleFunc("/", pages.ServeIndex)

	// --- SVELTE UI ---
	protectedMux.HandleFunc("/v2", pages.ServeSvelteUI)
	svelteAssetsFS, _ := fs.Sub(config.V1UIFS, "SSUI/onboard_bundled/v2/assets")
	protectedMux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(svelteAssetsFS))))
	protectedMux.HandleFunc("/api/v2/loader/reloadbackend", HandleReloadBackend)

	// SSE routes
	protectedMux.HandleFunc("/console", sseapi.GetLogOutput)
	protectedMux.HandleFunc("/events", sseapi.GetEventOutput)
	protectedMux.HandleFunc("/logs/debug", sseapi.GetDebugLogOutput)
	protectedMux.HandleFunc("/logs/info", sseapi.GetInfoLogOutput)
	protectedMux.HandleFunc("/logs/warn", sseapi.GetWarnLogOutput)
	protectedMux.HandleFunc("/logs/error", sseapi.GetErrorLogOutput)
	protectedMux.HandleFunc("/logs/backend", sseapi.GetBackendLogOutput)

	// Server Control
	protectedMux.HandleFunc("/start", legacyapi.StartServer)
	protectedMux.HandleFunc("/stop", legacyapi.StopServer)
	protectedMux.HandleFunc("/api/v2/server/start", legacyapi.StartServer) // TODO: should return json & get their own functions
	protectedMux.HandleFunc("/api/v2/server/stop", legacyapi.StopServer)   // TODO: should return json & get their own functions
	protectedMux.HandleFunc("/api/v2/server/status", GetGameServerRunState)
	protectedMux.HandleFunc("/api/v2/server/status/connectedplayers", legacyapi.HandleConnectedPlayersList)

	backupHandler := backupmgr.NewHTTPHandler(backupmgr.GlobalBackupManager)
	protectedMux.HandleFunc("/api/v2/backups", backupHandler.ListBackupsHandler)
	protectedMux.HandleFunc("/api/v2/backups/restore", backupHandler.RestoreBackupHandler)

	// Configuration
	protectedMux.HandleFunc("/saveconfigasjson", configchanger.SaveConfigForm)     // legacy, used on config page
	protectedMux.HandleFunc("/api/v2/saveconfig", configchanger.SaveConfigRestful) // used on twoboxform
	protectedMux.HandleFunc("/api/v2/SSCM/run", sscmapi.HandleCommand)             // Command execution via SSCM (needs to be enable, config.IsSSCMEnabled)
	protectedMux.HandleFunc("/api/v2/SSCM/enabled", sscmapi.HandleIsSSCMEnabled)   // Check if SSCM is enabled
	protectedMux.HandleFunc("/api/v2/steamcmd/run", HandleRunSteamCMD)             // Run SteamCMD

	// Custom Detections
	protectedMux.HandleFunc("/api/v2/custom-detections", detectionmgr.HandleCustomDetection)
	protectedMux.HandleFunc("/api/v2/custom-detections/delete/", detectionmgr.HandleDeleteCustomDetection)
	// Authentication
	protectedMux.HandleFunc("/changeuser", pages.ServeTwoBoxFormTemplate)
	protectedMux.HandleFunc("/api/v2/auth/adduser", httpauth.RegisterUserHandler)        // user registration and change password
	protectedMux.HandleFunc("/api/v2/auth/setup/apikey", httpauth.RegisterAPIKeyHandler) // apikey registration and change password
	protectedMux.HandleFunc("/api/v2/auth/whoami", httpauth.WhoAmIHandler)

	// Setup
	protectedMux.HandleFunc("/setup", pages.ServeTwoBoxFormTemplate)
	protectedMux.HandleFunc("/api/v2/auth/setup/register", httpauth.RegisterUserHandler) // user registration
	protectedMux.HandleFunc("/api/v2/auth/setup/finalize", httpauth.ActivateAuthHandler)

	// SteamServerUI

	// --- RUNFILE ---
	protectedMux.HandleFunc("/api/v2/runfile/groups", runfileapi.HandleRunfileGroups)
	protectedMux.HandleFunc("/api/v2/runfile/args", runfileapi.HandleRunfileArgs)
	protectedMux.HandleFunc("/api/v2/runfile/args/update", runfileapi.HandleRunfileArgUpdate)
	protectedMux.HandleFunc("/api/v2/runfile", runfileapi.HandleRunfile)
	protectedMux.HandleFunc("/api/v2/runfile/save", runfileapi.HandleRunfileSave)
	protectedMux.HandleFunc("/api/v2/runfile/hardreset", runfileapi.HandleSetRunfileGame)
	// --- LOADER ---
	protectedMux.HandleFunc("/api/v2/loader/reloadrunfile", runfileapi.HandleReloadRunfile)
	// --- SETTINGS ---
	protectedMux.HandleFunc("/api/v2/settings/save", settings.SaveSetting)
	protectedMux.HandleFunc("/api/v2/settings", settings.RetrieveSettings)
	// --- OS STATS ---
	protectedMux.HandleFunc("/api/v2/osstats", sysinfoapi.HandleGetOsStats)
	// --- RUNFILE GALLERY ---
	protectedMux.HandleFunc("/api/v2/gallery", runfileapi.GalleryHandler)
	protectedMux.HandleFunc("/api/v2/gallery/select", runfileapi.GallerySelectHandler)

	return mux, protectedMux
}
