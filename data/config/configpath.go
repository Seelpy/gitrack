package config

import (
	stderrors "errors"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var (
	ErrPathNotExist = stderrors.New("gitrack config path not found")
)

type ConfigPath struct {
	configPath string
}

type ConfigPathYAML struct {
	ConfigFilePath string `yaml:"config_path"`
}

func SetConfigPath(configPath string, input string) error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)

	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(exeDir, configPath)
	}

	if !filepath.IsAbs(input) {
		input, err = filepath.Abs(input)
		if err != nil {
			return err
		}
	}

	config := ConfigPathYAML{
		ConfigFilePath: input,
	}

	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(yamlData)
	if err != nil {
		return err
	}
	return nil
}

func GetConfigPath(configPath string) (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exePath)

	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(exeDir, configPath)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return "", errors.WithStack(ErrPathNotExist)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	var config ConfigPathYAML
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return "", err
	}

	return config.ConfigFilePath, nil
}
