/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"time"

	"github.com/chasinglogic/dfm/internal/logger"
	"github.com/chasinglogic/dfm/internal/state"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var debugMode bool

var RootCmd = &cobra.Command{
	Use:          "dfm",
	Short:        "A dotfile manager for pair programmers and lazy people",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if debugMode {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		}

		logger.SetLogger(
			zerolog.New(
				zerolog.ConsoleWriter{
					Out:        os.Stderr,
					TimeFormat: time.RFC3339,
				},
			).With().Timestamp().Logger(),
		)

		return state.Load()
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return state.Save()
	},
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(
		&debugMode,
		"debug",
		"d",
		os.Getenv("DFM_DEBUG") != "",
		"Turn on debug logging",
	)
}
