package main

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

var currentConfig *config

type watchConfig struct {
	Units              []string
	FilterActiveStates []string
}

type notifyTeamsConfig struct {
	Enabled    bool
	WebhookURL string
}

type notifyConfig struct {
	Teams notifyTeamsConfig
}

type config struct {
	Watch  watchConfig
	Notify notifyConfig
}

var defaultConfig = `# systemd-notify configuration
[watch]
Units = [ "important-thing.service" ]
FilterActiveStates = []

[notify.teams]
Enabled = false
WebhookURL = ""
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
