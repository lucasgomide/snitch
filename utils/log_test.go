package utils

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestPrintMessage(t *testing.T) {
	var buffer bytes.Buffer
	log.SetOutput(&buffer)

	LogError("Message to print")

	msg := buffer.String()
	expected := "Snitch: Message to print"

	if !strings.Contains(msg, expected) {
		t.Errorf("Exptected to contain message %#v, got %#v", expected, msg)
	}
}
