package pages

import (
	"io"
	"io/fs"
	"net/http"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
)

func ServeSvelteUI(w http.ResponseWriter, r *http.Request) {
	htmlFS, err := fs.Sub(config.V1UIFS, "SSUI/onboard_bundled/v2")
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
