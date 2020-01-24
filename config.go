package main

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

var currentConfig *config

type config struct {
	Watch watchConfig
}

type watchConfig struct {
	Units              []string
	FilterActiveStates []string
}

var defaultConfig = `# systemd-notify configuration
[watch]
Units = [ "important-thing.service" ]
FilterActiveStates = []
`

func createDefault() error {
	err := ioutil.WriteFile("config.toml", []byte(defaultConfig), 0644)
	if err != nil {
		return err
	}
	return nil
}

func loadConfig() error {
	if _, err := os.Stat("config.toml"); err != nil {
		createDefault()
	}
	if _, err := toml.DecodeFile("config.toml", &currentConfig); err != nil {
		return err
	}

	return nil
}
