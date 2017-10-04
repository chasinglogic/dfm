// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the GPLv3 license that can be found in
// the LICENSE file.

package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/blang/semver"
	"github.com/chasinglogic/dfm/commands"
	"github.com/spf13/cobra"
)

var (
	version = "0.0.0"
	commit  string
	date    string
)

func init() {
	commands.Root.AddCommand(versionCmd)
	commands.Root.AddCommand(update)
}

type githubResponse struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func findDownloadURL(r githubResponse) string {
	suf := fmt.Sprintf("%s_%s.tar.gz", runtime.GOOS, runtime.GOARCH)
	for _, a := range r.Assets {
		if strings.HasSuffix(a.Name, suf) {
			return a.DownloadURL
		}
	}

	return ""
}

// update updates dfm
var update = &cobra.Command{
	Use:   "update",
	Short: "downlaod and install dfm updates",
	Long:  "uses the github API to determine the latest version and install it",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := http.Get("https://api.github.com/repos/chasinglogic/dfm/releases/latest")
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}

		var ghRes githubResponse

		dec := json.NewDecoder(res.Body)
		err = dec.Decode(&ghRes)
		if err != nil {
			fmt.Println("ERROR Reading Github API:", err.Error())
			os.Exit(1)
		}

		localVersion, err := semver.Make(version)
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}

		remoteVersion, err := semver.Make(ghRes.TagName[1:])
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(1)
		}

		if localVersion.GTE(remoteVersion) {
			fmt.Println("Up to date!")
			os.Exit(0)
		}

		url := findDownloadURL(ghRes)
		res, err = http.Get(url)
		if err != nil {
			fmt.Println("ERROR Downloading Tarball:", err.Error())
			os.Exit(1)
		}

		tmpFile, err := ioutil.TempFile("", "")
		if err != nil {
			fmt.Println("ERROR Creating Temp File:", err.Error())
			os.Exit(1)
		}

		_, err = io.Copy(tmpFile, res.Body)
		if err != nil {
			fmt.Println("ERROR Creating Temp File:", err.Error())
			os.Exit(1)
		}

		tmpFile, err = os.Open(tmpFile.Name())
		if err != nil {
			fmt.Println("ERROR Creating Temp File:", err.Error())
			os.Exit(1)
		}

		gz, err := gzip.NewReader(tmpFile)
		if err != nil {
			fmt.Println("ERROR Reading Gzip:", err.Error())
			os.Exit(1)
		}

		arc := tar.NewReader(gz)
		for {
			h, err := arc.Next()
			if err == io.EOF {
				break
			}

			if err != nil {
				fmt.Println("ERROR:", err.Error())
				os.Exit(1)
			}

			if h == nil || h.Name != "dfm" {
				continue
			}

			path, err := exec.LookPath("dfm")
			if err != nil {
				fmt.Println("Looks like DFM isn't installed...", err.Error())
				break
			}

			dfmFile, err := os.Create(path)
			if err != nil {
				fmt.Println("Error Opening DFM:", err.Error())
				break
			}

			_, err = io.Copy(dfmFile, arc)
			if err != nil {
				fmt.Println("Error Extracting DFM:", err.Error())
				break
			}
		}

		fmt.Printf("Updated from %s to %s!\n", localVersion, remoteVersion)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print version information for dfm",
	Run: func(cmd *cobra.Command, args []string) {
		if strings.HasPrefix(version, "SNAPSHOT") {
			version = "DEV-" + commit
		}

		fmt.Printf("DFM Dotfile Manager %s %s/%s BuildDate: %s\n",
			version, runtime.GOOS, runtime.GOARCH, date)
	},
}

func main() {
	if err := commands.Execute(); err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}
}
