package install

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// ExtractorFunc is a type that represents a function for extracting archives.
// It takes an io.ReaderAt, the size of the content, and the destination directory.
type ExtractorFunc func(io.ReaderAt, int64, string) error

// Constants for repeated strings
const (
	SteamCMDLinuxURL   = "https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz"
	SteamCMDWindowsURL = "https://steamcdn-a.akamaihd.net/client/installer/steamcmd.zip"
	SteamCMDLinuxDir   = "./steamcmd"
	SteamCMDWindowsDir = "C:\\SteamCMD"
)

// Color codes for terminal
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

// Verbose mode flag
var verbose = true // Set this to false to disable verbose logging

// logVerbose prints a message only if verbose mode is enabled.
func logVerbose(message string) {
	if verbose {
		fmt.Print(message)
	}
}

// logError prints an error message regardless of verbose mode.
func logError(message string) {
	fmt.Print(ColorRed + message + ColorReset)
}

// logSuccess prints a success message only if verbose mode is enabled.
func logSuccess(message string) {
	if verbose {
		fmt.Print(ColorGreen + message + ColorReset)
	}
}

// InstallAndRunSteamCMD installs and runs SteamCMD based on the platform (Windows/Linux).
// It automatically detects the OS and calls the appropriate installation function.
func InstallAndRunSteamCMD() {
	if runtime.GOOS == "windows" {
		installSteamCMDWindows()
	} else if runtime.GOOS == "linux" {
		installSteamCMDLinux()
	} else {
		logError("❌ SteamCMD installation is not supported on this OS.\n")
		return
	}
}

// installSteamCMD downloads and installs SteamCMD for the given platform.
func installSteamCMD(platform string, steamCMDDir string, downloadURL string, extractFunc ExtractorFunc) {
	// Check if SteamCMD is already installed
	if _, err := os.Stat(steamCMDDir); os.IsNotExist(err) {
		logVerbose(ColorYellow + "⚠️ SteamCMD not found for " + platform + ", downloading...\n" + ColorReset)

		// Create SteamCMD directory
		if err := os.MkdirAll(steamCMDDir, os.ModePerm); err != nil {
			logError("❌ Error creating SteamCMD directory: " + err.Error() + "\n")
			return
		}
		logVerbose("✅ Created SteamCMD directory: " + steamCMDDir + "\n")

		// Ensure cleanup on failure
		success := false
		defer func() {
			if !success {
				logVerbose(ColorYellow + "⚠️ Cleaning up due to failure...\n" + ColorReset)
				os.RemoveAll(steamCMDDir)
			}
		}()

		// Validate download URL
		if err := validateURL(downloadURL); err != nil {
			logError("❌ Invalid download URL: " + err.Error() + "\n")
			return
		}
		logVerbose("✅ Validated download URL: " + downloadURL + "\n")

		// Download SteamCMD with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
		if err != nil {
			logError("❌ Error creating HTTP request: " + err.Error() + "\n")
			return
		}
		logVerbose("✅ Created HTTP request for download.\n")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logError("❌ Error downloading SteamCMD: " + err.Error() + "\n")
			return
		}
		defer resp.Body.Close()
		logVerbose("✅ Successfully downloaded SteamCMD.\n")

		// Check for successful HTTP response
		if resp.StatusCode != http.StatusOK {
			logError("❌ Failed to download SteamCMD: HTTP status " + resp.Status + "\n")
			return
		}
		logVerbose("✅ Received HTTP status: " + resp.Status + "\n")

		// Read the downloaded content into memory
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			logError("❌ Error reading SteamCMD content: " + err.Error() + "\n")
			return
		}
		logVerbose("✅ Read SteamCMD content into memory.\n")

		// Create a reader for the content
		contentReader := bytes.NewReader(content)

		// Extract the content using the provided extractor function
		if err := extractFunc(contentReader, int64(len(content)), steamCMDDir); err != nil {
			logError("❌ Error extracting SteamCMD: " + err.Error() + "\n")
			return
		}
		logVerbose("✅ Successfully extracted SteamCMD.\n")

		// Mark installation as successful
		success = true
		logSuccess("✅ SteamCMD installed successfully.\n")
	} else {
		logVerbose("✅ SteamCMD is already installed.\n")
	}

	// Run SteamCMD
	runSteamCMD(steamCMDDir)
}

