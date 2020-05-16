package config

import (
	"io/ioutil"
	"os"

	"github.com/hneis/web_begin/lesson8/homework/server"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Logger LoggerConfig        `yaml:"logger"`
	Server server.ServerConfig `yaml:"server"`
}

func ReadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	conf := Config{}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, err
}
