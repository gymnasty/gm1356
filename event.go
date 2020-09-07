package gm1356

import (
	"fmt"
	"time"
)

// EventType is type of event from GM1356 device
type EventType int

const (
	// EventTypeConfigured is configured event type
	EventTypeConfigured EventType = iota
	// EventTypeMeasured is measured event type
	EventTypeMeasured
)

// Event from GM1356 device
type Event interface {
	Type() EventType
	String() string
}

// ConfiguredEvent from GM1356
type ConfiguredEvent struct {
}

// Type returns event type
func (e ConfiguredEvent) Type() EventType {
	return EventTypeConfigured
}

// String returns string describes event
func (e ConfiguredEvent) String() string {
	return "ConfiguredEvent{}"
}

// MeasuredEvent from GM1356
type MeasuredEvent struct {
	Time       time.Time
	SoundLevel float32
	Config     Config
}

// Type returns event type
func (e MeasuredEvent) Type() EventType {
	return EventTypeMeasured
}

// String returns string describes event
func (e MeasuredEvent) String() string {
	return fmt.Sprintf("MeasuredEvent{Time: %s, SoundLevel: %.1fdB, Config: %s}", e.Time.Format("2006/01/02 15:04:05"), e.SoundLevel, e.Config.String())
}
