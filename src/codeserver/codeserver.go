package codeserver

import (
	"runtime"
	"strings"
)

func DownloadInstallCodeServer() string {

	if strings.ToLower(runtime.GOOS) == "linux" {
		return "Code Server is only supported on Linux"
	}

	// TODO: Download & Install Code Server in current directory
	// possibly with this? curl -fsSL https://code-server.dev/install.sh | sh

	return "Successfully installed Code Server"
}

func StartCodeServer() {

	// TODO: Start Code Server with something like this, bound to a local socket only: (AI recommendation lets see if it works)

	//Use exec.Command to run the code-server binary (e.g., ./bin/code-server --socket-path=/tmp/code-server.sock --folder-uri=/gameserver/ --disable-telemetry)

}

func StopCodeServer() {

	// TODO: Stop Code Server

}
