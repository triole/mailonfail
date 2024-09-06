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

type tCmdReturn struct {
	RunStart    time.Time
	RunEnd      time.Time
	RunDuration time.Duration
	Output      []byte
	Exitcode    int
	Error       error
	Success     bool
}

func (conf tConf) runCmd() (cr tCmdReturn) {
	cr.RunStart = time.Now()
	var stdBuffer bytes.Buffer
	if !conf.DryRun {
		cmd := exec.Command(conf.Cmd[0], conf.Cmd[1:]...)
		// mw := io.MultiWriter(&stdBuffer)
		mw := io.MultiWriter(os.Stdout, &stdBuffer)

		_, cr.Error = exec.LookPath(conf.Cmd[0])
		if cr.Error == nil {
			cmd.Stdout = mw
			cmd.Stderr = mw
			if cr.Error = cmd.Run(); cr.Error != nil {
				if exiterr, ok := cr.Error.(*exec.ExitError); ok {
					// the program has exited with an exit code != 0
					if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
						cr.Exitcode = status.ExitStatus()
					}
				}
			}
		} else {
			cr.Exitcode = 127
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
	cr.RunEnd = time.Now()
	cr.RunDuration = cr.RunEnd.Sub(cr.RunStart)
	cr.Success = (cr.Error == nil) && (cr.Exitcode == 0)
	cr.Output = stdBuffer.Bytes()
	return
}
