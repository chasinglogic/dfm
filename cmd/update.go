package cmd

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/chasinglogic/dfm/logger"
	"github.com/inconshreveable/go-update"
	"github.com/spf13/cobra"
)

type asset struct {
	URL                string `json:"url"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Name               string `json:"name"`
}

type release struct {
	Assets []asset `json:"assets"`
}

func getDFMBinary(gzipStream io.Reader) (*os.File, error) {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return nil, err
	}

	tarReader := tar.NewReader(uncompressedStream)
	var header *tar.Header

	tmpFile, err := os.CreateTemp("", "dfm")
	if err != nil {
		return tmpFile, err
	}

	for header, err = tarReader.Next(); err == nil; header, err = tarReader.Next() {
		if header.Typeflag == tar.TypeReg && header.Name == "dfm" {
			if _, err := io.Copy(tmpFile, tarReader); err != nil {
				// outFile.Close error omitted as Copy error is more interesting at this point
				tmpFile.Close()
				return tmpFile, fmt.Errorf("extracting dfm binary failed: %w", err)
			}
		}
	}

	return tmpFile, nil
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update dfm to the latest release",
	RunE: func(cmd *cobra.Command, args []string) error {
		platform := strings.Title(runtime.GOOS)
		arch := runtime.GOARCH
		url := "https://api.github.com/repos/chasinglogic/dfm/releases/latest"

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Accepts", "application/json")
		req.Header.Add("Content-Type", "application/json")

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		rel := release{}

		err = json.NewDecoder(resp.Body).Decode(&rel)
		if err != nil {
			return err
		}

		downloadURL := ""
		for _, asset := range rel.Assets {
			isForPlatform := strings.Contains(asset.Name, platform) && strings.Contains(asset.Name, arch)
			logger.Debug.Printf("checking if %s contains %s and %s: %s", asset.Name, platform, arch, isForPlatform)
			if isForPlatform {
				downloadURL = asset.BrowserDownloadURL
				break
			}
		}

		if downloadURL == "" {
			return fmt.Errorf("unable to find a download url for your platform: %s %s", platform, arch)
		}

		downloadResp, err := http.Get(downloadURL)
		if err != nil {
			return err
		}

		defer downloadResp.Body.Close()

		dfmBinary, err := getDFMBinary(downloadResp.Body)
		if err != nil {
			return err
		}

		return update.Apply(dfmBinary, update.Options{})
	},
}
