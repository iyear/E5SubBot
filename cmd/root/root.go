package root

import (
	"github.com/fatih/color"
	"github.com/iyear/E5SubBot/cmd/run"
	"github.com/iyear/E5SubBot/global"
	"github.com/spf13/cobra"
)

var version bool

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "E5SubBot",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			color.Blue("%s\n%s", global.Version, global.GetRuntime())
		}
	},
}

func init() {
	Cmd.AddCommand(run.Cmd)
	Cmd.PersistentFlags().BoolVarP(&version, "version", "v", false, "check the version of E5SubBot")
}

func Execute() {
	cobra.CheckErr(Cmd.Execute())
}
