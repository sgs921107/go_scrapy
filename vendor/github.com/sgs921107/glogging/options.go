package glogging

import (
	"time"
)


// Options 配置选项
type Options struct {
	Level			string
	FilePath		string
	Formater		string
	RotationMaxAge	time.Duration
	RotationTime	time.Duration
	NoLock			bool
	// ForceNewFile	bool
}
