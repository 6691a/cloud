package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

var setting *Setting

type Setting struct {
	Server     Server     `yaml:"server"`
	DNS        DNS        `yaml:"dns"`
	Hypervisor Hypervisor `yaml:"hypervisor"`
	Router     Router     `yaml:"router"`
}

func GetSetting() *Setting {
	if setting == nil {
		panic("Setting not initialized")
	}
	return setting
}

func NewSetting(path string) *Setting {
	if setting != nil {
		return GetSetting()
	}

	file, err := os.ReadFile(path)

	if err != nil {
		panic("Failed to read YAML file: " + err.Error())
	}

	err = yaml.Unmarshal(file, &setting)

	if err != nil {
		panic("Failed to unmarshal YAML file: " + err.Error())
	}
	return setting
}
