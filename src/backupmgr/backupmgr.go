package backupmgr

import (
	"archive/tar"
	"compress/gzip"
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
)

type Bckupcfg struct {
	BackupContentDir   string
	StoredBackupsDir   string
	BackupLoopInterval time.Duration
	BackupMode         string
	MaxFileSize        int64
	UseCompression     bool
	KeepSnapshot       bool
}

var cfg = Bckupcfg{
	BackupContentDir:   "./backupContentDir",
	StoredBackupsDir:   "./storedBackupsDir",
	BackupLoopInterval: 5 * time.Minute,
	BackupMode:         "tar",              // Changed from "zip" to "tar"
	MaxFileSize:        1024 * 1024 * 1024, // 1GB
	UseCompression:     true,
	KeepSnapshot:       false,
}

func InitBackupMgr() {
	logger.Backup.Debug("Initializing Backup Manager")

	if err := ensureDirectories(); err != nil {
		logger.Backup.Error("Failed to create directories: " + err.Error())
		return
	}
	StartBackupLoop()

	logger.Backup.Info("Backup Manager Initialized")
	logger.Backup.Info("Content Directory: " + cfg.BackupContentDir)
	logger.Backup.Info("Backup Directory: " + cfg.StoredBackupsDir)
	logger.Backup.Info("Backup Interval: " + cfg.BackupLoopInterval.String())
	logger.Backup.Info("Backup Mode: " + cfg.BackupMode)
	logger.Backup.Info("Max File Size: " + fmt.Sprintf("%d MB", cfg.MaxFileSize/(1024*1024)))
	logger.Backup.Info("Use Compression: " + fmt.Sprintf("%t", cfg.UseCompression))
	logger.Backup.Info("Keep Snapshots: " + fmt.Sprintf("%t", cfg.KeepSnapshot))
}

func ensureDirectories() error {
	dirs := []string{cfg.BackupContentDir, cfg.StoredBackupsDir}
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

	ticker := time.NewTicker(cfg.BackupLoopInterval)
	defer ticker.Stop()

	// Perform initial backup
	CreateBackup(cfg.BackupMode)

	for {
		select {
		case <-ctx.Done():
			logger.Backup.Debug("Backup loop cancelled")
			return
		case <-ticker.C:
			CreateBackup(cfg.BackupMode)
		}
	}
}

func CreateBackup(mode string) error {
	start := time.Now()
	logger.Backup.Debug("Starting backup operation")

	// Check if content directory exists and has content
	if !hasContent() {
		logger.Backup.Debug("No content to backup, skipping")
		return nil
	}

	// Generate timestamp for this backup
	timestamp := time.Now().Format("2006-01-02_15-04-05")

	// Step 1: Always create a fast snapshot first
	snapshotPath := filepath.Join(cfg.StoredBackupsDir, "snapshot_"+timestamp)
	logger.Backup.Debug("Creating snapshot: " + snapshotPath)

	if err := createSnapshot(snapshotPath); err != nil {
		logger.Backup.Error("Failed to create snapshot: " + err.Error())
		return err
	}

	snapshotDuration := time.Since(start)
	logger.Backup.Info("Snapshot created successfully: " + filepath.Base(snapshotPath) + " (took " + snapshotDuration.String() + ")")

	// Step 2: Handle different backup modes
	switch mode {
	case "copy":
		// For copy mode, we're done - the snapshot IS the backup
		if !cfg.KeepSnapshot {
			// Rename snapshot to final backup name
			finalPath := filepath.Join(cfg.StoredBackupsDir, "backup_"+timestamp)
			if err := os.Rename(snapshotPath, finalPath); err != nil {
				logger.Backup.Error("Failed to rename snapshot: " + err.Error())
				return err
			}
			logger.Backup.Info("Copy backup completed: " + filepath.Base(finalPath))
		}
	case "tar":
		// Create compressed tar in background
		go func() {
			finalPath := filepath.Join(cfg.StoredBackupsDir, "backup_"+timestamp+".tar.gz")
			if err := createCompressedTarFromSnapshot(snapshotPath, finalPath); err != nil {
				logger.Backup.Error("Background tar compression failed: " + err.Error())
			} else {
				duration := time.Since(start)
				compressionNote := ""
				if cfg.UseCompression {
					compressionNote = " (with compression enabled)"
				}
				logger.Backup.Info("Tar backup completed: " + filepath.Base(finalPath) + " (took " + duration.String() + compressionNote + ")")
			}

			// Cleanup snapshot if not keeping it
			if !cfg.KeepSnapshot {
				if err := os.RemoveAll(snapshotPath); err != nil {
					logger.Backup.Warn("Failed to cleanup snapshot: " + err.Error())
				}
			}
		}()
	default:
		// Cleanup snapshot for unsupported modes
		os.RemoveAll(snapshotPath)
		logger.Backup.Error("Unsupported backup mode:" + mode)
		return fmt.Errorf("unsupported backup mode: %s", mode)
	}

	return nil
}

