// detector.go
package detectionmgr

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

/*
Core Event Detection Engine
- Analyzes the log stream using regex and keyword matching
- Maintains states
- Triggers events with contextual payloads
- Supports extensible pattern matching with custom user defined rules
- Implements multi-stage processing pipeline:
  1. Basic keyword checks
  2. Complex regex pattern matching
  3. Custom rule evaluation
- Handles event distribution to registered handlers
*/

// NewDetector creates a new instance of Detector
func NewDetector() *Detector {
	return &Detector{
		handlers:         make(map[EventType][]Handler),
		connectedPlayers: make(map[string]string),
	}
}

// RegisterHandler registers a handler for a specific event type
func (d *Detector) RegisterHandler(eventType EventType, handler Handler) {
	if _, ok := d.handlers[eventType]; !ok {
		d.handlers[eventType] = []Handler{}
	}
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

// ProcessLogMessage analyzes a log message and triggers appropriate handlers
func (d *Detector) ProcessLogMessage(logMessage string) {
	// Check for simple keyword patterns
	keywordPatterns := map[string]EventType{
		"Ready":                               EventServerReady,
		"Unloading 1 Unused Serialized files": EventServerStarting,
		"EXCEPTION":                           EventServerError,
		"Initialize engine version":           EventServerRunning,
	}

	for keyword, eventType := range keywordPatterns {
		if strings.Contains(logMessage, keyword) {
			d.triggerEvent(Event{
				Type:      eventType,
				Message:   "Server event detected: " + string(eventType),
				RawLog:    logMessage,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		}
	}
	// Process regex patterns for more complex detections
	d.processRegexPatterns(logMessage)

	// Process CUSTOM PATTERNS (both regex and keywords)
	for _, cp := range d.customPatterns {
		if cp.IsRegex {
			matches := cp.Pattern.FindStringSubmatch(logMessage)
			if matches != nil {
				// Format message with {0}, {1} placeholders
				message := formatMessage(cp.MessageTmpl, matches)
				d.triggerEvent(Event{
					Type:      cp.EventType,
					Message:   message,
					RawLog:    logMessage,
					Timestamp: time.Now().Format(time.RFC3339),
				})
			}
		} else {
			// Keyword matching
			if strings.Contains(logMessage, cp.Keyword) {
				d.triggerEvent(Event{
					Type:      cp.EventType,
					Message:   cp.MessageTmpl,
					RawLog:    logMessage,
					Timestamp: time.Now().Format(time.RFC3339),
				})
			}
		}
	}
}

func formatMessage(template string, matches []string) string {
	for i, match := range matches {
		placeholder := fmt.Sprintf("{%d}", i)
		template = strings.ReplaceAll(template, placeholder, match)
	}
	return template
}

// processRegexPatterns handles more complex pattern matching with regex
func (d *Detector) processRegexPatterns(logMessage string) {
	patterns := []struct {
		pattern *regexp.Regexp
		handler func(matches []string, logMessage string)
	}{
		{
			// Player ready pattern
			pattern: regexp.MustCompile(`Client\s+(.+)\s+\((\d+)\)\s+is\s+ready!`),
			handler: func(matches []string, logMessage string) {
				username := matches[1]
				steamID := matches[2]

				d.triggerEvent(Event{
					Type:      EventPlayerReady,
					Message:   "Player is ready",
					RawLog:    logMessage,
					Timestamp: time.Now().Format(time.RFC3339),
					PlayerInfo: &PlayerInfo{
						Username: username,
						SteamID:  steamID,
					},
				})
			},
		},
		{
			// Player connecting pattern
			pattern: regexp.MustCompile(`Client:?\s+(.+?)\s+\((\d+)\)\.\s+Receiving`),
			handler: func(matches []string, logMessage string) {
				username := matches[1]
				steamID := matches[2]

				d.triggerEvent(Event{
					Type:      EventPlayerConnecting,
					Message:   "Player is connecting",
					RawLog:    logMessage,
					Timestamp: time.Now().Format(time.RFC3339),
					PlayerInfo: &PlayerInfo{
						Username: username,
						SteamID:  steamID,
					},
				})
			},
		},
		{
			// Player disconnect pattern
			pattern: regexp.MustCompile(`Client\s+disconnected:\s+\d+\s+\|\s+(.+)\s+connectTime:\s+\d+,\d+s,\s+ClientId:\s+(\d+)`),
			handler: func(matches []string, logMessage string) {
				username := matches[1]
				steamID := matches[2]

				d.triggerEvent(Event{
					Type:      EventPlayerDisconnect,
					Message:   "Player disconnected",
					RawLog:    logMessage,
					Timestamp: time.Now().Format(time.RFC3339),
					PlayerInfo: &PlayerInfo{
						Username: username,
						SteamID:  steamID,
					},
				})
			},
		},
		{
			// World saved pattern
			pattern: regexp.MustCompile(`World Saved:\s.*,\sBackupIndex:\s(\d+)`),
			handler: func(matches []string, logMessage string) {
				backupIndex := matches[1]

				d.triggerEvent(Event{
					Type:      EventWorldSaved,
					Message:   "World saved",
					RawLog:    logMessage,
					Timestamp: time.Now().Format(time.RFC3339),
					BackupInfo: &BackupInfo{
						BackupIndex: backupIndex,
					},
				})
			},
		},
		{
			// Exception pattern
			pattern: regexp.MustCompile(`(?m)^\s*>\s*\d{2}:\d{2}:\d{2}:.*Exception.*|>\s+\d{2}:\d{2}:\d{2}:.*StackTrace`),
			handler: func(matches []string, logMessage string) {
				d.triggerEvent(Event{
					Type:      EventException,
					Message:   "Exception detected",
					RawLog:    logMessage,
					Timestamp: time.Now().Format(time.RFC3339),
					ExceptionInfo: &ExceptionInfo{
						StackTrace: logMessage, // Using the full log message as the stack trace
					},
				})
			},
		},
		{
			pattern: regexp.MustCompile(`\d{2}:\d{2}:\d{2}: Changed setting '(.+?)' from '(.+?)' to '(.+?)'`),
			handler: func(matches []string, logMessage string) {
				settingName := matches[1]
				oldValue := matches[2]
				newValue := matches[3]
				d.triggerEvent(Event{
					Type:      EventSettingsChanged,
					Message:   fmt.Sprintf("Setting %s changed from %s to %s", settingName, oldValue, newValue),
					RawLog:    logMessage,
					Timestamp: time.Now().Format(time.RFC3339),
				})
			},
		},
		{
			pattern: regexp.MustCompile(`RocketNet Succesfully hosted with Address: (.+?) Port: (\d+)`),
			handler: func(matches []string, logMessage string) {
				address := matches[1]
				port := matches[2]
				d.triggerEvent(Event{
					Type:      EventServerHosted,
					Message:   fmt.Sprintf("RocketNet Server hosted at %s:%s", address, port),
					RawLog:    logMessage,
					Timestamp: time.Now().Format(time.RFC3339),
				})
			},
		},
		{
			pattern: regexp.MustCompile(`Started new game in world (.+)`),
			handler: func(matches []string, logMessage string) {
				worldName := matches[1]
				d.triggerEvent(Event{
					Type:      EventNewGameStarted,
					Message:   fmt.Sprintf("New game started in world %s", worldName),
					RawLog:    logMessage,
					Timestamp: time.Now().Format(time.RFC3339),
				})
			},
		},
	}

	for _, p := range patterns {
		if matches := p.pattern.FindStringSubmatch(logMessage); matches != nil {
			p.handler(matches, logMessage)
		}
	}
}

// triggerEvent calls all registered handlers for an event type
func (d *Detector) triggerEvent(event Event) {
	if handlers, ok := d.handlers[event.Type]; ok {
		for _, handler := range handlers {
			handler(event)
		}
	}
}

// GetConnectedPlayers returns a copy of the connected players map
func (d *Detector) GetConnectedPlayers() map[string]string {
	players := make(map[string]string)
	for k, v := range d.connectedPlayers {
		players[k] = v
	}
	return players
}

func (d *Detector) SetCustomPatterns(patterns []CustomPattern) {
	d.customPatterns = patterns
}
