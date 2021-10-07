package configuration

import (
	"fmt"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

// ConfigureLogger ...
func ConfigureLogger() *logrus.Logger {
	logger := logrus.New()

	logger.SetReportCaller(true)
	logger.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}

	logger.SetLevel(logrus.InfoLevel)
	return logger
}
