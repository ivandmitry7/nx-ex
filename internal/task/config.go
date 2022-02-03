package task

import (
	"gopkg.in/yaml.v3"
	"os"
)

type ParserCfg struct {
	Name   string
	Source string
	Reason string
	Times  []struct {
		Search  string
		Replace string
	}
	Coords string
}

type Config struct {
	Parsers []ParserCfg
}

func NewConfig(cfgPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config.Parsers); err != nil {
		return nil, err
	}
	return config, nil
}
