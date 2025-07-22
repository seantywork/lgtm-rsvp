package glob

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig() error {

	conf := Config{}

	file_b, err := os.ReadFile("./config/config.yaml")

	if err != nil {

		return err
	}

	err = yaml.Unmarshal(file_b, &conf)

	if err != nil {

		return err
	}

	G_CONF = &conf

	return nil

}
