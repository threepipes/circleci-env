package main

import "github.com/kelseyhightower/envconfig"

type Config struct {
	CircleciApiToken string `split_words:"true"`
	ProjectSlug      string `split_words:"true"`
}

func SetConfigFromEnv() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, err
	}
	return &c, nil
}
