package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/triole/logseal"
)

func initExpReturn(output, errString string, exitcode int, success bool) (cr tCmdReturn) {
	cr.Output = []byte(output)
	cr.Exitcode = exitcode
	cr.Success = success
	err := errors.New(errString)
	cr.Error = err
	return
}

func TestRunCmd(t *testing.T) {
	validateRunCmd(
		[]string{"echo", "hello world"},
		initExpReturn(
			"hello world\n", "", 0, true,
		), t,
	)
	validateRunCmd(
		[]string{"echooo"},
		initExpReturn(
			"", "executable file not found in $PATH", 127, false,
		),
		t,
	)
	validateRunCmd(
		[]string{"ls", "/this_folder_does_not_exist"},
		initExpReturn(
			"", "exit status 2", 2, false,
		),
		t,
	)
}

func validateRunCmd(cmd []string, exp tCmdReturn, t *testing.T) {
	lg = logseal.Init("trace")
	conf := tConf{
		Cmd:    cmd,
		DryRun: false,
	}
	cr := conf.runCmd()
	if !strings.Contains(string(cr.Output), string(exp.Output)) {
		t.Errorf("fail runCmd %q output %q != %q", cmd, cr.Output, exp.Output)
	}
	if !strings.Contains(
		fmt.Sprintf("%s", cr.Error), fmt.Sprintf("%s", exp.Error),
	) {
		t.Errorf("fail runCmd %q error: '%s' != '%s'", cmd, cr.Error, exp.Error)
	}
	if cr.Exitcode != exp.Exitcode {
		t.Errorf("fail runCmd %q exitcode %d != %d", cmd, cr.Exitcode, exp.Exitcode)
	}
	if cr.Success != exp.Success {
		t.Errorf("fail runCmd %q success %v != %v", cmd, cr.Success, exp.Success)
	}
}
