package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/triole/logseal"
)

func (conf tConf) RunCmd() ([]byte, int, error) {
	var stdBuffer bytes.Buffer
	var exitcode int
	var err error
	if !conf.DryRun {
		cmd := exec.Command(conf.Cmd[0], conf.Cmd[1:]...)
		// mw := io.MultiWriter(&stdBuffer)
		mw := io.MultiWriter(os.Stdout, &stdBuffer)

		cmd.Stdout = mw
		cmd.Stderr = mw
		if err = cmd.Run(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				// the program has exited with an exit code != 0
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					exitcode = status.ExitStatus()
				}
			}
		}
		if err != nil {
			lg.IfErrError(
				"exec failed", logseal.F{"cmd": conf.Cmd, "error": err},
			)
		}
	} else {
		lg.Info("would have run", logseal.F{"cmd": conf.Cmd})
	}
	return stdBuffer.Bytes(), exitcode, err
}
