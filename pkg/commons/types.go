// Package commons provides data structures used between packages
package commons

import (
	"time"
)

// Options provides overall configuration data
type Options struct {
	ServerPort       string
	DatabaseFileName string
	LogLevel         string
}

// BodyMeasurement provides generic data struct for storing measurement data
type BodyMeasurement struct {
	ID        int    `storm:"id,increment"`
	Type      string `storm:"index"`
	Value     float64
	UserID    string    `storm:"index"`
	Timestamp time.Time `storm:"index"`
}
