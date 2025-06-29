package backupmgr

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
)

func processFilesStreamingParallel(snapshotPath string, tarWriter *tar.Writer) error {
	// For large files, we collect metadata in parallel but stream content sequentially
	type FileMetadata struct {
		SrcPath string
		RelPath string
		Info    os.FileInfo
	}

	fileMetadata := make(chan FileMetadata, 100)
	var wg sync.WaitGroup

	// Collect file metadata in parallel
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(fileMetadata)

		filepath.Walk(snapshotPath, func(srcPath string, info os.FileInfo, err error) error {
			if err != nil {
				logger.Backup.Warn("Skipping file due to error: " + srcPath + " - " + err.Error())
				return nil
			}

			if srcPath == snapshotPath {
				return nil
			}

			relPath, err := filepath.Rel(snapshotPath, srcPath)
			if err != nil {
				logger.Backup.Warn("Failed to get relative path for: " + srcPath)
				return nil
			}

			fileMetadata <- FileMetadata{
				SrcPath: srcPath,
				RelPath: relPath,
				Info:    info,
			}
			return nil
		})
	}()

	// Process files sequentially but with pre-collected metadata
	for meta := range fileMetadata {
		header, err := tar.FileInfoHeader(meta.Info, "")
		if err != nil {
			logger.Backup.Warn("Failed to create tar header for: " + meta.SrcPath)
			continue
		}

		header.Name = strings.ReplaceAll(meta.RelPath, "\\", "/")

		if err := tarWriter.WriteHeader(header); err != nil {
			logger.Backup.Warn("Failed to write tar header for: " + meta.SrcPath)
			continue
		}

		if meta.Info.Mode().IsRegular() {
			file, err := os.Open(meta.SrcPath)
			if err != nil {
				logger.Backup.Warn("Failed to open file: " + meta.SrcPath)
				continue
			}

			// Use larger buffer for better I/O performance
			buffer := make([]byte, 64*1024) // 64KB buffer
			_, err = io.CopyBuffer(tarWriter, file, buffer)
			file.Close()

			if err != nil {
				logger.Backup.Warn("Failed to copy file to tar: " + meta.SrcPath)
			}
		}
	}

	wg.Wait()
	return nil
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

		//logger.Backup.Debug("Processing tar entry:" + header.Name + ", Typeflag:" + fmt.Sprintf("%d", header.Typeflag))

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
