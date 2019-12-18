package utils

// Event is a status message
type Event struct {
	Code    int    // Status code of the event; convention is to use `0` for startup, `1` for update and `2` for shutdown handlers
	Message string // Message of the event
}
