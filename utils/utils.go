package utils

import (
	"os"
	"strings"
)

func GetEnv(env string, defaultVal string) (key string) {
	if key = os.Getenv(env); key == "" {
		key = defaultVal
	}
	return
}

func GetArgs(env string) []string {
	return strings.Split(env, ";")
}
