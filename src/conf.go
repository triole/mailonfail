package main

import (
	"errors"
	"os"

	"github.com/triole/logseal"
	yaml "gopkg.in/yaml.v3"
)

type tConf struct {
	Cmd           []string
	DryRun        bool
	ConfFile      string
	SmtpHost      string `yaml:"smtp_host"`
	SmtpPort      int    `yaml:"smtp_port"`
	SmtpUser      string `yaml:"smtp_user"`
	SmtpPass      string `yaml:"smtp_pass"`
	MailFrom      string `yaml:"mail_from"`
	MailTo        string `yaml:"mail_to"`
	MailOnSuccess bool   `yaml:"mail_on_success"`
}

func initConf(cmd []string, configFile string, dryRun bool) (conf tConf) {
	if configFile != "" {
		conf = loadConfFile(configFile)
	} else {
		lg.Debug("skip load config file, not defined")
	}
	conf.Cmd = cmd
	conf.DryRun = dryRun
	return
}

func loadConfFile(configFile string) (conf tConf) {
	var err error
	if configFile != "" {
		by, err := os.ReadFile(configFile)
		lg.IfErrFatal(
			"can not read file", logseal.F{"path": configFile, "error": err},
		)
		err = yaml.Unmarshal(by, &conf)
		lg.IfErrFatal(
			"can not unmarshal config", logseal.F{"path": configFile, "error": err},
		)
	} else if errors.Is(err, os.ErrNotExist) {
		lg.Error("conf file does not exist", logseal.F{"file": conf.ConfFile})
	}
	return
}
