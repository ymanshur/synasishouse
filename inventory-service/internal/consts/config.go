package consts

import "time"

const (
	DefaultServerTimeout  = time.Duration(30 * time.Second)
	DefaultConfigFilename = "app"
	DefaultConfigExt      = "env"
)
