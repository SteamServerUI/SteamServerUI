package backupmgr

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"

	"github.com/fsnotify/fsnotify"
)

/*
The BackupManager manages backup operations. Each instance is independent with its own config and context.
Background routines (file watching and cleanup) only start when Start() is called. Multiple instances
can coexist but may conflict if configured with overlapping directories.
*/

// Initialize sets up required directories
func (m *BackupManager) Initialize() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	logger.Backup.Info("Initializing backup manager...")
	logger.Backup.Debug("Backup directory: " + m.config.BackupDir)
	logger.Backup.Debug("Safe backup directory: " + m.config.SafeBackupDir)
	if err := os.MkdirAll(m.config.BackupDir, os.ModePerm); err != nil { // Create backup directory
		return err
	}
	return os.MkdirAll(m.config.SafeBackupDir, os.ModePerm)
}

// Start begins the backup monitoring and cleanup routines
func (m *BackupManager) Start() error {
	if err := m.Initialize(); err != nil { // Not really needed here (double check, we use m.Initialize() before m.Start() ), but just in case this gets called from somwhere else in the future
		return fmt.Errorf("failed to initialize backup directories: %w", err)
	}

	// Start file watcher
	watcher, err := newFsWatcher(m.config.BackupDir)
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}
	m.watcher = watcher

	go m.watchBackups()

	if config.GetIsCleanupEnabled() {
		go m.startCleanupRoutine()
	}

	return nil
}

// watchBackups monitors the backup directory for new files
func (m *BackupManager) watchBackups() {
	m.wg.Add(1)
	defer m.wg.Done()

	logger.Backup.Debug("Starting backup file watcher...")
	defer logger.Backup.Debug("Backup file watcher stopped")

	for {
		select {
		case <-m.ctx.Done():
			return
		case event, ok := <-m.watcher.events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				logger.Backup.Info("New file detected: " + event.Name)
				m.handleNewBackup()
			}
		case err, ok := <-m.watcher.errors:
			if !ok {
				return
			}
			logger.Backup.Error("Backup watcher error: " + err.Error())
		}
	}
}

// handleNewBackup processes a newly created file by creating a zip archive
func (m *BackupManager) handleNewBackup() {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()

		time.Sleep(m.config.WaitTime)

		m.mu.Lock()
		defer m.mu.Unlock()

		// Generate a new backup index
		groups, err := m.getBackupGroups()
		if err != nil {
			logger.Backup.Error("Error getting backup groups: " + err.Error())
			return
		}
		index := 1
		if len(groups) > 0 {
			index = groups[0].Index + 1
		}

		// Create zip archive
		zipName := fmt.Sprintf("backup_%d.zip", index)
		dstPath := filepath.Join(m.config.SafeBackupDir, zipName)

		if err := zipDirectory(m.config.BackupDir, dstPath); err != nil {
			logger.Backup.Error("Error creating backup archive " + zipName + ": " + err.Error())
			return
		}

		logger.Backup.Info("Backup archive successfully created: " + dstPath)
	}()
}

// startCleanupRoutine runs periodic backup cleanup
func (m *BackupManager) startCleanupRoutine() {
	m.wg.Add(1)
	defer m.wg.Done()

	ticker := time.NewTicker(m.config.RetentionPolicy.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			if err := m.Cleanup(); err != nil {
				logger.Backup.Error("Backup cleanup error: " + err.Error())
			}
		}
	}
}

// getBackupGroups collects all backup archives
func (m *BackupManager) getBackupGroups() ([]BackupGroup, error) {
	dir, err := os.Open(m.config.SafeBackupDir)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var groups []BackupGroup
	for _, file := range files {
		if !isValidBackupFile(file.Name()) {
			continue
		}

		index := parseBackupIndex(file.Name())
		if index == -1 {
			continue
		}

		groups = append(groups, BackupGroup{
			Index:   index,
			ZipFile: file.Name(),
			ModTime: file.ModTime(),
		})
	}

	// Sort by index (newest first)
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Index > groups[j].Index
	})

	return groups, nil
}

// ListBackups returns information about available backups
// limit: number of recent backups to return (0 for all)
func (m *BackupManager) ListBackups(limit int) ([]BackupGroup, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	groups, err := m.getBackupGroups()
	if err != nil {
		return nil, err
	}

	if limit > 0 && limit < len(groups) {
		groups = groups[:limit]
	}

	return groups, nil
}

// Shutdown stops all backup operations
func (m *BackupManager) Shutdown() {
	logger.Backup.Debug("Shutting down backup manager...")

	m.mu.Lock()
	if m.cancel != nil {
		m.cancel()
		m.cancel = nil
	}

	if m.watcher != nil {
		m.watcher.close()
		m.watcher = nil
	}
	m.mu.Unlock()

	// Wait for all goroutines to finish
	logger.Backup.Debug("Waiting for background tasks to complete...")
	m.wg.Wait()

	logger.Backup.Debug("Backup manager shut down completely")
}

// NewBackupManager creates a new BackupManager instance
func NewBackupManager(cfg BackupConfig) *BackupManager {
	ctx, cancel := context.WithCancel(context.Background())

	if cfg.WaitTime == 0 {
		cfg.WaitTime = defaultWaitTime
	}

	return &BackupManager{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
	}
}
