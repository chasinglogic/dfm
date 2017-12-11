// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package hooks

import (
	"fmt"
	"os/exec"
)

// TODO: write a runCommand for windows

// RunCommands will run the given slice of strings each as their own command
func RunCommands(commands []string) {
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
