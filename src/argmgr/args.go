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
type GameTemplate struct {
	Meta map[string]interface{} `json:"meta"`
	Args map[string][]GameArg   `json:"args"`
}

func LoadRunfile(gameName, runFilesFolder string) (*GameTemplate, error) {
	filePath := filepath.Join(runFilesFolder, fmt.Sprintf("run%s.ssui", gameName))
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read runfile: %w", err)
	}

	var runfile GameTemplate
	if err := json.Unmarshal(fileData, &runfile); err != nil {
		return nil, fmt.Errorf("failed to parse runfile: %w", err)
	}

	// Initialize runtime values
	for category := range runfile.Args {
		for i := range runfile.Args[category] {
			runfile.Args[category][i].RuntimeValue = runfile.Args[category][i].DefaultValue
		}
	}

	return &runfile, nil
}

func SetArgValue(runfile *GameTemplate, flag string, value string) error {
	for category := range runfile.Args {
		for i := range runfile.Args[category] {
			if runfile.Args[category][i].Flag == flag {
				// Validate based on type
				switch runfile.Args[category][i].Type {
				case "int":
					if _, err := strconv.Atoi(value); err != nil {
						return fmt.Errorf("invalid integer value for %s", flag)
					}
				case "bool":
					if value != "true" && value != "false" {
						return fmt.Errorf("invalid boolean value for %s", flag)
					}
				}
				runfile.Args[category][i].RuntimeValue = value
				return nil
			}
		}
	}
	return fmt.Errorf("argument %s not found", flag)
}

func BuildCommandArgs(runfile *GameTemplate) ([]string, error) {
	var args []string
	allArgs := GetAllArgs(runfile)

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
