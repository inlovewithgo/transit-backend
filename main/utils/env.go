package utils

import (
	"os"
)

func GetENV(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}