package main

import (
	"fmt"
	"os"

	"github.com/chasinglogic/dfm/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}
}
