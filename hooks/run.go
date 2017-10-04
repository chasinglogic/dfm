package hooks

import (
	"fmt"
	"os/exec"
)

// TODO: write a runCommand for windows

func runCommands(commands []string) {
	for _, cmd := range commands {
		c := exec.Command("bash", "-c", cmd)
		out, err := c.CombinedOutput()
		if err != nil {
			fmt.Println("ERROR Running Command:", cmd, err.Error())
			return
		}
		fmt.Print(string(out))
	}
}
