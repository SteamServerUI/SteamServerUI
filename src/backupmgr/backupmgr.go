package backupmgr

import (
	"archive/tar"
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
	backupMode         = "zip"                     // "copy", "tar", "zip"
	maxFileSize        = int64(1024 * 1024 * 1024) // 1GB limit (configurable)
	useCompression     = true                      // Only compress if explicitly needed
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
	logger.Backup.Info("Max File Size: " + fmt.Sprintf("%d MB", maxFileSize/(1024*1024)))
	logger.Backup.Info("Use Compression: " + fmt.Sprintf("%t", useCompression))
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
	start := time.Now()
	logger.Backup.Debug("Starting backup operation")

	// Check if content directory exists and has content
	if !hasContent() {
		logger.Backup.Debug("No content to backup, skipping")
		return nil
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	var backupPath string
	var err error

	// Perform the backup based on mode
	switch backupMode {
	case "copy":
		backupPath = filepath.Join(storedBackupsDir, "backup_"+timestamp)
		err = createCopyBackup(backupPath)
	case "tar":
		backupPath = filepath.Join(storedBackupsDir, "backup_"+timestamp+".tar")
		err = createTarBackup(backupPath)
	case "zip":
		backupPath = filepath.Join(storedBackupsDir, "backup_"+timestamp+".zip")
		err = createZipBackup(backupPath)
	default:
		err = fmt.Errorf("unsupported backup mode: %s", backupMode)
	}

	if err != nil {
		logger.Backup.Error("Backup failed: " + err.Error())
		return err
	}

	duration := time.Since(start)
	backupName := filepath.Base(backupPath)
	logger.Backup.Info("Backup completed successfully: " + backupName + " (took " + duration.String() + ")")
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

// COPY MODE - Fastest for local storage, creates directory snapshot
func createCopyBackup(backupPath string) error {
	logger.Backup.Debug("Creating copy backup to: " + backupPath)

	// Create backup directory
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Copy all files and directories
	err := filepath.Walk(backupContentDir, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Backup.Warn("Skipping file due to error: " + srcPath + " - " + err.Error())
			return nil
		}

		// Skip the root directory itself
		if srcPath == backupContentDir {
			return nil
		}

		// Check file size limit
		if !info.IsDir() && info.Size() > maxFileSize {
			logger.Backup.Warn("Skipping large file: " + srcPath + " (size: " + fmt.Sprintf("%d MB", info.Size()/(1024*1024)) + ")")
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(backupContentDir, srcPath)
		if err != nil {
			logger.Backup.Warn("Failed to get relative path for: " + srcPath)
			return nil
		}

		destPath := filepath.Join(backupPath, relPath)

		if info.IsDir() {
			// Create directory
			return os.MkdirAll(destPath, info.Mode())
		}

		// Copy file
		return copyFile(srcPath, destPath, info.Mode())
	})

	return err
}

// TAR MODE - Fast, uncompressed archive
func createTarBackup(backupPath string) error {
	logger.Backup.Debug("Creating tar backup to: " + backupPath)

	// Create tar file
	tarFile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create tar file: %w", err)
	}
	defer tarFile.Close()

	// Create tar writer
	tarWriter := tar.NewWriter(tarFile)
	defer tarWriter.Close()

	// Walk and add files to tar
	err = filepath.Walk(backupContentDir, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Backup.Warn("Skipping file due to error: " + srcPath + " - " + err.Error())
			return nil
		}

		// Skip the root directory itself
		if srcPath == backupContentDir {
			return nil
		}

		// Check file size limit
		if !info.IsDir() && info.Size() > maxFileSize {
			logger.Backup.Warn("Skipping large file: " + srcPath + " (size: " + fmt.Sprintf("%d MB", info.Size()/(1024*1024)) + ")")
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(backupContentDir, srcPath)
		if err != nil {
			logger.Backup.Warn("Failed to get relative path for: " + srcPath)
			return nil
		}

		// Create tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			logger.Backup.Warn("Failed to create tar header for: " + srcPath)
			return nil
		}
		header.Name = strings.ReplaceAll(relPath, "\\", "/")

		// Write header
		if err := tarWriter.WriteHeader(header); err != nil {
			logger.Backup.Warn("Failed to write tar header for: " + srcPath)
			return nil
		}

		// Write file content if it's a regular file
		if info.Mode().IsRegular() {
			file, err := os.Open(srcPath)
			if err != nil {
				logger.Backup.Warn("Failed to open file: " + srcPath)
				return nil
			}
			defer file.Close()

			if _, err := io.Copy(tarWriter, file); err != nil {
				logger.Backup.Warn("Failed to copy file to tar: " + srcPath)
			}
		}

		return nil
	})

	return err
}

