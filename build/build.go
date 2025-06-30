// build/build.go
//go:build ignore
// +build ignore

// run from root with `go run build/build.go`

// build/build.go is supposed to be used A L O N G S I D E the VSCode B U I L D  T A S K  "Build Release". This task will:

// 1. Build the frontend Svelte for production
// 2  Build the frontend Electron for production
// 3  Execute this build script which will:
// 4. Increment the version in config.go
// 5. Build the backend for each platform
// 6. Copy Electron and Go Executables to the build directory

// Then, manually Create a new release on GitHub and upload the files in the build directory, except for the .version file and build.go

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
)

const (
	// ANSI color codes for styling terminal output
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
)

func main() {
	fmt.Printf("%s=== Starting Build-Release Pipeline ===%s\n", colorCyan, colorReset)

	// Increment the version
	newVersion := incrementGoVersion("./src/config/config.go")
	os.Create("build/.version")

	err := os.WriteFile("build/.version", []byte(newVersion), 0644)
	if err != nil {
		fmt.Printf("%s✗ Failed to write version to file: %s%s\n", colorRed, err, colorReset)
	}
	// Platforms to build for
	platforms := []struct {
		os   string
		arch string
	}{
		{"windows", "amd64"},
		{"linux", "amd64"},
	}

	// Clean up old executables
	cleanupOldExecutables(newVersion)

	// Build for each platform
	for _, platform := range platforms {
		fmt.Printf("%s\nBuilding for %s/%s...%s\n", colorBlue, platform.os, platform.arch, colorReset)

		// Set OS and architecture for cross-compilation
		os.Setenv("GOOS", platform.os)
		os.Setenv("GOARCH", platform.arch)

		// Prepare the output file name with the new version, branch, and platform
		var outputName string
		if config.Branch == "release" {
			outputName = fmt.Sprintf("SSUI%s", newVersion)
		} else {
			outputName = fmt.Sprintf("SSUI%s_%s", newVersion, config.Branch)
		}

		// Append appropriate extension based on platform
		if platform.os == "windows" {
			outputName += ".exe"
		}
		if platform.os == "linux" {
			outputName += ".x86_64"
		}

		// Output to /build
		outputPath := filepath.Join("build/release", outputName)

		// Run the go build command targeting server.go at root
		cmd := exec.Command("go", "build", "-ldflags=-s -w", "-gcflags=-l=4", "-o", outputPath, "server.go")

		// Capture any output or errors
		cmdOutput, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("%s✗ Build failed for %s/%s:%s %s\nOutput: %s\n",
				colorRed, platform.os, platform.arch, colorReset, err, string(cmdOutput))
			log.Fatalf("Build process terminated")
		}

		fmt.Printf("%s✓ Build successful!%s Created: %s%s%s\n",
			colorGreen, colorReset, colorYellow, outputPath, colorReset)
	}

	err = copyElectronFiles()
	if err != nil {
		log.Fatalf("Failed to copy Electron files: %v", err)
	}

	fmt.Printf("%s\n=== Build Pipeline Completed ===%s\n", colorCyan, colorReset)
}

// incrementVersion function to increment the version in config.go
func incrementGoVersion(configFile string) string {
	fmt.Printf("%sUpdating version...%s\n", colorBlue, colorReset)

	// Read the content of the config.go file
	content, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read config.go: %s", err)
	}

	// Use regex to find and increment the patch version (assuming version format is x.y.z)
	versionRegex := regexp.MustCompile(`Version\s*=\s*"(\d+)\.(\d+)\.(\d+)"`)
	matches := versionRegex.FindStringSubmatch(string(content))
	if len(matches) != 4 {
		log.Fatalf("Failed to find version in config.go")
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	//patch++ // Soft disabled auto increment for now
	newVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)

	// Replace the old version with the new version
	newContent := versionRegex.ReplaceAllString(string(content), fmt.Sprintf(`Version = "%s"`, newVersion))

	// Write the updated content back to config.go
	err = os.WriteFile(configFile, []byte(newContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write updated version to config.go: %s", err)
	}

	fmt.Printf("%s✓ Version updated from %s.%s.%s to %s%s\n",
		colorGreen, matches[1], matches[2], matches[3], newVersion, colorReset)
	fmt.Printf("%s✓ Version increment is soft disabled for now%s\n",
		colorGreen, colorReset)
	return newVersion
}

// Modified cleanupOldExecutables to handle both Windows and Linux executables in /build
func cleanupOldExecutables(buildVersion string) {
	fmt.Printf("%s\nCleaning up old executables...%s\n", colorBlue, colorReset)

	currentVersion := buildVersion
	dir := "./build/release"

	// Ensure build directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("%sNo build/release directory found, skipping cleanup%s\n", colorYellow, colorReset)
		return
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read build directory: %s", err)
	}

	deletedCount := 0
	for _, file := range files {
		filename := file.Name()
		if filepath.Ext(filename) == ".exe" || filepath.Ext(filename) == ".x86_64" {
			match, _ := filepath.Match("SSUI*", filename)
			if match && !strings.Contains(filename, currentVersion) {
				exePath := filepath.Join(dir, filename)
				fmt.Printf("%s- Removing: %s%s%s\n", colorMagenta, colorYellow, exePath, colorReset)

				err := os.Remove(exePath)
				if err != nil {
					fmt.Printf("%s✗ Failed to delete %s: %s%s\n", colorRed, exePath, err, colorReset)
				} else {
					fmt.Printf("%s✓ Deleted successfully%s\n", colorGreen, colorReset)
					deletedCount++
				}
			}
		}
	}

	if deletedCount == 0 {
		fmt.Printf("%sNo old executables found to clean up%s\n", colorYellow, colorReset)
	} else {
		fmt.Printf("%s✓ Cleaned up %d old executable(s)%s\n", colorGreen, deletedCount, colorReset)
	}
}

func copyElectronFiles() error {
	fmt.Printf("%s\nCopying Electron files...%s\n", colorBlue, colorReset)

	// Source and destination directories
	srcDir := "./frontend/dist_electron"
	destDir := "./build/release"

	// Ensure destination directory exists
	if err := os.MkdirAll(destDir, 0755); err != nil {
		fmt.Printf("Error creating build directory: %v\n", err)
		return err
	}

	// Patterns for files to copy
	patterns := []string{"SSUI-Desktop*.deb", "SSUI-Desktop*.exe", "latest-linux.yml", "latest.yml"}

	for _, pattern := range patterns {
		// Find files matching the pattern
		files, err := filepath.Glob(filepath.Join(srcDir, pattern))
		if err != nil {
			fmt.Printf("Error finding files for pattern %s: %v\n", pattern, err)
			continue
		}

		for _, srcFile := range files {
			// Get the base filename
			fileName := filepath.Base(srcFile)
			destFile := filepath.Join(destDir, fileName)

			// Open source file
			srcData, err := os.ReadFile(srcFile)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", srcFile, err)
				continue
			}

			// Write to destination
			err = os.WriteFile(destFile, srcData, 0644)
			if err != nil {
				fmt.Printf("Error copying file %s to %s: %v\n", srcFile, destFile, err)
				continue
			}

			fmt.Printf("Copied %s to %s\n", srcFile, destFile)
		}
	}

	return nil
}
