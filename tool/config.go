package tool

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	ErrFileNotExist = errors.New("config file not exist")
	ErrFileParse    = errors.New("config file parse error")
	ErrFileRead     = errors.New("config file read error")
)

func ReadConfig(confPath string) (Conf, error) {
	var err error
	var Config Conf
	var content []byte

	// check file exist
	if _, err = os.Stat(confPath); err != nil {
		return Conf{}, ErrFileNotExist
	}

	if content, err = os.ReadFile(confPath); err != nil {
		return Conf{}, ErrFileRead
	}

	if err = yaml.Unmarshal(content, &Config); err != nil {
		return Conf{}, ErrFileParse
	}

	return Config, nil
}
