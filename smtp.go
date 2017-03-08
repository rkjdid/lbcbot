package main

import (
	"fmt"
	"net/smtp"
	"strings"
)

type SMTPConfig struct {
	Host string
	Port string
	User string
	Pass string
	From string
	To   []string
}

func (cfg *SMTPConfig) SendMail(subject, body string) error {
	b := buildMailBytes(cfg.buildHeader(subject), body)
	return cfg.sendBytes(b)
}

func (cfg *SMTPConfig) sendBytes(msg []byte) error {
	auth := smtp.PlainAuth("", cfg.User, cfg.Pass, cfg.Host)
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	return smtp.SendMail(addr, auth, cfg.From, cfg.To, msg)
}

func (cfg *SMTPConfig) SendTestMessage() error {
	return cfg.SendMail("hallo", "some body")
}

func (cfg *SMTPConfig) buildHeader(subject string) map[string]string {
	return map[string]string{
		"From":         cfg.From,
		"To":           strings.Join(cfg.To, ";"),
		"Subject":      subject,
		"MIME-version": "1.0",
		"Content-Type": "text/html; charset=\"UTF-8\"",
	}
}

func buildMailBytes(header map[string]string, body string) []byte {
	var s string
	for k, v := range header {
		s += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	s += body
	return []byte(s)
}
