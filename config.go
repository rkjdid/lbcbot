package main

import (
	"encoding/json"
	"github.com/soul9/errors"
	"io"
	"os"
)

type Config struct {
	SMTPConfig      SMTPConfig
	WatchList       []Query
	PollIntervalMin int
}

func NewConfig() *Config {
	return &Config{
		SMTPConfig:      SMTPConfig{},
		PollIntervalMin: 10,
		WatchList:       make([]Query, 1),
	}
}

func LoadConfigFile(cfile string) (cfg *Config, err error) {
	var fd *os.File
	fd, err = os.Open(cfile)
	if err != nil {
		return nil, errors.New(err)
	}
	defer fd.Close()
	cfg, err = LoadConfig(fd)
	if err != nil {
		return cfg, errors.New(err)
	}
	return cfg, nil
}

func LoadConfig(rd io.Reader) (cfg *Config, err error) {
	dec := json.NewDecoder(rd)
	cfg = &Config{}
	err = dec.Decode(cfg)
	if (err != nil) && (err != io.EOF) {
		return cfg, errors.NewError(err)
	}
	return cfg, nil
}

func (cfg *Config) WriteFile(path string) error {
	var (
		fd  *os.File
		err error
	)
	fd, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.New(err)
	}
	defer fd.Close()
	if err = cfg.Write(fd); err != nil {
		return errors.NewError(err)
	}
	return nil
}

func (cfg *Config) Write(w io.Writer) error {
	prettyJson, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return errors.New(err)
	}
	w.Write(prettyJson)
	w.Write([]byte{'\n'})
	return nil
}
