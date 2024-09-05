package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/alecthomas/kong"
)

var (
	BUILDTAGS      string
	appName        = "mailonfail"
	appDescription = "If a command fails, send a mail...\n\nUsage: mailonfail [flags] [<command> ...]"
	appMainversion = "0.1"
)

var CLI struct {
	Command      []string `help:"command to run, flags always have to be in front" arg:"" optional:"" passthrough:""`
	ConfigFile   string   `help:"config file to load, values can be overwritten by env vars" short:"c"`
	LogFile      string   `help:"log file" default:"/dev/stdout"`
	LogLevel     string   `help:"log level" default:"info" enum:"trace,debug,info,error"`
	LogNoColors  bool     `help:"disable output colours, print plain text"`
	LogJSON      bool     `help:"enable json log, instead of text one"`
	ValidateConf bool     `help:"validate configuration and pretty print it"`
	DryRun       bool     `help:"dry run, just print operations that would run" short:"n"`
	VersionFlag  bool     `help:"display version" short:"V"`
}

func parseArgs() {
	ctx := kong.Parse(&CLI,
		kong.Name(appName),
		kong.Description(appDescription),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			NoAppSummary: true,
			Compact:      true,
			Summary:      true,
			FlagsLast:    false,
		}),
	)
	_ = ctx.Run()
	if CLI.VersionFlag {
		printBuildTags(BUILDTAGS)
		os.Exit(0)
	}
	if len(CLI.Command) < 1 {
		ctx.FatalIfErrorf(errors.New("command required"))
	}
}

func printBuildTags(buildtags string) {
	regexp, _ := regexp.Compile(`({|}|,)`)
	s := regexp.ReplaceAllString(buildtags, "\n")
	s = strings.Replace(s, "_subversion: ", "version: "+appMainversion+".", -1)
	fmt.Printf("%s\n", s)
}
