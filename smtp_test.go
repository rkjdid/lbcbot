package main

import (
	"testing"
)

func TestSendMessage(t *testing.T) {
	cfg, err := LoadConfigFile("smtp_test.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := cfg.SMTPConfig.SendTestMessage(); err != nil {
		t.Error(err)
	}
}