// installSteamCMDLinux downloads and installs SteamCMD on Linux.
func installSteamCMDLinux() {
	installSteamCMD("Linux", SteamCMDLinuxDir, SteamCMDLinuxURL, untarWrapper)
}

// installSteamCMDWindows downloads and installs SteamCMD on Windows.
func installSteamCMDWindows() {
	installSteamCMD("Windows", SteamCMDWindowsDir, SteamCMDWindowsURL, unzip)
}

// runSteamCMD runs the SteamCMD command to update the game.
func runSteamCMD(steamCMDDir string) {
	currentDir, err := os.Getwd()
	if err != nil {
		logError("❌ Error getting current working directory: " + err.Error() + "\n")
		return
	}
	logVerbose("✅ Current working directory: " + currentDir + "\n")

	// Build SteamCMD command
	cmd := buildSteamCMDCommand(steamCMDDir, currentDir)

	// Set output to stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	logVerbose(ColorBlue + "🕑 Running SteamCMD...\n" + ColorReset)
	err = cmd.Run()
	if err != nil {
		logError("❌ Error running SteamCMD: " + err.Error() + "\n")
		return
	}
	logSuccess("✅ SteamCMD executed successfully.\n")
}

// buildSteamCMDCommand constructs the SteamCMD command based on the OS.
func buildSteamCMDCommand(steamCMDDir, currentDir string) *exec.Cmd {
	var cmdPath string
	if runtime.GOOS == "windows" {
		cmdPath = filepath.Join(steamCMDDir, "steamcmd.exe")
	} else if runtime.GOOS == "linux" {
		cmdPath = filepath.Join(steamCMDDir, "steamcmd.sh")
	}
	logVerbose("✅ SteamCMD command path: " + cmdPath + "\n")

	return exec.Command(cmdPath, "+force_install_dir", currentDir, "+login", "anonymous", "+app_update", "600760", "+quit")
}

// untarWrapper adapts the untar function to match the ExtractorFunc signature.
func untarWrapper(r io.ReaderAt, _ int64, dest string) error {
	return untar(dest, io.NewSectionReader(r, 0, 1<<63-1)) // Use a large size for the section reader
}

// unzip extracts a zip archive.
func unzip(zipReader io.ReaderAt, size int64, dest string) error {
	reader, err := zip.NewReader(zipReader, size)
	if err != nil {
		return fmt.Errorf("failed to create zip reader: %w", err)
	}

	// Ensure the destination directory exists
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	for _, f := range reader.File {
		// Sanitize the file path to prevent path traversal
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			// Create directory with the same permissions as in the zip file
			if err := os.MkdirAll(fpath, f.Mode()); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		// Create the file with the same permissions as in the zip file
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer outFile.Close()

		// Open the file in the zip archive
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip: %w", err)
		}
		defer rc.Close()

		// Copy the file contents using a buffer for better performance
		buffer := make([]byte, 32*1024) // 32KB buffer
		if _, err := io.CopyBuffer(outFile, rc, buffer); err != nil {
			return fmt.Errorf("failed to copy file contents: %w", err)
		}
	}

	return nil
}

// untar extracts a tar.gz archive.
func untar(dest string, r io.Reader) error {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dest, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", target, err)
			}
		case tar.TypeReg:
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %v", target, err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tr); err != nil {
				return fmt.Errorf("failed to write file %s: %v", target, err)
			}
		case tar.TypeSymlink:
			if err := os.Symlink(header.Linkname, target); err != nil {
				return fmt.Errorf("failed to create symlink %s: %v", target, err)
			}
		default:
			return fmt.Errorf("unknown type: %v in %s", header.Typeflag, header.Name)
		}
	}

	return nil
}

// validateURL checks if a URL is valid.
func validateURL(rawURL string) error {
	_, err := url.ParseRequestURI(rawURL)
	return err
}
