package cmd

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		replacer := strings.NewReplacer(".", "_")
		viper.SetEnvKeyReplacer(replacer)

		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			logrus.Error(err)
		}

		viper.AutomaticEnv()

		config := apiserver.NewConfig()
		err := viper.Unmarshal(&config)
		if err != nil {
			logrus.Error(err)
		}

		if err := apiserver.Start(config); err != nil {
			logrus.Fatal(err)
		}
	},
}
