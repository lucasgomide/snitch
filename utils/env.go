package utils

import (
	"os"
	"strings"
)

func GetEnv(s string) string {
	sliceWords := strings.Split(s, "$")
	if len(sliceWords) == 1 || len(sliceWords) > 2 {
		return s
	} else {
		env := os.Getenv(sliceWords[1])
		if env == "" {
			return s
		} else {
			return env
		}
	}
}
