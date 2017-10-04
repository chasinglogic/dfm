package hooks

import (
	"fmt"
	"os/exec"
)

// TODO: write a runCommand for windows

// runCommand is used to run the command in a platform specific way
func runCommand(cmd string) {
	c := exec.Command("bash", "-c", cmd)
	out, err := c.CombinedOutput()
	if err != nil {
		fmt.Println("ERROR Running Command:", cmd, err.Error())
		return
	}
	fmt.Print(string(out))
}
