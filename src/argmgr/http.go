package argmgr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// ArgUpdateRequest represents a request to update one or more arguments
type ArgUpdateRequest struct {
	Args map[string]string `json:"args"`
}

// ArgUpdateResponse represents the response to an arg update operation
type ArgUpdateResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Updated map[string]string `json:"updated,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// HandleArgUpdate processes HTTP requests to update arguments in a game runfile
func HandleArgUpdate(w http.ResponseWriter, r *http.Request, gameTemplate *GameTemplate) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}
	defer r.Body.Close()

	// Parse JSON request
	var updateReq ArgUpdateRequest
	if err := json.Unmarshal(body, &updateReq); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if len(updateReq.Args) == 0 {
		sendErrorResponse(w, http.StatusBadRequest, "No arguments provided to update")
		return
	}

	// Process each argument update
	response := ArgUpdateResponse{
		Success: true,
		Updated: make(map[string]string),
		Errors:  make(map[string]string),
	}

	for argName, argValue := range updateReq.Args {
		if err := SetArgValue(gameTemplate, argName, argValue); err != nil {
			response.Errors[argName] = err.Error()
			response.Success = false
		} else {
			response.Updated[argName] = argValue
		}
	}

	// Set appropriate message based on results
	if !response.Success {
		if len(response.Updated) > 0 {
			response.Message = "Some arguments were updated but others failed"
		} else {
			response.Message = "Failed to update any arguments"
		}
	} else {
		response.Message = "All arguments updated successfully"
	}

	// If no arguments were successfully updated, clean up the map
	if len(response.Updated) == 0 {
		response.Updated = nil
	}
	if len(response.Errors) == 0 {
		response.Errors = nil
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	if !response.Success {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(response)
}

// SaveRunfileHandler handles saving the updated run file to disk
func SaveRunfileHandler(w http.ResponseWriter, r *http.Request, gameTemplate *GameTemplate, runFilesFolder string) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	gameName := ""
	if meta, ok := gameTemplate.Meta["name"].(string); ok {
		gameName = meta
	} else {
		sendErrorResponse(w, http.StatusInternalServerError, "Game name not found in runfile")
		return
	}

	if err := SaveGameTemplate(gameTemplate, gameName, runFilesFolder); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to save run file: "+err.Error())
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Run file saved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SaveGameTemplate saves the game runfile back to the run file
func SaveGameTemplate(runfile *GameTemplate, gameName, runFilesFolder string) error {
	if gameName == "" {
		return errors.New("game name cannot be empty")
	}

	filePath := filepath.Join(runFilesFolder, fmt.Sprintf("run%s.ssui", gameName))

	// Convert to JSON
	data, err := json.MarshalIndent(runfile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize runfile: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write run file: %w", err)
	}

	return nil
}

// Helper function to send error responses
func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"message": message,
	})
}
