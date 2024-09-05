package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/triole/logseal"
)

type cmdReturn struct {
	RunTime  time.Time
	Out      []byte
	Exitcode int
	Error    error
	Success  bool
}

func (conf tConf) runCmd() (cr cmdReturn) {
	cr.RunTime = time.Now()
	var stdBuffer bytes.Buffer
	if !conf.DryRun {
		cmd := exec.Command(conf.Cmd[0], conf.Cmd[1:]...)
		// mw := io.MultiWriter(&stdBuffer)
		mw := io.MultiWriter(os.Stdout, &stdBuffer)

		cmd.Stdout = mw
		cmd.Stderr = mw
		cr.Exitcode = 127 // default exitcode path not found
		if cr.Error = cmd.Run(); cr.Error != nil {
			if exiterr, ok := cr.Error.(*exec.ExitError); ok {
				// the program has exited with an exit code != 0
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					cr.Exitcode = status.ExitStatus()
				}
			}
		}
		if cr.Error != nil {
			lg.Info(
				"run command failed",
				logseal.F{
					"cmd": conf.Cmd, "error": cr.Error, "exitcode": cr.Exitcode,
				},
			)
		}
	} else {
		lg.Info("would have run", logseal.F{"cmd": conf.Cmd})
	}
	cr.Success = cr.Exitcode == 0
	return
}
