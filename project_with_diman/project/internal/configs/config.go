package configs

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

var (
	ErrData = errors.New("Empty data")
)

type Config struct {
	Port    int `yaml:"port"`
	Project struct {
		Db struct {
			Login    string `yaml:"login"`
			Password string `yaml:"password"`
		} `yaml:"db"`
		Port string `yaml:"port"`
	} `yaml:"project"`
}

func New() (*Config, error) {
	file, err := os.Open("../config.yaml")
	if err != nil {
		return nil, err
	}
	yamlFile, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Read config file error: %s", err)
	}
	var config Config
	fmt.Println(string(yamlFile))
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal config file error: %s", err)
	}
	if config.Project.Db.Login == "" || config.Project.Db.Password == "" || config.Project.Port == "" {
		return nil, fmt.Errorf("Empty data %w", ErrData)
	}
	return &config, nil
}
