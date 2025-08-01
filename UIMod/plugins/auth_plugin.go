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
		"message":      "Plugin executed successfully",
	}
}
