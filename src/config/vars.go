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
	DefaultUserLevel     = "user"
	IsTelemetryEnabled   bool
	BackendUUID          uuid.UUID
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

// Backup settings
var (
	BackupContentDir     string
	BackupsStoreDir      string
	BackupLoopInterval   time.Duration
	BackupMode           string
	BackupMaxFileSize    int64 = 20 * 1024 * 1024 * 1024
	BackupUseCompression bool
	BackupKeepSnapshot   bool
)

// Authentication and security
var (
	AuthEnabled       bool
	JwtKey            string
	AuthTokenLifetime int
	Users             map[string]string
	UserLevels        map[string]string
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
	UIModFolder                = "./UIMod/"
	TLSDir                     = UIModFolder + "config/tls"
	TLSCertPath                = "./UIMod/config/tls/cert.pem"
	TLSKeyPath                 = "./UIMod/config/tls/key.pem"
	ConfigPath                 = "./UIMod/config/settings.json"
	CustomDetectionsFilePath   = "./UIMod/config/customdetections.json"
	LogFolder                  = "./UIMod/logs/"
	SSCMWebDir                 = "./UIMod/sscm/"
	SSCMFilePath               = "./BepInEx/plugins/SSCM/SSCM.socket"
	SSCMPluginDir              = "./BepInEx/plugins/SSCM/"
	RunFilesFolder             = "./UIMod/runfiles/"
	CodeServerPath             = UIModFolder + "CodeServer/"
	CodeServerInstallDir       = UIModFolder + "CodeServer/standalone"
	CodeServerSocketPath       = CodeServerPath + "codeserver.sock"
	CodeServerBinaryPath       = CodeServerPath + "standalone/bin/code-server"
	CodeServerInstallScriptURL = "https://code-server.dev/install.sh"
	CodeServerConfigFilePath   = CodeServerPath + "config.yaml"
	CodeServerUserDataDir      = CodeServerPath + "userdata"
	CodeServerExtensionsDir    = CodeServerPath + "extensions"
	CodeServerSettingsFilePath = CodeServerUserDataDir + "/User/settings.json"
)

// Bundled Assets

var V2UIFS embed.FS
var V1UIFS embed.FS
var TWOBOXFS embed.FS
