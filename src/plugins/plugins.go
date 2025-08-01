// main.go - Your main application
package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

/*
Tried using https://github.com/yaegi/yaegi but it didn't work out for this project. Yaegi requires to expose everything as AppSymbols, which is not really what we want.
Ideally, we want to be able to trigger functions / config variables from plugins, but we can't do that with yaegi. Go Plugins are a better option, but only work on Linux - bummer.
Maybe will use yaegi later if there is a neeed to run self contained plugins.

*/

// Plugin interface that your plugins should implement
type Plugin interface {
	Run() interface{}
	GetName() string
}

// Export your config functions to plugins
var AppSymbols = map[string]map[string]reflect.Value{
	"github.com/SteamServerUI/SteamServerUI/v6/src/config/config": {
		"GetAuthEnabled": reflect.ValueOf(config.GetAuthEnabled),
	},
}

func InitPlugins() {
	pluginDir := "./UIMod/plugins"
	files, err := os.ReadDir(pluginDir)
	if err != nil {
		panic(err)
	}

	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)
	i.Use(AppSymbols) // Make your app functions available to plugins

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".go" {
			fullPath := filepath.Join(pluginDir, file.Name())
			fmt.Printf("Loading Plugin: %s\n", file.Name())

			_, err := i.EvalPath(fullPath)
			if err != nil {
				fmt.Printf("Error in %s: %v\n", file.Name(), err)
				continue
			}

			// Try to get the plugin name first
			nameVal, err := i.Eval("main.GetName()")
			if err != nil {
				fmt.Printf("Plugin %s doesn't implement GetName(): %v\n", file.Name(), err)
				continue
			}

			// Run the plugin
			result, err := i.Eval("main.Run()")
			if err != nil {
				fmt.Printf("Error running plugin %s: %v\n", file.Name(), err)
				continue
			}

			fmt.Printf("Plugin '%s' returned: %v\n", nameVal, result)
		}
	}
}

// Example plugin (save as UIMod/plugins/auth_plugin.go)
/*
package main

import (
	"fmt"
	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
)

func GetName() string {
	return "Auth Checker Plugin"
}

func Run() interface{} {
	fmt.Println("👋 Hello from auth plugin!")

	// Now this will work because we exported the function
	authEnabled := config.GetAuthEnabled()

	if authEnabled {
		fmt.Println("Authentication is enabled")
	} else {
		fmt.Println("Authentication is disabled")
	}

	return map[string]interface{}{
		"auth_enabled": authEnabled,
		"message": "Plugin executed successfully",
	}
}
*/
