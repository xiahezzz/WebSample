package tools

import (
	"bytes"
	"os/exec"
)

func RunCommand(cmdStr string) string {
	cmd := exec.Command("bash", "-c", cmdStr)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stderr.String()
	} else {
		return out.String()
	}
}
