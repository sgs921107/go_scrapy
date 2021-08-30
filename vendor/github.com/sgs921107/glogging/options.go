package glogging

import (
	"time"
)

// Options	logger options
type Options struct {
	Level			string
	FilePath		string
	Formatter		string
	RotationMaxAge	time.Duration
	RotationTime	time.Duration
	// 仅logrus
	NoLock			bool
}
