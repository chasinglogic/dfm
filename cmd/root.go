/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/chasinglogic/dfm/internal/logger"
	"github.com/chasinglogic/dfm/internal/state"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var debug bool

var rootCmd = &cobra.Command{
	Use:          "dfm",
	Short:        "A dotfile manager for pair programmers and lazy people",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if debug {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		}

		logger.SetLogger(
			zerolog.New(
				zerolog.ConsoleWriter{
					Out:          os.Stderr,
					PartsExclude: []string{"time", "level"},
				},
			).With().Timestamp().Logger(),
		)

		return state.Load()
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return state.Save()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Turn on debug logging")
}
