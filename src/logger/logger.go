package logger

import (
	"fmt"
	"sync"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/ssestream"
)

// Logger instances
var (
	Main       = &Logger{prefix: SYS_MAIN}
	Web        = &Logger{prefix: SYS_WEB}
	Discord    = &Logger{prefix: SYS_DISCORD}
	Backup     = &Logger{prefix: SYS_BACKUP}
	Detection  = &Logger{prefix: SYS_DETECT}
	Core       = &Logger{prefix: SYS_CORE}
	Config     = &Logger{prefix: SYS_CONFIG}
	Install    = &Logger{prefix: SYS_INSTALL}
	SSE        = &Logger{prefix: SYS_SSE}
	Security   = &Logger{prefix: SYS_SECURITY}
	Runfile    = &Logger{prefix: SYS_RUNFILE}
	Codeserver = &Logger{prefix: SYS_CODESERVER}
)

// Severity Levels
const (
	DEBUG = 10 // Fine-grained debugging
	INFO  = 20 // Normal operations
	WARN  = 30 // Potential issues
	ERROR = 40 // Critical errors
)

// Subsystems
const (
	SYS_MAIN       = "MAIN"
	SYS_WEB        = "WEB"
	SYS_DISCORD    = "DISCORD"
	SYS_BACKUP     = "BACKUP"
	SYS_DETECT     = "DETECT"
	SYS_CORE       = "CORE"
	SYS_CONFIG     = "CONFIG"
	SYS_INSTALL    = "INSTALL"
	SYS_SSE        = "SSE"
	SYS_SECURITY   = "SECURITY"
	SYS_RUNFILE    = "RUNFILE"
	SYS_CODESERVER = "CODESERVER"
)

const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
)

// Subsystem color map (distinct colors, cohesive vibe)
var subsystemColors = map[string]string{
	SYS_MAIN:       colorBlue,    // Calm, default system
	SYS_WEB:        colorCyan,    // Clean, UI-related
	SYS_DISCORD:    colorMagenta, // Flashy, chatty subsystem
	SYS_BACKUP:     colorGreen,   // Safe, reliable vibe
	SYS_DETECT:     colorYellow,  // Attention-grabbing for detection
	SYS_CORE:       colorMagenta, // Critical, stands out
	SYS_CONFIG:     colorYellow,  // Warning-like, config tweaks
	SYS_INSTALL:    colorBlue,    // Matches MAIN, setup-related
	SYS_SSE:        colorCyan,    // Matches WEB, streaming vibe
	SYS_SECURITY:   colorRed,     // Screams "pay attention"
	SYS_RUNFILE:    colorMagenta, // Matches CORE, runfile-related
	SYS_CODESERVER: colorCyan,    // TODO
}

type Logger struct {
	mu     sync.Mutex
	prefix string // Subsystem identifier (e.g., "DISCORD")
}

type logEntry struct {
	severity int
	prefix   string // Log type (e.g., "INFO", "CORE"), now a suffix (prints INSTALL/WARN instead of WARN/INSTALL since v6.4.1)
	color    string
	message  string
}

func (l *Logger) Debug(format string, args ...any) {
	l.log(logEntry{DEBUG, "DEBUG", colorReset, fmt.Sprintf(format, args...)}) // Subsystem color
}

func (l *Logger) Info(format string, args ...any) {
	l.log(logEntry{INFO, "INFO", colorReset, fmt.Sprintf(format, args...)}) // Subsystem color
}

func (l *Logger) Warn(format string, args ...any) {
	l.log(logEntry{WARN, "WARN", colorYellow, fmt.Sprintf(format, args...)}) // Yellow for warnings
}

func (l *Logger) Error(format string, args ...any) {
	l.log(logEntry{ERROR, "ERROR", colorRed, fmt.Sprintf(format, args...)}) // Red for errors
}

// log handles the core logging logic
func (l *Logger) log(entry logEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	// Use subsystem color by default, override with severity color if set
	entryColor := subsystemColors[l.prefix]
	if entry.color != colorReset {
		entryColor = entry.color
	}
	// Console version with colors
	consoleLine := fmt.Sprintf("%s%s [%s/%s] %s%s\n", entryColor, timestamp, l.prefix, entry.prefix, entry.message, colorReset)
	// File version without colors
	fileLine := fmt.Sprintf("%s [%s/%s] %s\n", timestamp, l.prefix, entry.prefix, entry.message)
	// Console output

	if entry.severity >= DEBUG {
		ssestream.BroadcastDebugLog(fileLine)
	}

	if entry.severity == INFO {
		ssestream.BroadcastInfoLog(fileLine)
	}

	if entry.severity == WARN {
		ssestream.BroadcastWarnLog(fileLine)
	}

	if entry.severity == ERROR {
		ssestream.BroadcastErrorLog(fileLine)
	}

	ssestream.BroadcastBackendLog(fileLine)

	if !l.shouldLog(entry.severity) {
		return
	}

	// File output if enabled
	if config.GetCreateSSUILogFile() {
		l.writeToFile(fileLine, l.prefix)
	}

	fmt.Print(consoleLine)
}

// shouldLog checks severity and subsystem filters
func (l *Logger) shouldLog(severity int) bool {
	// Subsystem filtering first
	if len(config.GetSubsystemFilters()) > 0 {
		allowed := false
		for _, sub := range config.GetSubsystemFilters() {
			if sub == l.prefix {
				allowed = true
				break
			}
		}
		if !allowed {
			return false // Subsystem not in filter, skip it
		}
	}

	effectiveLevel := config.GetLogLevel()
	return severity >= effectiveLevel
}
