package utils

import "log"

// LogError appends the given string with the message: "Snitch:"
// and print on prompt. Example: "Snitch: Here is the given string"
func LogError(s string) {
	log.Println("Snitch:", s)
}
