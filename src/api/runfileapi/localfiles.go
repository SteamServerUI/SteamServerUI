package runfileapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamserverui/runfile"
)

// FileRequest represents the JSON body for file operations
type FileRequest struct {
	Filename string `json:"filename"`
	Content  string `json:"content,omitempty"` // Used for save operations
}

// FileResponse represents the JSON response structure
type FileResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// FileInfo represents the file information returned by GetFileList
type FileInfo struct {
	Filename    string `json:"filename"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// GetFileList handles GET requests to list all available files with their details
func GetFileList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendFileError(w, http.StatusMethodNotAllowed, "only GET requests are allowed")
		return
	}

	files := runfile.GetFiles()
	if files == nil {
		sendFileError(w, http.StatusNotFound, "no runfile loaded or no files available")
		return
	}

	fileInfos := make([]FileInfo, 0, len(files))
	for _, file := range files {
		// if filename is inside the SSUI subdirectory, dont add file to the list
		if strings.HasPrefix(file.Filepath, "./SSUI") {
			continue
		}

		// if file does not exist, dont add it to the list
		if _, err := os.Stat(file.Filepath); os.IsNotExist(err) {
			continue
		}

		fileInfos = append(fileInfos, FileInfo{
			Filename:    file.Filename,
			Type:        file.Type,
			Description: file.Description,
		})
	}

	sendFileResponse(w, http.StatusOK, FileResponse{
		Success: true,
		Data:    fileInfos,
	})
}

// GetFile handles GET requests to retrieve a specific file's contents
func GetFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendFileError(w, http.StatusMethodNotAllowed, "only GET requests are allowed")
		return
	}

	var req FileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendFileError(w, http.StatusBadRequest, fmt.Sprintf("failed to parse request JSON body: %v", err))
		return
	}

	if req.Filename == "" {
		sendFileError(w, http.StatusBadRequest, "filename is required")
		return
	}

	// Find the file in runfile
	files := runfile.GetFiles()
	var targetFile *runfile.File
	for _, file := range files {
		if file.Filename == req.Filename {
			targetFile = &file
			break
		}
	}

	if targetFile == nil {
		sendFileError(w, http.StatusNotFound, fmt.Sprintf("file %s not found in runfile", req.Filename))
		return
	}

	// Check if file is in SSUI subdirectory
	if strings.HasPrefix(targetFile.Filepath, "./SSUI") {
		sendFileError(w, http.StatusForbidden, fmt.Sprintf("access to a file %s in the SSUI subdirectory is forbidden", req.Filename))
		return
	}

	// Stat the file
	fileInfo, err := os.Stat(targetFile.Filepath)
	if err != nil {
		if os.IsNotExist(err) {
			sendFileError(w, http.StatusNotFound, fmt.Sprintf("file %s does not exist at %s", req.Filename, targetFile.Filepath))
		} else {
			sendFileError(w, http.StatusInternalServerError, fmt.Sprintf("failed to stat file %s: %v", req.Filename, err))
		}
		return
	}

	// Check if file is writable
	if fileInfo.Mode().Perm()&0222 == 0 {
		sendFileError(w, http.StatusForbidden, fmt.Sprintf("file %s is not writable", req.Filename))
		return
	}

	// Read file contents
	content, err := os.ReadFile(targetFile.Filepath)
	if err != nil {
		sendFileError(w, http.StatusInternalServerError, fmt.Sprintf("failed to read file %s: %v", req.Filename, err))
		return
	}

	// Set content type based on file type
	contentType := "text/plain"
	switch strings.ToLower(targetFile.Type) {
	case "json":
		contentType = "application/json"
	case "xml":
		contentType = "application/xml"
	case "yaml":
		contentType = "application/yaml"
	case "ini":
		contentType = "text/plain"
	}

	// Send file contents directly
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(content); err != nil {
		logger.Runfile.Error(fmt.Sprintf("failed to write file response for %s: %v", req.Filename, err))
	}
}

// SaveFile handles POST requests to save an edited file
func SaveFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendFileError(w, http.StatusMethodNotAllowed, "only POST requests are allowed")
		return
	}

	// Get filename from query parameter
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		sendFileError(w, http.StatusBadRequest, "filename query parameter is required")
		return
	}

	// Read raw content from request body
	content, err := io.ReadAll(r.Body)
	if err != nil {
		sendFileError(w, http.StatusBadRequest, fmt.Sprintf("failed to read request body: %v", err))
		return
	}
	if len(content) == 0 {
		sendFileError(w, http.StatusBadRequest, "content is required")
		return
	}

	// Find the file in runfile
	files := runfile.GetFiles()
	var targetFile *runfile.File
	for _, file := range files {
		if file.Filename == filename {
			targetFile = &file
			break
		}
	}

	if targetFile == nil {
		sendFileError(w, http.StatusNotFound, fmt.Sprintf("file %s not found in runfile", filename))
		return
	}

	// Check if file is in SSUI subdirectory
	if strings.HasPrefix(targetFile.Filepath, "./SSUI") {
		sendFileError(w, http.StatusForbidden, fmt.Sprintf("saving file %s in the SSUI subdirectory is forbidden", filename))
		return
	}

	// Stat the file
	fileInfo, err := os.Stat(targetFile.Filepath)
	if err != nil && !os.IsNotExist(err) {
		sendFileError(w, http.StatusInternalServerError, fmt.Sprintf("failed to stat file %s: %v", filename, err))
		return
	}

	// Check if file is writable (or would be writable if it exists)
	if fileInfo != nil && fileInfo.Mode().Perm()&0222 == 0 {
		sendFileError(w, http.StatusForbidden, fmt.Sprintf("file %s is not writable", filename))
		return
	}

	// Write file with retries
	const maxRetries = 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := os.WriteFile(targetFile.Filepath, content, 0644); err != nil {
			logger.Runfile.Warn(fmt.Sprintf("failed to write file %s: attempt=%d, error=%v", filename, attempt, err))
			if attempt == maxRetries {
				sendFileError(w, http.StatusInternalServerError, fmt.Sprintf("failed to write file %s after %d attempts: %v", filename, maxRetries, err))
				return
			}
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}

	sendFileResponse(w, http.StatusOK, FileResponse{
		Success: true,
		Message: fmt.Sprintf("file %s saved successfully", filename),
	})
}

// sendFileResponse sends a JSON response
func sendFileResponse(w http.ResponseWriter, status int, resp FileResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Runfile.Error(fmt.Sprintf("failed to encode response: %v", err))
	}
}

// sendFileError sends an error response as JSON
func sendFileError(w http.ResponseWriter, status int, message string) {
	logger.API.Debug(message)
	sendFileResponse(w, status, FileResponse{
		Success: false,
		Message: message,
	})
}
