package telemetry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
)

type TelemetryData struct {
	Software       string
	BackendUUID    string
	OSType         string
	BackendVersion string
	BackendBranch  string
	Runfile        string
}

// InitTelemetry will be called once on startup to initialize and send anonymized telemetry data to the master server.
func InitTelemetry() {
	if !config.GetIsTelemetryEnabled() {
		return
	}

	//masterServerFQDN := "https://io.jxsn.dev/api/v1/telemetry" // does not exist yet
	masterServerFQDN := "http://localhost:8080/api/v1/telemetry" // for testing
	harvestData := harvestTelemetry()
	logger.Security.Infof("Initializing telemetry module...")
	logger.Security.Debugf("Supports Sprintf formatting")
	logger.Security.Warnf("This is how to use the logger")
	logger.Security.Errorf("some sting with %s", masterServerFQDN)

	sendTelemetry(harvestData, masterServerFQDN)
}

func harvestTelemetry() *TelemetryData {
	telemetryData := TelemetryData{}
	telemetryData.Software = "SSUI"
	telemetryData.BackendUUID = config.GetBackendUUID().String()
	telemetryData.BackendVersion = config.GetBackendVersion()
	telemetryData.BackendBranch = config.GetBackendBranch()
	telemetryData.Runfile = config.GetRunfileGame()
	telemetryData.OSType = getOSType()
	return &telemetryData
}

func sendTelemetry(harvestData *TelemetryData, endpointURI string) {
	logger.Security.Warnf("Sending telemetry data to master server...")
	jsonData, _ := json.Marshal(harvestData)
	// send data to master server
	resp, err := http.Post(endpointURI, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Security.Error("Error sending telemetry data: " + err.Error())
		return
	}
	fmt.Print(resp.StatusCode)
	if resp.StatusCode != 200 {
		logger.Security.Error("Error sending telemetry data: non-200 status code returned.")
		return
	}
	logger.Security.Info("Telemetry data sent successfully.")
}

func getOSType() string {
	switch runtime.GOOS {
	case "windows":
		return "windows"
	case "linux":
		return "linux"
	default:
		return "unknown"
	}
}
