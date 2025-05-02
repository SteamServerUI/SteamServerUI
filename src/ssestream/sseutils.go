// sseutils.go
package ssestream

import (
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
)

// Global managers for SSE streams
var (
	ConsoleStreamManager = NewSSEManager(config.GetMaxSSEConnections(), config.GetSSEMessageBufferSize())
	EventStreamManager   = NewSSEManager(config.GetMaxSSEConnections(), config.GetSSEMessageBufferSize())
	InternalLogStream    = NewSSEManager(config.GetMaxSSEConnections(), config.GetSSEMessageBufferSize())
)

// BroadcastConsoleOutput sends log to all connected console log clients
func BroadcastConsoleOutput(message string) {
	ConsoleStreamManager.Broadcast(message)
}

// BroadcastDetectionEvent sends an event to all connected clients
func BroadcastDetectionEvent(message string) {
	EventStreamManager.Broadcast(message)
}

// BroadcastInternalLog sends an event to all connected clients
func BroadcastInternalLog(message string) {
	InternalLogStream.Broadcast(message)
}
