package backupmgr

import (
	"context"
	"sync"
	"time"
)

const (
	defaultWaitTime = 30 * time.Second
)

// BackupConfig holds configuration for backup operations
type BackupConfig struct {
	BackupDir       string
	SafeBackupDir   string
	RetentionPolicy RetentionPolicy
	WaitTime        time.Duration
}

// RetentionPolicy defines backup retention rules
type RetentionPolicy struct {
	KeepLastN       int           // Keep last N backups regardless of age
	KeepDailyFor    time.Duration // Keep daily backups for this duration
	KeepWeeklyFor   time.Duration // Keep weekly backups for this duration
	KeepMonthlyFor  time.Duration // Keep monthly backups for this duration
	CleanupInterval time.Duration // How often to run cleanup
}

// BackupGroup represents a backup archive
type BackupGroup struct {
	Index   int
	ZipFile string
	ModTime time.Time
}

// BackupManager manages backup operations
type BackupManager struct {
	config  BackupConfig
	mu      sync.Mutex
	watcher *fsWatcher
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup // waitGroup tracking of goroutines wg
}
