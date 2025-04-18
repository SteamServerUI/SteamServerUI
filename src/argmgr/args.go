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
	Flag         string `json:"flag"`
	DefaultValue string `json:"default"`
	RuntimeValue string `json:"-"`
	Required     bool   `json:"required"`
	Description  string `json:"description"`
	Type         string `json:"type"`
	Special      string `json:"special,omitempty"`
	UILabel      string `json:"ui_label"`
	UIGroup      string `json:"ui_group"`
	Min          int    `json:"min,omitempty"`
	Max          int    `json:"max,omitempty"`
}

type GameTemplate struct {
	Meta map[string]interface{} `json:"meta"`
	Args map[string][]GameArg   `json:"args"`
}

func LoadGameTemplate(gameName, runFilesFolder string) (*GameTemplate, error) {
	filePath := filepath.Join(runFilesFolder, fmt.Sprintf("run%s.ssui", gameName))
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template: %w", err)
	}

	var template GameTemplate
	if err := json.Unmarshal(fileData, &template); err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Initialize runtime values
	for category := range template.Args {
		for i := range template.Args[category] {
			template.Args[category][i].RuntimeValue = template.Args[category][i].DefaultValue
		}
	}

	return &template, nil
}

func GetAllArgs(template *GameTemplate) []GameArg {
	var allArgs []GameArg
	for _, category := range []string{"basic", "network", "advanced"} {
		if args, exists := template.Args[category]; exists {
			allArgs = append(allArgs, args...)
		}
	}
	return allArgs
}

func GetArgByFlag(template *GameTemplate, flag string) (*GameArg, error) {
	for _, arg := range GetAllArgs(template) {
		if arg.Flag == flag {
			return &arg, nil
		}
	}
	return nil, fmt.Errorf("argument %s not found", flag)
}

func SetArgValue(template *GameTemplate, flag string, value string) error {
	for category := range template.Args {
		for i := range template.Args[category] {
			if template.Args[category][i].Flag == flag {
				// Validate based on type
				switch template.Args[category][i].Type {
				case "int":
					if _, err := strconv.Atoi(value); err != nil {
						return fmt.Errorf("invalid integer value for %s", flag)
					}
				case "bool":
					if value != "true" && value != "false" {
						return fmt.Errorf("invalid boolean value for %s", flag)
					}
				}
				template.Args[category][i].RuntimeValue = value
				return nil
			}
		}
	}
	return fmt.Errorf("argument %s not found", flag)
}

func BuildCommandArgs(template *GameTemplate) ([]string, error) {
	var args []string
	allArgs := GetAllArgs(template)

	// Sort by category order (basic -> network -> advanced)
	sort.Slice(allArgs, func(i, j int) bool {
		return getCategoryWeight(allArgs[i].UIGroup) < getCategoryWeight(allArgs[j].UIGroup)
	})

	for _, arg := range allArgs {
		// Skip optional empty args
		if !arg.Required && arg.RuntimeValue == "" {
			continue
		}

		// Validate required args
		if arg.Required && arg.RuntimeValue == "" {
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

		// Normal value
		if arg.RuntimeValue != "" {
			args = append(args, arg.RuntimeValue)
		}
	}

	return args, nil
}

func getCategoryWeight(group string) int {
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

func GetUIGroups(template *GameTemplate) []string {
	groups := make(map[string]bool)
	for _, arg := range GetAllArgs(template) {
		groups[arg.UIGroup] = true
	}

	var result []string
	for group := range groups {
		result = append(result, group)
	}
	return result
}

func GetArgsByGroup(template *GameTemplate, group string) []GameArg {
	var result []GameArg
	for _, arg := range GetAllArgs(template) {
		if arg.UIGroup == group {
			result = append(result, arg)
		}
	}
	return result
}
