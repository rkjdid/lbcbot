package main

type Config struct {
	SMTPConfig      SMTPConfig
	WatchList       []*Query
	PollIntervalMin int
	HtmlRoot        string
	cache           map[string]bool
}

func NewConfig() *Config {
	return &Config{
		SMTPConfig:      SMTPConfig{},
		PollIntervalMin: 10,
		WatchList:       make([]*Query, 1),
	}
}
