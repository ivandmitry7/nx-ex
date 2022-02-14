package task

import (
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type ReplacerCfg struct {
	Search  []string
	Replace string
}

type ParserCfg struct {
	Name   string
	Source string
	Reason string
	Times  []struct {
		Search  string
		Replace string
	}
	Coords []struct {
		Type    string
		Search  string
		Replace string
	}
}

type Config struct {
	Replacers []ReplacerCfg
	Parsers   []ParserCfg
}

func (c *Config) applyReplacers() {
	var list []string
	for _, rp := range c.Replacers {
		for _, s := range rp.Search {
			list = append(list, s)
			list = append(list, rp.Replace)
		}
	}

	replacer := strings.NewReplacer(list...)
	for _, p := range c.Parsers {
		p.Source = replacer.Replace(p.Source)
		p.Reason = replacer.Replace(p.Reason)
		for _, t := range p.Times {
			t.Search = replacer.Replace(t.Search)
			t.Replace = replacer.Replace(t.Replace)
		}
		for _, c := range p.Coords {
			c.Search = replacer.Replace(c.Search)
			c.Replace = replacer.Replace(c.Replace)
		}
	}
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
	if err := d.Decode(config); err != nil {
		return nil, err
	}
	config.applyReplacers()
	return config, nil
}