func hasContent() bool {
	entries, err := os.ReadDir(cfg.BackupContentDir)
	if err != nil {
		logger.Backup.Warn("Unable to read content directory: " + err.Error())
		return false
	}
	return len(entries) > 0
}

// Create a fast snapshot of the content directory
func createSnapshot(snapshotPath string) error {
	// Create snapshot directory
	if err := os.MkdirAll(snapshotPath, 0755); err != nil {
		return fmt.Errorf("failed to create snapshot directory: %w", err)
	}

	// Copy all files and directories quickly
	err := filepath.Walk(cfg.BackupContentDir, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Backup.Warn("Skipping file due to error: " + srcPath + " - " + err.Error())
			return nil
		}

		// Skip the root directory itself
		if srcPath == cfg.BackupContentDir {
			return nil
		}

		// Check file size limit with improved logging
		if !info.IsDir() && info.Size() > cfg.MaxFileSize {
			logger.Backup.Warn("Skipping file due to size limit:" + srcPath + " (size:" + fmt.Sprintf("%d MB", info.Size()/(1024*1024)) + ", limit:" + fmt.Sprintf("%d MB", cfg.MaxFileSize/(1024*1024)) + ")")
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(cfg.BackupContentDir, srcPath)
		if err != nil {
			logger.Backup.Warn("Failed to get relative path for: " + srcPath)
			return nil
		}

		destPath := filepath.Join(snapshotPath, relPath)

		if info.IsDir() {
			// Create directory
			return os.MkdirAll(destPath, info.Mode())
		}

		// Copy file
		return copyFile(srcPath, destPath, info.Mode())
	})

	return err
}

// Create compressed tar from snapshot with improved error handling and validation
func createCompressedTarFromSnapshot(snapshotPath, backupPath string) error {
	logger.Backup.Debug("Creating compressed tar from snapshot: " + backupPath)

	// Create tar file
	tarFile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create tar file: %w", err)
	}
	defer tarFile.Close()

	var writer io.Writer = tarFile

	// Add gzip compression
	var gzipWriter *gzip.Writer = gzip.NewWriter(tarFile)
	writer = gzipWriter
	defer func() {
		if gzipWriter != nil {
			if err := gzipWriter.Close(); err != nil {
				logger.Backup.Error("Failed to close gzip writer:" + err.Error())
			}
		}
	}()

	// Create tar writer
	tarWriter := tar.NewWriter(writer)
	defer func() {
		if err := tarWriter.Close(); err != nil {
			logger.Backup.Error("Failed to close tar writer:" + err.Error())
		}
	}()

	// Walk snapshot and add files to tar
	err = filepath.Walk(snapshotPath, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Backup.Warn("Skipping file due to error: " + srcPath + " - " + err.Error())
			return nil
		}

		// Skip the root directory itself
		if srcPath == snapshotPath {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(snapshotPath, srcPath)
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
		// Ensure forward slashes for cross-platform compatibility
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

	if err != nil {
		return err
	}

	// Close writers to ensure data is flushed
	if err := tarWriter.Close(); err != nil {
		logger.Backup.Error("Failed to close tar writer:" + err.Error())
		return fmt.Errorf("failed to close tar writer: %w", err)
	}
	if gzipWriter != nil {
		if err := gzipWriter.Close(); err != nil {
			logger.Backup.Error("Failed to close gzip writer:" + err.Error())
			return fmt.Errorf("failed to close gzip writer: %w", err)
		}
	}

	// Validate the created file
	if err := validateTarFile(backupPath); err != nil {
		logger.Backup.Error("Created tar file validation failed:" + err.Error())
		return fmt.Errorf("created tar file validation failed: %w", err)
	}

	return nil
}

func validateTarFile(backupPath string) error {
	file, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open tar file for validation: %w", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil || info.Size() == 0 {
		return fmt.Errorf("invalid tar file size: %d", info.Size())
	}

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader during validation: %w", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	_, err = tarReader.Next()
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read first tar header during validation: %w", err)
	}

	return nil
}

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
	entries, err := os.ReadDir(cfg.StoredBackupsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []string
	for _, entry := range entries {
		name := entry.Name()
		// Include all backup formats, but exclude temporary snapshots
		if strings.HasPrefix(name, "backup_") || (strings.HasPrefix(name, "snapshot_") && cfg.KeepSnapshot) {
			backups = append(backups, name)
		}
	}

	return backups, nil
}

func RestoreBackup(backupName string, skipPreBackup bool) error {
	logger.Backup.Info("Starting restore operation for: " + backupName)

	backupPath := filepath.Join(cfg.StoredBackupsDir, backupName)
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup not found: %s", backupName)
	}

	// Create pre-restore backup using standard CreateBackup
	if !skipPreBackup {
		logger.Backup.Info("Creating backup before restore")
		if err := CreateBackup(cfg.BackupMode); err != nil {
			return fmt.Errorf("pre-restore backup failed: %w", err)
		}
	}

	// Determine backup type and restore accordingly
	if strings.HasSuffix(backupName, ".tar.gz") || strings.HasSuffix(backupName, ".tar") {
		return restoreTarBackup(backupPath)
	} else {
		// Assume it's a copy backup (directory)
		return restoreCopyBackup(backupPath)
	}
}