// ZIP MODE - Compressed, slower but smaller (kept for compatibility)
func createZipBackup(backupPath string) error {
	logger.Backup.Debug("Creating zip backup to: " + backupPath)

	zipFile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(backupContentDir, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Backup.Warn("Skipping file due to error: " + srcPath + " - " + err.Error())
			return nil
		}

		if srcPath == backupContentDir {
			return nil
		}

		// Check file size limit
		if !info.IsDir() && info.Size() > maxFileSize {
			logger.Backup.Warn("Skipping large file: " + srcPath + " (size: " + fmt.Sprintf("%d MB", info.Size()/(1024*1024)) + ")")
			return nil
		}

		relPath, err := filepath.Rel(backupContentDir, srcPath)
		if err != nil {
			logger.Backup.Warn("Failed to get relative path for: " + srcPath)
			return nil
		}

		relPath = strings.ReplaceAll(relPath, "\\", "/")

		if info.IsDir() {
			_, err := zipWriter.Create(relPath + "/")
			if err != nil {
				logger.Backup.Warn("Failed to create directory in zip: " + relPath)
			}
			return nil
		}

		// Create zip entry
		method := zip.Store // No compression by default
		if useCompression {
			method = zip.Deflate
		}

		zipEntry, err := zipWriter.CreateHeader(&zip.FileHeader{
			Name:   relPath,
			Method: method,
		})
		if err != nil {
			logger.Backup.Warn("Failed to create zip entry for: " + srcPath)
			return nil
		}

		file, err := os.Open(srcPath)
		if err != nil {
			logger.Backup.Warn("Failed to open file: " + srcPath)
			return nil
		}
		defer file.Close()

		_, err = io.Copy(zipEntry, file)
		if err != nil {
			logger.Backup.Warn("Failed to copy file to zip: " + srcPath)
		}

		return nil
	})

	return err
}

// Fast file copy using OS-level operations
func copyFile(src, dst string, mode os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination directory if needed
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy file content
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// Set file permissions
	return os.Chmod(dst, mode)
}

func GetBackupList() ([]string, error) {
	entries, err := os.ReadDir(storedBackupsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []string
	for _, entry := range entries {
		name := entry.Name()
		// Include all backup formats
		if strings.HasPrefix(name, "backup_") {
			backups = append(backups, name)
		}
	}

	return backups, nil
}

func RestoreBackup(backupName string) error {
	logger.Backup.Info("Starting restore operation for: " + backupName)

	backupPath := filepath.Join(storedBackupsDir, backupName)
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup not found: %s", backupName)
	}

	// Determine backup type and restore accordingly
	if strings.HasSuffix(backupName, ".zip") {
		return restoreZipBackup(backupPath)
	} else if strings.HasSuffix(backupName, ".tar") {
		return restoreTarBackup(backupPath)
	} else {
		// Assume it's a copy backup (directory)
		return restoreCopyBackup(backupPath)
	}
}

func restoreCopyBackup(backupPath string) error {
	return restoreWithSafety(func(tempDir string) error {
		// Copy backup directory to temp location
		return copyDirectory(backupPath, tempDir)
	})
}

func restoreTarBackup(backupPath string) error {
	return restoreWithSafety(func(tempDir string) error {
		return extractTarBackup(backupPath, tempDir)
	})
}

func restoreZipBackup(backupPath string) error {
	return restoreWithSafety(func(tempDir string) error {
		return extractZipBackup(backupPath, tempDir)
	})
}

func restoreWithSafety(extractFunc func(string) error) error {
	tempDir := backupContentDir + "_temp_restore"
	if err := os.RemoveAll(tempDir); err != nil {
		return fmt.Errorf("failed to clean temp directory: %w", err)
	}

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Extract to temp directory
	if err := extractFunc(tempDir); err != nil {
		os.RemoveAll(tempDir)
		return err
	}

	// Backup current content
	backupCurrentDir := backupContentDir + "_backup_before_restore"
	if _, err := os.Stat(backupContentDir); err == nil {
		if err := os.RemoveAll(backupCurrentDir); err != nil {
			logger.Backup.Warn("Failed to remove old backup directory: " + err.Error())
		}
		if err := os.Rename(backupContentDir, backupCurrentDir); err != nil {
			os.RemoveAll(tempDir)
			return fmt.Errorf("failed to backup current content: %w", err)
		}
		logger.Backup.Info("Current content backed up to: " + backupCurrentDir)
	}

	// Move extracted content to final location
	if err := os.Rename(tempDir, backupContentDir); err != nil {
		if _, backupErr := os.Stat(backupCurrentDir); backupErr == nil {
			os.Rename(backupCurrentDir, backupContentDir)
		}
		os.RemoveAll(tempDir)
		return fmt.Errorf("failed to restore content: %w", err)
	}

	logger.Backup.Info("Backup restored successfully")
	return nil
}

func copyDirectory(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return copyFile(path, destPath, info.Mode())
	})
}

func extractTarBackup(backupPath, destDir string) error {
	file, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open tar file: %w", err)
	}
	defer file.Close()

	tarReader := tar.NewReader(file)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		destPath := filepath.Join(destDir, header.Name)

		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(destPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		outFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}

		if _, err := io.Copy(outFile, tarReader); err != nil {
			outFile.Close()
			return fmt.Errorf("failed to extract file: %w", err)
		}

		outFile.Close()
		os.Chmod(destPath, os.FileMode(header.Mode))
	}

	return nil
}

func extractZipBackup(backupPath, destDir string) error {
	zipReader, err := zip.OpenReader(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		destPath := filepath.Join(destDir, file.Name)

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		if file.FileInfo().IsDir() {
			continue
		}

		srcFile, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip: %w", err)
		}

		destFile, err := os.Create(destPath)
		if err != nil {
			srcFile.Close()
			return fmt.Errorf("failed to create destination file: %w", err)
		}

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
