// build.go
package main

import (
	"StationeersServerUI/src/config"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
)

type ServerConfig struct {
	ExePath  string `xml:"exePath"`
	Settings string `xml:"settings"`
}

type Config struct {
	Server       ServerConfig `xml:"server"`
	SaveFileName string       `xml:"saveFileName"`
}

func main() {
	// Update the config file with the correct executable path based on the OS
	err := updateConfigExePath()
	if err != nil {
		log.Fatalf("Error updating config executable path: %v", err)
	}

	// Load the config to access Version and Branch
	config.LoadConfig("./UIMod/config.json")

	// Increment the version
	newVersion := incrementVersion("src/config/config.go")

	// Prepare the output file name with the new version and branch
	outputName := fmt.Sprintf("StationeersServerControl%s_%s", newVersion, config.Branch)

	// Append .exe only on Windows
	if runtime.GOOS == "windows" {
		outputName += ".exe"
	}

	// Run the go build command with the custom output name
	cmd := exec.Command("go", "build", "-ldflags=-s -w", "-gcflags=-l=4", "-o", outputName, "./src")

	// Capture any output or errors
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Build failed: %s\nOutput: %s", err, string(cmdOutput))
	}

	fmt.Printf("Build successful! Output: %s\n", outputName)

	// Clean up old .exe files that follow the pattern "StationeersServerControl*"
	cleanupOldExecutables(outputName)
}

func updateConfigExePath() error {
	// Determine the executable path based on the operating system
	var exePath string
	if runtime.GOOS == "windows" {
		exePath = "./rocketstation_DedicatedServer.exe"
	} else {
		exePath = "./rocketstation_DedicatedServer"
	}

	// Load the existing config file
	configPath := "./UIMod/config.xml"
	xmlFile, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("error opening config file: %v", err)
	}
	defer xmlFile.Close()

	byteValue, err := io.ReadAll(xmlFile)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	err = xml.Unmarshal(byteValue, &config)
	if err != nil {
		return fmt.Errorf("error unmarshalling config file: %v", err)
	}

	// Update the ExePath
	config.Server.ExePath = exePath

	// Write the updated config back to the file
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("error creating config file: %v", err)
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("error encoding config: %v", err)
	}

	return nil
}

// incrementVersion function to increment the version in config.go
func incrementVersion(configFile string) string {
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

	// Increment the patch version
	patch++

	// Construct the new version
	newVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)

	// Replace the old version with the new version
	newContent := versionRegex.ReplaceAllString(string(content), fmt.Sprintf(`Version = "%s"`, newVersion))

	// Write the updated content back to config.go
	err = os.WriteFile(configFile, []byte(newContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write updated version to config.go: %s", err)
	}

	fmt.Printf("Version updated to %s\n", newVersion)
	return newVersion
}

// cleanupOldExecutables deletes all .exe files matching the pattern "StationeersServerControl*" except the current version
func cleanupOldExecutables(currentExe string) {
	// Get the directory of the current executable
	dir := filepath.Dir(currentExe)

	// Get a list of all .exe files in the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read directory: %s", err)
	}

	// Loop through the files and delete any matching .exe that doesn't match the current one
	for _, file := range files {
		// Check if the file is a .exe and matches the "StationeersServerControl*" pattern
		if filepath.Ext(file.Name()) == ".exe" && filepath.Base(file.Name()) != filepath.Base(currentExe) {
			match, _ := filepath.Match("StationeersServerControl*.exe", file.Name())
			if match {
				exePath := filepath.Join(dir, file.Name())
				fmt.Printf("Deleting old executable: %s\n", exePath)
				err := os.Remove(exePath)
				if err != nil {
					log.Printf("Failed to delete %s: %s", exePath, err)
				} else {
					fmt.Printf("Successfully deleted: %s\n", exePath)
				}
			}
		}
	}
}
