package backupmgr

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
)

var (
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	isRunning bool
	mu        sync.Mutex

	// Config variables - will be replaced with config calls later
	backupContentDir   = "./backupContentDir"
	storedBackupsDir   = "./storedBackupsDir"
	backupLoopInterval = 5 * time.Minute
	backupMode         = "zip"
)

func InitBackupMgr() {
	logger.Backup.Debug("Initializing Backup Manager")

	// Create directories if they don't exist
	if err := ensureDirectories(); err != nil {
		logger.Backup.Error("Failed to create directories: " + err.Error())
		return
	}

	// Start the backup loop
	StartBackupLoop()

	logger.Backup.Info("Backup Manager Initialized")
	logger.Backup.Info("Content Directory: " + backupContentDir)
	logger.Backup.Info("Backup Directory: " + storedBackupsDir)
	logger.Backup.Info("Backup Interval: " + backupLoopInterval.String())
	logger.Backup.Info("Backup Mode: " + backupMode)
}

func ensureDirectories() error {
	dirs := []string{backupContentDir, storedBackupsDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

func StartBackupLoop() {
	mu.Lock()
	defer mu.Unlock()

	if isRunning {
		logger.Backup.Warn("Backup loop is already running")
		return
	}

	ctx, cancel = context.WithCancel(context.Background())
	isRunning = true
	wg.Add(1)

	go backupLoop()
	logger.Backup.Info("Backup loop started")
}

func StopBackupLoop() {
	mu.Lock()
	wasRunning := isRunning
	isRunning = false
	mu.Unlock()

	if !wasRunning {
		return
	}

	cancel()
	wg.Wait()
	logger.Backup.Info("Backup loop stopped")
}

func backupLoop() {
	defer wg.Done()

	ticker := time.NewTicker(backupLoopInterval)
	defer ticker.Stop()

	// Perform initial backup
	CreateBackup()

	for {
		select {
		case <-ctx.Done():
			logger.Backup.Debug("Backup loop cancelled")
			return
		case <-ticker.C:
			CreateBackup()
		}
	}
}

func CreateBackup() error {
	logger.Backup.Debug("Starting backup operation")

	// Check if content directory exists and has content
	if !hasContent() {
		logger.Backup.Debug("No content to backup, skipping")
		return nil
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupFilename := fmt.Sprintf("backup_%s.zip", timestamp)
	backupPath := filepath.Join(storedBackupsDir, backupFilename)

	// Perform the backup based on mode
	var err error
	switch backupMode {
	case "zip":
		err = createZipBackup(backupPath)
	default:
		err = fmt.Errorf("unsupported backup mode: %s", backupMode)
	}

	if err != nil {
		logger.Backup.Error("Backup failed: " + err.Error())
		return err
	}

	logger.Backup.Info("Backup completed successfully: " + backupFilename)
	return nil
}

func hasContent() bool {
	entries, err := os.ReadDir(backupContentDir)
	if err != nil {
		logger.Backup.Warn("Unable to read content directory: " + err.Error())
		return false
	}
	return len(entries) > 0
}

func createZipBackup(backupPath string) error {
	// Create the zip file
	zipFile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	// Create zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through the content directory and add files to zip
	err = filepath.Walk(backupContentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == backupContentDir {
			return nil
		}

		// Get relative path for zip entry
		relPath, err := filepath.Rel(backupContentDir, path)
		if err != nil {
			return err
		}

		// Normalize path separators for zip
		relPath = strings.ReplaceAll(relPath, "\\", "/")

		if info.IsDir() {
			// Create directory entry in zip
			_, err := zipWriter.Create(relPath + "/")
			return err
		}

		// Create file entry in zip
		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		// Open source file
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// Copy file content to zip
		_, err = io.Copy(zipEntry, srcFile)
		return err
	})

	return err
}

func GetBackupList() ([]string, error) {
	entries, err := os.ReadDir(storedBackupsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".zip") {
			backups = append(backups, entry.Name())
		}
	}

	return backups, nil
}

func RestoreBackup(backupName string) error {
	logger.Backup.Info("Starting restore operation for: " + backupName)

	// Validate backup file exists
	backupPath := filepath.Join(storedBackupsDir, backupName)
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backupName)
	}

	// Create temporary directory for extraction
	tempDir := backupContentDir + "_temp_restore"
	if err := os.RemoveAll(tempDir); err != nil {
		return fmt.Errorf("failed to clean temp directory: %w", err)
	}

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Extract backup to temp directory
	if err := extractZipBackup(backupPath, tempDir); err != nil {
		os.RemoveAll(tempDir) // Clean up on error
		return fmt.Errorf("failed to extract backup: %w", err)
	}

	// Backup current content (if exists) before replacing
	backupCurrentDir := backupContentDir + "_backup_before_restore"
	if _, err := os.Stat(backupContentDir); err == nil {
		if err := os.RemoveAll(backupCurrentDir); err != nil {
			logger.Backup.Warn("Failed to remove old backup directory: " + err.Error())
		}
		if err := os.Rename(backupContentDir, backupCurrentDir); err != nil {
			os.RemoveAll(tempDir) // Clean up on error
			return fmt.Errorf("failed to backup current content: %w", err)
		}
		logger.Backup.Info("Current content backed up to: " + backupCurrentDir)
	}

	// Move extracted content to final location
	if err := os.Rename(tempDir, backupContentDir); err != nil {
		// Try to restore the backup if move failed
		if _, backupErr := os.Stat(backupCurrentDir); backupErr == nil {
			os.Rename(backupCurrentDir, backupContentDir)
		}
		os.RemoveAll(tempDir) // Clean up
		return fmt.Errorf("failed to restore content: %w", err)
	}

	logger.Backup.Info("Backup restored successfully: " + backupName)
	return nil
}

func extractZipBackup(backupPath, destDir string) error {
	// Open zip file
	zipReader, err := zip.OpenReader(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer zipReader.Close()

	// Extract each file
	for _, file := range zipReader.File {
		// Create full path for extraction
		destPath := filepath.Join(destDir, file.Name)

		// Ensure the destination directory exists
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Skip if it's a directory
		if file.FileInfo().IsDir() {
			continue
		}

		// Open file in zip
		srcFile, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip: %w", err)
		}

		// Create destination file
		destFile, err := os.Create(destPath)
		if err != nil {
			srcFile.Close()
			return fmt.Errorf("failed to create destination file: %w", err)
		}

		// Copy content
		_, err = io.Copy(destFile, srcFile)
		srcFile.Close()
		destFile.Close()

		if err != nil {
			return fmt.Errorf("failed to copy file content: %w", err)
		}
	}

	return nil
}

func IsBackupRunning() bool {
	mu.Lock()
	defer mu.Unlock()
	return isRunning
}
