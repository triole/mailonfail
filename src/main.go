package main

import (
	"github.com/triole/logseal"
)

var (
	lg logseal.Logseal
)

func main() {
	parseArgs()
	lg = logseal.Init(CLI.LogLevel, CLI.LogFile, CLI.LogNoColors, CLI.LogJSON)
	lg.Info(
		"run "+appName,
		logseal.F{"log-level": CLI.LogLevel, "log-file": CLI.LogFile},
	)
	conf := initConf(CLI.Command, CLI.ConfigFile, CLI.DryRun)

	commandReturn := conf.runCmd()
	if !commandReturn.Success || conf.MailOnSuccess || conf.DryRun {
		conf.sendMail(commandReturn)
	}
}
