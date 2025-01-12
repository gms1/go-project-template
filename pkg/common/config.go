package common

import "os"

func Getenv(name, defaultValue string) string {
	result := os.Getenv(name)
	if result == "" {
		return defaultValue
	}
	return result
}
