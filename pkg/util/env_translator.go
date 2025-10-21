package util

import (
	"strings"
)

var envs = map[string]string{
	"production":  "production",
	"staging":     "staging",
	"development": "development",
	"prod":        "production",
	"stg":         "staging",
	"dev":         "development",
	"prd":         "production",
	"green":       "green",
	"blue":        "blue",
}

// TranslateEnv transfors environment
func TranslateEnv(s string) string {
	cleaned := strings.Trim(s, " ")
	for env, v := range envs {
		if strings.Contains(cleaned, env) {
			return v
		}
	}

	return s
}
