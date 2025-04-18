package argmgr

import (
	"fmt"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
)

func SetupRunfile() (*GameTemplate, error) {
	gameTemplate, err := LoadRunfile(config.RunfileGame, config.RunFilesFolder)
	if err != nil {
		return nil, fmt.Errorf("failed to load game runfile: %w", err)
	}
	return gameTemplate, nil
}
