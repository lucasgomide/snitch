package utils

import (
	"os"
	"strings"
)

// GetEnv interprets and split with "$" the given string to find
// the ENV value
// It returns ENV value if split of given string lenght 2
// and the ENV value is not empty
// It returns the given name if split of given string not is lenght 2
// or if given string length is 2 but the ENV value is empty
func GetEnv(s string) string {
	sliceWords := strings.Split(s, "$")
	if len(sliceWords) == 2 {
		env := os.Getenv(sliceWords[1])
		if env != "" {
			return env
		}
	}
	return s
}
