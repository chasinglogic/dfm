package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

func getInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return text[:len(text)-1]
}

func setupWizard() {
	confirmation := "no"

	for strings.HasPrefix(strings.ToLower(confirmation), "n") {
		global.DefaultBackend = getInput(fmt.Sprintf("Which backend would you like to use? [git, dropbox] (Default: git) "))
		if global.DefaultBackend == "" {
			global.DefaultBackend = "git"
		}

		fmt.Println("\nNOTE: Some backends will change this automatically and override your setting. (i.e. dropbox)")
		global.Dir = getInput(fmt.Sprintf("Where would you like to store profiles? (Default: %s) ", GetDefaultConfigDir()))
		if global.Dir == "" {
			global.Dir = GetDefaultConfigDir()
		}

		yml, _ := yaml.Marshal(global)
		fmt.Println(string(yml))

		confirmation = getInput("Does this look correct? Y/n: ")
	}

	err := SaveConfig()
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
}
