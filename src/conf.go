package main

import (
	"errors"
	"os"
	"strconv"
	"strings"

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
	getEnvVars(&conf)
	return
}

func loadConfFile(configFile string) (conf tConf) {
	var err error
	if configFile != "" {
		by, err := os.ReadFile(configFile)
		lg.IfErrFatal(
			"can not read file",
			logseal.F{"path": configFile, "error": err},
		)
		err = yaml.Unmarshal(by, &conf)
		lg.IfErrFatal(
			"can not unmarshal config",
			logseal.F{"path": configFile, "error": err},
		)
	} else if errors.Is(err, os.ErrNotExist) {
		lg.Error(
			"conf file does not exist",
			logseal.F{"path": conf.ConfFile, "error": err})
	}
	return
}

func getEnvVars(conf *tConf) *tConf {
	msg := "set conf val from env var"
	for _, envVar := range os.Environ() {
		key, val, err := splitEnvVar(envVar)
		if err == nil {
			switch key {
			case "MOF_SMTP_HOST":
				lg.Debug(msg, logseal.F{"key": "smtp_host", "val": val})
				conf.SmtpHost = val
			case "MOF_SMTP_PORT":
				lg.Debug(msg, logseal.F{"key": "smtp_port", "val": val})
				p, err := strconv.Atoi(val)
				if err == nil {
					conf.SmtpPort = p
				} else {
					lg.Error(
						"port invalid, must be an integer",
						logseal.F{"port": val},
					)
				}
			case "MOF_SMTP_USER":
				lg.Debug(msg, logseal.F{"key": "smtp_user", "val": val})
				conf.SmtpUser = val
			case "MOF_SMTP_PASS":
				lg.Debug(msg, logseal.F{"key": "smtp_pass", "val": val})
				conf.SmtpPass = val
			case "MOF_MAIL_FROM":
				lg.Debug(msg, logseal.F{"key": "mail_from", "val": val})
				conf.MailFrom = val
			case "MOF_MAIL_TO":
				lg.Debug(msg, logseal.F{"key": "mail_to", "val": val})
				conf.MailTo = val
			case "MOF_MAIL_ON_SUCCESS":
				lg.Debug(msg, logseal.F{"key": "mail_on_success", "val": val})
				conf.MailOnSuccess = stringToBool(val)
			}
		}
	}
	return conf
}

func splitEnvVar(envVar string) (key, val string, err error) {
	pair := strings.SplitN(envVar, "=", 2)
	if len(pair) > 1 {
		key = pair[0]
		val = pair[1]
	} else {
		err = errors.New("failed to parse env var string: " + envVar)
	}
	return
}

func stringToBool(s string) (b bool) {
	var err error
	b, err = strconv.ParseBool(s)
	if err != nil {
		lg.Error(
			"can not parse string to bool value",
			logseal.F{"string": s},
		)
	}
	return b
}
