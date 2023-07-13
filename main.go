package main

import (
	"fmt"

	"github.com/chasinglogic/dfm/cmd"
)

var version string
var commit string
var date string

func main() {
	fullVersion := fmt.Sprintf("%s-%s %s", version, commit, date)
	cmd.Execute(fullVersion)
}
