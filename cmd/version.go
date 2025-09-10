package cmd

import (
	"fan2go-tui/cmd/global"
	"fan2go-tui/internal/logging"

	"github.com/spf13/cobra"
)

var long bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of fan2go-tui",
	Long:  `All software has versions. This is fan2go-tui's`,
	Run: func(cmd *cobra.Command, args []string) {
		if global.Verbose {
			logging.Printfln("%s-%s-%s", global.Version, global.Commit, global.Date)
		} else if long {
			logging.Printfln("%s-%s", global.Version, global.Commit)
		} else {
			logging.Printfln("%s", global.Version)
		}
	},
}

func init() {
	versionCmd.Flags().BoolVarP(&long, "long", "l", false, "Show the long version")

	rootCmd.AddCommand(versionCmd)
}
