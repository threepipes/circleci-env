package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ApiToken         string `split_words:"true"`
	OrganizationName string `split_words:"true"`
}

func getConfigPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get config path: %w", err)
	}
	return dir + "/.config/ccienv/config.yml", nil
}

func ReadConfig() (*Config, error) {
	cp, err := getConfigPath()
	cp = filepath.Dir(cp)
	if err != nil {
		return nil, err
	}
	viper.SetConfigType("yaml")
	viper.AddConfigPath(cp)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Errorln("No config is found. Please execute `ccienv config init`")
			return nil, fmt.Errorf("no settings found")
		} else {
			return nil, fmt.Errorf("read config from a config file: %w", err)
		}
	}
	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("read config from a config file: %w", err)
	}
	return &c, nil
}

func WriteConfig(conf *Config) error {
	cp, err := getConfigPath()
	if err != nil {
		return err
	}
	bt, err := json.Marshal(&conf)
	if err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	viper.SetConfigType("json")
	if err := viper.ReadConfig(bytes.NewBuffer(bt)); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	viper.SetConfigType("yaml")
	dirpath := filepath.Dir(cp)
	if err := os.MkdirAll(dirpath, os.ModePerm); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	if err := viper.WriteConfigAs(cp); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}
