package backupmgr

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
)

// Cleanup performs backup cleanup according to retention policy
func (m *BackupManager) Cleanup() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Clean safe backup dir with retention policy
	if err := m.cleanSafeBackupDir(); err != nil {
		return fmt.Errorf("safe backup dir cleanup failed: %w", err)
	}

	return nil
}

// cleanSafeBackupDir cleans the safe backup directory with retention policy
func (m *BackupManager) cleanSafeBackupDir() error {
	groups, err := m.getBackupGroups()
	if err != nil {
		return err
	}

	// Sort newest first
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].ModTime.After(groups[j].ModTime)
	})

	now := time.Now()
	var (
		lastKeptDaily  time.Time
		lastKeptWeekly time.Time
	)

	for i, group := range groups {
		age := now.Sub(group.ModTime)

		// Always keep the most recent N backups
		if i < m.config.RetentionPolicy.KeepLastN {
			continue
		}

		// Keep daily backups for specified duration
		if age < m.config.RetentionPolicy.KeepDailyFor {
			if lastKeptDaily.IsZero() || group.ModTime.Day() != lastKeptDaily.Day() {
				lastKeptDaily = group.ModTime
				continue
			}
		}

		// Keep weekly backups for specified duration
		if age < m.config.RetentionPolicy.KeepWeeklyFor {
			year1, week1 := group.ModTime.ISOWeek()
			year2, week2 := lastKeptWeekly.ISOWeek()
			if lastKeptWeekly.IsZero() || year1 != year2 || week1 != week2 {
				lastKeptWeekly = group.ModTime
				continue
			}
		}

		// Keep monthly backups for specified duration
		if age < m.config.RetentionPolicy.KeepMonthlyFor {
			if group.ModTime.Month() != lastKeptWeekly.Month() || group.ModTime.Year() != lastKeptWeekly.Year() {
				lastKeptWeekly = group.ModTime
				continue
			}
		}

		// If we get here, the backup should be deleted
		m.deleteBackupGroup(group)
	}

	return nil
}

// deleteBackupGroup removes a backup archive
func (m *BackupManager) deleteBackupGroup(group BackupGroup) {
	if group.ZipFile != "" {
		filePath := filepath.Join(m.config.SafeBackupDir, group.ZipFile)
		if err := os.Remove(filePath); err != nil {
			logger.Backup.Error("Failed to delete backup archive " + filePath + ": " + err.Error())
		}
	}
}
