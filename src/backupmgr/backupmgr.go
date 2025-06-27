package backupmgr

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
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
	BackupContentDir:   config.GetBackupContentDir(),
	StoredBackupsDir:   config.GetStoredBackupsDir(),
	BackupLoopInterval: config.GetBackupLoopInterval(),
	BackupMode:         config.GetBackupMode(),
	MaxFileSize:        config.GetMaxFileSize(),
	UseCompression:     config.GetUseCompression(),
	KeepSnapshot:       config.GetKeepSnapshot(),
}

type FileEntry struct {
	SrcPath string
	RelPath string
	Info    os.FileInfo
	Data    []byte // Pre-read file data for regular files
	Err     error
}

func CreateBackup(mode string) error {
	start := time.Now()
	logger.Backup.Debug("Starting backup operation")

	// Check if content directory exists and has content
	if !hasContent() {
		logger.Backup.Debug("No content to backup, skipping")
		return fmt.Errorf("no content to backup, skipping")
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
		// Create a compressed tar in background
		go func() {
			finalPath := filepath.Join(cfg.StoredBackupsDir, "backup_"+timestamp+".tar.gz")
			if err := createCompressedTarFromSnapshotStreaming(snapshotPath, finalPath); err != nil {
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

// tar mode (memory-conscious streaming of files)
func createCompressedTarFromSnapshotStreaming(snapshotPath, backupPath string) error {
	logger.Backup.Debug("Creating compressed tar from snapshot: " + backupPath)

	tarFile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create tar file: %w", err)
	}
	defer tarFile.Close()

	gzipWriter, err := gzip.NewWriterLevel(tarFile, gzip.BestSpeed)
	if err != nil {
		return fmt.Errorf("failed to create gzip writer: %w", err)
	}
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	return processFilesStreamingParallel(snapshotPath, tarWriter)
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
