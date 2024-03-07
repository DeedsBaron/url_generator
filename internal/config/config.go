package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	GrpcPort int `yaml:"grpcPort"`
	HttpPort int `yaml:"httpPort"`
	UrlLen   int `yaml:"urlLen"`
}

var Data Config

func New() error {
	rawYAML, err := os.ReadFile("config.yaml")
	if err != nil {
		return errors.Wrap(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &Data)
	if err != nil {
		return errors.Wrap(err, "parsing yaml")
	}

	if Data.UrlLen > 100 {
		return errors.New("len of generated url can't be greater than 100")
	}

	return nil
}
