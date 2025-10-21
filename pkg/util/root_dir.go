package util

import (
	"path/filepath"
	"runtime"
)

// RootDir get an absolute root dir of current project
func RootDir() string {
	_, b, _, ok := runtime.Caller(0)
	if ok {
		return filepath.Join(filepath.Dir(b), "..", "..")
	}
	return ""
}

// RootDirWithSlash get an absolute root dir of current project with leading '/'
func RootDirWithSlash() string {
	return RootDir() + "/"
}
