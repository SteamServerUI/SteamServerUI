package argmgr

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

// Package-level mutex for file operations
var runfileMutex sync.Mutex

var CurrentRunfile *RunFile

type GameArg struct {
	Flag          string `json:"flag"`
	DefaultValue  string `json:"default"`
	RuntimeValue  string `json:"-"`
	Required      bool   `json:"required"`
	RequiresValue bool   `json:"requires_value"`
	Description   string `json:"description"`
	Type          string `json:"type"`
	Special       string `json:"special,omitempty"`
	UILabel       string `json:"ui_label"`
	UIGroup       string `json:"ui_group"`
	Weight        int    `json:"weight"`
	Min           int    `json:"min,omitempty"`
	Max           int    `json:"max,omitempty"`
	Disabled      bool   `json:"disabled,omitempty"`
}

type Meta struct {
	Name    string `json:"name"`    // SSUI Specific Game Identifier, must match the one in the filename.
	Version string `json:"version"` // Runfile version
}

type RunFile struct {
	Meta              Meta                 `json:"meta"`
	Architecture      string               `json:"architecture,omitempty"`
	SteamAppID        string               `json:"steam_app_id"`
	WindowsExecutable string               `json:"windows_executable"`
	LinuxExecutable   string               `json:"linux_executable"`
	Args              map[string][]GameArg `json:"args"`
}

// LoadRunfile loads the runfile and stores it in CurrentRunfile
func LoadRunfile(gameName, runFilesFolder string) error {
	runfileMutex.Lock()
	defer runfileMutex.Unlock()

	filePath := filepath.Join(runFilesFolder, fmt.Sprintf("run%s.ssui", gameName))
	logger.Runfile.Debug(fmt.Sprintf("loading runfile: %s", filePath))

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		logger.Runfile.Error(fmt.Sprintf("failed to read runfile: %v", err))
		return fmt.Errorf("failed to read runfile: %w", err)
	}

	var runfile RunFile
	if err := json.Unmarshal(fileData, &runfile); err != nil {
		logger.Runfile.Error(fmt.Sprintf("failed to parse runfile: %v", err))
		return fmt.Errorf("failed to parse runfile: %w", err)
	}

	// Check architecture compatibility
	if runfile.Architecture != "" {
		goos := strings.ToLower(runtime.GOOS)
		arch := strings.ToLower(runfile.Architecture)
		if arch != "windows" && arch != "linux" {
			logger.Runfile.Error(fmt.Sprintf("invalid architecture in runfile: %s", runfile.Architecture))
			return fmt.Errorf("invalid architecture in runfile: %s", runfile.Architecture)
		}
		if arch != goos {
			logger.Runfile.Error(fmt.Sprintf("runfile architecture %s does not match current OS %s", arch, goos))
			return fmt.Errorf("runfile architecture %s does not match current OS %s", arch, goos)
		}
	}

	// Initialize runtime values
	for category := range runfile.Args {
		for i := range runfile.Args[category] {
			runfile.Args[category][i].RuntimeValue = runfile.Args[category][i].DefaultValue
		}
	}

	CurrentRunfile = &runfile
	logger.Runfile.Info(fmt.Sprintf("runfile loaded: %s", filePath))
	return nil
}

