package main

import (
	"fmt"

	"github.com/triole/logseal"
)

var (
	lg logseal.Logseal
)

func main() {
	parseArgs()
	lg = logseal.Init(CLI.LogLevel, CLI.LogFile, CLI.LogNoColors, CLI.LogJSON)
	conf := initConf(CLI.Command, CLI.ConfigFile, CLI.DryRun)

	// execute(conf)
	fmt.Printf("%+v\n", conf)
}
