package main

import (
	"flag"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/apiserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/devopsschoolbot.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	if _, err := toml.DecodeFile(configPath, config); err != nil {
		logrus.Fatal(err)
	}

	if level, err := logrus.ParseLevel(config.LogLevel); err != nil {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(level)
	}

	if err := apiserver.Start(config); err != nil {
		logrus.Fatal(err)
	}

}
