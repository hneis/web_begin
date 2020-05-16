package config

import (
	"log/syslog"
	"os"

	"github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

var LOG_TAG = "serv"

type LoggerConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
	Syslog bool   `yaml:"syslog"`
}

func NewLogger(conf *LoggerConfig) (*logrus.Logger, error) {
	lg := logrus.New()
	level, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		return nil, err
	}
	lg.SetLevel(level)

	if conf.Syslog {
		hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")
		if err != nil {
			return nil, err
		}
		lg.Hooks.Add(hook)
	} else {
		if conf.Output != "" {
			f, err := os.Create(conf.Output)
			if err != nil {
				return nil, err
			}
			lg.SetOutput(f)
		}
	}

	return lg, nil
}