// SaveRunfile persists the current RunFile to disk
func SaveRunfile() error {
	runfileMutex.Lock()
	defer runfileMutex.Unlock()

	if CurrentRunfile == nil {
		logger.Runfile.Error("runfile not loaded")
		return fmt.Errorf("runfile not loaded")
	}

	// Build filepath
	filePath := filepath.Join(config.RunFilesFolder, fmt.Sprintf("run%s.ssui", config.RunfileGame))
	logger.Runfile.Debug(fmt.Sprintf("saving runfile: %s", filePath))

	// Update DefaultValue from RuntimeValue
	for category := range CurrentRunfile.Args {
		for i := range CurrentRunfile.Args[category] {
			CurrentRunfile.Args[category][i].DefaultValue = CurrentRunfile.Args[category][i].RuntimeValue
		}
	}

	// Validate state
	if err := validateRunfile(); err != nil {
		logger.Runfile.Error(fmt.Sprintf("runfile validation failed: %v", err))
		return fmt.Errorf("runfile validation failed: %w", err)
	}

	// Serialize to JSON
	data, err := json.MarshalIndent(CurrentRunfile, "", "  ")
	if err != nil {
		logger.Runfile.Error(fmt.Sprintf("failed to serialize runfile: %v", err))
		return fmt.Errorf("failed to serialize runfile: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		logger.Runfile.Error(fmt.Sprintf("failed to write runfile: %v", err))
		return fmt.Errorf("failed to write runfile: %w", err)
	}

	logger.Runfile.Info(fmt.Sprintf("runfile saved: %s", filePath))
	return nil
}

// validateRunfile checks the RunFile state before saving
func validateRunfile() error {
	// Validate SteamAppID: non-empty, numeric
	if CurrentRunfile.SteamAppID == "" {
		return fmt.Errorf("SteamAppID is required")
	}
	if _, err := strconv.Atoi(CurrentRunfile.SteamAppID); err != nil {
		return fmt.Errorf("SteamAppID must be numeric, got %s", CurrentRunfile.SteamAppID)
	}

	// Validate WindowsExecutable: if non-empty, must end with .exe
	if CurrentRunfile.WindowsExecutable != "" {
		if !strings.HasSuffix(strings.ToLower(CurrentRunfile.WindowsExecutable), ".exe") {
			return fmt.Errorf("WindowsExecutable must end with .exe, got %s", CurrentRunfile.WindowsExecutable)
		}
	}

	// Validate LinuxExecutable: if non-empty, must not end with .exe
	if CurrentRunfile.LinuxExecutable != "" {
		if strings.HasSuffix(strings.ToLower(CurrentRunfile.LinuxExecutable), ".exe") {
			return fmt.Errorf("LinuxExecutable must not end with .exe, got %s", CurrentRunfile.LinuxExecutable)
		}
	}

	// Validate Meta: ensure Name is non-empty
	if CurrentRunfile.Meta.Name == "" {
		return fmt.Errorf("Meta.Name is required")
	}

	// Existing arg validation
	for _, arg := range GetAllArgs() {
		if arg.Disabled {
			continue
		}
		if arg.Required && arg.RequiresValue && arg.RuntimeValue == "" {
			return fmt.Errorf("required argument %s has no value", arg.Flag)
		}
		switch arg.Type {
		case "int":
			if arg.RuntimeValue != "" {
				if _, err := strconv.Atoi(arg.RuntimeValue); err != nil {
					return fmt.Errorf("invalid integer value for %s", arg.Flag)
				}
			}
		case "bool":
			if arg.RuntimeValue != "" && arg.RuntimeValue != "true" && arg.RuntimeValue != "false" {
				return fmt.Errorf("invalid boolean value for %s", arg.Flag)
			}
		}
	}
	return nil
}

// SetArgValue updates an argument's runtime value and saves the runfile
func SetArgValue(flag string, value string) error {
	if CurrentRunfile == nil {
		logger.Runfile.Error("runfile not loaded")
		return fmt.Errorf("runfile not loaded")
	}

	for category := range CurrentRunfile.Args {
		for i := range CurrentRunfile.Args[category] {
			if CurrentRunfile.Args[category][i].Flag == flag {
				switch CurrentRunfile.Args[category][i].Type {
				case "int":
					if _, err := strconv.Atoi(value); err != nil {
						logger.Runfile.Error(fmt.Sprintf("invalid integer value for %s", flag))
						return fmt.Errorf("invalid integer value for %s", flag)
					}
				case "bool":
					if value != "true" && value != "false" {
						logger.Runfile.Error(fmt.Sprintf("invalid boolean value for %s", flag))
						return fmt.Errorf("invalid boolean value for %s", flag)
					}
				}
				CurrentRunfile.Args[category][i].RuntimeValue = value
				logger.Runfile.Debug(fmt.Sprintf("set arg %s to %s", flag, value))
				if err := SaveRunfile(); err != nil {
					logger.Runfile.Error(fmt.Sprintf("failed to save runfile after setting %s: %v", flag, err))
					return fmt.Errorf("failed to save runfile after setting %s: %w", flag, err)
				}
				return nil
			}
		}
	}

	logger.Runfile.Error(fmt.Sprintf("argument %s not found", flag))
	return fmt.Errorf("argument %s not found", flag)
}

// BuildCommandArgs builds the command-line arguments
func BuildCommandArgs() ([]string, error) {
	if CurrentRunfile == nil {
		logger.Runfile.Error("no runfile is currently loaded")
		return nil, fmt.Errorf("no runfile is currently loaded")
	}

	var args []string
	allArgs := GetAllArgs()

	// Sort by weight (primary) and UIGroup (secondary)
	sort.Slice(allArgs, func(i, j int) bool {
		if allArgs[i].Weight != allArgs[j].Weight {
			return allArgs[i].Weight < allArgs[j].Weight
		}
		return switchCategoryWeight(allArgs[i].UIGroup) < switchCategoryWeight(allArgs[j].UIGroup)
	})

	for _, arg := range allArgs {
		// Skip disabled args
		if arg.Disabled {
			continue
		}

		// Skip optional empty args that require values
		if !arg.Required && arg.RequiresValue && arg.RuntimeValue == "" {
			continue
		}

		// Validate required args that need values
		if arg.Required && arg.RequiresValue && arg.RuntimeValue == "" {
			logger.Runfile.Error(fmt.Sprintf("required argument %s has no value", arg.Flag))
			return nil, fmt.Errorf("required argument %s has no value", arg.Flag)
		}

		// Add flag
		args = append(args, arg.Flag)

		// Special handling
		if arg.Special == "space_delimited" {
			parts := strings.Split(arg.RuntimeValue, " ")
			for _, part := range parts {
				if part != "" {
					args = append(args, part)
				}
			}
			continue
		}

		// Only add value if the argument requires one
		if arg.RequiresValue && arg.RuntimeValue != "" {
			args = append(args, arg.RuntimeValue)
		}
	}

	return args, nil
}

// switchCategoryWeight maps UIGroup to a weight for sorting
func switchCategoryWeight(group string) int {
	switch group {
	case "Basic":
		return 0
	case "Network":
		return 1
	case "Advanced":
		return 2
	default:
		return 3
	}
}
