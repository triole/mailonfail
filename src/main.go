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
	lg.Info("run " + appName)
	conf := initConf(CLI.Command, CLI.ConfigFile, CLI.DryRun)

	commandReturn := conf.runCmd()
	if commandReturn.Error != nil || commandReturn.Exitcode != 0 || conf.DryRun {
		conf.sendMail(commandReturn)
	}
}
