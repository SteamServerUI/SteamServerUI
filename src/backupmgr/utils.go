package backupmgr

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
)

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		return err
	}

	return destination.Sync()
}

// zipDirectory creates a zip archive of the source directory at the destination path,
// excluding specified files and folders.
func zipDirectory(srcDir, dstZip string) error {
	// Get the current executable path
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeName := filepath.Base(exePath)    // Current executable
	UIModPath := config.GetUIModFolder() // UIMod folder

	// Define paths to exclude
	exclusions := []string{
		exeName,                          // ./UIMod folder
		filepath.Join(srcDir, UIModPath), // Ensure full path to UIMod is covered
	}

	archive, err := os.Create(dstZip)
	if err != nil {
		return err
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	return filepath.Walk(srcDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Compute relative path
		relPath, err := filepath.Rel(srcDir, filePath)
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if relPath == "." {
			return nil
		}

		// Check if the file or folder should be excluded
		for _, exclude := range exclusions {
			// Check if the filePath matches or is under an excluded path
			if strings.HasPrefix(filePath, exclude) || relPath == exclude || filepath.Base(filePath) == exclude {
				if info.IsDir() {
					return filepath.SkipDir // Skip the entire directory
				}
				return nil // Skip the file
			}
		}

		// Create zip header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = relPath
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Copy file contents if it's not a directory
		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			return err
		}
		return nil
	})
}

// unzipDirectory extracts a zip archive to the destination directory
func unzipDirectory(srcZip, dstDir string) error {
	reader, err := zip.OpenReader(srcZip)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		filePath := filepath.Join(dstDir, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		srcFile, err := file.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return err
		}
	}
	return nil
}

// parseBackupIndex extracts the backup index from a filename
func parseBackupIndex(filename string) int {
	re := regexp.MustCompile(`backup_(\d+)\.zip`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) < 2 {
		return -1
	}

	index, err := strconv.Atoi(matches[1])
	if err != nil {
		return -1
	}

	return index
}

// isValidBackupFile checks if a file is a valid backup archive
func isValidBackupFile(filename string) bool {
	return strings.HasPrefix(filename, "backup_") && strings.HasSuffix(filename, ".zip")
}
