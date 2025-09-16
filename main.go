/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/chasinglogic/dfm/cmd"
)

var Version string

func main() {
	if info, ok := debug.ReadBuildInfo(); ok {
		if Version == "" {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					Version = setting.Value[0:10]
					break
				}
			}
		}

		cmd.RootCmd.Version = fmt.Sprintf(
			"%s built with %s",
			Version,
			info.GoVersion,
		)
	}

	err := cmd.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
