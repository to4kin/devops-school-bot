package cmd

import (
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/apiserver"
)

var (
	configPath string
)

func init() {
	startCmd.PersistentFlags().StringVarP(&configPath, "config-path", "c", "configs/devopsschoolbot.toml", "path to config file")
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start DevOps School Bot",
	Long: `Start DevOps School Bot with config file
Simply execute devops-school-bot start -c path/to/config/file.toml
or skip this flag to use default path`,
	Run: func(cmd *cobra.Command, args []string) {
		config := apiserver.NewConfig()
		if _, err := toml.DecodeFile(configPath, config); err != nil {
			logrus.Fatal(err)
		}

		if err := apiserver.Start(config); err != nil {
			logrus.Fatal(err)
		}
	},
}
