package main

import (
	"os"
	"os/exec"
	"strings"
)

func Run(prog string, args ...string) ([]byte, error) {
	cmd := strings.Split(prog, " ")
	if len(cmd) > 0 {
		args = append(cmd[1:], args...)
	}

	execCmd := exec.Command(cmd[0], args...)
	ret, err := execCmd.CombinedOutput()
	if len(ret) > 0 || (err != nil && err.Error() == "signal: aborted") {
		err = nil //Hack to get around programs that exit non-zero or abort, we always want the output
	}
	//log("RUN: %s %v\n%s", cmd[0], cmdArgs, string(ret))
	//log("RUN: %s %v", cmd[0], cmdArgs)
	return ret, err
}

func RunRealtime(prog string, args ...string) error {
	cmd := strings.Split(prog, " ")
	if len(cmd) > 0 {
		args = append(cmd[1:], args...)
	}

	execCmd := exec.Command(cmd[0], args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	err := execCmd.Run()
	if err != nil && err.Error() == "signal: aborted" {
		err = nil
	}
	return err
}
