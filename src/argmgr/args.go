package argmgr

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var CurrentRunfile *RunFile

type GameArg struct {
	Flag          string `json:"flag"`
	DefaultValue  string `json:"default"`
	RuntimeValue  string `json:"-"`
	Required      bool   `json:"required"`
	RequiresValue bool   `json:"requires_value"` // New field to distinguish between flags that need values
	Description   string `json:"description"`
	Type          string `json:"type"`
	Special       string `json:"special,omitempty"`
	UILabel       string `json:"ui_label"`
	UIGroup       string `json:"ui_group"`
	Weight        int    `json:"weight"` // New field for precise ordering
	Min           int    `json:"min,omitempty"`
	Max           int    `json:"max,omitempty"`
	Disabled      bool   `json:"disabled,omitempty"`
}
type RunFile struct {
	Meta map[string]interface{} `json:"meta"`
	Args map[string][]GameArg   `json:"args"`
}

// LoadRunfile loads the runfile and stores it in CurrentRunfile
func LoadRunfile(gameName, runFilesFolder string) error {
	filePath := filepath.Join(runFilesFolder, fmt.Sprintf("run%s.ssui", gameName))
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read runfile: %w", err)
	}

	var runfile RunFile
	if err := json.Unmarshal(fileData, &runfile); err != nil {
		return fmt.Errorf("failed to parse runfile: %w", err)
	}

	// Initialize runtime values
	for category := range runfile.Args {
		for i := range runfile.Args[category] {
			runfile.Args[category][i].RuntimeValue = runfile.Args[category][i].DefaultValue
		}
	}

	CurrentRunfile = &runfile
	return nil
}

func SetArgValue(flag string, value string) error {
	if CurrentRunfile == nil {
		return fmt.Errorf("runfile not loaded")
	}

	for category := range CurrentRunfile.Args {
		for i := range CurrentRunfile.Args[category] {
			if CurrentRunfile.Args[category][i].Flag == flag {
				// Validate based on type
				switch CurrentRunfile.Args[category][i].Type {
				case "int":
					if _, err := strconv.Atoi(value); err != nil {
						return fmt.Errorf("invalid integer value for %s", flag)
					}
				case "bool":
					if value != "true" && value != "false" {
						return fmt.Errorf("invalid boolean value for %s", flag)
					}
				}
				CurrentRunfile.Args[category][i].RuntimeValue = value
				return nil
			}
		}
	}
	return fmt.Errorf("argument %s not found", flag)
}

// BuildCommandArgs also doesn't need a runfile parameter
func BuildCommandArgs() ([]string, error) {
	if CurrentRunfile == nil {
		return nil, fmt.Errorf("runfile not loaded")
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
