package pkg

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Conf struct {
	Consul Consul `yaml:"consul"`
}

type Consul struct {
	Address string `yaml:"address"`
	Scheme  string `yaml:"scheme"`
}

var conf Conf

func InitConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &conf); err != nil {
		return err
	}

	return nil
}

func GetConfig() Conf {
	return conf
}
