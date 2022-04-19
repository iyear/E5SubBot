package run

import (
	"github.com/iyear/E5SubBot/app/bot"
	"github.com/spf13/cobra"
)

var (
	cfg      string
	tmplPath string
	dataPath string
)

var Cmd = &cobra.Command{
	Use:   "run",
	Short: "Start the bot",
	Long:  `Start the bot`,
	Run: func(cmd *cobra.Command, args []string) {
		bot.Run(cfg, tmplPath, dataPath)
	},
}

func init() {
	Cmd.PersistentFlags().StringVarP(&cfg, "config", "c", "config.yaml", "bot config file")
	Cmd.PersistentFlags().StringVarP(&tmplPath, "template", "t", "./template", "template files path")
	Cmd.PersistentFlags().StringVarP(&dataPath, "data", "d", "./data", "data files path")
}
