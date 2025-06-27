package web

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/SteamServerUI/SteamServerUI/v6/src/backupmgr"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
)

type BackupResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type BackupListResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Backups []string `json:"backups"`
}

type BackupStatusResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	IsRunning bool   `json:"isRunning"`
}

type RestoreRequest struct {
	BackupName    string `json:"backupName"`
	SkipPreBackup bool   `json:"skipPreBackup"`
}

type BackupCreateRequest struct {
	Mode string `json:"mode"`
}

func HandleBackupCreate(w http.ResponseWriter, r *http.Request) {
	logger.Web.Debug("API: Create backup requested")

	//accept only POST requests
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		respondBackupError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req BackupCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Web.Error("API: Failed to parse backup create request: " + err.Error())
		respondBackupError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate backup mode
	if req.Mode != "copy" && req.Mode != "tar" && req.Mode != "zip" {
		respondBackupError(w, "Invalid backup mode", http.StatusBadRequest)
		return
	}

	logger.Web.Info("API: Creating backup with mode: " + req.Mode)

	// Trigger backup creation
	err := backupmgr.CreateBackup(req.Mode)
	if err != nil {
		logger.Web.Error("API: Failed to create backup: " + err.Error())
		respondBackupError(w, "Failed to create backup: "+err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Web.Info("API: Backup created successfully")
	respondBackupSuccess(w, "Backup created successfully", nil)
}

func HandleBackupList(w http.ResponseWriter, r *http.Request) {
	logger.Web.Debug("API: Backup list requested")

	if r.Method != http.MethodGet {
		respondBackupError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get list of backups
	backups, err := backupmgr.GetBackupList()
	if err != nil {
		logger.Web.Error("API: Failed to get backup list: " + err.Error())
		respondBackupError(w, "Failed to get backup list: "+err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Web.Debug("API: Retrieved backup list successfully")
	respondBackupList(w, "Backup list retrieved successfully", backups)
}

func HandleBackupRestore(w http.ResponseWriter, r *http.Request) {
	logger.Web.Debug("API: Backup restore requested")

	if r.Method != http.MethodPost {
		respondBackupError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req RestoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Web.Error("API: Failed to parse restore request: " + err.Error())
		respondBackupError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate backup name
	if strings.TrimSpace(req.BackupName) == "" {
		respondBackupError(w, "Backup name is required", http.StatusBadRequest)
		return
	}

	// Validate backup name format (basic security check)
	if strings.Contains(req.BackupName, "..") || strings.Contains(req.BackupName, "/") || strings.Contains(req.BackupName, "\\") {
		respondBackupError(w, "Invalid backup name", http.StatusBadRequest)
		return
	}

	logger.Web.Info("API: Restoring backup: " + req.BackupName)
	if req.SkipPreBackup {
		logger.Web.Info("API: Skipping pre-restore backup")
	}

	// Perform restore
	err := backupmgr.RestoreBackup(req.BackupName, req.SkipPreBackup)
	if err != nil {
		logger.Web.Error("API: Failed to restore backup: " + err.Error())
		respondBackupError(w, "Failed to restore backup: "+err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Web.Info("API: Backup restored successfully: " + req.BackupName)
	respondBackupSuccess(w, "Backup restored successfully: "+req.BackupName, nil)
}

func HandleBackupStatus(w http.ResponseWriter, r *http.Request) {
	logger.Web.Debug("API: Backup status requested")

	if r.Method != http.MethodGet {
		respondBackupError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get backup manager status
	isRunning := backupmgr.IsBackupRunning()

	logger.Web.Debug("API: Backup status retrieved successfully")
	respondBackupStatus(w, "Backup status retrieved", isRunning)
}

// Helper functions for consistent API responses
func respondBackupSuccess(w http.ResponseWriter, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := BackupResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func respondBackupError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := BackupResponse{
		Success: false,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

func respondBackupList(w http.ResponseWriter, message string, backups []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := BackupListResponse{
		Success: true,
		Message: message,
		Backups: backups,
	}

	json.NewEncoder(w).Encode(response)
}

func respondBackupStatus(w http.ResponseWriter, message string, isRunning bool) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := BackupStatusResponse{
		Success:   true,
		Message:   message,
		IsRunning: isRunning,
	}

	json.NewEncoder(w).Encode(response)
}
