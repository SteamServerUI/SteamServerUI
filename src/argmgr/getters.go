package argmgr

import (
	"fmt"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

// GetAllArgs returns all GameArgs from the runfile
func GetAllArgs() []GameArg {
	if CurrentRunfile == nil {
		return nil
	}

	var allArgs []GameArg
	for _, category := range []string{"basic", "network", "advanced"} {
		if args, exists := CurrentRunfile.Args[category]; exists {
			allArgs = append(allArgs, args...)
		}
	}
	return allArgs
}

func GetUIGroups() []string {
	if CurrentRunfile == nil {
		return nil
	}

	groups := make(map[string]bool)
	for _, arg := range GetAllArgs() {
		groups[arg.UIGroup] = true
	}

	var result []string
	for group := range groups {
		result = append(result, group)
	}
	return result
}

func GetArgsByGroup(group string) []GameArg {
	if CurrentRunfile == nil {
		return nil
	}

	var result []GameArg
	for _, arg := range GetAllArgs() {
		if arg.UIGroup == group {
			result = append(result, arg)
		}
	}
	return result
}

// GetSingleArg retrieves a specific GameArg by its flag
func GetSingleArg(flag string) (*GameArg, error) {
	if CurrentRunfile == nil {
		logger.Runfile.Error("runfile not loaded")
		return nil, fmt.Errorf("runfile not loaded")
	}

	for _, arg := range GetAllArgs() {
		if arg.Flag == flag {
			return &arg, nil
		}
	}

	logger.Runfile.Error(fmt.Sprintf("argument %s not found", flag))
	return nil, fmt.Errorf("argument %s not found", flag)
}
