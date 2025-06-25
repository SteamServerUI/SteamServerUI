package config

import (
	"embed"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

/*
config.Version and config.Branch can be found in config.go

ConfigMu protects all config variables. Lock it for writes; reads are safe
if writes only happen via applyConfig or with ConfigMu locked.

WARNING: Do NOT set any config vars without locking ConfigMu:
config.ConfigMu.Lock()
config.SomeConfigVar = newValue
config.ConfigMu.Unlock()
*/

// DeferredAction is a function type for actions to be run after a setter completes
type DeferredAction func()

var ConfigMu sync.Mutex

// Game Server configuration
var (
	WorldName       string
	BackupWorldName string
)

// Logging, debugging and misc
var (
	IsDebugMode          bool //only used for pprof server, keep it like this and check the log level instead. Debug = 10
	CreateSSUILogFile    bool
	LogLevel             int
	LogMessageBuffer     string
	IsFirstTimeSetup     bool
	BufferFlushTicker    *time.Ticker
	SSEMessageBufferSize = 2000
	MaxSSEConnections    = 20
	GameServerAppID      int
	GameBranch           string
	SubsystemFilters     []string
	GameServerUUID       uuid.UUID // Assigned at startup to the current instance of the server we are managing. Currently unused.
	BackendEndpointPort  string
	BackendEndpointIP    string
	LegacyLogFile        string
)

// Discord integration
var (
	DiscordToken            string
	DiscordSession          *discordgo.Session
	IsDiscordEnabled        bool
	ControlChannelID        string
	StatusChannelID         string
	LogChannelID            string
	ErrorChannelID          string
	ConnectionListChannelID string
	SaveChannelID           string
	ControlPanelChannelID   string
	DiscordCharBufferSize   int
	ControlMessageID        string
	ExceptionMessageID      string
	BlackListFilePath       string
)

// Backup and cleanup settings
var (
	IsCleanupEnabled        bool
	BackupKeepLastN         int
	BackupKeepDailyFor      time.Duration
	BackupKeepWeeklyFor     time.Duration
	BackupKeepMonthlyFor    time.Duration
	BackupCleanupInterval   time.Duration
	ConfiguredBackupDir     string
	ConfiguredSafeBackupDir string
	BackupWaitTime          time.Duration
)

// Authentication and security
var (
	AuthEnabled       bool
	JwtKey            string
	AuthTokenLifetime int
	Users             map[string]string
)

// SSUI Updates
var (
	IsUpdateEnabled        bool
	AllowPrereleaseUpdates bool
	AllowMajorUpdates      bool
)

// SSCM (Stationeers Server Command Manager) settings

var (
	IsSSCMEnabled       bool
	IsCodeServerEnabled bool
)

// runfile Settings

var (
	RunfileGame string
)

// steamcmd Settings

var (
	SteamCMDLinuxDir   string = "./steamcmd"
	SteamCMDWindowsDir string = "C:\\SteamCMD"
	SteamCMDLinuxURL   string = "https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz"
	SteamCMDWindowsURL string = "https://steamcdn-a.akamaihd.net/client/installer/steamcmd.zip"
)

// File paths
var (
	TLSCertPath              = "./UIMod/tls/cert.pem"
	TLSKeyPath               = "./UIMod/tls/key.pem"
	ConfigPath               = "./UIMod/config/config.json"
	CustomDetectionsFilePath = "./UIMod/detectionmanager/customdetections.json"
	LogFolder                = "./UIMod/logs/"
	UIModFolder              = "./UIMod/"
	ConfigHtmlPath           = "./UIMod/ui/config.html"
	DetectionManagerHtmlPath = "./UIMod/ui/detectionmanager.html"
	IndexHtmlPath            = "./UIMod/ui/index.html"
	SSCMWebDir               = "./UIMod/sscm/"
	SSCMFilePath             = "./BepInEx/plugins/SSCM/SSCM.socket"
	SSCMPluginDir            = "./BepInEx/plugins/SSCM/"
	RunFilesFolder           = "./UIMod/runfiles/"
	CodeServerPath           = UIModFolder + "/CodeServer/"
	CodeServerSocketPath     = CodeServerPath + "/codeserver.sock"
	CodeServerBinaryPath     = "/usr/bin/code-server"
	InstallScriptURL         = "https://code-server.dev/install.sh"
	ConfigFilePath           = CodeServerPath + "/config.yaml"
)

// Bundled Assets

var V2UIFS embed.FS
var V1UIFS embed.FS
var TWOBOXFS embed.FS
