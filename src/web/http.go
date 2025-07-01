package web

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"text/template"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/gamemgr"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v6/src/steammgr"
)

/*
The http and http2 file are monolithic and will eventually be refactored into smaller, more descriptive files Soon™
http.go WAS used for everything added BEFORE the Svelte UI was added.

DO NOT ADD NEW HTTP ENDPOINTS TO THIS FILE, PLEASE use simpler descriptive files instead.

*/

// TemplateData holds data to be passed to templates
type TemplateData struct {
	Version string
	Branch  string
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {

	htmlFS, err := fs.Sub(config.V1UIFS, "UIMod/onboard_bundled/v1/ui")
	if err != nil {
		http.Error(w, "Error accessing Virt FS: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFS(htmlFS, "index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Core.Error("failed to serve v1 Index.html")
		return
	}

	data := TemplateData{
		Version: config.Version,
		Branch:  config.Branch,
	}
	if data.Version == "" {
		data.Version = "unknown"
	}
	if data.Branch == "" {
		data.Branch = "unknown"
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ServeDetectionManager(w http.ResponseWriter, r *http.Request) {
	htmlFS, err := fs.Sub(config.V1UIFS, "UIMod/onboard_bundled/v1/ui")
	if err != nil {
		http.Error(w, "Error accessing Virt FS: "+err.Error(), http.StatusInternalServerError)
		return
	}

	htmlFile, err := htmlFS.Open("detectionmanager.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading detectionmanager.html: %v", err), http.StatusInternalServerError)
		return
	}
	defer htmlFile.Close()

	htmlContent, err := io.ReadAll(htmlFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading detectionmanager.html content: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(htmlContent)
}

func ServeSvelteUI(w http.ResponseWriter, r *http.Request) {
	htmlFS, err := fs.Sub(config.V2UIFS, "UIMod/onboard_bundled/v2")
	if err != nil {
		http.Error(w, "Error accessing Svelte UI: "+err.Error(), http.StatusInternalServerError)
		return
	}

	htmlFile, err := htmlFS.Open("index.html")
	if err != nil {
		http.Error(w, "Error reading Svelte UI: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer htmlFile.Close()

	// Stream the file content to the response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = io.Copy(w, htmlFile)
	if err != nil {
		http.Error(w, "Error writing Svelte UI: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func ServeConfigPage(w http.ResponseWriter, r *http.Request) {
	htmlFS, err := fs.Sub(config.V1UIFS, "UIMod/onboard_bundled/v1/ui")
	if err != nil {
		http.Error(w, "Error accessing Virt FS: "+err.Error(), http.StatusInternalServerError)
		return
	}

	htmlFile, err := htmlFS.Open("config.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading config.html: %v", err), http.StatusInternalServerError)
		return
	}
	defer htmlFile.Close()

	htmlContentBytes, err := io.ReadAll(htmlFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading config.html content: %v", err), http.StatusInternalServerError)
		return
	}

	htmlContent := string(htmlContentBytes)

	// Determine selected attributes for boolean fields
	discordTrueSelected := ""
	discordFalseSelected := ""
	if config.GetIsDiscordEnabled() {
		discordTrueSelected = "selected"
	} else {
		discordFalseSelected = "selected"
	}

	// Replace placeholders in the HTML with actual config values
	replacements := map[string]string{
		"{{discordToken}}":                  config.GetDiscordToken(),
		"{{controlChannelID}}":              config.GetControlChannelID(),
		"{{statusChannelID}}":               config.GetStatusChannelID(),
		"{{connectionListChannelID}}":       config.GetConnectionListChannelID(),
		"{{logChannelID}}":                  config.GetLogChannelID(),
		"{{saveChannelID}}":                 config.GetSaveChannelID(),
		"{{controlPanelChannelID}}":         config.GetControlPanelChannelID(),
		"{{blackListFilePath}}":             config.GetBlackListFilePath(),
		"{{errorChannelID}}":                config.GetErrorChannelID(),
		"{{isDiscordEnabled}}":              fmt.Sprintf("%v", config.GetIsDiscordEnabled()),
		"{{IsDiscordEnabledTrueSelected}}":  discordTrueSelected,
		"{{IsDiscordEnabledFalseSelected}}": discordFalseSelected,
		"{{gameBranch}}":                    config.GetGameBranch(),
	}

	for placeholder, value := range replacements {
		htmlContent = strings.ReplaceAll(htmlContent, placeholder, value)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, htmlContent)
}

// StartServer HTTP handler
func StartServer(w http.ResponseWriter, r *http.Request) {
	logger.Web.Debug("Received start request from API")
	if err := gamemgr.InternalStartServer(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Web.Info("Error starting server: " + err.Error())
		return
	}
	fmt.Fprint(w, "Server started.")
	logger.Web.Info("Server started.")
}

// StopServer HTTP handler
func StopServer(w http.ResponseWriter, r *http.Request) {
	logger.Web.Debug("Received stop request from API")
	if err := gamemgr.InternalStopServer(); err != nil {
		if err.Error() == "server not running" {
			fmt.Fprint(w, "Server was not running or was already stopped")
			logger.Web.Info("Server not running or was already stopped")
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Web.Info("Error stopping server: " + err.Error())
		return
	}
	fmt.Fprint(w, "Server stopped.")
	logger.Web.Info("Server stopped.")
}

func GetGameServerRunState(w http.ResponseWriter, r *http.Request) {
	runState := gamemgr.InternalIsServerRunning()
	response := map[string]interface{}{
		"isRunning": runState,
		"uuid":      config.GetGameServerUUID().String(),
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to respond with Game Server status", http.StatusInternalServerError)
		return
	}
}

func ServeTwoBoxCss(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	http.ServeFile(w, r, config.GetUIModFolder()+"twoboxform/twoboxform.css")
}

func ServeTwoBoxJs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	http.ServeFile(w, r, config.GetUIModFolder()+"twoboxform/twoboxform.js")
}

func ServeSSCMJs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	http.ServeFile(w, r, config.GetSSCMWebDir()+"sscm.js")
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

func HandleRunSteamCMD(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	steammgr.RunSteamCMD()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "SteamCMD run started"})

}
