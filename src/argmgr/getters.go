package argmgr

import (
	"fmt"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

// GetAllArgs returns all GameArgs from the runfile
func GetAllArgs() []GameArg {
	if CurrentRunfile == nil {
		logger.Runfile.Error("runfile not loaded")
		return nil
	}
	return CurrentRunfile.getAllArgs()
}

func GetUIGroups() []string {
	if CurrentRunfile == nil {
		logger.Runfile.Error("runfile not loaded")
		return nil
	}

	groups := make(map[string]bool)
	for _, arg := range CurrentRunfile.getAllArgs() {
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
		logger.Runfile.Error("runfile not loaded")
		return nil
	}

	var result []GameArg
	for _, arg := range CurrentRunfile.getAllArgs() {
		if arg.UIGroup == group {
			result = append(result, arg)
		}
	}
	return result
}

// GetSingleArg retrieves a specific GameArg by its flag
func GetSingleArg(flag string) (*GameArg, error) {
	if CurrentRunfile == nil {
		err := ErrRunfileNotLoaded{Msg: "runfile not loaded"}
		logger.Runfile.Error(err.Error())
		return nil, err
	}

	for _, arg := range CurrentRunfile.getAllArgs() {
		if arg.Flag == flag {
			return &arg, nil
		}
	}

	err := ErrArgNotFound{Flag: flag}
	logger.Runfile.Error(fmt.Sprintf("arg not found: flag=%s", flag))
	return nil, err
}
