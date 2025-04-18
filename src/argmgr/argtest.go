package argmgr

import (
	"fmt"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
)

func Test() {
	// Load gameTemplate
	gameTemplate, err := LoadRunfile("Stationeers", config.RunFilesFolder)
	if err != nil {
		panic(err)
	}

	// Discover all available arguments
	allArgs := GetAllArgs(gameTemplate)
	fmt.Println("Available arguments:")
	for _, arg := range allArgs {
		fmt.Printf("%s (%s): %s\n", arg.Flag, arg.Type, arg.Description)
	}

	// Get arguments for UI display
	groups := GetUIGroups(gameTemplate)
	for _, group := range groups {
		fmt.Printf("\n%s Settings:\n", group)
		for _, arg := range GetArgsByGroup(gameTemplate, group) {
			fmt.Printf("- %s: %s (Current: %s)\n", arg.UILabel, arg.Description, arg.RuntimeValue)
		}
	}

	// Update a parameter
	if err := SetArgValue(gameTemplate, "GamePort", "28015"); err != nil {
		panic(err)
	}

	// Build final command line
	args, err := BuildCommandArgs(gameTemplate)
	if err != nil {
		panic(err)
	}
	fmt.Println("\nCommand line:", args)

}
