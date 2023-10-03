package cmd

import (
	"fan2go-tui/cmd/global"
	"fan2go-tui/internal"
	"fan2go-tui/internal/configuration"
	"fan2go-tui/internal/logging"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fan2go-tui",
	Short: "Easily recover snapshotted versions of files on your ZFS dataset.",
	Long:  ``,
	Args:  cobra.MaximumNArgs(1),
	// this is the default command to run when no subcommand is specified
	Run: func(cmd *cobra.Command, args []string) {
		configPath := configuration.DetectAndReadConfigFile()
		logging.Info("Using configuration file at: %s", configPath)
		configuration.LoadConfig()
		err := configuration.Validate(configPath)
		if err != nil {
			logging.Error("Config Validation Error: %v", err.Error())
			return
		}

		var path string
		if len(args) > 0 {
			path = args[0]
		} else {
			currentWorkingDirectory, err := os.Getwd()
			if err != nil {
				logging.Fatal("Couldn't find current working dir: %v", err)
			}
			path = currentWorkingDirectory
		}

		internal.RunApplication(path)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&global.CfgFile, "config", "c", "", "config file (default is $HOME/.fan2go-tui.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&global.NoColor, "no-color", "", false, "Disable all terminal output coloration")
	rootCmd.PersistentFlags().BoolVarP(&global.NoStyle, "no-style", "", false, "Disable all terminal output styling")
	rootCmd.PersistentFlags().BoolVarP(&global.Verbose, "verbose", "v", false, "More verbose output")
}

func setupUi() {
	logging.SetDebugEnabled(global.Verbose)

	if global.NoColor {
		pterm.DisableColor()
	}
	if global.NoStyle {
		pterm.DisableStyling()
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.OnInitialize(func() {
		configuration.InitConfig(global.CfgFile)
		setupUi()
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
