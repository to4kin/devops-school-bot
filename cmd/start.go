package cmd

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/configuration"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/server"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/server/apiserver"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/server/awslambda"
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

		config := configuration.NewConfig()
		err := viper.Unmarshal(&config)
		if err != nil {
			logrus.Error(err)
		}

		config.BuildDate = buildDate
		config.Version = version

		var srv server.Server
		if config.AWSLambda.Enabled {
			srv = awslambda.New(config)
		} else {
			srv = apiserver.New(config)
		}

		if err := srv.Start(); err != nil {
			logrus.Fatal(err)
		}
	},
}
