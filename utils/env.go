package utils

import "os"

func GetEnvWithDefault(key string, defaultval string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	} else {
		return defaultval
	}
}
