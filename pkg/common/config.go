package common

import "os"

func GetEnv(name, defaultValue string) string {
	result := os.Getenv(name)
	if result == "" {
		return defaultValue
	}
	return result
}