func restoreCopyBackup(backupPath string) error {
	// Clear content directory
	if err := os.RemoveAll(cfg.BackupContentDir); err != nil {
		return fmt.Errorf("failed to clear content directory: %w", err)
	}
	if err := os.MkdirAll(cfg.BackupContentDir, 0755); err != nil {
		return fmt.Errorf("failed to create content directory: %w", err)
	}

	// Copy backup directory to content directory
	return copyDirectory(backupPath, cfg.BackupContentDir)
}

func restoreTarBackup(backupPath string) error {
	// Clear content directory
	if err := os.RemoveAll(cfg.BackupContentDir); err != nil {
		return fmt.Errorf("failed to clear content directory: %w", err)
	}
	if err := os.MkdirAll(cfg.BackupContentDir, 0755); err != nil {
		return fmt.Errorf("failed to create content directory: %w", err)
	}

	return extractTarBackup(backupPath, cfg.BackupContentDir)
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
	logger.Backup.Debug("Opening tar file:" + backupPath)

	// Validate file before opening
	info, err := os.Stat(backupPath)
	if err != nil || info.Size() == 0 {
		logger.Backup.Error("Invalid tar file:" + backupPath + ", size:" + fmt.Sprintf("%d", info.Size()))
		return fmt.Errorf("invalid tar file: %w", err)
	}

	file, err := os.Open(backupPath)
	if err != nil {
		logger.Backup.Error("Failed to open tar file:" + err.Error())
		return fmt.Errorf("failed to open tar file: %w", err)
	}
	defer file.Close()

	var reader io.Reader = file

	// Check if it's a gzipped tar
	if strings.HasSuffix(backupPath, ".tar.gz") {
		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			logger.Backup.Error("Failed to create gzip reader:" + err.Error())
			return fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	tarReader := tar.NewReader(reader)

	// Ensure destDir has trailing slash for proper prefix checking
	destDir = filepath.Clean(destDir)
	if !strings.HasSuffix(destDir, string(filepath.Separator)) {
		destDir += string(filepath.Separator)
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Backup.Error("Failed to read tar header:" + err.Error())
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		logger.Backup.Debug("Processing tar entry:" + header.Name + ", Typeflag:" + fmt.Sprintf("%d", header.Typeflag))

		// Handle unsupported header types
		if header.Typeflag != tar.TypeReg && header.Typeflag != tar.TypeDir {
			logger.Backup.Warn("Skipping unsupported type in tar:" + header.Name + ", Typeflag:" + fmt.Sprintf("%d", header.Typeflag))
			continue
		}

		// Security check: reject paths with ".." components
		if strings.Contains(header.Name, "..") {
			logger.Backup.Warn("Skipping path with .. component in tar:" + header.Name)
			continue
		}

		// Build destination path
		destPath := filepath.Join(destDir, header.Name)
		destPath = filepath.Clean(destPath)

		// Security check: ensure the cleaned path is still within destDir
		if !strings.HasPrefix(destPath+string(filepath.Separator), destDir) && destPath != strings.TrimSuffix(destDir, string(filepath.Separator)) {
			logger.Backup.Warn("Skipping path outside destination directory in tar:" + header.Name + " -> " + destPath)
			continue
		}

		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(destPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		// Create parent directory if needed
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

	logger.Backup.Info("Backup restored successfully")
	return nil
}

func IsBackupRunning() bool {
	mu.Lock()
	defer mu.Unlock()
	return isRunning
}

func CleanupOldSnapshots(maxAge time.Duration) error {
	entries, err := os.ReadDir(cfg.StoredBackupsDir)
	if err != nil {
		return fmt.Errorf("failed to read backup directory: %w", err)
	}

	cutoff := time.Now().Add(-maxAge)

	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), "snapshot_") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			snapshotPath := filepath.Join(cfg.StoredBackupsDir, entry.Name())
			if err := os.RemoveAll(snapshotPath); err != nil {
				logger.Backup.Warn("Failed to cleanup old snapshot: " + entry.Name() + " - " + err.Error())
			} else {
				logger.Backup.Info("Cleaned up old snapshot: " + entry.Name())
			}
		}
	}

	return nil
}

func CleanupOldBackups(maxAge time.Duration) error {
	entries, err := os.ReadDir(cfg.StoredBackupsDir)
	if err != nil {
		return fmt.Errorf("failed to read backup directory: %w", err)
	}

	cutoff := time.Now().Add(-maxAge)

	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), "backup_") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			backupPath := filepath.Join(cfg.StoredBackupsDir, entry.Name())
			if err := os.RemoveAll(backupPath); err != nil {
				logger.Backup.Warn("Failed to cleanup old backup: " + entry.Name() + " - " + err.Error())
			} else {
				logger.Backup.Info("Cleaned up old backup: " + entry.Name())
			}
		}
	}

	return nil
}
