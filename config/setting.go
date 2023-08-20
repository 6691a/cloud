package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Setting struct {
	Server     Server     `yaml:"server"`
	DNS        DNS        `yaml:"dns"`
	Hypervisor Hypervisor `yaml:"hypervisor"`
}

func NewSetting(path string) *Setting {
	file, err := os.ReadFile(path)

	if err != nil {
		panic("Failed to read YAML file: " + err.Error())
	}

	if err != nil {
		panic("Failed to unmarshal YAML file: " + err.Error())
	}

	var Setting Setting
	err = yaml.Unmarshal(file, &Setting)

	if err != nil {
		panic("Failed to unmarshal YAML file: " + err.Error())
	}

	return &Setting
}
