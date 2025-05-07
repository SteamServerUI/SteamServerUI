package backupmgr

import (
	"fmt"
	"os"
	"path/filepath"
)

// RestoreBackup restores a backup with the given index
func (m *BackupManager) RestoreBackup(index int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Find the backup archive
	backupFile := fmt.Sprintf("backup_%d.zip", index)
	backupPath := filepath.Join(m.config.SafeBackupDir, backupFile)

	// Verify the backup file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup archive %s does not exist", backupFile)
	}

	// Clear the backup directory
	if err := os.RemoveAll(m.config.BackupDir); err != nil {
		return fmt.Errorf("failed to clear backup directory: %w", err)
	}
	if err := os.MkdirAll(m.config.BackupDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to recreate backup directory: %w", err)
	}

	// Extract the archive
	if err := unzipDirectory(backupPath, m.config.BackupDir); err != nil {
		return fmt.Errorf("failed to restore backup %s: %w", backupFile, err)
	}

	return nil
}

// revertRestore is no longer needed as restore is atomic
