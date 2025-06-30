package backupmgr

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
)

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
