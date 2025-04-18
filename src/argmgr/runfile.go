package argmgr

import (
	"fmt"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
)

func SetupRunfile() (*RunFile, error) {
	RunFile, err := LoadRunfile(config.RunfileGame, config.RunFilesFolder)
	if err != nil {
		return nil, fmt.Errorf("failed to load game runfile: %w", err)
	}
	return RunFile, nil
}
