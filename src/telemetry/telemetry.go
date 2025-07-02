package telemetry

import (
	"encoding/json"
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

	masterServerFQDN := "https://io.jxsn.dev/api/v1/telemetry" // does not exist yet
	harvestData := harvestTelemetry()
	logger.Security.Infof("Initializing telemetry module...")
	logger.Security.Debugf("Supports Sprintf formatting")
	logger.Security.Warnf("This is how to use the logger")
	logger.Security.Errorf("some sting with %s", masterServerFQDN)

	sendTelemetry(harvestData)
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

func sendTelemetry(harvestData *TelemetryData) {
	// TODO: Send telemetry data to actual master server. for now, we just log it.
	logger.Security.Warnf("Sending telemetry data to master server...")
	jsonData, _ := json.Marshal(harvestData)
	logger.Security.Infof("Telemetry data: %s", string(jsonData))
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
