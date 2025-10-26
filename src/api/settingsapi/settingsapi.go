package settingsapi

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/google/uuid"
)

// uploadFile is a reusable function to handle file uploads
func uploadFile(w http.ResponseWriter, r *http.Request, targetDir, filename string, allowedExts []string) {
	// Ensure the target directory exists
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Failed to create directory"})
		return
	}

	// Get the file from the form
	file, header, err := r.FormFile("file")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Error reading file"})
		return
	}
	defer file.Close()

	// Validate file extension if provided
	if len(allowedExts) > 0 {
		ext := strings.ToLower(filepath.Ext(header.Filename))
		valid := false
		for _, allowed := range allowedExts {
			if ext == allowed {
				valid = true
				break
			}
		}
		if !valid {
			json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Invalid file extension"})
			return
		}
	}

	// Determine the filename
	finalFilename := filename
	if finalFilename == "" {
		finalFilename = header.Filename
		if finalFilename == "" {
			finalFilename = uuid.NewString()
		}
	}

	// Create the output file
	filePath := filepath.Join(targetDir, finalFilename)
	out, err := os.Create(filePath)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Error creating file"})
		return
	}
	defer out.Close()

	// Copy the file content
	if _, err := io.Copy(out, file); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Error writing file"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "File uploaded successfully"})
}

// HandleTLSCertUpload handles uploads for cert.pem and key.pem
func HandleTLSCertUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := r.URL.Query().Get("type")
	if filename != "cert.pem" && filename != "key.pem" {
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Invalid file type. Must be cert.pem or key.pem"})
		return
	}

	targetDir := filepath.Join(config.GetSSUIFolder(), "tls")
	uploadFile(w, r, targetDir, filename, []string{".pem"})
}

// HandleFileUpload handles general file uploads
func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	targetDir := filepath.Join(config.GetSSUIFolder(), "config", "files")
	uploadFile(w, r, targetDir, "", nil) // Empty filename means use original or generate UUID
}

// HandleBackgroundUpload handles dashboard-background.png upload
func HandleBackgroundUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	targetDir := filepath.Join(config.GetSSUIFolder(), "config", "files")
	uploadFile(w, r, targetDir, "dashboard-background.png", []string{".png"})
}
