package telemetry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

type TelemetryData struct {
	// Core
	Software       string `json:"software"`
	BackendUUID    string `json:"backend_uuid"`
	BackendVersion string `json:"backend_version"`
	BackendBranch  string `json:"backend_branch"`
	Runfile        string `json:"runfile"`

	// System
	OSType         string `json:"os_type"`
	CPUArch        string `json:"cpu_arch"`
	CPUCores       int    `json:"cpu_cores"`
	TotalUsedRAMMB uint64 `json:"total_used_ram_mb"`
	IsContainer    bool   `json:"is_container"`

	// Features
	DiscordEnabled bool `json:"discord_enabled"`
	AuthEnabled    bool `json:"auth_enabled"`
	UpdatesEnabled bool `json:"updates_enabled"`
	BepInExEnabled bool `json:"bepinex_enabled"`
	SSCMEnabled    bool `json:"sscm_enabled"`

	// Network
	Port string `json:"port"`
}

// InitTelemetry will be called once on startup to initialize and send anonymized telemetry data to the master server.
func InitTelemetry() {
	if !config.GetIsTelemetryEnabled() {
		return
	}

	masterServerFQDN := "https://jxsn.dev/api/v1/telemetry"
	harvestData := harvestTelemetry()
	logger.Security.Infof("Initializing telemetry module...")
	sendTelemetry(harvestData, masterServerFQDN)
}

func harvestTelemetry() *TelemetryData {
	telemetryData := TelemetryData{}
	telemetryData.Software = "SteamServerUI"
	telemetryData.BackendUUID = "unset"
	telemetryData.BackendVersion = config.GetVersion()
	telemetryData.BackendBranch = config.GetBranch()
	telemetryData.Runfile = config.GetRunfileIdentifier()
	if telemetryData.Runfile == "" {
		telemetryData.Runfile = "none"
	}
	telemetryData.OSType = getOSType()
	telemetryData.CPUArch = runtime.GOARCH
	telemetryData.CPUCores = runtime.NumCPU()
	telemetryData.TotalUsedRAMMB = uint64(getCurrentHeapMB())
	telemetryData.IsContainer = config.GetIsDockerContainer()
	telemetryData.DiscordEnabled = config.GetIsDiscordEnabled()
	telemetryData.AuthEnabled = config.GetAuthEnabled()
	telemetryData.UpdatesEnabled = config.GetIsUpdateEnabled()
	telemetryData.BepInExEnabled = config.GetIsBepInExEnabled()
	telemetryData.SSCMEnabled = config.GetIsSSCMEnabled()
	telemetryData.Port = config.GetBackendEndpointPort()
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
	if resp.StatusCode != 202 {
		logger.Security.Error("Error sending telemetry data: non-200 status code returned.")
		return
	}
	logger.Security.Info("Telemetry data sent successfully.")
	logger.Security.Debug("Sent Telemetry data:" + string(jsonData))
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

func getCurrentHeapMB() uint64 {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return mem.HeapAlloc / 1024 / 1024
}
