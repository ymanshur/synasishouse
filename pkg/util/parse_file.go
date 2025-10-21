package util

import (
	"os"

	"github.com/goccy/go-yaml"
)

// ReadFromYAML reads the YAML file and pass to the object
// Args:
//
//	path: file path location
//	target: object which will hold the value
//
// Returns:
//
//	error: operation state error
func ReadFromYAML(path string, target any) error {
	f, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(f, target)
}
